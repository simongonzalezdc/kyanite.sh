# kyanite.sh Architecture

## Overview

kyanite.sh is a Go workspace (`go.work`) containing four independent TUI applications built on the Charm stack (Bubble Tea, Lipgloss, Bubbles). Each app is a standalone module under `apps/` with its own `go.mod`, `cmd/`, `internal/`, and test directories. A shared design module at `pkg/design/` provides unified themes, style tokens, and WCAG validation. A shared AI module at `pkg/ai/` provides a unified inference brain (LLM + STT + memory) consumed by all four apps.

```
kyanite.sh/
├── go.work              # workspace definition (go 1.26.0)
├── pkg/design/          # shared design system (themes, tokens, typography, WCAG, icons)
├── pkg/ai/              # unified inference brain (LLM, STT, memory)
├── apps/
│   ├── focus/           # journaling & productivity (cobra + bubbletea + huh)
│   ├── noise/           # music creation & voice capture (bubbletea + cgo/whisper.cpp)
│   ├── syntax/          # interactive fiction editor (bubbletea + custom editor)
│   └── prism/           # color palette & WCAG contrast tool (bubbletea + lipgloss)
├── cmd/kyanite/         # desktop launcher (bubbletea TUI menu)
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
| AI | `internal/ai` | BrainProvider → pkg/ai → NUCBox Ollama, prompt management |
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
| Infrastructure | `internal/infra/*` | db, files, voice (whisper.cpp via cgo), sync (websocket), brain (pkg/ai wrapper) |
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
| AI | `internal/ai` | BrainClient → pkg/ai → NUCBox Ollama, writing assistance |
| Support | `internal/spellcheck`, `internal/theme`, `internal/utils` | spell checking, theming, utilities |

### prism — Color Palette & WCAG Contrast

| Layer | Package | Responsibility |
|---|---|---|
| CLI | `cmd/prism` | entry point |
| Presentation | `internal/ui` | bubbletea views |
| Application | `internal/app` | app orchestration |
| Domain | `internal/color`, `internal/palette`, `internal/wcag` | color theory, palette generation, contrast checking |
| Infrastructure | `internal/storage`, `internal/clipboard`, `internal/export` | file persistence, clipboard, export |
| AI | `internal/ai` | Brain → pkg/ai → NUCBox Ollama, AI palette generation, accessible color suggestions |

## Cross-Module Patterns

### What's Shared (Convention, Not Code)

| Pattern | Implementation |
|---|---|
| TUI Framework | All apps use Bubble Tea with Lipgloss styling |
| CLI Framework | focus uses cobra; noise/syntax/prism use custom bubbletea entry |
| Storage | focus & syntax use file-based JSON; noise uses sqlite3; prism uses file-based JSON |
| AI Integration | All 4 apps use pkg/ai Brain → NUCBox Ollama (LLM) + local whisper.cpp (STT) + NUCBox PostgreSQL (memory) |
| Theming | Each app has its own theme registry with divergent structures |
| Error Handling | Raw `fmt.Errorf` with `%w` wrapping (no structured error library) |
| Testing | Standard `testing` package + `t.TempDir()` for isolation |

### Known Divergence

| Area | Issue |
|---|---|
| Theme registries | NOW UNIFIED via `pkg/design/` — 11 fields, lipgloss.Color, 6 themes, WCAG AA enforced |
| Layer boundaries | 20+ Presentation→Infrastructure import violations (noise: 14 files, focus: 6+ files) |
| Config loading | focus uses viper; others use ad-hoc flag/env parsing |
| Error types | No shared error types or wrapping conventions |
| Storage paths | focus uses `~/.focus`; no XDG compliance |

## Shared Design Module (`pkg/design/`)

The `pkg/design/` module (`github.com/kyanite/design`) provides a unified design system consumed by all 4 apps:

| File | Contents |
|---|---|
| `theme.go` | `Theme` struct (11 `lipgloss.Color` fields), `RegisterBuiltIn`/`RegisterCustom`/`Get`/`List` registry |
| `themes.go` | 6 curated dark-first themes: monochrome, amber-night, indigo-depths, forest-path, cyan-wave, electric-rose |
| `tokens.go` | `TokenSet` — 10 semantic style tokens + `Extensions` map for app-specific styles |
| `typography.go` | `TypographySet` — 4 text levels: Title (bold+accent), Heading (bold), Body, Muted (faint) |
| `wcag.go` | WCAG 2.1 contrast ratio calculation + `ValidateThemeAA` enforcement |
| `border.go` | `RoundedBox`/`Panel` helpers with rounded borders + spacing constants (0-4) |
| `icons/` | 3-tier icon system (ASCII/Unicode/NerdFont) + string width utilities |

## Shared AI Module (`pkg/ai/`)

The `pkg/ai/` module (`github.com/kyanite/ai`) provides a unified inference brain consumed by all 4 apps:

| File | Contents |
|---|---|
| `brain.go` | `Brain` struct — main entry point with `New(cfg)`, `Generate()`, `Chat()`, `TranscribeFile()`, `TranscribePCM()`, `SaveContext()`, `LoadContext()`, `SaveMessage()`, `LoadConversation()` |
| `config.go` | `Config` struct with `DefaultConfig(app)` — reads `KYANITE_OLLAMA_URL`, `KYANITE_MODEL`, `KYANITE_WHISPER_BIN`, `KYANITE_DB_HOST` etc. |
| `llm.go` | `LLMClient` — Ollama `/api/chat` client with `Generate()`, `Chat()`, `IsAvailable()` |
| `stt.go` | `STTClient` — whisper.cpp subprocess client with `TranscribeFile()`, `TranscribePCM()`, `IsAvailable()` |
| `memory.go` | `MemoryStore` — PostgreSQL store with `kyanite_context` and `kyanite_conversations` tables |
| `prompts.go` | Centralized prompt templates: `FocusParsePrompt`, `FocusSuggestPrompt`, `FocusSummarizePrompt`, `NoiseBrainstormPrompt`, `NoiseContextPrompt`, `SyntaxSuggestPrompt`, `PrismPalettePrompt`, `PrismContrastPrompt` |
| `types.go` | `Message`, `Transcription`, `GenerateOptions`, `ContextEntry` |

### Inference Architecture

```
Local machine (MacBook Air / Mac Mini):
  ┌─────────────────────────────┐
  │  kyanite TUI app             │
  │  STT: whisper.cpp (local)    │
  │  LLM: HTTP → nucbox:11434    │
  └──────────────┬──────────────┘
                 │ tailnet (encrypted)
                 ▼
  ┌──────────────────────────────┐
  │  NUCBox (Inference Server)   │
  │  Ollama (gemma4:12b on GPU)  │
  │  PostgreSQL (memory store)   │
  └──────────────────────────────┘
```

### Per-App Integration

| App | Integration Layer | LLM Features | STT |
|---|---|---|---|
| focus | `internal/ai/brain_provider.go` (implements Provider interface) | Task parsing, suggestions, summaries, chat | Planned |
| noise | `internal/infra/brain/client.go` (implements QuickLLMClient interface) | Songwriting brainstorming, quick ideas | Voice capture → Brain.TranscribePCM() |
| syntax | `internal/ai/brain_client.go` (implements Provider interface) | Writing suggestions (continue/improve/dialogue/description) | Planned |
| prism | `internal/ai/client.go` (standalone) | Palette generation from descriptions, accessible color suggestions | N/A |

### Environment Variables

| Variable | Default | Purpose |
|---|---|---|
| `KYANITE_OLLAMA_URL` | `http://nucbox:11434` | Ollama API endpoint |
| `KYANITE_MODEL` | `gemma4:12b` | LLM model name |
| `KYANITE_WHISPER_BIN` | `whisper-stream` | whisper.cpp binary |
| `KYANITE_WHISPER_MODEL` | auto-detected (ggml-large-v3-turbo.bin) | GGML model path |
| `KYANITE_WHISPER_LANG` | `en` | STT language |
| `KYANITE_DB_HOST` | `nucbox` | PostgreSQL host |
| `KYANITE_DB_NAME` | `kyanite` | Database name |
| `KYANITE_DB_USER` | `kyanite` | Database user |

### Graceful Degradation

All apps work fully offline:
- If NUCBox is unreachable: `IsLLMAvailable()` returns false, AI features disabled
- If whisper.cpp not found: `IsSTTAvailable()` returns false, voice disabled
- If PostgreSQL unreachable: memory silently skipped, no conversation history
- Focus falls back through BrainProvider → legacy Ollama → OpenRouter → rule-based parsing

### Remaining Extraction Candidates

| Module | Contents | Priority |
|---|---|---|
| `shared/config` | koanf-based config loading with XDG paths | MEDIUM |
| `shared/errors` | Structured error types (eris/fault) | MEDIUM |
| `shared/storage` | Repository interface + file/json/sqlite implementations | LOW |
| `pkg/ai` | Unified inference brain: Ollama LLM client, whisper.cpp STT, PostgreSQL memory, centralized prompts | **DONE** |

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

## Charm Ecosystem Integration

### Currently Used

| Library | Version | Apps | Purpose |
|---|---|---|---|
| bubbletea | v1.3.10 | all 4 | Elm-architecture TUI framework |
| lipgloss | v1.1.1 | all 4 | Terminal styling and layout |
| bubbles | v1.0.0 | focus, noise, prism | TUI components (list, input, spinner) |
| huh | v1.0.0 | focus | Terminal forms and prompts |
| glamour | v1.0.0 | noise | Markdown rendering |
| harmonica | v0.2.0 | noise | Physics-based animations |
| colorprofile | v0.4.1 | all (indirect) | Terminal color capability detection |
| x/ansi | v0.11.6 | all (indirect) | ANSI sequence handling |
| x/cellbuf | v0.0.15 | all (indirect) | Cell buffer management |
| x/term | v0.2.2 | all (indirect) | Terminal control |

### Recommended Adoption

| Library | Stars | Why | Priority | Apps | Status |
|---|---|---|---|---|---|
| **charmbracelet/log** | 3.3k | Replace noise's 425-line hand-rolled logger + focus/syntax bare `log`. Structured, colorful, lipgloss-native, slog handler, TUI-aware. | **HIGH** | all 4 | ADOPTED (commit 2128914) |
| **charmbracelet/fantasy** | 799 | Multi-provider AI agent framework. Replaced by `pkg/ai` which provides unified Brain with Ollama LLM, whisper.cpp STT, and PostgreSQL memory. | ~~HIGH~~ | — | **SUPPLANTED by pkg/ai** |
| **charmbracelet/ultraviolet** | 351 | Powers Bubble Tea v2. Cell-based diffing renderer, constraint-based layout (Cassowary), universal input. Useful for syntax's editor buffer or prism's color rendering. | MEDIUM | syntax, prism | Watch for bubbletea v2 migration |
| **charmbracelet/sequin** | 805 | Human-readable ANSI sequence library. Could simplify syntax's editor ANSI handling and prism's terminal output. | MEDIUM | syntax, prism | Evaluate |
| **charmbracelet/glow** | 26k | Markdown rendering on the CLI. Focus's journal viewer could use rich Markdown rendering. Already have glamour in noise. | LOW | focus | Consider for journal view |
| **charmbracelet/vhs** | 20k | Terminal recording for demos. Could generate GIF demos for README. | LOW | all (docs) | CI integration |
| **charmbracelet/catwalk** | 732 | LLM model/provider database. Companion to fantasy. | LOW | focus, noise | If fantasy adopted |
| **charmbracelet/fang** | 1.9k | CLI starter kit. Not needed — already using cobra + bubbletea. | SKIP | — | N/A |
| **charmbracelet/skate** | 1.8k | Personal KV store. Not applicable. | SKIP | — | N/A |
| **charmbracelet/soft-serve** | 7k | Self-hosted Git server. Not applicable. | SKIP | — | N/A |

### Adoption Roadmap

1. **Phase 1: charmbracelet/log** — ✅ DONE. Replaced all logging across 4 apps.
2. **Phase 2: Unified AI** — ✅ DONE. Built `pkg/ai/` with Ollama LLM client, whisper.cpp STT, PostgreSQL memory, centralized prompts. All 4 apps migrated. Supersedes charmbracelet/fantasy.
3. **Phase 3: charmbracelet/ultraviolet** — If bubbletea v2 migration happens, ultraviolet provides the rendering primitives. Syntax's editor buffer could use the cell-based renderer directly.
4. **Phase 4: charmbracelet/sequin + glow** — Polish pass on syntax's ANSI handling and focus's journal rendering.
