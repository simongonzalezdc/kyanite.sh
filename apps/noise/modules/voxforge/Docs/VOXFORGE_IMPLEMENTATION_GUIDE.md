# VOXFORGE - IMPLEMENTATION GUIDE FOR AI CODING AGENTS
**Version:** 2.0  
**Last Updated:** November 9, 2025  
**Purpose:** Step-by-step implementation instructions with complete code examples

---

## TABLE OF CONTENTS

1. [Project Setup](#project-setup)
2. [Day 1: Audio Recording](#day-1-audio-recording)
3. [Day 2: Pitch Detection](#day-2-pitch-detection)
4. [Day 3: BPM & Key Detection](#day-3-bpm--key-detection)
5. [Day 4: Drum Generator](#day-4-drum-generator)
6. [Day 5: Bass & Chords](#day-5-bass--chords)
7. [Day 6: Instrument Selection](#day-6-instrument-selection)
8. [Day 7: Section Management](#day-7-section-management)
9. [Day 8: Arrangement Modes](#day-8-arrangement-modes)
10. [Day 9: Stem Export](#day-9-stem-export)
11. [Day 10: MIDI Export](#day-10-midi-export)
12. [Day 11: Lyrics Assistant](#day-11-lyrics-assistant-optional)
13. [Day 12: Testing & Deployment](#day-12-testing--deployment)

---

## PROJECT SETUP

### Step 1: Initialize Next.js Project

```bash
npx create-next-app@latest voxforge --typescript --tailwind --app --no-src-dir
cd voxforge
```

**Selections when prompted:**
- TypeScript: Yes
- ESLint: Yes
- Tailwind CSS: Yes
- App Router: Yes
- Import alias: No (use default)

### Step 2: Install Dependencies

```bash
# Audio processing
npm install tone pitchfinder realtime-bpm-analyzer tonal @tonejs/midi

# UI & Visualization
npm install wavesurfer.js lucide-react

# Type definitions
npm install -D @types/node

# Optional: Lyrics feature
npm install openai
```

### Step 3: Create Folder Structure

```bash
mkdir -p app/components
mkdir -p lib/audio
mkdir -p lib/types
mkdir -p lib/utils
mkdir -p public/test-recordings
mkdir -p docs
```

### Step 4: Configure TypeScript

Create or update `tsconfig.json`:

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "lib": ["ES2020", "DOM", "DOM.Iterable", "webworker"],
    "allowJs": true,
    "skipLibCheck": true,
    "strict": true,
    "noEmit": true,
    "esModuleInterop": true,
    "module": "esnext",
    "moduleResolution": "bundler",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "jsx": "preserve",
    "incremental": true,
    "plugins": [
      {
        "name": "next"
      }
    ],
    "paths": {
      "@/*": ["./*"]
    }
  },
  "include": ["next-env.d.ts", "**/*.ts", "**/*.tsx", ".next/types/**/*.ts"],
  "exclude": ["node_modules"]
}
```

### Step 5: Configure Tailwind

Update `tailwind.config.ts`:

```typescript
import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#EEF2FF',
          100: '#E0E7FF',
          500: '#3B82F6',
          600: '#2563EB',
          700: '#1D4ED8',
        },
        secondary: {
          500: '#8B5CF6',
          600: '#7C3AED',
        }
      },
      animation: {
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      }
    },
  },
  plugins: [],
}
export default config
```

### Step 6: Update Global Styles

Update `app/globals.css`:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

:root {
  --background: #0a0a0a;
  --foreground: #ffffff;
  --primary: #3B82F6;
  --secondary: #8B5CF6;
}

body {
  color: var(--foreground);
  background: var(--background);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #1a1a1a;
}

::-webkit-scrollbar-thumb {
  background: #3B82F6;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #2563EB;
}
```

---

## DAY 1: AUDIO RECORDING

### Goal: Record voice and display waveform

### Step 1: Create Type Definitions

Create `lib/types/index.ts`:

```typescript
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
  key: string;          // "C Major", "A Minor"
  tonic: string;        // "C", "A"
  scale: string[];      // ["C", "D", "E", "F", "G", "A", "B"]
  chordProgression?: string[];
}

export type InstrumentType = 'drums' | 'bass' | 'chords';

export interface Section {
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

export interface Project {
  name: string;
  sections: Section[];
  arrangementMode: 'sequential' | 'layered';
  masterBPM: number;
  masterKey: string;
  createdAt: Date;
  duration: number;
}
```

### Step 2: Create Audio Recorder

Create `lib/audio/recorder.ts`:

```typescript
export class AudioRecorder {
  private mediaRecorder: MediaRecorder | null = null;
  private audioChunks: Blob[] = [];
  private audioContext: AudioContext;
  private analyser: AnalyserNode;
  private stream: MediaStream | null = null;

  constructor() {
    this.audioContext = new AudioContext({ sampleRate: 44100 });
    this.analyser = this.audioContext.createAnalyser();
    this.analyser.fftSize = 2048;
  }

  async requestPermission(): Promise<boolean> {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({
        audio: {
          sampleRate: 44100,
          channelCount: 1,
          echoCancellation: false,
          noiseSuppression: false,
          autoGainControl: false
        }
      });
      
      // Test and close immediately
      stream.getTracks().forEach(track => track.stop());
      return true;
    } catch (error) {
      console.error('Microphone permission denied:', error);
      return false;
    }
  }

  async startRecording(): Promise<void> {
    this.stream = await navigator.mediaDevices.getUserMedia({
      audio: {
        sampleRate: 44100,
        channelCount: 1,
        echoCancellation: false,
        noiseSuppression: false,
        autoGainControl: false
      }
    });

    this.mediaRecorder = new MediaRecorder(this.stream, {
      mimeType: 'audio/webm'
    });
    
    this.audioChunks = [];

    this.mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        this.audioChunks.push(event.data);
      }
    };

    // Connect to analyser for real-time visualization
    const source = this.audioContext.createMediaStreamSource(this.stream);
    source.connect(this.analyser);

    this.mediaRecorder.start(100); // Collect data every 100ms
  }

  stopRecording(): Promise<AudioBuffer> {
    return new Promise((resolve, reject) => {
      if (!this.mediaRecorder) {
        reject(new Error('No recording in progress'));
        return;
      }

      this.mediaRecorder.onstop = async () => {
        try {
          // Create blob from chunks
          const audioBlob = new Blob(this.audioChunks, { type: 'audio/webm' });
          
          // Convert to array buffer
          const arrayBuffer = await audioBlob.arrayBuffer();
          
          // Decode to audio buffer
          const audioBuffer = await this.audioContext.decodeAudioData(arrayBuffer);
          
          // Cleanup
          this.cleanup();
          
          resolve(audioBuffer);
        } catch (error) {
          reject(error);
        }
      };

      this.mediaRecorder.stop();
    });
  }

  getWaveformData(): Uint8Array {
    const dataArray = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteTimeDomainData(dataArray);
    return dataArray;
  }

  getFrequencyData(): Uint8Array {
    const dataArray = new Uint8Array(this.analyser.frequencyBinCount);
    this.analyser.getByteFrequencyData(dataArray);
    return dataArray;
  }

  private cleanup(): void {
    if (this.stream) {
      this.stream.getTracks().forEach(track => track.stop());
      this.stream = null;
    }
  }

  dispose(): void {
    this.cleanup();
    this.audioContext.close();
  }
}
```

### Step 3: Create Waveform Visualizer

Create `app/components/Waveform.tsx`:

```typescript
'use client';

import { useEffect, useRef } from 'react';

interface WaveformProps {
  isRecording: boolean;
  getWaveformData?: () => Uint8Array;
  audioBuffer?: AudioBuffer | null;
}

export default function Waveform({ isRecording, getWaveformData, audioBuffer }: WaveformProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationRef = useRef<number>();

  useEffect(() => {
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    // Set canvas size
    canvas.width = canvas.offsetWidth * 2; // Retina
    canvas.height = canvas.offsetHeight * 2;
    ctx.scale(2, 2);

    if (isRecording && getWaveformData) {
      // Real-time visualization
      const draw = () => {
        const data = getWaveformData();
        const width = canvas.offsetWidth;
        const height = canvas.offsetHeight;

        // Clear canvas
        ctx.fillStyle = '#0a0a0a';
        ctx.fillRect(0, 0, width, height);

        // Draw waveform
        ctx.lineWidth = 2;
        ctx.strokeStyle = '#3B82F6';
        ctx.beginPath();

        const sliceWidth = width / data.length;
        let x = 0;

        for (let i = 0; i < data.length; i++) {
          const v = data[i] / 128.0; // Normalize to 0-2
          const y = (v * height) / 2;

          if (i === 0) {
            ctx.moveTo(x, y);
          } else {
            ctx.lineTo(x, y);
          }

          x += sliceWidth;
        }

        ctx.stroke();

        animationRef.current = requestAnimationFrame(draw);
      };

      draw();

      return () => {
        if (animationRef.current) {
          cancelAnimationFrame(animationRef.current);
        }
      };
    } else if (audioBuffer) {
      // Static visualization of recorded audio
      const data = audioBuffer.getChannelData(0);
      const width = canvas.offsetWidth;
      const height = canvas.offsetHeight;

      ctx.fillStyle = '#0a0a0a';
      ctx.fillRect(0, 0, width, height);

      ctx.lineWidth = 1;
      ctx.strokeStyle = '#8B5CF6';
      ctx.beginPath();

      const step = Math.ceil(data.length / width);
      const amp = height / 2;

      for (let i = 0; i < width; i++) {
        let min = 1.0;
        let max = -1.0;

        for (let j = 0; j < step; j++) {
          const datum = data[(i * step) + j];
          if (datum < min) min = datum;
          if (datum > max) max = datum;
        }

        ctx.moveTo(i, (1 + min) * amp);
        ctx.lineTo(i, (1 + max) * amp);
      }

      ctx.stroke();
    }
  }, [isRecording, getWaveformData, audioBuffer]);

  return (
    <canvas
      ref={canvasRef}
      className="w-full h-32 rounded-lg bg-black/50 border border-gray-800"
    />
  );
}
```

### Step 4: Create Recorder Component

Create `app/components/Recorder.tsx`:

```typescript
'use client';

import { useState, useRef } from 'react';
import { AudioRecorder } from '@/lib/audio/recorder';
import Waveform from './Waveform';
import { Mic, Square, Play } from 'lucide-react';

interface RecorderProps {
  onRecordingComplete: (audioBuffer: AudioBuffer) => void;
}

export default function Recorder({ onRecordingComplete }: RecorderProps) {
  const [isRecording, setIsRecording] = useState(false);
  const [hasPermission, setHasPermission] = useState(false);
  const [recordedAudio, setRecordedAudio] = useState<AudioBuffer | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  
  const recorderRef = useRef<AudioRecorder | null>(null);
  const audioContextRef = useRef<AudioContext | null>(null);

  const requestPermission = async () => {
    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }
    
    const granted = await recorderRef.current.requestPermission();
    setHasPermission(granted);
    
    if (!granted) {
      alert('Microphone permission is required to record audio.');
    }
  };

  const startRecording = async () => {
    if (!hasPermission) {
      await requestPermission();
      if (!hasPermission) return;
    }

    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }

    try {
      await recorderRef.current.startRecording();
      setIsRecording(true);
      setRecordedAudio(null);
    } catch (error) {
      console.error('Failed to start recording:', error);
      alert('Failed to start recording. Please check microphone permissions.');
    }
  };

  const stopRecording = async () => {
    if (!recorderRef.current || !isRecording) return;

    try {
      const audioBuffer = await recorderRef.current.stopRecording();
      setIsRecording(false);
      setRecordedAudio(audioBuffer);
      onRecordingComplete(audioBuffer);
    } catch (error) {
      console.error('Failed to stop recording:', error);
      alert('Failed to process recording.');
    }
  };

  const playRecording = async () => {
    if (!recordedAudio) return;

    if (!audioContextRef.current) {
      audioContextRef.current = new AudioContext();
    }

    const source = audioContextRef.current.createBufferSource();
    source.buffer = recordedAudio;
    source.connect(audioContextRef.current.destination);
    
    source.onended = () => {
      setIsPlaying(false);
    };

    source.start();
    setIsPlaying(true);
  };

  return (
    <div className="space-y-4">
      <Waveform
        isRecording={isRecording}
        getWaveformData={isRecording ? () => recorderRef.current!.getWaveformData() : undefined}
        audioBuffer={recordedAudio}
      />

      <div className="flex gap-4 justify-center">
        {!isRecording ? (
          <button
            onClick={startRecording}
            className="flex items-center gap-2 px-6 py-3 bg-primary-500 hover:bg-primary-600 rounded-lg font-medium transition-colors"
          >
            <Mic size={20} />
            Start Recording
          </button>
        ) : (
          <button
            onClick={stopRecording}
            className="flex items-center gap-2 px-6 py-3 bg-red-500 hover:bg-red-600 rounded-lg font-medium transition-colors"
          >
            <Square size={20} />
            Stop Recording
          </button>
        )}

        {recordedAudio && !isRecording && (
          <button
            onClick={playRecording}
            disabled={isPlaying}
            className="flex items-center gap-2 px-6 py-3 bg-secondary-500 hover:bg-secondary-600 rounded-lg font-medium transition-colors disabled:opacity-50"
          >
            <Play size={20} />
            Play Back
          </button>
        )}
      </div>

      {recordedAudio && (
        <div className="text-center text-sm text-gray-400">
          Duration: {recordedAudio.duration.toFixed(2)}s
        </div>
      )}
    </div>
  );
}
```

### Step 5: Update Main Page

Update `app/page.tsx`:

```typescript
'use client';

import { useState } from 'react';
import Recorder from './components/Recorder';

export default function Home() {
  const [audioBuffer, setAudioBuffer] = useState<AudioBuffer | null>(null);

  const handleRecordingComplete = (buffer: AudioBuffer) => {
    setAudioBuffer(buffer);
    console.log('Recording complete:', buffer);
  };

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <header className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            VoxForge
          </h1>
          <p className="text-gray-400">
            Transform your voice into music
          </p>
        </header>

        <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
          <h2 className="text-xl font-semibold mb-4">Step 1: Record Your Voice</h2>
          <Recorder onRecordingComplete={handleRecordingComplete} />
        </div>

        {audioBuffer && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
            <h2 className="text-xl font-semibold mb-4">Audio Recorded!</h2>
            <p className="text-gray-400">
              Next: Analysis will be added in Day 2
            </p>
          </div>
        )}
      </div>
    </main>
  );
}
```

### Day 1 Testing Checklist

- [ ] Click "Start Recording" button
- [ ] Browser prompts for microphone permission
- [ ] Grant permission
- [ ] Waveform displays in real-time
- [ ] Speak/sing into microphone
- [ ] Click "Stop Recording"
- [ ] Static waveform displays
- [ ] "Play Back" button appears
- [ ] Click play, hear recording
- [ ] Console logs audio buffer details

---

## DAY 2: PITCH DETECTION

### Goal: Extract melody as MIDI notes

### Step 1: Create Pitch Detector

Create `lib/audio/pitch-detector.ts`:

```typescript
import Pitchfinder from 'pitchfinder';
import { PitchPoint } from '../types';

export class PitchDetector {
  private detectPitch: (input: Float32Array) => number | null;
  private sampleRate: number;

  constructor(sampleRate: number = 44100) {
    this.sampleRate = sampleRate;
    
    // YIN algorithm - best for voice
    this.detectPitch = Pitchfinder.YIN({
      sampleRate: this.sampleRate,
      threshold: 0.1 // Lower = more sensitive
    });
  }

  analyze(audioBuffer: AudioBuffer): PitchPoint[] {
    const channelData = audioBuffer.getChannelData(0);
    const pitches: PitchPoint[] = [];
    
    const windowSize = 2048;
    const hopSize = 512; // ~11.6ms at 44.1kHz
    
    console.log('Starting pitch analysis...');
    console.log('Audio duration:', audioBuffer.duration, 'seconds');
    console.log('Sample rate:', this.sampleRate);
    console.log('Total samples:', channelData.length);

    for (let i = 0; i < channelData.length - windowSize; i += hopSize) {
      const window = channelData.slice(i, i + windowSize);
      const frequency = this.detectPitch(window);

      if (frequency && frequency > 0 && frequency < 2000) {
        // Valid human voice range: ~80Hz - 2000Hz
        const time = i / this.sampleRate;
        const midi = this.frequencyToMidi(frequency);

        pitches.push({
          frequency,
          time,
          midi,
          confidence: 1.0 // YIN doesn't provide confidence
        });
      }
    }

    console.log('Pitch analysis complete. Detected', pitches.length, 'pitch points');
    return pitches;
  }

  private frequencyToMidi(frequency: number): number {
    // MIDI note 69 (A4) = 440 Hz
    return Math.round(69 + 12 * Math.log2(frequency / 440));
  }

  midiToNoteName(midi: number): string {
    const noteNames = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'];
    const octave = Math.floor(midi / 12) - 1;
    const note = noteNames[midi % 12];
    return `${note}${octave}`;
  }

  // Get simplified melody (quantized to nearest note)
  getSimplifiedMelody(pitches: PitchPoint[]): number[] {
    if (pitches.length === 0) return [];

    // Remove duplicate consecutive notes
    const simplified: number[] = [];
    let lastNote = -1;

    pitches.forEach(pitch => {
      const note = Math.round(pitch.midi);
      if (note !== lastNote) {
        simplified.push(note);
        lastNote = note;
      }
    });

    return simplified;
  }

  // Get note statistics
  getStats(pitches: PitchPoint[]) {
    if (pitches.length === 0) {
      return {
        averageFrequency: 0,
        minFrequency: 0,
        maxFrequency: 0,
        averageMidi: 0,
        range: 'unknown' as 'low' | 'mid' | 'high' | 'unknown'
      };
    }

    const frequencies = pitches.map(p => p.frequency);
    const midis = pitches.map(p => p.midi);

    const avgFreq = frequencies.reduce((a, b) => a + b, 0) / frequencies.length;
    const avgMidi = midis.reduce((a, b) => a + b, 0) / midis.length;

    let range: 'low' | 'mid' | 'high' = 'mid';
    if (avgMidi < 60) range = 'low';  // Below C4
    else if (avgMidi > 72) range = 'high'; // Above C5

    return {
      averageFrequency: avgFreq,
      minFrequency: Math.min(...frequencies),
      maxFrequency: Math.max(...frequencies),
      averageMidi: avgMidi,
      range
    };
  }
}
```

### Step 2: Create Analysis Display Component

Create `app/components/AnalysisDisplay.tsx`:

```typescript
'use client';

import { PitchPoint } from '@/lib/types';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { Music, Activity } from 'lucide-react';

interface AnalysisDisplayProps {
  pitches: PitchPoint[];
  bpm?: number | null;
  musicalKey?: string | null;
}

export default function AnalysisDisplay({ pitches, bpm, musicalKey }: AnalysisDisplayProps) {
  if (pitches.length === 0) {
    return (
      <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
        <p className="text-gray-400 text-center">No analysis data yet</p>
      </div>
    );
  }

  const detector = new PitchDetector();
  const stats = detector.getStats(pitches);
  const simplified = detector.getSimplifiedMelody(pitches);
  const noteNames = simplified.slice(0, 20).map(midi => detector.midiToNoteName(midi));

  return (
    <div className="bg-gray-900 rounded-xl p-8 border border-gray-800 space-y-6">
      <h2 className="text-xl font-semibold flex items-center gap-2">
        <Activity size={24} className="text-primary-500" />
        Audio Analysis
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Pitch Info */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <div className="flex items-center gap-2 mb-2">
            <Music size={20} className="text-secondary-500" />
            <h3 className="font-medium">Pitch Range</h3>
          </div>
          <p className="text-2xl font-bold text-secondary-500 capitalize">{stats.range}</p>
          <p className="text-sm text-gray-400 mt-1">
            {stats.averageFrequency.toFixed(1)} Hz average
          </p>
        </div>

        {/* BPM */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <h3 className="font-medium mb-2">Tempo</h3>
          <p className="text-2xl font-bold text-primary-500">
            {bpm ? `${bpm} BPM` : 'Detecting...'}
          </p>
          <p className="text-sm text-gray-400 mt-1">
            {bpm ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>

        {/* Key */}
        <div className="bg-black/30 rounded-lg p-4 border border-gray-800">
          <h3 className="font-medium mb-2">Musical Key</h3>
          <p className="text-2xl font-bold text-primary-500">
            {musicalKey || 'Detecting...'}
          </p>
          <p className="text-sm text-gray-400 mt-1">
            {musicalKey ? 'Auto-detected' : 'Analysis pending'}
          </p>
        </div>
      </div>

      {/* Detected Notes */}
      <div>
        <h3 className="font-medium mb-3">Detected Melody (first 20 notes)</h3>
        <div className="flex flex-wrap gap-2">
          {noteNames.map((note, i) => (
            <span
              key={i}
              className="px-3 py-1 bg-primary-500/20 border border-primary-500/50 rounded-md text-sm font-mono"
            >
              {note}
            </span>
          ))}
          {simplified.length > 20 && (
            <span className="px-3 py-1 text-gray-400 text-sm">
              +{simplified.length - 20} more
            </span>
          )}
        </div>
      </div>

      {/* Stats */}
      <div className="text-sm text-gray-400 space-y-1">
        <p>Total pitch points detected: {pitches.length}</p>
        <p>Unique notes: {simplified.length}</p>
        <p>Frequency range: {stats.minFrequency.toFixed(1)} - {stats.maxFrequency.toFixed(1)} Hz</p>
      </div>
    </div>
  );
}
```

### Step 3: Update Main Page for Analysis

Update `app/page.tsx` to add pitch analysis:

```typescript
'use client';

import { useState } from 'react';
import Recorder from './components/Recorder';
import AnalysisDisplay from './components/AnalysisDisplay';
import { PitchDetector } from '@/lib/audio/pitch-detector';
import { PitchPoint } from '@/lib/types';

export default function Home() {
  const [audioBuffer, setAudioBuffer] = useState<AudioBuffer | null>(null);
  const [pitches, setPitches] = useState<PitchPoint[]>([]);
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  const handleRecordingComplete = async (buffer: AudioBuffer) => {
    setAudioBuffer(buffer);
    setIsAnalyzing(true);

    // Run pitch detection
    try {
      const detector = new PitchDetector(buffer.sampleRate);
      const detectedPitches = detector.analyze(buffer);
      setPitches(detectedPitches);
    } catch (error) {
      console.error('Pitch detection failed:', error);
      alert('Failed to analyze pitch. Please try again.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-4xl mx-auto space-y-8">
        <header className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            VoxForge
          </h1>
          <p className="text-gray-400">
            Transform your voice into music
          </p>
        </header>

        <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
          <h2 className="text-xl font-semibold mb-4">Step 1: Record Your Voice</h2>
          <Recorder onRecordingComplete={handleRecordingComplete} />
        </div>

        {isAnalyzing && (
          <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
            <div className="flex items-center justify-center gap-3">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500"></div>
              <p className="text-gray-400">Analyzing audio...</p>
            </div>
          </div>
        )}

        {pitches.length > 0 && !isAnalyzing && (
          <>
            <AnalysisDisplay pitches={pitches} />
            
            <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
              <p className="text-gray-400 text-center">
                Next: BPM and Key detection will be added in Day 3
              </p>
            </div>
          </>
        )}
      </div>
    </main>
  );
}
```

### Day 2 Testing Checklist

- [ ] Record voice saying "Do Re Mi Fa Sol La Ti Do"
- [ ] See "Analyzing audio..." message
- [ ] Analysis completes within 2 seconds
- [ ] Detected notes show: C, D, E, F, G, A, B, C (or similar)
- [ ] Pitch range shows correctly (low/mid/high)
- [ ] Average frequency displays
- [ ] Console shows pitch detection logs

---

## DAY 3: BPM & KEY DETECTION

*[Content continues with BPM detection implementation using realtime-bpm-analyzer, and key detection using tonal.js - similar detailed format as above]*

---

**Note:** This document would continue with Days 4-12 following the same detailed format. Each day would include:
- Complete code implementations
- Step-by-step instructions
- Testing checklists
- Console commands
- Troubleshooting tips

The full document would be approximately 3000-4000 lines covering all 12 days of development.

**Would you like me to:**
1. Continue with the remaining days (4-12)?
2. Focus on a specific day in more detail?
3. Create a separate file for each day's implementation?

---

**END OF DAY 1-2 IMPLEMENTATION GUIDE**
