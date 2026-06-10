package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/ai"
)

// AISuggestionMsg is a message that carries AI suggestion results
type AISuggestionMsg struct {
	suggestion *ai.Suggestion
	err        error
}

func (m Model) viewAISuggestion() string {
	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(" AI Writing Assistant ")
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Heading.Render("🤖 AI Suggestions"))
	b.WriteString("\n\n")

	if m.AIGenerating {
		b.WriteString(m.Styles.Text.Render("⏳ Generating suggestion... Please wait.\n\n"))
		b.WriteString(m.Styles.Text.Faint(true).Render("This may take a few seconds depending on your Ollama setup.\n"))
		return b.String()
	}

	if m.AISuggestion != nil {
		// Show the suggestion
		b.WriteString(m.Styles.Text.Bold(true).Render("Suggestion:"))
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Text.Render(m.AISuggestion.Content))
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Text.Faint(true).Render(fmt.Sprintf("Generated at: %s\n\n",
			m.AISuggestion.Timestamp.Format("15:04:05"))))

		b.WriteString(m.Styles.Text.Render("Press Enter to insert, R to regenerate, Esc to cancel\n"))
	} else {
		// Show suggestion type menu
		b.WriteString(m.Styles.Text.Render("Select suggestion type:\n\n"))

		options := []struct {
			name string
			desc string
		}{
			{"Continue Story", "AI will continue writing from where you left off"},
			{"Improve Text", "AI will enhance the selected/recent text"},
			{"Suggest Dialogue", "AI will generate realistic dialogue for the scene"},
			{"Add Description", "AI will add vivid sensory details"},
			{"Character Action", "AI will suggest character behaviors and reactions"},
		}

		for i, opt := range options {
			prefix := "  "
			if i == m.SelectedIndex {
				prefix = "> "
			}

			style := m.Styles.Text
			if i == m.SelectedIndex {
				style = m.Styles.Accent
			}

			b.WriteString(style.Render(fmt.Sprintf("%s%s\n", prefix, opt.name)))
			if i == m.SelectedIndex {
				b.WriteString(m.Styles.Text.Faint(true).Render(fmt.Sprintf("     %s\n", opt.desc)))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(m.Styles.Text.Render("Press Enter to generate, Esc to cancel\n"))
	}

	// Footer
	if m.Message != "" {
		b.WriteString("\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	if m.Error != nil {
		b.WriteString("\n")
		b.WriteString(m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Error)))
	}

	return b.String()
}

func (m Model) handleAISuggestionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.AIGenerating {
		// Don't allow input while generating
		return m, nil
	}

	if m.AISuggestion != nil {
		// Suggestion is shown
		switch msg.String() {
		case "enter":
			// Insert suggestion at cursor
			if m.Buffer != nil {
				m.Buffer.InsertRune('\n')
				m.Buffer.InsertRune('\n')
				for _, r := range m.AISuggestion.Content {
					m.Buffer.InsertRune(r)
				}
				m.Message = "Suggestion inserted"
			}
			m.AISuggestion = nil
			m.CurrentScreen = m.PreviousScreen
			return m, nil

		case "r":
			// Regenerate
			m.AISuggestion = nil
			return m, nil

		case "esc":
			// Cancel
			m.AISuggestion = nil
			m.CurrentScreen = m.PreviousScreen
			return m, nil
		}
	} else {
		// Menu is shown
		switch msg.String() {
		case "up", "k":
			if m.SelectedIndex > 0 {
				m.SelectedIndex--
			}
			return m, nil

		case "down", "j":
			if m.SelectedIndex < 4 {
				m.SelectedIndex++
			}
			return m, nil

		case "enter":
			// Generate suggestion
			return m.generateAISuggestion()

		case "esc":
			m.CurrentScreen = m.PreviousScreen
			return m, nil
		}
	}

	return m, nil
}

func (m Model) generateAISuggestion() (Model, tea.Cmd) {
	if m.AIClient == nil || !m.AIClient.IsEnabled() {
		m.Error = fmt.Errorf("AI assistant is not enabled")
		return m, nil
	}

	if m.Buffer == nil {
		m.Error = fmt.Errorf("no content to generate from")
		return m, nil
	}

	// Map selected index to suggestion type
	var suggestionType ai.SuggestionType
	switch m.SelectedIndex {
	case 0:
		suggestionType = ai.SuggestionContinue
	case 1:
		suggestionType = ai.SuggestionImprove
	case 2:
		suggestionType = ai.SuggestionDialogue
	case 3:
		suggestionType = ai.SuggestionDescription
	case 4:
		suggestionType = ai.SuggestionCharacter
	default:
		suggestionType = ai.SuggestionContinue
	}

	// Get recent content for context
	content := m.Buffer.GetContent()
	lines := strings.Split(content, "\n")

	// Use last 20 lines or all if fewer
	startLine := max(0, len(lines)-20)
	contextLines := lines[startLine:]
	contentContext := strings.Join(contextLines, "\n")

	m.AIGenerating = true
	m.Error = nil
	m.Message = ""

	// Return command to generate suggestion asynchronously
	return m, func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		suggestion, err := m.AIClient.GetSuggestion(ctx, suggestionType, contentContext, "")
		return AISuggestionMsg{
			suggestion: suggestion,
			err:        err,
		}
	}
}

// HandleAISuggestionMsg handles the AI suggestion response
func (m Model) HandleAISuggestionMsg(msg AISuggestionMsg) Model {
	m.AIGenerating = false

	if msg.err != nil {
		m.Error = msg.err
		m.Message = ""
	} else {
		m.AISuggestion = msg.suggestion
		m.Error = nil
		m.Message = ""
	}

	return m
}
