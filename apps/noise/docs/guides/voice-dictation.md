# Voice Dictation Guide

noise.sh includes built-in voice-to-text dictation powered by whisper.cpp. Everything sets up automatically on first use.

## Quick Start

1. **Start Recording:** Press `Ctrl+D`
2. **Speak:** Say your lyrics naturally
3. **Stop Recording:** Press `Ctrl+D` again (or `Esc` to cancel)
4. **Done:** Your words appear at the cursor position

## First-Time Setup

On your first use of voice dictation, noise.sh will:
1. Download the whisper model (~142MB for base English)
2. Show download progress in the voice indicator
3. Initialize the speech recognition engine

This only happens once. After that, dictation works instantly.

## Keyboard Shortcuts

| Action | Shortcut |
|--------|----------|
| Start/Stop Dictation | `Ctrl+D` |
| Cancel Recording | `Esc` |

## Voice Indicator

While recording, you'll see:
- Recording duration
- Audio level meter (shows if your mic is picking up sound)
- Current state (Recording, Processing, etc.)

## Settings

Access voice settings through the settings menu to:
- **Change Model:** Choose between tiny (fast), base (balanced), or small (accurate)
- **Test Microphone:** Check your audio levels
- **Download Models:** Pre-download additional models

## Models

| Model | Size | Speed | Accuracy |
|-------|------|-------|----------|
| tiny.en | ~75MB | Fastest | Good |
| base.en | ~142MB | Balanced | Better |
| small.en | ~466MB | Slower | Best |

The default is `base.en` which offers a good balance.

## Troubleshooting

### No Microphone Detected
- Check system permissions for microphone access
- Ensure your microphone is connected and working
- Try a different audio input device

### Model Download Fails
- Check your internet connection
- Try again later (Hugging Face servers may be busy)
- Use `make download-model` to pre-download manually

### Poor Transcription Quality
- Speak clearly and at a moderate pace
- Reduce background noise
- Try the "small" model for better accuracy
- Ensure you're speaking in English (multilingual models coming soon)

### Recording Doesn't Stop
- Press `Esc` to force cancel
- Check if the voice service initialized correctly

## Technical Details

- **Engine:** whisper.cpp (local, privacy-preserving)
- **Audio Capture:** malgo (cross-platform)
- **Sample Rate:** 16kHz
- **Format:** 32-bit float mono
- **Max Duration:** 30 seconds (configurable)

## Privacy

All voice processing happens locally on your machine. No audio is ever sent to external servers. The whisper models are downloaded once and run entirely offline.
