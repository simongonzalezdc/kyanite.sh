# Installation Guide

This guide will help you install Focus.sh on your system and get it running with AI capabilities.

## System Requirements

- **Go**: 1.21.0 or higher
- **Operating System**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **Memory**: 4GB RAM minimum (8GB+ recommended for AI models)
- **Storage**: 500MB free space for the application
- **AI Models**: 2-4GB for local AI models (Ollama)

## Installation Methods

### Method 1: Build from Source (Recommended)

#### Prerequisites

1. **Install Go**
   
   - **Windows**: Download from [golang.org](https://golang.org/dl/) or use winget: `winget install GoLang.Go`
   - **macOS**: `brew install go`
   - **Linux**: `sudo apt install golang-go` (Ubuntu/Debian) or `sudo dnf install golang` (Fedora)

2. **Install Ollama** (for local AI processing)
   
   - **Windows**: Download from [ollama.ai](https://ollama.ai/download)
   - **macOS**: `brew install ollama`
   - **Linux**: `curl -fsSL https://ollama.ai/install.sh | sh`

#### Installation Steps

1. **Clone the repository**
   ```bash
   git clone https://github.com/kyanite/focus.git
   cd focus
   ```

2. **Download dependencies**
   ```bash
   go mod download
   ```

3. **Pull the recommended AI model**
   ```bash
   ollama pull qwen2.5:1.5b
   ```

4. **Build the application**
   ```bash
   # Windows
   build.bat
   
   # Unix (macOS/Linux)
   chmod +x build.sh
   ./build.sh
   ```

5. **Verify installation**
   ```bash
   ./focus --version
   ./focus --help
   ```

### Method 2: Run Directly with Go

If you prefer not to build, you can run directly:

```bash
# Clone the repository
git clone https://github.com/kyanite/focus.git
cd focus

# Run directly
go run cmd/focus/main.go --help
```

### Method 3: Pre-built Binaries (When Available)

Download pre-built binaries from [GitHub Releases](https://github.com/kyanite/focus/releases):

```bash
# Example for Linux
wget https://github.com/kyanite/focus/releases/download/v1.0.0/focus-linux-amd64.tar.gz
tar -xzf focus-linux-amd64.tar.gz
sudo mv focus /usr/local/bin/
```

## AI Configuration

### Option 1: Local AI (Ollama) - Recommended

1. **Start Ollama**
   ```bash
   # Background service (runs automatically on most systems)
   ollama serve &
   ```

2. **Verify Ollama is running**
   ```bash
   ollama list
   ```

3. **Test the model**
   ```bash
   ollama run qwen2.5:1.5b "Hello, can you help me with task management?"
   ```

### Option 2: Remote AI (OpenRouter) - Fallback

If you can't run Ollama locally, configure OpenRouter:

1. **Get API key** from [OpenRouter.ai](https://openrouter.ai/)

2. **Configure environment**
   ```bash
   # Create .env file
   echo "OPENROUTER_API_KEY=<your-openrouter-key>" > .env
   
   # Set default fallback model
   echo "FALLBACK_MODEL=mistralai/mistral-7b-instruct" >> .env
   ```

### Option 3: Enhanced Configuration Wizard

Run the configuration wizard for guided setup:

```bash
go run cmd/focus/main.go enhanced-config
```

This will:
- Detect available AI providers
- Test connectivity
- Set up default configuration
- Create necessary directories

## Platform-Specific Setup

### Windows

1. **Install Git for Windows** (if not already installed)
2. **Use PowerShell or Command Prompt**
3. **Add to PATH** (optional):
   ```powershell
   # Add to current session
   $env:PATH += ";$(pwd)"
   
   # Add permanently (run as Administrator)
   [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$(pwd)", "Machine")
   ```

4. **Create desktop shortcut** (optional):
   ```batch
   # Use the provided script
   create-desktop-shortcut.bat
   ```

### macOS

1. **Install Homebrew** (if not installed):
   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install dependencies**:
   ```bash
   brew install go ollama
   ```

3. **Add to PATH** (if needed):
   ```bash
   echo 'export PATH=$PATH:~/focus' >> ~/.zshrc
   source ~/.zshrc
   ```

### Linux (Ubuntu/Debian)

1. **Install system dependencies**:
   ```bash
   sudo apt update
   sudo apt install golang git curl
   ```

2. **Install Ollama**:
   ```bash
   curl -fsSL https://ollama.ai/install.sh | sh
   sudo systemctl enable ollama
   sudo systemctl start ollama
   ```

3. **Add to PATH** (optional):
   ```bash
   echo 'export PATH=$PATH:~/focus' >> ~/.bashrc
   source ~/.bashrc
   ```

## Verification

After installation, verify everything is working:

```bash
# Check application
./focus --version

# Test AI functionality
./focus add "Test task to verify AI parsing is working"

# List tasks
./focus list

# Test calendar
./focus calendar today
```

Expected output should show:
- Application version information
- AI successfully parsing your task
- Tasks displayed with proper formatting
- Calendar view loading correctly

## Troubleshooting

### Common Issues

#### "command not found: focus"
- **Solution**: Add the focus directory to your PATH or use `./focus`

#### "Ollama not available"
- **Solution**: Install Ollama or configure OpenRouter fallback
- **Check**: `ollama list` to verify Ollama is running

#### "Permission denied"
- **Linux/macOS**: `chmod +x focus` or `sudo chmod +x focus`
- **Windows**: Run PowerShell as Administrator

#### "Go modules not found"
- **Solution**: `go mod download` in the focus directory

#### "Model not found"
- **Solution**: `ollama pull qwen2.5:1.5b`

#### "Port 11434 already in use"
- **Solution**: Another Ollama instance is running, stop it first

### Getting Help

If you encounter issues:

1. **Check logs**: Look for error messages in the terminal output
2. **Verify installation**: Run the verification steps above
3. **Check dependencies**: Ensure Go and Ollama are properly installed
4. **Search issues**: [GitHub Issues](https://github.com/kyanite/focus/issues)
5. **Ask for help**: [GitHub Discussions](https://github.com/kyanite/focus/discussions)

## Next Steps

Once installed:

1. **Read the [Configuration Guide](configuration.md)** for detailed setup
2. **Explore [Commands](commands.md)** to learn all available features
3. **Customize [Themes](themes.md)** for your preferred visual style
4. **Check the [FAQ](faq.md)** for common questions

Happy task managing! 🚀
