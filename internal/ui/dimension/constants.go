package dimension

import "time"

// Scroll configuration
const (
	ScrollStepSize       = 3
	PageScrollMultiplier = 10
)

// Editor pane constants
const (
	LineNumberPadding = 4
	CursorBlinkRate   = 500 * time.Millisecond
)

// Status bar constants
const (
	MinStatusBarWidth  = 40
	MaxShortcutDisplay = 5
)
