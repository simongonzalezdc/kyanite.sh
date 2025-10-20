package ai

import (
	"context"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app/knowledge"
)

// BenchmarkContextDetector_AnalyzeContent benchmarks context detection performance
func BenchmarkContextDetector_AnalyzeContent(b *testing.B) {
	detector := NewContextDetector()

	tests := []struct {
		name    string
		content string
	}{
		{
			name: "Short lyric content",
			content: `[Verse]
I'm walking down the street tonight
The city lights are shining bright`,
		},
		{
			name: "Medium pattern content",
			content: `Verse: C - G - Am - F
Chorus: F - C - G - Am
Bridge: Dm - Am - E - Am
Tempo: 120 BPM
Key: C Major`,
		},
		{
			name: "Large mixed content",
			content: func() string {
				var content string
				for i := 0; i < 50; i++ {
					content += "[Verse]\n"
					content += "C G Am F\n"
					content += "This is line " + string(rune(i)) + " of the song\n"
					content += "With some lyrical content and emotion\n\n"
				}
				return content
			}(),
		},
		{
			name: "Very large content",
			content: func() string {
				var content string
				for i := 0; i < 500; i++ {
					content += "[Verse]\n"
					content += "This is line " + string(rune(i%26)) + " of the song\n"
					content += "With some lyrical content and emotion\n"
					content += "Love and heartbreak and feeling strong\n\n"
				}
				return content
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				detector.AnalyzeContent(tt.content)
			}
		})
	}
}

// BenchmarkContextDetector_GetContextAnalysis benchmarks detailed context analysis
func BenchmarkContextDetector_GetContextAnalysis(b *testing.B) {
	detector := NewContextDetector()

	tests := []struct {
		name    string
		content string
	}{
		{
			name: "Short content",
			content: `[Verse]
I'm walking down the street tonight`,
		},
		{
			name: "Medium content",
			content: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright
You left me here all alone
Now I'm trying to find my way home

[Chorus]
Love is like a burning fire
Desire takes me higher
Baby can't you see the light
You're my one and only delight`,
		},
		{
			name: "Large content",
			content: func() string {
				var content string
				for i := 0; i < 100; i++ {
					content += "[Verse " + string(rune('A'+(i%26))) + "]\n"
					content += "This is line " + string(rune(i)) + " of the song\n"
					content += "With some lyrical content and emotion\n"
					content += "Love and heartbreak and feeling strong\n\n"
				}
				return content
			}(),
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				detector.GetContextAnalysis(tt.content)
			}
		})
	}
}

// BenchmarkQuickIdeaAgent_Generate benchmarks quick idea generation performance
func BenchmarkQuickIdeaAgent_Generate(b *testing.B) {
	mockClient := NewMockLLMClient()
	mockClient.SetDelay(0) // No delay for benchmarking

	tests := []struct {
		name    string
		mode    QuickIdeaMode
		context string
		options map[string]string
	}{
		{
			name:    "Unstick with lyrics",
			mode:    QuickIdeaModeUnstick,
			context: "I'm walking down the street tonight",
			options: map[string]string{},
		},
		{
			name:    "Unstick with patterns",
			mode:    QuickIdeaModeUnstick,
			context: "C - G - Am - F",
			options: map[string]string{},
		},
		{
			name:    "Spark with theme",
			mode:    QuickIdeaModeSpark,
			context: "",
			options: map[string]string{"theme": "love"},
		},
		{
			name:    "Tweak with lyrics",
			mode:    QuickIdeaModeTweak,
			context: "Love is in the air tonight",
			options: map[string]string{},
		},
		{
			name:    "Check with patterns",
			mode:    QuickIdeaModeCheck,
			context: "C - G - Am - F",
			options: map[string]string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := agent.Generate(context.Background(), req)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkQuickIdeaAgent_Fallback benchmarks fallback generation performance
func BenchmarkQuickIdeaAgent_Fallback(b *testing.B) {
	mockClient := NewMockLLMClient()
	mockClient.SetFailure(true) // Force fallback

	tests := []struct {
		name    string
		mode    QuickIdeaMode
		context string
		options map[string]string
	}{
		{
			name:    "Fallback unstick",
			mode:    QuickIdeaModeUnstick,
			context: "I'm walking down the street",
			options: map[string]string{},
		},
		{
			name:    "Fallback spark",
			mode:    QuickIdeaModeSpark,
			context: "",
			options: map[string]string{"theme": "love"},
		},
		{
			name:    "Fallback tweak",
			mode:    QuickIdeaModeTweak,
			context: "Original line here",
			options: map[string]string{},
		},
		{
			name:    "Fallback check",
			mode:    QuickIdeaModeCheck,
			context: "Some content",
			options: map[string]string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := agent.Generate(context.Background(), req)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkContextAwarePrompts_GetPrompt benchmarks prompt retrieval performance
func BenchmarkContextAwarePrompts_GetPrompt(b *testing.B) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
	}{
		{
			name:        "Lyric unstick",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
		},
		{
			name:        "Pattern unstick",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
		},
		{
			name:        "Mixed unstick",
			contentType: ContentTypeMixed,
			mode:        QuickIdeaModeUnstick,
		},
		{
			name:        "Lyric spark",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
		},
		{
			name:        "Pattern spark",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeSpark,
		},
		{
			name:        "Lyric tweak",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeTweak,
		},
		{
			name:        "Pattern tweak",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeTweak,
		},
		{
			name:        "Lyric check",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeCheck,
		},
		{
			name:        "Pattern check",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeCheck,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				prompts.GetPrompt(tt.contentType, tt.mode)
			}
		})
	}
}

// BenchmarkContextAwarePrompts_RenderPrompt benchmarks prompt rendering performance
func BenchmarkContextAwarePrompts_RenderPrompt(b *testing.B) {
	prompts := NewContextAwarePrompts()
	prompts.Initialize()

	tests := []struct {
		name        string
		contentType ContentType
		mode        QuickIdeaMode
		context     string
		options     map[string]string
	}{
		{
			name:        "Render lyric unstick",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeUnstick,
			context:     "I'm walking down the street tonight",
			options:     map[string]string{},
		},
		{
			name:        "Render pattern unstick",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeUnstick,
			context:     "C - G - Am - F progression",
			options:     map[string]string{},
		},
		{
			name:        "Render spark with theme",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeSpark,
			context:     "",
			options:     map[string]string{"theme": "love"},
		},
		{
			name:        "Render tweak with context",
			contentType: ContentTypeLyrics,
			mode:        QuickIdeaModeTweak,
			context:     "Love is in the air tonight",
			options:     map[string]string{},
		},
		{
			name:        "Render check with context",
			contentType: ContentTypePatterns,
			mode:        QuickIdeaModeCheck,
			context:     "C - G - Am - F",
			options:     map[string]string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				prompts.RenderPrompt(tt.contentType, tt.mode, tt.context, tt.options)
			}
		})
	}
}

// BenchmarkKnowledgeBase_Search benchmarks knowledge base search performance
func BenchmarkKnowledgeBase_Search(b *testing.B) {
	mockKB := NewMockKnowledgeBase()

	// Add test cards
	for i := 0; i < 100; i++ {
		mockKB.AddTestCard(
			"card"+string(rune(i)),
			"Card "+string(rune(i)),
			"Content for card "+string(rune(i))+" with search terms",
			"inspiration",
			[]string{"tag" + string(rune(i)), "search"},
		)
	}

	tests := []struct {
		name    string
		query   string
		options knowledge.SearchOptions
	}{
		{
			name:  "Simple search",
			query: "search",
			options: knowledge.SearchOptions{
				Limit:        10,
				MinRelevance: 0.5,
				UseCache:     true,
			},
		},
		{
			name:  "Category filtered search",
			query: "search",
			options: knowledge.SearchOptions{
				Limit:        5,
				Categories:   []string{"inspiration"},
				MinRelevance: 0.7,
				UseCache:     true,
			},
		},
		{
			name:  "High relevance search",
			query: "search",
			options: knowledge.SearchOptions{
				Limit:        3,
				MinRelevance: 0.9,
				UseCache:     true,
			},
		},
		{
			name:  "No cache search",
			query: "search",
			options: knowledge.SearchOptions{
				Limit:        10,
				MinRelevance: 0.5,
				UseCache:     false,
			},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := mockKB.Search(context.Background(), tt.query, tt.options)
				if err != nil {
					b.Fatalf("Search failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkQuickIdeaAgent_KnowledgeBaseIntegration benchmarks knowledge base integration performance
func BenchmarkQuickIdeaAgent_KnowledgeBaseIntegration(b *testing.B) {
	mockClient := NewMockLLMClient()
	mockClient.SetDelay(0)

	mockKB := NewMockKnowledgeBase()
	for i := 0; i < 50; i++ {
		mockKB.AddTestCard(
			"card"+string(rune(i)),
			"Card "+string(rune(i)),
			"Content for card "+string(rune(i)),
			"inspiration",
			[]string{"tag" + string(rune(i))},
		)
	}

	mockProvider := NewMockEnhancementProvider()
	mockProvider.SetKnowledgeBase(mockKB)

	agent := NewQuickIdeaAgent()
	agent = agent.WithClient(mockClient, 5*time.Second)
	agent = agent.WithKnowledgeBase(mockProvider)

	tests := []struct {
		name    string
		mode    QuickIdeaMode
		context string
		options map[string]string
	}{
		{
			name:    "Unstick with knowledge base",
			mode:    QuickIdeaModeUnstick,
			context: "I'm writing about love",
			options: map[string]string{},
		},
		{
			name:    "Spark with knowledge base",
			mode:    QuickIdeaModeSpark,
			context: "",
			options: map[string]string{"theme": "creativity"},
		},
		{
			name:    "Tweak with knowledge base",
			mode:    QuickIdeaModeTweak,
			context: "Love is in the air",
			options: map[string]string{},
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			req := QuickRequest{
				Mode:    tt.mode,
				Context: tt.context,
				Options: tt.options,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := agent.Generate(context.Background(), req)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}

// BenchmarkQuickIdeaAgent_ConcurrentGeneration benchmarks concurrent generation performance
func BenchmarkQuickIdeaAgent_ConcurrentGeneration(b *testing.B) {
	mockClient := NewMockLLMClient()
	mockClient.SetDelay(0)

	agent := NewQuickIdeaAgent()
	agent = agent.WithClient(mockClient, 5*time.Second)

	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "I'm walking down the street tonight",
		Options: map[string]string{},
	}

	b.Run("Sequential", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := agent.Generate(context.Background(), req)
			if err != nil {
				b.Fatalf("Generate failed: %v", err)
			}
		}
	})

	b.Run("Concurrent", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := agent.Generate(context.Background(), req)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	})
}

// BenchmarkContextDetector_MultipleContentTypes benchmarks performance across different content types
func BenchmarkContextDetector_MultipleContentTypes(b *testing.B) {
	detector := NewContextDetector()

	contentTypes := []struct {
		name    string
		content string
	}{
		{
			name: "Lyrics",
			content: `[Verse 1]
I'm walking down the street tonight
The city lights are shining bright
You left me here all alone
Now I'm trying to find my way home

[Chorus]
Love is like a burning fire
Desire takes me higher
Baby can't you see the light
You're my one and only delight`,
		},
		{
			name: "Patterns",
			content: `Tempo: 120 BPM
Key: C Major
Time Signature: 4/4
Verse: C - G - Am - F
Chorus: F - C - G - Am
Bridge: Dm - Am - E - Am`,
		},
		{
			name: "Mixed",
			content: `[Verse]
C G Am F
I'm walking down the street tonight
The city lights are shining bright
Am F C G
You left me here all alone`,
		},
	}

	for _, ct := range contentTypes {
		b.Run(ct.name, func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				detector.AnalyzeContent(ct.content)
			}
		})
	}
}

// BenchmarkQuickIdeaAgent_ResponseTime benchmarks response time with different client delays
func BenchmarkQuickIdeaAgent_ResponseTime(b *testing.B) {
	delays := []time.Duration{
		0,
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
	}

	for _, delay := range delays {
		b.Run("Delay_"+delay.String(), func(b *testing.B) {
			mockClient := NewMockLLMClient()
			mockClient.SetDelay(delay)

			agent := NewQuickIdeaAgent()
			agent = agent.WithClient(mockClient, 5*time.Second)

			req := QuickRequest{
				Mode:    QuickIdeaModeUnstick,
				Context: "Some context",
				Options: map[string]string{},
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := agent.Generate(context.Background(), req)
				if err != nil {
					b.Fatalf("Generate failed: %v", err)
				}
			}
		})
	}
}
