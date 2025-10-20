@echo off
title focus.sh System
echo.
echo 🚀 Starting focus.sh with AI integration...
echo.

REM Try program files first
if exist "%PROGRAMFILES%\focus\focus.exe" (
    echo ✅ Found installed version, launching...
    "%PROGRAMFILES%\focus\focus.exe" %*
    goto :end
)

echo 📂 Trying project versions...

REM Try local versions
if exist "focus.exe" (
    echo ✅ Found local version, launching...
    focus.exe %*
    goto :end
)

if exist "focus-ai-fixed.exe" (
    echo ✅ Found AI version, launching...
    focus-ai-fixed.exe %*
    goto :end
)

echo ❌ focus.sh not found!
echo 💡 Run: install-focus.bat for system installation
echo 💫 Or copy focus.exe to this folder

:end
echo.
echo 🌌 focus.sh session ended.
timeout /t 3 >nul
