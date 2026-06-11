# Contributing to kyanite.sh

## Prerequisites

- **Go 1.26.0+** is required (workspace floor version)
- **CGO** is required for the noise app (whisper.cpp, malgo, sqlite3)
- **whisper.cpp** is optional (for STT features; apps work without it)

## Development

```bash
# Build all apps + launcher
make build

# Test a specific app
cd apps/focus && go test -race ./...

# Test the AI module
cd pkg/ai && go test -race ./...

# Workspace-wide sync
go work sync
```

## Structure

Each app is an independent Go module under `apps/`. The workspace (`go.work`) composes them for unified builds and releases. Shared modules live under `pkg/`:

| Module | Path | Purpose |
|---|---|---|
| `pkg/design` | `github.com/kyanite/design` | Themes, tokens, typography, WCAG, icons |
| `pkg/ai` | `github.com/kyanite/ai` | Unified inference brain (LLM, STT, memory) |
| `pkg/config` | `github.com/kyanite/config` | Koanf-based YAML config with env var overrides |
| `pkg/tui` | `github.com/kyanite/tui` | Shared TUI components (aipanel, stt) |

Each app's `go.mod` uses a `replace` directive to reference shared modules locally.

## AI Module

The `pkg/ai/` module provides a unified Brain for all four apps:
- **LLM**: Ollama on NUCBox over tailnet (`http://nucbox:11434`)
- **STT**: Local whisper.cpp subprocess
- **Memory**: PostgreSQL on NUCBox

Apps use `ai.DefaultConfig("appname")` and `ai.New(cfg)` to create a Brain. All apps degrade gracefully when services are unavailable.

### Adding AI to a new feature

1. Import `ai "github.com/kyanite/ai"`
2. Create a Brain with `ai.New(ai.DefaultConfig("yourapp"))`
3. Call `brain.Generate(ctx, prompt)` for LLM, `brain.TranscribeFile()` for STT
4. Check `brain.IsLLMAvailable()` before relying on AI features

### Adding AI panels to a new app

The `pkg/tui/aipanel` package provides a reusable Bubble Tea sub-model for streaming AI responses in a side panel:

1. Import `aipanel "github.com/kyanite/tui/aipanel"`
2. Create a panel with `aipanel.New(brain, width, height)` — pass your `*ai.Brain`
3. Wire into your app's `Update()` method — handle `aipanel.StreamChunk` and `aipanel.ErrorMsg` messages
4. Toggle visibility with `panel.Toggle()` (typically bound to a key like `Ctrl+A`)
5. Stream a response with `panel.Generate(prompt)` — returns a `tea.Cmd`
6. Continue reading chunks with `panel.ContinueStream(ch)` after each `StreamChunk` with `Done=false`
7. Render with `panel.View()` — returns a styled string for your layout
8. Add your prompt template to `pkg/ai/prompts.go` alongside existing ones

See existing integrations in `apps/focus/internal/tui/`, `apps/noise/internal/ui/`, `apps/syntax/internal/editor/`, and `apps/prism/internal/ui/`.

### Cross-app context

Apps share context through `brain.SaveCrossAppContext()` and `brain.GetCrossAppContext()`:

```go
// When Noise detects a mood, share it with Prism and Syntax
brain.SaveCrossAppContext(ctx, "prism", "mood", "melancholy, rainy afternoon", 0.8)
brain.SaveCrossAppContext(ctx, "syntax", "mood", "melancholy, rainy afternoon", 0.6)

// When Focus generates a daily briefing, pull in context from other apps
contexts, _ := brain.GetCrossAppContext(ctx, 10)
for _, c := range contexts {
    // c.SourceApp, c.ContextType, c.Summary, c.RelevanceScore
}
```

The `CrossAppContext` struct stores: source app, target app, context type, summary, and a relevance score (0.0–1.0). Context is persisted in PostgreSQL (`kyanite_cross_app_context` table) and degrades gracefully when the database is unreachable.

### Session save/resume

Apps persist their state on exit and restore it on launch:

```go
// On exit — save current state
brain.SaveSession(ctx, sessionID, "Editing Chapter 3", appState)

// On launch — restore previous session
session, _ := brain.LoadSession(ctx, sessionID)
if session != nil {
    // session.State contains JSON-encoded app state
}

// List recent sessions (for the launcher)
sessions, _ := brain.GetRecentSessions(ctx, 5)        // this app only
allSessions, _ := brain.GetAllRecentSessions(ctx, 10)  // across all apps
```

Sessions are stored in the `kyanite_sessions` PostgreSQL table with upsert semantics (same app + session_id updates in place).

## Config File

The `pkg/config/` module provides centralized configuration loaded from `~/.config/kyanite/config.yaml` with env var overrides via koanf.

### Adding new config fields

1. Add a field to the appropriate struct in `pkg/config/config.go`:
   - `BrainConfig` for inference settings
   - `AppConfig` for per-app settings (theme, directories)
   - Add a new struct + field on `Root` for entirely new config sections
2. Add a default value in the `defaults()` function
3. Add a koanf tag (e.g., `koanf:"new_field"`) — env vars follow the pattern `KYANITE_<SECTION>_<FIELD>`
4. Update the `Init()` function's template string with the new field
5. Update the `Show()` function's output to display the new field

CLI commands for users:
- `kyanite config init` — creates `~/.config/kyanite/config.yaml` with defaults
- `kyanite config show` — prints resolved config (defaults + file + env vars)

## Release

Tag-based releases via goreleaser. Push a `v*` tag to trigger the release pipeline. Binaries: focus, focus-mcp, noise, syntax, prism, kyanite (launcher).
