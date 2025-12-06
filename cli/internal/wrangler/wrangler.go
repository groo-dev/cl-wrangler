package wrangler

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/groo-dev/cl-wrangler/cli/internal/store"
)

type WhoamiInfo struct {
	Email       string
	AccountID   string
	AccountName string
}

// DetectWrangler tries to find wrangler and returns the command to use
// It checks: wrangler, npx wrangler
// If not found, prompts user
func DetectWrangler() (string, error) {
	// Try direct wrangler first
	if tryWranglerCmd("wrangler") {
		return "wrangler", nil
	}

	// Try npx wrangler
	if tryWranglerCmd("npx wrangler") {
		return "npx wrangler", nil
	}

	// Not found - ask user
	return promptForWrangler()
}

func tryWranglerCmd(cmd string) bool {
	parts := strings.Fields(cmd)
	args := append(parts[1:], "--version")

	c := exec.Command(parts[0], args...)
	err := c.Run()
	return err == nil
}

func promptForWrangler() (string, error) {
	var choice string

	err := huh.NewSelect[string]().
		Title("Wrangler not found. How would you like to run wrangler?").
		Options(
			huh.NewOption("Use npx (downloads wrangler on demand)", "npx"),
			huh.NewOption("Enter custom path", "custom"),
		).
		Value(&choice).
		Run()

	if err != nil {
		return "", err
	}

	if choice == "npx" {
		return "npx wrangler", nil
	}

	// Custom path
	var customPath string
	err = huh.NewInput().
		Title("Enter wrangler binary path:").
		Value(&customPath).
		Run()

	if err != nil {
		return "", err
	}

	// Verify it works
	if !tryWranglerCmd(customPath) {
		return "", fmt.Errorf("could not run wrangler at: %s", customPath)
	}

	return customPath, nil
}

// EnsureWranglerCmd makes sure we have a working wrangler command configured
func EnsureWranglerCmd(db *store.AccountsDB) (string, error) {
	// If already configured and works, use it
	if db.Settings.WranglerCmd != "" {
		if tryWranglerCmd(db.Settings.WranglerCmd) {
			return db.Settings.WranglerCmd, nil
		}
		// Configured but doesn't work anymore, re-detect
		fmt.Printf("Configured wrangler command '%s' no longer works. Re-detecting...\n", db.Settings.WranglerCmd)
	}

	// Detect or prompt
	cmd, err := DetectWrangler()
	if err != nil {
		return "", err
	}

	// Save to settings
	db.Settings.WranglerCmd = cmd
	if err := store.SaveDB(db); err != nil {
		return "", fmt.Errorf("failed to save wrangler command: %w", err)
	}

	return cmd, nil
}

// Whoami runs wrangler whoami and parses the output
func Whoami(wranglerCmd string) (*WhoamiInfo, error) {
	parts := strings.Fields(wranglerCmd)
	args := append(parts[1:], "whoami")

	cmd := exec.Command(parts[0], args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run %s whoami: %w\nOutput: %s", wranglerCmd, err, string(output))
	}

	return parseWhoamiOutput(string(output))
}

func parseWhoamiOutput(output string) (*WhoamiInfo, error) {
	info := &WhoamiInfo{}

	// Parse email: "associated with the email user@example.com"
	emailRegex := regexp.MustCompile(`associated with the email (\S+)`)
	if match := emailRegex.FindStringSubmatch(output); len(match) > 1 {
		info.Email = strings.TrimSuffix(match[1], ".")
	}

	// Parse table row: │ AccountName │ account_id │
	rowRegex := regexp.MustCompile(`│\s*([^│]+?)\s*│\s*([^│]+?)\s*│`)
	matches := rowRegex.FindAllStringSubmatch(output, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		name := strings.TrimSpace(match[1])
		id := strings.TrimSpace(match[2])

		// Skip header row
		if name == "Account Name" && id == "Account ID" {
			continue
		}

		// Found data row
		if name != "" && id != "" {
			info.AccountName = name
			info.AccountID = id
			break
		}
	}

	if info.AccountID == "" {
		return nil, fmt.Errorf("could not parse account ID from wrangler whoami output")
	}

	return info, nil
}

// Login runs wrangler login interactively
func Login(wranglerCmd string) error {
	parts := strings.Fields(wranglerCmd)
	args := append(parts[1:], "login")

	cmd := exec.Command(parts[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Logout runs wrangler logout
func Logout(wranglerCmd string) error {
	parts := strings.Fields(wranglerCmd)
	args := append(parts[1:], "logout")

	cmd := exec.Command(parts[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
