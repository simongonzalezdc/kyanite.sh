// Package icons provides a 3-tier icon system (ASCII, Unicode, NerdFont)
// with string width utilities and per-app icon registration.
package icons

import "sync"

// Style represents the icon rendering tier.
type Style int

const (
	// ASCII uses plain ASCII characters for maximum terminal compatibility.
	ASCII Style = iota
	// Unicode uses Unicode symbols for terminals with Unicode support.
	Unicode
	// NerdFont uses Nerd Font icons for terminals with Nerd Font installed.
	NerdFont
)

// Icon represents a named icon with variants for each tier.
type Icon struct {
	Name     string
	ASCII    string
	Unicode  string
	NerdFont string
}

// Get returns the icon string for the given style tier.
func (i Icon) Get(s Style) string {
	switch s {
	case NerdFont:
		if i.NerdFont != "" {
			return i.NerdFont
		}
		fallthrough
	case Unicode:
		if i.Unicode != "" {
			return i.Unicode
		}
		fallthrough
	default:
		return i.ASCII
	}
}

// DefaultStyle is the default icon style used when none is specified.
var DefaultStyle = Unicode

// baseIcons holds the shared base icon set.
var baseIcons = map[string]Icon{}
var baseIconsMu sync.RWMutex

// GetIcon returns the icon string for the given name using DefaultStyle,
// or empty string if not found.
func GetIcon(name string) string {
	baseIconsMu.RLock()
	defer baseIconsMu.RUnlock()
	icon, ok := baseIcons[name]
	if !ok {
		return ""
	}
	return icon.Get(DefaultStyle)
}

// RegisterIcons registers app-specific icons into the shared registry.
func RegisterIcons(appIcons map[string]Icon) {
	baseIconsMu.Lock()
	defer baseIconsMu.Unlock()
	for name, icon := range appIcons {
		baseIcons[name] = icon
	}
}

// currentStyle holds the global icon style.
var currentStyle = Unicode
var styleMu sync.RWMutex

// SetStyle changes the active icon style.
func SetStyle(s Style) {
	styleMu.Lock()
	defer styleMu.Unlock()
	currentStyle = s
}

// GetStyle returns the current icon style.
func GetStyle() Style {
	styleMu.RLock()
	defer styleMu.RUnlock()
	return currentStyle
}

// CurrentIcon returns the icon string for the given name using the current style.
func CurrentIcon(name string) string {
	styleMu.RLock()
	s := currentStyle
	styleMu.RUnlock()
	baseIconsMu.RLock()
	icon, ok := baseIcons[name]
	baseIconsMu.RUnlock()
	if !ok {
		return ""
	}
	return icon.Get(s)
}

func init() {
	RegisterIcons(baseIconDefs)
}

// baseIconDefs contains all shared icon definitions.
var baseIconDefs = map[string]Icon{
	// Status indicators
	"success":    {Name: "success", ASCII: "[OK]", Unicode: "✓", NerdFont: "\uf00c"},
	"error":      {Name: "error", ASCII: "[X]", Unicode: "✗", NerdFont: "\uf00d"},
	"warning":    {Name: "warning", ASCII: "[!]", Unicode: "⚠", NerdFont: "\uf071"},
	"info":       {Name: "info", ASCII: "[i]", Unicode: "ℹ", NerdFont: "\uf05a"},
	"recording":  {Name: "recording", ASCII: "[REC]", Unicode: "●", NerdFont: "\uf111"},
	"processing": {Name: "processing", ASCII: "[...]", Unicode: "⟳", NerdFont: "\uf110"},

	// Checkboxes/Selection
	"check_on":  {Name: "check_on", ASCII: "[x]", Unicode: "☑", NerdFont: "\uf046"},
	"check_off": {Name: "check_off", ASCII: "[ ]", Unicode: "☐", NerdFont: "\uf096"},
	"radio_on":  {Name: "radio_on", ASCII: "[*]", Unicode: "●", NerdFont: "\uf192"},
	"radio_off": {Name: "radio_off", ASCII: "[ ]", Unicode: "○", NerdFont: "\uf10c"},

	// Navigation
	"arrow_left":  {Name: "arrow_left", ASCII: "<-", Unicode: "←", NerdFont: "\uf060"},
	"arrow_right": {Name: "arrow_right", ASCII: "->", Unicode: "→", NerdFont: "\uf061"},
	"arrow_up":    {Name: "arrow_up", ASCII: "Up", Unicode: "↑", NerdFont: "\uf062"},
	"arrow_down":  {Name: "arrow_down", ASCII: "Down", Unicode: "↓", NerdFont: "\uf063"},
	"chevron_l":   {Name: "chevron_l", ASCII: "<<", Unicode: "◀", NerdFont: "\uf053"},
	"chevron_r":   {Name: "chevron_r", ASCII: ">>", Unicode: "▶", NerdFont: "\uf054"},

	// Bullets/Lists
	"bullet":    {Name: "bullet", ASCII: "-", Unicode: "•", NerdFont: "\uf054"},
	"bullet_alt": {Name: "bullet_alt", ASCII: "*", Unicode: "◦", NerdFont: "\uf105"},
	"separator": {Name: "separator", ASCII: "-", Unicode: "•", NerdFont: "\uf101"},

	// Progress
	"progress_full":  {Name: "progress_full", ASCII: "#", Unicode: "█", NerdFont: "\u2588"},
	"progress_empty": {Name: "progress_empty", ASCII: "-", Unicode: "░", NerdFont: "\u2591"},

	// Presence/Status
	"online":  {Name: "online", ASCII: "[*]", Unicode: "●", NerdFont: "\uf111"},
	"away":    {Name: "away", ASCII: "[~]", Unicode: "●", NerdFont: "\uf111"},
	"busy":    {Name: "busy", ASCII: "[!]", Unicode: "●", NerdFont: "\uf111"},
	"offline": {Name: "offline", ASCII: "[ ]", Unicode: "○", NerdFont: "\uf10c"},

	// Application
	"music":    {Name: "music", ASCII: "[~]", Unicode: "♪", NerdFont: "\uf001"},
	"folder":   {Name: "folder", ASCII: "[D]", Unicode: "📁", NerdFont: "\uf07b"},
	"file":     {Name: "file", ASCII: "[F]", Unicode: "📄", NerdFont: "\uf15b"},
	"settings": {Name: "settings", ASCII: "[S]", Unicode: "⚙", NerdFont: "\uf013"},
	"help":     {Name: "help", ASCII: "[?]", Unicode: "?", NerdFont: "\uf059"},
	"search":   {Name: "search", ASCII: "[/]", Unicode: "🔍", NerdFont: "\uf002"},
	"ai":       {Name: "ai", ASCII: "[AI]", Unicode: "🤖", NerdFont: "\uf544"},
	"mic":      {Name: "mic", ASCII: "[MIC]", Unicode: "🎤", NerdFont: "\uf130"},
	"photo":    {Name: "photo", ASCII: "[IMG]", Unicode: "📷", NerdFont: "\uf03e"},
	"export":   {Name: "export", ASCII: "[E]", Unicode: "📤", NerdFont: "\uf093"},
	"sync":     {Name: "sync", ASCII: "[SYNC]", Unicode: "🔄", NerdFont: "\uf021"},

	// Dashboard panels
	"stats":       {Name: "stats", ASCII: "[STATS]", Unicode: "📊", NerdFont: "\uf080"},
	"performance": {Name: "performance", ASCII: "[PERF]", Unicode: "⚡", NerdFont: "\uf0e7"},
	"storage":     {Name: "storage", ASCII: "[DISK]", Unicode: "💾", NerdFont: "\uf0a0"},
	"tools":       {Name: "tools", ASCII: "[TOOLS]", Unicode: "🛠", NerdFont: "\uf0ad"},
	"tip":         {Name: "tip", ASCII: "[TIP]", Unicode: "💡", NerdFont: "\uf0eb"},
}
