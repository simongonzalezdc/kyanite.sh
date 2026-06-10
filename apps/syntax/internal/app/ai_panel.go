package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// continueWriting generates 3 continuation paragraphs via streaming AI.
// Prompt: "Continue this story naturally. Match the tone, style, and character voices. Write 3 paragraphs."
func (m Model) continueWriting() (tea.Model, tea.Cmd) {
	if m.AIClient == nil || !m.AIClient.IsEnabled() {
		m.Message = "AI assistant is not enabled"
		return m, nil
	}
	if m.Buffer == nil {
		m.Message = "No content to continue from"
		return m, nil
	}

	content := m.Buffer.GetContent()
	lines := strings.Split(content, "\n")
	startLine := len(lines) - 20
	if startLine < 0 {
		startLine = 0
	}
	recentText := strings.Join(lines[startLine:], "\n")

	prompt := fmt.Sprintf(
		"Continue this story naturally. Match the tone, style, and character voices. Write 3 paragraphs.\nCurrent text:\n%s",
		recentText,
	)

	m.AIPanel = m.AIPanel.StartStream("Continue Writing")
	m.AIPanelInput = ""

	cmd := m.AIPanel.Generate(prompt)
	return m, cmd
}

// submitAIPanelCommand processes check:/voice: commands typed into the AI panel input.
func (m Model) submitAIPanelCommand() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.AIPanelInput)
	m.AIPanelInput = ""

	if input == "" {
		return m, nil
	}

	if m.Buffer == nil {
		m.Message = "No content to analyze"
		return m, nil
	}

	switch {
	case strings.HasPrefix(input, "check:"):
		return m.consistencyCheck(strings.TrimPrefix(input, "check:"))
	case strings.HasPrefix(input, "voice:"):
		charName := strings.TrimSpace(strings.TrimPrefix(input, "voice:"))
		return m.characterVoice(charName)
	default:
		m.Message = "Unknown command. Use 'check:' or 'voice:name'"
		return m, nil
	}
}

// consistencyCheck checks the current text for contradictions with earlier chapters.
func (m Model) consistencyCheck(query string) (tea.Model, tea.Cmd) {
	if m.AIClient == nil || !m.AIClient.IsEnabled() {
		m.Message = "AI assistant is not enabled"
		return m, nil
	}

	content := m.Buffer.GetContent()

	// Build story context from project scenes if available
	context := ""
	if m.CurrentProject != nil && len(m.CurrentProject.Scenes) > 0 {
		var ctxParts []string
		for _, s := range m.CurrentProject.Scenes {
			if s.Content != "" && s.ID != m.CurrentScene.ID {
				ctxParts = append(ctxParts, s.Content)
			}
		}
		if len(ctxParts) > 0 {
			// Limit context to avoid overly large prompts
			context = strings.Join(ctxParts, "\n---\n")
			if len(context) > 2000 {
				context = context[:2000] + "..."
			}
		}
	}

	prompt := fmt.Sprintf(
		"Review this passage for internal consistency. Check character names, locations, timeline. Flag any contradictions.\nCurrent chapter:\n%s\n\nStory context: %s",
		content,
		context,
	)

	m.AIPanel = m.AIPanel.StartStream("Consistency Check")

	cmd := m.AIPanel.Generate(prompt)
	return m, cmd
}

// characterVoice rewrites dialogue in the given character's voice.
func (m Model) characterVoice(charName string) (tea.Model, tea.Cmd) {
	if m.AIClient == nil || !m.AIClient.IsEnabled() {
		m.Message = "AI assistant is not enabled"
		return m, nil
	}

	if charName == "" {
		m.Message = "Usage: voice:CharacterName"
		return m, nil
	}

	content := m.Buffer.GetContent()

	// Build character context from project characters if available
	charContext := ""
	if m.CurrentProject != nil {
		for _, ch := range m.CurrentProject.Characters {
			if ch.Name == charName {
				parts := []string{fmt.Sprintf("Name: %s", ch.Name)}
				if ch.Role != "" {
					parts = append(parts, fmt.Sprintf("Role: %s", ch.Role))
				}
				if ch.Background != "" {
					parts = append(parts, fmt.Sprintf("Background: %s", ch.Background))
				}
				if ch.Bio != "" {
					parts = append(parts, fmt.Sprintf("Bio: %s", ch.Bio))
				}
				charContext = strings.Join(parts, "\n")
				break
			}
		}
	}

	prompt := fmt.Sprintf(
		"Rewrite this dialogue in %s's voice based on their established speech patterns:\n%s",
		charName,
		content,
	)
	if charContext != "" {
		prompt = fmt.Sprintf(
			"Character profile:\n%s\n\nRewrite this dialogue in %s's voice based on their established speech patterns:\n%s",
			charContext,
			charName,
			content,
		)
	}

	m.AIPanel = m.AIPanel.StartStream(fmt.Sprintf("Voice: %s", charName))

	cmd := m.AIPanel.Generate(prompt)
	return m, cmd
}

// renderWithAIPanel renders the main view alongside the AI writing partner panel.
func (m Model) renderWithAIPanel(mainView string) string {
	// Resize the AI panel to fit the right portion of the screen
	panelWidth := max(m.Width/3, 40)
	m.AIPanel = m.AIPanel.SetSize(panelWidth, m.Height)

	panelView := m.AIPanel.View()

	// Build the input line for the panel
	inputLine := ""
	if m.AIPanelInput != "" {
		inputLine = fmt.Sprintf("  > %s█", m.AIPanelInput)
	} else if !m.AIPanel.IsStreaming() {
		inputLine = "  type check: or voice:name, then Enter"
	}

	// Combine main view and panel side by side
	mainLines := strings.Split(mainView, "\n")
	panelLines := strings.Split(panelView, "\n")

	// Calculate main view width
	mainWidth := m.Width - panelWidth - 1
	if mainWidth < 20 {
		mainWidth = 20
	}

	maxLines := max(len(mainLines), len(panelLines))
	var b strings.Builder
	for i := 0; i < maxLines && i < m.Height-1; i++ {
		var mainLine, panelLine string
		if i < len(mainLines) {
			mainLine = mainLines[i]
		}
		if i < len(panelLines) {
			panelLine = panelLines[i]
		}

		// Truncate/pad main line
		if len(mainLine) > mainWidth {
			mainLine = mainLine[:mainWidth]
		}

		b.WriteString(mainLine)
		b.WriteString("│")
		b.WriteString(panelLine)
		b.WriteString("\n")
	}

	// Add input line at bottom of panel
	if inputLine != "" {
		b.WriteString(strings.Repeat(" ", mainWidth))
		b.WriteString("│")
		b.WriteString(inputLine)
		b.WriteString("\n")
	}

	return b.String()
}
