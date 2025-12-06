package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/fatih/color"
	"github.com/groo-dev/cl-wrangler/cli/internal/store"
	"github.com/groo-dev/cl-wrangler/cli/internal/wrangler"
	"github.com/sahilm/fuzzy"
	"github.com/spf13/cobra"
)

const addNewAccountOption = "__add_new__"

var switchCmd = &cobra.Command{
	Use:   "switch [account-name-or-id]",
	Short: "Switch to a saved account",
	Long: `Switch to a saved Cloudflare/Wrangler account.
If no argument is provided, shows an interactive list to select from.
Supports fuzzy matching for account names and IDs.`,
	RunE:              runSwitch,
	ValidArgsFunction: completeAccountNames,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

// completeAccountNames provides shell completion for account names
func completeAccountNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	db, err := store.LoadDB()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var completions []string
	for _, acc := range db.Accounts {
		completions = append(completions, acc.Name)
		completions = append(completions, acc.ID)
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

func runSwitch(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	var targetID string

	if len(args) == 0 {
		// Interactive selection
		targetID, err = selectAccountInteractive(db)
		if err != nil {
			return err
		}

		// Handle "Add new account" option
		if targetID == addNewAccountOption {
			return addNewAccount(db)
		}
	} else {
		if len(db.Accounts) == 0 {
			return fmt.Errorf("no accounts saved. Use 'cl add' or 'cl switch' to add an account")
		}
		// Fuzzy match
		query := strings.Join(args, " ")
		targetID, err = findAccountFuzzy(db, query)
		if err != nil {
			return err
		}
	}

	// Save current account before switching (if changed)
	if db.Current != "" && db.Current != targetID {
		currentAcc := db.GetAccount(db.Current)
		if currentAcc != nil {
			changed, newHash, err := store.SaveAccountConfigIfChanged(db.Current, currentAcc.ConfigHash)
			if err == nil && changed {
				currentAcc.ConfigHash = newHash
				db.AddAccount(*currentAcc)
			}
		}
	}

	// Switch to the account
	if err := store.RestoreAccountConfig(targetID); err != nil {
		return fmt.Errorf("failed to restore account config: %w", err)
	}

	db.Current = targetID
	if err := store.SaveDB(db); err != nil {
		return fmt.Errorf("failed to save database: %w", err)
	}

	acc := db.GetAccount(targetID)
	color.Green("✓ Switched to: %s (%s)", acc.Name, acc.Email)

	return nil
}

func selectAccountInteractive(db *store.AccountsDB) (string, error) {
	var options []huh.Option[string]

	for _, acc := range db.Accounts {
		label := fmt.Sprintf("%s (%s)", acc.Name, acc.Email)
		if acc.ID == db.Current {
			label = fmt.Sprintf("%s (%s) [current]", acc.Name, acc.Email)
		}
		options = append(options, huh.NewOption(label, acc.ID))
	}

	// Add "Login with new account" option
	options = append(options, huh.NewOption("+ Login with new account", addNewAccountOption))

	var selected string
	err := huh.NewSelect[string]().
		Title("Select account:").
		Options(options...).
		Value(&selected).
		Run()

	if err != nil {
		return "", err
	}

	return selected, nil
}

func addNewAccount(db *store.AccountsDB) error {
	// Ensure we have a working wrangler command
	wranglerCmd, err := wrangler.EnsureWranglerCmd(db)
	if err != nil {
		return fmt.Errorf("failed to find wrangler: %w", err)
	}

	// Save current account first (if there is one and changed)
	if db.Current != "" {
		currentAcc := db.GetAccount(db.Current)
		if currentAcc != nil {
			changed, newHash, err := store.SaveAccountConfigIfChanged(db.Current, currentAcc.ConfigHash)
			if err == nil && changed {
				fmt.Println("Saving current account before login...")
				currentAcc.ConfigHash = newHash
				db.AddAccount(*currentAcc)
				store.SaveDB(db)
			}
		}
	}

	// Run wrangler login
	fmt.Println("Opening browser for Cloudflare login...")
	if err := wrangler.Login(wranglerCmd); err != nil {
		return fmt.Errorf("wrangler login failed: %w", err)
	}

	// Get new account info
	fmt.Println("Getting new account info...")
	info, err := wrangler.Whoami(wranglerCmd)
	if err != nil {
		return fmt.Errorf("failed to get account info after login: %w", err)
	}

	// Save the new account
	configHash, err := store.SaveAccountConfig(info.AccountID)
	if err != nil {
		return fmt.Errorf("failed to save account config: %w", err)
	}

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

	color.Green("✓ Logged in and saved: %s (%s)", info.AccountName, info.Email)

	return nil
}

// accountSearchable implements fuzzy.Source for accounts
type accountSearchable struct {
	accounts []store.Account
}

func (a accountSearchable) String(i int) string {
	acc := a.accounts[i]
	return fmt.Sprintf("%s %s %s", acc.Name, acc.Email, acc.ID)
}

func (a accountSearchable) Len() int {
	return len(a.accounts)
}

func findAccountFuzzy(db *store.AccountsDB, query string) (string, error) {
	source := accountSearchable{accounts: db.Accounts}
	matches := fuzzy.FindFrom(query, source)

	if len(matches) == 0 {
		return "", fmt.Errorf("no account found matching: %s", query)
	}

	// Return the best match
	return db.Accounts[matches[0].Index].ID, nil
}
