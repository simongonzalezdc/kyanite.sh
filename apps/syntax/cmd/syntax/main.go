package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/ai"
	"github.com/kyanite/syntax/internal/app"
)

func main() {
	// Create the root model
	m := app.NewModel()

	// Initialize AI client (pkg/ai Brain backend)
	m.AIClient = ai.NewBrainClient()

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
