package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/leaharmstrong/heroku-calc/internal/analysis"
	"github.com/leaharmstrong/heroku-calc/internal/config"
	"github.com/leaharmstrong/heroku-calc/internal/heroku"
	"github.com/leaharmstrong/heroku-calc/internal/pricing"
)

// Messages for async operations
type loadedDataMsg struct {
	client       *heroku.Client
	appInfo      *heroku.AppInfo
	envVars      []config.HerokuEnvVar
	dynos        []config.DynoFormation
	addons       []config.Addon
	pricingData  *pricing.Data
	cfg          *config.Config
	err          error
}

type analysisCompleteMsg struct {
	result *config.AnalysisResult
	err    error
}

type applyCompleteMsg struct {
	success bool
	err     error
}

type errMsg struct {
	err error
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadData(m.projectPath, m.appName),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		if m.state == StateLoading || m.state == StateAnalyzing {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case loadedDataMsg:
		if msg.err != nil {
			m.state = StateError
			m.err = msg.err
			return m, nil
		}

		m.herokuClient = msg.client
		m.appInfo = msg.appInfo
		m.envVars = msg.envVars
		m.dynos = msg.dynos
		m.addons = msg.addons
		m.pricingData = msg.pricingData
		m.cfg = msg.cfg

		// Initialize selected env vars from config
		if m.cfg != nil {
			for _, varName := range m.cfg.SafeEnvVars {
				m.selectedEnvVars[varName] = true
			}
		}

		// Move to analyzing state
		m.state = StateAnalyzing
		m.statusMessage = "Running analysis..."
		return m, runAnalysis(m.herokuClient, m.pricingData)

	case analysisCompleteMsg:
		if msg.err != nil {
			m.state = StateError
			m.err = msg.err
			return m, nil
		}

		m.analysis = msg.result
		m.state = StateReady
		m.statusMessage = "Analysis complete"
		return m, nil

	case applyCompleteMsg:
		if msg.err != nil {
			m.statusMessage = fmt.Sprintf("Error: %v", msg.err)
		} else {
			m.statusMessage = "Changes applied successfully"
		}
		m.state = StateReady
		return m, nil

	case errMsg:
		m.state = StateError
		m.err = msg.err
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "tab", "right":
		if m.state == StateReady {
			m.currentTab = (m.currentTab + 1) % 6
			m.cursorPos = 0
		}
		return m, nil

	case "shift+tab", "left":
		if m.state == StateReady {
			if m.currentTab == 0 {
				m.currentTab = 5
			} else {
				m.currentTab--
			}
			m.cursorPos = 0
		}
		return m, nil

	case "up", "k":
		if m.cursorPos > 0 {
			m.cursorPos--
		}
		return m, nil

	case "down", "j":
		maxPos := m.getMaxCursorPos()
		if m.cursorPos < maxPos {
			m.cursorPos++
		}
		return m, nil

	case "enter", " ":
		return m.handleSelection()

	case "a":
		// Apply selected actions (if in appropriate mode and on Actions tab)
		if m.currentTab == TabActions && m.mode != ModeReadOnly {
			return m.applySelectedActions()
		}
		return m, nil

	case "e":
		// Export markdown report
		if m.state == StateReady && m.analysis != nil {
			return m.exportReport()
		}
		return m, nil

	case "?":
		// Toggle help
		return m, nil
	}

	return m, nil
}

// handleSelection handles item selection based on current tab
func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	switch m.currentTab {
	case TabEnvVars:
		// Toggle env var selection
		if m.cursorPos < len(m.envVars) {
			varName := m.envVars[m.cursorPos].Name
			m.selectedEnvVars[varName] = !m.selectedEnvVars[varName]

			// Update config
			if m.selectedEnvVars[varName] {
				m.cfg.AddSafeEnvVar(varName)
			} else {
				m.cfg.RemoveSafeEnvVar(varName)
			}

			// Save config
			_ = config.Save(m.cfg, m.projectPath)
		}

	case TabActions:
		// Toggle action selection
		if m.analysis != nil && m.cursorPos < len(m.analysis.Recommendations) {
			m.selectedActions[m.cursorPos] = !m.selectedActions[m.cursorPos]
		}
	}

	return m, nil
}

// getMaxCursorPos returns the maximum cursor position for the current tab
func (m Model) getMaxCursorPos() int {
	switch m.currentTab {
	case TabEnvVars:
		return len(m.envVars) - 1
	case TabActions:
		if m.analysis != nil {
			return len(m.analysis.Recommendations) - 1
		}
	}
	return 0
}

// loadData loads all necessary data from Heroku
func loadData(projectPath, appName string) tea.Cmd {
	return func() tea.Msg {
		// Auto-detect app name from git if not provided
		if appName == "" {
			detectedName, _, err := heroku.DetectHerokuApp(projectPath)
			if err != nil {
				return loadedDataMsg{err: fmt.Errorf("failed to detect Heroku app: %w", err)}
			}
			appName = detectedName
		}

		// Create Heroku client
		client, err := heroku.NewClient(appName)
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to create Heroku client: %w", err)}
		}

		// Test connection
		if err := client.TestConnection(); err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to connect to Heroku: %w", err)}
		}

		// Load app info
		appInfo, err := client.GetAppInfo()
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to load app info: %w", err)}
		}

		// Load env vars
		envVars, err := client.GetEnvVars()
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to load env vars: %w", err)}
		}

		// Load dynos
		dynos, err := client.GetDynos()
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to load dynos: %w", err)}
		}

		// Load addons
		addons, err := client.GetAddons()
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to load addons: %w", err)}
		}

		// Load pricing data
		pricingData, err := pricing.Get()
		if err != nil {
			return loadedDataMsg{err: fmt.Errorf("failed to load pricing data: %w", err)}
		}

		// Load or create config
		var cfg *config.Config
		if config.Exists(projectPath) {
			cfg, err = config.Load(projectPath)
			if err != nil {
				return loadedDataMsg{err: fmt.Errorf("failed to load config: %w", err)}
			}
		} else {
			cfg = config.New(appName, projectPath)
			_ = config.Save(cfg, projectPath)
		}

		return loadedDataMsg{
			client:      client,
			appInfo:     appInfo,
			envVars:     envVars,
			dynos:       dynos,
			addons:      addons,
			pricingData: pricingData,
			cfg:         cfg,
		}
	}
}

// runAnalysis performs the configuration analysis
func runAnalysis(client *heroku.Client, pricingData *pricing.Data) tea.Cmd {
	return func() tea.Msg {
		analyzer := analysis.NewAnalyzer(client, pricingData)

		if err := analyzer.LoadData(); err != nil {
			return analysisCompleteMsg{err: err}
		}

		result, err := analyzer.Analyze()
		if err != nil {
			return analysisCompleteMsg{err: err}
		}

		return analysisCompleteMsg{result: result}
	}
}
