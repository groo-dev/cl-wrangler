package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/groo-dev/cl-wranger/cli/internal/config"
	"github.com/groo-dev/cl-wranger/cli/internal/store"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit configuration",
	Long:  `View current configuration or edit settings like wrangler command path.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) error {
	db, err := store.LoadDB()
	if err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	configDir, _ := config.GetConfigDir()
	wranglerPath, _ := config.GetWranglerConfigPath()

	fmt.Println("Current configuration:")
	fmt.Printf("  Config directory:  %s\n", configDir)
	fmt.Printf("  Wrangler config:   %s\n", wranglerPath)
	fmt.Printf("  Wrangler command:  %s\n", db.Settings.WranglerCmd)
	fmt.Printf("  Saved accounts:    %d\n", len(db.Accounts))

	var edit bool
	err = huh.NewConfirm().
		Title("Edit wrangler command?").
		Value(&edit).
		Run()

	if err != nil || !edit {
		return nil
	}

	var newCmd string
	err = huh.NewInput().
		Title("Wrangler command:").
		Value(&newCmd).
		Placeholder(db.Settings.WranglerCmd).
		Run()

	if err != nil {
		return err
	}

	if newCmd != "" && newCmd != db.Settings.WranglerCmd {
		db.Settings.WranglerCmd = newCmd
		if err := store.SaveDB(db); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Printf("Updated wrangler command to: %s\n", newCmd)
	}

	return nil
}
