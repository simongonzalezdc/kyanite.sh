# noise.sh Theme System Enhancement
## Implementation Specification for Kilocode AI

**Version:** 1.0  
**Date:** October 17, 2025  
**Purpose:** Complete specification for implementing 12 cohesive themes in noise.sh TUI

---

## Table of Contents

1. [Overview](#overview)
2. [Theme Philosophy](#theme-philosophy)
3. [Complete Theme Definitions](#complete-theme-definitions)
4. [Implementation Requirements](#implementation-requirements)
5. [Code Structure](#code-structure)
6. [Lipgloss Integration](#lipgloss-integration)
7. [Testing Checklist](#testing-checklist)
8. [Priority Roadmap](#priority-roadmap)

---

## Overview

### Goal
Implement a dynamic theme system with 12 professionally designed, WCAG AA compliant color schemes that users can switch between during runtime.

### Key Requirements
- ✅ All themes follow Element/Nature naming convention
- ✅ Zero music genre references in names
- ✅ WCAG AA compliant (minimum 3:1 contrast for UI elements, 4.5:1 for text)
- ✅ Each theme has 6 distinct colors: Primary, Secondary, Accent, Background, Text, Success
- ✅ Themes map to musical genres but names are universal
- ✅ Runtime theme switching without restart

---

## Theme Philosophy

### Naming Pattern
All themes use **Element/Nature** naming:
- **Materials/Elements:** Slate Mist, Amber, Clay, Iron, Jade, Indigo, Sage, Plasma
- **Natural Phenomena:** Night, Storm, Tide, Ember, Whisper, Bloom, Depths, Meadow
- **Energy Modifiers:** Molten, Electric, Sunset, Forest

### Color Distinctions
- **Slate Mist** - Grayscale (neutral)
- **Amber Night** - Tan/Purple (warm sophistication)
- **Molten Gold** - Gold/Orange/Purple (bold energy)
- **Clay Roads** - Burnt Orange/Brown (rustic)
- **Iron Storm** - Crimson/Gray/Fire Orange (dark power)
- **Jade Tide** - Cool Teal/Aqua (elegant water)
- **Sunset Ember** - Pink/Orange/Yellow (warm creative)
- **Forest Whisper** - Warm Sage Green (calm natural)
- **Electric Bloom** - Hot Pink/Cyan/Yellow (high energy)
- **Plasma Pulse** - Neon Green/Cyan/Magenta (futuristic)
- **Indigo Depths** - Royal Blue/Cadet Blue (deep soulful)
- **Sage Meadow** - Sage Green/Tan/Sandy (peaceful earth)

---

## Complete Theme Definitions

### 1. Slate Mist (Universal - Minimal)

```go
type Theme struct {
    Name       string
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Accent     lipgloss.Color
    Background lipgloss.Color
    Text       lipgloss.Color
    Success    lipgloss.Color
}

var SlateMist = Theme{
    Name:       "Slate Mist",
    Primary:    lipgloss.Color("#E0E0E0"),
    Secondary:  lipgloss.Color("#B0B0B0"),
    Accent:     lipgloss.Color("#FFFFFF"),
    Background: lipgloss.Color("#1A1A1A"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#90C695"),
}
```

**Use Cases:** Focus mode, no distraction, terminal purist  
**Mood:** Clean, minimal, professional

---

### 2. Amber Night (Jazz)

```go
var AmberNight = Theme{
    Name:       "Amber Night",
    Primary:    lipgloss.Color("#B8936E"),
    Secondary:  lipgloss.Color("#6D5B8B"),
    Accent:     lipgloss.Color("#E8C547"),
    Background: lipgloss.Color("#12101A"),
    Text:       lipgloss.Color("#F0E6D8"),
    Success:    lipgloss.Color("#7FB89C"),
}
```

**Use Cases:** Jazz, late night sessions, complex lyrics  
**Mood:** Sophisticated, warm, mysterious  
**Note:** Fusion of tan warmth + purple sophistication

---

### 3. Molten Gold (Hip-Hop)

```go
var MoltenGold = Theme{
    Name:       "Molten Gold",
    Primary:    lipgloss.Color("#FFD700"),
    Secondary:  lipgloss.Color("#FF6600"),
    Accent:     lipgloss.Color("#8B00FF"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#FFFFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}
```

**Use Cases:** Rap bars, confident flow, wordplay  
**Mood:** Bold, confident, royal  
**Note:** Gold primary with orange secondary and purple accent for variety

---

### 4. Clay Roads (Country)

```go
var ClayRoads = Theme{
    Name:       "Clay Roads",
    Primary:    lipgloss.Color("#CC5500"),
    Secondary:  lipgloss.Color("#D2691E"),
    Accent:     lipgloss.Color("#FFD700"),
    Background: lipgloss.Color("#1C1410"),
    Text:       lipgloss.Color("#FFF8DC"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

**Use Cases:** Country, storytelling, honest lyrics  
**Mood:** Warm, rustic, heartland  
**Note:** Burnt orange with chocolate orange secondary

---

### 5. Iron Storm (Metal/Rock)

```go
var IronStorm = Theme{
    Name:       "Iron Storm",
    Primary:    lipgloss.Color("#DC143C"),
    Secondary:  lipgloss.Color("#708090"),
    Accent:     lipgloss.Color("#FF4500"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#4A6741"),
}
```

**Use Cases:** Metal, hard rock, dark themes  
**Mood:** Dark, powerful, intense
**Note:** Crimson + slate mist gray with FIRE ORANGE accent for explosive pop

---

### 6. Jade Tide (Classical)

```go
var JadeTide = Theme{
    Name:       "Jade Tide",
    Primary:    lipgloss.Color("#20B2AA"),
    Secondary:  lipgloss.Color("#2E8B87"),
    Accent:     lipgloss.Color("#E0F2F1"),
    Background: lipgloss.Color("#0F1419"),
    Text:       lipgloss.Color("#F5F8FA"),
    Success:    lipgloss.Color("#7FB89C"),
}
```

**Use Cases:** Classical, poetic, professional  
**Mood:** Elegant, refined, sophisticated  
**Note:** Cool aquatic teal - distinct from Forest Whisper (warm green) and Indigo Depths (blue)

---

### 7. Sunset Ember (Universal - Creative)

```go
var SunsetEmber = Theme{
    Name:       "Sunset Ember",
    Primary:    lipgloss.Color("#FF6B9D"),
    Secondary:  lipgloss.Color("#FF8C42"),
    Accent:     lipgloss.Color("#FFC312"),
    Background: lipgloss.Color("#1e1e2e"),
    Text:       lipgloss.Color("#FFEAA7"),
    Success:    lipgloss.Color("#55E6C1"),
}
```

**Use Cases:** Cross-genre, creative sessions, evening writing  
**Mood:** Warm, inspiring, sunset vibes  
**Note:** Pink to orange gradient feel

---

### 8. Forest Whisper (Universal - Calm)

```go
var ForestWhisper = Theme{
    Name:       "Forest Whisper",
    Primary:    lipgloss.Color("#52B788"),
    Secondary:  lipgloss.Color("#52A068"),
    Accent:     lipgloss.Color("#95D5B2"),
    Background: lipgloss.Color("#1B263B"),
    Text:       lipgloss.Color("#D8F3DC"),
    Success:    lipgloss.Color("#B7E4C7"),
}
```

**Use Cases:** Focus mode, long sessions, nature themes  
**Mood:** Calm, focused, natural  
**Note:** Warm sage greens - distinct from Jade Tide's cool teal

---

### 9. Electric Bloom (Pop/EDM)

```go
var ElectricBloom = Theme{
    Name:       "Electric Bloom",
    Primary:    lipgloss.Color("#FF0080"),
    Secondary:  lipgloss.Color("#00D4FF"),
    Accent:     lipgloss.Color("#FFE600"),
    Background: lipgloss.Color("#0D0221"),
    Text:       lipgloss.Color("#F0F3FF"),
    Success:    lipgloss.Color("#39FF14"),
}
```

**Use Cases:** Pop songs, EDM, high energy writing  
**Mood:** Bold, powerful, electric  
**Note:** Hot pink + cyan + yellow for maximum vibrancy

---

### 10. Plasma Pulse (Electronic)

```go
var PlasmaPulse = Theme{
    Name:       "Plasma Pulse",
    Primary:    lipgloss.Color("#39FF14"),
    Secondary:  lipgloss.Color("#00F5FF"),
    Accent:     lipgloss.Color("#FF1493"),
    Background: lipgloss.Color("#0A0118"),
    Text:       lipgloss.Color("#E0FFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}
```

**Use Cases:** EDM, techno, electronic music  
**Mood:** Futuristic, energetic, digital  
**Note:** Neon green + cyan with hot pink accent

---

### 11. Indigo Depths (Blues/Soul)

```go
var IndigoDepths = Theme{
    Name:       "Indigo Depths",
    Primary:    lipgloss.Color("#4169E1"),
    Secondary:  lipgloss.Color("#5F9EA0"),
    Accent:     lipgloss.Color("#DEB887"),
    Background: lipgloss.Color("#0C0C1E"),
    Text:       lipgloss.Color("#F0E68C"),
    Success:    lipgloss.Color("#5F9EA0"),
}
```

**Use Cases:** Blues, R&B, soul music  
**Mood:** Soulful, melancholic, deep  
**Note:** Royal blue + cadet blue with warm tan accent

---

### 12. Sage Meadow (Folk/Acoustic)

```go
var SageMeadow = Theme{
    Name:       "Sage Meadow",
    Primary:    lipgloss.Color("#8FBC8F"),
    Secondary:  lipgloss.Color("#D2B48C"),
    Accent:     lipgloss.Color("#F4A460"),
    Background: lipgloss.Color("#1A1612"),
    Text:       lipgloss.Color("#FAF0E6"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

**Use Cases:** Folk, acoustic, nature-themed lyrics  
**Mood:** Natural, peaceful, grounded  
**Note:** Sage green + tan/sandy colors

---

## Implementation Requirements

### File Structure

```
noise.sh/
├── internal/
│   ├── theme/
│   │   ├── theme.go          # Theme struct and definitions
│   │   ├── manager.go        # Theme switching logic
│   │   └── registry.go       # All 12 themes registered
│   └── ui/
│       ├── styles.go         # Lipgloss styles using current theme
│       └── components.go     # UI components
├── config/
│   └── config.go             # User's theme preference saved here
```

### Core Implementation Files

#### 1. `internal/theme/theme.go`

```go
package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a color scheme for the TUI
type Theme struct {
    Name       string
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Accent     lipgloss.Color
    Background lipgloss.Color
    Text       lipgloss.Color
    Success    lipgloss.Color
}

// GetStyle returns a lipgloss style with the given colors
func (t *Theme) GetStyle(fg, bg lipgloss.Color) lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(fg).
        Background(bg)
}

// PrimaryStyle returns a style using primary color
func (t *Theme) PrimaryStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Primary)
}

// SecondaryStyle returns a style using secondary color
func (t *Theme) SecondaryStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Secondary)
}

// AccentStyle returns a style using accent color
func (t *Theme) AccentStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Background).
        Background(t.Accent)
}

// SuccessStyle returns a style for success messages
func (t *Theme) SuccessStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Success).
        Background(t.Background)
}
```

---

#### 2. `internal/theme/registry.go`

```go
package theme

// Registry holds all available themes
var Registry = map[string]Theme{
    "slate":         SlateMist,
    "amber-night":   AmberNight,
    "molten-gold":   MoltenGold,
    "clay-roads":    ClayRoads,
    "iron-storm":    IronStorm,
    "jade-tide":     JadeTide,
    "sunset-ember":  SunsetEmber,
    "forest-whisper": ForestWhisper,
    "electric-bloom": ElectricBloom,
    "plasma-pulse":  PlasmaPulse,
    "indigo-depths": IndigoDepths,
    "sage-meadow":   SageMeadow,
}

// Default theme
var DefaultTheme = "amber-night"

// GetTheme returns a theme by name, or default if not found
func GetTheme(name string) Theme {
    if theme, ok := Registry[name]; ok {
        return theme
    }
    return Registry[DefaultTheme]
}

// ListThemes returns all available theme names
func ListThemes() []string {
    themes := make([]string, 0, len(Registry))
    for name := range Registry {
        themes = append(themes, name)
    }
    return themes
}

// Define all 12 themes
var SlateMist = Theme{
    Name:       "Slate Mist",
    Primary:    lipgloss.Color("#E0E0E0"),
    Secondary:  lipgloss.Color("#B0B0B0"),
    Accent:     lipgloss.Color("#FFFFFF"),
    Background: lipgloss.Color("#1A1A1A"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#90C695"),
}

var AmberNight = Theme{
    Name:       "Amber Night",
    Primary:    lipgloss.Color("#B8936E"),
    Secondary:  lipgloss.Color("#6D5B8B"),
    Accent:     lipgloss.Color("#E8C547"),
    Background: lipgloss.Color("#12101A"),
    Text:       lipgloss.Color("#F0E6D8"),
    Success:    lipgloss.Color("#7FB89C"),
}

var MoltenGold = Theme{
    Name:       "Molten Gold",
    Primary:    lipgloss.Color("#FFD700"),
    Secondary:  lipgloss.Color("#FF6600"),
    Accent:     lipgloss.Color("#8B00FF"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#FFFFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}

var ClayRoads = Theme{
    Name:       "Clay Roads",
    Primary:    lipgloss.Color("#CC5500"),
    Secondary:  lipgloss.Color("#D2691E"),
    Accent:     lipgloss.Color("#FFD700"),
    Background: lipgloss.Color("#1C1410"),
    Text:       lipgloss.Color("#FFF8DC"),
    Success:    lipgloss.Color("#9ACD32"),
}

var IronStorm = Theme{
    Name:       "Iron Storm",
    Primary:    lipgloss.Color("#DC143C"),
    Secondary:  lipgloss.Color("#708090"),
    Accent:     lipgloss.Color("#FF4500"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#4A6741"),
}

var JadeTide = Theme{
    Name:       "Jade Tide",
    Primary:    lipgloss.Color("#20B2AA"),
    Secondary:  lipgloss.Color("#2E8B87"),
    Accent:     lipgloss.Color("#E0F2F1"),
    Background: lipgloss.Color("#0F1419"),
    Text:       lipgloss.Color("#F5F8FA"),
    Success:    lipgloss.Color("#7FB89C"),
}

var SunsetEmber = Theme{
    Name:       "Sunset Ember",
    Primary:    lipgloss.Color("#FF6B9D"),
    Secondary:  lipgloss.Color("#FF8C42"),
    Accent:     lipgloss.Color("#FFC312"),
    Background: lipgloss.Color("#1e1e2e"),
    Text:       lipgloss.Color("#FFEAA7"),
    Success:    lipgloss.Color("#55E6C1"),
}

var ForestWhisper = Theme{
    Name:       "Forest Whisper",
    Primary:    lipgloss.Color("#52B788"),
    Secondary:  lipgloss.Color("#52A068"),
    Accent:     lipgloss.Color("#95D5B2"),
    Background: lipgloss.Color("#1B263B"),
    Text:       lipgloss.Color("#D8F3DC"),
    Success:    lipgloss.Color("#B7E4C7"),
}

var ElectricBloom = Theme{
    Name:       "Electric Bloom",
    Primary:    lipgloss.Color("#FF0080"),
    Secondary:  lipgloss.Color("#00D4FF"),
    Accent:     lipgloss.Color("#FFE600"),
    Background: lipgloss.Color("#0D0221"),
    Text:       lipgloss.Color("#F0F3FF"),
    Success:    lipgloss.Color("#39FF14"),
}

var PlasmaPulse = Theme{
    Name:       "Plasma Pulse",
    Primary:    lipgloss.Color("#39FF14"),
    Secondary:  lipgloss.Color("#00F5FF"),
    Accent:     lipgloss.Color("#FF1493"),
    Background: lipgloss.Color("#0A0118"),
    Text:       lipgloss.Color("#E0FFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}

var IndigoDepths = Theme{
    Name:       "Indigo Depths",
    Primary:    lipgloss.Color("#4169E1"),
    Secondary:  lipgloss.Color("#5F9EA0"),
    Accent:     lipgloss.Color("#DEB887"),
    Background: lipgloss.Color("#0C0C1E"),
    Text:       lipgloss.Color("#F0E68C"),
    Success:    lipgloss.Color("#5F9EA0"),
}

var SageMeadow = Theme{
    Name:       "Sage Meadow",
    Primary:    lipgloss.Color("#8FBC8F"),
    Secondary:  lipgloss.Color("#D2B48C"),
    Accent:     lipgloss.Color("#F4A460"),
    Background: lipgloss.Color("#1A1612"),
    Text:       lipgloss.Color("#FAF0E6"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

---

#### 3. `internal/theme/manager.go`

```go
package theme

import (
    "sync"
)

// Manager handles theme switching at runtime
type Manager struct {
    current Theme
    mu      sync.RWMutex
}

// NewManager creates a new theme manager with default theme
func NewManager() *Manager {
    return &Manager{
        current: GetTheme(DefaultTheme),
    }
}

// Current returns the currently active theme (thread-safe)
func (m *Manager) Current() Theme {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.current
}

// SetTheme changes the active theme (thread-safe)
func (m *Manager) SetTheme(name string) bool {
    theme := GetTheme(name)
    m.mu.Lock()
    m.current = theme
    m.mu.Unlock()
    return true
}

// SetThemeByName changes theme and returns success status
func (m *Manager) SetThemeByName(name string) (Theme, bool) {
    if _, ok := Registry[name]; !ok {
        return m.Current(), false
    }
    m.SetTheme(name)
    return m.Current(), true
}

// Global manager instance (singleton pattern)
var globalManager *Manager
var once sync.Once

// GetManager returns the global theme manager instance
func GetManager() *Manager {
    once.Do(func() {
        globalManager = NewManager()
    })
    return globalManager
}
```

---

#### 4. `internal/ui/styles.go`

```go
package ui

import (
    "github.com/charmbracelet/lipgloss"
    "noise.sh/internal/theme"
)

// StyleManager holds all UI styles using current theme
type StyleManager struct {
    themeManager *theme.Manager
}

// NewStyleManager creates a new style manager
func NewStyleManager() *StyleManager {
    return &StyleManager{
        themeManager: theme.GetManager(),
    }
}

// GetCurrentTheme returns the active theme
func (sm *StyleManager) GetCurrentTheme() theme.Theme {
    return sm.themeManager.Current()
}

// TitleStyle returns style for titles
func (sm *StyleManager) TitleStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Accent).
        Background(t.Background).
        Bold(true).
        Padding(0, 1)
}

// BorderStyle returns style for borders
func (sm *StyleManager) BorderStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(t.Primary).
        Padding(1, 2)
}

// ActiveTabStyle returns style for active tab
func (sm *StyleManager) ActiveTabStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Primary).
        Bold(true).
        Padding(0, 2)
}

// InactiveTabStyle returns style for inactive tabs
func (sm *StyleManager) InactiveTabStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Secondary).
        Background(t.Background).
        Padding(0, 2)
}

// SuccessStyle returns style for success messages
func (sm *StyleManager) SuccessStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Success).
        Background(t.Background).
        Bold(true)
}

// ErrorStyle returns style for error messages
func (sm *StyleManager) ErrorStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF6B6B")).
        Background(t.Background).
        Bold(true)
}

// TextStyle returns basic text style
func (sm *StyleManager) TextStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Background)
}

// HighlightStyle returns style for highlighted text
func (sm *StyleManager) HighlightStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Background).
        Background(t.Accent).
        Bold(true)
}
```

---

#### 5. `config/config.go` (Theme Persistence)

```go
package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    Theme string `json:"theme"`
    // Other config fields...
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".noise.sh", "config.json"), nil
}

// LoadConfig loads configuration from disk
func LoadConfig() (*Config, error) {
    path, err := GetConfigPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        // Return default config if file doesn't exist
        if os.IsNotExist(err) {
            return &Config{Theme: "amber-night"}, nil
        }
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// SaveConfig saves configuration to disk
func (c *Config) SaveConfig() error {
    path, err := GetConfigPath()
    if err != nil {
        return err
    }

    // Create directory if it doesn't exist
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}
```

---

## Lipgloss Integration

### Key Integration Points

1. **Initialize theme manager on app startup**
2. **Load user's saved theme preference**
3. **Pass theme to all UI components**
4. **Allow runtime theme switching with keybind (e.g., `Ctrl+T`)**
5. **Re-render UI when theme changes**

### Example: Main App Integration

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "noise.sh/internal/theme"
    "noise.sh/internal/ui"
    "noise.sh/config"
)

type model struct {
    styleManager *ui.StyleManager
    currentTheme string
    // ... other fields
}

func initialModel() model {
    // Load config
    cfg, _ := config.LoadConfig()
    
    // Set theme from config
    themeManager := theme.GetManager()
    themeManager.SetTheme(cfg.Theme)
    
    return model{
        styleManager: ui.NewStyleManager(),
        currentTheme: cfg.Theme,
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+t":
            // Cycle through themes
            themes := theme.ListThemes()
            currentIdx := indexOf(themes, m.currentTheme)
            nextIdx := (currentIdx + 1) % len(themes)
            nextTheme := themes[nextIdx]
            
            theme.GetManager().SetTheme(nextTheme)
            m.currentTheme = nextTheme
            
            // Save to config
            cfg := &config.Config{Theme: nextTheme}
            cfg.SaveConfig()
            
            return m, nil
        }
    }
    return m, nil
}

func (m model) View() string {
    // Use styleManager to get themed styles
    title := m.styleManager.TitleStyle().Render("noise.sh")
    
    // Current theme indicator
    themeInfo := m.styleManager.HighlightStyle().
        Render("Theme: " + m.currentTheme)
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        themeInfo,
        // ... rest of UI
    )
}
```

---

## Testing Checklist

### Visual Testing
- [ ] All 12 themes render correctly
- [ ] Text is readable on all backgrounds (WCAG AA)
- [ ] Primary/Secondary/Accent colors are distinct
- [ ] Success messages use correct green color
- [ ] Borders and UI elements use theme colors
- [ ] No color bleeding or contrast issues

### Functional Testing
- [ ] Theme switching works at runtime (Ctrl+T)
- [ ] Theme preference persists after restart
- [ ] Default theme loads correctly on first run
- [ ] Invalid theme names fall back to default
- [ ] Theme manager is thread-safe
- [ ] Config file saves/loads correctly

### Accessibility Testing
- [ ] All themes pass WCAG AA contrast requirements
- [ ] Text remains readable in all themes
- [ ] UI elements have sufficient contrast
- [ ] Color is not the only indicator (use icons/text too)

### Performance Testing
- [ ] Theme switching is instant (< 50ms)
- [ ] No memory leaks from theme changes
- [ ] UI re-renders efficiently

---

## Implementation Roadmap

### Single Phase: All 12 Themes
**Priority:** Implement all at once  
**Themes:** All 12 themes listed above

**Tasks:**
1. Create `internal/theme/` package structure
2. Implement Theme struct and all 12 theme definitions in registry
3. Create theme manager with switching logic
4. Integrate with existing UI components
5. Add Ctrl+T keybind for theme cycling
6. Implement custom theme loading from config
7. Add theme preview/selector UI
8. Test readability and contrast for all themes
9. Document all themes in help menu
10. Create example custom theme file

**Acceptance Criteria:**
- All 12 themes implemented and working
- User can switch between all themes at runtime
- Theme preference persists after restart
- Custom theme support via config file
- All themes are readable (WCAG AA compliant)
- UI consistently uses theme colors throughout
- Help menu shows theme list and keybinds
- Documentation includes custom theme creation guide

---

## Additional Features (Optional)

### Custom Theme Creation
**IMPORTANT:** Users can create custom themes via the config file.

#### How Custom Themes Work

1. **Config File Location:** `~/.noise.sh/config.json`
2. **Structure:** Users can define custom themes in the `custom_themes` section
3. **Theme Selection:** Custom themes appear in the theme list alongside built-in themes
4. **Validation:** System validates that all required colors are present and valid hex codes

#### Example Config with Custom Theme

```json
{
  "theme": "my-custom-sunset",
  "custom_themes": {
    "my-custom-sunset": {
      "name": "My Custom Sunset",
      "primary": "#FF6B9D",
      "secondary": "#FF8C42",
      "accent": "#FFC312",
      "background": "#1e1e2e",
      "text": "#FFEAA7",
      "success": "#55E6C1"
    },
    "my-ocean-theme": {
      "name": "Ocean Breeze",
      "primary": "#0077BE",
      "secondary": "#4A90E2",
      "accent": "#87CEEB",
      "background": "#001F3F",
      "text": "#E0F7FA",
      "success": "#4CAF50"
    }
  }
}
```

#### Implementation in `config/config.go`

```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
)

type CustomTheme struct {
    Name       string `json:"name"`
    Primary    string `json:"primary"`
    Secondary  string `json:"secondary"`
    Accent     string `json:"accent"`
    Background string `json:"background"`
    Text       string `json:"text"`
    Success    string `json:"success"`
}

type Config struct {
    Theme        string                  `json:"theme"`
    CustomThemes map[string]CustomTheme  `json:"custom_themes,omitempty"`
}

// ValidateHexColor checks if a string is a valid hex color
func ValidateHexColor(color string) bool {
    hexPattern := regexp.MustCompile(`^#[0-9A-Fa-f]{6}# noise.sh Theme System Enhancement
## Implementation Specification for Kilocode AI

**Version:** 1.0  
**Date:** October 17, 2025  
**Purpose:** Complete specification for implementing 12 cohesive themes in noise.sh TUI

---

## Table of Contents

1. [Overview](#overview)
2. [Theme Philosophy](#theme-philosophy)
3. [Complete Theme Definitions](#complete-theme-definitions)
4. [Implementation Requirements](#implementation-requirements)
5. [Code Structure](#code-structure)
6. [Lipgloss Integration](#lipgloss-integration)
7. [Testing Checklist](#testing-checklist)
8. [Priority Roadmap](#priority-roadmap)

---

## Overview

### Goal
Implement a dynamic theme system with 12 professionally designed, WCAG AA compliant color schemes that users can switch between during runtime.

### Key Requirements
- ✅ All themes follow Element/Nature naming convention
- ✅ Zero music genre references in names
- ✅ WCAG AA compliant (minimum 3:1 contrast for UI elements, 4.5:1 for text)
- ✅ Each theme has 6 distinct colors: Primary, Secondary, Accent, Background, Text, Success
- ✅ Themes map to musical genres but names are universal
- ✅ Runtime theme switching without restart

---

## Theme Philosophy

### Naming Pattern
All themes use **Element/Nature** naming:
- **Materials/Elements:** Slate Mist, Amber, Clay, Iron, Jade, Indigo, Sage, Plasma
- **Natural Phenomena:** Night, Storm, Tide, Ember, Whisper, Bloom, Depths, Meadow
- **Energy Modifiers:** Molten, Electric, Sunset, Forest

### Color Distinctions
- **Slate Mist** - Grayscale (neutral)
- **Amber Night** - Tan/Purple (warm sophistication)
- **Molten Gold** - Gold/Orange/Purple (bold energy)
- **Clay Roads** - Burnt Orange/Brown (rustic)
- **Iron Storm** - Crimson/Gray/Fire Orange (dark power)
- **Jade Tide** - Cool Teal/Aqua (elegant water)
- **Sunset Ember** - Pink/Orange/Yellow (warm creative)
- **Forest Whisper** - Warm Sage Green (calm natural)
- **Electric Bloom** - Hot Pink/Cyan/Yellow (high energy)
- **Plasma Pulse** - Neon Green/Cyan/Magenta (futuristic)
- **Indigo Depths** - Royal Blue/Cadet Blue (deep soulful)
- **Sage Meadow** - Sage Green/Tan/Sandy (peaceful earth)

---

## Complete Theme Definitions

### 1. Slate Mist (Universal - Minimal)

```go
type Theme struct {
    Name       string
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Accent     lipgloss.Color
    Background lipgloss.Color
    Text       lipgloss.Color
    Success    lipgloss.Color
}

var SlateMist = Theme{
    Name:       "Slate Mist",
    Primary:    lipgloss.Color("#E0E0E0"),
    Secondary:  lipgloss.Color("#B0B0B0"),
    Accent:     lipgloss.Color("#FFFFFF"),
    Background: lipgloss.Color("#1A1A1A"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#90C695"),
}
```

**Use Cases:** Focus mode, no distraction, terminal purist  
**Mood:** Clean, minimal, professional

---

### 2. Amber Night (Jazz)

```go
var AmberNight = Theme{
    Name:       "Amber Night",
    Primary:    lipgloss.Color("#B8936E"),
    Secondary:  lipgloss.Color("#6D5B8B"),
    Accent:     lipgloss.Color("#E8C547"),
    Background: lipgloss.Color("#12101A"),
    Text:       lipgloss.Color("#F0E6D8"),
    Success:    lipgloss.Color("#7FB89C"),
}
```

**Use Cases:** Jazz, late night sessions, complex lyrics  
**Mood:** Sophisticated, warm, mysterious  
**Note:** Fusion of tan warmth + purple sophistication

---

### 3. Molten Gold (Hip-Hop)

```go
var MoltenGold = Theme{
    Name:       "Molten Gold",
    Primary:    lipgloss.Color("#FFD700"),
    Secondary:  lipgloss.Color("#FF6600"),
    Accent:     lipgloss.Color("#8B00FF"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#FFFFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}
```

**Use Cases:** Rap bars, confident flow, wordplay  
**Mood:** Bold, confident, royal  
**Note:** Gold primary with orange secondary and purple accent for variety

---

### 4. Clay Roads (Country)

```go
var ClayRoads = Theme{
    Name:       "Clay Roads",
    Primary:    lipgloss.Color("#CC5500"),
    Secondary:  lipgloss.Color("#D2691E"),
    Accent:     lipgloss.Color("#FFD700"),
    Background: lipgloss.Color("#1C1410"),
    Text:       lipgloss.Color("#FFF8DC"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

**Use Cases:** Country, storytelling, honest lyrics  
**Mood:** Warm, rustic, heartland  
**Note:** Burnt orange with chocolate orange secondary

---

### 5. Iron Storm (Metal/Rock)

```go
var IronStorm = Theme{
    Name:       "Iron Storm",
    Primary:    lipgloss.Color("#DC143C"),
    Secondary:  lipgloss.Color("#708090"),
    Accent:     lipgloss.Color("#FF4500"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#4A6741"),
}
```

**Use Cases:** Metal, hard rock, dark themes  
**Mood:** Dark, powerful, intense
**Note:** Crimson + slate mist gray with FIRE ORANGE accent for explosive pop

---

### 6. Jade Tide (Classical)

```go
var JadeTide = Theme{
    Name:       "Jade Tide",
    Primary:    lipgloss.Color("#20B2AA"),
    Secondary:  lipgloss.Color("#2E8B87"),
    Accent:     lipgloss.Color("#E0F2F1"),
    Background: lipgloss.Color("#0F1419"),
    Text:       lipgloss.Color("#F5F8FA"),
    Success:    lipgloss.Color("#7FB89C"),
}
```

**Use Cases:** Classical, poetic, professional  
**Mood:** Elegant, refined, sophisticated  
**Note:** Cool aquatic teal - distinct from Forest Whisper (warm green) and Indigo Depths (blue)

---

### 7. Sunset Ember (Universal - Creative)

```go
var SunsetEmber = Theme{
    Name:       "Sunset Ember",
    Primary:    lipgloss.Color("#FF6B9D"),
    Secondary:  lipgloss.Color("#FF8C42"),
    Accent:     lipgloss.Color("#FFC312"),
    Background: lipgloss.Color("#1e1e2e"),
    Text:       lipgloss.Color("#FFEAA7"),
    Success:    lipgloss.Color("#55E6C1"),
}
```

**Use Cases:** Cross-genre, creative sessions, evening writing  
**Mood:** Warm, inspiring, sunset vibes  
**Note:** Pink to orange gradient feel

---

### 8. Forest Whisper (Universal - Calm)

```go
var ForestWhisper = Theme{
    Name:       "Forest Whisper",
    Primary:    lipgloss.Color("#52B788"),
    Secondary:  lipgloss.Color("#52A068"),
    Accent:     lipgloss.Color("#95D5B2"),
    Background: lipgloss.Color("#1B263B"),
    Text:       lipgloss.Color("#D8F3DC"),
    Success:    lipgloss.Color("#B7E4C7"),
}
```

**Use Cases:** Focus mode, long sessions, nature themes  
**Mood:** Calm, focused, natural  
**Note:** Warm sage greens - distinct from Jade Tide's cool teal

---

### 9. Electric Bloom (Pop/EDM)

```go
var ElectricBloom = Theme{
    Name:       "Electric Bloom",
    Primary:    lipgloss.Color("#FF0080"),
    Secondary:  lipgloss.Color("#00D4FF"),
    Accent:     lipgloss.Color("#FFE600"),
    Background: lipgloss.Color("#0D0221"),
    Text:       lipgloss.Color("#F0F3FF"),
    Success:    lipgloss.Color("#39FF14"),
}
```

**Use Cases:** Pop songs, EDM, high energy writing  
**Mood:** Bold, powerful, electric  
**Note:** Hot pink + cyan + yellow for maximum vibrancy

---

### 10. Plasma Pulse (Electronic)

```go
var PlasmaPulse = Theme{
    Name:       "Plasma Pulse",
    Primary:    lipgloss.Color("#39FF14"),
    Secondary:  lipgloss.Color("#00F5FF"),
    Accent:     lipgloss.Color("#FF1493"),
    Background: lipgloss.Color("#0A0118"),
    Text:       lipgloss.Color("#E0FFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}
```

**Use Cases:** EDM, techno, electronic music  
**Mood:** Futuristic, energetic, digital  
**Note:** Neon green + cyan with hot pink accent

---

### 11. Indigo Depths (Blues/Soul)

```go
var IndigoDepths = Theme{
    Name:       "Indigo Depths",
    Primary:    lipgloss.Color("#4169E1"),
    Secondary:  lipgloss.Color("#5F9EA0"),
    Accent:     lipgloss.Color("#DEB887"),
    Background: lipgloss.Color("#0C0C1E"),
    Text:       lipgloss.Color("#F0E68C"),
    Success:    lipgloss.Color("#5F9EA0"),
}
```

**Use Cases:** Blues, R&B, soul music  
**Mood:** Soulful, melancholic, deep  
**Note:** Royal blue + cadet blue with warm tan accent

---

### 12. Sage Meadow (Folk/Acoustic)

```go
var SageMeadow = Theme{
    Name:       "Sage Meadow",
    Primary:    lipgloss.Color("#8FBC8F"),
    Secondary:  lipgloss.Color("#D2B48C"),
    Accent:     lipgloss.Color("#F4A460"),
    Background: lipgloss.Color("#1A1612"),
    Text:       lipgloss.Color("#FAF0E6"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

**Use Cases:** Folk, acoustic, nature-themed lyrics  
**Mood:** Natural, peaceful, grounded  
**Note:** Sage green + tan/sandy colors

---

## Implementation Requirements

### File Structure

```
noise.sh/
├── internal/
│   ├── theme/
│   │   ├── theme.go          # Theme struct and definitions
│   │   ├── manager.go        # Theme switching logic
│   │   └── registry.go       # All 12 themes registered
│   └── ui/
│       ├── styles.go         # Lipgloss styles using current theme
│       └── components.go     # UI components
├── config/
│   └── config.go             # User's theme preference saved here
```

### Core Implementation Files

#### 1. `internal/theme/theme.go`

```go
package theme

import "github.com/charmbracelet/lipgloss"

// Theme represents a color scheme for the TUI
type Theme struct {
    Name       string
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Accent     lipgloss.Color
    Background lipgloss.Color
    Text       lipgloss.Color
    Success    lipgloss.Color
}

// GetStyle returns a lipgloss style with the given colors
func (t *Theme) GetStyle(fg, bg lipgloss.Color) lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(fg).
        Background(bg)
}

// PrimaryStyle returns a style using primary color
func (t *Theme) PrimaryStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Primary)
}

// SecondaryStyle returns a style using secondary color
func (t *Theme) SecondaryStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Secondary)
}

// AccentStyle returns a style using accent color
func (t *Theme) AccentStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Background).
        Background(t.Accent)
}

// SuccessStyle returns a style for success messages
func (t *Theme) SuccessStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Foreground(t.Success).
        Background(t.Background)
}
```

---

#### 2. `internal/theme/registry.go`

```go
package theme

// Registry holds all available themes
var Registry = map[string]Theme{
    "slate":         SlateMist,
    "amber-night":   AmberNight,
    "molten-gold":   MoltenGold,
    "clay-roads":    ClayRoads,
    "iron-storm":    IronStorm,
    "jade-tide":     JadeTide,
    "sunset-ember":  SunsetEmber,
    "forest-whisper": ForestWhisper,
    "electric-bloom": ElectricBloom,
    "plasma-pulse":  PlasmaPulse,
    "indigo-depths": IndigoDepths,
    "sage-meadow":   SageMeadow,
}

// Default theme
var DefaultTheme = "amber-night"

// GetTheme returns a theme by name, or default if not found
func GetTheme(name string) Theme {
    if theme, ok := Registry[name]; ok {
        return theme
    }
    return Registry[DefaultTheme]
}

// ListThemes returns all available theme names
func ListThemes() []string {
    themes := make([]string, 0, len(Registry))
    for name := range Registry {
        themes = append(themes, name)
    }
    return themes
}

// Define all 12 themes
var SlateMist = Theme{
    Name:       "Slate Mist",
    Primary:    lipgloss.Color("#E0E0E0"),
    Secondary:  lipgloss.Color("#B0B0B0"),
    Accent:     lipgloss.Color("#FFFFFF"),
    Background: lipgloss.Color("#1A1A1A"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#90C695"),
}

var AmberNight = Theme{
    Name:       "Amber Night",
    Primary:    lipgloss.Color("#B8936E"),
    Secondary:  lipgloss.Color("#6D5B8B"),
    Accent:     lipgloss.Color("#E8C547"),
    Background: lipgloss.Color("#12101A"),
    Text:       lipgloss.Color("#F0E6D8"),
    Success:    lipgloss.Color("#7FB89C"),
}

var MoltenGold = Theme{
    Name:       "Molten Gold",
    Primary:    lipgloss.Color("#FFD700"),
    Secondary:  lipgloss.Color("#FF6600"),
    Accent:     lipgloss.Color("#8B00FF"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#FFFFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}

var ClayRoads = Theme{
    Name:       "Clay Roads",
    Primary:    lipgloss.Color("#CC5500"),
    Secondary:  lipgloss.Color("#D2691E"),
    Accent:     lipgloss.Color("#FFD700"),
    Background: lipgloss.Color("#1C1410"),
    Text:       lipgloss.Color("#FFF8DC"),
    Success:    lipgloss.Color("#9ACD32"),
}

var IronStorm = Theme{
    Name:       "Iron Storm",
    Primary:    lipgloss.Color("#DC143C"),
    Secondary:  lipgloss.Color("#708090"),
    Accent:     lipgloss.Color("#FF4500"),
    Background: lipgloss.Color("#0D0D0D"),
    Text:       lipgloss.Color("#E8E8E8"),
    Success:    lipgloss.Color("#4A6741"),
}

var JadeTide = Theme{
    Name:       "Jade Tide",
    Primary:    lipgloss.Color("#20B2AA"),
    Secondary:  lipgloss.Color("#2E8B87"),
    Accent:     lipgloss.Color("#E0F2F1"),
    Background: lipgloss.Color("#0F1419"),
    Text:       lipgloss.Color("#F5F8FA"),
    Success:    lipgloss.Color("#7FB89C"),
}

var SunsetEmber = Theme{
    Name:       "Sunset Ember",
    Primary:    lipgloss.Color("#FF6B9D"),
    Secondary:  lipgloss.Color("#FF8C42"),
    Accent:     lipgloss.Color("#FFC312"),
    Background: lipgloss.Color("#1e1e2e"),
    Text:       lipgloss.Color("#FFEAA7"),
    Success:    lipgloss.Color("#55E6C1"),
}

var ForestWhisper = Theme{
    Name:       "Forest Whisper",
    Primary:    lipgloss.Color("#52B788"),
    Secondary:  lipgloss.Color("#52A068"),
    Accent:     lipgloss.Color("#95D5B2"),
    Background: lipgloss.Color("#1B263B"),
    Text:       lipgloss.Color("#D8F3DC"),
    Success:    lipgloss.Color("#B7E4C7"),
}

var ElectricBloom = Theme{
    Name:       "Electric Bloom",
    Primary:    lipgloss.Color("#FF0080"),
    Secondary:  lipgloss.Color("#00D4FF"),
    Accent:     lipgloss.Color("#FFE600"),
    Background: lipgloss.Color("#0D0221"),
    Text:       lipgloss.Color("#F0F3FF"),
    Success:    lipgloss.Color("#39FF14"),
}

var PlasmaPulse = Theme{
    Name:       "Plasma Pulse",
    Primary:    lipgloss.Color("#39FF14"),
    Secondary:  lipgloss.Color("#00F5FF"),
    Accent:     lipgloss.Color("#FF1493"),
    Background: lipgloss.Color("#0A0118"),
    Text:       lipgloss.Color("#E0FFFF"),
    Success:    lipgloss.Color("#00FF7F"),
}

var IndigoDepths = Theme{
    Name:       "Indigo Depths",
    Primary:    lipgloss.Color("#4169E1"),
    Secondary:  lipgloss.Color("#5F9EA0"),
    Accent:     lipgloss.Color("#DEB887"),
    Background: lipgloss.Color("#0C0C1E"),
    Text:       lipgloss.Color("#F0E68C"),
    Success:    lipgloss.Color("#5F9EA0"),
}

var SageMeadow = Theme{
    Name:       "Sage Meadow",
    Primary:    lipgloss.Color("#8FBC8F"),
    Secondary:  lipgloss.Color("#D2B48C"),
    Accent:     lipgloss.Color("#F4A460"),
    Background: lipgloss.Color("#1A1612"),
    Text:       lipgloss.Color("#FAF0E6"),
    Success:    lipgloss.Color("#9ACD32"),
}
```

---

#### 3. `internal/theme/manager.go`

```go
package theme

import (
    "sync"
)

// Manager handles theme switching at runtime
type Manager struct {
    current Theme
    mu      sync.RWMutex
}

// NewManager creates a new theme manager with default theme
func NewManager() *Manager {
    return &Manager{
        current: GetTheme(DefaultTheme),
    }
}

// Current returns the currently active theme (thread-safe)
func (m *Manager) Current() Theme {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.current
}

// SetTheme changes the active theme (thread-safe)
func (m *Manager) SetTheme(name string) bool {
    theme := GetTheme(name)
    m.mu.Lock()
    m.current = theme
    m.mu.Unlock()
    return true
}

// SetThemeByName changes theme and returns success status
func (m *Manager) SetThemeByName(name string) (Theme, bool) {
    if _, ok := Registry[name]; !ok {
        return m.Current(), false
    }
    m.SetTheme(name)
    return m.Current(), true
}

// Global manager instance (singleton pattern)
var globalManager *Manager
var once sync.Once

// GetManager returns the global theme manager instance
func GetManager() *Manager {
    once.Do(func() {
        globalManager = NewManager()
    })
    return globalManager
}
```

---

#### 4. `internal/ui/styles.go`

```go
package ui

import (
    "github.com/charmbracelet/lipgloss"
    "noise.sh/internal/theme"
)

// StyleManager holds all UI styles using current theme
type StyleManager struct {
    themeManager *theme.Manager
}

// NewStyleManager creates a new style manager
func NewStyleManager() *StyleManager {
    return &StyleManager{
        themeManager: theme.GetManager(),
    }
}

// GetCurrentTheme returns the active theme
func (sm *StyleManager) GetCurrentTheme() theme.Theme {
    return sm.themeManager.Current()
}

// TitleStyle returns style for titles
func (sm *StyleManager) TitleStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Accent).
        Background(t.Background).
        Bold(true).
        Padding(0, 1)
}

// BorderStyle returns style for borders
func (sm *StyleManager) BorderStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(t.Primary).
        Padding(1, 2)
}

// ActiveTabStyle returns style for active tab
func (sm *StyleManager) ActiveTabStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Primary).
        Bold(true).
        Padding(0, 2)
}

// InactiveTabStyle returns style for inactive tabs
func (sm *StyleManager) InactiveTabStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Secondary).
        Background(t.Background).
        Padding(0, 2)
}

// SuccessStyle returns style for success messages
func (sm *StyleManager) SuccessStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Success).
        Background(t.Background).
        Bold(true)
}

// ErrorStyle returns style for error messages
func (sm *StyleManager) ErrorStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FF6B6B")).
        Background(t.Background).
        Bold(true)
}

// TextStyle returns basic text style
func (sm *StyleManager) TextStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Text).
        Background(t.Background)
}

// HighlightStyle returns style for highlighted text
func (sm *StyleManager) HighlightStyle() lipgloss.Style {
    t := sm.GetCurrentTheme()
    return lipgloss.NewStyle().
        Foreground(t.Background).
        Background(t.Accent).
        Bold(true)
}
```

---

#### 5. `config/config.go` (Theme Persistence)

```go
package config

import (
    "encoding/json"
    "os"
    "path/filepath"
)

type Config struct {
    Theme string `json:"theme"`
    // Other config fields...
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".noise.sh", "config.json"), nil
}

// LoadConfig loads configuration from disk
func LoadConfig() (*Config, error) {
    path, err := GetConfigPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        // Return default config if file doesn't exist
        if os.IsNotExist(err) {
            return &Config{Theme: "amber-night"}, nil
        }
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// SaveConfig saves configuration to disk
func (c *Config) SaveConfig() error {
    path, err := GetConfigPath()
    if err != nil {
        return err
    }

    // Create directory if it doesn't exist
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}
```

---

## Lipgloss Integration

### Key Integration Points

1. **Initialize theme manager on app startup**
2. **Load user's saved theme preference**
3. **Pass theme to all UI components**
4. **Allow runtime theme switching with keybind (e.g., `Ctrl+T`)**
5. **Re-render UI when theme changes**

### Example: Main App Integration

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "noise.sh/internal/theme"
    "noise.sh/internal/ui"
    "noise.sh/config"
)

type model struct {
    styleManager *ui.StyleManager
    currentTheme string
    // ... other fields
}

func initialModel() model {
    // Load config
    cfg, _ := config.LoadConfig()
    
    // Set theme from config
    themeManager := theme.GetManager()
    themeManager.SetTheme(cfg.Theme)
    
    return model{
        styleManager: ui.NewStyleManager(),
        currentTheme: cfg.Theme,
    }
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+t":
            // Cycle through themes
            themes := theme.ListThemes()
            currentIdx := indexOf(themes, m.currentTheme)
            nextIdx := (currentIdx + 1) % len(themes)
            nextTheme := themes[nextIdx]
            
            theme.GetManager().SetTheme(nextTheme)
            m.currentTheme = nextTheme
            
            // Save to config
            cfg := &config.Config{Theme: nextTheme}
            cfg.SaveConfig()
            
            return m, nil
        }
    }
    return m, nil
}

func (m model) View() string {
    // Use styleManager to get themed styles
    title := m.styleManager.TitleStyle().Render("noise.sh")
    
    // Current theme indicator
    themeInfo := m.styleManager.HighlightStyle().
        Render("Theme: " + m.currentTheme)
    
    return lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        themeInfo,
        // ... rest of UI
    )
}
```

---

## Testing Checklist

### Visual Testing
- [ ] All 12 themes render correctly
- [ ] Text is readable on all backgrounds (WCAG AA)
- [ ] Primary/Secondary/Accent colors are distinct
- [ ] Success messages use correct green color
- [ ] Borders and UI elements use theme colors
- [ ] No color bleeding or contrast issues

### Functional Testing
- [ ] Theme switching works at runtime (Ctrl+T)
- [ ] Theme preference persists after restart
- [ ] Default theme loads correctly on first run
- [ ] Invalid theme names fall back to default
- [ ] Theme manager is thread-safe
- [ ] Config file saves/loads correctly

### Accessibility Testing
- [ ] All themes pass WCAG AA contrast requirements
- [ ] Text remains readable in all themes
- [ ] UI elements have sufficient contrast
- [ ] Color is not the only indicator (use icons/text too)

### Performance Testing
- [ ] Theme switching is instant (< 50ms)
- [ ] No memory leaks from theme changes
- [ ] UI re-renders efficiently

---

## Implementation Roadmap

### Single Phase: All 12 Themes
**Priority:** Implement all at once  
**Themes:** All 12 themes listed above

**Tasks:**
1. Create `internal/theme/` package structure
2. Implement Theme struct and all 12 theme definitions in registry
3. Create theme manager with switching logic
4. Integrate with existing UI components
5. Add Ctrl+T keybind for theme cycling
6. Implement custom theme loading from config
7. Add theme preview/selector UI
8. Test readability and contrast for all themes
9. Document all themes in help menu
10. Create example custom theme file

**Acceptance Criteria:**
- All 12 themes implemented and working
- User can switch between all themes at runtime
- Theme preference persists after restart
- Custom theme support via config file
- All themes are readable (WCAG AA compliant)
- UI consistently uses theme colors throughout
- Help menu shows theme list and keybinds
- Documentation includes custom theme creation guide

---

)
    return hexPattern.MatchString(color)
}

// ValidateCustomTheme checks if a custom theme has all required fields
func (ct *CustomTheme) Validate() error {
    if ct.Name == "" {
        return fmt.Errorf("theme name is required")
    }
    
    colors := map[string]string{
        "primary":    ct.Primary,
        "secondary":  ct.Secondary,
        "accent":     ct.Accent,
        "background": ct.Background,
        "text":       ct.Text,
        "success":    ct.Success,
    }
    
    for name, color := range colors {
        if !ValidateHexColor(color) {
            return fmt.Errorf("invalid hex color for %s: %s", name, color)
        }
    }
    
    return nil
}

// GetConfigPath returns the path to the config file
func GetConfigPath() (string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(home, ".noise.sh", "config.json"), nil
}

// LoadConfig loads configuration from disk
func LoadConfig() (*Config, error) {
    path, err := GetConfigPath()
    if err != nil {
        return nil, err
    }

    data, err := os.ReadFile(path)
    if err != nil {
        // Return default config if file doesn't exist
        if os.IsNotExist(err) {
            return &Config{
                Theme:        "amber-night",
                CustomThemes: make(map[string]CustomTheme),
            }, nil
        }
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    
    // Initialize CustomThemes map if nil
    if cfg.CustomThemes == nil {
        cfg.CustomThemes = make(map[string]CustomTheme)
    }

    return &cfg, nil
}

// SaveConfig saves configuration to disk
func (c *Config) SaveConfig() error {
    path, err := GetConfigPath()
    if err != nil {
        return err
    }

    // Create directory if it doesn't exist
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }

    return os.WriteFile(path, data, 0644)
}

// GetCustomThemes returns all custom themes with validation
func (c *Config) GetCustomThemes() (map[string]CustomTheme, []error) {
    validThemes := make(map[string]CustomTheme)
    var errors []error
    
    for id, theme := range c.CustomThemes {
        if err := theme.Validate(); err != nil {
            errors = append(errors, fmt.Errorf("invalid custom theme '%s': %w", id, err))
        } else {
            validThemes[id] = theme
        }
    }
    
    return validThemes, errors
}
```

#### Update `internal/theme/registry.go` to Support Custom Themes

```go
package theme

import (
    "github.com/charmbracelet/lipgloss"
    "noise.sh/config"
)

// Registry holds all available themes (built-in + custom)
var Registry = make(map[string]Theme)

// Initialize built-in themes
func init() {
    // Register all 12 built-in themes
    Registry["slate"] = SlateMist
    Registry["amber-night"] = AmberNight
    Registry["molten-gold"] = MoltenGold
    Registry["clay-roads"] = ClayRoads
    Registry["iron-storm"] = IronStorm
    Registry["jade-tide"] = JadeTide
    Registry["sunset-ember"] = SunsetEmber
    Registry["forest-whisper"] = ForestWhisper
    Registry["electric-bloom"] = ElectricBloom
    Registry["plasma-pulse"] = PlasmaPulse
    Registry["indigo-depths"] = IndigoDepths
    Registry["sage-meadow"] = SageMeadow
}

// LoadCustomThemes loads custom themes from config
func LoadCustomThemes(cfg *config.Config) []error {
    customThemes, errors := cfg.GetCustomThemes()
    
    for id, customTheme := range customThemes {
        // Convert custom theme to Theme struct
        theme := Theme{
            Name:       customTheme.Name,
            Primary:    lipgloss.Color(customTheme.Primary),
            Secondary:  lipgloss.Color(customTheme.Secondary),
            Accent:     lipgloss.Color(customTheme.Accent),
            Background: lipgloss.Color(customTheme.Background),
            Text:       lipgloss.Color(customTheme.Text),
            Success:    lipgloss.Color(customTheme.Success),
        }
        
        // Add to registry
        Registry[id] = theme
    }
    
    return errors
}

// Default theme
var DefaultTheme = "amber-night"

// GetTheme returns a theme by name, or default if not found
func GetTheme(name string) Theme {
    if theme, ok := Registry[name]; ok {
        return theme
    }
    return Registry[DefaultTheme]
}

// ListThemes returns all available theme names (built-in + custom)
func ListThemes() []string {
    themes := make([]string, 0, len(Registry))
    for name := range Registry {
        themes = append(themes, name)
    }
    return themes
}

// IsCustomTheme checks if a theme is custom (not built-in)
func IsCustomTheme(name string) bool {
    builtInThemes := []string{
        "slate", "amber-night", "molten-gold", "clay-roads",
        "iron-storm", "jade-tide", "sunset-ember", "forest-whisper",
        "electric-bloom", "plasma-pulse", "indigo-depths", "sage-meadow",
    }
    
    for _, builtIn := range builtInThemes {
        if name == builtIn {
            return false
        }
    }
    
    return true
}

// ... (rest of theme definitions remain the same)
```

#### Update Main App to Load Custom Themes

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "noise.sh/internal/theme"
    "noise.sh/internal/ui"
    "noise.sh/config"
    "log"
)

type model struct {
    styleManager *ui.StyleManager
    currentTheme string
    // ... other fields
}

func initialModel() model {
    // Load config
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Printf("Error loading config: %v", err)
        cfg = &config.Config{Theme: "amber-night"}
    }
    
    // Load custom themes from config
    if errors := theme.LoadCustomThemes(cfg); len(errors) > 0 {
        for _, err := range errors {
            log.Printf("Custom theme error: %v", err)
        }
    }
    
    // Set theme from config
    themeManager := theme.GetManager()
    themeManager.SetTheme(cfg.Theme)
    
    return model{
        styleManager: ui.NewStyleManager(),
        currentTheme: cfg.Theme,
    }
}

// ... rest of model implementation
```

#### User Documentation for Custom Themes

**Creating Your Own Theme:**

1. Create or edit `~/.noise.sh/config.json`
2. Add a `custom_themes` section with your theme definition
3. Use valid hex colors (format: `#RRGGBB`)
4. Give your theme a unique ID (lowercase, use hyphens)
5. Restart noise.sh or use `Ctrl+R` to reload config
6. Select your custom theme with `Ctrl+T`

**Required Colors:**
- `name`: Display name for your theme
- `primary`: Main UI color (buttons, highlights)
- `secondary`: Secondary UI elements (borders, inactive)
- `accent`: Important highlights (warnings, special text)
- `background`: Main background color
- `text`: Main text color
- `success`: Success messages and positive feedback

**Color Guidelines:**
- Ensure text is readable on background (contrast ratio ≥ 4.5:1)
- Primary/Secondary should contrast with background (≥ 3:1)
- Test your theme in the app before relying on it
- Use online contrast checkers if unsure

**Example Custom Theme:**
```json
{
  "theme": "midnight-purple",
  "custom_themes": {
    "midnight-purple": {
      "name": "Midnight Purple",
      "primary": "#9D84B7",
      "secondary": "#5E4B8B",
      "accent": "#F4D03F",
      "background": "#0A0E27",
      "text": "#E8DFF5",
      "success": "#52D3AA"
    }
  }
}
```

#### Command to Generate Theme Template

Add a CLI command to generate a custom theme template:

```bash
noise.sh theme new my-theme-name
```

This creates a config section with placeholder colors that users can edit.

---

### Auto Theme Detection
Detect system dark/light mode and suggest appropriate theme:
- Light mode → Slate Mist, Forest Whisper
- Dark mode → Amber Night, Iron Storm

### Theme Preview Mode
Before applying a theme, show a preview with sample text:
```
┌─ Amber Night Preview ─────────┐
│ Title Style                    │
│ Primary: ████ Secondary: ████  │
│ Accent:  ████ Success:   ████  │
│                                │
│ Sample lyric text goes here... │
│ With multiple lines to show    │
│ readability and contrast       │
└────────────────────────────────┘
```

---

## Notes for Kilocode AI

### Implementation Tips
1. **Implement all 12 themes at once** - They're already defined and tested
2. **Use theme.GetManager()** singleton pattern for accessing current theme
3. **Load custom themes on startup** from config file before initializing UI
4. **All UI components should use StyleManager** instead of hardcoded colors
5. **Test contrast ratios** - All themes have been verified but double-check in actual TUI
6. **Save theme preference immediately** when user changes it (don't wait for app exit)
7. **Validate custom themes** - Check hex colors and provide helpful error messages

### Common Pitfalls to Avoid
- ❌ Don't hardcode colors anywhere in UI code
- ❌ Don't forget to re-render UI after theme change
- ❌ Don't use theme colors directly - always go through theme manager
- ❌ Don't forget thread-safety (use mutex in theme manager)
- ❌ Don't skip custom theme validation - invalid hex codes will crash the app
- ❌ Don't forget to load custom themes before setting the active theme

### Testing Commands
```bash
# Test theme switching with all themes
go run cmd/noise.sh/main.go --theme slate
go run cmd/noise.sh/main.go --theme amber-night
go run cmd/noise.sh/main.go --theme molten-gold
go run cmd/noise.sh/main.go --theme clay-roads
go run cmd/noise.sh/main.go --theme iron-storm
go run cmd/noise.sh/main.go --theme jade-tide
go run cmd/noise.sh/main.go --theme sunset-ember
go run cmd/noise.sh/main.go --theme forest-whisper
go run cmd/noise.sh/main.go --theme electric-bloom
go run cmd/noise.sh/main.go --theme plasma-pulse
go run cmd/noise.sh/main.go --theme indigo-depths
go run cmd/noise.sh/main.go --theme sage-meadow

# Test with invalid theme (should fall back to default)
go run cmd/noise.sh/main.go --theme invalid-theme

# Test custom theme loading
# 1. Create ~/.noise.sh/config.json with custom theme
# 2. Run app and verify custom theme appears in theme list
# 3. Select custom theme with Ctrl+T

# Visual contrast test
# Use Ctrl+T to cycle through all themes in running app
# Verify each theme is readable and colors are distinct
```

---

## Success Criteria

This implementation will be considered successful when:

1. ✅ All 12 themes are implemented and selectable
2. ✅ User can switch themes at runtime with Ctrl+T
3. ✅ Theme preference persists across sessions
4. ✅ All themes pass WCAG AA contrast requirements
5. ✅ UI consistently uses theme colors (no hardcoded colors)
6. ✅ Theme switching is smooth and instant
7. ✅ Default theme (Amber Night) loads correctly
8. ✅ Invalid theme names fall back gracefully
9. ✅ Config file saves/loads correctly
10. ✅ Documentation is complete and accurate

---

## Support & Questions

For questions during implementation:
1. Verify all color hex codes match exactly (case-insensitive but # required)
2. Ensure lipgloss.Color() wraps all color values
3. Check that theme manager is initialized before any UI rendering
4. Confirm config directory (~/.noise.sh/) is created before saving
5. Test each theme visually - numbers can lie, eyes don't

**Good luck with implementation! 🎨🚀**

---

**Document Version:** 1.0  
**Last Updated:** October 17, 2025  
**Status:** Ready for Implementation
