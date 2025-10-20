#!/bin/bash
# noise.sh Build and Launch Script for Linux/macOS
# This script builds and launches the noise.sh application with theme testing options

set -e

echo "========================================"
echo "noise.sh Build and Launch Script"
echo "========================================"
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH"
    echo "Please install Go 1.21+ and try again"
    exit 1
fi

echo "Go installation found:"
go version
echo

# Clean previous builds
echo "Cleaning previous builds..."
rm -rf bin
rm -f noise

# Create build directory
echo "Creating build directory..."
mkdir -p bin

# Build the application
echo
echo "Building noise.sh..."
go build -trimpath -ldflags "-s -w" -o bin/noise ./cmd/noise

echo "Build successful!"
echo

# Copy to root for easier access
cp bin/noise noise

# Display launch options
echo "========================================"
echo "LAUNCH OPTIONS"
echo "========================================"
echo
echo "1. Launch normally"
echo "2. Launch with debug mode"
echo "3. Launch in quick mode (scratch mode)"
echo "4. Launch with theme testing"
echo "5. Exit"
echo

read -p "Select launch option (1-5): " choice

case $choice in
    1)
        echo
        echo "Launching noise.sh..."
        ./noise
        ;;
    2)
        echo
        echo "Launching noise.sh with debug mode..."
        ./noise --debug
        ;;
    3)
        echo
        echo "Launching noise.sh in quick mode..."
        ./noise quick
        ;;
    4)
        echo
        echo "========================================"
        echo "THEME TESTING MODE"
        echo "========================================"
        echo
        echo "Testing all 10 Kyanite themes..."
        echo
        echo "Available themes:"
        echo "1. Monochrome"
        echo "2. Amber Night (default)"
        echo "3. Twilight Mist"
        echo "4. Indigo Depths"
        echo "5. Forest Path"
        echo "6. Clay Earth"
        echo "7. Iron Forge"
        echo "8. Sunlight"
        echo "9. Cyan Wave"
        echo "10. Electric Rose"
        echo
        echo "Use Ctrl+Shift+T to cycle through themes"
        echo "Press F1 for help while in the application"
        echo
        read -p "Press Enter to launch with theme testing..."
        ./noise --debug
        ;;
    5)
        echo
        echo "Exiting..."
        exit 0
        ;;
    *)
        echo "Invalid choice. Exiting..."
        exit 1
        ;;
esac

echo
echo "Application closed."