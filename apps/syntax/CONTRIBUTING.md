# Contributing to syntax.sh

Thank you for your interest in contributing to syntax.sh! This document provides guidelines and instructions for contributing to the project.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Setup](#development-setup)
4. [How to Contribute](#how-to-contribute)
5. [Coding Standards](#coding-standards)
6. [Testing](#testing)
7. [Pull Request Process](#pull-request-process)
8. [Documentation](#documentation)
9. [Issue Guidelines](#issue-guidelines)
10. [Community](#community)

---

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of:
- Experience level
- Background
- Identity
- Location

### Our Standards

**Expected Behavior:**
- Be respectful and considerate
- Welcome newcomers and help them learn
- Accept constructive criticism gracefully
- Focus on what's best for the project
- Show empathy towards others

**Unacceptable Behavior:**
- Harassment, discrimination, or personal attacks
- Trolling, insulting comments, or off-topic discussions
- Publishing others' private information
- Spam or self-promotion

**Enforcement:**
- First violation: Warning
- Second violation: Temporary ban
- Third violation: Permanent ban

**Reporting:** Email conduct@kyanite.sh

---

## Getting Started

### Prerequisites

Before contributing, ensure you have:

1. **Go 1.21+** installed
   ```bash
   go version  # Should show 1.21 or higher
   ```

2. **Git** installed
   ```bash
   git --version
   ```

3. **Basic Go knowledge**
   - Understanding of Go syntax
   - Familiarity with Go modules
   - Knowledge of goroutines and channels (for async operations)

4. **Terminal emulator** for testing TUI
   - Linux: GNOME Terminal, Alacritty, or kitty
   - macOS: Terminal.app or iTerm2
   - Windows: Windows Terminal

### First-Time Contributors

**Good First Issues:**
- Look for issues labeled `good first issue` or `beginner-friendly`
- These are typically small, well-defined tasks
- Great for learning the codebase

**Recommended Reading:**
1. README.md - Project overview
2. ARCHITECTURE.md - Code structure
3. KYANITE-STANDARDS.md - Coding standards
4. This file - Contribution process

---

## Development Setup

### 1. Fork and Clone

```bash
# Fork the repository on GitHub first, then:

git clone https://github.com/YOUR_USERNAME/syntax.git
cd syntax
```

### 2. Install Dependencies

```bash
# Initialize Go modules
go mod download

# Verify everything works
go build ./cmd/syntax
```

### 3. Create a Branch

```bash
# Create feature branch
git checkout -b feature/my-feature-name

# Or bug fix branch
git checkout -b fix/issue-123-description
```

**Branch Naming:**
- Features: `feature/short-description`
- Bug fixes: `fix/issue-number-description`
- Documentation: `docs/what-you-changed`
- Tests: `test/what-you-tested`

### 4. Set Up Development Environment

```bash
# Enable debug logging
export DEBUG=1

# Run from source
go run ./cmd/syntax

# Or build and run
go build -o bin/syntax ./cmd/syntax
./bin/syntax
```

---

## How to Contribute

### Types of Contributions

We welcome:

1. **Bug Fixes**
   - Fix reported issues
   - Add tests to prevent regression

2. **New Features**
   - Check roadmap first (ROADMAP.md)
   - Discuss in GitHub Discussions before large features
   - Follow PRD for v1.0 features

3. **Documentation**
   - Fix typos, clarify explanations
   - Add examples
   - Improve README, guides

4. **Tests**
   - Increase test coverage
   - Add integration tests
   - Add benchmarks

5. **Performance**
   - Optimize slow operations
   - Reduce memory usage
   - Include benchmarks showing improvement

6. **UI/UX**
   - Improve terminal UI
   - Better keyboard shortcuts
   - Accessibility improvements

### Contribution Workflow

1. **Find or Create Issue**
   - Search existing issues first
   - If new, create issue describing problem/feature
   - Wait for maintainer approval for large features

2. **Discuss Approach**
   - Comment on issue with proposed solution
   - Get feedback before coding
   - Avoid duplicate work

3. **Write Code**
   - Follow coding standards (see below)
   - Add tests
   - Update documentation

4. **Test Locally**
   - Run all tests: `go test ./...`
   - Run on 80x24 terminal
   - Test on your OS

5. **Submit Pull Request**
   - See Pull Request Process below

---

## Coding Standards

### Go Style

**Follow Official Guidelines:**
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

**Key Rules:**

1. **Formatting:**
   ```bash
   # Format all code before committing
   go fmt ./...

   # Vet for common mistakes
   go vet ./...
   ```

2. **Naming:**
   ```go
   // ✅ GOOD: Clear, concise names
   func LoadCharacter(id string) (*Character, error)
   var characterCount int
   type EditorState struct

   // ❌ BAD: Unclear, abbreviated
   func LdChr(i string) (*Char, error)
   var cnt int
   type EdSt struct
   ```

3. **Error Handling:**
   ```go
   // ✅ GOOD: Always handle errors
   data, err := os.ReadFile(path)
   if err != nil {
       return fmt.Errorf("failed to read file: %w", err)
   }

   // ❌ BAD: Ignored errors
   data, _ := os.ReadFile(path)
   ```

4. **Comments:**
   ```go
   // ✅ GOOD: Explain WHY, not WHAT
   // Cache results to avoid recalculating on every render.
   // This improves performance from 100ms to <1ms per frame.
   if m.needsRender {
       m.cached = calculateExpensiveResult()
   }

   // ❌ BAD: Obvious comments
   // Set x to 5
   x := 5
   ```

5. **Function Length:**
   - Keep functions short (<50 lines preferred)
   - Extract helper functions for complex logic
   - One responsibility per function

### Project-Specific Standards

**1. Bubble Tea Patterns:**

```go
// ✅ GOOD: State changes only in Update()
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        m.counter++  // State change here
        return m, nil
    }
    return m, nil
}

// ❌ BAD: State changes in View()
func (m Model) View() string {
    m.counter++  // NEVER modify state in View()
    return fmt.Sprintf("Count: %d", m.counter)
}
```

**2. File Operations:**

```go
// ✅ GOOD: Atomic writes
func SaveScene(scene *Scene, path string) error {
    return AtomicWriteFile(path, data, 0600)
}

// ❌ BAD: Direct write (can corrupt data)
func SaveScene(scene *Scene, path string) error {
    return os.WriteFile(path, data, 0600)
}
```

**3. Error Messages:**

```go
// ✅ GOOD: User-friendly messages
if err != nil {
    return errors.New("Failed to save character. Check disk space and permissions.")
}

// ❌ BAD: Technical jargon
if err != nil {
    return errors.New("ENOSPC write() syscall failed")
}
```

**4. Performance:**

```go
// ✅ GOOD: Cache expensive operations
type Model struct {
    cached       string
    needsRender  bool
}

func (m Model) View() string {
    if m.needsRender {
        m.cached = expensiveRender()
        m.needsRender = false
    }
    return m.cached
}

// ❌ BAD: Recalculate every time
func (m Model) View() string {
    return expensiveRender()  // Called 60 times/second!
}
```

### Kyanite Suite Standards

Follow **KYANITE-STANDARDS.md** for:
- Universal keyboard shortcuts
- Theme system usage
- File storage conventions
- Cross-platform compatibility

---

## Testing

### Test Requirements

**All code changes must include tests.**

**Coverage Targets:**
- Critical paths (editor, storage): 90%+
- Business logic (character, scene): 80%+
- UI layer: 50%+
- Overall: 75%+

### Writing Tests

**Unit Tests:**
```go
func TestBufferInsert(t *testing.T) {
    // Arrange
    buf := editor.NewBuffer("Hello")

    // Act
    buf.Insert(5, " World")

    // Assert
    if buf.GetContent() != "Hello World" {
        t.Errorf("got %q, want %q", buf.GetContent(), "Hello World")
    }
}
```

**Table-Driven Tests:**
```go
func TestSanitizeFilename(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"simple", "test.txt", "test.txt"},
        {"spaces", "my file.txt", "my file.txt"},
        {"path traversal", "../../../etc/passwd", "etcpasswd"},
        {"special chars", "file@#$.txt", "file.txt"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := SanitizeFilename(tt.input)
            if got != tt.expected {
                t.Errorf("got %q, want %q", got, tt.expected)
            }
        })
    }
}
```

**Benchmarks:**
```go
func BenchmarkBufferInsert(b *testing.B) {
    buf := editor.NewBuffer("initial")
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        buf.Insert(5, "x")
    }
}
```

### Running Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run benchmarks
go test -bench=. -benchmem ./...

# Race condition detection
go test -race ./...

# Specific package
go test -v ./internal/editor/
```

### Test Files Location

```
internal/
├── editor/
│   ├── buffer.go
│   ├── buffer_test.go       # Unit tests
│   └── buffer_bench_test.go # Benchmarks
└── character/
    ├── character.go
    └── character_test.go

tests/
├── integration_test.go       # End-to-end tests
└── fixtures/                # Test data
```

---

## Pull Request Process

### Before Submitting

**Checklist:**
- [ ] Code follows style guidelines
- [ ] Tests added and passing
- [ ] Documentation updated
- [ ] Ran `go fmt ./...`
- [ ] Ran `go vet ./...`
- [ ] Ran `go test ./...`
- [ ] Tested on 80x24 terminal
- [ ] No merge conflicts with main

### PR Description Template

```markdown
## Description
Brief description of changes

## Related Issue
Fixes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Performance improvement
- [ ] Refactoring

## Testing
How was this tested?
- Unit tests added
- Manually tested on macOS/Linux/Windows
- Tested on 80x24 terminal

## Screenshots (if UI changes)
![Screenshot](url)

## Checklist
- [ ] Tests pass
- [ ] Documentation updated
- [ ] No breaking changes
```

### PR Review Process

1. **Automated Checks:**
   - CI runs tests on all platforms
   - Linters check code quality
   - Coverage report generated

2. **Maintainer Review:**
   - Code quality
   - Test coverage
   - Documentation
   - Performance impact

3. **Feedback:**
   - Address review comments
   - Push changes to same branch
   - Re-request review

4. **Merge:**
   - Maintainer merges when approved
   - PR closed, branch can be deleted

### What to Expect

**Timeline:**
- **Initial Review:** Within 3 business days
- **Follow-up:** 1-2 days for responses
- **Merge:** After approval and passing checks

**Possible Outcomes:**
- ✅ **Approved:** Merged immediately
- 🔄 **Changes Requested:** Address feedback and resubmit
- ❌ **Closed:** Doesn't fit project direction (rare, usually discussed first)

---

## Documentation

### What Needs Documentation

**Code Documentation:**
```go
// Package editor provides text editing functionality for syntax.sh.
//
// The editor supports:
// - Multi-line editing with undo/redo
// - Syntax highlighting for markdown
// - Efficient handling of large documents (10,000+ lines)
package editor

// Buffer represents an in-memory text buffer with edit history.
//
// Buffer is not thread-safe. All modifications must happen from the
// main UI goroutine.
type Buffer struct {
    // ...
}

// Insert adds text at the specified position.
//
// Position is the byte offset from the start of the buffer.
// Returns error if position is out of bounds.
func (b *Buffer) Insert(pos int, text string) error {
    // ...
}
```

**User Documentation:**
- README.md - Getting started
- User guides (if adding features)
- Help system (in-app)

**Developer Documentation:**
- Architecture decisions
- API changes
- Migration guides

### Documentation Standards

- Use clear, simple language
- Include code examples
- Add screenshots for UI features
- Keep README concise (link to detailed docs)
- Update CHANGELOG.md for user-facing changes

---

## Issue Guidelines

### Before Opening an Issue

1. **Search Existing Issues:**
   - Check if already reported
   - Add comment to existing issue instead of duplicating

2. **Check Documentation:**
   - README, ARCHITECTURE, FAQ
   - Might answer your question

3. **Verify on Latest Version:**
   - Ensure issue still exists in latest release

### Bug Report Template

```markdown
## Bug Description
Clear description of the bug

## Steps to Reproduce
1. Open syntax
2. Create new project
3. Click on...
4. See error

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: macOS 14.0 / Ubuntu 22.04 / Windows 11
- Terminal: iTerm2 / GNOME Terminal / Windows Terminal
- syntax.sh version: v1.0.0
- Go version: 1.21

## Screenshots
If applicable

## Logs
Error logs from ~/.cache/syntax/errors.log
```

### Feature Request Template

```markdown
## Feature Description
Clear description of proposed feature

## Use Case
Why is this needed? What problem does it solve?

## Proposed Solution
How might this work?

## Alternatives Considered
Other approaches you've thought of

## Additional Context
Screenshots, mockups, examples from other tools
```

### Issue Labels

- `bug` - Something isn't working
- `feature` - New feature request
- `documentation` - Documentation improvements
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `performance` - Performance improvement
- `question` - Question about usage
- `wontfix` - Will not be fixed (with explanation)

---

## Community

### Communication Channels

**GitHub Discussions:**
- Questions about usage
- Feature discussions
- Showcase your projects
- General chat

**GitHub Issues:**
- Bug reports
- Specific feature requests
- Implementation details

**Email:**
- Security issues: security@kyanite.sh
- Code of conduct: conduct@kyanite.sh

### Getting Help

**For Contribution Questions:**
- Comment on the issue you're working on
- Ask in GitHub Discussions
- Tag maintainers if urgent

**For Usage Questions:**
- Check documentation first
- Search GitHub Discussions
- Create new discussion if not found

### Recognition

**Contributors are recognized in:**
- CONTRIBUTORS.md (all contributors)
- Release notes (for specific releases)
- GitHub contributors page

**Significant Contributors:**
- Maintainer status (commit access)
- Decision-making participation
- Listed on project website

---

## Development Workflow

### Typical Contribution Flow

```bash
# 1. Sync with main
git checkout main
git pull upstream main

# 2. Create branch
git checkout -b fix/issue-123

# 3. Make changes
# ... edit code ...

# 4. Test
go test ./...
go run ./cmd/syntax

# 5. Commit
git add .
git commit -m "Fix character search bug (#123)"

# 6. Push
git push origin fix/issue-123

# 7. Open PR on GitHub

# 8. Address review feedback
# ... make changes ...
git add .
git commit -m "Address review comments"
git push origin fix/issue-123

# 9. Maintainer merges PR

# 10. Sync and delete branch
git checkout main
git pull upstream main
git branch -d fix/issue-123
```

### Commit Messages

**Format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation
- `test:` Tests
- `perf:` Performance
- `refactor:` Code restructuring
- `style:` Formatting

**Examples:**
```
feat: add character relationship map visualization

Implements ASCII art visualization of character relationships
using Unicode box-drawing characters. Supports tension levels
(low, medium, high) with different line styles.

Closes #42
```

```
fix: prevent crash on empty project load

Add validation to check if metadata.yaml exists before parsing.
Show user-friendly error message if file is missing.

Fixes #67
```

---

## Release Process

**For Maintainers:**

1. **Prepare Release:**
   - Update version in code
   - Update CHANGELOG.md
   - Run full test suite
   - Test on all platforms

2. **Create Tag:**
   ```bash
   git tag v1.0.1
   git push origin v1.0.1
   ```

3. **Build Binaries:**
   - GitHub Actions builds automatically
   - Attach to release

4. **Publish Release:**
   - Write release notes
   - Highlight key changes
   - Link to CHANGELOG

5. **Announce:**
   - GitHub Discussions
   - Project website
   - Social media

---

## Questions?

**Not sure where to start?**
- Look for `good first issue` label
- Ask in GitHub Discussions
- Read ARCHITECTURE.md to understand codebase

**Found a mistake in this guide?**
- Open an issue or PR to fix it!

---

**Thank you for contributing to syntax.sh! Your efforts help make this tool better for everyone. 🎉**

<!-- EMPOWER_ORCHESTRATOR:START -->
## Agent-law contribution rule

This repository follows the Empower Orchestrator law in `docs/agent-law/empower-orchestrator.md`.

If a change exposes a repeated task or repeated agent failure, contributors and agents should either ship the smallest durable prevention artifact or explain why this PR is intentionally one-off.

Automation and durable system changes require the scale/severity/reversibility/predictability blast-radius check before dispatch.
<!-- EMPOWER_ORCHESTRATOR:END -->
