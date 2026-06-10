# AI-Powered To-Do Assistant - Technical Design Document

## Overview
This document outlines the technical architecture and design decisions for the AI-Powered To-Do Assistant CLI application. The system combines Charm libraries for the TUI with local AI capabilities through Ollama.

## System Architecture

### Component Diagram
```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Parser    │────┤    Task Engine   │────┤   AI Manager    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │    │                  │    │
                    ┌─────────┘    └─────────┐        │    └───────┐
                    ▼                        ▼        ▼            ▼
         ┌──────────────────┐      ┌────────────────────┐  ┌──────────────┐
         │   Task Store     │      │  Bubble Tea TUI    │  │ Ollama API   │
         └──────────────────┘      └────────────────────┘  └──────────────┘
                                                             │
                                                   ┌─────────▼──────────┐
                                                   │  OpenRouter API    │
                                                   └────────────────────┘
```

### Core Modules

#### 1. CLI Parser (`internal/cli`)
Handles command line argument parsing and routing to appropriate handlers.

**Commands:**
- `add <natural language input>` - Add new task
- `list [--filter=active|completed|all]` - List tasks
- `complete <task_id>` - Mark task as complete
- `delete <task_id>` - Delete task
- `suggest` - Get AI-powered suggestions
- `summary` - Get AI-generated summary
- `pomodoro [task_id]` - Start Pomodoro timer for focused work

#### 2. Task Engine (`internal/engine`)
Core business logic for task management operations.

**Functions:**
- Task CRUD operations
- Task filtering and sorting
- State management
- Validation logic

#### 3. AI Manager (`internal/ai`)
Handles all AI interactions with fallback logic.

**Features:**
- Local model inference via Ollama
- Remote fallback to OpenRouter
- Response validation and hallucination detection
- Caching mechanism
- Model benchmarking and selection

#### 4. Task Store (`internal/store`)
Persistent storage for tasks.

**Structure:**
```go
type Task struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Status      string    `json:"status"` // pending, completed
    Priority    string    `json:"priority"` // low, medium, high
    Deadline    time.Time `json:"deadline,omitempty"`
    Categories  []string  `json:"categories,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### 5. TUI (`internal/tui`)
Interactive terminal interface using Bubble Tea.

**Components:**
- Main view with task list
- Input forms using Huh
- Status indicators
- Navigation controls

## AI Integration Design

### Model Selection Strategy
1. **Primary**: Llama3 8B (local via Ollama) - Best balance of performance and resource usage
2. **Alternative**: Mistral 7B - For comparison in parsing tasks
3. **Fallback**: OpenRouter with GPT-3.5-turbo - For internet-connected scenarios

### Prompt Engineering
```
System Prompt:
"You are a task parsing assistant. Extract the following from user input:
1. Task description
2. Deadline (if mentioned)
3. Priority (low/medium/high)
4. Categories (if mentioned)

Format response as JSON:
{
  "description": "task description",
  "deadline": "ISO date format or null",
  "priority": "low|medium|high",
  "categories": ["category1", "category2"]
}

If information is not provided, omit the field or use null."
```

### Hallucination Detection
1. **Validation Rules**:
   - Description must be non-empty
   - Deadline must be reasonable (> today)
   - Priority must be one of: low, medium, high
   - Categories must be meaningful (no gibberish)

2. **Confidence Scoring**:
   - Parse uncertainty tokens in responses
   - Flag responses with low confidence
   - Request clarification for ambiguous inputs

### Caching Strategy
1. **Local Cache**: Store recent AI responses in memory
2. **Persistent Cache**: File-based cache for common queries
3. **Invalidation**: Clear cache on model updates or after 24 hours

## Data Flow

### Task Addition Flow
1. User enters natural language task via CLI or TUI
2. AI Manager processes input through Ollama
3. Response validated for hallucinations
4. Structured task sent to Task Engine
5. Task Engine stores in Task Store
6. Confirmation displayed to user

### Fallback Flow
1. Ollama health check fails
2. AI Manager switches to OpenRouter
3. Same validation applied to remote responses
4. User notified of remote processing (for API cost awareness)
5. Response cached locally for future use

## Error Handling

### AI Errors
- **Model Unavailable**: Gracefully degrade to manual input
- **Invalid Response**: Request rephrasing or use default parsing
- **Hallucination Detected**: Flag for review and request clarification

### TUI Errors
- **Navigation Issues**: Clear error messaging with recovery options
- **Input Validation**: Real-time feedback on task constraints
- **Storage Failures**: Preserve data in memory and retry

## Security Considerations
- No sensitive data sent to remote APIs without explicit user consent
- Local storage encryption option for sensitive task data
- API key management through environment variables
- Rate limiting for remote API calls

## Performance Requirements
- Startup time: <2 seconds
- Local AI response: <1 second
- Remote AI response: <3 seconds
- Memory usage: <100MB
- Storage usage: <10MB

## Testing Strategy

### Unit Tests
- AI parsing accuracy (>85%)
- Hallucination detection effectiveness
- TUI component behavior
- Storage operations

### Integration Tests
- Ollama API integration
- OpenRouter fallback
- CLI command handling
- TUI navigation

### End-to-End Tests
- Complete task addition workflow
- Offline/online transition
- Cross-platform compatibility
- Performance benchmarks

## Deployment

### Build Process
1. Cross-platform compilation (Windows, macOS, Linux)
2. Single binary artifact
3. Embedded default configuration
4. Automatic Ollama model checking

### Installation
1. Download binary
2. Make executable
3. Run initialization command to check/setup Ollama
4. Optional: Install as system command

### Updates
1. Version checking against GitHub releases
2. Automatic download and replacement
3. Backup of previous version
4. Rollback capability

## Future Extensibility
- Plugin system for additional AI providers
- Custom task types (habits, projects, etc.)
- Sync capabilities with external services
- Advanced analytics and insights