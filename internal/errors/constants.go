package errors

import "time"

// Error history and statistics
const (
	// History limits
	DefaultHistoryLimit = 1000
	MaxHistoryLimit     = 10000

	// Time windows
	RecentErrorWindow = 1 * time.Hour
	StatsTimeWindow   = 24 * time.Hour
)

// Recovery constants
const (
	// Retry configuration
	DefaultMaxRetries     = 3
	DefaultRetryDelay     = 1 * time.Second
	DefaultMaxRecoveryTime = 30 * time.Second

	// Backoff configuration
	ExponentialBackoffBase = 2
	MaxBackoffDelay        = 30 * time.Second
)

// Notification constants
const (
	NotificationDisplayDuration = 5 * time.Second
	ErrorNotificationTimeout    = 10 * time.Second
)
