// Package main is the entry point for the noise.sh TUI application.
package main

import (
	"fmt"
	"os"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/logging"
	"github.com/kyanite/noise/internal/plugins"
	"github.com/kyanite/noise/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Initialize logging
	logger := logging.GetDefaultLogger()
	logger.Infof("Starting noise.sh version=%s commit=%s date=%s", version, commit, date)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Warnf("Failed to load config, using defaults: %v", err)
		cfg = config.DefaultConfig()
	}

	// Initialize plugin manager with config and logger
	pluginManager := plugins.NewManager(cfg, logger)

	// Create root model
	model := ui.NewRootModel(pluginManager)

	// Create Bubble Tea program with proper options for TUI rendering
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support for better UX
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running noise.sh: %v\n", err)
		os.Exit(1)
	}
}
