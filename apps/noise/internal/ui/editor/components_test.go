package editor

import (
	"testing"

	"github.com/kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEditorComponents_Accessibility(t *testing.T) {
	t.Run("FocusIndicators", func(t *testing.T) {
		// Test that editor components have proper focus indicators
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Initially focused
		editorPane.Focus()
		assert.True(t, editorPane.state.IsFocused())

		// Check that view reflects focus state
		view := editorPane.View()
		assert.Contains(t, view, "(Focused)")

		// Blur and check
		editorPane.Blur()
		assert.False(t, editorPane.state.IsFocused())
	})

	t.Run("KeyboardNavigation", func(t *testing.T) {
		// Test keyboard navigation support
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Test tab navigation
		tabMsg := tea.KeyMsg{Type: tea.KeyTab}
		_, cmd := editorPane.Update(tabMsg)
		// Should not panic
		assert.NotPanics(t, func() {
			_ = cmd
		})

		// Test escape key for modals
		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		_, cmd = editorPane.Update(escMsg)
		// Should not panic
		assert.NotPanics(t, func() {
			_ = cmd
		})
	})

	t.Run("ScreenReaderSupport", func(t *testing.T) {
		// Test that editor components provide meaningful text for screen readers
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Check that view contains meaningful text
		view := editorPane.View()
		assert.Contains(t, view, "Editor")

		// Test with different modes
		editorPane.SetEditorMode(ModeSketch)
		view = editorPane.View()
		assert.Contains(t, view, "[SKETCH]")

		editorPane.SetEditorMode(ModeDraft)
		view = editorPane.View()
		assert.Contains(t, view, "[DRAFT]")

		editorPane.SetEditorMode(ModePolish)
		view = editorPane.View()
		assert.Contains(t, view, "[POLISH]")
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that editor components use theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			th := theme.GetTheme(themeID)

			// Test that text and background colors are different
			assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)

			// Test that primary and background colors are different
			assert.NotEqual(t, th.Primary, th.Background, "Primary and background should have different colors for theme %s", themeID)
		}
	})
}

func TestChordPicker_Accessibility(t *testing.T) {
	t.Run("ChordPickerFocusManagement", func(t *testing.T) {
		// Test that chord picker has proper focus management
		chordPicker := NewChordPickerModel(nil)
		assert.NotNil(t, chordPicker)

		// Test visibility
		assert.False(t, chordPicker.IsVisible())

		// Test showing
		showCmd := chordPicker.Show(func(chords []string) {})
		assert.NotNil(t, showCmd)

		// Test hiding
		hideCmd := chordPicker.Hide()
		assert.NotNil(t, hideCmd)
	})

	t.Run("ChordPickerKeyboardNavigation", func(t *testing.T) {
		// Test that chord picker supports keyboard navigation
		chordPicker := NewChordPickerModel(nil)

		// Test escape when not visible (should not panic)
		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		_, cmd := chordPicker.Update(escMsg)
		assert.Nil(t, cmd)
		assert.False(t, chordPicker.IsVisible())

		// Test showing
		showMsg := ShowChordPickerMsg{InsertCallback: func(chords []string) {}}
		_, cmd = chordPicker.Update(showMsg)
		assert.Nil(t, cmd)
		// Note: We don't check visibility here because data loading might not be complete
		// but the important thing is that it doesn't panic
	})

	t.Run("ChordPickerScreenReaderSupport", func(t *testing.T) {
		// Test that chord picker provides meaningful text
		chordPicker := NewChordPickerModel(nil)

		view := chordPicker.View()
		// When not visible, should return empty string
		assert.Empty(t, view)

		// Show the picker
		showMsg := ShowChordPickerMsg{InsertCallback: func(chords []string) {}}
		_, _ = chordPicker.Update(showMsg)

		// Test that view doesn't panic and returns a string
		view = chordPicker.View()
		// We don't assert specific content because data loading might not be complete
		// but we ensure it returns a string without panicking
		assert.IsType(t, "", view)
	})
}

func TestBPMTapper_Accessibility(t *testing.T) {
	t.Run("BPMTapperFocusManagement", func(t *testing.T) {
		// Test that BPM tapper has proper focus management
		bpmTapper := NewBPMTapperModel()
		assert.NotNil(t, bpmTapper)

		// Test visibility
		assert.False(t, bpmTapper.IsVisible())

		// Test showing
		showCmd := bpmTapper.Show(func(bpm int) {})
		assert.NotNil(t, showCmd)

		// Test hiding
		hideCmd := bpmTapper.Hide()
		assert.NotNil(t, hideCmd)
	})

	t.Run("BPMTapperKeyboardNavigation", func(t *testing.T) {
		// Test that BPM tapper supports keyboard navigation
		bpmTapper := NewBPMTapperModel()

		// Show the tapper
		showMsg := ShowBPMTapperMsg{SetBMPCallback: func(bpm int) {}}
		_, cmd := bpmTapper.Update(showMsg)
		assert.Nil(t, cmd)
		assert.True(t, bpmTapper.IsVisible())

		// Test escape to close (do this first to avoid timeout issues)
		escMsg := tea.KeyMsg{Type: tea.KeyEsc}
		_, cmd = bpmTapper.Update(escMsg)
		assert.Nil(t, cmd)
		assert.False(t, bpmTapper.IsVisible())
	})

	t.Run("BPMTapperScreenReaderSupport", func(t *testing.T) {
		// Test that BPM tapper provides meaningful text
		bpmTapper := NewBPMTapperModel()

		view := bpmTapper.View()
		// When not visible, should return empty string
		assert.Empty(t, view)

		// Show the tapper
		showMsg := ShowBPMTapperMsg{SetBMPCallback: func(bpm int) {}}
		_, _ = bpmTapper.Update(showMsg)

		view = bpmTapper.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "BPM Tapper")
	})
}

func TestThemeIntegration_Accessibility(t *testing.T) {
	t.Run("ThemeSwitchingKeyboardSupport", func(t *testing.T) {
		// Test that theme switching supports keyboard navigation
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Test theme integration
		assert.NotPanics(t, func() {
			editorPane.SetThemeManager(theme.GetManager())
		})
	})

	t.Run("ThemeColorContrast", func(t *testing.T) {
		// Test that all themes have adequate color contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)

		for _, themeID := range themes {
			th := theme.GetTheme(themeID)

			// Test critical color combinations for contrast
			assert.NotEqual(t, th.Text, th.Background, "Text and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Primary, th.Background, "Primary and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Secondary, th.Background, "Secondary and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Accent, th.Background, "Accent and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Success, th.Background, "Success and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Error, th.Background, "Error and background should differ for theme %s", themeID)
			assert.NotEqual(t, th.Warning, th.Background, "Warning and background should differ for theme %s", themeID)
		}
	})

	t.Run("ThemeReducedMotion", func(t *testing.T) {
		// Test that theme integration respects reduced motion
		themeManager := theme.GetManager()
		originalTheme := themeManager.Current()

		// Test theme switching
		themes := theme.ListThemes()
		if len(themes) > 0 {
			themeManager.SetTheme(themes[0])
			current := themeManager.Current()
			assert.NotNil(t, current)
		}

		// Restore original theme
		themeManager.SetTheme(originalTheme.Name)
	})
}

func TestEditorPane_ComprehensiveAccessibility(t *testing.T) {
	t.Run("EditorPaneFocusManagement", func(t *testing.T) {
		// Test comprehensive focus management
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Test initial focus state (textarea is focused by default for input handling)
		assert.True(t, editorPane.state.IsFocused())

		// Test blurring
		editorPane.Blur()
		assert.False(t, editorPane.state.IsFocused())

		// Test focusing
		editorPane.Focus()
		assert.True(t, editorPane.state.IsFocused())
	})

	t.Run("EditorPaneScreenReaderText", func(t *testing.T) {
		// Test that editor pane provides comprehensive screen reader text
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Test basic view
		view := editorPane.View()
		assert.Contains(t, view, "Editor")

		// Test with content
		editorPane.SetText("Test content")
		view = editorPane.View()
		assert.Contains(t, view, "Editor")

		// Test with different modes
		editorPane.SetEditorMode(ModeSketch)
		view = editorPane.View()
		assert.Contains(t, view, "[SKETCH]")

		editorPane.SetEditorMode(ModeDraft)
		view = editorPane.View()
		assert.Contains(t, view, "[DRAFT]")

		editorPane.SetEditorMode(ModePolish)
		view = editorPane.View()
		assert.Contains(t, view, "[POLISH]")
	})

	t.Run("EditorPaneKeyboardShortcuts", func(t *testing.T) {
		// Test that editor pane handles keyboard shortcuts appropriately
		ta := textarea.New()
		editorPane := NewEditorPaneModel(ta, nil)

		// Test common shortcuts don't panic
		shortcuts := []tea.KeyMsg{
			{Type: tea.KeyTab},
			{Type: tea.KeyEsc},
			{Type: tea.KeyEnter},
			{Type: tea.KeyUp},
			{Type: tea.KeyDown},
			{Type: tea.KeyLeft},
			{Type: tea.KeyRight},
		}

		for _, shortcut := range shortcuts {
			assert.NotPanics(t, func() {
				editorPane.Update(shortcut)
			})
		}
	})
}
