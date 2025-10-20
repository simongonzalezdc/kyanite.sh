# 🌌 NEON Focus - AI-Powered Task Manager

An intelligent CLI task management application with natural language processing capabilities, built with Charm libraries and cyberpunk aesthetics.

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
neon config list

# Set theme
neon config set theme synthwave

# Enhanced configuration wizard
neon enhanced-config
```

## AI Configuration

The application uses Ollama by default for local AI processing with fallback options:

1. **Primary**: Ollama (local) with `qwen2.5:1.5b` model
2. **Fallback**: OpenRouter (remote) when Ollama is unavailable
3. **Optional**: OpenAI (requires API key)

Configure with:
```bash
# Enhanced configuration wizard
neon enhanced-config

# Or manual configuration
neon config set ai.provider ollama
neon config set ai.model qwen2.5:1.5b
```

## Data Storage

Tasks are stored in `~/.focus/tasks.json` on Unix systems or `%USERPROFILE%\.focus\tasks.json` on Windows.

## Development

To run directly without building:
```bash
go run cmd/neon/main.go [command]
```

To run tests:
```bash
go test ./...
```

## Documentation

- [CALENDAR_PLAN.md](CALENDAR_PLAN.md) - Calendar implementation status
- [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md) - Code organization  
- [NEON_DEVELOPMENT_PLAN.md](NEON_DEVELOPMENT_PLAN.md) - Development progress
- [ROADMAP.md](ROADMAP.md) - Future development plans
- [AI_TODO_AGENT_PRD.md](AI_TODO_AGENT_PRD.md) - Product requirements
- [AI_TODO_AGENT_TDD.md](AI_TODO_AGENT_TDD.md) - Technical design

## 🎯 Commands Overview

### **Core Task Management**
- `neon add` - Add tasks with AI parsing
- `neon list` - View tasks with beautiful formatting
- `neon done` - Mark tasks as complete
- `neon remove` - Delete tasks

### **AI-Powered Features**
- `neon chat` - AI assistance for tasks and usage
- `neon inspire` - AI-powered task suggestions
- `neon wizard` - Advanced task creation wizard

### **Calendar Features**
- `neon calendar show [month|week|day]` - Calendar views
- `neon calendar today` - Today's tasks
- `neon calendar add [task] [date]` - Add task with date
- `neon calendar list` - List tasks by date
- `neon calendar navigate [date]` - Navigate to date

### **Interactive Features**
- `neon unified` - Complete dashboard (recommended)
- `neon dashboard` - TUI dashboard with AI
- `neon interactive` - Gum-powered task creation
- `neon filter` - Interactive task filtering
- `neon config` - Gum-based configuration

### **Configuration**
- `neon enhanced-config` - Advanced configuration wizard
- `neon config list/get/set/reset/path` - Viper configuration
- `neon theme` - Visual theme switching

### **Advanced Wizards**
- `neon wizard` - Task creation with Huh forms
- `neon config-wizard` - Configuration with Huh forms
- `neon edit-wizard` - Edit tasks with advanced wizard

## 🌈 Themes

NEON Focus supports multiple visual themes:
- **synthwave** (default) - Cyberpunk synthwave aesthetics
- **light** - Clean light theme
- **plain** - Simple monochrome theme

```bash
# Change theme
neon theme synthwave
neon theme light
neon theme plain
```

## 🚀 Performance

- **Fast AI Parsing**: Local Ollama processing (~26,500 ops/sec validation)
- **Efficient Rendering**: Optimized TUI with Bubble Tea
- **Smart Caching**: AI response caching for speed
- **Cross-Platform**: Native performance on all platforms