package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) viewStats() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Statistics ", m.CurrentProject.Title))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Heading.Render("📊 Writing Statistics"))
	b.WriteString("\n\n")

	// Calculate total words
	totalWords := 0
	for _, scene := range m.CurrentProject.Scenes {
		totalWords += scene.WordCount
	}

	// Project stats
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Words:      %d\n", totalWords)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Target Words:     %d\n", m.CurrentProject.TargetWordCount)))

	progress := 0.0
	if m.CurrentProject.TargetWordCount > 0 {
		progress = float64(totalWords) / float64(m.CurrentProject.TargetWordCount) * 100
	}
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Progress:         %.1f%%\n", progress)))
	b.WriteString("\n")

	// Progress bar
	barWidth := 40
	filled := int(progress / 100 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
	b.WriteString(m.Styles.Accent.Render(bar))
	b.WriteString("\n\n")

	// Content stats
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Scenes:           %d\n", len(m.CurrentProject.Scenes))))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Characters:       %d\n", len(m.CurrentProject.Characters))))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Locations:        %d\n", len(m.CurrentProject.Locations))))
	b.WriteString("\n")

	// Session stats
	if m.CurrentProject.TotalSessions > 0 {
		b.WriteString(m.Styles.Heading.Render("Session History"))
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Sessions:   %d\n", m.CurrentProject.TotalSessions)))
		b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Current Streak:   %d days\n", m.CurrentProject.CurrentStreak)))

		hours := m.CurrentProject.TotalTimeSeconds / 3600
		minutes := (m.CurrentProject.TotalTimeSeconds % 3600) / 60
		b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Time:       %dh %dm\n", hours, minutes)))
	}

	b.WriteString("\n\n")
	b.WriteString(m.Styles.Text.Render("Press Esc to return"))

	return b.String()
}

func (m Model) handleStatsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentScreen = ScreenEditor
		return m, nil
	}
	return m, nil
}
