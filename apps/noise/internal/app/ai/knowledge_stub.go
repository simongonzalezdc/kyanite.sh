package ai

import "context"

// EnhancementProvider is a minimal stub replacing the deleted knowledge package.
type EnhancementProvider interface {
	GetSuggestions(context string) []string
	GetInspirationCards(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
	GetKnowledgeBase() *KnowledgeBase
}

type stubProvider struct{}

func newStubProvider() EnhancementProvider {
	return &stubProvider{}
}

func (s *stubProvider) GetSuggestions(ctx string) []string {
	return []string{"Try a different chord progression", "Experiment with syncopation"}
}

func (s *stubProvider) GetKnowledgeBase() *KnowledgeBase {
	return &KnowledgeBase{}
}

func (s *stubProvider) GetInspirationCards(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
	return &SearchResult{Cards: []Card{}, Total: 0, Query: query}, nil
}

// Types replacing the deleted knowledge package
type Card struct {
	ID         string
	Title      string
	Content    string
	Category   string
	Tags       []string
	Relevance  float64
	Metadata   map[string]string
}

type KnowledgeStatus struct {
	TotalCards  int
	Categories  []string
	LastUpdated string
	Available   bool
	Error       string
	CardCount   int
}

type SearchOptions struct {
	Query       string
	Category    string
	Categories  []string
	Limit       int
	MinRelevance float64
	UseCache    bool
}

type SearchResult struct {
	Cards []Card
	Total int
	Query string
}

type LyricSuggestion struct {
	Line     string
	Metadata map[string]string
}

type PatternSuggestion struct {
	Name     string
	Pattern  string
	Metadata map[string]string
}

type knowledgeEnhancedContext struct {
	context string
	cards   []Card
	tip     string
}

// KnowledgeBase stub for test compatibility
type KnowledgeBase struct{}

func (kb *KnowledgeBase) GetStatus(ctx context.Context) *KnowledgeStatus {
	return &KnowledgeStatus{Available: false}
}

func (kb *KnowledgeBase) IsAvailable(ctx context.Context) bool {
	return false
}
