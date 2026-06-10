package app

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/backup"
)

func (m Model) viewBackups() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Backup Management ", m.CurrentProject.Title))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Heading.Render("💾 Backups"))
	b.WriteString("\n\n")

	// List backups
	backups, err := backup.ListBackups(m.CurrentProject.Directory)
	if err != nil {
		b.WriteString(m.Styles.Error.Render(fmt.Sprintf("Error loading backups: %v", err)))
	} else if len(backups) == 0 {
		b.WriteString(m.Styles.Text.Render("No backups found. Press 'c' to create a backup."))
	} else {
		for i, bkp := range backups {
			prefix := "  "
			style := m.Styles.Text

			if i == m.SelectedIndex {
				prefix = "> "
				style = m.Styles.Accent
			}

			// Format timestamp
			timeStr := bkp.Timestamp.Format("2006-01-02 15:04:05")
			age := time.Since(bkp.Timestamp)
			ageStr := ""
			if age < time.Hour {
				ageStr = fmt.Sprintf("%d minutes ago", int(age.Minutes()))
			} else if age < 24*time.Hour {
				ageStr = fmt.Sprintf("%d hours ago", int(age.Hours()))
			} else {
				ageStr = fmt.Sprintf("%d days ago", int(age.Hours()/24))
			}

			// Format size
			sizeStr := formatSize(bkp.Size)

			line := fmt.Sprintf("%s%s (%s) - %s", prefix, timeStr, ageStr, sizeStr)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.Styles.Text.Render("c - Create Backup  |  r - Restore  |  d - Delete  |  Esc - Back"))

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

func (m Model) handleBackupsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	backups, err := backup.ListBackups(m.CurrentProject.Directory)
	if err != nil {
		m.Error = err
		return m, nil
	}

	backupCount := len(backups)

	switch msg.String() {
	case "up", "k":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}
		return m, nil

	case "down", "j":
		if backupCount > 0 && m.SelectedIndex < backupCount-1 {
			m.SelectedIndex++
		}
		return m, nil

	case "c":
		// Create backup
		m.Message = "Creating backup..."
		if err := backup.CreateBackup(m.CurrentProject.Directory); err != nil {
			m.Error = err
			m.Message = ""
		} else {
			m.Message = "Backup created successfully"
		}
		return m, nil

	case "r":
		// Restore from selected backup
		if backupCount > 0 && m.SelectedIndex < backupCount {
			bkp := backups[m.SelectedIndex]
			m.Message = "Restoring backup..."
			if err := backup.RestoreBackup(bkp.Path, m.CurrentProject.Directory); err != nil {
				m.Error = err
				m.Message = ""
			} else {
				m.Message = "Backup restored successfully. Restart to see changes."
			}
		}
		return m, nil

	case "d":
		// Delete selected backup
		if backupCount > 0 && m.SelectedIndex < backupCount {
			bkp := backups[m.SelectedIndex]
			if err := backup.DeleteBackup(bkp.Path); err != nil {
				m.Error = err
			} else {
				m.Message = "Backup deleted"
				// Adjust selected index if needed
				if m.SelectedIndex >= backupCount-1 {
					m.SelectedIndex = max(0, backupCount-2)
				}
			}
		}
		return m, nil

	case "esc":
		m.SelectedIndex = 0
		m.Message = ""
		m.Error = nil
		m.CurrentScreen = ScreenEditor
		return m, nil
	}

	return m, nil
}

// formatSize formats bytes into human-readable format
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
