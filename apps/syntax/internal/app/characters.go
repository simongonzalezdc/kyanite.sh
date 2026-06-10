package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/storage"
)

func (m Model) viewCharacters() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Characters ", m.CurrentProject.Title))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	// Data is loaded in Update via ensureDataLoaded()

	b.WriteString(m.Styles.Heading.Render("👥 Characters"))
	b.WriteString("\n\n")

	if len(m.CurrentProject.Characters) == 0 {
		b.WriteString(m.Styles.Text.Render("No characters yet. Press 'n' to create one."))
	} else {
		for _, char := range m.CurrentProject.Characters {
			b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("• %s", char.Name)))
			if char.Role != "" {
				b.WriteString(m.Styles.Text.Render(fmt.Sprintf(" (%s)", char.Role)))
			}
			b.WriteString("\n")
			if char.Occupation != "" {
				b.WriteString(m.Styles.Text.Render(fmt.Sprintf("  %s", char.Occupation)))
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n")

	// Show input prompt if in input mode
	if m.InputMode {
		b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("Enter character name: %s█", m.InputValue)))
		b.WriteString("\n")
	} else {
		b.WriteString(m.Styles.Text.Render("n - New Character | r - Relationship Map | Esc - Back"))
	}

	if m.Message != "" {
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	return b.String()
}

func (m Model) handleCharactersKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input mode for character creation
	if m.InputMode {
		switch msg.String() {
		case "enter":
			// Create character with user-provided name
			if m.InputValue != "" {
				char := &character.Character{
					Name: m.InputValue,
					Role: "", // Let user fill in editor
				}

				err := storage.SaveCharacter(m.CurrentProject.Directory, char)
				if err != nil {
					m.Error = err
				} else {
					if m.CurrentProject.Characters == nil {
						m.CurrentProject.Characters = make(map[string]*character.Character)
					}
					m.CurrentProject.Characters[char.ID] = char
					m.CurrentProject.TotalCharacters++
					m.Message = fmt.Sprintf("Created character: %s", char.Name)

					// Save updated project metadata
					if err := storage.SaveProjectMetadata(m.CurrentProject); err != nil {
						m.Error = fmt.Errorf("failed to save project: %w", err)
					}
				}
			}
			m.InputMode = false
			m.InputValue = ""
			return m, nil

		case "esc":
			m.InputMode = false
			m.InputValue = ""
			return m, nil

		case "backspace":
			if len(m.InputValue) > 0 {
				m.InputValue = m.InputValue[:len(m.InputValue)-1]
			}
			return m, nil

		default:
			// Add character to input
			if len(msg.String()) == 1 {
				m.InputValue += msg.String()
			}
			return m, nil
		}
	}

	// Normal mode
	switch msg.String() {
	case "n":
		// Enter input mode for character creation
		m.InputMode = true
		m.InputValue = ""
		return m, nil

	case "r":
		// Show relationship map
		m.SelectedIndex = 0
		m.CurrentScreen = ScreenRelationshipMap
		return m, nil

	case "esc":
		m.CurrentScreen = ScreenEditor
		m.Message = ""
		return m, nil
	}

	return m, nil
}
