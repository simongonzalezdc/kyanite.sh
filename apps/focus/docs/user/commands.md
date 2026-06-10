# Commands Reference

Complete reference for all Focus.sh commands, including options, examples, and usage patterns.

## Command Structure

```bash
focus [global_flags] <command> [subcommand] [arguments] [flags]
```

### Global Flags
- `--help, -h` - Show help information
- `--version` - Show version information
- `--config` - Specify config file location
- `--debug` - Enable debug logging
- `--quiet` - Suppress non-error output

## Core Task Management

### `focus add` - Add Task
Add new tasks with natural language processing.

#### Syntax
```bash
focus add [flags] <task_description>
```

#### Examples
```bash
# Simple task
focus add "Buy groceries"

# Task with priority
focus add --priority high "Complete project proposal"

# Task with due date
focus add --due 2024-01-25 "Submit tax documents"

# Task with category
focus add --category work "Team meeting at 2 PM"

# Natural language parsing
focus add "Schedule dentist appointment for next Tuesday morning"
```

#### Flags
- `--priority, -p` - Priority level: low, medium, high
- `--category, -c` - Task category
- `--due, -d` - Due date (YYYY-MM-DD or natural language)
- `--tags, -t` - Comma-separated tags

---

### `focus list` - List Tasks
Display tasks with various filtering and sorting options.

#### Syntax
```bash
focus list [flags]
```

#### Examples
```bash
# Show all tasks
focus list

# Filter by status
focus list --status pending
focus list --status completed

# Filter by category
focus list --category work
focus list --category personal

# Sort by priority
focus list --sort priority

# Show with calendar view
focus list --calendar

# Compact view
focus list --compact
```

#### Flags
- `--status, -s` - Filter by status: all, pending, completed
- `--category, -c` - Filter by category
- `--priority, -p` - Filter by priority
- `--sort` - Sort by: created, due, priority, category
- `--format, -f` - Output format: table, list, json
- `--compact` - Compact display mode
- `--calendar` - Calendar view integration

---

### `focus complete` - Complete Task
Mark tasks as completed.

#### Syntax
```bash
focus complete [flags] <task_id|task_pattern>
```

#### Examples
```bash
# Complete by ID
focus complete 5

# Complete by exact match
focus complete "Buy groceries"

# Complete multiple tasks
focus complete 3 7 12

# Complete by pattern
focus complete --pattern "*meeting*"

# Complete with confirmation
focus complete --confirm "Project proposal"
```

#### Flags
- `--confirm` - Ask for confirmation before completing
- `--pattern` - Use pattern matching instead of exact ID
- `--all` - Complete all matching tasks

---

### `focus delete` - Delete Task
Remove tasks from the system.

#### Syntax
```bash
focus delete [flags] <task_id|task_pattern>
```

#### Examples
```bash
# Delete by ID
focus delete 5

# Delete by exact match
focus delete "Old completed task"

# Delete with confirmation
focus delete --confirm "Cancel subscription"

# Force delete without confirmation
focus delete --force 15
```

#### Flags
- `--confirm` - Ask for confirmation before deleting
- `--force` - Skip confirmation prompts
- `--pattern` - Use pattern matching

---

## Calendar Integration

### `focus calendar` - Calendar Commands
Interact with calendar views and scheduled tasks.

#### `focus calendar show` - Display Calendar
```bash
# Month view
focus calendar show month

# Week view
focus calendar show week

# Day view
focus calendar show day

# Today's view
focus calendar today
```

#### `focus calendar add` - Add Scheduled Task
```bash
# Add task with date
focus calendar add "Team meeting" 2024-01-25

# Add with time
focus calendar add "Doctor appointment" "2024-01-25 14:30"

# Natural language
focus calendar add "Lunch with Sarah next Friday at noon"
```

#### `focus calendar navigate` - Navigate Calendar
```bash
# Go to specific date
focus calendar navigate 2024-01-25

# Next/Previous
focus calendar navigate next
focus calendar navigate prev

# Today
focus calendar navigate today
```

---

## AI-Powered Features

### `focus chat` - AI Assistant
Interactive AI chat for task management assistance.

#### Examples
```bash
# Start chat mode
focus chat

# Ask for help
focus chat "How do I prioritize my tasks?"

# Get suggestions
focus chat "Suggest some tasks for my learning goals"

# Task analysis
focus chat "Analyze my current workload"
```

#### Features
- Natural language conversation
- Task suggestions and recommendations
- Productivity advice
- Help with Focus.sh features

---

### `focus inspire` - AI Task Suggestions
Get AI-powered task suggestions based on your context.

#### Examples
```bash
# General suggestions
focus inspire

# Category-specific
focus inspire --category learning

# Time-based
focus inspire --duration 30

# Goal-oriented
focus inspire --goal "learn Go programming"
```

#### Flags
- `--category, -c` - Suggest tasks for specific category
- `--duration, -d` - Task duration in minutes
- `--goal, -g` - Goal-based suggestions

---

## Interactive Features

### `focus unified` - Unified Dashboard
Launch the complete interactive dashboard.

#### Examples
```bash
# Launch unified dashboard
focus unified

# With specific starting view
focus unified --view tasks

# Auto-refresh mode
focus unified --auto-refresh 30
```

#### Features
- Task management
- Calendar integration
- AI chat interface
- Real-time updates
- Theme switching

---

### `focus interactive` - Interactive Task Creation
Guided task creation with interactive forms.

#### Examples
```bash
# Launch interactive mode
focus interactive

# Start with pre-filled data
focus interactive --category work --priority high
```

#### Features
- Step-by-step task creation
- AI-powered suggestions
- Category and priority selection
- Date and time picking

---

### `focus dashboard` - TUI Dashboard
Terminal-based dashboard interface.

#### Examples
```bash
# Launch dashboard
focus dashboard

# Specific dashboard type
focus dashboard --type minimal
```

#### Dashboard Types
- `unified` - Complete feature set
- `minimal` - Task-focused view
- `calendar` - Calendar-centric view

---

## Configuration Commands

### `focus config` - Configuration Management
Manage Focus.sh configuration.

#### Subcommands
```bash
# List all configuration
focus config list

# Get specific value
focus config get ai.provider

# Set configuration value
focus config set theme synthwave

# Reset configuration
focus config reset

# Show config file location
focus config path

# Export configuration
focus config export

# Import configuration
focus config import <file>
```

---

### `focus enhanced-config` - Configuration Wizard
Interactive configuration setup wizard.

#### Examples
```bash
# Launch wizard
focus enhanced-config

# Configure specific section
focus enhanced-config --section ai
```

#### Features
- AI provider detection and setup
- Theme selection and preview
- Storage configuration
- Feature enable/disable

---

### `focus theme` - Theme Management
Manage visual themes.

#### Subcommands
```bash
# List available themes
focus theme list

# Set theme
focus theme synthwave

# Preview theme
focus theme preview synthwave

# Save current settings as theme
focus theme save my-theme
```

#### Available Themes
- `synthwave` - Cyberpunk aesthetics (default)
- `light` - Clean, minimal design
- `plain` - Monochrome styling
- Custom themes from `~/.config/focus/themes/`

---

## Advanced Features

### `focus wizard` - Task Creation Wizard
Advanced task creation with Huh forms.

#### Examples
```bash
# Launch wizard
focus wizard

# Wizard for specific category
focus wizard --category work

# Quick wizard mode
focus wizard --quick
```

---

### `focus filter` - Interactive Filtering
Interactive task filtering and search.

#### Examples
```bash
# Launch filter interface
focus filter

# Filter with initial query
focus filter --query "work"

# Filter by date range
focus filter --from 2024-01-01 --to 2024-01-31
```

---

### `focus suggest` - Command Suggestions
Get suggestions for commands and features.

#### Examples
```bash
# General suggestions
focus suggest

# Context-sensitive suggestions
focus suggest --context task-creation

# Feature suggestions
focus suggest --feature calendar
```

---

## Data Management

### `focus export` - Export Data
Export tasks and configuration.

#### Examples
```bash
# Export all tasks
focus export --output tasks.json

# Export filtered tasks
focus export --category work --output work-tasks.json

# Export with format
focus export --format csv --output tasks.csv
```

#### Flags
- `--output, -o` - Output file path
- `--format, -f` - Export format: json, csv, markdown
- `--category, -c` - Export specific category
- `--status, -s` - Export by status

---

### `focus import` - Import Data
Import tasks from external sources.

#### Examples
```bash
# Import from JSON
focus import tasks.json

# Import from CSV
focus import --format csv tasks.csv

# Import with category mapping
focus import --map-category "Work:work" tasks.json
```

---

## Utility Commands

### `focus status` - System Status
Show system and configuration status.

#### Examples
```bash
# Full status report
focus status

# AI provider status
focus status --ai

# Configuration status
focus status --config
```

#### Output Includes
- AI provider connectivity
- Model availability
- Storage locations
- Configuration validation
- Version information

---

### `focus version` - Version Information
Show version and build information.

```bash
focus version
# Focus.sh v1.0.0 (build: abc123)
# Go: 1.21.0
# Platform: linux/amd64
```

---

## Shortcuts and Aliases

Common command shortcuts for quick usage:

```bash
# Quick task addition
focus a "task description"    # alias for add

# Quick list
focus l                       # alias for list

# Quick completion
focus c 5                     # alias for complete

# Quick delete
focus d 5                     # alias for delete

# Unified dashboard
focus u                       # alias for unified

# Configuration
focus config theme synthwave  # shorthand
focus c t synthwave           # abbreviated
```

---

## Output Formats

### Table Format (Default)
```
┌────┬─────────────────────────┬──────────┬──────────┬────────────┐
│ ID │ Task                    │ Category │ Priority │ Due        │
├────┼─────────────────────────┼──────────┼──────────┼────────────┤
│ 1  │ Buy groceries           │ Personal │ Medium   │ 2024-01-25 │
│ 2  │ Complete project        │ Work     │ High     │ 2024-01-20 │
└────┴─────────────────────────┴──────────┴──────────┴────────────┘
```

### List Format
```
1. [Work] Complete project (High) - Due: 2024-01-20
2. [Personal] Buy groceries (Medium) - Due: 2024-01-25
```

### JSON Format
```json
{
  "tasks": [
    {
      "id": 1,
      "description": "Buy groceries",
      "category": "Personal",
      "priority": "Medium",
      "due_date": "2024-01-25",
      "status": "pending"
    }
  ]
}
```

---

## Exit Codes

- `0` - Success
- `1` - General error
- `2` - Command usage error
- `3` - Configuration error
- `4` - AI provider error
- `5` - File system error

---

## Getting Help

For more help with any command:

```bash
# Show command help
focus <command> --help

# Show all commands
focus --help

# Interactive help
focus help

# Get suggestions
focus suggest
```

See the [Troubleshooting Guide](troubleshooting.md) for common issues and solutions.