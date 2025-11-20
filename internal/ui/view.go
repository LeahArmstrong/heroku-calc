package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the application view
func (m Model) View() string {
	if m.state == StateError {
		return m.renderError()
	}

	if m.state == StateLoading || m.state == StateAnalyzing {
		return m.renderLoading()
	}

	var content strings.Builder

	// Header
	content.WriteString(m.renderHeader())
	content.WriteString("\n")

	// Tabs
	content.WriteString(m.renderTabs())
	content.WriteString("\n")

	// Tab content
	content.WriteString(m.renderTabContent())
	content.WriteString("\n")

	// Status bar
	content.WriteString(m.renderStatusBar())

	return content.String()
}

// renderHeader renders the application header
func (m Model) renderHeader() string {
	var header strings.Builder

	// App title
	title := appTitleStyle.Render("Heroku Config Analyzer")

	// App info
	appInfo := ""
	if m.appInfo != nil {
		appInfo = fmt.Sprintf("App: %s | Region: %s", m.appInfo.Name, m.appInfo.Region)
	}

	// Dyno summary
	dynoSummary := ""
	if len(m.dynos) > 0 {
		webCount := 0
		workerCount := 0
		for _, dyno := range m.dynos {
			if dyno.Type == "web" {
				webCount = dyno.Quantity
			} else if dyno.Type == "worker" {
				workerCount = dyno.Quantity
			}
		}
		dynoSummary = fmt.Sprintf(" | Dynos: %d web, %d worker", webCount, workerCount)
	}

	header.WriteString(title)
	header.WriteString("  ")
	header.WriteString(baseStyle.Render(appInfo + dynoSummary))

	return header.String()
}

// renderTabs renders the tab navigation
func (m Model) renderTabs() string {
	var tabs []string

	for i := Tab(0); i < 6; i++ {
		name := m.GetTabName(i)
		if i == m.currentTab {
			tabs = append(tabs, activeTabStyle.Render(name))
		} else {
			tabs = append(tabs, inactiveTabStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderTabContent renders the current tab's content
func (m Model) renderTabContent() string {
	switch m.currentTab {
	case TabOverview:
		return m.renderOverviewTab()
	case TabEnvVars:
		return m.renderEnvVarsTab()
	case TabDynos:
		return m.renderDynosTab()
	case TabAddons:
		return m.renderAddonsTab()
	case TabAnalysis:
		return m.renderAnalysisTab()
	case TabActions:
		return m.renderActionsTab()
	default:
		return "Unknown tab"
	}
}

// renderStatusBar renders the bottom status bar
func (m Model) renderStatusBar() string {
	helpText := "Tab ←→  ↑↓ Navigate  Enter Select  e Export  q Quit"
	if m.currentTab == TabActions && m.mode != ModeReadOnly {
		helpText = "Tab ←→  ↑↓ Navigate  Enter Select  a Apply  e Export  q Quit"
	}
	leftSection := helpStyle.Render(helpText)

	mode := fmt.Sprintf("Mode: %s", m.GetModeString())
	if m.statusMessage != "" {
		mode = m.statusMessage
	}
	rightSection := statusBarStyle.Render(mode)

	// Calculate padding to fill the width
	padding := m.width - lipgloss.Width(leftSection) - lipgloss.Width(rightSection)
	if padding < 0 {
		padding = 0
	}

	return leftSection + strings.Repeat(" ", padding) + rightSection
}

// renderLoading renders the loading screen
func (m Model) renderLoading() string {
	var content strings.Builder

	content.WriteString("\n\n")
	content.WriteString(titleStyle.Render("Heroku Config Analyzer"))
	content.WriteString("\n\n")

	message := "Loading data from Heroku..."
	if m.state == StateAnalyzing {
		message = "Analyzing configuration..."
	}

	content.WriteString(fmt.Sprintf("  %s %s\n\n", m.spinner.View(), message))

	return content.String()
}

// renderError renders the error screen
func (m Model) renderError() string {
	var content strings.Builder

	content.WriteString("\n\n")
	content.WriteString(statusCriticalStyle.Render("Error"))
	content.WriteString("\n\n")

	if m.err != nil {
		content.WriteString(fmt.Sprintf("  %s\n\n", m.err.Error()))
	} else {
		content.WriteString("  An unknown error occurred\n\n")
	}

	content.WriteString(helpStyle.Render("  Press q to quit"))
	content.WriteString("\n\n")

	return content.String()
}
