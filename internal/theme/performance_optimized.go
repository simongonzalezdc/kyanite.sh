package theme

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
)

// PerformanceOptimizedManager provides theme management with performance optimizations
type PerformanceOptimizedManager struct {
	// Core theme management
	*Manager

	// Performance optimization
	preloadedThemes map[string]Theme
	themeCache      map[string]*CachedTheme
	renderCache     map[string]*CachedRender
	mutex           sync.RWMutex

	// Async operations
	saveQueue chan *ThemeSaveRequest
	stopChan  chan struct{}

	// Performance metrics
	metrics *ThemeMetrics

	// Configuration
	config ThemePerformanceConfig
}

// ThemePerformanceConfig defines theme performance optimization settings
type ThemePerformanceConfig struct {
	PreloadThemes     []string      `json:"preload_themes"`
	CacheSize         int           `json:"cache_size"`
	EnableAsyncSave   bool          `json:"enable_async_save"`
	SaveBatchSize     int           `json:"save_batch_size"`
	SaveInterval      time.Duration `json:"save_interval"`
	EnableRenderCache bool          `json:"enable_render_cache"`
	RenderCacheSize   int           `json:"render_cache_size"`
	EnableMetrics     bool          `json:"enable_metrics"`
}

// CachedTheme represents a cached theme with metadata
type CachedTheme struct {
	Theme       Theme     `json:"theme"`
	LoadTime    time.Time `json:"load_time"`
	AccessCount int       `json:"access_count"`
	LastAccess  time.Time `json:"last_access"`
	Preloaded   bool      `json:"preloaded"`
}

// CachedRender represents a cached rendered theme component
type CachedRender struct {
	Rendered    string    `json:"rendered"`
	ThemeName   string    `json:"theme_name"`
	Component   string    `json:"component"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Hash        string    `json:"hash"`
	CreateTime  time.Time `json:"create_time"`
	AccessCount int       `json:"access_count"`
}

// ThemeSaveRequest represents an asynchronous theme save request
type ThemeSaveRequest struct {
	ThemeID string
	Timeout time.Duration
	Result  chan error
}

// ThemeMetrics tracks theme performance metrics
type ThemeMetrics struct {
	ThemeSwitches     int64         `json:"theme_switches"`
	AverageSwitchTime time.Duration `json:"average_switch_time"`
	CacheHits         int64         `json:"cache_hits"`
	CacheMisses       int64         `json:"cache_misses"`
	RenderCacheHits   int64         `json:"render_cache_hits"`
	RenderCacheMisses int64         `json:"render_cache_misses"`
	AsyncSaves        int64         `json:"async_saves"`
	FailedSaves       int64         `json:"failed_saves"`
	mutex             sync.RWMutex
}

// NewPerformanceOptimizedManager creates a new performance-optimized theme manager
func NewPerformanceOptimizedManager(config ThemePerformanceConfig) *PerformanceOptimizedManager {
	// Set defaults if not provided
	if config.CacheSize == 0 {
		config.CacheSize = 50
	}
	if config.SaveBatchSize == 0 {
		config.SaveBatchSize = 5
	}
	if config.SaveInterval == 0 {
		config.SaveInterval = 1 * time.Second
	}
	if config.RenderCacheSize == 0 {
		config.RenderCacheSize = 200
	}
	if len(config.PreloadThemes) == 0 {
		config.PreloadThemes = []string{"default", "dark", "light"}
	}

	manager := &PerformanceOptimizedManager{
		Manager:         GetManager(),
		preloadedThemes: make(map[string]Theme),
		themeCache:      make(map[string]*CachedTheme),
		renderCache:     make(map[string]*CachedRender),
		saveQueue:       make(chan *ThemeSaveRequest, 100),
		stopChan:        make(chan struct{}),
		metrics:         &ThemeMetrics{},
		config:          config,
	}

	// Preload common themes
	manager.preloadThemes()

	// Start background processes
	if config.EnableAsyncSave {
		go manager.startAsyncSaveWorker()
	}

	logging.GetDefaultLogger().Info("Performance-optimized theme manager initialized",
		"preload_themes", len(config.PreloadThemes),
		"cache_size", config.CacheSize,
		"async_save", config.EnableAsyncSave)

	return manager
}

// SetThemeOptimized sets the current theme with performance optimization
func (m *PerformanceOptimizedManager) SetThemeOptimized(id string) time.Duration {
	start := time.Now()

	// Update metrics
	m.metrics.mutex.Lock()
	m.metrics.ThemeSwitches++
	m.metrics.mutex.Unlock()

	// Check if theme is already cached
	if cached, found := m.getThemeFromCache(id); found {
		// Apply cached theme
		m.Manager.mu.Lock()
		m.Manager.current = cached.Theme
		m.Manager.mu.Unlock()

		// Update cache access
		cached.AccessCount++
		cached.LastAccess = time.Now()

		duration := time.Since(start)
		m.updateAverageSwitchTime(duration)

		logging.GetDefaultLogger().Debug("Theme switched from cache", "id", id, "duration", duration)
		return duration
	}

	// Load theme from registry
	theme := GetTheme(id)
	if theme.Name == "" {
		// Fallback to default theme
		theme = Default()
		id = "default"
	}

	// Cache the theme
	m.cacheTheme(id, theme)

	// Apply theme
	m.Manager.mu.Lock()
	m.Manager.current = theme
	m.Manager.mu.Unlock()

	// Save preference asynchronously if enabled
	if m.config.EnableAsyncSave {
		m.queueAsyncSave(id)
	} else {
		go func() {
			if err := m.SaveThemePreference(); err != nil {
				logging.GetDefaultLogger().Warn("Failed to save theme preference", "error", err)
			}
		}()
	}

	duration := time.Since(start)
	m.updateAverageSwitchTime(duration)

	logging.GetDefaultLogger().Debug("Theme switched and cached", "id", id, "duration", duration)
	return duration
}

// GetRenderedComponentOptimized gets a rendered component with caching
func (m *PerformanceOptimizedManager) GetRenderedComponentOptimized(component string, width, height int) string {
	if !m.config.EnableRenderCache {
		// Fallback to direct rendering
		return m.renderComponent(component, width, height)
	}

	currentTheme := m.Manager.Current()
	cacheKey := m.generateRenderCacheKey(currentTheme.Name, component, width, height)

	// Check render cache
	m.mutex.RLock()
	if cached, found := m.renderCache[cacheKey]; found {
		cached.AccessCount++
		m.mutex.RUnlock()

		m.metrics.mutex.Lock()
		m.metrics.RenderCacheHits++
		m.metrics.mutex.Unlock()

		return cached.Rendered
	}
	m.mutex.RUnlock()

	// Render component
	rendered := m.renderComponent(component, width, height)

	// Cache the rendered result
	m.mutex.Lock()
	if len(m.renderCache) >= m.config.RenderCacheSize {
		m.evictRenderCacheLRU()
	}

	m.renderCache[cacheKey] = &CachedRender{
		Rendered:    rendered,
		ThemeName:   currentTheme.Name,
		Component:   component,
		Width:       width,
		Height:      height,
		Hash:        cacheKey,
		CreateTime:  time.Now(),
		AccessCount: 1,
	}
	m.mutex.Unlock()

	m.metrics.mutex.Lock()
	m.metrics.RenderCacheMisses++
	m.metrics.mutex.Unlock()

	return rendered
}

// preloadThemes preloads commonly used themes
func (m *PerformanceOptimizedManager) preloadThemes() {
	for _, themeID := range m.config.PreloadThemes {
		if theme := GetTheme(themeID); theme.Name != "" {
			m.mutex.Lock()
			m.preloadedThemes[themeID] = theme
			m.themeCache[themeID] = &CachedTheme{
				Theme:       theme,
				LoadTime:    time.Now(),
				AccessCount: 0,
				LastAccess:  time.Now(),
				Preloaded:   true,
			}
			m.mutex.Unlock()

			logging.GetDefaultLogger().Debug("Preloaded theme", "id", themeID)
		}
	}
}

// getThemeFromCache retrieves a theme from cache
func (m *PerformanceOptimizedManager) getThemeFromCache(id string) (*CachedTheme, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cached, exists := m.themeCache[id]
	if !exists {
		m.metrics.mutex.Lock()
		m.metrics.CacheMisses++
		m.metrics.mutex.Unlock()
		return nil, false
	}

	m.metrics.mutex.Lock()
	m.metrics.CacheHits++
	m.metrics.mutex.Unlock()

	return cached, true
}

// cacheTheme stores a theme in cache
func (m *PerformanceOptimizedManager) cacheTheme(id string, theme Theme) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we need to evict entries
	if len(m.themeCache) >= m.config.CacheSize {
		m.evictThemeCacheLRU()
	}

	m.themeCache[id] = &CachedTheme{
		Theme:       theme,
		LoadTime:    time.Now(),
		AccessCount: 1,
		LastAccess:  time.Now(),
		Preloaded:   false,
	}
}

// evictThemeCacheLRU removes the least recently used theme from cache
func (m *PerformanceOptimizedManager) evictThemeCacheLRU() {
	if len(m.themeCache) == 0 {
		return
	}

	var oldestID string
	var oldestTime time.Time
	var lowestAccess int
	first := true

	for id, cached := range m.themeCache {
		if cached.Preloaded {
			continue // Don't evict preloaded themes
		}

		if first || cached.AccessCount < lowestAccess ||
			(cached.AccessCount == lowestAccess && cached.LastAccess.Before(oldestTime)) {
			oldestID = id
			oldestTime = cached.LastAccess
			lowestAccess = cached.AccessCount
			first = false
		}
	}

	if oldestID != "" {
		delete(m.themeCache, oldestID)
		logging.GetDefaultLogger().Debug("Evicted theme from cache", "id", oldestID)
	}
}

// evictRenderCacheLRU removes the least recently used render from cache
func (m *PerformanceOptimizedManager) evictRenderCacheLRU() {
	if len(m.renderCache) == 0 {
		return
	}

	var oldestKey string
	var oldestTime time.Time
	var lowestAccess int
	first := true

	for key, cached := range m.renderCache {
		if first || cached.AccessCount < lowestAccess ||
			(cached.AccessCount == lowestAccess && cached.CreateTime.Before(oldestTime)) {
			oldestKey = key
			oldestTime = cached.CreateTime
			lowestAccess = cached.AccessCount
			first = false
		}
	}

	if oldestKey != "" {
		delete(m.renderCache, oldestKey)
	}
}

// generateRenderCacheKey generates a unique key for render cache
func (m *PerformanceOptimizedManager) generateRenderCacheKey(themeName, component string, width, height int) string {
	return fmt.Sprintf("%s_%s_%dx%d", themeName, component, width, height)
}

// renderComponent renders a specific theme component (simplified)
func (m *PerformanceOptimizedManager) renderComponent(component string, width, height int) string {
	currentTheme := m.Manager.Current()

	// This is a simplified implementation
	// In a real application, this would render actual theme components
	switch component {
	case "header":
		return fmt.Sprintf("[%s] Header - %dx%d", currentTheme.Name, width, height)
	case "footer":
		return fmt.Sprintf("[%s] Footer - %dx%d", currentTheme.Name, width, height)
	case "status":
		return fmt.Sprintf("[%s] Status - %dx%d", currentTheme.Name, width, height)
	default:
		return fmt.Sprintf("[%s] %s - %dx%d", currentTheme.Name, component, width, height)
	}
}

// queueAsyncSave queues an asynchronous theme save request
func (m *PerformanceOptimizedManager) queueAsyncSave(themeID string) {
	request := &ThemeSaveRequest{
		ThemeID: themeID,
		Timeout: 5 * time.Second,
		Result:  make(chan error, 1),
	}

	select {
	case m.saveQueue <- request:
		m.metrics.mutex.Lock()
		m.metrics.AsyncSaves++
		m.metrics.mutex.Unlock()
	default:
		// Queue full, log warning
		logging.GetDefaultLogger().Warn("Async save queue full, dropping save request", "theme_id", themeID)
	}
}

// startAsyncSaveWorker starts the background async save worker
func (m *PerformanceOptimizedManager) startAsyncSaveWorker() {
	batch := make([]*ThemeSaveRequest, 0, m.config.SaveBatchSize)
	ticker := time.NewTicker(m.config.SaveInterval)
	defer ticker.Stop()

	for {
		select {
		case request := <-m.saveQueue:
			batch = append(batch, request)

			if len(batch) >= m.config.SaveBatchSize {
				m.processSaveBatch(batch)
				batch = batch[:0] // Reset batch
			}

		case <-ticker.C:
			if len(batch) > 0 {
				m.processSaveBatch(batch)
				batch = batch[:0] // Reset batch
			}

		case <-m.stopChan:
			// Process remaining requests before exiting
			if len(batch) > 0 {
				m.processSaveBatch(batch)
			}
			return
		}
	}
}

// processSaveBatch processes a batch of save requests
func (m *PerformanceOptimizedManager) processSaveBatch(batch []*ThemeSaveRequest) {
	// Group by theme ID to avoid duplicate saves
	themeMap := make(map[string][]*ThemeSaveRequest)
	for _, req := range batch {
		themeMap[req.ThemeID] = append(themeMap[req.ThemeID], req)
	}

	// Process each unique theme
	for themeID, requests := range themeMap {
		err := m.saveThemePreferenceSync(themeID)

		// Send result to all requesters
		for _, req := range requests {
			select {
			case req.Result <- err:
			default:
				// Requester not listening
			}
		}

		if err != nil {
			m.metrics.mutex.Lock()
			m.metrics.FailedSaves += int64(len(requests))
			m.metrics.mutex.Unlock()

			logging.GetDefaultLogger().Warn("Failed to save theme preference", "theme_id", themeID, "error", err)
		}
	}
}

// saveThemePreferenceSync saves theme preference synchronously
func (m *PerformanceOptimizedManager) saveThemePreferenceSync(themeID string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "noise")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configFile := filepath.Join(configDir, "theme.json")

	pref := ThemePreference{
		ThemeID: themeID,
	}

	data, err := json.MarshalIndent(pref, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// updateAverageSwitchTime updates the rolling average theme switch time
func (m *PerformanceOptimizedManager) updateAverageSwitchTime(duration time.Duration) {
	m.metrics.mutex.Lock()
	defer m.metrics.mutex.Unlock()

	// Simple rolling average calculation
	if m.metrics.AverageSwitchTime == 0 {
		m.metrics.AverageSwitchTime = duration
	} else {
		// Weighted average (90% old, 10% new)
		m.metrics.AverageSwitchTime = time.Duration(
			float64(m.metrics.AverageSwitchTime)*0.9 + float64(duration)*0.1,
		)
	}
}

// GetMetrics returns current theme performance metrics
func (m *PerformanceOptimizedManager) GetMetrics() *ThemeMetrics {
	m.metrics.mutex.RLock()
	defer m.metrics.mutex.RUnlock()

	return m.metrics
}

// GetPerformanceReport returns a comprehensive performance report
func (m *PerformanceOptimizedManager) GetPerformanceReport() map[string]interface{} {
	metrics := m.GetMetrics()

	m.mutex.RLock()
	themeCacheSize := len(m.themeCache)
	renderCacheSize := len(m.renderCache)
	preloadedCount := len(m.preloadedThemes)
	m.mutex.RUnlock()

	report := map[string]interface{}{
		"metrics": metrics,
		"config":  m.config,
		"cache_info": map[string]interface{}{
			"theme_cache_size":  themeCacheSize,
			"render_cache_size": renderCacheSize,
			"preloaded_themes":  preloadedCount,
		},
	}

	// Calculate cache hit rates
	if metrics.CacheHits+metrics.CacheMisses > 0 {
		report["theme_cache_hit_rate"] = float64(metrics.CacheHits) / float64(metrics.CacheHits+metrics.CacheMisses)
	}

	if metrics.RenderCacheHits+metrics.RenderCacheMisses > 0 {
		report["render_cache_hit_rate"] = float64(metrics.RenderCacheHits) / float64(metrics.RenderCacheHits+metrics.RenderCacheMisses)
	}

	return report
}

// ClearCaches clears all caches
func (m *PerformanceOptimizedManager) ClearCaches() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.themeCache = make(map[string]*CachedTheme)
	m.renderCache = make(map[string]*CachedRender)

	logging.GetDefaultLogger().Info("Theme caches cleared")
}

// Close cleans up resources
func (m *PerformanceOptimizedManager) Close() error {
	logging.GetDefaultLogger().Info("Performance-optimized theme manager shutting down")

	// Stop async save worker
	close(m.stopChan)

	// Clear caches
	m.ClearCaches()

	return nil
}
