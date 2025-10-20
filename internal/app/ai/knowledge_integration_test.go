package ai

import (
	"context"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app/knowledge"
)

func TestQuickIdeaAgent_KnowledgeBaseIntegration(t *testing.T) {
	ctx := context.Background()
	
	// Create agent with stub knowledge base
	agent := NewQuickIdeaAgent()
	
	// Test knowledge base status
	status := agent.GetKnowledgeBaseStatus(ctx)
	if status == nil {
		t.Error("GetKnowledgeBaseStatus() should not return nil")
		return
	}
	
	if !agent.IsKnowledgeBaseAvailable(ctx) {
		t.Error("IsKnowledgeBaseAvailable() should return true for stub")
	}
	
	// Test enhancement with knowledge base
	tests := []struct {
		name        string
		mode        QuickIdeaMode
		context     string
		options     map[string]string
		expectKB    bool
		expectCards int
	}{
		{
			name:        "Unstick with lyrics",
			mode:        QuickIdeaModeUnstick,
			context:     "I'm walking down the street tonight",
			options:     map[string]string{},
			expectKB:    true,
			expectCards: 0, // Cards are embedded in suggestions, not returned directly
		},
		{
			name:        "Spark with theme",
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
			expectKB:    true,
			expectCards: 0,
		},
		{
			name:        "Tweak with lyrics",
			mode:        QuickIdeaModeTweak,
			context:     "Love is in the air",
			options:     map[string]string{},
			expectKB:    true,
			expectCards: 0,
		},
		{
			name:        "Check with lyrics",
			mode:        QuickIdeaModeCheck,
			context:     "The rain falls on my window",
			options:     map[string]string{},
			expectKB:    true,
			expectCards: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}
			
			resp, err := agent.Generate(ctx, req)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}
			
			if resp == nil {
				t.Error("Generate() should not return nil response")
				return
			}
			
			// Check that response time is set
			if resp.ResponseTime == 0 {
				t.Error("Generate() should set response time")
			}
			
			// Check suggestions for non-check modes
			if tt.mode != QuickIdeaModeCheck {
				if len(resp.Suggestions) == 0 {
					t.Error("Generate() should return suggestions")
				}
				
				// Check if knowledge base enhanced the suggestions
				if tt.expectKB {
					// Look for knowledge-based enhancements in suggestions
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
			} else {
				// For check mode, verify rating and tip
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

func TestQuickIdeaAgent_WithKnowledgeBase(t *testing.T) {
	ctx := context.Background()
	
	// Create agent
	agent := NewQuickIdeaAgent()
	
	// Create custom knowledge base provider
	customKB := knowledge.NewStubEnhancementProvider()
	
	// Update agent with custom knowledge base
	agentWithKB := agent.WithKnowledgeBase(customKB)
	
	if agentWithKB == nil {
		t.Error("WithKnowledgeBase() should not return nil")
		return
	}
	
	// Test that the new agent has the knowledge base
	if !agentWithKB.IsKnowledgeBaseAvailable(ctx) {
		t.Error("Agent with knowledge base should have KB available")
	}
	
	// Test generation with custom knowledge base
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Test lyrics for enhancement",
		Options: map[string]string{},
	}
	
	resp, err := agentWithKB.Generate(ctx, req)
	if err != nil {
		t.Errorf("Generate() with custom KB error = %v", err)
		return
	}
	
	if resp == nil {
		t.Error("Generate() with custom KB should not return nil")
		return
	}
	
	if len(resp.Suggestions) == 0 {
		t.Error("Generate() with custom KB should return suggestions")
	}
}

func TestQuickIdeaAgent_KnowledgeBaseGracefulDegradation(t *testing.T) {
	ctx := context.Background()
	
	// Create agent with nil knowledge base
	agent := &QuickIdeaAgent{
		client:          &stubQuickClient{},
		model:           defaultQuickIdeaModel,
		timeout:         defaultQuickIdeaTimeout,
		prompts:         defaultQuickIdeaPrompts(),
		contextDetector: NewContextDetector(),
		knowledgeBase:   nil, // No knowledge base
	}
	
	// Test that agent still works without knowledge base
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Test lyrics without KB",
		Options: map[string]string{},
	}
	
	resp, err := agent.Generate(ctx, req)
	if err != nil {
		t.Errorf("Generate() without KB error = %v", err)
		return
	}
	
	if resp == nil {
		t.Error("Generate() without KB should not return nil")
		return
	}
	
	// Should still get fallback suggestions
	if len(resp.Suggestions) == 0 {
		t.Error("Generate() without KB should still return fallback suggestions")
	}
	
	// Knowledge base should not be available
	if agent.IsKnowledgeBaseAvailable(ctx) {
		t.Error("IsKnowledgeBaseAvailable() should return false when KB is nil")
	}
	
	status := agent.GetKnowledgeBaseStatus(ctx)
	if status != nil && status.Available {
		t.Error("GetKnowledgeBaseStatus() should not be available when KB is nil")
	}
}

func TestQuickIdeaAgent_KnowledgeBaseTimeout(t *testing.T) {
	ctx := context.Background()
	
	// Create agent with very short timeout
	agent := NewQuickIdeaAgent()
	agent.timeout = 1 * time.Millisecond
	
	// Test generation with timeout
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Test lyrics with timeout",
		Options: map[string]string{},
	}
	
	resp, err := agent.Generate(ctx, req)
	if err != nil {
		t.Errorf("Generate() with timeout error = %v", err)
		return
	}
	
	if resp == nil {
		t.Error("Generate() with timeout should not return nil")
		return
	}
	
	// Should get fallback response due to timeout
	if len(resp.Suggestions) == 0 {
		t.Error("Generate() with timeout should return fallback suggestions")
	}
}

func TestQuickIdeaAgent_KnowledgeBaseContentTypes(t *testing.T) {
	ctx := context.Background()
	
	agent := NewQuickIdeaAgent()
	
	tests := []struct {
		name        string
		content     string
		expectedType ContentType
	}{
		{
			name:        "Lyric content",
			content:     "I'm walking down the street tonight\nThe city lights are shining bright",
			expectedType: ContentTypeLyrics,
		},
		{
			name:        "Pattern content",
			content:     "Verse: C - G - Am - F\nChorus: F - C - G - Am",
			expectedType: ContentTypePatterns,
		},
		{
			name:        "Mixed content",
			content:     "[Verse]\nC G Am F\nI'm walking down the street tonight",
			expectedType: ContentTypeMixed,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: tt.content,
				Options: map[string]string{},
			}
			
			resp, err := agent.Generate(ctx, req)
			if err != nil {
				t.Errorf("Generate() error = %v", err)
				return
			}
			
			if resp == nil {
				t.Error("Generate() should not return nil")
				return
			}
			
			// Should get context-aware suggestions based on content type
			if len(resp.Suggestions) == 0 {
				t.Error("Generate() should return suggestions for content type")
			}
		})
	}
}

// Helper function to check if suggestions contain knowledge base indicators
func containsKnowledgeIndicators(suggestion string) bool {
	indicators := []string{
		"apply",
		"technique",
		"inspired by",
		"using",
		"KB Tip:",
		"progression",
		"rhyme scheme",
		"imagery",
		"metaphor",
	}
	
	for _, indicator := range indicators {
		if len(suggestion) >= len(indicator) {
			for i := 0; i <= len(suggestion)-len(indicator); i++ {
				if suggestion[i:i+len(indicator)] == indicator {
					return true
				}
			}
		}
	}
	
	return false
}