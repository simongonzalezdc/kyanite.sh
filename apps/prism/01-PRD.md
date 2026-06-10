# prism.sh - Product Requirements Document

**Version:** 1.0  
**Date:** November 2025  
**Owner:** Simon | Kyanite Suite  
**Status:** READY FOR IMPLEMENTATION  
**Target:** Independent GitHub Repository + v1.0 Release

---

## Executive Summary

**prism.sh** is a terminal-based color palette design and theory tool. It enables designers, artists, and creatives to explore color relationships, generate harmonious palettes, and understand color theory—entirely within the command line.

### Value Proposition

- **For Designers:** Quickly explore color relationships and generate harmonious palettes
- **For Creatives:** Match visual concepts to color choices instantly
- **For Developers:** Generate color palettes programmatically

### Target Users

1. **Primary:** Terminal-native designers and visual artists
2. **Secondary:** Developers building color systems
3. **Tertiary:** Creative professionals exploring terminal workflows

### Positioning

> "Color design in your terminal—beautiful, fast, and perfectly on-brand."

---

## Product Vision

### What It Is

An interactive color design tool with educational components. Create palettes, understand color theory, validate accessibility, and export in multiple formats—all without leaving the terminal.

### What It Is NOT

- Not a full design suite (use Figma/Illustrator for comprehensive design)
- Not a photo color extractor (use dedicated image analysis tools)
- Not a complete color accessibility suite (has basic WCAG checks only)

### Kyanite Suite Integration

**prism.sh is standalone but designed to integrate with future Kyanite tools:**
- Export palettes as themes for other tools
- Use as a shared library for color workflows
- Reference in documentation and tutorials

---

## Core Features

### Feature 1: Interactive Color Wheel / Circle of Colors

**Purpose:** Visual exploration of hue relationships

**Core Mechanics:**
- 360-degree hue circle displayed in terminal
- Navigate with arrow keys or vim keys (hjkl)
- Display current color: hex, RGB, HSL, HSV
- Show complementary, analogous, triadic colors automatically
- Animated highlight showing active hue
- Export current color to clipboard (platform-specific) or file

**Clipboard Support:**
- Linux: X11 via xclip/xsel, Wayland via wl-clipboard
- macOS: pbcopy
- Windows: clip.exe
- Fallback: Save to temporary file and display path

**Technical Requirements:**
- ASCII/Unicode visualization (no graphics)
- Real-time color calculations
- Terminal-aware color rendering:
  - Truecolor (24-bit): Full RGB support
  - 256-color: Closest color approximation
  - 16-color: Display hex codes only with warning
  - Monochrome: Fallback to text-only mode
- Responsive to terminal resize
- Auto-detect terminal capabilities via COLORTERM and TERM environment variables

**Success Criteria:**
- [ ] Displays full hue spectrum in ≤180 characters width
- [ ] Navigate 360 degrees with <50ms keystroke response
- [ ] Display all color formats instantly (<100ms)
- [ ] Functional UI in 80x24 terminal (minimum viable layout)

**80x24 Terminal Layout:**
```
╭────────────────────────────────────────────────────────────────────────────╮
│ PRISM - Color Wheel                                                  [?]  │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                            │
│   [████████████████ Color Wheel Ring ████████████████]                   │
│   [Red→Orange→Yellow→Green→Cyan→Blue→Magenta→Red]                        │
│                          ▲ Current: 180°                                  │
│                                                                            │
│   Current Color:  ████  #00D4FF  Cyan                                     │
│   RGB: (0, 212, 255)   HSL: (180°, 100%, 50%)                            │
│                                                                            │
│   Complementary:  ████  #FF2B00  Red-Orange                               │
│   Analogous:      ████  #00FF85  Spring Green                             │
│                   ████  #0085FF  Azure                                    │
│                                                                            │
├────────────────────────────────────────────────────────────────────────────┤
│ ←→: Navigate | Enter: Select | C: Copy | S: Save | Q: Quit | ?: Help     │
╰────────────────────────────────────────────────────────────────────────────╯
```

**Status:** Core feature, MUST ship v1.0

---

### Feature 2: Palette Generator

**Purpose:** Generate harmonious color palettes based on mathematical rules

**Core Mechanics:**
- **Input:** Base color (hex/name) + harmony rule
- **Harmony Rules:**
  - Monochromatic (tints & shades of single hue)
  - Complementary (opposite hue)
  - Analogous (adjacent hues ±30°)
  - Triadic (3 evenly spaced hues 120° apart)
  - Tetradic (4 hues, 2 pairs opposite 180° apart)
  - Split-Complementary (complementary + one adjacent)
  - Square (4 evenly spaced hues 90° apart)
- **Output:** 3-7 color palette with visual swatches
- **Export:** JSON, CSS variables, TOML, Kyanite theme format

**Technical Requirements:**
- HSL/HSV color space calculations
- Harmonic relationships precise to ±5° hue error
- Real-time generation (<500ms)
- Generated palettes maintain sufficient contrast ratios and balanced saturation/lightness distribution

**Success Criteria:**
- [ ] Generates palettes for all 7 rules
- [ ] Generated palettes maintain minimum 3:1 contrast between adjacent colors
- [ ] Visual display shows all palette colors with hex codes
- [ ] Export works for all 4 formats (JSON, CSS, TOML, Kyanite)

**Status:** Core feature, MUST ship v1.0

---

### Feature 3: Named Color Database

**Purpose:** Quick reference to named colors and variations

**Core Mechanics:**
- Search by name ("red", "ocean", "forest")
- Display exact color values and variations
- Show color family (warm/cool, saturated/desaturated)
- Generate tints (lighter), shades (darker), tones (desaturated)
- Support CSS color names + extended set
- Quick insertion into active palette

**Technical Requirements:**
- Pre-loaded database (~500 named colors)
  - Source: CSS Color Module Level 4 (147 colors) + X11 colors (public domain)
  - License: Public domain / W3C standard
- Fast fuzzy search (<100ms, matches with 1-2 character differences)
- Variations calculated in HSL space
- Works offline

**Success Criteria:**
- [ ] Search finds 5-10 matching colors for common queries
- [ ] Display name, hex, RGB, HSL for each color
- [ ] Show 5 tint/shade variations (±20%, ±40% lightness)
- [ ] Fuzzy search matches colors with 1-2 character typos

**Status:** Core feature, MUST ship v1.0

---

### Feature 4: Color Theory Education

**Purpose:** Teach color relationships and theory in context

**Core Mechanics:**
- **Interactive Lessons:** "What are complementary colors?", "Show analogous relationships", "Warm vs cool"
- **Quick Tips:** Display theory snippets inline while working
- **Visual Examples:** Show concepts with actual color swatches
- **Glossary:** Color theory terms with definitions
- **No Internet:** All content bundled (offline first)

**Content Required (v1.0):**
- 5+ foundational lessons
- 15+ glossary entries
- 50+ visual examples

**Success Criteria:**
- [ ] At least 5 lessons implemented
- [ ] Each lesson includes 3+ visual examples with color swatches
- [ ] Each lesson is 200-400 words (2-3 minute read time)
- [ ] Content uses simple language (8th grade reading level or below)

**Status:** Core feature, MUST ship v1.0

---

### Feature 5: WCAG Accessibility Validation

**Purpose:** Validate color combinations meet accessibility standards

**Core Mechanics:**
- **Input:** Two colors (foreground + background)
- **Output:** 
  - Contrast ratio (e.g., "4.5:1")
  - WCAG compliance level (AA/AAA/Fail)
  - Recommendation if failing
- **Batch Checking:** Validate multiple color pairs
- **Warnings:** Flag problematic combinations in palette

**Technical Requirements:**
- WCAG 2.1 contrast formula implementation
- Accurate to ±0.05 contrast ratio
- Works for all color formats

**Success Criteria:**
- [ ] Contrast ratio accurate within ±0.05 compared to WebAIM checker
- [ ] WCAG AA/AAA determination matches official WCAG 2.1 spec
- [ ] Provides specific recommendations (e.g., "Darken background by 15%")
- [ ] Supports batch checking of 10+ color pairs

**Status:** Core feature, MUST ship v1.0

---

### Feature 6: Palette Management

**Purpose:** Save, organize, and recall color palettes

**Core Mechanics:**
- **Save:** Store palette with name/description/tags
- **List:** Show saved palettes with color preview
- **Load:** Recall and edit previous palettes
- **Delete:** Remove saved palettes with confirmation
- **Export:** Save to file (JSON/CSS/TOML)
- **Import:** Load palette from file

**Storage:**
- Location (cross-platform):
  - Linux: `~/.config/prism/palettes/`
  - macOS: `~/Library/Application Support/prism/palettes/`
  - Windows: `%APPDATA%/prism/palettes/`
- Format: JSON with metadata
- Auto-organize by creation date
- File locking: Advisory locks prevent concurrent write conflicts

**Technical Requirements:**
- Fast file I/O (<100ms)
- Graceful error handling with user-friendly messages
- Support multiple export formats
- Handle corrupted files gracefully (skip and log, don't crash)
- Advisory file locking to prevent concurrent modification
- Atomic writes (write to .tmp, then rename) to prevent corruption

**Success Criteria:**
- [ ] Save/load cycle preserves all color values exactly (hex match)
- [ ] List shows 20+ previous palettes with pagination
- [ ] Export generates valid, parseable JSON/CSS/TOML files
- [ ] Import validates and handles malformed files gracefully

**Status:** Core feature, MUST ship v1.0

---

### Feature 7: Kyanite Suite Integration & Export

**Purpose:** Export colors for use in other Kyanite tools

**Core Mechanics:**
- **Export to Theme:** Convert palette to Kyanite theme format
- **Standard Format:** All exports use consistent structure
- **Interoperability:** Exported colors work in other tools
- **Documentation:** Format is documented and versioned

**Export Formats:**

**JSON:**
```json
{
  "id": "palette_20251115_103000",
  "name": "Electric Dream",
  "kyanite_version": "1.0",
  "created_at": "2025-11-15T10:30:00Z",
  "harmony_rule": "triadic",
  "colors": [
    {"hex": "#FF0080", "name": "Electric Pink", "role": "primary"},
    {"hex": "#00D4FF", "name": "Cyan", "role": "secondary"},
    {"hex": "#FFE600", "name": "Yellow", "role": "accent"}
  ]
}
```

**CSS Variables:**
```css
/* Generated by prism.sh - Electric Dream */
:root {
  --color-primary: #FF0080;
  --color-secondary: #00D4FF;
  --color-accent: #FFE600;
  --color-primary-rgb: 255, 0, 128;
  --color-secondary-rgb: 0, 212, 255;
  --color-accent-rgb: 255, 230, 0;
}
```

**TOML:**
```toml
# Generated by prism.sh - Electric Dream
name = "Electric Dream"
harmony_rule = "triadic"
created_at = "2025-11-15T10:30:00Z"

[[colors]]
name = "Electric Pink"
hex = "#FF0080"
role = "primary"

[[colors]]
name = "Cyan"
hex = "#00D4FF"
role = "secondary"

[[colors]]
name = "Yellow"
hex = "#FFE600"
role = "accent"
```

**Kyanite Theme Format:**
```json
{
  "name": "Electric Dream",
  "kyanite_version": "1.0",
  "created_at": "2025-11-15T10:30:00Z",
  "theme": {
    "primary": "#FF0080",
    "secondary": "#00D4FF",
    "accent": "#FFE600",
    "background": "#0D0221",
    "text": "#F0F3FF",
    "success": "#39FF14"
  }
}
```

**Success Criteria:**
- [ ] Exports are valid JSON per RFC 8259
- [ ] Format documented with version number and schema
- [ ] Can be imported by other tools (validation via JSON schema)

**Status:** Core feature, MUST ship v1.0

---

## Feature Matrix

| Feature | v1.0 | v1.1 | v2.0 |
|---------|------|------|------|
| **Color Wheel** | ✅ | | |
| **Palette Generator** | ✅ | | |
| **Named Colors** | ✅ | | |
| **Color Theory** | ✅ | | |
| **WCAG Checker** | ✅ | | |
| **Palette Manager** | ✅ | | |
| **Suite Export** | ✅ | | |
| **Color Mixing** | | ✅ | |
| **Mood-to-Palette** | | ✅ | |
| **Colorblind Sim** | | | ✅ |
| **Image Analysis** | | | ✅ |

---

## User Workflows

### Workflow 1: Designer Creating Brand Palette

```
1. Start prism
2. Navigate to hue (electric blue)
3. Generate tetradic palette
4. View 4-color palette
5. Check WCAG accessibility
6. Export as JSON
Time: 2-3 minutes
```

---

### Workflow 2: Learning Color Theory

```
1. Open lessons menu
2. Read "Complementary Colors"
3. See interactive example
4. Explore with color wheel
5. Understand concept
Time: 5-10 minutes per lesson
```

---

### Workflow 3: Checking Accessibility

```
1. Open WCAG checker
2. Input foreground: #FF0080
3. Input background: #0D0221
4. See: "5.2:1 - WCAG AAA ✓"
5. Adjust if needed
Time: 30 seconds
```

---

## Out of Scope

Explicitly NOT v1.0:

- Photo/image analysis (extract colors from images)
- Gradient generation (smooth color transitions)
- Colorblind simulation (deferred to v2.0)
- Real-time collaboration (multi-user editing)
- Cloud sync (all storage is local)
- Animation or transitions beyond simple highlights

---

## Success Metrics

### Technical Metrics

- **Startup:** <1 second
- **Color wheel render:** <100ms
- **Palette generation:** <500ms
- **Search:** <100ms
- **Memory:** <50MB idle
- **File I/O:** <100ms

### Feature Completion

- [ ] All 7 core features implemented
- [ ] All acceptance criteria met
- [ ] 0 known critical bugs

### Code Quality

- [ ] No TODO comments in main code paths
- [ ] Error handling complete with user-friendly messages
- [ ] Tests for critical paths (>70% coverage)
- [ ] Documentation comprehensive (README, ARCHITECTURE, inline docs)

---

## Acceptance Criteria

### Must Have (v1.0)

- [ ] Color wheel displays all 360 hues visually
- [ ] Arrow key navigation smooth (<50ms response)
- [ ] All 7 harmony rules work correctly (±5° accuracy)
- [ ] Generated palettes maintain 3:1 minimum contrast between adjacent colors
- [ ] Save/load preserves colors exactly (hex value match)
- [ ] WCAG checker accurate (±0.05 ratio compared to WebAIM)
- [ ] Export to JSON, CSS, TOML, Kyanite format works
- [ ] 10 Kyanite app themes applied for UI styling
- [ ] Ctrl+Shift+T theme switcher works
- [ ] All universal shortcuts implemented (except undo/redo - not in v1.0)
- [ ] Help system complete with all shortcuts documented
- [ ] No panics in application code (all panics recovered gracefully)
- [ ] Functional UI in 80x24 terminal minimum
- [ ] README complete with installation and usage examples
- [ ] ARCHITECTURE.md complete with module breakdown
- [ ] LICENSE file included (MIT)
- [ ] CONTRIBUTING.md with guidelines
- [ ] Clipboard support works on Linux/macOS/Windows

### Performance Targets

- [ ] Startup: <1s
- [ ] All operations: <500ms
- [ ] Memory: <50MB idle

### Quality Targets

- [ ] 0 panics in application code (all external panics recovered)
- [ ] 100% error handling with user-friendly messages
- [ ] All user workflows completable in <5 minutes

---

## Timeline

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| **Foundation** | 2-3 days | Project structure, theme system, core color math |
| **Features** | 5-6 days | All 7 features with UI, cross-platform support |
| **Content** | 2-3 days | Color theory lessons, named colors DB, help system |
| **Polish** | 2-3 days | Testing, docs, platform testing, refinement |
| **Release** | 1 day | GitHub repo, CI/CD, v1.0 release |
| **TOTAL** | 12-16 days | Launch ready |

**Note:** Timeline assumes single developer working full-time. Adjust accordingly for part-time or multiple contributors.

---

## Approval

**Product Owner:** Simon  
**Status:** ✅ APPROVED FOR IMPLEMENTATION  
**Start Date:** When ready  
**Target Launch:** v1.0 on GitHub

---

**This is a standalone, independent product. prism.sh can be built and released completely independently from syntax.sh or any other Kyanite tool.**

**Next step:** Review Technical Design Document (TDD)
