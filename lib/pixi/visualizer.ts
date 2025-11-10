import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { PitchPoint } from '../types';
import { VisualizerConfig } from '../types/pixi';
import { MusicParticleSystem } from './particle-system';

export class MusicVisualizer {
  private app: PIXI.Application;
  private container!: PIXI.Container;
  private noteLayer!: PIXI.Container;
  private gridLayer!: PIXI.Container;
  private particleLayer!: PIXI.Container;
  private particleSystem!: MusicParticleSystem;
  private config: VisualizerConfig;
  private currentTime: number = 0;
  private notes: Map<string, PIXI.Graphics> = new Map();
  private initPromise: Promise<void>;

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

    // Create Pixi application (v8 API - init separately)
    this.app = new PIXI.Application();
    this.initPromise = this.app.init({
      width: this.config.width,
      height: this.config.height,
      backgroundColor: this.config.backgroundColor,
      antialias: true,
      resolution: window.devicePixelRatio || 1,
      autoDensity: true
    }).then(() => {
      // Create layers after init
      this.container = new PIXI.Container();
      this.gridLayer = new PIXI.Container();
      this.noteLayer = new PIXI.Container();
      this.particleLayer = new PIXI.Container();

      this.container.addChild(this.gridLayer);
      this.container.addChild(this.noteLayer);
      this.container.addChild(this.particleLayer);
      this.app.stage.addChild(this.container);

      // Initialize particle system
      this.particleSystem = new MusicParticleSystem(this.particleLayer, 2000);
      
      // Start update loop - ticker should be available after init
      if (this.app.ticker) {
        this.app.ticker.add(() => {
          const deltaTime = this.app.ticker?.deltaMS ? this.app.ticker.deltaMS / 1000 : 0.016; // Default to ~60fps
          this.particleSystem.update(deltaTime);
        });
      }

      // Draw background grid
      this.drawGrid();
    });
  }

  /**
   * Wait for initialization to complete
   */
  async waitForInit(): Promise<void> {
    await this.initPromise;
  }

  /**
   * Get the Pixi canvas element to add to DOM
   * Note: Wait for init to complete before calling this
   */
  getCanvas(): HTMLCanvasElement {
    // Use app.canvas instead of app.view in v8
    const canvas = (this.app.canvas || this.app.view) as HTMLCanvasElement;
    if (!canvas) {
      throw new Error('Canvas not available. Wait for initialization to complete.');
    }
    // Ensure canvas doesn't block other elements
    canvas.style.position = 'relative';
    canvas.style.zIndex = '1';
    return canvas;
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
      
      graphics.setStrokeStyle({ width: 1, color: 0x333333, alpha: 0.5 });
      graphics.moveTo(0, y);
      graphics.lineTo(this.config.width, y);
    }

    // Vertical lines (beat markers)
    const beatsToShow = Math.ceil(this.config.width / this.config.pixelsPerSecond);
    for (let i = 0; i <= beatsToShow; i++) {
      const x = i * this.config.pixelsPerSecond;
      
      graphics.setStrokeStyle({ width: 1, color: 0x222222, alpha: 0.3 });
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
    noteGraphics.filters = [new PIXI.BlurFilter({ strength: 2, quality: 4 })];

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
      
      // Create particle burst at note position
      this.particleSystem.createBurst(x, y, {
        color,
        size: 3,
        lifetime: 1
      });
    }
  }
  
  /**
   * Create beat explosion
   */
  addBeatEffect(time: number, intensity: number = 1): void {
    const x = this.timeToX(time);
    const y = this.config.height / 2;
    
    this.particleSystem.createExplosion(x, y, 0x3B82F6, intensity * 50);
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
    cursor.setStrokeStyle({ width: 2, color: 0x3B82F6, alpha: 1 });
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
    if (this.particleSystem) {
      this.particleSystem.destroy();
    }
    this.clear();
    if (this.app) {
      this.app.destroy(true, { children: true, texture: true });
    }
  }
}

