package dashboard

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestRecentWorkModel_Functionality tests the RecentWorkModel implementation
func TestRecentWorkModel_Functionality(t *testing.T) {
	t.Run("NewRecentWorkModel initializes correctly", func(t *testing.T) {
		model := NewRecentWorkModel()
		if model == nil {
			t.Fatal("NewRecentWorkModel returned nil")
		}
		if model.selected != 0 {
			t.Errorf("Expected selected to be 0, got %d", model.selected)
		}
		if model.hovered != -1 {
			t.Errorf("Expected hovered to be -1, got %d", model.hovered)
		}
		if !model.loading {
			t.Error("Expected loading to be true initially")
		}
	})

	t.Run("Update handles WindowSizeMsg", func(t *testing.T) {
		model := NewRecentWorkModel()
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		if model.width != 100 {
			t.Errorf("Expected width 100, got %d", model.width)
		}
		if model.height != 50 {
			t.Errorf("Expected height 50, got %d", model.height)
		}
	})

	t.Run("View renders without panic", func(t *testing.T) {
		model := NewRecentWorkModel()
		model.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

		// Should not panic
		view := model.View()
		if view == "" {
			t.Error("View returned empty string")
		}
	})
}

// TestSystemInfoModel_Functionality tests the SystemInfoModel implementation
func TestSystemInfoModel_Functionality(t *testing.T) {
	t.Run("NewSystemInfoModel initializes correctly", func(t *testing.T) {
		model := NewSystemInfoModel()
		if model == nil {
			t.Fatal("NewSystemInfoModel returned nil")
		}
		if !model.loading {
			t.Error("Expected loading to be true initially")
		}
		if model.hovered != -1 {
			t.Errorf("Expected hovered to be -1, got %d", model.hovered)
		}
	})

	t.Run("Update handles WindowSizeMsg", func(t *testing.T) {
		model := NewSystemInfoModel()
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		if model.width != 100 {
			t.Errorf("Expected width 100, got %d", model.width)
		}
		if model.height != 50 {
			t.Errorf("Expected height 50, got %d", model.height)
		}
	})

	t.Run("View renders without panic", func(t *testing.T) {
		model := NewSystemInfoModel()
		model.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

		view := model.View()
		if view == "" {
			t.Error("View returned empty string")
		}
	})

	t.Run("Init returns command", func(t *testing.T) {
		model := NewSystemInfoModel()
		cmd := model.Init()
		// Without a database, Init should still return a command (to load stats)
		if cmd == nil {
			t.Error("Init should return a command")
		}
	})
}

// TestMusicToolsModel_Functionality tests the MusicToolsModel implementation
func TestMusicToolsModel_Functionality(t *testing.T) {
	t.Run("NewMusicToolsModel initializes correctly", func(t *testing.T) {
		model := NewMusicToolsModel()
		if model == nil {
			t.Fatal("NewMusicToolsModel returned nil")
		}
		if len(model.tools) != 4 {
			t.Errorf("Expected 4 tools, got %d", len(model.tools))
		}
		if model.selected != 0 {
			t.Errorf("Expected selected to be 0, got %d", model.selected)
		}
		if model.hovered != -1 {
			t.Errorf("Expected hovered to be -1, got %d", model.hovered)
		}
	})

	t.Run("Update handles WindowSizeMsg", func(t *testing.T) {
		model := NewMusicToolsModel()
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
		if model.width != 100 {
			t.Errorf("Expected width 100, got %d", model.width)
		}
		if model.height != 50 {
			t.Errorf("Expected height 50, got %d", model.height)
		}
	})

	t.Run("Update handles keyboard navigation", func(t *testing.T) {
		model := NewMusicToolsModel()
		model.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

		// Move down
		model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		if model.selected != 1 {
			t.Errorf("Expected selected to be 1 after 'j', got %d", model.selected)
		}

		// Move up
		model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		if model.selected != 0 {
			t.Errorf("Expected selected to be 0 after 'k', got %d", model.selected)
		}
	})

	t.Run("View renders without panic", func(t *testing.T) {
		model := NewMusicToolsModel()
		model.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

		view := model.View()
		if view == "" {
			t.Error("View returned empty string")
		}
	})
}

// TestAIAssistantModel_Functionality tests the AIAssistantModel implementation
func TestAIAssistantModel_Functionality(t *testing.T) {
	t.Run("NewAIAssistantModel initializes correctly", func(t *testing.T) {
		model := NewAIAssistantModel()
		if model == nil {
			t.Fatal("NewAIAssistantModel returned nil")
		}
		if model.selected != 0 {
			t.Errorf("Expected selected to be 0, got %d", model.selected)
		}
		if model.hovered != -1 {
			t.Errorf("Expected hovered to be -1, got %d", model.hovered)
		}
	})

	t.Run("GetAIStatus returns correct status", func(t *testing.T) {
		model := NewAIAssistantModel()
		status := model.GetAIStatus()
		if status != "Not configured" {
			t.Errorf("Expected 'Not configured', got '%s'", status)
		}
	})

	t.Run("SetFocused works", func(t *testing.T) {
		model := NewAIAssistantModel()
		model.SetFocused(true)
		if !model.focused {
			t.Error("Expected focused to be true")
		}
		model.SetFocused(false)
		if model.focused {
			t.Error("Expected focused to be false")
		}
	})

	t.Run("View renders without panic", func(t *testing.T) {
		model := NewAIAssistantModel()
		model.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

		view := model.View()
		if view == "" {
			t.Error("View returned empty string")
		}
	})
}

// TestDashboardModel_MouseForwarding tests that mouse events are forwarded to children
func TestDashboardModel_MouseForwarding(t *testing.T) {
	t.Run("Dashboard forwards mouse events", func(t *testing.T) {
		dm := NewDashboardModel()
		dm.Update(tea.WindowSizeMsg{Width: 200, Height: 50})

		// Send a mouse event - should not panic
		cmd := dm.Update(tea.MouseMsg{
			X:      10,
			Y:      10,
			Button: tea.MouseButtonLeft,
			Action: tea.MouseActionPress,
		})
		// The command might be nil, but we're testing it doesn't panic
		_ = cmd
	})
}

// TestScreenChangeMsg tests the screen change message type
func TestScreenChangeMsg(t *testing.T) {
	t.Run("ScreenChangeMsg has correct screen value", func(t *testing.T) {
		msg := ScreenChangeMsg{Screen: 5}
		if msg.Screen != 5 {
			t.Errorf("Expected Screen to be 5, got %d", msg.Screen)
		}
	})
}

// TestTriggerBrainstormMsg tests the brainstorm trigger message
func TestTriggerBrainstormMsg(t *testing.T) {
	t.Run("TriggerBrainstormMsg has theme", func(t *testing.T) {
		msg := TriggerBrainstormMsg{Theme: "inspiration"}
		if msg.Theme != "inspiration" {
			t.Errorf("Expected Theme to be 'inspiration', got '%s'", msg.Theme)
		}
	})
}
