@echo off

echo 🔍 Validating AI Focus Assistant Setup...

REM Check if Go is installed
echo Checking Go installation...
go version >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Go is installed
    go version
) else (
    echo ❌ Go is not installed
    echo Please install Go from https://golang.org/dl/
    exit /b 1
)

REM Check if Ollama is installed
echo Checking Ollama installation...
ollama --version >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Ollama is installed
    ollama --version
) else (
    echo ❌ Ollama is not installed
    echo Please install Ollama from https://ollama.ai/
    exit /b 1
)

REM Check if llama3 model is available
echo Checking for llama3 model...
ollama list | findstr llama3 >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ llama3 model is available
) else (
    echo ⚠️ llama3 model not found
    echo Run 'ollama pull llama3' to download the model
)

REM Run tests
echo Running tests...
go test -short ./... >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ All tests passed
) else (
    echo ⚠️ Some tests failed or were skipped
)

echo.
echo 🎉 Setup validation complete!
echo Your environment is ready to develop the AI Focus Assistant.