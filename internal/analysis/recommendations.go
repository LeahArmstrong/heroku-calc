package analysis

import (
	"fmt"

	"github.com/egg/heroku-calc/internal/config"
)

// generateRecommendations creates actionable recommendations based on analysis
func (a *Analyzer) generateRecommendations(
	dbAnalysis *config.DatabaseAnalysis,
	redisAnalysis *config.RedisAnalysis,
	webAnalysis *config.WebTierAnalysis,
) []config.Recommendation {
	recommendations := []config.Recommendation{}

	// Database recommendations
	if dbAnalysis != nil {
		recommendations = append(recommendations, a.generateDatabaseRecommendations(dbAnalysis)...)
	}

	// Redis recommendations
	if redisAnalysis != nil {
		recommendations = append(recommendations, a.generateRedisRecommendations(redisAnalysis)...)
	}

	// Web tier recommendations
	if webAnalysis != nil {
		recommendations = append(recommendations, a.generateWebTierRecommendations(webAnalysis)...)
	}

	return recommendations
}

func (a *Analyzer) generateDatabaseRecommendations(analysis *config.DatabaseAnalysis) []config.Recommendation {
	recommendations := []config.Recommendation{}

	if analysis.Status == config.StatusCritical || analysis.Status == config.StatusWarning {
		// Recommend upgrading Postgres plan
		if analysis.BufferPercent < 50 && analysis.PostgresPlan != "unknown" {
			// Find next tier up
			suggestedPlan := a.suggestNextPostgresPlan(analysis.PostgresPlan, analysis.TotalRequired)

			if suggestedPlan != "" {
				var severity config.RecommendationSeverity
				if analysis.BufferPercent < 20 {
					severity = config.SeverityCritical
				} else {
					severity = config.SeverityHigh
				}

				recommendations = append(recommendations, config.Recommendation{
					Category:    "database",
					Severity:    severity,
					Title:       "Upgrade Postgres Plan",
					Description: fmt.Sprintf("Current plan has only %.1f%% buffer. Recommend upgrading to ensure connection capacity.", analysis.BufferPercent),
					Current:     analysis.PostgresPlan,
					Suggested:   suggestedPlan,
					Impact:      a.calculatePostgresCostImpact(analysis.PostgresPlan, suggestedPlan),
					AutoApply:   false, // Plan upgrades require manual intervention
				})
			}
		}

		// Recommend reducing connections
		if analysis.BufferPercent < 20 {
			recommendations = append(recommendations, config.Recommendation{
				Category:    "database",
				Severity:    config.SeverityHigh,
				Title:       "Reduce Database Connections",
				Description: "Consider reducing WEB_CONCURRENCY, RAILS_MAX_THREADS, or SIDEKIQ_CONCURRENCY to lower connection usage",
				Current:     fmt.Sprintf("%d total connections required", analysis.TotalRequired),
				Suggested:   fmt.Sprintf("Target <%d connections (50%% of %d max)", analysis.MaxConnections/2, analysis.MaxConnections),
				Impact:      "Lower resource utilization, better burst capacity",
				AutoApply:   false,
			})
		}
	}

	return recommendations
}

func (a *Analyzer) generateRedisRecommendations(analysis *config.RedisAnalysis) []config.Recommendation {
	recommendations := []config.Recommendation{}

	// Recommend setting REDIS_POOL_SIZE if not set
	if analysis.RedisURL != "unknown" && !a.hasEnvVar("REDIS_POOL_SIZE") {
		webDynos := a.getDynosByType("web")
		if webDynos != nil {
			// Recommend explicit pool size based on concurrency
			threadsPerWorker := a.getEnvVarInt("RAILS_MAX_THREADS", 5)
			suggestedPoolSize := threadsPerWorker + 2 // Add small buffer

			recommendations = append(recommendations, config.Recommendation{
				Category:    "redis",
				Severity:    config.SeverityMedium,
				Title:       "Set REDIS_POOL_SIZE",
				Description: "Explicitly configure Redis connection pool size to match thread count",
				Current:     "not set (using default)",
				Suggested:   fmt.Sprintf("%d", suggestedPoolSize),
				EnvVarName:  "REDIS_POOL_SIZE",
				Impact:      "Prevents connection exhaustion and improves performance",
				AutoApply:   true,
			})
		}
	}

	// Recommend upgrading Redis plan if near capacity
	if analysis.MaxConnections > 0 && analysis.EstimatedUsage > 0 {
		utilizationPercent := float64(analysis.EstimatedUsage) / float64(analysis.MaxConnections) * 100

		if utilizationPercent > 80 {
			suggestedPlan := a.suggestNextRedisPlan(analysis.RedisPlan, analysis.EstimatedUsage)

			if suggestedPlan != "" {
				recommendations = append(recommendations, config.Recommendation{
					Category:    "redis",
					Severity:    config.SeverityHigh,
					Title:       "Upgrade Redis Plan",
					Description: fmt.Sprintf("Redis utilization at %.1f%% - upgrade for more connection capacity", utilizationPercent),
					Current:     analysis.RedisPlan,
					Suggested:   suggestedPlan,
					Impact:      a.calculateRedisCostImpact(analysis.RedisPlan, suggestedPlan),
					AutoApply:   false,
				})
			}
		}
	}

	return recommendations
}

func (a *Analyzer) generateWebTierRecommendations(analysis *config.WebTierAnalysis) []config.Recommendation {
	recommendations := []config.Recommendation{}

	// Recommend setting WEB_CONCURRENCY if not set
	if !a.hasEnvVar("WEB_CONCURRENCY") {
		recommendations = append(recommendations, config.Recommendation{
			Category:    "web",
			Severity:    config.SeverityMedium,
			Title:       "Set WEB_CONCURRENCY",
			Description: "Explicitly configure Puma worker count for better performance tuning",
			Current:     "not set (using default 2)",
			Suggested:   a.suggestWebConcurrency(analysis.DynoMemoryMB),
			EnvVarName:  "WEB_CONCURRENCY",
			Impact:      "Optimizes worker count for dyno size",
			AutoApply:   true,
		})
	}

	// Recommend setting RAILS_MAX_THREADS if not set
	if !a.hasEnvVar("RAILS_MAX_THREADS") {
		recommendations = append(recommendations, config.Recommendation{
			Category:    "web",
			Severity:    config.SeverityMedium,
			Title:       "Set RAILS_MAX_THREADS",
			Description: "Explicitly configure Puma thread count per worker",
			Current:     "not set (using default 5)",
			Suggested:   "5",
			EnvVarName:  "RAILS_MAX_THREADS",
			Impact:      "Prevents unexpected behavior from default changes",
			AutoApply:   true,
		})
	}

	// Recommend adjusting concurrency if memory per thread is too low
	if analysis.MemoryPerThread > 0 && analysis.MemoryPerThread < 50 {
		suggestedThreads := analysis.DynoMemoryMB / 80 // 80MB per thread target
		if suggestedThreads < 1 {
			suggestedThreads = 1
		}

		recommendations = append(recommendations, config.Recommendation{
			Category:    "web",
			Severity:    config.SeverityCritical,
			Title:       "Reduce Thread Count",
			Description: fmt.Sprintf("Only %d MB per thread - risk of memory exhaustion and dyno crashes", analysis.MemoryPerThread),
			Current:     fmt.Sprintf("WEB_CONCURRENCY=%d, RAILS_MAX_THREADS=%d (%d total threads)", analysis.WebConcurrency, analysis.RailsMaxThreads, analysis.TotalThreads),
			Suggested:   fmt.Sprintf("Reduce to ~%d total threads or upgrade dyno size", suggestedThreads),
			Impact:      "Prevents R14 memory errors and dyno restarts",
			AutoApply:   false, // Requires manual decision
		})
	}

	return recommendations
}

// Helper functions to suggest next tier plans

func (a *Analyzer) suggestNextPostgresPlan(currentPlan string, requiredConnections int) string {
	// List of plans in order
	plans := []string{"mini", "basic", "standard-0", "standard-2", "standard-3", "standard-4", "standard-5", "standard-6"}

	targetConnections := int(float64(requiredConnections) * 2.0) // 100% buffer

	for _, plan := range plans {
		if price, err := a.pricingData.GetPostgresPrice(plan); err == nil {
			if price.MaxConnections >= targetConnections {
				return plan
			}
		}
	}

	return ""
}

func (a *Analyzer) suggestNextRedisPlan(currentPlan string, requiredConnections int) string {
	plans := []string{"mini", "premium-0", "premium-1", "premium-2", "premium-3", "premium-4", "premium-5"}

	targetConnections := int(float64(requiredConnections) * 1.5) // 50% buffer

	for _, plan := range plans {
		if price, err := a.pricingData.GetRedisPrice(plan); err == nil {
			if price.MaxConnections >= targetConnections {
				return plan
			}
		}
	}

	return ""
}

func (a *Analyzer) suggestWebConcurrency(dynoMemoryMB int) string {
	// Conservative recommendations for WEB_CONCURRENCY based on dyno memory
	// Assuming ~200-300MB per worker process
	if dynoMemoryMB <= 512 {
		return "1"
	} else if dynoMemoryMB <= 1024 {
		return "2"
	} else if dynoMemoryMB <= 2560 {
		return "2-3"
	} else {
		return "3-4"
	}
}

func (a *Analyzer) calculatePostgresCostImpact(currentPlan, suggestedPlan string) string {
	currentPrice, err1 := a.pricingData.GetPostgresPrice(currentPlan)
	suggestedPrice, err2 := a.pricingData.GetPostgresPrice(suggestedPlan)

	if err1 != nil || err2 != nil {
		return "Cost impact unknown"
	}

	diff := suggestedPrice.PriceMonthly - currentPrice.PriceMonthly
	return fmt.Sprintf("+$%.2f/month ($%.2f → $%.2f)", diff, currentPrice.PriceMonthly, suggestedPrice.PriceMonthly)
}

func (a *Analyzer) calculateRedisCostImpact(currentPlan, suggestedPlan string) string {
	currentPrice, err1 := a.pricingData.GetRedisPrice(currentPlan)
	suggestedPrice, err2 := a.pricingData.GetRedisPrice(suggestedPlan)

	if err1 != nil || err2 != nil {
		return "Cost impact unknown"
	}

	diff := suggestedPrice.PriceMonthly - currentPrice.PriceMonthly
	return fmt.Sprintf("+$%.2f/month ($%.2f → $%.2f)", diff, currentPrice.PriceMonthly, suggestedPrice.PriceMonthly)
}
