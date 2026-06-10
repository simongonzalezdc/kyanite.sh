package knowledge

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/constants"
)

// StubKnowledgeBase provides a stub implementation of the KnowledgeBase interface
// It returns pre-defined songwriting knowledge and gracefully handles unavailability
type StubKnowledgeBase struct {
	cards      map[string]Card
	categories []string
	tags       []string
	available  bool
	status     KnowledgeStatus

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

// initializeSongwritingKnowledge populates the stub with comprehensive songwriting information
func (kb *StubKnowledgeBase) initializeSongwritingKnowledge() {
	now := time.Now()

	// Initializing phases
	kb.initHarmonicKnowledge(now)
	kb.initLyricalKnowledge(now)
	kb.initStructureKnowledge(now)
	kb.initRhythmicKnowledge(now)
	kb.initInspirationKnowledge(now)

	// Update categories and tags
	kb.updateCategoriesAndTags()
}

// initHarmonicKnowledge adds advanced harmonic and theory concepts
func (kb *StubKnowledgeBase) initHarmonicKnowledge(now time.Time) {
	// Royal Road Progression
	kb.addCard(Card{
		ID:        "progression-royal-road",
		Title:     "Royal Road Progression (IV-V-iii-vi)",
		Content:   "The IV-V-iii-vi progression, often called the 'Royal Road' or 'Oudou' progression in Japanese pop music. It provides a sense of continuous emotional lift and resolution. Widely used in anime themes and sophisticated pop for its balance of major and minor gravity.",
		Category:  "theory-harmony",
		Tags:      []string{"chords", "progression", "pop", "j-pop", "emotional"},
		Relevance: 0.95,
		Metadata: map[string]string{
			"example_c": "Fmaj7 - G - Em7 - Am",
			"feeling":   "hopeful, nostalgic, driving",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Andalusian Cadence
	kb.addCard(Card{
		ID:        "progression-andalusian",
		Title:     "Andalusian Cadence (i-VII-VI-V)",
		Content:   "A descending four-chord progression that uses modal interchange. Historically rooted in flamenco and classical music, it creates a sense of dark mystery and relentless movement. In a minor key, it descends from the tonic to the dominant.",
		Category:  "theory-harmony",
		Tags:      []string{"chords", "progression", "minor", "flamenco", "dark"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"example_am": "Am - G - F - E",
			"feeling":    "mysterious, tragic, intense",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Modal Interchange
	kb.addCard(Card{
		ID:        "theory-modal-interchange",
		Title:     "Modal Interchange (Borrowed Chords)",
		Content:   "The technique of borrowing chords from a parallel scale (e.g., borrowing from C Minor while playing in C Major). A classic example is the iv chord (Fm in C Major), which adds a 'bittersweet' or cinematic quality to pop songs.",
		Category:  "theory-harmony",
		Tags:      []string{"theory", "harmony", "advanced", "cinematic"},
		Relevance: 0.95,
		Metadata: map[string]string{
			"tip":     "Try using a flat-VI or b-VII chord to break out of a standard major key feel.",
			"example": "I - IV - iv - I (C - F - Fm - C)",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Tritone Substitution
	kb.addCard(Card{
		ID:        "theory-tritone-sub",
		Title:     "Tritone Substitution",
		Content:   "A jazz harmony technique where a dominant 7th chord is replaced by another dominant 7th chord whose root is a tritone away (e.g., substituting G7 with Db7 in the key of C). This creates chromatic bass movement and sophisticated tension.",
		Category:  "theory-harmony",
		Tags:      []string{"theory", "jazz", "harmony", "tension"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"tip":     "Use this to create smooth chromatic bass lines moving toward the tonic.",
			"example": "ii - bII7 - I (Dm7 - Db7 - Cmaj7)",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
}

// initLyricalKnowledge adds deep lyrical craft concepts
func (kb *StubKnowledgeBase) initLyricalKnowledge(now time.Time) {
	// Mosaic Rhymes
	kb.addCard(Card{
		ID:        "technique-mosaic-rhyme",
		Title:     "Mosaic Rhymes (Multi-word Rhymes)",
		Content:   "A sophisticated rhyming technique where a single multi-syllabic word is rhymed against multiple smaller words (e.g., 'Antigravity' rhymed with 'Plan to have it be'). Commonly used in rap and complex musical theater for rhythmic density.",
		Category:  "lyrical-techniques",
		Tags:      []string{"rhyme", "advanced", "prosody", "rap"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"tip": "Match the stressed syllables exactly across the multiple words for maximum impact.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Anaphora
	kb.addCard(Card{
		ID:        "technique-anaphora",
		Title:     "Anaphora (Repetition for Emphasis)",
		Content:   "The repetition of a word or phrase at the beginning of successive lines. It creates a sense of rhythmic urgency and builds emotional momentum. Think of 'I have a dream' or many folk protest songs.",
		Category:  "lyrical-techniques",
		Tags:      []string{"rhetoric", "structure", "emphasis", "momentum"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"example": "I want to feel the rain (A)\nI want to see the sun (A)\nI want to find my way (A)",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Metonymy
	kb.addCard(Card{
		ID:        "technique-metonymy",
		Title:     "Metonymy in Lyrics",
		Content:   "Using a related concept or object to represent something larger (e.g., 'The crown' for the monarchy or 'The street' for urban life). This allows for evocative, grounded descriptions without being literal.",
		Category:  "lyrical-techniques",
		Tags:      []string{"figurative", "descriptive", "imagery"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"tip": "Use objects that carry strong cultural or emotional baggage.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
}

// initStructureKnowledge adds advanced structural concepts
func (kb *StubKnowledgeBase) initStructureKnowledge(now time.Time) {
	// AABA Form
	kb.addCard(Card{
		ID:        "structure-aaba",
		Title:     "AABA Form (Thirty-two-bar form)",
		Content:   "A classic song structure common in jazz and early pop (the Great American Songbook). It consists of two verses (A), a bridge (B), and a final verse (A). This form focuses heavily on melodic development within a concise frame.",
		Category:  "song-structure",
		Tags:      []string{"structure", "jazz", "standard", "classic"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"example": "Verse 1 - Verse 2 - Bridge - Verse 3",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Pre-Chorus Tension
	kb.addCard(Card{
		ID:        "structure-pre-chorus",
		Title:     "Pre-Chorus Dynamic Build",
		Content:   "The pre-chorus serves to transition musically and emotionally from the verse to the chorus. It often uses rising melodic lines and increased rhythmic density (or sudden minimalist drops) to make the chorus 'explode' upon arrival.",
		Category:  "song-structure",
		Tags:      []string{"structure", "pop", "dynamics", "tension"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"tip": "Try starting the pre-chorus on a secondary dominant or a subdominant chord to pull the listener forward.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
}

// initRhythmicKnowledge adds advanced rhythmic concepts
func (kb *StubKnowledgeBase) initRhythmicKnowledge(now time.Time) {
	// Syncopation
	kb.addCard(Card{
		ID:        "theory-syncopation",
		Title:     "Syncopation and Weak Beats",
		Content:   "Rhythmic displacement that emphasizes the 'off-beats' or weak beats of a measure. Syncopation is the heart of groove in Funk, Jazz, and Latin music, creating a sense of danceable tension.",
		Category:  "theory-rhythm",
		Tags:      []string{"rhythm", "groove", "funk", "dance"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"tip": "Try placing a snare hit on the 'and' of beat 4 for a classic rhythmic push.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Neo-Soul Unquantized Feel
	kb.addCard(Card{
		ID:        "genre-neo-soul-groove",
		Title:     "Neo-Soul 'Dilla' Feel (Unquantized)",
		Content:   "A rhythmic style popularized by J Dilla, where drums (especially the kick and snare) are slightly 'behind the beat' or unquantized. This creates a human, swaying feel that is foundational to modern Neo-Soul and Lo-fi Hip-hop.",
		Category:  "theory-rhythm",
		Tags:      []string{"rhythm", "groove", "soul", "lo-fi", "modern"},
		Relevance: 0.95,
		Metadata: map[string]string{
			"tip": "Avoid strict grid alignment to achieve this 'lazy' but intentional pocket.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Polyrhythms
	kb.addCard(Card{
		ID:        "theory-polyrhythm",
		Title:     "Polyrhythms (3 against 4)",
		Content:   "The simultaneous use of two or more conflicting rhythms that are not readily perceived as deriving from one another. A common example is '3 against 4', which creates complex, overlapping textures.",
		Category:  "theory-rhythm",
		Tags:      []string{"rhythm", "advanced", "texture", "complexity"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"example": "Hi-hats playing triplets while the kick stays on the quarter note pulse.",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
}

// initInspirationKnowledge adds deep philosophical and thematic inspiration
func (kb *StubKnowledgeBase) initInspirationKnowledge(now time.Time) {
	// Surveillance Capitalism
	kb.addCard(Card{
		ID:        "theme-surveillance-caps",
		Title:     "Theme: Surveillance Capitalism",
		Content:   "A modern lyrical theme exploring the loss of privacy and the monetization of every human interaction by data-driven corporations. Ideal for dark synth-pop or cerebral indie songs.",
		Category:  "inspiration",
		Tags:      []string{"theme", "modern", "dystopian", "technology"},
		Relevance: 0.85,
		Metadata: map[string]string{
			"prompt": "How does it feel to be watched by a machine that knows your desires before you do?",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Post-Industrial Decay
	kb.addCard(Card{
		ID:        "theme-industrial-decay",
		Title:     "Theme: Post-Industrial Decay",
		Content:   "Exploring the aesthetics of abandoned factories, rusted machinery, and the silence of once-booming industrial centers. This theme evokes nostalgia for a lost physical reality in an increasingly digital world.",
		Category:  "inspiration",
		Tags:      []string{"theme", "nostalgia", "urban", "rust-belt"},
		Relevance: 0.8,
		Metadata: map[string]string{
			"feeling": "cold, heavy, silent, grounded",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})

	// Digital Connection vs Isolation
	kb.addCard(Card{
		ID:        "theme-digital-isolation",
		Title:     "Theme: The Paradox of Digital Connection",
		Content:   "Exploring how we are more 'connected' than ever through screens, yet often feel profoundly more isolated. A powerful theme for intimate ballads or high-energy glitch-pop.",
		Category:  "inspiration",
		Tags:      []string{"theme", "isolation", "connection", "modern-life"},
		Relevance: 0.9,
		Metadata: map[string]string{
			"prompt": "What's the loneliest thing about having 1,000 friends you've never met?",
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
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
	time.Sleep(constants.StubDelay)

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
			queryWords := strings.Fields(queryLower)
			if len(queryWords) == 0 {
				queryWords = []string{""}
			}

			anyWordMatch := false
			for _, word := range queryWords {
				if word == "" {
					anyWordMatch = true
					break
				}
				if strings.Contains(strings.ToLower(card.Title), word) ||
					strings.Contains(strings.ToLower(card.Content), word) ||
					strings.Contains(strings.ToLower(card.Category), word) {
					anyWordMatch = true
					break
				}
				for _, tag := range card.Tags {
					if strings.Contains(strings.ToLower(tag), word) {
						anyWordMatch = true
						break
					}
				}
				if anyWordMatch {
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

			if anyWordMatch && categoryMatchAllowed && tagMatchAllowed {
				matchingCards = append(matchingCards, card)
			}
		}
	}

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

	// Prefer a card from the requested categories if present in results
	var chosen *Card
	if len(result.Cards) > 0 && len(options.Categories) > 0 {
		for _, cat := range options.Categories {
			for _, rc := range result.Cards {
				if rc.Category == cat {
					tmp := rc
					chosen = &tmp
					break
				}
			}
			if chosen != nil {
				break
			}
		}
	}

	// If no preferred card found, fall back to first result (if any)
	if chosen == nil && len(result.Cards) > 0 {
		tmp := result.Cards[0]
		chosen = &tmp
	}

	// If still no result but category filters were provided, try to pick a representative
	if chosen == nil && len(options.Categories) > 0 {
		// Sort IDs for deterministic selection in tests
		cardIDs := make([]string, 0, len(sep.kb.cards))
		for id := range sep.kb.cards {
			cardIDs = append(cardIDs, id)
		}
		for i := 0; i < len(cardIDs); i++ {
			for j := i + 1; j < len(cardIDs); j++ {
				if cardIDs[i] > cardIDs[j] {
					cardIDs[i], cardIDs[j] = cardIDs[j], cardIDs[i]
				}
			}
		}

		for _, cat := range options.Categories {
			for _, id := range cardIDs {
				c := sep.kb.cards[id]
				if c.Category == cat {
					tmp := c
					chosen = &tmp
					break
				}
			}
			if chosen != nil {
				break
			}
		}
	}

	if chosen != nil {
		if chosen.Category == "lyrical-techniques" {
			suggestion.Suggestion = fmt.Sprintf("Consider applying %s to: %s", chosen.Title, lyrics)
			suggestion.Reason = chosen.Content
			suggestion.Confidence = 0.8
		} else if chosen.Category == "inspiration" {
			suggestion.Suggestion = fmt.Sprintf("Inspired by %s: %s", chosen.Title, lyrics)
			suggestion.Reason = chosen.Content
			suggestion.Confidence = 0.75
		} else {
			suggestion.Suggestion = fmt.Sprintf("Enhanced with %s insight: %s", chosen.Category, lyrics)
			suggestion.Reason = chosen.Content
		}
		// Ensure suggestion list reflects chosen card
		suggestion.Cards = []Card{*chosen}
	} else {
		// Fallback generic suggestion
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
		if card.Category == "theory-harmony" || card.Category == "theory-rhythm" || card.Category == "chord-progressions" {
			suggestion.Suggestion = card.Metadata["example_c"]
			if suggestion.Suggestion == "" {
				suggestion.Suggestion = card.Metadata["example_am"]
			}
			if suggestion.Suggestion == "" {
				suggestion.Suggestion = card.Metadata["example"]
			}
			if suggestion.Suggestion == "" {
				suggestion.Suggestion = fmt.Sprintf("Enhanced pattern: %s", pattern)
			}
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
