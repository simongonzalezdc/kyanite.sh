# kyanite.sh

A unified suite of terminal-first creative tools built with Go and the Charm stack.

## Apps

| App | Description |
|---|---|
| **focus** | ADHD-friendly terminal focus & task manager |
| **noise** | Local-first AI-assisted songwriting studio |
| **syntax** | Terminal fiction-writing & worldbuilding TUI |
| **prism** | Terminal color-palette design with WCAG contrast validation |

## Prerequisites

- Go 1.26.0+
- CGO (for noise — whisper.cpp, malgo, sqlite3)
- whisper.cpp (for STT — optional, apps work without it)

## Quick Start

```bash
# Build all apps + launcher
make build

# Run a specific app
go run ./apps/focus/cmd/focus
go run ./apps/noise/cmd/noise
go run ./apps/syntax/cmd/syntax
go run ./apps/prism/cmd/prism

# Run the desktop launcher
go run ./cmd/kyanite

# Run all tests
make test

# Lint
make lint
```

## AI Features

All four apps share a unified inference brain (`pkg/ai/`) that provides:

- **LLM**: Ollama on the NUCBox inference server (gemma4:12b on AMD ROCm GPU)
- **STT**: Local whisper.cpp (whisper-large-v3-turbo on Metal/CPU)
- **Memory**: PostgreSQL on NUCBox (conversation history, per-app context)

Configuration via environment variables:

```bash
export KYANITE_OLLAMA_URL=http://nucbox:11434  # Ollama endpoint
export KYANITE_MODEL=gemma4:12b                # LLM model
export KYANITE_WHISPER_MODEL=~/.local/share/.../ggml-large-v3-turbo.bin  # STT model
```

All apps work fully offline when the inference server is unreachable.

## Project Structure

```
kyanite.sh/
  go.work              # Go workspace definition
  apps/
    focus/             # github.com/kyanite/focus
    noise/             # github.com/kyanite/noise
    syntax/            # github.com/kyanite/syntax
    prism/             # github.com/kyanite/prism
  pkg/
    design/            # github.com/kyanite/design (themes, tokens, WCAG, icons)
    ai/                # github.com/kyanite/ai (unified inference brain)
  cmd/
    kyanite/           # Desktop launcher TUI
  desktop/             # macOS .app bundle + Linux .desktop entry
  .goreleaser.yaml     # Multi-binary release config

## Development

Each app is an independent Go module with its own `go.mod`. The workspace (`go.work`) composes them for unified builds and releases.

```bash
# Work on a single app
cd apps/focus
go test ./...

# Workspace-wide operations
go work sync
```
