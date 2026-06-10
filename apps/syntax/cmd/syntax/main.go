package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/ai"
	"github.com/kyanite/syntax/internal/app"
	"github.com/kyanite/syntax/internal/storage"
)

func main() {
	// Load configuration
	config, err := storage.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config, using defaults: %v\n", err)
		config = &storage.AppConfig{
			AI: ai.DefaultConfig(),
		}
	}

	// Create the root model
	m := app.NewModel()

	// Initialize AI client if enabled
	if config.AI.Enabled {
		m.AIClient = ai.NewClient(config.AI)
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running syntax.sh: %v\n", err)
		os.Exit(1)
	}
}
