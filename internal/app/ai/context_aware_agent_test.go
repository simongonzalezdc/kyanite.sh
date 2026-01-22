package ai

import (
	"context"
	"testing"
	"time"
)

// mockLLMClient implements QuickLLMClient for testing
type mockLLMClient struct {
	responses  map[string]string
	shouldFail bool
}

func (m *mockLLMClient) Generate(ctx context.Context, _, prompt string, options map[string]any) (string, error) {
	if m.shouldFail {
		return "", nil // Simulate failure with empty response
	}

	// Return different responses based on prompt content
	if containsString(prompt, "lyric") {
		return "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through", nil
	} else if containsString(prompt, "pattern") || containsString(prompt, "chord") {
		return "1. Am - G - C - F\n2. I - V - vi - IV\n3. Dm - G - C - Am", nil
	} else if containsString(prompt, "spark") {
		return "1. In the heart of creativity, I found my way\n2. creativity whispers through the window pane\n3. Chasing creativity through the pouring rain", nil
	} else if containsString(prompt, "tweak") {
		return "1. Rewrite with stronger imagery and emotion\n2. Replace clichés with fresh, specific details\n3. Enhance both lyrical and musical flow", nil
	} else if containsString(prompt, "check") {
		return "STRONG\nAdd vivid sensory details", nil
	}

	// Default response
	return "1. suggestion one\n2. suggestion two\n3. suggestion three", nil
}

func TestQuickIdeaAgent_ContextAwareGeneration(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockLLMClient{shouldFail: false}
	agent = agent.WithClient(mockClient, 5*time.Second)

	tests := []struct {
		name        string
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		expectType  string
		expectCount int
	}{
		{
			name: "Generate lyric continuations",
			mode: QuickIdeaModeUnstick,
			context: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright`,
			options:     map[string]string{},
			expectType:  "lyric",
			expectCount: 3,
		},
		{
			name: "Generate pattern continuations",
			mode: QuickIdeaModeUnstick,
			context: `Verse: C - G - Am - F
Chorus: F - C - G - Am`,
			options:     map[string]string{},
			expectType:  "pattern",
			expectCount: 3,
		},
		{
			name:        "Generate lyric sparks",
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
			expectType:  "lyric",
			expectCount: 3,
		},
		{
			name:        "Generate pattern sparks",
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "rhythm"},
			expectType:  "pattern",
			expectCount: 3,
		},
		{
			name:        "Generate lyric variations",
			mode:        QuickIdeaModeTweak,
			context:     "Love is in the air tonight",
			options:     map[string]string{},
			expectType:  "lyric",
			expectCount: 3,
		},
		{
			name:        "Generate pattern variations",
			mode:        QuickIdeaModeTweak,
			context:     "C - G - Am - F",
			options:     map[string]string{},
			expectType:  "pattern",
			expectCount: 3,
		},
		{
			name:        "Generate lyric quality check",
			mode:        QuickIdeaModeCheck,
			context:     "Love is in the air tonight",
			options:     map[string]string{},
			expectType:  "lyric",
			expectCount: 0, // Quality checks return rating, not suggestions
		},
		{
			name:        "Generate pattern quality check",
			mode:        QuickIdeaModeCheck,
			context:     "C - G - Am - F",
			options:     map[string]string{},
			expectType:  "pattern",
			expectCount: 0, // Quality checks return rating, not suggestions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			resp, err := agent.Generate(context.Background(), req)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}

			if resp == nil {
				t.Error("Generate() returned nil response")
				return
			}

			// Check suggestion count for non-check modes
			if tt.mode != QuickIdeaModeCheck {
				if len(resp.Suggestions) != tt.expectCount {
					t.Errorf("Generate() suggestions count = %d, expected %d", len(resp.Suggestions), tt.expectCount)
				}
			}

			// For check mode, verify rating and tip
			if tt.mode == QuickIdeaModeCheck {
				if resp.Rating == "" {
					t.Error("Generate() check mode should return a rating")
				}
				if resp.Tip == "" {
					t.Error("Generate() check mode should return a tip")
				}
			}

			// Verify response time is set
			if resp.ResponseTime == 0 {
				t.Error("Generate() should set response time")
			}
		})
	}
}

func TestQuickIdeaAgent_ContextAwareFallback(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockLLMClient{shouldFail: true} // Force fallback
	agent = agent.WithClient(mockClient, 5*time.Second)

	tests := []struct {
		name        string
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		expectCount int
	}{
		{
			name:        "Fallback lyric continuations",
			mode:        QuickIdeaModeUnstick,
			context:     "I'm walking down the street",
			options:     map[string]string{},
			expectCount: 3,
		},
		{
			name:        "Fallback pattern continuations",
			mode:        QuickIdeaModeUnstick,
			context:     "C - G - Am - F",
			options:     map[string]string{},
			expectCount: 3,
		},
		{
			name:        "Fallback sparks",
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
			expectCount: 3,
		},
		{
			name:        "Fallback variations",
			mode:        QuickIdeaModeTweak,
			context:     "Love is in the air",
			options:     map[string]string{},
			expectCount: 3,
		},
		{
			name:        "Fallback quality check",
			mode:        QuickIdeaModeCheck,
			context:     "Some content",
			options:     map[string]string{},
			expectCount: 0, // Quality checks don't return suggestions
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			resp, err := agent.Generate(context.Background(), req)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}

			if resp == nil {
				t.Error("Generate() returned nil response")
				return
			}

			// Check suggestion count for non-check modes
			if tt.mode != QuickIdeaModeCheck {
				if len(resp.Suggestions) != tt.expectCount {
					t.Errorf("Generate() suggestions count = %d, expected %d", len(resp.Suggestions), tt.expectCount)
				}
			}

			// For check mode, verify rating and tip
			if tt.mode == QuickIdeaModeCheck {
				if resp.Rating == "" {
					t.Error("Generate() check mode should return a rating")
				}
				if resp.Tip == "" {
					t.Error("Generate() check mode should return a tip")
				}
			}
		})
	}
}

func TestQuickIdeaAgent_MixedContentHandling(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockLLMClient{shouldFail: false}
	agent = agent.WithClient(mockClient, 5*time.Second)

	mixedContent := `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright
Am F C G
You left me here all alone`

	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: mixedContent,
		Options: map[string]string{},
	}

	resp, err := agent.Generate(context.Background(), req)
	if err != nil {
		t.Errorf("Generate() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("Generate() returned nil response")
		return
	}

	if len(resp.Suggestions) == 0 {
		t.Error("Generate() should return suggestions for mixed content")
	}
}

func TestQuickIdeaAgent_ErrorHandling(t *testing.T) {
	agent := NewQuickIdeaAgent()

	tests := []struct {
		name        string
		mode        QuickIdeaMode
		context     string
		expectError bool
	}{
		{
			name:        "Empty mode",
			mode:        "",
			context:     "Some content",
			expectError: true,
		},
		{
			name:        "Valid mode",
			mode:        QuickIdeaModeUnstick,
			context:     "Some content",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: map[string]string{},
			}

			resp, err := agent.Generate(context.Background(), req)

			if tt.expectError {
				if err == nil {
					t.Error("Generate() expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Generate() unexpected error = %v", err)
				}

				if resp == nil {
					t.Error("Generate() should return response even on fallback")
				}
			}
		})
	}
}

func TestQuickIdeaAgent_ContextDetection(t *testing.T) {
	agent := NewQuickIdeaAgent()

	tests := []struct {
		name         string
		content      string
		expectedType ContentType
	}{
		{
			name: "Detect lyric content",
			content: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright`,
			expectedType: ContentTypeLyrics,
		},
		{
			name: "Detect pattern content",
			content: `Verse: C - G - Am - F
Chorus: F - C - G - Am`,
			expectedType: ContentTypePatterns,
		},
		{
			name: "Detect mixed content",
			content: `[Verse]
C G Am F
I'm walking down the street tonight`,
			expectedType: ContentTypeMixed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detectedType := agent.contextDetector.AnalyzeContent(tt.content)
			if detectedType != tt.expectedType {
				t.Errorf("Context detection = %v, expected %v", detectedType, tt.expectedType)
			}
		})
	}
}

func TestQuickIdeaAgent_TimeoutHandling(t *testing.T) {
	agent := NewQuickIdeaAgent()

	// Create a slow client that will timeout
	slowClient := &mockLLMClient{shouldFail: false}
	agent = agent.WithClient(slowClient, 1*time.Millisecond) // Very short timeout

	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Some content",
		Options: map[string]string{},
	}

	// This should not error due to fallback
	resp, err := agent.Generate(context.Background(), req)
	if err != nil {
		t.Errorf("Generate() should not error on timeout due to fallback, got %v", err)
	}

	if resp == nil {
		t.Error("Generate() should return fallback response on timeout")
	}
}
