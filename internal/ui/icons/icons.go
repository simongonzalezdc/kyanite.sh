// Package icons provides configurable icon sets for the TUI.
// It supports ASCII (default, universal), Unicode, and Nerd Font icon sets.
// Users can switch between sets based on their terminal capabilities.
package icons

import "sync"

// IconSet contains all icons used throughout the application
type IconSet struct {
	// Status indicators
	Success    string
	Error      string
	Warning    string
	Info       string
	Recording  string
	Processing string

	// Checkboxes/Selection
	CheckOn  string
	CheckOff string
	RadioOn  string
	RadioOff string

	// Navigation
	ArrowLeft  string
	ArrowRight string
	ArrowUp    string
	ArrowDown  string
	ChevronL   string
	ChevronR   string

	// Bullets/Lists
	Bullet     string
	BulletAlt  string
	Separator  string

	// Progress
	ProgressFull  string
	ProgressEmpty string

	// Presence/Status
	Online  string
	Away    string
	Busy    string
	Offline string

	// Application
	Music    string
	Folder   string
	File     string
	Settings string
	Help     string
	Search   string
	AI       string
	Mic      string
	Photo    string
	Export   string
	Sync     string

	// Dashboard panels
	Stats       string
	Performance string
	Storage     string
	Tools       string
	Tip         string
}

// ASCII is the default icon set - works on all terminals
var ASCII = IconSet{
	// Status
	Success:    "[OK]",
	Error:      "[X]",
	Warning:    "[!]",
	Info:       "[i]",
	Recording:  "[REC]",
	Processing: "[...]",

	// Checkboxes
	CheckOn:  "[x]",
	CheckOff: "[ ]",
	RadioOn:  "[*]",
	RadioOff: "[ ]",

	// Navigation
	ArrowLeft:  "<-",
	ArrowRight: "->",
	ArrowUp:    "Up",
	ArrowDown:  "Down",
	ChevronL:   "<<",
	ChevronR:   ">>",

	// Bullets
	Bullet:    "-",
	BulletAlt: "*",
	Separator: "-",

	// Progress
	ProgressFull:  "#",
	ProgressEmpty: "-",

	// Presence
	Online:  "[*]",
	Away:    "[~]",
	Busy:    "[!]",
	Offline: "[ ]",

	// Application
	Music:    "[~]",
	Folder:   "[D]",
	File:     "[F]",
	Settings: "[S]",
	Help:     "[?]",
	Search:   "[/]",
	AI:       "[AI]",
	Mic:      "[MIC]",
	Photo:    "[IMG]",
	Export:   "[E]",
	Sync:     "[SYNC]",

	// Dashboard
	Stats:       "[STATS]",
	Performance: "[PERF]",
	Storage:     "[DISK]",
	Tools:       "[TOOLS]",
	Tip:         "[TIP]",
}

// Unicode uses standard Unicode symbols - works on most modern terminals
var Unicode = IconSet{
	// Status
	Success:    "✓",
	Error:      "✗",
	Warning:    "⚠",
	Info:       "ℹ",
	Recording:  "●",
	Processing: "⟳",

	// Checkboxes
	CheckOn:  "☑",
	CheckOff: "☐",
	RadioOn:  "●",
	RadioOff: "○",

	// Navigation
	ArrowLeft:  "←",
	ArrowRight: "→",
	ArrowUp:    "↑",
	ArrowDown:  "↓",
	ChevronL:   "◀",
	ChevronR:   "▶",

	// Bullets
	Bullet:    "•",
	BulletAlt: "◦",
	Separator: "•",

	// Progress
	ProgressFull:  "█",
	ProgressEmpty: "░",

	// Presence
	Online:  "●",
	Away:    "●",
	Busy:    "●",
	Offline: "○",

	// Application
	Music:    "♪",
	Folder:   "📁",
	File:     "📄",
	Settings: "⚙",
	Help:     "?",
	Search:   "🔍",
	AI:       "🤖",
	Mic:      "🎤",
	Photo:    "📷",
	Export:   "📤",
	Sync:     "🔄",

	// Dashboard
	Stats:       "📊",
	Performance: "⚡",
	Storage:     "💾",
	Tools:       "🛠",
	Tip:         "💡",
}

// NerdFont uses Nerd Font icons - requires a patched font
// Icon codes from https://www.nerdfonts.com/cheat-sheet
var NerdFont = IconSet{
	// Status
	Success:    "\uf00c", // nf-fa-check
	Error:      "\uf00d", // nf-fa-times
	Warning:    "\uf071", // nf-fa-exclamation_triangle
	Info:       "\uf05a", // nf-fa-info_circle
	Recording:  "\uf111", // nf-fa-circle
	Processing: "\uf110", // nf-fa-spinner

	// Checkboxes
	CheckOn:  "\uf046", // nf-fa-check_square_o
	CheckOff: "\uf096", // nf-fa-square_o
	RadioOn:  "\uf192", // nf-fa-dot_circle_o
	RadioOff: "\uf10c", // nf-fa-circle_o

	// Navigation
	ArrowLeft:  "\uf060", // nf-fa-arrow_left
	ArrowRight: "\uf061", // nf-fa-arrow_right
	ArrowUp:    "\uf062", // nf-fa-arrow_up
	ArrowDown:  "\uf063", // nf-fa-arrow_down
	ChevronL:   "\uf053", // nf-fa-chevron_left
	ChevronR:   "\uf054", // nf-fa-chevron_right

	// Bullets
	Bullet:    "\uf054", // nf-fa-chevron_right
	BulletAlt: "\uf105", // nf-fa-angle_right
	Separator: "\uf101", // nf-fa-angle_double_right

	// Progress
	ProgressFull:  "\u2588", // Full block (works in most fonts)
	ProgressEmpty: "\u2591", // Light shade

	// Presence
	Online:  "\uf111", // nf-fa-circle (green)
	Away:    "\uf111", // nf-fa-circle (yellow)
	Busy:    "\uf111", // nf-fa-circle (red)
	Offline: "\uf10c", // nf-fa-circle_o

	// Application
	Music:    "\uf001", // nf-fa-music
	Folder:   "\uf07b", // nf-fa-folder
	File:     "\uf15b", // nf-fa-file
	Settings: "\uf013", // nf-fa-cog
	Help:     "\uf059", // nf-fa-question_circle
	Search:   "\uf002", // nf-fa-search
	AI:       "\uf544", // nf-mdi-robot
	Mic:      "\uf130", // nf-fa-microphone
	Photo:    "\uf03e", // nf-fa-image
	Export:   "\uf093", // nf-fa-upload
	Sync:     "\uf021", // nf-fa-refresh

	// Dashboard
	Stats:       "\uf080", // nf-fa-bar_chart
	Performance: "\uf0e7", // nf-fa-bolt
	Storage:     "\uf0a0", // nf-fa-hdd_o
	Tools:       "\uf0ad", // nf-fa-wrench
	Tip:         "\uf0eb", // nf-fa-lightbulb_o
}

// IconStyle represents the available icon styles
type IconStyle string

const (
	StyleASCII    IconStyle = "ascii"
	StyleUnicode  IconStyle = "unicode"
	StyleNerdFont IconStyle = "nerdfonts"
)

// manager holds the global icon configuration
type manager struct {
	mu      sync.RWMutex
	current *IconSet
	style   IconStyle
}

var globalManager = &manager{
	current: &ASCII,
	style:   StyleASCII,
}

// SetStyle changes the active icon style
func SetStyle(style IconStyle) {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	globalManager.style = style
	switch style {
	case StyleUnicode:
		globalManager.current = &Unicode
	case StyleNerdFont:
		globalManager.current = &NerdFont
	default:
		globalManager.current = &ASCII
	}
}

// GetStyle returns the current icon style
func GetStyle() IconStyle {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	return globalManager.style
}

// Get returns the current icon set
func Get() *IconSet {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	return globalManager.current
}

// Current is an alias for Get() for convenience
func Current() *IconSet {
	return Get()
}
