package constants

import "time"

// Global durations and timeouts for the application
const (
	// Default timeouts
	DefaultContextTimeout = 30 * time.Second
	ShortContextTimeout   = 5 * time.Second
	LongContextTimeout    = 60 * time.Second

	// Editor specific timeouts
	SaveOperationTimeout = 30 * time.Second
	AutoSaveTimeout      = 15 * time.Second

	// AI timeouts
	BrainstormTimeout  = 30 * time.Second
	AIContextTimeout   = 3 * time.Second
	AIStreamingTimeout = 60 * time.Second

	// Database timeouts
	DBQueryTimeout      = 10 * time.Second
	DBValidationTimeout = 5 * time.Second

	// Retry delays
	DefaultRetryDelay  = 100 * time.Millisecond
	LockRetryDelay     = 50 * time.Millisecond
	SaveRetryDelayBase = 100 * time.Millisecond
	AutoSaveRetryDelay = 500 * time.Millisecond

	// Retry counts
	MaxSaveRetries     = 3
	MaxAutoSaveRetries = 3
	MaxLockRetries     = 3
	MaxDeleteRetries   = 3
	MaxDBRetries       = 3

	// UI Delays
	StatusIdleDelay = 2 * time.Second

	// Stub/Test delays
	StubDelay = 10 * time.Millisecond
)
