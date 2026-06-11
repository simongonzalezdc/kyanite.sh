// Package cache provides a simple file-backed LRU cache for AI responses.
package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Entry represents a cached response.
type Entry struct {
	Key       string      `json:"key"`
	Value     interface{} `json:"value"`
	CreatedAt time.Time   `json:"created_at"`
}

// LRU is a simple in-memory LRU cache with file persistence.
type LRU struct {
	mu         sync.Mutex
	entries    map[string]*Entry
	order      []string // insertion order for LRU eviction
	maxEntries int
	ttl        time.Duration
	filePath   string
	done       chan struct{}
}

// NewLRU creates a new LRU cache. If filePath is non-empty, the cache
// is loaded from and saved to that file.
func NewLRU(maxEntries int, ttl time.Duration, filePath string) *LRU {
	c := &LRU{
		entries:    make(map[string]*Entry),
		order:      make([]string, 0, maxEntries),
		maxEntries: maxEntries,
		ttl:        ttl,
		filePath:   filePath,
		done:       make(chan struct{}),
	}
	if filePath != "" {
		c.load()
		go c.periodicSave()
	}
	return c
}

// Get retrieves a cached value by key. Returns nil if not found or expired.
func (c *LRU) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	if time.Since(entry.CreatedAt) > c.ttl {
		delete(c.entries, key)
		return nil, false
	}
	return entry.Value, true
}

// Set stores a value with the given key.
func (c *LRU) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.entries[key]; exists {
		// Update existing entry and move to front of order
		c.entries[key] = &Entry{
			Key:       key,
			Value:     value,
			CreatedAt: time.Now(),
		}
		// Move to front of order for LRU
		for i, k := range c.order {
			if k == key {
				c.order = append(c.order[:i], c.order[i+1:]...)
				break
			}
		}
		c.order = append(c.order, key)
		return
	}

	if len(c.entries) >= c.maxEntries {
		c.evictOldest()
	}

	c.entries[key] = &Entry{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
	}
	c.order = append(c.order, key)
}

// Clear removes all entries from the cache.
func (c *LRU) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[string]*Entry)
	c.order = make([]string, 0, c.maxEntries)
}

// Close stops the background saver and flushes to disk.
func (c *LRU) Close() {
	close(c.done)
	if c.filePath != "" {
		c.save()
	}
}

// GenerateKey creates a cache key from operation and input.
func GenerateKey(operation, input string) string {
	h := md5.Sum([]byte(operation + ":" + input))
	return hex.EncodeToString(h[:])
}

func (c *LRU) evictOldest() {
	if len(c.order) == 0 {
		return
	}
	oldest := c.order[0]
	c.order = c.order[1:]
	delete(c.entries, oldest)
}

func (c *LRU) load() {
	data, err := os.ReadFile(c.filePath)
	if err != nil {
		return
	}
	var entries []*Entry
	if err := json.Unmarshal(data, &entries); err != nil {
		return
	}
	for _, e := range entries {
		if time.Since(e.CreatedAt) > c.ttl {
			continue
		}
		c.entries[e.Key] = e
		c.order = append(c.order, e.Key)
	}
}

func (c *LRU) save() {
	c.mu.Lock()
	defer c.mu.Unlock()

	var entries []*Entry
	for _, e := range c.entries {
		entries = append(entries, e)
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(c.filePath), 0755)
	_ = os.WriteFile(c.filePath, data, 0644)
}

func (c *LRU) periodicSave() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-c.done:
			return
		case <-ticker.C:
			c.save()
		}
	}
}
