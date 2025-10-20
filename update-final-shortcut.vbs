Set objShell = CreateObject("WScript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

REM Use final version with synthwave personality
neonPath = objShell.CurrentDirectory & "\neon-final.exe"

REM Update desktop shortcut
Set shortcut = objShell.CreateShortcut(objShell.SpecialFolders("Desktop") & "\NEON Focus.lnk")
shortcut.TargetPath = neonPath
shortcut.Description = "NEON Focus - Synthwave Cyberpunk AI + Fixed Calendar + Notes"
shortcut.Save()

MsgBox "✅ NEON Focus updated to final synthwave version!", "NEON Focus", 0
