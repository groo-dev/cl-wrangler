package cmd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/groo-dev/cl-wranger/cli/internal/store"
	"github.com/groo-dev/cl-wranger/cli/internal/wrangler"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Save current wrangler account",
	Long:  `Saves the current wrangler authentication as a profile that can be switched to later.`,
	RunE:  runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	// Ensure we have a working wrangler command
	wranglerCmd, err := wrangler.EnsureWranglerCmd(db)
	if err != nil {
		return fmt.Errorf("failed to find wrangler: %w", err)
	}

	// Get current account info from wrangler
	fmt.Println("Getting account info from wrangler...")
	info, err := wrangler.Whoami(wranglerCmd)
	if err != nil {
		return err
	}

	// Check if account already exists
	existing := db.GetAccount(info.AccountID)
	if existing != nil {
		fmt.Printf("Account '%s' already saved. Updating...\n", existing.Name)
	}

	// Save the config file
	configHash, err := store.SaveAccountConfig(info.AccountID)
	if err != nil {
		return fmt.Errorf("failed to save account config: %w", err)
	}

	// Add to database
	account := store.Account{
		ID:         info.AccountID,
		Name:       info.AccountName,
		Email:      info.Email,
		AddedAt:    time.Now(),
		ConfigHash: configHash,
	}
	db.AddAccount(account)
	db.Current = info.AccountID

	if err := store.SaveDB(db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	color.Green("✓ Account saved: %s (%s)", info.AccountName, info.Email)
	color.Cyan("  Account ID: %s", info.AccountID)

	return nil
}
