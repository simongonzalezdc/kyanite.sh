package editor

import (
	"context"
	"fmt"
	"strings"

	"github.com/kyanite/noise/internal/app"
	"github.com/kyanite/noise/internal/app/ai"
	"github.com/kyanite/noise/internal/constants"
	"github.com/kyanite/noise/internal/logging"
	"github.com/kyanite/noise/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// EditorAI implementation methods

// StartRapidBrainstorm starts the rapid brainstorm flow with the given theme
func (a *EditorAI) StartRapidBrainstorm(theme string) {
	a.rapidBrainstorm = true
	a.brainstormTheme = theme

	// Use the AI agent to generate brainstorm angles
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback brainstorm angles")
		a.brainstormAngles = []string{
			"Explore " + theme + " through personal memories",
			"Use nature imagery to symbolize " + theme,
			"Focus on sensory details related to " + theme,
		}
		return
	}

	// Create a request for brainstorm angles
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeSpark,
		Context: theme,
		Options: map[string]string{"theme": theme},
	}

	// Call the AI agent with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.AIContextTimeout)
	defer cancel()

	resp, err := a.aiAgent.Generate(ctx, req)
	if err != nil {
		logging.Errorf("Failed to generate brainstorm angles: %v", err)
		// Use fallback angles
		a.brainstormAngles = []string{
			"Explore " + theme + " through personal memories",
			"Use nature imagery to symbolize " + theme,
			"Focus on sensory details related to " + theme,
		}
		return
	}

	if len(resp.Suggestions) > 0 {
		a.brainstormAngles = resp.Suggestions
		logging.Debugf("Generated %d brainstorm angles for theme: %s", len(resp.Suggestions), theme)
	} else {
		// Use fallback angles
		a.brainstormAngles = []string{
			"Explore " + theme + " through personal memories",
			"Use nature imagery to symbolize " + theme,
			"Focus on sensory details related to " + theme,
		}
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

	// Use the AI agent to generate an opening line based on the selected angle
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback opening line")
		a.insertOpeningLine(state, "Opening line for: "+selectedAngle)
		a.clearBrainstormState()
		return
	}

	// Create a request for an opening line
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeSpark,
		Context: selectedAngle,
		Options: map[string]string{"theme": a.brainstormTheme, "angle": selectedAngle},
	}

	// Call the AI agent with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.AIContextTimeout)
	defer cancel()

	resp, err := a.aiAgent.Generate(ctx, req)
	if err != nil {
		logging.Errorf("Failed to generate opening line: %v", err)
		// Use fallback opening line
		a.insertOpeningLine(state, "Opening line for: "+selectedAngle)
	} else if len(resp.Suggestions) > 0 {
		// Use the first suggestion as the opening line
		a.insertOpeningLine(state, resp.Suggestions[0])
		logging.Debugf("Generated opening line: %s", resp.Suggestions[0])
	} else {
		// Use fallback opening line
		a.insertOpeningLine(state, "Opening line for: "+selectedAngle)
	}

	a.clearBrainstormState()
}

// insertOpeningLine inserts an opening line at the appropriate position
func (a *EditorAI) insertOpeningLine(state StateManagerInterface, openingLine string) {
	// Insert the opening line at the current cursor position
	currentContent := state.GetText()
	if currentContent != "" {
		openingLine = "\n" + openingLine
	}
	state.SetText(currentContent + openingLine)
}

// clearBrainstormState clears the brainstorm state
func (a *EditorAI) clearBrainstormState() {
	a.rapidBrainstorm = false
	a.brainstormTheme = ""
	a.brainstormAngles = nil
}

// StartContinueMode starts the continue writing mode
func (a *EditorAI) StartContinueMode() {
	a.continueMode = true

	// Use the AI agent to generate continue suggestions
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback continue suggestions")
		a.continueSuggestions = []string{
			"Continue with this line...",
			"Or try this alternative...",
			"Perhaps this direction...",
		}
		return
	}

	// Get the current content from the state (will be provided in SelectContinueSuggestion)
	// For now, we'll generate suggestions based on an empty context
	// The actual content will be provided when a suggestion is selected
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

// GenerateContinueSuggestions generates continue suggestions based on current content
func (a *EditorAI) GenerateContinueSuggestions(state StateManagerInterface) {
	currentContent := state.GetText()
	if currentContent == "" {
		a.continueSuggestions = []string{
			"Start with a compelling first line...",
			"Begin with a vivid image...",
			"Open with an intriguing question...",
		}
		return
	}

	// Use the AI agent to generate continue suggestions
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback continue suggestions")
		a.continueSuggestions = []string{
			"Continue with this line...",
			"Or try this alternative...",
			"Perhaps this direction...",
		}
		return
	}

	// Create a request for continue suggestions
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeUnstick,
		Context: currentContent,
		Options: map[string]string{},
	}

	// Call the AI agent with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.AIContextTimeout)
	defer cancel()

	resp, err := a.aiAgent.Generate(ctx, req)
	if err != nil {
		logging.Errorf("Failed to generate continue suggestions: %v", err)
		// Use fallback suggestions
		a.continueSuggestions = []string{
			"Continue with this line...",
			"Or try this alternative...",
			"Perhaps this direction...",
		}
		return
	}

	if len(resp.Suggestions) > 0 {
		a.continueSuggestions = resp.Suggestions
		logging.Debugf("Generated %d continue suggestions", len(resp.Suggestions))
	} else {
		// Use fallback suggestions
		a.continueSuggestions = []string{
			"Continue with this line...",
			"Or try this alternative...",
			"Perhaps this direction...",
		}
	}
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

	// Use the AI agent to generate variations
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback variations")
		a.variationOptions = []string{
			"Variation 1 of: " + selectedText,
			"Variation 2 of: " + selectedText,
			"Variation 3 of: " + selectedText,
		}
		return
	}

	// Create a request for variations
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeTweak,
		Context: selectedText,
		Options: map[string]string{},
	}

	// Call the AI agent with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.AIContextTimeout)
	defer cancel()

	resp, err := a.aiAgent.Generate(ctx, req)
	if err != nil {
		logging.Errorf("Failed to generate variations: %v", err)
		// Use fallback variations
		a.variationOptions = []string{
			"Variation 1 of: " + selectedText,
			"Variation 2 of: " + selectedText,
			"Variation 3 of: " + selectedText,
		}
		return
	}

	if len(resp.Suggestions) > 0 {
		a.variationOptions = resp.Suggestions
		logging.Debugf("Generated %d variations", len(resp.Suggestions))
	} else {
		// Use fallback variations
		a.variationOptions = []string{
			"Variation 1 of: " + selectedText,
			"Variation 2 of: " + selectedText,
			"Variation 3 of: " + selectedText,
		}
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
	content := state.GetText()
	if content == "" {
		return
	}

	// Use the AI agent to perform quality check
	if a.aiAgent == nil {
		logging.Warnf("AI agent not initialized, using fallback quality check")
		a.addQualityCheckResult(state, "OKAY", "Add vivid sensory image")
		return
	}

	// Create a request for quality check
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeCheck,
		Context: content,
		Options: map[string]string{},
	}

	// Call the AI agent with timeout
	ctx, cancel := context.WithTimeout(context.Background(), constants.AIContextTimeout)
	defer cancel()

	resp, err := a.aiAgent.Generate(ctx, req)
	if err != nil {
		logging.Errorf("Failed to perform quality check: %v", err)
		// Use fallback quality check
		a.addQualityCheckResult(state, "OKAY", "Add vivid sensory image")
		return
	}

	// Add quality check result
	a.addQualityCheckResult(state, resp.Rating, resp.Tip)
	logging.Debugf("Quality check result: %s - %s", resp.Rating, resp.Tip)
}

// addQualityCheckResult adds the quality check result to the content
func (a *EditorAI) addQualityCheckResult(state StateManagerInterface, rating, tip string) {
	// Create a quality check comment
	qualityComment := fmt.Sprintf("\n\n<!-- Quality Check: %s - %s -->", rating, tip)

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

// SetAIService sets the AI service for this editor component
func (a *EditorAI) SetAIService(service *app.AIService) {
	a.aiService = service
	// Also sync the underlying agent from the service
	if service != nil {
		a.aiAgent = service.GetQuickAgent()
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

	t := theme.GetManager().Current()
	// Create overlay style
	overlayStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Padding(0, 1).
		Width(width - 8).
		Background(t.Background)

	// Create title
	titleStyle := lipgloss.Style{}.
		Bold(true).
		Foreground(t.Primary).
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
	instructionStyle := lipgloss.Style{}.
		Foreground(t.Secondary).
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

	t := theme.GetManager().Current()
	// Create overlay style
	overlayStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Accent).
		Padding(0, 1).
		Width(width - 8).
		Background(t.Background)

	// Create title
	titleStyle := lipgloss.Style{}.
		Bold(true).
		Foreground(t.Accent).
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
	instructionStyle := lipgloss.Style{}.
		Foreground(t.Secondary).
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

	t := theme.GetManager().Current()
	// Create overlay style
	overlayStyle := lipgloss.Style{}.
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Success).
		Padding(0, 1).
		Width(width - 8).
		Background(t.Background)

	// Create title
	titleStyle := lipgloss.Style{}.
		Bold(true).
		Foreground(t.Success).
		Align(lipgloss.Center)

	title := titleStyle.Render("Variations for:")

	// Show original text
	originalStyle := lipgloss.Style{}.
		Foreground(t.Secondary).
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
	instructionStyle := lipgloss.Style{}.
		Foreground(t.Secondary).
		Align(lipgloss.Center).
		Italic(true)

	instructionsText := instructionStyle.Render(instructions)

	// Combine all elements
	content := lipgloss.JoinVertical(lipgloss.Left, title, "", originalText, "", variationsText, "", instructionsText)

	return overlayStyle.Render(content)
}
