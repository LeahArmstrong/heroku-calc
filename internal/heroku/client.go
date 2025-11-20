package heroku

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/egg/heroku-calc/internal/config"
)

// Client represents a Heroku API/CLI client
type Client struct {
	appName    string
	useCLI     bool
	apiToken   string
	apiBaseURL string
}

// NewClient creates a new Heroku client
// It will attempt to use the Heroku CLI by default, falling back to API if needed
func NewClient(appName string) (*Client, error) {
	client := &Client{
		appName:    appName,
		useCLI:     true,
		apiBaseURL: "https://api.heroku.com",
	}

	// Check if Heroku CLI is available
	if err := exec.Command("heroku", "version").Run(); err != nil {
		// CLI not available, we'll need API token
		client.useCLI = false
	}

	return client, nil
}

// SetAPIToken sets the API token for direct API calls
func (c *Client) SetAPIToken(token string) {
	c.apiToken = token
	if token != "" {
		c.useCLI = false
	}
}

// IsUsingCLI returns true if the client is configured to use the Heroku CLI
func (c *Client) IsUsingCLI() bool {
	return c.useCLI
}

// GetEnvVars retrieves all environment variables for the app
func (c *Client) GetEnvVars() ([]config.HerokuEnvVar, error) {
	if c.useCLI {
		return c.getEnvVarsCLI()
	}
	return c.getEnvVarsAPI()
}

func (c *Client) getEnvVarsCLI() ([]config.HerokuEnvVar, error) {
	cmd := exec.Command("heroku", "config", "-a", c.appName, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get config vars via CLI: %w", err)
	}

	var vars map[string]string
	if err := json.Unmarshal(output, &vars); err != nil {
		return nil, fmt.Errorf("failed to parse config vars: %w", err)
	}

	result := make([]config.HerokuEnvVar, 0, len(vars))
	for name, value := range vars {
		result = append(result, config.HerokuEnvVar{
			Name:  name,
			Value: value,
		})
	}

	return result, nil
}

func (c *Client) getEnvVarsAPI() ([]config.HerokuEnvVar, error) {
	// TODO: Implement direct API call
	// This would use the Heroku Platform API with the token
	return nil, fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// GetDynos retrieves dyno formation for the app
func (c *Client) GetDynos() ([]config.DynoFormation, error) {
	if c.useCLI {
		return c.getDynosCLI()
	}
	return c.getDynosAPI()
}

func (c *Client) getDynosCLI() ([]config.DynoFormation, error) {
	cmd := exec.Command("heroku", "ps", "-a", c.appName, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get dyno info via CLI: %w", err)
	}

	var dynos []struct {
		Type string `json:"type"`
		Size string `json:"size"`
	}
	if err := json.Unmarshal(output, &dynos); err != nil {
		return nil, fmt.Errorf("failed to parse dyno info: %w", err)
	}

	// Count dynos by type
	formationMap := make(map[string]*config.DynoFormation)
	for _, dyno := range dynos {
		key := dyno.Type
		if formation, exists := formationMap[key]; exists {
			formation.Quantity++
		} else {
			formationMap[key] = &config.DynoFormation{
				Type:     dyno.Type,
				Quantity: 1,
				Size:     dyno.Size,
			}
		}
	}

	result := make([]config.DynoFormation, 0, len(formationMap))
	for _, formation := range formationMap {
		result = append(result, *formation)
	}

	return result, nil
}

func (c *Client) getDynosAPI() ([]config.DynoFormation, error) {
	// TODO: Implement direct API call
	return nil, fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// SetEnvVar sets an environment variable on Heroku
func (c *Client) SetEnvVar(name, value string) error {
	if c.useCLI {
		return c.setEnvVarCLI(name, value)
	}
	return c.setEnvVarAPI(name, value)
}

func (c *Client) setEnvVarCLI(name, value string) error {
	configStr := fmt.Sprintf("%s=%s", name, value)
	cmd := exec.Command("heroku", "config:set", configStr, "-a", c.appName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set config var %s: %w\nOutput: %s", name, err, string(output))
	}
	return nil
}

func (c *Client) setEnvVarAPI(name, value string) error {
	// TODO: Implement direct API call
	return fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// UnsetEnvVar removes an environment variable from Heroku
func (c *Client) UnsetEnvVar(name string) error {
	if c.useCLI {
		return c.unsetEnvVarCLI(name)
	}
	return c.unsetEnvVarAPI(name)
}

func (c *Client) unsetEnvVarCLI(name string) error {
	cmd := exec.Command("heroku", "config:unset", name, "-a", c.appName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to unset config var %s: %w\nOutput: %s", name, err, string(output))
	}
	return nil
}

func (c *Client) unsetEnvVarAPI(name string) error {
	// TODO: Implement direct API call
	return fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// TestConnection tests if we can connect to Heroku
func (c *Client) TestConnection() error {
	if c.useCLI {
		cmd := exec.Command("heroku", "apps:info", "-a", c.appName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to connect to Heroku app %s: %w", c.appName, err)
		}
		return nil
	}

	// TODO: Implement API test
	return fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// GetAppInfo retrieves basic app information
func (c *Client) GetAppInfo() (*AppInfo, error) {
	if c.useCLI {
		return c.getAppInfoCLI()
	}
	return c.getAppInfoAPI()
}

func (c *Client) getAppInfoCLI() (*AppInfo, error) {
	cmd := exec.Command("heroku", "apps:info", "-a", c.appName, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get app info via CLI: %w", err)
	}

	var info struct {
		App struct {
			Name   string `json:"name"`
			Region struct {
				Name string `json:"name"`
			} `json:"region"`
			Stack struct {
				Name string `json:"name"`
			} `json:"stack"`
		} `json:"app"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse app info: %w", err)
	}

	return &AppInfo{
		Name:   info.App.Name,
		Region: info.App.Region.Name,
		Stack:  info.App.Stack.Name,
	}, nil
}

func (c *Client) getAppInfoAPI() (*AppInfo, error) {
	// TODO: Implement direct API call
	return nil, fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// AppInfo contains basic Heroku app information
type AppInfo struct {
	Name   string
	Region string
	Stack  string
}

// SanitizeEnvVarValue masks sensitive parts of environment variable values
func SanitizeEnvVarValue(name, value string) string {
	// For URLs, show only the protocol and host
	if strings.HasPrefix(value, "postgres://") || strings.HasPrefix(value, "redis://") || strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") {
		parts := strings.SplitN(value, "@", 2)
		if len(parts) == 2 {
			protocol := strings.Split(parts[0], "://")[0]
			hostParts := strings.Split(parts[1], "/")
			return fmt.Sprintf("%s://***@%s/***", protocol, hostParts[0])
		}
	}

	// For API keys, tokens, secrets - show only first 4 chars
	if strings.Contains(strings.ToLower(name), "key") ||
		strings.Contains(strings.ToLower(name), "token") ||
		strings.Contains(strings.ToLower(name), "secret") ||
		strings.Contains(strings.ToLower(name), "password") {
		if len(value) > 4 {
			return value[:4] + strings.Repeat("*", len(value)-4)
		}
		return strings.Repeat("*", len(value))
	}

	// For other vars, show as-is if short, otherwise truncate
	if len(value) > 50 {
		return value[:47] + "..."
	}
	return value
}
