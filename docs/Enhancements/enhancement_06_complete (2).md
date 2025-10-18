# Enhancement #6: Rapid Capture Suite + AI Right-Sizing
**noise.sh → Strategic Realignment for Rapid Ideation**

**Version:** 6.0  
**Date:** October 2025  
**Status:** READY FOR IMPLEMENTATION  
**Priority:** HIGH - Realigns entire tool with new vision  
**Estimated Effort:** 3-4 weeks  
**Owner:** Simon | Kyanite Suite

---

## Executive Summary

### Purpose
Realign noise.sh from "comprehensive songwriting suite" to "rapid capture tool" by:
- Adding live coding mode for quick musical patterns
- Simplifying AI (4 agents → 1 fast agent)
- Creating quick chord picker (instant suggestions)
- Establishing export pipeline to wave.sh
- Unifying documentation and deprecating old files

### Core Principle
**noise.sh = sketchbook, not studio. Capture ideas in <60 seconds.**

### Success Metrics
- ✅ Pattern sketching: <60 seconds
- ✅ AI response time: <2 seconds
- ✅ Chord picker: instant suggestions
- ✅ Export to wave.sh: working
- ✅ Documentation: unified and clear

---

## Part 1: New Features

### 1.1 Live Music Coding Mode

**Purpose:** Code musical patterns like a programmer (the "music coding" feature)

**This is the live coding/algorave-inspired feature - write music as code**

**Command:**
```bash
noise live [optional-name]
```

**UI Layout:**
```
┌─ noise.sh | LIVE MODE ────────────────────── 120 BPM ─┐
│                                                         │
│  [PATTERN CODE]              [VISUALIZATION]           │
│  s("bd sd bd sd")            █ . █ .                   │
│  s("hh*8").gain(0.5)         ████████                  │
│  n("0 3 7 5")                ▁▃▅▇                      │
│                                                         │
│  [Type patterns here]        [Live preview]            │
│                                                         │
│  ▶ Playing | Ctrl+P Play/Pause | Ctrl+S Save | Export │
└─────────────────────────────────────────────────────────┘
```

**Pattern Syntax (Minimal):**
```javascript
s("bd sd")           // Sample pattern (kick, snare)
n("0 3 7")           // MIDI note pattern
gain("1 0.8 0.6")    // Volume envelope
speed("1 1.5")       // Playback speed
-- comment line      // Comments
```

**Features:**
- TidalCycles-inspired syntax (simplified)
- Reuse existing split-pane editor
- Visual pattern display (right pane)
- **Basic audio playback** (play patterns to hear them)
- Save as `.noise` files
- Export PATTERN DATA as JSON (not audio - wave.sh generates audio)
- Simple sample-based playback engine

**Implementation Files:**
```
cmd/live.go                          (entry point)
internal/app/live/live_mode.go       (main logic)
internal/pattern/parser.go           (syntax parser)
internal/pattern/types.go            (data structures)
internal/ui/editor/pattern_view.go   (visualization)
internal/audio/simple_player.go      (basic audio playback)
internal/audio/sample_loader.go      (load drum samples)
data/samples/                        (basic drum kit samples)
```

---

### 1.2 Basic Audio Playback

**Purpose:** Hear patterns while coding (simple sample playback)

**Features:**
- Play/pause patterns with Ctrl+P
- Load basic drum kit samples (kick, snare, hat, etc.)
- Simple timing engine (no complex synthesis)
- Mono output, basic mixing
- **Keep lightweight** - just enough to hear ideas

**What noise.sh plays:**
- Sample-based drums (bd, sd, hh, cp, etc.)
- Basic sine wave for MIDI notes
- Volume/gain adjustments

**What goes to wave.sh:**
- Full synthesis (oscillators, ADSR, filters)
- Advanced mixing & effects
- MIDI file export
- Professional audio export

**Audio Library Choice:**
Use `github.com/hajimehoshi/oto/v2` - pure Go, cross-platform, low latency, no CGO required.

**Full Audio Implementation:**

```go
// internal/audio/simple_player.go
package audio

import (
    "encoding/binary"
    "math"
    "os"
    "time"
    
    "github.com/hajimehoshi/oto/v2"
)

const (
    sampleRate   = 44100
    channelCount = 2
    bitDepth     = 2
)

type SimplePlayer struct {
    context  *oto.Context
    samples  map[string][]float32 // Pre-loaded samples
    bpm      int
    playing  bool
    stopChan chan bool
}

func NewSimplePlayer(bpm int) (*SimplePlayer, error) {
    ctx, ready, err := oto.NewContext(sampleRate, channelCount, bitDepth)
    if err != nil {
        return nil, err
    }
    <-ready
    
    return &SimplePlayer{
        context:  ctx,
        samples:  make(map[string][]float32),
        bpm:      bpm,
        stopChan: make(chan bool),
    }, nil
}

func (p *SimplePlayer) LoadSample(name, filepath string) error {
    // Read WAV file
    data, err := os.ReadFile(filepath)
    if err != nil {
        return err
    }
    
    // Skip WAV header (44 bytes) and convert to float32
    audioData := data[44:]
    samples := make([]float32, len(audioData)/2)
    
    for i := 0; i < len(samples); i++ {
        sample := int16(binary.LittleEndian.Uint16(audioData[i*2:]))
        samples[i] = float32(sample) / 32768.0
    }
    
    p.samples[name] = samples
    return nil
}

func (p *SimplePlayer) PlayPattern(pattern *Pattern) error {
    if p.playing {
        return nil
    }
    
    p.playing = true
    go p.playLoop(pattern)
    return nil
}

func (p *SimplePlayer) playLoop(pattern *Pattern) {
    beatDuration := time.Duration(60.0/float64(p.bpm)*1000) * time.Millisecond
    ticker := time.NewTicker(beatDuration / 4) // 16th note resolution
    defer ticker.Stop()
    
    step := 0
    
    for {
        select {
        case <-p.stopChan:
            return
        case <-ticker.C:
            for _, line := range pattern.Lines {
                if p.shouldTrigger(line, step) {
                    p.playSample(line.Function, line.Effects)
                }
            }
            step++
        }
    }
}

func (p *SimplePlayer) playSample(sampleName string, effects []Effect) {
    samples, ok := p.samples[sampleName]
    if !ok {
        return
    }
    
    // Apply gain
    gain := 1.0
    for _, effect := range effects {
        if effect.Name == "gain" {
            gain = effect.Value
        }
    }
    
    // Convert to bytes
    buf := make([]byte, len(samples)*4)
    for i, sample := range samples {
        sample *= float32(gain)
        left := int16(sample * 32767)
        right := left
        
        binary.LittleEndian.PutUint16(buf[i*4:], uint16(left))
        binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(right))
    }
    
    // Play
    player := p.context.NewPlayer(bytes.NewReader(buf))
    player.Play()
}

func (p *SimplePlayer) Stop() {
    if p.playing {
        p.stopChan <- true
        p.playing = false
    }
}

func (p *SimplePlayer) Close() {
    p.Stop()
    p.context.Suspend()
}

// Simple sine wave generator for MIDI notes
func GenerateSineWave(frequency float64, duration time.Duration) []float32 {
    numSamples := int(float64(sampleRate) * duration.Seconds())
    samples := make([]float32, numSamples)
    
    for i := 0; i < numSamples; i++ {
        t := float64(i) / float64(sampleRate)
        samples[i] = float32(math.Sin(2 * math.Pi * frequency * t))
    }
    
    return samples
}
```

```go
// internal/audio/sample_loader.go
package audio

import "path/filepath"

func LoadDefaultSamples(player *SimplePlayer, sampleDir string) error {
    samples := map[string]string{
        "bd": "bd.wav",     // Kick
        "sd": "sd.wav",     // Snare
        "hh": "hh.wav",     // Hi-hat
        "cp": "cp.wav",     // Clap
        "oh": "oh.wav",     // Open hat
        "ch": "ch.wav",     // Closed hat
        "rs": "rs.wav",     // Rimshot
        "lt": "lt.wav",     // Low tom
        "mt": "mt.wav",     // Mid tom
        "ht": "ht.wav",     // High tom
    }
    
    for name, filename := range samples {
        path := filepath.Join(sampleDir, filename)
        if err := player.LoadSample(name, path); err != nil {
            return err
        }
    }
    
    return nil
}
```

**Sample Library:**
```
data/samples/
├── bd.wav       (kick drum)
├── sd.wav       (snare)
├── hh.wav       (hi-hat closed)
├── oh.wav       (hi-hat open)
├── ch.wav       (closed hat)
├── cp.wav       (clap)
├── rs.wav       (rimshot)
├── lt.wav       (low tom)
├── mt.wav       (mid tom)
├── ht.wav       (high tom)
└── README.md    (sample credits + license)
```

**Integration with Live Mode:**
```go
// internal/app/live/live_mode.go
func (m Model) Init() tea.Cmd {
    player, _ := audio.NewSimplePlayer(m.bpm)
    audio.LoadDefaultSamples(player, "data/samples")
    m.player = player
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "ctrl+p" {
            if m.player.playing {
                m.player.Stop()
            } else {
                pattern, _ := pattern.Parse(m.editor.Value())
                m.player.PlayPattern(pattern)
            }
        }
    }
    return m, nil
}
```

---

### 1.2 Quick Chord Picker

**Purpose:** Instant chord suggestions to unstick writer's block

**Trigger:** `Ctrl+F` (for "Fifth")

**UI:**
```
┌─ Quick Chords ────────────────────────────┐
│ Mood: [1] Happy [2] Sad [3] Tense [4] Chill│
│                                            │
│ Common Progressions:                       │
│                                            │
│ > 1. C-G-Am-F      (Pop Classic)          │
│   2. Am-F-C-G      (Emotional Journey)    │
│   3. Em-C-G-D      (Sad Ballad)           │
│   4. Dm-G7-Cmaj7   (Jazz Smooth)          │
│   5. C-F-G-C       (Happy Major)          │
│                                            │
│ [R] Random  [1-9] Select  [Esc] Cancel    │
└────────────────────────────────────────────┘
```

**Features:**
- 20 preset progressions (no calculation)
- Mood filters: happy, sad, tense, chill
- Random suggestion button
- One-click insert to cursor
- Genre hints (Pop, Jazz, Blues, Folk)

**Data File:**
```json
// data/chord_progressions.json
{
  "progressions": [
    {
      "name": "Pop Classic",
      "chords": ["C", "G", "Am", "F"],
      "mood": "happy",
      "description": "Upbeat, singalong friendly"
    }
  ]
}
```

**Implementation Files:**
```
internal/ui/quick_chord_picker.go    (UI component)
internal/ui/theory.go                (add quick mode)
data/chord_progressions.json         (presets)
```

---

### 1.3 Simple BPM Tapper

**Purpose:** Quick tempo detection while jamming

**Trigger:** `Ctrl+T` (Tempo)

**UI:**
```
┌─ BPM Tapper ──────────────────┐
│                                │
│   Tap spacebar to detect BPM  │
│                                │
│        Current: 120 BPM        │
│        Taps: █ █ █ █          │
│                                │
│   [Space] Tap  [Enter] Set    │
└────────────────────────────────┘
```

**Features:**
- Tap spacebar to detect tempo
- Average last 4-8 taps
- Display current BPM
- Set BPM for current pattern
- Simple, fast, lightweight

**Keep it simple:**
- No complex timing features → wave.sh
- No subdivision trainer → wave.sh
- Just: tap to get BPM

**Implementation Files:**
```
internal/ui/bpm_tapper.go
internal/app/tempo/simple_tapper.go
```

---

### 1.4 Export Pattern Data to wave.sh

**Purpose:** Send pattern DATA to wave.sh for audio production

**IMPORTANT:** This exports pattern CODE, not audio files. wave.sh generates the actual sound.

**Commands:**
```bash
noise export pattern --output my-idea.json
noise export draft --output lyrics.md
```

**Export Format (JSON):**
```json
{
  "type": "noise-pattern",
  "version": "1.0",
  "metadata": {
    "title": "Funky Groove",
    "bpm": 120,
    "created": "2025-10-18T14:30:00Z",
    "exported_at": "2025-10-18T15:45:00Z"
  },
  "patterns": [
    "s('bd sd bd sd')",
    "s('hh*8').gain(0.5)",
    "n('0 3 7 5').speed(1.2)"
  ],
  "lyrics": "Optional lyrics text here",
  "chords": ["C", "G", "Am", "F"],
  "notes": "Optional production notes"
}
```

**What this JSON contains:**
- Pattern code (text)
- Tempo/BPM data
- Chord progressions
- Lyrics
- Metadata

**What it does NOT contain:**
- Audio files
- Rendered sound
- WAV/MP3 data

**wave.sh generates the audio from this data**

**Implementation Files:**
```
cmd/export.go                    (CLI command)
internal/export/pattern.go       (export logic)
internal/export/format.go        (JSON schema)
internal/export/types.go         (data structures)
```

---

## Part 2: AI Right-Sizing

### Current Problem
Enhancement #2 AI is too complex for rapid capture:
- ❌ 4 specialized agents (continuation, variation, quality, brainstorm)
- ❌ Complex prompts (500+ tokens each)
- ❌ Slow responses (2-5 seconds)
- ❌ Deep critique scoring (overkill for sketching)

### New Approach: Single Lightweight Agent

**Consolidation:**
```
OLD (4 agents):
├─ ContinuationAgent
├─ VariationAgent
├─ QualityAgent
└─ RapidBrainstormAgent

NEW (1 agent, 4 modes):
└─ QuickIdeaAgent (LIGHTWEIGHT 7B MODEL)
   ├─ unstick mode
   ├─ spark mode
   ├─ tweak mode
   └─ check mode
```

### QuickIdeaAgent Modes

**1. Unstick Mode** (Ctrl+G)
- Generate next line suggestions
- 3 options, 8-12 syllables each
- Context: last 3-5 lines
- Response time: <2 sec

**2. Spark Mode** (Ctrl+R)
- Random creative starting point
- 3 first-line suggestions
- Based on theme/mood keyword
- Response time: <2 sec

**3. Tweak Mode** (Ctrl+V)
- Rewrite selected line 3 ways
- Keep similar meaning
- Different phrasing
- Response time: <2 sec

**4. Check Mode** (Ctrl+Shift+C)
- Quick quality check
- Returns: STRONG / OKAY / WEAK
- One 5-word improvement tip
- Response time: <1 sec

### Technical Optimizations

**Speed Improvements:**
1. Shorter prompts (100 tokens vs 500)
2. Lower temperature (0.7 vs 0.9)
3. Smaller model (7B vs 13B)
4. No streaming (complete response)
5. 2-second timeout
6. Cached system prompts

**Implementation Files:**
```
internal/app/ai/quick_idea_agent.go     (single agent)
internal/app/ai/prompts.go              (4 mode prompts)
internal/app/ai/types.go                (request/response)
```

### Archive Old Agents
```
docs/archive/ai/continuation_agent.go
docs/archive/ai/variation_agent.go
docs/archive/ai/quality_agent.go
docs/archive/ai/rapid_brainstorm_agent.go
```

---

## Part 3: Theme System Integration

### Purpose
All Kyanite tools share the same 12 themes + custom theme support.

### Core Themes (12)
1. Slate Mist
2. Violet Dusk  
3. Amber Night (default)
4. Molten Gold
5. Clay Roads
6. Iron Storm
7. Jade Tide
8. Sunset Ember
9. Forest Whisper
10. Electric Bloom
11. Plasma Pulse
12. Indigo Depths
13. Sage Meadow

### Implementation
**Create:** `internal/theme/` package (will migrate to `kyanite-common` later)

```go
// internal/theme/types.go
type Theme struct {
    Name       string
    Primary    lipgloss.Color
    Secondary  lipgloss.Color
    Accent     lipgloss.Color
    Background lipgloss.Color
    Text       lipgloss.Color
    Success    lipgloss.Color
    Warning    lipgloss.Color
    Error      lipgloss.Color
}
```

```go
// internal/theme/registry.go
var Registry = map[string]Theme{
    "slate-mist":     SlateMist,
    "violet-dusk":    VioletDusk,
    "amber-night":    AmberNight,
    // ... all 12 themes
}
```

### Custom Themes
**Format:** `~/.config/noise/themes/my-theme.toml`

```toml
name = "My Theme"

[colors]
primary    = "#FF6B9D"
secondary  = "#C44569"
accent     = "#FFC312"
background = "#1e1e2e"
text       = "#FFEAA7"
success    = "#55E6C1"
warning    = "#FFA502"
error      = "#EA2027"
```

### Runtime Switching
**Keybinding:** `Ctrl+Shift+T`

**UI:**
```
┌─ Select Theme ────────────┐
│ > Amber Night            │
│   Slate Mist             │
│   Violet Dusk            │
│   ...                    │
│ Custom:                  │
│   My Theme               │
└───────────────────────────┘
```

### Configuration
**File:** `~/.config/noise/config.toml`

```toml
[preferences]
theme = "amber-night"
```

### Implementation Files
```
internal/theme/types.go           (Theme struct)
internal/theme/registry.go        (12 themes)
internal/theme/manager.go         (runtime switching)
internal/theme/custom.go          (load .toml)
internal/theme/validator.go       (WCAG contrast check)
internal/config/theme_loader.go   (persistence)
```

---

## Part 4: Documentation Overhaul

### Create New Docs

**1. ROADMAP.md** (Replace old PRD)
```markdown
# noise.sh Roadmap

## Vision
Fast, simple, exportable.

## Current Status (v2.0)
✅ Editor complete
✅ Theme system complete
🔄 Enhancement #6 in progress

## Active Development (3-4 weeks)
- Live coding mode
- Quick chord picker
- AI right-sizing
- Export to wave.sh

## Future Phases
- Phase 2: Mobile PWA companion
- Phase 3: i18n support
- Phase 4: Plugin system

## Deferred Features
- Complex AI critique → Later
- Full music theory → Moved to wave.sh
- Quality scoring → Not needed for capture

## Integration
noise.sh (capture) → export → wave.sh (refine)
```

**2. ARCHITECTURE.md**
```markdown
# noise.sh Architecture

## Purpose
Rapid capture tool for musical ideas.
Export to wave.sh for refinement.

## Stack
- Go 1.21+
- Bubble Tea (TUI)
- Lipgloss (styling)
- Ollama (local AI - 7B models)
- SQLite (storage)

## Principles
1. Speed > completeness
2. Capture > polish
3. Export > perfection
4. Simple > sophisticated

## Modules
- Editor: Split-pane, markdown, reusable
- Live Mode: Pattern parser, visualization
- AI: Single QuickIdeaAgent
- Theory: Quick chord picker
- Export: JSON for wave.sh

## Integration
noise.sh → wave.sh → production
```

**3. DEPRECATED.md**
```markdown
# Deprecated Documentation

## Archived Files (Old Vision)
- ~~PRD.md~~ → Replaced by ROADMAP.md
- ~~TDD.md~~ → Replaced by ARCHITECTURE.md
- ~~Development_Roadmap.md~~ → Replaced by ROADMAP.md
- ~~enhancement_02_rapid.txt~~ → Superseded by Enhancement #6

## Why Archived
Vision changed: comprehensive suite → rapid capture tool

## Current Docs
- ROADMAP.md - Vision and plan
- ARCHITECTURE.md - Technical overview
- enhancement_06.md - Current work
```

### File Operations
```bash
# Create archive directory
mkdir -p docs/archive/

# Move old docs
mv docs/PRD.md docs/archive/
mv docs/TDD.md docs/archive/
mv docs/Development_Roadmap.md docs/archive/

# Create new docs
touch docs/ROADMAP.md
touch docs/ARCHITECTURE.md
touch docs/DEPRECATED.md
```

---

## Part 4: Implementation Plan

### Week 1: Live Coding Foundation + Basic Audio

**Goals:**
- Pattern parser working
- Live mode launches
- Basic visualization
- **Simple audio playback working**

**Tasks:**
1. Create pattern parser (simple BNF grammar)
2. Implement live mode command
3. Reuse editor for pattern input
4. Add basic visualization (right pane)
5. **Implement simple audio player (oto/v2)**
6. **Load basic drum kit samples**
7. **Add play/pause controls (Ctrl+P)**
8. Save/load `.noise` files

**Deliverables:**
- `noise live` command works
- Can type patterns
- Visual display updates
- **Ctrl+P plays/pauses patterns**
- **Basic drum sounds work**
- Save patterns

**Testing:**
- [ ] Parser handles valid syntax
- [ ] Parser rejects invalid syntax
- [ ] Visualization updates in real-time
- [ ] **Audio plays without glitches**
- [ ] **Samples load correctly**
- [ ] **Play/pause responsive**
- [ ] Files save/load correctly

---

### Week 2: Quick Tools, Export & Themes

**Goals:**
- Chord picker working
- BPM tapper working
- Export system complete
- **12 themes implemented**
- **Custom theme loading**

**Tasks:**
1. Implement quick chord picker UI
2. Load chord progressions data
3. Implement BPM tapper
4. Add export commands
5. Create JSON export format
6. **Implement theme system (12 themes)**
7. **Add theme switcher (Ctrl+Shift+T)**
8. **Add custom theme loader**
9. Test wave.sh import (manual)

**Deliverables:**
- Ctrl+F shows chord picker
- Can select/insert progressions
- Ctrl+T shows BPM tapper
- Tapping detects tempo accurately
- `noise export pattern` works
- JSON format documented
- **All 12 themes working**
- **Theme switching instant**
- **Custom themes load from ~/.config/noise/themes/**

**Testing:**
- [ ] Chord picker UI responds
- [ ] Mood filters work
- [ ] Random selection works
- [ ] BPM tapper calculates correctly
- [ ] Tapper UI is responsive
- [ ] Export produces valid JSON
- [ ] wave.sh can read format
- [ ] **All 12 themes render correctly**
- [ ] **Theme switching works (Ctrl+Shift+T)**
- [ ] **Custom themes load**
- [ ] **Theme persists across restarts**

---

### Week 3: AI Right-Sizing

**Goals:**
- Single agent working
- All 4 modes functional
- Response times <2 sec

**Tasks:**
1. Create QuickIdeaAgent (single file)
2. Write 4 mode prompts
3. Remove old agent files
4. Update keybindings
5. Test response times

**Deliverables:**
- Ctrl+G unstick works
- Ctrl+R spark works
- Ctrl+V tweak works
- Ctrl+Shift+C check works
- All responses <2 sec

**Testing:**
- [ ] Unstick generates 3 suggestions
- [ ] Spark returns creative prompts
- [ ] Tweak rephrases correctly
- [ ] Check gives rating + tip
- [ ] No timeouts or hangs

---

### Week 4: Documentation & Polish

**Goals:**
- Clean documentation
- All features tested
- Ready to ship

**Tasks:**
1. Create ROADMAP.md
2. Create ARCHITECTURE.md
3. Create DEPRECATED.md
4. Move old files to archive/
5. Update README.md
6. Test all features end-to-end

**Deliverables:**
- Documentation unified
- No confusing old docs
- README reflects new vision
- All features working

**Testing:**
- [ ] ROADMAP clear and accurate
- [ ] ARCHITECTURE correct
- [ ] No outdated references
- [ ] README updated
- [ ] Full integration test passes

---

## Part 5: Technical Specifications

### 5.1 Pattern Language Grammar (BNF)

```bnf
<pattern>       ::= <function_call> | <comment> | <empty_line>
<function_call> ::= <identifier> "(" <argument> ")" [ "." <method_chain> ]
<method_chain>  ::= <method> "(" <argument> ")" [ "." <method_chain> ]

<identifier>    ::= "s" | "n" | "cc" | "gain" | "pan" | "speed"
<method>        ::= "gain" | "pan" | "speed" | "delay" | "room"
<argument>      ::= <string> | <number> | <pattern_expr>

<string>        ::= '"' <any_chars> '"'
<number>        ::= <digit>+ [ "." <digit>+ ]
<pattern_expr>  ::= <value> [ <operator> <pattern_expr> ]
<operator>      ::= "*" | "/" | "," | " "

<comment>       ::= "--" <any_chars>
```

**Valid Examples:**
```javascript
s("bd sd bd sd")              // Sample pattern
s("hh*8").gain(0.5)          // With effect chain
n("0 3 7").speed(1.2)        // Note pattern
-- This is a comment
```

---

### 5.2 Data Structures

```go
// internal/pattern/types.go
package pattern

type Pattern struct {
    ID        string    `json:"id"`
    Lines     []Line    `json:"lines"`
    BPM       int       `json:"bpm"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Line struct {
    Raw        string   `json:"raw"`
    Function   string   `json:"function"`
    Argument   string   `json:"argument"`
    Effects    []Effect `json:"effects"`
    LineNumber int      `json:"line_number"`
}

type Effect struct {
    Name  string  `json:"name"`
    Value float64 `json:"value"`
}
```

```go
// internal/ui/chord_types.go
package ui

type ChordProgression struct {
    Name        string   `json:"name"`
    Chords      []string `json:"chords"`
    Mood        string   `json:"mood"`
    Description string   `json:"description"`
}
```

```go
// internal/export/types.go
package export

type NoiseExport struct {
    Type     string         `json:"type"`
    Version  string         `json:"version"`
    Metadata ExportMetadata `json:"metadata"`
    Patterns []string       `json:"patterns"`
    Lyrics   string         `json:"lyrics,omitempty"`
    Chords   []string       `json:"chords,omitempty"`
    Notes    string         `json:"notes,omitempty"`
}

type ExportMetadata struct {
    Title      string    `json:"title"`
    BPM        int       `json:"bpm"`
    Created    time.Time `json:"created"`
    ExportedAt time.Time `json:"exported_at"`
}
```

---

### 5.3 AI Agent Implementation

```go
// internal/app/ai/quick_idea_agent.go
package ai

type QuickIdeaAgent struct {
    client  *ollama.Client
    model   string        // "qwen2.5:3b" - ULTRA-LIGHTWEIGHT 3B MODEL
    timeout time.Duration // 2 seconds max
}

type QuickRequest struct {
    Mode    string // "unstick", "spark", "tweak", "check"
    Context string
    Options map[string]interface{}
}

type QuickResponse struct {
    Suggestions  []string
    ResponseTime time.Duration
}

func NewQuickIdeaAgent() *QuickIdeaAgent {
    return &QuickIdeaAgent{
        client:  ollama.NewClient(),
        model:   "qwen2.5:3b", // Ultra-lightweight 3B model
        timeout: 2 * time.Second,
    }
}

func (a *QuickIdeaAgent) Generate(req QuickRequest) (*QuickResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
    defer cancel()
    
    prompt := a.buildPrompt(req)
    start := time.Now()
    
    result, err := a.client.Generate(ctx, &ollama.GenerateRequest{
        Model:  a.model,
        Prompt: prompt,
        Options: map[string]interface{}{
            "temperature": 0.7,  // Lower = more consistent
            "top_p":       0.9,
            "num_predict": 100,  // Short responses
        },
    })
    
    if err != nil {
        return nil, err
    }
    
    return &QuickResponse{
        Suggestions:  a.parseResponse(result.Response),
        ResponseTime: time.Since(start),
    }, nil
}
```

---

### 5.4 AI Prompt Templates

```go
// internal/app/ai/prompts.go
package ai

const (
    UnstickPrompt = `Rapid songwriting assistant. Generate next line.

Context:
%s

Section: %s
Style: Concrete imagery, conversational, no clichés

Generate 3 options (8-12 syllables):
1.
2.
3.

Lines only, no explanations.`

    SparkPrompt = `Generate 3 creative song starting points.

Theme: %s

Format:
1. [First line]
2. [First line]
3. [First line]

Concrete images only.`

    TweakPrompt = `Rewrite this line 3 ways:

Original: %s

Variations:
1.
2.
3.

No explanations.`

    CheckPrompt = `Rate this lyric line:

"%s"

Reply ONE word: STRONG, OKAY, or WEAK
Then one 5-word tip.

Format:
RATING
[tip]`
)
```

---

### 5.5 Chord Progressions Data

```json
// data/chord_progressions.json
{
  "progressions": [
    {
      "name": "Pop Classic",
      "chords": ["C", "G", "Am", "F"],
      "mood": "happy",
      "description": "Upbeat, singalong friendly"
    },
    {
      "name": "Emotional Journey",
      "chords": ["Am", "F", "C", "G"],
      "mood": "sad",
      "description": "Melancholic, introspective"
    },
    {
      "name": "Blues Standard",
      "chords": ["E7", "A7", "E7", "B7"],
      "mood": "chill",
      "description": "Classic blues progression"
    },
    {
      "name": "Tension Builder",
      "chords": ["Dm", "G7", "Cmaj7", "Fmaj7"],
      "mood": "tense",
      "description": "Jazz-influenced"
    },
    {
      "name": "Folk Simplicity",
      "chords": ["G", "C", "D"],
      "mood": "chill",
      "description": "Simple, acoustic-friendly"
    },
    {
      "name": "Minor Drama",
      "chords": ["Am", "Dm", "E", "Am"],
      "mood": "tense",
      "description": "Dark, dramatic"
    },
    {
      "name": "Happy Major",
      "chords": ["C", "F", "G", "C"],
      "mood": "happy",
      "description": "Bright, optimistic"
    },
    {
      "name": "Sad Ballad",
      "chords": ["Em", "C", "G", "D"],
      "mood": "sad",
      "description": "Emotional storytelling"
    },
    {
      "name": "Jazz Smooth",
      "chords": ["Cmaj7", "Am7", "Dm7", "G7"],
      "mood": "chill",
      "description": "Sophisticated lounge"
    },
    {
      "name": "Rock Energy",
      "chords": ["A", "D", "E", "A"],
      "mood": "happy",
      "description": "Driving, powerful"
    }
  ]
}
```

---

### 5.6 Keybinding Updates

```go
// internal/ui/editor/editor.go - Add to KeyMap
type KeyMap struct {
    // ... existing
    LiveMode    key.Binding
    QuickChords key.Binding
    BPMTapper   key.Binding
    AIUnstick   key.Binding
    AITweak     key.Binding
    AISpark     key.Binding
    AICheck     key.Binding
}

var DefaultKeyMap = KeyMap{
    // ... existing
    LiveMode: key.NewBinding(
        key.WithKeys("ctrl+l"),
        key.WithHelp("ctrl+l", "live music coding"),
    ),
    QuickChords: key.NewBinding(
        key.WithKeys("ctrl+f"),
        key.WithHelp("ctrl+f", "quick chords"),
    ),
    BPMTapper: key.NewBinding(
        key.WithKeys("ctrl+t"),
        key.WithHelp("ctrl+t", "bpm tapper"),
    ),
    AIUnstick: key.NewBinding(
        key.WithKeys("ctrl+g"),
        key.WithHelp("ctrl+g", "ai unstick"),
    ),
    AITweak: key.NewBinding(
        key.WithKeys("ctrl+v"),
        key.WithHelp("ctrl+v", "ai tweak"),
    ),
    AISpark: key.NewBinding(
        key.WithKeys("ctrl+r"),
        key.WithHelp("ctrl+r", "ai spark"),
    ),
    AICheck: key.NewBinding(
        key.WithKeys("ctrl+shift+c"),
        key.WithHelp("ctrl+shift+c", "ai check"),
    ),
}
```

---

### 5.7 Dependencies

```bash
# Install required packages
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/ollama/ollama-go@latest
go get github.com/spf13/cobra@latest
go get modernc.org/sqlite@latest

# Audio libraries (lightweight)
go get github.com/hajimehoshi/oto/v2@latest
go get github.com/hajimehoshi/go-mp3@latest
```

```go
// go.mod
module kyanite/noise

go 1.21

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/ollama/ollama-go v0.1.0
    github.com/spf13/cobra v1.8.0
    modernc.org/sqlite v1.27.0
    
    // Audio (lightweight, pure Go)
    github.com/hajimehoshi/oto/v2 v2.4.0
    github.com/hajimehoshi/go-mp3 v0.3.4
)
```

---

## Part 6: Testing & Validation

### Testing Checklist

**Live Coding Mode:**
- [ ] `noise live` launches successfully
- [ ] Pattern input works
- [ ] Visualization updates in real-time
- [ ] **Audio plays patterns correctly**
- [ ] **Play/pause (Ctrl+P) works**
- [ ] **No audio glitches or latency**
- [ ] Syntax errors displayed clearly
- [ ] Save/load works correctly
- [ ] Export to JSON works

**Quick Chord Picker:**
- [ ] Ctrl+F opens picker
- [ ] Mood filters work
- [ ] Can select progression
- [ ] Chords insert at cursor
- [ ] Random button works

**BPM Tapper:**
- [ ] Ctrl+T opens tapper
- [ ] Spacebar registers taps
- [ ] BPM calculation accurate
- [ ] Sets pattern tempo
- [ ] UI responsive

**AI Right-Sizing:**
- [ ] Ctrl+G responds <2 sec
- [ ] Ctrl+R generates sparks
- [ ] Ctrl+V tweaks lines
- [ ] Ctrl+Shift+C checks quality
- [ ] All modes return 3 suggestions

**Export System:**
- [ ] `noise export pattern` works
- [ ] JSON is valid format
- [ ] Metadata included
- [ ] wave.sh can import (future test)

**Documentation:**
- [ ] ROADMAP.md clear
- [ ] ARCHITECTURE.md accurate
- [ ] Old docs archived
- [ ] No confusing references
- [ ] README updated

---

### Performance Targets

| Metric | Target | Critical |
|--------|--------|----------|
| Pattern parsing | <100ms | Yes |
| AI unstick | <2 sec | Yes |
| AI spark | <2 sec | Yes |
| AI tweak | <2 sec | Yes |
| AI check | <1 sec | No |
| Chord picker | <50ms | Yes |
| Export JSON | <200ms | No |

---

## Part 7: File Structure Summary

### New Files
```
cmd/live.go
cmd/export.go
internal/app/live/live_mode.go
internal/pattern/parser.go
internal/pattern/types.go
internal/ui/editor/pattern_view.go
internal/ui/quick_chord_picker.go
internal/ui/bpm_tapper.go
internal/app/tempo/simple_tapper.go
internal/audio/simple_player.go
internal/audio/sample_loader.go
internal/audio/mixer.go
internal/theme/types.go
internal/theme/registry.go
internal/theme/manager.go
internal/theme/custom.go
internal/theme/validator.go
internal/config/theme_loader.go
internal/app/ai/quick_idea_agent.go
internal/app/ai/prompts.go
internal/app/ai/types.go
internal/export/pattern.go
internal/export/format.go
internal/export/types.go
data/chord_progressions.json
data/samples/bd.wav
data/samples/sd.wav
data/samples/hh.wav
data/samples/cp.wav
data/samples/oh.wav
data/samples/ch.wav
data/samples/rs.wav
data/samples/lt.wav
data/samples/mt.wav
data/samples/ht.wav
data/samples/README.md
docs/ROADMAP.md
docs/ARCHITECTURE.md
docs/DEPRECATED.md
docs/THEMES.md
```

### Modified Files
```
internal/ui/theory.go             (add quick mode)
internal/ui/editor/editor.go      (update keybindings)
cmd/quick.go                      (use QuickIdeaAgent)
README.md                         (update vision)
```

### Archived Files
```
docs/archive/PRD.md
docs/archive/TDD.md
docs/archive/Development_Roadmap.md
docs/archive/enhancement_02_rapid.txt
internal/app/ai/continuation_agent.go
internal/app/ai/variation_agent.go
internal/app/ai/quality_agent.go
internal/app/ai/rapid_brainstorm_agent.go
```

---

## Part 8: Notes for Kilocode (AI Coding Agent)

### Context
- noise.sh purpose changed: comprehensive tool → rapid capture
- Enhancement #2 was 70% done but too complex
- This enhancement right-sizes everything for speed

### Implementation Priorities
1. **Speed over sophistication**
2. **Simple over complete**
3. **Export over perfection**
4. **Clear docs over comprehensive**

### Critical Rules
- **Reuse editor** - Don't rebuild split-pane
- **Single AI agent** - Not four separate agents
- **Ultra-lightweight 3B models ONLY** - qwen2.5:3b or phi3:3.8b (MAX 3B, never 7B+)
- **Preset chords** - No Circle of Fifths calculation
- **Simple audio in noise.sh** - Sample playback only, basic timing
- **Advanced audio in wave.sh** - Synthesis, mixing, mastering, effects
- **Clean export** - wave.sh handles refinement
- **Fast responses** - AI <2 sec, parser <100ms, audio <50ms latency
- **Data files** - All in `data/` directory
- **Internal packages** - All in `internal/`
- **Audio division:**
  - **noise.sh:** Play samples to hear ideas (kick, snare, hat)
  - **wave.sh:** Full synthesis, mixing, professional audio
- **Pattern data export** - JSON contains code/data, NOT rendered audio files

### Week-by-Week Execution
1. **Week 1:** Pattern parser + live mode + data structures
2. **Week 2:** Chord picker + export + JSON format
3. **Week 3:** QuickIdeaAgent + prompts + keybindings
4. **Week 4:** Documentation + testing + polish

### Validation Requirements
- Every feature must have tests
- AI responses timed and logged
- Export format schema-validated
- Documentation updated simultaneously
- No feature merges without tests

---

## Success Criteria

**After Enhancement #6 completes:**

✅ User sketches musical patterns in <60 seconds  
✅ AI suggestions return in <2 seconds  
✅ Chord picker unsticks writer's block instantly  
✅ Export to wave.sh works seamlessly  
✅ Documentation is unified and clear  
✅ No confusing legacy documentation  
✅ All features tested and working  

**Vision Achieved:**  
noise.sh = Fast sketchbook for musical ideas

---

**Status:** READY FOR IMPLEMENTATION  
**Owner:** Simon | Kyanite Suite  
**Implementation:** Hand to Kilocode, execute week by week  
**Questions:** Ping Simon on Discord