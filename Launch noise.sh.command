#!/bin/bash
# noise.sh Desktop Launcher for macOS
# Place this file on your Desktop and double-click to launch noise.sh

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Change to the project directory
cd "$SCRIPT_DIR" || {
    osascript -e 'display dialog "Error: Could not find noise.sh directory" buttons {"OK"} default button 1'
    exit 1
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    osascript -e 'display dialog "Error: Go is not installed or not in PATH.\n\nPlease install Go 1.21+ from https://go.dev" buttons {"OK"} default button 1'
    exit 1
fi

# Check if binary exists, if not build it
if [ ! -f "bin/noise" ] && [ ! -f "noise" ]; then
    echo "Building noise.sh..."
    osascript -e 'display notification "Building noise.sh..." with title "noise.sh"'
    
    # Create build directory
    mkdir -p bin
    
    # Build the application
    if go build -trimpath -ldflags "-s -w" -o bin/noise ./cmd/noise; then
        osascript -e 'display notification "Build successful!" with title "noise.sh"'
    else
        osascript -e 'display dialog "Build failed. Please check the Terminal for errors." buttons {"OK"} default button 1'
        exit 1
    fi
fi

# Determine which binary to use
BINARY=""
if [ -f "bin/noise" ]; then
    BINARY="bin/noise"
elif [ -f "noise" ]; then
    BINARY="./noise"
else
    osascript -e 'display dialog "Error: Could not find noise binary" buttons {"OK"} default button 1'
    exit 1
fi

# Launch the application in Terminal
osascript <<EOF
tell application "Terminal"
    activate
    do script "cd '$SCRIPT_DIR' && '$BINARY'"
end tell
EOF

