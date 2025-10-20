package editor

import (
	"github.com/Kyanite/noise/internal/ui/dimension"
	"github.com/Kyanite/noise/internal/theme"
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
	// Polish mode: Full suite with all tools
	// 40% editor, 20% preview, 20% theory, 20% critique

	editorWidth := l.width * 40 / 100
	previewWidth := l.width * 20 / 100
	theoryWidth := l.width * 20 / 100
	critiqueWidth := l.width - editorWidth - previewWidth - theoryWidth - 3 // -3 for dividers

	// Create editor pane (40% width)
	editorStyle := lipgloss.NewStyle().
		Width(editorWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary)

	editorPane := editorStyle.Render(editorContent)

	// Create preview pane (20% width)
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

	// Create theory pane (20% width)
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

	// Create critique pane (20% width)
	critiqueStyle := lipgloss.NewStyle().
		Width(critiqueWidth).
		Height(l.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Warning).
		Background(t.Background)

	critiqueTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Warning).
		Align(lipgloss.Center).
		Width(critiqueWidth - 2).
		Render("AI Critique")

	critiquePaneContent := lipgloss.NewStyle().
		Width(critiqueWidth - 2).
		Height(l.height - 4).
		Render(critiqueContent)

	critiquePane := critiqueStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left, critiqueTitle, "", critiquePaneContent),
	)

	// Create dividers
	dividerStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Render("â”‚")

	// Combine panes
	return lipgloss.JoinHorizontal(lipgloss.Top, editorPane, dividerStyle, previewPane, dividerStyle, theoryPane, dividerStyle, critiquePane)
}

// GetPreviewPanelWidth returns the width of the preview panel in polish mode
func (l *PolishLayout) GetPreviewPanelWidth() int {
	return l.width * 20 / 100
}

// GetTheoryPanelWidth returns the width of the theory panel in polish mode
func (l *PolishLayout) GetTheoryPanelWidth() int {
	return l.width * 20 / 100
}

// GetCritiquePanelWidth returns the width of the critique panel in polish mode
func (l *PolishLayout) GetCritiquePanelWidth() int {
	return l.width * 20 / 100
}

// GetEditorWidth returns the width of the editor pane in polish mode
func (l *PolishLayout) GetEditorWidth() int {
	return l.width * 40 / 100
}
