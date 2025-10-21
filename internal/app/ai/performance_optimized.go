package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// PerformanceOptimizedAI provides AI services with caching and performance optimizations
type PerformanceOptimizedAI struct {
	// Core AI components
	contextDetector *ContextDetector
	quickAgent      *QuickIdeaAgent
	knowledgeBase   interface{} // Using interface{} to avoid import issues

	// Performance optimization
	responseCache   *ResponseCache
	batchProcessor  *BatchProcessor
	metrics         *AIMetrics
	
	// Configuration
	config OptimizationConfig
}

// OptimizationConfig defines AI performance optimization settings
type OptimizationConfig struct {
	CacheEnabled          bool          `json:"cache_enabled"`
	CacheMaxSize          int           `json:"cache_max_size"`
	CacheTTL              time.Duration `json:"cache_ttl"`
	BatchProcessing       bool          `json:"batch_processing"`
	MaxConcurrentRequests int           `json:"max_concurrent_requests"`
	RequestTimeout        time.Duration `json:"request_timeout"`
	EnableMetrics         bool          `json:"enable_metrics"`
}

// ResponseCache provides intelligent caching for AI responses
type ResponseCache struct {
	entries map[string]*CacheEntry
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
	
	// Cache statistics
	hits   int64
	misses int64
}

// CacheEntry represents a cached AI response
type CacheEntry struct {
	Response    interface{} `json:"response"`
	Timestamp   time.Time   `json:"timestamp"`
	AccessCount int         `json:"access_count"`
	LastAccess  time.Time   `json:"last_access"`
	ContentType string      `json:"content_type"`
}

// BatchProcessor handles batch processing of AI requests
type BatchProcessor struct {
	requestQueue chan *BatchRequest
	workers      int
	mutex        sync.Mutex
	batchSize    int
	batchTimeout time.Duration
}

// BatchRequest represents a request to be processed in a batch
type BatchRequest struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Data     map[string]interface{} `json:"data"`
	Response chan *BatchResponse    `json:"-"`
	Context  context.Context        `json:"-"`
}

// BatchResponse represents the result of a batched request
type BatchResponse struct {
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   error       `json:"error"`
	Latency time.Duration `json:"latency"`
}

// AIMetrics tracks AI performance metrics
type AIMetrics struct {
	TotalRequests       int64         `json:"total_requests"`
	CacheHits           int64         `json:"cache_hits"`
	CacheMisses         int64         `json:"cache_misses"`
	AverageResponseTime time.Duration `json:"average_response_time"`
	FailedRequests      int64         `json:"failed_requests"`
	ConcurrentRequests  int64         `json:"concurrent_requests"`
	mutex               sync.RWMutex
}

// NewPerformanceOptimizedAI creates a new performance-optimized AI service
func NewPerformanceOptimizedAI(config OptimizationConfig) *PerformanceOptimizedAI {
	// Set defaults if not provided
	if config.CacheMaxSize == 0 {
		config.CacheMaxSize = 1000
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 30 * time.Minute
	}
	if config.MaxConcurrentRequests == 0 {
		config.MaxConcurrentRequests = 10
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = 10 * time.Second
	}

	ai := &PerformanceOptimizedAI{
		contextDetector: NewContextDetector(),
		quickAgent:      NewQuickIdeaAgent(),
		responseCache:   NewResponseCache(config.CacheMaxSize, config.CacheTTL),
		batchProcessor:  NewBatchProcessor(config.BatchProcessing),
		metrics:         &AIMetrics{},
		config:          config,
	}

	// Start background processes
	if config.CacheEnabled {
		go ai.responseCache.startCleanup()
	}

	if config.BatchProcessing {
		go ai.batchProcessor.start()
	}

	logging.GetDefaultLogger().Info("Performance-optimized AI service initialized",
		"cache_enabled", config.CacheEnabled,
		"cache_max_size", config.CacheMaxSize,
		"batch_processing", config.BatchProcessing,
		"max_concurrent", config.MaxConcurrentRequests)

	return ai
}

// NewResponseCache creates a new response cache
func NewResponseCache(maxSize int, ttl time.Duration) *ResponseCache {
	return &ResponseCache{
		entries: make(map[string]*CacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}
}

// GenerateWithContextOptimized generates AI responses with caching and optimization
func (ai *PerformanceOptimizedAI) GenerateWithContextOptimized(ctx context.Context, contentType ContentType, mode QuickIdeaMode, contentText string, options map[string]string) (interface{}, error) {
	start := time.Now()
	
	// Update metrics
	ai.metrics.mutex.Lock()
	ai.metrics.TotalRequests++
	ai.metrics.ConcurrentRequests++
	ai.metrics.mutex.Unlock()

	defer func() {
		ai.metrics.mutex.Lock()
		ai.metrics.ConcurrentRequests--
		ai.metrics.mutex.Unlock()
		
		latency := time.Since(start)
		ai.updateAverageResponseTime(latency)
	}()

	// Check cache first if enabled
	if ai.config.CacheEnabled {
		cacheKey := ai.generateCacheKey(contentType, mode, contentText, options)
		if cached, found := ai.responseCache.Get(cacheKey); found {
			ai.metrics.mutex.Lock()
			ai.metrics.CacheHits++
			ai.metrics.mutex.Unlock()
			
			logging.GetDefaultLogger().Debug("AI response cache hit", "key", cacheKey[:16]+"...")
			return cached, nil
		}
		
		ai.metrics.mutex.Lock()
		ai.metrics.CacheMisses++
		ai.metrics.mutex.Unlock()
	}

	// Apply request timeout
	if ai.config.RequestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ai.config.RequestTimeout)
		defer cancel()
	}

	// Generate response
	var response interface{}
	var err error

	// Use context-aware generation
	prompts := NewContextAwarePrompts()
	prompts.Initialize()
	
	switch mode {
	case QuickIdeaModeUnstick:
		response, err = ai.generateUnstickResponse(ctx, contentType, contentText, options, prompts)
	case QuickIdeaModeSpark:
		response, err = ai.generateSparkResponse(ctx, contentType, contentText, options, prompts)
	case QuickIdeaModeTweak:
		response, err = ai.generateTweakResponse(ctx, contentType, contentText, options, prompts)
	case QuickIdeaModeCheck:
		response, err = ai.generateCheckResponse(ctx, contentType, contentText, options, prompts)
	default:
		err = fmt.Errorf("unsupported mode: %s", mode)
	}

	if err != nil {
		ai.metrics.mutex.Lock()
		ai.metrics.FailedRequests++
		ai.metrics.mutex.Unlock()
		
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	// Cache the response if enabled
	if ai.config.CacheEnabled && response != nil {
		cacheKey := ai.generateCacheKey(contentType, mode, contentText, options)
		ai.responseCache.Set(cacheKey, response, string(contentType))
	}

	latency := time.Since(start)
	logging.GetDefaultLogger().Debug("AI response generated", 
		"mode", mode, 
		"content_type", contentType,
		"latency", latency,
		"cached", false)

	return response, nil
}

// generateCacheKey creates a unique cache key for AI requests
func (ai *PerformanceOptimizedAI) generateCacheKey(contentType ContentType, mode QuickIdeaMode, contentText string, options map[string]string) string {
	data := map[string]interface{}{
		"content_type": contentType,
		"mode":         mode,
		"context":      contentText,
		"options":      options,
	}
	
	jsonBytes, _ := json.Marshal(data)
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}

// Get returns a cached response if available and not expired
func (rc *ResponseCache) Get(key string) (interface{}, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	entry, exists := rc.entries[key]
	if !exists {
		return nil, false
	}

	// Check TTL
	if time.Since(entry.Timestamp) > rc.ttl {
		// Entry expired, will be cleaned up by background process
		return nil, false
	}

	// Update access statistics
	entry.AccessCount++
	entry.LastAccess = time.Now()
	
	rc.hits++
	return entry.Response, true
}

// Set stores a response in the cache
func (rc *ResponseCache) Set(key string, response interface{}, contentType string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Check if we need to evict entries
	if len(rc.entries) >= rc.maxSize {
		rc.evictLRU()
	}

	rc.entries[key] = &CacheEntry{
		Response:    response,
		Timestamp:   time.Now(),
		AccessCount: 1,
		LastAccess:  time.Now(),
		ContentType: contentType,
	}
}

// evictLRU removes the least recently used entries
func (rc *ResponseCache) evictLRU() {
	if len(rc.entries) == 0 {
		return
	}

	// Find the least recently used entry
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, entry := range rc.entries {
		if first || entry.LastAccess.Before(oldestTime) {
			oldestKey = key
			oldestTime = entry.LastAccess
			first = false
		}
	}

	if oldestKey != "" {
		delete(rc.entries, oldestKey)
	}
}

// startCleanup starts background cleanup of expired entries
func (rc *ResponseCache) startCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rc.cleanup()
	}
}

// cleanup removes expired entries
func (rc *ResponseCache) cleanup() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	now := time.Now()
	for key, entry := range rc.entries {
		if now.Sub(entry.Timestamp) > rc.ttl {
			delete(rc.entries, key)
		}
	}
}

// GetCacheStats returns cache performance statistics
func (rc *ResponseCache) GetCacheStats() map[string]interface{} {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	total := rc.hits + rc.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(rc.hits) / float64(total)
	}

	return map[string]interface{}{
		"entries":    len(rc.entries),
		"max_size":   rc.maxSize,
		"hits":       rc.hits,
		"misses":     rc.misses,
		"hit_rate":   hitRate,
		"ttl_minutes": rc.ttl.Minutes(),
	}
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(enabled bool) *BatchProcessor {
	if !enabled {
		return &BatchProcessor{
			requestQueue: nil,
			workers:      0,
		}
	}

	return &BatchProcessor{
		requestQueue: make(chan *BatchRequest, 100),
		workers:      3,
		batchSize:    5,
		batchTimeout: 100 * time.Millisecond,
	}
}

// start begins the batch processing workers
func (bp *BatchProcessor) start() {
	if bp.requestQueue == nil {
		return
	}

	for i := 0; i < bp.workers; i++ {
		go bp.worker(i)
	}
}

// worker processes batches of requests
func (bp *BatchProcessor) worker(id int) {
	batch := make([]*BatchRequest, 0, bp.batchSize)
	ticker := time.NewTicker(bp.batchTimeout)
	defer ticker.Stop()

	for {
		select {
		case request := <-bp.requestQueue:
			batch = append(batch, request)
			
			if len(batch) >= bp.batchSize {
				bp.processBatch(batch)
				batch = batch[:0] // Reset batch
			}

		case <-ticker.C:
			if len(batch) > 0 {
				bp.processBatch(batch)
				batch = batch[:0] // Reset batch
			}
		}
	}
}

// processBatch processes a batch of requests
func (bp *BatchProcessor) processBatch(batch []*BatchRequest) {
	start := time.Now()
	
	for _, request := range batch {
		if request == nil {
			continue
		}
		// Process each request in the batch
		// This is a simplified implementation
		response := &BatchResponse{
			ID:      request.ID,
			Result:  map[string]interface{}{"batch_processed": true},
			Error:   nil,
			Latency: time.Since(start),
		}
		
		if request.Response != nil {
			select {
			case request.Response <- response:
			case <-request.Context.Done():
				// Request timed out
			}
		}
	}
}

// updateAverageResponseTime updates the rolling average response time
func (ai *PerformanceOptimizedAI) updateAverageResponseTime(latency time.Duration) {
	ai.metrics.mutex.Lock()
	defer ai.metrics.mutex.Unlock()

	// Simple rolling average calculation
	if ai.metrics.AverageResponseTime == 0 {
		ai.metrics.AverageResponseTime = latency
	} else {
		// Weighted average (90% old, 10% new)
		ai.metrics.AverageResponseTime = time.Duration(
			float64(ai.metrics.AverageResponseTime)*0.9 + float64(latency)*0.1,
		)
	}
}

// GetMetrics returns current AI performance metrics
func (ai *PerformanceOptimizedAI) GetMetrics() AIMetrics {
	ai.metrics.mutex.RLock()
	defer ai.metrics.mutex.RUnlock()

	return *ai.metrics
}

// GetPerformanceReport returns a comprehensive performance report
func (ai *PerformanceOptimizedAI) GetPerformanceReport() map[string]interface{} {
	metrics := ai.GetMetrics()
	cacheStats := ai.responseCache.GetCacheStats()

	report := map[string]interface{}{
		"metrics": metrics,
		"cache":   cacheStats,
		"config":  ai.config,
	}

	// Calculate cache hit rate
	if cacheStats["hits"].(int64)+cacheStats["misses"].(int64) > 0 {
		report["cache_hit_rate"] = cacheStats["hit_rate"].(float64)
	}

	return report
}

// Helper methods for different AI modes
func (ai *PerformanceOptimizedAI) generateUnstickResponse(ctx context.Context, contentType ContentType, contentText string, options map[string]string, prompts *ContextAwarePrompts) (interface{}, error) {
	prompt := prompts.RenderPrompt(contentType, QuickIdeaModeUnstick, contentText, options)
	
	// Create a mock response for demonstration
	// In real implementation, this would call the actual AI service
	return map[string]interface{}{
		"type":    "unstick",
		"prompt":  prompt,
		"content": "Here are some ideas to get you unstuck...",
	}, nil
}

func (ai *PerformanceOptimizedAI) generateSparkResponse(ctx context.Context, contentType ContentType, contentText string, options map[string]string, prompts *ContextAwarePrompts) (interface{}, error) {
	prompt := prompts.RenderPrompt(contentType, QuickIdeaModeSpark, contentText, options)
	
	return map[string]interface{}{
		"type":    "spark",
		"prompt":  prompt,
		"content": "Here are some creative sparks...",
	}, nil
}

func (ai *PerformanceOptimizedAI) generateTweakResponse(ctx context.Context, contentType ContentType, contentText string, options map[string]string, prompts *ContextAwarePrompts) (interface{}, error) {
	prompt := prompts.RenderPrompt(contentType, QuickIdeaModeTweak, contentText, options)
	
	return map[string]interface{}{
		"type":    "tweak",
		"prompt":  prompt,
		"content": "Here are some tweaks to improve your content...",
	}, nil
}

func (ai *PerformanceOptimizedAI) generateCheckResponse(ctx context.Context, contentType ContentType, contentText string, options map[string]string, prompts *ContextAwarePrompts) (interface{}, error) {
	prompt := prompts.RenderPrompt(contentType, QuickIdeaModeCheck, contentText, options)
	
	return map[string]interface{}{
		"type":    "check",
		"prompt":  prompt,
		"content": "Here's an analysis of your content...",
	}, nil
}

// SetKnowledgeBase sets the knowledge base provider
func (ai *PerformanceOptimizedAI) SetKnowledgeBase(kb interface{}) {
	ai.knowledgeBase = kb
	if ai.quickAgent != nil {
		// This would need to be implemented in the actual QuickIdeaAgent
		// ai.quickAgent.WithKnowledgeBase(kb)
	}
}

// Close cleans up resources
func (ai *PerformanceOptimizedAI) Close() error {
	logging.GetDefaultLogger().Info("Performance-optimized AI service shutting down")
	
	// Close batch processor
	if ai.batchProcessor != nil && ai.batchProcessor.requestQueue != nil {
		close(ai.batchProcessor.requestQueue)
	}
	
	return nil
}