# Changelog

All notable changes to noise.sh will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Voice-to-Text Dictation** - Push-to-talk transcription with whisper.cpp
  - `Ctrl+D` shortcut for start/stop recording
  - Automatic model download on first use (~142MB base English model)
  - Multiple model support: tiny (fast), base (balanced), small (accurate)
  - Real-time audio level feedback and recording duration display
  - Voice settings screen for model management and microphone testing
- **PWA Sync Infrastructure** - Local-first sync for companion mobile app
  - Embedded HTTP/WebSocket server on configurable port
  - Secure device pairing with 6-digit codes
  - Media storage for voice memos and photos
  - Idea inbox for reviewing captured ideas
  - Real-time sync status display
  - Database schema for `captured_ideas` and `paired_devices`
- Comprehensive stress detection in prosody engine with syllable splitting
- Harmony normalization with enharmonic and chord type standardization
- SQLite WAL mode and performance PRAGMAs for improved database performance
- Retry logic with exponential backoff for AI requests
- `UnwrapOrElse` and `UnwrapOrZero` methods for Result/Option types
- Security validation framework for plugins

### Changed
- Voice and sync services auto-initialize on startup (no manual setup required)
- Makefile targets for voice/sync are now optional (for advanced users only)
- Updated all documentation to reflect Go 1.25+ requirement
- Updated PROJECT_SNAPSHOT.md to reflect Active Development status
- Improved menu tests with proper dimension setup and model capturing

### Fixed
- Menu tests now properly capture returned model from Update()
- Suppressed noisy schema migration warnings during test runs

## [0.1.0] - 2026-01-22

### Added
- Initial project structure following Go best practices
- Bubble Tea TUI framework with Elm Architecture pattern
- Lipgloss styling system with gradient support
- Multiple theme support (Midnight Jazz, Sunset Glow, etc.)
- QuickIdeaAgent for AI-assisted songwriting
  - Unstick mode for continuation suggestions
  - Spark mode for inspiration
  - Tweak mode for variations
  - Check mode for quality feedback
  - Harmony mode for chord progressions
- Context-aware prompts based on content type detection
- Knowledge base integration with stub provider
- SQLite database layer with repository pattern
- Domain models (Song, Section, Line, Version, Project)
- Music theory library
  - ProsodyEngine for syllable counting and meter analysis
  - HarmonyLib for chord validation and scale suggestions
  - Rhyme analysis tools
- Comprehensive test suite
  - Unit tests for core components
  - Integration tests for workflows
  - E2E testing framework
  - Load testing infrastructure
- GitHub Actions CI/CD pipeline
  - Automated testing across platforms
  - Security scanning with gosec and govulncheck
  - Dependency management with Dependabot
- Project documentation
  - Product Requirements Document (PRD)
  - Technical Design Document (TDD)
  - Development Roadmap
  - Design System Reference
  - Architecture documentation

### Technical Stack
- Go 1.25.3
- Bubble Tea v1.3.10 (TUI framework)
- Lipgloss v1.1.1 (styling)
- modernc.org/sqlite v1.27.0 (pure Go SQLite)
- Ollama (local LLM integration)

---

## Version History Summary

| Version | Date | Highlights |
|---------|------|------------|
| 0.1.0 | 2026-01-22 | Initial release with core TUI, AI integration, and music theory tools |

[Unreleased]: https://github.com/Kyanite/noise/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/Kyanite/noise/releases/tag/v0.1.0
