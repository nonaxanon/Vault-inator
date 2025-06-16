package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

var (
	config     *Config
	configOnce sync.Once
	configLock sync.RWMutex
)

// Config holds all configuration values
type Config struct {
	MasterPasswordHash string `json:"master_password_hash"`
	Salt               string `json:"salt"`
	MasterPassword     string `json:"-"` // Not stored in config file
}

// GetConfig returns the singleton config instance
func GetConfig() *Config {
	configOnce.Do(func() {
		config = &Config{}
		if err := config.load(); err != nil {
			panic(fmt.Sprintf("Failed to load config: %v", err))
		}
	})
	return config
}

// LoadConfig loads the configuration from environment variables and config file
func LoadConfig() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist
	_ = godotenv.Load()

	// Get master password from environment variable
	masterPassword := os.Getenv("MASTER_PASSWORD")
	if masterPassword == "" {
		return nil, fmt.Errorf("MASTER_PASSWORD environment variable is required")
	}

	config := GetConfig()
	config.MasterPassword = masterPassword
	return config, nil
}

// load reads the configuration from config file
func (c *Config) load() error {
	// Get the config file path
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".vaultinator", "config.json")
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Read existing config if it exists
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, c); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	configLock.Lock()
	defer configLock.Unlock()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".vaultinator", "config.json")
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// UpdateMasterPassword updates the master password hash and salt
func (c *Config) UpdateMasterPassword(hash, salt string) error {
	configLock.Lock()
	defer configLock.Unlock()

	c.MasterPasswordHash = hash
	c.Salt = salt
	return c.Save()
}

// GetMasterPasswordHash returns the stored master password hash
func (c *Config) GetMasterPasswordHash() string {
	configLock.RLock()
	defer configLock.RUnlock()
	return c.MasterPasswordHash
}

// GetSalt returns the stored salt
func (c *Config) GetSalt() string {
	configLock.RLock()
	defer configLock.RUnlock()
	return c.Salt
}
