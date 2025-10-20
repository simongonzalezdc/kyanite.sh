# Kyanite Suite Implementation Master Document
**Complete Implementation Guide for Coding Agents**

**Version:** 1.0  
**Date:** October 2025  
**Target:** focus.sh and noise.sh (apply to all future tools)  
**Status:** READY FOR IMPLEMENTATION

---

## Table of Contents

1. [Overview](#overview)
2. [Theme System (10 Themes)](#theme-system-10-themes)
3. [Complete Code Implementation](#complete-code-implementation)
4. [Suite Standards](#suite-standards)
5. [File Structure](#file-structure)
6. [Migration Logic](#migration-logic)
7. [Testing Checklist](#testing-checklist)
8. [Implementation Steps](#implementation-steps)

---

## Overview

### What Changed

**Themes:** 13 → 10 themes (consolidated, no redundancy)  
**Suite Standards:** Unified UI/UX, keyboard shortcuts, architecture

### What This Document Contains

- ✅ Complete 10-theme system with all hex codes
- ✅ Copy-paste ready Go code (registry.go, manager.go, validator.go)
- ✅ Universal keyboard shortcuts
- ✅ UI/UX standards for all tools
- ✅ Migration logic for old theme IDs
- ✅ Testing checklist
- ✅ Step-by-step implementation guide

---

## Theme System (10 Themes)

### Theme Consolidation Summary

**REMOVED/MERGED:**
- Slate Mist + Violet Dusk → **Twilight Mist** (cool purple/gray)
- Molten Gold + Sunset Ember → **Sunlight** (golden yellow)
- Forest Whisper + Sage Meadow → **Forest Path** (sage green)
- Electric Bloom + Sunset Ember → **Electric Rose** (hot pink)
- Plasma Pulse → REMOVED (accessibility fail)
- Jade Tide → Replaced by **Cyan Wave** (clearer cyan)
- Clay Roads → **Clay Earth** (true brown)
- Iron Storm → **Iron Forge** (renamed)

**ADDED:**
- **Monochrome** (pure black/white for accessibility)
- **Sunlight** (golden yellow for daytime)
- **Cyan Wave** (technical cyan)
- **Twilight Mist** (cool purple/gray)

### The 10 Final Themes

| # | Theme ID | Name | Primary | Default? |
|---|----------|------|---------|----------|
| 1 | monochrome | Monochrome | #FFFFFF | No |
| 2 | amber-night | Amber Night | #D4A574 | ✅ YES |
| 3 | twilight-mist | Twilight Mist | #B8A3C9 | No |
| 4 | indigo-depths | Indigo Depths | #4169E1 | No |
| 5 | forest-path | Forest Path | #8FBC8F | No |
| 6 | clay-earth | Clay Earth | #A0522D | No |
| 7 | iron-forge | Iron Forge | #DC143C | No |
| 8 | sunlight | Sunlight | #FFD700 | No |
| 9 | cyan-wave | Cyan Wave | #00CED1 | No |
| 10 | electric-rose | Electric Rose | #FF1493 | No |

---

## Complete Code Implementation

### File Structure

```
internal/theme/
├── types.go          # Theme struct definition
├── registry.go       # All 10 theme definitions
├── manager.go        # Runtime switching logic
├── custom.go         # Custom theme loader
└── validator.go      # WCAG contrast checker
```

---

### 1. types.go

```go
package theme

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the application
type Theme struct {
    Name       string
    Primary    lipgloss.Color  // Main UI elements
    Secondary  lipgloss.Color  // Supporting elements
    Accent     lipgloss.Color  // Highlights, focus state
    Background lipgloss.Color  // Base background
    Text       lipgloss.Color  // Primary text
    Success    lipgloss.Color  // Success states
    Warning    lipgloss.Color  // Warning states
    Error      lipgloss.Color  // Error states
}

// Manager handles theme selection and switching
type Manager struct {
    current Theme
}
```

---

### 2. registry.go (COMPLETE - Copy/Paste Ready)

```go
package theme

import "github.com/charmbracelet/lipgloss"

// All 10 Kyanite Suite themes
var (
    Monochrome = Theme{
        Name:       "Monochrome",
        Primary:    lipgloss.Color("#FFFFFF"),
        Secondary:  lipgloss.Color("#999999"),
        Accent:     lipgloss.Color("#FFFFFF"),
        Background: lipgloss.Color("#000000"),
        Text:       lipgloss.Color("#FFFFFF"),
        Success:    lipgloss.Color("#CCCCCC"),
        Warning:    lipgloss.Color("#888888"),
        Error:      lipgloss.Color("#666666"),
    }

    AmberNight = Theme{
        Name:       "Amber Night",
        Primary:    lipgloss.Color("#D4A574"),
        Secondary:  lipgloss.Color("#9D84B7"),
        Accent:     lipgloss.Color("#F4D03F"),
        Background: lipgloss.Color("#0A0E27"),
        Text:       lipgloss.Color("#E8DFF5"),
        Success:    lipgloss.Color("#52D3AA"),
        Warning:    lipgloss.Color("#FFA502"),
        Error:      lipgloss.Color("#EA2027"),
    }

    TwilightMist = Theme{
        Name:       "Twilight Mist",
        Primary:    lipgloss.Color("#B8A3C9"),
        Secondary:  lipgloss.Color("#8E7B9D"),
        Accent:     lipgloss.Color("#D4C5E0"),
        Background: lipgloss.Color("#151520"),
        Text:       lipgloss.Color("#E8E4F0"),
        Success:    lipgloss.Color("#90C695"),
        Warning:    lipgloss.Color("#C9A87C"),
        Error:      lipgloss.Color("#C77777"),
    }

    IndigoDepths = Theme{
        Name:       "Indigo Depths",
        Primary:    lipgloss.Color("#4169E1"),
        Secondary:  lipgloss.Color("#5F9EA0"),
        Accent:     lipgloss.Color("#87CEEB"),
        Background: lipgloss.Color("#0A0A1A"),
        Text:       lipgloss.Color("#E6F2FF"),
        Success:    lipgloss.Color("#52D3AA"),
        Warning:    lipgloss.Color("#FFB84D"),
        Error:      lipgloss.Color("#FF6B6B"),
    }

    ForestPath = Theme{
        Name:       "Forest Path",
        Primary:    lipgloss.Color("#8FBC8F"),
        Secondary:  lipgloss.Color("#6B8E6B"),
        Accent:     lipgloss.Color("#B4D7B4"),
        Background: lipgloss.Color("#1A1F1A"),
        Text:       lipgloss.Color("#E8F5E8"),
        Success:    lipgloss.Color("#90EE90"),
        Warning:    lipgloss.Color("#DAA520"),
        Error:      lipgloss.Color("#CD5C5C"),
    }

    ClayEarth = Theme{
        Name:       "Clay Earth",
        Primary:    lipgloss.Color("#A0522D"),
        Secondary:  lipgloss.Color("#8B4513"),
        Accent:     lipgloss.Color("#DEB887"),
        Background: lipgloss.Color("#1A1410"),
        Text:       lipgloss.Color("#F5E6D3"),
        Success:    lipgloss.Color("#8FBC8F"),
        Warning:    lipgloss.Color("#CD853F"),
        Error:      lipgloss.Color("#CD5C5C"),
    }

    IronForge = Theme{
        Name:       "Iron Forge",
        Primary:    lipgloss.Color("#DC143C"),
        Secondary:  lipgloss.Color("#4A4A4A"),
        Accent:     lipgloss.Color("#FF6347"),
        Background: lipgloss.Color("#1A0A0A"),
        Text:       lipgloss.Color("#FFE6E6"),
        Success:    lipgloss.Color("#90C695"),
        Warning:    lipgloss.Color("#FFB84D"),
        Error:      lipgloss.Color("#FF4444"),
    }

    Sunlight = Theme{
        Name:       "Sunlight",
        Primary:    lipgloss.Color("#FFD700"),
        Secondary:  lipgloss.Color("#DAA520"),
        Accent:     lipgloss.Color("#FFF8DC"),
        Background: lipgloss.Color("#0F0F0A"),
        Text:       lipgloss.Color("#FFFACD"),
        Success:    lipgloss.Color("#98D982"),
        Warning:    lipgloss.Color("#FF9800"),
        Error:      lipgloss.Color("#D32F2F"),
    }

    CyanWave = Theme{
        Name:       "Cyan Wave",
        Primary:    lipgloss.Color("#00CED1"),
        Secondary:  lipgloss.Color("#4682B4"),
        Accent:     lipgloss.Color("#7FFFD4"),
        Background: lipgloss.Color("#0A1418"),
        Text:       lipgloss.Color("#E0F7FA"),
        Success:    lipgloss.Color("#52D3AA"),
        Warning:    lipgloss.Color("#FFB84D"),
        Error:      lipgloss.Color("#FF6B6B"),
    }

    ElectricRose = Theme{
        Name:       "Electric Rose",
        Primary:    lipgloss.Color("#FF1493"),
        Secondary:  lipgloss.Color("#C71585"),
        Accent:     lipgloss.Color("#00CED1"),
        Background: lipgloss.Color("#1A0A1A"),
        Text:       lipgloss.Color("#FFF0F8"),
        Success:    lipgloss.Color("#52D3AA"),
        Warning:    lipgloss.Color("#FFB84D"),
        Error:      lipgloss.Color("#FF4444"),
    }
)

// Registry maps theme IDs to Theme structs
var Registry = map[string]Theme{
    "monochrome":    Monochrome,
    "amber-night":   AmberNight,
    "twilight-mist": TwilightMist,
    "indigo-depths": IndigoDepths,
    "forest-path":   ForestPath,
    "clay-earth":    ClayEarth,
    "iron-forge":    IronForge,
    "sunlight":      Sunlight,
    "cyan-wave":     CyanWave,
    "electric-rose": ElectricRose,
}

// Default returns the default theme (Amber Night)
func Default() Theme {
    return AmberNight
}

// GetTheme returns a theme by ID, falling back to default
func GetTheme(id string) Theme {
    // Check if migration needed
    id = migrateThemeID(id)
    
    if theme, ok := Registry[id]; ok {
        return theme
    }
    return Default()
}

// ListThemes returns all theme IDs in order
func ListThemes() []string {
    return []string{
        "monochrome",
        "amber-night",
        "twilight-mist",
        "indigo-depths",
        "forest-path",
        "clay-earth",
        "iron-forge",
        "sunlight",
        "cyan-wave",
        "electric-rose",
    }
}

// migrateThemeID migrates old theme IDs to new ones
func migrateThemeID(oldID string) string {
    migrations := map[string]string{
        "slate-mist":     "twilight-mist",
        "violet-dusk":    "twilight-mist",
        "molten-gold":    "sunlight",
        "clay-roads":     "clay-earth",
        "iron-storm":     "iron-forge",
        "jade-tide":      "cyan-wave",
        "sunset-ember":   "electric-rose",
        "forest-whisper": "forest-path",
        "electric-bloom": "electric-rose",
        "plasma-pulse":   "amber-night",  // Fallback to default
        "sage-meadow":    "forest-path",
    }
    
    if newID, ok := migrations[oldID]; ok {
        return newID
    }
    return oldID
}
```

---

### 3. manager.go

```go
package theme

import (
    "sync"
)

var (
    globalManager *Manager
    once          sync.Once
)

// Manager handles theme selection and switching
type Manager struct {
    mu      sync.RWMutex
    current Theme
}

// GetManager returns the global theme manager (singleton)
func GetManager() *Manager {
    once.Do(func() {
        globalManager = &Manager{
            current: Default(),
        }
    })
    return globalManager
}

// SetTheme sets the current theme by ID
func (m *Manager) SetTheme(id string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.current = GetTheme(id)
}

// Current returns the current theme
func (m *Manager) Current() Theme {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.current
}

// Next cycles to the next theme in the registry
func (m *Manager) Next() Theme {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    themes := ListThemes()
    currentID := ""
    
    // Find current theme ID
    for id, theme := range Registry {
        if theme.Name == m.current.Name {
            currentID = id
            break
        }
    }
    
    // Find next theme
    for i, id := range themes {
        if id == currentID {
            nextIndex := (i + 1) % len(themes)
            m.current = Registry[themes[nextIndex]]
            break
        }
    }
    
    return m.current
}

// Previous cycles to the previous theme in the registry
func (m *Manager) Previous() Theme {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    themes := ListThemes()
    currentID := ""
    
    // Find current theme ID
    for id, theme := range Registry {
        if theme.Name == m.current.Name {
            currentID = id
            break
        }
    }
    
    // Find previous theme
    for i, id := range themes {
        if id == currentID {
            prevIndex := (i - 1 + len(themes)) % len(themes)
            m.current = Registry[themes[prevIndex]]
            break
        }
    }
    
    return m.current
}
```

---

### 4. custom.go

```go
package theme

import (
    "os"
    "path/filepath"
    
    "github.com/BurntSushi/toml"
    "github.com/charmbracelet/lipgloss"
)

// CustomTheme represents a user-defined theme from TOML
type CustomTheme struct {
    Name   string       `toml:"name"`
    Colors ThemeColors  `toml:"colors"`
}

// ThemeColors holds hex color values
type ThemeColors struct {
    Primary    string `toml:"primary"`
    Secondary  string `toml:"secondary"`
    Accent     string `toml:"accent"`
    Background string `toml:"background"`
    Text       string `toml:"text"`
    Success    string `toml:"success"`
    Warning    string `toml:"warning"`
    Error      string `toml:"error"`
}

// LoadCustomThemes loads custom themes from ~/.config/[tool]/themes/
func LoadCustomThemes(toolName string) (map[string]Theme, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    
    themesDir := filepath.Join(homeDir, ".config", toolName, "themes")
    
    // Check if themes directory exists
    if _, err := os.Stat(themesDir); os.IsNotExist(err) {
        return map[string]Theme{}, nil
    }
    
    customThemes := make(map[string]Theme)
    
    // Read all .toml files in themes directory
    entries, err := os.ReadDir(themesDir)
    if err != nil {
        return nil, err
    }
    
    for _, entry := range entries {
        if entry.IsDir() || filepath.Ext(entry.Name()) != ".toml" {
            continue
        }
        
        var ct CustomTheme
        filePath := filepath.Join(themesDir, entry.Name())
        
        if _, err := toml.DecodeFile(filePath, &ct); err != nil {
            continue // Skip invalid files
        }
        
        // Convert to Theme
        theme := Theme{
            Name:       ct.Name,
            Primary:    lipgloss.Color(ct.Colors.Primary),
            Secondary:  lipgloss.Color(ct.Colors.Secondary),
            Accent:     lipgloss.Color(ct.Colors.Accent),
            Background: lipgloss.Color(ct.Colors.Background),
            Text:       lipgloss.Color(ct.Colors.Text),
            Success:    lipgloss.Color(ct.Colors.Success),
            Warning:    lipgloss.Color(ct.Colors.Warning),
            Error:      lipgloss.Color(ct.Colors.Error),
        }
        
        // Use filename without extension as ID
        themeID := entry.Name()[:len(entry.Name())-5]
        customThemes[themeID] = theme
    }
    
    return customThemes, nil
}
```

---

### 5. validator.go

```go
package theme

import (
    "fmt"
    "image/color"
    "math"
    "strconv"
    "strings"
)

// ContrastRatio calculates the WCAG contrast ratio between two colors
func ContrastRatio(c1, c2 color.Color) float64 {
    l1 := relativeLuminance(c1)
    l2 := relativeLuminance(c2)
    
    if l1 > l2 {
        return (l1 + 0.05) / (l2 + 0.05)
    }
    return (l2 + 0.05) / (l1 + 0.05)
}

// relativeLuminance calculates the relative luminance of a color
func relativeLuminance(c color.Color) float64 {
    r, g, b, _ := c.RGBA()
    
    // Convert to 0-1 range
    rs := float64(r) / 65535.0
    gs := float64(g) / 65535.0
    bs := float64(b) / 65535.0
    
    // Apply gamma correction
    r8 := gammaCorrect(rs)
    g8 := gammaCorrect(gs)
    b8 := gammaCorrect(bs)
    
    // Calculate luminance
    return 0.2126*r8 + 0.7152*g8 + 0.0722*b8
}

func gammaCorrect(v float64) float64 {
    if v <= 0.03928 {
        return v / 12.92
    }
    return math.Pow((v+0.055)/1.055, 2.4)
}

// ValidateTheme checks if a theme meets WCAG AA standards
func ValidateTheme(t Theme) []string {
    warnings := []string{}
    
    // Parse colors
    bgColor := parseHexColor(string(t.Background))
    textColor := parseHexColor(string(t.Text))
    primaryColor := parseHexColor(string(t.Primary))
    
    // Check text contrast (should be 4.5:1 for AA)
    textContrast := ContrastRatio(textColor, bgColor)
    if textContrast < 4.5 {
        warnings = append(warnings, 
            fmt.Sprintf("Text contrast too low: %.2f:1 (need 4.5:1)", textContrast))
    }
    
    // Check UI element contrast (should be 3:1 for AA)
    uiContrast := ContrastRatio(primaryColor, bgColor)
    if uiContrast < 3.0 {
        warnings = append(warnings, 
            fmt.Sprintf("UI contrast too low: %.2f:1 (need 3:1)", uiContrast))
    }
    
    return warnings
}

// parseHexColor converts hex string to color.Color
func parseHexColor(hex string) color.Color {
    hex = strings.TrimPrefix(hex, "#")
    
    if len(hex) != 6 {
        return color.White
    }
    
    r, _ := strconv.ParseUint(hex[0:2], 16, 8)
    g, _ := strconv.ParseUint(hex[2:4], 16, 8)
    b, _ := strconv.ParseUint(hex[4:6], 16, 8)
    
    return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
}
```

---

## Suite Standards

### Universal Keyboard Shortcuts

**MUST implement in ALL Kyanite tools:**

```go
const (
    // Universal shortcuts (all tools)
    KeyQuit          = "ctrl+q"
    KeyHelp          = "ctrl+h"
    KeyTheme         = "ctrl+shift+t"
    KeySave          = "ctrl+s"
    KeyUndo          = "ctrl+z"
    KeyRedo          = "ctrl+y"
    KeyCommandPalette = "ctrl+/"
    
    // AI shortcuts (if tool has AI)
    KeyAIGenerate = "ctrl+g"
    KeyAISpark    = "ctrl+r"
    KeyAITweak    = "ctrl+v"
    KeyAICheck    = "ctrl+shift+c"
)
```

### Theme Switching Implementation

**Add to your main model's Update() function:**

```go
case key.Matches(msg, keys.ThemeSwitch):  // Ctrl+Shift+T
    theme.GetManager().Next()
    
    // Save to config
    cfg := config.Load()
    cfg.Theme = getCurrentThemeID()
    cfg.Save()
    
    // Trigger UI re-render
    return m, nil
```

### Using Themes in Components

**Pattern for all UI components:**

```go
import "your-tool/internal/theme"

func (m model) View() string {
    t := theme.GetManager().Current()
    
    titleStyle := lipgloss.NewStyle().
        Foreground(t.Primary).
        Background(t.Background).
        Bold(true)
    
    return titleStyle.Render("Hello World")
}
```

### Configuration File Format

**Location:** `~/.config/[tool]/config.toml`

**Required theme section:**

```toml
[preferences]
theme = "amber-night"  # Must be one of 10 theme IDs

[theme_settings]
custom_themes_enabled = true
```

### Auto-Save Pattern

**MUST implement in all tools:**

```go
// Auto-save every 30 seconds
func (m model) autoSave() tea.Cmd {
    return tea.Tick(30*time.Second, func(t time.Time) tea.Msg {
        return autoSaveMsg{}
    })
}

// Show subtle notification
func showSaveNotification() string {
    style := lipgloss.NewStyle().
        Foreground(theme.GetManager().Current().Success).
        Faint(true)
    
    return style.Render("✓ Saved")
}
```

---

## File Structure

### Directory Structure (Standard)

```
[tool]/
├── cmd/
│   ├── root.go
│   ├── version.go
│   └── [feature].go
├── internal/
│   ├── app/
│   │   └── ai/              # If tool has AI
│   ├── ui/
│   │   ├── editor/          # If tool has editor
│   │   ├── styles/
│   │   └── screens/
│   ├── theme/              # ← NEW (or update existing)
│   │   ├── types.go
│   │   ├── registry.go
│   │   ├── manager.go
│   │   ├── custom.go
│   │   └── validator.go
│   ├── storage/
│   ├── config/
│   └── export/
├── data/
├── docs/
│   ├── ROADMAP.md
│   ├── ARCHITECTURE.md
│   └── archive/
├── go.mod
├── go.sum
└── README.md
```

---

## Migration Logic

### Config Migration Function

**Add to your config loading logic:**

```go
package config

import (
    "your-tool/internal/theme"
)

func Load() *Config {
    cfg := loadFromFile()
    
    // Migrate old theme IDs to new
    if needsMigration(cfg.Theme) {
        oldTheme := cfg.Theme
        cfg.Theme = theme.GetTheme(oldTheme).Name
        
        // Convert name back to ID
        cfg.Theme = getThemeID(cfg.Theme)
        
        // Save migrated config
        cfg.Save()
    }
    
    return cfg
}

func getThemeID(themeName string) string {
    for id, t := range theme.Registry {
        if t.Name == themeName {
            return id
        }
    }
    return "amber-night"
}
```

### Old to New Theme ID Mapping

**Already handled in registry.go's migrateThemeID() function:**

```
slate-mist     → twilight-mist
violet-dusk    → twilight-mist
molten-gold    → sunlight
clay-roads     → clay-earth
iron-storm     → iron-forge
jade-tide      → cyan-wave
sunset-ember   → electric-rose
forest-whisper → forest-path
electric-bloom → electric-rose
plasma-pulse   → amber-night (fallback)
sage-meadow    → forest-path
```

---

## Testing Checklist

### Visual Testing

```bash
# Run these tests after implementation

- [ ] All 10 themes render correctly
- [ ] Text readable in all themes
- [ ] UI elements properly colored
- [ ] No color bleeding or artifacts
- [ ] Borders use theme colors
- [ ] Focus states visible in all themes
```

### Functional Testing

```bash
- [ ] Ctrl+Shift+T cycles through themes
- [ ] Theme preference saves to config
- [ ] Theme preference loads on startup
- [ ] Invalid theme IDs fall back to default
- [ ] Old theme IDs auto-migrate
- [ ] Custom themes load (if implemented)
```

### Accessibility Testing

```bash
- [ ] Run contrast validator on all themes
- [ ] Test in grayscale mode (terminal settings)
- [ ] Test with colorblind simulators
- [ ] Verify Monochrome theme works perfectly
- [ ] Check WCAG AA compliance (4.5:1 text)
```

### Performance Testing

```bash
- [ ] Theme switching is instant (<50ms)
- [ ] No memory leaks from theme changes
- [ ] UI re-renders efficiently
- [ ] No flicker during theme switch
```

---

## Implementation Steps

### Phase 1: Create Theme System (Day 1)

**1. Create directory structure:**

```bash
mkdir -p internal/theme
```

**2. Copy files:**
- Copy `types.go` from this document → `internal/theme/types.go`
- Copy `registry.go` from this document → `internal/theme/registry.go`
- Copy `manager.go` from this document → `internal/theme/manager.go`
- Copy `custom.go` from this document → `internal/theme/custom.go`
- Copy `validator.go` from this document → `internal/theme/validator.go`

**3. Update imports:**

```bash
# Update go.mod if needed
go get github.com/BurntSushi/toml
go mod tidy
```

**4. Test compilation:**

```bash
go build ./internal/theme
```

---

### Phase 2: Update Existing Code (Day 2)

**1. Find all theme references:**

```bash
# Search for old theme usage
grep -r "lipgloss.Color" internal/ui/
grep -r "theme" internal/config/
```

**2. Update UI components:**

**Before:**
```go
titleStyle := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#9D84B7"))
```

**After:**
```go
import "your-tool/internal/theme"

t := theme.GetManager().Current()
titleStyle := lipgloss.NewStyle().
    Foreground(t.Primary)
```

**3. Update config loading:**

Add theme migration logic (see Migration Logic section)

**4. Add theme switching keybind:**

```go
case key.Matches(msg, keys.ThemeSwitch):
    theme.GetManager().Next()
    // Save theme preference
    return m, nil
```

---

### Phase 3: Update Configuration (Day 3)

**1. Update config.toml format:**

```toml
[preferences]
theme = "amber-night"

[theme_settings]
custom_themes_enabled = true
```

**2. Add theme loading on startup:**

```go
func main() {
    cfg := config.Load()
    theme.GetManager().SetTheme(cfg.Theme)
    // ... rest of app initialization
}
```

**3. Test theme persistence:**

```bash
# Change theme with Ctrl+Shift+T
# Restart app
# Verify theme persists
```

---

### Phase 4: Testing & Validation (Day 4)

**1. Run visual tests:**

```bash
# Start app
# Press Ctrl+Shift+T 10 times (cycle all themes)
# Verify each theme looks correct
```

**2. Run migration test:**

```bash
# Edit config.toml, set theme = "plasma-pulse"
# Restart app
# Verify it migrates to "amber-night"
```

**3. Run contrast validator:**

```go
// Add to main_test.go
func TestThemeContrast(t *testing.T) {
    for id, theme := range theme.Registry {
        warnings := theme.ValidateTheme(theme)
        if len(warnings) > 0 {
            t.Logf("Theme %s warnings: %v", id, warnings)
        }
    }
}
```

**4. Test edge cases:**

```bash
- [ ] No config file exists (should use default)
- [ ] Invalid theme ID in config (should fall back)
- [ ] Config file corrupted (should use default)
- [ ] Theme switching during active use
```

---

### Phase 5: Documentation (Day 5)

**1. Update README.md:**

```markdown
## Themes

[Tool] supports 10 professionally designed themes. Switch themes with `Ctrl+Shift+T`.

Available themes:
- Monochrome (accessibility)
- Amber Night (default)
- Twilight Mist
- Indigo Depths
- Forest Path
- Clay Earth
- Iron Forge
- Sunlight
- Cyan Wave
- Electric Rose

Configure default theme in `~/.config/[tool]/config.toml`
```

**2. Update help text:**

Add to `--help` output:

```
Keyboard Shortcuts:
  Ctrl+Shift+T    Switch theme
  Ctrl+H          Show help
  Ctrl+Q          Quit
```

**3. Create THEMES.md:**

Document all 10 themes with screenshots/descriptions

---

## Verification Checklist

**Before marking complete:**

```bash
✅ All 10 themes defined in registry.go
✅ Theme manager implemented
✅ Universal shortcuts added (Ctrl+Shift+T)
✅ Config migration logic added
✅ All UI components use theme.GetManager().Current()
✅ Theme preference persists across restarts
✅ Old theme IDs auto-migrate
✅ All tests pass
✅ Documentation updated
✅ README mentions themes
✅ Help text shows theme shortcut
```

---

## Troubleshooting

### Theme not persisting

**Check:**
- Config file location: `~/.config/[tool]/config.toml`
- Config save function called after theme change
- File permissions on config directory

### Theme colors not showing

**Check:**
- All UI components import `internal/theme`
- Using `theme.GetManager().Current()` not hardcoded colors
- Re-rendering UI after theme change

### Old themes still showing

**Check:**
- Migration logic in `registry.go` includes all old IDs
- Config loading calls migration function
- No cached old theme references

---

## Summary

### What You're Implementing

1. ✅ **10-theme system** (complete code provided)
2. ✅ **Theme manager** (runtime switching)
3. ✅ **Universal shortcuts** (Ctrl+Shift+T)
4. ✅ **Config migration** (old → new theme IDs)
5. ✅ **Suite standards** (UI patterns, file structure)

### Estimated Time

- **Day 1:** Copy theme files, test compilation
- **Day 2:** Update UI components to use themes
- **Day 3:** Update config, add persistence
- **Day 4:** Testing, validation, bug fixes
- **Day 5:** Documentation, polish

**Total:** ~5 days per tool (focus.sh, noise.sh)

### Files to Create/Modify

**NEW FILES:**
- `internal/theme/types.go`
- `internal/theme/registry.go`
- `internal/theme/manager.go`
- `internal/theme/custom.go`
- `internal/theme/validator.go`

**MODIFY:**
- All files in `internal/ui/` (use theme manager)
- `internal/config/` (add migration logic)
- Main entry point (initialize theme on startup)
- Help text (add Ctrl+Shift+T)
- README.md (document themes)

---

## Questions?

**If something is unclear:**
1. Re-read the relevant section
2. Check the code examples
3. Refer to the Testing Checklist
4. Look at troubleshooting section

**This document is complete.** Everything needed for implementation is included.

---

**Status:** ✅ READY FOR IMPLEMENTATION  
**Owner:** Coding Agent  
**Target Tools:** focus.sh, noise.sh (then all future tools)  
**Timeline:** 5 days per tool  

**Next Step:** Copy theme files to `internal/theme/` and begin Phase 1.