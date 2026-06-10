package ai

import "fmt"

// PromptRegistry provides centralized prompt templates for all kyanite apps.
// Each app has specific prompt constructors that produce the system + user
// messages for the LLM. The actual model inference goes through the shared
// LLMClient pointing at the NUCBox.

// ── Syntax Prompts ──────────────────────────────────────────────────

// SyntaxSuggestPrompt builds a prompt for writing suggestions.
func SyntaxSuggestPrompt(suggestionType, content, context string) string {
	base := `You are a creative writing assistant for syntax.sh, a terminal-based Markdown editor.`
	switch suggestionType {
	case "continue":
		return base + ` Continue the following story naturally. Match the style and tone. Write 2-3 paragraphs.\n\n` + content
	case "improve":
		return base + ` Improve the following passage. Enhance clarity, flow, and impact. Keep the original meaning.\n\n` + content
	case "dialogue":
		return base + ` Generate natural dialogue for the scene below. Make characters distinctive.\n\nContext: ` + context + `\n\n` + content
	case "description":
		return base + ` Add vivid sensory description to the following passage. Engage multiple senses.\n\n` + content
	case "character":
		return base + ` Suggest character behaviors and reactions for the scene below.\n\nContext: ` + context + `\n\n` + content
	default:
		return base + ` Help improve the following writing.\n\n` + content
	}
}

// ── Prism Prompts ───────────────────────────────────────────────────

// PrismPalettePrompt builds a prompt for generating color palettes from text.
func PrismPalettePrompt(description string) string {
	return `You are a color design assistant for prism.sh, a terminal-based color palette tool.
Generate a 5-color palette from the following description.
Respond with JSON only: {"colors": ["#hex1", "#hex2", "#hex3", "#hex4", "#hex5"], "name": "palette name"}

Description: ` + description
}

// PrismContrastPrompt builds a prompt for suggesting accessible color adjustments.
func PrismContrastPrompt(fg, bg string, ratio float64) string {
	return fmt.Sprintf(`The colors %s (foreground) and %s (background) have a contrast ratio of %.1f:1.
This %s WCAG AA requirements (minimum 4.5:1 for normal text).
Suggest an adjusted foreground color that meets WCAG AA while staying close to the original hue.
Respond with just the hex color code.`, fg, bg, ratio,
		map[bool]string{true: "meets", false: "does not meet"}[ratio >= 4.5])
}
