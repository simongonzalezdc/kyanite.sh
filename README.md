# VoxForge

Transform your voice into music - A browser-based voice-to-song tool that transforms vocal recordings (singing, humming, beatboxing) into complete music arrangements with drums, bass, chords, and AI-generated lyrics.

## Features

- 🎤 **Voice Recording** - Record directly in browser with real-time waveform visualization
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

- **Next.js 15** - React framework with App Router
- **React 19** - UI library with React Compiler
- **TypeScript 5.6** - Type safety
- **Tailwind CSS 3.4** - Styling
- **Tone.js 15.1** - Music generation and synthesis
- **pitchfinder 2.3** - Pitch detection (YIN algorithm)
- **realtime-bpm-analyzer 4.0** - Tempo detection
- **Tonal.js 6.4** - Music theory (key/scale detection)
- **@tonejs/midi 2.0** - MIDI file export
- **wavesurfer.js 7.11** - Waveform display
- **lucide-react 0.553** - Icons
- **OpenAI SDK 6.8** - AI lyrics generation (OpenAI-compatible APIs)

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

# Create environment file
cp .env.example .env.local

# Add your OpenAI-compatible API key
# OPENAI_API_KEY=your-api-key-here
# OPENAI_API_BASE_URL=https://openrouter.ai/api/v1  # Optional
# OPENAI_MODEL=gpt-4o-mini  # Optional

# Start development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

## Environment Variables

Create a `.env.local` file with the following variables:

```bash
# Required: OpenAI-compatible API key
OPENAI_API_KEY=your-api-key-here

# Optional: Custom API base URL (for OpenRouter, Chutes, Z.ai, etc.)
OPENAI_API_BASE_URL=https://openrouter.ai/api/v1

# Optional: Model identifier
OPENAI_MODEL=gpt-4o-mini
```

### Supported AI Providers

- **OpenRouter** - Base URL: `https://openrouter.ai/api/v1`
- **Chutes** - Provider-specific URL
- **Z.ai** - Provider-specific URL (GLM 4.6)
- Any OpenAI-compatible API endpoint

## Usage

1. **Record** - Click "Start Recording" and sing, hum, or beatbox a melody
2. **Analyze** - The app automatically detects pitch, BPM, and key
3. **Generate** - Click "Generate Music" to create drums, bass, and chords
4. **Arrange** - Record multiple sections and arrange them sequentially or layered
5. **Export** - Download individual stems or MIDI files
6. **Lyrics** - Use AI Lyric Assist to generate lyrics matching your melody

## Project Structure

```
voxforge/
├── app/
│   ├── api/lyrics/          # Lyrics generation API route
│   ├── components/          # React components
│   ├── layout.tsx           # Root layout
│   ├── page.tsx             # Main page
│   └── globals.css          # Global styles
├── lib/
│   ├── audio/               # Audio processing classes
│   ├── types/               # TypeScript types
│   └── utils/               # Utility functions
├── public/                  # Static assets
└── Docs/                    # Documentation
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
- [WaveSurfer.js](https://wavesurfer-js.org/) - Audio visualization

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

