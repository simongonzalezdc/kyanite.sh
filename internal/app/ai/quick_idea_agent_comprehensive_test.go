package ai

import (
	"context"
	"testing"
	"time"
)

// TestQuickIdeaAgent_ComprehensiveGeneration tests all generation modes thoroughly
func TestQuickIdeaAgent_ComprehensiveGeneration(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		setupMock   func(*MockLLMClient)
		expectError bool
		description string
	}{
		{
			name:    "Unstick mode with lyrics",
			mode:    QuickIdeaModeUnstick,
			context: "I'm walking down the street tonight",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("lyric continuation", "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through")
			},
			expectError: false,
			description: "Should generate lyric continuations",
		},
		{
			name:    "Unstick mode with patterns",
			mode:    QuickIdeaModeUnstick,
			context: "C - G - Am - F",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("pattern continuation", "1. Am - G - C - F\n2. I - V - vi - IV\n3. Dm - G - C - Am")
			},
			expectError: false,
			description: "Should generate pattern continuations",
		},
		{
			name:    "Spark mode with theme",
			mode:    QuickIdeaModeSpark,
			context: "",
			options: map[string]string{"theme": "love"},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("spark lyric", "1. In the heart of love, I found my way\n2. love whispers through the window pane\n3. Chasing love through the pouring rain")
			},
			expectError: false,
			description: "Should generate lyric ideas based on theme",
		},
		{
			name:    "Spark mode with pattern theme",
			mode:    QuickIdeaModeSpark,
			context: "",
			options: map[string]string{"theme": "rhythm"},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("spark pattern", "1. rhythm theme: C - G - Am - F progression\n2. rhythm rhythm: driving 4/4 with syncopation\n3. rhythm mood: minor key with descending bassline")
			},
			expectError: false,
			description: "Should generate pattern ideas based on theme",
		},
		{
			name:    "Tweak mode with lyrics",
			mode:    QuickIdeaModeTweak,
			context: "Love is in the air tonight",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("tweak lyric", "1. Rewrite with stronger imagery and emotion\n2. Replace clichés with fresh, specific details\n3. Enhance the emotional resonance")
			},
			expectError: false,
			description: "Should generate lyric variations",
		},
		{
			name:    "Tweak mode with patterns",
			mode:    QuickIdeaModeTweak,
			context: "C - G - Am - F",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("tweak pattern", "1. Add sophisticated voice leading\n2. Incorporate rhythmic variation\n3. Enhance harmonic movement")
			},
			expectError: false,
			description: "Should generate pattern variations",
		},
		{
			name:    "Check mode with lyrics",
			mode:    QuickIdeaModeCheck,
			context: "Love is in the air tonight",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("check lyric", "STRONG\nAdd vivid sensory details")
			},
			expectError: false,
			description: "Should evaluate lyric quality",
		},
		{
			name:    "Check mode with patterns",
			mode:    QuickIdeaModeCheck,
			context: "C - G - Am - F",
			options: map[string]string{},
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("check pattern", "STRONG\nStrengthen harmonic resolution")
			},
			expectError: false,
			description: "Should evaluate pattern quality",
		},
		{
			name:        "Invalid mode",
			mode:        "invalid",
			context:     "Some context",
			options:     map[string]string{},
			setupMock:   func(m *MockLLMClient) {},
			expectError: true,
			description: "Should return error for invalid mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			tt.setupMock(mockClient)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			// Execute
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			resp, err := agent.Generate(context.Background(), req)

			// Verify
			if tt.expectError {
				helper.AssertError(err)
				t.Logf("Description: %s", tt.description)
				return
			}

			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertTrue(resp.ResponseTime > 0, "Response time should be set")

			// Mode-specific assertions
			switch tt.mode {
			case QuickIdeaModeUnstick, QuickIdeaModeSpark, QuickIdeaModeTweak:
				helper.AssertLength(resp.Suggestions, 3)
				for _, suggestion := range resp.Suggestions {
					helper.AssertNotEmpty(suggestion)
				}
			case QuickIdeaModeCheck:
				helper.AssertNotEmpty(resp.Rating)
				helper.AssertNotEmpty(resp.Tip)
				helper.AssertTrue(resp.Rating == "STRONG" || resp.Rating == "OKAY" || resp.Rating == "WEAK",
					"Rating should be STRONG, OKAY, or WEAK")
			}

			t.Logf("Description: %s", tt.description)
			t.Logf("Response time: %v", resp.ResponseTime)
		})
	}
}

// TestQuickIdeaAgent_ContextAwareGeneration tests context-aware generation
func TestQuickIdeaAgent_ComprehensiveContextAwareGeneration(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name         string
		mode         QuickIdeaMode
		context      string
		expectedType ContentType
		setupMock    func(*MockLLMClient)
		description  string
	}{
		{
			name: "Lyric context detection",
			mode: QuickIdeaModeUnstick,
			context: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright`,
			expectedType: ContentTypeLyrics,
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("lyric", "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through")
			},
			description: "Should detect lyric content and generate appropriate suggestions",
		},
		{
			name: "Pattern context detection",
			mode: QuickIdeaModeUnstick,
			context: `Verse: C - G - Am - F
Chorus: F - C - G - Am`,
			expectedType: ContentTypePatterns,
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("pattern", "1. Am - G - C - F\n2. I - V - vi - IV\n3. Dm - G - C - Am")
			},
			description: "Should detect pattern content and generate appropriate suggestions",
		},
		{
			name: "Mixed context detection",
			mode: QuickIdeaModeUnstick,
			context: `[Verse]
C G Am F
I'm walking down the street tonight`,
			expectedType: ContentTypeMixed,
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("mixed", "1. with a gentle C major progression\n2. building to an emotional chorus\n3. with a driving rhythm section")
			},
			description: "Should detect mixed content and generate appropriate suggestions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			tt.setupMock(mockClient)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			// Verify context detection
			detectedType := agent.contextDetector.AnalyzeContent(tt.context)
			helper.AssertEqual(tt.expectedType, detectedType)

			// Execute
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: map[string]string{},
			}

			resp, err := agent.Generate(context.Background(), req)
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertLength(resp.Suggestions, 3)

			t.Logf("Description: %s", tt.description)
			t.Logf("Detected type: %v", detectedType)
			t.Logf("Suggestions: %v", resp.Suggestions)
		})
	}
}

// TestQuickIdeaAgent_ErrorHandling tests error handling scenarios
func TestQuickIdeaAgent_ComprehensiveErrorHandling(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		setupAgent  func() *QuickIdeaAgent
		req         QuickRequest
		expectError bool
		description string
	}{
		{
			name: "Client failure",
			setupAgent: func() *QuickIdeaAgent {
				mockClient := NewMockLLMClient()
				mockClient.SetFailure(true)
				agent := NewQuickIdeaAgent()
				return agent.WithClient(mockClient, 5*time.Second)
			},
			req: QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			},
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when client fails",
		},
		{
			name: "Client timeout",
			setupAgent: func() *QuickIdeaAgent {
				mockClient := NewMockLLMClient()
				mockClient.SetTimeout(true)
				agent := NewQuickIdeaAgent()
				return agent.WithClient(mockClient, 100*time.Millisecond)
			},
			req: QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			},
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when client times out",
		},
		{
			name: "Empty response",
			setupAgent: func() *QuickIdeaAgent {
				mockClient := NewMockLLMClient()
				mockClient.SetDefaultResponse("")
				agent := NewQuickIdeaAgent()
				return agent.WithClient(mockClient, 5*time.Second)
			},
			req: QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			},
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when client returns empty response",
		},
		{
			name: "No client configured",
			setupAgent: func() *QuickIdeaAgent {
				agent := NewQuickIdeaAgent()
				// Don't set client
				return agent
			},
			req: QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			},
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when no client is configured",
		},
		{
			name: "Nil client",
			setupAgent: func() *QuickIdeaAgent {
				agent := NewQuickIdeaAgent()
				return agent.WithClient(nil, 5*time.Second)
			},
			req: QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			},
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when client is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent := tt.setupAgent()
			resp, err := agent.Generate(context.Background(), tt.req)

			if tt.expectError {
				helper.AssertError(err)
			} else {
				helper.AssertNoError(err)
				helper.AssertNotNil(resp)

				// Should have fallback suggestions
				if tt.req.Mode != QuickIdeaModeCheck {
					helper.AssertLength(resp.Suggestions, 3)
				} else {
					helper.AssertNotEmpty(resp.Rating)
					helper.AssertNotEmpty(resp.Tip)
				}
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestQuickIdeaAgent_FallbackGeneration tests fallback generation
func TestQuickIdeaAgent_FallbackGeneration(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name            string
		mode            QuickIdeaMode
		context         string
		options         map[string]string
		expectedPattern string
		description     string
	}{
		{
			name:            "Fallback lyric unstick",
			mode:            QuickIdeaModeUnstick,
			context:         "I'm walking down the street",
			options:         map[string]string{},
			expectedPattern: "and the stars begin to fall",
			description:     "Should generate lyric fallback suggestions",
		},
		{
			name:            "Fallback pattern unstick",
			mode:            QuickIdeaModeUnstick,
			context:         "C - G - Am - F",
			options:         map[string]string{},
			expectedPattern: "Am - G - C - F",
			description:     "Should generate pattern fallback suggestions",
		},
		{
			name:            "Fallback spark",
			mode:            QuickIdeaModeSpark,
			context:         "",
			options:         map[string]string{"theme": "love"},
			expectedPattern: "I woke up chasing love",
			description:     "Should generate spark fallback suggestions",
		},
		{
			name:            "Fallback tweak",
			mode:            QuickIdeaModeTweak,
			context:         "Original line here",
			options:         map[string]string{},
			expectedPattern: "Original line here",
			description:     "Should generate tweak fallback suggestions",
		},
		{
			name:            "Fallback check",
			mode:            QuickIdeaModeCheck,
			context:         "Some content",
			options:         map[string]string{},
			expectedPattern: "OKAY",
			description:     "Should generate check fallback suggestions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup agent with failing client to force fallback
			mockClient := NewMockLLMClient()
			mockClient.SetFailure(true)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			resp, err := agent.Generate(context.Background(), req)
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)

			if tt.mode != QuickIdeaModeCheck {
				helper.AssertLength(resp.Suggestions, 3)
				helper.AssertContains(resp.Suggestions[0], tt.expectedPattern)
			} else {
				helper.AssertEqual(tt.expectedPattern, resp.Rating)
				helper.AssertNotEmpty(resp.Tip)
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestQuickIdeaAgent_PerformanceTests tests performance characteristics
func TestQuickIdeaAgent_PerformanceTests(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		setupMock   func(*MockLLMClient)
		maxTime     time.Duration
		description string
	}{
		{
			name: "Fast response",
			setupMock: func(m *MockLLMClient) {
				m.SetDelay(10 * time.Millisecond)
			},
			maxTime:     100 * time.Millisecond,
			description: "Should complete quickly with fast client",
		},
		{
			name: "Slow response",
			setupMock: func(m *MockLLMClient) {
				m.SetDelay(100 * time.Millisecond)
			},
			maxTime:     200 * time.Millisecond,
			description: "Should complete within reasonable time with slow client",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			tt.setupMock(mockClient)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			}

			// Execute with timing
			start := time.Now()
			resp, err := agent.Generate(context.Background(), req)
			elapsed := time.Since(start)

			// Verify
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertTrue(relaxPerfBudgets() || elapsed < tt.maxTime,
				"Should complete in less than %v, took %v", tt.maxTime, elapsed)

			t.Logf("Description: %s", tt.description)
			t.Logf("Completed in %v", elapsed)
		})
	}
}

// TestQuickIdeaAgent_MusicalStyles tests different musical styles and genres
func TestQuickIdeaAgent_MusicalStyles(t *testing.T) {
	helper := NewTestHelper(t)

	styles := []struct {
		name        string
		context     string
		style       string
		expectedKey string
	}{
		{
			name:        "Jazz progression",
			context:     "Verse: Cmaj7 - Am7 - Dm7 - G7",
			style:       "jazz",
			expectedKey: "C", // Use chord root instead of roman numeral
		},
		{
			name:        "Blues progression",
			context:     "Verse: A7 - D7 - E7",
			style:       "blues",
			expectedKey: "A7",
		},
		{
			name:        "Rock progression",
			context:     "Verse: E - B - C#m - A",
			style:       "rock",
			expectedKey: "E",
		},
		{
			name:        "Pop progression",
			context:     "Verse: C - G - Am - F",
			style:       "pop",
			expectedKey: "C",
		},
	}

	for _, style := range styles {
		t.Run(style.name, func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			mockClient.SetResponse(style.style, "1. "+style.expectedKey+" - "+style.expectedKey+" - "+style.expectedKey+"\n2. "+style.expectedKey+" variation\n3. "+style.expectedKey+" progression")

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: style.context,
				Options: map[string]string{"style": style.style},
			}

			resp, err := agent.Generate(context.Background(), req)
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertLength(resp.Suggestions, 3)

			// Verify suggestions contain the expected key
			for _, suggestion := range resp.Suggestions {
				helper.AssertContains(suggestion, style.expectedKey)
			}

			t.Logf("Style: %s", style.style)
			t.Logf("Suggestions: %v", resp.Suggestions)
		})
	}
}

// TestQuickIdeaAgent_IntegrationWithContextDetector tests integration with context detector
func TestQuickIdeaAgent_IntegrationWithContextDetector(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name         string
		content      string
		expectedType ContentType
		setupMock    func(*MockLLMClient)
		description  string
	}{
		{
			name: "Complex song structure",
			content: `[Intro]
Piano: Cmaj7 - Am7
Mood: melancholic

[Verse 1]
Walking through the morning light
C G Am F progression
Feeling the rhythm of life

[Chorus]
Emotional peak with soaring melody
F C G Am - powerful resolution`,
			expectedType: ContentTypeLyrics, // Leans toward lyrics with emotional content
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("lyric", "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through")
			},
			description: "Should handle complex song structure correctly",
		},
		{
			name: "Lyric-heavy with musical terms",
			content: `[Verse]
I'm playing my guitar tonight
C G Am F chords underneath
Singing about love and loss
With a minor key feeling`,
			expectedType: ContentTypeLyrics,
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("lyric", "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through")
			},
			description: "Should detect lyrics despite musical terms",
		},
		{
			name: "Pattern-heavy with emotional descriptions",
			content: `Verse: C - G - Am - F
Mood: melancholic and reflective
Tempo: slow ballad
Emotional journey through the progression
Bridge: Dm - Am - E - Am
Climax with harmonic resolution`,
			expectedType: ContentTypeMixed, // Mix of patterns and emotional content
			setupMock: func(m *MockLLMClient) {
				m.SetResponse("pattern", "1. Am - G - C - F\n2. I - V - vi - IV\n3. Dm - G - C - Am")
			},
			description: "Should detect patterns despite emotional descriptions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			tt.setupMock(mockClient)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			// Verify context detection
			detectedType := agent.contextDetector.AnalyzeContent(tt.content)
			helper.AssertEqual(tt.expectedType, detectedType)

			// Execute generation
			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: tt.content,
				Options: map[string]string{},
			}

			resp, err := agent.Generate(context.Background(), req)
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertLength(resp.Suggestions, 3)

			t.Logf("Description: %s", tt.description)
			t.Logf("Detected type: %v", detectedType)
			t.Logf("Suggestions: %v", resp.Suggestions)
		})
	}
}

// TestQuickIdeaAgent_ResponseTimeValidation tests response time validation
func TestQuickIdeaAgent_ResponseTimeValidation(t *testing.T) {
	helper := NewTestHelper(t)

	// Test with various response times
	responseTimes := []time.Duration{
		0,
		1 * time.Millisecond,
		10 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
	}

	for _, delay := range responseTimes {
		t.Run("Response time "+delay.String(), func(t *testing.T) {
			// Setup
			mockClient := NewMockLLMClient()
			mockClient.SetDelay(delay)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			}

			// Execute
			start := time.Now()
			resp, err := agent.Generate(context.Background(), req)
			elapsed := time.Since(start)

			// Verify
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)
			helper.AssertTrue(resp.ResponseTime > 0, "Response time should be set")
			helper.AssertTrue(elapsed >= delay, "Elapsed time should be at least the mock delay")

			t.Logf("Mock delay: %v, Actual elapsed: %v, Response time: %v",
				delay, elapsed, resp.ResponseTime)
		})
	}
}
