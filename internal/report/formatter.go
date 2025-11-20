package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// GenerateFileName creates a timestamped filename for the report
func GenerateFileName(appName string) string {
	timestamp := time.Now().Format("2006-01-02")
	return fmt.Sprintf("heroku-analysis-%s-%s.md", appName, timestamp)
}

// Save writes the report content to a file
func Save(content, filename string) error {
	// Ensure the file has .md extension
	if filepath.Ext(filename) != ".md" {
		filename += ".md"
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write report: %w", err)
	}

	return nil
}

// SaveInProjectDir writes the report to the project directory
func SaveInProjectDir(content, projectPath, appName string) (string, error) {
	filename := GenerateFileName(appName)
	fullPath := filepath.Join(projectPath, filename)

	if err := Save(content, fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
