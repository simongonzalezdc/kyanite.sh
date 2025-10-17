package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// AudioModel handles the audio tools screen
type AudioModel struct {
	width  int
	height int
}

// NewAudioModel creates a new audio model
func NewAudioModel() *AudioModel {
	return &AudioModel{}
}

// Init initializes the audio model
func (m *AudioModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the audio tools
func (m *AudioModel) Update(msg tea.Msg) (*AudioModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the audio tools screen
func (m *AudioModel) View() string {
	return "Audio tools not implemented yet"
}
