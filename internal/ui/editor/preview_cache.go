package editor

import (
	"sync"
	"time"
)

// PreviewCache manages cached preview content for performance
type PreviewCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	maxSize int
	hits    int64
	misses  int64
}

// CacheEntry represents a single cached preview
type CacheEntry struct {
	Content     string
	Rendered    string
	Timestamp   time.Time
	AccessCount int
}

// NewPreviewCache creates a new preview cache
func NewPreviewCache(maxSize int) *PreviewCache {
	return &PreviewCache{
		entries: make(map[string]*CacheEntry),
		maxSize: maxSize,
	}
}

// Get retrieves a cached preview
func (c *PreviewCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if entry, ok := c.entries[key]; ok {
		c.hits++
		entry.AccessCount++
		entry.Timestamp = time.Now()
		return entry.Rendered, true
	}

	c.misses++
	return "", false
}

// Set stores a preview in the cache
func (c *PreviewCache) Set(key, content, rendered string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Evict if cache is full
	if len(c.entries) >= c.maxSize {
		c.evictLRU()
	}

	c.entries[key] = &CacheEntry{
		Content:     content,
		Rendered:    rendered,
		Timestamp:   time.Now(),
		AccessCount: 1,
	}
}

// evictLRU evicts the least recently used entry
func (c *PreviewCache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range c.entries {
		if first || entry.Timestamp.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.Timestamp
			first = false
		}
	}

	if oldestKey != "" {
		delete(c.entries, oldestKey)
	}
}

// Clear removes all cached entries
func (c *PreviewCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}

// Stats returns cache statistics
func (c *PreviewCache) Stats() (hits, misses int64, size int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.hits, c.misses, len(c.entries)
}
