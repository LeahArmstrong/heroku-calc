package analysis

import (
	"fmt"

	"github.com/leaharmstrong/heroku-calc/internal/config"
)

// analyzeDatabase analyzes database connection configuration
func (a *Analyzer) analyzeDatabase() *config.DatabaseAnalysis {
	analysis := &config.DatabaseAnalysis{
		DatabaseURL:      "unknown",
		PostgresPlan:     "unknown",
		MaxConnections:   0,
		CurrentUsage:     0,
		WebDynos:         0,
		WorkersPerDyno:   0,
		ThreadsPerWorker: 0,
		SidekiqDynos:     0,
		SidekiqThreads:   0,
		TotalRequired:    0,
		BufferPercent:    0,
		Status:           config.StatusUnknown,
		Issues:           []string{},
	}

	// Check if DATABASE_URL exists
	if !a.hasEnvVar("DATABASE_URL") {
		analysis.Issues = append(analysis.Issues, "DATABASE_URL not found")
		analysis.Status = config.StatusWarning
		return analysis
	}

	analysis.DatabaseURL = "present"

	// Get Postgres plan
	postgresPlan := a.extractPostgresPlan()
	analysis.PostgresPlan = postgresPlan

	// Get max connections from pricing data
	if postgresPlan != "unknown" {
		if pgPrice, err := a.pricingData.GetPostgresPrice(postgresPlan); err == nil {
			analysis.MaxConnections = pgPrice.MaxConnections
		} else {
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Unknown Postgres plan: %s", postgresPlan))
		}
	}

	// Calculate connection requirements

	// Web dynos
	webDynos := a.getDynosByType("web")
	if webDynos != nil {
		analysis.WebDynos = webDynos.Quantity

		// WEB_CONCURRENCY (Puma workers per dyno)
		analysis.WorkersPerDyno = a.getEnvVarInt("WEB_CONCURRENCY", 2)

		// RAILS_MAX_THREADS (threads per worker)
		analysis.ThreadsPerWorker = a.getEnvVarInt("RAILS_MAX_THREADS", 5)

		// Total web connections = dynos × workers × threads
		webConnections := analysis.WebDynos * analysis.WorkersPerDyno * analysis.ThreadsPerWorker
		analysis.CurrentUsage += webConnections
	}

	// Worker/Sidekiq dynos
	workerDynos := a.getDynosByType("worker")
	if workerDynos != nil {
		analysis.SidekiqDynos = workerDynos.Quantity

		// Sidekiq concurrency
		concurrency := a.getEnvVarInt("SIDEKIQ_CONCURRENCY", 10)
		analysis.SidekiqThreads = concurrency

		// Total worker connections = dynos × concurrency
		workerConnections := analysis.SidekiqDynos * concurrency
		analysis.CurrentUsage += workerConnections
	}

	// Check DB_POOL setting
	dbPool := a.getEnvVarInt("DB_POOL", 0)
	if dbPool > 0 {
		// If DB_POOL is set, it overrides RAILS_MAX_THREADS for pool size
		analysis.Issues = append(analysis.Issues, fmt.Sprintf("DB_POOL is set to %d (overrides RAILS_MAX_THREADS)", dbPool))
	}

	analysis.TotalRequired = analysis.CurrentUsage

	// Calculate buffer percentage
	if analysis.MaxConnections > 0 {
		analysis.BufferPercent = float64(analysis.MaxConnections-analysis.TotalRequired) / float64(analysis.MaxConnections) * 100

		// Determine status
		if analysis.TotalRequired >= analysis.MaxConnections {
			analysis.Status = config.StatusCritical
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Connection exhaustion: %d required >= %d max", analysis.TotalRequired, analysis.MaxConnections))
		} else if analysis.BufferPercent < 20 {
			analysis.Status = config.StatusCritical
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Very low buffer: only %.1f%% available", analysis.BufferPercent))
		} else if analysis.BufferPercent < 50 {
			analysis.Status = config.StatusWarning
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Low buffer: %.1f%% available (recommend 50%+ for bursts)", analysis.BufferPercent))
		} else {
			analysis.Status = config.StatusOptimal
		}
	}

	return analysis
}
