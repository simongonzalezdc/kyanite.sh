package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/ui/editor"
)

// GetSplitPane returns the split pane model for external access
func (m *EditorModel) GetSplitPane() *editor.SplitPaneModel {
	return m.splitPane
}

// EditorModel handles the song editor screen with split-pane layout
type EditorModel struct {
	splitPane *editor.SplitPaneModel
	database  *db.DB
}

// NewEditorModel creates a new editor model with split-pane layout
func NewEditorModel(database *db.DB) *EditorModel {
	return &EditorModel{
		splitPane: editor.NewSplitPaneModel(database),
		database:  database,
	}
}

// Init initializes the editor model
func (m *EditorModel) Init() tea.Cmd {
	return m.splitPane.Init()
}

// Update handles messages for the editor
func (m *EditorModel) Update(msg tea.Msg) (*EditorModel, tea.Cmd) {
	var cmd tea.Cmd
	m.splitPane, cmd = m.splitPane.Update(msg)
	return m, cmd
}

// View renders the editor
func (m *EditorModel) View() string {
	return m.splitPane.View()
}

// GetEditorText returns the current editor text (for compatibility)
func (m *EditorModel) GetEditorText() string {
	return m.splitPane.GetEditorText()
}

// SetEditorText sets the editor text (for compatibility)
func (m *EditorModel) SetEditorText(text string) {
	m.splitPane.SetEditorText(text)
}
