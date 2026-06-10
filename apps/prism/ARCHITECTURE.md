# Prism.sh Architecture

**Version:** 1.0
**Date:** November 2025

## Overview

Prism.sh is a terminal-based color palette design tool built with Go, Bubble Tea (TUI framework), and Lipgloss (styling). This document describes the architecture, design decisions, and implementation details.

## Technology Stack

- **Language:** Go 1.21+
- **TUI Framework:** Bubble Tea (v1.3.10)
- **Styling:** Lipgloss (v1.1.0)
- **Architecture Pattern:** Elm Architecture (Model-Update-View)

## Project Structure

```
prism/
├── cmd/prism/              # Application entry point
│   └── main.go             # Initializes Bubble Tea program
├── internal/               # Private application code
│   ├── app/                # Root application model
│   │   └── model.go        # Main Bubble Tea model & routing
│   ├── ui/                 # User interface screens
│   │   ├── menu.go         # Main menu
│   │   ├── wheel.go        # Color wheel
│   │   ├── generator.go    # Palette generator
│   │   ├── theory.go       # Color theory lessons
│   │   ├── checker.go      # WCAG checker
│   │   ├── manager.go      # Palette manager
│   │   ├── help.go         # Help overlay
│   │   ├── styles.go       # Lipgloss styles
│   │   └── types.go        # Shared types
│   ├── color/              # Color mathematics
│   │   ├── types.go        # Color struct & types
│   │   ├── convert.go      # Color space conversions
│   │   └── names.go        # Named colors database
│   ├── palette/            # Palette generation
│   │   ├── types.go        # Palette struct & rules
│   │   └── generator.go    # Harmony algorithms
│   ├── wcag/               # Accessibility
│   │   └── contrast.go     # WCAG 2.1 calculations
│   ├── theme/              # UI theming
│   │   ├── types.go        # Theme struct
│   │   ├── registry.go     # 10 Kyanite themes
│   │   └── manager.go      # Theme switching
│   ├── storage/            # File persistence
│   │   ├── config.go       # Configuration
│   │   ├── palettes.go     # Palette save/load
│   │   └── lock.go         # File locking
│   └── export/             # Export formats
│       ├── json.go         # JSON export
│       ├── css.go          # CSS variables
│       ├── toml.go         # TOML format
│       └── theme.go        # Kyanite theme format
├── tests/                  # Test suite
├── data/                   # Static data
│   └── colors.json         # Named colors database
└── .github/workflows/      # CI/CD pipelines
```

## Design Patterns

### 1. Elm Architecture (Bubble Tea)

Prism.sh follows the Elm Architecture pattern enforced by Bubble Tea:

```
┌─────────────────────────────────────┐
│            User Input               │
│         (keyboard, mouse)           │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      Update(msg) -> (Model, Cmd)    │
│   (Pure function, no side effects)  │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│        View() -> string             │
│    (Renders current model state)    │
└─────────────────────────────────────┘
```

**Benefits:**
- Predictable state management
- Easy testing (pure functions)
- No race conditions
- Clear data flow

### 2. Screen-Based Navigation

Each screen is an independent Bubble Tea model:

```go
type Model struct {
    CurrentScreen Screen
    menuModel     ui.MenuModel
    wheelModel    ui.WheelModel
    generatorModel ui.GeneratorModel
    // ... other screens
}
```

Navigation uses message passing:

```go
type NavigateMsg struct {
    Screen int
}
```

### 3. Theme System

Centralized theme management with 10 pre-defined themes:

```go
type Theme struct {
    Name       string
    Primary    string
    Secondary  string
    Accent     string
    Background string
    Text       string
    Success    string
}
```

All UI components reference the current theme, enabling runtime theme switching.

## Module Breakdown

### Color Module (`internal/color/`)

**Responsibilities:**
- Color space conversions (RGB ↔ HSL ↔ HSV ↔ Hex)
- Color operations (lighten, darken, saturate, etc.)
- Named color database and search
- Temperature detection (warm/cool)

**Key Algorithms:**
- **RGB to HSL:** Uses standard HSL conversion formula with proper handling of edge cases
- **HSL to RGB:** Implements piecewise function based on hue sector
- **Fuzzy Search:** Levenshtein distance algorithm for 1-2 character typo tolerance

### Palette Module (`internal/palette/`)

**Responsibilities:**
- Generate harmonious color palettes
- Implement 7 harmony rules
- Validate palette contrast

**Harmony Rules:**
1. **Monochromatic:** Same hue, varying lightness (±20%, ±40%)
2. **Complementary:** Base + opposite (180°)
3. **Analogous:** Base ± 30°
4. **Triadic:** 120° spacing
5. **Tetradic:** 90° spacing (two complementary pairs)
6. **Split-Complementary:** 150°, 210°
7. **Square:** Same as tetradic

### WCAG Module (`internal/wcag/`)

**Responsibilities:**
- Calculate contrast ratios
- Validate WCAG 2.1 compliance
- Provide AA/AAA level determination

**Formula (WCAG 2.1):**
```
1. Linearize RGB: if c <= 0.03928 then c/12.92 else ((c+0.055)/1.055)^2.4
2. Luminance: 0.2126*R + 0.7152*G + 0.0722*B
3. Contrast: (L1 + 0.05) / (L2 + 0.05)  [L1 is lighter]
```

**Thresholds:**
- AA Small Text: ≥ 4.5:1
- AA Large Text: ≥ 3:1
- AAA Small Text: ≥ 7:1
- AAA Large Text: ≥ 4.5:1

### Storage Module (`internal/storage/`)

**Responsibilities:**
- Cross-platform file paths
- Atomic file writes
- File locking for concurrent access
- Configuration management

**File Structure:**
```
{ConfigDir}/prism/
├── config.toml          # User preferences
├── palettes/
│   ├── {id}.json        # Saved palettes
│   └── ...
└── history.json         # Recent colors
```

**Platform Support:**
- Linux: `~/.config/prism/`
- macOS: `~/Library/Application Support/prism/`
- Windows: `%APPDATA%/prism/`

### UI Module (`internal/ui/`)

**Responsibilities:**
- Render all 6 screens
- Handle user input
- Apply themes
- Manage navigation

**Styling Strategy:**
- Centralized styles in `styles.go`
- Theme-aware color application
- Responsive layouts using Lipgloss
- Consistent component styling

## Data Flow

### Example: Generating a Palette

```
User Input (Enter key)
    │
    ▼
GeneratorModel.Update()
    │
    ├── Parse base color
    ├── Call palette.Generate()
    │   │
    │   ├── Calculate harmony angles
    │   ├── Generate colors for each angle
    │   └── Return Palette struct
    │
    └── Update model state with palette
        │
        ▼
GeneratorModel.View()
    │
    ├── Render palette colors
    ├── Display swatches
    └── Return string output
```

### Example: Theme Switching

```
User Input (Ctrl+Shift+T)
    │
    ▼
RootModel.Update()
    │
    ├── Call ThemeManager.NextTheme()
    ├── Update all screen models
    └── Return updated model
        │
        ▼
Screens re-render with new theme
```

## Performance Considerations

### Optimizations

1. **View Caching:** Expensive calculations cached in Update(), not recalculated in View()
2. **Lazy Loading:** Named colors loaded once on first access
3. **Atomic Writes:** File operations use temp files + rename for safety
4. **Efficient Color Conversions:** Pre-calculated lookup tables where possible

### Targets Met

- ✅ Startup: <1s
- ✅ UI response: <100ms
- ✅ Palette generation: <500ms
- ✅ Memory: <50MB idle

## Error Handling

**Philosophy:** Never panic in application code. All errors are handled gracefully with user-friendly messages.

**Strategy:**
```go
if err != nil {
    log.Printf("Operation failed: %v", err)
    m.err = "Unable to save palette. Check file permissions."
    return m, nil  // Continue running
}
```

## Testing Strategy

### Test Coverage

- **Unit Tests:** 80%+ coverage for core packages (color, wcag, palette)
- **Integration Tests:** End-to-end workflows
- **TUI Tests:** State transitions and message handling

### Test Files

```
tests/
├── color_test.go     # Color math & conversions
├── palette_test.go   # Harmony rules & angles
├── wcag_test.go      # Contrast calculations
├── export_test.go    # Format validation
└── storage_test.go   # File I/O
```

## Security

### File Operations

- Atomic writes prevent corruption
- Advisory locks prevent concurrent access
- Path traversal prevented via filepath.Clean()
- Permissions: 0755 for directories, 0644 for files

### Input Validation

- Hex color validation with regex
- File path validation
- JSON schema validation for imports

## Future Enhancements

### Planned for v1.1+

- Clipboard support (cross-platform)
- Real-time color preview during wheel navigation
- Export to more formats (Swift, Kotlin)
- Colorblind simulation mode
- Gradient generation

### Deferred to v2.0+

- Image color extraction
- AI-powered palette suggestions
- Cloud sync (optional)
- Plugin system

## Dependencies

### Direct Dependencies

```
github.com/charmbracelet/bubbletea v1.3.10
github.com/charmbracelet/lipgloss v1.1.0
```

### Transitive Dependencies

All managed via `go.mod`. No CGO dependencies for maximum portability.

## Build & Deployment

### Build

```bash
go build -o bin/prism ./cmd/prism
```

### Cross-Platform Build

```bash
GOOS=linux GOARCH=amd64 go build -o bin/prism-linux-amd64 ./cmd/prism
GOOS=darwin GOARCH=arm64 go build -o bin/prism-darwin-arm64 ./cmd/prism
GOOS=windows GOARCH=amd64 go build -o bin/prism-windows-amd64.exe ./cmd/prism
```

### Binary Size

- Uncompressed: ~4.9MB
- Compressed (UPX): ~2MB

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup and guidelines.

## License

MIT License - See [LICENSE](LICENSE) for details.

---

**Last Updated:** November 2025
**Maintained by:** Kyanite Suite Team
