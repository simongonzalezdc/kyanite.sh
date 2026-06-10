package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/editor"
	"github.com/kyanite/syntax/internal/storage"
)

// LocationEditorField represents which field is being edited
type LocationEditorField int

const (
	FieldLocationName LocationEditorField = iota
	FieldLocationType
	FieldLocationRegion
	FieldLocationClimate
	FieldLocationPopulation
	FieldLocationSignificance
	FieldLocationDescription
)

func (m Model) viewLocationEditor() string {
	if m.CurrentProject == nil || m.CurrentLocation == nil {
		return "No location loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Editing Location: %s ",
		m.CurrentProject.Title, m.CurrentLocation.Name))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	// Field editor interface
	b.WriteString(m.Styles.Heading.Render("📍 Location Details"))
	b.WriteString("\n\n")

	// Determine which field is selected based on SelectedIndex
	fields := []struct {
		label string
		value string
		field LocationEditorField
	}{
		{"Name", m.CurrentLocation.Name, FieldLocationName},
		{"Type", m.CurrentLocation.Type, FieldLocationType},
		{"Region", m.CurrentLocation.Region, FieldLocationRegion},
		{"Climate", m.CurrentLocation.Climate, FieldLocationClimate},
		{"Population", fmt.Sprintf("%d", m.CurrentLocation.Population), FieldLocationPopulation},
		{"Significance", m.CurrentLocation.Significance, FieldLocationSignificance},
		{"Description", m.CurrentLocation.Description, FieldLocationDescription},
	}

	for i, field := range fields {
		prefix := "  "
		style := m.Styles.Text

		if i == m.SelectedIndex {
			prefix = "> "
			style = m.Styles.Accent
		}

		// If in input mode on this field, show input
		if m.InputMode && i == m.SelectedIndex {
			line := fmt.Sprintf("%s%s: %s█", prefix, field.label, m.InputValue)
			b.WriteString(style.Render(line))
		} else {
			displayValue := field.value
			if displayValue == "" {
				displayValue = "(not set)"
			}
			// Truncate long descriptions
			if field.field == FieldLocationDescription && len(displayValue) > 60 {
				displayValue = displayValue[:57] + "..."
			}
			line := fmt.Sprintf("%s%s: %s", prefix, field.label, displayValue)
			b.WriteString(style.Render(line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.InputMode {
		b.WriteString(m.Styles.Text.Render("Enter - Save field | Esc - Cancel"))
	} else {
		b.WriteString(m.Styles.Text.Render("Enter - Edit field | d - Edit full description | Ctrl+S - Save & Exit | Esc - Cancel"))
	}

	// Footer
	if m.Message != "" {
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	if m.Error != nil {
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Error)))
	}

	return b.String()
}

func (m Model) handleLocationEditorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.CurrentLocation == nil {
		return m, nil
	}

	// Handle input mode
	if m.InputMode {
		switch msg.String() {
		case "enter":
			// Save the field
			m = m.saveLocationField()
			m.InputMode = false
			m.InputValue = ""
			return m, nil

		case "esc":
			m.InputMode = false
			m.InputValue = ""
			m.Message = ""
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

	// Normal navigation mode
	switch msg.String() {
	case "up", "k":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}
		return m, nil

	case "down", "j":
		if m.SelectedIndex < 6 { // 7 fields total (0-6)
			m.SelectedIndex++
		}
		return m, nil

	case "enter":
		// Enter edit mode for selected field
		m.InputMode = true
		m.InputValue = m.getCurrentFieldValue()
		return m, nil

	case "d":
		// Edit full description in text editor
		if m.SelectedIndex == int(FieldLocationDescription) || msg.String() == "d" {
			m.Buffer = editor.NewBuffer(m.CurrentLocation.Description)
			m.EditorMode = EditorModeNormal
			m.PreviousScreen = ScreenLocationEditor
			// Note: We can reuse text editor but will need special handling to save back to location
			// For now, use inline editing only
			m.Message = "Use Enter to edit description inline, or Ctrl+S to save all changes"
		}
		return m, nil

	case "ctrl+s":
		// Save and exit
		if err := storage.SaveLocation(m.CurrentProject.Directory, m.CurrentLocation); err != nil {
			m.Error = err
		} else {
			m.Message = "Location saved"
			// Update in project's locations map
			m.CurrentProject.Locations[m.CurrentLocation.ID] = m.CurrentLocation

			// Return to locations list
			m.CurrentLocation = nil
			m.Buffer = nil
			m.SelectedIndex = 0
			m.CurrentScreen = ScreenLocations
		}
		return m, nil

	case "esc":
		// Exit without explicit save (changes are already saved per-field)
		m.CurrentLocation = nil
		m.Buffer = nil
		m.SelectedIndex = 0
		m.Message = ""
		m.Error = nil
		m.CurrentScreen = ScreenLocations
		return m, nil
	}

	return m, nil
}

// getCurrentFieldValue returns the current value of the selected field
func (m Model) getCurrentFieldValue() string {
	if m.CurrentLocation == nil {
		return ""
	}

	switch LocationEditorField(m.SelectedIndex) {
	case FieldLocationName:
		return m.CurrentLocation.Name
	case FieldLocationType:
		return m.CurrentLocation.Type
	case FieldLocationRegion:
		return m.CurrentLocation.Region
	case FieldLocationClimate:
		return m.CurrentLocation.Climate
	case FieldLocationPopulation:
		if m.CurrentLocation.Population > 0 {
			return fmt.Sprintf("%d", m.CurrentLocation.Population)
		}
		return ""
	case FieldLocationSignificance:
		return m.CurrentLocation.Significance
	case FieldLocationDescription:
		return m.CurrentLocation.Description
	default:
		return ""
	}
}

// saveLocationField saves the current input to the selected field
func (m Model) saveLocationField() Model {
	if m.CurrentLocation == nil {
		return m
	}

	value := m.InputValue

	switch LocationEditorField(m.SelectedIndex) {
	case FieldLocationName:
		m.CurrentLocation.Name = value
		m.Message = "Name updated"

	case FieldLocationType:
		m.CurrentLocation.Type = value
		m.Message = "Type updated"

	case FieldLocationRegion:
		m.CurrentLocation.Region = value
		m.Message = "Region updated"

	case FieldLocationClimate:
		m.CurrentLocation.Climate = value
		m.Message = "Climate updated"

	case FieldLocationPopulation:
		// Parse population as integer
		var pop int
		fmt.Sscanf(value, "%d", &pop)
		m.CurrentLocation.Population = pop
		m.Message = "Population updated"

	case FieldLocationSignificance:
		m.CurrentLocation.Significance = value
		m.Message = "Significance updated"

	case FieldLocationDescription:
		m.CurrentLocation.Description = value
		m.Message = "Description updated"
	}

	// Auto-save after each field change
	if err := storage.SaveLocation(m.CurrentProject.Directory, m.CurrentLocation); err != nil {
		m.Error = err
		m.Message = ""
	}

	return m
}
