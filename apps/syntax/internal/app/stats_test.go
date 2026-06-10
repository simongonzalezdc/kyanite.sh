package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

func TestHandleStatsKeys(t *testing.T) {
	tests := []struct {
		name           string
		key            string
		expectedScreen Screen
	}{
		{"esc returns to editor", "esc", ScreenEditor},
		{"unknown key stays on stats", "x", ScreenStats},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenStats
			m.CurrentProject = &story.Project{Title: "Test"}

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			result, _ := m.handleStatsKeys(msg)
			resultModel := result.(Model)

			if resultModel.CurrentScreen != tt.expectedScreen {
				t.Errorf("CurrentScreen = %v, expected %v", resultModel.CurrentScreen, tt.expectedScreen)
			}
		})
	}
}

func TestViewStats_NoProject(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenStats

	view := m.viewStats()

	if view != "No project loaded" {
		t.Errorf("viewStats() = %q, expected 'No project loaded'", view)
	}
}

func TestViewStats_WithProject(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenStats
	m.Width = 80
	m.Height = 24
	m.CurrentProject = &story.Project{
		Title:            "Test Novel",
		TargetWordCount:  50000,
		TotalSessions:    0,
		TotalTimeSeconds: 0,
		Scenes: map[string]*scene.Scene{
			"scene1": {WordCount: 1000},
			"scene2": {WordCount: 1500},
			"scene3": {WordCount: 2500},
		},
		Characters: make(map[string]*character.Character),
		Locations:  make(map[string]*location.Location),
	}

	view := m.viewStats()

	expectedContent := []string{
		"Test Novel",
		"Statistics",
		"5000", // Total words (1000 + 1500 + 2500)
		"50000", // Target word count
		"10.0%", // Progress (5000 / 50000 * 100)
		"3", // Number of scenes
		"Press Esc",
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewStats should contain %q", content)
		}
	}
}

func TestViewStats_ProgressCalculation(t *testing.T) {
	tests := []struct {
		name            string
		totalWords      int
		targetWordCount int
		expectedProg    string
	}{
		{"no target", 1000, 0, "0.0%"},
		{"halfway", 5000, 10000, "50.0%"},
		{"complete", 10000, 10000, "100.0%"},
		{"over target", 12000, 10000, "120.0%"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenStats
			m.Width = 80
			m.Height = 24
			m.CurrentProject = &story.Project{
				Title:            "Test",
				TargetWordCount:  tt.targetWordCount,
				Scenes:           make(map[string]*scene.Scene),
				Characters:       make(map[string]*character.Character),
				Locations:        make(map[string]*location.Location),
			}

			// Add scenes to reach target word count
			if tt.totalWords > 0 {
				m.CurrentProject.Scenes["scene1"] = &scene.Scene{WordCount: tt.totalWords}
			}

			view := m.viewStats()

			if !contains(view, tt.expectedProg) {
				t.Errorf("viewStats should contain progress %q", tt.expectedProg)
			}
		})
	}
}

func TestViewStats_WithSessionHistory(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenStats
	m.Width = 80
	m.Height = 24
	m.CurrentProject = &story.Project{
		Title:            "Test Novel",
		TotalSessions:    25,
		CurrentStreak:    7,
		TotalTimeSeconds: 7265, // 2 hours, 1 minute, 5 seconds
		Scenes:           make(map[string]*scene.Scene),
		Characters:       make(map[string]*character.Character),
		Locations:        make(map[string]*location.Location),
	}

	view := m.viewStats()

	expectedContent := []string{
		"Session History",
		"25", // Total sessions
		"7 days", // Current streak
		"2h 1m", // Total time (7265 seconds = 2h 1m 5s, showing hours and minutes)
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewStats should contain %q", content)
		}
	}
}

func TestViewStats_TimeFormatting(t *testing.T) {
	tests := []struct {
		name         string
		seconds      int
		expectedTime string
	}{
		{"less than hour", 1800, "0h 30m"},
		{"exactly one hour", 3600, "1h 0m"},
		{"multiple hours", 9000, "2h 30m"},
		{"many hours", 36000, "10h 0m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewModel()
			m.CurrentScreen = ScreenStats
			m.Width = 80
			m.Height = 24
			m.CurrentProject = &story.Project{
				Title:            "Test",
				TotalSessions:    1,
				TotalTimeSeconds: tt.seconds,
				Scenes:           make(map[string]*scene.Scene),
				Characters:       make(map[string]*character.Character),
				Locations:        make(map[string]*location.Location),
			}

			view := m.viewStats()

			if !contains(view, tt.expectedTime) {
				t.Errorf("viewStats should contain time %q, got view: %s", tt.expectedTime, view)
			}
		})
	}
}

func TestViewStats_EmptyProject(t *testing.T) {
	m := NewModel()
	m.CurrentScreen = ScreenStats
	m.Width = 80
	m.Height = 24
	m.CurrentProject = &story.Project{
		Title:      "Empty Project",
		Scenes:     make(map[string]*scene.Scene),
		Characters: make(map[string]*character.Character),
		Locations:  make(map[string]*location.Location),
	}

	view := m.viewStats()

	expectedContent := []string{
		"Empty Project",
		"0", // Total words
		"0.0%", // Progress
	}

	for _, content := range expectedContent {
		if !contains(view, content) {
			t.Errorf("viewStats should contain %q", content)
		}
	}
}
