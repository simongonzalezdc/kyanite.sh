@echo off
title NEON Focus System
echo.
echo 🚀 Starting NEON Focus with AI integration...
echo.

REM Try program files first
if exist "%PROGRAMFILES%\NEON\neon.exe" (
    echo ✅ Found installed version, launching...
    "%PROGRAMFILES%\NEON\neon.exe" %*
    goto :end
)

echo 📂 Trying project versions...

REM Try local versions
if exist "neon.exe" (
    echo ✅ Found local version, launching...
    neon.exe %*
    goto :end
)

if exist "neon-ai-fixed.exe" (
    echo ✅ Found AI version, launching...
    neon-ai-fixed.exe %*
    goto :end
)

echo ❌ NEON not found!
echo 💡 Run: install-neon.bat for system installation
echo 💫 Or copy neon.exe to this folder

:end
echo.
echo 🌌 NEON session ended.
timeout /t 3 >nul
