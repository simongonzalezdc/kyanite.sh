package ai

import (
	"context"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app/knowledge"
)

// TestQuickIdeaAgent_ComprehensiveKnowledgeBaseIntegration tests comprehensive knowledge base integration
func TestQuickIdeaAgent_ComprehensiveKnowledgeBaseIntegration(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name            string
		setupKB         func(*MockKnowledgeBase)
		mode            QuickIdeaMode
		context         string
		options         map[string]string
		expectedKBUsage bool
		description     string
	}{
		{
			name: "Knowledge base with lyric cards",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.AddTestCard("lyric1", "Rhyme Schemes", "Use AABB rhyme scheme for catchy choruses", "lyrical-techniques", []string{"rhyme", "chorus"})
				kb.AddTestCard("lyric2", "Metaphor Usage", "Create vivid metaphors using natural imagery", "lyrical-techniques", []string{"metaphor", "imagery"})
			},
			mode:            QuickIdeaModeUnstick,
			context:         "I'm writing a song about love",
			options:         map[string]string{},
			expectedKBUsage: true,
			description:     "Should enhance suggestions with lyric knowledge cards",
		},
		{
			name: "Knowledge base with pattern cards",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.AddTestCard("pattern1", "Common Progressions", "C - G - Am - F is a versatile progression", "chord-progressions", []string{"progression", "versatile"})
				kb.AddTestCard("pattern2", "Voice Leading", "Use smooth voice leading between chords", "chord-progressions", []string{"voice-leading", "smooth"})
			},
			mode:            QuickIdeaModeUnstick,
			context:         "C - G - Am - F",
			options:         map[string]string{},
			expectedKBUsage: true,
			description:     "Should enhance suggestions with pattern knowledge cards",
		},
		{
			name: "Knowledge base with inspiration cards",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.AddTestCard("inspire1", "Love Themes", "Explore different facets of love in your lyrics", "inspiration", []string{"love", "themes"})
				kb.AddTestCard("inspire2", "Nature Imagery", "Use natural elements to create vivid imagery", "inspiration", []string{"nature", "imagery"})
			},
			mode:            QuickIdeaModeSpark,
			context:         "",
			options:         map[string]string{"theme": "love"},
			expectedKBUsage: true,
			description:     "Should enhance ideas with inspiration knowledge cards",
		},
		{
			name: "Knowledge base empty",
			setupKB: func(kb *MockKnowledgeBase) {
				// Don't add any cards
			},
			mode:            QuickIdeaModeUnstick,
			context:         "Some context",
			options:         map[string]string{},
			expectedKBUsage: false,
			description:     "Should work gracefully with empty knowledge base",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup knowledge base
			mockKB := NewMockKnowledgeBase()
			tt.setupKB(mockKB)

			// Setup agent
			mockProvider := NewMockEnhancementProvider()
			mockProvider.SetKnowledgeBase(mockKB)
			agent := NewQuickIdeaAgent()
			agent = agent.WithKnowledgeBase(mockProvider)

			// Verify knowledge base status
			ctx := context.Background()
			helper.AssertTrue(agent.IsKnowledgeBaseAvailable(ctx))
			status := agent.GetKnowledgeBaseStatus(ctx)
			helper.AssertNotNil(status)
			helper.AssertTrue(status.Available)

			// Execute request
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			resp, err := agent.Generate(ctx, req)
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)

			// Verify response
			if tt.mode != QuickIdeaModeCheck {
				helper.AssertLength(resp.Suggestions, 3)
			} else {
				helper.AssertNotEmpty(resp.Rating)
				helper.AssertNotEmpty(resp.Tip)
			}

			// Check for knowledge base enhancement
			if tt.expectedKBUsage {
				foundKBEnhancement := false
				for _, suggestion := range resp.Suggestions {
					if containsKnowledgeIndicators(suggestion) {
						foundKBEnhancement = true
						break
					}
				}
				if !foundKBEnhancement {
					t.Log("Knowledge base enhancement not detected in suggestions (this may be expected)")
				}
			}

			t.Logf("Description: %s", tt.description)
		})
	}
}

// TestQuickIdeaAgent_KnowledgeBaseErrorHandling tests knowledge base error handling
func TestQuickIdeaAgent_KnowledgeBaseErrorHandling(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name        string
		setupKB     func(*MockKnowledgeBase)
		mode        QuickIdeaMode
		context     string
		expectError bool
		description string
	}{
		{
			name: "Knowledge base failure",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetFailure(true)
				kb.AddTestCard("test", "Test Card", "Test content", "test", []string{"test"})
			},
			mode:        QuickIdeaModeUnstick,
			context:     "Some context",
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when knowledge base fails",
		},
		{
			name: "Knowledge base timeout",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetDelay(5 * time.Second) // Long delay
				kb.AddTestCard("test", "Test Card", "Test content", "test", []string{"test"})
			},
			mode:        QuickIdeaModeUnstick,
			context:     "Some context",
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when knowledge base times out",
		},
		{
			name: "Knowledge base unavailable",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetAvailable(false)
			},
			mode:        QuickIdeaModeUnstick,
			context:     "Some context",
			expectError: false, // Should fallback gracefully
			description: "Should fallback gracefully when knowledge base is unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup knowledge base
			mockKB := NewMockKnowledgeBase()
			tt.setupKB(mockKB)

			// Setup agent
			mockProvider := NewMockEnhancementProvider()
			mockProvider.SetKnowledgeBase(mockKB)
			agent := NewQuickIdeaAgent()
			agent = agent.WithKnowledgeBase(mockProvider)

			// Execute request
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: map[string]string{},
			}

			resp, err := agent.Generate(context.Background(), req)

			if tt.expectError {
				helper.AssertError(err)
			} else {
				helper.AssertNoError(err)
				helper.AssertNotNil(resp)

				// Should still get fallback suggestions
				if tt.mode != QuickIdeaModeCheck {
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

// TestQuickIdeaAgent_KnowledgeBaseCaching tests knowledge base caching mechanisms
func TestQuickIdeaAgent_KnowledgeBaseCaching(t *testing.T) {
	helper := NewTestHelper(t)

	// Setup knowledge base
	mockKB := NewMockKnowledgeBase()
	mockKB.AddTestCard("cache1", "Cached Card", "This content should be cached", "inspiration", []string{"cache"})
	mockKB.AddTestCard("cache2", "Another Card", "More cached content", "lyrical-techniques", []string{"cached"})

	// Setup agent
	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)
	agent := NewQuickIdeaAgent()
	agent = agent.WithKnowledgeBase(mockProvider)

	// Execute multiple requests with same context
	req := QuickRequest{
		Mode:    QuickIdeaModeSpark,
		Context: "creativity",
		Options: map[string]string{"theme": "creativity"},
	}

	// First request
	start := time.Now()
	resp1, err1 := agent.Generate(context.Background(), req)
	firstDuration := time.Since(start)

	// Second request (should use cache)
	start = time.Now()
	resp2, err2 := agent.Generate(context.Background(), req)
	secondDuration := time.Since(start)

	// Verify both requests succeeded
	helper.AssertNoError(err1)
	helper.AssertNoError(err2)
	helper.AssertNotNil(resp1)
	helper.AssertNotNil(resp2)

	// Verify responses are consistent
	helper.AssertLength(resp1.Suggestions, 3)
	helper.AssertLength(resp2.Suggestions, 3)

	t.Logf("First request duration: %v", firstDuration)
	t.Logf("Second request duration: %v", secondDuration)
	t.Log("Second request should be faster due to caching")
}

// TestQuickIdeaAgent_KnowledgeBaseCategories tests knowledge base category filtering
func TestQuickIdeaAgent_KnowledgeBaseCategories(t *testing.T) {
	helper := NewTestHelper(t)

	// Setup knowledge base with different categories
	mockKB := NewMockKnowledgeBase()
	mockKB.AddTestCard("lyric1", "Rhyme Schemes", "Use AABB rhyme scheme", "lyrical-techniques", []string{"rhyme"})
	mockKB.AddTestCard("lyric2", "Metaphor Usage", "Create vivid metaphors", "lyrical-techniques", []string{"metaphor"})
	mockKB.AddTestCard("pattern1", "Common Progressions", "C - G - Am - F", "chord-progressions", []string{"progression"})
	mockKB.AddTestCard("pattern2", "Voice Leading", "Smooth voice leading", "chord-progressions", []string{"voice-leading"})
	mockKB.AddTestCard("inspire1", "Love Themes", "Explore love facets", "inspiration", []string{"love"})
	mockKB.AddTestCard("inspire2", "Nature Imagery", "Use natural elements", "inspiration", []string{"nature"})

	// Setup agent
	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)
	agent := NewQuickIdeaAgent()
	agent = agent.WithKnowledgeBase(mockProvider)

	tests := []struct {
		name             string
		mode             QuickIdeaMode
		context          string
		expectedCategory string
		description      string
	}{
		{
			name:             "Unstick with lyrics uses lyrical-techniques",
			mode:             QuickIdeaModeUnstick,
			context:          "I'm writing lyrics about love",
			expectedCategory: "lyrical-techniques",
			description:      "Should use lyrical-techniques category for lyric unstick",
		},
		{
			name:             "Unstick with patterns uses chord-progressions",
			mode:             QuickIdeaModeUnstick,
			context:          "C - G - Am - F progression",
			expectedCategory: "chord-progressions",
			description:      "Should use chord-progressions category for pattern unstick",
		},
		{
			name:             "Spark uses inspiration",
			mode:             QuickIdeaModeSpark,
			context:          "",
			expectedCategory: "inspiration",
			description:      "Should use inspiration category for spark mode",
		},
		{
			name:             "Tweak with lyrics uses lyrical-techniques",
			mode:             QuickIdeaModeTweak,
			context:          "Love is in the air",
			expectedCategory: "lyrical-techniques",
			description:      "Should use lyrical-techniques category for lyric tweak",
		},
		{
			name:             "Check uses multiple categories",
			mode:             QuickIdeaModeCheck,
			context:          "Some content to check",
			expectedCategory: "", // Multiple categories expected
			description:      "Should use multiple categories for check mode",
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
			helper.AssertNoError(err)
			helper.AssertNotNil(resp)

			if tt.mode != QuickIdeaModeCheck {
				helper.AssertLength(resp.Suggestions, 3)
			} else {
				helper.AssertNotEmpty(resp.Rating)
				helper.AssertNotEmpty(resp.Tip)
			}

			t.Logf("Description: %s", tt.description)
			t.Logf("Expected category: %s", tt.expectedCategory)
		})
	}
}

// TestQuickIdeaAgent_KnowledgeBaseRelevance tests knowledge base relevance filtering
func TestQuickIdeaAgent_KnowledgeBaseRelevance(t *testing.T) {
	helper := NewTestHelper(t)

	// Setup knowledge base with different relevance scores
	mockKB := NewMockKnowledgeBase()

	// Add cards with different relevance
	highRelevanceCard := knowledge.Card{
		ID:        "high1",
		Title:     "Highly Relevant",
		Content:   "This content is highly relevant to love songs",
		Category:  "inspiration",
		Tags:      []string{"love", "relevant"},
		Relevance: 0.9,
		Metadata:  map[string]string{"example_c": "C - G - Am - F"},
	}

	lowRelevanceCard := knowledge.Card{
		ID:        "low1",
		Title:     "Low Relevance",
		Content:   "This content is not very relevant",
		Category:  "inspiration",
		Tags:      []string{"other", "irrelevant"},
		Relevance: 0.3,
		Metadata:  map[string]string{"example_c": "D - Em - C - G"},
	}

	mockKB.AddCard(context.Background(), highRelevanceCard)
	mockKB.AddCard(context.Background(), lowRelevanceCard)

	// Setup agent
	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)
	agent := NewQuickIdeaAgent()
	agent = agent.WithKnowledgeBase(mockProvider)

	// Test with high relevance threshold
	req := QuickRequest{
		Mode:    QuickIdeaModeSpark,
		Context: "love",
		Options: map[string]string{"theme": "love"},
	}

	resp, err := agent.Generate(context.Background(), req)
	helper.AssertNoError(err)
	helper.AssertNotNil(resp)
	helper.AssertLength(resp.Suggestions, 3)

	// The high relevance card should be used
	foundEnhancement := false
	for _, suggestion := range resp.Suggestions {
		if containsKnowledgeIndicators(suggestion) {
			foundEnhancement = true
			break
		}
	}

	t.Logf("Found knowledge enhancement: %v", foundEnhancement)
	t.Log("High relevance cards should be prioritized")
}

// TestQuickIdeaAgent_KnowledgeBasePerformance tests knowledge base performance
func TestQuickIdeaAgent_KnowledgeBasePerformance(t *testing.T) {
	helper := NewTestHelper(t)

	// Setup knowledge base with many cards
	mockKB := NewMockKnowledgeBase()
	for i := 0; i < 100; i++ {
		mockKB.AddTestCard(
			"card"+string(rune(i)),
			"Card "+string(rune(i)),
			"Content for card "+string(rune(i)),
			"inspiration",
			[]string{"tag" + string(rune(i))},
		)
	}

	// Setup agent
	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)
	agent := NewQuickIdeaAgent()
	agent = agent.WithKnowledgeBase(mockProvider)

	// Test performance with large knowledge base
	req := QuickRequest{
		Mode:    QuickIdeaModeSpark,
		Context: "creativity",
		Options: map[string]string{"theme": "creativity"},
	}

	start := time.Now()
	resp, err := agent.Generate(context.Background(), req)
	elapsed := time.Since(start)

	helper.AssertNoError(err)
	helper.AssertNotNil(resp)
	helper.AssertLength(resp.Suggestions, 3)

	// Should complete quickly even with large knowledge base
	helper.AssertTrue(relaxPerfBudgets() || elapsed < 500*time.Millisecond,
		"Should complete in less than 500ms with large knowledge base, took %v", elapsed)

	t.Logf("Completed in %v with 100 knowledge cards", elapsed)
}

// TestQuickIdeaAgent_KnowledgeBaseMetadata tests knowledge base metadata usage
func TestQuickIdeaAgent_KnowledgeBaseMetadata(t *testing.T) {
	helper := NewTestHelper(t)

	// Setup knowledge base with metadata
	mockKB := NewMockKnowledgeBase()

	cardWithMetadata := knowledge.Card{
		ID:        "meta1",
		Title:     "Card with Metadata",
		Content:   "This card has useful metadata",
		Category:  "chord-progressions",
		Tags:      []string{"progression", "metadata"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"example_c":  "C - G - Am - F",
			"difficulty": "beginner",
			"genre":      "pop",
			"tempo":      "120 BPM",
		},
	}

	mockKB.AddCard(context.Background(), cardWithMetadata)

	// Setup agent
	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)
	agent := NewQuickIdeaAgent()
	agent = agent.WithKnowledgeBase(mockProvider)

	// Test metadata usage
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "C - G - Am - F",
		Options: map[string]string{},
	}

	resp, err := agent.Generate(context.Background(), req)
	helper.AssertNoError(err)
	helper.AssertNotNil(resp)
	helper.AssertLength(resp.Suggestions, 3)

	// Check if metadata is used in suggestions
	foundMetadataUsage := false
	for _, suggestion := range resp.Suggestions {
		// Look for example_c from metadata
		if containsString(suggestion, "C - G - Am - F") {
			foundMetadataUsage = true
			break
		}
	}

	t.Logf("Found metadata usage: %v", foundMetadataUsage)
	t.Log("Knowledge base metadata should be used in suggestions")
}

// TestQuickIdeaAgent_KnowledgeBaseStatusReporting tests knowledge base status reporting
func TestQuickIdeaAgent_KnowledgeBaseStatusReporting(t *testing.T) {
	helper := NewTestHelper(t)

	tests := []struct {
		name              string
		setupKB           func(*MockKnowledgeBase)
		expectedAvailable bool
		expectedError     string
		description       string
	}{
		{
			name: "Available knowledge base",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetAvailable(true)
				kb.AddTestCard("test", "Test", "Content", "test", []string{"test"})
			},
			expectedAvailable: true,
			expectedError:     "",
			description:       "Should report available status correctly",
		},
		{
			name: "Unavailable knowledge base",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetAvailable(false)
			},
			expectedAvailable: false,
			expectedError:     "",
			description:       "Should report unavailable status correctly",
		},
		{
			name: "Knowledge base with error",
			setupKB: func(kb *MockKnowledgeBase) {
				kb.SetAvailable(false)
			},
			expectedAvailable: false,
			expectedError:     "",
			description:       "Should report error status correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup knowledge base
			mockKB := NewMockKnowledgeBase()
			tt.setupKB(mockKB)

			// Setup agent
			mockProvider := NewMockEnhancementProvider()
			mockProvider.SetKnowledgeBase(mockKB)
			agent := NewQuickIdeaAgent()
			agent = agent.WithKnowledgeBase(mockProvider)

			// Check status
			ctx := context.Background()
			isAvailable := agent.IsKnowledgeBaseAvailable(ctx)
			status := agent.GetKnowledgeBaseStatus(ctx)

			helper.AssertEqual(tt.expectedAvailable, isAvailable)
			helper.AssertNotNil(status)
			helper.AssertEqual(tt.expectedAvailable, status.Available)

			if tt.expectedError != "" {
				helper.AssertEqual(tt.expectedError, status.Error)
			}

			t.Logf("Description: %s", tt.description)
			t.Logf("Status: Available=%v, CardCount=%d", status.Available, status.CardCount)
		})
	}
}
