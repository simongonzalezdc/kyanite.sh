package app

import "time"

// Editor constants
const (
	// Retry configuration
	MaxSaveRetries       = 3
	MaxAutoSaveRetries   = 2
	SaveRetryDelayBase   = 1 * time.Second
	AutoSaveRetryDelay   = 500 * time.Millisecond
	SaveOperationTimeout = 30 * time.Second

	// Auto-save configuration
	MinAutoSaveInterval     = 5  // seconds
	MaxAutoSaveInterval     = 300 // seconds
	MinAutoSaveDebounce     = 100 // milliseconds
	MaxAutoSaveDebounce     = 10000 // milliseconds
	DefaultAutoSaveInterval = 30 * time.Second
	DefaultDebounceDelay    = 2000 * time.Millisecond
	DefaultMaxVersions      = 10

	// Status display timing
	StatusMessageDuration = 2 * time.Second
	IdleStatusDelay       = 2 * time.Second

	// Dictionary constants
	MinSyllableWordLength = 2
	DefaultSyllableCount  = 1
	PhoneticEndingLength  = 3

	// Content size limits
	MaxDescriptionLength = 1000 // characters
)
