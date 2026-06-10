# Changelog

All notable changes to Prism.sh will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Clipboard support (cross-platform)
- Color name search screen
- Export history tracking
- Gradient generation
- Colorblind simulation mode

## [1.0.0] - 2025-11-11

### Added

**Core Features**
- Interactive color wheel with 360° hue navigation
- 7 harmony rule palette generators (monochromatic, complementary, analogous, triadic, tetradic, split-complementary, square)
- WCAG 2.1 accessibility checker with AA/AAA validation
- Color theory education system with 5 lessons
- Palette manager (save, load, delete)
- Named colors database (147 CSS colors)
- Fuzzy search for color names

**User Interface**
- Full Bubble Tea TUI with 6 screens
- 10 Kyanite themes for UI customization
- Theme switching (Ctrl+Shift+T)
- Help system (Ctrl+H)
- Responsive layouts (80x24 minimum)
- Real-time color preview with swatches

**Technical**
- Complete color space conversions (RGB ↔ HSL ↔ HSV ↔ Hex)
- WCAG 2.1 contrast calculation with gamma correction
- Cross-platform file storage (Linux, macOS, Windows)
- Atomic file writes for data safety
- Export to 4 formats (JSON, CSS, TOML, Kyanite)

**Documentation**
- Comprehensive README
- ARCHITECTURE.md with technical details
- CONTRIBUTING.md with development guidelines
- SECURITY.md with security policy
- TODO.md with implementation checklist
- Inline code documentation

**Development**
- Complete test suite (>80% coverage)
- CI/CD with GitHub Actions
- Multi-platform builds (Linux, macOS, Windows)
- Makefile for common tasks

### Global Keyboard Shortcuts
- `Ctrl+Q` - Quit application
- `Ctrl+H` - Toggle help
- `Ctrl+Shift+T` - Cycle themes
- `Esc` - Back to menu / Quit
- `↑/↓` or `k/j` - Navigate up/down
- `←/→` or `h/l` - Navigate left/right
- `Enter` or `Space` - Select/Confirm

### Color Wheel Shortcuts
- `←/→` - Adjust hue
- `↑/↓` - Adjust lightness
- `+/-` - Adjust saturation

### Palette Manager Shortcuts
- `D` - Delete selected palette
- `R` - Refresh list

## [0.2.0] - 2025-11-11 (Development)

### Added
- Bubble Tea TUI implementation
- All 6 screens functional
- Theme system with 10 themes
- Navigation between screens
- Help overlay

## [0.1.0] - 2025-11-11 (Foundation)

### Added
- Project structure
- Core color mathematics
- WCAG contrast calculations
- Palette generation algorithms
- Theme registry
- Basic CLI demo

---

## Version History Summary

- **v1.0.0** - First stable release with all core features
- **v0.2.0** - TUI implementation
- **v0.1.0** - Foundation and core algorithms

## Upgrade Guide

### From Development Builds to v1.0.0

If you were using development builds:

1. **Configuration:** Config format is stable, no changes needed
2. **Palettes:** All saved palettes are compatible
3. **New Features:** See v1.0.0 added features above

### Fresh Installation

See [README.md](README.md) for installation instructions.

## Breaking Changes

### v1.0.0
- None (first stable release)

## Deprecations

### v1.0.0
- None

## Known Issues

### v1.0.0
- Clipboard support not yet implemented (planned for v1.1)
- Terminal must support true color for best experience
- Windows: Color display may vary by terminal emulator

### Workarounds
- **16-color terminals:** Colors shown as hex codes with warning
- **Clipboard:** Use export to file instead

## Migration Guides

### From Pre-release to v1.0.0

No migration needed - fresh installation recommended.

## Performance Improvements

### v1.0.0
- Startup time: <1s
- UI response: <100ms
- Palette generation: <500ms
- Memory usage: <50MB idle
- Binary size: ~4.9MB

## Dependencies

### v1.0.0
- Go 1.21+
- github.com/charmbracelet/bubbletea v1.3.10
- github.com/charmbracelet/lipgloss v1.1.0

## Contributors

### v1.0.0
- Initial implementation by Kyanite Suite Team

---

**Note:** This changelog is manually maintained. For a complete list of changes, see the [commit history](https://github.com/kyanite/prism/commits/main).

## Reporting Issues

Found a bug? Have a feature request? Open an issue on [GitHub](https://github.com/kyanite/prism/issues).

## Stay Updated

- **Watch** the repository for new releases
- **Star** to show your support
- **Follow** [@KyaniteSuite](https://github.com/kyanite) for updates

---

**Last Updated:** November 11, 2025
