# NEON Focus - Project Structure

```
.
├── cmd/
│   └── neon/
│       └── main.go              # Application entry point
├── internal/
│   ├── ai/
│   │   ├── manager.go           # AI manager with Ollama/OpenRouter
│   │   ├── openrouter.go        # OpenRouter API integration
│   │   └── manager_test.go     # AI manager tests
│   ├── cli/
│   │   ├── root.go              # Root command setup
│   │   ├── add.go               # Add task command
│   │   ├── list.go              # List tasks command
│   │   ├── complete.go          # Complete task command
│   │   ├── delete.go            # Delete task command
│   │   ├── suggest.go           # AI suggestions command
│   │   ├── chat.go              # AI chat command
│   │   ├── calendar.go          # Calendar commands (NEW!)
│   │   ├── config.go            # Configuration command
│   │   ├── theme.go             # Theme switching command
│   │   ├── dashboard.go         # TUI dashboard command
│   │   ├── interactive.go       # Gum interactive commands
│   │   ├── filter.go            # Gum filter command
│   │   ├── notes.go             # Notes management command
│   │   ├── wizard.go            # Huh wizard command
│   │   ├── config-wizard.go     # Configuration wizard
│   │   ├── edit-wizard.go       # Edit task wizard
│   │   ├── enhanced_config.go   # Enhanced configuration
│   │   ├── unified_dashboard.go  # Unified dashboard command
│   │   └── viper_commands.go     # Viper config commands
│   ├── engine/
│   │   └── engine.go            # Task management logic
│   ├── store/
│   │   └── store.go             # Task persistence
│   ├── tui/
│   │   ├── main.go              # Main TUI application
│   │   ├── dashboard.go         # TUI dashboard component
│   │   ├── tasklist.go          # Task selection interface
│   │   └── unified_dashboard.go  # Unified TUI dashboard
│   └── wizards/
│       └── task_creation.go     # Task creation wizards
├── pkg/
│   ├── calendar/
│   │   └── calendar.go          # Calendar system (NEW!)
│   ├── config/
│   │   └── config.go            # Configuration management (NEW!)
│   ├── models/
│   │   └── task.go              # Data structures
│   ├── validation/
│   │   └── validation.go        # Input validation (NEW!)
│   ├── styles/
│   │   └── styles.go            # Styling utilities
│   └── utils/
│       └── storage.go           # Storage utilities
├── test/
│   ├── task_test.go            # Integration tests
│   └── integration_test.go     # Integration tests
├── go.mod                        # Go module dependencies
├── go.sum                        # Go module checksums
├── README.md                     # Project documentation
├── build.bat                     # Build script (Windows)
├── CALENDAR_PLAN.md              # Calendar implementation plan
└── NEON_DEVELOPMENT_PLAN.md      # Development progress
```

## Component Overview

### cmd/neon/main.go
The entry point of the NEON CLI application.

### internal/ai/
Handles all AI interactions:
- **manager.go**: Core AI logic with Ollama/OpenRouter integration and smart caching
- **openrouter.go**: Remote AI API integration
- **manager_test.go**: Comprehensive AI functionality tests

### internal/cli/
Command-line interface with 20+ commands:
- **Core Commands**: add, list, complete, delete, suggest, chat
- **Calendar Commands**: calendar, show, today, add, list, navigate (NEW!)
- **Interactive Commands**: interactive, filter, config, theme, notes
- **Wizard Commands**: wizard, config-wizard, edit-wizard (Huh forms)
- **Configuration Commands**: config, enhanced-config, viper commands
- **UI Commands**: dashboard, unified-dashboard (TUI interfaces)

### internal/engine
Business logic layer:
- **engine.go**: Core task management with advanced filtering and operations

### internal/store
Data persistence layer:
- **store.go**: Atomic file operations with JSON task storage

### internal/tui
Terminal User Interface components:
- **dashboard.go**: Advanced TUI dashboard with AI integration
- **unified_dashboard.go**: Complete feature integration in single interface

### internal/wizards
Interactive form components:
- **task_creation.go**: Advanced Huh-based task creation wizards

### pkg/calendar
Calendar system (NEW!):
- **calendar.go**: Full calendar implementation with month/week/day views

### pkg/config
Configuration management (NEW!):
- **config.go**: Viper-based configuration with persistence

### pkg/validation
Input validation (NEW!):
- **validation.go**: SQL injection protection and input sanitization

### pkg/models
Data structures:
- **task.go**: Complete task and parsed task models with validation

### pkg/styles
Visual styling utilities:
- **styles.go**: Synthwave-themed styling with Lip Gloss

### pkg/utils
Utility functions:
- **storage.go**: Cross-platform storage path management

## Key Features Implemented

### 🚀 AI Integration
- Local AI processing with Ollama (qwen2.5:1.5b)
- Remote AI fallback with OpenRouter
- Smart caching for performance
- Natural language task and date parsing

### 🎨 User Interfaces
- Advanced TUI dashboard with Bubble Tea
- Interactive forms with Huh
- Command-line Gum integration
- Calendar system with multiple views
- Unified dashboard for all features

### 📅 Calendar System
- Month/week/day calendar views
- Natural language date parsing
- Task integration with priority indicators
- CLI calendar commands

### ⚙️ Configuration Management
- Viper-based configuration system
- Persistent settings storage
- Interactive configuration wizards
- Theme switching capabilities

### 🔒 Security & Validation
- SQL injection protection
- Input sanitization
- Type-safe operations
- Comprehensive error handling

### 🧪 Testing
- Comprehensive test suite
- Performance benchmarks
- Integration tests
- CLI command testing

## Architecture

The application follows clean architecture principles:

1. **CLI Layer**: Command-line interface handling
2. **Business Logic**: Task management and operations
3. **Data Layer**: Persistence and storage
4. **AI Layer**: Natural language processing
5. **UI Layer**: Terminal user interface

Each layer has clear separation of concerns and minimal dependencies, ensuring maintainability and scalability.