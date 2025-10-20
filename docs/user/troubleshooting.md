# Troubleshooting Guide

This guide covers common issues, error messages, and solutions for Focus.sh.

## Quick Diagnosis

Start with these commands to diagnose issues:

```bash
# Check system status
focus status

# Verify AI connectivity
focus status --ai

# Check configuration
focus config list

# Enable debug logging
export FOCUS_DEBUG=true
focus --debug list
```

## Installation Issues

### "command not found: focus"

**Problem**: Focus.sh not in PATH or not installed

**Solutions**:
```bash
# Use full path
./path/to/focus --help

# Add to PATH (temporary)
export PATH=$PATH:/path/to/focus

# Add to PATH (permanent - add to ~/.bashrc or ~/.zshrc)
echo 'export PATH=$PATH:/path/to/focus' >> ~/.bashrc

# Reinstall with proper setup
./build.sh  # Unix
build.bat   # Windows
```

### "Go not found"

**Problem**: Go is not installed or not in PATH

**Solutions**:
```bash
# Check Go installation
go version

# Install Go
# Ubuntu/Debian
sudo apt install golang-go

# macOS
brew install go

# Windows
winget install GoLang.Go
```

### "permission denied"

**Problem**: Insufficient permissions to run Focus.sh

**Solutions**:
```bash
# Make executable (Unix/macOS)
chmod +x focus

# Run with appropriate permissions
sudo ./focus --help  # Only if necessary

# Check file ownership
ls -la focus
```

## AI Provider Issues

### "Ollama not available"

**Problem**: Ollama is not running or not installed

**Diagnosis**:
```bash
# Check if Ollama is running
curl http://localhost:11434/api/tags

# Check Ollama service status
ollama list
```

**Solutions**:
```bash
# Start Ollama service
ollama serve &

# Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh  # Linux
brew install ollama                           # macOS

# Download and install from https://ollama.ai/download  # Windows
```

### "Model not found"

**Problem**: Required AI model is not available

**Diagnosis**:
```bash
# List available models
ollama list
```

**Solutions**:
```bash
# Pull recommended model
ollama pull qwen2.5:1.5b

# Or alternative models
ollama pull llama3.2:3b
ollama pull mistral:7b

# Update configuration
focus config set ai.model qwen2.5:1.5b
```

### "OpenRouter API key invalid"

**Problem**: OpenRouter API authentication failing

**Diagnosis**:
```bash
# Check environment variable
echo $OPENROUTER_API_KEY

# Test API connectivity
curl -H "Authorization: Bearer $OPENROUTER_API_KEY" \
     https://openrouter.ai/api/v1/models
```

**Solutions**:
```bash
# Set correct API key
export OPENROUTER_API_KEY="your_correct_api_key"

# Add to .env file
echo "OPENROUTER_API_KEY=your_correct_api_key" >> .env

# Get new API key from https://openrouter.ai/keys
```

## Configuration Issues

### "Configuration file not found"

**Problem**: Focus.sh cannot find or create configuration

**Diagnosis**:
```bash
# Show config file location
focus config path

# Check directory exists
ls -la ~/.config/focus/
```

**Solutions**:
```bash
# Run configuration wizard
focus enhanced-config

# Create directories manually
mkdir -p ~/.config/focus
mkdir -p ~/.local/share/focus
mkdir -p ~/.local/log/focus

# Reset configuration
focus config reset
```

### "Invalid configuration value"

**Problem**: Configuration contains invalid values

**Diagnosis**:
```bash
# Check current configuration
focus config list

# Look for validation errors in debug output
focus --debug list
```

**Solutions**:
```bash
# Reset specific setting
focus config reset ai.provider

# Set valid value
focus config set ai.provider ollama
focus config set theme synthwave

# Reset entire configuration
focus config reset
```

## Task Management Issues

### "Task not found"

**Problem**: Cannot find task with specified ID or pattern

**Diagnosis**:
```bash
# List all tasks
focus list

# Check task IDs
focus list --format json | jq '.tasks[] | {id, description}'
```

**Solutions**:
```bash
# Use correct task ID
focus complete 5

# Search by pattern
focus complete --pattern "*groceries*"

# List and choose interactively
focus filter
```

### "Cannot parse task description"

**Problem**: AI cannot understand task input

**Diagnosis**:
```bash
# Test AI directly
focus chat "Parse this task: buy milk"

# Check AI status
focus status --ai
```

**Solutions**:
```bash
# Use simpler description
focus add "Buy milk"

# Use structured input
focus add --category shopping "Buy milk"

# Disable AI parsing temporarily
focus add --no-ai "Buy milk"
```

## Calendar Issues

### "Google Calendar authentication failed"

**Problem**: Cannot authenticate with Google Calendar

**Diagnosis**:
```bash
# Check calendar configuration
focus config get calendar

# Test calendar status
focus status --calendar
```

**Solutions**:
```bash
# Re-authenticate
focus enhanced-config --section calendar

# Check credentials file
ls -la ~/.config/focus/google-credentials.json

# Reset calendar settings
focus config reset calendar
```

## Performance Issues

### "Slow AI response"

**Problem**: AI processing is taking too long

**Diagnosis**:
```bash
# Check AI provider performance
focus status --ai

# Test model directly
ollama run qwen2.5:1.5b "test response"
```

**Solutions**:
```bash
# Use faster model
focus config set ai.model qwen2.5:1.5b

# Reduce timeout
focus config set ai.timeout 15

# Enable caching
focus config set performance.cache_enabled true

# Use fallback provider
focus config set ai.provider openrouter
```

### "High memory usage"

**Problem**: Focus.sh is using excessive memory

**Solutions**:
```bash
# Clear cache
focus config set performance.cache_enabled false

# Reduce context window
focus config set ai.context_window 5

# Use smaller model
focus config set ai.model qwen2.5:1.5b

# Restart application
```

## Theme and Display Issues

### "Theme not found"

**Problem**: Specified theme doesn't exist

**Diagnosis**:
```bash
# List available themes
focus theme list

# Check current theme
focus config get theme
```

**Solutions**:
```bash
# Use default theme
focus config set theme synthwave

# Reset theme
focus theme synthwave

# Check theme files
ls -la ~/.config/focus/themes/
```

### "Display formatting issues"

**Problem**: Terminal output is garbled or unreadable

**Solutions**:
```bash
# Reset terminal
reset

# Use plain theme
focus theme plain

# Disable colors
export NO_COLOR=1
focus list

# Check terminal compatibility
echo $TERM
```

## Network Issues

### "Cannot connect to AI provider"

**Problem**: Network connectivity problems

**Diagnosis**:
```bash
# Test internet connectivity
curl -I https://httpbin.org/status/200

# Test Ollama connectivity
curl http://localhost:11434/api/tags

# Test OpenRouter connectivity
curl -I https://openrouter.ai/api/v1/models
```

**Solutions**:
```bash
# Check firewall settings
# Allow localhost:11434 for Ollama
# Allow https://openrouter.ai for remote AI

# Use local AI only
focus config set ai.provider ollama
focus config set ai.fallback_enabled false

# Configure proxy if needed
export HTTPS_PROXY=http://proxy.company.com:8080
```

## File System Issues

### "Cannot write to data directory"

**Problem**: Permission or disk space issues

**Diagnosis**:
```bash
# Check directory permissions
ls -la ~/.local/share/focus/

# Check disk space
df -h

# Test write permissions
touch ~/.local/share/focus/test.txt
```

**Solutions**:
```bash
# Fix permissions
chmod 755 ~/.local/share/focus/
chmod 644 ~/.local/share/focus/tasks.json

# Free up disk space
# Clean up old logs
rm ~/.local/log/focus/*.log

# Use alternative data directory
focus config set storage.data_path ~/focus-data
```

## Debug Mode

When issues persist, enable comprehensive debugging:

```bash
# Enable debug environment
export FOCUS_DEBUG=true
export FOCUS_LOG_LEVEL=debug
export FOCUS_LOG_FILE=~/focus-debug.log

# Run with debug output
focus --debug --verbose list

# Follow log file
tail -f ~/focus-debug.log
```

### Debug Information to Collect

When reporting issues, include:

```bash
# System information
focus status --all
focus version
go version
uname -a  # Unix
ver       # Windows

# Configuration
focus config list
focus config path

# Debug output
focus --debug add "test task" 2>&1 | tee debug-output.log

# Log files
cat ~/.local/log/focus/focus.log
```

## Getting Help

### Self-Service Resources

1. **Documentation**: [docs/](../docs/)
2. **Commands Reference**: [commands.md](user/commands.md)
3. **Configuration Guide**: [configuration.md](user/configuration.md)
4. **FAQ**: [faq.md](user/faq.md)

### Community Support

1. **GitHub Issues**: [Report bugs](https://github.com/kyanite/focus/issues)
2. **GitHub Discussions**: [Ask questions](https://github.com/kyanite/focus/discussions)
3. **Security Issues**: [Security Policy](../../SECURITY.md)

### When Reporting Issues

Include this information:

```bash
# Run this command and include output
focus status --all > system-info.txt
echo "Focus.sh version: $(focus --version)" >> system-info.txt
echo "Go version: $(go version)" >> system-info.txt
echo "OS: $(uname -a)" >> system-info.txt
echo "Current config:" >> system-info.txt
focus config list >> system-info.txt

# Include error messages and reproduction steps
```

## Common Solutions Summary

| Issue | Quick Fix |
|-------|-----------|
| "command not found" | Use `./focus` or add to PATH |
| "Ollama not available" | Run `ollama serve` or install Ollama |
| "Model not found" | Run `ollama pull qwen2.5:1.5b` |
| "Config not found" | Run `focus enhanced-config` |
| "AI not working" | Check `focus status --ai` |
| "Display issues" | Try `focus theme plain` |
| "Slow performance" | Enable caching or use faster model |
| "Network issues" | Check firewall or use local AI |

---

Need more help? Check the [FAQ](faq.md) or open an issue on GitHub.