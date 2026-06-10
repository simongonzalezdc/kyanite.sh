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
| **Ctrl+Z** | Undo |
| **Ctrl+Y** | Redo |
| **Esc** | Back/Cancel |

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

### Config File Location

```
~/.config/{tool}/config.toml
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

### Cross-Platform Directory Paths

**Platform-Specific Locations:**

**Linux/macOS:**
```
Config:  ~/.config/{tool}/
Data:    ~/.local/share/{tool}/
Cache:   ~/.cache/{tool}/
```

**Windows:**
```
Config:  %APPDATA%/{tool}/           (e.g., C:\Users\Name\AppData\Roaming\syntax)
Data:    %LOCALAPPDATA%/{tool}/      (e.g., C:\Users\Name\AppData\Local\syntax)
Cache:   %TEMP%/{tool}/              (e.g., C:\Users\Name\AppData\Local\Temp\syntax)
```

### Cross-Platform Implementation

**Recommended Library:** `github.com/adrg/xdg`

```go
import (
    "path/filepath"
    "github.com/adrg/xdg"
)

// Cross-platform config directory
func GetConfigDir(tool string) string {
    return filepath.Join(xdg.ConfigHome, tool)
}

func GetDataDir(tool string) string {
    return filepath.Join(xdg.DataHome, tool)
}

func GetCacheDir(tool string) string {
    return filepath.Join(xdg.CacheHome, tool)
}

// Storage layer example
func SaveToConfig(tool string, filename string, data []byte) error {
    configDir := GetConfigDir(tool)

    // Create directory with proper permissions
    if err := os.MkdirAll(configDir, 0700); err != nil {
        return fmt.Errorf("failed to create config dir: %w", err)
    }

    filePath := filepath.Join(configDir, filename)

    // Write with restricted permissions (owner only)
    if err := os.WriteFile(filePath, data, 0600); err != nil {
        return fmt.Errorf("failed to write config: %w", err)
    }

    return nil
}
```

### File Permissions

- **Config files:** 0600 (owner read/write only)
- **Data directories:** 0700 (owner only)
- **Cache files:** 0600 (owner read/write only)
- **Exported files:** 0644 (owner write, all read)

**Windows Note:** Permissions are handled differently on Windows (ACLs). The `os` package handles this automatically.

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

## Security Standards

### Sensitive Data Storage

**API Keys and Secrets:**

1. **Primary:** Use OS-native keyring/credential manager
   ```go
   import "github.com/zalando/go-keyring"

   // Store API key securely
   func StoreAPIKey(service, username, key string) error {
       return keyring.Set(service, username, key)
   }

   // Retrieve API key
   func GetAPIKey(service, username string) (string, error) {
       return keyring.Get(service, username)
   }
   ```

2. **Fallback:** Encrypted config file with user warning
   ```go
   // Only if keyring unavailable (rare)
   func encryptAndSave(key string) error {
       // Use AES-256 encryption
       // Password derived from machine-specific data
       // Show warning: "API key stored in encrypted file. Use at your own risk."
       return nil
   }
   ```

3. **Never:** Plain text storage

### Input Sanitization

**File Paths:**
```go
func SanitizeFilename(input string) string {
    // Remove path separators
    safe := strings.ReplaceAll(input, "/", "-")
    safe = strings.ReplaceAll(safe, "\\", "-")
    safe = strings.ReplaceAll(safe, "..", "")

    // Remove special characters
    safe = strings.Map(func(r rune) rune {
        if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' || r == ' ' {
            return r
        }
        return -1
    }, safe)

    // Limit length
    if len(safe) > 255 {
        safe = safe[:255]
    }

    return filepath.Clean(safe)
}

// Reject path traversal attempts
func ValidateProjectPath(path string) error {
    cleaned := filepath.Clean(path)

    if strings.Contains(cleaned, "..") {
        return errors.New("path traversal detected")
    }

    if filepath.IsAbs(path) && !isWithinAllowedDir(path) {
        return errors.New("absolute path outside allowed directories")
    }

    return nil
}
```

**User Input:**
```go
// Character names, scene titles, etc.
func SanitizeUserInput(input string, maxLen int) string {
    // Trim whitespace
    input = strings.TrimSpace(input)

    // Limit length
    if len(input) > maxLen {
        input = input[:maxLen]
    }

    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")

    return input
}
```

### File Permissions (Security)

**On File Creation:**
```go
// Config files (may contain sensitive data)
os.WriteFile(path, data, 0600)  // Owner read/write only

// Data directories
os.MkdirAll(path, 0700)  // Owner only

// Exported files (meant to be shared)
os.WriteFile(path, data, 0644)  // Standard permissions
```

### Denial of Service Prevention

**Resource Limits:**
```go
const (
    MaxFileSize      = 50 * 1024 * 1024  // 50MB
    MaxCharacters    = 1000              // Per project
    MaxScenes        = 10000             // Per project
    MaxSearchResults = 100               // Search result limit
)

func LoadScene(path string) (*Scene, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }

    if info.Size() > MaxFileSize {
        return nil, errors.New("file too large")
    }

    // Read with size limit
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    return ParseScene(data)
}
```

### Dependency Security

**Regular Audits:**
```bash
# Use govulncheck (official Go vulnerability scanner)
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Alternative: nancy (third-party)
go list -json -m all | nancy sleuth
```

**Update Policy:**
- Review dependencies monthly
- Update immediately for CVE fixes
- Test thoroughly before updating major versions
- Avoid unmaintained dependencies

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
| Save | <100ms | Make non-blocking async |
| Search | <100ms | Index or cache results |
| Memory | <50MB | Profile with pprof |

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

## Logging Strategy

### Log Levels

- **ERROR:** Critical failures requiring user attention
- **WARN:** Recoverable issues, logged silently
- **INFO:** Major operations (startup, save, export)
- **DEBUG:** Detailed tracing (development only, disabled in production)

### Production Logging

**Error Logging Only:**
```go
import (
    "log"
    "os"
    "path/filepath"
    "github.com/adrg/xdg"
)

var errorLog *log.Logger

func InitLogging(tool string) error {
    // Log errors to cache directory
    logDir := filepath.Join(xdg.CacheHome, tool)
    os.MkdirAll(logDir, 0700)

    logPath := filepath.Join(logDir, "errors.log")

    // Rotate if > 10MB
    if info, err := os.Stat(logPath); err == nil {
        if info.Size() > 10*1024*1024 {
            os.Rename(logPath, logPath+".old")
        }
    }

    file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
    if err != nil {
        return err
    }

    errorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
    return nil
}

func LogError(format string, v ...interface{}) {
    if errorLog != nil {
        errorLog.Printf(format, v...)
    }
}
```

### Development Logging

```go
// Set via environment variable
func InitDevLogging() {
    if os.Getenv("DEBUG") == "1" {
        logFile, _ := os.Create("/tmp/syntax-debug.log")
        log.SetOutput(logFile)
        log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
        log.Println("Debug logging enabled")
    }
}
```

### User Access to Logs

```bash
# Add CLI command
syntax --show-logs  # Opens log directory in file manager
syntax --debug      # Run with debug output to stderr
```

---

## Accessibility

### Keyboard Navigation

**Requirements:**
- ALL features must work without mouse
- Tab navigation through UI elements
- Vim-style shortcuts encouraged (h/j/k/l for navigation)
- Clear focus indicators

**Example:**
```go
// Handle keyboard-only navigation
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "tab":
            m.focusNext()
        case "shift+tab":
            m.focusPrev()
        case "j", "down":
            m.selectNext()
        case "k", "up":
            m.selectPrev()
        }
    }
    return m, nil
}
```

### Screen Readers

**v1.0:**
- Limited support (terminal-based, inherits terminal's reader)
- Use semantic text output
- Avoid ASCII art for critical info

**v2.0 (Future):**
- Consider explicit accessibility hints
- Alternative text for ASCII visualizations
- ARIA-like attributes for TUI

### Visual Accessibility

**Color Contrast:**
- All themes must pass WCAG AA contrast requirements
- Minimum contrast ratio: 4.5:1 (text/background)
- Test with: https://webaim.org/resources/contrastchecker/

**High Contrast Mode:**
- Force "monochrome" theme for maximum contrast
- Triggered by environment variable or flag
  ```bash
  SYNTAX_HIGH_CONTRAST=1 syntax
  ```

**Text Size:**
- Respect terminal font size settings
- No hardcoded small fonts
- UI scales with terminal size

---

## Semantic Versioning

### Format

**vMAJOR.MINOR.PATCH**

- **MAJOR:** Breaking changes (data format, API, major features removed)
- **MINOR:** New features (backward compatible)
- **PATCH:** Bug fixes only (no new features)

### Examples

- `v1.0.0 → v1.1.0`: Added spell check feature (new feature)
- `v1.1.0 → v1.1.1`: Fixed crash bug (bug fix)
- `v1.5.0 → v2.0.0`: Changed storage format (breaking)

### Version Consistency

**In Code:**
```go
const Version = "1.0.0"
```

**In Git:**
```bash
git tag v1.0.0
git push origin v1.0.0
```

**In Releases:**
- GitHub release: `v1.0.0`
- Binary name: `syntax-v1.0.0-linux-amd64`

### Deprecation Policy

1. Features marked deprecated in v1.x
2. Warnings shown to users for 1+ minor releases
3. Removed in next major version (v2.0)
4. Migration guide provided

**Example:**
```go
// v1.5: Mark as deprecated
func OldFunction() {
    log.Println("Warning: OldFunction is deprecated. Use NewFunction instead. Will be removed in v2.0")
    // ... implementation
}

// v2.0: Remove completely
```

---

## Dependency Management

### Security Scanning

**Before Each Release:**
```bash
# Official Go vulnerability scanner
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Check for known vulnerabilities
go list -json -m all | nancy sleuth
```

### Update Policy

1. **Security Patches:** Immediate update
2. **Minor Updates:** Monthly review
3. **Major Updates:** Test thoroughly in development branch

### Minimal Dependencies

**Prefer:**
- Standard library first
- Well-maintained packages (active commits, responsive maintainers)
- Popular packages with security track record
- Pure Go packages (no CGO)

**Avoid:**
- Complex frameworks
- CGO dependencies (platform-specific)
- Huge dependency trees
- Unmaintained libraries (>1 year no commits)

**Check Dependency Tree:**
```bash
go mod graph | grep -v 'go.mod'
go mod why -m <package>
```

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

### Profile Performance (Extended)

**CPU Profiling:**
```bash
# Profile tests
go test -cpuprofile=cpu.prof -bench=. ./...

# Visualize in browser
go tool pprof -http=:8080 cpu.prof

# Text output
go tool pprof cpu.prof
> top10
> list functionName
```

**Memory Profiling:**
```bash
# Profile memory allocation
go test -memprofile=mem.prof -bench=. ./...

# Analyze
go tool pprof -http=:8080 mem.prof

# Check for leaks
go tool pprof -alloc_space mem.prof
> top10
```

**Memory Leak Detection:**
```bash
# Run long-running test
go test -memprofile=mem.prof -run=TestLongSession

# Check if memory grows over time
go tool pprof -alloc_space mem.prof

# Look for growing allocations
> top10
> list suspiciousFunction
```

**Benchmarking Standards:**
```go
func BenchmarkBufferInsert(b *testing.B) {
    buf := editor.NewBuffer("initial")
    b.ResetTimer()  // Exclude setup time

    for i := 0; i < b.N; i++ {
        buf.Insert(5, "text")
    }
}

// Run benchmarks
// go test -bench=. -benchmem ./...

// Target performance:
// BenchmarkBufferInsert  1000000  1000 ns/op  64 B/op  1 allocs/op
//                        ^        ^           ^        ^
//                        ops      ns/op       bytes    allocations
```

### Terminal Debugging

```bash
# Set terminal to 80x24
resize -s 24 80

# Run tool
./bin/{tool}

# In another terminal, monitor
watch -n 1 'ps aux | grep {tool}'

# Check file descriptor leaks
lsof -p $(pgrep {tool})

# Monitor system calls
strace -p $(pgrep {tool})  # Linux
dtruss -p $(pgrep {tool})  # macOS
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
- [ ] 10 themes working
- [ ] Universal shortcuts working
- [ ] Performance targets met
- [ ] All tests passing (>70% coverage)
- [ ] Documentation complete
- [ ] Works on 80x24 terminal
- [ ] No panics or crashes

### GitHub Release

```bash
# Tag release
git tag v1.0
git push origin v1.0

# Create release on GitHub
# Upload binaries for multiple platforms
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

## CI/CD Pipeline

### GitHub Actions Workflow

**Create:** `.github/workflows/test.yml`

```yaml
name: Test and Build

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    name: Test
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22']

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go }}

      - name: Run tests
        run: go test -v -cover -race ./...

      - name: Run vet
        run: go vet ./...

      - name: Check formatting
        run: |
          gofmt -l .
          test -z "$(gofmt -l .)"

      - name: Vulnerability scan
        run: |
          go install golang.org/x/vuln/cmd/govulncheck@latest
          govulncheck ./...

  build:
    name: Build
    runs-on: ubuntu-latest
    needs: test

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build
        run: |
          go build -v -o bin/syntax ./cmd/syntax

      - name: Upload artifact
        uses: actions/upload-artifact@v3
        with:
          name: syntax-binary
          path: bin/syntax
```

### Pre-commit Hooks

**Setup:** `.git/hooks/pre-commit`

```bash
#!/bin/bash

echo "Running pre-commit checks..."

# Format code
echo "Formatting code..."
go fmt ./...

# Vet code
echo "Running go vet..."
go vet ./...
if [ $? -ne 0 ]; then
    echo "go vet failed"
    exit 1
fi

# Run tests
echo "Running tests..."
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed"
    exit 1
fi

# Check for TODO/FIXME in staged files
echo "Checking for TODO/FIXME..."
git diff --cached --name-only | xargs grep -l "TODO\|FIXME" && {
    echo "Found TODO/FIXME in staged files"
    exit 1
}

echo "Pre-commit checks passed!"
exit 0
```

**Make executable:**
```bash
chmod +x .git/hooks/pre-commit
```

### Release Automation

**Create:** `.github/workflows/release.yml`

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Build binaries
        run: |
          # Linux
          GOOS=linux GOARCH=amd64 go build -o syntax-linux-amd64 ./cmd/syntax
          GOOS=linux GOARCH=arm64 go build -o syntax-linux-arm64 ./cmd/syntax

          # macOS
          GOOS=darwin GOARCH=amd64 go build -o syntax-macos-amd64 ./cmd/syntax
          GOOS=darwin GOARCH=arm64 go build -o syntax-macos-arm64 ./cmd/syntax

          # Windows
          GOOS=windows GOARCH=amd64 go build -o syntax-windows-amd64.exe ./cmd/syntax

      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            syntax-linux-amd64
            syntax-linux-arm64
            syntax-macos-amd64
            syntax-macos-arm64
            syntax-windows-amd64.exe
          draft: false
          prerelease: false
```

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
