package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/groo-dev/cl-wranger/cli/internal/store"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current account",
	Long:  `Shows the currently active Cloudflare/Wrangler account.`,
	RunE:  runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	if db.Current == "" {
		fmt.Println("No current account set. Use 'cl add' to save your current wrangler account.")
		return nil
	}

	acc := db.GetAccount(db.Current)
	if acc == nil {
		fmt.Println("Current account not found in database. Use 'cl add' to re-add it.")
		return nil
	}

	color.Green("Current account:")
	fmt.Printf("  Name:  %s\n", acc.Name)
	fmt.Printf("  Email: %s\n", acc.Email)
	fmt.Printf("  ID:    %s\n", acc.ID)

	return nil
}
