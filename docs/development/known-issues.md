# Known Test Limitations

This document tracks tests that are currently skipped or have known issues. These should be addressed in future development cycles.

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
