package ai

import (
	"fmt"
	"strings"
)

// ContextAwarePrompts manages prompts that adapt based on content type
type ContextAwarePrompts struct {
	lyricPrompts   map[QuickIdeaMode]string
	patternPrompts map[QuickIdeaMode]string
	mixedPrompts   map[QuickIdeaMode]string
	templates      map[QuickIdeaMode]string
}

// NewContextAwarePrompts creates a new context-aware prompt manager
func NewContextAwarePrompts() *ContextAwarePrompts {
	return &ContextAwarePrompts{
		lyricPrompts:   make(map[QuickIdeaMode]string),
		patternPrompts: make(map[QuickIdeaMode]string),
		mixedPrompts:   make(map[QuickIdeaMode]string),
		templates:      make(map[QuickIdeaMode]string),
	}
}

// Initialize prompts for different content types
func (cap *ContextAwarePrompts) Initialize() {
	// Lyric-specific prompts
	cap.lyricPrompts[QuickIdeaModeUnstick] = `You help songwriters stay in motion with creative lyric ideas.
Context (recent lines):
%s

Return three concise lyric continuations (8-12 syllables) that match the tone, imagery, and emotional content.
Focus on:
- Maintaining the established mood and theme
- Using sensory details and vivid imagery
- Creating natural-sounding rhymes and rhythm
- Building on the narrative or emotional arc

Format exactly:
1. first continuation
2. second continuation
3. third continuation
`

	cap.lyricPrompts[QuickIdeaModeSpark] = `You generate creative first lines for songs based on themes.
Theme or mood: %s

Return three vivid opening lines (8-12 syllables) with concrete imagery and emotional resonance.
Focus on:
- Establishing a clear mood or atmosphere
- Using specific, sensory details
- Creating intrigue or emotional connection
- Setting up potential narrative directions

Format:
1. first idea
2. second idea
3. third idea
`

	cap.lyricPrompts[QuickIdeaModeTweak] = `You rewrite a single lyric line to make it more powerful and evocative.
Original line:
%s

Provide three alternate phrasings that maintain the core meaning but enhance the emotional impact.
Focus on:
- Stronger imagery and sensory details
- More authentic or specific language
- Better rhythm and flow
- Deeper emotional resonance

Format:
1. variation
2. variation
3. variation
`

	cap.lyricPrompts[QuickIdeaModeCheck] = `You evaluate a lyric line for emotional impact and poetic quality.
Line:
%s

Respond with exactly:
RATING (STRONG, OKAY, or WEAK)
Tip (5 words, actionable for lyric improvement)

Focus on:
- Imagery and sensory details
- Emotional authenticity
- Rhythm and flow
- Originality and specificity

Example:
STRONG
Add vivid sensory image
`

	// Pattern-specific prompts
	cap.patternPrompts[QuickIdeaModeUnstick] = `You help musicians and producers develop musical patterns and progressions.
Context (recent patterns):
%s

Return three musical continuations that build on the established harmonic or rhythmic foundation.
Focus on:
- Harmonic consistency and progression logic
- Rhythmic variation that maintains groove
- Musical tension and release
- Practical playability and musicality

Format exactly:
1. first continuation
2. second continuation
3. third continuation
`

	cap.patternPrompts[QuickIdeaModeSpark] = `You generate creative musical patterns based on a theme or mood.
Theme or mood: %s

Return three distinct musical patterns (chord progressions, rhythmic patterns, or melodic ideas).
Focus on:
- Appropriate harmonic language for the mood
- Clear rhythmic character
- Musical interest and development potential
- Practical implementation

Format:
1. first pattern
2. second pattern
3. third pattern
`

	cap.patternPrompts[QuickIdeaModeTweak] = `You refine and enhance musical patterns for better musical impact.
Original pattern:
%s

Provide three alternate versions that improve the musical effectiveness.
Focus on:
- Stronger harmonic movement
- Better rhythmic interest
- More sophisticated voice leading
- Enhanced musical character

Format:
1. variation
2. variation
3. variation
`

	cap.patternPrompts[QuickIdeaModeCheck] = `You evaluate musical patterns for technical and musical quality.
Pattern:
%s

Respond with exactly:
RATING (STRONG, OKAY, or WEAK)
Tip (5 words, actionable for pattern improvement)

Focus on:
- Harmonic logic and progression
- Rhythmic clarity and interest
- Practical playability
- Musical effectiveness

Example:
STRONG
Strengthen harmonic resolution
`

	// Mixed content prompts
	cap.mixedPrompts[QuickIdeaModeUnstick] = `You help songwriters develop both lyrics and musical patterns.
Context (recent content):
%s

Analyze whether this is primarily lyrical or musical content, then return three appropriate continuations.
For lyrics: focus on imagery, emotion, and narrative
For patterns: focus on harmony, rhythm, and musical development
For mixed: provide suggestions that enhance both elements

Format exactly:
1. first continuation
2. second continuation
3. third continuation
`

	cap.mixedPrompts[QuickIdeaModeSpark] = `You generate creative ideas that can work for both lyrics and musical patterns.
Theme or mood: %s

Return three creative starting points that could inspire either lyrics or music.
Focus on:
- Evocative imagery and mood
- Rhythmic and melodic potential
- Emotional resonance
- Creative flexibility

Format:
1. first idea
2. second idea
3. third idea
`

	cap.mixedPrompts[QuickIdeaModeTweak] = `You refine content that combines lyrical and musical elements.
Original content:
%s

Provide three alternate versions that enhance both the lyrical and musical qualities.
Focus on:
- Better flow and rhythm
- Stronger imagery and emotion
- Improved musicality
- Enhanced creativity

Format:
1. variation
2. variation
3. variation
`

	cap.mixedPrompts[QuickIdeaModeCheck] = `You evaluate content that combines lyrical and musical elements.
Content:
%s

Respond with exactly:
RATING (STRONG, OKAY, or WEAK)
Tip (5 words, actionable for improvement)

Focus on:
- Balance between lyrics and music
- Overall creative effectiveness
- Rhythmic and flow qualities
- Emotional and musical impact

Example:
OKAY
Strengthen musical-rhythmic connection
`
}

// GetPrompt returns the appropriate prompt based on content type and mode
func (cap *ContextAwarePrompts) GetPrompt(contentType ContentType, mode QuickIdeaMode) string {
	switch contentType {
	case ContentTypeLyrics:
		if prompt, exists := cap.lyricPrompts[mode]; exists {
			return prompt
		}
	case ContentTypePatterns:
		if prompt, exists := cap.patternPrompts[mode]; exists {
			return prompt
		}
	case ContentTypeMixed:
		if prompt, exists := cap.mixedPrompts[mode]; exists {
			return prompt
		}
	}

	// Fallback to a generic prompt
	return cap.getFallbackPrompt(mode)
}

// getFallbackPrompt returns a generic prompt when no specific one is available
func (cap *ContextAwarePrompts) getFallbackPrompt(mode QuickIdeaMode) string {
	switch mode {
	case QuickIdeaModeUnstick:
		return `You help creators stay in motion.
Context:
%s

Return three creative continuations that build on the established content.
Format:
1. first continuation
2. second continuation
3. third continuation
`
	case QuickIdeaModeSpark:
		return `You generate creative starting ideas.
Theme: %s

Return three distinct creative ideas.
Format:
1. first idea
2. second idea
3. third idea
`
	case QuickIdeaModeTweak:
		return `You refine and enhance creative content.
Original:
%s

Provide three improved variations.
Format:
1. variation
2. variation
3. variation
`
	case QuickIdeaModeCheck:
		return `You evaluate creative content quality.
Content:
%s

Respond with exactly:
RATING (STRONG, OKAY, or WEAK)
Tip (5 words, actionable)
`
	default:
		return "Provide three creative suggestions."
	}
}

// RenderPrompt renders the prompt with the provided context and options
func (cap *ContextAwarePrompts) RenderPrompt(contentType ContentType, mode QuickIdeaMode, context string, options map[string]string) string {
	prompt := cap.GetPrompt(contentType, mode)

	if prompt == "" {
		return "Provide three creative suggestions."
	}

	switch mode {
	case QuickIdeaModeSpark:
		theme := strings.TrimSpace(options["theme"])
		if theme == "" {
			theme = "creativity"
		}
		return fmt.Sprintf(prompt, theme)
	case QuickIdeaModeTweak, QuickIdeaModeCheck:
		return fmt.Sprintf(prompt, strings.TrimSpace(context))
	default:
		return fmt.Sprintf(prompt, strings.TrimSpace(context))
	}
}

// GetContentTypeDescription returns a human-readable description of the content type
func (cap *ContextAwarePrompts) GetContentTypeDescription(contentType ContentType) string {
	switch contentType {
	case ContentTypeLyrics:
		return "Songwriting - focusing on lyrics, imagery, and emotion"
	case ContentTypePatterns:
		return "Musical patterns - focusing on harmony, rhythm, and structure"
	case ContentTypeMixed:
		return "Mixed content - combining lyrical and musical elements"
	case ContentTypeUnknown:
		return "Unknown content type - using general creative assistance"
	default:
		return "General creative assistance"
	}
}

// GetModeDescription returns a human-readable description of the AI mode
func (cap *ContextAwarePrompts) GetModeDescription(mode QuickIdeaMode) string {
	switch mode {
	case QuickIdeaModeUnstick:
		return "Continue writing - get suggestions to keep moving forward"
	case QuickIdeaModeSpark:
		return "Generate ideas - create new starting points and concepts"
	case QuickIdeaModeTweak:
		return "Refine content - improve and vary existing material"
	case QuickIdeaModeCheck:
		return "Quality check - evaluate and get improvement suggestions"
	default:
		return "AI assistance"
	}
}
