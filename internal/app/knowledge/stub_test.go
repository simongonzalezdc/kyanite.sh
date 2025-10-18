package knowledge

import (
	"context"
	"testing"
)

func TestStubKnowledgeBase_Search(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		query    string
		options  SearchOptions
		expected int
	}{
		{
			name:     "Search for rhyme schemes",
			query:    "rhyme",
			options:  SearchOptions{Limit: 5},
			expected: 2, // Should find AABB and ABAB rhyme schemes
		},
		{
			name:     "Search for chord progressions",
			query:    "chord",
			options:  SearchOptions{Limit: 5},
			expected: 2, // Should find 50s and four-chord progressions
		},
		{
			name:     "Search for lyrical techniques",
			query:    "imagery",
			options:  SearchOptions{Limit: 5},
			expected: 1, // Should find sensory imagery card
		},
		{
			name:     "Search with category filter",
			query:    "",
			options:  SearchOptions{Limit: 10, Categories: []string{"inspiration"}},
			expected: 2, // Should find love and adversity cards
		},
		{
			name:     "Search with tag filter",
			query:    "",
			options:  SearchOptions{Limit: 10, Tags: []string{"pop"}},
			expected: 3, // Should find AABB, ABAB, and four-chord progression
		},
		{
			name:     "Search with no results",
			query:    "nonexistent",
			options:  SearchOptions{Limit: 5},
			expected: 0,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := kb.Search(ctx, tt.query, tt.options)
			if err != nil {
				t.Errorf("Search() error = %v", err)
				return
			}
			
			if len(result.Cards) != tt.expected {
				t.Errorf("Search() returned %d cards, expected %d", len(result.Cards), tt.expected)
			}
			
			if result.Query != tt.query {
				t.Errorf("Search() query = %s, expected %s", result.Query, tt.query)
			}
			
			if result.Duration == 0 {
				t.Error("Search() duration should be > 0")
			}
		})
	}
}

func TestStubKnowledgeBase_CardManagement(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	// Test adding a card
	newCard := Card{
		ID:       "test-card",
		Title:    "Test Card",
		Content:  "This is a test card",
		Category: "test",
		Tags:     []string{"test", "example"},
		Relevance: 0.8,
	}
	
	err := kb.AddCard(ctx, newCard)
	if err != nil {
		t.Errorf("AddCard() error = %v", err)
	}
	
	// Test getting the card
	retrieved, err := kb.GetCard(ctx, "test-card")
	if err != nil {
		t.Errorf("GetCard() error = %v", err)
		return
	}
	
	if retrieved.Title != "Test Card" {
		t.Errorf("GetCard() title = %s, expected Test Card", retrieved.Title)
	}
	
	// Test updating the card
	updatedCard := newCard
	updatedCard.Content = "Updated content"
	err = kb.UpdateCard(ctx, updatedCard)
	if err != nil {
		t.Errorf("UpdateCard() error = %v", err)
	}
	
	// Verify the update
	retrieved, err = kb.GetCard(ctx, "test-card")
	if err != nil {
		t.Errorf("GetCard() after update error = %v", err)
		return
	}
	
	if retrieved.Content != "Updated content" {
		t.Errorf("GetCard() content = %s, expected Updated content", retrieved.Content)
	}
	
	// Test deleting the card
	err = kb.DeleteCard(ctx, "test-card")
	if err != nil {
		t.Errorf("DeleteCard() error = %v", err)
	}
	
	// Verify deletion
	_, err = kb.GetCard(ctx, "test-card")
	if err == nil {
		t.Error("GetCard() after deletion should return error")
	}
}

func TestStubKnowledgeBase_CategoriesAndTags(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	categories, err := kb.GetCategories(ctx)
	if err != nil {
		t.Errorf("GetCategories() error = %v", err)
		return
	}
	
	expectedCategories := []string{
		"rhyme-schemes",
		"chord-progressions", 
		"lyrical-techniques",
		"song-structure",
		"inspiration",
	}
	
	if len(categories) != len(expectedCategories) {
		t.Errorf("GetCategories() returned %d categories, expected %d", len(categories), len(expectedCategories))
	}
	
	for _, expected := range expectedCategories {
		found := false
		for _, actual := range categories {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetCategories() missing expected category: %s", expected)
		}
	}
	
	tags, err := kb.GetTags(ctx)
	if err != nil {
		t.Errorf("GetTags() error = %v", err)
		return
	}
	
	if len(tags) == 0 {
		t.Error("GetTags() should return tags")
	}
	
	// Check for some expected tags
	expectedTags := []string{"rhyme", "chords", "imagery", "love"}
	for _, expected := range expectedTags {
		found := false
		for _, actual := range tags {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetTags() missing expected tag: %s", expected)
		}
	}
}

func TestStubKnowledgeBase_Status(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	status := kb.GetStatus(ctx)
	if status == nil {
		t.Error("GetStatus() should not return nil")
		return
	}
	
	if status.Available {
		t.Error("Stub knowledge base should not be available as real KB")
	}
	
	if status.CardCount == 0 {
		t.Error("Stub knowledge base should have cards")
	}
	
	if status.Version != "stub-1.0.0" {
		t.Errorf("Status version = %s, expected stub-1.0.0", status.Version)
	}
	
	if status.Error == "" {
		t.Error("Stub knowledge base should have error message indicating it's a stub")
	}
}

func TestStubKnowledgeBase_IsAvailable(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	// Stub should always return true for IsAvailable (graceful degradation)
	if !kb.IsAvailable(ctx) {
		t.Error("Stub knowledge base should be available for graceful degradation")
	}
}

func TestStubEnhancementProvider_EnhanceLyrics(t *testing.T) {
	provider := NewStubEnhancementProvider()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		lyrics   string
		options  SearchOptions
		expected string
	}{
		{
			name:     "Enhance simple lyrics",
			lyrics:   "I love you so much",
			options:  SearchOptions{Limit: 3},
			expected: "Enhanced",
		},
		{
			name:     "Enhance with imagery technique",
			lyrics:   "The rain falls down",
			options:  SearchOptions{Limit: 3, Categories: []string{"lyrical-techniques"}},
			expected: "imagery",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestion, err := provider.EnhanceLyrics(ctx, tt.lyrics, tt.options)
			if err != nil {
				t.Errorf("EnhanceLyrics() error = %v", err)
				return
			}
			
			if suggestion.Original != tt.lyrics {
				t.Errorf("EnhanceLyrics() original = %s, expected %s", suggestion.Original, tt.lyrics)
			}
			
			if suggestion.Suggestion == "" {
				t.Error("EnhanceLyrics() suggestion should not be empty")
			}
			
			if suggestion.Reason == "" {
				t.Error("EnhanceLyrics() reason should not be empty")
			}
			
			if suggestion.Confidence <= 0 {
				t.Error("EnhanceLyrics() confidence should be > 0")
			}
			
			// Check if expected content is in the suggestion
			if tt.expected != "" && !containsString(suggestion.Suggestion, tt.expected) {
				t.Errorf("EnhanceLyrics() suggestion should contain %s", tt.expected)
			}
		})
	}
}

func TestStubEnhancementProvider_EnhancePatterns(t *testing.T) {
	provider := NewStubEnhancementProvider()
	ctx := context.Background()
	
	tests := []struct {
		name     string
		pattern  string
		options  SearchOptions
		expected string
	}{
		{
			name:     "Enhance simple pattern",
			pattern:  "C - G - Am",
			options:  SearchOptions{Limit: 3},
			expected: "Enhanced",
		},
		{
			name:     "Enhance with chord progression",
			pattern:  "Simple progression",
			options:  SearchOptions{Limit: 3, Categories: []string{"chord-progressions"}},
			expected: "C - G - Am - F", // Should return the four-chord progression
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			suggestion, err := provider.EnhancePatterns(ctx, tt.pattern, tt.options)
			if err != nil {
				t.Errorf("EnhancePatterns() error = %v", err)
				return
			}
			
			if suggestion.Original != tt.pattern {
				t.Errorf("EnhancePatterns() original = %s, expected %s", suggestion.Original, tt.pattern)
			}
			
			if suggestion.Suggestion == "" {
				t.Error("EnhancePatterns() suggestion should not be empty")
			}
			
			if suggestion.Reason == "" {
				t.Error("EnhancePatterns() reason should not be empty")
			}
			
			if suggestion.Confidence <= 0 {
				t.Error("EnhancePatterns() confidence should be > 0")
			}
			
			if suggestion.PatternType == "" {
				t.Error("EnhancePatterns() pattern type should not be empty")
			}
			
			// Check if expected content is in the suggestion
			if tt.expected != "" && !containsString(suggestion.Suggestion, tt.expected) {
				t.Errorf("EnhancePatterns() suggestion should contain %s", tt.expected)
			}
		})
	}
}

func TestStubEnhancementProvider_GetInspirationCards(t *testing.T) {
	provider := NewStubEnhancementProvider()
	ctx := context.Background()
	
	result, err := provider.GetInspirationCards(ctx, "love", SearchOptions{Limit: 5})
	if err != nil {
		t.Errorf("GetInspirationCards() error = %v", err)
		return
	}
	
	if len(result.Cards) == 0 {
		t.Error("GetInspirationCards() should return cards")
	}
	
	// All cards should be from inspiration category
	for _, card := range result.Cards {
		if card.Category != "inspiration" {
			t.Errorf("GetInspirationCards() card category = %s, expected inspiration", card.Category)
		}
	}
}

func TestStubKnowledgeBase_InitializeAndClose(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	// Test initialize
	err := kb.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize() error = %v", err)
	}
	
	// Test close
	err = kb.Close(ctx)
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}
}

func TestStubKnowledgeBase_PerformanceTracking(t *testing.T) {
	kb := NewStubKnowledgeBase()
	ctx := context.Background()
	
	// Initial state should have zero searches
	if kb.searchCount != 0 {
		t.Errorf("Initial search count = %d, expected 0", kb.searchCount)
	}
	
	// Perform a search
	_, err := kb.Search(ctx, "test", SearchOptions{Limit: 5})
	if err != nil {
		t.Errorf("Search() error = %v", err)
		return
	}
	
	// Check that search count increased
	if kb.searchCount != 1 {
		t.Errorf("Search count after search = %d, expected 1", kb.searchCount)
	}
	
	// Check that last search time was set
	if kb.lastSearchTime == 0 {
		t.Error("Last search time should be set after search")
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}