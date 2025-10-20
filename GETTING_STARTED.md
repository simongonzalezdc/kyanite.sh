# Getting Started Guide

## Prerequisites Installation

### Step 1: Install Go
1. Visit [https://golang.org/dl/](https://golang.org/dl/)
2. Download "go1.21.x.windows-amd64.msi"
3. Run the installer as Administrator
4. Follow the installation wizard
5. Restart your command prompt
6. Verify with: `go version`

### Step 2: Install Ollama
1. Visit [https://ollama.ai/](https://ollama.ai/)
2. Click "Download for Windows"
3. Run the installer
4. Follow the installation wizard
5. Restart your command prompt
6. Verify with: `ollama --version`

### Step 3: Download AI Model
1. Open Command Prompt
2. Run: `ollama pull llama3`
3. Wait for download to complete (4.7GB)
4. Verify with: `ollama list`

## Project Setup

### Step 4: Initialize Project
1. Open Command Prompt in project directory
2. Run: `go mod tidy`
3. This installs all required dependencies

### Step 5: Verify Setup
1. Run: `verify_setup.bat`
2. This checks all components are working

## Testing

### Step 6: Run Core Tests
1. Run: `test_core.bat`
2. This builds and tests basic functionality

### Step 7: Run All Tests
1. Run: `go test ./...`
2. This runs all unit and integration tests

## Usage

### Step 8: Build Application
1. Run: `build.bat`
2. This creates `todo.exe`

### Step 9: Test Commands
```bash
# Add a task with natural language
todo.exe add "Buy milk and bread tomorrow"

# List all tasks
todo.exe list

# List only active tasks
todo.exe list --filter=active

# Mark a task as complete (use actual task ID from list)
todo.exe complete <task_id>

# Delete a task (use actual task ID from list)
todo.exe delete <task_id>
```

## Troubleshooting

### Common Issues

1. **"go is not recognized"**
   - Restart Command Prompt after Go installation
   - Check if Go was added to PATH

2. **"ollama is not recognized"**
   - Restart Command Prompt after Ollama installation
   - Check if Ollama service is running

3. **AI Parsing Issues**
   - Ensure Ollama is running: `ollama serve`
   - Check model is loaded: `ollama list`
   - Try: `ollama run llama3 "test prompt"`

4. **Build Errors**
   - Run: `go mod tidy` to update dependencies
   - Check Go version is 1.21+

### Getting Help

1. Check logs in `~/.focus/tasks.json`
2. Run tests with verbose output: `go test -v ./...`
3. File issues on project repository