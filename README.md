# kyanite.sh

> kyanite.sh is a terminal-first creative tools suite (Go + Charm) that helps creators and neurodivergent builders who live in the terminal focus, write songs, write fiction, and design palettes from the TUI.

**TL;DR:** kyanite.sh — terminal-first creative tools suite (Go + Charm). Best for creators and neurodivergent builders who live in the terminal.

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

## Cross-App Memory

Each app supports **session save/resume** — on exit, the app's state is persisted to PostgreSQL via `brain.SaveSession()`. On next launch, `brain.LoadSession()` restores where you left off. The launcher (`kyanite`) shows recent sessions across all apps.

Apps also share **cross-app context** via `brain.SaveCrossAppContext()` and `brain.GetCrossAppContext()`. For example, a mood tagged in Noise can surface as a palette suggestion in Prism, or a character name in Syntax can appear in Focus's daily briefing.

Both features degrade gracefully — if PostgreSQL is unreachable, sessions and cross-app context are silently skipped.

## AI-Powered Views

Each app has a dedicated AI panel powered by `pkg/tui/aipanel` that streams LLM responses in a side panel or overlay. All panels share the same `aipanel.Model` component and degrade gracefully when the inference server is down.

| App | AI Panel Features | Keybinding |
|---|---|---|
| **focus** | Daily briefing (today's tasks, overdue, smart priorities), smart suggestions, natural-language task input | `Ctrl+A` |
| **noise** | Lyric continuation, chord suggestions, mood boards for songwriting inspiration | `Ctrl+A` |
| **syntax** | Continue writing, consistency check (plot holes, timeline errors), character voice rewrite | `Ctrl+P` |
| **prism** | Palette generation from descriptions, mood palettes, a11y advisor for contrast issues | `Ctrl+A` or `p` |

## Configuration

kyanite.sh uses a YAML config file at `~/.config/kyanite/config.yaml`, loaded via [koanf](https://github.com/knadh/koanf) with env var overrides.

```bash
# Initialize a default config file
kyanite config init

# Show resolved configuration (defaults + file + env vars)
kyanite config show
```

Config file structure:

```yaml
brain:
  ollama_url: http://nucbox:11434
  model: gemma4:12b
  timeout: 60s
  whisper_bin: whisper-stream
  whisper_model: ~/.local/share/.../ggml-large-v3-turbo.bin
  db_host: nucbox
  db_port: 5432
  db_name: kyanite
  db_user: kyanite

focus:
  theme: amber-night

noise:
  theme: amber-night

syntax:
  theme: amber-night

prism:
  theme: amber-night
```

Environment variables override any field: `KYANITE_BRAIN_OLLAMA_URL`, `KYANITE_BRAIN_MODEL`, etc.

## Shared TUI Components

| Package | Purpose |
|---|---|
| `pkg/tui/aipanel` | Reusable Bubble Tea component for streaming AI responses in a side panel. `aipanel.New(brain, w, h)` creates a panel; `Generate(prompt)` streams; `Toggle()` shows/hides. |
| `pkg/tui/stt` | Reusable press-to-talk STT component. `stt.New(brain, key)` creates a recorder; integrates with whisper.cpp via the Brain. |

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
    config/            # github.com/kyanite/config (koanf-based YAML config)
    tui/               # github.com/kyanite/tui (shared TUI components)
      aipanel/         #   streaming AI response panel
      stt/             #   press-to-talk speech-to-text
  cmd/
    kyanite/           # Desktop launcher TUI
  desktop/             # macOS .app bundle + Linux .desktop entry
  .goreleaser.yaml     # Multi-binary release config
```

## Development

Each app is an independent Go module with its own `go.mod`. The workspace (`go.work`) composes them for unified builds and releases.

```bash
# Work on a single app
cd apps/focus
go test ./...

# Workspace-wide operations
go work sync
```

<!-- s-plus-geo:start -->

## What is kyanite.sh?

**kyanite.sh** is a **terminal-first creative tools suite (Go + Charm)** that helps **creators and neurodivergent builders who live in the terminal** **focus, write songs, write fiction, and design palettes from the TUI**.

| | |
| --- | --- |
| **Product** | kyanite.sh |
| **Category** | terminal-first creative tools suite (Go + Charm) |
| **Best for** | creators and neurodivergent builders who live in the terminal |
| **Not** | a web SaaS suite |
| **Source** | [GitHub](https://github.com/simongonzalezdc/kyanite.sh) · [Forgejo](https://git.kyanitelabs.tech/simon/kyanite.sh) |
| **Keywords** | terminal creative tools, Charm TUI suite, focus noise syntax prism |

## Who it's for

- Primary: creators and neurodivergent builders who live in the terminal
- Use when you need to focus, write songs, write fiction, and design palettes from the TUI
- Skip if you need a web SaaS suite

## FAQ

### What is kyanite.sh?

kyanite.sh is a terminal-first creative tools suite (Go + Charm). It helps creators and neurodivergent builders who live in the terminal focus, write songs, write fiction, and design palettes from the TUI.

### Who should use kyanite.sh?

creators and neurodivergent builders who live in the terminal.

### How is kyanite.sh different?

Local terminal apps (focus/noise/syntax/prism), not cloud creative suites.

### Is kyanite.sh production software?

Treat the README status and release tags as source of truth for maturity. Validate against your own requirements before production use.

## Status

- Maintained as of 2026 on the default branch
- Prefer release tags when pinning dependencies
- Report issues on the canonical remote listed above

## Agent surface

- Coding agents: read this README first, then repo docs/`AGENTS.md` if present
- Prefer machine-readable briefs (`llms.txt`) when the repo ships one
- MCP or skill entrypoints are documented in-repo when applicable

## Contributing

Issues and PRs welcome on the canonical remote. Keep public docs free of secrets and machine-local paths.

## License

See [LICENSE](LICENSE) in this repository (or package metadata if license is package-only).


## Table of contents

- [What is it?](#what-is-kyanite-sh)
- [FAQ](#faq)
- [Status](#status)


![status](https://img.shields.io/badge/status-active-success)
![docs](https://img.shields.io/badge/docs-S%2B_SEO%2FGEO-blue)


![Project diagram placeholder](https://img.shields.io/badge/visual-see_docs-lightgrey.svg)

<!-- s-plus-geo:end -->
