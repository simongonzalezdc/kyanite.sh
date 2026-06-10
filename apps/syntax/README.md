# syntax.sh

> A terminal-based fiction writing and worldbuilding tool

**syntax.sh** is a beautiful, distraction-free writing environment for authors. Write stories, manage characters, organize scenes, track plots, and build your world—all in your terminal.

## Features (MVP v0.1.0)

This is an initial MVP implementation demonstrating the core architecture:

- ✅ **10 Beautiful Themes** - Cycle through Kyanite color themes with Ctrl+Shift+T
- ✅ **Project Management** - Create and manage multiple writing projects
- ✅ **Character Database** - Create and track story characters
- ✅ **Scene Organization** - Organize your story into chapters and scenes
- ✅ **Cross-Platform Storage** - Works on Linux, macOS, and Windows
- ✅ **Atomic File Operations** - Never lose your work

### Coming Soon (Full v1.0)

- 📝 Split-pane editor with markdown preview
- 🔄 Undo/redo functionality
- 🗺️  ASCII character relationship maps
- 📍 Location database with connections
- 📊 Writing statistics and progress tracking
- 📤 Export to PDF, DOCX, HTML, Markdown
- 🤖 AI story assistant (via Ollama)
- 📖 Story outline/plot tracker

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/simongonzalezdc/Syntax.sh.git
cd Syntax.sh

# Build the application
go build -o bin/syntax ./cmd/syntax

# Run it!
./bin/syntax
```

### Requirements

- Go 1.21 or higher
- A terminal with UTF-8 support
- Minimum 80x24 terminal size (120x30 recommended)

### Supported Platforms

- Linux (GNOME Terminal, Alacritty, kitty)
- macOS (Terminal.app, iTerm2)
- Windows (Windows Terminal with UTF-8 enabled)

## Usage

### Keyboard Shortcuts

**Global:**
- `Ctrl+Q` - Quit application
- `Ctrl+Shift+T` - Cycle through themes
- `?` - Show help (coming soon)
- `Esc` - Go back/cancel

**Welcome Screen:**
- `n` - Create new project
- `o` - Open existing project
- `q` - Quit

**Project Dashboard:**
- `c` - View characters
- `s` - View scenes
- `l` - View locations
- `e` - Export (coming soon)
- `Esc` - Back to welcome screen

**Character/Scene Lists:**
- `n` - Create new item
- `Esc` - Back to dashboard

### Project Structure

Projects are stored in:
- **Linux/macOS:** `~/.local/share/syntax/projects/`
- **Windows:** `%LOCALAPPDATA%\syntax\projects\`

Each project contains:
```
project-id/
├── metadata.yaml          # Project information
├── characters/            # Character files
├── locations/             # Location files
├── scenes/                # Scene files
├── outline/               # Story outline
├── stats/                 # Writing statistics
└── .backups/              # Automatic backups
```

### File Format

All data is stored in human-readable YAML frontmatter + Markdown format:

```markdown
---
id: char_abc123def4567890
name: "Jane Doe"
role: "protagonist"
occupation: "Detective"
---

# Jane Doe - Character Biography

Jane is a determined homicide detective...
```

## Themes

syntax.sh includes 10 beautiful Kyanite themes:

1. **Monochrome** - Classic black and white
2. **Amber Night** - Warm earth tones
3. **Twilight Mist** - Purple and lavender
4. **Indigo Depths** - Deep blues
5. **Forest Path** - Natural greens
6. **Clay Earth** - Rustic browns
7. **Iron Forge** - Industrial grays
8. **Sunlight** - Bright yellows
9. **Cyan Wave** - Ocean blues
10. **Electric Rose** - Vibrant pink

Press `Ctrl+Shift+T` to cycle through themes anytime.

## Architecture

This implementation follows the technical specifications in:
- **01-PRD-2.md** - Product Requirements
- **02-TDD.md** - Technical Design
- **KYANITE-STANDARDS.md** - Coding Standards
- **STORAGE-SPEC.md** - File Format Specification
- **SECURITY.md** - Security Policies

### Technology Stack

- **Go 1.21+** - Programming language
- **Bubble Tea** - Terminal UI framework
- **Lipgloss** - Styling and theming
- **Cross-platform** - Uses XDG Base Directory spec

### Project Structure

```
syntax/
├── cmd/syntax/            # Main entry point
├── internal/
│   ├── app/               # Bubble Tea UI
│   ├── character/         # Character types
│   ├── scene/             # Scene types
│   ├── location/          # Location types
│   ├── story/             # Project types
│   ├── storage/           # File I/O
│   └── theme/             # Theme system
├── tests/                 # Test files
├── go.mod                 # Dependencies
└── README.md              # This file
```

## Development

### Building from Source

```bash
# Install dependencies
go mod download

# Build
go build -o bin/syntax ./cmd/syntax

# Run tests
go test ./...

# Run with debug logging
DEBUG=1 ./bin/syntax
```

### Running Tests

```bash
# All tests
go test -v ./...

# With coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please see **CONTRIBUTING.md** for guidelines.

### Code of Conduct

Be respectful, welcoming, and constructive. See **CONTRIBUTING.md** for details.

## Security

Found a security vulnerability? Please email security@kyanite.sh.

See **SECURITY.md** for our full security policy.

## Roadmap

See **01-PRD-2.md** for the complete feature roadmap.

### Current Status: MVP (v0.1.0)

- [x] Project structure
- [x] Theme system
- [x] Basic storage layer
- [x] Character management
- [x] Scene management
- [x] Project navigation

### Next Release: v1.0.0

- [ ] Split-pane editor
- [ ] Markdown preview
- [ ] Undo/redo
- [ ] Writing statistics
- [ ] Export functionality
- [ ] AI assistant
- [ ] Complete test coverage

## License

MIT License - see LICENSE file for details

## Part of Kyanite Suite

syntax.sh is part of the Kyanite Suite of terminal productivity tools.

**Other tools:**
- **prism.sh** - Color palette generator (coming soon)
- **focus.sh** - Pomodoro timer (coming soon)

## Support

- **Documentation:** See the `*.md` files in this repository
- **Issues:** https://github.com/simongonzalezdc/Syntax.sh/issues
- **Discussions:** https://github.com/simongonzalezdc/Syntax.sh/discussions

## Acknowledgments

Built with:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) by Charm
- [Lipgloss](https://github.com/charmbracelet/lipgloss) by Charm
- [xdg](https://github.com/adrg/xdg) by Adrian-George Bostan

---

**"Your novel deserves a beautiful editor. syntax.sh is built for the story first, complexity second."**
