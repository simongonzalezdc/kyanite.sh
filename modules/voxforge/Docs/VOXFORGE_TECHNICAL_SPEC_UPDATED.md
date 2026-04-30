# VOXFORGE - TECHNICAL SPECIFICATION (UPDATED)
**Version:** 2.0  
**Last Updated:** November 9, 2025  
**Target Completion:** November 21, 2025 (12 days)

---

## PROJECT OVERVIEW

**Name:** VoxForge  
**Description:** Browser-based voice-to-song tool that transforms vocal recordings (singing, humming, beatboxing) into full song arrangements with separate instrument stems and MIDI export.

**Core Value Proposition:** Turn a 30-second vocal idea into a complete demo with drums, bass, keys, and arrangement in under 60 seconds.

**Unique Positioning:**
- Preserves user's original melody (doesn't transform voice)
- Generates instrumental backing around user input
- 100% browser-based (no uploads, no backend for MVP)
- Exports professional stems + MIDI for DAW finishing

---

## SYSTEM ARCHITECTURE

```
┌─────────────────────────────────────────────┐
│         USER BROWSER (100% Client-Side)     │
├─────────────────────────────────────────────┤
│  Next.js 15 App Router                      │
│  ├─ Microphone Input (Web Audio API)        │
│  ├─ Real-time Waveform (WaveSurfer.js)      │
│  │                                           │
│  ├─ AUDIO ANALYSIS PIPELINE                 │
│  │   ├─ Pitch Detection (pitchfinder)       │
│  │   ├─ BPM Detection (realtime-bpm)        │
│  │   └─ Key Detection (tonal.js)            │
│  │                                           │
│  ├─ MUSIC GENERATION                        │
│  │   ├─ Drums (Tone.js)                     │
│  │   ├─ Bass (Tone.js)                      │
│  │   ├─ Chords (Tone.js)                    │
│  │   └─ Metronome (Tone.js)                 │
│  │                                           │
│  ├─ EXPORT ENGINE                           │
│  │   ├─ Stem Export (Tone.Recorder)         │
│  │   ├─ Mix Export (Tone.Recorder)          │
│  │   └─ MIDI Export (@tonejs/midi)          │
│  │                                           │
│  └─ UI Components (React + Tailwind)        │
└─────────────────────────────────────────────┘

          │ Optional: Only for lyrics
          ▼
┌─────────────────────────────────────────────┐
│  BACKEND API (Next.js API Route)            │
│  └─ OpenAI/Anthropic (Lyrics Assist)        │
└─────────────────────────────────────────────┘
```

**Key Architectural Decisions:**
- **Client-side only for MVP** - No backend required for core features
- **No database** - All state in React (ephemeral)
- **No user accounts** - Direct download of files
- **Optional backend** - Only needed for lyrics assistance (Day 11+)

---

## UPDATED TECHNOLOGY STACK

### Frontend Core
```json
{
  "next": "^15.0.0",
  "react": "^19.0.0",
  "react-dom": "^19.0.0",
  "typescript": "^5.6.0",
  "tailwindcss": "^3.4.0"
}
```

### Audio Processing Libraries (NEW)
```json
{
  "tone": "^15.1.22",
  "pitchfinder": "^3.0.0",
  "realtime-bpm-analyzer": "^4.0.2",
  "tonal": "^6.0.0",
  "@tonejs/midi": "^2.0.28"
}
```

### UI & Visualization
```json
{
  "wavesurfer.js": "^7.0.0",
  "lucide-react": "^0.445.0"
}
```

### Optional Backend (Lyrics Only)
```json
{
  "openai": "^4.67.0"
}
```

---

## CORE FEATURES (MVP)

### 1. VOICE RECORDING
**Implementation:** Web Audio API + MediaRecorder

**Features:**
- Click to record from microphone
- Real-time waveform visualization (WaveSurfer.js)
- Max recording length: 3 minutes (adjustable)
- Audio format: 44.1kHz, mono
- Count-in metronome (optional)

**User Flow:**
1. User clicks "Record"
2. Browser requests mic permission
3. Optional: Metronome plays 4 beats count-in
4. Recording starts with live waveform
5. User clicks "Stop" or max time reached
6. Audio buffer ready for analysis

---

### 2. AUDIO ANALYSIS PIPELINE

#### A. PITCH DETECTION
**Library:** pitchfinder (YIN algorithm)  
**Purpose:** Extract melody as MIDI notes

**Process:**
```typescript
import Pitchfinder from 'pitchfinder';

// Initialize detector
const detectPitch = Pitchfinder.YIN({
  sampleRate: 44100,
  threshold: 0.1 // Confidence threshold
});

// Analyze audio in windows
const windowSize = 2048;
const hopSize = 512; // ~11.6ms at 44.1kHz
const pitches = [];

for (let i = 0; i < audioData.length - windowSize; i += hopSize) {
  const window = audioData.slice(i, i + windowSize);
  const pitch = detectPitch(window);
  
  if (pitch && pitch > 0) {
    pitches.push({
      frequency: pitch,
      time: i / sampleRate,
      midi: frequencyToMidi(pitch)
    });
  }
}
```

**Output:**
- Array of pitch values (Hz)
- Converted to MIDI note numbers (0-127)
- Timestamps for each detected note
- Confidence scores

---

#### B. BPM DETECTION
**Library:** realtime-bpm-analyzer  
**Purpose:** Auto-detect tempo

**Process:**
```typescript
import { createRealTimeBpmProcessor } from 'realtime-bpm-analyzer';

const audioContext = new AudioContext();
const bpmProcessor = await createRealTimeBpmProcessor(audioContext, {
  continuousAnalysis: false,
  stabilizationTime: 10000 // 10 seconds
});

// Connect audio source
const source = audioContext.createBufferSource();
source.buffer = audioBuffer;
source.connect(bpmProcessor);

// Listen for BPM
bpmProcessor.port.onmessage = (event) => {
  if (event.data.message === 'BPM_STABLE') {
    const detectedBPM = Math.round(event.data.data.bpm);
    console.log('Detected BPM:', detectedBPM);
  }
};

source.start();
```

**Features:**
- Real-time BPM detection
- Stable BPM after analysis period
- Falls back to 120 BPM if unclear
- User can manually override

**Output:**
- Detected BPM (integer)
- Confidence level
- Manual override option in UI

---

#### C. KEY DETECTION
**Library:** tonal.js  
**Purpose:** Identify musical key

**Process:**
```typescript
import { Note, Key, Interval } from 'tonal';

// Convert MIDI notes to note names
const noteNames = midiNotes.map(midi => Note.fromMidi(midi));

// Count note occurrences (pitch classes)
const pitchClasses = noteNames.map(note => Note.chroma(note));
const counts = new Array(12).fill(0);
pitchClasses.forEach(pc => counts[pc]++);

// Find most common note (likely tonic)
const tonic = Note.fromMidi(
  midiNotes[counts.indexOf(Math.max(...counts))]
);

// Try major and minor keys
const majorKey = Key.majorKey(tonic);
const minorKey = Key.minorKey(tonic);

// Compare detected notes against key scales
const majorMatches = countMatchingNotes(noteNames, majorKey.scale);
const minorMatches = countMatchingNotes(noteNames, minorKey.scale);

const detectedKey = majorMatches > minorMatches 
  ? majorKey.tonic + ' Major'
  : minorKey.tonic + ' Minor';
```

**Output:**
- Detected key (e.g., "C Major", "A Minor")
- Scale notes
- Chord progression suggestions
- User can manually override

---

### 3. MUSIC GENERATION ENGINE

**Library:** Tone.js  
**Purpose:** Generate instrumental backing

#### A. DRUM GENERATOR

**Instruments:**
- Kick drum (low boom)
- Snare drum (mid crack)
- Hi-hat (high tick)

**Patterns:**
```typescript
import * as Tone from 'tone';

// Create drum synths
const kick = new Tone.MembraneSynth({
  pitchDecay: 0.05,
  octaves: 10,
  oscillator: { type: 'sine' },
  envelope: { attack: 0.001, decay: 0.4, sustain: 0.01, release: 1.4 }
}).toDestination();

const snare = new Tone.NoiseSynth({
  noise: { type: 'white' },
  envelope: { attack: 0.005, decay: 0.1, sustain: 0 }
}).toDestination();

const hihat = new Tone.MetalSynth({
  frequency: 200,
  envelope: { attack: 0.001, decay: 0.1, release: 0.01 },
  harmonicity: 5.1,
  modulationIndex: 32,
  resonance: 4000,
  octaves: 1.5
}).toDestination();

// Basic pattern (1/8 notes)
const drumPattern = new Tone.Pattern((time, note) => {
  if (note === 'kick') kick.triggerAttackRelease('C1', '8n', time);
  if (note === 'snare') snare.triggerAttack(time);
  if (note === 'hihat') hihat.triggerAttack(time);
}, ['kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat']);

drumPattern.start(0);
```

**Patterns Available:**
- Simple: Kick on 1 & 3, Snare on 2 & 4, Hi-hat on 8ths
- Moderate: Add kick variations, hi-hat on 16ths
- Busy: Complex fills, syncopation

---

#### B. BASS GENERATOR

**Purpose:** Root notes following chord progression

```typescript
const bass = new Tone.Synth({
  oscillator: { type: 'sawtooth' },
  envelope: { attack: 0.01, decay: 0.1, sustain: 0.5, release: 0.5 }
}).toDestination();

// Follow chord roots
const bassLine = new Tone.Sequence((time, note) => {
  bass.triggerAttackRelease(note, '4n', time);
}, ['C2', 'C2', 'G2', 'G2', 'A2', 'A2', 'F2', 'F2']);

bassLine.start(0);
```

---

#### C. CHORD GENERATOR

**Purpose:** Harmonic support

```typescript
import { Chord } from 'tonal';

const chordSynth = new Tone.PolySynth(Tone.Synth, {
  oscillator: { type: 'triangle' },
  envelope: { attack: 0.02, decay: 0.1, sustain: 0.8, release: 1 }
}).toDestination();

// Generate I-V-vi-IV progression in detected key
const key = 'C Major';
const progression = ['I', 'V', 'vi', 'IV'];

const chordNotes = progression.map(roman => {
  const chord = Chord.getChord('major', 'C', roman); // Simplified
  return chord.notes.map(n => n + '3'); // Octave 3
});

const chordPattern = new Tone.Part((time, notes) => {
  chordSynth.triggerAttackRelease(notes, '2n', time);
}, chordNotes.map((notes, i) => [i * 2, notes]));

chordPattern.start(0);
```

---

#### D. METRONOME

**Purpose:** Click track for recording

```typescript
const click = new Tone.MembraneSynth({
  pitchDecay: 0.008,
  envelope: { attack: 0.001, decay: 0.3, sustain: 0 }
}).toDestination();

const metronome = new Tone.Loop((time) => {
  click.triggerAttackRelease('C5', '32n', time);
}, '4n'); // Quarter notes

// Start with 4-beat count-in
Tone.Transport.scheduleOnce((time) => {
  metronome.start(time);
  // Stop after count-in
  metronome.stop(time + Tone.Time('4 * 4n').toSeconds());
}, '+0.1');
```

---

### 4. SECTION MANAGEMENT SYSTEM

**Purpose:** Build songs with multiple parts

**Features:**
- Record multiple sections (verse, chorus, bridge)
- Name each section
- Choose instrument arrangement per section
- Arrange sections in sequence OR layer simultaneously

**Data Structure:**
```typescript
interface Section {
  id: string;
  name: string; // "Verse 1", "Chorus", etc.
  audioBuffer: AudioBuffer;
  duration: number;
  
  // Analysis results
  pitches: Array<{ frequency: number; time: number; midi: number }>;
  detectedBPM: number;
  detectedKey: string;
  
  // Arrangement
  instruments: ('drums' | 'bass' | 'chords')[];
  
  // Generated audio
  generatedPart?: Tone.Part;
}

interface Project {
  sections: Section[];
  arrangementMode: 'sequential' | 'layered';
  masterBPM: number; // Used for layered mode
  masterKey: string; // Used for layered mode
}
```

**Section Operations:**
- Add new section
- Delete section
- Reorder sections (drag & drop)
- Toggle instruments per section
- Solo/mute sections

---

### 5. EXPORT SYSTEM

**Purpose:** Download stems and MIDI for DAW

#### A. STEM EXPORT

**What Gets Exported:**
- Vocal (original recording, dry)
- Drums (generated)
- Bass (generated)
- Chords (generated)
- Full Mix (all combined)

**Implementation:**
```typescript
import * as Tone from 'tone';

async function exportStem(
  instrumentType: 'vocal' | 'drums' | 'bass' | 'chords'
): Promise<Blob> {
  const recorder = new Tone.Recorder();
  
  // Route only the target instrument to recorder
  if (instrumentType === 'drums') {
    kick.connect(recorder);
    snare.connect(recorder);
    hihat.connect(recorder);
  } else if (instrumentType === 'vocal') {
    // Play back original AudioBuffer
    const player = new Tone.Player(audioBuffer);
    player.connect(recorder);
    player.start();
  }
  // ... similar for bass, chords
  
  recorder.start();
  Tone.Transport.start();
  
  // Wait for full duration
  await new Promise(resolve => 
    setTimeout(resolve, totalDuration * 1000)
  );
  
  const recording = await recorder.stop();
  return recording;
}
```

**File Format:** WAV (uncompressed) or MP3 (compressed)
**Naming Convention:**
- `voxforge-vocal-[timestamp].wav`
- `voxforge-drums-[timestamp].wav`
- `voxforge-bass-[timestamp].wav`
- `voxforge-chords-[timestamp].wav`
- `voxforge-mix-[timestamp].mp3`

---

#### B. MIDI EXPORT

**Purpose:** Export melody as MIDI for editing in DAW

**Implementation:**
```typescript
import { Midi } from '@tonejs/midi';

function exportMIDI(sections: Section[]): Blob {
  const midi = new Midi();
  midi.header.setTempo(project.masterBPM);
  
  // Create track for melody
  const track = midi.addTrack();
  track.name = 'VoxForge Melody';
  
  sections.forEach((section, idx) => {
    const startTime = idx * section.duration;
    
    section.pitches.forEach((pitch, i) => {
      const noteTime = startTime + pitch.time;
      const noteDuration = 0.25; // Quarter note default
      
      track.addNote({
        midi: pitch.midi,
        time: noteTime,
        duration: noteDuration,
        velocity: 0.8
      });
    });
  });
  
  // Convert to binary
  const midiArray = midi.toArray();
  return new Blob([midiArray], { type: 'audio/midi' });
}
```

**File Format:** Standard MIDI File (.mid)
**Naming:** `voxforge-melody-[timestamp].mid`

---

#### C. METADATA EXPORT

**Purpose:** Save project settings for reference

**File:** `voxforge-metadata-[timestamp].json`

```json
{
  "projectName": "My Song Idea",
  "createdAt": "2025-11-09T12:00:00Z",
  "bpm": 120,
  "key": "C Major",
  "sections": [
    {
      "name": "Verse 1",
      "duration": 8,
      "instruments": ["drums", "bass"]
    },
    {
      "name": "Chorus",
      "duration": 8,
      "instruments": ["drums", "bass", "chords"]
    }
  ],
  "totalDuration": 16
}
```

---

## REPOSITORY STRUCTURE

```
voxforge/
├── app/
│   ├── layout.tsx                 # Root layout
│   ├── page.tsx                   # Main app page
│   ├── globals.css                # Global styles
│   │
│   ├── api/                       # Server-side (optional)
│   │   └── lyrics/
│   │       └── route.ts           # Lyrics generation endpoint
│   │
│   └── components/                # React components
│       ├── Recorder.tsx           # Recording UI
│       ├── Waveform.tsx           # Visualization
│       ├── AnalysisDisplay.tsx    # Show detected BPM/key
│       ├── SectionManager.tsx     # Section list
│       ├── InstrumentPicker.tsx   # Select instruments
│       ├── Metronome.tsx          # Visual metronome
│       ├── PlaybackControls.tsx   # Play/pause/stop
│       ├── ExportPanel.tsx        # Export options
│       └── LyricsAssistant.tsx    # Lyrics UI (optional)
│
├── lib/
│   ├── audio/
│   │   ├── recorder.ts            # Web Audio recorder
│   │   ├── pitch-detector.ts     # pitchfinder wrapper
│   │   ├── bpm-detector.ts       # realtime-bpm wrapper
│   │   ├── key-detector.ts       # tonal.js wrapper
│   │   ├── music-generator.ts    # Tone.js generator
│   │   ├── stem-exporter.ts      # Export stems
│   │   └── midi-exporter.ts      # MIDI export
│   │
│   ├── types/
│   │   └── index.ts               # TypeScript interfaces
│   │
│   └── utils/
│       ├── audio-utils.ts         # Helper functions
│       └── music-theory.ts        # Theory helpers
│
├── public/
│   ├── test-recordings/           # Test audio files
│   └── assets/                    # Static assets
│
├── docs/
│   ├── TECHNICAL_SPEC.md          # This file
│   ├── IMPLEMENTATION_GUIDE.md    # Step-by-step guide
│   ├── API_REFERENCE.md           # Code reference
│   ├── DAY_BY_DAY_CHECKLIST.md   # Daily tasks
│   └── HUMAN_TASKS.md             # Non-coding tasks
│
├── package.json
├── next.config.js
├── tsconfig.json
├── tailwind.config.js
└── README.md
```

---

## TYPESCRIPT TYPE DEFINITIONS

```typescript
// lib/types/index.ts

export interface PitchPoint {
  frequency: number;  // Hz
  time: number;       // seconds
  midi: number;       // 0-127
  confidence?: number;
}

export interface BPMAnalysis {
  bpm: number;
  confidence: number;
  stable: boolean;
}

export interface KeyAnalysis {
  key: string;          // "C Major", "A Minor", etc.
  tonic: string;        // "C", "A", etc.
  scale: string[];      // ["C", "D", "E", "F", "G", "A", "B"]
  chordProgression: string[]; // ["I", "V", "vi", "IV"]
}

export type InstrumentType = 'drums' | 'bass' | 'chords';

export interface Section {
  id: string;
  name: string;
  
  // Original recording
  audioBuffer: AudioBuffer;
  duration: number;
  
  // Analysis
  pitches: PitchPoint[];
  bpm: BPMAnalysis;
  key: KeyAnalysis;
  
  // Arrangement
  instruments: InstrumentType[];
  
  // Generated (Tone.js)
  generatedParts?: {
    drums?: Tone.Part;
    bass?: Tone.Part;
    chords?: Tone.Part;
  };
}

export interface Project {
  name: string;
  sections: Section[];
  
  // Arrangement
  arrangementMode: 'sequential' | 'layered';
  
  // Global settings (for layered mode)
  masterBPM: number;
  masterKey: string;
  
  // Metadata
  createdAt: Date;
  duration: number;
}

export interface ExportOptions {
  format: 'wav' | 'mp3';
  stems: {
    vocal: boolean;
    drums: boolean;
    bass: boolean;
    chords: boolean;
    mix: boolean;
  };
  midi: boolean;
  metadata: boolean;
}

export interface LyricsRequest {
  syllableCount: number;
  phraseLengths: number[];
  detectedWords: string[];
  musicalKey: string;
  mood?: string;
  theme?: string;
}

export interface LyricsResponse {
  variations: Array<{
    lines: string[];
    syllableCounts: number[];
  }>;
}
```

---

## DEVELOPMENT TIMELINE (12 DAYS)

### Phase 1: Foundation (Days 1-3)

#### Day 1: Project Setup + Recording
**Goal:** User can record voice and see waveform

**Tasks:**
- [ ] Initialize Next.js project with TypeScript
- [ ] Install all dependencies
- [ ] Set up Tailwind CSS
- [ ] Create project structure (folders)
- [ ] Build Recorder component with Web Audio API
- [ ] Add WaveSurfer.js visualization
- [ ] Test mic recording and playback

**Deliverables:**
- User can click "Record" button
- Browser asks for mic permission
- Live waveform displays during recording
- User can stop recording
- Recording plays back

**Files Created:**
- `app/page.tsx`
- `app/components/Recorder.tsx`
- `app/components/Waveform.tsx`
- `lib/audio/recorder.ts`

---

#### Day 2: Pitch Detection
**Goal:** Extract melody from recording

**Tasks:**
- [ ] Install pitchfinder library
- [ ] Create pitch-detector.ts wrapper
- [ ] Implement sliding window analysis
- [ ] Convert Hz to MIDI notes
- [ ] Display detected notes in UI
- [ ] Test with various vocal inputs

**Deliverables:**
- After recording, pitch analysis runs automatically
- UI shows detected notes (C4, D4, E4, etc.)
- Console logs MIDI values for debugging

**Files Created:**
- `lib/audio/pitch-detector.ts`
- `app/components/AnalysisDisplay.tsx`

**Testing:**
- Record simple scale (Do Re Mi)
- Verify notes detected correctly
- Check for spurious detections

---

#### Day 3: BPM & Key Detection
**Goal:** Auto-detect tempo and musical key

**Tasks:**
- [ ] Install realtime-bpm-analyzer
- [ ] Create bpm-detector.ts wrapper
- [ ] Implement stable BPM detection
- [ ] Install tonal.js
- [ ] Create key-detector.ts
- [ ] Implement key detection algorithm
- [ ] Add manual override inputs in UI

**Deliverables:**
- BPM detected and displayed
- Key detected and displayed
- User can manually override both
- Default to 120 BPM if unclear

**Files Created:**
- `lib/audio/bpm-detector.ts`
- `lib/audio/key-detector.ts`

**Testing:**
- Sing at different tempos
- Verify BPM accuracy (±5 BPM acceptable)
- Test different keys

---

### Phase 2: Music Generation (Days 4-6)

#### Day 4: Drum Generator
**Goal:** Generate drum patterns

**Tasks:**
- [ ] Create drum synths (kick, snare, hihat)
- [ ] Implement 3 pattern types (simple/moderate/busy)
- [ ] Create music-generator.ts
- [ ] Connect to detected BPM
- [ ] Add play/stop controls
- [ ] Add visual beat indicator

**Deliverables:**
- Click "Generate Drums" button
- Drums play in sync with detected BPM
- User can switch between pattern types
- Visual metronome shows beats

**Files Created:**
- `lib/audio/music-generator.ts`
- `app/components/PlaybackControls.tsx`
- `app/components/Metronome.tsx`

---

#### Day 5: Bass & Chords
**Goal:** Add harmonic instruments

**Tasks:**
- [ ] Create bass synth
- [ ] Generate bass line from chord roots
- [ ] Create chord synth (poly)
- [ ] Generate I-V-vi-IV progression
- [ ] Adjust to detected key
- [ ] Mix all instruments together

**Deliverables:**
- Full arrangement plays: drums + bass + chords
- Arrangement follows detected key
- Sounds musical and in sync

**Testing:**
- Different keys sound correct
- Instruments balanced (not too loud)
- No clipping or distortion

---

#### Day 6: Instrument Selection
**Goal:** User can choose instruments per section

**Tasks:**
- [ ] Create InstrumentPicker component
- [ ] Add checkboxes: Drums, Bass, Chords
- [ ] Update music-generator to respect selection
- [ ] Test various combinations

**Deliverables:**
- User can toggle instruments on/off
- Changes apply immediately to playback
- At least one instrument must be selected

**Files Created:**
- `app/components/InstrumentPicker.tsx`

---

### Phase 3: Section Management (Days 7-8)

#### Day 7: Multi-Section Recording
**Goal:** Record multiple song sections

**Tasks:**
- [ ] Create Section data structure
- [ ] Build SectionManager component
- [ ] Add "New Section" button
- [ ] Name sections (Verse, Chorus, etc.)
- [ ] Store multiple sections in state
- [ ] Display section list

**Deliverables:**
- User can record multiple sections
- Each section shows: name, duration, BPM, key
- User can delete sections
- Sections persist during session

**Files Created:**
- `app/components/SectionManager.tsx`
- `lib/types/index.ts` (updated)

---

#### Day 8: Arrangement Modes
**Goal:** Sequential vs layered playback

**Tasks:**
- [ ] Add arrangement mode toggle
- [ ] Implement sequential playback (one after another)
- [ ] Implement layered playback (simultaneous)
- [ ] For layered: use master BPM/key
- [ ] Add drag-and-drop section reordering

**Deliverables:**
- "Sequential" mode: sections play in order
- "Layered" mode: sections play together
- User can reorder sections
- Smooth transitions between sections

**Testing:**
- Record 2 sections: verse + chorus
- Play sequentially: verse → chorus
- Play layered: both at once

---

### Phase 4: Export (Days 9-10)

#### Day 9: Stem Export
**Goal:** Download individual instrument tracks

**Tasks:**
- [ ] Create stem-exporter.ts
- [ ] Implement Tone.Recorder for each instrument
- [ ] Generate WAV files
- [ ] Add download buttons in UI
- [ ] Create ExportPanel component
- [ ] Test file downloads

**Deliverables:**
- Export buttons for: Vocal, Drums, Bass, Chords, Mix
- Each downloads as WAV file
- Files open in music players
- Stems are time-aligned

**Files Created:**
- `lib/audio/stem-exporter.ts`
- `app/components/ExportPanel.tsx`

---

#### Day 10: MIDI Export
**Goal:** Export melody as MIDI

**Tasks:**
- [ ] Install @tonejs/midi
- [ ] Create midi-exporter.ts
- [ ] Convert pitch data to MIDI format
- [ ] Include tempo and key signature
- [ ] Add download button
- [ ] Test in DAW (GarageBand, Ableton, etc.)

**Deliverables:**
- "Export MIDI" button downloads .mid file
- File opens in DAWs
- Notes match recorded melody
- Tempo set correctly

**Files Created:**
- `lib/audio/midi-exporter.ts`

---

### Phase 5: Polish & Optional Features (Days 11-12)

#### Day 11: Lyrics Assistant (Optional)
**Goal:** AI-powered lyrics suggestions

**Tasks:**
- [ ] Create Next.js API route
- [ ] Set up OpenAI API integration
- [ ] Create LyricsAssistant component
- [ ] Syllable counting
- [ ] Generate lyrics based on melody
- [ ] Display suggestions

**Deliverables:**
- "Get Lyrics" button
- User inputs theme/mood
- AI suggests lyrics matching melody rhythm
- User can accept/reject

**Files Created:**
- `app/api/lyrics/route.ts`
- `app/components/LyricsAssistant.tsx`

**Note:** Requires OpenAI API key (costs $)

---

#### Day 12: Testing & Deployment
**Goal:** Ship it!

**Tasks:**
- [ ] End-to-end testing of full workflow
- [ ] Test on different browsers (Chrome, Firefox, Safari)
- [ ] Test on mobile (recording may not work)
- [ ] Fix critical bugs
- [ ] Add loading states
- [ ] Deploy to Railway/Vercel
- [ ] Write README
- [ ] Record demo video

**Deliverables:**
- Live app at public URL
- All features working
- Smooth UX (no errors)
- README with screenshots

---

## INSTALLATION INSTRUCTIONS

### Prerequisites
- Node.js 18+ 
- npm or yarn
- Modern browser (Chrome 90+, Firefox 88+, Safari 14+)

### Setup

```bash
# 1. Create Next.js app
npx create-next-app@latest voxforge
cd voxforge

# 2. Install core dependencies
npm install tone pitchfinder realtime-bpm-analyzer tonal @tonejs/midi

# 3. Install UI dependencies
npm install wavesurfer.js lucide-react

# 4. Install Tailwind (if not already)
npm install -D tailwindcss postcss autoprefixer
npx tailwindcss init -p

# 5. Optional: Install OpenAI (for lyrics)
npm install openai

# 6. Start development server
npm run dev
```

### Configuration

**next.config.js:**
```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Allow audio file types
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(mp3|wav|ogg)$/,
      type: 'asset/resource'
    });
    return config;
  }
}

module.exports = nextConfig;
```

**tailwind.config.js:**
```javascript
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: '#3B82F6',
        secondary: '#8B5CF6',
      }
    },
  },
  plugins: [],
}
```

**tsconfig.json:**
```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable"],
    "jsx": "preserve",
    "module": "ESNext",
    "moduleResolution": "bundler",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "paths": {
      "@/*": ["./*"]
    }
  }
}
```

---

## ENVIRONMENT VARIABLES

Create `.env.local`:

```bash
# Optional: For lyrics feature
OPENAI_API_KEY=sk-...your-key-here

# Optional: For Anthropic API
ANTHROPIC_API_KEY=sk-ant-...your-key-here
```

---

## TESTING STRATEGY

### Unit Tests
- Audio recorder (mic permission, buffer creation)
- Pitch detector (accuracy with known frequencies)
- BPM detector (known tempo audio files)
- Key detector (known key audio files)
- MIDI exporter (correct note values)

### Integration Tests
- Full pipeline: record → analyze → generate → export
- Section management: add, delete, reorder
- Playback: sequential and layered modes

### Manual Testing
1. **Simple melody test:**
   - Sing Do-Re-Mi-Fa-Sol-La-Ti-Do
   - Verify correct notes detected
   - Verify key detected as C Major

2. **Rhythm test:**
   - Clap steady beat at 120 BPM
   - Verify BPM detected within ±5

3. **Multi-section test:**
   - Record verse (8 bars)
   - Record chorus (8 bars)
   - Play sequentially
   - Play layered
   - Verify both sound good

4. **Export test:**
   - Export all stems
   - Import into GarageBand/Ableton
   - Verify stems align perfectly
   - Verify MIDI matches melody

---

## PERFORMANCE TARGETS

| Metric | Target | Acceptable |
|--------|--------|------------|
| Recording latency | < 30ms | < 50ms |
| Pitch detection time | < 1s per 30s | < 2s per 30s |
| BPM detection time | < 2s | < 5s |
| Music generation time | < 1s | < 3s |
| Stem export time | < 5s per stem | < 10s per stem |
| UI responsiveness | 60 FPS | 30 FPS |

---

## BROWSER COMPATIBILITY

| Browser | Status | Notes |
|---------|--------|-------|
| Chrome 90+ | ✅ Fully supported | Recommended |
| Firefox 88+ | ✅ Fully supported | Works great |
| Safari 14+ | ⚠️ Mostly supported | Web Audio API quirks |
| Edge 90+ | ✅ Fully supported | Chromium-based |
| Mobile Chrome | ⚠️ Limited | Recording may not work |
| Mobile Safari | ⚠️ Limited | Recording may not work |

**Known Issues:**
- Mobile browsers may not support microphone recording
- Safari has stricter autoplay policies
- Some browsers require user gesture before audio

---

## DEPLOYMENT

### Option 1: Railway (Recommended)

```bash
# 1. Install Railway CLI
npm i -g @railway/cli

# 2. Login
railway login

# 3. Initialize
railway init

# 4. Deploy
railway up
```

**Environment:**
- Node.js 18+
- Build: `npm run build`
- Start: `npm start`
- Port: Auto-detected

### Option 2: Vercel

```bash
# 1. Install Vercel CLI
npm i -g vercel

# 2. Deploy
vercel
```

**Note:** Vercel may have request timeout limits for lyrics API

### Option 3: Netlify

```bash
# 1. Install Netlify CLI
npm i -g netlify-cli

# 2. Deploy
netlify deploy --prod
```

---

## FUTURE ENHANCEMENTS (Post-MVP)

### Phase 2 Features
- **User accounts:** Save projects to cloud
- **Project library:** Browse past recordings
- **Collaboration:** Share project links
- **More instruments:** Guitar, piano, strings
- **Effects:** Reverb, delay, EQ
- **Advanced arrangement:** Intro, outro, bridge sections
- **Chord suggestions:** AI-powered chord progressions
- **Voice effects:** Autotune, harmonizer
- **Mobile app:** Native iOS/Android

### Monetization Ideas
- **Free tier:** 3 exports per day
- **Pro tier ($9/mo):** Unlimited exports, more instruments
- **Studio tier ($29/mo):** Advanced features, priority support
- **One-time purchase:** $49 for lifetime access

---

## API REFERENCE (Optional Lyrics Feature)

### POST /api/lyrics/route.ts

**Request:**
```json
{
  "syllableCount": 32,
  "phraseLengths": [8, 8, 8, 8],
  "detectedWords": [],
  "musicalKey": "C Major",
  "mood": "happy",
  "theme": "summer"
}
```

**Response:**
```json
{
  "variations": [
    {
      "lines": [
        "Walking down the sunny street",
        "Feeling happy on my feet", 
        "Summer days are here to stay",
        "Let's go out and dance and play"
      ],
      "syllableCounts": [8, 8, 8, 8]
    }
  ]
}
```

---

## TROUBLESHOOTING

### Common Issues

**1. Microphone not working:**
- Check browser permissions
- Use HTTPS (required for mic access)
- Try different browser

**2. Pitch detection inaccurate:**
- Sing louder and clearer
- Reduce background noise
- Try humming instead of singing

**3. BPM detection wrong:**
- Manually override detected BPM
- Sing with clearer rhythm
- Use metronome during recording

**4. Export files corrupted:**
- Check browser console for errors
- Try smaller recording (< 1 minute)
- Clear browser cache

**5. Audio quality poor:**
- Use better microphone
- Reduce background noise
- Adjust gain/input volume

---

## SECURITY & PRIVACY

**No Data Storage:**
- All processing happens in browser
- No audio uploaded to server
- No user data collected

**Lyrics API (Optional):**
- Only text sent to OpenAI
- No voice recordings sent
- API key stored securely in environment variables

**Microphone Permission:**
- Only requested when user clicks record
- Can be revoked anytime
- Audio never leaves device (except exports)

---

## LICENSE

**Recommended:** MIT License

**Dependencies:**
- Tone.js: MIT
- pitchfinder: MIT
- realtime-bpm-analyzer: Apache 2.0
- Tonal.js: MIT
- @tonejs/midi: MIT
- WaveSurfer.js: BSD-3-Clause

All dependencies are permissively licensed for commercial use.

---

## CREDITS & ACKNOWLEDGMENTS

**Libraries:**
- [Tone.js](https://tonejs.github.io/) - Web Audio framework
- [pitchfinder](https://github.com/peterkhayes/pitchfinder) - Pitch detection
- [realtime-bpm-analyzer](https://github.com/dlepaux/realtime-bpm-analyzer) - BPM detection
- [Tonal.js](https://github.com/tonaljs/tonal) - Music theory
- [WaveSurfer.js](https://wavesurfer-js.org/) - Waveform visualization

**Inspiration:**
- Joe Sullivan's [Beat Detection Article](http://joesul.li/van/beat-detection-using-web-audio/)
- Spotify's [Basic Pitch](https://github.com/spotify/basic-pitch)
- Various voice-to-music experiments

---

**END OF TECHNICAL SPECIFICATION**
