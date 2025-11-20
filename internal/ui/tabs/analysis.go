package tabs

import (
	"fmt"
	"strings"

	"github.com/leaharmstrong/heroku-calc/internal/config"
)

// RenderAnalysis renders the analysis tab
func RenderAnalysis(analysis *config.AnalysisResult) string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("ANALYSIS RESULTS\n\n")

	if analysis == nil {
		content.WriteString("  No analysis available\n")
		return content.String()
	}

	// Database Analysis
	if analysis.DatabaseAnalysis != nil {
		content.WriteString(renderDatabaseAnalysis(analysis.DatabaseAnalysis))
		content.WriteString("\n")
	}

	// Redis Analysis
	if analysis.RedisAnalysis != nil {
		content.WriteString(renderRedisAnalysis(analysis.RedisAnalysis))
		content.WriteString("\n")
	}

	// Web Tier Analysis
	if analysis.WebTierAnalysis != nil {
		content.WriteString(renderWebTierAnalysis(analysis.WebTierAnalysis))
	}

	return content.String()
}

func renderDatabaseAnalysis(analysis *config.DatabaseAnalysis) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("DATABASE CONNECTIONS - %s\n", formatStatus(analysis.Status)))
	content.WriteString(fmt.Sprintf("  Plan: %s\n", analysis.PostgresPlan))
	content.WriteString(fmt.Sprintf("  Max connections: %d\n", analysis.MaxConnections))
	content.WriteString(fmt.Sprintf("  Current usage: %d dynos × %d workers × %d threads = %d connections\n",
		analysis.WebDynos, analysis.WorkersPerDyno, analysis.ThreadsPerWorker,
		analysis.WebDynos*analysis.WorkersPerDyno*analysis.ThreadsPerWorker))

	if analysis.SidekiqDynos > 0 {
		content.WriteString(fmt.Sprintf("  Sidekiq: %d dynos × %d threads = %d connections\n",
			analysis.SidekiqDynos, analysis.SidekiqThreads, analysis.SidekiqDynos*analysis.SidekiqThreads))
	}

	content.WriteString(fmt.Sprintf("  Total required: %d / %d available (%.1f%% buffer)\n",
		analysis.TotalRequired, analysis.MaxConnections, analysis.BufferPercent))

	if len(analysis.Issues) > 0 {
		content.WriteString("\n  Issues:\n")
		for _, issue := range analysis.Issues {
			content.WriteString(fmt.Sprintf("  • %s\n", issue))
		}
	}

	return content.String()
}

func renderRedisAnalysis(analysis *config.RedisAnalysis) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("REDIS CONFIGURATION - %s\n", formatStatus(analysis.Status)))

	if analysis.RedisURL == "unknown" {
		content.WriteString("  Not configured\n")
		return content.String()
	}

	content.WriteString(fmt.Sprintf("  Plan: %s\n", analysis.RedisPlan))
	content.WriteString(fmt.Sprintf("  Max connections: %d\n", analysis.MaxConnections))
	content.WriteString(fmt.Sprintf("  Sidekiq concurrency: %d\n", analysis.SidekiqConcurrency))
	content.WriteString(fmt.Sprintf("  Redis pool size: %d\n", analysis.RedisPoolSize))
	content.WriteString(fmt.Sprintf("  Estimated usage: %d\n", analysis.EstimatedUsage))

	if analysis.MaxConnections > 0 {
		utilizationPercent := float64(analysis.EstimatedUsage) / float64(analysis.MaxConnections) * 100
		content.WriteString(fmt.Sprintf("  Utilization: %.1f%%\n", utilizationPercent))
	}

	if len(analysis.Issues) > 0 {
		content.WriteString("\n  Issues:\n")
		for _, issue := range analysis.Issues {
			content.WriteString(fmt.Sprintf("  • %s\n", issue))
		}
	}

	return content.String()
}

func renderWebTierAnalysis(analysis *config.WebTierAnalysis) string {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("WEB TIER - %s\n", formatStatus(analysis.Status)))
	content.WriteString(fmt.Sprintf("  Dyno type: %s (%d MB RAM)\n", analysis.DynoType, analysis.DynoMemoryMB))
	content.WriteString(fmt.Sprintf("  WEB_CONCURRENCY: %d workers\n", analysis.WebConcurrency))
	content.WriteString(fmt.Sprintf("  RAILS_MAX_THREADS: %d threads per worker\n", analysis.RailsMaxThreads))
	content.WriteString(fmt.Sprintf("  Total threads: %d\n", analysis.TotalThreads))

	if analysis.MemoryPerThread > 0 {
		content.WriteString(fmt.Sprintf("  Memory per thread: %d MB\n", analysis.MemoryPerThread))
	}

	if len(analysis.Issues) > 0 {
		content.WriteString("\n  Issues:\n")
		for _, issue := range analysis.Issues {
			content.WriteString(fmt.Sprintf("  • %s\n", issue))
		}
	}

	return content.String()
}
