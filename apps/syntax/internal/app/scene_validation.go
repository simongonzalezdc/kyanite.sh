package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/scene"
)

func (m Model) viewSceneValidation() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	// Data is loaded in Update via ensureDataLoaded()

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(" Scene Compilation Report ")
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Heading.Render("📋 Scene Validation"))
	b.WriteString("\n\n")

	// Run compilation
	report := scene.CompileScenes(m.CurrentProject.Scenes)

	// Summary
	b.WriteString(m.Styles.Text.Bold(true).Render("Summary:"))
	b.WriteString("\n")
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Scenes: %d\n", report.TotalScenes)))
	b.WriteString(m.Styles.Text.Render(fmt.Sprintf("Total Chapters: %d\n", report.ChapterCount)))
	b.WriteString("\n")

	summaryStyle := m.Styles.Success
	if report.HasErrors() {
		summaryStyle = m.Styles.Error
	} else if report.HasWarnings() {
		summaryStyle = m.Styles.Text
	}
	b.WriteString(summaryStyle.Render(report.Summary()))
	b.WriteString("\n\n")

	// Issues
	if len(report.Issues) > 0 {
		b.WriteString(m.Styles.Text.Bold(true).Render("Issues:"))
		b.WriteString("\n\n")

		for i, issue := range report.Issues {
			if i >= 20 { // Limit display to 20 issues
				b.WriteString(m.Styles.Text.Faint(true).Render(fmt.Sprintf("... and %d more issues\n", len(report.Issues)-20)))
				break
			}

			var icon string
			var style = m.Styles.Text

			switch issue.Severity {
			case "error":
				icon = "✗"
				style = m.Styles.Error
			case "warning":
				icon = "⚠"
				style = m.Styles.Text
			case "info":
				icon = "ℹ"
				style = m.Styles.Text.Faint(true)
			}

			b.WriteString(style.Render(fmt.Sprintf("  %s %s\n", icon, issue.Message)))
		}
		b.WriteString("\n")
	}

	// Orphans
	if len(report.Orphans) > 0 {
		b.WriteString(m.Styles.Text.Bold(true).Render(fmt.Sprintf("Orphan Scenes (%d):", len(report.Orphans))))
		b.WriteString("\n")
		for _, sc := range report.Orphans {
			b.WriteString(m.Styles.Text.Render(fmt.Sprintf("  • %s (ID: %s)\n", sc.Name, sc.ID)))
		}
		b.WriteString("\n")
	}

	// Gaps
	if len(report.Gaps) > 0 {
		b.WriteString(m.Styles.Text.Bold(true).Render(fmt.Sprintf("Missing Scenes (%d):", len(report.Gaps))))
		b.WriteString("\n")
		for i, gap := range report.Gaps {
			if i >= 10 { // Limit display
				b.WriteString(m.Styles.Text.Faint(true).Render(fmt.Sprintf("  ... and %d more gaps\n", len(report.Gaps)-10)))
				break
			}
			b.WriteString(m.Styles.Text.Render(fmt.Sprintf("  • %s\n", gap.Message)))
		}
		b.WriteString("\n")
	}

	// Duplicates
	if len(report.Duplicates) > 0 {
		b.WriteString(m.Styles.Text.Bold(true).Render(fmt.Sprintf("Duplicate Scenes (%d):", len(report.Duplicates))))
		b.WriteString("\n")
		for _, dup := range report.Duplicates {
			b.WriteString(m.Styles.Error.Render(fmt.Sprintf("  • %s\n", dup.Message)))
			for _, id := range dup.SceneIDs {
				b.WriteString(m.Styles.Text.Faint(true).Render(fmt.Sprintf("    - %s\n", id)))
			}
		}
		b.WriteString("\n")
	}

	b.WriteString(m.Styles.Text.Render("Press Esc to go back"))

	return b.String()
}

func (m Model) handleSceneValidationKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.CurrentScreen = ScreenScenes
		return m, nil
	}

	return m, nil
}
