package ui

import (
	"sync"
	"time"
)

// CacheEntry holds a cached value with metadata
type CacheEntry[T any] struct {
	Value      T
	Expiration time.Time
}

// GenericCache is a type-safe generic cache with expiration
type GenericCache[K comparable, V any] struct {
	mu      sync.RWMutex
	entries map[K]CacheEntry[V]
	ttl     time.Duration
}

// NewGenericCache creates a new generic cache
func NewGenericCache[K comparable, V any](ttl time.Duration) *GenericCache[K, V] {
	return &GenericCache[K, V]{
		entries: make(map[K]CacheEntry[V]),
		ttl:     ttl,
	}
}

// Get retrieves a value from the cache
func (c *GenericCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[key]
	if !ok {
		var zero V
		return zero, false
	}

	// Check if expired
	if time.Now().After(entry.Expiration) {
		var zero V
		return zero, false
	}

	return entry.Value, true
}

// Set stores a value in the cache
func (c *GenericCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = CacheEntry[V]{
		Value:      value,
		Expiration: time.Now().Add(c.ttl),
	}
}

// Delete removes a value from the cache
func (c *GenericCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Clear removes all entries
func (c *GenericCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[K]CacheEntry[V])
}

// Cleanup removes expired entries
func (c *GenericCache[K, V]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for key, entry := range c.entries {
		if now.After(entry.Expiration) {
			delete(c.entries, key)
		}
	}
}

// Len returns the number of cached entries
func (c *GenericCache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
