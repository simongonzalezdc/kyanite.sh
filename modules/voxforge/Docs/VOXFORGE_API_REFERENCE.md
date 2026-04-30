# VOXFORGE - API REFERENCE
**Version:** 2.0  
**Last Updated:** November 9, 2025  
**Purpose:** Complete code reference for all libraries and custom modules

---

## TABLE OF CONTENTS

1. [Custom Modules API](#custom-modules-api)
2. [Third-Party Libraries Reference](#third-party-libraries-reference)
3. [Type Definitions](#type-definitions)
4. [Utility Functions](#utility-functions)
5. [Configuration](#configuration)

---

## CUSTOM MODULES API

### AudioRecorder

**Location:** `lib/audio/recorder.ts`

#### Constructor

```typescript
const recorder = new AudioRecorder();
```

Creates a new audio recorder instance with Web Audio API.

**No parameters required.**

#### Methods

##### `requestPermission(): Promise<boolean>`

Request microphone permission from the browser.

**Returns:** Promise resolving to `true` if permission granted, `false` otherwise.

**Example:**
```typescript
const hasPermission = await recorder.requestPermission();
if (!hasPermission) {
  alert('Microphone access required');
}
```

##### `startRecording(): Promise<void>`

Start recording audio from the microphone.

**Throws:** Error if mic permission denied or device unavailable.

**Example:**
```typescript
try {
  await recorder.startRecording();
  console.log('Recording started');
} catch (error) {
  console.error('Failed to start recording:', error);
}
```

##### `stopRecording(): Promise<AudioBuffer>`

Stop recording and return the audio buffer.

**Returns:** Promise resolving to AudioBuffer containing the recording.

**Example:**
```typescript
const audioBuffer = await recorder.stopRecording();
console.log('Duration:', audioBuffer.duration);
console.log('Sample rate:', audioBuffer.sampleRate);
```

##### `getWaveformData(): Uint8Array`

Get real-time waveform data for visualization (while recording).

**Returns:** Uint8Array of time-domain audio data (0-255).

**Example:**
```typescript
const waveform = recorder.getWaveformData();
// Use for canvas visualization
```

##### `getFrequencyData(): Uint8Array`

Get real-time frequency data for visualization (while recording).

**Returns:** Uint8Array of frequency-domain audio data (0-255).

**Example:**
```typescript
const frequencies = recorder.getFrequencyData();
// Use for spectrum analyzer
```

##### `dispose(): void`

Clean up resources (stop tracks, close audio context).

**Example:**
```typescript
recorder.dispose();
```

---

### PitchDetector

**Location:** `lib/audio/pitch-detector.ts`

#### Constructor

```typescript
const detector = new PitchDetector(sampleRate?: number);
```

**Parameters:**
- `sampleRate` (optional): Audio sample rate in Hz. Default: 44100

**Example:**
```typescript
const detector = new PitchDetector(44100);
```

#### Methods

##### `analyze(audioBuffer: AudioBuffer): PitchPoint[]`

Analyze an audio buffer and extract pitch information.

**Parameters:**
- `audioBuffer`: AudioBuffer to analyze

**Returns:** Array of PitchPoint objects

**Example:**
```typescript
const pitches = detector.analyze(audioBuffer);
console.log('Detected', pitches.length, 'pitch points');

pitches.forEach(pitch => {
  console.log(`Time: ${pitch.time}s, Note: ${detector.midiToNoteName(pitch.midi)}, Freq: ${pitch.frequency}Hz`);
});
```

##### `midiToNoteName(midi: number): string`

Convert MIDI note number to note name (e.g., 60 → "C4").

**Parameters:**
- `midi`: MIDI note number (0-127)

**Returns:** Note name with octave

**Example:**
```typescript
detector.midiToNoteName(60);  // "C4"
detector.midiToNoteName(69);  // "A4"
detector.midiToNoteName(72);  // "C5"
```

##### `getSimplifiedMelody(pitches: PitchPoint[]): number[]`

Get simplified melody with consecutive duplicates removed.

**Parameters:**
- `pitches`: Array of PitchPoint objects

**Returns:** Array of MIDI note numbers

**Example:**
```typescript
const simplified = detector.getSimplifiedMelody(pitches);
const noteNames = simplified.map(midi => detector.midiToNoteName(midi));
console.log('Melody:', noteNames.join(' - '));
```

##### `getStats(pitches: PitchPoint[]): PitchStats`

Get statistical summary of pitch analysis.

**Parameters:**
- `pitches`: Array of PitchPoint objects

**Returns:** Object containing:
- `averageFrequency`: number
- `minFrequency`: number
- `maxFrequency`: number
- `averageMidi`: number
- `range`: 'low' | 'mid' | 'high' | 'unknown'

**Example:**
```typescript
const stats = detector.getStats(pitches);
console.log('Average pitch:', stats.averageFrequency, 'Hz');
console.log('Vocal range:', stats.range);
```

---

### BPMDetector

**Location:** `lib/audio/bpm-detector.ts`

#### Constructor

```typescript
const bpmDetector = new BPMDetector();
```

#### Methods

##### `analyze(audioBuffer: AudioBuffer): Promise<BPMAnalysis>`

Analyze audio buffer and detect BPM.

**Parameters:**
- `audioBuffer`: AudioBuffer to analyze

**Returns:** Promise resolving to BPMAnalysis object

**Example:**
```typescript
const bpmAnalysis = await bpmDetector.analyze(audioBuffer);
console.log('Detected BPM:', bpmAnalysis.bpm);
console.log('Confidence:', bpmAnalysis.confidence);
console.log('Stable:', bpmAnalysis.stable);
```

---

### KeyDetector

**Location:** `lib/audio/key-detector.ts`

#### Constructor

```typescript
const keyDetector = new KeyDetector();
```

#### Methods

##### `analyze(pitches: PitchPoint[]): KeyAnalysis`

Detect musical key from pitch data.

**Parameters:**
- `pitches`: Array of PitchPoint objects from PitchDetector

**Returns:** KeyAnalysis object

**Example:**
```typescript
const keyAnalysis = keyDetector.analyze(pitches);
console.log('Detected key:', keyAnalysis.key);
console.log('Tonic:', keyAnalysis.tonic);
console.log('Scale notes:', keyAnalysis.scale);
console.log('Chord progression:', keyAnalysis.chordProgression);
```

---

### MusicGenerator

**Location:** `lib/audio/music-generator.ts`

#### Constructor

```typescript
const generator = new MusicGenerator();
```

#### Methods

##### `initialize(): Promise<void>`

Initialize Tone.js (required before use).

**Example:**
```typescript
await generator.initialize();
```

##### `setBPM(bpm: number): void`

Set the tempo for music generation.

**Parameters:**
- `bpm`: Beats per minute (40-240)

**Example:**
```typescript
generator.setBPM(120);
```

##### `setKey(key: string): void`

Set the musical key.

**Parameters:**
- `key`: Key signature (e.g., "C Major", "A Minor")

**Example:**
```typescript
generator.setKey('C Major');
```

##### `generateDrums(pattern: 'simple' | 'moderate' | 'busy'): void`

Generate drum pattern.

**Parameters:**
- `pattern`: Complexity level

**Example:**
```typescript
generator.generateDrums('moderate');
```

##### `generateBass(): void`

Generate bass line following chord progression.

**Example:**
```typescript
generator.generateBass();
```

##### `generateChords(): void`

Generate chord progression (I-V-vi-IV).

**Example:**
```typescript
generator.generateChords();
```

##### `start(): void`

Start playback of all generated instruments.

**Example:**
```typescript
generator.start();
```

##### `stop(): void`

Stop playback.

**Example:**
```typescript
generator.stop();
```

##### `dispose(): void`

Clean up all Tone.js resources.

**Example:**
```typescript
generator.dispose();
```

---

### StemExporter

**Location:** `lib/audio/stem-exporter.ts`

#### Constructor

```typescript
const exporter = new StemExporter(generator: MusicGenerator);
```

**Parameters:**
- `generator`: MusicGenerator instance

#### Methods

##### `exportStem(type: 'vocal' | 'drums' | 'bass' | 'chords' | 'mix', audioBuffer?: AudioBuffer): Promise<Blob>`

Export individual instrument or full mix.

**Parameters:**
- `type`: Which stem to export
- `audioBuffer`: (Optional) Original vocal recording for 'vocal' export

**Returns:** Promise resolving to audio Blob (WAV format)

**Example:**
```typescript
const drumStem = await exporter.exportStem('drums');
const url = URL.createObjectURL(drumStem);
const a = document.createElement('a');
a.href = url;
a.download = 'voxforge-drums.wav';
a.click();
```

---

### MidiExporter

**Location:** `lib/audio/midi-exporter.ts`

#### Constructor

```typescript
const midiExporter = new MidiExporter();
```

#### Methods

##### `export(pitches: PitchPoint[], bpm: number, key?: string): Blob`

Export melody as MIDI file.

**Parameters:**
- `pitches`: Array of PitchPoint objects
- `bpm`: Tempo in beats per minute
- `key`: (Optional) Musical key for key signature

**Returns:** MIDI file as Blob

**Example:**
```typescript
const midiBlob = midiExporter.export(pitches, 120, 'C Major');
const url = URL.createObjectURL(midiBlob);
const a = document.createElement('a');
a.href = url;
a.download = 'voxforge-melody.mid';
a.click();
```

---

## THIRD-PARTY LIBRARIES REFERENCE

### Tone.js

**Documentation:** https://tonejs.github.io/

#### Essential Classes

##### Tone.Transport

Global transport for timing and scheduling.

```typescript
import * as Tone from 'tone';

// Set BPM
Tone.Transport.bpm.value = 120;

// Start/stop
Tone.Transport.start();
Tone.Transport.stop();

// Schedule events
Tone.Transport.schedule((time) => {
  synth.triggerAttackRelease('C4', '8n', time);
}, '1m'); // At 1 measure
```

##### Tone.Synth

Monophonic synthesizer.

```typescript
const synth = new Tone.Synth({
  oscillator: { type: 'sine' },
  envelope: {
    attack: 0.01,
    decay: 0.1,
    sustain: 0.5,
    release: 0.5
  }
}).toDestination();

synth.triggerAttackRelease('C4', '8n');
```

##### Tone.PolySynth

Polyphonic synthesizer (for chords).

```typescript
const polySynth = new Tone.PolySynth(Tone.Synth).toDestination();

polySynth.triggerAttackRelease(['C4', 'E4', 'G4'], '2n');
```

##### Tone.MembraneSynth

For kick drums.

```typescript
const kick = new Tone.MembraneSynth({
  pitchDecay: 0.05,
  octaves: 10,
  oscillator: { type: 'sine' }
}).toDestination();

kick.triggerAttackRelease('C1', '8n');
```

##### Tone.NoiseSynth

For snare drums.

```typescript
const snare = new Tone.NoiseSynth({
  noise: { type: 'white' },
  envelope: { attack: 0.005, decay: 0.1, sustain: 0 }
}).toDestination();

snare.triggerAttack();
```

##### Tone.MetalSynth

For hi-hats.

```typescript
const hihat = new Tone.MetalSynth({
  frequency: 200,
  envelope: { attack: 0.001, decay: 0.1, release: 0.01 }
}).toDestination();

hihat.triggerAttack();
```

##### Tone.Recorder

Record audio output.

```typescript
const recorder = new Tone.Recorder();
synth.connect(recorder);

recorder.start();
// ... play audio ...
const recording = await recorder.stop();

// Download
const url = URL.createObjectURL(recording);
const a = document.createElement('a');
a.href = url;
a.download = 'recording.webm';
a.click();
```

##### Tone.Part

Schedule multiple events.

```typescript
const part = new Tone.Part((time, note) => {
  synth.triggerAttackRelease(note, '8n', time);
}, [
  ['0:0', 'C4'],
  ['0:2', 'E4'],
  ['1:0', 'G4']
]);

part.start(0);
Tone.Transport.start();
```

##### Tone.Pattern

Create repeating patterns.

```typescript
const pattern = new Tone.Pattern((time, note) => {
  synth.triggerAttackRelease(note, '8n', time);
}, ['C4', 'D4', 'E4', 'F4']);

pattern.start(0);
Tone.Transport.start();
```

---

### pitchfinder

**Documentation:** https://github.com/peterkhayes/pitchfinder

#### Algorithms

##### YIN (Recommended)

Best balance of accuracy and speed.

```typescript
import Pitchfinder from 'pitchfinder';

const detectPitch = Pitchfinder.YIN({
  sampleRate: 44100,
  threshold: 0.1
});

const pitch = detectPitch(float32Array);
if (pitch) {
  console.log('Detected frequency:', pitch, 'Hz');
}
```

##### AMDF

Slower but more consistent.

```typescript
const detectPitch = Pitchfinder.AMDF({
  sampleRate: 44100
});
```

##### Dynamic Wavelet

Very fast for high frequencies.

```typescript
const detectPitch = Pitchfinder.DynamicWavelet({
  sampleRate: 44100
});
```

#### Multiple Detectors

Use multiple detectors for better accuracy.

```typescript
const detectors = [
  Pitchfinder.YIN({ sampleRate: 44100 }),
  Pitchfinder.AMDF({ sampleRate: 44100 })
];

const frequencies = Pitchfinder.frequencies(
  detectors,
  float32Array,
  {
    tempo: 120,
    quantization: 4 // 16th notes
  }
);
```

---

### realtime-bpm-analyzer

**Documentation:** https://www.realtime-bpm-analyzer.com/

#### Basic Usage

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
source.connect(audioContext.destination);

// Listen for BPM
bpmProcessor.port.onmessage = (event) => {
  if (event.data.message === 'BPM') {
    console.log('Current BPM:', event.data.data.bpm);
  }
  
  if (event.data.message === 'BPM_STABLE') {
    console.log('Stable BPM:', event.data.data.bpm);
  }
};

source.start();
```

#### With Biquad Filter

```typescript
import { getBiquadFilter } from 'realtime-bpm-analyzer';

const lowpass = getBiquadFilter(audioContext);

source
  .connect(lowpass)
  .connect(bpmProcessor)
  .connect(audioContext.destination);
```

---

### Tonal.js

**Documentation:** https://github.com/tonaljs/tonal

#### Note Operations

```typescript
import { Note } from 'tonal';

// MIDI conversion
Note.midi('A4');        // 69
Note.fromMidi(60);      // "C4"
Note.fromMidiSharps(61); // "C#4"

// Frequency
Note.freq('A4');        // 440
Note.freqToMidi(440);   // 69

// Properties
Note.get('C#4');        // { name: "C#4", pc: "C#", oct: 4, ... }
Note.octave('C4');      // 4
Note.pc('C#4');         // "C#"
Note.chroma('C');       // 0 (C=0, C#=1, D=2, etc.)

// Transposition
Note.transpose('C4', '5P');  // "G4"
Note.transpose('C4', 'M3');  // "E4"
```

#### Key Detection

```typescript
import { Key } from 'tonal';

// Get key information
const key = Key.majorKey('C');
console.log(key.tonic);           // "C"
console.log(key.scale);           // ["C", "D", "E", "F", "G", "A", "B"]
console.log(key.chords);          // ["Cmaj7", "Dm7", "Em7", ...]
console.log(key.chordsHarmonicFunction); // ["T", "SD", "T", ...]

// Minor key
const minorKey = Key.minorKey('A');
console.log(minorKey.scale);      // ["A", "B", "C", "D", "E", "F", "G"]
console.log(minorKey.natural);    // Natural minor scale
console.log(minorKey.harmonic);   // Harmonic minor scale
console.log(minorKey.melodic);    // Melodic minor scale
```

#### Scale Operations

```typescript
import { Scale } from 'tonal';

// Get scale
const scale = Scale.get('C major');
console.log(scale.notes);         // ["C", "D", "E", "F", "G", "A", "B"]
console.log(scale.intervals);     // ["1P", "2M", "3M", "4P", "5P", "6M", "7M"]

// Scale degrees
Scale.degrees('C major')(['1', '3', '5']);  // ["C", "E", "G"]

// Detect scale
Scale.detect(['C', 'D', 'E', 'F', 'G', 'A', 'B']);  // ["C major", "G mixolydian", ...]
```

#### Chord Operations

```typescript
import { Chord } from 'tonal';

// Get chord
const chord = Chord.get('Cmaj7');
console.log(chord.notes);         // ["C", "E", "G", "B"]
console.log(chord.intervals);     // ["1P", "3M", "5P", "7M"]

// Detect chord
Chord.detect(['C', 'E', 'G']);    // ["C", "Cm", "Am", ...]

// Chord degrees
Chord.degrees('Cmaj7')(['1', '3', '5', '7']);  // ["C", "E", "G", "B"]
```

#### Interval Operations

```typescript
import { Interval } from 'tonal';

// Interval properties
Interval.semitones('5P');         // 7
Interval.semitones('M3');         // 4

// Interval from notes
Interval.distance('C4', 'G4');    // "5P"
Interval.distance('C4', 'E4');    // "M3"

// Add/subtract intervals
Interval.add('M3', 'M3');         // "A5"
Interval.substract('5P', 'M3');   // "m3"
```

---

### @tonejs/midi

**Documentation:** https://tonejs.github.io/midi/

#### Create MIDI

```typescript
import { Midi } from '@tonejs/midi';

// Create new MIDI file
const midi = new Midi();

// Set tempo
midi.header.setTempo(120);

// Add track
const track = midi.addTrack();
track.name = 'Melody';

// Add notes
track.addNote({
  midi: 60,           // C4
  time: 0,            // Start time in seconds
  duration: 0.5,      // Duration in seconds
  velocity: 0.8       // 0-1
});

track.addNote({
  midi: 64,           // E4
  time: 0.5,
  duration: 0.5,
  velocity: 0.8
});

// Export
const midiArray = midi.toArray();
const blob = new Blob([midiArray], { type: 'audio/midi' });
```

#### Read MIDI

```typescript
import { Midi } from '@tonejs/midi';

// Load MIDI file
const midiFile = await Midi.fromUrl('path/to/file.mid');

console.log('Tempo:', midiFile.header.tempos[0].bpm);
console.log('Tracks:', midiFile.tracks.length);

midiFile.tracks.forEach(track => {
  console.log('Track:', track.name);
  track.notes.forEach(note => {
    console.log(`Note: ${note.name}, Time: ${note.time}s, Duration: ${note.duration}s`);
  });
});
```

---

### WaveSurfer.js

**Documentation:** https://wavesurfer-js.org/

#### Basic Usage

```typescript
import WaveSurfer from 'wavesurfer.js';

const wavesurfer = WaveSurfer.create({
  container: '#waveform',
  waveColor: '#3B82F6',
  progressColor: '#8B5CF6',
  cursorColor: '#ffffff',
  barWidth: 2,
  barRadius: 3,
  cursorWidth: 1,
  height: 128,
  barGap: 2
});

// Load audio
wavesurfer.load('path/to/audio.mp3');

// Or load from AudioBuffer
wavesurfer.loadBlob(audioBlob);

// Or from audio element
wavesurfer.loadMediaElement(audioElement);

// Controls
wavesurfer.play();
wavesurfer.pause();
wavesurfer.stop();

// Events
wavesurfer.on('ready', () => {
  console.log('Waveform ready');
});

wavesurfer.on('play', () => {
  console.log('Playing');
});

wavesurfer.on('finish', () => {
  console.log('Finished');
});
```

---

## TYPE DEFINITIONS

See `lib/types/index.ts` for complete type definitions.

### PitchPoint

```typescript
interface PitchPoint {
  frequency: number;   // Hz
  time: number;        // seconds
  midi: number;        // 0-127
  confidence?: number; // 0-1
}
```

### BPMAnalysis

```typescript
interface BPMAnalysis {
  bpm: number;         // beats per minute
  confidence: number;  // 0-1
  stable: boolean;     // true if detection is stable
}
```

### KeyAnalysis

```typescript
interface KeyAnalysis {
  key: string;                // "C Major", "A Minor"
  tonic: string;              // "C", "A"
  scale: string[];            // ["C", "D", "E", "F", "G", "A", "B"]
  chordProgression?: string[]; // ["I", "V", "vi", "IV"]
}
```

### Section

```typescript
type InstrumentType = 'drums' | 'bass' | 'chords';

interface Section {
  id: string;
  name: string;
  audioBuffer: AudioBuffer | null;
  duration: number;
  pitches: PitchPoint[];
  bpm: BPMAnalysis | null;
  key: KeyAnalysis | null;
  instruments: InstrumentType[];
  generatedParts?: {
    drums?: any;
    bass?: any;
    chords?: any;
  };
}
```

### Project

```typescript
interface Project {
  name: string;
  sections: Section[];
  arrangementMode: 'sequential' | 'layered';
  masterBPM: number;
  masterKey: string;
  createdAt: Date;
  duration: number;
}
```

---

## UTILITY FUNCTIONS

### Audio Utils

**Location:** `lib/utils/audio-utils.ts`

#### frequencyToMidi

```typescript
function frequencyToMidi(frequency: number): number {
  return Math.round(69 + 12 * Math.log2(frequency / 440));
}

// Example
frequencyToMidi(440);  // 69 (A4)
frequencyToMidi(261.6); // 60 (C4)
```

#### midiToFrequency

```typescript
function midiToFrequency(midi: number): number {
  return 440 * Math.pow(2, (midi - 69) / 12);
}

// Example
midiToFrequency(69);  // 440 (A4)
midiToFrequency(60);  // 261.6 (C4)
```

#### downloadBlob

```typescript
function downloadBlob(blob: Blob, filename: string): void {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

// Example
downloadBlob(audioBlob, 'recording.wav');
```

#### formatDuration

```typescript
function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

// Example
formatDuration(125);  // "2:05"
formatDuration(65);   // "1:05"
```

### Music Theory Utils

**Location:** `lib/utils/music-theory.ts`

#### getChordProgression

```typescript
function getChordProgression(key: string): string[][] {
  // Returns I-V-vi-IV progression
  // Returns array of chord notes
}

// Example
getChordProgression('C Major');
// [
//   ['C', 'E', 'G'],  // I (C)
//   ['G', 'B', 'D'],  // V (G)
//   ['A', 'C', 'E'],  // vi (Am)
//   ['F', 'A', 'C']   // IV (F)
// ]
```

#### transposeNote

```typescript
function transposeNote(note: string, semitones: number): string {
  // Transpose a note by semitones
}

// Example
transposeNote('C4', 4);  // "E4"
transposeNote('A3', -2); // "G3"
```

---

## CONFIGURATION

### Tailwind Config

**Location:** `tailwind.config.ts`

```typescript
const config: Config = {
  content: ['./app/**/*.{js,ts,jsx,tsx,mdx}'],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#EEF2FF',
          500: '#3B82F6',
          600: '#2563EB',
        },
        secondary: {
          500: '#8B5CF6',
          600: '#7C3AED',
        }
      }
    }
  }
}
```

### Next.js Config

**Location:** `next.config.js`

```javascript
const nextConfig = {
  reactStrictMode: true,
  webpack: (config) => {
    config.module.rules.push({
      test: /\.(mp3|wav|ogg)$/,
      type: 'asset/resource'
    });
    return config;
  }
}
```

### Environment Variables

**Location:** `.env.local`

```bash
# Optional: For lyrics feature
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
```

---

**END OF API REFERENCE**
