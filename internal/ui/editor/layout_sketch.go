package editor

import (
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// SketchLayout represents the layout for sketch mode
type SketchLayout struct {
	width  int
	height int
}

// NewSketchLayout creates a new sketch layout
func NewSketchLayout() *SketchLayout {
	return &SketchLayout{}
}

// SetDimensions sets the dimensions for the layout
func (l *SketchLayout) SetDimensions(width, height int) {
	dimension.Set(&l.width, &l.height, width, height)
}

func (l *SketchLayout) GetDimensions() (int, int) {
	return l.width, l.height
}

// Render renders the sketch mode layout
func (l *SketchLayout) Render(editorContent string, brainstormContent string) string {
	t := theme.GetManager().Current()
	// Sketch mode: Minimal UI with editor + AI panel only
	// 80% editor, 20% AI panel

	editorWidth := l.width * 80 / 100
	aiPanelWidth := l.width - editorWidth - 1 // -1 for divider

	// Create editor pane (80% width)
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary)

	editorPane := editorStyle.Render(editorContent)

	// Create AI panel (20% width)
	aiPanelStyle := lipgloss.NewStyle().
		Width(aiPanelWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Background(t.Background)

	aiPanelTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent).
		Align(lipgloss.Center).
		Width(aiPanelWidth - 2).
		Render("AI Assistant")

	aiPanelContent := lipgloss.NewStyle().
		Width(aiPanelWidth - 2).
		Height(l.height - 4).
		Render(brainstormContent)

	aiPanel := aiPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, aiPanelTitle, "", aiPanelContent),
	)

	// Create divider
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("â”‚")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, aiPanel)
}

// GetBrainstormPanelWidth returns the width of the brainstorm panel in sketch mode
func (l *SketchLayout) GetBrainstormPanelWidth() int {
	return l.width * 20 / 100
}

// GetEditorWidth returns the width of the editor pane in sketch mode
func (l *SketchLayout) GetEditorWidth() int {
	return l.width * 80 / 100
}
