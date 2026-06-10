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

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `KYANITE_OLLAMA_URL` | `http://nucbox:11434` | Ollama API endpoint |
| `KYANITE_MODEL` | `gemma4:12b` | LLM model |
| `KYANITE_WHISPER_BIN` | `whisper-stream` | whisper.cpp binary |
| `KYANITE_WHISPER_MODEL` | auto-detected | GGML model path |
| `KYANITE_DB_HOST` | `nucbox` | PostgreSQL host |

## Release

Tag-based releases via goreleaser. Push a `v*` tag to trigger the release pipeline. Binaries: focus, focus-mcp, noise, syntax, prism, kyanite (launcher).
