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

// ── Focus Prompts ───────────────────────────────────────────────────

// FocusDailyBriefingPrompt builds a prompt for daily task briefing.
func FocusDailyBriefingPrompt(taskList string, completed, due, overdue int) string {
	return fmt.Sprintf(`You are a productivity assistant for focus.sh, an ADHD-friendly task manager.
Generate a concise daily briefing (3-5 lines).
Tasks: %s
Completed today: %d
Due today: %d
Overdue: %d
Highlight what needs attention. Be direct and encouraging.`, taskList, completed, due, overdue)
}

// FocusSmartSuggestPrompt builds a prompt for smart task suggestions.
func FocusSmartSuggestPrompt(context string) string {
	return `You are a productivity assistant for focus.sh. Based on the user's recent activity, suggest one specific actionable improvement. Be concise (2-3 sentences max).

Recent activity: ` + context
}

// FocusNLInputPrompt builds a prompt for natural language task queries.
func FocusNLInputPrompt(query, context string) string {
	return fmt.Sprintf(`You are a helpful productivity assistant for focus.sh. Answer concisely.
Context: %s

User question: %s`, context, query)
}

// ── Noise Prompts ───────────────────────────────────────────────────

// NoiseLyricContinuationPrompt builds a prompt for lyric continuation.
func NoiseLyricContinuationPrompt(lyrics string) string {
	return `You are a songwriting assistant for noise.sh. Continue these lyrics with 4 more lines that match the style, mood, and rhythm. Do not repeat existing lines.

Current lyrics:
` + lyrics
}

// NoiseChordSuggestPrompt builds a prompt for chord suggestions.
func NoiseChordSuggestPrompt(progression string) string {
	return `You are a music theory assistant for noise.sh. Suggest 4 chords that complement this progression. Explain briefly why they work.

Progression: ` + progression
}

// NoiseMoodBoardPrompt builds a prompt for mood-based song suggestions.
func NoiseMoodBoardPrompt(mood string) string {
	return fmt.Sprintf(`You are a creative music assistant for noise.sh.
For a song described as "%s", suggest:
- Key signature
- Tempo (BPM)
- 3 reference songs
- A 4-chord progression
- 2 lyric fragments that match the mood
Be creative and specific.`, mood)
}

// ── Syntax Prompts (additional) ──────────────────────────────────────

// SyntaxContinueWritingPrompt builds a prompt for story continuation.
func SyntaxContinueWritingPrompt(text, storyContext string) string {
	return fmt.Sprintf(`You are a creative writing assistant for syntax.sh. Continue this story naturally.
Match the tone, style, and character voices. Write 3 paragraphs.

Story context: %s

Current text:
%s`, storyContext, text)
}

// SyntaxConsistencyCheckPrompt builds a prompt for consistency checking.
func SyntaxConsistencyCheckPrompt(text, storyContext string) string {
	return fmt.Sprintf(`You are a story editor for syntax.sh. Review this passage for internal consistency.
Check character names, locations, timeline. Flag contradictions.

Story context: %s

Current passage:
%s`, storyContext, text)
}

// SyntaxCharacterVoicePrompt builds a prompt for character voice rewriting.
func SyntaxCharacterVoicePrompt(character, dialogue, storyContext string) string {
	return fmt.Sprintf(`You are a dialogue editor for syntax.sh. Rewrite this dialogue in %s's voice based on their established speech patterns.

Story context: %s

Original dialogue:
%s`, character, storyContext, dialogue)
}

// ── Prism Prompts (additional) ──────────────────────────────────────

// PrismMoodPalettePrompt builds a prompt for mood-based palette generation.
func PrismMoodPalettePrompt(mood string) string {
	return fmt.Sprintf(`You are a color design assistant for prism.sh.
Create 3 different 5-color palettes for the mood: "%s".
Style each differently. Name each palette.
Respond with JSON: {"palettes": [{"name": "...", "colors": ["#hex1",...]}]}`, mood)
}

// PrismA11yCheckPrompt builds a prompt for accessibility analysis.
func PrismA11yCheckPrompt(colors string) string {
	return `You are an accessibility expert for prism.sh. Analyze this palette for WCAG 2.1 AA compliance.
Flag contrast issues and suggest fixes. Be specific with hex codes.

Colors:
` + colors
}
