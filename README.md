# noise.sh — rapid-capture songwriting studio

## Overview
noise.sh is a local-first, AI-assisted TUI that turns inspiration into captured lyrics in under a minute. The terminal interface combines embedded songwriting pedagogy, deterministic AI prompts, and a modular monolith architecture so contributors and artists can create, iterate, and export without ever leaving their machine or leaking creative work.

- ⚡ **Instant idea-to-lyrics workflow** – launch straight into scratch mode with live autosave.
- 🧠 **Right-sized AI** – the QuickIdeaAgent delivers unstick, spark, tweak, and check flows in <2 seconds while respecting context.
- 🎨 **Theme + export pipeline** – twelve curated themes, custom theme loading, and JSON/Markdown exports ready for downstream tooling.
- 🔒 **Local-only privacy** – all processing happens on-device; Ollama models are optional and configurable.

For the strategic roadmap and long-term milestones, see [`docs/roadmap.md`](docs/roadmap.md).

## Key Features
1. **Instant capture** – Split-pane editor, live mode pattern parser, BPM tapper, and chord picker shipped in Enhancement 6 Weeks 1-2.
2. **Smart suggestions** – QuickIdeaAgent modes (unstick, spark, tweak, check) replace the legacy multi-agent stack with faster, context-aware responses.
3. **Theme-ready exports** – Markdown + JSON export services deliver wave.sh-compatible payloads and theme-preserving drafts.
4. **Low-latency AI** – Cached prompts, 3B local models, and 2s timeouts keep writing momentum intact.

## Quickstart (15 seconds)
```bash
# Launch scratch mode with instant brainstorming
noise.exe quick "lost my dog"

# Start pattern live coding
noise live

# Export current draft to Markdown + JSON bundle
noise export draft --output drafts/lost-my-dog.md
```

### Editor shortcuts
| Action | Keys |
| --- | --- |
| AI unstick / next line | `Ctrl+G` |
| Spark ideas | `Ctrl+R` |
| Tweak selected line | `Ctrl+V` |
| Quick chord picker | `Ctrl+F` |
| BPM tapper | `Ctrl+T` |
| Theme switcher | `Ctrl+Shift+T` |
| Keep scratch draft | `Ctrl+K` |

## Installation
Clone the repository and build the single binary:

```bash
git clone https://github.com/puente-labs/noise.sh.git
cd noise.sh
go build -o noise.exe ./cmd/noise
```

### Requirements
- Go 1.21+
- (Optional) Ollama running locally for AI features
- SQLite bundled via `modernc.org/sqlite`

## Workflow essentials
1. **Capture** – `noise.exe` or `noise.exe quick` drops you into the split-pane editor with autosave.
2. **Shape** – Use QuickIdeaAgent shortcuts, chord picker, BPM tapper, and theme manager to iterate rapidly.
3. **Export** – `noise export pattern` or `noise export draft` generates JSON + Markdown ready for wave.sh or sharing.
4. **Extend** – Plugins live under `internal/plugins/`; see [`internal/plugins/README.md`](internal/plugins/README.md) for capabilities.

## Roadmap snapshot
- 📍 **Week 4 (current):** Documentation unification, README refresh, archival cleanup.
- 🧭 **Mid-term:** Theme accessibility QA, QuickIdeaAgent observability, end-to-end workflow stitching.
- 🚀 **Long-term:** Mobile companion capture app, i18n expansion, plugin ecosystem growth.

Details and cross-team dependencies are tracked in [`docs/roadmap.md`](docs/roadmap.md) and [`docs/Enhancements/enhancement_06_complete (2).md`](docs/Enhancements/enhancement_06_complete%20(2).md).

## Privacy & telemetry
noise.sh is local-first by design:
- No analytics or outbound network traffic.
- AI models run through Ollama locally (optional).
- Exports stay on disk unless you explicitly share them.

## Contributing
1. Review the architecture reference in [`docs/architecture.md`](docs/architecture.md) for layering and design conventions.
2. Submit enhancements or documentation updates via pull requests.
3. Tests: `go test ./...`
4. Markdown lint (manual review): ensure internal links remain relative after doc moves.

## Support
- Issues & feature requests: [GitHub Issues](https://github.com/puente-labs/noise.sh/issues)
- Documentation index: [`docs/`](docs/)
- Archived legacy material: [`docs/archive/`](docs/archive/)

noise.sh keeps creativity private, fast, and repeatable—capture on the spot, ship when ready.
