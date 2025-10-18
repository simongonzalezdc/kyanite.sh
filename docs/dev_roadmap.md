# Development Roadmap & Implementation Guide
## noise.sh: AI-Powered Songwriting TUI

**Version:** 1.0  
**Date:** October 17, 2025  
**Document Owner:** Simon (Puente Labs)  
**For:** Solo developer with AI-assisted coding (Kilocode/VSCode)

---

## Quick Start: First Steps

### 1. Environment Setup (Day 1)

```bash
# Create project
mkdir noise.sh && cd noise.sh
git init

# Initialize Go module
go mod init github.com/puente-labs/noise.sh

# Install core dependencies
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/glamour@latest
go get github.com/charmbracelet/huh@latest
go get github.com/charmbracelet/harmonica@latest

# Music/audio
go get github.com/go-music-theory/music-theory@latest
go get github.com/gen2brain/beeep@latest
go get gitlab.com/gomidi/midi/v2@latest
go get github.com/DylanMeeus/GoAudio@latest

# AI/NLP
go get github.com/ollama/ollama/api@latest
# Note: ChromaDB will be accessed via HTTP

# Data/persistence
go get gopkg.in/yaml.v2
go get github.com/mattn/go-sqlite3@latest
go get github.com/go-git/go-git/v5@latest

# CLI framework
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest

# Testing
go get github.com/stretchr/testify@latest
```

### 2. Project Structure

```
noise.sh/
├── cmd/
│   └── noise.sh/
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

### 3. Initial Code Scaffolding

**Create `cmd/noise.sh/main.go`:**

```go
package main

import (
    "fmt"
    "os"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/puente-labs/noise.sh/internal/ui"
)

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    // Initialize root model
    m := ui.NewRootModel()
    
    // Create Bubble Tea program
    p := tea.NewProgram(
        m,
        tea.WithAltScreen(),       // Use alternate screen buffer
        tea.WithMouseCellMotion(), // Enable mouse support
    )
    
    // Run the program
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

**Create `internal/ui/root.go`:**

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type screen int

const (
    screenSplash screen = iota
    screenMenu
    screenEditor
    screenTheory
    screenAudio
    screenManager
    screenSettings
)

type RootModel struct {
    currentScreen screen
    width         int
    height        int
    
    // Child models
    // ... will add as we build
}

func NewRootModel() RootModel {
    return RootModel{
        currentScreen: screenSplash,
    }
}

func (m RootModel) Init() tea.Cmd {
    return tea.Batch(
        tea.EnterAltScreen,
        checkOllamaCmd(), // Check if Ollama is running
    )
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        }
    }
    
    // Route to child models based on current screen
    // ... implement routing
    
    return m, nil
}

func (m RootModel) View() string {
    switch m.currentScreen {
    case screenSplash:
        return m.renderSplash()
    case screenMenu:
        return m.renderMenu()
    // ... other screens
    default:
        return "Loading..."
    }
}

func (m RootModel) renderSplash() string {
    style := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF69B4")).
        Align(lipgloss.Center).
        Width(m.width).
        Height(m.height)
    
    return style.Render("🎵 noise.sh 🎵\n\nLoading...")
}

func checkOllamaCmd() tea.Cmd {
    return func() tea.Msg {
        // Check if Ollama is accessible
        // Return ollamaStatusMsg
        return nil
    }
}
```

---

## Phase-by-Phase Implementation Guide

### Phase 1: Foundation (Week 1-2)

#### Day 1-2: Core Data Models

**Priority: CRITICAL**

Create `internal/domain/song.go`:

```go
package domain

import (
    "time"
)

type Song struct {
    ID         int          `json:"id"`
    Filepath   string       `json:"filepath"`
    Metadata   SongMetadata `json:"metadata"`
    Sections   []Section    `json:"sections"`
    RawContent string       `json:"-"`
}

type SongMetadata struct {
    Title         string    `yaml:"title" json:"title"`
    Artist        string    `yaml:"artist,omitempty" json:"artist,omitempty"`
    Key           string    `yaml:"key,omitempty" json:"key,omitempty"`
    Tempo         int       `yaml:"tempo,omitempty" json:"tempo,omitempty"`
    TimeSignature string    `yaml:"time_signature,omitempty" json:"time_signature,omitempty"`
    Structure     string    `yaml:"structure,omitempty" json:"structure,omitempty"`
    Tags          []string  `yaml:"tags,omitempty" json:"tags,omitempty"`
    CreatedAt     time.Time `yaml:"created_at" json:"created_at"`
    UpdatedAt     time.Time `yaml:"updated_at" json:"updated_at"`
}

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

type Line struct {
    Text          string   `json:"text"`
    Syllables     int      `json:"syllables"`
    RhymeScheme   string   `json:"rhyme_scheme,omitempty"`
    StressPattern string   `json:"stress_pattern,omitempty"`
    Flags         []string `json:"flags,omitempty"`
}
```

**AI Coding Prompt for Kilocode:**

```
Create comprehensive unit tests for the Song domain model in internal/domain/song_test.go. 
Include tests for:
- SongMetadata validation
- Section type constants
- Line syllable counting
- JSON marshaling/unmarshaling
- Edge cases (empty strings, nil slices)

Use testify/assert for assertions.
```

#### Day 3-4: File I/O & Parsing

**Priority: CRITICAL**

Create `internal/infra/files/parser.go`:

```go
package files

import (
    "bufio"
    "bytes"
    "errors"
    "io"
    "strings"
    
    "gopkg.in/yaml.v2"
    "github.com/puente-labs/noise.sh/internal/domain"
)

type Parser struct{}

func NewParser() *Parser {
    return &Parser{}
}

// ParseSong reads markdown file and extracts metadata + sections
func (p *Parser) ParseSong(content string) (*domain.Song, error) {
    song := &domain.Song{}
    
    // Extract YAML frontmatter
    metadata, body, err := p.extractFrontmatter(content)
    if err != nil {
        return nil, err
    }
    
    // Parse metadata
    if err := yaml.Unmarshal([]byte(metadata), &song.Metadata); err != nil {
        return nil, err
    }
    
    // Parse sections
    song.Sections, err = p.parseSections(body)
    if err != nil {
        return nil, err
    }
    
    song.RawContent = content
    return song, nil
}

func (p *Parser) extractFrontmatter(content string) (metadata, body string, err error) {
    if !strings.HasPrefix(content, "---\n") {
        return "", "", errors.New("missing YAML frontmatter")
    }
    
    // Find closing ---
    parts := strings.SplitN(content[4:], "\n---\n", 2)
    if len(parts) != 2 {
        return "", "", errors.New("malformed YAML frontmatter")
    }
    
    return parts[0], parts[1], nil
}

func (p *Parser) parseSections(body string) ([]domain.Section, error) {
    var sections []domain.Section
    scanner := bufio.NewScanner(strings.NewReader(body))
    
    var currentSection *domain.Section
    
    for scanner.Scan() {
        line := scanner.Text()
        
        // Check for section header (## Verse 1, ## Chorus, etc.)
        if strings.HasPrefix(line, "## ") {
            // Save previous section
            if currentSection != nil {
                sections = append(sections, *currentSection)
            }
            
            // Parse new section
            sectionName := strings.TrimPrefix(line, "## ")
            currentSection = p.parseSection(sectionName)
            continue
        }
        
        // Add line to current section
        if currentSection != nil && strings.TrimSpace(line) != "" {
            currentSection.Lines = append(currentSection.Lines, domain.Line{
                Text: line,
            })
        }
    }
    
    // Add last section
    if currentSection != nil {
        sections = append(sections, *currentSection)
    }
    
    return sections, scanner.Err()
}

func (p *Parser) parseSection(name string) *domain.Section {
    section := &domain.Section{}
    
    // Parse "Verse 1", "Chorus", "Bridge", etc.
    parts := strings.Fields(name)
    if len(parts) == 0 {
        return section
    }
    
    sectionType := strings.ToLower(parts[0])
    switch sectionType {
    case "verse":
        section.Type = domain.SectionVerse
    case "chorus":
        section.Type = domain.SectionChorus
    case "pre-chorus":
        section.Type = domain.SectionPreChorus
    case "bridge":
        section.Type = domain.SectionBridge
    case "intro":
        section.Type = domain.SectionIntro
    case "outro":
        section.Type = domain.SectionOutro
    default:
        section.Type = domain.SectionVerse
    }
    
    // Parse number if present
    if len(parts) > 1 {
        // Parse number (e.g., "Verse 1" -> 1)
        // Simple implementation for now
        section.Number = 1
    }
    
    return section
}

// SerializeSong converts Song to markdown format
func (p *Parser) SerializeSong(song *domain.Song) (string, error) {
    var buf bytes.Buffer
    
    // Write YAML frontmatter
    buf.WriteString("---\n")
    yamlData, err := yaml.Marshal(song.Metadata)
    if err != nil {
        return "", err
    }
    buf.Write(yamlData)
    buf.WriteString("---\n\n")
    
    // Write sections
    for _, section := range song.Sections {
        // Write section header
        buf.WriteString("## ")
        buf.WriteString(p.formatSectionName(section))
        buf.WriteString("\n\n")
        
        // Write lines
        for _, line := range section.Lines {
            buf.WriteString(line.Text)
            buf.WriteString("\n")
        }
        buf.WriteString("\n")
    }
    
    return buf.String(), nil
}

func (p *Parser) formatSectionName(section domain.Section) string {
    name := string(section.Type)
    name = strings.Title(name)
    if section.Number > 0 {
        name += " " + strconv.Itoa(section.Number)
    }
    return name
}
```

**AI Coding Prompt:**

```
Create comprehensive tests for the file parser in internal/infra/files/parser_test.go.
Include test cases for:
- Valid markdown with YAML frontmatter
- Missing frontmatter
- Malformed YAML
- Multiple sections (verse, chorus, bridge)
- Empty sections
- Roundtrip (parse then serialize)

Provide sample test data in testdata/songs/ directory.
```

#### Day 5-7: SQLite Database

**Priority: HIGH**

Create `internal/infra/db/schema.go`:

```go
package db

const schema = `
CREATE TABLE IF NOT EXISTS songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    filepath TEXT UNIQUE NOT NULL,
    title TEXT NOT NULL,
    artist TEXT,
    key TEXT,
    tempo INTEGER,
    time_signature TEXT,
    structure TEXT,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    word_count INTEGER DEFAULT 0,
    verse_count INTEGER DEFAULT 0,
    chorus_count INTEGER DEFAULT 0,
    quality_score REAL DEFAULT 0.0
);

CREATE TABLE IF NOT EXISTS versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    song_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    is_milestone BOOLEAN DEFAULT FALSE,
    milestone_name TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS writing_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL UNIQUE,
    words_written INTEGER DEFAULT 0,
    songs_created INTEGER DEFAULT 0,
    songs_edited INTEGER DEFAULT 0,
    ai_requests INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS kb_entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding BLOB,
    metadata TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_songs_updated ON songs(updated_at);
CREATE INDEX IF NOT EXISTS idx_versions_song ON versions(song_id, created_at);
CREATE INDEX IF NOT EXISTS idx_stats_date ON writing_stats(date);
CREATE INDEX IF NOT EXISTS idx_kb_topic ON kb_entries(topic);
`
```

Create `internal/infra/db/db.go`:

```go
package db

import (
    "database/sql"
    "path/filepath"
    
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    conn *sql.DB
}

func New(dataDir string) (*DB, error) {
    dbPath := filepath.Join(dataDir, "noise.sh.db")
    
    conn, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Enable foreign keys
    if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
        return nil, err
    }
    
    // Create schema
    if _, err := conn.Exec(schema); err != nil {
        return nil, err
    }
    
    return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
    return db.conn.Close()
}

// Repository methods will be added here
```

**AI Coding Prompt:**

```
Create a complete repository pattern implementation in internal/infra/db/repository.go with methods for:
- InsertSong(song *domain.Song) error
- GetSong(id int) (*domain.Song, error)
- UpdateSong(song *domain.Song) error
- DeleteSong(id int) error
- ListSongs(limit, offset int) ([]*domain.Song, error)
- SaveVersion(songID int, content string, isMilestone bool, name string) error
- GetVersions(songID int, limit int) ([]Version, error)

Include proper error handling and SQL injection prevention.
Create tests using in-memory SQLite (:memory:).
```

---

### Phase 2: Editor & UI (Week 3-4)

#### Day 8-10: Basic Editor with Bubble Tea

**Priority: CRITICAL**

Create `internal/ui/editor/model.go`:

```go
package editor

import (
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/viewport"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/glamour"
    
    "github.com/puente-labs/noise.sh/internal/domain"
)

type Model struct {
    song     *domain.Song
    textarea textarea.Model
    viewport viewport.Model
    renderer *glamour.TermRenderer
    
    activePane int // 0 = editor, 1 = preview
    width      int
    height     int
    
    autoSaveTimer time.Time
}

func New() Model {
    ta := textarea.New()
    ta.Placeholder = "Start writing your lyrics..."
    ta.Focus()
    
    vp := viewport.New(0, 0)
    
    renderer, _ := glamour.NewTermRenderer(
        glamour.WithAutoStyle(),
        glamour.WithWordWrap(80),
    )
    
    return Model{
        textarea:   ta,
        viewport:   vp,
        renderer:   renderer,
        activePane: 0,
    }
}

func (m Model) Init() tea.Cmd {
    return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmds []tea.Cmd
    var cmd tea.Cmd
    
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        
        // Split width for two panes
        halfWidth := m.width / 2
        m.textarea.SetWidth(halfWidth - 2)
        m.textarea.SetHeight(m.height - 4)
        
        m.viewport.Width = halfWidth - 2
        m.viewport.Height = m.height - 4
        
        // Re-render preview
        cmd = m.updatePreview()
        cmds = append(cmds, cmd)
        
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            // Switch active pane
            m.activePane = (m.activePane + 1) % 2
            if m.activePane == 0 {
                cmd = m.textarea.Focus()
            } else {
                m.textarea.Blur()
            }
            cmds = append(cmds, cmd)
            
        case "ctrl+s":
            // Manual save
            cmd = m.saveCmd()
            cmds = append(cmds, cmd)
        }
    }
    
    // Update active component
    if m.activePane == 0 {
        m.textarea, cmd = m.textarea.Update(msg)
        cmds = append(cmds, cmd)
        
        // Update preview on text change
        if msg, ok := msg.(tea.KeyMsg); ok && msg.Type == tea.KeyRunes {
            cmd = m.updatePreview()
            cmds = append(cmds, cmd)
        }
    } else {
        m.viewport, cmd = m.viewport.Update(msg)
        cmds = append(cmds, cmd)
    }
    
    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    // Create split-pane layout
    editorStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62")).
        Width(m.width/2 - 2).
        Height(m.height - 2)
    
    previewStyle := lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("170")).
        Width(m.width/2 - 2).
        Height(m.height - 2)
    
    if m.activePane == 0 {
        editorStyle = editorStyle.BorderForeground(lipgloss.Color("212"))
    } else {
        previewStyle = previewStyle.BorderForeground(lipgloss.Color("212"))
    }
    
    editor := editorStyle.Render(m.textarea.View())
    preview := previewStyle.Render(m.viewport.View())
    
    return lipgloss.JoinHorizontal(lipgloss.Top, editor, preview)
}

func (m Model) updatePreview() tea.Cmd {
    return func() tea.Msg {
        content := m.textarea.Value()
        rendered, _ := m.renderer.Render(content)
        m.viewport.SetContent(rendered)
        return nil
    }
}

func (m Model) saveCmd() tea.Cmd {
    return func() tea.Msg {
        // Save song to file
        // Return savedMsg
        return nil
    }
}
```

**AI Coding Prompt:**

```
Enhance the editor model with:
1. Auto-save every 30 seconds (use tea.Tick)
2. Status bar showing word count, line number, cursor position
3. Keyboard shortcuts (Ctrl+B for bold, Ctrl+I for italic in markdown)
4. Visual feedback for save operations (spinner, success message)
5. Handle unsaved changes warning

Create tests using teatest for keyboard interactions and state changes.
```

#### Day 11-12: Beautiful Lipgloss Styling

**Priority: HIGH**

Create `internal/ui/styles/theme.go`:

```go
package styles

import (
    "github.com/charmbracelet/lipgloss"
)

var (
    // Color palette (beautiful, colorful)
    Primary   = lipgloss.Color("#FF69B4") // Hot pink
    Secondary = lipgloss.Color("#9370DB") // Medium purple
    Accent    = lipgloss.Color("#FFD700") // Gold
    Success   = lipgloss.Color("#00FA9A") // Medium spring green
    Warning   = lipgloss.Color("#FFA500") // Orange
    Error     = lipgloss.Color("#FF6347") // Tomato
    Info      = lipgloss.Color("#87CEEB") // Sky blue
    
    // Adaptive colors
    Text = lipgloss.AdaptiveColor{
        Light: "#000000",
        Dark:  "#FFFFFF",
    }
    
    Background = lipgloss.AdaptiveColor{
        Light: "#FFFFFF",
        Dark:  "#000000",
    }
    
    Subtle = lipgloss.AdaptiveColor{
        Light: "#D9DCCF",
        Dark:  "#383838",
    }
)

// Common styles
var (
    TitleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        Align(lipgloss.Center).
        MarginBottom(1)
    
    SubtitleStyle = lipgloss.NewStyle().
        Foreground(Secondary).
        Italic(true)
    
    BorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Primary).
        Padding(1, 2)
    
    HighlightStyle = lipgloss.NewStyle().
        Background(Primary).
        Foreground(lipgloss.Color("#000000")).
        Bold(true).
        Padding(0, 1)
    
    // Gradient helper
    GradientColors = []lipgloss.Color{
        lipgloss.Color("#FF69B4"), // Pink
        lipgloss.Color("#FF1493"), // Deep pink
        lipgloss.Color("#C71585"), // Medium violet red
        lipgloss.Color("#9370DB"), // Medium purple
        lipgloss.Color("#8A2BE2"), // Blue violet
    }
)

func GradientText(text string) string {
    var result string
    colorCount := len(GradientColors)
    
    for i, char := range text {
        color := GradientColors[i%colorCount]
        style := lipgloss.NewStyle().Foreground(color)
        result += style.Render(string(char))
    }
    
    return result
}
```

**AI Coding Prompt:**

```
Create additional styling components:
1. Spinner animations with Harmonica spring physics
2. Progress bar with gradient fill
3. Toast notification system (success, error, info)
4. Modal dialog style
5. Tab bar for switching views
6. Badge/tag styles for metadata

Use the color palette defined in theme.go.
Export all styles from internal/ui/styles/ package.
```

---

### Phase 3: Music Theory Tools (Week 5-6)

#### Day 13-15: Circle of Fifths

**Priority: MEDIUM**

Create `internal/ui/theory/circle.go`:

```go
package theory

import (
    "math"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    theory "github.com/go-music-theory/music-theory"
)

type CircleModel struct {
    keys          []string
    selectedIndex int
    width         int
    height        int
}

func NewCircleModel() CircleModel {
    keys := []string{"C", "G", "D", "A", "E", "B", "F#/Gb", "Db", "Ab", "Eb", "Bb", "F"}
    return CircleModel{
        keys:          keys,
        selectedIndex: 0,
    }
}

func (m CircleModel) Update(msg tea.Msg) (CircleModel, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "left":
            m.selectedIndex = (m.selectedIndex - 1 + len(m.keys)) % len(m.keys)
        case "right":
            m.selectedIndex = (m.selectedIndex + 1) % len(m.keys)
        }
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    
    return m, nil
}

func (m CircleModel) View() string {
    // Calculate positions for circular layout
    centerX := m.width / 2
    centerY := m.height / 2
    radius := min(m.width, m.height) / 3
    
    // Create 2D grid for positioning
    grid := make([][]string, m.height)
    for i := range grid {
        grid[i] = make([]string, m.width)
        for j := range grid[i] {
            grid[i][j] = " "
        }
    }
    
    // Place keys around circle
    for i, key := range m.keys {
        angle := 2 * math.Pi * float64(i) / float64(len(m.keys))
        angle -= math.Pi / 2 // Start at top (12 o'clock)
        
        x := int(float64(centerX) + float64(radius)*math.Cos(angle))
        y := int(float64(centerY) + float64(radius)*math.Sin(angle))
        
        // Style based on selection
        style := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
        if i == m.selectedIndex {
            style = lipgloss.NewStyle().
                Background(lipgloss.Color("#FF69B4")).
                Foreground(lipgloss.Color("#000000")).
                Bold(true).
                Padding(0, 1)
        }
        
        // Place in grid
        rendered := style.Render(key)
        if y >= 0 && y < m.height && x >= 0 && x < m.width-len(rendered) {
            for j, char := range rendered {
                if x+j < m.width {
                    grid[y][x+j] = string(char)
                }
            }
        }
    }
    
    // Convert grid to string
    var result string
    for _, row := range grid {
        for _, cell := range row {
            result += cell
        }
        result += "\n"
    }
    
    return result
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

**AI Coding Prompt:**

```
Enhance the Circle of Fifths with:
1. Show related chords for selected key (I, ii, iii, IV, V, vi, vii°)
2. Display scale notes for selected key
3. Animated transitions using Harmonica when changing keys
4. Color-code major vs minor keys
5. Show parallel major/minor pairs
6. Add connection lines between keys (with Unicode box-drawing)

Create info panel showing:
- Key signature (sharps/flats)
- Diatonic chords
- Common progressions (I-V-vi-IV, etc.)
```

#### Day 16-18: Rhyme Dictionary & Syllable Counter

**Priority: HIGH**

Create `internal/app/theory/rhyme.go`:

```go
package theory

import (
    "github.com/abadojack/whatlanggo"
    "github.com/kward/go-cmudict"
)

type RhymeService struct {
    dict *cmudict.Dict
}

func NewRhymeService() (*RhymeService, error) {
    // Load CMU Pronouncing Dictionary
    dict, err := cmudict.LoadDefault()
    if err != nil {
        return nil, err
    }
    
    return &RhymeService{dict: dict}, nil
}

type RhymeType int

const (
    RhymePerfect RhymeType = iota
    RhymeSlant
    RhymeAssonance
    RhymeConsonance
)

func (s *RhymeService) FindRhymes(word string, rhymeType RhymeType) ([]string, error) {
    // Get pronunciation for input word
    pronunciations, ok := s.dict.Pronunciations(word)
    if !ok || len(pronunciations) == 0 {
        return nil, errors.New("word not found in dictionary")
    }
    
    targetPhonemes := pronunciations[0].Phonemes
    
    // Find rhyming words
    var rhymes []string
    
    s.dict.ForEach(func(w string, prons []cmudict.Pronunciation) {
        if w == word {
            return // Skip same word
        }
        
        for _, pron := range prons {
            if s.isRhyme(targetPhonemes, pron.Phonemes, rhymeType) {
                rhymes = append(rhymes, w)
                break
            }
        }
    })
    
    return rhymes, nil
}

func (s *RhymeService) isRhyme(a, b []string, rhymeType RhymeType) bool {
    switch rhymeType {
    case RhymePerfect:
        return s.isPerfectRhyme(a, b)
    case RhymeSlant:
        return s.isSlantRhyme(a, b)
    // ... implement other types
    }
    return false
}

func (s *RhymeService) isPerfectRhyme(a, b []string) bool {
    // Perfect rhyme: same sounds from stressed vowel to end
    // Implementation here
    return false
}

func (s *RhymeService) CountSyllables(word string) int {
    pronunciations, ok := s.dict.Pronunciations(word)
    if !ok || len(pronunciations) == 0 {
        return estimateSyllables(word) // Fallback
    }
    
    // Count vowel phonemes
    count := 0
    for _, phoneme := range pronunciations[0].Phonemes {
        if isVowel(phoneme) {
            count++
        }
    }
    
    return count
}

func isVowel(phoneme string) bool {
    // Check if phoneme represents a vowel sound
    // CMU dict vowels end with 0, 1, or 2 (stress markers)
    if len(phoneme) == 0 {
        return false
    }
    lastChar := phoneme[len(phoneme)-1]
    return lastChar == '0' || lastChar == '1' || lastChar == '2'
}

func estimateSyllables(word string) int {
    // Simple heuristic: count vowel groups
    word = strings.ToLower(word)
    count := 0
    lastWasVowel := false
    
    for _, char := range word {
        isCurrentVowel := strings.ContainsRune("aeiouy", char)
        if isCurrentVowel && !lastWasVowel {
            count++
        }
        lastWasVowel = isCurrentVowel
    }
    
    // Adjust for silent 'e'
    if strings.HasSuffix(word, "e") && count > 1 {
        count--
    }
    
    if count == 0 {
        count = 1
    }
    
    return count
}
```

**AI Coding Prompt:**

```
Create UI for rhyme dictionary in internal/ui/theory/rhyme.go:
1. Search input with real-time suggestions
2. Results grouped by rhyme type (perfect, slant, assonance)
3. Show syllable count for each result
4. Keyboard navigation through results
5. Quick insert into editor (press Enter)
6. Recently used words history

Add tests for rhyme detection algorithms.
Include edge cases (words not in dictionary, multi-syllable, compound words).
```

---

### Phase 4: Audio Features (Week 7)

#### Day 19-21: Metronome & Audio Feedback

**Priority: MEDIUM**

Create `internal/ui/audio/metronome.go`:

```go
package audio

import (
    "time"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/harmonica"
    "github.com/gen2brain/beeep"
)

type MetronomeModel struct {
    bpm      int
    running  bool
    beatNum  int
    maxBeats int
    
    // Animation
    spring harmonica.Spring
}

func NewMetronomeModel() MetronomeModel {
    return MetronomeModel{
        bpm:      120,
        running:  false,
        beatNum:  0,
        maxBeats: 4,
        spring:   harmonica.NewSpring(harmonica.FPS(60), 10.0, 0.5),
    }
}

type beatMsg struct{}
type animateMsg struct{}

func tickCmd(bpm int) tea.Cmd {
    interval := time.Minute / time.Duration(bpm)
    return tea.Tick(interval, func(t time.Time) tea.Msg {
        return beatMsg{}
    })
}

func beepCmd(frequency int) tea.Cmd {
    return func() tea.Msg {
        beeep.Beep(frequency, 50) // 50ms beep
        return nil
    }
}

func (m MetronomeModel) Update(msg tea.Msg) (MetronomeModel, tea.Cmd) {
    var cmds []tea.Cmd
    
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            m.running = !m.running
            if m.running {
                m.beatNum = 0
                cmds = append(cmds, tickCmd(m.bpm))
            }
        case "+":
            m.bpm = min(m.bpm+5, 300)
        case "-":
            m.bpm = max(m.bpm-5, 30)
        }
        
    case beatMsg:
        m.beatNum = (m.beatNum + 1) % m.maxBeats
        
        // Trigger spring animation
        m.spring.SetTarget(1.0)
        
        // Play sound (higher pitch on downbeat)
        frequency := 800
        if m.beatNum == 0 {
            frequency = 1200 // Downbeat
        }
        
        cmds = append(cmds,
            beepCmd(frequency),
            tickCmd(m.bpm),
            animate(),
        )
        
    case animateMsg:
        m.spring.Update()
        if !m.spring.Done() {
            cmds = append(cmds, animate())
        } else {
            m.spring.SetTarget(0.0)
        }
    }
    
    return m, tea.Batch(cmds...)
}

func (m MetronomeModel) View() string {
    // BPM display
    bpmStyle := lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#FF69B4")).
        Align(lipgloss.Center)
    
    bpmDisplay := bpmStyle.Render(fmt.Sprintf("%d BPM", m.bpm))
    
    // Beat indicators
    var beatIndicators string
    for i := 0; i < m.maxBeats; i++ {
        style := lipgloss.NewStyle().
            Width(4).
            Align(lipgloss.Center).
            Border(lipgloss.RoundedBorder())
        
        if i == m.beatNum && m.running {
            // Animated active beat
            intensity := m.spring.Value()
            color := interpolateColor("#FF69B4", "#FFD700", intensity)
            style = style.
                Background(lipgloss.Color(color)).
                BorderForeground(lipgloss.Color("#FFD700"))
            beatIndicators += style.Render("●")
        } else {
            style = style.
                BorderForeground(lipgloss.Color("240"))
            beatIndicators += style.Render("○")
        }
        
        beatIndicators += " "
    }
    
    // Controls
    controls := lipgloss.NewStyle().
        Foreground(lipgloss.Color("240")).
        Align(lipgloss.Center).
        Render("Space: Start/Stop | +/-: Adjust BPM")
    
    return lipgloss.JoinVertical(
        lipgloss.Center,
        bpmDisplay,
        beatIndicators,
        controls,
    )
}

func animate() tea.Cmd {
    return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
        return animateMsg{}
    })
}

func interpolateColor(start, end string, t float64) string {
    // Simple color interpolation
    // Convert hex to RGB, interpolate, convert back
    return start // Simplified for example
}

func min(a, b int) int {
    if a < b { return a }
    return b
}

func max(a, b int) int {
    if a > b { return a }
    return b
}
```

**AI Coding Prompt:**

```
Add chord playback feature using DylanMeeus/GoAudio:
1. Generate sine waves for each note in chord
2. Apply ADSR envelope for natural sound
3. Mix multiple oscillators for chord
4. Play progression with tempo sync
5. Visual waveform display (ASCII art)

Create MIDI input capture:
1. Detect connected MIDI devices
2. Listen for note on/off events
3. Display pressed notes in real-time
4. Capture chord progressions
5. Insert captured chords into song

Add tests for audio timing accuracy.
```

---

### Phase 5: AI Integration (Week 8-10)

This is the most complex phase. I'll create a separate detailed guide for AI integration.

#### Day 22-24: Ollama Client & RAG Setup

**Priority: CRITICAL**

Create `internal/infra/ollama/client.go`:

```go
package ollama

import (
    "context"
    "fmt"
    
    "github.com/ollama/ollama/api"
)

type Client struct {
    client *api.Client
    config Config
}

type Config struct {
    Host         string
    DefaultModel string
    Timeout      time.Duration
}

func NewClient(cfg Config) (*Client, error) {
    client, err := api.ClientFromEnvironment()
    if err != nil {
        return nil, fmt.Errorf("failed to create Ollama client: %w", err)
    }
    
    // Test connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := client.Heartbeat(ctx); err != nil {
        return nil, fmt.Errorf("Ollama not reachable: %w", err)
    }
    
    return &Client{
        client: client,
        config: cfg,
    }, nil
}

func (c *Client) Generate(ctx context.Context, req *api.GenerateRequest) (<-chan string, error) {
    chunkChan := make(chan string, 10)
    
    go func() {
        defer close(chunkChan)
        
        err := c.client.Generate(ctx, req, func(resp api.GenerateResponse) error {
            if resp.Response != "" {
                select {
                case chunkChan <- resp.Response:
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
            return nil
        })
        
        if err != nil {
            // Send error through channel or log
            fmt.Printf("Generate error: %v\n", err)
        }
    }()
    
    return chunkChan, nil
}

func (c *Client) Chat(ctx context.Context, messages []api.Message) (<-chan string, error) {
    req := &api.ChatRequest{
        Model:    c.config.DefaultModel,
        Messages: messages,
        Stream:   true,
    }
    
    chunkChan := make(chan string, 10)
    
    go func() {
        defer close(chunkChan)
        
        err := c.client.Chat(ctx, req, func(resp api.ChatResponse) error {
            if resp.Message.Content != "" {
                select {
                case chunkChan <- resp.Message.Content:
                case <-ctx.Done():
                    return ctx.Err()
                }
            }
            return nil
        })
        
        if err != nil {
            fmt.Printf("Chat error: %v\n", err)
        }
    }()
    
    return chunkChan, nil
}
```

**AI Coding Prompt:**

```
Create RAG knowledge base system in pkg/kb/:
1. Knowledge card structure (JSON files)
2. ChromaDB integration for vector storage
3. Embedding generation using nomic-embed-text
4. Semantic search with similarity threshold
5. Query expansion for better retrieval

Create initial knowledge base with 25-40 cards covering:
- Prosody principles (Pat Pattison)
- Rhyme schemes and patterns
- Section-specific techniques (verse, chorus, bridge)
- Anti-cliché lists
- Sensory vocabulary
- Object writing examples
- Power positioning strategies

Store in pkg/kb/data/ as JSON files.
Create ingestion script to load into ChromaDB.
```

This is getting very long! Would you like me to:

1. **Continue with the remaining phases** (AI agents, project management, polish) in this same artifact, or
2. **Create a separate artifact** for the AI Integration Deep Dive Guide?
3. **Create another artifact** for the VS Code + Kilocode setup and workflow tips?

Let me know and I'll continue!
