package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
)

// renderTaskEntry renders the task entry modal
func (m *MainModel) renderTaskEntry() string {
	var b strings.Builder

	if m.taskInput == "PRIORITY_MODE" {
		// Priority change mode
		header := lipgloss.NewStyle().
			Foreground(synthYellow).
			Bold(true).
			Align(lipgloss.Center).
			Render("🎚️ CHANGE TASK PRIORITY")
		b.WriteString(header)
		b.WriteString("\n\n")
		b.WriteString(priorityInputStyle.Render("Select Priority:"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("  [1] Low 💤"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("  [2] Medium ⚡"))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthRed).Render("  [3] High 🔥"))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthYellow).Render("[Esc] Cancel"))
	} else {
		// Task entry mode
		header := lipgloss.NewStyle().
			Foreground(synthGreen).
			Bold(true).
			Align(lipgloss.Center).
			Render("➕ ADD NEW MISSION")
		b.WriteString(header)
		b.WriteString("\n\n")
		b.WriteString(taskInputStyle.Render("Mission: " + m.taskInput + "█"))
		b.WriteString("\n\n")
		b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("[Enter] Confirm  [Esc] Cancel"))
	}

	return b.String()
}

// renderNotesEditor renders the notes editing interface
func (m *MainModel) renderNotesEditor() string {
	var b strings.Builder

	if m.editingTask == nil {
		return "No task selected for notes"
	}

	// Notes editor header
	header := lipgloss.NewStyle().
		Foreground(synthCyan).
		Bold(true).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("📝 EDITING NOTES FOR: %s", m.editingTask.Description))
	b.WriteString(header)
	b.WriteString("\n\n")

	// Current notes display with Glow markdown rendering
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Bold(true).Render("📝 Current notes (Markdown):"))
	b.WriteString("\n")

	// Show notes using Glow markdown renderer
	currentNotes := m.editingTask.Notes
	if currentNotes == "" {
		currentNotes = "*No notes yet - start writing in markdown!*\n\n## Examples:\n- **Bold text**\n- *Italic text*\n- `Code snippets`\n- # Headers\n- [Links](url)"
	}

	// Use Glow for enhanced markdown rendering
	if m.glowStyler != nil {
		glowContent := m.glowStyler.RenderSectionWithGlow(
			"📝 Notes Content",
			currentNotes,
			"#00FFF0", // Cyan accent
		)
		b.WriteString(glowContent)
	} else {
		// Fallback to basic rendering
		notesStyle := lipgloss.NewStyle().
			Foreground(synthYellow).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Width(70).
			Render(currentNotes)
		b.WriteString(notesStyle)
	}
	b.WriteString("\n\n")

	// Input area
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("Enter new notes:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(synthCyan).
		Background(styles.GetBoxStyle().GetBackground()).
		Padding(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(synthGreen).
		Width(70).
		Height(6).
		Render(m.notesInput + "█"))
	b.WriteString("\n\n")

	// Instructions
	instructions := []string{
		"[Enter] Save notes",
		"[Esc] Cancel",
		"Type to add notes",
	}

	for _, instruction := range instructions {
		b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  " + instruction))
		b.WriteString("\n")
	}

	return b.String()
}

// renderFilterEditor renders the filter configuration interface
func (m *MainModel) renderFilterEditor() string {
	var b strings.Builder

	// Filter editor header
	header := lipgloss.NewStyle().
		Foreground(synthYellow).
		Bold(true).
		Align(lipgloss.Center).
		Render("🔍 TASK FILTER")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Current filter status
	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Current filter:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).
		Render(fmt.Sprintf("Status: %s | Priority: %s", m.filterStatus, m.filterPriority)))
	b.WriteString("\n\n")

	// Filter options
	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Status options:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [1] All tasks"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [2] Pending only"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [3] Completed only"))
	b.WriteString("\n\n")

	b.WriteString(lipgloss.NewStyle().Foreground(synthCyan).Render("Priority options:"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [H] High priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [M] Medium priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [L] Low priority"))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("  [A] All priorities"))
	b.WriteString("\n\n")

	// Instructions
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("[Enter] Apply filter  [Esc] Cancel"))

	return b.String()
}

// renderSettingsView renders the settings configuration interface
func (m *MainModel) renderSettingsView() string {
	var b strings.Builder

	// Settings header
	header := lipgloss.NewStyle().
		Foreground(synthPink).
		Bold(true).
		Align(lipgloss.Center).
		Render("⚙️ SYNTHWAVE SETTINGS")
	b.WriteString(header)
	b.WriteString("\n\n")

	// Current settings display
	currentTheme := styles.GetTheme().Name
	settings := []string{
		fmt.Sprintf("🎨 Theme: %s", currentTheme),
		fmt.Sprintf("🔊 Audio: %t", m.audioEnabled),
		fmt.Sprintf("⏰ Work Duration: %v", m.workDuration),
		fmt.Sprintf("☕ Break Duration: %v", m.breakDuration),
		fmt.Sprintf("📅 Calendar View: %s", m.calViewMode),
	}

	for _, setting := range settings {
		b.WriteString(lipgloss.NewStyle().
			Foreground(synthCyan).
			Background(styles.GetBoxStyle().GetBackground()).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(synthBlue).
			Render(setting))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Controls info
	controls := []string{
		"[T] Cycle theme",
		"[M] Toggle audio",
		"[W] Change work duration",
		"[B] Change break duration",
		"[Tab] Switch view",
	}

	b.WriteString(sectionTitleStyle.Render("⚡ SETTINGS CONTROLS:"))
	b.WriteString("\n")
	for _, control := range controls {
		b.WriteString(lipgloss.NewStyle().Foreground(synthGreen).Render("  " + control))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(synthBlue).Render("Press 'Tab' to return to dashboard"))

	return b.String()
}
