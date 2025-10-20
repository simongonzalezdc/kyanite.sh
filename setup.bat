@echo off

REM Setup script for AI Focus application on Windows

echo 🔧 Setting up AI Focus Assistant...

REM Check if Go is installed
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go 1.21+ from https://golang.org/dl/
    exit /b 1
)

echo ✅ Go is installed

REM Check if Ollama is installed
where ollama >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Ollama is not installed. Please install Ollama from https://ollama.ai/
    exit /b 1
)

echo ✅ Ollama is installed

REM Pull required model
echo 📥 Pulling llama3 model...
ollama pull llama3

REM Install Go dependencies
echo 📦 Installing Go dependencies...
go mod tidy

echo ✅ Setup complete!
echo.
echo To build the application, run:
echo   go build -o focus.exe cmd/focus/main.go
echo.
echo To run directly, use:
echo   go run cmd/focus/main.go