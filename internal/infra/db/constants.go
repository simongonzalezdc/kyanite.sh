package db

import "time"

// Connection pool constants
const (
	// Default pool configuration
	DefaultMaxOpenConns    = 10
	DefaultMaxIdleConns    = 5
	DefaultConnMaxLifetime = 1 * time.Hour
	DefaultConnMaxIdleTime = 5 * time.Minute

	// Optimized pool configuration
	OptimizedMaxOpenConns    = 25
	OptimizedMaxIdleConns    = 10
	OptimizedConnMaxLifetime = 30 * time.Minute
	OptimizedConnMaxIdleTime = 5 * time.Minute
)

// Query retry constants
const (
	MaxQueryRetries    = 3
	RetryBackoffBase   = 100 * time.Millisecond
	QueryTimeout       = 30 * time.Second
	TransactionTimeout = 30 * time.Second
)

// Performance thresholds
const (
	SlowQueryThreshold  = 1 * time.Second
	HealthCheckInterval = 30 * time.Second
)

// Batch operation constants
const (
	DefaultBatchSize = 100
	MaxBatchSize     = 1000
)
