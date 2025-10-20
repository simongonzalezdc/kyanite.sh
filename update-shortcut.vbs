Set objShell = CreateObject("WScript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

REM Use TUI-first version in current directory
neonPath = objShell.CurrentDirectory & "\neon-tui.exe"

REM Update desktop shortcut
Set shortcut = objShell.CreateShortcut(objShell.SpecialFolders("Desktop") & "\NEON Focus.lnk")
shortcut.TargetPath = neonPath
shortcut.Description = "NEON Focus - TUI Dashboard with AI Integration"
shortcut.Save()

MsgBox "✅ NEON Focus shortcut updated to TUI-first version!", "NEON Focus", 0
