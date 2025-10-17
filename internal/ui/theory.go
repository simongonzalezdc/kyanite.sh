package ui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TheoryModel handles the music theory tools screen
type TheoryModel struct {
	width  int
	height int
}

// NewTheoryModel creates a new theory model
func NewTheoryModel() *TheoryModel {
	return &TheoryModel{}
}

// Init initializes the theory model
func (m *TheoryModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the theory tools
func (m *TheoryModel) Update(msg tea.Msg) (*TheoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the theory tools screen
func (m *TheoryModel) View() string {
	return "Theory tools not implemented yet"
}
