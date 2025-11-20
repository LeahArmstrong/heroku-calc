package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".heroku-calc.yml"

// Load reads the configuration from the specified project path
func Load(projectPath string) (*Config, error) {
	configPath := filepath.Join(projectPath, ConfigFileName)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file not found: %s (run with --init to create)", configPath)
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set project path if not already set
	if cfg.ProjectPath == "" {
		cfg.ProjectPath = projectPath
	}

	return &cfg, nil
}

// Exists checks if a config file exists at the given path
func Exists(projectPath string) bool {
	configPath := filepath.Join(projectPath, ConfigFileName)
	_, err := os.Stat(configPath)
	return err == nil
}

// GetConfigPath returns the full path to the config file
func GetConfigPath(projectPath string) string {
	return filepath.Join(projectPath, ConfigFileName)
}
