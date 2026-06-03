package collaboration

import "time"

// Session constants
const (
	// Session limits
	MaxSessionNameLength  = 100
	MaxSessionUsers       = 50
	MaxSessionDescription = 500

	// Timeouts
	SessionInactivityTimeout = 30 * time.Minute
	SessionCleanupInterval   = 5 * time.Minute

	// Heartbeat
	HeartbeatInterval    = 10 * time.Second
	MissedHeartbeatLimit = 3
)

// Event and message constants
const (
	// Channel buffer sizes
	EventChannelBuffer     = 100
	BroadcastChannelBuffer = 100
	MessageChannelBuffer   = 1000

	// Processing intervals
	EventProcessInterval     = 100 * time.Millisecond
	BroadcastProcessInterval = 100 * time.Millisecond

	// Batch processing
	BatchProcessInterval = 100 * time.Millisecond
	BatchTimeout         = 5 * time.Second
)

// Cache constants
const (
	// Session cache
	MaxSessionCacheSize    = 100
	SessionCacheExpiration = 10 * time.Minute
	CacheCleanupInterval   = 1 * time.Minute

	// Eviction configuration
	CacheEvictionThreshold = 90 // percent
	MinCacheEvictionCount  = 10
)

// Performance optimization constants
const (
	// Average tracking
	AverageOperationWindow = 100

	// Metrics
	MetricsUpdateInterval = 1 * time.Second
)
