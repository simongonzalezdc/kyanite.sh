package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/syntax/internal/storage"
)

func (m Model) viewWelcome() string {
	var b strings.Builder

	title := m.Styles.Title.Render("✎ syntax.sh")
	subtitle := m.Styles.Text.Render("Terminal-Based Fiction Writing Tool")

	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, title))
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, subtitle))
	b.WriteString("\n\n")

	// Menu
	menu := []string{
		"n - New Project",
		"o - Open Project",
		"? - Help",
		"q - Quit",
	}

	for _, item := range menu {
		b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.Styles.Text.Render(item)))
		b.WriteString("\n")
	}

	// Message
	if m.Message != "" {
		b.WriteString("\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.Styles.Success.Render(m.Message)))
	}

	if m.Error != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.Styles.Error.Render(m.Error.Error())))
	}

	return b.String()
}

func (m Model) handleWelcomeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "n":
		// Create new project (simplified for MVP)
		project, err := storage.CreateProject("My Novel", "Author", "Fiction")
		if err != nil {
			m.Error = err
			return m, nil
		}
		m.CurrentProject = project
		m.CurrentScreen = ScreenEditor
		m.Message = fmt.Sprintf("Created project: %s", project.Title)
		return m, nil

	case "o":
		// Load project list
		projects, err := storage.ListProjects()
		if err != nil {
			m.Error = err
			return m, nil
		}
		m.Projects = projects
		m.CurrentScreen = ScreenProjectList
		m.SelectedIndex = 0
		return m, nil

	case "q":
		return m, tea.Quit
	}

	return m, nil
}
