package editor

import (
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/dimension"
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

	editorWidthPercent := 80
	editorWidth := l.width * editorWidthPercent / 100
	aiPanelWidth := l.width - editorWidth - 1 // -1 for divider

	// Create editor pane (80% width) - NO border since editorContent already has one
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height)

	editorPane := editorStyle.Render(editorContent)

	// Create AI panel (20% width) with border
	aiPanelStyle := lipgloss.NewStyle().
		Width(aiPanelWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent)

	aiPanelTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent).
		Align(lipgloss.Center).
		Width(aiPanelWidth - 4).
		Render("AI Assistant")

	aiPanelContent := lipgloss.NewStyle().
		Width(aiPanelWidth - 4).
		Height(l.height - 6).
		Render(brainstormContent)

	aiPanel := aiPanelStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, aiPanelTitle, "", aiPanelContent),
	)

	// Create divider
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("|")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, aiPanel)
}

// GetBrainstormPanelWidth returns the width of the brainstorm panel in sketch mode
func (l *SketchLayout) GetBrainstormPanelWidth() int {
	aiPanelWidthPercent := 20
	return l.width * aiPanelWidthPercent / 100
}

// GetEditorWidth returns the width of the editor pane in sketch mode
func (l *SketchLayout) GetEditorWidth() int {
	editorWidthPercent := 80
	return l.width * editorWidthPercent / 100
}
