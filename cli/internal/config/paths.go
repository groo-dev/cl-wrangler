package config

import (
	"os"
	"path/filepath"
)

// GetConfigDir returns the cl-wrangler config directory path
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "cl-wrangler"), nil
}

// GetAccountsDir returns the directory where account TOML files are stored
func GetAccountsDir() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "accounts"), nil
}

// GetAccountsDBPath returns the path to accounts.json
func GetAccountsDBPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "accounts.json"), nil
}

// GetWranglerConfigPath returns the path to wrangler's default.toml
// On macOS, wrangler uses ~/Library/Preferences/.wrangler/config/default.toml
func GetWranglerConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, "Library", "Preferences", ".wrangler", "config", "default.toml"), nil
}

// EnsureConfigDirs creates the config directories if they don't exist
func EnsureConfigDirs() error {
	accountsDir, err := GetAccountsDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(accountsDir, 0755)
}
