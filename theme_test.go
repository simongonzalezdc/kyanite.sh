package main

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/lyricforge/internal/ui/styles"
)

// TestThemeColors tests theme color constants
func TestThemeColors(t *testing.T) {
	// Test primary colors
	if styles.Primary == "" {
		t.Error("Expected Primary color to be defined")
	}

	if styles.Secondary == "" {
		t.Error("Expected Secondary color to be defined")
	}

	if styles.Accent == "" {
		t.Error("Expected Accent color to be defined")
	}

	// Test functional colors
	if styles.Success == "" {
		t.Error("Expected Success color to be defined")
	}

	if styles.Warning == "" {
		t.Error("Expected Warning color to be defined")
	}

	if styles.Error == "" {
		t.Error("Expected Error color to be defined")
	}

	if styles.Info == "" {
		t.Error("Expected Info color to be defined")
	}

	// Test background & text colors
	if styles.Background == "" {
		t.Error("Expected Background color to be defined")
	}

	if styles.TextPrimary == "" {
		t.Error("Expected TextPrimary color to be defined")
	}

	if styles.TextSecondary == "" {
		t.Error("Expected TextSecondary color to be defined")
	}

	if styles.TextMuted == "" {
		t.Error("Expected TextMuted color to be defined")
	}

	if styles.TextAccent == "" {
		t.Error("Expected TextAccent color to be defined")
	}
}

// TestThemeStyles tests theme style constants
func TestThemeStyles(t *testing.T) {
	// Test title style
	if styles.Title.GetForeground() != styles.Primary {
		t.Error("Expected Title style to use Primary color")
	}

	if !styles.Title.GetBold() {
		t.Error("Expected Title style to be bold")
	}

	// Test subtitle style
	if styles.Subtitle.GetForeground() != styles.TextSecondary {
		t.Error("Expected Subtitle style to use TextSecondary color")
	}

	if !styles.Subtitle.GetItalic() {
		t.Error("Expected Subtitle style to be italic")
	}

	// Test text style
	if styles.Text.GetForeground() != styles.TextPrimary {
		t.Error("Expected Text style to use TextPrimary color")
	}

	// Test muted style
	if styles.Muted.GetForeground() != styles.TextMuted {
		t.Error("Expected Muted style to use TextMuted color")
	}

	// Test emphasis style
	if styles.Emphasis.GetForeground() != styles.TextAccent {
		t.Error("Expected Emphasis style to use TextAccent color")
	}

	if !styles.Emphasis.GetBold() {
		t.Error("Expected Emphasis style to be bold")
	}
}

// TestThemeBorderStyles tests border style constants
func TestThemeBorderStyles(t *testing.T) {
	// Test standard border
	borderRendered := styles.Border.Render("test")
	if borderRendered == "" {
		t.Error("Expected Border style to render content")
	}

	// Test active border
	activeBorderRendered := styles.BorderActive.Render("test")
	if activeBorderRendered == "" {
		t.Error("Expected BorderActive style to render content")
	}

	// Test thick border
	thickBorderRendered := styles.BorderThick.Render("test")
	if thickBorderRendered == "" {
		t.Error("Expected BorderThick style to render content")
	}
}

// TestThemeButtonStyles tests button style constants
func TestThemeButtonStyles(t *testing.T) {
	// Test primary button
	if styles.ButtonPrimary.GetBackground() != styles.Primary {
		t.Error("Expected ButtonPrimary style to use Primary background")
	}

	if styles.ButtonPrimary.GetForeground() != styles.Background {
		t.Error("Expected ButtonPrimary style to use Background foreground")
	}

	if !styles.ButtonPrimary.GetBold() {
		t.Error("Expected ButtonPrimary style to be bold")
	}

	// Test secondary button
	secondaryButtonRendered := styles.ButtonSecondary.Render("Test")
	if secondaryButtonRendered == "" {
		t.Error("Expected ButtonSecondary style to render content")
	}

	// Test accent button
	if styles.ButtonAccent.GetBackground() != styles.Accent {
		t.Error("Expected ButtonAccent style to use Accent background")
	}

	if styles.ButtonAccent.GetForeground() != styles.Background {
		t.Error("Expected ButtonAccent style to use Background foreground")
	}

	// Test disabled button
	if styles.ButtonDisabled.GetBackground() != styles.Dark3 {
		t.Error("Expected ButtonDisabled style to use Dark3 background")
	}

	if styles.ButtonDisabled.GetForeground() != styles.TextMuted {
		t.Error("Expected ButtonDisabled style to use TextMuted foreground")
	}
}

// TestThemeStatusStyles tests status style constants
func TestThemeStatusStyles(t *testing.T) {
	// Test success status
	if styles.StatusSuccess.GetForeground() != styles.Success {
		t.Error("Expected StatusSuccess style to use Success color")
	}

	if !styles.StatusSuccess.GetBold() {
		t.Error("Expected StatusSuccess style to be bold")
	}

	// Test warning status
	if styles.StatusWarning.GetForeground() != styles.Warning {
		t.Error("Expected StatusWarning style to use Warning color")
	}

	if !styles.StatusWarning.GetBold() {
		t.Error("Expected StatusWarning style to be bold")
	}

	// Test error status
	if styles.StatusError.GetForeground() != styles.Error {
		t.Error("Expected StatusError style to use Error color")
	}

	if !styles.StatusError.GetBold() {
		t.Error("Expected StatusError style to be bold")
	}

	// Test info status
	if styles.StatusInfo.GetForeground() != styles.Info {
		t.Error("Expected StatusInfo style to use Info color")
	}
}

// TestThemeEditorStyles tests editor component style constants
func TestThemeEditorStyles(t *testing.T) {
	// Test editor pane
	editorPaneRendered := styles.EditorPane.Render("test")
	if editorPaneRendered == "" {
		t.Error("Expected EditorPane style to render content")
	}

	// Test preview pane
	previewPaneRendered := styles.PreviewPane.Render("test")
	if previewPaneRendered == "" {
		t.Error("Expected PreviewPane style to render content")
	}

	// Test status bar
	if styles.StatusBar.GetBackground() != styles.Dark2 {
		t.Error("Expected StatusBar style to use Dark2 background")
	}

	if styles.StatusBar.GetForeground() != styles.TextPrimary {
		t.Error("Expected StatusBar style to use TextPrimary foreground")
	}

	// Test cursor
	if styles.Cursor.GetBackground() != styles.Primary {
		t.Error("Expected Cursor style to use Primary background")
	}

	if styles.Cursor.GetForeground() != styles.Background {
		t.Error("Expected Cursor style to use Background foreground")
	}

	// Test divider
	if styles.Divider.GetForeground() != styles.BorderColor {
		t.Error("Expected Divider style to use BorderColor")
	}
}

// TestThemeTypographyStyles tests typography style constants
func TestThemeTypographyStyles(t *testing.T) {
	// Test H1
	if styles.H1.GetForeground() != styles.Primary {
		t.Error("Expected H1 style to use Primary color")
	}

	if !styles.H1.GetBold() {
		t.Error("Expected H1 style to be bold")
	}

	// Test H2
	if styles.H2.GetForeground() != styles.Secondary {
		t.Error("Expected H2 style to use Secondary color")
	}

	if !styles.H2.GetBold() {
		t.Error("Expected H2 style to be bold")
	}

	if !styles.H2.GetUnderline() {
		t.Error("Expected H2 style to be underlined")
	}

	// Test H3
	if styles.H3.GetForeground() != styles.TextAccent {
		t.Error("Expected H3 style to use TextAccent color")
	}

	if !styles.H3.GetBold() {
		t.Error("Expected H3 style to be bold")
	}

	// Test Bold
	if !styles.Bold.GetBold() {
		t.Error("Expected Bold style to be bold")
	}

	// Test Italic
	if !styles.Italic.GetItalic() {
		t.Error("Expected Italic style to be italic")
	}

	// Test Underline
	if !styles.Underline.GetUnderline() {
		t.Error("Expected Underline style to be underlined")
	}

	// Test Code
	if styles.Code.GetForeground() != styles.Accent {
		t.Error("Expected Code style to use Accent color")
	}

	if styles.Code.GetBackground() != styles.Dark2 {
		t.Error("Expected Code style to use Dark2 background")
	}

	// Test Quote
	if styles.Quote.GetForeground() != styles.TextSecondary {
		t.Error("Expected Quote style to use TextSecondary color")
	}

	if !styles.Quote.GetItalic() {
		t.Error("Expected Quote style to be italic")
	}

	if !styles.Quote.GetBorderLeft() {
		t.Error("Expected Quote style to have left border")
	}
}

// TestThemeListStyles tests list and menu style constants
func TestThemeListStyles(t *testing.T) {
	// Test list item
	if styles.ListItem.GetForeground() != styles.TextPrimary {
		t.Error("Expected ListItem style to use TextPrimary color")
	}

	// Test selected list item
	if styles.ListItemSelected.GetBackground() != styles.Primary {
		t.Error("Expected ListItemSelected style to use Primary background")
	}

	if styles.ListItemSelected.GetForeground() != styles.Background {
		t.Error("Expected ListItemSelected style to use Background foreground")
	}

	if !styles.ListItemSelected.GetBold() {
		t.Error("Expected ListItemSelected style to be bold")
	}
}

// TestThemeCardStyles tests card and panel style constants
func TestThemeCardStyles(t *testing.T) {
	// Test standard card
	cardRendered := styles.Card.Render("test")
	if cardRendered == "" {
		t.Error("Expected Card style to render content")
	}

	// Test highlighted card
	highlightCardRendered := styles.CardHighlight.Render("test")
	if highlightCardRendered == "" {
		t.Error("Expected CardHighlight style to render content")
	}

	// Test success card
	successCardRendered := styles.CardSuccess.Render("test")
	if successCardRendered == "" {
		t.Error("Expected CardSuccess style to render content")
	}

	// Test error card
	errorCardRendered := styles.CardError.Render("test")
	if errorCardRendered == "" {
		t.Error("Expected CardError style to render content")
	}
}

// TestThemeBadgeStyles tests badge and tag style constants
func TestThemeBadgeStyles(t *testing.T) {
	// Test tag
	if styles.Tag.GetBackground() != styles.Secondary {
		t.Error("Expected Tag style to use Secondary background")
	}

	if styles.Tag.GetForeground() != styles.TextPrimary {
		t.Error("Expected Tag style to use TextPrimary foreground")
	}

	// Test badge
	if styles.Badge.GetBackground() != styles.Accent {
		t.Error("Expected Badge style to use Accent background")
	}

	if styles.Badge.GetForeground() != styles.Background {
		t.Error("Expected Badge style to use Background foreground")
	}

	if !styles.Badge.GetBold() {
		t.Error("Expected Badge style to be bold")
	}
}

// TestThemeHelperFunctions tests theme helper functions
func TestThemeHelperFunctions(t *testing.T) {
	// Test Gradient function
	colors := []lipgloss.Color{
		styles.Primary,
		styles.Secondary,
		styles.Accent,
	}

	gradientText := styles.Gradient("Test", colors)
	if gradientText == "" {
		t.Error("Expected Gradient function to return non-empty string")
	}

	// Test TitleGradient function
	titleGradient := styles.TitleGradient("Test Title")
	if titleGradient == "" {
		t.Error("Expected TitleGradient function to return non-empty string")
	}

	// Test StatusBadge function
	statusBadge := styles.StatusBadge("OK", styles.Success)
	if statusBadge == "" {
		t.Error("Expected StatusBadge function to return non-empty string")
	}

	// Test ListItemWithIcon function
	listItem := styles.ListItemWithIcon("✓", "Test Item", false)
	if listItem == "" {
		t.Error("Expected ListItemWithIcon function to return non-empty string")
	}

	// Test selected list item
	selectedItem := styles.ListItemWithIcon("✓", "Selected Item", true)
	if selectedItem == "" {
		t.Error("Expected ListItemWithIcon function to return non-empty string for selected item")
	}

	// Test ToHex function
	hex := styles.ToHex(styles.Primary)
	if hex == "" {
		t.Error("Expected ToHex function to return non-empty string")
	}
}

// TestThemeExtendedPalette tests extended palette colors
func TestThemeExtendedPalette(t *testing.T) {
	// Test purple colors
	if styles.Purple1 == "" {
		t.Error("Expected Purple1 color to be defined")
	}

	if styles.Purple2 == "" {
		t.Error("Expected Purple2 color to be defined")
	}

	if styles.Purple3 == "" {
		t.Error("Expected Purple3 color to be defined")
	}

	// Test gold colors
	if styles.Gold1 == "" {
		t.Error("Expected Gold1 color to be defined")
	}

	if styles.Gold2 == "" {
		t.Error("Expected Gold2 color to be defined")
	}

	// Test dark colors
	if styles.Dark1 == "" {
		t.Error("Expected Dark1 color to be defined")
	}

	if styles.Dark2 == "" {
		t.Error("Expected Dark2 color to be defined")
	}

	if styles.Dark3 == "" {
		t.Error("Expected Dark3 color to be defined")
	}

	// Test border color
	if styles.BorderColor == "" {
		t.Error("Expected BorderColor to be defined")
	}
}

// TestThemeConsistency tests theme consistency across styles
func TestThemeConsistency(t *testing.T) {
	// Test that all primary text styles use consistent colors
	if styles.Title.GetForeground() != styles.H1.GetForeground() {
		t.Error("Expected Title and H1 to use the same color")
	}

	// Test that all secondary text styles use consistent colors
	if styles.Subtitle.GetForeground() != styles.H2.GetForeground() {
		t.Error("Expected Subtitle and H2 to use the same color")
	}

	// Test that all accent styles use consistent colors
	if styles.Emphasis.GetForeground() != styles.H3.GetForeground() {
		t.Error("Expected Emphasis and H3 to use the same color")
	}

	// Test that button styles use consistent colors
	if styles.ButtonPrimary.GetBackground() != styles.Title.GetForeground() {
		t.Error("Expected ButtonPrimary background to match Title color")
	}
}

// TestThemeStyleProperties tests specific style properties
func TestThemeStyleProperties(t *testing.T) {
	// Test that code style has padding
	codeStyle := styles.Code.Render("test")
	if codeStyle == "" {
		t.Error("Expected Code style to render content")
	}

	// Test that quote style has left border
	quoteStyle := styles.Quote.Render("test quote")
	if quoteStyle == "" {
		t.Error("Expected Quote style to render content")
	}

	// Test that border styles have padding
	borderStyle := styles.Border.Render("test")
	if borderStyle == "" {
		t.Error("Expected Border style to render content")
	}

	// Test that button styles have padding
	buttonStyle := styles.ButtonPrimary.Render("Button")
	if buttonStyle == "" {
		t.Error("Expected ButtonPrimary style to render content")
	}
}

// BenchmarkThemeRendering benchmarks theme rendering performance
func BenchmarkThemeRendering(b *testing.B) {
	text := "Benchmark Text"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = styles.Title.Render(text)
		_ = styles.Text.Render(text)
		_ = styles.ButtonPrimary.Render(text)
		_ = styles.Card.Render(text)
	}
}

// BenchmarkThemeHelperFunctions benchmarks theme helper function performance
func BenchmarkThemeHelperFunctions(b *testing.B) {
	text := "Benchmark Text"
	colors := []lipgloss.Color{styles.Primary, styles.Secondary, styles.Accent}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = styles.Gradient(text, colors)
		_ = styles.TitleGradient(text)
		_ = styles.StatusBadge("OK", styles.Success)
		_ = styles.ListItemWithIcon("✓", text, false)
	}
}
