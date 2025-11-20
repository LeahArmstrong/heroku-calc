package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leaharmstrong/heroku-calc/internal/config"
	"github.com/leaharmstrong/heroku-calc/internal/heroku"
	"github.com/leaharmstrong/heroku-calc/internal/report"
)

// applySelectedActions applies the selected recommendations
func (m Model) applySelectedActions() (tea.Model, tea.Cmd) {
	if m.analysis == nil {
		return m, nil
	}

	// Collect selected recommendations
	selectedRecs := []config.Recommendation{}
	for i, selected := range m.selectedActions {
		if selected && i < len(m.analysis.Recommendations) {
			selectedRecs = append(selectedRecs, m.analysis.Recommendations[i])
		}
	}

	if len(selectedRecs) == 0 {
		m.statusMessage = "No actions selected"
		return m, nil
	}

	// Filter based on mode
	applicableRecs := []config.Recommendation{}
	for _, rec := range selectedRecs {
		switch m.mode {
		case ModeDryRun:
			// In dry-run mode, just show what would be applied
			applicableRecs = append(applicableRecs, rec)
		case ModeInteractive:
			// In interactive mode, only apply auto-apply recommendations
			// (interactive confirmation would be implemented here)
			if rec.AutoApply {
				applicableRecs = append(applicableRecs, rec)
			}
		case ModeApply:
			// In apply mode, apply all selected recommendations
			if rec.AutoApply {
				applicableRecs = append(applicableRecs, rec)
			}
		}
	}

	if len(applicableRecs) == 0 {
		m.statusMessage = "No auto-applicable actions selected (manual changes required)"
		return m, nil
	}

	m.state = StateApplying
	m.statusMessage = fmt.Sprintf("Applying %d change(s)...", len(applicableRecs))

	return m, applyRecommendations(m.herokuClient, applicableRecs, m.mode)
}

// exportReport exports the analysis as a markdown report
func (m Model) exportReport() (tea.Model, tea.Cmd) {
	if m.analysis == nil {
		m.statusMessage = "No analysis to export"
		return m, nil
	}

	appName := m.appName
	if m.appInfo != nil {
		appName = m.appInfo.Name
	}

	// Generate markdown
	markdown := report.GenerateMarkdown(appName, m.analysis)

	// Save to file
	filename, err := report.SaveInProjectDir(markdown, m.projectPath, appName)
	if err != nil {
		m.statusMessage = fmt.Sprintf("Export failed: %v", err)
		return m, nil
	}

	m.statusMessage = fmt.Sprintf("Report exported to: %s", filename)
	return m, nil
}

// applyRecommendations applies the given recommendations
func applyRecommendations(clientInterface interface{}, recommendations []config.Recommendation, mode AppMode) tea.Cmd {
	return func() tea.Msg {
		if mode == ModeDryRun {
			// In dry-run mode, just return success without actually applying
			return applyCompleteMsg{
				success: true,
				err:     nil,
			}
		}

		// Cast the client
		client, ok := clientInterface.(*heroku.Client)
		if !ok || client == nil {
			return applyCompleteMsg{
				success: false,
				err:     fmt.Errorf("invalid heroku client"),
			}
		}

		// Apply each recommendation
		for _, rec := range recommendations {
			if !rec.AutoApply {
				continue // Skip manual recommendations
			}

			if rec.EnvVarName == "" {
				continue // Can't auto-apply without env var name
			}

			// Set the environment variable
			err := client.SetEnvVar(rec.EnvVarName, rec.Suggested)
			if err != nil {
				return applyCompleteMsg{
					success: false,
					err:     fmt.Errorf("failed to set %s: %w", rec.EnvVarName, err),
				}
			}
		}

		return applyCompleteMsg{
			success: true,
			err:     nil,
		}
	}
}
