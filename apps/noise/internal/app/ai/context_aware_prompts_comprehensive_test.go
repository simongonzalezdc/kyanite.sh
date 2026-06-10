package ai

import (
	"testing"
)

// TestContextAwarePrompts_ComprehensiveInitialization tests thorough initialization
func TestContextAwarePrompts_ComprehensiveInitialization(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		setupFunc   func(*ContextAwarePrompts)
		description string
	}{
		{
			name: "Default initialization",
			setupFunc: func(cap *ContextAwarePrompts) {
				cap.Initialize()
			},
			description: "Should initialize all prompt types correctly",
		},
		{
			name: "Multiple initializations",
			setupFunc: func(cap *ContextAwarePrompts) {
				cap.Initialize()
				cap.Initialize() // Should not cause issues
			},
			description: "Should handle multiple initializations gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompts := NewContextAwarePrompts()
			tt.setupFunc(prompts)

			// Test that all prompt types are initialized
			testModes := []QuickIdeaMode{
				QuickIdeaModeUnstick,
				QuickIdeaModeSpark,
				QuickIdeaModeTweak,
				QuickIdeaModeCheck,
			}

			for _, mode := range testModes {
				// Test lyric prompts
				lyricPrompt := prompts.GetPrompt(ContentTypeLyrics, mode)
				helper.AssertNotEmpty(lyricPrompt)

				// Test pattern prompts
				patternPrompt := prompts.GetPrompt(ContentTypePatterns, mode)
				helper.AssertNotEmpty(patternPrompt)

				// Test mixed prompts
				mixedPrompt := prompts.GetPrompt(ContentTypeMixed, mode)
				helper.AssertNotEmpty(mixedPrompt)
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextAwarePrompts_PromptContentValidation tests prompt content validation
func TestContextAwarePrompts_PromptContentValidation(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		expectTerms []string
		description string
	}{
		{
			name:        "Lyric unstick prompt",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"lyric", "continuations", "imagery", "rhyme", "rhythm"},
			description: "Should contain lyric-specific terms",
		},
		{
			name:        "Pattern unstick prompt",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"musical", "harmonic", "rhythmic", "progression"},
			description: "Should contain pattern-specific terms",
		},
		{
			name:        "Mixed unstick prompt",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"lyric", "musical", "continuations"},
			description: "Should contain both lyric and musical terms",
		},
		{
			name:        "Lyric spark prompt",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			expectTerms: []string{"opening", "imagery", "emotional", "narrative"},
			description: "Should contain creative opening terms",
		},
		{
			name:        "Pattern spark prompt",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeSpark,
			expectTerms: []string{"patterns", "harmonic", "rhythmic", "musical"},
			description: "Should contain musical pattern terms",
		},
		{
			name:        "Lyric tweak prompt",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeTweak,
			expectTerms: []string{"rewrite", "imagery", "emotional", "rhythm"},
			description: "Should contain lyric improvement terms",
		},
		{
			name:        "Pattern tweak prompt",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeTweak,
			expectTerms: []string{"refine", "harmonic", "rhythmic", "musical"},
			description: "Should contain pattern improvement terms",
		},
		{
			name:        "Lyric check prompt",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeCheck,
			expectTerms: []string{"evaluate", "poetic", "emotional", "specificity"},
			description: "Should contain lyric evaluation terms",
		},
		{
			name:        "Pattern check prompt",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeCheck,
			expectTerms: []string{"evaluate", "harmonic", "playability", "effective"},
			description: "Should contain pattern evaluation terms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := prompts.GetPrompt(tt.contentType, tt.mode)
			helper.AssertNotEmpty(prompt)

			// Check for expected terms
			for _, term := range tt.expectTerms {
				helper.AssertContains(prompt, term)
			}

			t.Logf("Description: %s", tt.description)
			t.Logf("Prompt contains expected terms: %v", tt.expectTerms)
		})
	}
}

// TestContextAwarePrompts_RenderingVariations tests different rendering scenarios
func TestContextAwarePrompts_RenderingVariations(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		description string
	}{
		{
			name:        "Render lyric unstick with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "I'm walking down the street tonight",
			options:     map[string]string{},
			description: "Should render context correctly for lyrics",
		},
		{
			name:        "Render pattern unstick with context",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
			context:     "C - G - Am - F progression",
			options:     map[string]string{},
			description: "Should render context correctly for patterns",
		},
		{
			name:        "Render spark with theme",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
			description: "Should render theme correctly for spark mode",
		},
		{
			name:        "Render spark without theme",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{},
			description: "Should use default theme when none provided",
		},
		{
			name:        "Render tweak with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeTweak,
			context:     "Love is in the air tonight",
			options:     map[string]string{},
			description: "Should render context correctly for tweak mode",
		},
		{
			name:        "Render check with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeCheck,
			context:     "The rain falls on my window",
			options:     map[string]string{},
			description: "Should render context correctly for check mode",
		},
		{
			name:        "Render with empty context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "",
			options:     map[string]string{},
			description: "Should handle empty context gracefully",
		},
		{
			name:        "Render with special characters",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "Special chars: !@#$%^&*()",
			options:     map[string]string{},
			description: "Should handle special characters in context",
		},
		{
			name:        "Render with unicode",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "Unicode: ♫ ♪ ♫ ♪ 音乐之歌",
			options:     map[string]string{},
			description: "Should handle unicode in context",
		},
		{
			name:        "Render with long context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "This is a very long context that should still be rendered correctly without any issues or problems",
			options:     map[string]string{},
			description: "Should handle long context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := prompts.RenderPrompt(tt.contentType, tt.mode, tt.context, tt.options)
			helper.AssertNotEmpty(rendered)

			// Verify context or theme is included
			if tt.mode == QuickIdeaModeSpark {
				theme := tt.options["theme"]
				if theme == "" {
					theme = "creativity" // Default theme
				}
				helper.AssertContains(rendered, theme)
			} else if tt.context != "" {
				helper.AssertContains(rendered, tt.context)
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextAwarePrompts_FallbackBehavior tests fallback behavior
func TestContextAwarePrompts_ComprehensiveFallbackBehavior(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		setupFunc   func(*ContextAwarePrompts)
		contentType ContentType
		mode        QuickIdeaMode
		description string
	}{
		{
			name: "Fallback without initialization",
			setupFunc: func(cap *ContextAwarePrompts) {
				// Don't initialize
			},
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			description: "Should provide fallback prompts when not initialized",
		},
		{
			name: "Fallback for unknown content type",
			setupFunc: func(cap *ContextAwarePrompts) {
				cap.Initialize()
			},
			contentType: ContentTypeUnknown,
			mode:        QuickIdeaModeUnstick,
			description: "Should provide fallback for unknown content type",
		},
		{
			name: "Fallback for unknown mode",
			setupFunc: func(cap *ContextAwarePrompts) {
				cap.Initialize()
			},
			contentType: ContentTypeLyrics,
			mode:        "unknown",
			description: "Should provide fallback for unknown mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompts := NewContextAwarePrompts()
			tt.setupFunc(prompts)

			// Test prompt retrieval
			prompt := prompts.GetPrompt(tt.contentType, tt.mode)
			helper.AssertNotEmpty(prompt)

			// Test rendering
			rendered := prompts.RenderPrompt(tt.contentType, tt.mode, "test context", map[string]string{})
			helper.AssertNotEmpty(rendered)

			t.Logf("Description: %s", tt.description)
			t.Logf("Fallback prompt: %s", prompt)
		})
	}
}

// TestContextAwarePrompts_DescriptionMethods tests description methods
func TestContextAwarePrompts_DescriptionMethods(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()

	tests := []struct {
		name        string
		contentType ContentType
		expected    string
		description string
	}{
		{
			name:        "Lyrics description",
			contentType: ContentTypeLyrics,
			expected:    "Songwriting - focusing on lyrics, imagery, and emotion",
			description: "Should return accurate description for lyrics",
		},
		{
			name:        "Patterns description",
			contentType: ContentTypePatterns,
			expected:    "Musical patterns - focusing on harmony, rhythm, and structure",
			description: "Should return accurate description for patterns",
		},
		{
			name:        "Mixed description",
			contentType: ContentTypeMixed,
			expected:    "Mixed content - combining lyrical and musical elements",
			description: "Should return accurate description for mixed content",
		},
		{
			name:        "Unknown description",
			contentType: ContentTypeUnknown,
			expected:    "Unknown content type - using general creative assistance",
			description: "Should return appropriate description for unknown content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			description := prompts.GetContentTypeDescription(tt.contentType)
			helper.AssertEqual(tt.expected, description)

			t.Logf("Description: %s", tt.description)
		})
	}

	// Test mode descriptions
	modeTests := []struct {
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

	for _, tt := range modeTests {
		t.Run(tt.name, func(t *testing.T) {
			description := prompts.GetModeDescription(tt.mode)
			helper.AssertEqual(tt.expected, description)
		})
	}
}

// TestContextAwarePrompts_PromptStructure tests prompt structure and format
func TestContextAwarePrompts_PromptStructure(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name         string
		contentType  ContentType
		mode         QuickIdeaMode
		expectFormat string
		description  string
	}{
		{
			name:         "Lyric unstick format",
			contentType:  ContentTypeLyrics,
			mode:         QuickIdeaModeUnstick,
			expectFormat: "1.",
			description:  "Should use numbered format for suggestions",
		},
		{
			name:         "Pattern unstick format",
			contentType:  ContentTypePatterns,
			mode:         QuickIdeaModeUnstick,
			expectFormat: "1.",
			description:  "Should use numbered format for pattern suggestions",
		},
		{
			name:         "Spark format",
			contentType:  ContentTypeLyrics,
			mode:         QuickIdeaModeSpark,
			expectFormat: "1.",
			description:  "Should use numbered format for ideas",
		},
		{
			name:         "Tweak format",
			contentType:  ContentTypeLyrics,
			mode:         QuickIdeaModeTweak,
			expectFormat: "1.",
			description:  "Should use numbered format for variations",
		},
		{
			name:         "Check format",
			contentType:  ContentTypeLyrics,
			mode:         QuickIdeaModeCheck,
			expectFormat: "RATING",
			description:  "Should use specific format for quality check",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := prompts.GetPrompt(tt.contentType, tt.mode)
			helper.AssertNotEmpty(prompt)
			helper.AssertContains(prompt, tt.expectFormat)

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextAwarePrompts_ContextSpecificTerms tests context-specific terminology
func TestContextAwarePrompts_ContextSpecificTerms(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		expectTerms []string
		avoidTerms  []string
		description string
	}{
		{
			name:        "Lyric-specific terms",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"lyric", "continuations", "imagery", "emotional", "narrative"},
			avoidTerms:  []string{"chord", "harmonic", "progression"},
			description: "Should contain lyric terms and avoid pattern terms",
		},
		{
			name:        "Pattern-specific terms",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"musical", "harmonic", "rhythmic", "progression"},
			avoidTerms:  []string{"lyric", "imagery", "emotional"},
			description: "Should contain pattern terms and avoid lyric terms",
		},
		{
			name:        "Mixed content terms",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeUnstick,
			expectTerms: []string{"lyric", "musical", "continuations"},
			avoidTerms:  []string{},
			description: "Should contain both lyric and musical terms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := prompts.GetPrompt(tt.contentType, tt.mode)
			helper.AssertNotEmpty(prompt)

			// Check for expected terms
			for _, term := range tt.expectTerms {
				helper.AssertContains(prompt, term)
			}

			// Check that avoided terms are not present
			for _, term := range tt.avoidTerms {
				helper.AssertNotContains(prompt, term)
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextAwarePrompts_Consistency tests prompt consistency
func TestContextAwarePrompts_Consistency(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		description string
	}{
		{
			name:        "Lyric prompt consistency",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			description: "Should return consistent prompts",
		},
		{
			name:        "Pattern prompt consistency",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
			description: "Should return consistent prompts",
		},
		{
			name:        "Mixed prompt consistency",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeUnstick,
			description: "Should return consistent prompts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get multiple prompts
			prompt1 := prompts.GetPrompt(tt.contentType, tt.mode)
			prompt2 := prompts.GetPrompt(tt.contentType, tt.mode)
			prompt3 := prompts.GetPrompt(tt.contentType, tt.mode)

			// Should be identical
			helper.AssertEqual(prompt1, prompt2)
			helper.AssertEqual(prompt2, prompt3)

			// Test rendering consistency
			rendered1 := prompts.RenderPrompt(tt.contentType, tt.mode, "test context", map[string]string{})
			rendered2 := prompts.RenderPrompt(tt.contentType, tt.mode, "test context", map[string]string{})
			rendered3 := prompts.RenderPrompt(tt.contentType, tt.mode, "test context", map[string]string{})

			helper.AssertEqual(rendered1, rendered2)
			helper.AssertEqual(rendered2, rendered3)

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestContextAwarePrompts_EdgeCases tests edge cases and error conditions
func TestContextAwarePrompts_EdgeCases(t *testing.T) {
	helper := NewTestHelper(t)

	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		description string
	}{
		{
			name:        "Empty context and options",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "",
			options:     map[string]string{},
			description: "Should handle empty context and options",
		},
		{
			name:        "Nil options",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     nil,
			description: "Should handle nil options",
		},
		{
			name:        "Very long context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     string(make([]byte, 10000)), // 10KB of context
			options:     map[string]string{},
			description: "Should handle very long context",
		},
		{
			name:        "Context with newlines",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "Line 1\nLine 2\nLine 3",
			options:     map[string]string{},
			description: "Should handle context with newlines",
		},
		{
			name:        "Options with special characters",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "Special chars: !@#$%^&*()"},
			description: "Should handle options with special characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := prompts.RenderPrompt(tt.contentType, tt.mode, tt.context, tt.options)
			helper.AssertNotEmpty(rendered)

			t.Logf("Description: %s", tt.description)
		})
	}
}
