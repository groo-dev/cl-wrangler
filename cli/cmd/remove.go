package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/groo-dev/cl-wrangler/cli/internal/store"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove [account-name-or-id]",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a saved account",
	Long: `Remove a saved Cloudflare/Wrangler account.
If no argument is provided, shows an interactive list to select from.
Supports fuzzy matching for account names and IDs.`,
	RunE:              runRemove,
	ValidArgsFunction: completeAccountNames,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	if len(db.Accounts) == 0 {
		return fmt.Errorf("no accounts saved")
	}

	var targetID string

	if len(args) == 0 {
		// Interactive selection
		targetID, err = selectAccountForRemoval(db)
		if err != nil {
			return err
		}
	} else {
		// Fuzzy match
		query := strings.Join(args, " ")
		targetID, err = findAccountForRemoval(db, query)
		if err != nil {
			return err
		}
	}

	acc := db.GetAccount(targetID)
	if acc == nil {
		return fmt.Errorf("account not found")
	}

	// Confirm deletion
	var confirm bool
	err = huh.NewConfirm().
		Title(fmt.Sprintf("Remove account '%s' (%s)?", acc.Name, acc.Email)).
		Value(&confirm).
		Run()

	if err != nil {
		return err
	}

	if !confirm {
		fmt.Println("Cancelled.")
		return nil
	}

	// Delete config file
	if err := store.DeleteAccountConfig(targetID); err != nil {
		return fmt.Errorf("failed to delete account config: %w", err)
	}

	// Remove from database
	db.RemoveAccount(targetID)
	if err := store.SaveDB(db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	color.Green("✓ Removed: %s (%s)", acc.Name, acc.Email)

	return nil
}

func selectAccountForRemoval(db *store.AccountsDB) (string, error) {
	var options []huh.Option[string]

	for _, acc := range db.Accounts {
		label := fmt.Sprintf("%s (%s)", acc.Name, acc.Email)
		options = append(options, huh.NewOption(label, acc.ID))
	}

	var selected string
	err := huh.NewSelect[string]().
		Title("Select account to remove:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", err
	}

	return selected, nil
}

func findAccountForRemoval(db *store.AccountsDB, query string) (string, error) {
	source := accountSearchable{accounts: db.Accounts}
	matches := fuzzy.FindFrom(query, source)

	if len(matches) == 0 {
		return "", fmt.Errorf("no account found matching: %s", query)
	}

	return db.Accounts[matches[0].Index].ID, nil
}
