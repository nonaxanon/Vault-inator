package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	MasterPassword string `json:"master_password"`
}

// LoadConfig loads the configuration from a file
func LoadConfig() (*Config, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Create the config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".vaultinator")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return nil, err
	}

	// Path to the config file
	configPath := filepath.Join(configDir, "config.json")

	// Check if a config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		cfg := &Config{
			MasterPassword: "changeme", // Default master password
		}

		// Save the default config
		if err := cfg.SaveConfig(); err != nil {
			return nil, err
		}

		return cfg, nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig saves the configuration to a file
func (c *Config) SaveConfig() error {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Path to the config file
	configPath := filepath.Join(homeDir, ".vaultinator", "config.json")

	// Marshal the config to JSON
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	// Write the config file
	return os.WriteFile(configPath, data, 0600)
}
