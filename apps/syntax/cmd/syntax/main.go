package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/ai"
	"github.com/kyanite/syntax/internal/app"
)

func main() {
	// Create the root model
	m := app.NewModel()

	// Initialize AI client (pkg/ai Brain backend)
	brainClient := ai.NewBrainClient()
	m.AIClient = brainClient

	// Generate session ID for this invocation
	sessionID := fmt.Sprintf("syntax-%d", time.Now().Unix())
	m.SessionID = sessionID

	// Attempt to resume the most recent session (best-effort)
	ctx := context.Background()
	sessions, _ := brainClient.GetRecentSessions(ctx, 1)
	if len(sessions) > 0 {
		fmt.Printf("Welcome back — resume '%s'?\n", sessions[0].Title)
	}

	// Create the Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running syntax.sh: %v\n", err)
		os.Exit(1)
	}

	// On shutdown: save session and cross-app context (best-effort)
	if fm, ok := finalModel.(app.Model); ok {
		title := fmt.Sprintf("session %s", time.Now().Format("2006-01-02"))
		if fm.CurrentProject != nil && fm.CurrentProject.Title != "" {
			title = fm.CurrentProject.Title
		}

		// Serialize serializable state for session persistence
		state := map[string]any{
			"screen":   int(fm.CurrentScreen),
			"project":  fm.CurrentProject != nil,
		}
		if fm.CurrentProject != nil {
			state["project_title"] = fm.CurrentProject.Title
			state["project_id"] = fm.CurrentProject.ID
		}

		_ = brainClient.SaveSession(ctx, fm.SessionID, title, state)

		// Cross-app context summary
		summary := "Active in syntax"
		if fm.CurrentProject != nil {
			summary = fmt.Sprintf("Editing '%s' in syntax", fm.CurrentProject.Title)
		}
		_ = brainClient.SaveCrossAppContext(ctx, "focus", "editing", summary, 0.7)
		_ = brainClient.SaveCrossAppContext(ctx, "noise", "editing", summary, 0.5)
	}

	// Clean up brain resources
	brainClient.Close()
}
