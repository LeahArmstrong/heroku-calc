package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	primaryColor   = lipgloss.Color("63")  // Blue
	successColor   = lipgloss.Color("42")  // Green
	warningColor   = lipgloss.Color("220") // Yellow
	criticalColor  = lipgloss.Color("196") // Red
	mutedColor     = lipgloss.Color("241") // Gray
	highlightColor = lipgloss.Color("212") // Pink

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Title styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	appTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("white")).
			Background(primaryColor).
			Padding(0, 2)

	// Status styles
	statusCriticalStyle = lipgloss.NewStyle().
				Foreground(criticalColor).
				Bold(true)

	statusWarningStyle = lipgloss.NewStyle().
				Foreground(warningColor).
				Bold(true)

	statusOptimalStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	// Tab styles
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("white")).
			Background(primaryColor).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 2)

	// Border styles
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Status bar style
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("white")).
			Background(primaryColor).
			Padding(0, 1)

	// Table styles
	tableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor)

	tableRowStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Button styles
	activeButtonStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("white")).
				Background(primaryColor).
				Padding(0, 2).
				MarginRight(1)

	inactiveButtonStyle = lipgloss.NewStyle().
				Foreground(mutedColor).
				Padding(0, 2).
				MarginRight(1)

	// Checkbox styles
	selectedCheckboxStyle = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true)

	unselectedCheckboxStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	// Loading spinner style
	spinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)
)
