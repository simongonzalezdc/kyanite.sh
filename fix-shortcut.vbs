Set objShell = CreateObject("WScript.Shell")
Set fso = CreateObject("Scripting.FileSystemObject")

REM Check if NEON exists in program files
neonPath = ""
If fso.FileExists(objShell.ExpandEnvironmentStrings("%PROGRAMFILES%\NEON\neon.exe")) Then
    neonPath = objShell.ExpandEnvironmentStrings("%PROGRAMFILES%\NEON\neon.exe")
Else
    REM Use current directory
    neonPath = objShell.CurrentDirectory & "\neon-ai-fixed.exe"
End If

REM Create desktop shortcut
Set shortcut = objShell.CreateShortcut(objShell.SpecialFolders("Desktop") & "\NEON Focus.lnk")
shortcut.TargetPath = neonPath
shortcut.Description = "NEON Focus - AI-Powered Task Manager"
shortcut.Save

MsgBox "✅ NEON Focus shortcut created on your Desktop!", "NEON Focus", 0
