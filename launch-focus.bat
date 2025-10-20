@echo off
title Kyanite focus.sh Launcher
color 0A

echo.
echo ╔════════════════════════════════════════════════════════════════╗
echo ║                    🌌 KYANITE FOCUS.SH LAUNCHER                ║
echo ║                      Professional Task Manager                   ║
echo ╚════════════════════════════════════════════════════════════════╝
echo.

:menu
echo Choose an option:
echo.
echo [1] Launch Dashboard (TUI Interface)
echo [2] Add New Task  
echo [3] List All Tasks
echo [4] View Help
echo [5] Change Theme
echo [6] Open Journal
echo [7] Chat with AI Assistant
echo [8] Exit
echo.
set /p choice="Enter your choice (1-8): "

if "%choice%"=="1" goto dashboard
if "%choice%"=="2" goto addtask
if "%choice%"=="3" goto listtasks
if "%choice%"=="4" goto help
if "%choice%"=="5" goto theme
if "%choice%"=="6" goto journal
if "%choice%"=="7" goto chat
if "%choice%"=="8" goto exit

echo Invalid choice. Please try again.
echo.
goto menu

:dashboard
echo.
echo 🚀 Launching focus.sh Dashboard...
echo.
cd /d "%~dp0"
focus.exe dashboard
goto menu

:addtask
echo.
set /p task="Enter your task: "
if "%task%"=="" goto menu
echo.
echo ✅ Adding task: %task%
cd /d "%~dp0"
focus.exe add "%task%"
pause
goto menu

:listtasks
echo.
echo 📋 Your Tasks:
echo.
cd /d "%~dp0"
focus.exe list
pause
goto menu

:help
echo.
echo 📚 focus.sh Help:
echo.
cd /d "%~dp0"
focus.exe --help
pause
goto menu

:theme
echo.
echo 🎨 Available Themes:
echo.
cd /d "%~dp0"
focus.exe theme --help
echo.
set /p themechoice="Enter theme name: 
cd /d "%~dp0"
focus.exe theme "%themechoice%"
pause
goto menu

:journal
echo.
echo 📝 Journal Options:
echo.
echo [N] New Entry
echo [L] List Entries  
echo [V] View Entry
echo [S] Search Entries
echo [B] Back to Main Menu
echo.
set /p journalchoice="Choose option: "

if /i "%journalchoice%"=="N" goto journalnew
if /i "%journalchoice%"=="L" goto journallist
if /i "%journalchoice%"=="V" goto journalview
if /i "%journalchoice%"=="S" goto journalsearch
if /i "%journalchoice%"=="B" goto menu

goto journal

:journalnew
echo.
cd /d "%~dp0"
focus.exe journal new
pause
goto menu

:journallist
echo.
cd /d "%~dp0"
focus.exe journal list
pause
goto menu

:journalview
echo.
set /p journaldate="Enter date (YYYY-MM-DD) or leave blank for today: 
cd /d "%~dp0"
if "%journaldate%"=="" (
    focus.exe journal view
) else (
    focus.exe journal view --date "%journaldate%"
)
pause
goto menu

:journalsearch
echo.
set /p searchquery="Enter search query: 
cd /d "%~dp0"
focus.exe journal search --query "%searchquery%"
pause
goto menu

:chat
echo.
echo 🤖 Starting AI Assistant...
echo Type 'exit' to return to menu.
echo.
cd /d "%~dp0"
focus.exe chat
pause
goto menu

:exit
echo.
echo 👋 Thanks for using Kyanite focus.sh!
echo.
timeout /t 2 >nul
exit