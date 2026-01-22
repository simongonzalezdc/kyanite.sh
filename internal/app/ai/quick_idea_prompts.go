package ai

import (
	"fmt"
	"strings"
)

// quickIdeaPrompts manages rendering and parsing of the lightweight QuickIdeaAgent prompts.
type quickIdeaPrompts struct {
	templates map[QuickIdeaMode]string
}

// defaultQuickIdeaPrompts configures prompt templates for all supported modes.
func defaultQuickIdeaPrompts() quickIdeaPrompts {
	return quickIdeaPrompts{
		templates: map[QuickIdeaMode]string{
			QuickIdeaModeUnstick: `You help songwriters stay in motion.
Context (recent lines):
%s

Return three concise continuations (8-12 syllables) matching the tone and imagery.
Format exactly:
1. first continuation
2. second continuation
3. third continuation
`,
			QuickIdeaModeSpark: `You generate creative first lines for songs.
Theme or mood: %s

Return three vivid opening lines (8-12 syllables) with concrete imagery.
Format:
1. first idea
2. second idea
3. third idea
`,
			QuickIdeaModeTweak: `You rewrite a single lyric line.
Original line:
%s

Provide three alternate phrasings with similar meaning, bolder imagery.
Format:
1. variation
2. variation
3. variation
`,
			QuickIdeaModeCheck: `You evaluate a lyric line briefly.
Original Line: %s

[GUARDRAILS]
- Do not invent quality metrics; use only the provided scale.
- Tips must be 5 words maximum.
- No conversational filler.

Respond with exactly:
RATING (STRONG, OKAY, or WEAK)
Tip (5 words, actionable)
`,
			QuickIdeaModeHarmony: `You are a music theory assistant. Generate chord progressions for a mood.
Mood: %s

[GUARDRAILS]
- Use only standard chord symbols (e.g., Cmaj7, Am9, Gsus4).
- Do not hallucinate non-existent chord extensions.
- Group chords into logical 4-bar or 8-bar progressions.

Return three distinct progressions.
Format:
1. Chord - Chord - Chord - Chord
2. Chord - Chord - Chord - Chord
3. Chord - Chord - Chord - Chord
`,
		},
	}
}

// render builds the final prompt string based on mode, context, and options.
func (p quickIdeaPrompts) render(mode QuickIdeaMode, context string, options map[string]string) string {
	template, ok := p.templates[mode]
	if !ok {
		return "Provide three concise lyric suggestions."
	}

	switch mode {
	case QuickIdeaModeSpark, QuickIdeaModeHarmony:
		return fmt.Sprintf(template, strings.TrimSpace(options["theme"]))
	case QuickIdeaModeTweak, QuickIdeaModeCheck:
		return fmt.Sprintf(template, strings.TrimSpace(context))
	case QuickIdeaModeUnstick:
		style := options["style"]
		if style != "" {
			return fmt.Sprintf("Style: %s\n", style) + fmt.Sprintf(template, strings.TrimSpace(context))
		}
		return fmt.Sprintf(template, strings.TrimSpace(context))
	default:
		return fmt.Sprintf(template, strings.TrimSpace(context))
	}
}

// parse interprets the model response into a structured QuickResponse.
func (p quickIdeaPrompts) parse(mode QuickIdeaMode, raw string) *QuickResponse {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	switch mode {
	case QuickIdeaModeCheck:
		return parseCheckResponse(raw)
	case QuickIdeaModeUnstick, QuickIdeaModeSpark, QuickIdeaModeTweak, QuickIdeaModeHarmony:
		return parseNumberedSuggestions(raw)
	default:
		return nil
	}
}

func parseNumberedSuggestions(raw string) *QuickResponse {
	lines := strings.Split(raw, "\n")
	suggestions := make([]string, 0, 3)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 3 {
			continue
		}

		if strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") || strings.HasPrefix(line, "3.") {
			parts := strings.SplitN(line, ".", 2)
			if len(parts) == 2 {
				suggestion := strings.TrimSpace(parts[1])
				if suggestion != "" {
					suggestions = append(suggestions, suggestion)
				}
			}
		} else if len(suggestions) < 3 && line != "" {
			suggestions = append(suggestions, line)
		}
	}

	if len(suggestions) == 0 {
		return nil
	}

	return &QuickResponse{Suggestions: suggestions}
}

func parseCheckResponse(raw string) *QuickResponse {
	lines := strings.Split(raw, "\n")

	var rating string
	var tip string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if rating == "" && isCheckRating(line) {
			rating = strings.ToUpper(line)
			continue
		}

		if tip == "" {
			tip = line
			break
		}
	}

	if rating == "" {
		return nil
	}

	return &QuickResponse{
		Rating: rating,
		Tip:    tip,
	}
}

func isCheckRating(value string) bool {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case "STRONG", "OKAY", "WEAK":
		return true
	default:
		return false
	}
}
