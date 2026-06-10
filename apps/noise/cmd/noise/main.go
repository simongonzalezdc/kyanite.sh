// Package main is the entry point for the noise.sh TUI application.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/infra/brain"
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

	// Create brain client for session lifecycle
	brainClient := brain.NewClient()
	sessionID := fmt.Sprintf("noise-%d", time.Now().Unix())
	b := brainClient.Brain() // may be nil if Brain init failed (offline mode)

	// Attempt to load the most recent session on startup
	if b != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if recent, err := b.GetRecentSessions(ctx, 1); err == nil && len(recent) > 0 {
			title := recent[0].Title
			if title == "" {
				title = "untitled"
			}
			fmt.Printf("Welcome back — last project: %s\n", title)
		} else if err != nil {
			logger.Debugf("session restore skipped: %v", err)
		}
		cancel()
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
	finalModel, runErr := p.Run()
	if runErr != nil {
		fmt.Fprintf(os.Stderr, "Error running noise.sh: %v\n", runErr)
	}

	// --- Shutdown: save session and cross-app context ---
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Derive session title from the current song if available
	sessionTitle := fmt.Sprintf("session %s", time.Now().Format("2006-01-02"))
	if root, ok := finalModel.(*ui.RootModel); ok {
		if songTitle := root.GetCurrentSongTitle(); songTitle != "" {
			sessionTitle = songTitle
		}
	}

	// Best-effort session save
	if b != nil {
		if err := b.SaveSession(shutdownCtx, sessionID, sessionTitle, nil); err != nil {
			logger.Warnf("failed to save session: %v", err)
		}

		// Best-effort cross-app context save
		summary := fmt.Sprintf("Used noise.sh — %s", sessionTitle)
		if err := b.SaveCrossAppContext(shutdownCtx, "syntax", "session_summary", summary, 0.5); err != nil {
			logger.Warnf("failed to save cross-app context: %v", err)
		}
		if err := b.SaveCrossAppContext(shutdownCtx, "focus", "session_summary", summary, 0.5); err != nil {
			logger.Warnf("failed to save cross-app context: %v", err)
		}
	}

	brainClient.Close()

	if runErr != nil {
		os.Exit(1)
	}
}
