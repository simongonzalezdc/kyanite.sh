# Kyanite Standards & Implementation Guide

**Shared documentation for all Kyanite Suite tools**  
**Version:** 1.0  
**Date:** November 2025

---

## Universal Shortcuts (All Tools)

| Shortcut | Action |
|----------|--------|
| **Ctrl+Q** | Quit |
| **Ctrl+H** | Help |
| **Ctrl+Shift+T** | Switch theme |
| **Ctrl+S** | Manual save |
| **Esc** | Back/Cancel |

**Note:** Undo/Redo (Ctrl+Z/Ctrl+Y) is planned for v2.0+ and not required for v1.0 releases.

---

## 10 Kyanite Themes (Copy for all tools)

All 10 themes use exact hex codes below:

### Theme System Implementation

Each tool must include:

1. **internal/theme/registry.go** - All 10 themes with hex codes
2. **internal/theme/manager.go** - Runtime theme switching
3. **internal/theme/types.go** - Theme struct

### The 10 Themes

```go
// 1. Monochrome
Primary: #E0E0E0
Secondary: #B0B0B0
Accent: #FFFFFF
Background: #1A1A1A
Text: #E8E8E8
Success: #90C695

// 2. Amber Night
Primary: #D4A574
Secondary: #8B7355
Accent: #C0956E
Background: #2A2421
Text: #F5E6D3
Success: #A8C686

// 3. Twilight Mist
Primary: #9D84B7
Secondary: #7B68A6
Accent: #C9AEE8
Background: #2D263E
Text: #E8DFF5
Success: #A8C686

// 4. Indigo Depths
Primary: #4169E1
Secondary: #5F9EA0
Accent: #DEB887
Background: #0C0C1E
Text: #F0E68C
Success: #5F9EA0

// 5. Forest Path
Primary: #52B788
Secondary: #52A068
Accent: #95D5B2
Background: #1B263B
Text: #D8F3DC
Success: #B7E4C7

// 6. Clay Earth
Primary: #A0644E
Secondary: #8B4513
Accent: #D2691E
Background: #2F1F1F
Text: #F5E6D3
Success: #9ACD32

// 7. Iron Forge
Primary: #6B4423
Secondary: #8B4513
Accent: #CD853F
Background: #1A1A1A
Text: #F5DEB3
Success: #98FB98

// 8. Sunlight
Primary: #FFB92C
Secondary: #FF8C00
Accent: #FFA500
Background: #1F1F1F
Text: #FFFACD
Success: #32CD32

// 9. Cyan Wave
Primary: #00D4FF
Secondary: #00B8D4
Accent: #00F5FF
Background: #0D1F26
Text: #E0FFFF
Success: #00FF7F

// 10. Electric Rose
Primary: #FF0080
Secondary: #D4005C
Accent: #FF1493
Background: #1A0E1A
Text: #FFE0F0
Success: #39FF14
```

---

## Architecture Patterns

### Bubble Tea Structure

```go
// Every tool uses this pattern
type RootModel struct {
    currentScreen Screen
    width         int
    height        int
    // Child models...
}

func (m RootModel) Init() tea.Cmd {
    return tea.EnterAltScreen
}

func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Route to child models
    return m, nil
}

func (m RootModel) View() string {
    // Render current screen
    return ""
}
```

**CRITICAL RULES:**
- State changes ONLY in Update()
- Async operations via tea.Cmd
- NO goroutines without message channel
- Cache expensive View() calculations

---

## Configuration Management

### Config File Location (Cross-Platform)

**Linux:**
```
~/.config/{tool}/config.toml
```

**macOS:**
```
~/Library/Application Support/{tool}/config.toml
```

**Windows:**
```
%APPDATA%/{tool}/config.toml
```

**Implementation:**
```go
import "os"

func GetConfigDir(tool string) string {
    configDir, err := os.UserConfigDir()
    if err != nil {
        // Fallback to home directory
        homeDir, _ := os.UserHomeDir()
        return filepath.Join(homeDir, ".config", tool)
    }
    return filepath.Join(configDir, tool)
}
```

### Standard Config Structure

```toml
[preferences]
theme = "amber-night"
auto_save = true
auto_save_interval = 30

[display]
show_help_on_startup = false
```

### Implementation

```go
func LoadConfig(tool string) (Config, error) {
    homeDir, _ := os.UserHomeDir()
    cfgPath := filepath.Join(homeDir, ".config", tool, "config.toml")
    
    // Create if not exists
    if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
        os.MkdirAll(filepath.Dir(cfgPath), 0755)
        // Create default config
    }
    
    // Load and parse
    return parseConfig(cfgPath)
}
```

---

## File Storage

### Principle

Use XDG Base Directory spec (with cross-platform support):

**Linux (XDG):**
```
~/.config/{tool}/         (Config)
~/.local/share/{tool}/    (Data)
~/.cache/{tool}/          (Temporary)
```

**macOS:**
```
~/Library/Application Support/{tool}/  (Config & Data)
~/Library/Caches/{tool}/               (Temporary)
```

**Windows:**
```
%APPDATA%/{tool}/          (Config & Data)
%LOCALAPPDATA%/{tool}/     (Temporary)
```

**Go Implementation:**
```go
// Use os.UserConfigDir() for config
// Use os.UserCacheDir() for cache
// Create custom function for data directory
```

### Storage Layer Example

```go
func SaveToConfig(tool string, key string, data interface{}) error {
    homeDir, _ := os.UserHomeDir()
    configDir := filepath.Join(homeDir, ".config", tool)
    os.MkdirAll(configDir, 0755)
    
    file := filepath.Join(configDir, "config.toml")
    // Write data
    return ioutil.WriteFile(file, data, 0644)
}
```

---

## Error Handling Pattern

**NEVER panic:**

```go
// ❌ WRONG
if err != nil {
    panic(err)
}

// ✅ CORRECT
if err != nil {
    log.Printf("Operation failed: %v", err)
    m.err = "Unable to save. Check permissions."
    return m, nil  // Continue gracefully
}
```

**User-Friendly Messages:**

```go
// ❌ WRONG: "JSON parse error at line 42"
// ✅ CORRECT: "Config file is invalid. Try: rm ~/.config/{tool}/config.toml"
```

---

## Testing Template

### Test File

```go
package {module}_test

import (
    "testing"
    "{module}"
)

func TestFeatureName(t *testing.T) {
    // Arrange
    input := "test data"
    expected := "expected result"
    
    // Act
    result := module.Function(input)
    
    // Assert
    if result != expected {
        t.Errorf("got %v, want %v", result, expected)
    }
}
```

### Running Tests

```bash
go test -v ./...
go test -cover ./...
go test -run TestName ./...
```

---

## Performance Guidelines

### Optimization Rules

| Operation | Target | Rule |
|-----------|--------|------|
| Startup | <1s | Only initialize essentials |
| Keystroke | <100ms | No heavy compute in View() |
| Save | <100ms | Use atomic writes, make async if needed |
| Search | <100ms | Index or cache results |
| Memory (idle) | <50MB | Profile with pprof |

### Common Optimizations

```go
// ❌ SLOW: Expensive computation in View()
func (m Model) View() string {
    expensiveResult := calculateSomething()  // SLOW!
    return expensiveResult
}

// ✅ FAST: Cache the result
type Model struct {
    cached string
    needsRender bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Only recalculate when model changes
    m.needsRender = true
    return m, nil
}

func (m Model) View() string {
    if m.needsRender {
        m.cached = calculateSomething()
        m.needsRender = false
    }
    return m.cached
}
```

---

## Dependencies

### Go Modules

```bash
go mod init github.com/kyanite/{tool}

go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
```

### No External Dependencies Allowed

❌ Avoid:
- Complex frameworks
- CGO dependencies
- Huge dependency trees
- Outdated libraries

✅ Use:
- Standard library first
- Charm Bracelet ecosystem
- Well-maintained, popular packages
- Pure Go packages

---

## Project Setup Checklist

- [ ] `go mod init github.com/kyanite/{tool}`
- [ ] Install required dependencies
- [ ] Create directory structure
- [ ] Copy theme system (registry.go, manager.go, types.go)
- [ ] Create main.go entry point
- [ ] Create go.mod and go.sum
- [ ] Create Makefile
- [ ] Create tests/ directory
- [ ] Create data/ directory (if needed)
- [ ] First compile: `go build -o bin/{tool} cmd/{tool}/main.go`

---

## Debugging

### Enable Logging

```go
import "log"

logFile, _ := os.Create("/tmp/{tool}-debug.log")
log.SetOutput(logFile)
log.Printf("Debug info: %v", value)
```

**View logs:**
```bash
tail -f /tmp/{tool}-debug.log
```

### Profile Performance

```bash
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof
```

### Terminal Debugging

```bash
# Set terminal to 80x24
resize -s 24 80

# Run tool
./bin/{tool}

# In another terminal, monitor
watch -n 1 'ps aux | grep {tool}'
```

---

## Git Repository Setup

### Repository Naming

```
github.com/kyanite/prism
github.com/kyanite/syntax
```

### .gitignore Template

```
bin/
*.out
*.test
.DS_Store
/.vscode
/vendor
*.prof
/tmp
```

### README Template

```markdown
# {tool}.sh

> One-line description

## Quick Start

\`\`\`bash
go install github.com/kyanite/{tool}@latest
{tool}
\`\`\`

## Features

- Feature 1
- Feature 2

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| Ctrl+Q | Quit |
| Ctrl+H | Help |

## Part of Kyanite Suite

See: https://kyanite.sh

## License

MIT
```

---

## Release Checklist

### Before v1.0

- [ ] All core features implemented
- [ ] All acceptance criteria met
- [ ] 0 critical bugs
- [ ] 10 themes working (for app UI styling)
- [ ] Universal shortcuts working (Ctrl+Q, Ctrl+H, Ctrl+Shift+T, Ctrl+S, Esc)
- [ ] Performance targets met (startup <1s, operations <500ms, memory <50MB)
- [ ] All tests passing (>70% coverage)
- [ ] Documentation complete (README, ARCHITECTURE, CONTRIBUTING, LICENSE)
- [ ] Functional UI in 80x24 terminal minimum
- [ ] No panics in application code (all external panics recovered gracefully)
- [ ] Cross-platform support (Linux, macOS, Windows)

### GitHub Release

```bash
# Tag release
git tag v1.0
git push origin v1.0

# Create release on GitHub
# Upload binaries for multiple platforms
```

---

## CI/CD Best Practices

### Required GitHub Actions

All Kyanite tools should include:

1. **Test workflow** - Run tests on every push/PR
2. **Build workflow** - Build for Linux, macOS, Windows (amd64 & arm64)
3. **Release workflow** - Automated releases on version tags

### Example Workflows Directory Structure

```
.github/
└── workflows/
    ├── test.yml       (Tests + coverage check)
    ├── build.yml      (Multi-platform builds)
    └── release.yml    (Create GitHub releases)
```

### Test Workflow Template

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -v -race -coverprofile=coverage.out ./...
      - run: |
          coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          if (( $(echo "$coverage < 70" | bc -l) )); then
            echo "Coverage below 70%"; exit 1
          fi
```

### Makefile Standard

All tools should include a Makefile with:

```makefile
.PHONY: test build install clean lint fmt run

test:
	go test -v -race -coverprofile=coverage.out ./...

build:
	go build -o bin/{tool} ./cmd/{tool}

install:
	go install ./cmd/{tool}

clean:
	rm -rf bin/ coverage.out

lint:
	golangci-lint run

fmt:
	go fmt ./...

run:
	go run ./cmd/{tool}
```

---

## Common Issues

### Problem: Colors Wrong in Terminal

**Solution:**
```bash
export TERM=xterm-256color
export LANG=en_US.UTF-8
```

### Problem: Build Fails

**Solution:**
```bash
go mod tidy
go mod download
go mod vendor  # If needed
go build ./...
```

### Problem: Slow Performance

**Solution:**
```bash
# Profile the code
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Common issues:
# - Too many allocations (use sync.Pool)
# - View() too expensive (cache results)
# - File I/O blocking (use async)
```

---

## Key Principles

1. **Speed First** - <100ms response time
2. **Simplicity** - One clear path, not ten options
3. **Export** - Tools flow into each other
4. **Offline** - Works without internet
5. **Local Files** - Data stays on user's computer
6. **ADHD-Friendly** - No overwhelming options
7. **Beautiful** - Terminal can be beautiful
8. **Useful** - Real workflow, not toy

---

## Questions?

If something is unclear:

1. Check the tool's PRD (feature details)
2. Check the tool's TDD (architecture)
3. Check this guide (standards)
4. Look at code examples
5. Run tests

---

**Every tool follows these standards. Consistency across the suite is important.**

**When in doubt, prioritize: speed, simplicity, and user experience.**
