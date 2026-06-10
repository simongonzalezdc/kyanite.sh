# VOXFORGE - PIXI.JS INTEGRATION GUIDE
**Version:** 1.0  
**Last Updated:** November 9, 2025  
**Purpose:** Complete implementation guide for all Pixi.js interactive features

---

## 📚 TABLE OF CONTENTS

1. [Overview](#overview)
2. [Installation & Setup](#installation--setup)
3. [Day 13: Real-Time Visualization](#day-13-real-time-visualization)
4. [Day 14: Interactive Piano Roll](#day-14-interactive-piano-roll)
5. [Day 15: Particle System](#day-15-particle-system)
6. [Day 16: Rhythm Game Mode](#day-16-rhythm-game-mode)
7. [Integration with Existing Code](#integration-with-existing-code)
8. [Testing & Debugging](#testing--debugging)
9. [Performance Optimization](#performance-optimization)

---

## OVERVIEW

### What We're Building

Adding Pixi.js to VoxForge will enable four major interactive features:

1. **Real-Time Pitch Visualization** - See notes as you sing them
2. **Interactive Piano Roll Editor** - Click/drag to edit melody
3. **Audio-Reactive Particle System** - Beautiful visual effects
4. **Rhythm Game Recording Mode** - Gamified music creation

### Why Pixi.js?

- **High Performance:** WebGL-accelerated 2D rendering
- **Lightweight:** Only 500KB (vs Phaser's 1.5MB)
- **Perfect for Music:** Ideal for audio visualizations
- **Battle-Tested:** Used by Disney, BBC, McDonald's
- **Great with Tone.js:** Common combination in music apps

### Technology Stack Addition

```json
{
  "pixi.js": "^8.0.0",
  "gsap": "^3.12.0",
  "@pixi/graphics-extras": "^8.0.0"
}
```

---

## INSTALLATION & SETUP

### Step 1: Install Dependencies

```bash
cd voxforge
npm install pixi.js gsap @pixi/graphics-extras
npm install -D @types/gsap
```

### Step 2: Create Folder Structure

```bash
mkdir -p lib/pixi
mkdir -p app/components/pixi
mkdir -p public/assets/particles
```

### Step 3: Update TypeScript Config

Add to `tsconfig.json`:

```json
{
  "compilerOptions": {
    "types": ["pixi.js", "gsap"]
  }
}
```

### Step 4: Create Pixi Types

Create `lib/types/pixi.ts`:

```typescript
import * as PIXI from 'pixi.js';
import { PitchPoint } from './index';

export interface PixiNote {
  sprite: PIXI.Graphics;
  midi: number;
  startTime: number;
  duration: number;
  velocity: number;
}

export interface ParticleOptions {
  color: number;
  size: number;
  lifetime: number;
  velocity: { x: number; y: number };
  gravity: number;
}

export interface VisualizerConfig {
  width: number;
  height: number;
  backgroundColor: number;
  noteHeight: number;
  pixelsPerSecond: number;
  minMidi: number;
  maxMidi: number;
}

export interface PianoRollConfig extends VisualizerConfig {
  gridColor: number;
  noteColor: number;
  activeNoteColor: number;
  snapToGrid: boolean;
  beatsPerBar: number;
  subdivision: number;
}

export interface RhythmGameConfig {
  fallSpeed: number;
  targetY: number;
  hitWindow: number;
  scoreMultiplier: number;
}
```

---

## DAY 13: REAL-TIME VISUALIZATION

### Goal: Display notes in real-time as user sings

### Step 1: Create Base Visualizer Class

Create `lib/pixi/visualizer.ts`:

```typescript
import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { PitchPoint } from '../types';
import { VisualizerConfig } from '../types/pixi';

export class MusicVisualizer {
  private app: PIXI.Application;
  private container: PIXI.Container;
  private noteLayer: PIXI.Container;
  private gridLayer: PIXI.Container;
  private config: VisualizerConfig;
  private currentTime: number = 0;
  private notes: Map<string, PIXI.Graphics> = new Map();

  constructor(config: Partial<VisualizerConfig> = {}) {
    this.config = {
      width: config.width || 800,
      height: config.height || 400,
      backgroundColor: config.backgroundColor || 0x0a0a0a,
      noteHeight: config.noteHeight || 6,
      pixelsPerSecond: config.pixelsPerSecond || 100,
      minMidi: config.minMidi || 36, // C2
      maxMidi: config.maxMidi || 84  // C6
    };

    // Create Pixi application
    this.app = new PIXI.Application({
      width: this.config.width,
      height: this.config.height,
      backgroundColor: this.config.backgroundColor,
      antialias: true,
      resolution: window.devicePixelRatio || 1,
      autoDensity: true
    });

    // Create layers
    this.container = new PIXI.Container();
    this.gridLayer = new PIXI.Container();
    this.noteLayer = new PIXI.Container();

    this.container.addChild(this.gridLayer);
    this.container.addChild(this.noteLayer);
    this.app.stage.addChild(this.container);

    // Draw background grid
    this.drawGrid();
  }

  /**
   * Get the Pixi canvas element to add to DOM
   */
  getCanvas(): HTMLCanvasElement {
    return this.app.view as HTMLCanvasElement;
  }

  /**
   * Draw piano roll grid lines
   */
  private drawGrid(): void {
    const graphics = new PIXI.Graphics();
    
    // Horizontal lines (one per octave)
    const midiRange = this.config.maxMidi - this.config.minMidi;
    const octaves = Math.floor(midiRange / 12);
    
    for (let i = 0; i <= octaves; i++) {
      const y = this.midiToY(this.config.minMidi + (i * 12));
      
      graphics.lineStyle(1, 0x333333, 0.5);
      graphics.moveTo(0, y);
      graphics.lineTo(this.config.width, y);
    }

    // Vertical lines (beat markers)
    const beatsToShow = Math.ceil(this.config.width / this.config.pixelsPerSecond);
    for (let i = 0; i <= beatsToShow; i++) {
      const x = i * this.config.pixelsPerSecond;
      
      graphics.lineStyle(1, 0x222222, 0.3);
      graphics.moveTo(x, 0);
      graphics.lineTo(x, this.config.height);
    }

    this.gridLayer.addChild(graphics);
  }

  /**
   * Convert MIDI note to Y position
   */
  private midiToY(midi: number): number {
    const midiRange = this.config.maxMidi - this.config.minMidi;
    const normalizedMidi = (midi - this.config.minMidi) / midiRange;
    return this.config.height - (normalizedMidi * this.config.height);
  }

  /**
   * Convert time to X position
   */
  private timeToX(time: number): number {
    return (time - this.currentTime) * this.config.pixelsPerSecond;
  }

  /**
   * Get color for MIDI note (rainbow across octaves)
   */
  private getColorForMidi(midi: number): number {
    const hue = ((midi % 12) / 12) * 360;
    return PIXI.utils.string2hex(`hsl(${hue}, 70%, 60%)`);
  }

  /**
   * Add a note to the visualization
   */
  addNote(pitch: PitchPoint, isRealtime: boolean = true): void {
    const noteGraphics = new PIXI.Graphics();
    const color = this.getColorForMidi(pitch.midi);
    const y = this.midiToY(pitch.midi);
    const x = this.timeToX(pitch.time);

    // Draw note rectangle
    noteGraphics.beginFill(color, 0.8);
    noteGraphics.drawRoundedRect(
      x,
      y - this.config.noteHeight / 2,
      5, // width (will grow as note sustains)
      this.config.noteHeight,
      2 // corner radius
    );
    noteGraphics.endFill();

    // Add glow effect
    noteGraphics.filters = [new PIXI.filters.BlurFilter(2, 4)];

    this.noteLayer.addChild(noteGraphics);

    // Store reference
    const noteId = `${pitch.time}-${pitch.midi}`;
    this.notes.set(noteId, noteGraphics);

    if (isRealtime) {
      // Fade in animation
      noteGraphics.alpha = 0;
      gsap.to(noteGraphics, {
        alpha: 1,
        duration: 0.1
      });
    }
  }

  /**
   * Update note duration (as user continues singing same note)
   */
  updateNoteDuration(noteId: string, endTime: number): void {
    const note = this.notes.get(noteId);
    if (!note) return;

    const duration = endTime - this.currentTime;
    const width = duration * this.config.pixelsPerSecond;
    
    // Redraw with new width
    note.clear();
    const color = note.tint || 0x3B82F6;
    note.beginFill(color, 0.8);
    note.drawRoundedRect(0, 0, width, this.config.noteHeight, 2);
    note.endFill();
  }

  /**
   * Scroll the view (for real-time recording)
   */
  scroll(deltaTime: number): void {
    this.currentTime += deltaTime;
    this.container.x -= deltaTime * this.config.pixelsPerSecond;

    // Remove notes that are off-screen
    const minX = -this.container.x - 100;
    this.notes.forEach((note, id) => {
      if (note.x < minX) {
        note.destroy();
        this.notes.delete(id);
      }
    });
  }

  /**
   * Clear all notes
   */
  clear(): void {
    this.notes.forEach(note => note.destroy());
    this.notes.clear();
    this.noteLayer.removeChildren();
    this.currentTime = 0;
    this.container.x = 0;
  }

  /**
   * Load and display pre-recorded notes
   */
  loadNotes(pitches: PitchPoint[]): void {
    this.clear();
    pitches.forEach(pitch => this.addNote(pitch, false));
  }

  /**
   * Animate playback cursor
   */
  playback(startTime: number, duration: number): void {
    const cursor = new PIXI.Graphics();
    cursor.lineStyle(2, 0x3B82F6, 1);
    cursor.moveTo(0, 0);
    cursor.lineTo(0, this.config.height);
    
    this.container.addChild(cursor);

    gsap.to(cursor, {
      x: duration * this.config.pixelsPerSecond,
      duration: duration,
      ease: 'none',
      onComplete: () => cursor.destroy()
    });
  }

  /**
   * Resize canvas
   */
  resize(width: number, height: number): void {
    this.config.width = width;
    this.config.height = height;
    this.app.renderer.resize(width, height);
    
    // Redraw grid
    this.gridLayer.removeChildren();
    this.drawGrid();
  }

  /**
   * Cleanup
   */
  destroy(): void {
    this.clear();
    this.app.destroy(true, { children: true, texture: true });
  }
}
```

### Step 2: Create React Component

Create `app/components/pixi/VisualizerCanvas.tsx`:

```typescript
'use client';

import { useEffect, useRef } from 'react';
import { MusicVisualizer } from '@/lib/pixi/visualizer';
import { PitchPoint } from '@/lib/types';

interface VisualizerCanvasProps {
  pitches?: PitchPoint[];
  isRecording?: boolean;
  onVisualizerReady?: (visualizer: MusicVisualizer) => void;
}

export default function VisualizerCanvas({ 
  pitches, 
  isRecording,
  onVisualizerReady 
}: VisualizerCanvasProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const visualizerRef = useRef<MusicVisualizer | null>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create visualizer
    const visualizer = new MusicVisualizer({
      width: containerRef.current.clientWidth,
      height: 400
    });

    // Add canvas to DOM
    containerRef.current.appendChild(visualizer.getCanvas());
    visualizerRef.current = visualizer;

    // Notify parent
    onVisualizerReady?.(visualizer);

    // Handle resize
    const handleResize = () => {
      if (containerRef.current) {
        visualizer.resize(
          containerRef.current.clientWidth,
          400
        );
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      visualizer.destroy();
    };
  }, []);

  // Load pitches when they change
  useEffect(() => {
    if (visualizerRef.current && pitches && !isRecording) {
      visualizerRef.current.loadNotes(pitches);
    }
  }, [pitches, isRecording]);

  // Auto-scroll during recording
  useEffect(() => {
    if (!visualizerRef.current || !isRecording) return;

    let animationFrame: number;
    let lastTime = performance.now();

    const scroll = (currentTime: number) => {
      const deltaTime = (currentTime - lastTime) / 1000; // Convert to seconds
      lastTime = currentTime;

      visualizerRef.current?.scroll(deltaTime);
      animationFrame = requestAnimationFrame(scroll);
    };

    animationFrame = requestAnimationFrame(scroll);

    return () => {
      cancelAnimationFrame(animationFrame);
    };
  }, [isRecording]);

  return (
    <div 
      ref={containerRef} 
      className="w-full rounded-lg overflow-hidden border border-gray-800"
    />
  );
}
```

### Step 3: Integrate with Recording

Update `app/components/Recorder.tsx` to include visualizer:

```typescript
'use client';

import { useState, useRef } from 'react';
import { AudioRecorder } from '@/lib/audio/recorder';
import { MusicVisualizer } from '@/lib/pixi/visualizer';
import VisualizerCanvas from './pixi/VisualizerCanvas';
import { Mic, Square } from 'lucide-react';

export default function Recorder({ onRecordingComplete }) {
  const [isRecording, setIsRecording] = useState(false);
  const recorderRef = useRef<AudioRecorder | null>(null);
  const visualizerRef = useRef<MusicVisualizer | null>(null);
  const pitchDetectionInterval = useRef<NodeJS.Timeout | null>(null);

  const handleVisualizerReady = (visualizer: MusicVisualizer) => {
    visualizerRef.current = visualizer;
  };

  const startRecording = async () => {
    if (!recorderRef.current) {
      recorderRef.current = new AudioRecorder();
    }

    await recorderRef.current.startRecording();
    setIsRecording(true);

    // Real-time pitch visualization
    pitchDetectionInterval.current = setInterval(() => {
      if (!visualizerRef.current || !recorderRef.current) return;

      // Get current audio data
      const waveform = recorderRef.current.getWaveformData();
      
      // Simple pitch detection (you'd use your actual pitch detector here)
      // This is just for real-time visualization
      const pitch = detectPitchFromWaveform(waveform);
      
      if (pitch) {
        visualizerRef.current.addNote({
          frequency: pitch,
          time: performance.now() / 1000,
          midi: frequencyToMidi(pitch),
          confidence: 1
        }, true);
      }
    }, 50); // Update every 50ms
  };

  const stopRecording = async () => {
    if (pitchDetectionInterval.current) {
      clearInterval(pitchDetectionInterval.current);
    }

    const audioBuffer = await recorderRef.current!.stopRecording();
    setIsRecording(false);
    onRecordingComplete(audioBuffer);
  };

  return (
    <div className="space-y-4">
      <VisualizerCanvas 
        isRecording={isRecording}
        onVisualizerReady={handleVisualizerReady}
      />

      <div className="flex gap-4 justify-center">
        {!isRecording ? (
          <button
            onClick={startRecording}
            className="flex items-center gap-2 px-6 py-3 bg-primary-500 hover:bg-primary-600 rounded-lg"
          >
            <Mic size={20} />
            Start Recording
          </button>
        ) : (
          <button
            onClick={stopRecording}
            className="flex items-center gap-2 px-6 py-3 bg-red-500 hover:bg-red-600 rounded-lg"
          >
            <Square size={20} />
            Stop Recording
          </button>
        )}
      </div>
    </div>
  );
}

// Helper functions
function frequencyToMidi(frequency: number): number {
  return Math.round(69 + 12 * Math.log2(frequency / 440));
}

function detectPitchFromWaveform(waveform: Uint8Array): number | null {
  // Simplified pitch detection
  // In production, use your actual PitchDetector
  // This is just for demo
  return null;
}
```

### Day 13 Testing Checklist

- [ ] Pixi.js canvas appears in UI
- [ ] Grid lines display correctly
- [ ] Notes appear as you sing
- [ ] Notes are color-coded by pitch
- [ ] Canvas scrolls during recording
- [ ] Notes fade in smoothly
- [ ] Canvas resizes with window
- [ ] No performance issues

---

## DAY 14: INTERACTIVE PIANO ROLL

### Goal: Click/drag to edit melody visually

### Step 1: Create Piano Roll Class

Create `lib/pixi/piano-roll.ts`:

```typescript
import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { PitchPoint } from '../types';
import { PianoRollConfig, PixiNote } from '../types/pixi';
import { Midi } from '@tonejs/midi';

export class InteractivePianoRoll {
  private app: PIXI.Application;
  private container: PIXI.Container;
  private gridLayer: PIXI.Container;
  private noteLayer: PIXI.Container;
  private pianoLayer: PIXI.Container;
  private config: PianoRollConfig;
  private notes: PixiNote[] = [];
  private selectedNote: PixiNote | null = null;
  private isDragging: boolean = false;
  private dragStartPos: { x: number; y: number } = { x: 0, y: 0 };

  constructor(config: Partial<PianoRollConfig> = {}) {
    this.config = {
      width: config.width || 1000,
      height: config.height || 600,
      backgroundColor: config.backgroundColor || 0x0a0a0a,
      noteHeight: config.noteHeight || 12,
      pixelsPerSecond: config.pixelsPerSecond || 100,
      minMidi: config.minMidi || 36,
      maxMidi: config.maxMidi || 84,
      gridColor: config.gridColor || 0x333333,
      noteColor: config.noteColor || 0x3B82F6,
      activeNoteColor: config.activeNoteColor || 0x8B5CF6,
      snapToGrid: config.snapToGrid ?? true,
      beatsPerBar: config.beatsPerBar || 4,
      subdivision: config.subdivision || 16 // 16th notes
    };

    // Create Pixi application
    this.app = new PIXI.Application({
      width: this.config.width,
      height: this.config.height,
      backgroundColor: this.config.backgroundColor,
      antialias: true
    });

    // Create layers
    this.container = new PIXI.Container();
    this.gridLayer = new PIXI.Container();
    this.pianoLayer = new PIXI.Container();
    this.noteLayer = new PIXI.Container();

    this.container.addChild(this.gridLayer);
    this.container.addChild(this.pianoLayer);
    this.container.addChild(this.noteLayer);
    this.app.stage.addChild(this.container);

    // Make noteLayer interactive
    this.noteLayer.interactive = true;
    this.noteLayer.hitArea = new PIXI.Rectangle(
      0, 0, 
      this.config.width, 
      this.config.height
    );

    // Setup interaction
    this.setupInteraction();

    // Draw grid and piano keys
    this.drawGrid();
    this.drawPianoKeys();
  }

  getCanvas(): HTMLCanvasElement {
    return this.app.view as HTMLCanvasElement;
  }

  /**
   * Draw piano roll grid
   */
  private drawGrid(): void {
    const graphics = new PIXI.Graphics();
    const midiRange = this.config.maxMidi - this.config.minMidi;

    // Horizontal lines (one per MIDI note)
    for (let i = 0; i <= midiRange; i++) {
      const midi = this.config.minMidi + i;
      const y = this.midiToY(midi);
      const isBlackKey = this.isBlackKey(midi);
      
      // Alternate colors for white/black keys
      graphics.lineStyle(
        1, 
        this.config.gridColor, 
        isBlackKey ? 0.3 : 0.2
      );
      graphics.moveTo(60, y);
      graphics.lineTo(this.config.width, y);
    }

    // Vertical lines (beat markers)
    const totalBeats = Math.ceil(this.config.width / this.config.pixelsPerSecond);
    for (let i = 0; i <= totalBeats; i++) {
      const x = 60 + (i * this.config.pixelsPerSecond);
      const isBarLine = i % this.config.beatsPerBar === 0;
      
      graphics.lineStyle(
        isBarLine ? 2 : 1,
        this.config.gridColor,
        isBarLine ? 0.5 : 0.2
      );
      graphics.moveTo(x, 0);
      graphics.lineTo(x, this.config.height);
    }

    // Subdivision lines
    const subdivisionWidth = this.config.pixelsPerSecond / this.config.subdivision;
    for (let i = 0; i <= (this.config.width / subdivisionWidth); i++) {
      if (i % this.config.subdivision === 0) continue; // Skip beat lines
      
      const x = 60 + (i * subdivisionWidth);
      graphics.lineStyle(1, this.config.gridColor, 0.1);
      graphics.moveTo(x, 0);
      graphics.lineTo(x, this.config.height);
    }

    this.gridLayer.addChild(graphics);
  }

  /**
   * Draw piano keys on the left
   */
  private drawPianoKeys(): void {
    const graphics = new PIXI.Graphics();
    const midiRange = this.config.maxMidi - this.config.minMidi;

    for (let i = 0; i <= midiRange; i++) {
      const midi = this.config.minMidi + i;
      const y = this.midiToY(midi);
      const isBlackKey = this.isBlackKey(midi);
      
      // Draw key
      graphics.beginFill(isBlackKey ? 0x1a1a1a : 0x2a2a2a);
      graphics.drawRect(
        0, 
        y - this.config.noteHeight / 2,
        60,
        this.config.noteHeight
      );
      graphics.endFill();

      // Draw border
      graphics.lineStyle(1, 0x444444, 0.5);
      graphics.drawRect(
        0,
        y - this.config.noteHeight / 2,
        60,
        this.config.noteHeight
      );

      // Add note label for C notes
      if (midi % 12 === 0) {
        const text = new PIXI.Text(
          this.midiToNoteName(midi),
          {
            fontSize: 10,
            fill: 0x888888
          }
        );
        text.x = 5;
        text.y = y - 5;
        this.pianoLayer.addChild(text);
      }
    }

    this.pianoLayer.addChild(graphics);
  }

  /**
   * Setup mouse/touch interaction
   */
  private setupInteraction(): void {
    this.noteLayer.on('pointerdown', (e: PIXI.FederatedPointerEvent) => {
      const pos = e.getLocalPosition(this.noteLayer);
      
      // Check if clicking on existing note
      const clickedNote = this.getNoteAt(pos.x, pos.y);
      
      if (clickedNote) {
        // Select note for dragging
        this.selectNote(clickedNote);
        this.isDragging = true;
        this.dragStartPos = { x: pos.x, y: pos.y };
      } else {
        // Create new note
        this.createNoteAt(pos.x, pos.y);
      }
    });

    this.noteLayer.on('pointermove', (e: PIXI.FederatedPointerEvent) => {
      if (!this.isDragging || !this.selectedNote) return;

      const pos = e.getLocalPosition(this.noteLayer);
      const deltaX = pos.x - this.dragStartPos.x;
      const deltaY = pos.y - this.dragStartPos.y;

      // Move note
      this.moveNote(this.selectedNote, deltaX, deltaY);
      this.dragStartPos = { x: pos.x, y: pos.y };
    });

    this.noteLayer.on('pointerup', () => {
      this.isDragging = false;
      if (this.selectedNote) {
        this.deselectNote(this.selectedNote);
        this.selectedNote = null;
      }
    });

    this.noteLayer.on('pointerupoutside', () => {
      this.isDragging = false;
      if (this.selectedNote) {
        this.deselectNote(this.selectedNote);
        this.selectedNote = null;
      }
    });

    // Right-click to delete
    this.noteLayer.on('rightdown', (e: PIXI.FederatedPointerEvent) => {
      e.preventDefault();
      const pos = e.getLocalPosition(this.noteLayer);
      const note = this.getNoteAt(pos.x, pos.y);
      if (note) {
        this.deleteNote(note);
      }
    });
  }

  /**
   * Create a new note at position
   */
  private createNoteAt(x: number, y: number): void {
    // Convert position to time and MIDI
    let time = this.xToTime(x);
    let midi = Math.round(this.yToMidi(y));

    // Snap to grid if enabled
    if (this.config.snapToGrid) {
      const subdivision = 1 / (this.config.subdivision / 4); // Quarter note subdivisions
      time = Math.round(time / subdivision) * subdivision;
    }

    // Default duration (1 beat)
    const duration = 1;

    // Create note sprite
    const noteSprite = this.createNoteSprite(time, midi, duration, 0.8);
    
    const pixiNote: PixiNote = {
      sprite: noteSprite,
      midi,
      startTime: time,
      duration,
      velocity: 0.8
    };

    this.notes.push(pixiNote);
    this.noteLayer.addChild(noteSprite);

    // Animate appearance
    noteSprite.alpha = 0;
    gsap.to(noteSprite, {
      alpha: 1,
      duration: 0.2
    });
  }

  /**
   * Create note sprite
   */
  private createNoteSprite(
    time: number,
    midi: number,
    duration: number,
    velocity: number
  ): PIXI.Graphics {
    const sprite = new PIXI.Graphics();
    const x = this.timeToX(time);
    const y = this.midiToY(midi);
    const width = duration * this.config.pixelsPerSecond;
    const color = this.getColorForMidi(midi);

    sprite.beginFill(color, 0.8);
    sprite.drawRoundedRect(
      x,
      y - this.config.noteHeight / 2,
      width,
      this.config.noteHeight,
      3
    );
    sprite.endFill();

    // Border
    sprite.lineStyle(1, 0xffffff, 0.3);
    sprite.drawRoundedRect(
      x,
      y - this.config.noteHeight / 2,
      width,
      this.config.noteHeight,
      3
    );

    // Make interactive
    sprite.interactive = true;
    sprite.buttonMode = true;
    sprite.cursor = 'pointer';

    return sprite;
  }

  /**
   * Get note at position
   */
  private getNoteAt(x: number, y: number): PixiNote | null {
    const time = this.xToTime(x);
    const midi = Math.round(this.yToMidi(y));

    return this.notes.find(note => {
      const inTimeRange = time >= note.startTime && 
                         time <= note.startTime + note.duration;
      const inMidiRange = Math.abs(note.midi - midi) < 1;
      return inTimeRange && inMidiRange;
    }) || null;
  }

  /**
   * Select note for editing
   */
  private selectNote(note: PixiNote): void {
    this.selectedNote = note;
    note.sprite.tint = this.config.activeNoteColor;
    
    gsap.to(note.sprite, {
      alpha: 1,
      duration: 0.1
    });
  }

  /**
   * Deselect note
   */
  private deselectNote(note: PixiNote): void {
    note.sprite.tint = this.getColorForMidi(note.midi);
    
    gsap.to(note.sprite, {
      alpha: 0.8,
      duration: 0.1
    });
  }

  /**
   * Move note by delta
   */
  private moveNote(note: PixiNote, deltaX: number, deltaY: number): void {
    // Update time
    let newTime = note.startTime + (deltaX / this.config.pixelsPerSecond);
    
    // Snap to grid
    if (this.config.snapToGrid) {
      const subdivision = 1 / (this.config.subdivision / 4);
      newTime = Math.round(newTime / subdivision) * subdivision;
    }

    // Update MIDI
    const deltaM MIDI = Math.round(-(deltaY / this.config.noteHeight));
    let newMidi = note.midi + deltaMidi;
    newMidi = Math.max(this.config.minMidi, Math.min(this.config.maxMidi, newMidi));

    // Update note data
    note.startTime = Math.max(0, newTime);
    note.midi = newMidi;

    // Update sprite position
    note.sprite.x = this.timeToX(note.startTime);
    note.sprite.y = this.midiToY(note.midi) - this.config.noteHeight / 2;
    note.sprite.tint = this.getColorForMidi(note.midi);
  }

  /**
   * Delete note
   */
  private deleteNote(note: PixiNote): void {
    // Animate removal
    gsap.to(note.sprite, {
      alpha: 0,
      scale: 0.5,
      duration: 0.2,
      onComplete: () => {
        this.noteLayer.removeChild(note.sprite);
        note.sprite.destroy();
      }
    });

    // Remove from array
    const index = this.notes.indexOf(note);
    if (index > -1) {
      this.notes.splice(index, 1);
    }
  }

  /**
   * Load notes from pitch data
   */
  loadNotes(pitches: PitchPoint[]): void {
    this.clear();
    
    // Convert pitches to notes
    let currentNote: PixiNote | null = null;
    
    pitches.forEach((pitch, i) => {
      if (!currentNote || pitch.midi !== currentNote.midi) {
        // Start new note
        if (currentNote) {
          currentNote.duration = pitch.time - currentNote.startTime;
        }
        
        const noteSprite = this.createNoteSprite(
          pitch.time,
          pitch.midi,
          0.25, // Default duration
          0.8
        );
        
        currentNote = {
          sprite: noteSprite,
          midi: pitch.midi,
          startTime: pitch.time,
          duration: 0.25,
          velocity: 0.8
        };
        
        this.notes.push(currentNote);
        this.noteLayer.addChild(noteSprite);
      } else {
        // Extend current note
        currentNote.duration = pitch.time - currentNote.startTime;
        
        // Update sprite width
        const width = currentNote.duration * this.config.pixelsPerSecond;
        currentNote.sprite.clear();
        currentNote.sprite.beginFill(this.getColorForMidi(currentNote.midi), 0.8);
        currentNote.sprite.drawRoundedRect(
          this.timeToX(currentNote.startTime),
          this.midiToY(currentNote.midi) - this.config.noteHeight / 2,
          width,
          this.config.noteHeight,
          3
        );
        currentNote.sprite.endFill();
      }
    });
  }

  /**
   * Export to MIDI
   */
  exportToMIDI(bpm: number): Midi {
    const midi = new Midi();
    midi.header.setTempo(bpm);
    
    const track = midi.addTrack();
    track.name = 'VoxForge Melody';
    
    this.notes.forEach(note => {
      track.addNote({
        midi: note.midi,
        time: note.startTime,
        duration: note.duration,
        velocity: note.velocity
      });
    });
    
    return midi;
  }

  /**
   * Export to PitchPoint array
   */
  exportToPitches(): PitchPoint[] {
    const pitches: PitchPoint[] = [];
    
    this.notes.forEach(note => {
      // Create pitch points for start and end of note
      pitches.push({
        frequency: this.midiToFrequency(note.midi),
        time: note.startTime,
        midi: note.midi,
        confidence: 1
      });
    });
    
    return pitches.sort((a, b) => a.time - b.time);
  }

  /**
   * Clear all notes
   */
  clear(): void {
    this.notes.forEach(note => {
      note.sprite.destroy();
    });
    this.notes = [];
    this.noteLayer.removeChildren();
  }

  // Helper methods
  private midiToY(midi: number): number {
    const midiRange = this.config.maxMidi - this.config.minMidi;
    const normalizedMidi = (midi - this.config.minMidi) / midiRange;
    return this.config.height - (normalizedMidi * this.config.height);
  }

  private yToMidi(y: number): number {
    const normalizedY = 1 - (y / this.config.height);
    const midiRange = this.config.maxMidi - this.config.minMidi;
    return this.config.minMidi + (normalizedY * midiRange);
  }

  private timeToX(time: number): number {
    return 60 + (time * this.config.pixelsPerSecond);
  }

  private xToTime(x: number): number {
    return (x - 60) / this.config.pixelsPerSecond;
  }

  private getColorForMidi(midi: number): number {
    const hue = ((midi % 12) / 12) * 360;
    return PIXI.utils.string2hex(`hsl(${hue}, 70%, 60%)`);
  }

  private isBlackKey(midi: number): boolean {
    const pitchClass = midi % 12;
    return [1, 3, 6, 8, 10].includes(pitchClass); // C#, D#, F#, G#, A#
  }

  private midiToNoteName(midi: number): string {
    const notes = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'];
    const octave = Math.floor(midi / 12) - 1;
    return `${notes[midi % 12]}${octave}`;
  }

  private midiToFrequency(midi: number): number {
    return 440 * Math.pow(2, (midi - 69) / 12);
  }

  resize(width: number, height: number): void {
    this.config.width = width;
    this.config.height = height;
    this.app.renderer.resize(width, height);
    
    // Redraw
    this.gridLayer.removeChildren();
    this.pianoLayer.removeChildren();
    this.drawGrid();
    this.drawPianoKeys();
  }

  destroy(): void {
    this.clear();
    this.app.destroy(true, { children: true, texture: true });
  }
}
```

### Step 2: Create Piano Roll Component

Create `app/components/pixi/PianoRollEditor.tsx`:

```typescript
'use client';

import { useEffect, useRef, useState } from 'react';
import { InteractivePianoRoll } from '@/lib/pixi/piano-roll';
import { PitchPoint } from '@/lib/types';
import { Download, Trash2, Grid3x3 } from 'lucide-react';

interface PianoRollEditorProps {
  initialPitches?: PitchPoint[];
  bpm: number;
  onExport?: (pitches: PitchPoint[]) => void;
}

export default function PianoRollEditor({ 
  initialPitches, 
  bpm,
  onExport 
}: PianoRollEditorProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const pianoRollRef = useRef<InteractivePianoRoll | null>(null);
  const [snapToGrid, setSnapToGrid] = useState(true);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create piano roll
    const pianoRoll = new InteractivePianoRoll({
      width: containerRef.current.clientWidth,
      height: 600,
      snapToGrid
    });

    containerRef.current.appendChild(pianoRoll.getCanvas());
    pianoRollRef.current = pianoRoll;

    // Load initial notes
    if (initialPitches) {
      pianoRoll.loadNotes(initialPitches);
    }

    // Handle resize
    const handleResize = () => {
      if (containerRef.current) {
        pianoRoll.resize(containerRef.current.clientWidth, 600);
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      pianoRoll.destroy();
    };
  }, [snapToGrid]);

  const handleExport = () => {
    if (!pianoRollRef.current) return;
    const pitches = pianoRollRef.current.exportToPitches();
    onExport?.(pitches);
  };

  const handleExportMIDI = () => {
    if (!pianoRollRef.current) return;
    const midi = pianoRollRef.current.exportToMIDI(bpm);
    const midiArray = midi.toArray();
    const blob = new Blob([midiArray], { type: 'audio/midi' });
    const url = URL.createObjectURL(blob);
    
    const a = document.createElement('a');
    a.href = url;
    a.download = `voxforge-edited-${Date.now()}.mid`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleClear = () => {
    if (!pianoRollRef.current) return;
    if (confirm('Clear all notes?')) {
      pianoRollRef.current.clear();
    }
  };

  const toggleSnapToGrid = () => {
    setSnapToGrid(!snapToGrid);
  };

  return (
    <div className="space-y-4">
      {/* Toolbar */}
      <div className="flex items-center justify-between bg-gray-900 p-4 rounded-lg border border-gray-800">
        <div className="flex items-center gap-4">
          <h3 className="text-lg font-semibold">Piano Roll Editor</h3>
          <button
            onClick={toggleSnapToGrid}
            className={`flex items-center gap-2 px-3 py-2 rounded-lg transition-colors ${
              snapToGrid 
                ? 'bg-primary-500 text-white' 
                : 'bg-gray-800 text-gray-400'
            }`}
          >
            <Grid3x3 size={16} />
            Snap to Grid
          </button>
        </div>

        <div className="flex gap-2">
          <button
            onClick={handleClear}
            className="flex items-center gap-2 px-4 py-2 bg-red-500/20 hover:bg-red-500/30 text-red-400 rounded-lg transition-colors"
          >
            <Trash2 size={16} />
            Clear
          </button>
          
          <button
            onClick={handleExport}
            className="flex items-center gap-2 px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg transition-colors"
          >
            <Download size={16} />
            Apply Changes
          </button>
          
          <button
            onClick={handleExportMIDI}
            className="flex items-center gap-2 px-4 py-2 bg-secondary-500 hover:bg-secondary-600 rounded-lg transition-colors"
          >
            <Download size={16} />
            Export MIDI
          </button>
        </div>
      </div>

      {/* Instructions */}
      <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-800">
        <p className="text-sm text-gray-400">
          <strong>Controls:</strong> Click to add notes • Drag to move • Right-click to delete • Scroll to zoom
        </p>
      </div>

      {/* Piano Roll Canvas */}
      <div 
        ref={containerRef}
        className="rounded-lg overflow-hidden border border-gray-800"
      />
    </div>
  );
}
```

### Day 14 Testing Checklist

- [ ] Piano roll displays with grid
- [ ] Piano keys show on left side
- [ ] Click to add note
- [ ] Drag to move note
- [ ] Right-click to delete note
- [ ] Snap to grid works
- [ ] Notes stay in bounds
- [ ] Export to MIDI works
- [ ] Apply changes updates melody

---

## DAY 15: PARTICLE SYSTEM

### Goal: Beautiful audio-reactive particle effects

### Step 1: Create Particle System

Create `lib/pixi/particle-system.ts`:

```typescript
import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { ParticleOptions } from '../types/pixi';

interface Particle {
  sprite: PIXI.Graphics;
  velocity: { x: number; y: number };
  life: number;
  maxLife: number;
}

export class MusicParticleSystem {
  private container: PIXI.Container;
  private particles: Particle[] = [];
  private particlePool: PIXI.Graphics[] = [];
  private maxParticles: number = 1000;

  constructor(container: PIXI.Container, maxParticles: number = 1000) {
    this.container = container;
    this.maxParticles = maxParticles;
    
    // Pre-create particle pool for performance
    this.initializePool();
  }

  /**
   * Pre-create particles for object pooling
   */
  private initializePool(): void {
    for (let i = 0; i < this.maxParticles; i++) {
      const particle = new PIXI.Graphics();
      particle.visible = false;
      this.container.addChild(particle);
      this.particlePool.push(particle);
    }
  }

  /**
   * Get particle from pool
   */
  private getParticle(): PIXI.Graphics | null {
    return this.particlePool.find(p => !p.visible) || null;
  }

  /**
   * Create particle burst (for note events)
   */
  createBurst(
    x: number,
    y: number,
    options: Partial<ParticleOptions> = {}
  ): void {
    const opts = {
      color: options.color || 0x3B82F6,
      size: options.size || 3,
      lifetime: options.lifetime || 1,
      velocity: options.velocity || { x: 0, y: -50 },
      gravity: options.gravity || 100
    };

    // Calculate particle count based on intensity
    const count = Math.floor(10 + Math.random() * 20);

    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      // Random angle
      const angle = Math.random() * Math.PI * 2;
      const speed = 50 + Math.random() * 100;
      
      const velocity = {
        x: Math.cos(angle) * speed,
        y: Math.sin(angle) * speed
      };

      // Draw particle
      particle.clear();
      particle.beginFill(opts.color);
      particle.drawCircle(0, 0, opts.size);
      particle.endFill();
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      // Add to active particles
      this.particles.push({
        sprite: particle,
        velocity,
        life: opts.lifetime,
        maxLife: opts.lifetime
      });
    }
  }

  /**
   * Create continuous stream (for sustained notes)
   */
  createStream(
    x: number,
    y: number,
    color: number,
    intensity: number = 1
  ): void {
    const count = Math.floor(intensity * 5);
    
    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      // Upward velocity with slight randomness
      const velocity = {
        x: (Math.random() - 0.5) * 30,
        y: -50 - Math.random() * 50
      };

      // Draw particle
      particle.clear();
      particle.beginFill(color, 0.8);
      particle.drawCircle(0, 0, 2 + Math.random() * 2);
      particle.endFill();
      
      particle.x = x + (Math.random() - 0.5) * 20;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity,
        life: 0.5 + Math.random() * 0.5,
        maxLife: 1
      });
    }
  }

  /**
   * Create explosion (for beat hits)
   */
  createExplosion(
    x: number,
    y: number,
    color: number,
    size: number = 50
  ): void {
    const count = Math.floor(30 + size);

    for (let i = 0; i < count; i++) {
      const particle = this.getParticle();
      if (!particle) continue;

      const angle = (i / count) * Math.PI * 2;
      const speed = size * 2 + Math.random() * size;
      
      const velocity = {
        x: Math.cos(angle) * speed,
        y: Math.sin(angle) * speed
      };

      // Draw particle with glow
      particle.clear();
      particle.beginFill(color, 0.9);
      particle.drawCircle(0, 0, 2 + Math.random() * 3);
      particle.endFill();
      
      // Add glow filter
      particle.filters = [new PIXI.filters.BlurFilter(2)];
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity,
        life: 0.8 + Math.random() * 0.4,
        maxLife: 1.2
      });
    }
  }

  /**
   * Create wave pattern (for harmonics)
   */
  createWave(
    startX: number,
    startY: number,
    endX: number,
    endY: number,
    color: number,
    amplitude: number = 50
  ): void {
    const distance = Math.sqrt(
      Math.pow(endX - startX, 2) + 
      Math.pow(endY - startY, 2)
    );
    const steps = Math.floor(distance / 10);

    for (let i = 0; i < steps; i++) {
      const t = i / steps;
      const particle = this.getParticle();
      if (!particle) continue;

      // Sine wave offset
      const waveOffset = Math.sin(t * Math.PI * 4) * amplitude;
      
      const x = startX + (endX - startX) * t;
      const y = startY + (endY - startY) * t + waveOffset;

      // Draw particle
      particle.clear();
      particle.beginFill(color, 0.6);
      particle.drawCircle(0, 0, 2);
      particle.endFill();
      
      particle.x = x;
      particle.y = y;
      particle.alpha = 1;
      particle.visible = true;

      this.particles.push({
        sprite: particle,
        velocity: { x: 0, y: 20 },
        life: 1 + t * 0.5,
        maxLife: 1.5
      });
    }
  }

  /**
   * Update particles (call in animation loop)
   */
  update(deltaTime: number): void {
    for (let i = this.particles.length - 1; i >= 0; i--) {
      const particle = this.particles[i];
      
      // Update life
      particle.life -= deltaTime;
      
      if (particle.life <= 0) {
        // Return to pool
        particle.sprite.visible = false;
        particle.sprite.filters = [];
        this.particles.splice(i, 1);
        continue;
      }

      // Update position
      particle.sprite.x += particle.velocity.x * deltaTime;
      particle.sprite.y += particle.velocity.y * deltaTime;

      // Apply gravity
      particle.velocity.y += 100 * deltaTime;

      // Fade out based on remaining life
      particle.sprite.alpha = particle.life / particle.maxLife;
    }
  }

  /**
   * Clear all particles
   */
  clear(): void {
    this.particles.forEach(p => {
      p.sprite.visible = false;
      p.sprite.filters = [];
    });
    this.particles = [];
  }

  /**
   * Cleanup
   */
  destroy(): void {
    this.clear();
    this.particlePool.forEach(p => p.destroy());
    this.particlePool = [];
  }
}
```

### Step 2: Integrate Particles with Visualizer

Update `lib/pixi/visualizer.ts` to add particles:

```typescript
import { MusicParticleSystem } from './particle-system';

export class MusicVisualizer {
  // ... existing code ...
  private particleSystem: MusicParticleSystem;
  private particleLayer: PIXI.Container;

  constructor(config: Partial<VisualizerConfig> = {}) {
    // ... existing setup ...
    
    // Add particle layer
    this.particleLayer = new PIXI.Container();
    this.container.addChild(this.particleLayer);
    
    // Initialize particle system
    this.particleSystem = new MusicParticleSystem(this.particleLayer, 2000);
    
    // Start update loop
    this.app.ticker.add((delta) => {
      this.particleSystem.update(delta / 60); // Convert to seconds
    });
  }

  /**
   * Add note with particle effect
   */
  addNote(pitch: PitchPoint, isRealtime: boolean = true): void {
    // ... existing note creation ...
    
    // Create particle burst at note position
    const x = this.timeToX(pitch.time);
    const y = this.midiToY(pitch.midi);
    const color = this.getColorForMidi(pitch.midi);
    
    this.particleSystem.createBurst(x, y, {
      color,
      size: 3,
      lifetime: 1
    });
  }

  /**
   * Create beat explosion
   */
  addBeatEffect(time: number, intensity: number = 1): void {
    const x = this.timeToX(time);
    const y = this.config.height / 2;
    
    this.particleSystem.createExplosion(x, y, 0x3B82F6, intensity * 50);
  }

  destroy(): void {
    this.particleSystem.destroy();
    // ... existing cleanup ...
  }
}
```

### Day 15 Testing Checklist

- [ ] Particles appear on note events
- [ ] Burst effect on new notes
- [ ] Stream effect on sustained notes
- [ ] Explosion on beat hits
- [ ] Particles fade out smoothly
- [ ] No performance degradation
- [ ] Particles stay within bounds

---

## DAY 16: RHYTHM GAME MODE

### Goal: Gamified recording interface

### Step 3: Create Rhythm Game Class

Create `lib/pixi/rhythm-game.ts`:

```typescript
import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import * as Tone from 'tone';
import { RhythmGameConfig } from '../types/pixi';

interface FallingNote {
  sprite: PIXI.Graphics;
  targetTime: number;
  midi: number;
  hit: boolean;
}

export class RhythmGameRecorder {
  private app: PIXI.Application;
  private container: PIXI.Container;
  private targetLine: PIXI.Graphics;
  private fallingNotes: FallingNote[] = [];
  private config: RhythmGameConfig;
  private score: number = 0;
  private combo: number = 0;
  private isPlaying: boolean = false;
  private startTime: number = 0;
  private recordedNotes: Array<{ time: number; midi: number }> = [];
  private bpm: number = 120;
  private onScoreUpdate?: (score: number, combo: number) => void;

  constructor(
    config: Partial<RhythmGameConfig> = {},
    onScoreUpdate?: (score: number, combo: number) => void
  ) {
    this.config = {
      fallSpeed: config.fallSpeed || 200, // pixels per second
      targetY: config.targetY || 500,
      hitWindow: config.hitWindow || 0.15, // seconds
      scoreMultiplier: config.scoreMultiplier || 10
    };

    this.onScoreUpdate = onScoreUpdate;

    // Create Pixi application
    this.app = new PIXI.Application({
      width: 800,
      height: 600,
      backgroundColor: 0x0a0a0a,
      antialias: true
    });

    this.container = new PIXI.Container();
    this.app.stage.addChild(this.container);

    // Draw target line
    this.drawTargetLine();

    // Setup keyboard input
    this.setupInput();

    // Start game loop
    this.app.ticker.add(() => this.update());
  }

  getCanvas(): HTMLCanvasElement {
    return this.app.view as HTMLCanvasElement;
  }

  /**
   * Draw the target hit line
   */
  private drawTargetLine(): void {
    this.targetLine = new PIXI.Graphics();
    this.targetLine.lineStyle(3, 0x3B82F6, 1);
    this.targetLine.moveTo(0, this.config.targetY);
    this.targetLine.lineTo(800, this.config.targetY);
    
    // Add glow
    this.targetLine.filters = [new PIXI.filters.BlurFilter(4)];
    
    this.container.addChild(this.targetLine);

    // Add hit zones
    const hitZone = new PIXI.Graphics();
    hitZone.lineStyle(2, 0x3B82F6, 0.3);
    hitZone.drawRect(
      200,
      this.config.targetY - 30,
      400,
      60
    );
    this.container.addChild(hitZone);
  }

  /**
   * Setup keyboard input
   */
  private setupInput(): void {
    window.addEventListener('keydown', (e) => {
      if (!this.isPlaying) return;

      // Map keys to MIDI notes (like Guitar Hero)
      const keyMap: { [key: string]: number } = {
        '1': 60, // C
        '2': 62, // D
        '3': 64, // E
        '4': 65, // F
        '5': 67, // G
        ' ': -1  // Any note (rhythm only)
      };

      const midi = keyMap[e.key];
      if (midi !== undefined) {
        this.handleInput(midi);
      }
    });
  }

  /**
   * Handle user input
   */
  private handleInput(midi: number): void {
    const currentTime = (performance.now() - this.startTime) / 1000;

    // Find closest falling note
    let closestNote: FallingNote | null = null;
    let minDistance = Infinity;

    this.fallingNotes.forEach(note => {
      if (note.hit) return;
      
      const distance = Math.abs(note.targetTime - currentTime);
      if (distance < minDistance && distance < this.config.hitWindow) {
        // If spacebar, accept any note. Otherwise, match MIDI
        if (midi === -1 || note.midi === midi) {
          minDistance = distance;
          closestNote = note;
        }
      }
    });

    if (closestNote) {
      // Hit!
      this.hitNote(closestNote, minDistance);
    } else {
      // Miss!
      this.miss();
    }
  }

  /**
   * Hit a note successfully
   */
  private hitNote(note: FallingNote, accuracy: number): void {
    note.hit = true;

    // Calculate score based on accuracy
    const accuracyMultiplier = 1 - (accuracy / this.config.hitWindow);
    const points = Math.floor(
      this.config.scoreMultiplier * 
      accuracyMultiplier * 
      (1 + this.combo * 0.1)
    );

    this.score += points;
    this.combo++;

    // Record the note
    this.recordedNotes.push({
      time: note.targetTime,
      midi: note.midi
    });

    // Visual feedback
    note.sprite.tint = 0x00ff00;
    gsap.to(note.sprite, {
      alpha: 0,
      scale: 1.5,
      duration: 0.3,
      onComplete: () => {
        this.container.removeChild(note.sprite);
        note.sprite.destroy();
      }
    });

    // Play sound
    this.playNoteSound(note.midi);

    // Update UI
    this.onScoreUpdate?.(this.score, this.combo);

    // Show accuracy text
    this.showAccuracyFeedback(
      note.sprite.x,
      note.sprite.y,
      accuracy < 0.05 ? 'PERFECT!' : 'GOOD!'
    );
  }

  /**
   * Miss a note
   */
  private miss(): void {
    this.combo = 0;
    this.onScoreUpdate?.(this.score, this.combo);
    
    // Flash target line red
    gsap.to(this.targetLine, {
      tint: 0xff0000,
      duration: 0.1,
      yoyo: true,
      repeat: 1,
      onComplete: () => {
        this.targetLine.tint = 0xffffff;
      }
    });
  }

  /**
   * Show accuracy feedback
   */
  private showAccuracyFeedback(x: number, y: number, text: string): void {
    const feedback = new PIXI.Text(text, {
      fontSize: 24,
      fill: text === 'PERFECT!' ? 0xffd700 : 0x3B82F6,
      fontWeight: 'bold'
    });
    
    feedback.x = x;
    feedback.y = y;
    feedback.anchor.set(0.5);
    
    this.container.addChild(feedback);
    
    gsap.to(feedback, {
      y: y - 50,
      alpha: 0,
      duration: 1,
      onComplete: () => {
        this.container.removeChild(feedback);
        feedback.destroy();
      }
    });
  }

  /**
   * Load note sequence to play
   */
  loadSequence(notes: Array<{ time: number; midi: number }>, bpm: number): void {
    this.bpm = bpm;
    
    // Create falling notes
    notes.forEach(note => {
      this.createFallingNote(note.time, note.midi);
    });
  }

  /**
   * Create a falling note
   */
  private createFallingNote(targetTime: number, midi: number): void {
    const sprite = new PIXI.Graphics();
    const color = this.getColorForMidi(midi);
    
    sprite.beginFill(color, 0.8);
    sprite.drawCircle(0, 0, 20);
    sprite.endFill();
    
    // Add border
    sprite.lineStyle(2, 0xffffff, 0.5);
    sprite.drawCircle(0, 0, 20);
    
    // Position at top, will fall down
    const laneX = 300 + ((midi - 60) * 50); // Spread across lanes
    sprite.x = laneX;
    sprite.y = 0;
    
    this.container.addChild(sprite);
    
    const fallingNote: FallingNote = {
      sprite,
      targetTime,
      midi,
      hit: false
    };
    
    this.fallingNotes.push(fallingNote);
  }

  /**
   * Start the game
   */
  start(): void {
    this.isPlaying = true;
    this.startTime = performance.now();
    this.score = 0;
    this.combo = 0;
    this.recordedNotes = [];
    
    // Start metronome
    Tone.Transport.start();
  }

  /**
   * Stop the game
   */
  stop(): Array<{ time: number; midi: number }> {
    this.isPlaying = false;
    Tone.Transport.stop();
    
    return this.recordedNotes;
  }

  /**
   * Update loop
   */
  private update(): void {
    if (!this.isPlaying) return;

    const currentTime = (performance.now() - this.startTime) / 1000;
    const deltaTime = this.app.ticker.deltaMS / 1000;

    // Update falling notes
    for (let i = this.fallingNotes.length - 1; i >= 0; i--) {
      const note = this.fallingNotes[i];
      
      if (note.hit) continue;

      // Calculate target Y position based on timing
      const timeUntilHit = note.targetTime - currentTime;
      const targetY = this.config.targetY - (timeUntilHit * this.config.fallSpeed);
      
      // Lerp toward target position
      note.sprite.y += (targetY - note.sprite.y) * 0.2;

      // Remove if past target and not hit
      if (currentTime > note.targetTime + this.config.hitWindow) {
        if (!note.hit) {
          this.miss();
        }
        
        this.container.removeChild(note.sprite);
        note.sprite.destroy();
        this.fallingNotes.splice(i, 1);
      }
    }
  }

  /**
   * Play note sound
   */
  private playNoteSound(midi: number): void {
    const synth = new Tone.Synth().toDestination();
    const note = Tone.Frequency(midi, 'midi').toNote();
    synth.triggerAttackRelease(note, '8n');
  }

  /**
   * Get color for MIDI note
   */
  private getColorForMidi(midi: number): number {
    const hue = ((midi % 12) / 12) * 360;
    return PIXI.utils.string2hex(`hsl(${hue}, 70%, 60%)`);
  }

  /**
   * Resize canvas
   */
  resize(width: number, height: number): void {
    this.app.renderer.resize(width, height);
  }

  /**
   * Cleanup
   */
  destroy(): void {
    this.fallingNotes.forEach(note => note.sprite.destroy());
    this.app.destroy(true, { children: true, texture: true });
  }
}
```

### Step 4: Create Rhythm Game Component

Create `app/components/pixi/RhythmGameMode.tsx`:

```typescript
'use client';

import { useEffect, useRef, useState } from 'react';
import { RhythmGameRecorder } from '@/lib/pixi/rhythm-game';
import { Play, Square, Trophy } from 'lucide-react';

interface RhythmGameModeProps {
  bpm: number;
  onComplete?: (notes: Array<{ time: number; midi: number }>) => void;
}

export default function RhythmGameMode({ bpm, onComplete }: RhythmGameModeProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const gameRef = useRef<RhythmGameRecorder | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [score, setScore] = useState(0);
  const [combo, setCombo] = useState(0);
  const [highScore, setHighScore] = useState(0);

  useEffect(() => {
    if (!containerRef.current) return;

    // Create game
    const game = new RhythmGameRecorder(
      { fallSpeed: 200 },
      (newScore, newCombo) => {
        setScore(newScore);
        setCombo(newCombo);
      }
    );

    containerRef.current.appendChild(game.getCanvas());
    gameRef.current = game;

    // Load demo sequence
    const demoNotes = generateDemoSequence(bpm);
    game.loadSequence(demoNotes, bpm);

    return () => {
      game.destroy();
    };
  }, [bpm]);

  const handleStart = () => {
    if (!gameRef.current) return;
    
    gameRef.current.start();
    setIsPlaying(true);
    setScore(0);
    setCombo(0);
  };

  const handleStop = () => {
    if (!gameRef.current) return;
    
    const recordedNotes = gameRef.current.stop();
    setIsPlaying(false);
    
    // Update high score
    if (score > highScore) {
      setHighScore(score);
    }
    
    onComplete?.(recordedNotes);
  };

  return (
    <div className="space-y-4">
      {/* Score Display */}
      <div className="grid grid-cols-3 gap-4">
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800">
          <div className="text-sm text-gray-400">Score</div>
          <div className="text-3xl font-bold text-primary-500">{score}</div>
        </div>
        
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800">
          <div className="text-sm text-gray-400">Combo</div>
          <div className="text-3xl font-bold text-secondary-500">
            {combo > 0 ? `×${combo}` : '0'}
          </div>
        </div>
        
        <div className="bg-gray-900 p-4 rounded-lg border border-gray-800 flex items-center gap-2">
          <Trophy className="text-yellow-500" size={20} />
          <div>
            <div className="text-sm text-gray-400">High Score</div>
            <div className="text-2xl font-bold text-yellow-500">{highScore}</div>
          </div>
        </div>
      </div>

      {/* Instructions */}
      <div className="bg-gray-900/50 p-4 rounded-lg border border-gray-800">
        <p className="text-sm text-gray-400">
          <strong>How to Play:</strong> Press keys 1-5 or spacebar when notes reach the target line. 
          Perfect timing = higher score & combo multiplier!
        </p>
      </div>

      {/* Game Canvas */}
      <div 
        ref={containerRef}
        className="rounded-lg overflow-hidden border border-gray-800"
      />

      {/* Controls */}
      <div className="flex justify-center">
        {!isPlaying ? (
          <button
            onClick={handleStart}
            className="flex items-center gap-2 px-8 py-4 bg-primary-500 hover:bg-primary-600 rounded-lg text-lg font-semibold"
          >
            <Play size={24} />
            Start Game
          </button>
        ) : (
          <button
            onClick={handleStop}
            className="flex items-center gap-2 px-8 py-4 bg-red-500 hover:bg-red-600 rounded-lg text-lg font-semibold"
          >
            <Square size={24} />
            Stop & Save
          </button>
        )}
      </div>
    </div>
  );
}

// Generate demo sequence for testing
function generateDemoSequence(bpm: number): Array<{ time: number; midi: number }> {
  const beatDuration = 60 / bpm;
  const notes: Array<{ time: number; midi: number }> = [];
  
  // Generate a simple melody
  const melody = [60, 62, 64, 65, 67, 65, 64, 62];
  
  melody.forEach((midi, i) => {
    notes.push({
      time: i * beatDuration + 2, // Start after 2 seconds
      midi
    });
  });
  
  return notes;
}
```

### Day 16 Testing Checklist

- [ ] Notes fall from top
- [ ] Target line displays
- [ ] Spacebar registers hits
- [ ] Score updates on hits
- [ ] Combo multiplier works
- [ ] Miss breaks combo
- [ ] Perfect/Good feedback shows
- [ ] Recorded notes save correctly

---

## INTEGRATION WITH EXISTING CODE

### Update Main Page

Update `app/page.tsx` to add tabs for different modes:

```typescript
'use client';

import { useState } from 'react';
import Recorder from './components/Recorder';
import VisualizerCanvas from './components/pixi/VisualizerCanvas';
import PianoRollEditor from './components/pixi/PianoRollEditor';
import RhythmGameMode from './components/pixi/RhythmGameMode';

export default function Home() {
  const [mode, setMode] = useState<'record' | 'edit' | 'game'>('record');
  const [pitches, setPitches] = useState([]);
  const [bpm, setBpm] = useState(120);

  return (
    <main className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto space-y-8">
        <header className="text-center space-y-2">
          <h1 className="text-4xl font-bold bg-gradient-to-r from-primary-500 to-secondary-500 bg-clip-text text-transparent">
            VoxForge
          </h1>
          <p className="text-gray-400">Transform your voice into music</p>
        </header>

        {/* Mode Tabs */}
        <div className="flex justify-center gap-4">
          <button
            onClick={() => setMode('record')}
            className={`px-6 py-3 rounded-lg font-medium transition-colors ${
              mode === 'record'
                ? 'bg-primary-500 text-white'
                : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
            }`}
          >
            Record
          </button>
          
          <button
            onClick={() => setMode('edit')}
            className={`px-6 py-3 rounded-lg font-medium transition-colors ${
              mode === 'edit'
                ? 'bg-primary-500 text-white'
                : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
            }`}
          >
            Piano Roll
          </button>
          
          <button
            onClick={() => setMode('game')}
            className={`px-6 py-3 rounded-lg font-medium transition-colors ${
              mode === 'game'
                ? 'bg-primary-500 text-white'
                : 'bg-gray-800 text-gray-400 hover:bg-gray-700'
            }`}
          >
            Rhythm Game
          </button>
        </div>

        {/* Content */}
        <div className="bg-gray-900 rounded-xl p-8 border border-gray-800">
          {mode === 'record' && (
            <Recorder onRecordingComplete={(buffer) => {
              // Process recording...
            }} />
          )}
          
          {mode === 'edit' && (
            <PianoRollEditor
              initialPitches={pitches}
              bpm={bpm}
              onExport={(newPitches) => {
                setPitches(newPitches);
              }}
            />
          )}
          
          {mode === 'game' && (
            <RhythmGameMode
              bpm={bpm}
              onComplete={(notes) => {
                // Convert game notes to pitches...
              }}
            />
          )}
        </div>
      </div>
    </main>
  );
}
```

---

## TESTING & DEBUGGING

### Debug Checklist

- [ ] Pixi.js console logs (enable with `PIXI.settings.DEBUG = true`)
- [ ] Check browser console for errors
- [ ] Verify canvas renders (inspect element)
- [ ] Test interaction events (click, drag)
- [ ] Monitor FPS (should stay above 55)
- [ ] Check memory usage (no leaks)
- [ ] Test on different browsers

### Common Issues

**Issue: Canvas doesn't appear**
```typescript
// Ensure container exists before creating Pixi app
if (!containerRef.current) {
  console.error('Container not found');
  return;
}
```

**Issue: Particles cause lag**
```typescript
// Reduce max particles
const particleSystem = new MusicParticleSystem(container, 500); // Lower number
```

**Issue: Click events not working**
```typescript
// Ensure interactive flag is set
sprite.interactive = true;
sprite.buttonMode = true;
```

**Issue: Canvas not resizing**
```typescript
// Force resize on window resize
window.addEventListener('resize', () => {
  visualizer.resize(container.clientWidth, container.clientHeight);
});
```

---

## PERFORMANCE OPTIMIZATION

### Best Practices

1. **Object Pooling**
```typescript
// Reuse objects instead of creating new ones
private particlePool: PIXI.Graphics[] = [];

getParticle(): PIXI.Graphics {
  return this.particlePool.find(p => !p.visible) || new PIXI.Graphics();
}
```

2. **Batch Rendering**
```typescript
// Use ParticleContainer for many similar sprites
const particles = new PIXI.ParticleContainer(10000, {
  scale: true,
  position: true,
  alpha: true
});
```

3. **Texture Atlas**
```typescript
// Combine multiple images into one texture
const spritesheet = PIXI.Spritesheet.from('atlas.json');
```

4. **Limit Updates**
```typescript
// Only update visible objects
if (sprite.x < -100 || sprite.x > width + 100) {
  sprite.visible = false;
}
```

5. **Use Filters Sparingly**
```typescript
// Filters are expensive, use only when needed
if (needsGlow) {
  sprite.filters = [new PIXI.filters.BlurFilter()];
} else {
  sprite.filters = [];
}
```

---

## COMPLETION CHECKLIST

### Day 13 Complete
- [ ] Real-time visualization working
- [ ] Notes display as you sing
- [ ] Canvas scrolls during recording
- [ ] Colors coded by pitch
- [ ] No performance issues

### Day 14 Complete
- [ ] Piano roll displays correctly
- [ ] Can add notes by clicking
- [ ] Can drag to move notes
- [ ] Can delete notes
- [ ] Snap to grid works
- [ ] Export to MIDI works

### Day 15 Complete
- [ ] Particle system initialized
- [ ] Bursts on note events
- [ ] Particles fade smoothly
- [ ] No lag with many particles
- [ ] Different effect types work

### Day 16 Complete
- [ ] Rhythm game displays
- [ ] Notes fall correctly
- [ ] Input registers
- [ ] Score calculates
- [ ] Combo system works
- [ ] Recorded notes save

---

## NEXT STEPS

After implementing all Pixi.js features:

1. **Polish UI** - Add tooltips, help text
2. **Add Presets** - Save/load particle configurations
3. **More Instruments** - Extend to other instruments
4. **Mobile Support** - Touch controls for piano roll
5. **Export Options** - Save as image/video
6. **Multiplayer** - Collaborate in rhythm game

---

**END OF PIXI.JS INTEGRATION GUIDE**

**Total Implementation Time:** 4 days  
**Lines of Code:** ~2000  
**Dependencies Added:** 3 (pixi.js, gsap, @pixi/graphics-extras)  
**Bundle Size Impact:** +500KB gzipped

🎉 **You now have a complete interactive music creation tool!**
