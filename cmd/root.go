package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/leaharmstrong/heroku-calc/internal/ui"
	"github.com/spf13/cobra"
)

var (
	// Flags
	projectPath  string
	appName      string
	dryRun       bool
	interactive  bool
	apply        bool
	exportReport string
)

var rootCmd = &cobra.Command{
	Use:   "heroku-calc",
	Short: "Heroku configuration analyzer for Rails applications",
	Long: `A BubbleTea TUI application for analyzing Heroku Rails application configurations,
providing performance recommendations, and optionally applying suggested changes.`,
	RunE: runRoot,
}

func init() {
	rootCmd.Flags().StringVarP(&projectPath, "project", "p", "", "Path to Rails project (default: current directory)")
	rootCmd.Flags().StringVarP(&appName, "app", "a", "", "Heroku app name (auto-detected from git if not specified)")
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would change without applying")
	rootCmd.Flags().BoolVar(&interactive, "interactive", false, "Interactively prompt for each change")
	rootCmd.Flags().BoolVar(&apply, "apply", false, "Apply all recommended changes (use with caution)")
	rootCmd.Flags().StringVarP(&exportReport, "export", "e", "", "Export markdown report to file (default: auto-generated filename)")
}

func Execute() error {
	return rootCmd.Execute()
}

func runRoot(cmd *cobra.Command, args []string) error {
	// Determine project path
	if projectPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		projectPath = cwd
	}

	// Resolve to absolute path
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		return fmt.Errorf("failed to resolve project path: %w", err)
	}
	projectPath = absPath

	// Launch the TUI
	return launchTUI()
}

func launchTUI() error {
	// Determine mode
	mode := determineMode()

	// Create and run the BubbleTea app
	p := tea.NewProgram(ui.NewModel(projectPath, appName, mode))

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	return nil
}

func determineMode() ui.AppMode {
	if apply {
		return ui.ModeApply
	}
	if interactive {
		return ui.ModeInteractive
	}
	if dryRun {
		return ui.ModeDryRun
	}
	return ui.ModeReadOnly
}
