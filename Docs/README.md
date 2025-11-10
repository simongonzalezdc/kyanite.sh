# VOXFORGE - PROJECT README
**Version:** 2.0  
**Last Updated:** November 9, 2025

---

## 🎵 WHAT IS VOXFORGE?

VoxForge is a browser-based voice-to-song tool that transforms your vocal recordings (singing, humming, or beatboxing) into complete music arrangements with drums, bass, and chords.

**Core Value:** Turn a 30-second vocal idea into a full song demo in under 60 seconds.

---

## ✨ KEY FEATURES

- 🎤 **Voice Recording** - Record directly in browser
- 🎹 **Pitch Detection** - Extract melody as MIDI notes
- ⏱️ **BPM Detection** - Auto-detect tempo
- 🎼 **Key Detection** - Identify musical key
- 🥁 **Drum Generation** - Multiple pattern styles
- 🎸 **Bass & Chords** - Generate harmonic backing
- 📊 **Section Management** - Build multi-part songs
- 🎚️ **Arrangement Modes** - Sequential or layered playback
- 💾 **Stem Export** - Individual instrument tracks (WAV)
- 📝 **MIDI Export** - Editable melody for DAWs
- 🤖 **AI Lyrics** - Optional lyrics generation (OpenAI/Anthropic)

---

## 📚 COMPLETE DOCUMENTATION

This project includes comprehensive documentation for both developers and humans:

### For AI Coding Agents

1. **[VOXFORGE_TECHNICAL_SPEC_UPDATED.md](VOXFORGE_TECHNICAL_SPEC_UPDATED.md)**
   - Complete technical specification
   - System architecture
   - Technology stack details
   - Development timeline
   - All type definitions
   - Browser requirements

2. **[VOXFORGE_IMPLEMENTATION_GUIDE.md](VOXFORGE_IMPLEMENTATION_GUIDE.md)**
   - Step-by-step code implementation
   - Complete code examples for Days 1-2
   - Pattern continues for all 12 days
   - Testing checklists
   - Troubleshooting guides

3. **[VOXFORGE_API_REFERENCE.md](VOXFORGE_API_REFERENCE.md)**
   - Complete API documentation
   - All custom modules documented
   - Third-party library usage examples
   - Type definitions
   - Utility functions

4. **[VOXFORGE_AI_AGENT_QUICK_START.md](VOXFORGE_AI_AGENT_QUICK_START.md)**
   - Quick reference guide
   - Copy-paste ready commands
   - Common patterns
   - Debug checklist
   - Performance tips

### For Humans

5. **[VOXFORGE_DAY_BY_DAY_CHECKLIST.md](VOXFORGE_DAY_BY_DAY_CHECKLIST.md)**
   - 12-day development checklist
   - Daily goals and tasks
   - Testing requirements
   - Progress tracking
   - Time estimates

6. **[VOXFORGE_HUMAN_TASKS.md](VOXFORGE_HUMAN_TASKS.md)**
   - All non-coding tasks
   - Testing procedures
   - Content creation guides
   - Deployment instructions
   - Marketing activities
   - Estimated time: 10-15 hours

---

## 🚀 QUICK START

### Prerequisites

- Node.js 18+
- Modern browser (Chrome 90+, Firefox 88+, Safari 14+)
- Microphone access

### Installation

```bash
# Clone or create project
npx create-next-app@latest voxforge --typescript --tailwind --app

# Install dependencies
cd voxforge
npm install tone pitchfinder realtime-bpm-analyzer tonal @tonejs/midi
npm install wavesurfer.js lucide-react

# Optional: For lyrics feature
npm install openai

# Start development server
npm run dev
```

Open http://localhost:3000

### Project Structure

```
voxforge/
├── app/
│   ├── page.tsx                   # Main interface
│   ├── components/                # React components
│   └── api/lyrics/                # Optional lyrics API
├── lib/
│   ├── audio/                     # Audio processing
│   ├── types/                     # TypeScript types
│   └── utils/                     # Utility functions
├── public/
│   └── test-recordings/           # Test audio files
├── docs/                          # Documentation (this folder)
└── README.md                      # This file
```

---

## 🛠️ TECHNOLOGY STACK

### Frontend Core
- **Next.js 15** - React framework
- **React 19** - UI library
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling

### Audio Processing
- **Tone.js** - Music generation and synthesis
- **pitchfinder** - Pitch detection (YIN algorithm)
- **realtime-bpm-analyzer** - Tempo detection
- **Tonal.js** - Music theory (key/scale detection)
- **@tonejs/midi** - MIDI file export

### UI & Visualization
- **WaveSurfer.js** - Waveform display
- **Lucide React** - Icons
- **Web Audio API** - Native browser audio

### Optional Backend
- **OpenAI API** - Lyrics generation
- **Anthropic API** - Alternative lyrics generation

---

## 📅 DEVELOPMENT TIMELINE

**Total Duration:** 12 days

- **Days 1-3:** Audio recording, pitch/BPM/key detection
- **Days 4-6:** Music generation (drums, bass, chords)
- **Days 7-8:** Section management & arrangement
- **Days 9-10:** Export functionality (stems & MIDI)
- **Day 11:** Optional lyrics assistant
- **Day 12:** Testing & deployment

See [VOXFORGE_DAY_BY_DAY_CHECKLIST.md](VOXFORGE_DAY_BY_DAY_CHECKLIST.md) for detailed breakdown.

---

## 🎯 UNIQUE POSITIONING

### What Makes VoxForge Different?

**Preserves Your Creativity:**
- Keeps your original melody (doesn't transform voice)
- Generates instrumental backing around your input
- You remain the artist

**Privacy First:**
- 100% browser-based (MVP)
- No audio uploads
- No user accounts required
- All processing happens locally

**Professional Output:**
- Export separate stems for DAW finishing
- MIDI export for editing
- High-quality synthesis

**Fast & Simple:**
- 30 seconds to record
- 30 seconds to generate
- 60 seconds total to demo

### Competitive Comparison

| Feature | VoxForge | Competitors |
|---------|----------|-------------|
| Voice Input | ✅ | ✅ |
| Preserves Original Melody | ✅ | ❌ (transform voice) |
| Browser-Based | ✅ | ❌ (cloud processing) |
| No Upload Required | ✅ | ❌ |
| Stem Export | ✅ | ❌ |
| MIDI Export | ✅ | ⚠️ (some) |
| Free (MVP) | ✅ | ❌ (subscription) |

---

## 📖 HOW TO USE

### Basic Workflow

1. **Record**
   - Click "Record" button
   - Sing, hum, or beatbox a melody
   - Click "Stop"

2. **Analyze**
   - App auto-detects pitch, BPM, and key
   - Review detected information
   - Manually adjust if needed

3. **Generate**
   - Click "Generate Music"
   - Select instruments (drums, bass, chords)
   - Choose pattern style

4. **Arrange**
   - Record multiple sections (verse, chorus, etc.)
   - Arrange sequentially or layer simultaneously
   - Adjust per-section instruments

5. **Export**
   - Export individual stems (vocal, drums, bass, chords)
   - Export full mix
   - Export MIDI for DAW editing

### Advanced Features

- **Multi-Section Songs** - Build verse/chorus structure
- **Layered Recording** - Record harmony parts
- **Custom Arrangements** - Per-section instrument selection
- **AI Lyrics** - Generate lyrics matching rhythm (optional)

---

## 🧪 TESTING

### Manual Testing

See [VOXFORGE_HUMAN_TASKS.md](VOXFORGE_HUMAN_TASKS.md) for complete testing procedures.

**Quick Test:**
1. Record voice saying "Do Re Mi Fa Sol La Ti Do"
2. Verify notes detected: C, D, E, F, G, A, B, C
3. Generate music
4. Export stems
5. Verify all files download correctly

### Browser Support

| Browser | Status | Notes |
|---------|--------|-------|
| Chrome 90+ | ✅ Full | Recommended |
| Firefox 88+ | ✅ Full | Works great |
| Safari 14+ | ⚠️ Mostly | Some audio quirks |
| Edge 90+ | ✅ Full | Chromium-based |
| Mobile Chrome | ⚠️ Limited | Recording may not work |
| Mobile Safari | ⚠️ Limited | Recording may not work |

---

## 🚢 DEPLOYMENT

### Quick Deploy

**Railway (Recommended):**
```bash
npm i -g @railway/cli
railway login
railway init
railway up
```

**Vercel:**
```bash
npm i -g vercel
vercel
```

**Netlify:**
```bash
npm i -g netlify-cli
netlify deploy --prod
```

### Environment Variables

If using lyrics feature, add to `.env.local`:

```bash
OPENAI_API_KEY=sk-...
# OR
ANTHROPIC_API_KEY=sk-ant-...
```

---

## 📄 LICENSE

MIT License - See [LICENSE](LICENSE) file for details.

**What this means:**
- ✅ Free to use
- ✅ Free to modify
- ✅ Free to distribute
- ✅ Commercial use allowed
- ✅ No attribution required (but appreciated!)

**Dependencies:**
All libraries used are permissively licensed (MIT, Apache 2.0, BSD).

---

## 🙏 ACKNOWLEDGMENTS

### Open Source Libraries

- [Tone.js](https://tonejs.github.io/) - Web Audio framework
- [pitchfinder](https://github.com/peterkhayes/pitchfinder) - Pitch detection algorithms
- [realtime-bpm-analyzer](https://github.com/dlepaux/realtime-bpm-analyzer) - BPM detection
- [Tonal.js](https://github.com/tonaljs/tonal) - Music theory library
- [WaveSurfer.js](https://wavesurfer-js.org/) - Audio visualization

### Inspiration

- Joe Sullivan's [Beat Detection Article](http://joesul.li/van/beat-detection-using-web-audio/)
- Spotify's [Basic Pitch](https://github.com/spotify/basic-pitch) research

---

## 🗺️ ROADMAP

### MVP (Current)
- [x] Voice recording
- [x] Pitch/BPM/key detection
- [x] Music generation
- [x] Multi-section arrangement
- [x] Stem & MIDI export
- [ ] Lyrics assistant (optional)

### Phase 2 (Future)
- [ ] User accounts
- [ ] Cloud project storage
- [ ] Collaboration features
- [ ] More instruments (guitar, piano, strings)
- [ ] Audio effects (reverb, delay, EQ)
- [ ] Mobile app (iOS/Android)
- [ ] Advanced arrangement tools

### Monetization Ideas
- **Free Tier:** 3 exports/day
- **Pro Tier ($9/mo):** Unlimited exports, more instruments
- **Studio Tier ($29/mo):** Advanced features, priority support

---

## 💡 CONTRIBUTING

### Ways to Contribute

1. **Report Bugs**
   - Open GitHub issue with:
     - Steps to reproduce
     - Expected vs actual behavior
     - Browser and OS info

2. **Suggest Features**
   - Open GitHub issue with "Feature Request" label
   - Describe use case
   - Explain why it's valuable

3. **Submit Pull Requests**
   - Fork the repository
   - Create feature branch
   - Make changes
   - Test thoroughly
   - Submit PR with description

4. **Share Feedback**
   - Twitter: @yourhandle
   - Email: your@email.com
   - Discord: [invite link]

---

## 📞 SUPPORT

### Getting Help

1. **Check Documentation**
   - Read all .md files in `/docs/`
   - Search for error messages

2. **Common Issues**
   - See [VOXFORGE_API_REFERENCE.md](VOXFORGE_API_REFERENCE.md) troubleshooting section
   - Check browser console for errors

3. **Ask Questions**
   - Open GitHub issue with "Question" label
   - Be specific about what you're trying to do

### FAQ

**Q: Why can't I record audio?**  
A: Check microphone permissions in browser settings. Must use HTTPS.

**Q: Pitch detection is inaccurate. What should I do?**  
A: Sing louder and clearer. Reduce background noise. Try humming instead.

**Q: BPM detection is wrong. Can I fix it?**  
A: Yes! Manually override the detected BPM in the UI.

**Q: Can I use this commercially?**  
A: Yes! MIT license allows commercial use.

**Q: Will there be a mobile app?**  
A: Maybe in Phase 2. Web version works on mobile browsers (with limitations).

**Q: How much does it cost?**  
A: MVP is free and open source. Future paid tiers may be added.

---

## 📊 PROJECT STATUS

**Current Version:** 2.0  
**Status:** ⬜ Planning | ⬜ In Development | ⬜ MVP Complete  
**Last Updated:** November 9, 2025

**Development Progress:**
- [ ] Day 1: Audio Recording
- [ ] Day 2: Pitch Detection
- [ ] Day 3: BPM & Key Detection
- [ ] Day 4: Drum Generation
- [ ] Day 5: Bass & Chords
- [ ] Day 6: Instrument Selection
- [ ] Day 7: Section Management
- [ ] Day 8: Arrangement Modes
- [ ] Day 9: Stem Export
- [ ] Day 10: MIDI Export
- [ ] Day 11: Lyrics Assistant (optional)
- [ ] Day 12: Testing & Deployment

---

## 🎨 DESIGN PHILOSOPHY

### Core Principles

1. **Simplicity** - One clear purpose, no feature bloat
2. **Speed** - Under 60 seconds from idea to demo
3. **Privacy** - No uploads, no tracking, no accounts (MVP)
4. **Quality** - Professional-grade output for DAW finishing
5. **Accessibility** - Works in any modern browser
6. **Open Source** - Learn, modify, improve, share

### UI/UX Goals

- **Minimalist** - Clean, distraction-free interface
- **Intuitive** - No tutorial needed
- **Responsive** - Instant feedback
- **Forgiving** - Easy to experiment and retry
- **Delightful** - Smooth animations, satisfying interactions

---

## 🔗 LINKS

- **Live Demo:** [Coming soon]
- **GitHub:** [Your repo URL]
- **Documentation:** You're reading it!
- **Twitter:** [@yourhandle]
- **Blog Post:** [Coming soon]
- **Demo Video:** [Coming soon]

---

## 📝 CHANGELOG

### Version 2.0 (November 9, 2025)
- Updated tech stack with new libraries
- Added pitchfinder for pitch detection
- Added realtime-bpm-analyzer for BPM
- Added @tonejs/midi for MIDI export
- Comprehensive documentation overhaul
- 12-day implementation guide
- Separated human and AI tasks

### Version 1.0 (November 7, 2025)
- Initial specification
- Original concept and architecture
- Basic feature set defined

---

## 🎉 GET STARTED

Ready to build VoxForge?

1. **Read This:** [VOXFORGE_TECHNICAL_SPEC_UPDATED.md](VOXFORGE_TECHNICAL_SPEC_UPDATED.md)
2. **For AI Agents:** [VOXFORGE_AI_AGENT_QUICK_START.md](VOXFORGE_AI_AGENT_QUICK_START.md)
3. **For Humans:** [VOXFORGE_HUMAN_TASKS.md](VOXFORGE_HUMAN_TASKS.md)
4. **Track Progress:** [VOXFORGE_DAY_BY_DAY_CHECKLIST.md](VOXFORGE_DAY_BY_DAY_CHECKLIST.md)

**Let's build something amazing! 🚀**

---

**Made with ❤️ by [Your Name]**

**Built with:**  
Next.js • React • TypeScript • Tone.js • Web Audio API

---

**END OF README**
