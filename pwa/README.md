# noise.sh PWA Companion

Mobile capture companion for noise.sh - capture song ideas on the go that sync to your desktop.

## Features

- **Text Capture** - Quick lyric ideas and notes
- **Voice Memos** - Record melody ideas and vocal sketches  
- **Photo Capture** - Snap photos of handwritten lyrics, setlists, inspiration
- **Tap Tempo** - Tap out the BPM of your song idea

## How It Works

1. Enable sync in noise.sh desktop (Settings → Sync)
2. Open this PWA on your phone
3. Enter the server URL and pairing code
4. Start capturing ideas - they sync automatically!

## Technical Stack

- **Framework:** Next.js 15 (App Router)
- **UI:** React 18 + TypeScript
- **Styling:** Tailwind CSS 4
- **Storage:** Dexie.js (IndexedDB) for offline queue
- **PWA:** Serwist service worker

## Development

```bash
# Install dependencies
npm install

# Start dev server
npm run dev

# Build for production
npm run build

# Lint
npm run lint
```

## Sync Architecture

The PWA connects to the embedded sync server in noise.sh desktop over your local WiFi network:

```
┌─────────────┐         ┌─────────────────┐
│   PWA       │  WiFi   │  noise.sh TUI   │
│   Mobile    │ ──────► │  Desktop        │
│             │         │                 │
│ ┌─────────┐ │         │ ┌─────────────┐ │
│ │ Capture │ │ HTTP    │ │ Sync Server │ │
│ │ Tools   │───────────►│ (port 8765) │ │
│ └─────────┘ │         │ └─────────────┘ │
│             │         │        │        │
│ ┌─────────┐ │         │ ┌──────▼──────┐ │
│ │ Offline │ │ WS      │ │ Idea Inbox  │ │
│ │ Queue   │◄──────────│ │             │ │
│ └─────────┘ │         │ └─────────────┘ │
└─────────────┘         └─────────────────┘
```

## Privacy

- All sync happens over your local network
- No data is sent to external servers
- Ideas are stored locally until synced
- You control what gets captured and kept

## API Endpoints

The PWA communicates with these endpoints on the Go sync server:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/status` | GET | Server status |
| `/pair` | POST | Device pairing |
| `/ideas` | GET/POST | Ideas CRUD |
| `/upload` | POST | Media upload |
| `/ws` | WS | Real-time sync |

## License

Part of noise.sh - see main repository for license.
