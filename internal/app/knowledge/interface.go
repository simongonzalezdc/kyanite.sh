package knowledge

import (
	"context"
	"time"
)

// Card represents a knowledge card with songwriting information
type Card struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Category  string            `json:"category"`
	Tags      []string          `json:"tags"`
	Relevance float64           `json:"relevance"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// SearchResult represents the result of a knowledge base search
type SearchResult struct {
	Cards     []Card        `json:"cards"`
	Query     string        `json:"query"`
	Total     int           `json:"total"`
	Duration  time.Duration `json:"duration"`
	FromCache bool          `json:"from_cache"`
}

// SearchOptions configures how searches are performed
type SearchOptions struct {
	Limit        int      `json:"limit"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
	MinRelevance float64  `json:"min_relevance"`
	UseCache     bool     `json:"use_cache"`
}

// KnowledgeStatus represents the current status of the knowledge base
type KnowledgeStatus struct {
	Available    bool          `json:"available"`
	CardCount    int           `json:"card_count"`
	LastSync     time.Time     `json:"last_sync"`
	Version      string        `json:"version"`
	Error        string        `json:"error"`
	ResponseTime time.Duration `json:"response_time"`
}

// KnowledgeBase defines the interface for knowledge base operations
type KnowledgeBase interface {
	// Search performs a search query against the knowledge base
	Search(ctx context.Context, query string, options SearchOptions) (*SearchResult, error)

	// AddCard adds a new card to the knowledge base
	AddCard(ctx context.Context, card Card) error

	// GetCard retrieves a specific card by ID
	GetCard(ctx context.Context, id string) (*Card, error)

	// UpdateCard updates an existing card
	UpdateCard(ctx context.Context, card Card) error

	// DeleteCard removes a card from the knowledge base
	DeleteCard(ctx context.Context, id string) error

	// GetCategories returns all available categories
	GetCategories(ctx context.Context) ([]string, error)

	// GetTags returns all available tags
	GetTags(ctx context.Context) ([]string, error)

	// GetStatus returns the current status of the knowledge base
	GetStatus(ctx context.Context) *KnowledgeStatus

	// IsAvailable returns whether the knowledge base is currently available
	IsAvailable(ctx context.Context) bool

	// Initialize sets up the knowledge base
	Initialize(ctx context.Context) error

	// Close performs cleanup operations
	Close(ctx context.Context) error
}

// LyricSuggestion represents a suggested lyric improvement
type LyricSuggestion struct {
	Original   string  `json:"original"`
	Suggestion string  `json:"suggestion"`
	Reason     string  `json:"reason"`
	Cards      []Card  `json:"cards"`
	Confidence float64 `json:"confidence"`
}

// PatternSuggestion represents a suggested musical pattern
type PatternSuggestion struct {
	Original    string  `json:"original"`
	Suggestion  string  `json:"suggestion"`
	Reason      string  `json:"reason"`
	Cards       []Card  `json:"cards"`
	Confidence  float64 `json:"confidence"`
	PatternType string  `json:"pattern_type"`
}

// EnhancementProvider defines the interface for knowledge-based enhancements
type EnhancementProvider interface {
	// EnhanceLyrics provides lyric suggestions based on knowledge
	EnhanceLyrics(ctx context.Context, lyrics string, options SearchOptions) (*LyricSuggestion, error)

	// EnhancePatterns provides pattern suggestions based on knowledge
	EnhancePatterns(ctx context.Context, pattern string, options SearchOptions) (*PatternSuggestion, error)

	// GetInspirationCards returns cards for creative inspiration
	GetInspirationCards(ctx context.Context, theme string, options SearchOptions) (*SearchResult, error)

	// GetKnowledgeBase returns the underlying knowledge base
	GetKnowledgeBase() KnowledgeBase
}
