package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/ai"
	"github.com/kyanite/appnames"
	"github.com/kyanite/config"
	"github.com/kyanite/syntax/internal/app"
	"github.com/kyanite/tui/aipanel"
)

func main() {
	m := app.NewModel()

	root, _ := config.Load()
	cfg := ai.ConfigFromRoot(root, "syntax")
	brain, err := ai.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: brain init failed (AI features offline): %v\n", err)
	}
	m.Brain = brain
	m.AIPanel = aipanel.New(brain, 40, 24)

	sessionID := fmt.Sprintf("syntax-%d", time.Now().Unix())
	m.SessionID = sessionID

	ctx := context.Background()
	if brain != nil {
		if sessions, _ := brain.GetRecentSessions(ctx, 1); len(sessions) > 0 {
			fmt.Printf("Welcome back — resume '%s'?\n", sessions[0].Title)
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error running syntax.sh: %v\n", err)
		os.Exit(1)
	}

	if fm, ok := finalModel.(app.Model); ok {
		if brain == nil {
			return
		}
		title := fmt.Sprintf("session %s", time.Now().Format("2006-01-02"))
		if fm.CurrentProject != nil && fm.CurrentProject.Title != "" {
			title = fm.CurrentProject.Title
		}
		state := map[string]any{
			"screen":  int(fm.CurrentScreen),
			"project": fm.CurrentProject != nil,
		}
		if fm.CurrentProject != nil {
			state["project_title"] = fm.CurrentProject.Title
			state["project_id"] = fm.CurrentProject.ID
		}
		_ = brain.SaveSession(ctx, fm.SessionID, title, state)
		summary := "Active in syntax"
		if fm.CurrentProject != nil {
			summary = fmt.Sprintf("Editing '%s' in syntax", fm.CurrentProject.Title)
		}
		_ = brain.SaveCrossAppContext(ctx, appnames.Focus, "editing", summary, 0.7)
		_ = brain.SaveCrossAppContext(ctx, appnames.Noise, "editing", summary, 0.5)
	}

	if brain != nil {
		brain.Close()
	}
}
