package ui

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/infra/db"
)

// EditorModel handles the song editor screen
type EditorModel struct {
	textarea textarea.Model
	database *db.DB
	width    int
	height   int
}

// NewEditorModel creates a new editor model
func NewEditorModel(database *db.DB) *EditorModel {
	ta := textarea.New()
	ta.Placeholder = "Start writing your lyrics..."
	ta.Focus()

	return &EditorModel{
		textarea: ta,
		database: database,
	}
}

// Init initializes the editor model
func (m *EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

// Update handles messages for the editor
func (m *EditorModel) Update(msg tea.Msg) (*EditorModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textarea.SetWidth(msg.Width - 4)
		m.textarea.SetHeight(msg.Height - 4)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+s":
			// Save functionality would go here
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

// View renders the editor
func (m *EditorModel) View() string {
	return m.textarea.View()
}
