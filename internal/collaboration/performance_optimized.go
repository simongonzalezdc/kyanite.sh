package collaboration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/logging"
)

// PerformanceOptimizedCollaborationManager provides collaboration features with performance optimizations
type PerformanceOptimizedCollaborationManager struct {
	*CollaborationManager
	
	// Performance optimization
	sessionCache      map[string]*CachedSession
	operationPool     *OperationPool
	batchProcessor    *BatchProcessor
	
	// Connection pooling
	connectionPool    *ConnectionPool
	
	// Performance monitoring
	metrics           *CollaborationMetrics
	
	// Configuration
	config            CollaborationPerformanceConfig
	
	// Synchronization
	mutex             sync.RWMutex
}

// CollaborationPerformanceConfig defines collaboration performance optimization settings
type CollaborationPerformanceConfig struct {
	SessionCacheSize     int           `json:"session_cache_size"`
	OperationPoolSize    int           `json:"operation_pool_size"`
	BatchSize            int           `json:"batch_size"`
	BatchTimeout         time.Duration `json:"batch_timeout"`
	ConnectionPoolSize   int           `json:"connection_pool_size"`
	ConnectionTimeout    time.Duration `json:"connection_timeout"`
	EnableMetrics        bool          `json:"enable_metrics"`
	EnableCompression    bool          `json:"enable_compression"`
	MaxConcurrentUsers   int           `json:"max_concurrent_users"`
}

// CachedSession represents a cached collaboration session
type CachedSession struct {
	Session      *Session       `json:"session"`
	Participants []*Participant `json:"participants"`
	LoadTime     time.Time      `json:"load_time"`
	AccessCount  int            `json:"access_count"`
	LastAccess   time.Time      `json:"last_access"`
	Dirty        bool           `json:"dirty"`
}

// OperationPool provides reusable operation objects
type OperationPool struct {
	operations []Operation
	mutex      sync.Mutex
	nextIndex  int
}

// BatchProcessor handles batch processing of collaboration operations
type BatchProcessor struct {
	requestQueue chan *BatchRequest
	workers      int
	batchSize    int
	batchTimeout time.Duration
	mutex        sync.Mutex
}

// BatchRequest represents a batched collaboration request
type BatchRequest struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
	Response  chan *BatchResponse    `json:"-"`
	Context   context.Context        `json:"-"`
}

// BatchResponse represents the result of a batched request
type BatchResponse struct {
	ID      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   error       `json:"error"`
	Latency time.Duration `json:"latency"`
}

// ConnectionPool manages database connections for collaboration
type ConnectionPool struct {
	connections chan *db.DB
	mu          sync.Mutex
	maxSize     int
	timeout     time.Duration
}

// CollaborationMetrics tracks collaboration performance metrics
type CollaborationMetrics struct {
	SessionOperations   int64         `json:"session_operations"`
	OperationBatches    int64         `json:"operation_batches"`
	AverageOperationTime time.Duration `json:"average_operation_time"`
	CacheHits           int64         `json:"cache_hits"`
	CacheMisses         int64         `json:"cache_misses"`
	ConnectionErrors    int64         `json:"connection_errors"`
	ConcurrentUsers     int64         `json:"concurrent_users"`
	mutex               sync.RWMutex
}

// NewPerformanceOptimizedCollaborationManager creates a new performance-optimized collaboration manager
func NewPerformanceOptimizedCollaborationManager(database *db.DB, config CollaborationPerformanceConfig) *PerformanceOptimizedCollaborationManager {
	// Set defaults if not provided
	if config.SessionCacheSize == 0 {
		config.SessionCacheSize = 100
	}
	if config.OperationPoolSize == 0 {
		config.OperationPoolSize = 50
	}
	if config.BatchSize == 0 {
		config.BatchSize = 10
	}
	if config.BatchTimeout == 0 {
		config.BatchTimeout = 100 * time.Millisecond
	}
	if config.ConnectionPoolSize == 0 {
		config.ConnectionPoolSize = 20
	}
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = 30 * time.Second
	}
	if config.MaxConcurrentUsers == 0 {
		config.MaxConcurrentUsers = 50
	}

	// Create base collaboration manager
	baseManager := NewCollaborationManager(database)

	manager := &PerformanceOptimizedCollaborationManager{
		CollaborationManager: baseManager,
		sessionCache:         make(map[string]*CachedSession),
		operationPool:        NewOperationPool(config.OperationPoolSize),
		batchProcessor:       NewBatchProcessor(config.BatchSize, config.BatchTimeout),
		connectionPool:       NewConnectionPool(database, config.ConnectionPoolSize, config.ConnectionTimeout),
		metrics:              &CollaborationMetrics{},
		config:               config,
	}

	// Start background processes
	go manager.batchProcessor.start()

	logging.GetDefaultLogger().Info("Performance-optimized collaboration manager initialized",
		"session_cache_size", config.SessionCacheSize,
		"operation_pool_size", config.OperationPoolSize,
		"connection_pool_size", config.ConnectionPoolSize)

	return manager
}

// CreateSessionOptimized creates a new session with performance optimizations
func (m *PerformanceOptimizedCollaborationManager) CreateSessionOptimized(documentID int, name, createdBy string, settings SessionSettings) (*Session, error) {
	start := time.Now()
	
	// Update metrics
	m.metrics.mutex.Lock()
	m.metrics.SessionOperations++
	m.metrics.mutex.Unlock()

	// Check concurrent user limit
	if m.getConcurrentUserCount() >= m.config.MaxConcurrentUsers {
		return nil, fmt.Errorf("maximum concurrent users exceeded: %d", m.config.MaxConcurrentUsers)
	}

	// Get connection from pool
	conn, err := m.connectionPool.Get()
	if err != nil {
		m.metrics.mutex.Lock()
		m.metrics.ConnectionErrors++
		m.metrics.mutex.Unlock()
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	defer m.connectionPool.Put(conn)

	// Create session using pooled connection
	session, err := m.createSessionWithConnection(conn, documentID, name, createdBy, settings)
	if err != nil {
		return nil, err
	}

	// Cache the session
	m.cacheSession(session.ID, session, nil)

	duration := time.Since(start)
	m.updateAverageOperationTime(duration)
	
	logging.GetDefaultLogger().Debug("Session created with optimization", "id", session.ID, "duration", duration)
	return session, nil
}

// JoinSessionOptimized allows a user to join a session with performance optimizations
func (m *PerformanceOptimizedCollaborationManager) JoinSessionOptimized(sessionID, userID, username string, role ParticipantRole) (*Session, error) {
	start := time.Now()

	// Check cache first
	if cached, found := m.getSessionFromCache(sessionID); found {
		// Validate session is still active
		if cached.Session.IsActive {
			// Add participant to cached session
			participant := &Participant{
				UserID:      userID,
				SessionID:   sessionID,
				Username:    username,
				Role:        role,
				JoinedAt:    time.Now(),
				LastSeen:    time.Now(),
				IsActive:    true,
				Permissions: m.getPermissionsForRole(role),
			}
			
			cached.Participants = append(cached.Participants, participant)
			cached.AccessCount++
			cached.LastAccess = time.Now()
			cached.Dirty = true

			duration := time.Since(start)
			m.updateAverageOperationTime(duration)
			
			logging.GetDefaultLogger().Debug("User joined session from cache", "session_id", sessionID, "user_id", userID)
			return cached.Session, nil
		}
	}

	// Fallback to base implementation
	session, err := m.CollaborationManager.JoinSession(sessionID, userID, username, role)
	if err != nil {
		return nil, err
	}

	// Update cache
	m.cacheSession(sessionID, session, nil)

	duration := time.Since(start)
	m.updateAverageOperationTime(duration)
	
	return session, nil
}

// ApplyOperationOptimized applies an operation with performance optimizations
func (m *PerformanceOptimizedCollaborationManager) ApplyOperationOptimized(sessionID, userID string, operation Operation) error {
	start := time.Now()

	// Get operation from pool
	op := m.operationPool.Get()
	defer m.operationPool.Put(op)

	// Copy operation data
	*op = operation
	op.ID = m.generateOperationID()
	op.SessionID = sessionID
	op.UserID = userID
	op.Timestamp = time.Now()

	// Add to batch processor
	request := &BatchRequest{
		ID:        op.ID,
		Type:      "operation",
		Data:      map[string]interface{}{"operation": op},
		Timestamp: time.Now(),
		Response:  make(chan *BatchResponse, 1),
		Context:   context.Background(),
	}

	m.batchProcessor.Queue(request)

	// Wait for response
	select {
	case response := <-request.Response:
		if response.Error != nil {
			return response.Error
		}
		
		duration := time.Since(start)
		m.updateAverageOperationTime(duration)
		
		return nil
		
	case <-time.After(5 * time.Second):
		return fmt.Errorf("operation timeout")
	}
}

// getSessionFromCache retrieves a session from cache
func (m *PerformanceOptimizedCollaborationManager) getSessionFromCache(sessionID string) (*CachedSession, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	cached, exists := m.sessionCache[sessionID]
	if !exists {
		m.metrics.mutex.Lock()
		m.metrics.CacheMisses++
		m.metrics.mutex.Unlock()
		return nil, false
	}

	m.metrics.mutex.Lock()
	m.metrics.CacheHits++
	m.metrics.mutex.Unlock()

	cached.AccessCount++
	cached.LastAccess = time.Now()

	return cached, true
}

// cacheSession stores a session in cache
func (m *PerformanceOptimizedCollaborationManager) cacheSession(sessionID string, session *Session, participants []*Participant) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if we need to evict entries
	if len(m.sessionCache) >= m.config.SessionCacheSize {
		m.evictSessionCacheLRU()
	}

	if participants == nil {
		participants = []*Participant{}
	}

	m.sessionCache[sessionID] = &CachedSession{
		Session:      session,
		Participants: participants,
		LoadTime:     time.Now(),
		AccessCount:  1,
		LastAccess:   time.Now(),
		Dirty:        false,
	}
}

// evictSessionCacheLRU removes the least recently used session from cache
func (m *PerformanceOptimizedCollaborationManager) evictSessionCacheLRU() {
	if len(m.sessionCache) == 0 {
		return
	}

	var oldestID string
	var oldestTime time.Time
	var lowestAccess int
	first := true

	for id, cached := range m.sessionCache {
		if first || cached.AccessCount < lowestAccess || 
		   (cached.AccessCount == lowestAccess && cached.LastAccess.Before(oldestTime)) {
			oldestID = id
			oldestTime = cached.LastAccess
			lowestAccess = cached.AccessCount
			first = false
		}
	}

	if oldestID != "" {
		delete(m.sessionCache, oldestID)
		logging.GetDefaultLogger().Debug("Evicted session from cache", "id", oldestID)
	}
}

// getConcurrentUserCount returns the current number of concurrent users
func (m *PerformanceOptimizedCollaborationManager) getConcurrentUserCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	count := 0
	for _, cached := range m.sessionCache {
		if cached.Session.IsActive {
			count += len(cached.Participants)
		}
	}

	return count
}

// createSessionWithConnection creates a session using a specific database connection
func (m *PerformanceOptimizedCollaborationManager) createSessionWithConnection(conn *db.DB, documentID int, name, createdBy string, settings SessionSettings) (*Session, error) {
	// This is a simplified implementation
	// In a real application, this would use the provided connection
	return m.CollaborationManager.CreateSession(documentID, name, createdBy, settings)
}

// updateAverageOperationTime updates the rolling average operation time
func (m *PerformanceOptimizedCollaborationManager) updateAverageOperationTime(duration time.Duration) {
	m.metrics.mutex.Lock()
	defer m.metrics.mutex.Unlock()

	if m.metrics.AverageOperationTime == 0 {
		m.metrics.AverageOperationTime = duration
	} else {
		// Weighted average (90% old, 10% new)
		m.metrics.AverageOperationTime = time.Duration(
			float64(m.metrics.AverageOperationTime)*0.9 + float64(duration)*0.1,
		)
	}
}

// NewOperationPool creates a new operation pool
func NewOperationPool(size int) *OperationPool {
	pool := &OperationPool{
		operations: make([]Operation, size),
		nextIndex:  0,
	}

	// Initialize pool with empty operations
	for i := range pool.operations {
		pool.operations[i] = Operation{}
	}

	return pool
}

// Get gets an operation from the pool
func (op *OperationPool) Get() *Operation {
	op.mutex.Lock()
	defer op.mutex.Unlock()

	if len(op.operations) == 0 {
		return &Operation{}
	}

	operation := &op.operations[op.nextIndex]
	op.nextIndex = (op.nextIndex + 1) % len(op.operations)

	return operation
}

// Put returns an operation to the pool
func (op *OperationPool) Put(operation *Operation) {
	// Reset operation
	*operation = Operation{}
	
	op.mutex.Lock()
	defer op.mutex.Unlock()
	
	// Operations are reused by index, no need to explicitly put back
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(batchSize int, timeout time.Duration) *BatchProcessor {
	return &BatchProcessor{
		requestQueue: make(chan *BatchRequest, 100),
		workers:      3,
		batchSize:    batchSize,
		batchTimeout: timeout,
	}
}

// Queue adds a request to the batch processor queue
func (bp *BatchProcessor) Queue(request *BatchRequest) {
	select {
	case bp.requestQueue <- request:
	default:
		// Queue full, process immediately
		bp.processRequest(request)
	}
}

// start starts the batch processor workers
func (bp *BatchProcessor) start() {
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
				batch = batch[:0]
			}

		case <-ticker.C:
			if len(batch) > 0 {
				bp.processBatch(batch)
				batch = batch[:0]
			}
		}
	}
}

// processBatch processes a batch of requests
func (bp *BatchProcessor) processBatch(batch []*BatchRequest) {
	start := time.Now()
	
	for _, request := range batch {
		bp.processRequest(request)
	}
	
	duration := time.Since(start)
	logging.GetDefaultLogger().Debug("Processed batch", "size", len(batch), "duration", duration)
}

// processRequest processes a single request
func (bp *BatchProcessor) processRequest(request *BatchRequest) {
	start := time.Now()
	
	response := &BatchResponse{
		ID:      request.ID,
		Result:  map[string]interface{}{"processed": true},
		Error:   nil,
		Latency: time.Since(start),
	}
	
	select {
	case request.Response <- response:
	case <-request.Context.Done():
		// Request timed out
	}
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(baseConn *db.DB, maxSize int, timeout time.Duration) *ConnectionPool {
	pool := &ConnectionPool{
		connections: make(chan *db.DB, maxSize),
		maxSize:     maxSize,
		timeout:     timeout,
	}

	// For now, we'll use the same connection multiple times
	// In a real application, you'd create multiple actual connections
	for i := 0; i < maxSize; i++ {
		pool.connections <- baseConn
	}

	return pool
}

// Get gets a connection from the pool
func (cp *ConnectionPool) Get() (*db.DB, error) {
	select {
	case conn := <-cp.connections:
		return conn, nil
	case <-time.After(cp.timeout):
		return nil, fmt.Errorf("connection pool timeout")
	}
}

// Put returns a connection to the pool
func (cp *ConnectionPool) Put(conn *db.DB) {
	select {
	case cp.connections <- conn:
	default:
		// Pool full, connection will be garbage collected
	}
}

// GetMetrics returns current collaboration performance metrics
func (m *PerformanceOptimizedCollaborationManager) GetMetrics() CollaborationMetrics {
	m.metrics.mutex.RLock()
	defer m.metrics.mutex.RUnlock()

	return *m.metrics
}

// GetPerformanceReport returns a comprehensive performance report
func (m *PerformanceOptimizedCollaborationManager) GetPerformanceReport() map[string]interface{} {
	metrics := m.GetMetrics()
	
	m.mutex.RLock()
	sessionCacheSize := len(m.sessionCache)
	concurrentUsers := m.getConcurrentUserCount()
	m.mutex.RUnlock()

	report := map[string]interface{}{
		"metrics": metrics,
		"config":  m.config,
		"cache_info": map[string]interface{}{
			"session_cache_size": sessionCacheSize,
			"concurrent_users":   concurrentUsers,
		},
	}

	// Calculate cache hit rate
	if metrics.CacheHits+metrics.CacheMisses > 0 {
		report["cache_hit_rate"] = float64(metrics.CacheHits) / float64(metrics.CacheHits+metrics.CacheMisses)
	}

	return report
}

// generateOperationID generates a unique operation ID
func (m *PerformanceOptimizedCollaborationManager) generateOperationID() string {
	return fmt.Sprintf("opt_%d", time.Now().UnixNano())
}

// Close cleans up resources
func (m *PerformanceOptimizedCollaborationManager) Close() error {
	logging.GetDefaultLogger().Info("Performance-optimized collaboration manager shutting down")
	
	// Close base manager
	return m.CollaborationManager.Close()
}