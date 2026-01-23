package editor

import (
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/charmbracelet/lipgloss"
)

// DraftLayout represents the layout for draft mode
type DraftLayout struct {
	width  int
	height int
}

// NewDraftLayout creates a new draft layout
func NewDraftLayout() *DraftLayout {
	return &DraftLayout{}
}

// SetDimensions sets the dimensions for the layout
func (l *DraftLayout) SetDimensions(width, height int) {
	dimension.Set(&l.width, &l.height, width, height)
}

func (l *DraftLayout) GetDimensions() (int, int) {
	return l.width, l.height
}

// Render renders the draft mode layout
func (l *DraftLayout) Render(editorContent string, previewContent string, theoryContent string) string {
	t := theme.GetManager().Current()
	// Draft mode: Editor + Preview + Theory tools
	// 50% editor, 25% preview, 25% theory

	editorWidthPercent := 50
	previewWidthPercent := 25

	editorWidth := l.width * editorWidthPercent / 100
	previewWidth := l.width * previewWidthPercent / 100
	theoryWidth := l.width - editorWidth - previewWidth - 2 // -2 for dividers

	// Create editor pane (50% width) - NO border since editorContent already has one
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height)

	editorPane := editorStyle.Render(editorContent)

	// Create preview pane (25% width) - NO border since previewContent already has one
	previewStyle := lipgloss.NewStyle().
		Width(previewWidth).
		Height(l.height)

	previewPane := previewStyle.Render(previewContent)

	// Create theory pane (25% width) - this needs a border since it's custom content
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

	// Create dividers
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("|")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, previewPane, dividerStyle, theoryPane)
}

// GetPreviewPanelWidth returns the width of the preview panel in draft mode
func (l *DraftLayout) GetPreviewPanelWidth() int {
	previewWidthPercent := 25
	return l.width * previewWidthPercent / 100
}

// GetTheoryPanelWidth returns the width of the theory panel in draft mode
func (l *DraftLayout) GetTheoryPanelWidth() int {
	theoryWidthPercent := 25
	return l.width * theoryWidthPercent / 100
}

// GetEditorWidth returns the width of the editor pane in draft mode
func (l *DraftLayout) GetEditorWidth() int {
	editorWidthPercent := 50
	return l.width * editorWidthPercent / 100
}
