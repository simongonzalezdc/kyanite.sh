# Product Requirements Document (PRD)
## noise.sh: AI-Powered Songwriting TUI

**Version:** 1.0  
**Date:** October 17, 2025  
**Document Owner:** Simon (Puente Labs)  
**Status:** Draft for Development

---

## Executive Summary

### Product Vision

noise.sh is a **local-first, AI-powered terminal user interface (TUI) application** that transforms the songwriting process by combining professional songwriting pedagogy with beautiful visual design and intelligent AI assistance. It empowers songwriters to create authentic, human-quality lyrics while maintaining complete privacy and ownership of their creative work.

### Key Differentiators

- ✅ **100% Local & Private:** No cloud, no data sharing, complete creative sovereignty
- 🎓 **Expert Knowledge Integration:** Encodes proven techniques from Pat Pattison, Andrea Stolpe, Jason Blume
- 🎨 **Beautiful TUI Design:** Showcases Charm Bracelet ecosystem with gradients, spinners, animations
- 🚫 **Anti-"AI Slop":** Systematic guardrails preventing generic, clichéd output
- 🔄 **Version Control Native:** Plain text format enables Git workflow from day one

### Success Criteria

- Generate lyrics scoring **70+** on quality rubric (specificity, originality, prosody)
- Sub-**2-second** response time for AI suggestions
- **Zero cliché flags** in polished output
- Portfolio-worthy showcase of technical & design skills
- GitHub star potential **500+** within 6 months of release

---

## Target User Personas

### Primary: The Technical Songwriter

**Demographics:**
- Age: 25-45
- Occupation: Developer/Creative professional
- Tech-savvy, CLI-comfortable
- Active on GitHub, uses terminal daily

**Pain Points:**
- GUI apps too distracting (notifications, chrome, visual clutter)
- Writer's block frustration, need creative spark
- No privacy in cloud tools (Google Docs, Notion)
- Can't version control lyrics like code
- Existing AI tools produce generic "slop"

**Goals:**
- Distraction-free creative environment
- AI assistance that understands songcraft, not just rhyming
- Private creative work that stays local
- Version control for lyric iterations
- Professional-grade output

**User Story:**
> "I want a distraction-free environment where I can write lyrics with AI assistance that actually understands songcraft, not just rhyming words. I need my creative work to stay private and version-controlled like my code."

### Secondary: The Portfolio Evaluator

**Demographics:**
- Hiring managers, recruiters
- Open source contributors
- Technical community members
- Dev tool enthusiasts

**Evaluation Criteria:**
- Code quality & architecture
- UI/UX polish and creativity
- Problem-solving approach
- Documentation quality
- Novel use of technology

### Tertiary: Terminal Enthusiasts

**Communities:**
- r/commandline, r/unixporn, r/vim
- HackerNews, Lobsters
- Terminal tool collectors

**Attraction Factors:**
- Beautiful design in terminal
- Novel creative use case
- Technical excellence
- Open source contribution opportunity

---

## Core Features & Functionality

### 🎨 Feature Category 1: Beautiful Writing Environment

| Feature | Description | Technology |
|---------|-------------|------------|
| **Split-Pane Editor** | Live markdown editing with real-time rendered preview | Bubble Tea textarea + Glamour |
| **YAML Frontmatter** | Song metadata (title, artist, key, tempo, tags, structure) | gopkg.in/yaml.v2 |
| **Syntax Highlighting** | Verse/Chorus/Bridge sections visually distinct | Lipgloss styling |
| **Distraction-Free Mode** | Fullscreen, centered text, focus mode with visual breathing | Lipgloss layouts + Harmonica |
| **Auto-Save** | Every 30 seconds with subtle notification | File I/O + Bubble Tea commands |
| **Keyboard Shortcuts** | Vim-inspired navigation, quick actions | Native Bubble Tea keybindings |

### 🎵 Feature Category 2: Music Theory Tools

| Feature | Description | Technology |
|---------|-------------|------------|
| **Circle of Fifths** | Interactive visualization with color-coded relationships | Lipgloss + Unicode + go-music-theory |
| **Chord Progressions** | Suggest common progressions, display with visual timeline | go-music-theory + animated highlights |
| **Rhyme Dictionary** | Perfect, slant, internal rhymes with syllable counts | pronouncing + CMU dict |
| **Syllable Counter** | Real-time per-line syllable count for meter matching | pronouncing library |
| **Scale Reference** | Visual scale displays for key selection | go-music-theory |
| **Chord Insertion** | Quick insertion of chord symbols above lyrics | Custom parser |

### 🎹 Feature Category 3: Audio & Interactive Elements

| Feature | Description | Technology |
|---------|-------------|------------|
| **Visual Metronome** | Animated beat indicator with color flash and sound | gen2brain/beeep + Harmonica |
| **Chord Playback** | Hear chord progressions with tone generation | DylanMeeus/GoAudio synthesis |
| **MIDI Input** | Capture chord progressions from keyboard/controller | gomidi/midi v2 (midicatdrv) |
| **Audio Notifications** | Gentle beeps for save, completion, achievements | gen2brain/beeep |
| **Reference Audio** | Play demo melodies for inspiration (future) | gopxl/beep |
| **Tempo Trainer** | Gradual BPM increases for practice | Custom timing + audio |

### 🤖 Feature Category 4: AI-Powered Assistance (Strategic Integration)

| Feature | Value Proposition | Implementation |
|---------|-------------------|----------------|
| **Brainstorming Assistant** | Generate 5-10 specific angles on theme, avoiding generic concepts | Ollama API + RAG KB query |
| **Object Writing Engine** | Pat Pattison's 7-sense exercise for concrete imagery generation | Specialized prompt + streaming UI |
| **Lyric Refinement** | Replace abstractions with concrete images, strengthen verbs | Multi-pass agent with diff view |
| **Prosody Critique** | Analyze stress-syllable alignment, suggest meter improvements | Hybrid: programmatic check + AI explanation |
| **Cliché Detection** | Flag overused phrases, suggest fresh alternatives | Blocklist + AI synonym generation |
| **Real-Time Suggestions** | Inline completions during typing (pause-based trigger) | Streaming completion + inline UI |
| **Quality Scoring** | 7-dimension rubric with actionable feedback | Critic agent + visualization |

**AI Design Philosophy:**

AI features are strategically deployed where they provide **irreplaceable value**: breaking writer's block, encoding expert knowledge, and automating tedious checks. The AI never *replaces* the songwriter—it amplifies their craft. All AI outputs are presented as suggestions with human approval required.

### 💾 Feature Category 5: Project Management

| Feature | Description | Technology |
|---------|-------------|------------|
| **File Browser** | Navigate song library with fuzzy search | Bubbles list + filtering |
| **Version History** | 30-day rolling auto-save + manual milestones | SQLite + visual timeline |
| **Git Integration** | Auto-init repo on first save, commit shortcuts | go-git library |
| **Export Formats** | Markdown, PDF, HTML, plain text | Glamour + Pandoc wrapper |
| **Statistics Dashboard** | Writing streaks, song count, productivity charts | SQLite queries + Lipgloss charts |
| **Tags & Search** | Organize songs by tags, full-text search | SQLite FTS5 |

---

## Success Metrics & KPIs

### Product Quality Metrics

| Metric | Target |
|--------|--------|
| Lyric Quality Score | ≥70/100 |
| Cliché Count | <5 per song |
| AI Response Time | <2 seconds |
| Prosody Flags | 0 (clean pass) |
| Imagery Ratio | >0.6 |

### Portfolio Impact Metrics

| Metric | Target |
|--------|--------|
| GitHub Stars (6mo) | 500+ |
| Reddit Upvotes | 100+ on showcase post |
| Documentation Quality | Excellent (comprehensive) |
| Code Test Coverage | >70% |
| Demo Video Views | 1000+ (YouTube) |

### User Experience Metrics

| Metric | Target |
|--------|--------|
| First-Use Success | Complete 1 verse in <10 min |
| Load Time | <1 second |
| Frame Rate | 60 FPS |
| Keyboard Shortcuts | 100% features accessible |
| Cross-Platform | Win/Mac/Linux identical |

### Technical Excellence

| Metric | Target |
|--------|--------|
| Binary Size | <50MB |
| Memory Usage | <200MB (without models) |
| Go Report Card | A+ grade |
| Linting | Zero golangci-lint issues |
| Dependencies | Minimal CVEs |

---

## Technical Constraints & Considerations

### Critical Constraints

- **Local-Only Architecture:** All processing must happen on user's machine. No cloud dependencies except optional Ollama model downloads.
- **Privacy First:** No telemetry, no analytics, no data collection. User's creative work never leaves their device.
- **Cross-Platform from Day 1:** Must work identically on Windows, macOS, Linux. No platform-specific code without abstraction.
- **Single Binary Distribution:** No installation scripts, no dependencies to install separately (except Ollama).
- **Ollama Required:** Users must have Ollama installed and running for AI features. Graceful degradation if not available.

### Development Constraints

- **Solo Developer + AI:** Architecture must be maintainable by one person with AI coding assistance
- **Clear Separation of Concerns:** Modular design enabling parallel feature development
- **Comprehensive Testing:** Unit tests, integration tests, golden file tests for TUI
- **Documentation as Code:** All features documented inline, README auto-generated

### Performance Requirements

- **Responsive UI:** 60 FPS rendering, no dropped frames during typing
- **Fast Startup:** <1 second from launch to ready state
- **Low Latency:** AI streaming must start within 500ms of request
- **Efficient Memory:** Total memory footprint <200MB without models loaded
- **Battery Friendly:** No busy-waiting, no unnecessary CPU spinning

### Portfolio Requirements

- **Visual Impact:** Must showcase advanced Charm Bracelet techniques (gradients, animations, layouts)
- **Code Quality:** Production-grade Go with idiomatic patterns, proper error handling
- **Documentation:** README, architecture docs, API docs, contribution guide
- **Demo Materials:** GIF demos, video walkthrough, blog post explaining architecture
- **Open Source Ready:** MIT license, clear contribution guidelines, issue templates

---

## Out of Scope (Not Version 1.0)

- Cloud sync or collaboration features
- Mobile apps or web version
- DAW integration or VST plugins
- Full music composition (MIDI sequencer)
- Voice recording or audio editing
- Social features or community integration
- Paid tiers or commercial licensing
- Translation to other languages (English only for v1)

---

## Design Principles

### Visual Design

- **Colorful & Beautiful:** Showcase Charm Bracelet's full potential with gradients, animations, smooth transitions
- **Not Minimalist:** Use color, visual hierarchy, and animation generously
- **Accessibility:** Support light/dark modes, readable fonts, clear contrast
- **Delightful Details:** Spinners, progress bars, success animations

### Interaction Design

- **Keyboard-First:** Every feature accessible via keyboard shortcuts
- **Progressive Disclosure:** Advanced features hidden until needed
- **Gentle Onboarding:** First-time tutorial, contextual help
- **Fast Feedback:** Instant visual response to all actions

### Information Architecture

- **Clear Mental Model:** Logical grouping of features
- **Consistent Navigation:** Same patterns across all screens
- **Escape Hatches:** Easy way to cancel/go back from any flow
- **Status Awareness:** Always know where you are, what's happening

---

## Document History

**Version 1.0** - October 17, 2025 - Initial draft for development

**Next Steps:**
1. Review and approve PRD
2. Create Technical Design Document
3. Set up development environment
4. Begin implementation

---

**Document Owner:** Simon (Puente Labs)  
**Status:** Ready for TDD creation  
**Next Review:** Upon completion of Technical Design Document
