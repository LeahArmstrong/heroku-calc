package analysis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/leaharmstrong/heroku-calc/internal/config"
	"github.com/leaharmstrong/heroku-calc/internal/heroku"
	"github.com/leaharmstrong/heroku-calc/internal/pricing"
)

// Analyzer performs configuration analysis
type Analyzer struct {
	client      *heroku.Client
	pricingData *pricing.Data
	envVars     map[string]string
	dynos       []config.DynoFormation
	addons      []config.Addon
}

// NewAnalyzer creates a new analyzer instance
func NewAnalyzer(client *heroku.Client, pricingData *pricing.Data) *Analyzer {
	return &Analyzer{
		client:      client,
		pricingData: pricingData,
		envVars:     make(map[string]string),
	}
}

// LoadData loads all necessary data from Heroku
func (a *Analyzer) LoadData() error {
	// Load environment variables
	envVars, err := a.client.GetEnvVars()
	if err != nil {
		return fmt.Errorf("failed to load env vars: %w", err)
	}

	a.envVars = make(map[string]string)
	for _, ev := range envVars {
		a.envVars[ev.Name] = ev.Value
	}

	// Load dynos
	dynos, err := a.client.GetDynos()
	if err != nil {
		return fmt.Errorf("failed to load dynos: %w", err)
	}
	a.dynos = dynos

	// Load addons
	addons, err := a.client.GetAddons()
	if err != nil {
		return fmt.Errorf("failed to load addons: %w", err)
	}
	a.addons = addons

	return nil
}

// Analyze performs comprehensive analysis
func (a *Analyzer) Analyze() (*config.AnalysisResult, error) {
	result := &config.AnalysisResult{
		Recommendations: []config.Recommendation{},
	}

	// Analyze database configuration
	dbAnalysis := a.analyzeDatabase()
	result.DatabaseAnalysis = dbAnalysis

	// Analyze Redis configuration
	redisAnalysis := a.analyzeRedis()
	result.RedisAnalysis = redisAnalysis

	// Analyze web tier configuration
	webAnalysis := a.analyzeWebTier()
	result.WebTierAnalysis = webAnalysis

	// Generate recommendations based on analysis
	result.Recommendations = a.generateRecommendations(dbAnalysis, redisAnalysis, webAnalysis)

	return result, nil
}

// getEnvVarInt retrieves an environment variable as an integer
func (a *Analyzer) getEnvVarInt(name string, defaultValue int) int {
	if val, ok := a.envVars[name]; ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// getDynosByType returns all dynos of a specific type
func (a *Analyzer) getDynosByType(dynoType string) *config.DynoFormation {
	for _, dyno := range a.dynos {
		if dyno.Type == dynoType {
			return &dyno
		}
	}
	return nil
}

// extractPostgresPlan extracts the Postgres plan from DATABASE_URL or addon name
func (a *Analyzer) extractPostgresPlan() string {
	// Look through addons for Postgres
	for _, addon := range a.addons {
		if strings.Contains(strings.ToLower(addon.Plan), "postgres") ||
		   strings.Contains(strings.ToLower(addon.Name), "postgres") {
			// Extract plan from something like "heroku-postgresql:standard-0"
			parts := strings.Split(addon.Plan, ":")
			if len(parts) == 2 {
				return parts[1]
			}
			return addon.Plan
		}
	}

	// If not found in addons, return "unknown"
	return "unknown"
}

// extractRedisPlan extracts the Redis plan from REDIS_URL or addon name
func (a *Analyzer) extractRedisPlan() string {
	// Look through addons for Redis
	for _, addon := range a.addons {
		if strings.Contains(strings.ToLower(addon.Plan), "redis") ||
		   strings.Contains(strings.ToLower(addon.Name), "redis") {
			// Extract plan from something like "heroku-redis:premium-0"
			parts := strings.Split(addon.Plan, ":")
			if len(parts) == 2 {
				return parts[1]
			}
			return addon.Plan
		}
	}

	return "unknown"
}

// hasEnvVar checks if an environment variable exists
func (a *Analyzer) hasEnvVar(name string) bool {
	_, ok := a.envVars[name]
	return ok
}
