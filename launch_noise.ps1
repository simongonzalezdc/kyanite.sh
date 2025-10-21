#!/usr/bin/env pwsh

# PowerShell wrapper script for noise.sh build and launch
# This script properly invokes the batch file from PowerShell

param(
    [Parameter(Position=0)]
    [string]$Mode = ""
)

# Get the directory of this script
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$BatchScript = Join-Path $ScriptDir "scripts\build_and_launch.bat"

# Check if the batch script exists
if (-not (Test-Path $BatchScript)) {
    Write-Error "Build and launch script not found: $BatchScript"
    exit 1
}

# Execute the batch script with proper invocation
switch ($Mode.ToLower()) {
    "quick" {
        # Launch in quick mode (option 3)
        & cmd /c "echo 3 | $BatchScript"
    }
    "debug" {
        # Launch in debug mode (option 2)
        & cmd /c "echo 2 | $BatchScript"
    }
    "theme" {
        # Launch with theme testing (option 4)
        & cmd /c "echo 4 | $BatchScript"
    }
    "normal" {
        # Launch normally (option 1)
        & cmd /c "echo 1 | $BatchScript"
    }
    default {
        # Run interactively
        & $BatchScript
    }
}