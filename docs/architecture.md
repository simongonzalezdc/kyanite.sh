# noise.sh Technical Architecture

## 1. Overview
noise.sh is delivered as a modular monolith that follows the Elm-inspired Model–Update–View loop provided by Bubble Tea. The terminal UI, application services, AI agents, and persistence utilities stay in a single repository and process space, which keeps operations deterministic while still enabling clean layering. The UI leverages Bubble Tea for event routing and Lip Gloss for styling so that keyboard-first workflows remain predictable across Windows, macOS, and Linux. Product positioning in [`PROJECT_SNAPSHOT.md`](../PROJECT_SNAPSHOT.md) emphasizes a rapid capture studio that keeps all processing on-device; the architecture reflects that by avoiding external dependencies beyond optional Ollama model downloads.

## 2. Technology Stack
| Layer | Technology | Notes |
| --- | --- | --- |
| Language | Go 1.21+ | Single static binary per platform. |
| TUI | Bubble Tea, Lip Gloss, Bubbles, Glamour, Harmonica | Bubble Tea implements the Elm pattern; Lip Gloss provides deterministic styling; Glamour renders previews; Harmonica supplies animations. |
| CLI | Cobra, Viper | Cobra powers command structure (e.g. `noise live`, `noise quick`); Viper manages configuration. |
| Persistence | SQLite (modernc.org/sqlite), local Markdown + YAML | Domain data is stored in Markdown with YAML frontmatter; SQLite tracks metadata, versions, knowledge base, and project stats; schema defined in [`internal/infra/db/schema.go`](../internal/infra/db/schema.go). |
| AI | Ollama + Go client | QuickIdeaAgent uses Ollama models with cached prompts and 2s timeouts. |
| Export | Custom JSON + Markdown formatters | Export services stream Markdown/JSON to disk and wave.sh. |

Dependencies are cataloged in [`docs/tdd_markdown.md`](./tdd_markdown.md) under Technology Stack, along with integration notes for each subsystem.

## 3. System Layers
1. **Presentation** – Bubble Tea models in `internal/ui/` manage screen routing, layout composition, and styling via Lip Gloss.
2. **Application** – Services in `internal/app/` orchestrate editor behaviors, AI agent workflows, export logic, and tempo/theme facilities.
3. **Domain** – Structs in `internal/domain/` define Songs, Sections, Lines, QualityScore, and related entities.
4. **Infrastructure** – Packages under `internal/infra/` provide persistence (`db` with repository patterns), file parsing, Ollama clients, and other adapters.
5. **Export** – `internal/export/` formats captured data for Markdown and pattern JSON, ensuring wave.sh interoperability.

The layering mirrors the diagrams in [`PROJECT_SNAPSHOT.md`](../PROJECT_SNAPSHOT.md) and keeps terminal rendering decoupled from persistence.

## 4. Data Models & Stores
Primary structs live in [`internal/domain/song.go`](../internal/domain/song.go):

- `Song`, `SongMetadata`, `Section`, `Line`, `QualityScore`, `Version`, `WritingStats`, and `KnowledgeCard`.
- Tags, structure, and timestamps are maintained in YAML frontmatter per song file.
- SQLite schema in [`internal/infra/db/schema.go`](../internal/infra/db/schema.go) tracks songs, versions, writing stats, knowledge base vectors, projects, and collaboration placeholders.
- Repository interfaces (implemented in `internal/infra/db`) persist domain entities and enforce foreign key relationships.
- Export serialization reuses the same models, preserving parity between Markdown, JSON, and the SQLite record.

## 5. UI Architecture
Key UI components (referenced in [`docs/design_system_reference.md`](./design_system_reference.md)):

- **Split-pane editor** – Markdown textarea (Bubble Tea’s textarea bubble) plus Glamour preview in `internal/ui/editor`.
- **Theme-aware status bar** and **mode indicator** – reflect scratch mode, word counts, AI status.
- **Live mode** UI – pattern editor and visualization panes rendered in Bubble Tea models under `internal/app/live` and `internal/ui/editor/pattern_view`.
- **Quick chord picker & BPM tapper** – modal components in `internal/ui/quick_chord_picker.go` and `internal/ui/bpm_tapper.go`.
- **TUI styling** – Achieved through Lip Gloss gradients, adaptive colors, and layout helpers documented in the design system reference.

The Elm pattern ensures every screen is a Bubble Tea model with `Init`, `Update`, and `View` functions.

## 6. AI Integration
The QuickIdeaAgent consolidates unstick, spark, tweak, and check flows into a single pipeline:

- Architecture, prompts, and response parsing are described in [`docs/context_aware_ai_guide.md`](./context_aware_ai_guide.md) and the Week 3 section of the enhancement plan.
- The agent uses Ollama via a lightweight Go client, caches system prompts, and enforces sub-2s timeouts with 3B models to maintain responsiveness.
- Context-aware detection toggles prompts based on lyric vs. pattern content, so suggestions stay grounded without cloud calls.

## 7. Theme & Styling
[`docs/design_system_reference.md`](./design_system_reference.md) defines the core palette, Lip Gloss styles, animation patterns, and component guidance:

- 12 named themes registered in `internal/theme/` with runtime switching (`Ctrl+Shift+T`).
- Custom theme loader reads TOML presets from `~/.config/noise/themes/`.
- Theme manager enforces WCAG contrast; `internal/theme/validator.go` provides checks.
- Status bars, modals, and editor panes all rehydrate theme tokens to keep CLI and documentation visuals aligned.

## 8. Plugin Interface
[`internal/plugins/README.md`](../internal/plugins/README.md) documents plugin lifecycle:

- `Plugin` interface exposes metadata, initialize/cleanup, enable/disable hooks.
- Capabilities include export formats, editor tools, menu extensions, and AI hooks.
- Plugin manager validates manifests from approved directories (`~/.noise.sh/plugins/`, project-local `plugins/`, internal examples).
- Security layer enforces path validation, size limits, and capability declarations before registration.

## 9. Export & File I/O
Export services turn captured ideas into shareable artifacts per [`docs/export_formats.md`](./export_formats.md):

- Markdown export preserves sections, chords, metadata headers.
- Plain text strips formatting for raw lyric sharing.
- ChordPro (and JSON pattern exports) keep tempo, chord progressions, and metadata aligned for wave.sh ingestion.
- CLI commands (`noise export pattern`, `noise export draft`) call handlers in `internal/export/`, which use domain structs to serialize content.

## 10. Performance & Observability
Targets sourced from [`PROJECT_SNAPSHOT.md`](../PROJECT_SNAPSHOT.md) and the enhancement plan:

- **AI latency** < 2 seconds; prompts cached and 3B models used.
- **Pattern parsing** < 100 ms; live mode response stays instant through incremental parsing.
- **Theme switching** < 50 ms; precomputed palettes avoid runtime slowness.
- **Logging** – lightweight timing logs collected around AI calls; error surfaces bubble up to status bar notifications instead of silent failures.

## 11. Security & Error Handling
Guidance summarized from [`docs/enhancement_risk_assessment.md`](./enhancement_risk_assessment.md):

- **Local-first privacy** – No telemetry or cloud storage; optional Ollama downloads only.
- **Dependency isolation** – Single binary distribution, no CGO dependencies beyond sqlite/oto as pure Go.
- **Error surfacing** – Application layer returns domain-specific errors that UI renders as banner notifications; CLI subcommands guard invalid inputs before execution.
- **Fallback behavior** – AI timeouts fall back to helpful messaging; theme manager defaults to Amber Night on invalid config; export commands validate targets before writing.
