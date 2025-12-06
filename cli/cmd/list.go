package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/groo-dev/cl-wranger/cli/internal/store"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List saved accounts",
	Long:    `Lists all saved Cloudflare/Wrangler accounts.`,
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	if len(db.Accounts) == 0 {
		fmt.Println("No accounts saved. Use 'cl add' to save your current wrangler account.")
		return nil
	}

	green := color.New(color.FgGreen)

	for _, acc := range db.Accounts {
		if acc.ID == db.Current {
			green.Printf("→ %s (%s)\n", acc.Name, acc.Email)
			green.Printf("  %s\n", acc.ID)
		} else {
			fmt.Printf("  %s (%s)\n", acc.Name, acc.Email)
			fmt.Printf("  %s\n", acc.ID)
		}
	}

	return nil
}
