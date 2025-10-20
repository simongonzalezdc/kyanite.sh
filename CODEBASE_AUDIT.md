# focus.sh Codebase Audit Report

## 1. Project Info
**Current Names:**
- Binary: `neon.exe` (Windows), `todo.exe` (Linux/Mac)
- Module: `github.com/crush-cli/crush`
- Package: `crush-cli/crush`
- Directory: `crush-cli`

**Go Version:** 1.23.6

**LOC (Approximate):** 47 Go files, ~8,000+ lines

---

## 2. Directory Structure
```
C:\Users\Simon\dev\crush-cli\
├── cmd\
│   ├── neon\main.go          # Main NEON application
│   ├── todo\main.go          # TODO application  
│   ├── mcp-server\main.go    # MCP server
│   └── glow_demo\main.go     # Demo application
├── internal\
│   ├── tui\                  # Terminal UI components
│   ├── cli\                  # CLI commands
│   ├── ai\                   # AI integration
│   ├── engine\               # Core business logic
│   ├── store\                # Data persistence
│   └── wizards\              # Interactive wizards
├── pkg\
│   ├── models\task.go        # Task data models
│   ├── calendar\             # Calendar functionality
│   ├── styles\               # Visual themes
│   ├── audio\                # Sound effects
│   ├── config\               # Configuration
│   └── utils\                # Utility functions
├── assets\                   # Static assets
├── test\                     # Test files
└── docs\                     # Documentation
```

---

## 3. Features Status

| Feature | Status | Files | DB Tables | Keybindings |
|---------|--------|-------|-----------|-------------|
| Tasks | ✅ Done | `internal/cli/add.go`, `list.go`, `complete.go`, `delete.go` | JSON | `neon add`, `neon list` |
| Notes | ✅ Done | `internal/cli/notes.go` | Embedded in Task | `neon notes <id>` |
| Calendar | ✅ Done | `pkg/calendar/calendar.go`, `internal/cli/calendar.go` | JSON | `neon calendar` |
| Themes | ✅ Done | `pkg/styles/themes.go`, `synthwave.go` | JSON | `neon theme` |
| AI Integration | ✅ Done | `internal/ai/manager.go` | Cache JSON | `neon chat`, `inspire` |
| Pomodoro | ✅ Done | `internal/tui/main.go` | JSON | `s`, `x`, `r`, `f` |
| TUI Dashboard | ✅ Done | `internal/tui/main.go` | JSON | Tab navigation |

---

## 4. Database

**Storage Method:** JSON files (no SQL database)

**File Locations:**
- Main storage: `~/.neon/tasks.json`
- Configuration: `~/.neon/config.yml`
- AI cache: `~/.todo/ai_cache.json`
- Test data: `test/test_tasks.json`

**Data Schema (Task):**
```json
{
  "id": "string",
  "description": "string", 
  "status": "pending|completed",
  "priority": "low|medium|high",
  "deadline": "time.Time",
  "categories": ["string"],
  "notes": "string",
  "created_at": "time.Time",
  "updated_at": "time.Time"
}
```

**No migrations** - Simple JSON append/overwrite pattern

---

## 5. UI Architecture

**Screens:**
- `dashboardView` - Main task list (`internal/tui/main.go`)
- `focusView` - Pomodoro timer (`internal/tui/main.go`)
- `chatView` - AI assistant (`internal/tui/main.go`)
- `calendarView` - Calendar view (`internal/tui/main.go`)
- `notesView` - Notes editing (`internal/tui/main.go`)
- `settingsView` - Settings configuration (`internal/tui/main.go`)

**Navigation:**
- Tab-based switching between views
- Arrow keys for task navigation
- Key bindings: `s` (start), `x` (stop), `r` (reset), `f` (focus), `c` (chat)

**Theming:**
- Lipgloss-based styling system
- Theme files in `pkg/styles/`
- Cyberpunk/synthwave aesthetic default
- JSON theme configuration in `assets/themes/`

---

## 6. AI Integration

**Provider:** Ollama (local) + OpenRouter (fallback)

**Models:**
- Primary: `qwen2.5:1.5b` (Ollama)
- Fallback: `gpt-3.5-turbo` (OpenRouter)

**AI Features:**
- Natural language task parsing (`neon add "..."`)
- Contextual task suggestions (`neon inspire`)
- AI-powered task summaries
- Cyberpunk-themed chat assistant (`neon chat`)
- Real-time AI status monitoring

**Prompts:** Defined in `internal/ai/manager.go` with cyberpunk persona

**Caching:** 24-hour TTL JSON cache at `~/.todo/ai_cache.json`

---

## 7. Dependencies

**Core Dependencies:**
```
github.com/charmbracelet/bubbles v0.21.1
github.com/charmbracelet/bubbletea v1.3.6  
github.com/charmbracelet/huh v0.8.0
github.com/charmbracelet/lipgloss v1.1.1
github.com/spf13/cobra v1.9.1
github.com/spf13/viper v1.21.0
github.com/pterm/pterm v0.12.76
```

**Total Dependencies:** 45+ (including indirect)

---

## 8. Build

**Build Commands:**
- Windows: `go build -o neon.exe ./cmd/neon`
- Linux/Mac: `go build -o todo.exe cmd/todo/main.go`

**Scripts:**
- `build.bat` - Windows build
- `build.sh` - Linux/Mac build  
- `setup.bat` - Environment setup
- `StartNEON.bat` - Application launcher

**Cross-Platform:** ✅ Windows, Linux, Mac supported

---

## 9. Naming Issues

**Multiple Project Identities:**

**"neon" occurrences:**
- `cmd/neon/main.go:93` - Categories: ["coding", "neon"]
- `internal/tui/unified_dashboard.go:345` - "Use: neon interactive"
- `internal/cli/theme.go:17` - "neon colors (default)"
- `internal/ai/manager.go:235` - "Try using 'neon inspire'"
- `internal/cli/root.go:12` - `neonBlue` color constant
- `pkg/config/config.go:69` - Config path: `~/.neon/`
- `internal/cli/add.go:23` - "Example: neon add"
- `pkg/utils/storage.go:12` - Storage path: `~/.neon/tasks.json`
- `pkg/styles/synthwave.go:278` - "Initialize with 'neon add'"
- `pkg/styles/styles.go:16` - `neonRed` color constant
- 10+ other command references and help text

**"focus" occurrences:**
- `internal/tui/main.go:152` - `focusView` function
- `internal/cli/complete.go:50` - "boost your focus karma"
- `internal/cli/calendar.go.bak:318` - "stay focused"

**"lyricforge" occurrences:**
- None found

**Inconsistencies:**
- Module: `github.com/crush-cli/crush` vs UI brand: "NEON"
- Storage: `~/.neon/` but AI cache: `~/.todo/`
- Build outputs: `neon.exe` vs `todo.exe`
- CLI commands: `neon` vs some references to `todo`

---

## 10. Red Flags

**Critical Issues:**
- ❌ **Naming chaos** - 3 different project identities (crush, neon, todo)
- ❌ **Inconsistent paths** - Mixed `~/.neon/` and `~/.todo/` directories
- ❌ **Build confusion** - Different binary names across platforms

**Performance Issues:**
- ⚠️ JSON file locking for concurrent access
- ⚠️ No database indexing for large task lists
- ⚠️ AI response caching could be more efficient

**Areas Needing Refactor Before Rename:**
1. **Unified branding** - Choose single project name
2. **Consistent paths** - Standardize storage locations  
3. **Build harmonization** - Same binary name across platforms
4. **Module path update** - Align with chosen brand
5. **Documentation sync** - Update all help text and examples
6. **Configuration cleanup** - Remove conflicting config paths

**Recommendation:** Resolve naming inconsistencies before any release to avoid user confusion and maintenance overhead.