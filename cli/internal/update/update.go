package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	repoOwner     = "groo-dev"
	repoName      = "cl-wrangler"
	checkInterval = 24 * time.Hour
	versionAPI    = "https://ops.groo.dev/v1/webhook/version?environment=production"
	apiToken      = "groo_b310854d80189784e3ed222a5860562f992587bf8b1f34d6677d0c1857812461"
)

type versionResponse struct {
	Version string `json:"version"`
}

// CheckForUpdate checks for a newer version
// Returns (newVersion, downloadURL, error)
func CheckForUpdate(currentVersion string) (string, string, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest("GET", versionAPI, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("failed to fetch version: status %d", resp.StatusCode)
	}

	var versionResp versionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versionResp); err != nil {
		return "", "", err
	}

	latestVersion := versionResp.Version

	if isNewerVersion(latestVersion, currentVersion) {
		releaseURL := fmt.Sprintf("https://github.com/%s/%s/releases/tag/v%s", repoOwner, repoName, latestVersion)
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
