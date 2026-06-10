package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/noise/internal/infra/sync"
	"github.com/kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// IdeaInboxModel manages the inbox of captured ideas from PWA
type IdeaInboxModel struct {
	ideas        []*sync.CapturedIdea
	selectedIdx  int
	width        int
	height       int
	scrollOffset int

	// Styles
	titleStyle    lipgloss.Style
	itemStyle     lipgloss.Style
	selectedStyle lipgloss.Style
	typeStyle     lipgloss.Style
	timeStyle     lipgloss.Style
	contentStyle  lipgloss.Style
	emptyStyle    lipgloss.Style
	helpStyle     lipgloss.Style
}

// NewIdeaInboxModel creates a new idea inbox model
func NewIdeaInboxModel() *IdeaInboxModel {
	t := theme.GetManager().Current()

	return &IdeaInboxModel{
		ideas: make([]*sync.CapturedIdea, 0),
		titleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(t.Primary).
			MarginBottom(1),
		itemStyle: lipgloss.NewStyle().
			Foreground(t.Text).
			PaddingLeft(2).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(t.Secondary),
		selectedStyle: lipgloss.NewStyle().
			Foreground(t.Background).
			Background(t.Primary).
			Bold(true).
			PaddingLeft(2),
		typeStyle: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),
		timeStyle: lipgloss.NewStyle().
			Foreground(t.Secondary),
		contentStyle: lipgloss.NewStyle().
			Foreground(t.Text),
		emptyStyle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true),
		helpStyle: lipgloss.NewStyle().
			Foreground(t.Secondary),
	}
}

// Init initializes the idea inbox model
func (m *IdeaInboxModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *IdeaInboxModel) Update(msg tea.Msg) (*IdeaInboxModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.updateScroll()
			}
		case "down", "j":
			if m.selectedIdx < len(m.ideas)-1 {
				m.selectedIdx++
				m.updateScroll()
			}
		case "enter":
			return m, m.assignIdea()
		case "delete":
			return m, m.deleteIdea()
		case "ctrl+p":
			return m, m.previewIdea()
		case "esc":
			return m, func() tea.Msg { return BackToSettingsMsg{} }
		}

	case IdeaReceivedMsg:
		m.ideas = append(m.ideas, msg.Idea)
	}

	return m, nil
}

// View renders the idea inbox
func (m *IdeaInboxModel) View() string {
	var b strings.Builder

	// Title
	b.WriteString(m.titleStyle.Render("[*] Idea Inbox"))
	b.WriteString(fmt.Sprintf(" (%d ideas)\n\n", len(m.ideas)))

	// Empty state
	if len(m.ideas) == 0 {
		b.WriteString(m.emptyStyle.Render("No ideas captured yet.\n"))
		b.WriteString(m.emptyStyle.Render("Use the PWA companion to capture ideas on the go."))
		b.WriteString("\n\n")
		b.WriteString(m.helpStyle.Render("Esc: Back"))
		return b.String()
	}

	// Calculate visible items
	visibleHeight := m.height - 8 // Account for title and help
	itemHeight := 3               // Lines per item
	visibleItems := 0
	if itemHeight > 0 {
		visibleItems = visibleHeight / itemHeight
	}
	if visibleItems < 1 {
		visibleItems = 1
	}

	// Render visible items
	start := m.scrollOffset
	end := start + visibleItems
	if end > len(m.ideas) {
		end = len(m.ideas)
	}

	for i := start; i < end; i++ {
		idea := m.ideas[i]
		b.WriteString(m.renderIdea(idea, i == m.selectedIdx))
		b.WriteString("\n")
	}

	// Scroll indicator
	if len(m.ideas) > visibleItems {
		scrollInfo := fmt.Sprintf("(%d-%d of %d)", start+1, end, len(m.ideas))
		b.WriteString(m.timeStyle.Render(scrollInfo))
		b.WriteString("\n")
	}

	// Help
	b.WriteString("\n")
	b.WriteString(m.helpStyle.Render("Up/Down: Navigate - Enter: Assign to song - Delete: Remove - Ctrl+P: Preview - Esc: Back"))

	return b.String()
}

// renderIdea renders a single idea item
func (m *IdeaInboxModel) renderIdea(idea *sync.CapturedIdea, selected bool) string {
	var style lipgloss.Style
	if selected {
		style = m.selectedStyle
	} else {
		style = m.itemStyle
	}

	// Type icon
	icon := m.getTypeIcon(idea.Type)

	// Time
	timeStr := formatRelativeTime(idea.CreatedAt)

	// Content preview
	content := idea.Content
	if len(content) > 60 {
		content = content[:57] + "..."
	}
	if idea.Type == sync.IdeaTypeTempo {
		content = fmt.Sprintf("%d BPM", idea.BPM)
	}
	if idea.Type == sync.IdeaTypeVoiceMemo {
		content = "[MIC] Voice memo"
	}
	if idea.Type == sync.IdeaTypePhoto {
		content = "[IMG] Photo"
	}

	// Format
	line1 := fmt.Sprintf("%s %s  %s", icon, m.typeStyle.Render(string(idea.Type)), m.timeStyle.Render(timeStr))
	line2 := m.contentStyle.Render(content)

	return style.Render(line1 + "\n" + line2)
}

// getTypeIcon returns an icon for the idea type
func (m *IdeaInboxModel) getTypeIcon(ideaType sync.IdeaType) string {
	switch ideaType {
	case sync.IdeaTypeText:
		return "[TXT]"
	case sync.IdeaTypeVoiceMemo:
		return "[MIC]"
	case sync.IdeaTypePhoto:
		return "[IMG]"
	case sync.IdeaTypeTempo:
		return "[BPM]"
	default:
		return "[*]"
	}
}

// updateScroll adjusts scroll offset to keep selected item visible
func (m *IdeaInboxModel) updateScroll() {
	visibleHeight := m.height - 8
	itemHeight := 3
	visibleItems := 0
	if itemHeight > 0 {
		visibleItems = visibleHeight / itemHeight
	}
	if visibleItems < 1 {
		visibleItems = 1
	}

	if m.selectedIdx < m.scrollOffset {
		m.scrollOffset = m.selectedIdx
	} else if m.selectedIdx >= m.scrollOffset+visibleItems {
		m.scrollOffset = m.selectedIdx - visibleItems + 1
	}
}

// assignIdea assigns the selected idea to a song
func (m *IdeaInboxModel) assignIdea() tea.Cmd {
	if m.selectedIdx >= len(m.ideas) {
		return nil
	}
	// TODO: Open song picker to assign idea
	return nil
}

// deleteIdea removes the selected idea
func (m *IdeaInboxModel) deleteIdea() tea.Cmd {
	if m.selectedIdx >= len(m.ideas) {
		return nil
	}

	// Remove from list
	m.ideas = append(m.ideas[:m.selectedIdx], m.ideas[m.selectedIdx+1:]...)

	// Adjust selection
	if m.selectedIdx >= len(m.ideas) && m.selectedIdx > 0 {
		m.selectedIdx--
	}

	return nil
}

// previewIdea shows a preview of the selected idea
func (m *IdeaInboxModel) previewIdea() tea.Cmd {
	if m.selectedIdx >= len(m.ideas) {
		return nil
	}
	// TODO: Open preview modal
	return nil
}

// AddIdea adds a new idea to the inbox
func (m *IdeaInboxModel) AddIdea(idea *sync.CapturedIdea) {
	m.ideas = append(m.ideas, idea)
}

// GetIdeas returns all ideas
func (m *IdeaInboxModel) GetIdeas() []*sync.CapturedIdea {
	return m.ideas
}

// formatRelativeTime formats a time relative to now
func formatRelativeTime(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "yesterday"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		return t.Format("Jan 2")
	}
}

// IdeaReceivedMsg is sent when a new idea is received
type IdeaReceivedMsg struct {
	Idea *sync.CapturedIdea
}
