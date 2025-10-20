package collaboration

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/Kyanite/noise/internal/ui/styles"
)

// PresenceIndicatorModel handles displaying user presence indicators
type PresenceIndicatorModel struct {
	// Presence data
	indicators []SessionPresenceIndicator
	width      int
	height     int

	// UI state
	showDetails bool
	selected    int

	// Styles
	containerStyle lipgloss.Style
	indicatorStyle lipgloss.Style
	selectedStyle  lipgloss.Style
	detailsStyle   lipgloss.Style
}

// SessionPresenceIndicator represents presence info for UI display
type SessionPresenceIndicator struct {
	UserID    string                          `json:"user_id"`
	Username  string                          `json:"username"`
	Indicator collaboration.PresenceIndicator `json:"indicator"`
	Cursor    collaboration.CursorPosition    `json:"cursor"`
	Role      collaboration.ParticipantRole   `json:"role"`
}

// NewPresenceIndicatorModel creates a new presence indicator model
func NewPresenceIndicatorModel() *PresenceIndicatorModel {
	return &PresenceIndicatorModel{
		indicators:  make([]SessionPresenceIndicator, 0),
		showDetails: false,
		selected:    0,
		containerStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1),
		indicatorStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1),
		selectedStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1).
			Background(styles.Accent).
			Foreground(styles.Background),
		detailsStyle: lipgloss.NewStyle().
			Foreground(styles.TextMuted).
			Italic(true),
	}
}

// SetDimensions sets the dimensions for the presence indicator
func (m *PresenceIndicatorModel) SetDimensions(width, height int) {
	dimension.Set(&m.width, &m.height, width, height)
}

func (m *PresenceIndicatorModel) GetDimensions() (int, int) {
	return m.width, m.height
}

// UpdateIndicators updates the presence indicators
func (m *PresenceIndicatorModel) UpdateIndicators(indicators []SessionPresenceIndicator) {
	m.indicators = indicators
}

// View renders the presence indicator
func (m *PresenceIndicatorModel) View() string {
	if len(m.indicators) == 0 {
		return ""
	}

	var content strings.Builder

	// Header
	title := "ðŸ‘¥ Collaborators"
	if len(m.indicators) == 1 {
		title += " (1 user)"
	} else {
		title += fmt.Sprintf(" (%d users)", len(m.indicators))
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		MarginBottom(1)

	content.WriteString(titleStyle.Render(title))
	content.WriteString("\n")

	// Show indicators
	for i, indicator := range m.indicators {
		m.renderIndicator(&content, indicator, i)
	}

	// Apply container styling
	return m.containerStyle.Width(m.width).Render(content.String())
}

// renderIndicator renders a single presence indicator
func (m *PresenceIndicatorModel) renderIndicator(content *strings.Builder, indicator SessionPresenceIndicator, index int) {
	var style lipgloss.Style
	if m.showDetails && index == m.selected {
		style = m.selectedStyle
	} else {
		style = m.indicatorStyle
	}

	// Create indicator symbol with color
	icon := m.getColoredIcon(indicator.Indicator)

	// Format user info
	userInfo := fmt.Sprintf("%s %s", icon, indicator.Username)

	// Add role indicator
	roleIcon := m.getRoleIcon(indicator.Role)
	if roleIcon != "" {
		userInfo += fmt.Sprintf(" %s", roleIcon)
	}

	// Add cursor position if available
	if indicator.Cursor.Line > 0 || indicator.Cursor.Column > 0 {
		userInfo += fmt.Sprintf(" (Ln %d, Col %d)", indicator.Cursor.Line+1, indicator.Cursor.Column+1)
	}

	// Render the indicator
	rendered := style.Render(userInfo)
	content.WriteString(rendered)
	content.WriteString(" ")

	// Add details if in details mode
	if m.showDetails {
		details := m.getUserDetails(indicator)
		if details != "" {
			content.WriteString("\n")
			content.WriteString(m.detailsStyle.Render("  " + details))
		}
	}

	content.WriteString("\n")
}

// getColoredIcon returns a colored icon for the presence status
func (m *PresenceIndicatorModel) getColoredIcon(indicator collaboration.PresenceIndicator) string {
	var color lipgloss.Color

	switch indicator.Color {
	case "green":
		color = lipgloss.Color("#00ff00")
	case "yellow":
		color = lipgloss.Color("#ffff00")
	case "red":
		color = lipgloss.Color("#ff0000")
	case "gray":
		color = lipgloss.Color("#666666")
	default:
		color = lipgloss.Color("#666666")
	}

	iconStyle := lipgloss.NewStyle().Foreground(color)
	return iconStyle.Render(indicator.Icon)
}

// getRoleIcon returns an icon representing the user's role
func (m *PresenceIndicatorModel) getRoleIcon(role collaboration.ParticipantRole) string {
	switch role {
	case collaboration.RoleOwner:
		return "ðŸ‘‘"
	case collaboration.RoleEditor:
		return "âœï¸"
	case collaboration.RoleViewer:
		return "ðŸ‘ï¸"
	default:
		return ""
	}
}

// getUserDetails returns detailed information about a user
func (m *PresenceIndicatorModel) getUserDetails(indicator SessionPresenceIndicator) string {
	var details []string

	// Add status details
	switch indicator.Indicator.Status {
	case collaboration.StatusOnline:
		details = append(details, "Currently editing")
	case collaboration.StatusAway:
		details = append(details, "Away")
	case collaboration.StatusBusy:
		details = append(details, "Busy")
	case collaboration.StatusOffline:
		details = append(details, "Offline")
	}

	// Add role details
	switch indicator.Role {
	case collaboration.RoleOwner:
		details = append(details, "Session owner")
	case collaboration.RoleEditor:
		details = append(details, "Can edit")
	case collaboration.RoleViewer:
		details = append(details, "Read-only")
	}

	return strings.Join(details, " â€¢ ")
}

// ToggleDetails toggles the display of detailed information
func (m *PresenceIndicatorModel) ToggleDetails() {
	m.showDetails = !m.showDetails
}

// SelectNext cycles to the next indicator
func (m *PresenceIndicatorModel) SelectNext() {
	if len(m.indicators) > 0 {
		m.selected = (m.selected + 1) % len(m.indicators)
	}
}

// SelectPrev cycles to the previous indicator
func (m *PresenceIndicatorModel) SelectPrev() {
	if len(m.indicators) > 0 {
		m.selected = (m.selected - 1 + len(m.indicators)) % len(m.indicators)
	}
}

// GetSelectedUser returns the currently selected user
func (m *PresenceIndicatorModel) GetSelectedUser() *SessionPresenceIndicator {
	if len(m.indicators) == 0 || m.selected >= len(m.indicators) {
		return nil
	}
	return &m.indicators[m.selected]
}
