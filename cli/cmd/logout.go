package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/groo-dev/cl-wrangler/cli/internal/store"
	"github.com/groo-dev/cl-wrangler/cli/internal/wrangler"
	"github.com/spf13/cobra"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from current account",
	Long:  `Runs wrangler logout and removes the current account from saved accounts.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	if db.Current == "" {
		return fmt.Errorf("no current account set")
	}

	acc := db.GetAccount(db.Current)
	if acc == nil {
		return fmt.Errorf("current account not found in database")
	}

	var confirm bool
	err = huh.NewConfirm().
		Title(fmt.Sprintf("Logout from '%s' and remove from saved accounts?", acc.Name)).
		Value(&confirm).
		Run()

	if err != nil || !confirm {
		fmt.Println("Cancelled.")
		return nil
	}

	// Get wrangler command
	wranglerCmd, err := wrangler.EnsureWranglerCmd(db)
	if err != nil {
		return fmt.Errorf("failed to find wrangler: %w", err)
	}

	// Run wrangler logout
	fmt.Println("Running wrangler logout...")
	if err := wrangler.Logout(wranglerCmd); err != nil {
		return fmt.Errorf("wrangler logout failed: %w", err)
	}

	// Remove from our storage
	if err := store.DeleteAccountConfig(db.Current); err != nil {
		// Not fatal - file might already be gone
		fmt.Printf("Warning: could not delete config file: %v\n", err)
	}

	accountName := acc.Name
	db.RemoveAccount(db.Current)

	if err := store.SaveDB(db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	color.Green("✓ Logged out and removed: %s", accountName)

	if len(db.Accounts) > 0 {
		fmt.Println("\nRemaining accounts:")
		for _, a := range db.Accounts {
			fmt.Printf("  • %s (%s)\n", a.Name, a.Email)
		}
		fmt.Println("\nUse 'cl switch' to login to another account.")
	}

	return nil
}
