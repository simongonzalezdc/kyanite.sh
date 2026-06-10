package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
)

func (m *MainModel) renderHelpBar() string {
	// CONSISTENT: Help bar styling
	helpStyle := lipgloss.NewStyle().
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).                         // CONSISTENT: 1 vertical, 2 horizontal
		BorderStyle(lipgloss.RoundedBorder()). // CONSISTENT: RoundedBorder
		BorderForeground(styles.GetBorder()).
		Width(m.width - 4)

	// Get shortcuts for current view
	var shortcuts []string

	switch m.currentView {
	case dashboardView:
		shortcuts = []string{
			"[↑/↓] Navigate",
			"[A] Add Task",
			"[D] Complete",
			"[N] Notes",
			"[G] Journal",
			"[K] Calendar",
			"[C] Chat",
			"[T] Theme",
			"[S] Settings",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case focusView:
		shortcuts = []string{
			"[Space] Switch Work/Break",
			"[P] Change Priority",
			"[ESC] Return",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case chatView:
		shortcuts = []string{
			"[↑/↓] Navigate History",
			"[Enter] Send Message",
			"[ESC] Return",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case calendarView:
		shortcuts = []string{
			"[←/→] Previous/Next Month",
			"[↑/↓] Navigate Days",
			"[Enter] Select Date",
			"[ESC] Return",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case notesView:
		shortcuts = []string{
			"[↑/↓] Navigate",
			"[Enter] Edit",
			"[Escape] Save & Exit",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case journalView:
		shortcuts = []string{
			"[↑/↓] Navigate Commands",
			"[Enter] Select Command",
			"[Escape] Return",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	case settingsView:
		shortcuts = []string{
			"[↑/↓] Navigate Options",
			"[Enter] Select",
			"[Escape] Save & Exit",
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	default:
		shortcuts = []string{
			"[Tab] Next View",
			"[Ctrl+Q] Quit",
		}
	}

	// Format shortcuts with proper spacing
	var shortcutStrings []string
	for _, shortcut := range shortcuts {
		shortcutStrings = append(shortcutStrings,
			lipgloss.NewStyle().
				Foreground(styles.GetForeground()).
				Render(shortcut))
	}

	helpContent := strings.Join(shortcutStrings, " • ")
	return helpStyle.Render(helpContent)
}
