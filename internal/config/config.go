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
	// Try to load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	// Get configuration from environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("Warning: DATABASE_URL not set, using default")
		dbURL = "postgres://admin:admin@192.168.1.8:5432/dev?sslmode=disable"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		fmt.Println("Warning: CONFIG_PATH not set, using default")
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		configPath = filepath.Join(homeDir, ".vaultinator", "config.json")
	}

	// Set environment variables for other parts of the application
	os.Setenv("DATABASE_URL", dbURL)
	os.Setenv("CONFIG_PATH", configPath)

	config := GetConfig()

	// Only set master password if it's provided
	if masterPassword := os.Getenv("MASTER_PASSWORD"); masterPassword != "" {
		fmt.Println("Info: MASTER_PASSWORD found in environment")
		config.MasterPassword = masterPassword
	} else {
		fmt.Println("Info: MASTER_PASSWORD not set, will be initialized through UI")
	}

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

// save writes the configuration to the config file without locking
func (c *Config) save() error {
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

// Save writes the configuration to the config file with locking
func (c *Config) Save() error {
	configLock.Lock()
	defer configLock.Unlock()
	return c.save()
}

// UpdateMasterPassword updates the master password hash and salt
func (c *Config) UpdateMasterPassword(hash, salt string) error {
	configLock.Lock()
	defer configLock.Unlock()

	c.MasterPasswordHash = hash
	c.Salt = salt
	return c.save() // Use save() instead of Save() to avoid double locking
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
