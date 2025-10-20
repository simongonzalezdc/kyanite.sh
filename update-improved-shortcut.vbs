Set objShell = CreateObject("WScript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

REM Use improved version in current directory
neonPath = objShell.CurrentDirectory & "\neon-improved.exe"

REM Update desktop shortcut
Set shortcut = objShell.CreateShortcut(objShell.SpecialFolders("Desktop") & "\NEON Focus.lnk")
shortcut.TargetPath = neonPath
shortcut.Description = "NEON Focus - Improved TUI with Better Chat + Calendar + Notes"
shortcut.Save()

MsgBox "✅ NEON Focus shortcut updated to improved version!", "NEON Focus", 0
