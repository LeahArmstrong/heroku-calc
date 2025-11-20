package tabs

import (
	"fmt"
	"strings"

	"github.com/leaharmstrong/heroku-calc/internal/config"
	"github.com/leaharmstrong/heroku-calc/internal/heroku"
)

// RenderOverview renders the overview tab
func RenderOverview(appInfo *heroku.AppInfo, dynos []config.DynoFormation, addons []config.Addon, analysis *config.AnalysisResult) string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("APPLICATION OVERVIEW\n\n")

	if appInfo != nil {
		content.WriteString(fmt.Sprintf("  Name:   %s\n", appInfo.Name))
		content.WriteString(fmt.Sprintf("  Region: %s\n", appInfo.Region))
		content.WriteString(fmt.Sprintf("  Stack:  %s\n\n", appInfo.Stack))
	}

	// Dyno summary
	content.WriteString("DYNOS\n")
	if len(dynos) > 0 {
		for _, dyno := range dynos {
			content.WriteString(fmt.Sprintf("  %s: %d × %s\n", dyno.Type, dyno.Quantity, dyno.Size))
		}
	} else {
		content.WriteString("  No dynos found\n")
	}
	content.WriteString("\n")

	// Addon summary
	content.WriteString("ADDONS\n")
	if len(addons) > 0 {
		for _, addon := range addons {
			content.WriteString(fmt.Sprintf("  %s (%s)\n", addon.Name, addon.Plan))
		}
	} else {
		content.WriteString("  No addons configured\n")
	}
	content.WriteString("\n")

	// Quick status summary
	if analysis != nil {
		content.WriteString("HEALTH STATUS\n")

		dbStatus := "⚪ Unknown"
		if analysis.DatabaseAnalysis != nil {
			dbStatus = formatStatus(analysis.DatabaseAnalysis.Status)
		}
		content.WriteString(fmt.Sprintf("  Database: %s\n", dbStatus))

		redisStatus := "⚪ Unknown"
		if analysis.RedisAnalysis != nil {
			redisStatus = formatStatus(analysis.RedisAnalysis.Status)
		}
		content.WriteString(fmt.Sprintf("  Redis:    %s\n", redisStatus))

		webStatus := "⚪ Unknown"
		if analysis.WebTierAnalysis != nil {
			webStatus = formatStatus(analysis.WebTierAnalysis.Status)
		}
		content.WriteString(fmt.Sprintf("  Web Tier: %s\n", webStatus))
	}

	return content.String()
}

func formatStatus(status config.AnalysisStatus) string {
	switch status {
	case config.StatusCritical:
		return "🔴 Critical"
	case config.StatusWarning:
		return "🟡 Warning"
	case config.StatusOptimal:
		return "🟢 Optimal"
	default:
		return "⚪ Unknown"
	}
}
