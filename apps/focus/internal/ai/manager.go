package ai

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	ai "github.com/kyanite/ai"
)

// Manager handles AI interactions via BrainProvider.
type Manager struct {
	brainProvider *BrainProvider
	validator     *TaskValidator

	// Cache for AI responses with LRU eviction
	cache        map[string]cacheEntry
	cacheHits    map[string]time.Time // Track last access time for LRU
	cacheMaxSize int
	cacheMutex   sync.RWMutex
	cachePath    string
	cacheDirty   bool // Track if cache needs saving

	done chan struct{}
}

// cacheEntry represents a cached AI response
type cacheEntry struct {
	Response  any       `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	ExpiresAt time.Time `json:"expires_at"`
	AccessCnt int       `json:"access_count"` // Track access count for LRU
}

// ParsedTask represents the structured output from AI
type ParsedTask struct {
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline,omitempty"`
	Priority    string    `json:"priority"`
	Categories  []string  `json:"categories,omitempty"`
}

// New creates a new AI manager backed by BrainProvider.
func New() *Manager {
	manager := &Manager{
		brainProvider: NewBrainProvider(),
		validator:     NewTaskValidator(),
		cache:         make(map[string]cacheEntry),
		cacheHits:     make(map[string]time.Time),
		cacheMaxSize:  500,
		done:          make(chan struct{}),
	}

	// Set cache path
	home, err := os.UserHomeDir()
	if err != nil {
		manager.cachePath = "./ai_cache.json"
	} else {
		manager.cachePath = filepath.Join(home, ".focus", "ai_cache.json")
	}

	// Load cache if it exists
	manager.loadCache()

	// Start background cache saver (saves every 5 minutes)
	go manager.periodicCacheSaver()

	return manager
}

// Close stops background goroutines and flushes the cache.
func (m *Manager) Close() {
	close(m.done)
	if m.brainProvider != nil {
		m.brainProvider.Close()
	}
	m.cacheMutex.RLock()
	needsSave := m.cacheDirty
	m.cacheMutex.RUnlock()
	if needsSave {
		m.cacheMutex.Lock()
		if m.cacheDirty {
			m.saveCache()
			m.cacheDirty = false
		}
		m.cacheMutex.Unlock()
	}
}

// ParseTask converts natural language to structured task using AI
func (m *Manager) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	// Check cache first
	cacheKey := m.generateCacheKey("parse", input)
	if cached, found := m.getFromCache(cacheKey); found {
		if task, ok := cached.(*ParsedTask); ok {
			return task, nil
		}
	}

	// Try BrainProvider
	if m.brainProvider.IsAvailable() {
		result, err := m.brainProvider.ParseTask(ctx, input)
		if err == nil {
			if validated, ok := m.validateResponse(result); ok {
				m.saveToCache(cacheKey, validated)
				return validated, nil
			}
		}
	}

	// Fallback: rule-based parsing
	basicResult := m.basicParse(input)
	m.saveToCache(cacheKey, basicResult)
	return basicResult, nil
}

// SuggestTasks generates contextual task suggestions
func (m *Manager) SuggestTasks(ctx context.Context, existingTasks []string) ([]string, error) {
	// Create cache key from tasks
	taskStr := strings.Join(existingTasks, "|")
	cacheKey := m.generateCacheKey("suggest", taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if suggestions, ok := cached.([]string); ok {
			return suggestions, nil
		}
	}

	// TODO: Implement SuggestTasks via BrainProvider when available

	// Deterministic fallback suggestions
	return []string{
		"Review today's highest-priority task",
		"Break one blocked item into a smaller next step",
	}, nil
}

// SummarizeTasks generates a summary of tasks using AI
func (m *Manager) SummarizeTasks(ctx context.Context, tasks []string) (string, error) {
	// Create cache key from tasks
	taskStr := strings.Join(tasks, "|")
	cacheKey := m.generateCacheKey("summary", taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if summary, ok := cached.(string); ok {
			return summary, nil
		}
	}

	// TODO: Implement SummarizeTasks via BrainProvider when available

	// Fallback: basic summary
	summary := m.basicSummary(tasks)
	m.saveToCache(cacheKey, summary)
	return summary, nil
}

// ChatAssistant provides help and answers questions about tasks and app usage
func (m *Manager) ChatAssistant(ctx context.Context, question string, tasks []string) (string, error) {
	// Create cache key from question and tasks
	taskStr := strings.Join(tasks, "|")
	cacheKey := m.generateCacheKey("chat", question+taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if response, ok := cached.(string); ok {
			return response, nil
		}
	}

	// Try BrainProvider
	if m.brainProvider.IsAvailable() {
		result, err := m.brainProvider.ChatAssistant(ctx, question, tasks)
		if err == nil {
			m.saveToCache(cacheKey, result)
			return result, nil
		}
	}

	// Fallback response
	if strings.Contains(strings.ToLower(question), "help") || strings.Contains(strings.ToLower(question), "hi") {
		return "Hello! I'm your focus.sh AI assistant. I can help you with:\n• Task management and organization\n• Productivity tips\n• App usage guidance\n• Smart suggestions\n\nTry asking me about your tasks or how to use specific features!", nil
	}

	if strings.Contains(strings.ToLower(question), "task") {
		if len(tasks) > 0 {
			return fmt.Sprintf("I see you have %d tasks. I can help you organize, prioritize, or complete them. Try using 'focus inspire' for suggestions or 'focus list' to see all your tasks!", len(tasks)), nil
		}
		return "You don't have any tasks yet! Start by adding a mission with 'focus add \"your task here\"' and I can help you manage them.", nil
	}

	return fmt.Sprintf("AI is currently unavailable. You have %d tasks. Try 'focus --help' to see available commands!", len(tasks)), nil
}

// GetCrossAppContext retrieves recent context from other kyanite apps (syntax, noise, prism).
// Returns nil if the brain is unavailable — callers should treat this as best-effort.
func (m *Manager) GetCrossAppContext(ctx context.Context, limit int) []ai.CrossAppContext {
	brain := m.brainProvider.Brain()
	if brain == nil {
		return nil
	}
	contexts, err := brain.GetCrossAppContext(ctx, limit)
	if err != nil {
		return nil
	}
	return contexts
}

// IsOllamaAvailable checks if the AI backend is reachable.
// Kept for backward compatibility — now checks BrainProvider availability.
func (m *Manager) IsOllamaAvailable() bool {
	return m.brainProvider.IsAvailable()
}

// LaunchOllama is a no-op kept for backward compatibility.
// BrainProvider connects to NUCBox automatically.
func (m *Manager) LaunchOllama() error {
	return nil
}

// validateResponse delegates to TaskValidator
func (m *Manager) validateResponse(task *ParsedTask) (*ParsedTask, bool) {
	return m.validator.Validate(task)
}

// basicParse provides fallback parsing when AI fails
func (m *Manager) basicParse(input string) *ParsedTask {
	task := &ParsedTask{
		Description: input,
		Priority:    "medium",
	}

	// Simple keyword-based priority detection
	lowerInput := strings.ToLower(input)
	if strings.Contains(lowerInput, "urgent") || strings.Contains(lowerInput, "asap") ||
		strings.Contains(lowerInput, "critical") || strings.Contains(lowerInput, "emergency") {
		task.Priority = "high"
	} else if strings.Contains(lowerInput, "low priority") || strings.Contains(lowerInput, "when possible") ||
		strings.Contains(lowerInput, "sometime") {
		task.Priority = "low"
	}

	// Simple category detection
	var categories []string
	if strings.Contains(lowerInput, "work") || strings.Contains(lowerInput, "job") {
		categories = append(categories, "work")
	}
	if strings.Contains(lowerInput, "personal") || strings.Contains(lowerInput, "home") {
		categories = append(categories, "personal")
	}
	if strings.Contains(lowerInput, "meeting") {
		categories = append(categories, "meetings")
	}
	task.Categories = categories

	// Simple deadline detection
	if strings.Contains(lowerInput, "today") {
		task.Deadline = time.Now()
	} else if strings.Contains(lowerInput, "tomorrow") {
		task.Deadline = time.Now().AddDate(0, 0, 1)
	} else if strings.Contains(lowerInput, "next week") {
		task.Deadline = time.Now().AddDate(0, 0, 7)
	}

	// Truncate long descriptions
	words := strings.Fields(task.Description)
	if len(words) > 10 {
		task.Description = strings.Join(words[:10], " ") + "..."
	}

	return task
}

// basicSummary provides fallback summarization when AI fails
func (m *Manager) basicSummary(tasks []string) string {
	completed := 0
	pending := 0

	for _, task := range tasks {
		if strings.Contains(task, "completed") {
			completed++
		} else {
			pending++
		}
	}

	return fmt.Sprintf("You have %d tasks (%d completed, %d pending). Keep up the good work!",
		len(tasks), completed, pending)
}

// generateCacheKey creates a unique cache key for a request
func (m *Manager) generateCacheKey(operation, input string) string {
	key := fmt.Sprintf("%s:%s", operation, input)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// getFromCache retrieves a response from cache if available and not expired
func (m *Manager) getFromCache(key string) (interface{}, bool) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	entry, exists := m.cache[key]
	if !exists {
		return nil, false
	}

	// Check if entry is expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Update access time for LRU
	m.cacheHits[key] = time.Now()
	entry.AccessCnt++
	m.cache[key] = entry

	return entry.Response, true
}

// saveToCache stores a response in cache with expiration
func (m *Manager) saveToCache(key string, response interface{}) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	// Evict LRU entry if cache is full
	if len(m.cache) >= m.cacheMaxSize {
		m.evictLRUEntryUnsafe()
	}

	now := time.Now()
	m.cache[key] = cacheEntry{
		Response:  response,
		Timestamp: now,
		ExpiresAt: now.Add(24 * time.Hour), // Cache for 24 hours
		AccessCnt: 1,
	}
	m.cacheHits[key] = now
	m.cacheDirty = true // Mark cache as needing save
}

// evictLRUEntryUnsafe removes the least recently used cache entry
// Must be called with cacheMutex held
func (m *Manager) evictLRUEntryUnsafe() {
	if len(m.cache) == 0 {
		return
	}

	// Find LRU entry
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, accessTime := range m.cacheHits {
		if first || accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = accessTime
			first = false
		}
	}

	// Remove oldest entry
	if oldestKey != "" {
		delete(m.cache, oldestKey)
		delete(m.cacheHits, oldestKey)
	}
}

// loadCache loads cached responses from file
func (m *Manager) loadCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	data, err := os.ReadFile(m.cachePath)
	if err != nil {
		return
	}

	var cache map[string]cacheEntry
	if err := json.Unmarshal(data, &cache); err != nil {
		return
	}

	// Filter out expired entries and initialize cacheHits
	now := time.Now()
	for key, entry := range cache {
		if now.After(entry.ExpiresAt) {
			delete(cache, key)
		} else {
			m.cacheHits[key] = entry.Timestamp
		}
	}

	m.cache = cache
	m.cacheDirty = false
}

// periodicCacheSaver saves cache to disk every 5 minutes.
// It exits when the done channel is closed (via Close()).
func (m *Manager) periodicCacheSaver() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-m.done:
			return
		case <-ticker.C:
			m.cacheMutex.RLock()
			needsSave := m.cacheDirty
			m.cacheMutex.RUnlock()

			if needsSave {
				m.cacheMutex.Lock()
				if m.cacheDirty {
					m.saveCache()
					m.cacheDirty = false
				}
				m.cacheMutex.Unlock()
			}
		}
	}
}

// saveCache persists cache to file
func (m *Manager) saveCache() {
	// Ensure directory exists
	dir := filepath.Dir(m.cachePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}

	data, err := json.Marshal(m.cache)
	if err != nil {
		return
	}

	_ = os.WriteFile(m.cachePath, data, 0o644)
}
