@echo off
setlocal

REM NEON LAUNCHER - Works from any directory
REM Automatically finds NEON installation and launches

echo.

REM Try to find NEON in common locations
set NEON_EXE=

REM Check if we're in the project directory
if exist "%~dp0neon.exe" (
    set NEON_EXE=%~dp0neon.exe
    goto :found
)

REM Check current directory
if exist "%CD%\neon.exe" (
    set NEON_EXE=%CD%\neon.exe
    goto :found
)

REM Check common installation paths
if exist "%USERPROFILE%\Downloads\neon.exe" (
    set NEON_EXE=%USERPROFILE%\Downloads\neon.exe
    goto :found
)

if exist "%USERPROFILE%\Desktop\neon.exe" (
    set NEON_EXE=%USERPROFILE%\Desktop\neon.exe
    goto :found
)

if exist "%PROGRAMFILES%\NEON\neon.exe" (
    set NEON_EXE=%PROGRAMFILES%\NEON\neon.exe
    goto :found
)

REM Search in common development folders
for %%d in ("%USERPROFILE%\Projects\*" "%USERPROFILE%\Desktop\*" "%USERPROFILE%\Documents\*") do (
    if exist "%%d\neon.exe" (
        set NEON_EXE=%%d\neon.exe
        goto :found
    )
)

echo.
echo 🤖 NEON not found! Let's set it up for you...
echo.

REM Offer to copy from current project if it exists
if exist "%~dp0neon.exe" (
    echo 💾 Found NEON in project folder. Copying to Desktop...
    copy "%~dp0neon.exe" "%USERPROFILE%\Desktop\" >nul
    set NEON_EXE=%USERPROFILE%\Desktop\neon.exe
    echo ✅ Copied to Desktop!
    goto :found
)

if exist "%~dp0neon-ai-fixed.exe" (
    echo 💾 Found NEON AI version in project folder. Copying to Desktop...
    copy "%~dp0neon-ai-fixed.exe" "%USERPROFILE%\Desktop\neon.exe" >nul
    set NEON_EXE=%USERPROFILE%\Desktop\neon.exe
    echo ✅ Copied to Desktop!
    goto :found
)

echo.
echo ❌ NEON not found anywhere!
echo 💡 Please run one of these commands from your project directory:
echo    1. go build -o neon.exe ./cmd/neon
echo    2. go build -o neon-ai-fixed.exe ./cmd/neon (AI version)
echo    3. Or re-run this launcher from the project folder
echo.
pause
exit /b 1

:found
echo 🚀 Launching NEON Focus System...
echo 📍 Found at: %NEON_EXE%
echo.

REM Start NEON
"%NEON_EXE%" %*

echo.
echo 🌌 NEON session ended. Thanks for using NEON Focus!
timeout /t 2 >nul
exit /b 0
