# kyanite.sh

A unified suite of terminal-first creative tools built with Go and the Charm stack.

## Apps

| App | Description |
|---|---|
| **focus** | ADHD-friendly terminal focus & task manager |
| **noise** | Local-first AI-assisted songwriting studio |
| **syntax** | Terminal fiction-writing & worldbuilding TUI |
| **prism** | Terminal color-palette design with WCAG contrast validation |

## Prerequisites

- Go 1.25.5+

## Quick Start

```bash
# Build all apps
make build

# Run a specific app
go run ./apps/focus/cmd/focus
go run ./apps/noise/cmd/noise
go run ./apps/syntax/cmd/syntax
go run ./apps/prism/cmd/prism

# Run all tests
make test

# Lint
make lint
```

## Project Structure

```
kyanite.sh/
  go.work              # Go workspace definition
  apps/
    focus/             # github.com/kyanite/focus
    noise/             # github.com/kyanite/noise
    syntax/            # github.com/kyanite/syntax
    prism/             # github.com/kyanite/prism
  tools/               # Shared release and lint tooling
  .goreleaser.yaml     # Multi-binary release config
```

## Development

Each app is an independent Go module with its own `go.mod`. The workspace (`go.work`) composes them for unified builds and releases.

```bash
# Work on a single app
cd apps/focus
go test ./...

# Workspace-wide operations
go work sync
```
