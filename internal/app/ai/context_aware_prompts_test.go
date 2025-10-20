package ai

import (
	"testing"
)

func TestContextAwarePrompts_Initialize(t *testing.T) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	// Test that all prompt types are initialized
	testModes := []QuickIdeaMode{
		QuickIdeaModeUnstick,
		QuickIdeaModeSpark,
		QuickIdeaModeTweak,
		QuickIdeaModeCheck,
	}

	for _, mode := range testModes {
		// Test lyric prompts
		if prompt, exists := prompts.lyricPrompts[mode]; !exists || prompt == "" {
			t.Errorf("Lyric prompt for mode %v should be initialized", mode)
		}

		// Test pattern prompts
		if prompt, exists := prompts.patternPrompts[mode]; !exists || prompt == "" {
			t.Errorf("Pattern prompt for mode %v should be initialized", mode)
		}

		// Test mixed prompts
		if prompt, exists := prompts.mixedPrompts[mode]; !exists || prompt == "" {
			t.Errorf("Mixed prompt for mode %v should be initialized", mode)
		}
	}
}

func TestContextAwarePrompts_GetPrompt(t *testing.T) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		expectEmpty bool
	}{
		{
			name:        "Lyrics unstick prompt",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			expectEmpty: false,
		},
		{
			name:        "Patterns spark prompt",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeSpark,
			expectEmpty: false,
		},
		{
			name:        "Mixed tweak prompt",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeTweak,
			expectEmpty: false,
		},
		{
			name:        "Unknown content type",
			contentType: ContentTypeUnknown,
			mode:        QuickIdeaModeCheck,
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := prompts.GetPrompt(tt.contentType, tt.mode)

			if tt.expectEmpty && prompt != "" {
				t.Errorf("Expected empty prompt for %s/%s, got: %s", tt.contentType, tt.mode, prompt)
			}

			if !tt.expectEmpty && prompt == "" {
				t.Errorf("Expected non-empty prompt for %s/%s", tt.contentType, tt.mode)
			}
		})
	}
}

func TestContextAwarePrompts_RenderPrompt(t *testing.T) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		expectEmpty bool
	}{
		{
			name:        "Render lyrics unstick with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "I'm walking down the street",
			options:     map[string]string{},
			expectEmpty: false,
		},
		{
			name:        "Render patterns spark with theme",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
			expectEmpty: false,
		},
		{
			name:        "Render mixed tweak with context",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeTweak,
			context:     "C - G - Am - F progression",
			options:     map[string]string{},
			expectEmpty: false,
		},
		{
			name:        "Render check with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeCheck,
			context:     "Love is in the air tonight",
			options:     map[string]string{},
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := prompts.RenderPrompt(tt.contentType, tt.mode, tt.context, tt.options)

			if tt.expectEmpty && rendered != "" {
				t.Errorf("Expected empty rendered prompt for %s/%s, got: %s", tt.contentType, tt.mode, rendered)
			}

			if !tt.expectEmpty && rendered == "" {
				t.Errorf("Expected non-empty rendered prompt for %s/%s", tt.contentType, tt.mode)
			}

			// Check that the rendered prompt contains the expected context or theme
			if !tt.expectEmpty && rendered != "" {
				if tt.mode == QuickIdeaModeSpark && tt.options["theme"] != "" {
					if !containsString(rendered, tt.options["theme"]) {
						t.Errorf("Rendered prompt should contain theme '%s', got: %s", tt.options["theme"], rendered)
					}
				} else if tt.context != "" {
					if !containsString(rendered, tt.context) {
						t.Errorf("Rendered prompt should contain context '%s', got: %s", tt.context, rendered)
					}
				}
			}
		})
	}
}

func TestContextAwarePrompts_GetContentTypeDescription(t *testing.T) {
	prompts := NewContextAwarePrompts()

	tests := []struct {
		name        string
		contentType ContentType
		expected    string
	}{
		{
			name:        "Lyrics description",
			contentType: ContentTypeLyrics,
			expected:    "Songwriting - focusing on lyrics, imagery, and emotion",
		},
		{
			name:        "Patterns description",
			contentType: ContentTypePatterns,
			expected:    "Musical patterns - focusing on harmony, rhythm, and structure",
		},
		{
			name:        "Mixed description",
			contentType: ContentTypeMixed,
			expected:    "Mixed content - combining lyrical and musical elements",
		},
		{
			name:        "Unknown description",
			contentType: ContentTypeUnknown,
			expected:    "Unknown content type - using general creative assistance",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description := prompts.GetContentTypeDescription(tt.contentType)
			if description != tt.expected {
				t.Errorf("GetContentTypeDescription() = %v, expected %v", description, tt.expected)
			}
		})
	}
}

func TestContextAwarePrompts_GetModeDescription(t *testing.T) {
	prompts := NewContextAwarePrompts()

	tests := []struct {
		name     string
		mode     QuickIdeaMode
		expected string
	}{
		{
			name:     "Unstick description",
			mode:     QuickIdeaModeUnstick,
			expected: "Continue writing - get suggestions to keep moving forward",
		},
		{
			name:     "Spark description",
			mode:     QuickIdeaModeSpark,
			expected: "Generate ideas - create new starting points and concepts",
		},
		{
			name:     "Tweak description",
			mode:     QuickIdeaModeTweak,
			expected: "Refine content - improve and vary existing material",
		},
		{
			name:     "Check description",
			mode:     QuickIdeaModeCheck,
			expected: "Quality check - evaluate and get improvement suggestions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description := prompts.GetModeDescription(tt.mode)
			if description != tt.expected {
				t.Errorf("GetModeDescription() = %v, expected %v", description, tt.expected)
			}
		})
	}
}

func TestContextAwarePrompts_PromptContent(t *testing.T) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	// Test that lyric prompts contain lyric-specific terms
	lyricPrompt := prompts.GetPrompt(ContentTypeLyrics, QuickIdeaModeUnstick)
	if !containsString(lyricPrompt, "lyric") && !containsString(lyricPrompt, "imagery") {
		t.Error("Lyric prompt should contain lyric-specific terms")
	}

	// Test that pattern prompts contain pattern-specific terms
	patternPrompt := prompts.GetPrompt(ContentTypePatterns, QuickIdeaModeUnstick)
	if !containsString(patternPrompt, "musical") && !containsString(patternPrompt, "harmonic") {
		t.Error("Pattern prompt should contain musical/pattern-specific terms")
	}

	// Test that mixed prompts contain both types of terms
	mixedPrompt := prompts.GetPrompt(ContentTypeMixed, QuickIdeaModeUnstick)
	if !containsString(mixedPrompt, "lyric") && !containsString(mixedPrompt, "musical") {
		t.Error("Mixed prompt should contain both lyric and musical terms")
	}
}

func TestContextAwarePrompts_FallbackBehavior(t *testing.T) {
	prompts := NewContextAwarePrompts()
	// Don't initialize to test fallback behavior

	// Test that fallback prompts are returned when not initialized
	fallbackPrompt := prompts.GetPrompt(ContentTypeLyrics, QuickIdeaModeUnstick)
	if fallbackPrompt == "" {
		t.Error("Fallback prompt should be returned when prompts are not initialized")
	}

	// Test that fallback rendering works
	rendered := prompts.RenderPrompt(ContentTypeLyrics, QuickIdeaModeSpark, "", map[string]string{"theme": "test"})
	if rendered == "" {
		t.Error("Fallback rendering should work when prompts are not initialized")
	}
}
