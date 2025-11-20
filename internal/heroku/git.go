package heroku

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// DetectHerokuApp attempts to detect the Heroku app name from git remotes
func DetectHerokuApp(projectPath string) (appName string, remoteName string, err error) {
	// Get git remotes
	cmd := exec.Command("git", "-C", projectPath, "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("failed to get git remotes: %w (is this a git repository?)", err)
	}

	// Parse remotes looking for heroku.com URLs
	// Format: remotename	https://git.heroku.com/appname.git (fetch)
	// or:    remotename	git@heroku.com:appname.git (fetch)
	lines := strings.Split(string(output), "\n")

	herokuRegex := regexp.MustCompile(`^(\S+)\s+(?:https://git\.heroku\.com/|git@heroku\.com:)([^.\s]+)(?:\.git)?\s+\(fetch\)`)

	for _, line := range lines {
		matches := herokuRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			remoteName = matches[1]
			appName = matches[2]
			return appName, remoteName, nil
		}
	}

	return "", "", fmt.Errorf("no Heroku git remote found")
}

// GetGitRemotes returns all git remotes for the project
func GetGitRemotes(projectPath string) (map[string]string, error) {
	cmd := exec.Command("git", "-C", projectPath, "remote", "-v")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git remotes: %w", err)
	}

	remotes := make(map[string]string)
	lines := strings.Split(string(output), "\n")

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			remotes[parts[0]] = parts[1]
		}
	}

	return remotes, nil
}

// IsGitRepository checks if the given path is a git repository
func IsGitRepository(projectPath string) bool {
	cmd := exec.Command("git", "-C", projectPath, "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}
