package collaboration

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/noise/internal/ui/dimension"
	"github.com/charmbracelet/lipgloss"
)

// CollaborationStatusBar displays collaboration status in the status bar
type CollaborationStatusBar struct {
	// Collaboration state
	sessionID        string
	isCollaborating  bool
	participantCount int
	currentUserRole  string
	hasConflicts     bool

	// UI state
	width       int
	height      int
	showDetails bool

	// Styles
	containerStyle lipgloss.Style
	activeStyle    lipgloss.Style
	inactiveStyle  lipgloss.Style
	conflictStyle  lipgloss.Style
	detailsStyle   lipgloss.Style
}

// NewCollaborationStatusBar creates a new collaboration status bar
func NewCollaborationStatusBar() *CollaborationStatusBar {
	return &CollaborationStatusBar{
		isCollaborating:  false,
		participantCount: 0,
		showDetails:      false,
		containerStyle: lipgloss.NewStyle().
			Padding(0, 1).
			Margin(0, 1),
		activeStyle: lipgloss.NewStyle().
			Bold(true),
		inactiveStyle: lipgloss.NewStyle(),
		conflictStyle: lipgloss.NewStyle().
			Bold(true),
		detailsStyle: lipgloss.NewStyle().
			Italic(true),
	}
}

// SetDimensions sets the dimensions for the status bar
func (csb *CollaborationStatusBar) SetDimensions(width, height int) {
	dimension.Set(&csb.width, &csb.height, width, height)
}

func (csb *CollaborationStatusBar) GetDimensions() (int, int) {
	return csb.width, csb.height
}

// UpdateCollaborationState updates the collaboration state
func (csb *CollaborationStatusBar) UpdateCollaborationState(sessionID string, isCollaborating bool, participantCount int, userRole string, hasConflicts bool) {
	csb.sessionID = sessionID
	csb.isCollaborating = isCollaborating
	csb.participantCount = participantCount
	csb.currentUserRole = userRole
	csb.hasConflicts = hasConflicts
}

// ToggleDetails toggles the display of detailed information
func (csb *CollaborationStatusBar) ToggleDetails() {
	csb.showDetails = !csb.showDetails
}

// View renders the collaboration status bar
func (csb *CollaborationStatusBar) View() string {
	if !csb.isCollaborating {
		return csb.inactiveStyle.Render("No active collaboration")
	}

	var content strings.Builder

	// Main status
	if csb.hasConflicts {
		content.WriteString(csb.conflictStyle.Render("[!] "))
	} else {
		content.WriteString(csb.activeStyle.Render("¤\u009d "))
	}

	// Session info
	if csb.sessionID != "" {
		sessionShort := csb.sessionID
		if len(sessionShort) > 8 {
			sessionShort = sessionShort[:8] + "..."
		}
		content.WriteString(fmt.Sprintf("Session: %s", sessionShort))
	}

	// Participant count
	if csb.participantCount > 0 {
		content.WriteString(fmt.Sprintf(" (%d users)", csb.participantCount))
	}

	// User role
	if csb.currentUserRole != "" {
		content.WriteString(fmt.Sprintf(" [%s]", csb.currentUserRole))
	}

	// Add details if requested
	if csb.showDetails {
		content.WriteString("\n")
		content.WriteString(csb.renderDetails())
	}

	return csb.containerStyle.Render(content.String())
}

// renderDetails renders detailed collaboration information
func (csb *CollaborationStatusBar) renderDetails() string {
	var details []string

	if csb.sessionID != "" {
		details = append(details, fmt.Sprintf("Session ID: %s", csb.sessionID))
	}

	if csb.participantCount > 0 {
		details = append(details, fmt.Sprintf("Participants: %d", csb.participantCount))
	}

	if csb.currentUserRole != "" {
		details = append(details, fmt.Sprintf("Your role: %s", csb.currentUserRole))
	}

	if csb.hasConflicts {
		details = append(details, "Conflicts detected - manual resolution may be needed")
	} else {
		details = append(details, "No conflicts")
	}

	details = append(details, fmt.Sprintf("Last updated: %s", time.Now().Format("15:04:05")))

	return csb.detailsStyle.Render(strings.Join(details, " - "))
}

// GetStatusText returns a short status text for compact display
func (csb *CollaborationStatusBar) GetStatusText() string {
	if !csb.isCollaborating {
		return "Solo"
	}

	if csb.hasConflicts {
		return fmt.Sprintf("Conflict (%d)", csb.participantCount)
	}

	return fmt.Sprintf("Collab (%d)", csb.participantCount)
}

// IsActive returns whether collaboration is currently active
func (csb *CollaborationStatusBar) IsActive() bool {
	return csb.isCollaborating
}

// HasConflicts returns whether there are active conflicts
func (csb *CollaborationStatusBar) HasConflicts() bool {
	return csb.hasConflicts
}
