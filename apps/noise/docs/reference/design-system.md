# noise.sh Design System & Styling Reference
## For Kilocode AI Coding Agent

**Version:** 1.0  
**Date:** October 17, 2025  
**Purpose:** Complete styling reference for AI-assisted TUI development

---

## Table of Contents

1. [Selected Theme: Violet Dusk](#selected-theme-violet-dusk)
2. [Color Palette](#color-palette)
3. [Icon System](#icon-system)
4. [Lipgloss Style Definitions](#lipgloss-style-definitions)
5. [Typography System](#typography-system)
6. [Component Library](#component-library)
7. [Border & Box Styles](#border--box-styles)
8. [Animation Patterns](#animation-patterns)
9. [Layout System](#layout-system)
10. [Usage Examples](#usage-examples)
11. [Best Practices](#best-practices)

---

## Selected Theme: Midnight Jazz

**Aesthetic:** Sophisticated, mysterious, creative  
**Mood:** Deep blues and purples with gold accents - like a late-night jazz club  
**Style:** Elegant with bold highlights  
**Use Case:** Professional portfolio piece with memorable visual impact

### Theme Rationale

- **Portfolio Impact:** Sophisticated and memorable, stands out in screenshots
- **Readability:** High contrast ensures long coding sessions won't cause eye strain
- **Creativity:** Purple/gold combination evokes artistic, musical atmosphere
- **Professional:** Not too "loud" - tasteful and polished

---

## Color Palette

### Primary Colors

```go
// Core brand colors - use for main UI elements
const (
    Primary   = lipgloss.Color("#9D84B7") // Soft purple - main brand
    Secondary = lipgloss.Color("#5E4B8B") // Deep purple - secondary actions
    Accent    = lipgloss.Color("#F4D03F") // Gold - highlights, important elements
)
```

### Functional Colors

```go
// Status and feedback colors
const (
    Success = lipgloss.Color("#52D3AA") // Mint green - success states
    Warning = lipgloss.Color("#FFA500") // Orange - warnings
    Error   = lipgloss.Color("#FF6347") // Tomato - errors
    Info    = lipgloss.Color("#87CEEB") // Sky blue - info messages
)
```

### Background & Text

```go
// Base colors for backgrounds and text
const (
    Background = lipgloss.Color("#0A0E27") // Deep navy - main background
    
    // Text colors
    TextPrimary   = lipgloss.Color("#E8DFF5") // Light lavender - main text
    TextSecondary = lipgloss.Color("#9D84B7") // Soft purple - secondary text
    TextMuted     = lipgloss.Color("#5E4B8B") // Deep purple - muted text
    TextAccent    = lipgloss.Color("#F4D03F") // Gold - emphasized text
)
```

### Extended Palette (for gradients & variations)

```go
// Additional shades for gradients and variety
const (
    Purple1 = lipgloss.Color("#9D84B7") // Lightest purple
    Purple2 = lipgloss.Color("#5E4B8B") // Medium purple
    Purple3 = lipgloss.Color("#3D2C8D") // Darkest purple
    
    Gold1   = lipgloss.Color("#F4D03F") // Bright gold
    Gold2   = lipgloss.Color("#D4AF37") // Muted gold
    
    Dark1   = lipgloss.Color("#0A0E27") // Main background
    Dark2   = lipgloss.Color("#1A1E37") // Lighter background
    Dark3   = lipgloss.Color("#2A2E47") // Border/divider
)
```

### Color Usage Guide

| Color | Use For | Don't Use For |
|-------|---------|---------------|
| **Primary (#9D84B7)** | Active focus states, selected items, primary buttons | Background fills, large text blocks |
| **Secondary (#5E4B8B)** | Secondary buttons, borders, inactive states | Primary actions, success messages |
| **Accent (#F4D03F)** | Call-to-action, highlights, section headers | Body text, backgrounds (too bright) |
| **Success (#52D3AA)** | Save confirmations, quality scores >70, completion | Warnings, neutral states |
| **Background (#0A0E27)** | Main canvas, card backgrounds | Text (invisible) |
| **TextPrimary (#E8DFF5)** | Body text, lyrics, user content | Backgrounds (too light) |

---

## Icon System

noise.sh uses a configurable icon system (`internal/ui/icons/`) that supports three rendering modes to ensure compatibility across different terminal environments.

### Icon Styles

| Style | Description | Use Case |
|-------|-------------|----------|
| **ASCII** (default) | Plain ASCII characters like `[OK]`, `[X]`, `[-]` | Universal compatibility, works everywhere |
| **Unicode** | Standard Unicode symbols like `✓`, `✗`, `•` | Modern terminals with UTF-8 support |
| **NerdFont** | Patched font icons using Nerd Fonts | Users with installed Nerd Fonts |

### Usage

```go
import "github.com/kyanite/noise/internal/ui/icons"

// Get current icon set (defaults to ASCII)
i := icons.Current()

// Use icons in your UI
fmt.Println(i.Success + " Saved!")      // [OK] Saved!
fmt.Println(i.Bullet + " List item")    // - List item

// Switch styles (typically from user config)
icons.SetStyle(icons.StyleUnicode)      // Use Unicode symbols
icons.SetStyle(icons.StyleNerdFont)     // Use Nerd Font icons
icons.SetStyle(icons.StyleASCII)        // Back to ASCII (default)
```

### Available Icons

```go
// Status indicators
Success, Error, Warning, Info, Recording, Processing

// Selection
CheckOn, CheckOff, RadioOn, RadioOff

// Navigation  
ArrowLeft, ArrowRight, ArrowUp, ArrowDown, ChevronL, ChevronR

// Lists
Bullet, BulletAlt, Separator

// Progress bars
ProgressFull, ProgressEmpty

// Presence/Status
Online, Away, Busy, Offline

// Application-specific
Music, Folder, File, Settings, Help, Search, AI, Mic, Photo, Export, Sync
Stats, Performance, Storage, Tools, Tip
```

### Icon Reference Table

| Icon | ASCII | Unicode | NerdFont |
|------|-------|---------|----------|
| Success | `[OK]` | `✓` | `` |
| Error | `[X]` | `✗` | `` |
| Warning | `[!]` | `⚠` | `` |
| Bullet | `-` | `•` | `` |
| CheckOn | `[x]` | `☑` | `` |
| CheckOff | `[ ]` | `☐` | `` |
| ArrowRight | `->` | `→` | `` |
| Music | `[~]` | `♪` | `` |

### String Width Utilities

The icons package includes Unicode-safe string utilities using `rivo/uniseg`:

```go
import "github.com/kyanite/noise/internal/ui/icons"

// Get display width (handles emoji, CJK, combining chars)
width := icons.StringWidth("Hello 世界!")  // Returns 11, not 9

// Safe truncation respecting grapheme clusters
s := icons.Truncate("Hello 世界!", 8, "...")  // Won't break in middle of 世

// Padding with proper width calculation
padded := icons.PadRight("名前", 10)  // Correctly pads to 10 columns
centered := icons.Center("Title", 20) // Centers properly
```

### Configuration

Users can configure their preferred icon style in `config.yaml`:

```yaml
ui:
  icon_style: "ascii"  # Options: "ascii", "unicode", "nerdfonts"
```

### Best Practices

**DO:**
- Default to ASCII for maximum compatibility
- Let users opt-in to Unicode/NerdFont via configuration
- Use `icons.StringWidth()` instead of `len()` for display calculations
- Use `icons.Truncate()` instead of `s[:n]` to avoid breaking characters

**DON'T:**
- Don't hardcode emoji or Unicode symbols directly
- Don't assume `len(s)` equals display width
- Don't truncate strings with `s[:n]` (breaks multi-byte characters)
- Don't auto-detect font support (unreliable)

---

## Lipgloss Style Definitions

### File: `internal/ui/styles/theme.go`

```go
package styles

import "github.com/charmbracelet/lipgloss"

// Colors
var (
    Primary   = lipgloss.Color("#9D84B7")
    Secondary = lipgloss.Color("#5E4B8B")
    Accent    = lipgloss.Color("#F4D03F")
    Success   = lipgloss.Color("#52D3AA")
    Warning   = lipgloss.Color("#FFA500")
    Error     = lipgloss.Color("#FF6347")
    Info      = lipgloss.Color("#87CEEB")
    
    Background    = lipgloss.Color("#0A0E27")
    TextPrimary   = lipgloss.Color("#E8DFF5")
    TextSecondary = lipgloss.Color("#9D84B7")
    TextMuted     = lipgloss.Color("#5E4B8B")
    TextAccent    = lipgloss.Color("#F4D03F")
    
    BorderColor = lipgloss.Color("#2A2E47")
)

// Base Styles
var (
    // Title style for app header
    Title = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        MarginBottom(1).
        Padding(0, 2)
    
    // Subtitle for sections
    Subtitle = lipgloss.NewStyle().
        Foreground(TextSecondary).
        Italic(true).
        MarginBottom(1)
    
    // Main content text
    Text = lipgloss.NewStyle().
        Foreground(TextPrimary)
    
    // Muted text for hints/secondary info
    Muted = lipgloss.NewStyle().
        Foreground(TextMuted)
    
    // Emphasized text
    Emphasis = lipgloss.NewStyle().
        Foreground(TextAccent).
        Bold(true)
)

// Border Styles
var (
    // Standard border for panels
    Border = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(BorderColor).
        Padding(1, 2)
    
    // Active/focused border
    BorderActive = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Primary).
        Padding(1, 2)
    
    // Thick border for emphasis
    BorderThick = lipgloss.NewStyle().
        Border(lipgloss.ThickBorder()).
        BorderForeground(Accent).
        Padding(1, 2)
)

// Button Styles
var (
    // Primary action button
    ButtonPrimary = lipgloss.NewStyle().
        Background(Primary).
        Foreground(lipgloss.Color("#000000")).
        Bold(true).
        Padding(0, 3).
        MarginRight(1)
    
    // Secondary button
    ButtonSecondary = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Secondary).
        Foreground(Secondary).
        Padding(0, 2).
        MarginRight(1)
    
    // Accent/highlight button
    ButtonAccent = lipgloss.NewStyle().
        Background(Accent).
        Foreground(Background).
        Bold(true).
        Padding(0, 3).
        MarginRight(1)
    
    // Disabled button
    ButtonDisabled = lipgloss.NewStyle().
        Background(lipgloss.Color("#2A2E47")).
        Foreground(TextMuted).
        Padding(0, 3).
        MarginRight(1)
)

// Status Styles
var (
    StatusSuccess = lipgloss.NewStyle().
        Foreground(Success).
        Bold(true)
    
    StatusWarning = lipgloss.NewStyle().
        Foreground(Warning).
        Bold(true)
    
    StatusError = lipgloss.NewStyle().
        Foreground(Error).
        Bold(true)
    
    StatusInfo = lipgloss.NewStyle().
        Foreground(Info)
)

// Helper Functions
func Gradient(text string, colors []lipgloss.Color) string {
    var result string
    textRunes := []rune(text)
    colorCount := len(colors)
    
    for i, char := range textRunes {
        if char == ' ' {
            result += " "
            continue
        }
        colorIdx := (i * colorCount) / len(textRunes)
        if colorIdx >= colorCount {
            colorIdx = colorCount - 1
        }
        style := lipgloss.NewStyle().Foreground(colors[colorIdx])
        result += style.Render(string(char))
    }
    
    return result
}

// Purple to gold gradient for titles
func TitleGradient(text string) string {
    colors := []lipgloss.Color{
        Primary,
        lipgloss.Color("#7D6497"),
        lipgloss.Color("#8D7447"),
        Accent,
    }
    return Gradient(text, colors)
}
```

---

## Typography System

### Hierarchy

```go
// Font sizes (relative via padding/spacing, not actual size)
const (
    SizeH1 = 2  // Title padding
    SizeH2 = 1  // Section padding
    SizeH3 = 1  // Subsection padding
    SizeBody = 0 // Normal text
    SizeSmall = 0 // Small text (use Muted style)
)

// Heading styles
var (
    H1 = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        MarginBottom(1).
        Padding(0, SizeH1)
    
    H2 = lipgloss.NewStyle().
        Bold(true).
        Foreground(Secondary).
        MarginBottom(1).
        Padding(0, SizeH2).
        Underline(true)
    
    H3 = lipgloss.NewStyle().
        Bold(true).
        Foreground(TextAccent).
        MarginBottom(1)
)
```

### Text Formatting

```go
var (
    Bold = lipgloss.NewStyle().Bold(true)
    
    Italic = lipgloss.NewStyle().Italic(true)
    
    Underline = lipgloss.NewStyle().Underline(true)
    
    Code = lipgloss.NewStyle().
        Foreground(Accent).
        Background(lipgloss.Color("#1A1E37")).
        Padding(0, 1)
    
    Quote = lipgloss.NewStyle().
        Foreground(TextSecondary).
        Italic(true).
        BorderLeft(true).
        BorderForeground(Primary).
        PaddingLeft(2)
)
```

---

## Component Library

### Editor Components

```go
// Split pane editor - left side (editing)
var EditorPane = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Primary).
    Padding(1, 2)

// Split pane preview - right side (rendered)
var PreviewPane = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Secondary).
    Padding(1, 2)

// Status bar at bottom
var StatusBar = lipgloss.NewStyle().
    Background(lipgloss.Color("#1A1E37")).
    Foreground(TextPrimary).
    Padding(0, 2)

// Cursor indicator
var Cursor = lipgloss.NewStyle().
    Background(Primary).
    Foreground(Background)
```

### List & Menu Items

```go
// Normal list item
var ListItem = lipgloss.NewStyle().
    Foreground(TextPrimary).
    Padding(0, 2)

// Selected list item
var ListItemSelected = lipgloss.NewStyle().
    Background(Primary).
    Foreground(lipgloss.Color("#000000")).
    Bold(true).
    Padding(0, 2)

// List item with icon
func ListItemWithIcon(icon string, text string, selected bool) string {
    iconStyle := lipgloss.NewStyle().Foreground(Accent)
    textStyle := ListItem
    if selected {
        textStyle = ListItemSelected
        iconStyle = iconStyle.Background(Primary)
    }
    return iconStyle.Render(icon) + " " + textStyle.Render(text)
}
```

### Progress & Loading

```go
// Progress bar container
var ProgressBar = lipgloss.NewStyle().
    Width(40).
    Height(1).
    Background(lipgloss.Color("#1A1E37"))

// Progress bar fill
var ProgressFill = lipgloss.NewStyle().
    Background(Primary)

// Spinner text
var SpinnerText = lipgloss.NewStyle().
    Foreground(Primary)

// Loading states
func RenderProgress(percent float64, width int) string {
    fillWidth := int(float64(width) * percent)
    emptyWidth := width - fillWidth
    
    fill := lipgloss.NewStyle().
        Background(Primary).
        Render(strings.Repeat(" ", fillWidth))
    
    empty := lipgloss.NewStyle().
        Background(lipgloss.Color("#1A1E37")).
        Render(strings.Repeat(" ", emptyWidth))
    
    return fill + empty
}
```

### Cards & Panels

```go
// Standard card for grouped content
var Card = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(BorderColor).
    Padding(2, 4).
    MarginBottom(1)

// Highlighted card (for important info)
var CardHighlight = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Accent).
    Padding(2, 4).
    MarginBottom(1)

// Success card
var CardSuccess = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Success).
    Padding(2, 4).
    MarginBottom(1)

// Error card
var CardError = lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(Error).
    Padding(2, 4).
    MarginBottom(1)
```

### Badges & Tags

```go
// Tag for metadata
var Tag = lipgloss.NewStyle().
    Background(Secondary).
    Foreground(TextPrimary).
    Padding(0, 1).
    MarginRight(1)

// Badge for counts/numbers
var Badge = lipgloss.NewStyle().
    Background(Accent).
    Foreground(Background).
    Bold(true).
    Padding(0, 1)

// Status badge (with color coding)
func StatusBadge(status string, color lipgloss.Color) string {
    return lipgloss.NewStyle().
        Background(color).
        Foreground(Background).
        Bold(true).
        Padding(0, 1).
        Render(status)
}
```

---

## Border & Box Styles

### Unicode Box Drawing Characters

```go
// Use these for custom borders and layouts
const (
    // Rounded corners
    CornerTL = "â•­"  // Top-left
    CornerTR = "â•®"  // Top-right
    CornerBL = "â•°"  // Bottom-left
    CornerBR = "â•¯"  // Bottom-right
    
    // Lines
    Horizontal = "â”€"
    Vertical   = "â”‚"
    
    // T-junctions
    JuncT = "â”¬"  // Top T
    JuncB = "â”´"  // Bottom T
    JuncL = "â”œ"  // Left T
    JuncR = "â”¤"  // Right T
    
    // Cross
    Cross = "â”¼"
    
    // Double lines (for emphasis)
    HorizontalDouble = "â•"
    VerticalDouble   = "â•‘"
    CornerTLDouble   = "â•”"
    CornerTRDouble   = "â•—"
    CornerBLDouble   = "â•š"
    CornerBRDouble   = "â•"
)
```

### Border Presets

```go
// Standard rounded border
lipgloss.RoundedBorder()

// Thick border for emphasis
lipgloss.ThickBorder()

// Double border for special sections
lipgloss.DoubleBorder()

// No border (just padding)
lipgloss.NormalBorder()
```

---

## Animation Patterns

### Using Harmonica for Smooth Animations

```go
import "github.com/charmbracelet/harmonica"

// Spring animation for smooth transitions
type AnimatedElement struct {
    spring harmonica.Spring
}

func NewAnimatedElement() AnimatedElement {
    return AnimatedElement{
        spring: harmonica.NewSpring(
            harmonica.FPS(60),  // 60 FPS
            10.0,               // Stiffness
            0.5,                // Damping
        ),
    }
}

// Trigger animation
func (a *AnimatedElement) Trigger() {
    a.spring.SetTarget(1.0)
}

// Update in your Update() function
func (a *AnimatedElement) Update() {
    a.spring.Update()
}

// Get current value (0.0 to 1.0)
func (a *AnimatedElement) Value() float64 {
    return a.spring.Value()
}
```

### Animation Use Cases

```go
// Pulsing indicator (AI thinking)
func PulseIndicator(value float64) string {
    // value cycles 0.0 -> 1.0 -> 0.0
    opacity := int(value * 100)
    color := lipgloss.Color(fmt.Sprintf("#9D84B7%02X", opacity))
    return lipgloss.NewStyle().
        Foreground(color).
        Render("â—")
}

// Slide-in panel
func SlideInPanel(value float64, width int) lipgloss.Style {
    currentWidth := int(float64(width) * value)
    return lipgloss.NewStyle().Width(currentWidth)
}

// Fade in text
func FadeIn(text string, value float64) string {
    opacity := int(value * 255)
    color := lipgloss.Color(fmt.Sprintf("#E8DFF5%02X", opacity))
    return lipgloss.NewStyle().
        Foreground(color).
        Render(text)
}
```

### Common Animation Patterns

```go
// Beat indicator (metronome)
type BeatAnimation struct {
    spring harmonica.Spring
    active bool
}

func (b *BeatAnimation) Trigger() {
    b.spring.SetTarget(1.0)
    b.active = true
}

func (b *BeatAnimation) Update() {
    b.spring.Update()
    if b.spring.Done() {
        b.spring.SetTarget(0.0)
        b.active = false
    }
}

func (b *BeatAnimation) Render() string {
    intensity := b.spring.Value()
    
    // Interpolate between two colors
    if intensity > 0.5 {
        return lipgloss.NewStyle().
            Background(Accent).
            Foreground(Background).
            Render("â—")
    } else {
        return lipgloss.NewStyle().
            Foreground(Primary).
            Render("â—‹")
    }
}
```

---

## Layout System

### Responsive Layouts

```go
// Calculate responsive widths
func CalculateSplitPaneWidth(terminalWidth int) (leftWidth, rightWidth int) {
    leftWidth = terminalWidth / 2
    rightWidth = terminalWidth - leftWidth
    return
}

// Minimum viable size
const (
    MinWidth  = 80
    MinHeight = 24
)

// Check if terminal is large enough
func IsTerminalSizeViable(width, height int) bool {
    return width >= MinWidth && height >= MinHeight
}
```

### Split Layouts

```go
// Horizontal split (side-by-side)
func SplitHorizontal(left, right string) string {
    return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

// Vertical split (stacked)
func SplitVertical(top, bottom string) string {
    return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
}

// Three-column layout
func ThreeColumn(left, center, right string) string {
    return lipgloss.JoinHorizontal(
        lipgloss.Top,
        left,
        center,
        right,
    )
}
```

### Centering & Alignment

```go
// Center content in available space
var Centered = lipgloss.NewStyle().
    Align(lipgloss.Center).
    Width(80)

// Full-width content
var FullWidth = lipgloss.NewStyle().
    Width(lipgloss.Width)

// Right-aligned (for status info)
var RightAlign = lipgloss.NewStyle().
    Align(lipgloss.Right)
```

---

## Usage Examples

### Example 1: Split-Pane Editor

```go
func RenderEditor(content, preview string, activePane int, width, height int) string {
    leftWidth := width / 2
    rightWidth := width - leftWidth
    
    // Style panes based on focus
    leftStyle := EditorPane.Width(leftWidth - 2).Height(height - 4)
    rightStyle := PreviewPane.Width(rightWidth - 2).Height(height - 4)
    
    if activePane == 0 {
        leftStyle = leftStyle.BorderForeground(Primary)
    } else {
        rightStyle = rightStyle.BorderForeground(Primary)
    }
    
    left := leftStyle.Render(content)
    right := rightStyle.Render(preview)
    
    return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}
```

### Example 2: Status Bar with Sections

```go
func RenderStatusBar(wordCount, lineCount int, saved bool, width int) string {
    // Left section: stats
    stats := fmt.Sprintf("Words: %d | Lines: %d", wordCount, lineCount)
    
    // Right section: save status
    var saveStatus string
    if saved {
        saveStatus = StatusSuccess.Render("â— Saved")
    } else {
        saveStatus = StatusWarning.Render("â— Unsaved")
    }
    
    // Center section: hints
    hints := Muted.Render("Ctrl+B: Brainstorm | Ctrl+R: Refine")
    
    // Calculate spacing
    leftStyle := lipgloss.NewStyle().Width(width / 3).Align(lipgloss.Left)
    centerStyle := lipgloss.NewStyle().Width(width / 3).Align(lipgloss.Center)
    rightStyle := lipgloss.NewStyle().Width(width / 3).Align(lipgloss.Right)
    
    left := leftStyle.Render(stats)
    center := centerStyle.Render(hints)
    right := rightStyle.Render(saveStatus)
    
    bar := lipgloss.JoinHorizontal(lipgloss.Top, left, center, right)
    
    return StatusBar.Width(width).Render(bar)
}
```

### Example 3: Quality Score Card

```go
func RenderQualityScore(score *QualityScore) string {
    title := H2.Render("Quality Analysis")
    
    // Score visualization
    scorePercent := score.Total / 70.0
    progressBar := RenderProgress(scorePercent, 40)
    
    scoreText := fmt.Sprintf("%.0f/100", (scorePercent * 100))
    var scoreStyle lipgloss.Style
    if scorePercent >= 0.7 {
        scoreStyle = StatusSuccess
    } else if scorePercent >= 0.5 {
        scoreStyle = StatusWarning
    } else {
        scoreStyle = StatusError
    }
    
    scoreDisplay := scoreStyle.Render(scoreText)
    
    // Individual dimensions
    dimensions := []string{
        fmt.Sprintf("Specificity:    %.0f/10", score.Specificity),
        fmt.Sprintf("Originality:    %.0f/10", score.Originality),
        fmt.Sprintf("Emotional:      %.0f/10", score.EmotionalResonance),
        fmt.Sprintf("Prosody:        %.0f/10", score.Prosody),
        fmt.Sprintf("Coherence:      %.0f/10", score.Coherence),
        fmt.Sprintf("Voice:          %.0f/10", score.VoiceConsistency),
        fmt.Sprintf("Surprise:       %.0f/10", score.SurpriseFactor),
    }
    
    content := lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        progressBar + " " + scoreDisplay,
        "",
        lipgloss.JoinVertical(lipgloss.Left, dimensions...),
    )
    
    return CardHighlight.Render(content)
}
```

### Example 4: AI Streaming Response

```go
func RenderAIResponse(chunks []string, thinking bool) string {
    title := H3.Render("âœ¨ AI Suggestion")
    
    var content string
    if thinking {
        spinner := SpinnerText.Render("â— Thinking...")
        content = spinner
    } else {
        content = Text.Render(strings.Join(chunks, ""))
    }
    
    return Card.Render(lipgloss.JoinVertical(
        lipgloss.Left,
        title,
        "",
        content,
    ))
}
```

---

## Best Practices

### Color Usage

âœ… **DO:**
- Use Primary for active/focused elements
- Use Accent sparingly for important calls-to-action
- Ensure 4.5:1 contrast ratio for readability (WCAG AA)
- Use semantic colors (Success for saves, Error for problems)
- Test in both light and dark terminal themes

âŒ **DON'T:**
- Don't use Accent for large background areas (too bright)
- Don't rely solely on color to convey information
- Don't use more than 3-4 colors in a single view
- Don't use gradients everywhere (visual fatigue)

### Typography

âœ… **DO:**
- Maintain clear hierarchy (H1 â†’ H2 â†’ H3 â†’ Body)
- Use Bold for emphasis, not ALL CAPS
- Use monospace for code, data, numbers
- Keep line length under 80 characters for readability

âŒ **DON'T:**
- Don't use italic for large blocks of text
- Don't use multiple typefaces (terminal is monospace)
- Don't underline everything (reserve for links/headers)

### Animation

âœ… **DO:**
- Keep animations under 500ms for responsiveness
- Use spring physics for natural feel (Harmonica)
- Animate state changes (expand/collapse, fade in/out)
- Target 60 FPS for smooth motion

âŒ **DON'T:**
- Don't animate constantly (battery drain, distraction)
- Don't use animations for critical information
- Don't animate during text input (jarring)

### Layout

âœ… **DO:**
- Support terminal sizes 80x24 and up
- Use percentage-based widths (not fixed)
- Provide minimum viable layout for small screens
- Test on multiple terminal emulators

âŒ **DON'T:**
- Don't assume terminal size
- Don't hardcode widths/heights
- Don't overflow beyond terminal bounds

---

## Quick Reference for Kilocode

### Copy-Paste Color Definitions

```go
// Midnight Jazz Theme - Paste this at the top of any styling file
const (
    Primary   = "#9D84B7"  // Soft purple
    Secondary = "#5E4B8B"  // Deep purple
    Accent    = "#F4D03F"  // Gold
    Success   = "#52D3AA"  // Mint green
    Background = "#0A0E27" // Deep navy
    Text      = "#E8DFF5"  // Light lavender
)
```

### Common Patterns

```go
// Focused element
lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#9D84B7")).
    Padding(1, 2)

// Highlighted button
lipgloss.NewStyle().
    Background(lipgloss.Color("#F4D03F")).
    Foreground(lipgloss.Color("#0A0E27")).
    Bold(true).
    Padding(0, 3)

// Success message
lipgloss.NewStyle().
    Foreground(lipgloss.Color("#52D3AA")).
    Bold(true).
    Render("âœ“ Saved successfully")

// Gradient title
TitleGradient("noise.sh") // Purple â†’ Gold gradient
```

---

**For Kilocode Agent:**  
When generating TUI components, always import and use these color constants. Never hardcode colors. Always use Lipgloss for styling. Prioritize readability and contrast. Keep animations smooth (60 FPS target).

---

**Document Owner:** Simon (Kyanite)  
**Status:** Ready for development  
**Last Updated:** October 17, 2025
