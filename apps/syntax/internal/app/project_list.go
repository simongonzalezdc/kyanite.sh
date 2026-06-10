package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/syntax/internal/storage"
)

func (m Model) viewProjectList() string {
	var b strings.Builder

	title := m.Styles.Title.Render("Select a Project")
	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, title))
	b.WriteString("\n\n")

	if len(m.Projects) == 0 {
		b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.Styles.Text.Render("No projects found. Press 'n' to create one.")))
	} else {
		for i, proj := range m.Projects {
			var line string
			if i == m.SelectedIndex {
				line = m.Styles.MenuSelected.Render(fmt.Sprintf("▸ %s", proj.Title))
			} else {
				line = m.Styles.MenuUnselected.Render(fmt.Sprintf("  %s", proj.Title))
			}
			b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.PlaceHorizontal(m.Width, lipgloss.Center, m.Styles.Text.Render("↑/↓: Navigate  Enter: Open  Esc: Back")))

	return b.String()
}

func (m Model) handleProjectListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}
		return m, nil

	case "down", "j":
		if m.SelectedIndex < len(m.Projects)-1 {
			m.SelectedIndex++
		}
		return m, nil

	case "enter":
		if len(m.Projects) > 0 {
			project, err := storage.LoadProject(m.Projects[m.SelectedIndex].ID)
			if err != nil {
				m.Error = err
				return m, nil
			}
			m.CurrentProject = project
			m.CurrentScreen = ScreenEditor
			return m, nil
		}
		return m, nil

	case "esc":
		m.CurrentScreen = ScreenWelcome
		return m, nil
	}

	return m, nil
}
