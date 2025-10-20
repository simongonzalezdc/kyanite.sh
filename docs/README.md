# Focus.sh Documentation

Welcome to the Focus.sh documentation. This AI-powered CLI task manager helps you stay productive with natural language task management and beautiful terminal interfaces.

## 📚 Documentation Categories

### [User Guide](user/)
- [Installation](user/installation.md) - Get Focus.sh running on your system
- [Configuration](user/configuration.md) - Set up AI providers and preferences
- [Commands](user/commands.md) - Complete command reference
- [Themes](user/themes.md) - Visual customization options
- [Troubleshooting](user/troubleshooting.md) - Common issues and solutions
- [FAQ](user/faq.md) - Frequently asked questions

### [Developer Guide](developer/)
- [Contributing](../CONTRIBUTING.md) - How to contribute to Focus.sh
- [Architecture](developer/architecture.md) - System design and components
- [API Reference](developer/api.md) - Internal APIs and interfaces
- [Testing](developer/testing.md) - Test suite and guidelines
- [Development Setup](developer/setup.md) - Local development environment

### [Project Information](project/)
- [Changelog](../CHANGELOG.md) - Version history and changes
- [Roadmap](project/roadmap.md) - Future development plans
- [Security](../SECURITY.md) - Security policy and reporting
- [License](../LICENSE) - MIT License details

### [Internal Documentation](internal/)
- Historical development plans and status reports
- Technical analysis documents
- Implementation notes and research

---

## 🚀 Quick Start

1. **Install dependencies**: [Go 1.21+](https://golang.org/dl/) and [Ollama](https://ollama.ai/)
2. **Pull AI model**: `ollama pull qwen2.5:1.5b`
3. **Run Focus.sh**: `go run cmd/focus/main.go unified`

## 🤖 AI Setup

Focus.sh uses AI for natural language task processing:

### Primary: Local AI (Ollama)
- **Model**: `qwen2.5:1.5b` (recommended)
- **Privacy**: All processing happens locally
- **Setup**: `ollama pull qwen2.5:1.5b`

### Fallback: Remote AI (OpenRouter)
- **Usage**: When local AI is unavailable
- **Models**: Various models available
- **Setup**: Requires API key

See [Configuration Guide](user/configuration.md) for detailed setup.

## 🎯 Key Features

- **Natural Language Tasks**: Add tasks with plain English
- **AI-Powered Insights**: Smart categorization and prioritization
- **Beautiful Terminal UI**: Synthwave themes and modern design
- **Calendar Integration**: Task scheduling and calendar views
- **Cross-Platform**: Windows, macOS, and Linux support
- **Interactive Dashboards**: Multiple TUI interfaces
- **Local First**: Data stored locally with optional cloud features

---

## 📞 Support

- 🐛 **Report Issues**: [GitHub Issues](https://github.com/kyanite/focus/issues)
- 💬 **Community**: [GitHub Discussions](https://github.com/kyanite/focus/discussions)
- 🔒 **Security**: See [Security Policy](../SECURITY.md)
- 📖 **Documentation**: Browse the sections above

---

*Last updated: 2024-01-XX*