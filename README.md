# noise.sh

AI-powered songwriting assistant - *From chaos to song*

## Overview

noise.sh is a terminal-based songwriting application that leverages AI to help musicians and lyricists create, refine, and perfect their lyrics. Built with Go and the Bubble Tea TUI framework, it provides a clean, efficient interface for creative work.

## Features

- **AI-Powered Lyrics**: Generate and refine lyrics with AI assistance
- **Song Management**: Organize songs with metadata and versioning
- **Auto-save**: Never lose your work with automatic saving
- **Plugin System**: Extensible architecture for custom functionality
- **Terminal UI**: Clean, responsive interface that works in any terminal
- **Real-time Preview**: See your changes as you type
- **Collaboration**: Work together with other songwriters

## Installation

### From Source

```bash
git clone https://github.com/puente-labs/noise.git
cd noise
go build -o noise cmd/noise/main.go
```

### Binary Download

Download the latest release from the [GitHub releases page](https://github.com/puente-labs/noise/releases).

## Usage

### Basic Usage

```bash
# Start noise.sh
noise

# Start with debug logging
noise --debug

# Open a specific song file
noise song.md

# Show version information
noise --version

# Show help
noise --help
```

### Configuration

noise.sh stores configuration in `~/.config/noise/noise.yaml` and data in `~/.noise/`.

Environment variables can be used to override configuration:

```bash
# Set data directory
export NOISE_APP_DATA_DIR="/path/to/data"

# Enable debug mode
export NOISE_DEV_DEBUG="true"

# Set AI API key
export NOISE_AI_API_KEY="your-api-key"
```

## Development

### Prerequisites

- Go 1.25.3 or later
- SQLite3

### Building

```bash
# Build for current platform
go build -o noise cmd/noise/main.go

# Build for multiple platforms
go build -o noise-linux-amd64 -GOOS=linux -GOARCH=amd64 cmd/noise/main.go
go build -o noise-darwin-amd64 -GOOS=darwin -GOARCH=amd64 cmd/noise/main.go
go build -o noise-windows-amd64 -GOOS=windows -GOARCH=amd64 cmd/noise/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...
```

## Configuration

noise.sh uses YAML configuration files. The default configuration looks like:

```yaml
app:
  name: "noise.sh"
  version: "1.0.0"
  data_dir: "~/.noise"
  auto_save: true
  auto_save_interval: "30s"

ui:
  theme: "violet-dusk"
  font_size: 12
  show_line_numbers: true
  word_wrap: true
  animations: true

ai:
  provider: "ollama"
  model: "qwen2.5:7b-instruct"
  base_url: "http://localhost:11434"
  temperature: 0.7
  max_tokens: 2048
  enabled: true
```

## Plugin Development

noise.sh supports plugins for extending functionality. Plugins can be placed in:

- `~/.noise/plugins/`
- `./plugins/`
- `./internal/plugins/examples/`

See `internal/plugins/examples/` for example plugins.

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

- **Documentation**: [Project Wiki](https://github.com/puente-labs/noise/wiki)
- **Issues**: [GitHub Issues](https://github.com/puente-labs/noise/issues)
- **Discussions**: [GitHub Discussions](https://github.com/puente-labs/noise/discussions)

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- UI styling with [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- Database operations with [SQLite](https://sqlite.org/)

---

**noise.sh** - From chaos to song 🎵
