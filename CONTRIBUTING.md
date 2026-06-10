# Contributing to kyanite.sh

## Prerequisites

- **Go 1.25.5+** is required (workspace floor version)

## Development

```bash
# Build all apps
make build

# Test a specific app
cd apps/focus && go test ./...

# Workspace-wide sync
go work sync
```

## Structure

Each app is an independent Go module under `apps/`. The workspace (`go.work`) composes them for unified builds and releases.

## Release

Tag-based releases via goreleaser. Push a `v*` tag to trigger the release pipeline.
