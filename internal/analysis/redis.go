package analysis

import (
	"fmt"

	"github.com/egg/heroku-calc/internal/config"
)

// analyzeRedis analyzes Redis/cache connection configuration
func (a *Analyzer) analyzeRedis() *config.RedisAnalysis {
	analysis := &config.RedisAnalysis{
		RedisURL:           "unknown",
		RedisPlan:          "unknown",
		MaxConnections:     0,
		SidekiqConcurrency: 0,
		RedisPoolSize:      0,
		EstimatedUsage:     0,
		Status:             config.StatusUnknown,
		Issues:             []string{},
	}

	// Check if REDIS_URL exists
	if !a.hasEnvVar("REDIS_URL") {
		// Redis is optional for Rails apps
		analysis.Status = config.StatusOptimal
		analysis.Issues = append(analysis.Issues, "REDIS_URL not configured (optional)")
		return analysis
	}

	analysis.RedisURL = "present"

	// Get Redis plan
	redisPlan := a.extractRedisPlan()
	analysis.RedisPlan = redisPlan

	// Get max connections from pricing data
	if redisPlan != "unknown" {
		if redisPrice, err := a.pricingData.GetRedisPrice(redisPlan); err == nil {
			analysis.MaxConnections = redisPrice.MaxConnections
		} else {
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Unknown Redis plan: %s", redisPlan))
		}
	}

	// Calculate connection requirements

	// Sidekiq connections
	workerDynos := a.getDynosByType("worker")
	if workerDynos != nil {
		concurrency := a.getEnvVarInt("SIDEKIQ_CONCURRENCY", 10)
		analysis.SidekiqConcurrency = workerDynos.Quantity * concurrency
		analysis.EstimatedUsage += analysis.SidekiqConcurrency
	}

	// Web dynos using Redis (for cache, sessions, etc.)
	webDynos := a.getDynosByType("web")
	if webDynos != nil {
		// Check if REDIS_POOL_SIZE is set
		redisPoolSize := a.getEnvVarInt("REDIS_POOL_SIZE", 0)

		if redisPoolSize > 0 {
			analysis.RedisPoolSize = redisPoolSize
			analysis.EstimatedUsage += webDynos.Quantity * redisPoolSize
		} else {
			// Default Redis pool size is typically 5 per web process
			workersPerDyno := a.getEnvVarInt("WEB_CONCURRENCY", 2)
			defaultPoolPerProcess := 5
			analysis.RedisPoolSize = defaultPoolPerProcess
			analysis.EstimatedUsage += webDynos.Quantity * workersPerDyno * defaultPoolPerProcess

			analysis.Issues = append(analysis.Issues, "REDIS_POOL_SIZE not set (using default estimate of 5 per Puma worker)")
		}
	}

	// Determine status
	if analysis.MaxConnections > 0 {
		utilizationPercent := float64(analysis.EstimatedUsage) / float64(analysis.MaxConnections) * 100

		if analysis.EstimatedUsage >= analysis.MaxConnections {
			analysis.Status = config.StatusCritical
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Redis connection exhaustion: %d estimated >= %d max", analysis.EstimatedUsage, analysis.MaxConnections))
		} else if utilizationPercent > 80 {
			analysis.Status = config.StatusWarning
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("High Redis utilization: %.1f%% (recommend <80%%)", utilizationPercent))
		} else {
			analysis.Status = config.StatusOptimal
		}
	}

	return analysis
}
