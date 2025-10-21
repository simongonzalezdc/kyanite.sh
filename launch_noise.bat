@echo off
REM Simple launch script for noise.sh on Windows
REM This is an alternative to scripts\build_and_launch.bat

echo ========================================
echo noise.sh Quick Launcher
echo ========================================
echo.

REM Check if noise.exe exists in current directory
if exist noise.exe (
    echo Found noise.exe in current directory
    echo Launching noise.sh...
    echo.
    noise.exe %*
    goto :end
)

REM Check if noise.exe exists in bin directory
if exist bin\noise.exe (
    echo Found noise.exe in bin directory
    echo Launching noise.sh...
    echo.
    bin\noise.exe %*
    goto :end
)

REM If not found, try to build it
echo noise.exe not found. Attempting to build...
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.21+ and try again
    pause
    exit /b 1
)

REM Create bin directory if it doesn't exist
if not exist bin mkdir bin

REM Build the application
echo Building noise.sh...
go build -trimpath -ldflags "-s -w" -o bin\noise.exe ./cmd/noise

if errorlevel 1 (
    echo ERROR: Build failed
    pause
    exit /b 1
)

echo Build successful!
echo.
echo Launching noise.sh...
bin\noise.exe %*

:end
echo.
echo Application closed.
pause