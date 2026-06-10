package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	palette "github.com/kyanite/prism/internal/palette"
	theme "github.com/kyanite/prism/internal/theme"
)

// TestUIModelInitialization tests UI model creation and initialization
func TestUIModelInitialization(t *testing.T) {
	tm := theme.NewManager()

	t.Run("GeneratorModel", func(t *testing.T) {
		model := NewGeneratorModel(tm)

		// Verify initial state
		if model.Init() != nil {
			t.Error("GeneratorModel.Init() should return nil")
		}

		// Verify the model can be rendered without panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("GeneratorModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("CheckerModel", func(t *testing.T) {
		model := NewCheckerModel(tm)

		if model.Init() != nil {
			t.Error("CheckerModel.Init() should return nil")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("CheckerModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("MenuModel", func(t *testing.T) {
		model := NewMenuModel(tm)

		if model.Init() != nil {
			t.Error("MenuModel.Init() should return nil")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MenuModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("SearchModel", func(t *testing.T) {
		model := NewSearchModel(tm)

		if model.Init() != nil {
			t.Error("SearchModel.Init() should return nil")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("SearchModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("ManagerModel", func(t *testing.T) {
		model := NewManagerModel(tm)

		// Init may return a command for loading palettes
		cmd := model.Init()
		if cmd != nil {
			// Execute the command if it exists
			_ = cmd()
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("ManagerModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("WheelModel", func(t *testing.T) {
		model := NewWheelModel(tm)

		if model.Init() != nil {
			t.Error("WheelModel.Init() should return nil")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("WheelModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})

	t.Run("TheoryModel", func(t *testing.T) {
		model := NewTheoryModel(tm)

		if model.Init() != nil {
			t.Error("TheoryModel.Init() should return nil")
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("TheoryModel.View() panicked: %v", r)
			}
		}()
		_ = model.View()
	})
}

// TestUIModelNavigation tests navigation in UI models
func TestUIModelNavigation(t *testing.T) {
	tm := theme.NewManager()

	t.Run("MenuNavigation", func(t *testing.T) {
		model := NewMenuModel(tm)

		// Simulate down key
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
		newModel, _ := model.Update(msg)

		// Model should update without panic
		// Note: We can't compare models directly as they contain uncomparable types
		_ = newModel

		// Simulate up key
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
		newModel, _ = newModel.Update(msg)
		_ = newModel
	})

	t.Run("GeneratorRuleSelection", func(t *testing.T) {
		model := NewGeneratorModel(tm)

		// Navigate through rules
		for i := 0; i < 3; i++ {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model, _ = model.Update(msg)
		}

		// Navigate back up
		for i := 0; i < 2; i++ {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			model, _ = model.Update(msg)
		}

		// Model should handle navigation without panic
		_ = model.View()
	})
}

// TestUIModelRendering tests that all models can render
func TestUIModelRendering(t *testing.T) {
	tm := theme.NewManager()

	models := []struct {
		name  string
		model interface{ View() string }
	}{
		{"GeneratorModel", NewGeneratorModel(tm)},
		{"CheckerModel", NewCheckerModel(tm)},
		{"MenuModel", NewMenuModel(tm)},
		{"SearchModel", NewSearchModel(tm)},
		{"ManagerModel", NewManagerModel(tm)},
		{"WheelModel", NewWheelModel(tm)},
		{"TheoryModel", NewTheoryModel(tm)},
	}

	for _, m := range models {
		t.Run(m.name+"_Render", func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s.View() panicked: %v", m.name, r)
				}
			}()

			output := m.model.View()
			if output == "" {
				t.Logf("Warning: %s.View() returned empty string", m.name)
			}
		})
	}
}

// TestUIThemeIntegration tests theme integration with UI
func TestUIThemeIntegration(t *testing.T) {
	tm := theme.NewManager()

	t.Run("ThemeManager", func(t *testing.T) {
		initialTheme := tm.CurrentTheme()
		if initialTheme.Name == "" {
			t.Error("CurrentTheme should have a name")
		}

		// Cycle to next theme
		tm.NextTheme()
		nextTheme := tm.CurrentTheme()

		// Theme should potentially change
		t.Logf("Initial theme: %s, Next theme: %s", initialTheme.Name, nextTheme.Name)
	})

	t.Run("StylesCreation", func(t *testing.T) {
		currentTheme := tm.CurrentTheme()
		styles := NewStyles(currentTheme)

		// Styles should be created without panic
		_ = styles
	})

	t.Run("AllThemes", func(t *testing.T) {
		themes := theme.AllThemes()

		if len(themes) == 0 {
			t.Error("AllThemes should return at least one theme")
		}

		for i := range themes {
			if themes[i].Name == "" {
				t.Error("Theme should have a name")
			}

			// Create styles for each theme
			_ = NewStyles(&themes[i])
		}
	})

	t.Run("GetTheme", func(t *testing.T) {
		themes := theme.AllThemes()
		if len(themes) > 0 {
			themeName := themes[0].Name
			retrieved := theme.GetTheme(themeName)

			if retrieved == nil {
				t.Errorf("GetTheme(%s) should return a theme", themeName)
			}

			if retrieved != nil && retrieved.Name != themeName {
				t.Errorf("Retrieved theme name = %s, want %s", retrieved.Name, themeName)
			}
		}
	})

	t.Run("SetTheme", func(t *testing.T) {
		themes := theme.AllThemes()
		if len(themes) >= 2 {
			// Set to second theme
			tm.SetTheme(&themes[1])
			current := tm.CurrentTheme()

			if current.Name != themes[1].Name {
				t.Errorf("SetTheme failed: current = %s, want %s", current.Name, themes[1].Name)
			}
		}
	})
}

// TestUIModelInteractions tests model interactions
func TestUIModelInteractions(t *testing.T) {
	tm := theme.NewManager()

	t.Run("GeneratorKeyHandling", func(t *testing.T) {
		model := NewGeneratorModel(tm)

		// Test various key presses
		keys := []string{"j", "k", "enter", "c", "s", "e", "q"}

		for _, key := range keys {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
			model, _ = model.Update(msg)

			// Should not panic
			_ = model.View()
		}
	})

	t.Run("CheckerKeyHandling", func(t *testing.T) {
		model := NewCheckerModel(tm)

		keys := []string{"enter", "c", "q"}

		for _, key := range keys {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
			model, _ = model.Update(msg)
			_ = model.View()
		}
	})
}

// TestUINavigateMessage tests the Navigate message
func TestUINavigateMessage(t *testing.T) {
	screens := []int{0, 1, 2, 3, 4, 5, 6}

	for _, screen := range screens {
		msg := Navigate(screen)

		// Should create a valid message
		if msg.Screen != screen {
			t.Errorf("Navigate(%d) screen = %d, want %d", screen, msg.Screen, screen)
		}
	}
}

// TestUIHelperFunctions tests UI helper functions
func TestUIHelperFunctions(t *testing.T) {
	tm := theme.NewManager()

	t.Run("RenderHelp", func(t *testing.T) {
		// RenderHelp requires theme manager and dimensions
		helpText := RenderHelp(tm, 80, 24)

		if helpText == "" {
			t.Error("RenderHelp should return non-empty string")
		}
	})
}

// TestUIModelStatePersistence tests that models maintain state correctly
func TestUIModelStatePersistence(t *testing.T) {
	tm := theme.NewManager()

	t.Run("SearchModelQueryPersistence", func(t *testing.T) {
		model := NewSearchModel(tm)

		// Simulate typing (this is a simplified test)
		// In real usage, the model would receive character inputs

		// Update multiple times
		for i := 0; i < 5; i++ {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
			model, _ = model.Update(msg)
		}

		// Model should maintain internal state
		_ = model.View()
	})

	t.Run("ManagerModelListPersistence", func(t *testing.T) {
		model := NewManagerModel(tm)

		// Initialize (may load palettes)
		cmd := model.Init()
		if cmd != nil {
			_ = cmd()
		}

		// Navigate through list
		for i := 0; i < 3; i++ {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model, _ = model.Update(msg)
		}

		// Model should maintain selection state
		_ = model.View()
	})
}

// TestUIModelEdgeCases tests edge cases in UI models
func TestUIModelEdgeCases(t *testing.T) {
	tm := theme.NewManager()

	t.Run("UnknownKeyPress", func(t *testing.T) {
		model := NewGeneratorModel(tm)

		// Send unknown key
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'~'}}
		newModel, _ := model.Update(msg)

		// Should handle gracefully
		_ = newModel.View()
	})

	t.Run("WindowSizeMessage", func(t *testing.T) {
		model := NewGeneratorModel(tm)

		// Send window size message
		msg := tea.WindowSizeMsg{Width: 80, Height: 24}
		newModel, _ := model.Update(msg)

		// Should handle without panic
		_ = newModel.View()
	})

	t.Run("RapidKeyPresses", func(t *testing.T) {
		model := NewMenuModel(tm)

		// Simulate rapid navigation
		for i := 0; i < 100; i++ {
			var msg tea.Msg
			if i%2 == 0 {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			} else {
				msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
			}
			model, _ = model.Update(msg)
		}

		// Should handle without panic or errors
		_ = model.View()
	})
}

// TestUIModelWithRealData tests models with actual data
func TestUIModelWithRealData(t *testing.T) {
	tm := theme.NewManager()

	t.Run("TheoryModelWithPalette", func(t *testing.T) {
		model := NewTheoryModel(tm)

		// Theory model should work with default data
		_ = model.View()

		// Navigate through different sections
		for i := 0; i < 5; i++ {
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
			model, _ = model.Update(msg)
			_ = model.View()
		}
	})

	t.Run("WheelModelRendering", func(t *testing.T) {
		model := NewWheelModel(tm)

		// Wheel model should render the color wheel
		output := model.View()

		// Should produce some output
		if len(output) < 10 {
			t.Error("WheelModel should produce substantial output")
		}
	})
}

// TestHarmonyRuleDisplay tests that all harmony rules can be displayed
func TestHarmonyRuleDisplay(t *testing.T) {
	tm := theme.NewManager()
	model := NewGeneratorModel(tm)

	allRules := palette.AllRules()

	for i, rule := range allRules {
		t.Run(string(rule), func(t *testing.T) {
			// Navigate to this rule
			for j := 0; j < i; j++ {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
				model, _ = model.Update(msg)
			}

			// Render
			output := model.View()

			// Should contain something
			if output == "" {
				t.Errorf("View with rule %s should not be empty", rule)
			}

			// Reset to beginning
			for j := 0; j < i; j++ {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
				model, _ = model.Update(msg)
			}
		})
	}
}
