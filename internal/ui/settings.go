package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// SettingsModel handles the settings screen
type SettingsModel struct {
	width  int
	height int
}

// NewSettingsModel creates a new settings model
func NewSettingsModel() *SettingsModel {
	return &SettingsModel{}
}

// Init initializes the settings model
func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the settings
func (m *SettingsModel) Update(msg tea.Msg) (*SettingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the settings screen
func (m *SettingsModel) View() string {
	return "Settings not implemented yet"
}
