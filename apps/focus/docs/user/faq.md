# Frequently Asked Questions

Common questions and answers about Focus.sh, organized by topic.

## Getting Started

### **Q: What is Focus.sh?**
A: Focus.sh is an AI-powered CLI task manager that helps you manage tasks with natural language input. It uses local AI (Ollama) for privacy and has remote AI (OpenRouter) as fallback.

### **Q: What do I need to run Focus.sh?**
A: You need:
- Go 1.21+ 
- Ollama (recommended) or OpenRouter API key
- Terminal/command prompt
- 4GB+ RAM (8GB+ recommended for AI models)

### **Q: How do I install Focus.sh?**
A: The recommended way is to build from source:
```bash
git clone https://github.com/simongonzalezdc/focus.sh.git
cd focus
go build ./cmd/focus
ollama pull qwen2.5:1.5b
```

See the [Installation Guide](installation.md) for detailed instructions.

## AI and Models

### **Q: Does Focus.sh send my data to the cloud?**
A: By default, no. Focus.sh uses Ollama for local AI processing, so all your data stays on your machine. Only if you configure OpenRouter as fallback will data be sent to their servers.

### **Q: What AI models work best?**
A: Recommended models:
- `qwen2.5:1.5b` - Fast, capable, low resource usage
- `llama3.2:3b` - Good balance of speed and capability  
- `mistral:7b` - Most capable but slower

### **Q: How much RAM do AI models need?**
A: Approximately:
- `qwen2.5:1.5b`: 2-3GB RAM
- `llama3.2:3b`: 4-6GB RAM
- `mistral:7b`: 6-8GB RAM

### **Q: Can I use Focus.sh without Ollama?**
A: Yes. You can use OpenRouter as your AI provider:
```bash
export OPENROUTER_API_KEY="your_key"
focus config set ai.provider openrouter
```

### **Q: Why is my AI response slow?**
A: Try these solutions:
1. Use a smaller model: `focus config set ai.model qwen2.5:1.5b`
2. Enable caching: `focus config set performance.cache_enabled true`
3. Reduce timeout: `focus config set ai.timeout 15`

## Task Management

### **Q: How do I add a task with due date?**
A: Multiple ways:
```bash
# Explicit date
focus add "Meeting" --due 2024-01-25

# Natural language
focus add "Meeting next Friday at 2 PM"

# Interactive mode
focus interactive
```

### **Q: Can I organize tasks into categories?**
A: Yes:
```bash
# Specify category
focus add "Write report" --category work

# Default categories: work, personal, learning, health
# You can create custom categories too
```

### **Q: How do I complete or delete tasks?**
A: Use the task ID shown in `focus list`:
```bash
focus complete 5    # Mark task #5 as complete
focus delete 7      # Delete task #7
```

### **Q: Can I import/export tasks?**
A: Yes:
```bash
# Export to JSON
focus export --output my-tasks.json

# Export specific category
focus export --category work --output work-tasks.json

# Import tasks
focus import my-tasks.json
```

## Calendar Integration

### **Q: Does Focus.sh integrate with calendars?**
A: Yes, Focus.sh has calendar views and Google Calendar integration:
```bash
focus calendar today      # Today's tasks
focus calendar week       # Week view
focus calendar add "Meeting" 2024-01-25
```

### **Q: Is Google Calendar setup required?**
A: No. Basic calendar functionality works without Google Calendar. Google Calendar sync is optional for two-way synchronization.

### **Q: How do I sync with Google Calendar?**
A: Run the configuration wizard:
```bash
focus enhanced-config --section calendar
```
Follow the OAuth authentication flow.

## Themes and Display

### **Q: What themes are available?**
A: Built-in themes:
- `synthwave` - Cyberpunk aesthetics (default)
- `light` - Clean, minimal design
- `plain` - Monochrome, terminal-friendly

### **Q: How do I change themes?**
A: 
```bash
focus theme synthwave
# or
focus config set theme light
```

### **Q: Can I create custom themes?**
A: Yes. Create JSON files in `~/.config/focus/themes/`:
```json
{
  "name": "my-theme",
  "colors": {
    "primary": "#00ff00",
    "secondary": "#ff00ff"
  }
}
```

## Performance and Resources

### **Q: Is Focus.sh resource-intensive?**
A: Focus.sh is designed to be lightweight:
- Base application: ~50MB RAM
- AI models: 2-8GB RAM (depending on model)
- Storage: ~10MB for data and configuration

### **Q: How can I improve performance?**
A: 
1. Use the `qwen2.5:1.5b` model
2. Enable response caching
3. Use local storage (default)
4. Close unnecessary applications

### **Q: Does Focus.sh work offline?**
A: Yes, with Ollama configured for local AI. Only OpenRouter requires internet connectivity.

## Privacy and Security

### **Q: Where is my data stored?**
A: All data is stored locally:
- **Config**: `~/.config/focus/`
- **Data**: `~/.local/share/focus/`
- **Logs**: `~/.local/log/focus/`

### **Q: Is my task data private?**
A: Yes, when using Ollama (local AI). Your tasks never leave your machine. Only OpenRouter usage involves external processing.

### **Q: Can I encrypt my data?**
A: Currently, data is stored in plain text JSON. Encryption is planned for a future release.

## Troubleshooting

### **Q: Focus.sh says "Ollama not available"**
A: Solutions:
1. Install Ollama: https://ollama.ai/download
2. Start Ollama: `ollama serve`
3. Pull a model: `ollama pull qwen2.5:1.5b`
4. Or use OpenRouter: `focus config set ai.provider openrouter`

### **Q: Commands aren't working as expected**
A: Try these steps:
1. Check syntax: `focus <command> --help`
2. Update Focus.sh: `git pull && go build`
3. Reset config: `focus config reset`
4. Enable debug: `focus --debug <command>`

### **Q: The terminal output looks garbled**
A: Try:
1. Use plain theme: `focus theme plain`
2. Disable colors: `export NO_COLOR=1`
3. Reset terminal: `reset`

## Advanced Usage

### **Q: Can I use Focus.sh in scripts?**
A: Yes. Focus.sh supports scripting:
```bash
#!/bin/bash
# Add tasks from file
while read task; do
    focus add "$task"
done < tasks.txt

# List tasks in JSON format
focus list --format json > tasks.json
```

### **Q: Does Focus.sh have an API?**
A: Focus.sh includes an MCP (Model Context Protocol) server for integration with other tools. See `cmd/focus-mcp/` for details.

### **Q: Can I customize the AI prompts?**
A: Yes. Edit the AI configuration:
```bash
focus config set ai.system_prompt "You are a helpful assistant for task management."
```

## Comparison with Other Tools

### **Q: How is Focus.sh different from Todo.txt?**
A: Focus.sh offers:
- AI-powered natural language input
- Interactive terminal UI with themes
- Calendar integration
- Task prioritization and categorization
- Built-in dashboards and analytics

### **Q: vs. other task managers?**
A: Focus.sh advantages:
- Privacy-focused (local AI)
- Terminal-based (fast, keyboard-driven)
- No subscription required
- Open source and extensible
- Works offline

## Development and Contributing

### **Q: How can I contribute to Focus.sh?**
A: See the [Contributing Guide](../CONTRIBUTING.md) for:
- Development setup
- Code style guidelines
- Pull request process
- Issue reporting

### **Q: What technology stack does Focus.sh use?**
A: 
- **Language**: Go 1.21+
- **UI Framework**: Bubble Tea (TUI)
- **Styling**: Lip Gloss, Glow
- **Configuration**: Viper
- **AI**: Ollama/OpenRouter integration
- **Calendar**: Google Calendar API

### **Q: Can I build Focus.sh for different platforms?**
A: Yes. Use Go's cross-compilation:
```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o focus-linux ./cmd/focus

# macOS  
GOOS=darwin GOARCH=amd64 go build -o focus-macos ./cmd/focus

# Windows
GOOS=windows GOARCH=amd64 go build -o focus-windows.exe ./cmd/focus
```

## Installation and Setup Issues

### **Q: "go: command not found"**
A: Install Go from https://golang.org/dl/ or use a package manager:
```bash
# Ubuntu/Debian
sudo apt install golang-go

# macOS
brew install go

# Windows
winget install GoLang.Go
```

### **Q: "permission denied" on Unix**
A: 
```bash
chmod +x focus
# or
sudo chmod +x focus
```

### **Q: Windows says "not recognized as internal command"**
A: Use the full path or add to PATH:
```powershell
# Use full path
.\focus --help

# Add to current session
$env:PATH += ";$(pwd)"

# Add permanently (as Administrator)
[Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";C:\path\to\focus", "Machine")
```

## Still Have Questions?

### **Where can I get help?**
- **Documentation**: Browse the [docs/](../) directory
- **GitHub Issues**: [Report bugs](https://github.com/simongonzalezdc/focus.sh/issues)
- **GitHub Discussions**: [Ask questions](https://github.com/simongonzalezdc/focus.sh/discussions)
- **Troubleshooting**: [Troubleshooting Guide](troubleshooting.md)

### **How do I report a bug?**
1. Check existing issues first
2. Use debug mode: `focus --debug <command>`
3. Include system info: `focus status --all`
4. Provide reproduction steps
5. Open an issue with detailed information

---

*Is your question not answered here? Check the full documentation or ask in GitHub Discussions.*
