package ai

import "fmt"

// PromptRegistry provides centralized prompt templates for all kyanite apps.
// Each app has specific prompt constructors that produce the system + user
// messages for the LLM. The actual model inference goes through the shared
// LLMClient pointing at the NUCBox.

// ── Focus Prompts ──────────────────────────────────────────────────

// FocusParsePrompt builds a prompt for parsing natural language into a task.
func FocusParsePrompt(input string) string {
	return `You are a task management assistant. Parse the following input into a structured task.
Respond with JSON only:
{
  "title": "string",
  "priority": "low|medium|high",
  "due_date": "YYYY-MM-DD or null",
  "tags": ["tag1", "tag2"],
  "notes": "string or empty"
}

Input: ` + input
}

// FocusSuggestPrompt builds a prompt for suggesting next actions.
func FocusSuggestPrompt(existingTasks string) string {
	return `You are a productivity assistant. Given the current task list, suggest 3-5 specific next actions.
Be concise and actionable. Respond with a JSON array of strings.

Current tasks:
` + existingTasks
}

// FocusSummarizePrompt builds a prompt for summarizing tasks.
func FocusSummarizePrompt(tasks string) string {
	return `Summarize the following tasks in 2-3 sentences. Focus on progress and blockers.

Tasks:
` + tasks
}

// FocusChatPrompt builds a system prompt for the focus chat assistant.
func FocusChatSystemPrompt() string {
	return `You are a friendly productivity assistant for focus.sh, a terminal-based task manager.
Help the user with task management, productivity tips, and planning. Be concise and actionable.
You have access to the user's task list for context.`
}

// ── Noise Prompts ──────────────────────────────────────────────────

// NoiseBrainstormPrompt builds a prompt for the QuickIdeaAgent.
func NoiseBrainstormPrompt(mode, songContext string) string {
	base := `You are a creative songwriting assistant for noise.sh, a terminal-based music production tool.`
	switch mode {
	case "unstick":
		return base + ` The user is stuck. Provide 3 creative angles to break through writer's block.
` + songContext
	case "spark":
		return base + ` Generate 3 lyric ideas inspired by the theme below. Each should be 2-4 lines.
` + songContext
	case "tweak":
		return base + ` Suggest 3 improvements to the existing lyrics/chords below. Be specific.
` + songContext
	case "check":
		return base + ` Analyze the following lyrics for rhythm, rhyme, and emotional impact. Give brief feedback.
` + songContext
	default:
		return base + ` Help the user with their songwriting. Be creative and specific.
` + songContext
	}
}

// NoiseContextPrompt formats song metadata into a context string.
func NoiseContextPrompt(title, genre, key, mood string) string {
	parts := []string{}
	if title != "" {
		parts = append(parts, "Title: "+title)
	}
	if genre != "" {
		parts = append(parts, "Genre: "+genre)
	}
	if key != "" {
		parts = append(parts, "Key: "+key)
	}
	if mood != "" {
		parts = append(parts, "Mood: "+mood)
	}
	if len(parts) == 0 {
		return "No song context available."
	}
	return "Song context:\n" + joinLines(parts)
}

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

// ── Shared ──────────────────────────────────────────────────────────

func joinLines(parts []string) string {
	result := ""
	for _, p := range parts {
		result += p + "\n"
	}
	return result
}
