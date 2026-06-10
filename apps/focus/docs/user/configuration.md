# Configuration Guide

This guide covers all aspects of configuring Focus.sh, from AI providers to visual themes and advanced settings.

## Configuration Overview

Focus.sh uses a layered configuration system:

1. **Default values** - Built-in sensible defaults
2. **Configuration file** - `~/.config/focus/config.json`
3. **Environment variables** - Override configuration
4. **Command-line flags** - Temporary overrides

## Quick Configuration

### Automated Setup (Recommended)

Run the enhanced configuration wizard for guided setup:

```bash
focus enhanced-config
```

This interactive wizard will:
- Detect available AI providers
- Test connectivity and model availability
- Set up default preferences
- Create necessary directories
- Generate initial configuration

### Basic Configuration

```bash
# View current configuration
focus config list

# Set AI provider
focus config set ai.provider ollama

# Set theme
focus config set theme synthwave

# Set default dashboard
focus config set dashboard unified
```

## AI Configuration

### Primary AI: Ollama (Local)

#### Setup
```bash
# Install Ollama
# https://ollama.ai/download

# Pull recommended model
ollama pull qwen2.5:1.5b

# Start Ollama service
ollama serve
```

#### Configuration
```json
{
  "ai": {
    "provider": "ollama",
    "model": "qwen2.5:1.5b",
    "base_url": "http://localhost:11434",
    "timeout": 30,
    "max_tokens": 1000
  }
}
```

#### Available Models
- `qwen2.5:1.5b` (recommended - fast, capable)
- `llama3.2:3b` (good balance)
- `mistral:7b` (more powerful, slower)

### Fallback AI: OpenRouter (Remote)

#### Setup
```bash
# Get API key from https://openrouter.ai/
export OPENROUTER_API_KEY="<your-openrouter-key>"
```

#### Configuration
```json
{
  "ai": {
    "fallback_provider": "openrouter",
    "fallback_model": "mistralai/mistral-7b-instruct",
    "openrouter_api_key": "${OPENROUTER_API_KEY}",
    "fallback_enabled": true
  }
}
```

#### Environment Variables
```bash
# .env file
OPENROUTER_API_KEY=<your-openrouter-key>
FALLBACK_MODEL=mistralai/mistral-7b-instruct
DEFAULT_MODEL=qwen2.5:1.5b
OLLAMA_BASE_URL=http://localhost:11434
```

## Visual Configuration

### Themes

Focus.sh includes 10 beautiful built-in themes:

#### Dark Themes
- **amber-night** (default) - Warm amber and purple accents
- **twilight-mist** - Soft purple and blue mist tones
- **indigo-depths** - Deep blue and indigo palette
- **forest-path** - Natural green forest colors
- **clay-earth** - Warm brown earth tones
- **iron-forge** - Industrial red and gray forge colors
- **cyan-wave** - Cool cyan and blue ocean tones
- **electric-rose** - Vibrant pink and electric blue

#### Light Theme  
- **sunlight** - Bright golden yellows and warm tones

#### Minimal Theme
- **monochrome** - Clean black and white design

#### Custom Themes
Create custom themes in `~/.config/focus/themes/`:

```json
{
  "name": "custom-theme",
  "colors": {
    "primary": "#00ff00",
    "secondary": "#ff00ff",
    "background": "#000000",
    "text": "#ffffff",
    "accent": "#00ffff"
  },
  "styles": {
    "border": "double",
    "padding": 2,
    "margin": 1
  }
}
```

### Dashboard Configuration

```json
{
  "dashboard": {
    "default_view": "unified",
    "auto_refresh": true,
    "refresh_interval": 30,
    "show_calendar": true,
    "show_tasks": true,
    "show_ai_chat": true
  }
}
```

## Data Storage Configuration

### Default Locations

#### Unix/macOS
- **Config**: `~/.config/focus/`
- **Data**: `~/.local/share/focus/`
- **Logs**: `~/.local/log/focus/`

#### Windows
- **Config**: `%USERPROFILE%\.config\focus\`
- **Data**: `%USERPROFILE%\.local\share\focus\`
- **Logs**: `%USERPROFILE%\.local\log\focus\`

### Custom Locations

```bash
# Set custom data directory
focus config set storage.data_path "/path/to/custom/data"

# Set custom config directory
focus config set storage.config_path "/path/to/custom/config"
```

### Configuration Files

#### `config.json` - Main configuration
```json
{
  "ai": {
    "provider": "ollama",
    "model": "qwen2.5:1.5b",
    "base_url": "http://localhost:11434"
  },
  "theme": "synthwave",
  "dashboard": {
    "default_view": "unified",
    "auto_refresh": true
  },
  "storage": {
    "data_path": "~/.local/share/focus",
    "config_path": "~/.config/focus"
  }
}
```

#### `.env` - Environment variables
```bash
# AI Configuration
OLLAMA_BASE_URL=http://localhost:11434
DEFAULT_MODEL=qwen2.5:1.5b
OPENROUTER_API_KEY=<your-openrouter-key>
FALLBACK_MODEL=mistralai/mistral-7b-instruct

# Google Calendar (optional)
GOOGLE_CALENDAR_ID=primary
GOOGLE_CREDENTIALS_FILE=~/.config/focus/google-credentials.json

# Debug
FOCUS_DEBUG=true
FOCUS_LOG_LEVEL=info
```

## Advanced Configuration

### Performance Settings

```json
{
  "performance": {
    "cache_enabled": true,
    "cache_ttl": 3600,
    "parallel_processing": true,
    "max_concurrent_requests": 5
  }
}
```

### AI Behavior

```json
{
  "ai": {
    "provider": "ollama",
    "model": "qwen2.5:1.5b",
    "temperature": 0.7,
    "max_tokens": 1000,
    "system_prompt": "You are a helpful task management assistant.",
    "cache_responses": true,
    "context_window": 10
  }
}
```

### Calendar Integration

```json
{
  "calendar": {
    "provider": "google",
    "enabled": true,
    "auto_sync": false,
    "sync_interval": 3600,
    "default_calendar": "primary",
    "task_reminder": true,
    "reminder_advance": 15
  }
}
```

### Audio Feedback

```json
{
  "audio": {
    "enabled": true,
    "completion_sound": true,
    "notification_sound": true,
    "volume": 0.7,
    "sound_theme": "synthwave"
  }
}
```

## Command Reference

### Configuration Commands

```bash
# View all configuration
focus config list

# Get specific value
focus config get ai.provider

# Set value
focus config set theme synthwave

# Reset to default
focus config reset

# Show config file path
focus config path

# Export configuration
focus config export > my-config.json

# Import configuration
focus config import my-config.json
```

### Theme Commands

```bash
# List available themes
focus theme list

# Set theme
focus theme synthwave

# Create theme from current settings
focus theme save my-custom-theme

# Preview theme
focus theme preview synthwave
```

## Environment Variables

### AI Configuration
```bash
OLLAMA_BASE_URL=http://localhost:11434
DEFAULT_MODEL=qwen2.5:1.5b
OPENROUTER_API_KEY=<your-openrouter-key>
FALLBACK_MODEL=mistralai/mistral-7b-instruct
AI_TIMEOUT=30
```

### Storage Configuration
```bash
FOCUS_CONFIG_DIR=~/.config/focus
FOCUS_DATA_DIR=~/.local/share/focus
FOCUS_LOG_DIR=~/.local/log/focus
```

### Debug Configuration
```bash
FOCUS_DEBUG=true
FOCUS_LOG_LEVEL=debug  # debug, info, warn, error
FOCUS_LOG_FILE=~/.local/log/focus/debug.log
```

## Troubleshooting Configuration

### Common Issues

#### "AI provider not found"
```bash
# Check provider configuration
focus config get ai.provider

# Test Ollama connectivity
curl http://localhost:11434/api/tags

# Verify model is available
ollama list
```

#### "Configuration file not found"
```bash
# Show config path
focus config path

# Create default config
focus enhanced-config

# Check file permissions
ls -la ~/.config/focus/
```

#### "Theme not found"
```bash
# List available themes
focus theme list

# Reset to default theme
focus config set theme synthwave

# Verify theme files
ls -la ~/.config/focus/themes/
```

### Debug Mode

Enable debug logging for troubleshooting:

```bash
# Set debug mode
export FOCUS_DEBUG=true
export FOCUS_LOG_LEVEL=debug

# Run with debug output
focus --debug list

# Check log file
tail -f ~/.local/log/focus/focus.log
```

## Backup and Restore

### Backup Configuration

```bash
# Backup all settings
tar -czf focus-backup-$(date +%Y%m%d).tar.gz ~/.config/focus/ ~/.local/share/focus/

# Backup just configuration
cp ~/.config/focus/config.json config-backup.json
```

### Restore Configuration

```bash
# Restore from backup
tar -xzf focus-backup-YYYYMMDD.tar.gz -C ~/

# Import configuration
focus config import config-backup.json
```

## Next Steps

After configuration:

1. **Test your setup**: `focus add "Test AI parsing"`
2. **Explore commands**: See [Commands Guide](commands.md)
3. **Customize themes**: See [Themes Guide](themes.md)
4. **Set up calendar**: Configure Google Calendar if needed

Need help? Check the [Troubleshooting Guide](troubleshooting.md) or open an issue on GitHub.
