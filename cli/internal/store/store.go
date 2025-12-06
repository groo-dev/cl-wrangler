package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/groo-dev/cl-wranger/cli/internal/config"
)

type Account struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	AddedAt    time.Time `json:"added_at"`
	ConfigHash string    `json:"config_hash,omitempty"`
}

type Settings struct {
	WranglerCmd     string    `json:"wrangler_cmd"`
	LastUpdateCheck time.Time `json:"last_update_check,omitempty"`
}

type AccountsDB struct {
	Accounts []Account `json:"accounts"`
	Current  string    `json:"current"`
	Settings Settings  `json:"settings"`
}

// LoadDB loads the accounts database from disk
func LoadDB() (*AccountsDB, error) {
	dbPath, err := config.GetAccountsDBPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(dbPath)
	if os.IsNotExist(err) {
		return &AccountsDB{Accounts: []Account{}}, nil
	}
	if err != nil {
		return nil, err
	}

	var db AccountsDB
	if err := json.Unmarshal(data, &db); err != nil {
		return nil, err
	}
	return &db, nil
}

// SaveDB saves the accounts database to disk
func SaveDB(db *AccountsDB) error {
	if err := config.EnsureConfigDirs(); err != nil {
		return err
	}

	dbPath, err := config.GetAccountsDBPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dbPath, data, 0644)
}

// AddAccount adds or updates an account in the database
func (db *AccountsDB) AddAccount(account Account) {
	for i, a := range db.Accounts {
		if a.ID == account.ID {
			db.Accounts[i] = account
			return
		}
	}
	db.Accounts = append(db.Accounts, account)
}

// RemoveAccount removes an account from the database
func (db *AccountsDB) RemoveAccount(id string) {
	for i, a := range db.Accounts {
		if a.ID == id {
			db.Accounts = append(db.Accounts[:i], db.Accounts[i+1:]...)
			if db.Current == id {
				db.Current = ""
			}
			return
		}
	}
}

// GetAccount finds an account by ID
func (db *AccountsDB) GetAccount(id string) *Account {
	for _, a := range db.Accounts {
		if a.ID == id {
			return &a
		}
	}
	return nil
}

// HashFile computes SHA256 hash of a file
func HashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// GetCurrentConfigHash returns the hash of wrangler's current default.toml
func GetCurrentConfigHash() (string, error) {
	path, err := config.GetWranglerConfigPath()
	if err != nil {
		return "", err
	}
	return HashFile(path)
}

// SaveAccountConfig copies the current wrangler config to our storage and returns the hash
func SaveAccountConfig(accountID string) (string, error) {
	if err := config.EnsureConfigDirs(); err != nil {
		return "", err
	}

	srcPath, err := config.GetWranglerConfigPath()
	if err != nil {
		return "", err
	}

	accountsDir, err := config.GetAccountsDir()
	if err != nil {
		return "", err
	}

	dstPath := filepath.Join(accountsDir, accountID+".toml")
	if err := copyFile(srcPath, dstPath); err != nil {
		return "", err
	}

	return HashFile(srcPath)
}

// SaveAccountConfigIfChanged saves config only if it has changed, returns (changed, newHash, error)
func SaveAccountConfigIfChanged(accountID string, currentHash string) (bool, string, error) {
	newHash, err := GetCurrentConfigHash()
	if err != nil {
		return false, "", err
	}

	if newHash == currentHash {
		return false, currentHash, nil
	}

	_, err = SaveAccountConfig(accountID)
	if err != nil {
		return false, "", err
	}

	return true, newHash, nil
}

// RestoreAccountConfig copies a saved config back to wrangler's location
func RestoreAccountConfig(accountID string) error {
	accountsDir, err := config.GetAccountsDir()
	if err != nil {
		return err
	}

	srcPath := filepath.Join(accountsDir, accountID+".toml")

	dstPath, err := config.GetWranglerConfigPath()
	if err != nil {
		return err
	}

	return copyFile(srcPath, dstPath)
}

// DeleteAccountConfig removes a saved account config file
func DeleteAccountConfig(accountID string) error {
	accountsDir, err := config.GetAccountsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(accountsDir, accountID+".toml")
	return os.Remove(path)
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
