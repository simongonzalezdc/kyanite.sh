# Known Issues & Lessons Learned

This document tracks tests that are currently skipped, known issues, and important lessons learned during development.

---

## Unicode & Terminal Rendering Issues (RESOLVED)

### The Problem: Mojibake

In January 2026, we discovered that Unicode symbols (✓, ✗, ⚠, 🎵, etc.) were rendering as garbled multi-byte sequences (mojibake) like `â€¢`, `âœ"`, `âš ` in some terminal environments.

### Root Causes Identified

1. **Encoding Mismatches**: UTF-8 characters being interpreted as Latin-1/ISO-8859-1
2. **Font Support**: Not all terminal fonts include glyphs for emoji/symbols
3. **Terminal Capability**: Some terminals don't properly support Unicode
4. **Stdout Corruption**: Logger writing to stdout interfered with Bubble Tea's terminal rendering

### Solutions Implemented

1. **Configurable Icon System** (`internal/ui/icons/`)
   - Three icon styles: ASCII (default), Unicode, NerdFont
   - ASCII default ensures universal compatibility
   - Users opt-in to Unicode/NerdFont via configuration

2. **Logger Fix** (`internal/logging/logger.go`)
   - Changed logger output from `os.Stdout` to `os.Stderr`
   - Added TUI mode that suppresses log output during rendering
   - Prevents log messages from corrupting terminal display

3. **Unicode Width Utilities** (`internal/ui/icons/width.go`)
   - Uses `rivo/uniseg` for proper grapheme cluster handling
   - `StringWidth()` returns correct display width
   - `Truncate()` respects character boundaries

### Files Changed

| File | Changes |
|------|---------|
| `internal/ui/toast.go` | ✓ → [OK], ✗ → [X], ⚠ → [!] |
| `internal/ui/settings.go` | ✔ → [OK], checkbox symbols |
| `internal/ui/sync_settings.go` | Various emoji → ASCII brackets |
| `internal/ui/voice_settings.go` | 🎤 → [MIC], progress bars |
| `internal/ui/voice_indicator.go` | Recording indicator, progress bars |
| `internal/ui/onboarding.go` | ●○ → [x][ ], bullets |
| `internal/ui/menu.go` | 🎵 → [~] |
| `internal/ui/idea_inbox.go` | Various type icons |
| `internal/ui/dashboard/*.go` | All emoji → ASCII |
| `internal/collaboration/presence.go` | Presence indicators |
| `internal/logging/logger.go` | stdout → stderr, TUI mode |

### Lessons Learned

1. **Never hardcode Unicode symbols** - Always use the icons package
2. **Default to ASCII** - Let users opt-in to fancy characters
3. **Use `rivo/uniseg` for width** - `len(s)` counts bytes, not display columns
4. **Don't truncate with `s[:n]`** - Breaks multi-byte characters
5. **Log to stderr** - stdout is owned by the TUI framework
6. **Test on multiple terminals** - iTerm, Terminal.app, Windows Terminal, etc.
7. **Font detection is unreliable** - Use configuration, not auto-detection

### Related Resources

- [Nerd Fonts](https://www.nerdfonts.com/) - Icon font collection
- [rivo/uniseg](https://github.com/rivo/uniseg) - Go Unicode segmentation
- [mattn/go-runewidth](https://github.com/mattn/go-runewidth) - Terminal width calculation
- [Charm TUI Best Practices](https://github.com/charmbracelet/bubbletea) - Bubble Tea documentation

---

## Platform-Specific Limitations

### Security Tests (internal/plugins)
- `TestMaliciousPlugin_PrivilegeEscalation` - World-writable file detection doesn't work reliably on macOS due to umask
- `TestSecurityManager_ScanForMaliciousPlugins` (world-writable portion) - Same umask issue

## Incomplete Feature Tests

### Graceful Degradation (internal/errors)
The graceful degradation feature is partially implemented. The following tests are skipped until the feature is complete:
- `TestShouldDegradeFeature` - Feature degradation window logic
- `TestHandleDependentFeatures` - Dependent feature cascading
- `TestRecoveryWhenServicesReturn` - Service recovery detection
- `TestDependencyCascade` - Multi-level dependency handling

### Advanced Backup/Recovery (internal/errors)
- `TestEnhancedBackupManager` - Enhanced backup with verification (partial implementation)
- `TestAutoRecovery` - Automatic recovery from corruption (partial implementation)
- `TestCleanupOldBackups` - Logging verification issue
- `TestBasicCorruptionCheck` - Empty file detection edge case
- `TestScanDirectory` - Markdown corruption detection

### UI Accessibility (internal/ui, internal/ui/dashboard)
These tests rely on accessibility features that need proper implementation:
- `TestAnimationManagerAccessibility`
- `TestDashboardModel_Accessibility`
- `TestDashboardModel_ScreenReaderText`
- `TestDashboardModel_Components`
- `TestDashboardModel_ResponsiveLayout`

### AI Editor Integration (internal/ui/editor)
- `TestEditorPane_ContinueModeSelection` - AI mode selection flow
- `TestEditorPane_VariationModeSelection` - AI mode selection flow
- `TestEditorPane_BrainstormModeSelection` - AI mode selection flow

## Collaboration Feature (Disabled)

All collaboration tests are gated behind the `collaboration` build tag as the feature is disabled for single-user mode. See [future/collaboration.md](../future/collaboration.md) for details.

Run collaboration tests with:
```bash
go test -tags collaboration ./...
```

## Code Quality Warnings (go vet)

The following pre-existing `go vet` warnings exist due to structs containing `sync.RWMutex` being copied by value:

- `internal/theme/performance_optimized.go` - ThemeMetrics copies mutex
- `internal/infra/db/performance_optimized.go` - PerformanceMetrics copies mutex
- `internal/collaboration/performance_optimized.go` - CollaborationMetrics copies mutex
- `internal/performance/integration.go` - Multiple mutex copy warnings
- `internal/performance/monitor.go` - AggregatedMetrics copies mutex

**Fix**: These structs should be passed by pointer, not by value.

## Action Items

1. **P1 - Error Handling**: Fix corruption detection for empty and markdown files
2. **P1 - Code Quality**: Fix mutex copy-by-value warnings in performance code
3. **P2 - Graceful Degradation**: Complete implementation of feature degradation
4. **P2 - Accessibility**: Implement proper screen reader support
5. **P3 - AI Integration**: Fix editor mode selection tests
6. **P3 - Platform**: Add CI matrix for Windows/Linux to catch platform-specific issues
7. **P3 - Autosave**: Fix race condition in periodic save callback test

---
*Last updated: January 22, 2026*
