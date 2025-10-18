package knowledge

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// StubKnowledgeBase provides a stub implementation of the KnowledgeBase interface
// It returns pre-defined songwriting knowledge and gracefully handles unavailability
type StubKnowledgeBase struct {
	cards     map[string]Card
	categories []string
	tags      []string
	available bool
	status    KnowledgeStatus
	
	// Performance tracking
	lastSearchTime time.Duration
	searchCount    int
	cacheHits      int
}

// NewStubKnowledgeBase creates a new stub knowledge base with pre-populated songwriting knowledge
func NewStubKnowledgeBase() *StubKnowledgeBase {
	kb := &StubKnowledgeBase{
		cards:       make(map[string]Card),
		available:   false, // Clearly indicates this is a stub
		searchCount: 0,
		cacheHits:   0,
	}
	
	// Initialize with basic songwriting knowledge
	kb.initializeSongwritingKnowledge()
	
	// Set initial status
	kb.status = KnowledgeStatus{
		Available:    false, // Stub is not a real KB
		CardCount:    len(kb.cards),
		LastSync:     time.Now(),
		Version:      "stub-1.0.0",
		Error:        "Stub implementation - full knowledge base not available",
		ResponseTime: 0,
	}
	
	return kb
}

// initializeSongwritingKnowledge populates the stub with basic songwriting information
func (kb *StubKnowledgeBase) initializeSongwritingKnowledge() {
	now := time.Now()
	
	// Rhyme schemes and patterns
	kb.addCard(Card{
		ID:       "rhyme-scheme-aabb",
		Title:    "AABB Rhyme Scheme",
		Content:  "The AABB rhyme scheme pairs rhyming lines consecutively. Common in folk music, nursery rhymes, and pop choruses. Creates a simple, memorable structure that's easy for listeners to follow.",
		Category: "rhyme-schemes",
		Tags:     []string{"rhyme", "structure", "pop", "folk"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"example": "Roses are red (A)\nViolets are blue (A)\nSugar is sweet (B)\nAnd so are you (B)",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	kb.addCard(Card{
		ID:       "rhyme-scheme-abab",
		Title:    "ABAB Rhyme Scheme",
		Content:  "The ABAB rhyme scheme alternates rhyming lines. Often found in ballads, hymns, and more sophisticated pop songs. Creates a flowing, conversational feel with less predictability than AABB.",
		Category: "rhyme-schemes",
		Tags:     []string{"rhyme", "structure", "ballad", "pop"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"example": "I wandered lonely as a cloud (A)\nThat floats on high o'er vales and hills (B)\nWhen all at once I saw a crowd (A)\nA host, of golden daffodils (B)",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	// Common chord progressions
	kb.addCard(Card{
		ID:       "progression-50s",
		Title:    "50s Progression (I-vi-IV-V)",
		Content:  "One of the most common progressions in Western pop music. The I-vi-IV-V progression creates emotional movement from tonic to submediant to subdominant to dominant. Used in countless hits from 'Stand By Me' to 'Every Breath You Take'.",
		Category: "chord-progressions",
		Tags:     []string{"chords", "progression", "pop", "common"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"example_c": "C - Am - F - G",
			"example_g": "G - Em - C - D",
			"feeling":   "nostalgic, emotional, familiar",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	kb.addCard(Card{
		ID:       "progression-four-chord",
		Title:    "Four-Chord Song Progression",
		Content:  "The I-V-vi-IV progression powers countless pop hits. This progression creates tension and release while remaining accessible. Its ubiquity makes it instantly familiar to listeners.",
		Category: "chord-progressions",
		Tags:     []string{"chords", "progression", "pop", "hit-song"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"example_c": "C - G - Am - F",
			"example_g": "G - D - Em - C",
			"feeling":   "uplifting, anthemic, radio-friendly",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	// Lyrical techniques
	kb.addCard(Card{
		ID:       "technique-imagery",
		Title:    "Sensory Imagery",
		Content:  "Use concrete sensory details to create vivid mental images. Engage multiple senses (sight, sound, smell, touch, taste) to make lyrics more immersive. Instead of 'I was sad,' try 'The rain traced cold patterns on my window pane.'",
		Category: "lyrical-techniques",
		Tags:     []string{"imagery", "sensory", "description", "show-dont-tell"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"tip": "Replace abstract emotions with specific sensory details",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	kb.addCard(Card{
		ID:       "technique-metaphor",
		Title:    "Metaphor and Simile",
		Content:  "Compare unlike things to create new understanding. Metaphors state direct comparison ('You are my sunshine') while similes use 'like' or 'as' ('Your love is like a summer day'). Effective metaphors feel fresh yet familiar.",
		Category: "lyrical-techniques",
		Tags:     []string{"metaphor", "simile", "comparison", "figurative-language"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"tip": "Choose metaphors that extend through multiple lines for deeper impact",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	// Song structure
	kb.addCard(Card{
		ID:       "structure-verse-chorus",
		Title:    "Verse-Chorus Structure",
		Content:  "The most common song structure in popular music. Verses develop the story or narrative while choruses deliver the main message and emotional hook. This balance of development and repetition creates memorable, effective songs.",
		Category: "song-structure",
		Tags:     []string{"structure", "verse", "chorus", "pop"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"example": "Verse - Chorus - Verse - Chorus - Bridge - Chorus",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	// Songwriting inspiration
	kb.addCard(Card{
		ID:       "inspiration-theme-love",
		Title:    "Writing About Love",
		Content:  "Love songs work best when they explore specific moments or feelings rather than general statements. Focus on the small details that reveal larger emotions: the way someone's laugh sounds, a shared memory, or a particular object that represents your relationship.",
		Category: "inspiration",
		Tags:     []string{"love", "theme", "emotion", "relationships"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"prompt": "What specific object or memory represents this feeling?",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	kb.addCard(Card{
		ID:       "inspiration overcoming-adversity",
		Title:    "Overcoming Adversity",
		Content:  "Songs about struggle and resilience resonate deeply. Use journey metaphors (climbing mountains, weathering storms) and contrast darkness with light. Show the process of overcoming, not just the victory. Vulnerability makes these songs powerful.",
		Category: "inspiration",
		Tags:     []string{"adversity", "resilience", "struggle", "empowerment"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"prompt": "What was the moment you realized you were stronger?",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	
	// Update categories and tags
	kb.updateCategoriesAndTags()
}

// addCard adds a card to the knowledge base
func (kb *StubKnowledgeBase) addCard(card Card) {
	kb.cards[card.ID] = card
}

// updateCategoriesAndTags refreshes the categories and tags lists
func (kb *StubKnowledgeBase) updateCategoriesAndTags() {
	categoriesSet := make(map[string]bool)
	tagsSet := make(map[string]bool)
	
	for _, card := range kb.cards {
		categoriesSet[card.Category] = true
		for _, tag := range card.Tags {
			tagsSet[tag] = true
		}
	}
	
	kb.categories = make([]string, 0, len(categoriesSet))
	for category := range categoriesSet {
		kb.categories = append(kb.categories, category)
	}
	
	kb.tags = make([]string, 0, len(tagsSet))
	for tag := range tagsSet {
		kb.tags = append(kb.tags, tag)
	}
}

// Search performs a search query against the stub knowledge base
func (kb *StubKnowledgeBase) Search(ctx context.Context, query string, options SearchOptions) (*SearchResult, error) {
	start := time.Now()
	kb.searchCount++
	
	// Simulate some processing time
	time.Sleep(10 * time.Millisecond)
	
	// Default options if not provided
	if options.Limit <= 0 {
		options.Limit = 5
	}
	if options.MinRelevance <= 0 {
		options.MinRelevance = 0.5
	}
	
	// Simple keyword matching for the stub
	queryLower := strings.ToLower(query)
	var matchingCards []Card
	
	for _, card := range kb.cards {
		if card.Relevance >= options.MinRelevance {
			// Check if query matches title, content, category, or tags
			titleMatch := strings.Contains(strings.ToLower(card.Title), queryLower)
			contentMatch := strings.Contains(strings.ToLower(card.Content), queryLower)
			categoryMatch := strings.Contains(strings.ToLower(card.Category), queryLower)
			
			tagMatch := false
			for _, tag := range card.Tags {
				if strings.Contains(strings.ToLower(tag), queryLower) {
					tagMatch = true
					break
				}
			}
			
			// Category filtering
			categoryMatchAllowed := len(options.Categories) == 0
			for _, category := range options.Categories {
				if card.Category == category {
					categoryMatchAllowed = true
					break
				}
			}
			
			// Tag filtering
			tagMatchAllowed := len(options.Tags) == 0
			for _, tag := range options.Tags {
				for _, cardTag := range card.Tags {
					if cardTag == tag {
						tagMatchAllowed = true
						break
					}
				}
			}
			
			if (titleMatch || contentMatch || categoryMatch || tagMatch) && 
			   categoryMatchAllowed && tagMatchAllowed {
				matchingCards = append(matchingCards, card)
			}
		}
	}
	
	// Shuffle results for variety in the stub
	rand.Shuffle(len(matchingCards), func(i, j int) {
		matchingCards[i], matchingCards[j] = matchingCards[j], matchingCards[i]
	})
	
	// Limit results
	if len(matchingCards) > options.Limit {
		matchingCards = matchingCards[:options.Limit]
	}
	
	duration := time.Since(start)
	kb.lastSearchTime = duration
	
	return &SearchResult{
		Cards:     matchingCards,
		Query:     query,
		Total:     len(matchingCards),
		Duration:  duration,
		FromCache: false, // Stub doesn't implement caching
	}, nil
}

// AddCard adds a new card to the stub knowledge base
func (kb *StubKnowledgeBase) AddCard(ctx context.Context, card Card) error {
	if card.ID == "" {
		return fmt.Errorf("card ID cannot be empty")
	}
	
	card.CreatedAt = time.Now()
	card.UpdatedAt = time.Now()
	kb.cards[card.ID] = card
	kb.updateCategoriesAndTags()
	
	// Update status
	kb.status.CardCount = len(kb.cards)
	
	return nil
}

// GetCard retrieves a specific card by ID
func (kb *StubKnowledgeBase) GetCard(ctx context.Context, id string) (*Card, error) {
	card, exists := kb.cards[id]
	if !exists {
		return nil, fmt.Errorf("card with ID %s not found", id)
	}
	return &card, nil
}

// UpdateCard updates an existing card
func (kb *StubKnowledgeBase) UpdateCard(ctx context.Context, card Card) error {
	if _, exists := kb.cards[card.ID]; !exists {
		return fmt.Errorf("card with ID %s not found", card.ID)
	}
	
	card.UpdatedAt = time.Now()
	kb.cards[card.ID] = card
	kb.updateCategoriesAndTags()
	
	return nil
}

// DeleteCard removes a card from the stub knowledge base
func (kb *StubKnowledgeBase) DeleteCard(ctx context.Context, id string) error {
	if _, exists := kb.cards[id]; !exists {
		return fmt.Errorf("card with ID %s not found", id)
	}
	
	delete(kb.cards, id)
	kb.updateCategoriesAndTags()
	
	// Update status
	kb.status.CardCount = len(kb.cards)
	
	return nil
}

// GetCategories returns all available categories
func (kb *StubKnowledgeBase) GetCategories(ctx context.Context) ([]string, error) {
	return kb.categories, nil
}

// GetTags returns all available tags
func (kb *StubKnowledgeBase) GetTags(ctx context.Context) ([]string, error) {
	return kb.tags, nil
}

// GetStatus returns the current status of the stub knowledge base
func (kb *StubKnowledgeBase) GetStatus(ctx context.Context) *KnowledgeStatus {
	kb.status.ResponseTime = kb.lastSearchTime
	return &kb.status
}

// IsAvailable returns whether the knowledge base is currently available
func (kb *StubKnowledgeBase) IsAvailable(ctx context.Context) bool {
	// Stub is always "available" but clearly indicates it's not a real KB
	return true
}

// Initialize sets up the stub knowledge base
func (kb *StubKnowledgeBase) Initialize(ctx context.Context) error {
	// Stub is already initialized during creation
	kb.status.LastSync = time.Now()
	return nil
}

// Close performs cleanup operations
func (kb *StubKnowledgeBase) Close(ctx context.Context) error {
	// Stub doesn't need cleanup
	return nil
}

// StubEnhancementProvider provides enhancement functionality using the stub knowledge base
type StubEnhancementProvider struct {
	kb *StubKnowledgeBase
}

// NewStubEnhancementProvider creates a new enhancement provider using the stub knowledge base
func NewStubEnhancementProvider() *StubEnhancementProvider {
	return &StubEnhancementProvider{
		kb: NewStubKnowledgeBase(),
	}
}

// EnhanceLyrics provides lyric suggestions based on stub knowledge
func (sep *StubEnhancementProvider) EnhanceLyrics(ctx context.Context, lyrics string, options SearchOptions) (*LyricSuggestion, error) {
	// Search for relevant knowledge cards
	result, err := sep.kb.Search(ctx, lyrics, options)
	if err != nil {
		return nil, err
	}
	
	// Generate a simple suggestion based on the first matching card
	suggestion := &LyricSuggestion{
		Original:   lyrics,
		Confidence: 0.7, // Stub has moderate confidence
		Cards:      result.Cards,
	}
	
	if len(result.Cards) > 0 {
		card := result.Cards[0]
		if card.Category == "lyrical-techniques" {
			suggestion.Suggestion = fmt.Sprintf("Consider applying %s to: %s", card.Title, lyrics)
			suggestion.Reason = card.Content
			suggestion.Confidence = 0.8
		} else if card.Category == "inspiration" {
			suggestion.Suggestion = fmt.Sprintf("Inspired by %s: %s", card.Title, lyrics)
			suggestion.Reason = card.Content
			suggestion.Confidence = 0.75
		} else {
			suggestion.Suggestion = fmt.Sprintf("Enhanced with %s insight: %s", card.Category, lyrics)
			suggestion.Reason = card.Content
		}
	} else {
		suggestion.Suggestion = fmt.Sprintf("Enhanced version: %s", lyrics)
		suggestion.Reason = "General lyrical enhancement - try adding more sensory details"
	}
	
	return suggestion, nil
}

// EnhancePatterns provides pattern suggestions based on stub knowledge
func (sep *StubEnhancementProvider) EnhancePatterns(ctx context.Context, pattern string, options SearchOptions) (*PatternSuggestion, error) {
	// Search for relevant knowledge cards
	result, err := sep.kb.Search(ctx, pattern, options)
	if err != nil {
		return nil, err
	}
	
	// Generate a simple suggestion based on the first matching card
	suggestion := &PatternSuggestion{
		Original:    pattern,
		PatternType: "chord-progression",
		Confidence:  0.7, // Stub has moderate confidence
		Cards:       result.Cards,
	}
	
	if len(result.Cards) > 0 {
		card := result.Cards[0]
		if card.Category == "chord-progressions" {
			suggestion.Suggestion = card.Metadata["example_c"]
			suggestion.Reason = card.Content
			suggestion.Confidence = 0.85
		} else {
			suggestion.Suggestion = fmt.Sprintf("Enhanced pattern: %s", pattern)
			suggestion.Reason = card.Content
		}
	} else {
		suggestion.Suggestion = "C - G - Am - F" // Default progression
		suggestion.Reason = "Classic pop progression that works in many contexts"
	}
	
	return suggestion, nil
}

// GetInspirationCards returns cards for creative inspiration
func (sep *StubEnhancementProvider) GetInspirationCards(ctx context.Context, theme string, options SearchOptions) (*SearchResult, error) {
	// Search with inspiration category filter
	options.Categories = []string{"inspiration"}
	
	// If theme is provided, include it in the search
	query := theme
	if query == "" {
		query = "creative inspiration"
	}
	
	return sep.kb.Search(ctx, query, options)
}

// GetKnowledgeBase returns the underlying knowledge base
func (sep *StubEnhancementProvider) GetKnowledgeBase() KnowledgeBase {
	return sep.kb
}