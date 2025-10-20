package glow

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	statusColor := "#FF71CE"
	switch status {
	case "completed":
		statusIcon = "✅"
		statusColor = "#00FF66"
	case "in-progress":
		statusIcon = "🔄"
		statusColor = "#FFD700"
	}

	// Priority indicator
	priorityIcon := "⚪"
	priorityColor := "#808080"
	switch priority {
	case "high":
		priorityIcon = "🔴"
		priorityColor = "#FF0040"
	case "medium":
		priorityIcon = "🟡"
		priorityColor = "#FFFF00"
	case "low":
		priorityIcon = "🟢"
		priorityColor = "#00FF00"
	}

	// Create enhanced task styling
	taskStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFF0")).
		Padding(0, 1).
		MarginBottom(1)

	content := fmt.Sprintf("%s %s %s | %s | %s",
		lipgloss.NewStyle().Foreground(lipgloss.Color(statusColor)).Render(statusIcon),
		lipgloss.NewStyle().Foreground(lipgloss.Color(priorityColor)).Render(priorityIcon),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF71CE")).Bold(true).Render(fmt.Sprintf("Task %d", index)),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFF0")).Render(id),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Italic(true).Render(description),
	)

	return taskStyle.Render(content)
}

// RenderHeaderWithGlow renders a header with enhanced Glow styling
func (gs *GlowStyler) RenderHeaderWithGlow(title, subtitle string) string {
	// Create enhanced header styling
	headerStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF71CE")).
		Padding(1, 2).
		MarginBottom(1)

	// Title styling
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF71CE")).
		Bold(true).
		Render(title)

	// Subtitle styling
	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFF0")).
		Render(subtitle)

	content := fmt.Sprintf("%s\n%s", titleStyle, subtitleStyle)
	return headerStyle.Render(content)
}

// RenderSectionWithGlow renders a section with Glow-enhanced styling
func (gs *GlowStyler) RenderSectionWithGlow(title string, content string, accentColor string) string {
	// Title bar styling
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(accentColor)).
		Bold(true).
		Render(fmt.Sprintf("─ %s ─", title))

	// Content styling
	contentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2D3748")). // Dark text for light theme
		Background(lipgloss.Color("#F7FAFC")). // Light panel background
		Padding(1, 2).
		Width(80).
		Render(content)

	sectionContent := fmt.Sprintf("%s\n%s\n%s", titleStyle, contentStyle, strings.Repeat("─", 30))
	return sectionContent
}

// RenderSuccessWithGlow renders success message with Glow styling
func (gs *GlowStyler) RenderSuccessWithGlow(message string) string {
	successStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FF66")).
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true).
		Padding(0, 2).
		MarginTop(1)

	return successStyle.Render(fmt.Sprintf("✅ %s", message))
}

// RenderErrorWithGlow renders error message with Glow styling
func (gs *GlowStyler) RenderErrorWithGlow(message string) string {
	errorStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF0040")).
		Foreground(lipgloss.Color("#FF0040")).
		Bold(true).
		Padding(0, 2).
		MarginTop(1)

	return errorStyle.Render(fmt.Sprintf("❌ %s", message))
}

// RenderInfoWithGlow renders info message with Glow styling
func (gs *GlowStyler) RenderInfoWithGlow(message string) string {
	infoStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#00FFF0")).
		Foreground(lipgloss.Color("#00FFF0")).
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
	progressBar.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Render(filled))

	// Empty portion
	empty := strings.Repeat("░", barWidth-filledWidth)
	progressBar.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CBD5E0")). // Light grey for light theme instead of #808080
		Render(empty))

	progressBar.WriteString("]")

	// Progress text
	progressText := fmt.Sprintf("%s %d/%d (%.1f%%)",
		lipgloss.NewStyle().Foreground(lipgloss.Color("#FF71CE")).Render(label),
		current, total, percentage*100)

	// Combine with styling
	barStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#00FFF0")).
		Padding(0, 1).
		MarginTop(1)

	barContent := fmt.Sprintf("%s\n%s", progressBar.String(), progressText)
	return barStyle.Render(barContent)
}

// RenderSeparatorWithGlow renders a separator with Glow styling
func (gs *GlowStyler) RenderSeparatorWithGlow(char string) string {
	separator := strings.Repeat(char, 50)
	sepStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFF0")).
		Bold(true)
	return sepStyle.Render(separator)
}
