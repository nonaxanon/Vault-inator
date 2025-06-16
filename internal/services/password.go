package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// Password represents a stored password entry
type Password struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL      string `json:"url"`
	Notes    string `json:"notes"`
}

// PasswordService handles password storage and retrieval
type PasswordService struct {
	filePath string
	mu       sync.RWMutex
	block    cipher.Block
}

// NewPasswordService creates a new password service instance
func NewPasswordService() *PasswordService {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}

	filePath := filepath.Join(homeDir, ".vaultinator", "passwords.json")
	return &PasswordService{
		filePath: filePath,
	}
}

// SetEncryptionKey sets the encryption key for the password service
func (s *PasswordService) SetEncryptionKey(key []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("failed to create cipher: %w", err)
	}
	s.block = block
	return nil
}

// loadPasswords reads and decrypts passwords from the storage file
func (s *PasswordService) loadPasswords() ([]Password, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(s.filePath), 0700); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Read file if it exists
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Password{}, nil
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return []Password{}, nil
	}

	// Decrypt data
	decrypted, err := s.decrypt(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	var passwords []Password
	if err := json.Unmarshal(decrypted, &passwords); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return passwords, nil
}

// savePasswords encrypts and writes passwords to the storage file
func (s *PasswordService) savePasswords(passwords []Password) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Marshal passwords to JSON
	data, err := json.Marshal(passwords)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Encrypt data
	encrypted, err := s.encrypt(data)
	if err != nil {
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(s.filePath, encrypted, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// encrypt encrypts data using AES-GCM
func (s *PasswordService) encrypt(data []byte) ([]byte, error) {
	// Create a new GCM
	gcm, err := cipher.NewGCM(s.block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and seal
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// decrypt decrypts data using AES-GCM
func (s *PasswordService) decrypt(data []byte) ([]byte, error) {
	// Create a new GCM
	gcm, err := cipher.NewGCM(s.block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt and open
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// GetAllPasswords returns all stored passwords
func (s *PasswordService) GetAllPasswords() ([]Password, error) {
	return s.loadPasswords()
}

// CreatePassword adds a new password entry
func (s *PasswordService) CreatePassword(password *Password) error {
	passwords, err := s.loadPasswords()
	if err != nil {
		return err
	}

	password.ID = uuid.New().String()
	passwords = append(passwords, *password)

	return s.savePasswords(passwords)
}

// UpdatePassword updates an existing password entry
func (s *PasswordService) UpdatePassword(password *Password) error {
	passwords, err := s.loadPasswords()
	if err != nil {
		return err
	}

	for i, p := range passwords {
		if p.ID == password.ID {
			passwords[i] = *password
			return s.savePasswords(passwords)
		}
	}

	return fmt.Errorf("password not found")
}

// DeletePassword removes a password entry
func (s *PasswordService) DeletePassword(id string) error {
	passwords, err := s.loadPasswords()
	if err != nil {
		return err
	}

	for i, p := range passwords {
		if p.ID == id {
			passwords = append(passwords[:i], passwords[i+1:]...)
			return s.savePasswords(passwords)
		}
	}

	return fmt.Errorf("password not found")
}
