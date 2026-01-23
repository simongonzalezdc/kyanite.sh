# noise.sh — rapid-capture songwriting studio

## Overview
noise.sh is a local-first, AI-assisted TUI that turns inspiration into captured lyrics in under a minute. The terminal interface combines embedded songwriting pedagogy, deterministic AI prompts, and a modular monolith architecture so contributors and artists can create, iterate, and export without ever leaving their machine or leaking creative work.

- **Instant idea-to-lyrics workflow** — launch straight into scratch mode with live autosave.
- **Voice-to-text dictation** — press `Ctrl+D` to dictate lyrics hands-free using local whisper.cpp (auto-downloads on first use).
- **Right-sized AI** — the QuickIdeaAgent delivers unstick, spark, tweak, and check flows in <2 seconds while respecting context.
- **Theme + export pipeline** — ten curated themes, custom theme loading, and JSON/Markdown exports ready for downstream tooling.
- **PWA companion ready** — local sync server for capturing ideas from your phone (voice memos, photos, tap tempo).
- **Local-only privacy** — all processing happens on-device; Ollama models are optional and configurable.

For the strategic roadmap and long-term milestones, see [`docs/roadmap.md`](docs/roadmap.md).

## Key Features
1. **Instant capture** — Split-pane editor, live mode pattern parser, BPM tapper, and chord picker.
2. **Voice dictation** — Push-to-talk transcription with whisper.cpp. Models download automatically on first use.
3. **Smart suggestions** — QuickIdeaAgent modes (unstick, spark, tweak, check) with fast, context-aware responses.
4. **Theme-ready exports** — Markdown + JSON export services deliver wave.sh-compatible payloads and theme-preserving drafts.
5. **Mobile sync** — Embedded HTTP/WebSocket server for companion PWA (voice memos, photos, tempo capture).
6. **Low-latency AI** — Cached prompts, 3B local models, and 2s timeouts keep writing momentum intact.

## Quickstart (15 seconds)
```bash
# Launch scratch mode with instant brainstorming
noise.exe quick "lost my dog"

# Start with default interface
noise.exe

# Open a specific song file
noise.exe song.md
```

### Editor shortcuts
| Action | Keys |
| --- | --- |
| Voice dictation (push-to-talk) | `Ctrl+D` |
| AI unstick / next line | `Alt+G` |
| Spark ideas | `Alt+R` |
| Tweak selected line | `Alt+V` |
| AI check | `Alt+C` |
| Quick chord picker | `Ctrl+F` |
| BPM tapper | `Ctrl+Shift+B` |
| Theme switcher | `Ctrl+Shift+T` |
| Export menu | `Ctrl+E` |
| Save | `Ctrl+S` |

## Installation

### Quick Start (Recommended)
Use the launcher script which handles everything automatically:

```bash
git clone https://github.com/Kyanite/noise.sh.git
cd noise.sh

# macOS/Linux
./launch.sh

# Windows
launch.bat
```

The launcher will:
1. Check for Go installation
2. Download dependencies
3. Build the binary
4. Set up data directories
5. Launch the app

### Launch with PWA Companion
To start both the TUI and the mobile companion PWA:

```bash
./launch-all.sh
```

This starts:
- TUI app in your terminal
- PWA dev server at http://localhost:3000 (access from your phone on the same network)

### Manual Build
If you prefer to build manually:

```bash
git clone https://github.com/Kyanite/noise.sh.git
cd noise.sh
go build -o noise.exe ./cmd/noise
```

### Requirements
- Go 1.25+
- (Optional) Ollama running locally for AI features
- SQLite bundled via `modernc.org/sqlite`
- Voice models download automatically on first use (~142MB for base English model)

## Workflow essentials
1. **Capture** — `noise.exe` or `noise.exe quick` drops you into the split-pane editor with autosave.
2. **Dictate** — Press `Ctrl+D` to start voice dictation. Speak your lyrics, press `Ctrl+D` again to stop and transcribe.
3. **Shape** — Use QuickIdeaAgent shortcuts (`Alt+G`, `Alt+R`, `Alt+V`, `Alt+C`), chord picker (`Ctrl+F`), BPM tapper (`Ctrl+Shift+B`), and theme manager (`Ctrl+Shift+T`) to iterate rapidly.
4. **Sync** — Capture ideas on your phone via the companion PWA. Ideas appear in your inbox for review.
5. **Export** — Press `Ctrl+E` to open the export menu, then select your format (Markdown, Plain Text, ChordPro, JSON, HTML, or PDF).
6. **Extend** — Plugins live under `internal/plugins/`; see [`internal/plugins/README.md`](internal/plugins/README.md) for capabilities.

## Roadmap snapshot
- **Current:** Voice-to-text integration complete, PWA sync infrastructure ready.
- **Near-term:** PWA companion app development, voice settings UI polish.
- **Mid-term:** Theme accessibility QA, QuickIdeaAgent observability, cross-device workflow stitching.
- **Long-term:** i18n expansion, plugin ecosystem growth.

Details and cross-team dependencies are tracked in [`docs/roadmap.md`](docs/roadmap.md).

## Privacy & telemetry
noise.sh is local-first by design:
- No analytics or outbound network traffic.
- AI models run through Ollama locally (optional).
- Exports stay on disk unless you explicitly share them.

## Contributing
1. Review the architecture reference in [`docs/reference/architecture.md`](docs/reference/architecture.md) for layering and design conventions.
2. Submit enhancements or documentation updates via pull requests.
3. Tests: `go test ./...`
4. Markdown lint (manual review): ensure internal links remain relative after doc moves.

## Support
- Issues & feature requests: [GitHub Issues](https://github.com/Kyanite/noise.sh/issues)
- Documentation index: [`docs/`](docs/)
- Archived legacy material: [`docs/archive/`](docs/archive/)

noise.sh keeps creativity private, fast, and repeatable—capture on the spot, ship when ready.
