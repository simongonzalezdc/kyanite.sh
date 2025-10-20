package editor

import (
	"context"
	"fmt"
	"strings"

	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

// EditorAI implementation methods

// StartRapidBrainstorm starts the rapid brainstorm flow with the given theme
func (a *EditorAI) StartRapidBrainstorm(theme string) {
	a.rapidBrainstorm = true
	a.brainstormTheme = theme

	// For now, we'll use placeholder angles
	// In a full implementation, this would call the AI service
	a.brainstormAngles = []string{
		"Explore " + theme + " through personal memories",
		"Use nature imagery to symbolize " + theme,
		"Focus on sensory details related to " + theme,
	}
}

// GetBrainstormAngles returns the current brainstorm angles
func (a *EditorAI) GetBrainstormAngles() []string {
	return a.brainstormAngles
}

// CancelBrainstormMode cancels the current brainstorm flow.
func (a *EditorAI) CancelBrainstormMode() {
	a.rapidBrainstorm = false
	a.brainstormTheme = ""
	a.brainstormAngles = nil
}

// SelectBrainstormAngle selects a brainstorm angle and generates an opening line
func (a *EditorAI) SelectBrainstormAngle(index int, state StateManagerInterface) {
	if index < 0 || index >= len(a.brainstormAngles) {
		return
	}

	selectedAngle := a.brainstormAngles[index]

	// For now, we'll use a placeholder opening line
	// In a full implementation, this would call the AI service
	openingLine := "Opening line for: " + selectedAngle

	// Insert the opening line at the current cursor position
	currentContent := state.GetText()
	if currentContent != "" {
		openingLine = "\n" + openingLine
	}
	state.SetText(currentContent + openingLine)

	// Clear brainstorm state
	a.rapidBrainstorm = false
	a.brainstormTheme = ""
	a.brainstormAngles = nil
}

// StartContinueMode starts the continue writing mode
func (a *EditorAI) StartContinueMode() {
	a.continueMode = true

	// For now, we'll use placeholder suggestions
	// In a full implementation, this would call the AI service
	a.continueSuggestions = []string{
		"Continue with this line...",
		"Or try this alternative...",
		"Perhaps this direction...",
	}
}

// GetContinueSuggestions returns the current continue suggestions
func (a *EditorAI) GetContinueSuggestions() []string {
	return a.continueSuggestions
}

// SelectContinueSuggestion selects a continue suggestion
func (a *EditorAI) SelectContinueSuggestion(index int, state StateManagerInterface) {
	if index < 0 || index >= len(a.continueSuggestions) {
		return
	}

	selectedLine := a.continueSuggestions[index]

	// Insert the selected line at the current cursor position
	currentContent := state.GetText()
	if currentContent != "" {
		selectedLine = "\n" + selectedLine
	}
	state.SetText(currentContent + selectedLine)

	// Clear continue mode
	a.continueMode = false
	a.continueSuggestions = nil
}

// CancelContinueMode cancels the continue writing mode
func (a *EditorAI) CancelContinueMode() {
	a.continueMode = false
	a.continueSuggestions = nil
}

// StartVariationMode starts the variation mode for the selected text
func (a *EditorAI) StartVariationMode(selectedText string) {
	a.variationMode = true
	a.variationOriginal = selectedText

	// For now, we'll use placeholder variations
	// In a full implementation, this would call the AI service
	a.variationOptions = []string{
		"Variation 1 of: " + selectedText,
		"Variation 2 of: " + selectedText,
		"Variation 3 of: " + selectedText,
	}
}

// GetVariationOptions returns the current variation options
func (a *EditorAI) GetVariationOptions() []string {
	return a.variationOptions
}

// SelectVariation selects a variation option
func (a *EditorAI) SelectVariation(index int, state StateManagerInterface) {
	if index < 0 || index >= len(a.variationOptions) {
		return
	}

	selectedVariation := a.variationOptions[index]

	// Replace the original text with the selected variation
	// For now, we'll just append it
	// In a full implementation, this would replace the selected text
	currentContent := state.GetText()
	if currentContent != "" {
		selectedVariation = "\n" + selectedVariation
	}
	state.SetText(currentContent + selectedVariation)

	// Clear variation mode
	a.variationMode = false
	a.variationOriginal = ""
	a.variationOptions = nil
}

// CancelVariationMode cancels the variation mode
func (a *EditorAI) CancelVariationMode() {
	a.variationMode = false
	a.variationOriginal = ""
	a.variationOptions = nil
}

// PerformQualityCheck performs a quality check on the current content
func (a *EditorAI) PerformQualityCheck(state StateManagerInterface) {
	// For now, we'll use a placeholder implementation
	// In a full implementation, this would call the AI service
	content := state.GetText()
	if content == "" {
		return
	}

	// Placeholder quality check results
	// In a full implementation, this would use the QuickIdeaAgent
	qualityRating := "OKAY"
	qualityTip := "Add vivid sensory image"

	// Create a simple overlay to show the quality check result
	// For now, we'll just add it as a comment
	qualityComment := fmt.Sprintf("\n\n<!-- Quality Check: %s - %s -->", qualityRating, qualityTip)

	// Get current content
	currentContent := state.GetText()

	// Add quality check result to content
	newContent := currentContent + qualityComment
	state.SetText(newContent)
}

// UpdateKnowledgeBaseStatus updates the knowledge base status in the status bar
func (a *EditorAI) UpdateKnowledgeBaseStatus(metrics *EditorMetrics) {
	if a.aiAgent != nil && metrics != nil {
		ctx := context.Background()
		available := a.aiAgent.IsKnowledgeBaseAvailable(ctx)
		status := a.aiAgent.GetKnowledgeBaseStatus(ctx)

		if status != nil {
			statusText := "KB: "
			if available {
				statusText += "Stub"
				if status.CardCount > 0 {
					statusText += fmt.Sprintf(" (%d cards)", status.CardCount)
				}
			} else {
				statusText += "Unavailable"
			}

			metrics.UpdateKnowledgeBaseStatus(available, statusText)
		} else {
			metrics.UpdateKnowledgeBaseStatus(false, "KB: Error")
		}
	}
}

// SetAIAgent sets the AI agent for this editor component
func (a *EditorAI) SetAIAgent(agent *ai.QuickIdeaAgent) {
	a.aiAgent = agent
}

// GetAIAgent returns the current AI agent
func (a *EditorAI) GetAIAgent() *ai.QuickIdeaAgent {
	return a.aiAgent
}

// AnalyzeContentType analyzes the content type
func (a *EditorAI) AnalyzeContentType(content string) string {
	if a.contextDetector != nil {
		contentType := a.contextDetector.AnalyzeContent(content)
		contentTypeStr := string(contentType)

		// Only update if content type has changed
		if contentTypeStr != a.lastContentType {
			a.lastContentType = contentTypeStr
			return contentTypeStr
		}
	}
	return a.lastContentType
}

// RenderOverlays renders AI-related overlays
func (a *EditorAI) RenderOverlays(width int) string {
	if a.rapidBrainstorm {
		return a.renderBrainstormOverlay(width)
	} else if a.continueMode {
		return a.renderContinueOverlay(width)
	} else if a.variationMode {
		return a.renderVariationOverlay(width)
	}
	return ""
}

// IsRapidBrainstorm returns whether rapid brainstorm mode is active
func (a *EditorAI) IsRapidBrainstorm() bool {
	return a.rapidBrainstorm
}

// IsContinueMode returns whether continue mode is active
func (a *EditorAI) IsContinueMode() bool {
	return a.continueMode
}

// IsVariationMode returns whether variation mode is active
func (a *EditorAI) IsVariationMode() bool {
	return a.variationMode
}

// GetBrainstormTheme returns the current brainstorm theme
func (a *EditorAI) GetBrainstormTheme() string {
	return a.brainstormTheme
}

// GetVariationOriginal returns the original text for variation
func (a *EditorAI) GetVariationOriginal() string {
	return a.variationOriginal
}

// Private helper methods

// renderBrainstormOverlay renders the brainstorm overlay UI
func (a *EditorAI) renderBrainstormOverlay(width int) string {
	if !a.rapidBrainstorm || len(a.brainstormAngles) == 0 {
		return ""
	}

	// Create overlay style
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Accent).
		Padding(0, 1).
		Width(width - 8).
		Background(styles.Dark2)

	// Create title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Primary).
		Align(lipgloss.Center)

	title := titleStyle.Render("Theme: " + a.brainstormTheme)

	// Create angle options
	var angleOptions []string
	for i, angle := range a.brainstormAngles {
		angleOption := fmt.Sprintf("[%d] %s", i+1, angle)
		angleOptions = append(angleOptions, angleOption)
	}

	anglesText := strings.Join(angleOptions, "\n")

	// Create instructions
	instructions := "Press 1-3 to select, Esc to cancel"
	instructionStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Italic(true)

	instructionsText := instructionStyle.Render(instructions)

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", anglesText, "", instructionsText)

	return overlayStyle.Render(content)
}

// renderContinueOverlay renders the continue writing overlay UI
func (a *EditorAI) renderContinueOverlay(width int) string {
	if !a.continueMode || len(a.continueSuggestions) == 0 {
		return ""
	}

	// Create overlay style
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Info).
		Padding(0, 1).
		Width(width - 8).
		Background(styles.Dark2)

	// Create title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Info).
		Align(lipgloss.Center)

	title := titleStyle.Render("Continue with:")

	// Create suggestion options
	var suggestionOptions []string
	for i, suggestion := range a.continueSuggestions {
		suggestionOption := fmt.Sprintf("[%d] %s", i+1, suggestion)
		suggestionOptions = append(suggestionOptions, suggestionOption)
	}

	suggestionsText := strings.Join(suggestionOptions, "\n")

	// Create instructions
	instructions := "Press 1-3 to select, Esc to write manually"
	instructionStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Italic(true)

	instructionsText := instructionStyle.Render(instructions)

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", suggestionsText, "", instructionsText)

	return overlayStyle.Render(content)
}

// renderVariationOverlay renders the variation overlay UI
func (a *EditorAI) renderVariationOverlay(width int) string {
	if !a.variationMode || len(a.variationOptions) == 0 {
		return ""
	}

	// Create overlay style
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Success).
		Padding(0, 1).
		Width(width - 8).
		Background(styles.Dark2)

	// Create title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(styles.Success).
		Align(lipgloss.Center)

	title := titleStyle.Render("Variations for:")

	// Show original text
	originalStyle := lipgloss.NewStyle().
		Foreground(styles.TextSecondary).
		Italic(true)

	originalText := originalStyle.Render(a.variationOriginal)

	// Create variation options
	var variationOptions []string
	for i, variation := range a.variationOptions {
		variationOption := fmt.Sprintf("[%d] %s", i+1, variation)
		variationOptions = append(variationOptions, variationOption)
	}

	variationsText := strings.Join(variationOptions, "\n")

	// Create instructions
	instructions := "Press 1-3 to replace, Enter to keep original, Esc to cancel"
	instructionStyle := lipgloss.NewStyle().
		Foreground(styles.TextMuted).
		Align(lipgloss.Center).
		Italic(true)

	instructionsText := instructionStyle.Render(instructions)

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", originalText, "", variationsText, "", instructionsText)

	return overlayStyle.Render(content)
}
