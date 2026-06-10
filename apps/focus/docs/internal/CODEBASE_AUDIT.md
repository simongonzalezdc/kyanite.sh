# focus.sh Codebase Audit Report

## 1. Project Info
- Current names
  - Binary: focus (focus.exe on Windows)
  - Module: github.com/kyanite/focus
  - Package/dir: cmd/focus, internal/*, pkg/* under C:/Users/Simon/dev/crush-cli
- Go version: 1.23.6
- LOC (approximate): ~9,000–11,000 (TUI main ~2k+)

## 2. Directory Structure
```
C:/Users/Simon/dev/crush-cli/
├─ cmd/
│  ├─ focus/
│  └─ focus-mcp/
├─ internal/
│  ├─ ai/  cli/  engine/  journal/  store/  theme/  tui/  wizards/
├─ pkg/
│  ├─ audio/ calendar/ config/ glow/ gum/ models/ styles/ utils/ validation/
├─ assets/
│  └─ themes/
├─ docs/
├─ test/
├─ main.go
├─ go.mod / go.sum
├─ build.sh / build.bat
└─ README.md
```

## 3. Features Status
- Tasks
  - Status: ✅ Done
  - Files: internal/engine/*, internal/store/*, internal/cli/{add,list,complete,delete}.go, pkg/models/task.go
  - DB tables: n/a (JSON storage)
  - Keybindings: dashboard up/down/enter, d complete, p priority, a add
- Notes
  - Status: ⏳ Partial (TUI notes editing; CLI command present)
  - Files: internal/cli/notes.go, internal/tui/main.go (notesView)
  - DB tables: n/a
  - Keybindings: n open notes; Enter/Esc save/cancel
- Calendar
  - Status: ✅ Done (month/week/day, integrated in TUI)
  - Files: pkg/calendar/*, internal/cli/calendar.go, internal/tui/main.go (calendarView)
  - DB tables: n/a
  - Keybindings: k open; ←/h prev; →/l next; Tab cycles
- Themes
  - Status: ✅ Done
  - Files: pkg/styles/*, internal/theme/registry.go, assets/themes/*
  - DB tables: n/a
  - Keybindings: ctrl+t (cycle); t (unified dashboard)
- AI
  - Status: ⏳ Partial→Strong (Ollama primary; OpenRouter fallback)
  - Files: internal/ai/manager.go, internal/cli/{chat,suggest,summary}.go, internal/tui/main.go (chatView)
  - DB tables: n/a; cache: ~/.focus/ai_cache.json
  - Keybindings: c chat; Enter send; Esc back
- Pomodoro/Timer
  - Status: ⏳ Partial (TUI Focus view timers)
  - Files: internal/tui/main.go (timer/stopwatch)
  - DB tables: n/a
  - Keybindings: s start, x stop, r reset, b/space switch
- Wizards
  - Status: ✅ Done (Huh/Gum wizards)
  - Files: internal/wizards/task_creation.go, internal/cli/{wizards,enhanced_config,interactive,simple-interactive}.go
- Journal/Export
  - Status: ⏳ Partial (export to syntax.sh)
  - Files: internal/journal/export.go, pkg/models/journal.go, internal/cli/journal.go
- Unified Dashboard
  - Status: ✅ Done
  - Files: internal/tui/unified_dashboard.go, internal/cli/unified_dashboard.go
- MCP server
  - Status: ⏳ Partial
  - Files: cmd/focus-mcp/main.go, internal/cli/mcp_server.go

## 4. Database
- Tables (SQL schema): none (JSON persistence)
- File location: ~/.focus/tasks.json
- Migrations: best-effort migration from ~/.neon and ~/.todo into ~/.focus (pkg/utils/storage.go)

## 5. UI Architecture
- Screens: dashboardView, focusView, chatView, calendarView, notesView, settingsView (internal/tui/main.go)
- Navigation: Tab cycles views; arrow keys navigate; global keys a,d,p,/,n,c,k,ctrl+t,q
- Theming: pkg/styles manages theme state/colors; internal/theme/registry.go; assets/themes presets

## 6. AI Integration
- Provider: Ollama (local), OpenRouter fallback (optional)
- Models: qwen2.5:1.5b (primary); OpenRouter uses gpt-3.5-turbo
- Features using AI: parse (planned in manager), suggest, summarize, chat; TUI ChatAssistant uses manager
- Prompts: internal/ai/manager.go (buildParsePrompt/buildSuggestPrompt/buildSummaryPrompt/buildChatPrompt)

## 7. Dependencies
```
(go list -m all | head -20 equivalent from go.mod)
github.com/kyanite/focus
github.com/charmbracelet/bubbles v0.21.1-...
github.com/charmbracelet/bubbletea v1.3.6
github.com/charmbracelet/huh v0.8.0
github.com/charmbracelet/lipgloss v1.1.1-...
github.com/pterm/pterm v0.12.76
github.com/spf13/cobra v1.9.1
github.com/spf13/viper v1.21.0
(see go.mod for full list)
```
Highlight: Bubble Tea, Lipgloss, Huh, Cobra/Viper; SQLite driver: none (JSON)

## 8. Build
- Build command: go build -o focus ./cmd/focus
- Binary name: focus
- Cross-platform: Win/Mac/Linux (scripts provided). Ollama required for AI features.

## 9. Naming Issues
- "neon"
  - Docs/scripts and migration mention neon (StartNEON.bat, install-neon.bat, .vbs shortcuts, assets/manifest.txt; pkg/utils/storage.go migration from ~/.neon)
- "focus"
  - Module and code consistently use github.com/kyanite/focus; CLI/TUI branding focus.sh across internal/* and pkg/*
- "lyricforge"
  - None found

Format: [file]: [context] - [type]
- pkg/utils/storage.go: copies from ~/.neon → ~/.focus - code
- StartNEON.bat/install-neon.bat/update-*.vbs: neon binary paths - scripts
- assets/manifest.txt: neon-glow.png - asset ref
- Multiple docs (NEON_DEVELOPMENT_PLAN.md, DEPLOYMENT_SUMMARY.md) - docs

## 10. Red Flags
- Critical
  - Tests reference undefined validation helpers (ValidateTask, ValidateTaskDeadline, ValidateTaskCategories, ValidateEmail, ValidateUrl, ValidatePhoneNumber) causing compile/test failures. Implement or update tests.
  - go 1.23.6 in go.mod may not match installed toolchains; verify environment.
- Performance
  - Large monolithic TUI (internal/tui/main.go ~2k+ LOC) suggests modularization need.
  - String-building heavy rendering; acceptable but watch large datasets.
- Pre-ship cleanup
  - Remove legacy neon scripts or align branding; ensure build outputs and docs consistently use focus.
  - Review automatic Ollama install/launch in cmd/focus/main.go for side effects in managed environments.
