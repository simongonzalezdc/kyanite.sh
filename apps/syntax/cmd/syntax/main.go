package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/ai"
	"github.com/kyanite/syntax/internal/app"
	"github.com/kyanite/tui/aipanel"
)

func main() {
	// Create the root model
	m := app.NewModel()

	// Initialize AI inference brain (pkg/ai). Brain may be nil if the NUCBox
	// is unreachable; the app degrades gracefully and works offline.
	cfg := ai.DefaultConfig("syntax")
	brain, _ := ai.New(cfg)
	m.Brain = brain
	// Initialize AI writing partner panel (aipanel)
	m.AIPanel = aipanel.New(brain, 40, 24)

	// Generate session ID for this invocation
	sessionID := fmt.Sprintf("syntax-%d", time.Now().Unix())
	m.SessionID = sessionID

	// Attempt to resume the most recent session (best-effort)
	ctx := context.Background()
	if brain != nil {
		if sessions, _ := brain.GetRecentSessions(ctx, 1); len(sessions) > 0 {
			fmt.Printf("Welcome back — resume '%s'?\n", sessions[0].Title)
		}
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
		if brain == nil {
			return
		}

		title := fmt.Sprintf("session %s", time.Now().Format("2006-01-02"))
		if fm.CurrentProject != nil && fm.CurrentProject.Title != "" {
			title = fm.CurrentProject.Title
		}

		// Serialize serializable state for session persistence
		state := map[string]any{
			"screen":  int(fm.CurrentScreen),
			"project": fm.CurrentProject != nil,
		}
		if fm.CurrentProject != nil {
			state["project_title"] = fm.CurrentProject.Title
			state["project_id"] = fm.CurrentProject.ID
		}

		_ = brain.SaveSession(ctx, fm.SessionID, title, state)

		// Cross-app context summary
		summary := "Active in syntax"
		if fm.CurrentProject != nil {
			summary = fmt.Sprintf("Editing '%s' in syntax", fm.CurrentProject.Title)
		}
		_ = brain.SaveCrossAppContext(ctx, "focus", "editing", summary, 0.7)
		_ = brain.SaveCrossAppContext(ctx, "noise", "editing", summary, 0.5)
	}

	// Clean up brain resources
	if brain != nil {
		brain.Close()
	}
}
