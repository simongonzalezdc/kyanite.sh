# noise.sh — rapid-capture songwriting studio

## TL;DR
noise.sh is a local-first, AI-assisted TUI that turns inspiration into captured lyrics in under a minute. It combines professional songwriting pedagogy with a high-performance terminal interface.

## Quick Start
```bash
git clone https://github.com/Kyanite/noise.sh.git
cd noise.sh
./launch.sh  # macOS/Linux
# OR
launch.bat   # Windows
```

## Architecture
The system follows a **modular monolith** architecture using the **Elm Architecture** (Bubble Tea) for the TUI and a Next.js PWA for mobile sync.

- **Presentation**: Bubble Tea (TUI) & Next.js (PWA)
- **Application**: Editor, Theory, and AI Services
- **Domain**: Song models, Prosody engine, Knowledge base
- **Infrastructure**: SQLite (WAL mode), whisper.cpp (Local VTT), Ollama (Local LLM)

## Key Features
- **Instant capture** — Split-pane editor with live pattern parsing.
- **Voice dictation** — Push-to-talk transcription with local whisper.cpp.
- **Right-sized AI** — QuickIdeaAgent modes (unstick, spark, tweak, check) in <2s.
- **Mobile sync** — Local network sync with companion PWA for on-the-go capture.
- **Privacy-first** — All processing happens on-device; no cloud dependencies.

## Contributing
1. Review [`docs/reference/architecture.md`](docs/reference/architecture.md) for design conventions.
2. Ensure all changes pass quality gates: `go test ./...`
3. Submit enhancements or documentation updates via pull requests.

## Requirements
- Go 1.25+
- (Optional) Ollama running locally for AI features
- SQLite (bundled)

---
*noise.sh keeps creativity private, fast, and repeatable.*
