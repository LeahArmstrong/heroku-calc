package ui

import (
	"fmt"
	"strings"

	"github.com/egg/heroku-calc/internal/heroku"
	"github.com/egg/heroku-calc/internal/ui/tabs"
)

// renderOverviewTab renders the overview tab
func (m Model) renderOverviewTab() string {
	return tabs.RenderOverview(m.appInfo, m.dynos, m.addons, m.analysis)
}

// renderEnvVarsTab renders the environment variables tab
func (m Model) renderEnvVarsTab() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("ENVIRONMENT VARIABLES\n")
	content.WriteString("Select variables to include in .heroku-calc.yml\n\n")

	if len(m.envVars) == 0 {
		content.WriteString("  No environment variables found\n")
		return content.String()
	}

	for i, envVar := range m.envVars {
		checkbox := "[ ]"
		if m.selectedEnvVars[envVar.Name] {
			checkbox = "[✓]"
		}

		cursor := "  "
		if i == m.cursorPos {
			cursor = "> "
		}

		// Sanitize the value for display
		displayValue := heroku.SanitizeEnvVarValue(envVar.Name, envVar.Value)

		content.WriteString(fmt.Sprintf("%s%s %s\n", cursor, checkbox, envVar.Name))

		// Show a preview of the value if cursor is on this item
		if i == m.cursorPos {
			content.WriteString(fmt.Sprintf("      %s\n", displayValue))
		}
	}

	content.WriteString(fmt.Sprintf("\nSelected: %d / %d variables\n", len(m.selectedEnvVars), len(m.envVars)))

	return content.String()
}

// renderDynosTab renders the dynos tab
func (m Model) renderDynosTab() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("DYNO FORMATION\n\n")

	if len(m.dynos) == 0 {
		content.WriteString("  No dynos found\n")
		return content.String()
	}

	// Header
	content.WriteString("  Type       Quantity  Size            Monthly Cost\n")
	content.WriteString("  ─────────  ────────  ──────────────  ────────────\n")

	totalCost := 0.0

	for _, dyno := range m.dynos {
		// Look up pricing
		monthlyCost := 0.0
		if m.pricingData != nil {
			if price, err := m.pricingData.GetDynoPrice(dyno.Size); err == nil {
				monthlyCost = price.PriceMonthly * float64(dyno.Quantity)
				totalCost += monthlyCost
			}
		}

		content.WriteString(fmt.Sprintf("  %-9s  %-8d  %-14s  $%.2f\n",
			dyno.Type, dyno.Quantity, dyno.Size, monthlyCost))
	}

	content.WriteString("  ─────────────────────────────────────────────────\n")
	content.WriteString(fmt.Sprintf("  Total monthly cost: $%.2f\n", totalCost))

	return content.String()
}

// renderAddonsTab renders the addons tab
func (m Model) renderAddonsTab() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("ADDONS\n\n")

	if len(m.addons) == 0 {
		content.WriteString("  No addons configured\n")
		return content.String()
	}

	for _, addon := range m.addons {
		content.WriteString(fmt.Sprintf("  %s\n", addon.Name))
		content.WriteString(fmt.Sprintf("    Plan: %s\n", addon.Plan))
		if addon.Price != "unknown" {
			content.WriteString(fmt.Sprintf("    Price: %s\n", addon.Price))
		}
		content.WriteString(fmt.Sprintf("    Added: %s\n\n", addon.AddedAt.Format("2006-01-02")))
	}

	return content.String()
}

// renderAnalysisTab renders the analysis tab
func (m Model) renderAnalysisTab() string {
	return tabs.RenderAnalysis(m.analysis)
}

// renderActionsTab renders the actions tab
func (m Model) renderActionsTab() string {
	var content strings.Builder

	content.WriteString("\n")
	content.WriteString("RECOMMENDED ACTIONS\n")
	content.WriteString("Select actions to apply\n\n")

	if m.analysis == nil || len(m.analysis.Recommendations) == 0 {
		content.WriteString("  No recommendations\n")
		return content.String()
	}

	for i, rec := range m.analysis.Recommendations {
		checkbox := "[ ]"
		if m.selectedActions[i] {
			checkbox = "[✓]"
		}

		cursor := "  "
		if i == m.cursorPos {
			cursor = "> "
		}

		// Severity indicator
		severity := ""
		switch rec.Severity {
		case "critical":
			severity = "🔴"
		case "high":
			severity = "🟠"
		case "medium":
			severity = "🟡"
		default:
			severity = "🟢"
		}

		autoApply := ""
		if rec.AutoApply {
			autoApply = " [AUTO]"
		} else {
			autoApply = " [MANUAL]"
		}

		content.WriteString(fmt.Sprintf("%s%s %s %s%s\n", cursor, checkbox, severity, rec.Title, autoApply))

		// Show details if cursor is on this item
		if i == m.cursorPos {
			content.WriteString(fmt.Sprintf("      %s\n", rec.Description))
			content.WriteString(fmt.Sprintf("      Current: %s\n", rec.Current))
			content.WriteString(fmt.Sprintf("      Suggested: %s\n", rec.Suggested))
			if rec.Impact != "" {
				content.WriteString(fmt.Sprintf("      Impact: %s\n", rec.Impact))
			}
			if rec.EnvVarName != "" {
				content.WriteString(fmt.Sprintf("      Env Var: %s\n", rec.EnvVarName))
			}
		}
	}

	selectedCount := 0
	for _, selected := range m.selectedActions {
		if selected {
			selectedCount++
		}
	}

	content.WriteString(fmt.Sprintf("\nSelected: %d / %d actions\n", selectedCount, len(m.analysis.Recommendations)))

	if m.mode != ModeReadOnly && selectedCount > 0 {
		content.WriteString("\nPress 'a' to apply selected actions\n")
	}

	return content.String()
}
