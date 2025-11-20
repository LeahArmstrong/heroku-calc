package analysis

import (
	"fmt"

	"github.com/leaharmstrong/heroku-calc/internal/config"
)

// analyzeWebTier analyzes web tier concurrency configuration
func (a *Analyzer) analyzeWebTier() *config.WebTierAnalysis {
	analysis := &config.WebTierAnalysis{
		DynoType:        "unknown",
		DynoMemoryMB:    0,
		WebConcurrency:  0,
		RailsMaxThreads: 0,
		TotalThreads:    0,
		MemoryPerThread: 0,
		Status:          config.StatusUnknown,
		Issues:          []string{},
	}

	// Get web dynos
	webDynos := a.getDynosByType("web")
	if webDynos == nil {
		analysis.Issues = append(analysis.Issues, "No web dynos found")
		analysis.Status = config.StatusWarning
		return analysis
	}

	analysis.DynoType = webDynos.Size

	// Get dyno memory from pricing data
	if dynoPrice, err := a.pricingData.GetDynoPrice(webDynos.Size); err == nil {
		analysis.DynoMemoryMB = dynoPrice.MemoryMB
	} else {
		analysis.Issues = append(analysis.Issues, fmt.Sprintf("Unknown dyno type: %s", webDynos.Size))
	}

	// Get concurrency settings
	analysis.WebConcurrency = a.getEnvVarInt("WEB_CONCURRENCY", 2)
	analysis.RailsMaxThreads = a.getEnvVarInt("RAILS_MAX_THREADS", 5)

	// Calculate total threads per dyno
	analysis.TotalThreads = analysis.WebConcurrency * analysis.RailsMaxThreads

	// Calculate memory per thread
	if analysis.DynoMemoryMB > 0 && analysis.TotalThreads > 0 {
		analysis.MemoryPerThread = analysis.DynoMemoryMB / analysis.TotalThreads
	}

	// Analyze configuration
	if analysis.DynoMemoryMB > 0 {
		// Recommended minimum memory per thread for Rails apps
		const minMemoryPerThread = 50 // MB
		const recommendedMemoryPerThread = 80 // MB

		if analysis.MemoryPerThread < minMemoryPerThread {
			analysis.Status = config.StatusCritical
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Too many threads for dyno size: %d MB per thread (recommend minimum %d MB)", analysis.MemoryPerThread, minMemoryPerThread))
		} else if analysis.MemoryPerThread < recommendedMemoryPerThread {
			analysis.Status = config.StatusWarning
			analysis.Issues = append(analysis.Issues, fmt.Sprintf("Tight memory allocation: %d MB per thread (recommend %d+ MB)", analysis.MemoryPerThread, recommendedMemoryPerThread))
		} else {
			analysis.Status = config.StatusOptimal
		}

		// Check if WEB_CONCURRENCY or RAILS_MAX_THREADS are not set
		if !a.hasEnvVar("WEB_CONCURRENCY") {
			analysis.Issues = append(analysis.Issues, "WEB_CONCURRENCY not explicitly set (using default 2)")
		}
		if !a.hasEnvVar("RAILS_MAX_THREADS") {
			analysis.Issues = append(analysis.Issues, "RAILS_MAX_THREADS not explicitly set (using default 5)")
		}

		// Recommendations for specific dyno types
		switch webDynos.Size {
		case "eco", "basic":
			if analysis.TotalThreads > 5 {
				analysis.Status = config.StatusWarning
				analysis.Issues = append(analysis.Issues, fmt.Sprintf("Eco/Basic dynos recommended for max 5 total threads, currently: %d", analysis.TotalThreads))
			}
		case "standard-1x":
			if analysis.TotalThreads > 5 {
				analysis.Status = config.StatusWarning
				analysis.Issues = append(analysis.Issues, fmt.Sprintf("Standard-1X recommended for max 5 total threads, currently: %d", analysis.TotalThreads))
			}
		case "standard-2x":
			if analysis.TotalThreads > 10 {
				analysis.Status = config.StatusWarning
				analysis.Issues = append(analysis.Issues, fmt.Sprintf("Standard-2X recommended for max 10 total threads, currently: %d", analysis.TotalThreads))
			}
		}
	}

	return analysis
}
