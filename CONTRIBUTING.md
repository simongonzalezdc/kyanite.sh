# Contributing to noise.sh

Thanks for your interest in contributing to **noise.sh**! This document explains how to set up your environment, run the test suite, and ship a release using the automated workflow introduced in this iteration.

## Getting started

- **Language & tooling**: noise.sh is written in Go (see `go.mod`).
- **Requirements**: Go 1.25+, GNU Make, and the GitHub CLI (for local release simulation).
- **Repository**: clone via SSH or HTTPS and create feature branches for changes.

```bash
git clone git@github.com:Kyanite/noise.sh.git
cd noise.sh
```

Install tooling:

```bash
go install github.com/onsi/ginkgo/v2/ginkgo@latest # optional test helper
```

## Development workflow

1. **Install dependencies**

   ```bash
   go mod tidy
   ```

2. **Lint & test**

   ```bash
   make lint
   go test ./...
   ```

3. **Run the TUI locally**

   ```bash
   make build
   ./bin/noise
   ```

4. **Submit changes**

   - Format with `gofmt` (or `goimports`).
   - Commit using Conventional Commit messages (e.g., `feat: add new pattern parser`).
   - Open a pull request targeting the default branch.

## Automated release process

Releases are fully automated through [`release.yml`](.github/workflows/release.yml) and the Go tooling under `tools/release`. Tags that match the semantic versioning pattern (`vX.Y.Z`) trigger the workflow.

### Versioning

- We follow [Semantic Versioning](https://semver.org/).
- The release tooling uses [`github.com/Masterminds/semver/v3`](go.mod) to validate new tags.
- Application version metadata is injected at build time via Go linker flags.

### Release workflow (CI)

When a tag like `v1.2.3` is pushed:

1. **Prepare job**
   - Validates that the tag is semver-compliant.
   - Determines the previous tag to compute diff ranges.

2. **Quality job**
   - Runs linting (`make lint`) and full unit tests (`go test ./...`).

3. **Package job**
   - Cross-compiles `noise` for Linux, macOS, and Windows (amd64/arm64) with `CGO_ENABLED=0`.
   - Produces tarballs (`*.tar.gz`) or zips (`*.zip`).

4. **Checksums job**
   - Generates a SHA-256 manifest (`checksums.txt`) for all packaged artifacts.

5. **Validate job**
   - Executes release-specific tests (`go test ./test -run Release`), which ensure version data is embedded correctly and tags are valid.

6. **Release job**
   - Produces release notes from commits using `tools/release` (`notes` command).
   - Updates `CHANGELOG.md` with the new entry and pushes it back to the default branch if it changed.
   - Verifies checksums.
  - Publishes a GitHub Release with packaged binaries and the checksum manifest.

Artifacts such as `build/release-notes.md` and `changeLOG.md` are uploaded for traceability.

### Manual steps before tagging

1. Ensure all tests pass and the default branch is green.
2. Update documentation or outstanding issues.
3. Decide the new semantic version:
   - `MAJOR` for incompatible changes.
   - `MINOR` for new features (backward-compatible).
   - `PATCH` for bug fixes.

### Tagging and releasing

```bash
VERSION=vX.Y.Z
git checkout main
git pull
git tag "${VERSION}"
git push origin "${VERSION}"
```

The workflow handles the rest: packaging, checksums, release notes, changelog, and publishing.

### Local release simulation (optional)

You can dry-run parts of the process locally:

```bash
# Generate release notes between tags
VERSION_TAG=v1.2.3
FROM_TAG=v1.2.2
make release-notes VERSION_TAG=${VERSION_TAG} FROM=${FROM_TAG} OUTPUT=build/release-notes.md
make release-changelog VERSION_TAG=${VERSION_TAG} NOTES=build/release-notes.md
make release-checksums CHECKSUM_OUTPUT=build/checksums.txt
make verify-checksums CHECKSUM_FILE=build/checksums.txt
```

These commands call the same Go tooling used by CI for semantic version validation, changelog updates, and checksum verification.

## Questions & support

If you encounter issues:

- **File a GitHub issue** with reproduction steps.
- **Join the discussions** tab for feature ideas.
- **Mention maintainers** if you need a release expedited.

Thanks for helping keep noise.sh stable and delightful!