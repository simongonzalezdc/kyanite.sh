# Technical Design Document (TDD)
## LyricForge: AI-Powered Songwriting TUI

**Version:** 1.0  
**Date:** October 17, 2025  
**Document Owner:** Simon (Puente Labs)  
**Status:** Draft for Development

---

## Table of Contents

1. [System Architecture Overview](#system-architecture-overview)
2. [Technology Stack](#technology-stack)
3. [Database Design & Data Models](#database-design--data-models)
4. [AI Agent Architecture](#ai-agent-architecture)
5. [Component Architecture & UI Flows](#component-architecture--ui-flows)
6. [API Specifications](#api-specifications)
7. [Security & Privacy](#security--privacy)
8. [Performance Optimization](#performance-optimization)
9. [Testing Strategy](#testing-strategy)
10. [Deployment & Distribution](#deployment--distribution)

---

## System Architecture Overview

### Architecture Philosophy

LyricForge follows a **modular monolith** architecture with clear separation between UI, business logic, AI orchestration, and data persistence. The system uses the **Model-Update-View (Elm Architecture)** pattern via Bubble Tea for predictable state management and side-effect handling.

### High-Level Architecture Diagram

```
┌────────────────────────────────────────────────────────────────────┐
│                         User Terminal                               │
│                    (Windows/macOS/Linux)                            │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                      PRESENTATION LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ Main Router  │──│ Screen Stack │──│ Bubble Tea Components    │ │
│  └──────────────┘  └──────────────┘  │ • Editor  • File Browser │ │
│                                       │ • Theory  • AI Panel     │ │
│                                       │ • Settings               │ │
│                                       └──────────────────────────┘ │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER                              │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────────┐  │
│  │ Editor Service  │  │ Theory Service  │  │  AI Orchestrator │  │
│  │ • Parse YAML    │  │ • Chords/Scales │  │  • Multi-Agent   │  │
│  │ • Validate      │  │ • Circle of 5ths│  │  • RAG Query     │  │
│  │ • Auto-save     │  │ • Rhyme Dict    │  │  • Streaming     │  │
│  └─────────────────┘  └─────────────────┘  └──────────────────┘  │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                       DOMAIN LAYER                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌────────────────────────┐  │
│  │ Song Model   │  │ Project      │  │ Knowledge Base         │  │
│  │ • Metadata   │  │ • Versions   │  │ • Pedagogy Cards       │  │
│  │ • Sections   │  │ • Stats      │  │ • Anti-Cliché Lists    │  │
│  │ • Quality    │  │ • Git Repo   │  │ • Example Patterns     │  │
│  └──────────────┘  └──────────────┘  └────────────────────────┘  │
└────────────────────────┬───────────────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────────────┐
│                    INFRASTRUCTURE LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐ │
│  │ File System  │  │ SQLite DB    │  │ External Integrations    │ │
│  │ • Markdown   │  │ • Versions   │  │ • Ollama API (Local)     │ │
│  │ • Git Ops    │  │ • Stats      │  │ • MIDI Devices           │ │
│  │ • Exports    │  │ • KB Vectors │  │ • Audio (beep/GoAudio)   │ │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘ │
└────────────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

**Presentation Layer:**
- Rendering terminal UI using Bubble Tea/Lipgloss
- Handling keyboard/mouse input
- Screen routing and navigation
- Animation and visual effects

**Application Layer:**
- Business logic orchestration
- Service coordination
- AI agent pipeline management
- File operations and data transformation

**Domain Layer:**
- Core business entities (Song, Section, Line)
- Domain rules and validation
- Quality scoring algorithms
- Knowledge base queries

**Infrastructure Layer:**
- File system I/O
- Database operations
- External API calls (Ollama, MIDI)
- Audio playback and generation

---

## Technology Stack

### Core Framework & UI

| Component | Technology | Version | Purpose & Rationale |
|-----------|------------|---------|---------------------|
| **Language** | Go | 1.21+ | Fast compilation, excellent concurrency, cross-platform, single binary |
| **TUI Framework** | Bubble Tea | v0.25+ | Elm Architecture, proven in production (GitHub CLI), excellent state management |
| **Styling** | Lipgloss | v0.9+ | CSS-like terminal styling, adaptive colors, gradients, borders |
| **Components** | Bubbles | v0.18+ | Pre-built textarea, viewport, list, table, spinner components |
| **Forms** | Huh | v0.3+ | Accessible forms with screen reader support, validation |
| **Animations** | Harmonica | v0.2+ | Spring-physics animations for smooth transitions |
| **Markdown** | Glamour | v0.6+ | Markdown rendering with themes, syntax highlighting |

### Music & Audio

| Component | Technology | Purpose & Rationale |
|-----------|------------|---------------------|
| **Music Theory** | go-music-theory | Scales, chords, keys, Circle of Fifths - zero dependencies |
| **Rhyme/Syllable** | pronouncing | CMU pronunciation dictionary, rhyme detection, syllable counting |
| **Audio Feedback** | gen2brain/beeep | Cross-platform beeps with frequency/duration control |
| **MIDI** | gomidi/midi v2 | Real-time MIDI I/O, midicatdrv backend (no CGO) |
| **Synthesis** | DylanMeeus/GoAudio | Oscillators, ADSR, tone generation for chord playback |
| **Audio Playback** | gopxl/beep | WAV/MP3 playback (future: reference audio) |

### AI & NLP

| Component | Technology | Purpose & Rationale |
|-----------|------------|---------------------|
| **LLM Runtime** | Ollama | Local model hosting, OpenAI-compatible API, user installs separately |
| **LLM Client** | ollama/ollama/api | Official Go SDK, streaming support, context management |
| **Vector DB** | ChromaDB | Embeddings storage for RAG knowledge base |
| **Embeddings** | nomic-embed-text | 274MB model, efficient, good for semantic search |
| **Primary Models** | qwen2.5:7b-instruct / llama3.1:8b | Qwen for creative/ideation, Llama for structured/constrained |

### Data & Persistence

| Component | Technology | Purpose & Rationale |
|-----------|------------|---------------------|
| **Database** | SQLite (go-sqlite3) | Version history, stats, KB vectors - single file, no server |
| **Config Format** | YAML (gopkg.in/yaml.v2) | Song metadata frontmatter, user preferences |
| **Primary Storage** | Markdown files | Human-readable, version-control friendly, plain text durability |
| **Git Integration** | go-git/go-git v5 | Pure Go Git implementation, auto-commit, history |

### Build & Distribution

| Component | Technology | Purpose & Rationale |
|-----------|------------|---------------------|
| **CLI Framework** | Cobra v1.8+ | Command structure, flags, help text generation |
| **Config Management** | Viper v1.18+ | Config file reading, env vars, defaults |
| **Release Automation** | GoReleaser | Multi-platform builds, checksums, GitHub releases |
| **Testing** | teatest + testify | TUI golden file testing + standard Go testing |
| **Linting** | golangci-lint | Comprehensive linting, enforces best practices |

---

## Database Design & Data Models

### Database Schema (SQLite)

```sql
-- Songs table (metadata index)
CREATE TABLE songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filepath TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    artist TEXT,
    key TEXT,
    tempo INTEGER,
    time_signature TEXT,
    structure TEXT,
    tags TEXT, -- JSON array
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    word_count INTEGER DEFAULT 0,
    verse_count INTEGER DEFAULT 0,
    chorus_count INTEGER DEFAULT 0,
    quality_score REAL DEFAULT 0.0
);

-- Version history (auto-save snapshots)
CREATE TABLE versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_milestone BOOLEAN DEFAULT FALSE,
    milestone_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- Writing statistics
CREATE TABLE writing_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL UNIQUE,
    words_written INTEGER DEFAULT 0,
    songs_created INTEGER DEFAULT 0,
    songs_edited INTEGER DEFAULT 0,
    ai_requests INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0
);

-- Knowledge Base vectors (RAG)
CREATE TABLE kb_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding BLOB, -- Vector embedding
    metadata TEXT, -- JSON (source, rules, examples)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_songs_updated ON songs(updated_at);
CREATE INDEX idx_versions_song ON versions(song_id, created_at);
CREATE INDEX idx_stats_date ON writing_stats(date);
CREATE INDEX idx_kb_topic ON kb_entries(topic);
```

### Core Data Structures (Go)

```go
// Song represents a complete song with metadata and content
type Song struct {
    Metadata   SongMetadata `yaml:"---"`
    Sections   []Section    `yaml:"-"`
    RawContent string       `yaml:"-"`
}

// SongMetadata is the YAML frontmatter
type SongMetadata struct {
    Title         string    `yaml:"title"`
    Artist        string    `yaml:"artist,omitempty"`
    Key           string    `yaml:"key,omitempty"`
    Tempo         int       `yaml:"tempo,omitempty"`
    TimeSignature string    `yaml:"time_signature,omitempty"`
    Structure     string    `yaml:"structure,omitempty"`
    Tags          []string  `yaml:"tags,omitempty"`
    CreatedAt     time.Time `yaml:"created_at"`
    UpdatedAt     time.Time `yaml:"updated_at"`
}

// Section represents verse, chorus, bridge, etc.
type Section struct {
    Type   SectionType `json:"type"`
    Number int         `json:"number"`
    Lines  []Line      `json:"lines"`
    Notes  string      `json:"notes,omitempty"`
}

type SectionType string

const (
    SectionVerse     SectionType = "verse"
    SectionChorus    SectionType = "chorus"
    SectionPreChorus SectionType = "pre-chorus"
    SectionBridge    SectionType = "bridge"
    SectionOutro     SectionType = "outro"
    SectionIntro     SectionType = "intro"
)

// Line is a single lyric line with analysis
type Line struct {
    Text          string   `json:"text"`
    Syllables     int      `json:"syllables"`
    RhymeScheme   string   `json:"rhyme_scheme,omitempty"`
    StressPattern string   `json:"stress_pattern,omitempty"`
    Flags         []string `json:"flags,omitempty"`
}

// QualityScore represents the 7-dimension rubric
type QualityScore struct {
    Specificity        float64           `json:"specificity"`
    Originality        float64           `json:"originality"`
    EmotionalResonance float64           `json:"emotional_resonance"`
    Prosody            float64           `json:"prosody"`
    Coherence          float64           `json:"coherence"`
    VoiceConsistency   float64           `json:"voice_consistency"`
    SurpriseFactor     float64           `json:"surprise_factor"`
    Total              float64           `json:"total"`
    Feedback           map[string]string `json:"feedback"`
}

// KnowledgeCard for RAG system
type KnowledgeCard struct {
    ID       string   `json:"id"`
    Topic    string   `json:"topic"`
    Summary  string   `json:"summary"`
    Rules    []string `json:"rules"`
    Examples []string `json:"examples"`
    Checks   []string `json:"checks"`
}
```

### File Format Example

```markdown
---
title: "Lost in Translation"
artist: "Simon"
key: "Am"
tempo: 95
time_signature: "4/4"
structure: "V-PC-C-V-PC-C-B-C"
tags: ["introspective", "alternative", "complete"]
created_at: 2025-10-17T14:30:00Z
updated_at: 2025-10-17T16:45:00Z
---

## Verse 1

Coffee's gone cold on the window sill
Your jacket still draped on the chair
I'm learning how silence can fill a room
Like you were never really there

## Pre-Chorus

And I'm...

## Chorus

Lost in translation between the lines we never said
Searching for meaning in the space you left instead
Every word we swallowed down
Now echoes through this empty town
Lost in translation, can't find my way back
```

---

## AI Agent Architecture

### Multi-Agent Orchestration Pattern

Based on research documents, we implement a **deterministic 6-stage pipeline** with specialized agents for each creative task.

### Agent Pipeline Flow

```
┌─────────────────────────────────────────────────────────────────────┐
│ Stage 1: IDEATION AGENT                                             │
│ Input:  User theme/mood/genre                                       │
│ Output: 5-10 specific angles, avoiding generic concepts             │
│ Model:  Qwen 2.5 7B (creative, good with brainstorming)            │
│ Temp:   0.85-0.90 (high creativity)                                │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│ Stage 2: STRUCTURE PLANNER AGENT                                    │
│ Input:  Selected angle from ideation                                │
│ Output: Song structure (V-C-V-C-B-C), emotional arc, key imagery   │
│ Model:  Llama 3.1 8B (good at structured thinking)                 │
│ Temp:   0.70 (balanced)                                            │
│ RAG:    Query KB for structure templates                           │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│ Stage 3: GENERATION AGENTS (Parallel)                               │
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌─────────────────┐  │
│  │ Verse Generator  │  │ Chorus Generator │  │ Bridge Generator│  │
│  │ • Storytelling   │  │ • Emotional peak │  │ • New angle     │  │
│  │ • Specific detail│  │ • Hook placement │  │ • Contrast      │  │
│  │ • Build tension  │  │ • Memorable      │  │ • Resolution    │  │
│  └──────────────────┘  └──────────────────┘  └─────────────────┘  │
│                                                                      │
│ Model:  Qwen 2.5 7B per section                                    │
│ Temp:   0.80-0.85                                                  │
│ RAG:    Query KB for section-specific techniques                   │
│ Output: Raw lyric sections with target rhyme scheme                │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│ Stage 4: HUMANIZER AGENT                                            │
│ Input:  Generated sections                                          │
│ Tasks:  • Replace abstractions with concrete images                │
│         • Strengthen verbs (trudged vs walked)                     │
│         • Remove clichés from blocklist                            │
│         • Add sensory details (7 senses)                           │
│ Model:  Llama 3.1 8B (precise edits)                               │
│ Temp:   0.60 (controlled refinement)                               │
│ RAG:    Query anti-cliché lists, sensory vocabulary                │
└─────────────────────────┬───────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────────┐
│ Stage 5: CRITIC AGENT (Iterative Loop)                              │
│ Input:  Humanized lyrics                                            │
│ Tasks:  • Run quality rubric (7 dimensions)                        │
│         • Programmatic checks (prosody, syllables, rhyme)         │
│         • Identify specific issues with line numbers              │
│ Output: QualityScore + actionable feedback                         │
│ Model:  Llama 3.1 8B (analytical)                                  │
│ Temp:   0.30 (factual assessment)                                 │
│                                                                      │
│ ┌─────────────────────────────────────────────────────────────┐   │
│ │ IF score < 70/100:                                          │   │
│ │   ↓                                                         │   │
│ │ Stage 6: REFINER AGENT                                      │   │
│ │ Input:  Critique feedback                                   │   │
│ │ Tasks:  • Address specific dimension weaknesses            │   │
│ │         • Re-generate problem lines                        │   │
│ │         • Maintain overall voice/theme                     │   │
│ │ Output: Improved lyrics                                     │   │
│ │   ↓                                                         │   │
│ │ LOOP BACK TO CRITIC (max 3 iterations)                     │   │
│ └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
│ ELSE: Output final polished lyrics                                 │
└─────────────────────────────────────────────────────────────────────┘
```

### RAG Knowledge Base Integration

Each agent queries the knowledge base before generation to retrieve relevant pedagogical rules and examples.

**RAG Query Process:**

1. Convert user intent → embedding vector
2. Semantic search in ChromaDB (top 5 matches)
3. Extract relevant rules/examples
4. Inject into agent system prompt
5. Agent generates with grounded context

**Example:**

```
Query:    "How to write a memorable chorus?"
Retrieved:
  - Hook placement strategies (Pat Pattison)
  - Power positioning (Andrea Stolpe)
  - Repetition techniques (Jason Blume)
  - Concrete examples from hit songs
```

### Streaming Implementation Pattern

```go
// Bubble Tea command for streaming AI response
func streamAIResponse(prompt string, agentType AgentType) tea.Cmd {
    return func() tea.Msg {
        chunkChan := make(chan string, 10)
        
        go func() {
            defer close(chunkChan)
            
            req := &api.GenerateRequest{
                Model:  getModelForAgent(agentType),
                Prompt: buildPromptWithRAG(prompt, agentType),
                Stream: true,
                Options: map[string]interface{}{
                    "temperature": getTemperatureForAgent(agentType),
                    "num_ctx":     4096,
                },
            }
            
            client.Generate(context.Background(), req, func(resp api.GenerateResponse) error {
                chunkChan <- resp.Response
                return nil
            })
        }()
        
        return AIChunkMsg{
            Agent: agentType,
            Chunk: <-chunkChan,
            Chan:  chunkChan,
        }
    }
}
```

### Model Selection Strategy

| Agent Type | Model | Temperature | Rationale |
|------------|-------|-------------|-----------|
| Ideation | Qwen 2.5 7B | 0.85-0.90 | Creative brainstorming, diverse ideas |
| Structure Planner | Llama 3.1 8B | 0.70 | Structured thinking, logical flow |
| Verse Generator | Qwen 2.5 7B | 0.80-0.85 | Creative storytelling, imagery |
| Chorus Generator | Qwen 2.5 7B | 0.80-0.85 | Memorable hooks, emotional peaks |
| Bridge Generator | Qwen 2.5 7B | 0.80-0.85 | Fresh perspective, contrast |
| Humanizer | Llama 3.1 8B | 0.60 | Precise edits, controlled changes |
| Critic | Llama 3.1 8B | 0.30 | Analytical, factual assessment |
| Refiner | Llama 3.1 8B | 0.65 | Targeted improvements |

---

## Component Architecture & UI Flows

### Bubble Tea Component Tree

```
RootModel (main router)
├── SplashScreen (initial load, Ollama check)
├── MainMenu (file browser, new song, settings)
├── EditorModel (main workspace)
│   ├── SplitPaneModel
│   │   ├── TextareaModel (left: editing)
│   │   └── PreviewViewport (right: Glamour render)
│   ├── StatusBar (word count, cursor pos, mode)
│   ├── AIPanel (collapsible sidebar)
│   │   ├── BrainstormView
│   │   ├── SuggestionView (inline completions)
│   │   ├── CritiqueView (quality feedback)
│   │   └── ChatView (conversational AI)
│   └── QuickActions (chord insert, rhyme lookup)
├── TheoryToolModel
│   ├── CircleOfFifthsView (interactive, animated)
│   ├── ChordProgressionView (visual timeline)
│   ├── RhymeDictionaryView (search + results)
│   └── ScaleReferenceView (key selection)
├── AudioToolModel
│   ├── MetronomeView (visual beat, tempo control)
│   ├── ChordPlayerView (play progression)
│   └── MIDIInputView (capture chords)
├── ProjectManagerModel
│   ├── FileListView (fuzzy search)
│   ├── VersionHistoryView (timeline with diffs)
│   └── StatsView (charts, streaks, insights)
└── SettingsModel
    ├── AppearanceView (themes, colors, fonts)
    ├── AIConfigView (model selection, temperature)
    └── KeybindingsView (customizable shortcuts)
```

### Key User Flows

#### Flow 1: Create New Song with AI Assistance

1. Launch app → Main Menu
2. Press **N** for New Song
3. Huh form: Title, Key, Tempo, Genre/Mood
4. Auto-navigate to Editor with blank file
5. Press **Ctrl+B** to open AI Brainstorm
6. View 5-10 angle suggestions (streaming)
7. Select angle with number key
8. AI generates structure outline
9. User starts typing verse in left pane
10. Real-time preview updates in right pane
11. Pause typing 60s → AI offers inline suggestion
12. **Tab** to accept, **Esc** to dismiss
13. Complete section → Press **Ctrl+R** for refinement
14. View critique scores + feedback
15. Accept/reject AI changes via diff view
16. Repeat for other sections
17. Auto-save every 30 seconds

#### Flow 2: Use Music Theory Tools

1. In Editor, press **Ctrl+T** → Theory Tools
2. Tab to Circle of Fifths view
3. Arrow keys navigate around circle
4. Enter on "G" → Shows G major scale, diatonic chords
5. Press **P** → Switch to Chord Progressions
6. View common progressions (I-V-vi-IV, etc.)
7. Press number to hear chord (beep synthesis)
8. Press **I** to insert chord symbols into lyrics
9. Press **R** → Rhyme Dictionary
10. Type word → View perfect/slant rhymes
11. Select rhyme → Insert at cursor
12. Press **Esc** → Back to Editor

#### Flow 3: Version History & Git

1. Press **Ctrl+H** → Version History
2. Timeline shows auto-saves + milestones
3. Arrow keys navigate timeline
4. Press **D** to see diff vs current
5. Press **R** to restore old version
6. Press **M** to mark current as milestone
7. Enter milestone name
8. Auto-commit to Git with message
9. Press **Esc** → Back to Editor

---

## API Specifications

### Internal Service APIs

#### EditorService

```go
type EditorService interface {
    LoadSong(filepath string) (*Song, error)
    SaveSong(song *Song) error
    ParseContent(content string) ([]Section, error)
    ValidateMetadata(meta *SongMetadata) error
    AutoSave(song *Song) <-chan error
    CreateVersion(songID int, content string, isMilestone bool) error
}
```

#### TheoryService

```go
type TheoryService interface {
    GetScale(key string, scaleType string) ([]string, error)
    GetChord(root string, chordType string) ([]string, error)
    GetProgression(key string, pattern string) ([]string, error)
    FindRhymes(word string, rhymeType RhymeType) ([]string, error)
    CountSyllables(word string) (int, error)
    AnalyzeProsody(line string) (*ProsodyAnalysis, error)
}
```

#### AIOrchestrator

```go
type AIOrchestrator interface {
    RunPipeline(ctx context.Context, input *PipelineInput) (*PipelineOutput, error)
    Brainstorm(ctx context.Context, theme string) ([]string, error)
    GenerateSection(ctx context.Context, sectionType SectionType, context string) (*Section, error)
    Humanize(ctx context.Context, section *Section) (*Section, error)
    Critique(ctx context.Context, song *Song) (*QualityScore, error)
    Refine(ctx context.Context, section *Section, feedback string) (*Section, error)
    StreamCompletion(ctx context.Context, prompt string) (<-chan string, error)
}
```

---

## Security & Privacy

### Data Protection

- **Local-First:** All data stored locally, no cloud uploads
- **No Telemetry:** No analytics, tracking, or data collection
- **File Permissions:** User-writable directories only (~/.lyricforge/)
- **Git Security:** All commits signed with user's Git identity

### Input Validation

- Sanitize file paths (prevent directory traversal)
- Validate YAML frontmatter (prevent code injection)
- Limit file sizes (prevent memory exhaustion)
- Never execute AI-generated code

---

## Performance Optimization

### Goals

- **Startup:** <1 second from launch to ready
- **Rendering:** 60 FPS, no dropped frames
- **AI Latency:** <500ms to first token, <2s total
- **Memory:** <200MB without models, <2GB with models

### Strategies

- Lazy-load heavy dependencies
- Cache Glamour renders when unchanged
- Keep Ollama models loaded (OLLAMA_KEEP_ALIVE=10m)
- Use smaller quantized models (Q4/Q5)
- Stream AI responses immediately

---

## Testing Strategy

### Unit Tests (70%+ coverage)

```go
func TestSongParser(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Song
        wantErr bool
    }{
        {
            name:    "valid song with frontmatter",
            input:   sampleValidSong,
            want:    expectedParsedSong,
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseSong(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            assert.Equal(t, tt.want, got)
        })
    }
}
```

### TUI Tests (Golden Files)

```go
func TestEditorView(t *testing.T) {
    m := newEditorModel()
    tm := teatest.NewTestModel(t, m)
    
    tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("Hello")})
    
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Hello"))
    })
    
    teatest.RequireEqualOutput(t, tm.Output())
}
```

---

## Deployment & Distribution

### Build Configuration (GoReleaser)

```yaml
builds:
  - binary: lyricforge
    goos: [linux, windows, darwin]
    goarch: [amd64, arm64]
    ldflags:
      - -s -w -X main.version={{.Version}}

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

release:
  github:
    owner: puente-labs
    name: lyricforge
```

### Installation Methods

**Homebrew:**
```bash
brew tap puente-labs/tap
brew install lyricforge
```

**Direct Download:**
```bash
curl -L https://github.com/puente-labs/lyricforge/releases/latest/download/lyricforge_Linux_x86_64.tar.gz | tar xz
sudo mv lyricforge /usr/local/bin/
```

**Go Install:**
```bash
go install github.com/puente-labs/lyricforge@latest
```

### Prerequisites

Users must install:
1. **Ollama** - `brew install ollama` or https://ollama.ai
2. **Recommended models:**
   ```bash
   ollama pull qwen2.5:7b-instruct
   ollama pull llama3.1:8b
   ollama pull nomic-embed-text
   ```

---

## Development Phases

### Phase 1: Foundation (Week 1-2)
- Project structure & CI/CD
- Core data models
- Basic Bubble Tea app
- File I/O and database

### Phase 2: Editor & UI (Week 3-4)
- Split-pane editor
- Glamour rendering
- Beautiful Lipgloss styling
- Animations

### Phase 3: Music Theory (Week 5-6)
- Circle of Fifths
- Rhyme dictionary
- Chord progressions
- Syllable counter

### Phase 4: Audio (Week 7)
- Visual metronome
- Chord playback
- MIDI input

### Phase 5: AI Integration (Week 8-10)
- Ollama client
- RAG knowledge base
- Agent pipeline
- Streaming UI

### Phase 6: Project Management (Week 11)
- Version history
- Git integration
- Statistics
- Export formats

### Phase 7: Polish & Release (Week 12-13)
- Testing (>70% coverage)
- Documentation
- Demo materials
- Performance optimization

---

**Document Owner:** Simon (Puente Labs)  
**Status:** Ready for development  
**Next Review:** After Phase 1 completion