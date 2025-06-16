package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/nonaxanon/vault-inator/internal/config"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNotInitialized     = errors.New("master password not initialized")
)

type AuthService struct {
	config *config.Config
}

func NewAuthService() *AuthService {
	return &AuthService{
		config: config.GetConfig(),
	}
}

// InitializeMasterPassword sets up the initial master password
func (s *AuthService) InitializeMasterPassword(password string) error {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password with the salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Store the hash and salt
	if err := s.config.UpdateMasterPassword(string(hashedPassword), base64.StdEncoding.EncodeToString(salt)); err != nil {
		return fmt.Errorf("failed to save master password: %w", err)
	}

	return nil
}

// ChangeMasterPassword updates the master password
func (s *AuthService) ChangeMasterPassword(currentPassword, newPassword string) error {
	// Verify current password
	if !s.VerifyMasterPassword(currentPassword) {
		return ErrInvalidCredentials
	}

	// Generate a new salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update the stored hash and salt
	if err := s.config.UpdateMasterPassword(string(hashedPassword), base64.StdEncoding.EncodeToString(salt)); err != nil {
		return fmt.Errorf("failed to update master password: %w", err)
	}

	return nil
}

// VerifyMasterPassword checks if the provided password matches the stored hash
func (s *AuthService) VerifyMasterPassword(password string) bool {
	storedHash := s.config.GetMasterPasswordHash()
	if storedHash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	return err == nil
}

// IsInitialized checks if the master password has been set up
func (s *AuthService) IsInitialized() bool {
	return s.config.GetMasterPasswordHash() != ""
}
