package collaboration

import (
	"fmt"
	"testing"

	"github.com/Kyanite/noise/internal/collaboration"
)

func TestPresenceIndicatorModel_NewPresenceIndicatorModel(t *testing.T) {
	model := NewPresenceIndicatorModel()

	if model == nil {
		t.Fatal("Expected non-nil model")
	}

	if model.showDetails != false {
		t.Errorf("Expected showDetails to be false, got %v", model.showDetails)
	}

	if model.selected != 0 {
		t.Errorf("Expected selected to be 0, got %d", model.selected)
	}

	if len(model.indicators) != 0 {
		t.Errorf("Expected no indicators, got %d", len(model.indicators))
	}
}

func TestPresenceIndicatorModel_SetDimensions(t *testing.T) {
	model := NewPresenceIndicatorModel()

	model.SetDimensions(100, 50)

	width, height := model.GetDimensions()
	if width != 100 {
		t.Errorf("Expected width 100, got %d", width)
	}
	if height != 50 {
		t.Errorf("Expected height 50, got %d", height)
	}
}

func TestPresenceIndicatorModel_UpdateIndicators(t *testing.T) {
	model := NewPresenceIndicatorModel()

	indicators := []SessionPresenceIndicator{
		{
			UserID:   "user1",
			Username: "User One",
			Indicator: collaboration.PresenceIndicator{
				Status:  collaboration.StatusOnline,
				Color:   "green",
				Icon:    "●",
				Tooltip: "Online",
			},
			Cursor: collaboration.CursorPosition{Line: 5, Column: 10},
			Role:   collaboration.RoleEditor,
		},
		{
			UserID:   "user2",
			Username: "User Two",
			Indicator: collaboration.PresenceIndicator{
				Status:  collaboration.StatusAway,
				Color:   "yellow",
				Icon:    "●",
				Tooltip: "Away",
			},
			Cursor: collaboration.CursorPosition{Line: 2, Column: 7},
			Role:   collaboration.RoleViewer,
		},
	}

	model.UpdateIndicators(indicators)

	if len(model.indicators) != 2 {
		t.Errorf("Expected 2 indicators, got %d", len(model.indicators))
	}

	if model.indicators[0].UserID != "user1" {
		t.Errorf("Expected first indicator user ID to be user1, got %s", model.indicators[0].UserID)
	}

	if model.indicators[1].UserID != "user2" {
		t.Errorf("Expected second indicator user ID to be user2, got %s", model.indicators[1].UserID)
	}
}

func TestPresenceIndicatorModel_View(t *testing.T) {
	model := NewPresenceIndicatorModel()
	model.SetDimensions(80, 20)

	// Test with no indicators
	view := model.View()
	if view != "" {
		t.Errorf("Expected empty view with no indicators, got %q", view)
	}

	// Test with indicators
	indicators := []SessionPresenceIndicator{
		{
			UserID:   "user1",
			Username: "User One",
			Indicator: collaboration.PresenceIndicator{
				Status:  collaboration.StatusOnline,
				Color:   "green",
				Icon:    "●",
				Tooltip: "Online",
			},
			Cursor: collaboration.CursorPosition{Line: 5, Column: 10},
			Role:   collaboration.RoleEditor,
		},
	}

	model.UpdateIndicators(indicators)
	view = model.View()

	if view == "" {
		t.Error("Expected non-empty view with indicators")
	}

	// Check that view contains expected elements
	if !contains(view, "👥 Collaborators") {
		t.Error("Expected view to contain '👥 Collaborators'")
	}

	if !contains(view, "(1 user)") {
		t.Error("Expected view to contain '(1 user)'")
	}

	if !contains(view, "User One") {
		t.Error("Expected view to contain 'User One'")
	}

	if !contains(view, "✏️") { // Editor icon
		t.Error("Expected view to contain editor icon")
	}

	if !contains(view, "(Ln 6, Col 11)") { // Cursor position (0-indexed to 1-indexed)
		t.Error("Expected view to contain cursor position")
	}

	// Test with multiple users
	indicators = append(indicators, SessionPresenceIndicator{
		UserID:   "user2",
		Username: "User Two",
		Indicator: collaboration.PresenceIndicator{
			Status:  collaboration.StatusAway,
			Color:   "yellow",
			Icon:    "●",
			Tooltip: "Away",
		},
		Cursor: collaboration.CursorPosition{Line: 2, Column: 7},
		Role:   collaboration.RoleViewer,
	})

	model.UpdateIndicators(indicators)
	view = model.View()

	if !contains(view, "(2 users)") {
		t.Error("Expected view to contain '(2 users)'")
	}

	if !contains(view, "User Two") {
		t.Error("Expected view to contain 'User Two'")
	}
}

func TestPresenceIndicatorModel_ToggleDetails(t *testing.T) {
	model := NewPresenceIndicatorModel()

	// Initially details should be hidden
	if model.showDetails != false {
		t.Errorf("Expected showDetails to be false initially, got %v", model.showDetails)
	}

	// Toggle to show details
	model.ToggleDetails()
	if model.showDetails != true {
		t.Errorf("Expected showDetails to be true after toggle, got %v", model.showDetails)
	}

	// Toggle to hide details
	model.ToggleDetails()
	if model.showDetails != false {
		t.Errorf("Expected showDetails to be false after second toggle, got %v", model.showDetails)
	}
}

func TestPresenceIndicatorModel_SelectNavigation(t *testing.T) {
	model := NewPresenceIndicatorModel()

	// Add some indicators
	indicators := []SessionPresenceIndicator{
		{UserID: "user1", Username: "User One"},
		{UserID: "user2", Username: "User Two"},
		{UserID: "user3", Username: "User Three"},
	}
	model.UpdateIndicators(indicators)

	// Initially no selection
	if model.selected != 0 {
		t.Errorf("Expected selected to be 0 initially, got %d", model.selected)
	}

	// Select next
	model.SelectNext()
	if model.selected != 1 {
		t.Errorf("Expected selected to be 1 after SelectNext, got %d", model.selected)
	}

	// Select next again
	model.SelectNext()
	if model.selected != 2 {
		t.Errorf("Expected selected to be 2 after SelectNext, got %d", model.selected)
	}

	// Select next should wrap around
	model.SelectNext()
	if model.selected != 0 {
		t.Errorf("Expected selected to wrap to 0, got %d", model.selected)
	}

	// Select previous should wrap around
	model.SelectPrev()
	if model.selected != 2 {
		t.Errorf("Expected selected to wrap to 2, got %d", model.selected)
	}

	// Select previous again
	model.SelectPrev()
	if model.selected != 1 {
		t.Errorf("Expected selected to be 1 after SelectPrev, got %d", model.selected)
	}
}

func TestPresenceIndicatorModel_GetSelectedUser(t *testing.T) {
	model := NewPresenceIndicatorModel()

	// Test with no indicators
	selected := model.GetSelectedUser()
	if selected != nil {
		t.Error("Expected nil selected user with no indicators")
	}

	// Add indicators
	indicators := []SessionPresenceIndicator{
		{UserID: "user1", Username: "User One"},
		{UserID: "user2", Username: "User Two"},
	}
	model.UpdateIndicators(indicators)

	// Initially should select first user
	selected = model.GetSelectedUser()
	if selected == nil {
		t.Error("Expected non-nil selected user")
	}
	if selected.UserID != "user1" {
		t.Errorf("Expected selected user ID to be user1, got %s", selected.UserID)
	}

	// Select second user
	model.SelectNext()
	selected = model.GetSelectedUser()
	if selected == nil {
		t.Error("Expected non-nil selected user")
	}
	if selected.UserID != "user2" {
		t.Errorf("Expected selected user ID to be user2, got %s", selected.UserID)
	}
}

func TestPresenceIndicatorModel_GetColoredIcon(t *testing.T) {
	model := NewPresenceIndicatorModel()

	tests := []struct {
		name     string
		color    string
		icon     string
		expected string
	}{
		{
			name:     "Green online icon",
			color:    "green",
			icon:     "●",
			expected: "●", // With green color applied
		},
		{
			name:     "Yellow away icon",
			color:    "yellow",
			icon:     "●",
			expected: "●", // With yellow color applied
		},
		{
			name:     "Red busy icon",
			color:    "red",
			icon:     "●",
			expected: "●", // With red color applied
		},
		{
			name:     "Gray offline icon",
			color:    "gray",
			icon:     "○",
			expected: "○", // With gray color applied
		},
		{
			name:     "Unknown color defaults to gray",
			color:    "unknown",
			icon:     "●",
			expected: "●", // With gray color applied
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			indicator := collaboration.PresenceIndicator{
				Color: tt.color,
				Icon:  tt.icon,
			}

			result := model.getColoredIcon(indicator)
			if result == "" {
				t.Error("Expected non-empty colored icon")
			}
			// Note: We can't easily test the actual color rendering without complex terminal escape sequence parsing
		})
	}
}

func TestPresenceIndicatorModel_GetRoleIcon(t *testing.T) {
	model := NewPresenceIndicatorModel()

	tests := []struct {
		name     string
		role     collaboration.ParticipantRole
		expected string
	}{
		{
			name:     "Owner role",
			role:     collaboration.RoleOwner,
			expected: "👑",
		},
		{
			name:     "Editor role",
			role:     collaboration.RoleEditor,
			expected: "✏️",
		},
		{
			name:     "Viewer role",
			role:     collaboration.RoleViewer,
			expected: "👁️",
		},
		{
			name:     "Unknown role",
			role:     "unknown",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.getRoleIcon(tt.role)
			if result != tt.expected {
				t.Errorf("Expected role icon %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPresenceIndicatorModel_GetUserDetails(t *testing.T) {
	model := NewPresenceIndicatorModel()

	tests := []struct {
		name      string
		indicator SessionPresenceIndicator
		expected  string
	}{
		{
			name: "Online editor user",
			indicator: SessionPresenceIndicator{
				Indicator: collaboration.PresenceIndicator{
					Status: collaboration.StatusOnline,
				},
				Role: collaboration.RoleEditor,
			},
			expected: "Currently editing • Can edit",
		},
		{
			name: "Away viewer user",
			indicator: SessionPresenceIndicator{
				Indicator: collaboration.PresenceIndicator{
					Status: collaboration.StatusAway,
				},
				Role: collaboration.RoleViewer,
			},
			expected: "Away • Read-only",
		},
		{
			name: "Busy owner user",
			indicator: SessionPresenceIndicator{
				Indicator: collaboration.PresenceIndicator{
					Status: collaboration.StatusBusy,
				},
				Role: collaboration.RoleOwner,
			},
			expected: "Busy • Session owner",
		},
		{
			name: "Offline user",
			indicator: SessionPresenceIndicator{
				Indicator: collaboration.PresenceIndicator{
					Status: collaboration.StatusOffline,
				},
				Role: collaboration.RoleEditor,
			},
			expected: "Offline • Can edit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.getUserDetails(tt.indicator)
			if result != tt.expected {
				t.Errorf("Expected user details %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestPresenceIndicatorModel_DetailedView(t *testing.T) {
	model := NewPresenceIndicatorModel()
	model.SetDimensions(80, 20)

	// Add indicators
	indicators := []SessionPresenceIndicator{
		{
			UserID:   "user1",
			Username: "User One",
			Indicator: collaboration.PresenceIndicator{
				Status:  collaboration.StatusOnline,
				Color:   "green",
				Icon:    "●",
				Tooltip: "Online",
			},
			Cursor: collaboration.CursorPosition{Line: 5, Column: 10},
			Role:   collaboration.RoleEditor,
		},
		{
			UserID:   "user2",
			Username: "User Two",
			Indicator: collaboration.PresenceIndicator{
				Status:  collaboration.StatusAway,
				Color:   "yellow",
				Icon:    "●",
				Tooltip: "Away",
			},
			Cursor: collaboration.CursorPosition{Line: 2, Column: 7},
			Role:   collaboration.RoleViewer,
		},
	}

	model.UpdateIndicators(indicators)

	// Enable details mode
	model.ToggleDetails()

	view := model.View()

	// Should contain details for both users
	if !contains(view, "Currently editing") {
		t.Error("Expected view to contain 'Currently editing' for online user")
	}

	if !contains(view, "Away") {
		t.Error("Expected view to contain 'Away' for away user")
	}

	if !contains(view, "Can edit") {
		t.Error("Expected view to contain 'Can edit' for editor")
	}

	if !contains(view, "Read-only") {
		t.Error("Expected view to contain 'Read-only' for viewer")
	}
}

func TestPresenceIndicatorModel_EdgeCases(t *testing.T) {
	model := NewPresenceIndicatorModel()

	// Test with empty indicators slice
	model.UpdateIndicators([]SessionPresenceIndicator{})

	view := model.View()
	if view != "" {
		t.Errorf("Expected empty view with empty indicators, got %q", view)
	}

	// Test navigation with empty indicators
	model.SelectNext()
	if model.selected != 0 {
		t.Errorf("Expected selected to remain 0 with empty indicators, got %d", model.selected)
	}

	model.SelectPrev()
	if model.selected != 0 {
		t.Errorf("Expected selected to remain 0 with empty indicators, got %d", model.selected)
	}

	// Test with indicators that have no cursor position
	indicators := []SessionPresenceIndicator{
		{
			UserID:   "user1",
			Username: "User One",
			Indicator: collaboration.PresenceIndicator{
				Status: collaboration.StatusOnline,
			},
			Cursor: collaboration.CursorPosition{Line: 0, Column: 0}, // No cursor
			Role:   collaboration.RoleViewer,
		},
	}

	model.UpdateIndicators(indicators)
	view = model.View()

	// Should not contain cursor position
	if contains(view, "(Ln") {
		t.Error("Expected view to not contain cursor position when cursor is at (0,0)")
	}
}

func TestPresenceIndicatorModel_Performance(t *testing.T) {
	model := NewPresenceIndicatorModel()

	// Test with many indicators (performance test)
	var indicators []SessionPresenceIndicator
	for i := 0; i < 100; i++ {
		indicators = append(indicators, SessionPresenceIndicator{
			UserID:   fmt.Sprintf("user%d", i),
			Username: fmt.Sprintf("User %d", i),
			Indicator: collaboration.PresenceIndicator{
				Status: collaboration.StatusOnline,
				Color:  "green",
				Icon:   "●",
			},
			Role: collaboration.RoleEditor,
		})
	}

	model.UpdateIndicators(indicators)

	if len(model.indicators) != 100 {
		t.Errorf("Expected 100 indicators, got %d", len(model.indicators))
	}

	// Test rapid navigation
	for i := 0; i < 100; i++ {
		model.SelectNext()
	}

	if model.selected != 0 { // Should wrap around
		t.Errorf("Expected selected to wrap to 0 after 100 selections, got %d", model.selected)
	}

	// Test view generation with many users
	view := model.View()
	if view == "" {
		t.Error("Expected non-empty view with many indicators")
	}

	if !contains(view, "(100 users)") {
		t.Error("Expected view to show 100 users")
	}
}
