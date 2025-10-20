@echo off
echo Creating Kyanite focus.sh Desktop Shortcut...

:: Set variables
set DESKTOP=%USERPROFILE%\Desktop
set ICON_PATH=%~dp0focus.exe
set WORKING_DIR=%~dp0
set SHORTCUT_NAME=Kyanite focus.lnk

:: Create VBS script to make shortcut
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%TEMP%\MakeShortcut.vbs"
echo sLinkFile = "%DESKTOP%\%SHORTCUT_NAME%" >> "%TEMP%\MakeShortcut.vbs"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%TEMP%\MakeShortcut.vbs"
echo oLink.TargetPath = "%WORKING_DIR%focus.exe" >> "%TEMP%\MakeShortcut.vbs"
echo oLink.WorkingDirectory = "%WORKING_DIR%" >> "%TEMP%\MakeShortcut.vbs"
echo oLink.Description = "Kyanite focus.sh - Professional Task Manager" >> "%TEMP%\MakeShortcut.vbs"
echo oLink.Save >> "%TEMP%\MakeShortcut.vbs"

:: Execute the VBS script
cscript //nologo "%TEMP%\MakeShortcut.vbs"

:: Clean up
del "%TEMP%\MakeShortcut.vbs"

echo.
echo ✅ Desktop shortcut created: "Kyanite focus.lnk"
echo.
echo Double-click the desktop shortcut to launch focus.sh TUI directly!
echo.
pause