package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) viewEditor() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s | %s ", m.CurrentProject.Title, m.CurrentTheme.Name))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	// Project dashboard
	b.WriteString(m.Styles.Heading.Render("📖 Project Dashboard"))
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Title: %s\n", m.CurrentProject.Title)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Author: %s\n", m.CurrentProject.Author)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Genre: %s\n", m.CurrentProject.Genre)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Status: %s\n", m.CurrentProject.Status)))
	b.WriteString("\n")

	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Words: %d\n", m.CurrentProject.TotalWords)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Scenes: %d\n", m.CurrentProject.TotalScenes)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Characters: %d\n", m.CurrentProject.TotalCharacters)))
	b.WriteString("\n")

	// Navigation menu
	b.WriteString(m.Styles.Heading.Render("Navigation"))
	b.WriteString("\n\n")

	menu := []string{
		"c - Characters",
		"s - Scenes",
		"l - Locations",
		"b - Backups",
		"t - Statistics",
		"e - Export",
		"h - Help",
		"q - Quit",
	}

	for _, item := range menu {
		b.WriteString(m.Styles.Text.Render(fmt.Sprintf("  %s\n", item)))
	}

	// Footer
	if m.Message != "" {
		b.WriteString("\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	return b.String()
}

func (m Model) handleEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "c":
		m.CurrentScreen = ScreenCharacters
		return m, nil

	case "s":
		m.CurrentScreen = ScreenScenes
		return m, nil

	case "l":
		m.CurrentScreen = ScreenLocations
		return m, nil

	case "b":
		// Backups - navigate to backup screen
		m.SelectedIndex = 0
		m.Message = ""
		m.Error = nil
		m.CurrentScreen = ScreenBackups
		return m, nil

	case "t":
		m.CurrentScreen = ScreenStats
		return m, nil

	case "e":
		// Export - navigate to export screen
		m.SelectedIndex = 0
		m.Message = ""
		m.Error = nil
		m.CurrentScreen = ScreenExport
		return m, nil

	case "h":
		m.PreviousScreen = ScreenEditor
		m.CurrentScreen = ScreenHelp
		return m, nil

	case "esc":
		m.CurrentScreen = ScreenWelcome
		return m, nil
	}

	return m, nil
}
