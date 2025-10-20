# 🌌 focus.sh - AI-Powered Task Manager

An intelligent CLI task management application within the Kyanite Suite, delivering natural language productivity with Charm libraries and modern professional design.

## Features

- 🤖 **AI-Powered Task Management**: Natural language input parsing with local AI (Ollama) and remote AI (OpenRouter)
- 🎨 **Beautiful Synthwave Interface**: Stunning terminal UI with cyberpunk aesthetics
- 📅 **Calendar Integration**: Full calendar system with month/week/day views
- 🧙‍♂️ **Interactive Wizards**: Advanced forms with Huh for task creation and configuration
- 🌿 **Gum Enhancement**: Interactive inputs and filtering with Gum
- ⚙️ **Viper Configuration**: Complete configuration management system
- 🎯 **Unified Dashboard**: Single interface accessing all features
- ✨ **Cross-Platform Support**: Works on Windows, macOS, and Linux

## Quick Start

1. Install [Go 1.21+](https://golang.org/dl/) and [Ollama](https://ollama.ai/)
2. Pull the required model: `ollama pull qwen2.5:1.5b`
3. Build with `build.bat` or run directly with `go run cmd/focus/main.go`

## Usage

```bash
# View all available commands
focus --help

# Add a task with natural language
focus add "Complete the synthwave project by Friday"

# View tasks in beautiful list format
focus list

# View tasks in calendar format
focus calendar today

# Add task with specific date
focus calendar add "Team meeting" 2025-10-20

# Launch unified dashboard (recommended)
focus unified

# Launch AI chat for assistance
focus chat
```

## Configuration

The application uses extensive configuration for AI providers, themes, and UI settings:

```bash
# View all configuration
focus config list

# Set theme
focus config set theme synthwave

# Enhanced configuration wizard
focus enhanced-config
```

## AI Configuration

The application uses Ollama by default for local AI processing with fallback options:

1. **Primary**: Ollama (local) with `qwen2.5:1.5b` model
2. **Fallback**: OpenRouter (remote) when Ollama is unavailable
3. **Optional**: OpenAI (requires API key)

Configure with:
```bash
# Enhanced configuration wizard
focus enhanced-config

# Or manual configuration
focus config set ai.provider ollama
focus config set ai.model qwen2.5:1.5b
```

## Data Storage

Tasks are stored in `~/.focus/tasks.json` on Unix systems or `%USERPROFILE%\.focus\tasks.json` on Windows.

## Development

To run directly without building:
```bash
go run cmd/focus/main.go [command]
```

To run tests:
```bash
go test ./...
```

## Documentation

📖 **[Complete Documentation](docs/)**

### Quick Links
- [Installation Guide](docs/user/installation.md) - Get Focus.sh running
- [Configuration Guide](docs/user/configuration.md) - Set up AI and preferences  
- [Commands Reference](docs/user/commands.md) - All commands and options
- [Troubleshooting](docs/user/troubleshooting.md) - Common issues and solutions
- [FAQ](docs/user/faq.md) - Frequently asked questions
- [Contributing Guide](CONTRIBUTING.md) - Development contribution guide
- [Roadmap](docs/project/roadmap.md) - Future development plans

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! 🎉 Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Support

- 📖 **Documentation**: [docs/](docs/)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/kyanite/focus/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/kyanite/focus/discussions)
- 🔒 **Security**: [Security Policy](SECURITY.md)

---

**Built with ❤️ using Go, Bubble Tea, and Charm libraries**

## 🎯 Commands Overview

### **Core Task Management**
- `focus add` - Add tasks with AI parsing
- `focus list` - View tasks with beautiful formatting
- `focus complete` - Mark tasks as complete
- `focus delete` - Delete tasks

### **AI-Powered Features**
- `focus chat` - AI assistance for tasks and usage
- `focus inspire` - AI-powered task suggestions
- `focus wizard` - Advanced task creation wizard

### **Calendar Features**
- `focus calendar show [month|week|day]` - Calendar views
- `focus calendar today` - Today's tasks
- `focus calendar add [task] [date]` - Add task with date
- `focus calendar list` - List tasks by date
- `focus calendar navigate [date]` - Navigate to date

### **Interactive Features**
- `focus unified` - Complete dashboard (recommended)
- `focus dashboard` - TUI dashboard with AI
- `focus interactive` - Gum-powered task creation
- `focus filter` - Interactive task filtering
- `focus config` - Gum-based configuration

### **Configuration**
- `focus enhanced-config` - Advanced configuration wizard
- `focus config list/get/set/reset/path` - Viper configuration
- `focus theme` - Visual theme switching

### **Advanced Wizards**
- `focus wizard` - Task creation with Huh forms
- `focus config-wizard` - Configuration with Huh forms
- `focus edit-wizard` - Edit tasks with advanced wizard

## 🌈 Themes

focus.sh supports multiple visual themes:
- **synthwave** (default) - Cyberpunk synthwave aesthetics
- **light** - Clean light theme
- **plain** - Simple monochrome theme

```bash
# Change theme
focus theme synthwave
focus theme light
focus theme plain
```

## 🚀 Performance

- **Fast AI Parsing**: Local Ollama processing (~26,500 ops/sec validation)
- **Efficient Rendering**: Optimized TUI with Bubble Tea
- **Smart Caching**: AI response caching for speed
- **Cross-Platform**: Native performance on all platforms