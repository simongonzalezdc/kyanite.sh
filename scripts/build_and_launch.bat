@echo off
REM noise.sh Build and Launch Script for Windows
REM This script builds and launches the noise.sh application with theme testing options

echo ========================================
echo noise.sh Build and Launch Script
echo ========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.21+ and try again
    pause
    exit /b 1
)

echo Go installation found:
go version
echo.

REM Clean previous builds
echo Cleaning previous builds...
if exist bin (
    rmdir /s /q bin
)
if exist noise.exe (
    del noise.exe
)

REM Create build directory
echo Creating build directory...
if not exist bin mkdir bin

REM Build the application
echo.
echo Building noise.sh...
go build -trimpath -ldflags "-s -w" -o bin\noise.exe ./cmd/noise

if errorlevel 1 (
    echo ERROR: Build failed
    pause
    exit /b 1
)

echo Build successful!
echo.

REM Copy to root for easier access
copy bin\noise.exe noise.exe >nul

REM Display launch options
echo ========================================
echo LAUNCH OPTIONS
echo ========================================
echo.
echo 1. Launch normally
echo 2. Launch with debug mode
echo 3. Launch in quick mode (scratch mode)
echo 4. Launch with theme testing
echo 5. Exit
echo.

set /p choice="Select launch option (1-5): "

if "%choice%"=="1" (
    echo.
    echo Launching noise.sh...
    noise.exe
) else if "%choice%"=="2" (
    echo.
    echo Launching noise.sh with debug mode...
    noise.exe --debug
) else if "%choice%"=="3" (
    echo.
    echo Launching noise.sh in quick mode...
    noise.exe quick
) else if "%choice%"=="4" (
    echo.
    echo ========================================
    echo THEME TESTING MODE
    echo ========================================
    echo.
    echo Testing all 10 Kyanite themes...
    echo.
    echo Available themes:
    echo 1. Monochrome
    echo 2. Amber Night (default)
    echo 3. Twilight Mist
    echo 4. Indigo Depths
    echo 5. Forest Path
    echo 6. Clay Earth
    echo 7. Iron Forge
    echo 8. Sunlight
    echo 9. Cyan Wave
    echo 10. Electric Rose
    echo.
    echo Use Ctrl+Shift+T to cycle through themes
    echo Press F1 for help while in the application
    echo.
    set /p continue="Press Enter to launch with theme testing..."
    noise.exe --debug
) else if "%choice%"=="5" (
    echo.
    echo Exiting...
    exit /b 0
) else (
    echo Invalid choice. Exiting...
    exit /b 1
)

echo.
echo Application closed.
pause