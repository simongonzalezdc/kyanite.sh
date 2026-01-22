# noise.sh Project Snapshot

## Project Overview

noise.sh is an AI-powered terminal user interface (TUI) application that transforms the songwriting process by combining professional songwriting pedagogy with beautiful visual design and intelligent AI assistance. It empowers songwriters to create authentic, human-quality lyrics while maintaining complete privacy and ownership of their creative work.

The application is designed as a **local-first, privacy-focused** tool that brings together:
- Expert songwriting techniques from renowned instructors (Pat Pattison, Andrea Stolpe, Jason Blume)
- Beautiful terminal interface using the Charm Bracelet ecosystem
- Strategic AI integration that enhances rather than replaces human creativity
- Professional-grade output suitable for portfolio showcase

## Current Development Status

**Status: Active Development**

The project has progressed beyond planning into active implementation. Core infrastructure is functional with ongoing feature development and quality improvements:

- âœ… **Product Requirements Document (PRD)** - Comprehensive feature specifications and success criteria
- âœ… **Technical Design Document (TDD)** - Detailed architecture and implementation strategy
- âœ… **Development Roadmap** - 13-week phased implementation plan
- âœ… **Design System Reference** - Complete styling and UI component specifications
- âœ… **Core TUI Framework** - Bubble Tea application with responsive layouts
- âœ… **Theme System** - Multiple themes with gradient support
- âœ… **AI Integration** - QuickIdeaAgent with retry logic and knowledge base
- âœ… **Database Layer** - SQLite with WAL mode and performance optimizations
- âœ… **Music Theory** - Prosody engine, harmony library, and rhyme tools

**Current Progress:**
- Project structure and documentation: **100% complete**
- Core TUI implementation: **Functional**
- AI agent framework: **Implemented with fallback support**
- Database persistence: **Operational with performance tuning**
- Music theory tools: **Core features complete**
- Test coverage: **Growing, with comprehensive test suites**

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
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                      PRESENTATION LAYER                             â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â” â”‚
â”‚  â”‚ Main Router  â”‚â”€â”€â”‚ Screen Stack â”‚â”€â”€â”‚ Bubble Tea Components    â”‚ â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â”‚ â€¢ Editor  â€¢ File Browser â”‚ â”‚
â”‚                                       â”‚ â€¢ Theory  â€¢ AI Panel     â”‚ â”‚
â”‚                                       â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜ â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
                         â”‚
                         â–¼
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                      APPLICATION LAYER                              â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”‚
â”‚  â”‚ Editor Service  â”‚  â”‚ Theory Service  â”‚  â”‚  AI Orchestrator â”‚  â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
                         â”‚
                         â–¼
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                       DOMAIN LAYER                                  â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”‚
â”‚  â”‚ Song Model   â”‚  â”‚ Project      â”‚  â”‚ Knowledge Base         â”‚  â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”¬â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
                         â”‚
                         â–¼
â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”
â”‚                    INFRASTRUCTURE LAYER                             â”‚
â”‚  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”  â”Œâ”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â” â”‚
â”‚  â”‚ File System  â”‚  â”‚ SQLite DB    â”‚  â”‚ External Integrations    â”‚ â”‚
â”‚  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜  â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜ â”‚
â””â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”€â”˜
```

## Project Structure

The planned project structure follows Go best practices with clear separation of concerns:

```
noise.sh/
â”œâ”€â”€ cmd/
â”‚   â””â”€â”€ noise.sh/
â”‚       â””â”€â”€ main.go              # Entry point
â”œâ”€â”€ internal/
â”‚   â”œâ”€â”€ app/                     # Application layer
â”‚   â”‚   â”œâ”€â”€ editor/              # Editor service
â”‚   â”‚   â”œâ”€â”€ theory/              # Music theory service
â”‚   â”‚   â””â”€â”€ ai/                  # AI orchestration
â”‚   â”œâ”€â”€ domain/                  # Domain models
â”‚   â”‚   â”œâ”€â”€ song.go
â”‚   â”‚   â”œâ”€â”€ section.go
â”‚   â”‚   â””â”€â”€ quality.go
â”‚   â”œâ”€â”€ infra/                   # Infrastructure
â”‚   â”‚   â”œâ”€â”€ db/                  # SQLite
â”‚   â”‚   â”œâ”€â”€ files/               # File I/O
â”‚   â”‚   â”œâ”€â”€ git/                 # Git operations
â”‚   â”‚   â””â”€â”€ ollama/              # Ollama client
â”‚   â””â”€â”€ ui/                      # TUI components
â”‚       â”œâ”€â”€ root.go              # Root model/router
â”‚       â”œâ”€â”€ editor/              # Editor screen
â”‚       â”œâ”€â”€ theory/              # Theory tools screen
â”‚       â”œâ”€â”€ audio/               # Audio tools screen
â”‚       â”œâ”€â”€ manager/             # Project manager screen
â”‚       â”œâ”€â”€ settings/            # Settings screen
â”‚       â””â”€â”€ common/              # Shared components
â”œâ”€â”€ pkg/                         # Public packages
â”‚   â”œâ”€â”€ kb/                      # Knowledge base
â”‚   â””â”€â”€ prompts/                 # AI prompts
â”œâ”€â”€ scripts/                     # Build/dev scripts
â”‚   â”œâ”€â”€ build.sh
â”‚   â”œâ”€â”€ test.sh
â”‚   â””â”€â”€ install-models.sh
â”œâ”€â”€ testdata/                    # Test fixtures
â”‚   â”œâ”€â”€ songs/
â”‚   â””â”€â”€ golden/                  # Golden file tests
â”œâ”€â”€ docs/                        # Documentation
â”‚   â”œâ”€â”€ architecture.md
â”‚   â”œâ”€â”€ api.md
â”‚   â””â”€â”€ user-guide.md
â”œâ”€â”€ .goreleaser.yaml             # Release config
â”œâ”€â”€ .golangci.yaml               # Linter config
â”œâ”€â”€ go.mod
â”œâ”€â”€ go.sum
â”œâ”€â”€ LICENSE
â””â”€â”€ README.md
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

### ðŸŽ¨ Beautiful Writing Environment
- **Split-Pane Editor:** Live markdown editing with real-time Glamour preview
- **YAML Frontmatter:** Rich song metadata (title, artist, key, tempo, structure)
- **Syntax Highlighting:** Visual distinction for verse/chorus/bridge sections
- **Distraction-Free Mode:** Fullscreen, centered text with focus animations
- **Auto-Save:** Every 30 seconds with subtle visual feedback

### ðŸŽµ Music Theory Tools
- **Circle of Fifths:** Interactive visualization with color-coded relationships
- **Chord Progressions:** Common patterns with visual timeline and playback
- **Rhyme Dictionary:** Perfect, slant, and internal rhymes with syllable counts
- **Syllable Counter:** Real-time per-line analysis for meter matching
- **Scale Reference:** Visual scale displays for key selection

### ðŸŽ¹ Audio & Interactive Elements
- **Visual Metronome:** Animated beat indicator with spring physics
- **Chord Playback:** Synthesized audio for chord progressions
- **MIDI Input:** Real-time chord capture from keyboards/controllers
- **Audio Notifications:** Gentle feedback for saves and completions

### ðŸ¤– AI-Powered Assistance (Strategic Integration)
- **Brainstorming Assistant:** 5-10 specific angles avoiding generic concepts
- **Object Writing Engine:** Pat Pattison's 7-sense exercise implementation
- **Lyric Refinement:** Concrete imagery replacement and verb strengthening
- **Prosody Critique:** Stress-syllable alignment analysis
- **ClichÃ© Detection:** Anti-clichÃ© lists with fresh alternative suggestions
- **Quality Scoring:** 7-dimension rubric with actionable feedback

### ðŸ’¾ Project Management
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
   mkdir noise.sh && cd noise.sh
   git init

   # Initialize Go module
   go mod init github.com/Kyanite/noise.sh

   # Install core dependencies (Bubble Tea ecosystem)
   go get github.com/charmbracelet/bubbletea@latest
   go get github.com/charmbracelet/lipgloss@latest
   go get github.com/charmbracelet/bubbles@latest
   go get github.com/charmbracelet/glamour@latest
   go get github.com/charmbracelet/huh@latest
   go get github.com/charmbracelet/harmonica@latest
   ```

2. **Create Initial Code Structure**
   - Implement `cmd/noise.sh/main.go` entry point
   - Create `internal/ui/root.go` with basic Bubble Tea setup
   - Set up `internal/domain/` with core data models
   - Initialize `internal/infra/` for file operations

3. **Development Environment**
   - Configure VS Code with Go extensions
   - Set up `golangci-lint` for code quality
   - Create initial test structure
   - Establish Git workflow and commit patterns

### Success Criteria for Phase 1
- âœ… Project compiles and runs without errors
- âœ… Basic TUI displays (splash screen, menu navigation)
- âœ… Core data models implemented with validation
- âœ… File I/O operations functional (load/save markdown)
- âœ… SQLite database operational with schema

### Portfolio Impact
This project is designed as a **portfolio showcase piece** with:
- **GitHub Star Potential:** 500+ stars within 6 months
- **Visual Impact:** Beautiful terminal interface with gradients and animations
- **Technical Excellence:** Production-grade Go with comprehensive testing
- **Innovation:** Novel combination of AI assistance with creative workflows
- **Documentation Quality:** Comprehensive technical documentation

---

**Document Status:** Active Development
**Last Updated:** January 22, 2026
**Version:** 1.1
