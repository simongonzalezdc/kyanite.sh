@echo off
REM noise.sh Launcher for Windows
REM Installs dependencies, builds, and launches the TUI app

setlocal EnableDelayedExpansion

echo.
echo   # noise.sh Launcher
echo   ===================================
echo.

cd /d "%~dp0"

REM -----------------------------------------------------------------------------
REM Step 1: Check Go installation
REM -----------------------------------------------------------------------------
echo [1/4] Checking Go installation...

where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo.
    echo ERROR: Go is not installed.
    echo.
    echo Please install Go 1.21 or later:
    echo   Download from https://go.dev/dl/
    echo.
    pause
    exit /b 1
)

for /f "tokens=3" %%v in ('go version') do set GO_VERSION=%%v
echo   [OK] Found %GO_VERSION%

REM -----------------------------------------------------------------------------
REM Step 2: Download Go dependencies
REM -----------------------------------------------------------------------------
echo [2/4] Downloading dependencies...

go mod download
if %ERRORLEVEL% neq 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)
echo   [OK] Dependencies ready

REM -----------------------------------------------------------------------------
REM Step 3: Build the application
REM -----------------------------------------------------------------------------
echo [3/4] Building noise.sh...

if not exist "bin" mkdir bin

set BINARY=bin\noise.exe

REM Check if rebuild needed
set NEEDS_BUILD=0
if not exist "%BINARY%" set NEEDS_BUILD=1

if %NEEDS_BUILD%==1 (
    echo   Building...
    go build -trimpath -ldflags "-s -w" -o %BINARY% ./cmd/noise
    if %ERRORLEVEL% neq 0 (
        echo ERROR: Build failed
        pause
        exit /b 1
    )
    echo   [OK] Build complete
) else (
    echo   [OK] Already built
)

REM -----------------------------------------------------------------------------
REM Step 4: Create data directories
REM -----------------------------------------------------------------------------
echo [4/4] Setting up data directories...

if not exist "%USERPROFILE%\.noise" mkdir "%USERPROFILE%\.noise"
if not exist "data\sync\media\voice" mkdir "data\sync\media\voice"
if not exist "data\sync\media\photos" mkdir "data\sync\media\photos"

echo   [OK] Data directories ready

REM -----------------------------------------------------------------------------
REM Launch
REM -----------------------------------------------------------------------------
echo.
echo ===================================
echo   Ready! Launching noise.sh...
echo ===================================
echo.
echo   Tip: Press 'q' to quit, '?' for help
echo.

timeout /t 2 >nul

REM Run the app
%BINARY% %*
