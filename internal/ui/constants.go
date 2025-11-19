package ui

import "time"

// UI timing constants
const (
	// Refresh rates
	UIRefreshInterval     = 100 * time.Millisecond
	AnimationTickInterval = 50 * time.Millisecond

	// Message display durations
	SaveMessageDuration  = 3 * time.Second
	ErrorMessageDuration = 5 * time.Second

	// Theme switching
	ThemeSwitchAnimationDuration = 200 * time.Millisecond
)

// Preview pane constants
const (
	// Content thresholds
	DefaultContentThreshold = 50000 // 50KB
	SmallContentThreshold   = 30000 // 30KB
	MaxContentLength        = 100000 // 100KB

	// Cache configuration
	MaxCacheSize       = 100
	CacheEvictionCount = 20

	// Scroll configuration
	ScrollStepSize      = 3
	PageScrollMultiplier = 10
)

// Editor pane constants
const (
	// Line number padding
	LineNumberPadding = 4

	// Cursor blink rate
	CursorBlinkRate = 500 * time.Millisecond
)

// Status bar constants
const (
	// Update intervals
	StatusBarUpdateInterval = 100 * time.Millisecond

	// Display widths
	MinStatusBarWidth = 40
	MaxShortcutDisplay = 5
)
