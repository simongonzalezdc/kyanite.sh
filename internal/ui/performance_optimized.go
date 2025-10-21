package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// PerformanceOptimizedUI provides UI components with performance optimizations
type PerformanceOptimizedUI struct {
	// Rendering optimization
	renderCache    *RenderCache
	frameRateLimiter *FrameRateLimiter
	animationPool  *AnimationPool
	
	// Theme optimization
	themeOptimizer *ThemeOptimizer
	
	// Performance monitoring
	uiMetrics      *UIMetrics
	
	// Configuration
	config UIPerformanceConfig
}

// UIPerformanceConfig defines UI performance optimization settings
type UIPerformanceConfig struct {
	MaxFrameRate       int           `json:"max_frame_rate"`
	EnableRenderCache  bool          `json:"enable_render_cache"`
	CacheMaxSize       int           `json:"cache_max_size"`
	EnableLazyLoading  bool          `json:"enable_lazy_loading"`
	AnimationPoolSize  int           `json:"animation_pool_size"`
	ThemePreloadCount int           `json:"theme_preload_count"`
	EnableMetrics      bool          `json:"enable_metrics"`
	RenderTimeout      time.Duration `json:"render_timeout"`
}

// RenderCache provides intelligent caching for rendered UI components
type RenderCache struct {
	entries map[string]*RenderEntry
	mutex   sync.RWMutex
	maxSize int
	
	// Cache statistics
	hits   int64
	misses int64
}

// RenderEntry represents a cached rendered component
type RenderEntry struct {
	Rendered string
	Width    int
	Height   int
	Theme    string
	Timestamp time.Time
	AccessCount int
}

// FrameRateLimiter controls the maximum frame rate for UI updates
type FrameRateLimiter struct {
	maxFrameRate int
	lastRender   time.Time
	minInterval  time.Duration
	mutex        sync.Mutex
}

// AnimationPool provides reusable animation components
type AnimationPool struct {
	spinners    []*AnimatedLoadingSpinner
	statusBars  []*AnimatedStatusBar
	notifications []*AnimatedNotification
	spinnerMutex sync.Mutex
	statusMutex  sync.Mutex
	notifMutex   sync.Mutex
	nextIndex    int
}

// ThemeOptimizer provides optimized theme switching and caching
type ThemeOptimizer struct {
	preloadedThemes map[string]theme.Theme
	currentTheme    string
	mutex           sync.RWMutex
	cacheSize       int
}

// UIMetrics tracks UI performance metrics
type UIMetrics struct {
	RenderCount      int64         `json:"render_count"`
	CacheHits        int64         `json:"cache_hits"`
	CacheMisses      int64         `json:"cache_misses"`
	AverageRenderTime time.Duration `json:"average_render_time"`
	ThemeSwitchTime  time.Duration `json:"theme_switch_time"`
	FrameDrops       int64         `json:"frame_drops"`
	mutex            sync.RWMutex
}

// NewPerformanceOptimizedUI creates a new performance-optimized UI manager
func NewPerformanceOptimizedUI(config UIPerformanceConfig) *PerformanceOptimizedUI {
	// Set defaults if not provided
	if config.MaxFrameRate == 0 {
		config.MaxFrameRate = 60 // 60 FPS target
	}
	if config.CacheMaxSize == 0 {
		config.CacheMaxSize = 500
	}
	if config.AnimationPoolSize == 0 {
		config.AnimationPoolSize = 10
	}
	if config.ThemePreloadCount == 0 {
		config.ThemePreloadCount = 3
	}
	if config.RenderTimeout == 0 {
		config.RenderTimeout = 16 * time.Millisecond // 60 FPS = 16.67ms per frame
	}

	ui := &PerformanceOptimizedUI{
		renderCache:      NewRenderCache(config.CacheMaxSize),
		frameRateLimiter: NewFrameRateLimiter(config.MaxFrameRate),
		animationPool:    NewAnimationPool(config.AnimationPoolSize),
		themeOptimizer:   NewThemeOptimizer(config.ThemePreloadCount),
		uiMetrics:        &UIMetrics{},
		config:           config,
	}

	// Preload common themes
	ui.themeOptimizer.preloadCommonThemes()

	return ui
}

// NewRenderCache creates a new render cache
func NewRenderCache(maxSize int) *RenderCache {
	return &RenderCache{
		entries: make(map[string]*RenderEntry),
		maxSize: maxSize,
	}
}

// NewFrameRateLimiter creates a new frame rate limiter
func NewFrameRateLimiter(maxFrameRate int) *FrameRateLimiter {
	minInterval := time.Second / time.Duration(maxFrameRate)
	return &FrameRateLimiter{
		maxFrameRate: maxFrameRate,
		minInterval:  minInterval,
		lastRender:   time.Time{},
	}
}

// NewAnimationPool creates a new animation pool
func NewAnimationPool(size int) *AnimationPool {
	pool := &AnimationPool{
		spinners:      make([]*AnimatedLoadingSpinner, 0, size),
		statusBars:    make([]*AnimatedStatusBar, 0, size),
		notifications: make([]*AnimatedNotification, 0, size),
		nextIndex:     0,
	}

	// Pre-allocate animation components
	for i := 0; i < size; i++ {
		pool.spinners = append(pool.spinners, NewAnimatedLoadingSpinner(""))
		pool.statusBars = append(pool.statusBars, NewAnimatedStatusBar(""))
		pool.notifications = append(pool.notifications, NewAnimatedNotification("", "info", 5*time.Second))
	}

	return pool
}

// NewThemeOptimizer creates a new theme optimizer
func NewThemeOptimizer(preloadCount int) *ThemeOptimizer {
	return &ThemeOptimizer{
		preloadedThemes: make(map[string]theme.Theme),
		currentTheme:    "",
		cacheSize:       preloadCount,
	}
}

// RenderOptimized renders a UI component with performance optimizations
func (ui *PerformanceOptimizedUI) RenderOptimized(component string, width, height int, themeName string) string {
	start := time.Now()
	
	// Update metrics
	ui.uiMetrics.mutex.Lock()
	ui.uiMetrics.RenderCount++
	ui.uiMetrics.mutex.Unlock()

	defer func() {
		duration := time.Since(start)
		ui.updateAverageRenderTime(duration)
	}()

	// Check render cache first if enabled
	if ui.config.EnableRenderCache {
		cacheKey := ui.generateRenderCacheKey(component, width, height, themeName)
		if cached, found := ui.renderCache.Get(cacheKey, width, height, themeName); found {
			ui.uiMetrics.mutex.Lock()
			ui.uiMetrics.CacheHits++
			ui.uiMetrics.mutex.Unlock()
			
			return cached
		}
		
		ui.uiMetrics.mutex.Lock()
		ui.uiMetrics.CacheMisses++
		ui.uiMetrics.mutex.Unlock()
	}

	// Apply frame rate limiting
	if !ui.frameRateLimiter.CanRender() {
		ui.uiMetrics.mutex.Lock()
		ui.uiMetrics.FrameDrops++
		ui.uiMetrics.mutex.Unlock()
		
		return "" // Skip this frame
	}

	// Render the component
	rendered := ui.renderComponent(component, width, height, themeName)

	// Cache the rendered result if enabled
	if ui.config.EnableRenderCache && rendered != "" {
		cacheKey := ui.generateRenderCacheKey(component, width, height, themeName)
		ui.renderCache.Set(cacheKey, rendered, width, height, themeName)
	}

	// Update frame rate limiter
	ui.frameRateLimiter.RecordRender()

	return rendered
}

// GetSpinnerFromPool gets a spinner from the pool
func (ui *PerformanceOptimizedUI) GetSpinnerFromPool(message string) *AnimatedLoadingSpinner {
	ui.animationPool.spinnerMutex.Lock()
	defer ui.animationPool.spinnerMutex.Unlock()

	if len(ui.animationPool.spinners) == 0 {
		// Pool exhausted, create new spinner
		return NewAnimatedLoadingSpinner(message)
	}

	// Get next spinner from pool
	spinner := ui.animationPool.spinners[ui.animationPool.nextIndex]
	ui.animationPool.nextIndex = (ui.animationPool.nextIndex + 1) % len(ui.animationPool.spinners)

	// Reset spinner with new message
	spinner.message = message
	return spinner
}

// GetStatusBarFromPool gets a status bar from the pool
func (ui *PerformanceOptimizedUI) GetStatusBarFromPool(message string) *AnimatedStatusBar {
	ui.animationPool.statusMutex.Lock()
	defer ui.animationPool.statusMutex.Unlock()

	if len(ui.animationPool.statusBars) == 0 {
		// Pool exhausted, create new status bar
		return NewAnimatedStatusBar(message)
	}

	// Get next status bar from pool
	statusBar := ui.animationPool.statusBars[ui.animationPool.nextIndex]
	ui.animationPool.nextIndex = (ui.animationPool.nextIndex + 1) % len(ui.animationPool.statusBars)

	// Reset status bar with new message
	statusBar.message = message
	statusBar.progress = 0.0
	return statusBar
}

// SwitchThemeOptimized switches themes with performance optimization
func (ui *PerformanceOptimizedUI) SwitchThemeOptimized(themeName string) time.Duration {
	start := time.Now()

	// Get theme from optimizer
	newTheme := ui.themeOptimizer.GetTheme(themeName)
	if newTheme == nil {
		// Fallback to default theme
		newTheme = ui.themeOptimizer.GetTheme("default")
	}

	// Apply theme atomically
	ui.themeOptimizer.SetCurrentTheme(themeName)

	// Clear render cache when theme changes
	if ui.config.EnableRenderCache {
		ui.renderCache.Clear()
	}

	duration := time.Since(start)
	
	// Update metrics
	ui.uiMetrics.mutex.Lock()
	ui.uiMetrics.ThemeSwitchTime = duration
	ui.uiMetrics.mutex.Unlock()

	return duration
}

// Get returns a cached rendered component if available
func (rc *RenderCache) Get(key string, width, height int, themeName string) (string, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	entry, exists := rc.entries[key]
	if !exists {
		return "", false
	}

	// Check if dimensions match
	if entry.Width != width || entry.Height != height || entry.Theme != themeName {
		return "", false
	}

	// Update access statistics
	entry.AccessCount++
	
	rc.hits++
	return entry.Rendered, true
}

// Set stores a rendered component in the cache
func (rc *RenderCache) Set(key string, rendered string, width, height int, themeName string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Check if we need to evict entries
	if len(rc.entries) >= rc.maxSize {
		rc.evictLRU()
	}

	rc.entries[key] = &RenderEntry{
		Rendered:    rendered,
		Width:       width,
		Height:      height,
		Theme:       themeName,
		Timestamp:   time.Now(),
		AccessCount: 1,
	}
}

// Clear clears all entries in the render cache
func (rc *RenderCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.entries = make(map[string]*RenderEntry)
}

// evictLRU removes the least recently used entries
func (rc *RenderCache) evictLRU() {
	if len(rc.entries) == 0 {
		return
	}

	// Find the least recently used entry
	var oldestKey string
	var oldestTime time.Time
	var lowestAccess int
	first := true

	for key, entry := range rc.entries {
		if first || entry.AccessCount < lowestAccess || 
		   (entry.AccessCount == lowestAccess && entry.Timestamp.Before(oldestTime)) {
			oldestKey = key
			oldestTime = entry.Timestamp
			lowestAccess = entry.AccessCount
			first = false
		}
	}

	if oldestKey != "" {
		delete(rc.entries, oldestKey)
	}
}

// CanRender checks if a new render is allowed based on frame rate limiting
func (frl *FrameRateLimiter) CanRender() bool {
	frl.mutex.Lock()
	defer frl.mutex.Unlock()

	now := time.Now()
	if now.Sub(frl.lastRender) >= frl.minInterval {
		return true
	}
	return false
}

// RecordRender records that a render has occurred
func (frl *FrameRateLimiter) RecordRender() {
	frl.mutex.Lock()
	defer frl.mutex.Unlock()

	frl.lastRender = time.Now()
}

// preloadCommonThemes preloads commonly used themes
func (to *ThemeOptimizer) preloadCommonThemes() {
	to.mutex.Lock()
	defer to.mutex.Unlock()

	// Preload common themes
	commonThemes := []string{"default", "dark", "light"}
	
	for _, themeName := range commonThemes {
		if themeValue := theme.GetTheme(themeName); themeValue.Name != "" {
			to.preloadedThemes[themeName] = themeValue
		}
	}
}

// GetTheme retrieves a theme from the optimizer
func (to *ThemeOptimizer) GetTheme(name string) *theme.Theme {
	to.mutex.RLock()
	defer to.mutex.RUnlock()

	if theme, exists := to.preloadedThemes[name]; exists {
		return &theme
	}

	// Fallback to getting theme from registry
	themeValue := theme.GetTheme(name)
	if themeValue.Name != "" && len(to.preloadedThemes) < to.cacheSize {
		// Cache this theme for future use
		to.preloadedThemes[name] = themeValue
		return &themeValue
	}

	return &themeValue
}

// SetCurrentTheme sets the current theme
func (to *ThemeOptimizer) SetCurrentTheme(name string) {
	to.mutex.Lock()
	defer to.mutex.Unlock()

	to.currentTheme = name
}

// GetCurrentTheme returns the current theme
func (to *ThemeOptimizer) GetCurrentTheme() string {
	to.mutex.RLock()
	defer to.mutex.RUnlock()

	return to.currentTheme
}

// generateRenderCacheKey creates a unique cache key for rendered components
func (ui *PerformanceOptimizedUI) generateRenderCacheKey(component string, width, height int, themeName string) string {
	return fmt.Sprintf("%s_%dx%d_%s", component, width, height, themeName)
}

// renderComponent renders a specific component (simplified implementation)
func (ui *PerformanceOptimizedUI) renderComponent(component string, width, height int, themeName string) string {
	// This is a simplified implementation
	// In a real application, this would render the actual component
	
	t := theme.GetManager().Current()
	if themeName != "" {
		if optTheme := ui.themeOptimizer.GetTheme(themeName); optTheme != nil {
			t = *optTheme
		}
	}

	switch component {
	case "loading":
		spinner := ui.GetSpinnerFromPool("Loading...")
		return spinner.View()
	case "status":
		status := ui.GetStatusBarFromPool("Ready")
		return status.View()
	default:
		style := lipgloss.NewStyle().
			Foreground(t.Text).
			Width(width).
			Height(height)
		return style.Render(component)
	}
}

// updateAverageRenderTime updates the rolling average render time
func (ui *PerformanceOptimizedUI) updateAverageRenderTime(duration time.Duration) {
	ui.uiMetrics.mutex.Lock()
	defer ui.uiMetrics.mutex.Unlock()

	// Simple rolling average calculation
	if ui.uiMetrics.AverageRenderTime == 0 {
		ui.uiMetrics.AverageRenderTime = duration
	} else {
		// Weighted average (90% old, 10% new)
		ui.uiMetrics.AverageRenderTime = time.Duration(
			float64(ui.uiMetrics.AverageRenderTime)*0.9 + float64(duration)*0.1,
		)
	}
}

// GetMetrics returns current UI performance metrics
func (ui *PerformanceOptimizedUI) GetMetrics() UIMetrics {
	ui.uiMetrics.mutex.RLock()
	defer ui.uiMetrics.mutex.RUnlock()

	return *ui.uiMetrics
}

// GetPerformanceReport returns a comprehensive performance report
func (ui *PerformanceOptimizedUI) GetPerformanceReport() map[string]interface{} {
	metrics := ui.GetMetrics()
	
	report := map[string]interface{}{
		"metrics": metrics,
		"config":  ui.config,
	}

	// Calculate cache hit rate
	if metrics.CacheHits+metrics.CacheMisses > 0 {
		report["cache_hit_rate"] = float64(metrics.CacheHits) / float64(metrics.CacheHits+metrics.CacheMisses)
	}

	// Calculate current FPS if we have render data
	if metrics.AverageRenderTime > 0 {
		report["current_fps"] = 1000 / metrics.AverageRenderTime.Milliseconds()
	}

	return report
}

// GetCacheStats returns render cache statistics
func (ui *PerformanceOptimizedUI) GetCacheStats() map[string]interface{} {
	ui.renderCache.mutex.RLock()
	defer ui.renderCache.mutex.RUnlock()

	total := ui.renderCache.hits + ui.renderCache.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(ui.renderCache.hits) / float64(total)
	}

	return map[string]interface{}{
		"entries":    len(ui.renderCache.entries),
		"max_size":   ui.renderCache.maxSize,
		"hits":       ui.renderCache.hits,
		"misses":     ui.renderCache.misses,
		"hit_rate":   hitRate,
	}
}

// Close cleans up resources
func (ui *PerformanceOptimizedUI) Close() error {
	// Clear caches
	if ui.renderCache != nil {
		ui.renderCache.Clear()
	}
	
	return nil
}