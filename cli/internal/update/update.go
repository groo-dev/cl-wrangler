package update

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	repoOwner     = "groo-dev"
	repoName      = "cl-wrangler"
	checkInterval = 24 * time.Hour
	versionFile   = "cli/VERSION"
)

// CheckForUpdate checks GitHub for a newer version
// Returns (newVersion, downloadURL, error)
func CheckForUpdate(currentVersion string) (string, string, error) {
	// Fetch VERSION file from repo
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/main/%s", repoOwner, repoName, versionFile)

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch VERSION file: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	latestVersion := strings.TrimSpace(string(body))

	if isNewerVersion(latestVersion, currentVersion) {
		releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/cli-v%s", repoOwner, repoName, latestVersion)
		return latestVersion, releaseURL, nil
	}

	return "", "", nil
}

// isNewerVersion compares semantic versions using proper semver parsing
func isNewerVersion(latest, current string) bool {
	// Skip check for dev versions
	if current == "dev" || current == "" {
		return false
	}

	currentVer, err := semver.NewVersion(current)
	if err != nil {
		return false
	}

	latestVer, err := semver.NewVersion(latest)
	if err != nil {
		return false
	}

	return latestVer.GreaterThan(currentVer)
}

// ShouldCheck returns true if enough time has passed since last check
func ShouldCheck(lastCheck time.Time) bool {
	return time.Since(lastCheck) > checkInterval
}
