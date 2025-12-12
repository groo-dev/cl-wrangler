package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/groo-dev/cl-wrangler/cli/internal/store"
	"github.com/groo-dev/cl-wrangler/cli/internal/update"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cl",
	Short:   "Cloudflare Wrangler account switcher",
	Long:    `A CLI tool to easily switch between multiple Cloudflare/Wrangler accounts.`,
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip update check for version and completion commands
		if cmd.Name() == "version" || cmd.Name() == "completion" {
			return
		}
		checkForUpdates()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func checkForUpdates() {
	db, err := store.LoadDB()
	if err != nil {
		return
	}

	// Only check once per day
	if !update.ShouldCheck(db.Settings.LastUpdateCheck) {
		return
	}

	// Update last check time
	db.Settings.LastUpdateCheck = time.Now()
	store.SaveDB(db)

	// Check for updates in background (don't block CLI)
	go func() {
		newVersion, url, err := update.CheckForUpdate(Version)
		if err != nil || newVersion == "" {
			return
		}

		// Print update notice
		fmt.Println()
		color.Yellow("A new version of cl is available: v%s → v%s", Version, newVersion)
		fmt.Printf("Download: %s\n\n", url)
	}()

	// Small delay to allow update message to print
	time.Sleep(100 * time.Millisecond)
}
