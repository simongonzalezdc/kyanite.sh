# kyanite.sh Architecture

## Overview

kyanite.sh is a Go workspace (`go.work`) containing four independent TUI applications built on the Charm stack (Bubble Tea, Lipgloss, Bubbles). Each app is a standalone module under `apps/` with its own `go.mod`, `cmd/`, `internal/`, and test directories.

```
kyanite.sh/
├── go.work              # workspace definition (go 1.26.0)
├── apps/
│   ├── focus/           # journaling & productivity (cobra + bubbletea + huh)
│   ├── noise/           # music creation & voice capture (bubbletea + cgo/whisper.cpp)
│   ├── syntax/          # interactive fiction editor (bubbletea + custom editor)
│   └── prism/           # color palette & WCAG contrast tool (bubbletea + lipgloss)
├── .github/workflows/   # CI (per-app jobs) + release (goreleaser)
├── .golangci.yml
├── .goreleaser.yaml
└── Makefile
```

## Applications

### focus — Journaling & Productivity

| Layer | Package | Responsibility |
|---|---|---|
| CLI | `cmd/focus`, `cmd/focus-mcp` | cobra commands, MCP server |
| Presentation | `internal/cli`, `internal/tui` | lipgloss rendering, bubbletea models |
| Application | `internal/engine`, `internal/command` | business logic orchestration |
| Domain | `pkg/models`, `pkg/validation` | data types, validation rules |
| AI | `internal/ai` | Ollama integration, prompt management |
| Infrastructure | `internal/store`, `internal/repository`, `internal/journal` | SQLite storage, file I/O |
| Shared | `pkg/config`, `pkg/styles`, `pkg/calendar`, `pkg/utils` | reusable utilities |

Key binaries: `focus` (CLI/TUI), `focus-mcp` (Model Context Protocol server).

### noise — Music Creation & Voice Capture

| Layer | Package | Responsibility |
|---|---|---|
| CLI | `cmd/noise` | entry point |
| Presentation | `internal/ui/*` | bubbletea views (dashboard, editor, settings) |
| Application | `internal/app/*` | AI, theory, knowledge, agent services |
| Domain | `internal/domain` | core music data types |
| Infrastructure | `internal/infra/*` | db, files, voice (whisper.cpp via cgo), sync (websocket), ollama |
| Collaboration | `internal/collaboration` | real-time multi-user editing |

**CGO dependency**: noise requires `CGO_ENABLED=1` for whisper.cpp, malgo, and sqlite3. This limits cross-compilation to linux/amd64 + linux/arm64.

### syntax — Interactive Fiction Editor

| Layer | Package | Responsibility |
|---|---|---|
| CLI | `cmd/syntax` | entry point |
| Presentation | `internal/editor` | custom terminal editor buffer with UTF-8 handling |
| Application | `internal/app`, `internal/scene` | story compilation, scene management |
| Domain | `internal/story`, `internal/character`, `internal/location` | narrative data model |
| Infrastructure | `internal/storage`, `internal/backup`, `internal/export` | file storage, zip backup, markdown export |
| AI | `internal/ai` | writing assistance |
| Support | `internal/spellcheck`, `internal/theme`, `internal/utils` | spell checking, theming, utilities |

### prism — Color Palette & WCAG Contrast

| Layer | Package | Responsibility |
|---|---|---|
| CLI | `cmd/prism` | entry point |
| Presentation | `internal/ui` | bubbletea views |
| Application | `internal/app` | app orchestration |
| Domain | `internal/color`, `internal/palette`, `internal/wcag` | color theory, palette generation, contrast checking |
| Infrastructure | `internal/storage`, `internal/clipboard`, `internal/export` | file persistence, clipboard, export |

## Cross-Module Patterns

### What's Shared (Convention, Not Code)

| Pattern | Implementation |
|---|---|
| TUI Framework | All apps use Bubble Tea with Lipgloss styling |
| CLI Framework | focus uses cobra; noise/syntax/prism use custom bubbletea entry |
| Storage | focus & syntax use file-based JSON; noise uses sqlite3; prism uses file-based JSON |
| AI Integration | focus & noise & syntax have Ollama-based AI features |
| Theming | Each app has its own theme registry with divergent structures |
| Error Handling | Raw `fmt.Errorf` with `%w` wrapping (no structured error library) |
| Testing | Standard `testing` package + `t.TempDir()` for isolation |

### Known Divergence

| Area | Issue |
|---|---|
| Theme registries | 4 independent copies: focus (11 fields, lipgloss.Color), noise (8 fields, string), syntax (6 fields, string), prism (10 fields, string) |
| Layer boundaries | 20+ Presentation→Infrastructure import violations (noise: 14 files, focus: 6+ files) |
| Config loading | focus uses viper; others use ad-hoc flag/env parsing |
| Error types | No shared error types or wrapping conventions |
| Storage paths | focus uses `~/.focus`; no XDG compliance |

## Shared Module Proposal (Deferred)

A future extraction could create `shared/` modules under the workspace:

| Module | Contents | Consumers |
|---|---|---|
| `shared/themes` | Unified Theme type, lipgloss.Color conversion, all palettes | all 4 apps |
| `shared/config` | koanf-based config loading with XDG paths | all 4 apps |
| `shared/errors` | Structured error types (eris/fault) | all 4 apps |
| `shared/storage` | Repository interface + file/json/sqlite implementations | focus, noise, syntax |
| `shared/ai` | Ollama client, prompt sanitization, response parsing | focus, noise, syntax |

**Status**: Deferred — requires coordinated changes across all 4 apps. The theme unification alone needs resolving lipgloss.Color vs string type differences.

## CI/CD

- **CI** (`.github/workflows/ci.yml`): 4 parallel per-app jobs (build, test, lint) + all-green release gate
- **Release** (`.github/workflows/release.yml`): goreleaser with 5 binaries (focus, focus-mcp, noise, syntax, prism)
- **Dependabot**: Monitors all 4 modules + github-actions
- **Govulncheck**: All stdlib vulns addressed via `toolchain go1.26.4` directive

## Security Posture

| Control | Status |
|---|---|
| GitHub Actions SHA-pinned | ✅ All actions pinned to commit SHA |
| Permissions blocks | ✅ CI: contents:read, Release: contents:write |
| Prompt injection | ✅ User content sanitized + delimited in AI prompts |
| Path traversal | ✅ validateProjectID (regex), filepath.Base, io.LimitReader |
| Test isolation | ✅ t.TempDir() throughout |
| Secret hygiene | ✅ .gitignore covers .env, .pem, .key, *.db |
| Goroutine leaks | ✅ Close() with done channels in AI manager, collaboration |
| Race conditions | ✅ sync.Mutex on store, RLock fix on websocket |
| Secrets scan | ✅ Clean — no leaked credentials |
