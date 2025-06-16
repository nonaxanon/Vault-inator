package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/bcrypt"

	"github.com/nonaxanon/vault-inator/internal/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotInitialized     = errors.New("master password not initialized")
)

// AuthService handles master password operations
type AuthService struct {
	passwordService *PasswordService
	mu              sync.RWMutex
}

// NewAuthService creates a new auth service instance
func NewAuthService(passwordService *PasswordService) *AuthService {
	return &AuthService{
		passwordService: passwordService,
	}
}

// InitializeMasterPassword sets up the initial master password
func (s *AuthService) InitializeMasterPassword(password string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate salt
	salt := make([]byte, 16)
	if _, err := os.ReadFile(filepath.Join(os.TempDir(), "salt")); err != nil {
		if err := os.WriteFile(filepath.Join(os.TempDir(), "salt"), salt, 0600); err != nil {
			return fmt.Errorf("failed to write salt: %w", err)
		}
	}

	// Hash password with salt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update config
	cfg := config.GetConfig()
	if err := cfg.UpdateMasterPassword(string(hash), base64.StdEncoding.EncodeToString(salt)); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Set encryption key for password service
	key := sha256.Sum256([]byte(password))
	if err := s.passwordService.SetEncryptionKey(key[:]); err != nil {
		return fmt.Errorf("failed to set encryption key: %w", err)
	}

	return nil
}

// ChangeMasterPassword updates the master password
func (s *AuthService) ChangeMasterPassword(currentPassword, newPassword string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verify current password
	if !s.VerifyMasterPassword(currentPassword) {
		return fmt.Errorf("invalid current password")
	}

	// Generate new salt
	salt := make([]byte, 16)
	if err := os.WriteFile(filepath.Join(os.TempDir(), "salt"), salt, 0600); err != nil {
		return fmt.Errorf("failed to write salt: %w", err)
	}

	// Hash new password with salt
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update config
	cfg := config.GetConfig()
	if err := cfg.UpdateMasterPassword(string(hash), base64.StdEncoding.EncodeToString(salt)); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	// Set new encryption key for password service
	key := sha256.Sum256([]byte(newPassword))
	if err := s.passwordService.SetEncryptionKey(key[:]); err != nil {
		return fmt.Errorf("failed to set encryption key: %w", err)
	}

	return nil
}

// VerifyMasterPassword checks if the provided password matches the stored hash
func (s *AuthService) VerifyMasterPassword(password string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cfg := config.GetConfig()
	hash := cfg.GetMasterPasswordHash()
	if hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsInitialized checks if the master password has been set
func (s *AuthService) IsInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cfg := config.GetConfig()
	return cfg.GetMasterPasswordHash() != ""
}
