# PWA Companion Sync Guide

noise.sh includes an embedded sync server for capturing ideas from your phone using a companion Progressive Web App (PWA).

## Overview

The sync system lets you:
- Capture voice memos on your phone
- Snap photos of lyric ideas
- Tap tempo to record BPM
- Jot quick text notes
- Review everything in the desktop app's Idea Inbox

All syncing happens over your local network—no cloud required.

## Setup

### 1. Enable Sync in Config

Add to your config file (`~/.config/noise/noise.yaml`):

```yaml
sync:
  enabled: true
  port: 8765
  auto_start: true
  media_path: "data/sync/media"
```

### 2. Start the Sync Server

The server starts automatically if `auto_start: true`. Otherwise, enable it from the settings menu.

### 3. Pair Your Device

1. Open sync settings in noise.sh
2. Select "Generate Pairing Code"
3. Note the 6-digit code and local URL (e.g., `http://192.168.1.100:8765`)
4. Open the PWA on your phone
5. Enter the pairing code

## Using the PWA

### Voice Memos
- Tap the microphone icon
- Record your idea
- It syncs automatically to your inbox

### Photos
- Tap the camera icon
- Capture a photo (setlist, napkin lyrics, etc.)
- Photo syncs to your inbox

### Tap Tempo
- Tap the tempo button rhythmically
- BPM is calculated and synced

### Text Notes
- Quick text entry for lyrics or ideas
- Syncs instantly

## Idea Inbox

In the desktop app, access the Idea Inbox to:
- Review captured ideas
- Assign ideas to songs
- Preview voice memos and photos
- Delete unwanted items

## Configuration Options

```yaml
sync:
  enabled: true          # Enable/disable sync server
  port: 8765             # Server port
  auto_start: true       # Start server on app launch
  media_path: "data/sync/media"  # Where media files are stored
```

## Network Requirements

- Desktop and phone must be on the same local network
- Port 8765 (default) must not be blocked by firewall
- No internet connection required

## Security

- **Local Network Only:** Server only accepts connections from local network
- **Pairing Required:** Devices must be paired with a 6-digit code
- **Code Expiry:** Pairing codes expire after 5 minutes
- **No Cloud:** All data stays on your local network

## Troubleshooting

### Can't Connect from Phone
- Ensure both devices are on the same WiFi network
- Check that the sync server is running (see status bar)
- Verify the IP address is correct
- Check firewall settings

### Pairing Code Doesn't Work
- Codes expire after 5 minutes—generate a new one
- Ensure you're entering the code correctly (numbers only)

### Media Not Syncing
- Check available disk space
- Verify the media_path directory exists and is writable
- Check the app logs for errors

### Server Won't Start
- Port may be in use—try a different port in config
- Check for permission issues with the media directory

## File Storage

Synced media is stored locally:
```
data/sync/media/
├── voice/    # Voice memos (.wav)
└── photos/   # Photos (.jpg, .png)
```

Files are named with timestamp, device ID, and random suffix to prevent conflicts.

## Database

Synced ideas are stored in the SQLite database:
- `captured_ideas` - All captured content
- `paired_devices` - Paired device records

## Privacy

- All sync happens over your local network
- No data is sent to external servers
- Media files are stored only on your machine
- You control what gets captured and kept
