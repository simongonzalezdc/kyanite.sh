package editor

import (
	"github.com/kyanite/noise/internal/theme"
	"github.com/kyanite/noise/internal/ui/dimension"
	"github.com/charmbracelet/lipgloss"
)

// PolishLayout represents the layout for polish mode
type PolishLayout struct {
	width  int
	height int
}

// NewPolishLayout creates a new polish layout
func NewPolishLayout() *PolishLayout {
	return &PolishLayout{}
}

// SetDimensions sets the dimensions for the layout
func (l *PolishLayout) SetDimensions(width, height int) {
	dimension.Set(&l.width, &l.height, width, height)
}

func (l *PolishLayout) GetDimensions() (int, int) {
	return l.width, l.height
}

// Render renders the polish mode layout
func (l *PolishLayout) Render(editorContent string, previewContent string, theoryContent string, critiqueContent string) string {
	t := theme.GetManager().Current()
	// Poland mode: Full suite with all tools
	// 40% editor, 20% preview, 20% theory, 20% critique

	editorWidthPercent := 40
	previewWidthPercent := 20
	theoryWidthPercent := 20

	editorWidth := l.width * editorWidthPercent / 100
	previewWidth := l.width * previewWidthPercent / 100
	theoryWidth := l.width * theoryWidthPercent / 100
	critiqueWidth := l.width - editorWidth - previewWidth - theoryWidth - 3 // -3 for dividers

	// Create editor pane (40% width) - NO border since editorContent already has one
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height)

	editorPane := editorStyle.Render(editorContent)

	// Create preview pane (20% width) - NO border since previewContent already has one
	previewStyle := lipgloss.NewStyle().
		Width(previewWidth).
		Height(l.height)

	previewPane := previewStyle.Render(previewContent)

	// Create theory pane (20% width) - needs border since it's custom content
	theoryStyle := lipgloss.NewStyle().
		Width(theoryWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Success)

	theoryTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Success).
		Align(lipgloss.Center).
		Width(theoryWidth - 4).
		Render("Theory Tools")

	theoryPaneContent := lipgloss.NewStyle().
		Width(theoryWidth - 4).
		Height(l.height - 6).
		Render(theoryContent)

	theoryPane := theoryStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, theoryTitle, "", theoryPaneContent),
	)

	// Create critique pane (20% width) - needs border since it's custom content
	critiqueStyle := lipgloss.NewStyle().
		Width(critiqueWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Warning)

	critiqueTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Warning).
		Align(lipgloss.Center).
		Width(critiqueWidth - 4).
		Render("AI Critique")

	critiquePaneContent := lipgloss.NewStyle().
		Width(critiqueWidth - 4).
		Height(l.height - 6).
		Render(critiqueContent)

	critiquePane := critiqueStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, critiqueTitle, "", critiquePaneContent),
	)

	// Create dividers
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("|")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, previewPane, dividerStyle, theoryPane, dividerStyle, critiquePane)
}

// GetPreviewPanelWidth returns the width of the preview panel in polish mode
func (l *PolishLayout) GetPreviewPanelWidth() int {
	previewWidthPercent := 20
	return l.width * previewWidthPercent / 100
}

// GetTheoryPanelWidth returns the width of the theory panel in polish mode
func (l *PolishLayout) GetTheoryPanelWidth() int {
	theoryWidthPercent := 20
	return l.width * theoryWidthPercent / 100
}

// GetCritiquePanelWidth returns the width of the critique panel in polish mode
func (l *PolishLayout) GetCritiquePanelWidth() int {
	critiqueWidthPercent := 20
	return l.width * critiqueWidthPercent / 100
}

// GetEditorWidth returns the width of the editor pane in polish mode
func (l *PolishLayout) GetEditorWidth() int {
	editorWidthPercent := 40
	return l.width * editorWidthPercent / 100
}
