# Exhaustive Architecture Deepening — kyanite.sh

Run the improve-codebase-architecture skill exhaustively: at least 5 full passes (or more until zero findings remain), then 3 consecutive clean runs with zero findings. Each pass explores the codebase, identifies shallow modules, applies the deletion test, deepens candidates, and verifies with tests.

## Constraints
- All changes must pass `go test ./... -race -count=1` across all modules
- All 5 binaries must build cleanly
- Every pass must leave the codebase architecturally better than it found it
- Use LANGUAGE.md vocabulary: module, interface, seam, adapter, depth, leverage, locality
- Apply deletion test rigorously before proposing any deepening

## Quality Gate
- Architect review: CLEAR on architecture, product, code
- Full test suite with -race: zero failures
- All binaries build
- Zero blocking slop-cleaner findings
- 3 consecutive clean runs with zero architecture findings

@goal: Pass 1 — Deepen Brain wrappers into unified Brain interface
Delete the 4 duplicate per-app Brain wrappers (focus BrainProvider, noise brain.Client, syntax BrainClient, prism ai.Client). Make pkg/ai.Brain the single interface all apps use directly. Add app-specific typed methods to Brain (ParseTask, ContinueWriting, GeneratePalette, etc.) so apps call brain.ParseTask(input) instead of wrapping Brain in another struct. Update all call sites. All 60+ tests pass.

@goal: Pass 2 — Centralize config loading to single init point
Replace all scattered ai.DefaultConfig() + ai.New() calls with single pkg/config.Load() → ai.ConfigFromRoot() → one Brain per app at startup. Delete 8+ duplicate Brain construction sites in focus CLI commands. Make Brain a constructor-injected dependency. All tests pass.

@goal: Pass 3 — Extract Focus Manager into focused modules
Extract LRU cache from focus/internal/ai/manager.go into pkg/cache/lru.go. Extract fallback parsing into pkg/ai/fallback.go. Delete backward-compat stubs (IsOllamaAvailable, LaunchOllama). Manager becomes thin orchestration or is deleted entirely. All tests pass.

@goal: Pass 4 — Extract shared session lifecycle
Create pkg/session/lifecycle.go with OnStartup(brain, app) and OnShutdown(brain, session, title, state). Replace 4 copies of session lifecycle in each app's main.go. Deep module hides session ID generation, welcome printing, cross-app context saving. All tests pass.

@goal: Pass 5 — Extract shared AI panel integration
Create pkg/tui/aipanel/integration.go with WrapModel that handles toggle, streaming, and rendering automatically. Replace 4 separate AI panel wiring implementations. All tests pass.

@goal: Pass 6 — Create shared test infrastructure
Create pkg/test/ with NewTestBrain(t) backed by httptest.Server mock of Ollama. Every app's AI tests use it without NUCBox. All tests pass without network.

@goal: Pass 7 — Exhaustive clean run 1
Run improve-codebase-architecture exploration across the entire codebase. Target: zero findings. If findings exist, deepen them and re-run. All tests pass.

@goal: Pass 8 — Exhaustive clean run 2
Second consecutive clean run. Zero findings. All tests pass.

@goal: Pass 9 — Exhaustive clean run 3 (final)
Third consecutive clean run. Zero findings. All tests pass. Declare exhaustive.
