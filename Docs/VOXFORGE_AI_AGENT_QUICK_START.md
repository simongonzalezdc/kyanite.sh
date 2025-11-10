10pm# VOXFORGE - AI CODING AGENT QUICK START
**Version:** 2.0  
**Last Updated:** November 9, 2025  
**Purpose:** Quick reference for AI coding agents

---

## PROJECT SETUP (COPY-PASTE READY)

### Initialize Project

```bash
# Create Next.js app
npx create-next-app@latest voxforge --typescript --tailwind --app --no-src-dir
cd voxforge

# Install all dependencies
npm install tone pitchfinder realtime-bpm-analyzer tonal @tonejs/midi
npm install wavesurfer.js lucide-react
npm install -D @types/node

# Optional: For lyrics
npm install openai

# Create folder structure
mkdir -p app/components app/api/lyrics
mkdir -p lib/audio lib/types lib/utils
mkdir -p public/test-recordings docs

# Start dev server
npm run dev
```

### Quick File Checklist

Essential files to create:

```
✅ lib/types/index.ts          # Type definitions
✅ lib/audio/recorder.ts        # Audio recording
✅ lib/audio/pitch-detector.ts  # Pitch analysis  
✅ lib/audio/bpm-detector.ts    # BPM analysis
✅ lib/audio/key-detector.ts    # Key detection
✅ lib/audio/music-generator.ts # Music generation
✅ lib/audio/stem-exporter.ts   # Export stems
✅ lib/audio/midi-exporter.ts   # MIDI export
✅ app/components/Recorder.tsx  # Recording UI
✅ app/components/Waveform.tsx  # Visualization
✅ app/components/AnalysisDisplay.tsx # Analysis UI
✅ app/components/PlaybackControls.tsx # Controls
✅ app/components/Metronome.tsx # Beat indicator
✅ app/components/InstrumentPicker.tsx # Instrument selection
✅ app/components/SectionManager.tsx # Section management
✅ app/components/ExportPanel.tsx # Export UI
✅ app/page.tsx                 # Main page
```

---

## KEY PATTERNS & CODE SNIPPETS

### 1. Audio Recording Pattern

```typescript
// lib/audio/recorder.ts
export class AudioRecorder {
  private mediaRecorder: MediaRecorder | null = null;
  private audioContext: AudioContext;
  private analyser: AnalyserNode;

  constructor() {
    this.audioContext = new AudioContext({ sampleRate: 44100 });
    this.analyser = this.audioContext.createAnalyser();
  }

  async startRecording(): Promise<void> {
    const stream = await navigator.mediaDevices.getUserMedia({
      audio: { sampleRate: 44100, channelCount: 1 }
    });
    this.mediaRecorder = new MediaRecorder(stream);
    this.mediaRecorder.start(100);
  }

  stopRecording(): Promise<AudioBuffer> {
    return new Promise((resolve) => {
      this.mediaRecorder!.onstop = async () => {
        const blob = new Blob(this.audioChunks, { type: 'audio/webm' });
        const arrayBuffer = await blob.arrayBuffer();
        const audioBuffer = await this.audioContext.decodeAudioData(arrayBuffer);
        resolve(audioBuffer);
      };
      this.mediaRecorder!.stop();
    });
  }
}
```

### 2. Pitch Detection Pattern

```typescript
// lib/audio/pitch-detector.ts
import Pitchfinder from 'pitchfinder';

export class PitchDetector {
  private detectPitch: (input: Float32Array) => number | null;

  constructor(sampleRate: number = 44100) {
    this.detectPitch = Pitchfinder.YIN({ sampleRate, threshold: 0.1 });
  }

  analyze(audioBuffer: AudioBuffer): PitchPoint[] {
    const data = audioBuffer.getChannelData(0);
    const pitches: PitchPoint[] = [];
    const windowSize = 2048;
    const hopSize = 512;

    for (let i = 0; i < data.length - windowSize; i += hopSize) {
      const window = data.slice(i, i + windowSize);
      const frequency = this.detectPitch(window);
      
      if (frequency && frequency > 0) {
        pitches.push({
          frequency,
          time: i / audioBuffer.sampleRate,
          midi: this.frequencyToMidi(frequency)
        });
      }
    }
    return pitches;
  }

  private frequencyToMidi(freq: number): number {
    return Math.round(69 + 12 * Math.log2(freq / 440));
  }
}
```

### 3. BPM Detection Pattern

```typescript
// lib/audio/bpm-detector.ts
import { createRealTimeBpmProcessor } from 'realtime-bpm-analyzer';

export class BPMDetector {
  async analyze(audioBuffer: AudioBuffer): Promise<BPMAnalysis> {
    const audioContext = new AudioContext();
    const processor = await createRealTimeBpmProcessor(audioContext);
    
    return new Promise((resolve) => {
      processor.port.onmessage = (event) => {
        if (event.data.message === 'BPM_STABLE') {
          resolve({
            bpm: Math.round(event.data.data.bpm),
            confidence: 1.0,
            stable: true
          });
        }
      };

      const source = audioContext.createBufferSource();
      source.buffer = audioBuffer;
      source.connect(processor);
      source.start();
    });
  }
}
```

### 4. Music Generation Pattern

```typescript
// lib/audio/music-generator.ts
import * as Tone from 'tone';

export class MusicGenerator {
  private kick: Tone.MembraneSynth;
  private snare: Tone.NoiseSynth;
  private bass: Tone.Synth;
  private chords: Tone.PolySynth;

  constructor() {
    this.kick = new Tone.MembraneSynth().toDestination();
    this.snare = new Tone.NoiseSynth().toDestination();
    this.bass = new Tone.Synth({ oscillator: { type: 'sawtooth' } }).toDestination();
    this.chords = new Tone.PolySynth(Tone.Synth).toDestination();
  }

  setBPM(bpm: number): void {
    Tone.Transport.bpm.value = bpm;
  }

  generateDrums(pattern: 'simple' | 'moderate' | 'busy'): void {
    const drumPattern = new Tone.Pattern((time, note) => {
      if (note === 'kick') this.kick.triggerAttackRelease('C1', '8n', time);
      if (note === 'snare') this.snare.triggerAttack(time);
    }, ['kick', 'hihat', 'snare', 'hihat']);
    
    drumPattern.start(0);
  }

  start(): void {
    Tone.Transport.start();
  }

  stop(): void {
    Tone.Transport.stop();
  }
}
```

### 5. Export Pattern

```typescript
// lib/audio/stem-exporter.ts
import * as Tone from 'tone';

export class StemExporter {
  async exportStem(type: 'drums' | 'bass' | 'chords'): Promise<Blob> {
    const recorder = new Tone.Recorder();
    
    // Connect instrument to recorder
    if (type === 'drums') {
      this.kick.connect(recorder);
      this.snare.connect(recorder);
    }
    
    recorder.start();
    Tone.Transport.start();
    
    await new Promise(resolve => setTimeout(resolve, duration * 1000));
    
    return await recorder.stop();
  }
}
```

---

## COMMON ISSUES & SOLUTIONS

### Issue 1: "AudioContext not allowed"

**Error:** `The AudioContext was not allowed to start`

**Solution:**
```typescript
// Must initialize Tone.js after user interaction
await Tone.start();
```

### Issue 2: Microphone permission denied

**Solution:**
```typescript
try {
  const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
} catch (error) {
  alert('Microphone permission required');
}
```

### Issue 3: Pitch detection too slow

**Solution:**
```typescript
// Increase hop size for faster analysis
const hopSize = 1024; // Instead of 512
```

### Issue 4: Export files too large

**Solution:**
```typescript
// Use MP3 instead of WAV
const recorder = new Tone.Recorder({
  audioBitsPerSecond: 128000 // 128kbps MP3
});
```

### Issue 5: BPM detection incorrect

**Solution:**
```typescript
// Use longer stabilization time
const processor = await createRealTimeBpmProcessor(audioContext, {
  stabilizationTime: 20000 // 20 seconds
});
```

---

## TESTING COMMANDS

```bash
# Run development server
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Type check
npx tsc --noEmit

# Lint
npm run lint
```

---

## DEPLOYMENT COMMANDS

### Railway

```bash
# Install CLI
npm install -g @railway/cli

# Login
railway login

# Deploy
railway up
```

### Vercel

```bash
# Install CLI
npm install -g vercel

# Deploy
vercel
```

### Netlify

```bash
# Install CLI
npm install -g netlify-cli

# Deploy
netlify deploy --prod
```

---

## CRITICAL REMINDERS

1. **Always await Tone.start()** before using Tone.js
2. **Dispose audio resources** when done
3. **Handle microphone permissions** gracefully
4. **Test on multiple browsers** (Chrome, Firefox, Safari)
5. **Use environment variables** for API keys
6. **Add loading states** for async operations
7. **Handle errors** with try-catch
8. **Optimize performance** (use requestAnimationFrame for viz)
9. **Clean up intervals/timeouts** on unmount
10. **Follow React hooks rules** (useEffect dependencies)

---

## LIBRARY IMPORT PATTERNS

```typescript
// Audio Processing
import Pitchfinder from 'pitchfinder';
import { createRealTimeBpmProcessor } from 'realtime-bpm-analyzer';
import * as Tone from 'tone';
import { Note, Key, Chord } from 'tonal';
import { Midi } from '@tonejs/midi';

// UI Components
import { Mic, Square, Play, Pause, Download } from 'lucide-react';

// React Hooks
import { useState, useEffect, useRef, useCallback } from 'react';

// Next.js
import type { Metadata } from 'next';
```

---

## FILE TEMPLATES

### Component Template

```typescript
'use client';

import { useState } from 'react';

interface Props {
  // Define props
}

export default function ComponentName({ }: Props) {
  const [state, setState] = useState();

  return (
    <div className="space-y-4">
      {/* Content */}
    </div>
  );
}
```

### Class Template

```typescript
export class ClassName {
  private property: Type;

  constructor(param: Type) {
    this.property = param;
  }

  public method(): ReturnType {
    // Implementation
  }

  dispose(): void {
    // Cleanup
  }
}
```

---

## DEBUG CHECKLIST

When something doesn't work:

1. ✅ Check console for errors
2. ✅ Verify all imports are correct
3. ✅ Check if audio context started (await Tone.start())
4. ✅ Verify microphone permission granted
5. ✅ Check if audio buffer is valid
6. ✅ Verify sample rate (should be 44100)
7. ✅ Check if Transport is started
8. ✅ Verify network tab for failed requests
9. ✅ Check React DevTools for state
10. ✅ Try different browser

---

## PERFORMANCE TIPS

1. **Use useMemo for expensive calculations**
```typescript
const analysis = useMemo(() => detector.analyze(audioBuffer), [audioBuffer]);
```

2. **Use useCallback for event handlers**
```typescript
const handleRecord = useCallback(() => {
  recorder.start();
}, [recorder]);
```

3. **Debounce frequent updates**
```typescript
const debouncedUpdate = debounce(updateValue, 300);
```

4. **Use Web Workers for heavy processing** (advanced)
```typescript
const worker = new Worker('pitch-worker.js');
worker.postMessage(audioData);
```

5. **Lazy load components**
```typescript
const ExportPanel = lazy(() => import('./ExportPanel'));
```

---

## DOCUMENTATION LINKS

- **Tone.js:** https://tonejs.github.io/
- **pitchfinder:** https://github.com/peterkhayes/pitchfinder
- **realtime-bpm:** https://www.realtime-bpm-analyzer.com/
- **Tonal.js:** https://github.com/tonaljs/tonal
- **@tonejs/midi:** https://tonejs.github.io/midi/
- **Next.js:** https://nextjs.org/docs
- **React:** https://react.dev/
- **Tailwind:** https://tailwindcss.com/docs

---

## SUCCESS CRITERIA

✅ User can record voice  
✅ Pitch detected accurately  
✅ BPM detected (±5)  
✅ Key detected correctly  
✅ Drums generate and play  
✅ Bass follows chords  
✅ Chords sound musical  
✅ Multiple sections work  
✅ Stems export successfully  
✅ MIDI exports correctly  
✅ No console errors  
✅ Smooth UX (no lag)  
✅ Works in Chrome/Firefox/Safari  

---

**END OF QUICK START GUIDE**
