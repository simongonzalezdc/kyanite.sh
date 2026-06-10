package ai

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/app/knowledge"
)

// MockLLMClient implements QuickLLMClient for testing
type MockLLMClient struct {
	responses       map[string]string
	defaultResponse string
	shouldFail      bool
	delay           time.Duration
	shouldTimeout   bool
}

// NewMockLLMClient creates a new mock LLM client
func NewMockLLMClient() *MockLLMClient {
	return &MockLLMClient{
		responses:       make(map[string]string),
		defaultResponse: "1. suggestion one\n2. suggestion two\n3. suggestion three",
		shouldFail:      false,
		delay:           0,
		shouldTimeout:   false,
	}
}

// SetResponse sets a response for a specific prompt pattern
func (m *MockLLMClient) SetResponse(pattern, response string) {
	m.responses[pattern] = response
}

// SetDefaultResponse sets the default response for unmatched prompts
func (m *MockLLMClient) SetDefaultResponse(response string) {
	m.defaultResponse = response
}

// SetFailure sets whether the client should fail
func (m *MockLLMClient) SetFailure(shouldFail bool) {
	m.shouldFail = shouldFail
}

// SetDelay sets the delay for the client response
func (m *MockLLMClient) SetDelay(delay time.Duration) {
	m.delay = delay
}

// SetTimeout sets whether the client should timeout (by delaying longer than typical test timeouts)
func (m *MockLLMClient) SetTimeout(shouldTimeout bool) {
	m.shouldTimeout = shouldTimeout
	if shouldTimeout {
		m.delay = 5 * time.Second
	}
}

// Generate implements QuickLLMClient
func (m *MockLLMClient) Generate(ctx context.Context, _, prompt string, options map[string]any) (string, error) {
	if m.shouldFail {
		return "", errors.New("mock client error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	// Check for specific pattern responses
	for pattern, response := range m.responses {
		if strings.Contains(strings.ToLower(prompt), strings.ToLower(pattern)) {
			return response, nil
		}
	}

	// Return different responses based on prompt content
	promptLower := strings.ToLower(prompt)
	switch {
	case strings.Contains(promptLower, "lyric") && strings.Contains(promptLower, "continuation"):
		return "1. and the stars begin to fall\n2. while the city sleeps below\n3. as the morning light breaks through", nil
	case strings.Contains(promptLower, "pattern") && strings.Contains(promptLower, "continuation"):
		return "1. Am - G - C - F\n2. I - V - vi - IV\n3. Dm - G - C - Am", nil
	case strings.Contains(promptLower, "spark") && strings.Contains(promptLower, "lyric"):
		return "1. In the heart of creativity, I found my way\n2. creativity whispers through the window pane\n3. Chasing creativity through the pouring rain", nil
	case strings.Contains(promptLower, "spark") && strings.Contains(promptLower, "pattern"):
		return "1. creativity theme: C - G - Am - F progression\n2. creativity rhythm: driving 4/4 with syncopation\n3. creativity mood: minor key with descending bassline", nil
	case strings.Contains(promptLower, "tweak") && strings.Contains(promptLower, "lyric"):
		return "1. Rewrite with stronger imagery and emotion\n2. Replace clichés with fresh, specific details\n3. Enhance the emotional resonance", nil
	case strings.Contains(promptLower, "tweak") && strings.Contains(promptLower, "pattern"):
		return "1. Add sophisticated voice leading\n2. Incorporate rhythmic variation\n3. Enhance harmonic movement", nil
	case strings.Contains(promptLower, "check") && strings.Contains(promptLower, "lyric"):
		return "STRONG\nAdd vivid sensory details", nil
	case strings.Contains(promptLower, "check") && strings.Contains(promptLower, "pattern"):
		return "STRONG\nStrengthen harmonic resolution", nil
	default:
		return m.defaultResponse, nil
	}
}

// MockKnowledgeBase implements knowledge.KnowledgeBase for testing
type MockKnowledgeBase struct {
	cards      map[string]*knowledge.Card
	available  bool
	categories []string
	tags       []string
	delay      time.Duration
	shouldFail bool
}

// NewMockKnowledgeBase creates a new mock knowledge base
func NewMockKnowledgeBase() *MockKnowledgeBase {
	return &MockKnowledgeBase{
		cards:      make(map[string]*knowledge.Card),
		available:  true,
		categories: []string{"lyrical-techniques", "chord-progressions", "inspiration", "song-structure"},
		tags:       []string{"rhyme", "metaphor", "harmony", "rhythm", "emotion"},
		delay:      0,
		shouldFail: false,
	}
}

// AddTestCard adds a test card to the knowledge base
func (m *MockKnowledgeBase) AddTestCard(id, title, content, category string, tags []string) {
	m.cards[id] = &knowledge.Card{
		ID:        id,
		Title:     title,
		Content:   content,
		Category:  category,
		Tags:      tags,
		Relevance: 0.8,
		Metadata:  map[string]string{"example_c": "C - G - Am - F"},
	}
}

// SetAvailable sets whether the knowledge base is available
func (m *MockKnowledgeBase) SetAvailable(available bool) {
	m.available = available
}

// SetDelay sets the delay for operations
func (m *MockKnowledgeBase) SetDelay(delay time.Duration) {
	m.delay = delay
}

// SetFailure sets whether operations should fail
func (m *MockKnowledgeBase) SetFailure(shouldFail bool) {
	m.shouldFail = shouldFail
}

// Search implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) Search(ctx context.Context, query string, options knowledge.SearchOptions) (*knowledge.SearchResult, error) {
	if m.shouldFail {
		return nil, errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	var results []knowledge.Card
	queryLower := strings.ToLower(query)

	for _, card := range m.cards {
		// Simple matching logic
		if strings.Contains(strings.ToLower(card.Content), queryLower) ||
			strings.Contains(strings.ToLower(card.Title), queryLower) {

			// Filter by category if specified
			if len(options.Categories) > 0 {
				found := false
				for _, cat := range options.Categories {
					if card.Category == cat {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			// Filter by relevance if specified
			if card.Relevance < options.MinRelevance {
				continue
			}

			results = append(results, *card)
		}
	}

	// Limit results
	if options.Limit > 0 && len(results) > options.Limit {
		results = results[:options.Limit]
	}

	return &knowledge.SearchResult{
		Cards:     results,
		Query:     query,
		Total:     len(results),
		Duration:  time.Millisecond * 10,
		FromCache: options.UseCache,
	}, nil
}

// AddCard implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) AddCard(ctx context.Context, card knowledge.Card) error {
	if m.shouldFail {
		return errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	m.cards[card.ID] = &card
	return nil
}

// GetCard implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) GetCard(ctx context.Context, id string) (*knowledge.Card, error) {
	if m.shouldFail {
		return nil, errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	card, exists := m.cards[id]
	if !exists {
		return nil, errors.New("card not found")
	}
	return card, nil
}

// UpdateCard implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) UpdateCard(ctx context.Context, card knowledge.Card) error {
	if m.shouldFail {
		return errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	if _, exists := m.cards[card.ID]; !exists {
		return errors.New("card not found")
	}
	m.cards[card.ID] = &card
	return nil
}

// DeleteCard implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) DeleteCard(ctx context.Context, id string) error {
	if m.shouldFail {
		return errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	delete(m.cards, id)
	return nil
}

// GetCategories implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) GetCategories(ctx context.Context) ([]string, error) {
	if m.shouldFail {
		return nil, errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	return m.categories, nil
}

// GetTags implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) GetTags(ctx context.Context) ([]string, error) {
	if m.shouldFail {
		return nil, errors.New("mock knowledge base error")
	}

	if m.delay > 0 {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(m.delay):
			// Continue after delay
		}
	}

	return m.tags, nil
}

// GetStatus implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) GetStatus(ctx context.Context) *knowledge.KnowledgeStatus {
	return &knowledge.KnowledgeStatus{
		Available:    m.available,
		CardCount:    len(m.cards),
		LastSync:     time.Now(),
		Version:      "test-1.0.0",
		Error:        "",
		ResponseTime: time.Millisecond * 5,
	}
}

// IsAvailable implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) IsAvailable(ctx context.Context) bool {
	return m.available
}

// Initialize implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) Initialize(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("mock knowledge base error")
	}
	return nil
}

// Close implements knowledge.KnowledgeBase
func (m *MockKnowledgeBase) Close(ctx context.Context) error {
	if m.shouldFail {
		return errors.New("mock knowledge base error")
	}
	return nil
}

// MockEnhancementProvider implements knowledge.EnhancementProvider for testing
type MockEnhancementProvider struct {
	kb *MockKnowledgeBase
}

// NewMockEnhancementProvider creates a new mock enhancement provider
func NewMockEnhancementProvider() *MockEnhancementProvider {
	return &MockEnhancementProvider{
		kb: NewMockKnowledgeBase(),
	}
}

// SetKnowledgeBase sets the underlying knowledge base
func (m *MockEnhancementProvider) SetKnowledgeBase(kb *MockKnowledgeBase) {
	m.kb = kb
}

// GetKnowledgeBase returns the underlying knowledge base
func (m *MockEnhancementProvider) GetKnowledgeBase() knowledge.KnowledgeBase {
	return m.kb
}

// EnhanceLyrics implements knowledge.EnhancementProvider
func (m *MockEnhancementProvider) EnhanceLyrics(ctx context.Context, lyrics string, options knowledge.SearchOptions) (*knowledge.LyricSuggestion, error) {
	if m.kb.shouldFail {
		return nil, errors.New("mock enhancement provider error")
	}

	result, err := m.kb.Search(ctx, lyrics, options)
	if err != nil {
		return nil, err
	}

	suggestion := &knowledge.LyricSuggestion{
		Original:   lyrics,
		Suggestion: "Enhanced: " + lyrics + " with vivid imagery",
		Reason:     "Added sensory details and emotional depth",
		Cards:      result.Cards,
		Confidence: 0.8,
	}

	return suggestion, nil
}

// EnhancePatterns implements knowledge.EnhancementProvider
func (m *MockEnhancementProvider) EnhancePatterns(ctx context.Context, pattern string, options knowledge.SearchOptions) (*knowledge.PatternSuggestion, error) {
	if m.kb.shouldFail {
		return nil, errors.New("mock enhancement provider error")
	}

	result, err := m.kb.Search(ctx, pattern, options)
	if err != nil {
		return nil, err
	}

	suggestion := &knowledge.PatternSuggestion{
		Original:    pattern,
		Suggestion:  "Enhanced: " + pattern + " with sophisticated voice leading",
		Reason:      "Improved harmonic movement and rhythmic interest",
		Cards:       result.Cards,
		Confidence:  0.8,
		PatternType: "chord-progression",
	}

	return suggestion, nil
}

// GetInspirationCards implements knowledge.EnhancementProvider
func (m *MockEnhancementProvider) GetInspirationCards(ctx context.Context, theme string, options knowledge.SearchOptions) (*knowledge.SearchResult, error) {
	return m.kb.Search(ctx, theme, options)
}

// TestHelper provides utility functions for AI testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertNoError asserts that an error is nil
func (h *TestHelper) AssertNoError(err error) {
	if err != nil {
		h.t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError asserts that an error is not nil
func (h *TestHelper) AssertError(err error) {
	if err == nil {
		h.t.Fatal("Expected error, got nil")
	}
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		h.t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func (h *TestHelper) AssertNotEqual(expected, actual interface{}) {
	if expected == actual {
		h.t.Fatalf("Expected %v to not equal %v", expected, actual)
	}
}

// AssertTrue asserts that a condition is true
func (h *TestHelper) AssertTrue(condition bool, msgAndArgs ...interface{}) {
	if !condition {
		if len(msgAndArgs) == 0 {
			h.t.Fatal("Expected true, got false")
		} else if len(msgAndArgs) == 1 {
			h.t.Fatalf("Expected true, got false: %s", msgAndArgs[0])
		} else {
			format := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			h.t.Fatalf("Expected true, got false: "+format, args...)
		}
	}
}

// AssertFalse asserts that a condition is false
func (h *TestHelper) AssertFalse(condition bool, msg string) {
	if condition {
		h.t.Fatalf("Expected false, got true: %s", msg)
	}
}

// AssertNotEmpty asserts that a string is not empty
func (h *TestHelper) AssertNotEmpty(s string) {
	if s == "" {
		h.t.Fatal("Expected non-empty string")
	}
}

// AssertContains asserts that a string contains a substring
func (h *TestHelper) AssertContains(s, substr string) {
	if !strings.Contains(s, substr) {
		h.t.Fatalf("Expected '%s' to contain '%s'", s, substr)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func (h *TestHelper) AssertNotContains(s, substr string) {
	if strings.Contains(s, substr) {
		h.t.Fatalf("Expected '%s' to not contain '%s'", s, substr)
	}
}

// AssertNotNil asserts that a value is not nil
func (h *TestHelper) AssertNotNil(value interface{}) {
	if value == nil {
		h.t.Fatal("Expected non-nil value")
	}
}

// AssertLength asserts that a slice has the expected length
func (h *TestHelper) AssertLength(slice interface{}, expected int) {
	switch v := slice.(type) {
	case []string:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []knowledge.Card:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []knowledge.LyricSuggestion:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []knowledge.PatternSuggestion:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	default:
		h.t.Fatalf("Unsupported slice type for AssertLength")
	}
}
