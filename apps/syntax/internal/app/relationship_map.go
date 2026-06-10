package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/visualization"
)

type MapViewMode int

const (
	MapViewDetailed MapViewMode = iota
	MapViewMatrix
)

func (m Model) viewRelationshipMap() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	// Data is loaded in Update via ensureDataLoaded()

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(" Character Relationship Map ")
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	if len(m.CurrentProject.Characters) == 0 {
		b.WriteString(m.Styles.Text.Render("No characters to visualize.\n\n"))
		b.WriteString(m.Styles.Text.Render("Create some characters first!\n"))
		b.WriteString("\n")
		b.WriteString(m.Styles.Text.Render("Press Esc to go back"))
		return b.String()
	}

	// Determine which view mode to show
	// Default to detailed view, use m.SelectedIndex to toggle
	var mapOutput string
	if m.SelectedIndex == 1 {
		mapOutput = visualization.GenerateRelationshipMatrix(m.CurrentProject.Characters)
	} else {
		mapOutput = visualization.GenerateRelationshipMap(m.CurrentProject.Characters)
	}

	b.WriteString(mapOutput)
	b.WriteString("\n\n")

	// Controls
	viewMode := "Detailed View"
	if m.SelectedIndex == 1 {
		viewMode = "Matrix View"
	}

	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Current: %s\n", viewMode)))
	b.WriteString(m.Styles.Text.Render("Press 'v' to toggle view | Esc to go back\n"))

	// Footer
	if m.Message != "" {
		b.WriteString("\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	return b.String()
}

func (m Model) handleRelationshipMapKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "v":
		// Toggle between views
		if m.SelectedIndex == 0 {
			m.SelectedIndex = 1
		} else {
			m.SelectedIndex = 0
		}
		return m, nil

	case "esc":
		m.SelectedIndex = 0
		m.CurrentScreen = ScreenCharacters
		return m, nil
	}

	return m, nil
}
