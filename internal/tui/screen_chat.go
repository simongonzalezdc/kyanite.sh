package tui

import (
	"strings"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
)

func (m MainModel) renderChatView() string {
	var b strings.Builder

	// CONSISTENT: Standard header styling
	header := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Background(styles.GetPanel()).
		Bold(true).
		Italic(true).
		Align(lipgloss.Center).
		Padding(1, 3).           // CONSISTENT: 1 vertical, 3 horizontal
		BorderStyle(lipgloss.RoundedBorder()).  // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Underline(true).
		Render("💬 CHAT ASSISTANT 🤖")
	b.WriteString(header)
	b.WriteString("\n\n")

	// CONSISTENT: Input box styling
	inputBox := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).           // CONSISTENT: 1 vertical, 2 horizontal
		Margin(1, 0, 1, 0).      // CONSISTENT: 1 top/bottom
		Border(lipgloss.RoundedBorder()).  // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Width(m.width - 4).
		Height(3)

	// Current message (with wrapping)
	maxWidth := m.width - 8
	message := m.chatInput
	if len(message) > maxWidth-2 {
		// Simple truncation for long messages
		message = message[:maxWidth-2] + ".."
	}

	currentMessage := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Render(fmt.Sprintf("> %s", message))

	b.WriteString(inputBox.Render(currentMessage))

	b.WriteString("\n\n")

	// CONSISTENT: Chat history styling
	historyBox := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).           // CONSISTENT: 1 vertical, 2 horizontal
		Margin(1, 0, 1, 0).      // CONSISTENT: 1 top/bottom
		Border(lipgloss.RoundedBorder()).  // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Width(m.width - 4).
		Height(m.height - 15)

	var historyContent strings.Builder
	for i, msg := range m.chatHistory {
		if i > len(m.chatHistory)-6 { // Show last 6 messages
			historyContent.WriteString(msg + "\n")
		}
	}

	b.WriteString(historyBox.Render(historyContent.String()))

	b.WriteString("\n\n")

	// CONSISTENT: Instructions styling
	instructions := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Italic(true).
		Align(lipgloss.Center).
		Render("Type your message and press Enter to send. Press Escape to return to dashboard.")

	b.WriteString(instructions)

	return b.String()
}