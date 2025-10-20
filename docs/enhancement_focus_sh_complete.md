# Enhancement: focus.sh Realignment to Kyanite Suite

**Status:** READY FOR IMPLEMENTATION  
**Priority:** HIGH  
**Complexity:** MODERATE (naming chaos requires careful refactor)  
**Estimated Effort:** 2-3 weeks  
**Target Ship:** Month 4-5  
**Owner:** Simon | Kyanite Suite

---

## Executive Summary

The existing "Neon Focus" codebase (`crush-cli/crush`) has severe **naming/branding inconsistencies** (neon/crush/todo/focus mixed throughout). This enhancement realigns it to the **Kyanite Suite** as **focus.sh** with:

✅ Single unified brand: **focus.sh**  
✅ Consistent module/binary/paths naming  
✅ Journal feature with syntax.sh export pipeline  
✅ Kyanite Suite integration strategy  
✅ Production-ready release  

---

## Part 1: Current State Assessment

### Naming Chaos Breakdown

| Component | Current | Target |
|-----------|---------|--------|
| **Module Path** | `github.com/crush-cli/crush` | `github.com/kyanite/focus` |
| **Binary Name** | `neon.exe` / `todo.exe` | `focus` (all platforms) |
| **Storage Dir** | `~/.neon/` + `~/.todo/` | `~/.focus/` |
| **Config File** | `~/.neon/config.yml` | `~/.focus/config.yml` |
| **AI Cache** | `~/.todo/ai_cache.json` | `~/.focus/ai_cache.json` |
| **CLI Brand** | "neon", "crush", "todo" | "focus" |
| **Package Name** | `crush-cli/crush` | `focus` |

### Current Features (All ✅ Complete)
- Tasks (add/list/complete/delete)
- Notes (embedded in tasks)
- Calendar (month/week view)
- Themes (synthwave/cyberpunk)
- AI (Ollama + OpenRouter fallback)
- Pomodoro (25-min sessions)
- TUI Dashboard (tab-based navigation)

---

## Part 2: Branding Realignment

### 2.1 Module & Package Renaming

**Go Module Path:**
```bash
# OLD: github.com/crush-cli/crush
# NEW: github.com/kyanite/focus
```

**In go.mod:**
```
OLD: module github.com/crush-cli/crush
NEW: module github.com/kyanite/focus
```

**All imports update:**
```go
// OLD
import "github.com/crush-cli/crush/internal/tui"

// NEW
import "github.com/kyanite/focus/internal/tui"
```

**Affected files:** All 47 Go files with imports

---

### 2.2 Directory Restructuring

**Rename directories:**
```
crush-cli/                    → focus/
cmd/neon/main.go             → cmd/focus/main.go
cmd/todo/main.go             → DELETE (duplicate, merge into focus)
cmd/mcp-server/main.go       → cmd/focus-mcp/main.go (rename for clarity)
cmd/glow_demo/main.go        → DELETE (internal demo only)
```

**Result:**
```
focus/
├── cmd/
│   ├── focus/main.go        # Main binary
│   └── focus-mcp/main.go    # MCP server (optional)
├── internal/                # (unchanged structure)
├── pkg/                      # (unchanged structure)
├── go.mod                   # Updated module path
└── go.sum
```

---

### 2.3 Storage Path Consolidation

**Migrate all storage to ~/.focus/**

**Current paths:**
- `~/.neon/tasks.json` → `~/.focus/tasks.json`
- `~/.neon/config.yml` → `~/.focus/config.yml`
- `~/.todo/ai_cache.json` → `~/.focus/ai_cache.json`

**Migration strategy:**
```go
func migrateStorage() error {
    // On first run of new version:
    // 1. Check for ~/.neon/ → copy to ~/.focus/
    // 2. Check for ~/.todo/ → merge into ~/.focus/
    // 3. Delete old directories (after backup)
}
```

**Code location:** `pkg/utils/storage.go` (new migration function)

---

### 2.4 String & Reference Updates

**Search & replace in ALL files:**

| Old Pattern | New Pattern | Files |
|-------------|-------------|-------|
| "neon" (brand ref) | "focus" | ~15 files |
| "crush-cli" | "kyanite" | ~8 files |
| "todo" (brand ref) | "focus" | ~10 files |
| "NEON" (display) | "focus.sh" | ~5 files |
| `neonBlue` (color var) | `focusBlue` | styles files |
| `~/.neon` | `~/.focus` | help text, docs |
| `~/.todo` | `~/.focus` | help text, docs |

**Critical files to audit:**
- `internal/cli/root.go` - Help text, examples
- `internal/ai/manager.go` - Prompts, examples
- `internal/tui/main.go` - Display text
- `pkg/styles/themes.go` - Theme names
- All help text in CLI commands
- README.md, docs/

---

### 2.5 Build Configuration

**Update build scripts:**
```bash
# build.bat (Windows)
OLD: go build -o neon.exe ./cmd/neon
NEW: go build -o focus.exe ./cmd/focus

# build.sh (Linux/Mac)
OLD: go build -o todo.exe cmd/todo/main.go
NEW: go build -o focus cmd/focus/main.go
```

**Result:** Same binary name on all platforms: `focus` (or `focus.exe` on Windows)

---

## Part 3: Journal Feature

### 3.1 Database Schema Addition

**New table in storage (JSON format):**
```json
{
  "id": "uuid",
  "date": "2025-10-19",
  "title": "optional string",
  "content": "markdown text",
  "mood": "emoji or string (😌, energized, etc)",
  "tags": ["tag1", "tag2"],
  "word_count": 150,
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "is_private": false,
  "template_used": "morning_pages|evening_reflection|daily_log"
}
```

**Storage:**
- File: `~/.focus/journal.json`
- Structure: Array of entries, searchable by date/tag

---

### 3.2 Journal CLI Commands

```bash
# Create entry
focus journal new [--date 2025-10-19] [--template morning_pages]

# List entries
focus journal list [--date 2025-10-19] [--tag productivity]

# View entry
focus journal view <date>

# Search
focus journal search "keyword"

# Export to syntax.sh
focus journal export <date> --to-syntax --type character|dialogue|scene|research
```

---

### 3.3 Journal TUI Screen

**New view: `journalView`**

```
┌─ focus.sh | JOURNAL ─────────────────────────────────────────┐
│ [2025-10-19] [Mood: 😌] [Tags: #productive #breakthrough]   │
│                                                               │
│ Had solid progress on focus.sh today. Architecture is        │
│ clicking into place. Realized the importance of separation   │
│ of concerns...                                                │
│                                                               │
│ [Word count: 147] [Ctrl+E] Edit | [Ctrl+X] Export           │
│                                                               │
│ Navigation: [←→] Prev/Next | [N]ew | [S]earch | [E]xport    │
└───────────────────────────────────────────────────────────────┘
```

**Navigation shortcuts:**
- `J` - Open journal
- `←/→` - Previous/next entry
- `N` - New entry
- `S` - Search
- `E` - Export to syntax.sh
- `M` - Set mood

---

### 3.4 Export to syntax.sh Pipeline

**New feature: Export journal entries as story material**

**Command:**
```bash
focus journal export <date> --to-syntax --type <type>
```

**Export types:**
- `character` → `~/syntax/imports/character-<date>.md`
- `dialogue` → `~/syntax/imports/dialogue-<date>.md`
- `scene` → `~/syntax/imports/scene-<date>.md`
- `research` → `~/syntax/imports/research-<date>.md`

**Generated file format:**
```markdown
# Imported from focus.sh
**Source:** Journal entry 2025-10-19  
**Type:** Character observation  
**Tags:** #breakthrough #architecture

---

## Raw Entry
[Full journal content]

---

## Story Application
[Space for user to develop this into story material]
```

**Implementation:**
```go
// internal/journal/export.go
func (e *Entry) ExportToSyntax(exportType string) error {
    filename := fmt.Sprintf("%s-%s.md", exportType, e.Date)
    content := formatForSyntaxImport(e, exportType)
    return writeToSyntaxImports(filename, content)
}
```

---

## Part 4: Kyanite Suite Integration

### 4.1 Suite Positioning

**focus.sh role in Kyanite:**
```
noise.sh (Music ideas)
  ↓ (task tracking)
focus.sh (Productivity & Reflection)
  ↓ (journal export)
syntax.sh (Story development)
  ↓ (worldbuilding)
wave.sh (Audio production)
  ↓ (final output)
prism.sh (Color/visual)
```

---

### 4.2 Cross-Tool Awareness

**In help text, reference other tools:**

```bash
focus help

Output should include:
"Part of Kyanite Creative Suite:
  noise.sh  - Music & songwriting
  syntax.sh - Fiction & story writing
  prism.sh  - Color & visual design
  wave.sh   - Audio production & synthesis
  focus.sh  - Your tool (productivity & journaling)
"
```

---

### 4.3 Future Integration Points (placeholder)

```go
// pkg/integration/kyanite.go (new file)

// Placeholder for future integrations:
// - Import tasks from noise.sh projects
// - Export journal to syntax.sh ✅ (Phase 1)
// - Link to prism.sh color palettes (Phase 2)
// - Sync with wave.sh session logs (Phase 3)
```

---

## Part 5: Implementation Checklist

### Phase 1: Renaming & Refactor (Week 1)

**Go Module & Imports:**
- [ ] Update go.mod: `github.com/crush-cli/crush` → `github.com/kyanite/focus`
- [ ] Run `find . -name "*.go" -exec sed -i 's|crush-cli/crush|kyanite/focus|g' {} \;`
- [ ] Verify all imports: `go build ./...` compiles

**Directory Structure:**
- [ ] Rename `crush-cli/` → `focus/`
- [ ] Rename `cmd/neon/` → `cmd/focus/`
- [ ] Merge `cmd/todo/main.go` into `cmd/focus/main.go` (delete duplicate)
- [ ] Rename `cmd/mcp-server/` → `cmd/focus-mcp/`
- [ ] Delete `cmd/glow_demo/`

**Build Configuration:**
- [ ] Update `build.bat`: `neon.exe` → `focus.exe`
- [ ] Update `build.sh`: `todo.exe` → `focus`
- [ ] Test builds on all platforms

**String Updates:**
- [ ] Replace "neon" brand refs with "focus" (~15 files)
- [ ] Replace "~/.neon" with "~/.focus" (help text, code, docs)
- [ ] Replace "~/.todo" with "~/.focus" (all references)
- [ ] Verify no remaining "crush-cli" references in user-facing text

**Tests & Verification:**
- [ ] `go build ./...` passes
- [ ] `go test ./...` passes
- [ ] Binary runs: `./focus --version`
- [ ] Help text shows "focus.sh" not "neon"

---

### Phase 2: Storage Migration (Week 1-2)

**Migration Logic:**
- [ ] Implement `migrateStorage()` in `pkg/utils/storage.go`
- [ ] On first run of v2.0, detect old `~/.neon/` and `~/.todo/` paths
- [ ] Copy/merge to `~/.focus/`
- [ ] Create backup: `~/.focus/backup-v1/` (preserve old data)
- [ ] Delete old directories after successful migration

**Path Updates:**
- [ ] Update all hardcoded paths: `~/.neon/` → `~/.focus/`
- [ ] Update `pkg/config/config.go`
- [ ] Update `pkg/utils/storage.go`
- [ ] Update help text and documentation

**Testing:**
- [ ] [ ] Run migration on test system with v1 data
- [ ] [ ] Verify all tasks/notes/settings migrate correctly
- [ ] [ ] Confirm backup created and accessible
- [ ] [ ] Test on fresh install (no migration needed)

---

### Phase 3: Journal Feature (Week 2-3)

**Data Model:**
- [ ] Create `models/journal.go` with Entry struct
- [ ] Add schema to `~/.focus/journal.json`

**CLI Commands:**
- [ ] Implement `cmd/focus/journal.go`
- [ ] Add `journal new` command
- [ ] Add `journal list` command
- [ ] Add `journal view` command
- [ ] Add `journal search` command
- [ ] Add `journal export --to-syntax` command

**TUI Screen:**
- [ ] Create `internal/tui/journal_view.go`
- [ ] Add entry creation screen
- [ ] Add entry viewing screen
- [ ] Add mood/tag picker
- [ ] Add navigation (`←/→` prev/next)
- [ ] Integrate into tab navigation

**Export Pipeline:**
- [ ] Create `internal/journal/export.go`
- [ ] Implement `ExportToSyntax()` function
- [ ] Create output formats: character, dialogue, scene, research
- [ ] Test export creates valid markdown in syntax.sh import dir

**Testing:**
- [ ] [ ] Create journal entry via CLI
- [ ] [ ] View entry in TUI
- [ ] [ ] Search entries
- [ ] [ ] Export to syntax.sh
- [ ] [ ] Verify syntax.sh receives import correctly

---

### Phase 4: Documentation & Polish (Week 3)

**Documentation:**
- [ ] Update README.md: All references to "neon" → "focus.sh"
- [ ] Add journal feature documentation
- [ ] Update quick start guide
- [ ] Document storage paths
- [ ] Add migration guide for v1 users
- [ ] Create example journal entries

**Testing:**
- [ ] [ ] Run full test suite: `go test ./...`
- [ ] [ ] Manual testing on all 3 platforms
- [ ] [ ] Test all keyboard shortcuts
- [ ] [ ] Verify AI integration still works
- [ ] [ ] Test migration from v1 → v2
- [ ] [ ] Test journal creation/export

**Final Checks:**
- [ ] [ ] No "neon", "crush", "todo" in user-facing text
- [ ] [ ] All help text updated
- [ ] [ ] Binary named "focus" on all platforms
- [ ] [ ] Storage path is ~/.focus/ everywhere
- [ ] [ ] AI provider still works (Ollama + OpenRouter)
- [ ] [ ] No import errors or dead code

---

## Part 6: File-by-File Changes Required

**Critical files to update (rename or modify):**

| File | Change | Priority |
|------|--------|----------|
| `cmd/neon/main.go` | Rename to `cmd/focus/main.go` | CRITICAL |
| `cmd/todo/main.go` | Merge into `cmd/focus/main.go` | CRITICAL |
| `go.mod` | Update module path | CRITICAL |
| `pkg/utils/storage.go` | Add migration logic + path updates | CRITICAL |
| `pkg/config/config.go` | Update all paths to ~/.focus/ | CRITICAL |
| `internal/ai/manager.go` | Update example text | HIGH |
| `internal/cli/root.go` | Update help text + examples | HIGH |
| `internal/cli/*.go` | Update command help text | HIGH |
| `pkg/styles/synthwave.go` | Update theme names/refs | MEDIUM |
| `pkg/styles/themes.go` | Replace with theme registry | CRITICAL |
| `assets/themes/` | Remove synthwave configs, add Kyanite themes | CRITICAL |
| `internal/theme/registry.go` | Create with all 12 themes | CRITICAL |
| `README.md` | Complete rewrite for focus.sh | HIGH |
| `build.bat`, `build.sh` | Update output binary names | CRITICAL |
| All help text | Replace "neon" → "focus" | HIGH |

---

## Part 7: Success Criteria

**focus.sh v2.0 ships when:**

✅ Zero "neon", "crush", "todo" references in user-facing text  
✅ Binary consistently named "focus" across all platforms  
✅ All storage at `~/.focus/`  
✅ Module path: `github.com/kyanite/focus`  
✅ Journal feature working (create/view/search/export)  
✅ Export to syntax.sh pipeline tested  
✅ **12 Kyanite themes implemented (same as noise.sh)**  
✅ **Default theme: Amber Night**  
✅ Theme switching with Ctrl+T  
✅ Migration from v1 data succeeds  
✅ All tests pass: `go test ./...`  
✅ README updated with focus.sh branding  
✅ Demo tested on Windows, Mac, Linux  
✅ AI integration still functional  

---

## Part 8: Deployment

**Release Package:**
```bash
# Outputs
focus-v2.0-windows-amd64.exe
focus-v2.0-darwin-amd64 (Mac)
focus-v2.0-linux-amd64

# With documentation
README.md
MIGRATION_GUIDE.md (for v1 → v2 users)
JOURNAL_GUIDE.md (new feature)
```

**Installation:**
```bash
# User downloads and runs
./focus --version  # Should show: focus.sh v2.0

# First run triggers migration
# Result: ~/.focus/ with all data migrated
```

---

## Part 9: Theme Color Definitions (All 13)

**For `internal/theme/registry.go`:**

```go
var SlateMist = Theme{
    Name:       "Slate Mist",
    Primary:    "#E0E0E0",
    Secondary:  "#B0B0B0",
    Accent:     "#FFFFFF",
    Background: "#1A1A1A",
    Text:       "#E8E8E8",
    Success:    "#90C695",
}

var VioletDusk = Theme{
    Name:       "Violet Dusk",
    Primary:    "#9D84B7",
    Secondary:  "#6D5B7F",
    Accent:     "#D4A5FF",
    Background: "#1A1527",
    Text:       "#E8DFF5",
    Success:    "#9FE2BF",
}

var AmberNight = Theme{
    Name:       "Amber Night",
    Primary:    "#B8936E",
    Secondary:  "#6D5B8B",
    Accent:     "#E8C547",
    Background: "#12101A",
    Text:       "#F0E6D8",
    Success:    "#52D3AA",
}

var MoltenGold = Theme{
    Name:       "Molten Gold",
    Primary:    "#D4AF37",
    Secondary:  "#8B5A2B",
    Accent:     "#FFD700",
    Background: "#0F0905",
    Text:       "#FFF8DC",
    Success:    "#98D8C8",
}

var ClayRoads = Theme{
    Name:       "Clay Roads",
    Primary:    "#8B4513",
    Secondary:  "#A0522D",
    Accent:     "#CD853F",
    Background: "#1A1410",
    Text:       "#F5DEB3",
    Success:    "#7FB89C",
}

var IronStorm = Theme{
    Name:       "Iron Storm",
    Primary:    "#DC143C",
    Secondary:  "#696969",
    Accent:     "#FF6347",
    Background: "#0B0C0E",
    Text:       "#F5F5F5",
    Success:    "#90EE90",
}

var JadeTide = Theme{
    Name:       "Jade Tide",
    Primary:    "#20B2AA",
    Secondary:  "#2E8B87",
    Accent:     "#E0F2F1",
    Background: "#0F1419",
    Text:       "#F5F8FA",
    Success:    "#7FB89C",
}

var SunsetEmber = Theme{
    Name:       "Sunset Ember",
    Primary:    "#FF6B9D",
    Secondary:  "#FF8C42",
    Accent:     "#FFC312",
    Background: "#1e1e2e",
    Text:       "#FFEAA7",
    Success:    "#55E6C1",
}

var ForestWhisper = Theme{
    Name:       "Forest Whisper",
    Primary:    "#52B788",
    Secondary:  "#52A068",
    Accent:     "#95D5B2",
    Background: "#1B263B",
    Text:       "#D8F3DC",
    Success:    "#B7E4C7",
}

var ElectricBloom = Theme{
    Name:       "Electric Bloom",
    Primary:    "#FF0080",
    Secondary:  "#00D4FF",
    Accent:     "#FFE600",
    Background: "#0D0221",
    Text:       "#F0F3FF",
    Success:    "#39FF14",
}

var PlasmaPulse = Theme{
    Name:       "Plasma Pulse",
    Primary:    "#39FF14",
    Secondary:  "#00F5FF",
    Accent:     "#FF1493",
    Background: "#0A0118",
    Text:       "#E0FFFF",
    Success:    "#00FF7F",
}

var IndigoDepths = Theme{
    Name:       "Indigo Depths",
    Primary:    "#4B0082",
    Secondary:  "#36648B",
    Accent:     "#9370DB",
    Background: "#0A0E27",
    Text:       "#E6E6FA",
    Success:    "#7FB89C",
}

var SageMeadow = Theme{
    Name:       "Sage Meadow",
    Primary:    "#9CAF88",
    Secondary:  "#B8B8A3",
    Accent:     "#D4D4AF",
    Background: "#1B1B1B",
    Text:       "#E8E8D0",
    Success:    "#A8D5BA",
}
```

---

## Part 10: Migration Code Template

**For `pkg/utils/storage.go`:**

```go
func migrateStorage() error {
    // Step 1: Check for old directories
    neonPath := filepath.Join(os.Getenv("HOME"), ".neon")
    todoPath := filepath.Join(os.Getenv("HOME"), ".todo")
    focusPath := filepath.Join(os.Getenv("HOME"), ".focus")
    
    // Create new focus directory if doesn't exist
    os.MkdirAll(focusPath, 0755)
    
    // Step 2: Backup existing .focus if it exists
    if _, err := os.Stat(focusPath); err == nil {
        backupPath := filepath.Join(focusPath, "backup-pre-migration")
        os.Rename(focusPath, backupPath)
        os.MkdirAll(focusPath, 0755)
    }
    
    // Step 3: Migrate from .neon
    if _, err := os.Stat(neonPath); err == nil {
        copyDir(neonPath, focusPath)
    }
    
    // Step 4: Merge .todo (especially ai_cache.json)
    if _, err := os.Stat(todoPath); err == nil {
        if _, err := os.Stat(filepath.Join(todoPath, "ai_cache.json")); err == nil {
            src := filepath.Join(todoPath, "ai_cache.json")
            dst := filepath.Join(focusPath, "ai_cache.json")
            copyFile(src, dst)
        }
    }
    
    return nil
}

func copyDir(src, dst string) error {
    return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        rel, _ := filepath.Rel(src, path)
        dstPath := filepath.Join(dst, rel)
        
        if info.IsDir() {
            os.MkdirAll(dstPath, 0755)
        } else {
            copyFile(path, dstPath)
        }
        return nil
    })
}
```

---

## Part 11: Journal Export Format Example

**Output file: `~/syntax/imports/character-2025-10-19.md`**

```markdown
# Imported from focus.sh

**Source:** Journal entry from 2025-10-19  
**Import Type:** Character observation  
**Tags:** #breakthrough #architecture  
**Mood:** Focused  

---

## Original Journal Entry

Had solid progress on focus.sh today. Architecture is clicking into place. 
Realized the importance of separation of concerns. Each tool does one thing 
well and integrates cleanly with others. This is how good systems work.

---

## Story Application

### Character Notes
How does this character approach systems? Do they:
- Break complex problems into simpler parts?
- Integrate different elements cohesively?
- Appreciate clarity and separation of concerns?
- Have "breakthrough" moments of clarity?

### Dialogue Seed
_"Each part has its purpose. Nothing more, nothing less."_

### Scene Potential
A moment where character realizes the elegance of a well-designed system.

---

**Next Steps:** Develop this into character notes, scene, or dialogue as needed.
```

---

## Part 12: Configuration File Examples

**`~/.focus/config.yml` (after migration):**

```yaml
# focus.sh Configuration

preferences:
  theme: "amber-night"
  default_view: "dashboard"
  
pomodoro:
  work_duration: 25
  break_duration: 5
  long_break_duration: 15
  sessions_before_long_break: 4

storage:
  tasks_file: "tasks.json"
  journal_file: "journal.json"
  config_file: "config.yml"
  ai_cache_file: "ai_cache.json"

ai:
  provider: "ollama"
  model: "qwen2.5:1.5b"
  temperature: 0.7
  timeout_seconds: 2
```

---

## Part 13: Find & Replace Patterns

**For coding agent to use:**

```bash
# Pattern 1: Module imports
OLD: "github.com/crush-cli/crush"
NEW: "github.com/kyanite/focus"
COMMAND: find . -name "*.go" -exec sed -i 's|github.com/crush-cli/crush|github.com/kyanite/focus|g' {} \;

# Pattern 2: Binary references in code
OLD: neon.exe / todo.exe
NEW: focus / focus.exe
COMMAND: grep -r "neon\.exe\|todo\.exe" --include="*.go" --include="*.sh" --include="*.bat"

# Pattern 3: Path references
OLD: ~/.neon/
NEW: ~/.focus/
COMMAND: grep -r "~/.neon" --include="*.go" --include="*.md"

# Pattern 4: String constants
OLD: "neon" / "crush" / "todo" (brand)
NEW: "focus"
MANUAL: Review each match - not all "neon" should be replaced (e.g., neonBlue var)

# Pattern 5: Help text
OLD: Use: neon add
NEW: Use: focus add
COMMAND: grep -r "Use: neon\|neon add\|neon list" --include="*.go"
```

---

## Part 14: Testing Checklist for Kilocode

**Before shipping, run these tests:**

### Migration Tests
- [ ] Backup created at ~/.focus/backup-pre-migration
- [ ] All tasks from ~/.neon/tasks.json copied to ~/.focus/
- [ ] AI cache from ~/.todo/ai_cache.json merged
- [ ] Old directories still exist (not deleted) for safety
- [ ] config.yml created with defaults
- [ ] Fresh install (no old data) still works

### Build Tests
- [ ] `go build ./...` compiles without errors
- [ ] `go test ./...` all tests pass
- [ ] Windows: `focus.exe --version` shows correct output
- [ ] Mac/Linux: `./focus --version` shows correct output
- [ ] No "neon", "crush", "todo" in version output

### Feature Tests
- [ ] Tasks: Add, list, complete, delete all work
- [ ] Notes: Create, view, edit functional
- [ ] Calendar: Navigate dates, view events
- [ ] Themes: Ctrl+T switches through all 13 themes
- [ ] AI: `focus chat` connects to Ollama, <2 sec response
- [ ] Pomodoro: Timer counts down, notifications work
- [ ] Journal: Create entry, add mood, add tags, search

### Journal-Specific Tests
- [ ] `focus journal new` creates entry
- [ ] Entry saved to ~/.focus/journal.json
- [ ] `focus journal list` shows all entries
- [ ] `focus journal search "keyword"` finds entry
- [ ] `focus journal export <date> --to-syntax --type character` creates file
- [ ] Exported file appears in ~/syntax/imports/

### Theme Tests
- [ ] All 13 themes load without error
- [ ] Each theme readable (good contrast)
- [ ] Theme change persists across restarts
- [ ] Amber Night is default
- [ ] Invalid theme name falls back to default

### Naming Tests
- [ ] Zero instances of "neon" in help text
- [ ] Zero instances of "crush-cli" in user-facing text
- [ ] All command examples show "focus" not "neon"
- [ ] Help text mentions Kyanite Suite
- [ ] Module imports show "kyanite/focus"

---

## Part 15: Deployment Checklist

**Before release:**

```
Version Bumping:
- [ ] Update go version (if needed)
- [ ] Version tag: v2.0.0 (major release due to rebranding)
- [ ] git tag -a v2.0.0 -m "Major: Rebrand to focus.sh, add journal, migrate themes"
- [ ] git push origin v2.0.0

Build Artifacts:
- [ ] focus-v2.0.0-windows-amd64.exe (Windows)
- [ ] focus-v2.0.0-darwin-amd64 (macOS Intel)
- [ ] focus-v2.0.0-darwin-arm64 (macOS Apple Silicon)
- [ ] focus-v2.0.0-linux-amd64 (Linux x86)
- [ ] focus-v2.0.0-linux-arm64 (Linux ARM - Raspberry Pi)

Documentation:
- [ ] README.md (complete rewrite for focus.sh)
- [ ] MIGRATION_GUIDE.md (v1 → v2 with migration steps)
- [ ] JOURNAL_GUIDE.md (how to use journal feature)
- [ ] THEMES_GUIDE.md (all 13 themes + custom theme setup)
- [ ] CHANGELOG.md (what changed, breaking changes, new features)

Quality Assurance:
- [ ] All tests passing
- [ ] Cross-platform tested (Win/Mac/Linux)
- [ ] Migration tested with real v1 data
- [ ] AI integration verified with Ollama running
- [ ] No compiler warnings
- [ ] No unused imports

GitHub Release:
- [ ] Tag created and pushed
- [ ] Release notes published
- [ ] Binaries attached
- [ ] Badges updated (if applicable)
- [ ] README links verified
```

---

## Notes

- **Naming complexity:** This refactor is non-trivial due to identity chaos (3 project names). Execute carefully.
- **Migration is critical:** Users have data in ~/.neon/ — losing it = disaster. Test thoroughly.
- **AI unchanged:** Ollama + OpenRouter integration should continue working.
- **Theme colors:** All 13 hex codes provided above - copy/paste directly into code.
- **Next phase:** Kyanite suite integration (linking to other tools) comes after shipping.
- **Version strategy:** This is v2.0.0 (major rebranding) not v1.x (minor). Update version in code.

---

**Status:** Ready for implementation  
**Owner:** Simon  
**Dependency:** Coding agent executes checklist above  
**Target Ship:** 3-4 weeks  
**Kilocode Notes:** All code snippets, color codes, patterns, and test cases provided. No guessing needed.