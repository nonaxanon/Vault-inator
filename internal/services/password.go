package services

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/nonaxanon/vault-inator/internal/encryption"
	"github.com/nonaxanon/vault-inator/internal/storage"
)

// Password represents a password entry in the service layer.
type Password struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	URL      string    `json:"url"`
	Notes    string    `json:"notes"`
}

// PasswordService handles password storage and retrieval
type PasswordService struct {
	db        *storage.DB
	mu        sync.RWMutex
	encryptor *encryption.Encryptor
}

// NewPasswordService creates a new password service instance
func NewPasswordService(db *storage.DB) *PasswordService {
	return &PasswordService{
		db: db,
	}
}

// SetEncryptionKey sets the encryption key for the password service
func (s *PasswordService) SetEncryptionKey(key []byte) error {
	encryptor, err := encryption.NewEncryptor(string(key))
	if err != nil {
		return fmt.Errorf("failed to create encryptor: %w", err)
	}
	s.encryptor = encryptor
	return nil
}

// GetAllPasswords returns all stored passwords
func (s *PasswordService) GetAllPasswords() ([]Password, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := s.db.GetAllPasswords()
	if err != nil {
		return nil, fmt.Errorf("failed to get passwords: %w", err)
	}

	passwords := make([]Password, len(entries))
	for i, entry := range entries {
		passwords[i] = Password{
			ID:       entry.ID,
			Title:    entry.Title,
			Username: entry.Username,
			Password: entry.Password,
			Notes:    entry.Notes,
		}
	}

	return passwords, nil
}

// CreatePassword adds a new password entry
func (s *PasswordService) CreatePassword(password *Password) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := storage.PasswordEntry{
		Title:    password.Title,
		Username: password.Username,
		Password: password.Password,
		Notes:    password.Notes,
	}

	if err := s.db.AddPassword(entry); err != nil {
		return fmt.Errorf("failed to add password: %w", err)
	}

	return nil
}

// UpdatePassword updates an existing password entry
func (s *PasswordService) UpdatePassword(password *Password) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.db.UpdatePassword(storage.PasswordEntry{
		ID:       password.ID,
		Title:    password.Title,
		Username: password.Username,
		Password: password.Password,
		Notes:    password.Notes,
	}); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeletePassword removes a password entry
func (s *PasswordService) DeletePassword(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.db.DeletePassword(id); err != nil {
		return fmt.Errorf("failed to delete password: %w", err)
	}

	return nil
}
