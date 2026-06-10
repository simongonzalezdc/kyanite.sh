# Windows Launch Guide for noise.sh

This guide explains how to properly launch the noise.sh application on Windows, especially when using PowerShell.

## The PowerShell "Script Not Recognized" Issue

When trying to run `.\scripts\build_and_launch.bat` in PowerShell, you might encounter this error:

```
.\scripts\build_and_launch.bat : The term '.\scripts\build_and_launch.bat' is not recognized as the name of a cmdlet, function, script file, or operable program.
```

### Root Cause

This error occurs because PowerShell treats the command as a PowerShell script rather than a batch file. PowerShell looks for a PowerShell command named `build_and_launch.bat`, but since the extension isn't a known PowerShell script type, it reports the error.

### Solution

Use one of these correct invocation methods:

1. **Using the call operator (`&`)** - Recommended:
   ```powershell
   & .\scripts\build_and_launch.bat
   ```

2. **Using `cmd /c` wrapper**:
   ```powershell
   cmd /c .\scripts\build_and_launch.bat
   ```

3. **Using the provided PowerShell wrapper script**:
   ```powershell
   .\launch_noise.ps1
   ```

## Alternative Launch Methods

### Quick Launch Options

1. **Interactive mode** (shows menu):
   ```powershell
   .\launch_noise.ps1
   ```

2. **Quick mode** (scratch mode):
   ```powershell
   .\launch_noise.ps1 quick
   ```

3. **Debug mode**:
   ```powershell
   .\launch_noise.ps1 debug
   ```

4. **Theme testing mode**:
   ```powershell
   .\launch_noise.ps1 theme
   ```

### Direct Batch Execution

If you prefer to run the batch file directly:

```powershell
# Interactive mode
& .\scripts\build_and_launch.bat

# Quick mode (option 3)
echo 3 | .\scripts\build_and_launch.bat

# Debug mode (option 2)
echo 2 | .\scripts\build_and_launch.bat

# Theme testing mode (option 4)
echo 4 | .\scripts\build_and_launch.bat
```

## Troubleshooting Tips

1. **Verify the script exists**:
   ```powershell
   Get-Item .\scripts\build_and_launch.bat | Format-List FullName, Exists, IsReadOnly
   ```

2. **Check execution policy**:
   ```powershell
   Get-ExecutionPolicy
   ```
   The `RemoteSigned` policy (default) does not block batch files.

3. **If you still get "not recognized" errors**:
   - Verify the path is correct (no extra spaces, correct case)
   - Ensure PowerShell's `PATHEXT` includes `.BAT`:
     ```powershell
     $Env:PATHEXT
     ```
   - If you have an alias that shadows the file, remove it:
     ```powershell
     Remove-Item Alias:build_and_launch
     ```

## Execution Policy Note

The PowerShell execution policy does not affect batch files. It only governs PowerShell (`.ps1`) scripts. The `RemoteSigned` policy allows locally created scripts to run, which is sufficient for our purposes.

## File Permissions

The batch files have proper permissions for execution. You can verify this with:

```powershell
Get-Item .\scripts\build_and_launch.bat | Format-List Mode, IsReadOnly, Attributes
```

Expected output shows `Mode: -a----` and `IsReadOnly: False`, indicating the file is readable and executable.