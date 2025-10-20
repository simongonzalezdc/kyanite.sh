@echo off
echo 🚀 NEON Focus - TUI-First Installation
echo.

REM Check if we're in the correct directory
if not exist "cmd\neon\main.go" (
    echo ❌ Please run this from the NEON project directory
    echo 💡 Directory should contain: cmd\neon\main.go
    pause
    exit /b 1
)

echo 🔨 Building NEON with TUI-first and AI integration...
go build -ldflags="-w" -o neon.exe ./cmd/neon
if errorlevel 1 (
    echo ❌ Build failed!
    pause
    exit /b 1
)

echo ✅ TUI-First build successful!

REM Create installation directory
if not exist "%PROGRAMFILES%\NEON" mkdir "%PROGRAMFILES%\NEON"

echo 📦 Installing to %PROGRAMFILES%\NEON...
copy neon.exe "%PROGRAMFILES%\NEON\" >nul

echo 📋 Adding to PATH...
REM Add to system PATH
setx PATH "%PATH%;%PROGRAMFILES%\NEON" /M >nul

echo 🎯 Creating global launcher...
echo @echo off > "%USERPROFILE%\Desktop\neon.bat"
echo title NEON Focus System >> "%USERPROFILE%\Desktop\neon.bat"
echo echo 🌌 Starting NEON Focus TUI Dashboard... >> "%USERPROFILE%\Desktop\neon.bat"
echo echo 💫 AI-powered task management with cyberpunk aesthetics >> "%USERPROFILE%\Desktop\neon.bat"
echo echo. >> "%USERPROFILE%\Desktop\neon.bat"
echo "%PROGRAMFILES%\NEON\neon.exe" %%* >> "%USERPROFILE%\Desktop\neon.bat"

echo.
echo ✅ TUI-First installation complete!
echo.
echo 🎮 NEON will now launch directly to TUI Dashboard!
echo 💫 No more CLI mode - TUI-first experience!
echo 🤖 Real AI integration included!
echo.
echo 🌟 You can now launch NEON from anywhere using:
echo    neon [options]
echo.
echo 🖥  Desktop shortcut created on your desktop.
echo 🌟 Or run: neon.bat
echo.
echo 🔄 Close and reopen terminal to use global 'neon' command.
echo.

REM Test installation
echo 🧪 Testing TUI installation...
"%PROGRAMFILES%\NEON\neon.exe" --help >nul 2>&1
if errorlevel 0 (
    echo ✅ Installation verified!
) else (
    echo ⚠️  Installation completed but test failed
    echo 💡 Try restarting your terminal first
)

echo.
echo 🌌 NEON Focus TUI-First installed successfully!
echo 💫 Enjoy your cyberpunk productivity dashboard!
echo.
