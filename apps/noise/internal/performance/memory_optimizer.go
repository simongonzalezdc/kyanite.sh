package performance

import (
	"runtime"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/logging"
)

// MemoryOptimizer provides memory usage optimization and garbage collection tuning
type MemoryOptimizer struct {
	// Configuration
	config MemoryOptimizerConfig

	// Memory monitoring
	currentMemory MemoryMetrics
	peakMemory    MemoryMetrics
	memoryHistory []MemorySnapshot
	mutex         sync.RWMutex

	// Garbage collection control
	gcController *GCController

	// Object pooling
	stringPool *StringPool
	bufferPool *BufferPool

	// Memory pressure detection
	pressureDetector *MemoryPressureDetector

	// Background processes
	stopChan chan struct{}
}

// MemoryOptimizerConfig defines memory optimization configuration
type MemoryOptimizerConfig struct {
	EnableMonitoring     bool          `json:"enable_monitoring"`
	EnableGCOptimization bool          `json:"enable_gc_optimization"`
	EnableObjectPooling  bool          `json:"enable_object_pooling"`
	MemoryCheckInterval  time.Duration `json:"memory_check_interval"`
	GCThreshold          float64       `json:"gc_threshold"`
	PressureThreshold    float64       `json:"pressure_threshold"`
	MaxHistorySize       int           `json:"max_history_size"`
	StringPoolSize       int           `json:"string_pool_size"`
	BufferPoolSize       int           `json:"buffer_pool_size"`
	EnableAdaptiveGC     bool          `json:"enable_adaptive_gc"`
	TargetMemoryUsage    uint64        `json:"target_memory_usage"`
}

// MemorySnapshot represents a point-in-time memory snapshot
type MemorySnapshot struct {
	Timestamp time.Time      `json:"timestamp"`
	Memory    MemoryMetrics  `json:"memory"`
	Pressure  MemoryPressure `json:"pressure"`
	GCStats   GCStats        `json:"gc_stats"`
}

// MemoryPressure represents memory pressure level
type MemoryPressure int

const (
	PressureLow MemoryPressure = iota
	PressureMedium
	PressureHigh
	PressureCritical
)

// GCStats represents garbage collection statistics
type GCStats struct {
	NumGC         uint32        `json:"num_gc"`
	NumForcedGC   uint32        `json:"num_forced_gc"`
	GCCPUFraction float64       `json:"gc_cpu_fraction"`
	TotalPause    time.Duration `json:"total_pause"`
	LastGC        time.Time     `json:"last_gc"`
}

// GCController provides garbage collection control
type GCController struct {
	enabled           bool
	threshold         float64
	adaptiveMode      bool
	targetMemoryUsage uint64
	lastGC            time.Time
	gcInterval        time.Duration
	mutex             sync.Mutex
}

// StringPool provides reusable string objects
type StringPool struct {
	pool    []string
	mutex   sync.Mutex
	maxSize int
}

// BufferPool provides reusable byte buffer objects
type BufferPool struct {
	pool    [][]byte
	mutex   sync.Mutex
	maxSize int
	bufSize int
}

// MemoryPressureDetector detects memory pressure conditions
type MemoryPressureDetector struct {
	threshold    float64
	currentLevel MemoryPressure
	history      []float64
	maxHistory   int
	mutex        sync.RWMutex
}

// NewMemoryOptimizer creates a new memory optimizer
func NewMemoryOptimizer(config MemoryOptimizerConfig) *MemoryOptimizer {
	// Set defaults if not provided
	if config.MemoryCheckInterval == 0 {
		config.MemoryCheckInterval = 5 * time.Second
	}
	if config.GCThreshold == 0 {
		config.GCThreshold = 0.8 // 80% of available memory
	}
	if config.PressureThreshold == 0 {
		config.PressureThreshold = 0.9 // 90% of available memory
	}
	if config.MaxHistorySize == 0 {
		config.MaxHistorySize = 100
	}
	if config.StringPoolSize == 0 {
		config.StringPoolSize = 1000
	}
	if config.BufferPoolSize == 0 {
		config.BufferPoolSize = 100
	}
	if config.TargetMemoryUsage == 0 {
		config.TargetMemoryUsage = 100 * 1024 * 1024 // 100MB
	}

	optimizer := &MemoryOptimizer{
		config:           config,
		memoryHistory:    make([]MemorySnapshot, 0, config.MaxHistorySize),
		gcController:     NewGCController(config.GCThreshold, config.EnableAdaptiveGC, config.TargetMemoryUsage),
		stringPool:       NewStringPool(config.StringPoolSize),
		bufferPool:       NewBufferPool(config.BufferPoolSize, 64*1024), // 64KB buffers
		pressureDetector: NewMemoryPressureDetector(config.PressureThreshold, config.MaxHistorySize),
		stopChan:         make(chan struct{}),
	}

	return optimizer
}

// Start begins memory optimization
func (mo *MemoryOptimizer) Start() {
	if !mo.config.EnableMonitoring {
		logging.GetDefaultLogger().Info("Memory optimization disabled")
		return
	}

	logging.GetDefaultLogger().Info("Starting memory optimization",
		"check_interval", mo.config.MemoryCheckInterval,
		"gc_threshold", mo.config.GCThreshold,
		"pressure_threshold", mo.config.PressureThreshold)

	// Start memory monitoring goroutine
	go mo.monitoringLoop()

	// Start GC optimization if enabled
	if mo.config.EnableGCOptimization {
		go mo.gcOptimizationLoop()
	}
}

// Stop stops memory optimization
func (mo *MemoryOptimizer) Stop() {
	logging.GetDefaultLogger().Info("Stopping memory optimization")
	close(mo.stopChan)
}

// monitoringLoop monitors memory usage at regular intervals
func (mo *MemoryOptimizer) monitoringLoop() {
	ticker := time.NewTicker(mo.config.MemoryCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mo.checkMemoryUsage()
		case <-mo.stopChan:
			return
		}
	}
}

// gcOptimizationLoop optimizes garbage collection
func (mo *MemoryOptimizer) gcOptimizationLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mo.optimizeGC()
		case <-mo.stopChan:
			return
		}
	}
}

// checkMemoryUsage checks current memory usage and takes action if needed
func (mo *MemoryOptimizer) checkMemoryUsage() {
	memory := mo.collectMemoryMetrics()
	pressure := mo.pressureDetector.DetectPressure(memory)
	gcStats := mo.collectGCStats()

	// Update current memory
	mo.mutex.Lock()
	mo.currentMemory = memory

	// Update peak memory
	if memory.HeapAlloc > mo.peakMemory.HeapAlloc {
		mo.peakMemory = memory
	}

	// Add to history
	snapshot := MemorySnapshot{
		Timestamp: time.Now(),
		Memory:    memory,
		Pressure:  pressure,
		GCStats:   gcStats,
	}

	mo.memoryHistory = append(mo.memoryHistory, snapshot)
	if len(mo.memoryHistory) > mo.config.MaxHistorySize {
		mo.memoryHistory = mo.memoryHistory[1:]
	}
	mo.mutex.Unlock()

	// Take action based on memory pressure
	switch pressure {
	case PressureHigh:
		logging.GetDefaultLogger().Warn("High memory pressure detected", "heap_alloc", memory.HeapAlloc)
		mo.triggerGC()
		mo.clearCaches()

	case PressureCritical:
		logging.GetDefaultLogger().Error("Critical memory pressure detected", "heap_alloc", memory.HeapAlloc)
		mo.triggerGC()
		mo.clearCaches()
		mo.releasePools()
	}

	// Log memory usage periodically
	if len(mo.memoryHistory)%10 == 0 {
		logging.GetDefaultLogger().Debug("Memory usage",
			"heap_alloc", memory.HeapAlloc,
			"heap_sys", memory.HeapSys,
			"pressure", pressure)
	}
}

// optimizeGC optimizes garbage collection based on memory usage patterns
func (mo *MemoryOptimizer) optimizeGC() {
	if !mo.config.EnableGCOptimization {
		return
	}

	mo.mutex.RLock()
	currentMemory := mo.currentMemory
	mo.mutex.RUnlock()

	// Calculate memory usage ratio
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memoryRatio := float64(currentMemory.HeapAlloc) / float64(memStats.Sys)

	// Trigger GC if memory usage exceeds threshold
	if memoryRatio > mo.config.GCThreshold {
		mo.triggerGC()
	}

	// Adaptive GC tuning
	if mo.config.EnableAdaptiveGC {
		mo.gcController.AdaptiveTune(currentMemory, memoryRatio)
	}
}

// triggerGC triggers garbage collection
func (mo *MemoryOptimizer) triggerGC() {
	start := time.Now()
	runtime.GC()
	duration := time.Since(start)

	logging.GetDefaultLogger().Debug("Manual GC triggered", "duration", duration)

	mo.gcController.RecordGC()
}

// clearCaches clears object pools to free memory
func (mo *MemoryOptimizer) clearCaches() {
	if mo.config.EnableObjectPooling {
		mo.stringPool.Clear()
		mo.bufferPool.Clear()
	}

	logging.GetDefaultLogger().Debug("Memory caches cleared")
}

// releasePools releases object pools to free memory
func (mo *MemoryOptimizer) releasePools() {
	if mo.config.EnableObjectPooling {
		mo.stringPool.Release()
		mo.bufferPool.Release()
	}

	logging.GetDefaultLogger().Debug("Memory pools released")
}

// collectMemoryMetrics collects current memory metrics
func (mo *MemoryOptimizer) collectMemoryMetrics() MemoryMetrics {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return MemoryMetrics{
		Alloc:      memStats.Alloc,
		TotalAlloc: memStats.TotalAlloc,
		Sys:        memStats.Sys,
		HeapAlloc:  memStats.HeapAlloc,
		HeapSys:    memStats.HeapSys,
		HeapInuse:  memStats.HeapInuse,
	}
}

// collectGCStats collects garbage collection statistics
func (mo *MemoryOptimizer) collectGCStats() GCStats {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return GCStats{
		NumGC:         memStats.NumGC,
		NumForcedGC:   memStats.NumForcedGC,
		GCCPUFraction: memStats.GCCPUFraction,
		TotalPause:    time.Duration(memStats.PauseTotalNs),
		LastGC:        time.Unix(0, int64(memStats.LastGC)),
	}
}

// GetString gets a string from the pool
func (mo *MemoryOptimizer) GetString() string {
	if !mo.config.EnableObjectPooling {
		return ""
	}
	return mo.stringPool.Get()
}

// PutString returns a string to the pool
func (mo *MemoryOptimizer) PutString(s string) {
	if !mo.config.EnableObjectPooling {
		return
	}
	mo.stringPool.Put(s)
}

// GetBuffer gets a buffer from the pool
func (mo *MemoryOptimizer) GetBuffer() []byte {
	if !mo.config.EnableObjectPooling {
		return make([]byte, 64*1024)
	}
	return mo.bufferPool.Get()
}

// PutBuffer returns a buffer to the pool
func (mo *MemoryOptimizer) PutBuffer(buf []byte) {
	if !mo.config.EnableObjectPooling {
		return
	}
	mo.bufferPool.Put(buf)
}

// GetMemoryMetrics returns current memory metrics
func (mo *MemoryOptimizer) GetMemoryMetrics() MemoryMetrics {
	mo.mutex.RLock()
	defer mo.mutex.RUnlock()
	return mo.currentMemory
}

// GetPeakMemoryMetrics returns peak memory metrics
func (mo *MemoryOptimizer) GetPeakMemoryMetrics() MemoryMetrics {
	mo.mutex.RLock()
	defer mo.mutex.RUnlock()
	return mo.peakMemory
}

// GetMemoryHistory returns memory usage history
func (mo *MemoryOptimizer) GetMemoryHistory() []MemorySnapshot {
	mo.mutex.RLock()
	defer mo.mutex.RUnlock()

	history := make([]MemorySnapshot, len(mo.memoryHistory))
	copy(history, mo.memoryHistory)
	return history
}

// GetMemoryReport returns a comprehensive memory report
func (mo *MemoryOptimizer) GetMemoryReport() map[string]interface{} {
	current := mo.GetMemoryMetrics()
	peak := mo.GetPeakMemoryMetrics()
	history := mo.GetMemoryHistory()
	pressure := mo.pressureDetector.GetCurrentLevel()
	gcStats := mo.collectGCStats()

	report := map[string]interface{}{
		"current_memory": current,
		"peak_memory":    peak,
		"pressure_level": pressure,
		"gc_stats":       gcStats,
		"config":         mo.config,
	}

	// Calculate memory efficiency
	if current.Sys > 0 {
		report["memory_efficiency"] = float64(current.HeapAlloc) / float64(current.Sys)
	}

	// Add recent history samples
	if len(history) > 10 {
		report["recent_history"] = history[len(history)-10:]
	}

	return report
}

// NewGCController creates a new garbage collection controller
func NewGCController(threshold float64, adaptiveMode bool, targetMemoryUsage uint64) *GCController {
	return &GCController{
		enabled:           true,
		threshold:         threshold,
		adaptiveMode:      adaptiveMode,
		targetMemoryUsage: targetMemoryUsage,
		gcInterval:        30 * time.Second,
	}
}

// RecordGC records a garbage collection event
func (gc *GCController) RecordGC() {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.lastGC = time.Now()
}

// AdaptiveTune adaptively tunes GC parameters
func (gc *GCController) AdaptiveTune(current MemoryMetrics, memoryRatio float64) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()

	if !gc.adaptiveMode {
		return
	}

	// Adjust GC interval based on memory usage
	if memoryRatio > gc.threshold {
		// High memory usage - more frequent GC
		gc.gcInterval = 10 * time.Second
	} else if memoryRatio < gc.threshold*0.7 {
		// Low memory usage - less frequent GC
		gc.gcInterval = 60 * time.Second
	} else {
		// Medium memory usage - normal GC interval
		gc.gcInterval = 30 * time.Second
	}
}

// NewStringPool creates a new string pool
func NewStringPool(maxSize int) *StringPool {
	return &StringPool{
		pool:    make([]string, 0, maxSize),
		maxSize: maxSize,
	}
}

// Get gets a string from the pool
func (sp *StringPool) Get() string {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if len(sp.pool) == 0 {
		return ""
	}

	// Get last string from pool
	s := sp.pool[len(sp.pool)-1]
	sp.pool = sp.pool[:len(sp.pool)-1]
	return s
}

// Put returns a string to the pool
func (sp *StringPool) Put(s string) {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()

	if len(sp.pool) >= sp.maxSize {
		return // Pool full
	}

	// Reset string and add to pool
	sp.pool = append(sp.pool, "")
}

// Clear clears the string pool
func (sp *StringPool) Clear() {
	sp.mutex.Lock()
	defer sp.mutex.Unlock()
	sp.pool = sp.pool[:0]
}

// Release releases the string pool
func (sp *StringPool) Release() {
	sp.Clear()
}

// NewBufferPool creates a new buffer pool
func NewBufferPool(maxSize, bufSize int) *BufferPool {
	return &BufferPool{
		pool:    make([][]byte, 0, maxSize),
		maxSize: maxSize,
		bufSize: bufSize,
	}
}

// Get gets a buffer from the pool
func (bp *BufferPool) Get() []byte {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	if len(bp.pool) == 0 {
		return make([]byte, bp.bufSize)
	}

	// Get last buffer from pool
	buf := bp.pool[len(bp.pool)-1]
	bp.pool = bp.pool[:len(bp.pool)-1]
	return buf
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf []byte) {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()

	if len(bp.pool) >= bp.maxSize || cap(buf) != bp.bufSize {
		return // Pool full or wrong size
	}

	// Clear buffer and add to pool
	for i := range buf {
		buf[i] = 0
	}
	bp.pool = append(bp.pool, buf)
}

// Clear clears the buffer pool
func (bp *BufferPool) Clear() {
	bp.mutex.Lock()
	defer bp.mutex.Unlock()
	bp.pool = bp.pool[:0]
}

// Release releases the buffer pool
func (bp *BufferPool) Release() {
	bp.Clear()
}

// NewMemoryPressureDetector creates a new memory pressure detector
func NewMemoryPressureDetector(threshold float64, maxHistory int) *MemoryPressureDetector {
	return &MemoryPressureDetector{
		threshold:  threshold,
		history:    make([]float64, 0, maxHistory),
		maxHistory: maxHistory,
	}
}

// DetectPressure detects memory pressure based on current memory usage
func (mpd *MemoryPressureDetector) DetectPressure(memory MemoryMetrics) MemoryPressure {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	memoryRatio := float64(memory.HeapAlloc) / float64(memStats.Sys)

	// Add to history
	mpd.mutex.Lock()
	mpd.history = append(mpd.history, memoryRatio)
	if len(mpd.history) > mpd.maxHistory {
		mpd.history = mpd.history[1:]
	}

	// Calculate average memory ratio
	var avgRatio float64
	if len(mpd.history) > 0 {
		sum := 0.0
		for _, ratio := range mpd.history {
			sum += ratio
		}
		avgRatio = sum / float64(len(mpd.history))
	}

	// Determine pressure level
	var pressure MemoryPressure
	switch {
	case avgRatio >= 0.95:
		pressure = PressureCritical
	case avgRatio >= 0.85:
		pressure = PressureHigh
	case avgRatio >= 0.7:
		pressure = PressureMedium
	default:
		pressure = PressureLow
	}

	mpd.currentLevel = pressure
	mpd.mutex.Unlock()

	return pressure
}

// GetCurrentLevel returns the current memory pressure level
func (mpd *MemoryPressureDetector) GetCurrentLevel() MemoryPressure {
	mpd.mutex.RLock()
	defer mpd.mutex.RUnlock()
	return mpd.currentLevel
}
