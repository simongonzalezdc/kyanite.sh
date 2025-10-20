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

	editorWidth := l.width * 50 / 100
	previewWidth := l.width * 25 / 100
	theoryWidth := l.width - editorWidth - previewWidth - 2 // -2 for dividers

	// Create editor pane (50% width)
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary)

	editorPane := editorStyle.Render(editorContent)

	// Create preview pane (25% width)
	previewStyle := lipgloss.NewStyle().
		Width(previewWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Background(t.Background)

	previewTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Accent).
		Align(lipgloss.Center).
		Width(previewWidth - 2).
		Render("Preview")

	previewPaneContent := lipgloss.NewStyle().
		Width(previewWidth - 2).
		Height(l.height - 4).
		Render(previewContent)

	previewPane := previewStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, previewTitle, "", previewPaneContent),
	)

	// Create theory pane (25% width)
	theoryStyle := lipgloss.NewStyle().
		Width(theoryWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Success).
		Background(t.Background)

	theoryTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Success).
		Align(lipgloss.Center).
		Width(theoryWidth - 2).
		Render("Theory Tools")

	theoryPaneContent := lipgloss.NewStyle().
		Width(theoryWidth - 2).
		Height(l.height - 4).
		Render(theoryContent)

	theoryPane := theoryStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, theoryTitle, "", theoryPaneContent),
	)

	// Create dividers
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("â”‚")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, previewPane, dividerStyle, theoryPane)
}

// GetPreviewPanelWidth returns the width of the preview panel in draft mode
func (l *DraftLayout) GetPreviewPanelWidth() int {
	return l.width * 25 / 100
}

// GetTheoryPanelWidth returns the width of the theory panel in draft mode
func (l *DraftLayout) GetTheoryPanelWidth() int {
	return l.width * 25 / 100
}

// GetEditorWidth returns the width of the editor pane in draft mode
func (l *DraftLayout) GetEditorWidth() int {
	return l.width * 50 / 100
}
