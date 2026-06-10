# Changelog

All notable changes to Focus.sh will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- AI-powered task categorization and prioritization
- Local LLM support via Ollama
- OpenRouter API fallback support
- Calendar integration with Google Calendar
- Journal system with AI-powered insights
- Multiple dashboard modes (unified, main, minimal)
- Theme system with cyberpunk and synthwave themes
- Configuration wizard for initial setup
- TUI-based interactive mode
- MCP (Model Context Protocol) server support
- Audio feedback for task completion
- Task completion tracking and analytics
- Streaming response support for AI interactions

### Changed
- Migrated from internal TODO implementation to external focus.sh CLI
- Improved error handling and user feedback
- Enhanced configuration management with Viper
- Better cross-platform support (Windows, macOS, Linux)

### Fixed
- Issues with AI model selection and fallback
- Configuration file location detection
- Command-line argument parsing
- Cross-platform path handling

## [1.0.0] - 2024-01-XX

### Added
- Initial release of Focus.sh
- Basic task management (add, list, complete, delete)
- Simple AI integration
- Command-line interface
- Basic configuration system
- Installation scripts for Windows and Unix

## [0.9.0] - 2023-12-XX

### Added
- Development and testing infrastructure
- Basic project structure
- Initial AI integration experiments
- Task management core functionality

---

## Version History

### v1.0.x Series
- **Focus**: Production-ready AI-powered task management
- **Features**: Local LLM support, calendar integration, themes, TUI
- **Stability**: Ready for daily use

### v0.9.x Series
- **Focus**: Beta testing and feature development
- **Features**: Basic task management and AI integration
- **Stability**: Experimental, for testing only

---

## Migration Guide

### From v0.9.x to v1.0.x
- Configuration format updated (automatic migration supported)
- Some command names changed:
  - `focus done` → `focus complete`
  - `focus remove` → `focus delete`
- New AI configuration options added

---

## Release Schedule

- **Major releases**: Significant new features, breaking changes
- **Minor releases**: New features, improvements
- **Patch releases**: Bug fixes, security updates

---

## Support

For questions about upgrades or issues with specific versions:
- Check [GitHub Issues](https://github.com/kyanite/focus/issues)
- Review [Documentation](https://github.com/kyanite/focus/docs)
- Join community discussions

---

*Note: This changelog covers changes from the initial focus.sh implementation. Previous TODO Manager history is documented separately.*