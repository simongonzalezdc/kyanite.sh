# Contributing to Prism.sh

Thank you for your interest in contributing to Prism.sh! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Documentation](#documentation)

## Code of Conduct

This project follows a code of conduct that all contributors are expected to adhere to:

- Be respectful and inclusive
- Welcome newcomers
- Focus on constructive feedback
- Prioritize the community's best interest

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- A terminal emulator with true color support (for testing)

### Finding Issues to Work On

- Check the [Issues](https://github.com/kyanite/prism/issues) page
- Look for issues labeled `good first issue` or `help wanted`
- Comment on an issue to let others know you're working on it

## Development Setup

1. **Fork the repository**

```bash
# Fork on GitHub, then clone your fork
git clone https://github.com/YOUR_USERNAME/prism.git
cd prism
```

2. **Install dependencies**

```bash
go mod download
```

3. **Build and run**

```bash
make build
./bin/prism
```

4. **Run tests**

```bash
make test
```

## Making Changes

### Branch Naming

Use descriptive branch names:

- `feature/color-picker` - New features
- `fix/wcag-calculation` - Bug fixes
- `docs/update-readme` - Documentation updates
- `refactor/palette-generator` - Code refactoring

### Development Workflow

1. Create a new branch from `main`:

```bash
git checkout -b feature/your-feature-name
```

2. Make your changes

3. Write or update tests

4. Run tests and ensure they pass:

```bash
make test
```

5. Format your code:

```bash
make fmt
```

6. Commit your changes:

```bash
git add .
git commit -m "feat: add color picker feature"
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./internal/color/
```

### Writing Tests

- Place tests in the `tests/` directory
- Name test files `*_test.go`
- Use table-driven tests for multiple cases
- Aim for >70% code coverage

**Example Test:**

```go
func TestColorConversion(t *testing.T) {
    tests := []struct {
        name string
        hex  string
        want color.RGB
    }{
        {"Red", "#FF0000", color.RGB{255, 0, 0}},
        {"Green", "#00FF00", color.RGB{0, 255, 0}},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, _ := color.ParseHex(tt.hex)
            if got.RGB != tt.want {
                t.Errorf("got %v, want %v", got.RGB, tt.want)
            }
        })
    }
}
```

### TUI Testing

For Bubble Tea components:

```go
func TestMenuNavigation(t *testing.T) {
    m := ui.NewMenuModel(theme.NewManager())

    // Simulate key press
    msg := tea.KeyMsg{Type: tea.KeyDown}
    newModel, _ := m.Update(msg)

    // Verify state changed
    if newModel.selected != 1 {
        t.Errorf("Expected selected=1, got %d", newModel.selected)
    }
}
```

## Submitting Changes

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**

```
feat(color): add RGB to CMYK conversion

fix(wcag): correct contrast calculation for edge cases

docs(readme): update installation instructions
```

### Pull Request Process

1. **Update your branch** with latest `main`:

```bash
git fetch upstream
git rebase upstream/main
```

2. **Push to your fork**:

```bash
git push origin feature/your-feature-name
```

3. **Create a Pull Request** on GitHub

4. **Fill out the PR template** with:
   - Description of changes
   - Related issues
   - Screenshots (if UI changes)
   - Testing performed

5. **Wait for review**

   - Address review comments
   - Update PR with changes
   - Request re-review

6. **Merge** (done by maintainers after approval)

### PR Checklist

Before submitting, ensure:

- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Documentation updated (if needed)
- [ ] No breaking changes (or documented if unavoidable)
- [ ] Commit messages follow conventional commits
- [ ] Branch is up-to-date with main

## Coding Standards

### Go Style

Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines:

- Use `gofmt` for formatting
- Keep functions small and focused
- Use descriptive variable names
- Add comments for exported functions
- Handle errors explicitly

### Project Conventions

1. **Package Organization**
   - Keep packages focused on single responsibility
   - Use internal/ for private packages
   - Export only what's necessary

2. **Error Handling**
   ```go
   // ✅ Good
   if err != nil {
       return fmt.Errorf("failed to load config: %w", err)
   }

   // ❌ Bad
   if err != nil {
       panic(err)
   }
   ```

3. **Bubble Tea Models**
   - Never mutate state outside Update()
   - Use messages for communication
   - Keep View() pure (no side effects)

4. **Naming**
   - Use camelCase for unexported names
   - Use PascalCase for exported names
   - Be descriptive: `calculateContrast` not `calc`

### Code Review Focus

Reviewers will check for:

- **Functionality:** Does it work as intended?
- **Tests:** Are there adequate tests?
- **Performance:** Are there any bottlenecks?
- **Security:** Are there security concerns?
- **Style:** Does it follow project conventions?
- **Documentation:** Is it well-documented?

## Documentation

### Code Documentation

- Add godoc comments to all exported functions:

```go
// ParseHex converts a hex color string to a Color struct.
// The hex string can be in the format "#RRGGBB" or "RRGGBB".
// Returns an error if the hex string is invalid.
func ParseHex(hex string) (Color, error) {
    // ...
}
```

### README Updates

Update README.md when:
- Adding new features
- Changing installation process
- Updating dependencies
- Changing CLI interface

### Architecture Documentation

Update ARCHITECTURE.md when:
- Adding new modules
- Changing data flow
- Modifying core algorithms

## Release Process

(For maintainers)

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create git tag: `git tag v1.1.0`
4. Push tag: `git push origin v1.1.0`
5. GitHub Actions will create release

## Getting Help

- **Questions:** Open a [Discussion](https://github.com/kyanite/prism/discussions)
- **Bugs:** Open an [Issue](https://github.com/kyanite/prism/issues)
- **Security:** See [SECURITY.md](SECURITY.md)

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in commit history

Thank you for contributing to Prism.sh! 🎨

---

**Last Updated:** November 2025

<!-- EMPOWER_ORCHESTRATOR:START -->
## Agent-law contribution rule

This repository follows the Empower Orchestrator law in `docs/agent-law/empower-orchestrator.md`.

If a change exposes a repeated task or repeated agent failure, contributors and agents should either ship the smallest durable prevention artifact or explain why this PR is intentionally one-off.

Automation and durable system changes require the scale/severity/reversibility/predictability blast-radius check before dispatch.
<!-- EMPOWER_ORCHESTRATOR:END -->
