package heroku

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/egg/heroku-calc/internal/config"
)

// GetAddons retrieves all addons for the app
func (c *Client) GetAddons() ([]config.Addon, error) {
	if c.useCLI {
		return c.getAddonsCLI()
	}
	return c.getAddonsAPI()
}

func (c *Client) getAddonsCLI() ([]config.Addon, error) {
	cmd := exec.Command("heroku", "addons", "-a", c.appName, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get addons via CLI: %w", err)
	}

	var addons []struct {
		Name      string    `json:"name"`
		PlanName  string    `json:"plan_name"`
		CreatedAt time.Time `json:"created_at"`
		// Note: Price is not directly available from heroku addons command
		// We'll need to look it up from pricing data
	}

	if err := json.Unmarshal(output, &addons); err != nil {
		return nil, fmt.Errorf("failed to parse addons: %w", err)
	}

	result := make([]config.Addon, len(addons))
	for i, addon := range addons {
		result[i] = config.Addon{
			Name:    addon.Name,
			Plan:    addon.PlanName,
			Price:   "unknown", // Will be populated from pricing data
			AddedAt: addon.CreatedAt,
		}
	}

	return result, nil
}

func (c *Client) getAddonsAPI() ([]config.Addon, error) {
	// TODO: Implement direct API call
	return nil, fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}

// AddonPlanInfo contains detailed information about an addon plan
type AddonPlanInfo struct {
	Service string
	Plan    string
	Limits  map[string]interface{}
}

// GetAddonPlanInfo retrieves detailed information about a specific addon plan
func (c *Client) GetAddonPlanInfo(addonName string) (*AddonPlanInfo, error) {
	if c.useCLI {
		return c.getAddonPlanInfoCLI(addonName)
	}
	return c.getAddonPlanInfoAPI(addonName)
}

func (c *Client) getAddonPlanInfoCLI(addonName string) (*AddonPlanInfo, error) {
	cmd := exec.Command("heroku", "addons:info", addonName, "-a", c.appName, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get addon info via CLI: %w", err)
	}

	var info struct {
		AddonService struct {
			Name string `json:"name"`
		} `json:"addon_service"`
		Plan struct {
			Name string `json:"name"`
		} `json:"plan"`
		ConfigVars []string `json:"config_vars"`
	}

	if err := json.Unmarshal(output, &info); err != nil {
		return nil, fmt.Errorf("failed to parse addon info: %w", err)
	}

	return &AddonPlanInfo{
		Service: info.AddonService.Name,
		Plan:    info.Plan.Name,
		Limits:  make(map[string]interface{}),
	}, nil
}

func (c *Client) getAddonPlanInfoAPI(addonName string) (*AddonPlanInfo, error) {
	// TODO: Implement direct API call
	return nil, fmt.Errorf("API mode not yet implemented - please install Heroku CLI")
}
