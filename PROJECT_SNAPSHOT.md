# LyricForge Project Snapshot

## Project Overview

LyricForge is an AI-powered terminal user interface (TUI) application that transforms the songwriting process by combining professional songwriting pedagogy with beautiful visual design and intelligent AI assistance. It empowers songwriters to create authentic, human-quality lyrics while maintaining complete privacy and ownership of their creative work.

The application is designed as a **local-first, privacy-focused** tool that brings together:
- Expert songwriting techniques from renowned instructors (Pat Pattison, Andrea Stolpe, Jason Blume)
- Beautiful terminal interface using the Charm Bracelet ecosystem
- Strategic AI integration that enhances rather than replaces human creativity
- Professional-grade output suitable for portfolio showcase

## Current Development Status

**Status: Planning/Documentation Phase**

The project is currently in the foundational planning and documentation stage. All core project documents have been created and approved:

- ✅ **Product Requirements Document (PRD)** - Comprehensive feature specifications and success criteria
- ✅ **Technical Design Document (TDD)** - Detailed architecture and implementation strategy
- ✅ **Development Roadmap** - 13-week phased implementation plan
- ✅ **Design System Reference** - Complete styling and UI component specifications

**Current Progress:**
- Project structure and documentation: **100% complete**
- Development environment setup: **Ready for implementation**
- Core architectural decisions: **Finalized**
- Technology stack: **Selected and validated**

## Technical Architecture

### Core Technology Stack

**Primary Framework & UI:**
- **Language:** Go 1.21+
- **TUI Framework:** Bubble Tea (Elm Architecture pattern)
- **Styling:** Lipgloss (CSS-like terminal styling with gradients and animations)
- **Components:** Bubbles (pre-built textarea, viewport, list components)
- **Forms:** Huh (accessible forms with validation)
- **Animations:** Harmonica (spring-physics animations)

**Music & Audio:**
- **Music Theory:** go-music-theory (scales, chords, Circle of Fifths)
- **Rhyme/Syllable:** pronouncing (CMU dictionary integration)
- **Audio Feedback:** gen2brain/beeep (cross-platform audio cues)
- **MIDI Support:** gomidi/midi v2 (real-time MIDI input)
- **Synthesis:** DylanMeeus/GoAudio (chord playback and tone generation)

**AI & NLP:**
- **LLM Runtime:** Ollama (local model hosting)
- **Vector Database:** ChromaDB (RAG knowledge base)
- **Embeddings:** nomic-embed-text (semantic search)
- **Primary Models:** Qwen 2.5 7B + Llama 3.1 8B (creativity + structured thinking)

**Data & Persistence:**
- **Database:** SQLite (version history, statistics, knowledge base)
- **File Format:** Markdown with YAML frontmatter (Git-friendly)
- **Git Integration:** go-git (automatic commits and version control)

### Architecture Philosophy

The system follows a **modular monolith** architecture with clear separation of concerns:

```
┌────────────────────────────────────────────────────────────────────┐
│                      PRESENTATION LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ Main Router  │──│ Screen Stack │──│ Bubble Tea Components    │ │
│  └──────────────┘  └──────────────┘  │ • Editor  • File Browser │ │
│                                       │ • Theory  • AI Panel     │ │
│                                       └──────────────────────────┘ │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐  │
│  │ Editor Service  │  │ Theory Service  │  │  AI Orchestrator │  │
│  └─────────────────┘  └─────────────────┘  └──────────────────┘  │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                       DOMAIN LAYER                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────────┐  │
│  │ Song Model   │  │ Project      │  │ Knowledge Base         │  │
│  └──────────────┘  └──────────────┘  └────────────────────────┘  │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ File System  │  │ SQLite DB    │  │ External Integrations    │ │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

## Project Structure

The planned project structure follows Go best practices with clear separation of concerns:

```
lyricforge/
├── cmd/
│   └── lyricforge/
│       └── main.go              # Entry point
├── internal/
│   ├── app/                     # Application layer
│   │   ├── editor/              # Editor service
│   │   ├── theory/              # Music theory service
│   │   └── ai/                  # AI orchestration
│   ├── domain/                  # Domain models
│   │   ├── song.go
│   │   ├── section.go
│   │   └── quality.go
│   ├── infra/                   # Infrastructure
│   │   ├── db/                  # SQLite
│   │   ├── files/               # File I/O
│   │   ├── git/                 # Git operations
│   │   └── ollama/              # Ollama client
│   └── ui/                      # TUI components
│       ├── root.go              # Root model/router
│       ├── editor/              # Editor screen
│       ├── theory/              # Theory tools screen
│       ├── audio/               # Audio tools screen
│       ├── manager/             # Project manager screen
│       ├── settings/            # Settings screen
│       └── common/              # Shared components
├── pkg/                         # Public packages
│   ├── kb/                      # Knowledge base
│   └── prompts/                 # AI prompts
├── scripts/                     # Build/dev scripts
│   ├── build.sh
│   ├── test.sh
│   └── install-models.sh
├── testdata/                    # Test fixtures
│   ├── songs/
│   └── golden/                  # Golden file tests
├── docs/                        # Documentation
│   ├── architecture.md
│   ├── api.md
│   └── user-guide.md
├── .goreleaser.yaml             # Release config
├── .golangci.yaml               # Linter config
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Documentation Status

**Complete Documentation Available:**

1. **`README.md`** - Project overview and setup instructions
2. **`docs/prd_markdown.md`** - Product Requirements Document (287 lines)
   - Executive summary and product vision
   - Target user personas and use cases
   - Detailed feature specifications across 5 categories
   - Success metrics and KPIs
   - Technical constraints and portfolio requirements

3. **`docs/tdd_markdown.md`** - Technical Design Document (828 lines)
   - System architecture overview with layer diagrams
   - Complete technology stack specifications
   - Database design and data models
   - AI agent architecture with 6-stage pipeline
   - Component architecture and UI flows
   - API specifications and testing strategy

4. **`docs/dev_roadmap.md`** - Development Roadmap (1,540 lines)
   - 13-week phased implementation plan
   - Day-by-day development guide with code examples
   - Dependency installation instructions
   - Detailed implementation roadmap for each phase

5. **`docs/design_system_reference.md`** - Design System (960 lines)
   - "Midnight Jazz" theme specification
   - Complete color palette and usage guidelines
   - Lipgloss style definitions and component library
   - Typography system and animation patterns
   - Layout system and best practices

**Documentation Quality:** **Excellent** - All documents are comprehensive, technically detailed, and ready for implementation.

## Development Roadmap

### 13-Week Development Plan

**Phase 1: Foundation (Week 1-2)** - Core infrastructure
- Day 1-2: Core data models and domain entities
- Day 3-4: File I/O and markdown parsing
- Day 5-7: SQLite database and repository pattern

**Phase 2: Editor & UI (Week 3-4)** - User interface
- Day 8-10: Basic Bubble Tea editor with split-pane
- Day 11-12: Beautiful Lipgloss styling and animations
- Day 13-14: Advanced editor features (auto-save, shortcuts)

**Phase 3: Music Theory Tools (Week 5-6)** - Creative assistance
- Day 15-17: Circle of Fifths interactive visualization
- Day 18-20: Rhyme dictionary and syllable counter
- Day 21-23: Chord progressions and scale reference

**Phase 4: Audio Features (Week 7)** - Interactive elements
- Day 24-26: Visual metronome with beat animation
- Day 27-28: Chord playback synthesis
- Day 29-30: MIDI input capture

**Phase 5: AI Integration (Week 8-10)** - Intelligent assistance
- Day 31-33: Ollama client and RAG knowledge base
- Day 34-37: Multi-agent pipeline implementation
- Day 38-40: Streaming UI and real-time suggestions

**Phase 6: Project Management (Week 11)** - Professional features
- Day 41-43: Version history and Git integration
- Day 44-46: Statistics dashboard and analytics
- Day 47-48: Export formats and file management

**Phase 7: Polish & Release (Week 12-13)** - Production readiness
- Day 49-52: Comprehensive testing (>70% coverage)
- Day 53-55: Performance optimization and documentation
- Day 56-57: Demo materials and release preparation

## Key Features

### 🎨 Beautiful Writing Environment
- **Split-Pane Editor:** Live markdown editing with real-time Glamour preview
- **YAML Frontmatter:** Rich song metadata (title, artist, key, tempo, structure)
- **Syntax Highlighting:** Visual distinction for verse/chorus/bridge sections
- **Distraction-Free Mode:** Fullscreen, centered text with focus animations
- **Auto-Save:** Every 30 seconds with subtle visual feedback

### 🎵 Music Theory Tools
- **Circle of Fifths:** Interactive visualization with color-coded relationships
- **Chord Progressions:** Common patterns with visual timeline and playback
- **Rhyme Dictionary:** Perfect, slant, and internal rhymes with syllable counts
- **Syllable Counter:** Real-time per-line analysis for meter matching
- **Scale Reference:** Visual scale displays for key selection

### 🎹 Audio & Interactive Elements
- **Visual Metronome:** Animated beat indicator with spring physics
- **Chord Playback:** Synthesized audio for chord progressions
- **MIDI Input:** Real-time chord capture from keyboards/controllers
- **Audio Notifications:** Gentle feedback for saves and completions

### 🤖 AI-Powered Assistance (Strategic Integration)
- **Brainstorming Assistant:** 5-10 specific angles avoiding generic concepts
- **Object Writing Engine:** Pat Pattison's 7-sense exercise implementation
- **Lyric Refinement:** Concrete imagery replacement and verb strengthening
- **Prosody Critique:** Stress-syllable alignment analysis
- **Cliché Detection:** Anti-cliché lists with fresh alternative suggestions
- **Quality Scoring:** 7-dimension rubric with actionable feedback

### 💾 Project Management
- **File Browser:** Fuzzy search across song library
- **Version History:** 30-day rolling auto-save with manual milestones
- **Git Integration:** Auto-initialization and commit shortcuts
- **Export Formats:** Markdown, PDF, HTML, plain text
- **Statistics Dashboard:** Writing streaks and productivity insights

## Next Steps

### Immediate Actions (Week 1)

1. **Environment Setup**
   ```bash
   # Create project structure
   mkdir lyricforge && cd lyricforge
   git init

   # Initialize Go module
   go mod init github.com/puente-labs/lyricforge

   # Install core dependencies (Bubble Tea ecosystem)
   go get github.com/charmbracelet/bubbletea@latest
   go get github.com/charmbracelet/lipgloss@latest
   go get github.com/charmbracelet/bubbles@latest
   go get github.com/charmbracelet/glamour@latest
   go get github.com/charmbracelet/huh@latest
   go get github.com/charmbracelet/harmonica@latest
   ```

2. **Create Initial Code Structure**
   - Implement `cmd/lyricforge/main.go` entry point
   - Create `internal/ui/root.go` with basic Bubble Tea setup
   - Set up `internal/domain/` with core data models
   - Initialize `internal/infra/` for file operations

3. **Development Environment**
   - Configure VS Code with Go extensions
   - Set up `golangci-lint` for code quality
   - Create initial test structure
   - Establish Git workflow and commit patterns

### Success Criteria for Phase 1
- ✅ Project compiles and runs without errors
- ✅ Basic TUI displays (splash screen, menu navigation)
- ✅ Core data models implemented with validation
- ✅ File I/O operations functional (load/save markdown)
- ✅ SQLite database operational with schema

### Portfolio Impact
This project is designed as a **portfolio showcase piece** with:
- **GitHub Star Potential:** 500+ stars within 6 months
- **Visual Impact:** Beautiful terminal interface with gradients and animations
- **Technical Excellence:** Production-grade Go with comprehensive testing
- **Innovation:** Novel combination of AI assistance with creative workflows
- **Documentation Quality:** Comprehensive technical documentation

---

**Document Status:** Complete and ready for implementation
**Last Updated:** October 17, 2025
**Version:** 1.0