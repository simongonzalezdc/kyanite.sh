import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { PitchPoint } from '../types';
import { PianoRollConfig, PixiNote } from '../types/pixi';
import { Midi } from '@tonejs/midi';

export class InteractivePianoRoll {
  private app: PIXI.Application | null = null;
  private container!: PIXI.Container;
  private gridLayer!: PIXI.Container;
  private noteLayer!: PIXI.Container;
  private pianoLayer!: PIXI.Container;
  private config: PianoRollConfig;
  private notes: PixiNote[] = [];
  private selectedNote: PixiNote | null = null;
  private isDragging: boolean = false;
  private dragStartPos: { x: number; y: number } = { x: 0, y: 0 };
  private initPromise: Promise<void>;
  private isInitialized: boolean = false;

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

    // Initialize promise that resolves when setup is complete
    this.initPromise = this.initializeApp();
  }

  /**
   * Initialize the PixiJS application
   */
  private async initializeApp(): Promise<void> {
    try {
      // Create Pixi application with v7 API
      this.app = new PIXI.Application({
        width: this.config.width,
        height: this.config.height,
        backgroundColor: this.config.backgroundColor,
        antialias: true
      });

      // Wait for app to be ready
      await new Promise((resolve) => {
        this.app!.ticker.addOnce(resolve);
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

      this.isInitialized = true;
    } catch (error) {
      console.error('Failed to initialize PixiJS application:', error);
      throw error;
    }
  }

  /**
   * Wait for initialization to complete
   */
  async waitForInit(): Promise<void> {
    if (this.isInitialized) return;
    await this.initPromise;
  }

  getCanvas(): HTMLCanvasElement {
    if (!this.app) {
      throw new Error('Application not initialized');
    }

    // Use app.view for v7
    const canvas = this.app.view as HTMLCanvasElement;
    if (!canvas) {
      throw new Error('Canvas not available. Wait for initialization to complete.');
    }
    
    // Ensure canvas doesn't block other elements
    canvas.style.position = 'relative';
    canvas.style.zIndex = '1';
    return canvas;
  }

  /**
   * Draw piano roll grid
   */
  private drawGrid(): void {
    const graphics = new PIXI.Graphics();
    const midiRange = this.config.maxMidi - this.config.minMidi;

    // Clear existing grid
    this.gridLayer.removeChildren();

    // Horizontal lines (one per MIDI note)
    for (let i = 0; i <= midiRange; i++) {
      const midi = this.config.minMidi + i;
      const y = this.midiToY(midi);
      const isBlackKey = this.isBlackKey(midi);
      
      // Alternate colors for white/black keys
      graphics.lineStyle(1, this.config.gridColor, isBlackKey ? 0.3 : 0.2);
      graphics.moveTo(60, y);
      graphics.lineTo(this.config.width, y);
    }

    // Vertical lines (beat markers)
    const totalBeats = Math.ceil(this.config.width / this.config.pixelsPerSecond);
    for (let i = 0; i <= totalBeats; i++) {
      const x = 60 + (i * this.config.pixelsPerSecond);
      const isBarLine = i % this.config.beatsPerBar === 0;
      
      graphics.lineStyle(isBarLine ? 2 : 1, this.config.gridColor, isBarLine ? 0.5 : 0.2);
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

    // Clear existing keys
    this.pianoLayer.removeChildren();

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
    const deltaMidi = Math.round(-(deltaY / this.config.noteHeight));
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
    if (this.noteLayer) {
      this.noteLayer.removeChildren();
    }
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
    return this.hslToHex(hue, 70, 60);
  }
  
  /**
   * Convert HSL to hex color
   */
  private hslToHex(h: number, s: number, l: number): number {
    s /= 100;
    l /= 100;
    const c = (1 - Math.abs(2 * l - 1)) * s;
    const x = c * (1 - Math.abs((h / 60) % 2 - 1));
    const m = l - c / 2;
    let r = 0, g = 0, b = 0;
    
    if (0 <= h && h < 60) {
      r = c; g = x; b = 0;
    } else if (60 <= h && h < 120) {
      r = x; g = c; b = 0;
    } else if (120 <= h && h < 180) {
      r = 0; g = c; b = x;
    } else if (180 <= h && h < 240) {
      r = 0; g = x; b = c;
    } else if (240 <= h && h < 300) {
      r = x; g = 0; b = c;
    } else if (300 <= h && h < 360) {
      r = c; g = 0; b = x;
    }
    
    r = Math.round((r + m) * 255);
    g = Math.round((g + m) * 255);
    b = Math.round((b + m) * 255);
    
    return (r << 16) | (g << 8) | b;
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
    if (!this.app) return;
    
    this.config.width = width;
    this.config.height = height;
    this.app.renderer.resize(width, height);
    
    // Redraw
    if (this.gridLayer) {
      this.gridLayer.removeChildren();
    }
    if (this.pianoLayer) {
      this.pianoLayer.removeChildren();
    }
    this.drawGrid();
    this.drawPianoKeys();
  }

  destroy(): void {
    if (!this.app) return;
    
    this.clear();
    this.app.destroy(true, { children: true, texture: true });
    this.app = null;
  }
}
