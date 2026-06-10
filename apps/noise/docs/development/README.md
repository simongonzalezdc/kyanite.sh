# Development Documentation

Documentation for contributors and developers working on noise.sh.

## Contents

| Document | Description |
|----------|-------------|
| [known-issues.md](known-issues.md) | Known test limitations and platform-specific issues |

## Getting Started

See the main [CONTRIBUTING.md](../../CONTRIBUTING.md) for:
- Development setup
- Code style guidelines
- Pull request process
- Release workflow

## Running Tests

```bash
# All tests
go test ./...

# Short tests (skip integration)
go test -short ./...

# Specific package
go test ./internal/infra/voice/...

# With coverage
go test -cover ./...
```

## Building

```bash
# Standard build
make build

# Build with voice support (requires CGO)
make build-native-voice

# Run
./bin/noise
```

## Project Structure

```
noise.sh/
├── cmd/noise/           # Application entry point
├── internal/
│   ├── app/             # Application services
│   ├── config/          # Configuration
│   ├── domain/          # Domain models
│   ├── infra/           # Infrastructure (db, voice, sync)
│   ├── ui/              # TUI components
│   └── theme/           # Theme system
├── test/                # Test suites
└── docs/                # Documentation
```

## Key Packages

| Package | Purpose |
|---------|---------|
| `internal/ui` | Bubble Tea TUI models |
| `internal/app` | Business logic services |
| `internal/infra/voice` | Voice capture and whisper.cpp |
| `internal/infra/sync` | PWA sync server |
| `internal/infra/db` | SQLite persistence |
| `internal/theme` | Theme management |

## Feature Flags

Some features can be toggled via configuration:

```yaml
features:
  enable_collaboration: false  # Collaboration (future)
  
voice:
  enabled: true               # Voice-to-text

sync:
  enabled: true               # PWA sync
```
