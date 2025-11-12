#!/bin/bash
# Script to create a desktop launcher for noise.sh on macOS

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
DESKTOP_DIR="$HOME/Desktop"

echo "Creating desktop launcher for noise.sh..."
echo "Project directory: $PROJECT_DIR"
echo "Desktop directory: $DESKTOP_DIR"
echo

# Create the launcher script
LAUNCHER_FILE="$PROJECT_DIR/Launch noise.sh.command"

if [ ! -f "$LAUNCHER_FILE" ]; then
    echo "Error: Launcher file not found at $LAUNCHER_FILE"
    exit 1
fi

# Make it executable
chmod +x "$LAUNCHER_FILE"

# Create a symlink on the desktop
if [ -d "$DESKTOP_DIR" ]; then
    SYMLINK="$DESKTOP_DIR/Launch noise.sh.command"
    if [ -L "$SYMLINK" ] || [ -f "$SYMLINK" ]; then
        echo "Removing existing launcher from desktop..."
        rm "$SYMLINK"
    fi
    
    echo "Creating desktop launcher..."
    ln -s "$LAUNCHER_FILE" "$SYMLINK"
    echo "✅ Desktop launcher created at: $SYMLINK"
    echo
    echo "You can now double-click 'Launch noise.sh.command' on your Desktop to launch the app!"
else
    echo "Warning: Desktop directory not found at $DESKTOP_DIR"
    echo "You can manually copy the launcher to your Desktop:"
    echo "  cp '$LAUNCHER_FILE' ~/Desktop/"
fi

