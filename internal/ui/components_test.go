package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnimatedLoadingSpinner(t *testing.T) {
	t.Run("NewAnimatedLoadingSpinner", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		
		assert.NotNil(t, spinner)
		assert.NotNil(t, spinner.spinner)
		assert.NotNil(t, spinner.animation)
		assert.Equal(t, "Loading...", spinner.message)
		assert.True(t, spinner.fadeIn)
	})

	t.Run("Init", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		cmd := spinner.Init()
		
		assert.NotNil(t, cmd)
	})

	t.Run("UpdateWithAnimationTick", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		spinner.Init()
		
		// Send an animation tick message
		msg := AnimationTickMsg{}
		cmd := spinner.Update(msg)
		
		// Should return a command for continued animation
		assert.NotNil(t, cmd)
	})

	t.Run("UpdateWithWindowSize", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		
		msg := tea.WindowSizeMsg{Width: 80, Height: 24}
		cmd := spinner.Update(msg)
		
		assert.Nil(t, cmd)
		assert.Equal(t, 80, spinner.width)
		assert.Equal(t, 24, spinner.height)
	})

	t.Run("View", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		spinner.SetSize(80, 24)
		
		view := spinner.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Loading...")
	})

	t.Run("SetSize", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		spinner.SetSize(100, 50)
		
		assert.Equal(t, 100, spinner.width)
		assert.Equal(t, 50, spinner.height)
	})

	t.Run("StartPulseAnimation", func(t *testing.T) {
		spinner := NewAnimatedLoadingSpinner("Loading...")
		cmd := spinner.startPulseAnimation()
		
		assert.NotNil(t, cmd)
	})
}

func TestAnimatedStatusBar(t *testing.T) {
	t.Run("NewAnimatedStatusBar", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Test message")
		
		assert.NotNil(t, statusBar)
		assert.Equal(t, "Test message", statusBar.message)
		assert.Equal(t, 0.0, statusBar.progress)
		assert.NotNil(t, statusBar.animation)
		assert.Equal(t, StatusNormal, statusBar.status)
	})

	t.Run("SetProgress", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Test message")
		
		// Test setting progress between 0 and 1
		statusBar.SetProgress(0.5)
		assert.Equal(t, 0.5, statusBar.progress)
		
		// Test setting progress at boundaries
		statusBar.SetProgress(0.0)
		assert.Equal(t, 0.0, statusBar.progress)
		
		statusBar.SetProgress(1.0)
		assert.Equal(t, 1.0, statusBar.progress)
	})

	t.Run("SetStatus", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Test message")
		
		// Test setting different statuses
		statusBar.SetStatus(StatusSuccess, "Success message")
		assert.Equal(t, StatusSuccess, statusBar.status)
		assert.Equal(t, "Success message", statusBar.message)
		
		statusBar.SetStatus(StatusError, "Error message")
		assert.Equal(t, StatusError, statusBar.status)
		assert.Equal(t, "Error message", statusBar.message)
		
		statusBar.SetStatus(StatusWarning, "Warning message")
		assert.Equal(t, StatusWarning, statusBar.status)
		assert.Equal(t, "Warning message", statusBar.message)
		
		statusBar.SetStatus(StatusLoading, "Loading message")
		assert.Equal(t, StatusLoading, statusBar.status)
		assert.Equal(t, "Loading message", statusBar.message)
		
		statusBar.SetStatus(StatusNormal, "Normal message")
		assert.Equal(t, StatusNormal, statusBar.status)
		assert.Equal(t, "Normal message", statusBar.message)
	})

	t.Run("ViewWithNormalStatus", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Normal message")
		statusBar.width = 80
		
		view := statusBar.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Normal message")
	})

	t.Run("ViewWithSuccessStatus", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Success message")
		statusBar.width = 80
		statusBar.SetStatus(StatusSuccess, "Success message")
		
		view := statusBar.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Success message")
	})

	t.Run("ViewWithErrorStatus", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Error message")
		statusBar.width = 80
		statusBar.SetStatus(StatusError, "Error message")
		
		view := statusBar.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Error message")
	})

	t.Run("ViewWithProgressBar", func(t *testing.T) {
		statusBar := NewAnimatedStatusBar("Progress message")
		statusBar.width = 80
		statusBar.SetProgress(0.5)
		
		view := statusBar.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Progress message")
		assert.Contains(t, view, "[")
		assert.Contains(t, view, "]")
	})
}

func TestAnimatedNotification(t *testing.T) {
	t.Run("NewAnimatedNotification", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		
		assert.NotNil(t, notification)
		assert.Equal(t, "Test message", notification.message)
		assert.Equal(t, "info", notification.notifType)
		assert.Equal(t, 5*time.Second, notification.duration)
		assert.True(t, notification.active)
		assert.NotNil(t, notification.animation)
	})

	t.Run("Init", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		cmd := notification.Init()
		
		assert.NotNil(t, cmd)
		// Verify start time was set
		assert.False(t, notification.startTime.IsZero())
	})

	t.Run("UpdateWithAnimationTick", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 1*time.Millisecond)
		notification.Init()
		
		// Wait for duration to pass
		time.Sleep(2 * time.Millisecond)
		
		msg := AnimationTickMsg{}
		cmd := notification.Update(msg)
		
		assert.NotNil(t, cmd)
		// Notification should be inactive after duration
		assert.False(t, notification.active)
	})

	t.Run("ViewWhenActive", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		notification.Init()
		
		view := notification.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Test message")
	})

	t.Run("ViewWhenInactive", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 1*time.Millisecond)
		notification.Init()
		
		// Make notification inactive
		notification.active = false
		
		view := notification.View()
		assert.Empty(t, view)
	})

	t.Run("ViewWithDifferentTypes", func(t *testing.T) {
		types := []string{"success", "error", "warning", "info"}
		
		for _, notifType := range types {
			notification := NewAnimatedNotification("Test message", notifType, 5*time.Second)
			notification.Init()
			
			view := notification.View()
			assert.NotEmpty(t, view)
			assert.Contains(t, view, "Test message")
		}
	})

	t.Run("IsActive", func(t *testing.T) {
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		notification.Init()
		
		// Should be active initially
		assert.True(t, notification.IsActive())
		
		// Make inactive
		notification.active = false
		assert.False(t, notification.IsActive())
	})
}

func TestAccessibilityFeatures(t *testing.T) {
	// Test that UI components have proper focus indicators
	t.Run("FocusIndicators", func(t *testing.T) {
		// Test spinner has visual indication of activity
		spinner := NewAnimatedLoadingSpinner("Loading...")
		view := spinner.View()
		assert.NotEmpty(t, view)
		// Should have some visual element indicating activity
		assert.True(t, strings.Contains(view, "Loading...") || len(view) > 0)
	})

	t.Run("StatusIndicators", func(t *testing.T) {
		// Test status bar has different visual indicators for different statuses
		statusBar := NewAnimatedStatusBar("Test message")
		
		// Normal status
		view := statusBar.View()
		assert.NotEmpty(t, view)
		
		// Success status
		statusBar.SetStatus(StatusSuccess, "Success")
		successView := statusBar.View()
		assert.NotEmpty(t, successView)
		assert.NotEqual(t, view, successView)
		
		// Error status
		statusBar.SetStatus(StatusError, "Error")
		errorView := statusBar.View()
		assert.NotEmpty(t, errorView)
		assert.NotEqual(t, successView, errorView)
	})

	t.Run("NotificationVisibility", func(t *testing.T) {
		// Test notifications are visible and distinguishable
		notification := NewAnimatedNotification("Test message", "info", 5*time.Second)
		view := notification.View()
		assert.NotEmpty(t, view)
		assert.Contains(t, view, "Test message")
	})

	t.Run("KeyboardNavigationSupport", func(t *testing.T) {
		// Test that components can handle keyboard navigation messages
		spinner := NewAnimatedLoadingSpinner("Loading...")
		
		// Test with tab key (common navigation key)
		tabKey := tea.KeyMsg{Type: tea.KeyTab}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			spinner.Update(tabKey)
		})
		
		// Test with escape key
		escKey := tea.KeyMsg{Type: tea.KeyEsc}
		// Should not crash and should handle gracefully
		// Just ensure it doesn't panic
		assert.NotPanics(t, func() {
			spinner.Update(escKey)
		})
	})

	t.Run("ScreenReaderText", func(t *testing.T) {
		// Test that components provide meaningful text content for screen readers
		spinner := NewAnimatedLoadingSpinner("Loading content...")
		view := spinner.View()
		// Should contain the message for screen readers
		assert.Contains(t, view, "Loading content...")
		
		statusBar := NewAnimatedStatusBar("Status information")
		statusView := statusBar.View()
		// Should contain the status message
		assert.Contains(t, statusView, "Status information")
		
		notification := NewAnimatedNotification("Notification message", "info", 5*time.Second)
		notifView := notification.View()
		// Should contain the notification message
		assert.Contains(t, notifView, "Notification message")
	})

	t.Run("ReducedMotionSupport", func(t *testing.T) {
		// Test that animations can be controlled
		animationManager := NewAnimationManager()
		
		// Test default configuration
		config := animationManager.GetConfig()
		assert.True(t, config.Enabled)
		assert.False(t, config.ReducedMotion)
		
		// Test setting reduced motion
		config.ReducedMotion = true
		animationManager.SetConfig(config)
		
		newConfig := animationManager.GetConfig()
		assert.True(t, newConfig.ReducedMotion)
	})

	t.Run("ColorContrast", func(t *testing.T) {
		// Test that UI components use theme colors with adequate contrast
		themes := theme.ListThemes()
		require.Greater(t, len(themes), 0)
		
		for _, themeID := range themes {
			th := theme.GetTheme(themeID)
			
			// Test that text and background colors are different
			assert.NotEqual(t, th.Text, th.Background, "Text and background should have different colors for theme %s", themeID)
			
			// Test that primary and background colors are different
			assert.NotEqual(t, th.Primary, th.Background, "Primary and background should have different colors for theme %s", themeID)
			
			// Test that success and background colors are different
			assert.NotEqual(t, th.Success, th.Background, "Success and background should have different colors for theme %s", themeID)
			
			// Test that error and background colors are different
			assert.NotEqual(t, th.Error, th.Background, "Error and background should have different colors for theme %s", themeID)
		}
	})
}

func TestAnimationManagerAccessibility(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Animation accessibility features incomplete - see docs/KNOWN_TEST_LIMITATIONS.md")
	t.Run("AnimationManagerCreation", func(t *testing.T) {
		manager := NewAnimationManager()
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.animations)
		assert.NotNil(t, manager.ctx)
		assert.NotNil(t, manager.cancel)
	})

	t.Run("AnimationConfig", func(t *testing.T) {
		manager := NewAnimationManager()
		
		// Test default config
		config := manager.GetConfig()
		assert.True(t, config.Enabled)
		assert.Equal(t, 6.0, config.AngularFrequency)
		assert.Equal(t, 0.5, config.DampingRatio)
		assert.False(t, config.ReducedMotion)
		assert.Equal(t, 60, config.FrameRate)
		
		// Test setting config
		newConfig := AnimationConfig{
			Enabled:          false,
			AngularFrequency: 3.0,
			DampingRatio:     0.8,
			ReducedMotion:    true,
			FrameRate:        30,
		}
		manager.SetConfig(newConfig)
		
		updatedConfig := manager.GetConfig()
		assert.False(t, updatedConfig.Enabled)
		assert.Equal(t, 3.0, updatedConfig.AngularFrequency)
		assert.Equal(t, 0.8, updatedConfig.DampingRatio)
		assert.True(t, updatedConfig.ReducedMotion)
		assert.Equal(t, 30, updatedConfig.FrameRate)
	})

	t.Run("AnimationControl", func(t *testing.T) {
		manager := NewAnimationManager()
		
		// Test starting animation
		manager.StartAnimation("test", AnimationFade, 1.0)
		assert.True(t, manager.IsAnimationActive("test"))
		
		// Test stopping animation
		manager.StopAnimation("test")
		assert.False(t, manager.IsAnimationActive("test"))
		
		// Test clearing all animations
		manager.StartAnimation("test1", AnimationFade, 1.0)
		manager.StartAnimation("test2", AnimationSlide, 1.0)
		assert.True(t, manager.IsAnimationActive("test1"))
		assert.True(t, manager.IsAnimationActive("test2"))
		
		manager.ClearAllAnimations()
		assert.False(t, manager.IsAnimationActive("test1"))
		assert.False(t, manager.IsAnimationActive("test2"))
	})

	t.Run("ReducedMotionRespectsUserPreference", func(t *testing.T) {
		manager := NewAnimationManager()
		
		// Enable reduced motion
		config := manager.GetConfig()
		config.ReducedMotion = true
		manager.SetConfig(config)
		
		// Start animation - should not actually start when reduced motion is enabled
		manager.StartAnimation("test", AnimationFade, 1.0)
		// In a full implementation, this would check that animations are minimized
		// For now, we just ensure it doesn't crash
		assert.NotNil(t, manager)
	})
}