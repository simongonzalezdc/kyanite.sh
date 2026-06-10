# prism.sh

> Terminal-based color palette design and theory tool

**Status:** ✅ v1.0.0 - Complete and Ready for Use

Prism.sh is a fully functional terminal-based color palette design tool with interactive TUI, comprehensive color theory education, and WCAG accessibility validation.

## Quick Start

```bash
# Install dependencies
go mod download

# Build and run
make build
./bin/prism

# Or run directly
make run
```

## Features ✅

### Interactive Terminal UI
- **Bubble Tea TUI** - Full interactive experience with 6 screens
  - Main menu navigation
  - Interactive color wheel with live preview
  - Palette generator with all 7 harmony rules
  - Color theory lessons (5 comprehensive lessons)
  - WCAG accessibility checker with live validation
  - Palette manager (save/load/delete)
- **Theme System** - 10 Kyanite themes with runtime switching (Ctrl+Shift+T)
- **Help System** - Context-sensitive help overlay (Ctrl+H)
- **80x24 Support** - Fully functional in minimum terminal size

### Color Mathematics & Theory
- **Color Conversions** - RGB ↔ HSL ↔ HSV ↔ Hex with ±1 accuracy
- **Color Operations** - Lighten, darken, saturate, desaturate, complement
- **Temperature Detection** - Warm/cool color classification
- **Named Colors** - 147 CSS Color Module Level 4 colors with fuzzy search
- **Theory Education** - Interactive lessons on complementary, analogous, triadic, warm/cool, tints/shades/tones

### Accessibility & Standards
- **WCAG 2.1 Validation** - Full gamma-corrected contrast calculations
- **AA/AAA Compliance** - Real-time validation with visual feedback
- **Accuracy** - ±0.05 contrast ratio precision
- **Color-Coded Results** - Green (AAA), Yellow (AA), Red (FAIL)

### Palette Generation
- **7 Harmony Rules** - Monochromatic, Complementary, Analogous, Triadic, Tetradic, Split-Complementary, Square
- **Angle Accuracy** - ±5° precision on all harmony calculations
- **Contrast Validation** - Ensures minimum 3:1 contrast in generated palettes

### Data Management
- **Save/Load Palettes** - Persistent storage with atomic writes
- **Export Formats** - JSON, CSS (variables + RGB), TOML, Kyanite theme format
- **Cross-Platform** - XDG Base Directory support (Linux/macOS/Windows)
- **Clipboard Support** - Copy colors to system clipboard (xclip/pbcopy/clip.exe)

### Quality Assurance
- **Comprehensive Tests** - 70%+ code coverage with race detection
- **CI/CD Pipelines** - Automated testing, building, and releases
- **Documentation** - ARCHITECTURE.md, CONTRIBUTING.md, SECURITY.md, CHANGELOG.md
- **Version Control** - Semantic versioning with git tags

## Project Structure

```
prism/
├── cmd/prism/           # ✅ Entry point with CLI flags
├── internal/
│   ├── app/             # ✅ Bubble Tea root model
│   ├── ui/              # ✅ All 6 TUI screens + components
│   ├── color/           # ✅ Color math + named colors
│   ├── palette/         # ✅ 7 harmony rule generators
│   ├── wcag/            # ✅ WCAG 2.1 validation
│   ├── theme/           # ✅ 10 Kyanite themes
│   ├── export/          # ✅ 4 export formats
│   ├── storage/         # ✅ Atomic file I/O
│   ├── clipboard/       # ✅ Cross-platform clipboard
│   └── data/            # ✅ Embedded named colors
├── tests/               # ✅ Comprehensive test suite
├── data/                # ✅ Named colors database
├── .github/workflows/   # ✅ CI/CD pipelines
└── docs/                # ✅ Complete documentation
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test
```

### Run
```bash
make run
```

### Format
```bash
make fmt
```

## Architecture

See [ARCHITECTURE.md](ARCHITECTURE.md) and [02-TDD-2.md](02-TDD-2.md) for complete technical design.

### Core Principles
- **Bubble Tea** (Elm architecture) for TUI
- **Lipgloss** for terminal styling with theme support
- **XDG Base Directory** specification for cross-platform config
- **Atomic writes** for data safety (write-to-temp-then-rename)
- **Message-based navigation** between screens
- **Embedded resources** using Go embed for portability

### Performance Metrics
- **Startup time**: <1s
- **Binary size**: 4.9MB (meets <10MB target)
- **Memory usage**: <50MB idle
- **UI responsiveness**: <100ms for all interactions
- **Search performance**: <100ms for fuzzy color search

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Key areas for enhancement:
1. **Additional Color Spaces** - CMYK, LAB, LCH support
2. **More Export Formats** - SCSS, LESS, Tailwind config
3. **Palette Analysis** - Color blindness simulation
4. **More Themes** - Community-contributed themes
5. **CLI Mode** - Non-interactive command-line interface

See detailed design in [01-PRD.md](01-PRD.md) and [02-TDD-2.md](02-TDD-2.md).

## Development Timeline

Total development: 12-16 days full-time

| Phase | Duration | Status |
|-------|----------|---------|
| Phase 1: Foundation | 2-3 days | ✅ Complete |
| Phase 2: TUI Implementation | 3-4 days | ✅ Complete |
| Phase 3: Data & Content | 2 days | ✅ Complete |
| Phase 4: Storage & Export | 1-2 days | ✅ Complete |
| Phase 5: Testing | 2-3 days | ✅ Complete |
| Phase 6: Documentation | 1-2 days | ✅ Complete |
| Phase 7: CI/CD | 1 day | ✅ Complete |
| Phase 8-9: Polish & Release | 1 day | ✅ Complete |

## Usage

### Command-Line Options

```bash
prism           # Start interactive TUI
prism -version  # Display version information
prism -help     # Display help message
```

### Keyboard Shortcuts

**Global:**
- `Ctrl+Q` - Quit application
- `Ctrl+H` - Toggle help overlay
- `Ctrl+Shift+T` - Cycle through themes
- `Esc` - Navigate back/cancel

**Navigation:**
- `↑/↓` or `j/k` - Move up/down in menus
- `←/→` or `h/l` - Adjust values, navigate wheel
- `Enter` - Select/confirm
- `Tab` - Next field
- `Shift+Tab` - Previous field

**Screen-Specific:**
- Color Wheel: Arrow keys to navigate hue, `c` to copy
- Generator: Number keys (1-7) to select harmony rule
- Manager: `d` to delete palette, `Enter` to load

### Configuration

Prism.sh stores configuration and palettes in:
- **Linux**: `~/.config/prism/`
- **macOS**: `~/Library/Application Support/prism/`
- **Windows**: `%APPDATA%\prism\`

Configuration file: `config.toml`

## License

MIT License - See [LICENSE](LICENSE)

## Part of Kyanite Suite

Prism.sh is an independent tool in the Kyanite Suite of terminal-based creative tools.

---

**v1.0.0** - Built with ❤️ using Go, Bubble Tea, and Lipgloss
