@echo off
echo 🔧 Fixing NEON Desktop Shortcut...
echo.

REM Create a proper PowerShell shortcut
echo 🎯 Creating Windows shortcut...
powershell -Command "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut('%USERPROFILE%\Desktop\NEON Focus.lnk'); $Shortcut.TargetPath = '%PROGRAMFILES%\NEON\neon.exe'; $Shortcut.WorkingDirectory = '%USERPROFILE%'; $Shortcut.Description = 'NEON Focus - AI-Powered Task Manager'; $Shortcut.Save()"

REM Create a simple batch file as backup
echo 📝 Creating batch file...
echo @echo off > "%USERPROFILE%\Desktop\Start NEON.bat"
echo title NEON Focus System >> "%USERPROFILE%\Desktop\Start NEON.bat"
echo echo 🚀 Starting NEON Focus... >> "%USERPROFILE%\Desktop\Start NEON.bat"
echo echo 🤖 Checking AI setup... >> "%USERPROFILE%\Desktop\Start NEON.bat"
echo "%PROGRAMFILES%\NEON\neon.exe" %%* >> "%USERPROFILE%\Desktop\Start NEON.bat"
echo if errorlevel 1 pause >> "%USERPROFILE%\Desktop\Start NEON.bat"

echo.
echo ✅ Desktop shortcuts created!
echo.
echo 🎮 Try these desktop shortcuts:
echo    1. NEON Focus.lnk (Windows shortcut)
echo    2. Start NEON.bat (Batch file)
echo.

REM Test if shortcut was created
if exist "%USERPROFILE%\Desktop\NEON Focus.lnk" (
    echo ✅ Windows shortcut created successfully!
) else (
    echo ⚠️  Windows shortcut creation may have failed
)

if exist "%USERPROFILE%\Desktop\Start NEON.bat" (
    echo ✅ Batch file created successfully!
) else (
    echo ⚠️  Batch file creation may have failed
)

echo.
echo 🧪 Testing NEON installation...
"%PROGRAMFILES%\NEON\neon.exe" --help >nul 2>&1
if errorlevel 0 (
    echo ✅ NEON installation verified!
    echo.
    echo 🎉 Ready to use your desktop shortcuts!
) else (
    echo ⚠️  NEON test failed
    echo 💡 Try running: "%PROGRAMFILES%\NEON\neon.exe" manually
)

pause
