# VoxForge

Transform your voice into music - A browser-based voice-to-song tool that transforms vocal recordings (singing, humming, beatboxing) into complete music arrangements with drums, bass, chords, and AI-generated lyrics.

## Features

- 🎤 **Voice Recording** - Record directly in browser with real-time canvas-based waveform visualization
- 🎹 **Pitch Detection** - Extract melody as MIDI notes using YIN algorithm
- ⏱️ **BPM Detection** - Auto-detect tempo from your recording
- 🎼 **Key Detection** - Identify musical key automatically
- 🥁 **Drum Generation** - Multiple pattern styles (simple, moderate, busy)
- 🎸 **Bass & Chords** - Generate harmonic backing following detected key
- 📊 **Section Management** - Build multi-part songs with multiple sections
- 🎚️ **Arrangement Modes** - Sequential or layered playback
- 💾 **Stem Export** - Individual instrument tracks (WAV)
- 📝 **MIDI Export** - Editable melody for DAWs
- 🤖 **AI Lyrics** - Generate lyrics matching melody rhythm (OpenAI-compatible APIs)

## Technology Stack

- **Next.js 15.0.0** - React framework with App Router
- **React 19.0.0** - UI library with React Compiler
- **TypeScript 5.6.0** - Type safety
- **Tailwind CSS 3.4.0** - Styling
- **Tone.js 15.1.22** - Music generation and synthesis
- **pitchfinder 2.3.2** - Pitch detection (YIN algorithm)
- **realtime-bpm-analyzer 4.0.2** - Tempo detection
- **Tonal.js 6.4.2** - Music theory (key/scale detection)
- **@tonejs/midi 2.0.28** - MIDI file export
- **lucide-react 0.553.0** - Icons
- **OpenAI SDK 6.8.1** - AI lyrics generation (OpenAI-compatible APIs)

## Installation

### Prerequisites

- Node.js 18+ 
- npm or yarn
- Modern browser (Chrome 90+, Firefox 88+, Safari 14+)
- Microphone access

### Setup

```bash
# Clone the repository
git clone <repository-url>
cd VoxForge

# Install dependencies
npm install

# Create environment file (optional - only needed for AI lyrics feature)
# Create .env.local file manually with your API key

# Start development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Environment Variables (Optional - for AI Lyrics feature)

Create a `.env.local` file in the root directory with the following variables:

```bash
# Required: OpenAI-compatible API key (for AI lyrics feature)
OPENAI_API_KEY=<your-openai-key>

# Optional: Custom API base URL (for OpenRouter, Chutes, Z.ai, etc.)
# If not specified, defaults to https://api.openai.com/v1
OPENAI_API_BASE_URL=https://openrouter.ai/api/v1

# Optional: Model identifier
# Examples: gpt-4o-mini, openrouter/anthropic/claude-3.5-sonnet, z.ai/glm-4.6
OPENAI_MODEL=gpt-4o-mini
```

### Supported AI Providers

- **OpenRouter** (recommended) - Base URL: `https://openrouter.ai/api/v1`
- **Chutes** - Provider-specific URL
- **Z.ai** - Provider-specific URL (supports GLM 4.6)
- Any OpenAI-compatible API endpoint

**Note:** The AI lyrics feature is optional. All other features work without an API key.

## Usage

1. **Record** - Click "Start Recording" and sing, hum, or beatbox a melody. Watch the real-time waveform visualization.
2. **Analyze** - The app automatically detects pitch, BPM, and key after recording stops.
3. **Generate Music** - Select instruments (drums, bass, chords) and click "Generate Music" to create backing tracks.
4. **Playback** - Use play/pause/stop controls to hear your arrangement. The metronome shows the beat.
5. **Multiple Sections** - Click "New Section" to record additional parts (verse, chorus, etc.).
6. **Arrange** - Choose between sequential (sections play in order) or layered (sections play together) arrangement modes.
7. **Export** - Download individual stems (vocal, drums, bass, chords, mix) as WAV files or export MIDI.
8. **AI Lyrics** - Enter a theme and mood, then click "Generate Lyrics" to get AI-generated lyrics matching your melody rhythm.

## Project Structure

```
voxforge/
├── app/
│   ├── api/
│   │   └── lyrics/
│   │       └── route.ts          # Lyrics generation API route
│   ├── components/                # React components
│   │   ├── AnalysisDisplay.tsx   # Analysis results display
│   │   ├── ExportPanel.tsx       # Export controls
│   │   ├── InstrumentPicker.tsx  # Instrument selection
│   │   ├── LyricsAssistant.tsx   # AI lyrics UI
│   │   ├── Metronome.tsx         # Visual metronome
│   │   ├── PlaybackControls.tsx  # Play/pause/stop
│   │   ├── Recorder.tsx          # Recording UI
│   │   ├── SectionManager.tsx    # Section management
│   │   └── Waveform.tsx         # Canvas waveform visualization
│   ├── layout.tsx                 # Root layout
│   ├── page.tsx                   # Main page
│   └── globals.css                # Global styles
├── lib/
│   ├── audio/                     # Audio processing classes
│   │   ├── bpm-detector.ts       # BPM detection
│   │   ├── key-detector.ts       # Key detection
│   │   ├── midi-exporter.ts      # MIDI export
│   │   ├── music-generator.ts    # Music generation
│   │   ├── pitch-detector.ts     # Pitch detection
│   │   ├── recorder.ts           # Audio recording
│   │   └── stem-exporter.ts      # Stem export
│   ├── types/
│   │   └── index.ts              # TypeScript type definitions
│   └── utils/
│       ├── audio-utils.ts        # Audio utility functions
│       └── music-theory.ts       # Music theory helpers
├── public/
│   └── test-recordings/          # Test audio files
├── Docs/                          # Project documentation
├── package.json
├── tsconfig.json
├── tailwind.config.ts
└── next.config.js
```

## Development

```bash
# Development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Lint
npm run lint
```

## Browser Support

| Browser | Status | Notes |
|---------|--------|-------|
| Chrome 90+ | ✅ Full | Recommended |
| Firefox 88+ | ✅ Full | Works great |
| Safari 14+ | ⚠️ Mostly | Some audio quirks |
| Edge 90+ | ✅ Full | Chromium-based |
| Mobile Chrome | ⚠️ Limited | Recording may not work |
| Mobile Safari | ⚠️ Limited | Recording may not work |

## License

MIT License - See LICENSE file for details.

## Acknowledgments

- [Tone.js](https://tonejs.github.io/) - Web Audio framework
- [pitchfinder](https://github.com/peterkhayes/pitchfinder) - Pitch detection algorithms
- [realtime-bpm-analyzer](https://github.com/dlepaux/realtime-bpm-analyzer) - BPM detection
- [Tonal.js](https://github.com/tonaljs/tonal) - Music theory library
- [Lucide React](https://lucide.dev/) - Icon library

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
