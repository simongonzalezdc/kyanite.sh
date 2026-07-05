package glow

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
)

// GlowStyler provides enhanced styling using Glow and Lip Gloss
type GlowStyler struct {
	theme string
}

// NewGlowStyler creates a new Glow styler with the specified theme
func NewGlowStyler(theme string) *GlowStyler {
	return &GlowStyler{
		theme: theme,
	}
}

// RenderTaskWithGlow renders a task with enhanced Glow styling
func (gs *GlowStyler) RenderTaskWithGlow(id, description, status, priority string, index int) string {
	// Status indicator
	statusIcon := "⏳"
	statusColor := styles.GetAccent()
	switch status {
	case "completed":
		statusIcon = "✅"
		statusColor = styles.GetSuccess()
	case "in-progress":
		statusIcon = "🔄"
		statusColor = styles.GetWarning()
	}

	// Priority indicator
	priorityIcon := "⚪"
	priorityColor := styles.GetBorder()
	switch priority {
	case "high":
		priorityIcon = "🔴"
		priorityColor = styles.GetError()
	case "medium":
		priorityIcon = "🟡"
		priorityColor = styles.GetWarning()
	case "low":
		priorityIcon = "🟢"
		priorityColor = styles.GetSuccess()
	}

	// Create enhanced task styling
	taskStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetAccent()).
		Padding(0, 1).
		MarginBottom(1)

	content := fmt.Sprintf("%s %s %s | %s | %s",
		lipgloss.Style{}.Foreground(statusColor).Render(statusIcon),
		lipgloss.Style{}.Foreground(priorityColor).Render(priorityIcon),
		lipgloss.Style{}.Foreground(styles.GetAccent()).Bold(true).Render(fmt.Sprintf("Task %d", index)),
		lipgloss.Style{}.Foreground(styles.GetAccent()).Render(id),
		lipgloss.Style{}.Foreground(styles.GetForeground()).Italic(true).Render(description),
	)

	return taskStyle.Render(content)
}

// RenderHeaderWithGlow renders a header with enhanced Glow styling
func (gs *GlowStyler) RenderHeaderWithGlow(title, subtitle string) string {
	// Create enhanced header styling
	headerStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetAccent()).
		Padding(1, 2).
		MarginBottom(1)

	// Title styling
	titleStyle := lipgloss.Style{}.
		Foreground(styles.GetAccent()).
		Bold(true).
		Render(title)

	// Subtitle styling
	subtitleStyle := lipgloss.Style{}.
		Foreground(styles.GetAccent()).
		Render(subtitle)

	content := fmt.Sprintf("%s\n%s", titleStyle, subtitleStyle)
	return headerStyle.Render(content)
}

// RenderSectionWithGlow renders a section with Glow-enhanced styling
func (gs *GlowStyler) RenderSectionWithGlow(title string, content string) string {
	// Title bar styling
	titleStyle := lipgloss.Style{}.
		Foreground(styles.GetAccent()).
		Bold(true).
		Render(fmt.Sprintf("─ %s ─", title))

	// Content styling
	contentStyle := lipgloss.Style{}.
		Foreground(styles.GetForeground()). // Dark text
		Background(styles.GetPanel()).      // Panel background
		Padding(1, 2).
		Width(80).
		Render(content)

	sectionContent := fmt.Sprintf("%s\n%s\n%s", titleStyle, contentStyle, strings.Repeat("─", 30))
	return sectionContent
}

// RenderSuccessWithGlow renders success message with Glow styling
func (gs *GlowStyler) RenderSuccessWithGlow(message string) string {
	successStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetSuccess()).
		Foreground(styles.GetSuccess()).
		Bold(true).
		Padding(0, 2).
		MarginTop(1)

	return successStyle.Render(fmt.Sprintf("✅ %s", message))
}

// RenderErrorWithGlow renders error message with Glow styling
func (gs *GlowStyler) RenderErrorWithGlow(message string) string {
	errorStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetError()).
		Foreground(styles.GetError()).
		Bold(true).
		Padding(0, 2).
		MarginTop(1)

	return errorStyle.Render(fmt.Sprintf("❌ %s", message))
}

// RenderInfoWithGlow renders info message with Glow styling
func (gs *GlowStyler) RenderInfoWithGlow(message string) string {
	infoStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetAccent()).
		Foreground(styles.GetAccent()).
		Padding(0, 2).
		MarginTop(1)

	return infoStyle.Render(fmt.Sprintf("💡 %s", message))
}

// RenderProgressBarWithGlow renders a progress bar with Glow styling
func (gs *GlowStyler) RenderProgressBarWithGlow(current, total int, label string) string {
	percentage := float64(current) / float64(total)
	barWidth := 30
	filledWidth := int(percentage * float64(barWidth))

	// Create progress bar
	var progressBar strings.Builder
	progressBar.WriteString("[")

	// Filled portion
	filled := strings.Repeat("█", filledWidth)
	progressBar.WriteString(lipgloss.Style{}.
		Foreground(styles.GetSuccess()).
		Render(filled))

	// Empty portion
	empty := strings.Repeat("░", barWidth-filledWidth)
	progressBar.WriteString(lipgloss.Style{}.
		Foreground(styles.GetBorder()). // Grey for empty portion
		Render(empty))

	progressBar.WriteString("]")

	// Progress text
	progressText := fmt.Sprintf("%s %d/%d (%.1f%%)",
		lipgloss.Style{}.Foreground(styles.GetAccent()).Render(label),
		current, total, percentage*100)

	// Combine with styling
	barStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetAccent()).
		Padding(0, 1).
		MarginTop(1)

	barContent := fmt.Sprintf("%s\n%s", progressBar.String(), progressText)
	return barStyle.Render(barContent)
}

// RenderSeparatorWithGlow renders a separator with Glow styling
func (gs *GlowStyler) RenderSeparatorWithGlow(char string) string {
	separator := strings.Repeat(char, 50)
	sepStyle := lipgloss.Style{}.
		Foreground(styles.GetAccent()).
		Bold(true)
	return sepStyle.Render(separator)
}
