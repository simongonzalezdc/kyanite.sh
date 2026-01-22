//go:build collaboration

// Collaboration UI tests are only run when the collaboration build tag is specified.
// Run with: go test -tags collaboration ./internal/ui/collaboration/...

package collaboration

import (
	"testing"
)

func TestCollaborationStatusBar_NewCollaborationStatusBar(t *testing.T) {
	csb := NewCollaborationStatusBar()

	if csb == nil {
		t.Fatal("Expected non-nil status bar")
	}

	if csb.isCollaborating != false {
		t.Errorf("Expected isCollaborating to be false, got %v", csb.isCollaborating)
	}

	if csb.participantCount != 0 {
		t.Errorf("Expected participantCount to be 0, got %d", csb.participantCount)
	}

	if csb.showDetails != false {
		t.Errorf("Expected showDetails to be false, got %v", csb.showDetails)
	}
}

func TestCollaborationStatusBar_SetDimensions(t *testing.T) {
	csb := NewCollaborationStatusBar()

	csb.SetDimensions(100, 5)

	width, height := csb.GetDimensions()
	if width != 100 {
		t.Errorf("Expected width 100, got %d", width)
	}
	if height != 5 {
		t.Errorf("Expected height 5, got %d", height)
	}
}

func TestCollaborationStatusBar_UpdateCollaborationState(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Update to active collaboration state
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		false,
	)

	if csb.sessionID != "session123" {
		t.Errorf("Expected sessionID to be session123, got %s", csb.sessionID)
	}

	if csb.isCollaborating != true {
		t.Errorf("Expected isCollaborating to be true, got %v", csb.isCollaborating)
	}

	if csb.participantCount != 3 {
		t.Errorf("Expected participantCount to be 3, got %d", csb.participantCount)
	}

	if csb.currentUserRole != "editor" {
		t.Errorf("Expected currentUserRole to be editor, got %s", csb.currentUserRole)
	}

	if csb.hasConflicts != false {
		t.Errorf("Expected hasConflicts to be false, got %v", csb.hasConflicts)
	}

	// Update to conflict state
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		true,
	)

	if csb.hasConflicts != true {
		t.Errorf("Expected hasConflicts to be true, got %v", csb.hasConflicts)
	}

	// Update to inactive state
	csb.UpdateCollaborationState(
		"",
		false,
		0,
		"",
		false,
	)

	if csb.sessionID != "" {
		t.Errorf("Expected sessionID to be empty, got %s", csb.sessionID)
	}

	if csb.isCollaborating != false {
		t.Errorf("Expected isCollaborating to be false, got %v", csb.isCollaborating)
	}

	if csb.participantCount != 0 {
		t.Errorf("Expected participantCount to be 0, got %d", csb.participantCount)
	}
}

func TestCollaborationStatusBar_ToggleDetails(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Initially details should be hidden
	if csb.showDetails != false {
		t.Errorf("Expected showDetails to be false initially, got %v", csb.showDetails)
	}

	// Toggle to show details
	csb.ToggleDetails()
	if csb.showDetails != true {
		t.Errorf("Expected showDetails to be true after toggle, got %v", csb.showDetails)
	}

	// Toggle to hide details
	csb.ToggleDetails()
	if csb.showDetails != false {
		t.Errorf("Expected showDetails to be false after second toggle, got %v", csb.showDetails)
	}
}

func TestCollaborationStatusBar_View_Inactive(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 5)

	// Test view when not collaborating
	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	if !contains(view, "No active collaboration") {
		t.Errorf("Expected view to contain 'No active collaboration', got %q", view)
	}
}

func TestCollaborationStatusBar_View_Active(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 5)

	// Update to active collaboration state
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		false,
	)

	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain collaboration icon
	if !contains(view, "🤝") {
		t.Errorf("Expected view to contain collaboration icon, got %q", view)
	}

	// Should contain session info
	if !contains(view, "Session:") {
		t.Errorf("Expected view to contain 'Session:', got %q", view)
	}

	// Should contain participant count
	if !contains(view, "(3 users)") {
		t.Errorf("Expected view to contain '(3 users)', got %q", view)
	}

	// Should contain user role
	if !contains(view, "[editor]") {
		t.Errorf("Expected view to contain '[editor]', got %q", view)
	}

	// Should not contain conflict warning
	if contains(view, "⚠️") {
		t.Errorf("Expected view to not contain conflict warning, got %q", view)
	}
}

func TestCollaborationStatusBar_View_WithConflicts(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 5)

	// Update to active collaboration state with conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		true,
	)

	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain conflict warning
	if !contains(view, "⚠️") {
		t.Errorf("Expected view to contain conflict warning, got %q", view)
	}

	// Should not contain normal collaboration icon
	if contains(view, "🤝") {
		t.Errorf("Expected view to not contain normal collaboration icon when conflicts exist, got %q", view)
	}
}

func TestCollaborationStatusBar_View_WithDetails(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 10)

	// Update to active collaboration state
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		false,
	)

	// Enable details
	csb.ToggleDetails()

	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain detailed information
	if !contains(view, "Session ID: session123") {
		t.Errorf("Expected view to contain session ID, got %q", view)
	}

	if !contains(view, "Participants: 3") {
		t.Errorf("Expected view to contain participant count, got %q", view)
	}

	if !contains(view, "Your role: editor") {
		t.Errorf("Expected view to contain user role, got %q", view)
	}

	if !contains(view, "No conflicts") {
		t.Errorf("Expected view to contain 'No conflicts', got %q", view)
	}

	if !contains(view, "Last updated:") {
		t.Errorf("Expected view to contain last updated time, got %q", view)
	}
}

func TestCollaborationStatusBar_View_WithDetailsAndConflicts(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 10)

	// Update to active collaboration state with conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		true,
	)

	// Enable details
	csb.ToggleDetails()

	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain conflict information
	if !contains(view, "Conflicts detected - manual resolution may be needed") {
		t.Errorf("Expected view to contain conflict message, got %q", view)
	}

	// Should not contain "No conflicts"
	if contains(view, "No conflicts") {
		t.Errorf("Expected view to not contain 'No conflicts' when conflicts exist, got %q", view)
	}
}

func TestCollaborationStatusBar_View_LongSessionID(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 5)

	// Update with very long session ID
	longSessionID := "session_very_long_id_that_should_be_truncated_in_display"
	csb.UpdateCollaborationState(
		longSessionID,
		true,
		2,
		"viewer",
		false,
	)

	view := csb.View()

	if view == "" {
		t.Error("Expected non-empty view")
	}

	// Should contain truncated session ID
	if !contains(view, "session_...") {
		t.Errorf("Expected view to contain truncated session ID, got %q", view)
	}

	// Should not contain full session ID
	if contains(view, longSessionID) {
		t.Errorf("Expected view to not contain full session ID, got %q", view)
	}
}

func TestCollaborationStatusBar_GetStatusText(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Test when not collaborating
	status := csb.GetStatusText()
	if status != "Solo" {
		t.Errorf("Expected status 'Solo', got %q", status)
	}

	// Test when collaborating without conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		false,
	)

	status = csb.GetStatusText()
	if status != "Collab (3)" {
		t.Errorf("Expected status 'Collab (3)', got %q", status)
	}

	// Test when collaborating with conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		true,
	)

	status = csb.GetStatusText()
	if status != "Conflict (3)" {
		t.Errorf("Expected status 'Conflict (3)', got %q", status)
	}

	// Test with single participant
	csb.UpdateCollaborationState(
		"session123",
		true,
		1,
		"owner",
		false,
	)

	status = csb.GetStatusText()
	if status != "Collab (1)" {
		t.Errorf("Expected status 'Collab (1)', got %q", status)
	}
}

func TestCollaborationStatusBar_IsActive(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Initially should be inactive
	if csb.IsActive() != false {
		t.Errorf("Expected IsActive to return false initially, got %v", csb.IsActive())
	}

	// Update to active state
	csb.UpdateCollaborationState(
		"session123",
		true,
		2,
		"editor",
		false,
	)

	if csb.IsActive() != true {
		t.Errorf("Expected IsActive to return true after update, got %v", csb.IsActive())
	}

	// Update to inactive state
	csb.UpdateCollaborationState(
		"",
		false,
		0,
		"",
		false,
	)

	if csb.IsActive() != false {
		t.Errorf("Expected IsActive to return false after deactivation, got %v", csb.IsActive())
	}
}

func TestCollaborationStatusBar_HasConflicts(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Initially should not have conflicts
	if csb.HasConflicts() != false {
		t.Errorf("Expected HasConflicts to return false initially, got %v", csb.HasConflicts())
	}

	// Update to state with conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		2,
		"editor",
		true,
	)

	if csb.HasConflicts() != true {
		t.Errorf("Expected HasConflicts to return true after update, got %v", csb.HasConflicts())
	}

	// Update to state without conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		2,
		"editor",
		false,
	)

	if csb.HasConflicts() != false {
		t.Errorf("Expected HasConflicts to return false after resolving conflicts, got %v", csb.HasConflicts())
	}
}

func TestCollaborationStatusBar_DifferentRoles(t *testing.T) {
	csb := NewCollaborationStatusBar()
	csb.SetDimensions(80, 5)

	roles := []string{"owner", "editor", "viewer", "moderator"}

	for _, role := range roles {
		csb.UpdateCollaborationState(
			"session123",
			true,
			2,
			role,
			false,
		)

		view := csb.View()

		if view == "" {
			t.Errorf("Expected non-empty view for role %s", role)
		}

		// Should contain role in brackets
		expectedRoleText := "[" + role + "]"
		if !contains(view, expectedRoleText) {
			t.Errorf("Expected view to contain %q for role %s, got %q", expectedRoleText, role, view)
		}
	}
}

func TestCollaborationStatusBar_EdgeCases(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Test with empty session ID but active collaboration
	csb.UpdateCollaborationState(
		"",
		true,
		1,
		"editor",
		false,
	)

	view := csb.View()
	if view == "" {
		t.Error("Expected non-empty view with empty session ID")
	}

	// Test with zero participants but active collaboration
	csb.UpdateCollaborationState(
		"session123",
		true,
		0,
		"owner",
		false,
	)

	view = csb.View()
	if view == "" {
		t.Error("Expected non-empty view with zero participants")
	}

	// Should show "(0 users)"
	if !contains(view, "(0 users)") {
		t.Errorf("Expected view to contain '(0 users)', got %q", view)
	}

	// Test with many participants
	csb.UpdateCollaborationState(
		"session123",
		true,
		999,
		"viewer",
		false,
	)

	view = csb.View()
	if view == "" {
		t.Error("Expected non-empty view with many participants")
	}

	// Should show "(999 users)"
	if !contains(view, "(999 users)") {
		t.Errorf("Expected view to contain '(999 users)', got %q", view)
	}
}

func TestCollaborationStatusBar_StateTransitions(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Test transitions between different states
	states := []struct {
		sessionID    string
		active       bool
		participants int
		role         string
		conflicts    bool
	}{
		{"", false, 0, "", false},             // Inactive
		{"session1", true, 1, "owner", false}, // Single user
		{"session1", true, 2, "owner", false}, // Multiple users
		{"session1", true, 2, "owner", true},  // With conflicts
		{"session1", true, 3, "owner", false}, // Conflicts resolved
		{"", false, 0, "", false},             // Back to inactive
	}

	for i, state := range states {
		csb.UpdateCollaborationState(
			state.sessionID,
			state.active,
			state.participants,
			state.role,
			state.conflicts,
		)

		// Verify state
		if csb.isCollaborating != state.active {
			t.Errorf("State %d: Expected isCollaborating to be %v, got %v", i, state.active, csb.isCollaborating)
		}

		if csb.participantCount != state.participants {
			t.Errorf("State %d: Expected participantCount to be %d, got %d", i, state.participants, csb.participantCount)
		}

		if csb.currentUserRole != state.role {
			t.Errorf("State %d: Expected currentUserRole to be %q, got %q", i, state.role, csb.currentUserRole)
		}

		if csb.hasConflicts != state.conflicts {
			t.Errorf("State %d: Expected hasConflicts to be %v, got %v", i, state.conflicts, csb.hasConflicts)
		}

		// Verify view is generated
		view := csb.View()
		if view == "" {
			t.Errorf("State %d: Expected non-empty view", i)
		}
	}
}

func TestCollaborationStatusBar_RenderDetails(t *testing.T) {
	csb := NewCollaborationStatusBar()

	// Test rendering details with no session
	details := csb.renderDetails()
	if details == "" {
		t.Error("Expected non-empty details even with no session")
	}

	// Test with full session info
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		false,
	)

	details = csb.renderDetails()

	// Should contain all expected fields
	expectedFields := []string{
		"Session ID: session123",
		"Participants: 3",
		"Your role: editor",
		"No conflicts",
		"Last updated:",
	}

	for _, field := range expectedFields {
		if !contains(details, field) {
			t.Errorf("Expected details to contain %q, got %q", field, details)
		}
	}

	// Test with conflicts
	csb.UpdateCollaborationState(
		"session123",
		true,
		3,
		"editor",
		true,
	)

	details = csb.renderDetails()

	if !contains(details, "Conflicts detected") {
		t.Errorf("Expected details to contain conflict message when conflicts exist, got %q", details)
	}
}
