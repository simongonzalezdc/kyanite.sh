#!/usr/bin/env bash
#
# noise.sh Launcher
# Installs dependencies, builds, and launches the TUI app
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}"
echo "  ♪ noise.sh Launcher"
echo "  ═══════════════════════════════════════"
echo -e "${NC}"

# -----------------------------------------------------------------------------
# Step 1: Check Go installation
# -----------------------------------------------------------------------------
echo -e "${YELLOW}[1/4]${NC} Checking Go installation..."

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo ""
    echo "Please install Go 1.21 or later:"
    echo "  macOS:   brew install go"
    echo "  Linux:   sudo apt install golang-go  (or download from https://go.dev)"
    echo "  Windows: Download from https://go.dev/dl/"
    echo ""
    exit 1
fi

GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | head -1)
echo -e "  ${GREEN}✓${NC} Found $GO_VERSION"

# -----------------------------------------------------------------------------
# Step 2: Download Go dependencies
# -----------------------------------------------------------------------------
echo -e "${YELLOW}[2/4]${NC} Downloading dependencies..."

if [ ! -f "go.sum" ] || [ "go.mod" -nt "go.sum" ]; then
    go mod download
    echo -e "  ${GREEN}✓${NC} Dependencies downloaded"
else
    echo -e "  ${GREEN}✓${NC} Dependencies up to date"
fi

# -----------------------------------------------------------------------------
# Step 3: Build the application
# -----------------------------------------------------------------------------
echo -e "${YELLOW}[3/4]${NC} Building noise.sh..."

BUILD_DIR="bin"
BINARY="$BUILD_DIR/noise"

# Add .exe suffix on Windows
if [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    BINARY="$BUILD_DIR/noise.exe"
fi

mkdir -p "$BUILD_DIR"

# Check if rebuild is needed
NEEDS_BUILD=false
if [ ! -f "$BINARY" ]; then
    NEEDS_BUILD=true
elif [ -n "$(find . -name '*.go' -newer "$BINARY" -not -path './deps/*' 2>/dev/null | head -1)" ]; then
    NEEDS_BUILD=true
fi

if [ "$NEEDS_BUILD" = true ]; then
    echo "  Building..."
    go build -trimpath -ldflags "-s -w" -o "$BINARY" ./cmd/noise
    echo -e "  ${GREEN}✓${NC} Build complete"
else
    echo -e "  ${GREEN}✓${NC} Already built (up to date)"
fi

# -----------------------------------------------------------------------------
# Step 4: Create data directories
# -----------------------------------------------------------------------------
echo -e "${YELLOW}[4/4]${NC} Setting up data directories..."

DATA_DIR="$HOME/.noise"
mkdir -p "$DATA_DIR"
mkdir -p "data/sync/media/voice" 2>/dev/null || true
mkdir -p "data/sync/media/photos" 2>/dev/null || true

echo -e "  ${GREEN}✓${NC} Data directory: $DATA_DIR"

# -----------------------------------------------------------------------------
# Launch
# -----------------------------------------------------------------------------
echo ""
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo -e "${GREEN}  Ready! Launching noise.sh...${NC}"
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo ""
echo -e "  ${BLUE}Tip:${NC} Press 'q' to quit, '?' for help"
echo ""

# Small delay so user can read the message
sleep 1

# Run the app
exec "$BINARY" "$@"
