# syntax.sh - Technical Design Document

**Version:** 1.0  
**Date:** November 2025  
**For:** Independent Implementation  
**Status:** READY FOR IMPLEMENTATION

---

## Architecture Overview

```
┌────────────────────────────────────────────┐
│        Bubble Tea Root Model               │
├────────────────────────────────────────────┤
│                                            │
│  ┌──────────────┐  ┌──────────────┐      │
│  │ Editor Pane  │  │ Library Pane │      │
│  │              │  │ (Chars, Locs)│      │
│  └──────┬───────┘  └──────┬───────┘      │
│         │                 │              │
│         └─────────┬───────┘              │
│                   │                      │
│            ┌──────▼────────┐             │
│            │ Navigation    │             │
│            │ Router        │             │
│            └──────┬────────┘             │
├───────────────────┼────────────────────┤
│  Core Packages    │                    │
├───────────────────┼────────────────────┤
│                   │                    │
│  ┌───────┐  ┌────▼─────┐  ┌────────┐ │
│  │story/ │  │character/ │  │scene/  │ │
│  └───────┘  └──────────┘  └────────┘ │
│                                       │
│  ┌────────┐  ┌────────┐  ┌────────┐ │
│  │location│  │outline/│  │export/ │ │
│  └────────┘  └────────┘  └────────┘ │
│                                       │
│  ┌────────┐  ┌────────┐  ┌────────┐ │
│  │storage/│  │editor/ │  │stats/  │ │
│  └────────┘  └────────┘  └────────┘ │
│                                       │
│  ┌────────┐  ┌────────┐              │
│  │ai/     │  │theme/  │              │
│  └────────┘  └────────┘              │
│                                       │
└────────────────────────────────────────┘
```

---

## Tech Stack

### Core

```
Go 1.21+
Bubble Tea (TUI)
Lipgloss (styling)
Glamour (markdown rendering)
```

### Supporting Libraries

```
frontmatter (YAML parsing)
markdown (stats calculation)
yaml (YAML generation)
cobra (CLI)
viper (config)
```

### Storage

```
~/.config/syntax/projects/
└── {project-id}/
    ├── metadata.yaml
    ├── characters/
    ├── locations/
    ├── scenes/
    └── exports/
```

---

## Project Structure

```
syntax/
├── cmd/
│   └── syntax/
│       └── main.go
├── internal/
│   ├── app/
│   │   ├── model.go
│   │   ├── nav.go
│   │   └── keys.go
│   ├── ui/
│   │   ├── editor/
│   │   │   ├── editor.go
│   │   │   ├── buffer.go
│   │   │   └── viewport.go
│   │   ├── library/
│   │   │   ├── character_list.go
│   │   │   ├── location_list.go
│   │   │   └── outline_view.go
│   │   ├── components/
│   │   │   ├── textarea.go
│   │   │   ├── modal.go
│   │   │   └── list.go
│   │   └── layout.go
│   ├── story/
│   │   ├── types.go
│   │   └── project.go
│   ├── character/
│   │   ├── types.go
│   │   ├── db.go
│   │   └── render.go
│   ├── location/
│   │   ├── types.go
│   │   ├── db.go
│   │   └── render.go
│   ├── scene/
│   │   ├── types.go
│   │   ├── db.go
│   │   ├── compile.go
│   │   └── stats.go
│   ├── outline/
│   │   ├── types.go
│   │   └── manager.go
│   ├── editor/
│   │   ├── buffer.go
│   │   ├── cursor.go
│   │   └── state.go
│   ├── storage/
│   │   ├── project.go
│   │   ├── character.go
│   │   ├── scene.go
│   │   ├── location.go
│   │   ├── outline.go
│   │   └── config.go
│   ├── export/
│   │   ├── markdown.go
│   │   ├── pdf.go
│   │   ├── docx.go
│   │   ├── html.go
│   │   └── stats_report.go
│   ├── ai/
│   │   ├── assistant.go
│   │   ├── modes.go
│   │   └── prompts.go
│   ├── stats/
│   │   ├── calculator.go
│   │   ├── tracker.go
│   │   └── goals.go
│   ├── theme/
│   │   ├── registry.go
│   │   ├── manager.go
│   │   └── types.go
│   └── config/
│       └── types.go
├── tests/
│   ├── editor_test.go
│   ├── story_test.go
│   ├── character_test.go
│   └── export_test.go
├── go.mod
├── go.sum
├── README.md
├── ARCHITECTURE.md
├── ROADMAP.md
└── Makefile
```

---

## Core Data Types

### Project Type

```go
package story

type Project struct {
    ID              string
    Title           string
    Author          string
    Genre           string
    Status          string  // draft, revising, complete
    TargetWordCount int
    CreatedAt       time.Time
    UpdatedAt       time.Time
    
    // Relationships
    Characters map[string]*character.Character
    Locations  map[string]*location.Location
    Outline    *outline.Outline
    Scenes     map[string]*scene.Scene
}

func (p *Project) NewCharacter(name, role string) (*character.Character, error)
func (p *Project) NewScene(chapter, name string) (*scene.Scene, error)
func (p *Project) WordCount() int
func (p *Project) Export(format string) ([]byte, error)
```

---

### Character Type

```go
package character

type Character struct {
    ID            string
    Name          string
    Aliases       []string
    Role          string      // protagonist, antagonist, etc
    Age           int
    Occupation    string
    Appearance    string
    Background    string
    Arc           string      // Character development
    Relationships map[string]Relationship
    CreatedAt     time.Time
    UpdatedAt     time.Time
    Bio           string      // Markdown content
}

type Relationship struct {
    CharacterID string
    Type        string  // love interest, rival, etc
    Tension     string  // Low, Medium, High
    Notes       string
}
```

### Character Storage Format

**File Path:** `.story-name/characters/{character-id}.md`

**ID Generation:**
- Format: `char_` + 16 hex characters (e.g., `char_abc123def4567890`)
- Generated using: `crypto/rand` for secure random bytes
- Collision check: Verify ID doesn't exist before creating file

**File Structure:**
```markdown
---
id: char_abc123def4567890
name: "Jane Doe"
aliases: ["JD", "The Detective"]
role: "protagonist"
age: 34
occupation: "Homicide Detective"
appearance: "Tall, athletic build, short dark hair"
background: "Former military, joined police force after discharge"
arc: "Learning to trust others again"
created_at: 2025-11-10T14:30:00Z
updated_at: 2025-11-11T09:15:00Z
relationships:
  - character_id: char_xyz789abc1234567
    type: "rival"
    tension: "high"
    notes: "Competing for promotion to captain"
  - character_id: char_def456ghi7890123
    type: "mentor"
    tension: "low"
    notes: "Former training officer, now retired"
---

# Jane Doe - Character Bio

Jane grew up in a military family, moving from base to base...

## Personality Traits

- Determined and focused
- Struggles with vulnerability
- Fiercely loyal to those she trusts

## Character Arc

Begins the story as isolated and self-reliant. Through her partnership
with [other character], she learns to open up and trust her team.
```

**Parsing:**
- Use `github.com/adrg/frontmatter` for YAML parsing
- Validate required fields: `id`, `name`, `created_at`
- Optional fields get default values (empty string, zero, etc.)

---

### Scene Type

```go
package scene

type Scene struct {
    ID           string
    Chapter      int
    SceneNumber  int
    Name         string
    POVCharacter string
    Location     string      // location ID
    TimeOfDay    string      // morning, evening, etc
    PlotPoints   []string
    Status       string      // draft, revising, done
    WordCount    int
    Content      string      // Markdown
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

func (s *Scene) UpdateContent(text string) error
func (s *Scene) GetWordCount() int
```

---

### Outline Type

```go
package outline

type Outline struct {
    Structure string  // three-act, hero-journey, custom
    Acts      []Act
}

type Act struct {
    Number int
    Name   string
    Goal   string
    Beats  []Beat
}

type Beat struct {
    ID        string
    Number    int
    Name      string
    Status    string  // todo, active, done
    SceneRef  string
}
```

---

### Stats Type

```go
package stats

type ProjectStats struct {
    TotalWords        int
    TotalSessions     int
    TotalTime         time.Duration
    DaysWithWrites    int
    CurrentStreak     int
    WordsByChapter    map[int]int
    WordsByCharacter  map[string]int
    DailyStats        map[time.Time]int
}

func (s *ProjectStats) CalculateProgress(goal int) float64
func (s *ProjectStats) GetStreak() int
```

---

## Concurrency Strategy

### Auto-Save Implementation

**Problem:** Auto-save must not corrupt document during user editing

**Solution:**
```go
type AutoSaver struct {
    mu           sync.Mutex
    lastSaveTime time.Time
    debounce     time.Duration  // 300ms
    saveTimer    *time.Timer
}

func (a *AutoSaver) TriggerSave(buffer *Buffer) {
    a.mu.Lock()
    defer a.mu.Unlock()

    // Reset debounce timer
    if a.saveTimer != nil {
        a.saveTimer.Stop()
    }

    a.saveTimer = time.AfterFunc(a.debounce, func() {
        // Create snapshot before save (avoid mid-edit corruption)
        snapshot := buffer.Snapshot()

        // Save in goroutine with error channel
        go func() {
            if err := saveToFile(snapshot); err != nil {
                log.Printf("Auto-save failed: %v", err)
                // Send error to UI via channel
            }
        }()
    })
}
```

**Key Points:**
- Debounce 300ms after last keystroke
- Snapshot buffer before save (immutable copy)
- Non-blocking save in goroutine
- Error reporting via channel to UI

### Stats Calculation

**Strategy:**
- Read-only operations, safe for concurrent access
- Use atomic counters for session stats

```go
type SessionStats struct {
    WordsWritten atomic.Int64
    KeyStrokes   atomic.Int64
    StartTime    time.Time
}

func (s *SessionStats) IncrementWords(n int) {
    s.WordsWritten.Add(int64(n))
}
```

### File Locks

**Current (v1.0):**
- Not required (single-user, single-instance tool)
- User warned if opening same project twice

**Future (v1.1+):**
- Add `.lock` file for safety
- Check for stale locks (>1 hour old)
- Offer to force-open or recover

### Thread-Safety Rules

1. **UI Updates:** Only in Bubble Tea Update() method
2. **File I/O:** Always use goroutines for writes
3. **Shared State:** Protect with mutex or use channels
4. **Buffer Modifications:** Single-threaded (UI goroutine only)

---

## Module Breakdown

### 1. editor/ - Text Editing

**Key Functions:**

```go
// Buffer
NewBuffer(initialContent string) *Buffer
func (b *Buffer) Insert(pos int, text string)
func (b *Buffer) Delete(start, end int)
func (b *Buffer) GetContent() string
func (b *Buffer) GetLine(lineNum int) string

// Cursor
func (b *Buffer) CursorUp()
func (b *Buffer) CursorDown()
func (b *Buffer) CursorLeft()
func (b *Buffer) CursorRight()
func (b *Buffer) GetCursorPos() (line, col int)

// Selection
func (b *Buffer) SelectAll()
func (b *Buffer) GetSelection() string

// Undo/Redo
func (b *Buffer) Undo()
func (b *Buffer) Redo()
```

**Undo/Redo Implementation:**

```go
// EditHistory manages undo/redo state
type EditHistory struct {
    past    []BufferState  // Max 100 states
    future  []BufferState  // Cleared on new edit
    maxSize int            // 100 (configurable)
}

// BufferState captures snapshot for undo
type BufferState struct {
    Content    string
    CursorLine int
    CursorCol  int
    Timestamp  time.Time
}

func (b *Buffer) RecordState() {
    state := BufferState{
        Content:    b.GetContent(),
        CursorLine: b.cursorLine,
        CursorCol:  b.cursorCol,
        Timestamp:  time.Now(),
    }

    // Add to history
    b.history.past = append(b.history.past, state)

    // Trim if exceeds max size (keep most recent 100)
    if len(b.history.past) > b.history.maxSize {
        b.history.past = b.history.past[1:]
    }

    // Clear future (new edit invalidates redo stack)
    b.history.future = nil
}

func (b *Buffer) Undo() error {
    if len(b.history.past) == 0 {
        return errors.New("nothing to undo")
    }

    // Save current state to future
    current := b.currentState()
    b.history.future = append(b.history.future, current)

    // Pop from past
    prev := b.history.past[len(b.history.past)-1]
    b.history.past = b.history.past[:len(b.history.past)-1]

    // Restore state
    b.restoreState(prev)
    return nil
}

func (b *Buffer) Redo() error {
    if len(b.history.future) == 0 {
        return errors.New("nothing to redo")
    }

    // Save current to past
    current := b.currentState()
    b.history.past = append(b.history.past, current)

    // Pop from future
    next := b.history.future[len(b.history.future)-1]
    b.history.future = b.history.future[:len(b.history.future)-1]

    // Restore state
    b.restoreState(next)
    return nil
}

// Memory limit: Max 100 undo states
// Avg doc size: 50KB
// Total memory: ~5MB for undo history (acceptable)
```

**Tests Required:**
- Insert/delete operations
- Cursor movement
- Selection operations
- Undo/redo
- Large documents (10,000+ lines)

---

### 2. character/ - Character Management

**Key Functions:**

```go
func CreateCharacter(name, role string) *Character
func LoadCharacter(projectDir, characterID string) (*Character, error)
func (c *Character) Save(projectDir string) error
func (c *Character) Delete(projectDir string) error
func GetAllCharacters(projectDir string) ([]*Character, error)
func SearchCharacters(projectDir, query string) ([]*Character, error)
func RenderRelationshipMap(characters []*Character) string
```

**ASCII Relationship Map:**

**Library:** Hand-drawn using Unicode box-drawing characters (U+2500 range)

**Example Output:**
```
Character Relationships
━━━━━━━━━━━━━━━━━━━━━━━

     Jane Doe ←━━[rivals/high]━━→ John Smith
         ↓                              ↓
    [mentor/low]                  [partner/medium]
         ↓                              ↓
     Alice Chen ←━━[love/low]━━━→ Bob Wilson
```

**Tension Level Visualization:**
- Low: Solid line `━━━`
- Medium: Dashed line `╌╌╌`
- High: Jagged line `╱╲╱`

**Implementation:**
```go
func RenderRelationshipMap(characters []*Character) string {
    // Build adjacency graph
    graph := buildGraph(characters)

    // Layout characters in grid (simple tree layout)
    positions := calculatePositions(graph)

    // Render with box-drawing characters
    canvas := make([][]rune, 50, 80)  // 50 lines, 80 cols

    // Draw character nodes
    for _, char := range characters {
        pos := positions[char.ID]
        drawBox(canvas, pos, char.Name)
    }

    // Draw relationship lines
    for _, char := range characters {
        for _, rel := range char.Relationships {
            drawLine(canvas, positions[char.ID],
                    positions[rel.CharacterID], rel.Tension)
        }
    }

    return canvasToString(canvas)
}

// Line styles based on tension
func getLineStyle(tension string) string {
    switch tension {
    case "low":    return "━"  // Solid
    case "medium": return "╌"  // Dashed
    case "high":   return "╱"  // Jagged
    default:       return "─"  // Default
    }
}
```

**Location Map (Similar):**
```
World Map
━━━━━━━━━

[Riverside Tavern]━━━━road━━━━[Castle Blackstone]
       │                            │
    [forest]                    [courtyard]
       │                            │
   [Deep Woods]━━━━river━━━━━[Eastern Bridge]
```

---

### 3. scene/ - Scene Organization

**Key Functions:**

```go
func NewScene(chapter, name string) *Scene
func LoadScene(projectDir, sceneID string) (*Scene, error)
func (s *Scene) Save(projectDir string) error
func CompileStory(projectDir, format string) ([]byte, error)
func GetAllScenes(projectDir string) ([]*Scene, error)
func GetScenesByCharacter(projectDir, charID string) ([]*Scene, error)
```

---

### 4. storage/ - File Persistence

**Key Functions:**

```go
func LoadProject(projectID string) (*Project, error)
func SaveProject(p *Project) error
func CreateProject(name, genre string) (*Project, error)
func LoadAllCharacters(projectDir string) (map[string]*Character, error)
func LoadAllScenes(projectDir string) (map[string]*Scene, error)
```

---

### 5. export/ - Output Formats

**Key Functions:**

```go
func ExportMarkdown(projectDir string) ([]byte, error)
func ExportPDF(projectDir string) ([]byte, error)
func ExportDOCX(projectDir string) ([]byte, error)
func ExportHTML(projectDir string) ([]byte, error)
func ExportStatsReport(projectDir string) (string, error)
```

**Error Handling & Size Limits:**

```go
// Custom errors
var (
    ErrProjectTooLarge = errors.New("project exceeds size limit for this format")
    ErrNoScenes        = errors.New("project has no scenes to export")
    ErrInvalidFormat   = errors.New("unsupported export format")
    ErrExportFailed    = errors.New("export generation failed")
)

// Size limits (in bytes)
const (
    MaxDOCXSize = 10 * 1024 * 1024  // 10MB (~500,000 words)
    MaxPDFSize  = 50 * 1024 * 1024  // 50MB (~2,500,000 words)
)

func ExportDOCX(projectDir string) ([]byte, error) {
    // Calculate project size
    size, err := calculateProjectSize(projectDir)
    if err != nil {
        return nil, fmt.Errorf("failed to calculate size: %w", err)
    }

    // Check size limit
    if size > MaxDOCXSize {
        return nil, fmt.Errorf("%w: %d bytes (max %d)",
            ErrProjectTooLarge, size, MaxDOCXSize)
    }

    // Load scenes
    scenes, err := scene.LoadAll(projectDir)
    if err != nil {
        return nil, fmt.Errorf("failed to load scenes: %w", err)
    }

    if len(scenes) == 0 {
        return nil, ErrNoScenes
    }

    // Generate DOCX
    doc := docx.New()
    for _, scene := range scenes {
        // Add chapter heading
        doc.AddHeading(scene.Chapter, 1)

        // Add scene content
        doc.AddParagraph(scene.Content)
    }

    // Write to buffer
    var buf bytes.Buffer
    if err := doc.Write(&buf); err != nil {
        return nil, fmt.Errorf("%w: %v", ErrExportFailed, err)
    }

    return buf.Bytes(), nil
}

func calculateProjectSize(projectDir string) (int64, error) {
    var totalSize int64

    scenesDir := filepath.Join(projectDir, "scenes")
    err := filepath.Walk(scenesDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() {
            totalSize += info.Size()
        }
        return nil
    })

    return totalSize, err
}
```

**Export Dependencies:**
- **PDF:** `github.com/jung-kurt/gofpdf` (2MB)
- **DOCX:** `github.com/nguyenthenguyen/docx` (3MB)
- Total binary size increase: ~5MB

---

### 6. ai/ - Story Assistant

**Key Functions:**

```go
func (a *Assistant) CheckCharacter(scene, character string) string
func (a *Assistant) PolishDialogue(dialogueText string) string
func (a *Assistant) EvaluateSceneEnergy(sceneText string) string
func (a *Assistant) CheckPlotConsistency(scene *Scene) string
```

---

### 7. stats/ - Progress Tracking

**Key Functions:**

```go
func GetProjectWordCount(projectDir string) int
func GetChapterWordCount(projectDir, chapterNum int) int
func StartSession()
func EndSession() SessionStats
func GetGoalProgress() (current, target int, percent float64)
```

---

## Implementation Phases

### Phase 1: Foundation (Days 1-2)

- [ ] Project structure
- [ ] Theme system
- [ ] Editor buffer and cursor
- [ ] File I/O layer

**Deliverable:** Project compiles and runs

---

### Phase 2: Features (Days 3-5)

- [ ] Split-pane editor UI
- [ ] Character database
- [ ] Scene organization
- [ ] Location database
- [ ] Outline editor
- [ ] Storage layer

**Deliverable:** All features functional

---

### Phase 3: Polish (Days 6-8)

- [ ] AI assistant
- [ ] Stats calculation
- [ ] Export formats
- [ ] Help system
- [ ] Testing
- [ ] Documentation

**Deliverable:** v1.0 release ready

---

## Complete Code Skeleton

### main.go

```go
package main

import (
    "fmt"
    "os"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kyanite/syntax/internal/app"
)

func main() {
    m := app.NewRootModel()
    p := tea.NewProgram(
        m,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )
    
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

### internal/app/model.go

```go
package app

import tea "github.com/charmbracelet/bubbletea"

type Screen int

const (
    ScreenWelcome Screen = iota
    ScreenEditor
    ScreenLibrary
)

type RootModel struct {
    CurrentScreen Screen
    Width         int
    Height        int
    CurrentProject *story.Project
}

func NewRootModel() RootModel {
    return RootModel{
        CurrentScreen: ScreenWelcome,
    }
}

func (m RootModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return m, nil
}

func (m RootModel) View() string {
    switch m.CurrentScreen {
    case ScreenEditor:
        return "Editor"
    default:
        return "Welcome"
    }
}
```

---

## Testing Strategy

### Test Coverage Targets

**Overall Target:** 75%+ coverage

**By Module:**
- **Critical Paths (90%+ required):**
  - `editor/` - Buffer operations, cursor movement, undo/redo
  - `storage/` - File save/load, corruption handling
  - `export/` - All format generation, error cases

- **Business Logic (80%+ required):**
  - `character/` - CRUD operations, search, relationships
  - `scene/` - Scene management, compilation
  - `stats/` - Calculation accuracy

- **UI Layer (50%+ required):**
  - `ui/` - View rendering, input handling
  - `app/` - Navigation, state management

### Test Files

```
tests/
├── editor_test.go           # Buffer, cursor, undo/redo
├── editor_bench_test.go     # Performance benchmarks
├── story_test.go            # Project management
├── character_test.go        # Character CRUD, search
├── scene_test.go            # Scene compilation
├── export_test.go           # All export formats
├── storage_test.go          # File I/O, corruption
└── integration_test.go      # End-to-end workflows
```

### Example Tests

**Unit Test:**
```go
func TestBufferInsert(t *testing.T) {
    buf := editor.NewBuffer("Hello")
    buf.Insert(5, " World")

    if buf.GetContent() != "Hello World" {
        t.Errorf("got %q, want %q", buf.GetContent(), "Hello World")
    }
}
```

**Error Case Test:**
```go
func TestExportDOCX_ProjectTooLarge(t *testing.T) {
    // Create project larger than 10MB
    projectDir := createLargeProject(t, 15*1024*1024)
    defer os.RemoveAll(projectDir)

    _, err := export.ExportDOCX(projectDir)

    if !errors.Is(err, export.ErrProjectTooLarge) {
        t.Errorf("expected ErrProjectTooLarge, got %v", err)
    }
}
```

**Benchmark:**
```go
func BenchmarkBufferInsert(b *testing.B) {
    buf := editor.NewBuffer("initial text")
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        buf.Insert(5, "x")
    }
}

// Target: <1000 ns/op
```

**Integration Test:**
```go
func TestWorkflow_CreateCharacterAndSave(t *testing.T) {
    // Create project
    proj, _ := story.CreateProject("Test Novel", "fantasy")
    defer os.RemoveAll(proj.Dir)

    // Create character
    char := character.CreateCharacter("Jane Doe", "protagonist")

    // Save
    err := char.Save(proj.Dir)
    if err != nil {
        t.Fatalf("failed to save: %v", err)
    }

    // Load and verify
    loaded, err := character.LoadCharacter(proj.Dir, char.ID)
    if err != nil {
        t.Fatalf("failed to load: %v", err)
    }

    if loaded.Name != "Jane Doe" {
        t.Errorf("got name %q, want %q", loaded.Name, "Jane Doe")
    }
}
```

### Running Tests

```bash
# All tests with coverage
go test -v -cover ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Benchmarks
go test -bench=. -benchmem ./...

# Specific module
go test -v ./internal/editor/

# Race condition detection
go test -race ./...
```

---

## Performance Targets

| Operation | Target | Must-Have | Measurement Method |
|-----------|--------|-----------|-------------------|
| Startup | <1s | ✅ | `time ./syntax open "project"` |
| Auto-save | <100ms | ✅ | pprof in tests |
| Scene switch | <200ms | ✅ | UI timing logs |
| Search | <100ms | ✅ | Benchmarks |
| Memory idle | <50MB | ✅ | `ps aux` / `top` |

### Performance Measurement

**Startup Time:**
```bash
# Measure cold start
time ./syntax open "test-project"

# Should complete in <1s (including rendering)
```

**Auto-Save:**
```go
func TestAutoSavePerformance(t *testing.T) {
    proj := createTestProject(t)
    buf := editor.NewBuffer(strings.Repeat("test ", 10000)) // Large doc

    start := time.Now()
    err := saveProject(proj, buf)
    duration := time.Since(start)

    if duration > 100*time.Millisecond {
        t.Errorf("Save took %v, want <100ms", duration)
    }
}
```

**Memory Profiling:**
```bash
# Run with memory profiling
go test -memprofile=mem.prof -run=TestLongSession

# Analyze
go tool pprof mem.prof
> top10
> list functionName

# Check for leaks
go tool pprof -alloc_space mem.prof
```

**CPU Profiling:**
```bash
# Profile specific operation
go test -cpuprofile=cpu.prof -bench=BenchmarkSceneSwitch

# Visualize
go tool pprof -http=:8080 cpu.prof
```

**Real-Time Monitoring:**
```bash
# Monitor memory during long session
watch -n 1 'ps aux | grep syntax'

# Expected: <50MB idle, <200MB during heavy use
```

---

## Theme System Integration

### Theme Loading Sequence

1. **Startup:**
   - Check `~/.config/syntax/config.toml` for `theme = "theme-name"`
   - Default to "monochrome" if not set or invalid
   - Load theme from `internal/theme/registry.go`
   - Apply colors via Lipgloss styles

2. **Runtime Switching:**
   - Ctrl+Shift+T cycles through all 10 themes
   - Updates config file immediately
   - Refreshes all UI components
   - Shows theme name briefly in status bar

3. **Theme Application:**
   ```go
   func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
       switch msg := msg.(type) {
       case tea.KeyMsg:
           if msg.String() == "ctrl+shift+t" {
               m.currentTheme = m.themeManager.NextTheme()
               m.refreshStyles()
               m.saveConfig()
               return m, showMessage("Theme: " + m.currentTheme.Name)
           }
       }
       return m, nil
   }

   func (m *Model) refreshStyles() {
       theme := m.currentTheme
       m.editorStyle = lipgloss.NewStyle().
           Foreground(lipgloss.Color(theme.Text)).
           Background(lipgloss.Color(theme.Background))

       m.headingStyle = lipgloss.NewStyle().
           Foreground(lipgloss.Color(theme.Accent)).
           Bold(true)

       // ... refresh all styles
   }
   ```

### User Customization

**v1.0:** No custom themes (10 built-in only)

**v1.1 (Future):**
- Support `~/.config/syntax/themes/custom.toml`
- Format:
  ```toml
  name = "My Theme"
  primary = "#FF0000"
  secondary = "#00FF00"
  accent = "#0000FF"
  background = "#000000"
  text = "#FFFFFF"
  success = "#00FF00"
  ```
- Validate colors (valid hex)
- Show in theme switcher

### Error Handling

```go
func loadTheme(name string) Theme {
    theme, ok := registry.GetTheme(name)
    if !ok {
        log.Printf("Invalid theme %q, using default", name)
        return registry.GetTheme("monochrome")
    }
    return theme
}

// Corrupted config file
func loadConfig() Config {
    cfg, err := parseConfig(configPath)
    if err != nil {
        log.Printf("Config parse error: %v, using defaults", err)
        return defaultConfig()
    }
    return cfg
}
```

---

## Validation Checklist

### Before Release

- [ ] All 8 features implemented
- [ ] All acceptance criteria met
- [ ] 0 critical bugs
- [ ] 10 themes working
- [ ] Universal shortcuts implemented
- [ ] Performance targets met
- [ ] All tests passing
- [ ] Documentation complete
- [ ] Works on 80x24 terminal
- [ ] No panics

---

**This is a completely independent, standalone tool. Everything needed is documented here.**

**Next step:** Review README for usage
