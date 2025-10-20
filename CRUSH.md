# CRUSH CLI Coding Lab

## Purpose
This directory is a coding lab for experimenting with AI-powered CLI tools using Charm libraries. Focus on learning development patterns with local LLMs (Ollama) and OpenRouter fallback.

## Essential Commands
- Build: `go build ./...`
- Lint: `golangci-lint run`
- Test: `go test ./...`
- Single test: `go test -run TestName ./path/to/package`
- Run CLI: `go run cmd/crush/main.go`

## AI Integration Requirements
- Primary: Local LLMs via Ollama
- Fallback: OpenRouter API (emergency only)
- Always check local model availability first
- Graceful degradation when models unavailable
- Cache responses to minimize API calls

## Learning Projects
1. **AI-Powered To-Do Assistant** (Bubble Tea + Huh + Ollama)
   - Natural language task entry and categorization
   - Smart task prioritization with local LLM
   - Context-aware suggestions

2. **Pomodoro Timer** (Bubble Tea + Lip Gloss)
   - Beautiful TUI-based Pomodoro timer
   - Work/break session tracking
   - Task-focused productivity sessions

2. **Intelligent CLI Dashboard** (Lip Gloss + Bubbles + Ollama)
   - AI-enhanced system monitoring insights
   - Natural language querying of system data
   - Predictive resource usage warnings

3. **Smart SSH File Assistant** (Wish + Bubble Tea + Ollama)
   - Natural language file operations
   - AI-powered file search and organization
   - Context-aware command suggestions

4. **AI Markdown Assistant** (Glamour + Bubble Tea + Ollama)
   - Natural language document querying
   - Smart summarization and tagging
   - Content generation and editing suggestions

5. **Adaptive Progress System** (Harmonica + Bubble Tea + Ollama)
   - AI-estimated time completion
   - Dynamic progress visualization adjustments
   - Context-aware status updates

6. **Intelligent Logger** (Log + Ollama)
   - AI-powered log analysis and pattern detection
   - Natural language log querying
   - Automated anomaly detection

## Development Workflow
- Make small, focused changes
- Test frequently with manual runs
- Commit often with descriptive messages
- Use feature branches for experimental work
- Always test both local and fallback AI modes

## Code Style
- Follow idiomatic Go conventions
- Use descriptive names for all identifiers
- Handle errors explicitly
- Write concise comments for exported functions/types
- Format code with `go fmt`
- Import grouping: stdlib, third-party, internal
- Prefer early returns to reduce nesting
- Separate CLI logic from core business logic
- Isolate AI integration points for easy testing

## Project Structure
- cmd/: CLI entry points
- internal/: Private application logic
- pkg/ai/: AI integration utilities
- pkg/cli/: Charm-based components
- test/: Integration and utility tests

## Experimental Development
- Try new patterns freely
- Document interesting approaches
- Keep experiments isolated
- Refactor working experiments into proper structure
- Always implement both local and remote AI paths