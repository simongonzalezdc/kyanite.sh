#!/usr/bin/env bash
#
# noise.sh Full Stack Launcher
# Launches both the TUI app and the PWA companion
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo -e "${BLUE}"
echo "  ♪ noise.sh Full Stack Launcher"
echo "  ═══════════════════════════════════════"
echo -e "${NC}"

# -----------------------------------------------------------------------------
# Check prerequisites
# -----------------------------------------------------------------------------
echo -e "${YELLOW}Checking prerequisites...${NC}"

# Check Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "  Install: brew install go (macOS) or visit https://go.dev"
    exit 1
fi
echo -e "  ${GREEN}✓${NC} Go found"

# Check Node.js
if ! command -v node &> /dev/null; then
    echo -e "${RED}Error: Node.js is not installed.${NC}"
    echo "  Install: brew install node (macOS) or visit https://nodejs.org"
    exit 1
fi
echo -e "  ${GREEN}✓${NC} Node.js found"

# Check npm
if ! command -v npm &> /dev/null; then
    echo -e "${RED}Error: npm is not installed.${NC}"
    exit 1
fi
echo -e "  ${GREEN}✓${NC} npm found"

# -----------------------------------------------------------------------------
# Build TUI
# -----------------------------------------------------------------------------
echo ""
echo -e "${YELLOW}Building TUI app...${NC}"

go mod download
mkdir -p bin
go build -trimpath -ldflags "-s -w" -o bin/noise ./cmd/noise
echo -e "  ${GREEN}✓${NC} TUI built"

# -----------------------------------------------------------------------------
# Setup PWA
# -----------------------------------------------------------------------------
echo ""
echo -e "${YELLOW}Setting up PWA...${NC}"

if [ -d "pwa" ]; then
    cd pwa
    if [ ! -d "node_modules" ]; then
        echo "  Installing dependencies..."
        npm install --silent
    fi
    echo -e "  ${GREEN}✓${NC} PWA ready"
    cd ..
else
    echo -e "  ${YELLOW}!${NC} PWA directory not found, skipping"
fi

# -----------------------------------------------------------------------------
# Create data directories
# -----------------------------------------------------------------------------
mkdir -p "$HOME/.noise"
mkdir -p "data/sync/media/voice" 2>/dev/null || true
mkdir -p "data/sync/media/photos" 2>/dev/null || true

# -----------------------------------------------------------------------------
# Launch
# -----------------------------------------------------------------------------
echo ""
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo -e "${GREEN}  Starting noise.sh ecosystem...${NC}"
echo -e "${GREEN}═══════════════════════════════════════${NC}"
echo ""

# Start PWA dev server in background
if [ -d "pwa" ]; then
    echo -e "${CYAN}Starting PWA server...${NC}"
    cd pwa
    npm run dev > /dev/null 2>&1 &
    PWA_PID=$!
    cd ..
    
    # Wait for server to start
    sleep 3
    
    # Get local IP for mobile access
    LOCAL_IP=$(ipconfig getifaddr en0 2>/dev/null || hostname -I 2>/dev/null | awk '{print $1}' || echo "localhost")
    
    echo -e "  ${GREEN}✓${NC} PWA running:"
    echo -e "      Local:   ${CYAN}http://localhost:3000${NC}"
    echo -e "      Network: ${CYAN}http://${LOCAL_IP}:3000${NC}"
    echo ""
fi

echo -e "${CYAN}Starting TUI app...${NC}"
echo ""
echo -e "  ${BLUE}Tips:${NC}"
echo "    - Press 'q' to quit the TUI"
echo "    - Press '?' for help"
echo "    - Go to Settings → Sync to pair with PWA"
echo ""

# Trap to kill PWA server on exit
cleanup() {
    if [ -n "$PWA_PID" ]; then
        kill $PWA_PID 2>/dev/null || true
    fi
}
trap cleanup EXIT

sleep 2

# Launch TUI in foreground
exec ./bin/noise "$@"
