package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/leaharmstrong/heroku-calc/internal/config"
	"github.com/leaharmstrong/heroku-calc/internal/heroku"
	"github.com/leaharmstrong/heroku-calc/internal/pricing"
)

// AppMode represents the current operation mode
type AppMode int

const (
	ModeReadOnly AppMode = iota
	ModeDryRun
	ModeInteractive
	ModeApply
)

// Tab represents the available tabs
type Tab int

const (
	TabOverview Tab = iota
	TabEnvVars
	TabDynos
	TabAddons
	TabAnalysis
	TabActions
)

// AppState represents the current state of the application
type AppState int

const (
	StateLoading AppState = iota
	StateReady
	StateAnalyzing
	StateError
	StateApplying
)

// Model represents the BubbleTea application model
type Model struct {
	// Configuration
	projectPath string
	appName     string
	mode        AppMode

	// Current state
	state      AppState
	currentTab Tab
	err        error

	// Data
	herokuClient *heroku.Client
	appInfo      *heroku.AppInfo
	envVars      []config.HerokuEnvVar
	dynos        []config.DynoFormation
	addons       []config.Addon
	cfg          *config.Config
	pricingData  *pricing.Data
	analysis     *config.AnalysisResult

	// UI state
	spinner         spinner.Model
	width           int
	height          int
	selectedEnvVars map[string]bool
	selectedActions map[int]bool
	cursorPos       int

	// Messages
	statusMessage string
}

// NewModel creates a new application model
func NewModel(projectPath, appName string, mode AppMode) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle

	return Model{
		projectPath:     projectPath,
		appName:         appName,
		mode:            mode,
		state:           StateLoading,
		currentTab:      TabOverview,
		spinner:         s,
		selectedEnvVars: make(map[string]bool),
		selectedActions: make(map[int]bool),
	}
}

// GetTabName returns the display name for a tab
func (m Model) GetTabName(tab Tab) string {
	switch tab {
	case TabOverview:
		return "Overview"
	case TabEnvVars:
		return "Env Vars"
	case TabDynos:
		return "Dynos"
	case TabAddons:
		return "Addons"
	case TabAnalysis:
		return "Analysis"
	case TabActions:
		return "Actions"
	default:
		return "Unknown"
	}
}

// GetModeString returns the display string for the current mode
func (m Model) GetModeString() string {
	switch m.mode {
	case ModeReadOnly:
		return "Read-Only"
	case ModeDryRun:
		return "Dry Run"
	case ModeInteractive:
		return "Interactive"
	case ModeApply:
		return "Apply Mode"
	default:
		return "Unknown"
	}
}
