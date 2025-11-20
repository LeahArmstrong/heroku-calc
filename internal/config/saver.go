package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Save writes the configuration to disk
func Save(cfg *Config, projectPath string) error {
	// Update timestamp
	cfg.LastUpdated = time.Now()

	// Ensure project path is set
	if cfg.ProjectPath == "" {
		cfg.ProjectPath = projectPath
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(projectPath, ConfigFileName)

	// Write with user-only read/write permissions
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// New creates a new empty configuration
func New(appName string, projectPath string) *Config {
	return &Config{
		AppName:         appName,
		SafeEnvVars:     []string{},
		ExcludedEnvVars: []string{},
		LastUpdated:     time.Now(),
		ProjectPath:     projectPath,
	}
}

// AddSafeEnvVar adds an environment variable to the safe list
func (c *Config) AddSafeEnvVar(name string) {
	// Check if already exists
	for _, v := range c.SafeEnvVars {
		if v == name {
			return
		}
	}
	c.SafeEnvVars = append(c.SafeEnvVars, name)
}

// RemoveSafeEnvVar removes an environment variable from the safe list
func (c *Config) RemoveSafeEnvVar(name string) {
	filtered := make([]string, 0, len(c.SafeEnvVars))
	for _, v := range c.SafeEnvVars {
		if v != name {
			filtered = append(filtered, v)
		}
	}
	c.SafeEnvVars = filtered
}

// IsSafe checks if an environment variable is marked as safe
func (c *Config) IsSafe(name string) bool {
	for _, v := range c.SafeEnvVars {
		if v == name {
			return true
		}
	}
	return false
}

// IsExcluded checks if an environment variable is marked as excluded
func (c *Config) IsExcluded(name string) bool {
	for _, v := range c.ExcludedEnvVars {
		if v == name {
			return true
		}
	}
	return false
}
