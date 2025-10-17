package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/infra/db"
)

// ManagerModel handles the project manager screen
type ManagerModel struct {
	database *db.DB
	width    int
	height   int
}

// NewManagerModel creates a new manager model
func NewManagerModel(database *db.DB) *ManagerModel {
	return &ManagerModel{
		database: database,
	}
}

// Init initializes the manager model
func (m *ManagerModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the project manager
func (m *ManagerModel) Update(msg tea.Msg) (*ManagerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View renders the project manager screen
func (m *ManagerModel) View() string {
	return "Project manager not implemented yet"
}
