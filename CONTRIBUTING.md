# Contributing to Focus.sh

Thank you for your interest in contributing to Focus.sh! This guide will help you get started with contributing to this AI-powered CLI productivity tool.

## Table of Contents

- [Development Setup](#development-setup)
- [Code Style Guidelines](#code-style-guidelines)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Bug Reports](#bug-reports)
- [Feature Requests](#feature-requests)

## Development Setup

### Prerequisites

- Go 1.21.0 or higher
- Git
- Ollama (for local AI testing) or OpenRouter API key
- Make (optional, for convenience scripts)

### Setup Steps

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/your-username/focus.git
   cd focus
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up AI configuration**
   
   Create a `.env` file in the root directory:
   ```env
   # Local AI (Ollama) - Primary
   OLLAMA_BASE_URL=http://localhost:11434
   DEFAULT_MODEL=qwen2.5:1.5b
   
   # Fallback (OpenRouter) - Emergency only
   OPENROUTER_API_KEY=your_openrouter_key
   FALLBACK_MODEL=mistralai/mistral-7b-instruct
   ```

4. **Install Ollama and pull the default model**
   ```bash
   # Install Ollama (https://ollama.ai/)
   ollama pull qwen2.5:1.5b
   ```

5. **Build and test**
   ```bash
   go build ./cmd/focus
   go test ./...
   ```

## Code Style Guidelines

### Go Conventions

- Follow standard Go formatting (`go fmt`)
- Use `golangci-lint` for linting (configuration in `.golangci.yml`)
- Prefer early returns and explicit error handling
- Use descriptive names with PascalCase for exported types
- Group imports: stdlib → third-party → internal

### Naming Conventions

- **Project name**: Use "focus.sh" or "Focus" consistently
- **Module**: `github.com/kyanite/focus`
- **Commands**: Use action-oriented names (`add`, `complete`, `delete`)
- **Files**: Use snake_case for non-Go files, PascalCase for Go files

### Testing

- Write tests for all new functionality
- Use table-driven tests with struct slices
- Test naming: `Test{TypeName}_{MethodName}`
- Use `t.TempDir()` for test files
- Use `if testing.Short() { t.Skip() }` for integration tests

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbosity
go test ./internal/... ./pkg/... -v

# Run specific test
go test -run TestAIManager_ParseTask ./internal/ai

# Run integration tests
go test ./test/...

# Run tests with coverage
go test -cover ./...
```

### Test Structure

- **Unit tests**: Co-located with source files (`*_test.go`)
- **Integration tests**: In `test/` directory
- **Manual test files**: In `test/` directory for interactive testing

### Essential Test Commands

Add these to your development workflow:

```bash
# Build and run tests
go build ./... && go test ./...

# Lint code
golangci-lint run

# Format code
go fmt ./...

# Full test suite
./test_deployment.sh  # Unix
test_core.bat         # Windows
```

## Submitting Changes

### Branch Strategy

1. **Create a feature branch** from `main`
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make focused, atomic commits**
   ```bash
   git add .
   git commit -m "feat: add new task prioritization feature"
   ```

3. **Keep your branch updated**
   ```bash
   git fetch origin
   git rebase origin/main
   ```

### Commit Message Format

Use conventional commits:

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `style:` Code style changes (no functional impact)
- `refactor:` Code refactoring
- `test:` Test additions/modifications
- `chore:` Maintenance tasks

### Pull Request Process

1. **Ensure all tests pass**
   ```bash
   go test ./...
   golangci-lint run
   ```

2. **Update documentation** if needed
3. **Create a pull request** with:
   - Clear title and description
   - Testing instructions
   - Any breaking changes highlighted

4. **Code review**: Address feedback promptly

## Bug Reports

### Report Format

When reporting bugs, include:

```markdown
## Bug Description
Brief description of the issue

## Steps to Reproduce
1. Run `focus add "test task"`
2. Run `focus list`
3. ...

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., Windows 11, macOS 14.0, Ubuntu 22.04]
- Go version: [e.g., 1.21.0]
- Focus.sh version: [e.g., v1.0.0]
- AI provider: [e.g., Ollama, OpenRouter]
- Model: [e.g., qwen2.5:1.5b]
```

## Feature Requests

### Request Format

```markdown
## Feature Description
Clear description of the proposed feature

## Problem Statement
What problem does this solve?

## Proposed Solution
How should it work?

## Alternatives Considered
Other approaches and why they weren't chosen

## Additional Context
Any relevant information or screenshots
```

## Development Areas

We welcome contributions in these areas:

### Core Features
- **AI Integration**: Local LLM support, fallback mechanisms
- **Task Management**: CRUD operations, prioritization, categorization
- **Calendar Integration**: Event management, reminders
- **TUI**: Bubble Tea interfaces, dashboard improvements

### CLI Features
- **Command Completion**: Better auto-completion and suggestions
- **Configuration**: Enhanced config management
- **Output Formatting**: Multiple output formats (JSON, table, etc.)

### Infrastructure
- **Testing**: Unit tests, integration tests, end-to-end tests
- **Documentation**: User guides, API docs, tutorials
- **CI/CD**: GitHub Actions improvements
- **Performance**: Optimization and benchmarking

## Getting Help

- **Documentation**: Check `docs/` directory and `CRUSH.md`
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions
- **Code Review**: Request reviews from maintainers

## Recognition

Contributors are recognized in:
- `README.md` contributors section
- Release notes for significant contributions
- Special thanks in major version releases

Thank you for contributing to Focus.sh! 🚀

<!-- EMPOWER_ORCHESTRATOR:START -->
## Agent-law contribution rule

This repository follows the Empower Orchestrator law in `docs/agent-law/empower-orchestrator.md`.

If a change exposes a repeated task or repeated agent failure, contributors and agents should either ship the smallest durable prevention artifact or explain why this PR is intentionally one-off.

Automation and durable system changes require the scale/severity/reversibility/predictability blast-radius check before dispatch.
<!-- EMPOWER_ORCHESTRATOR:END -->
