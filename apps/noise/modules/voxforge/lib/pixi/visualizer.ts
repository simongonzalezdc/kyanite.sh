import * as PIXI from 'pixi.js';
import { gsap } from 'gsap';
import { PitchPoint } from '../types';
import { VisualizerConfig } from '../types/pixi';

export class MusicVisualizer {
  private app: PIXI.Application | null = null;
  private container: PIXI.Container | null = null;
  private noteLayer: PIXI.Container | null = null;
  private gridLayer: PIXI.Container | null = null;
  private config: VisualizerConfig;
  private currentTime: number = 0;
  private notes: Map<string, PIXI.Graphics> = new Map();
  private initPromise: Promise<void>;
  private isInitialized: boolean = false;

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

    // Initialize the application properly
    this.initPromise = this.initializeApp();
  }

  /**
   * Initialize the Pixi Application
   */
  private async initializeApp(): Promise<void> {
    try {
      // Create the application with v7-compatible API
      this.app = new PIXI.Application();
      
      // Initialize the application with fallback options
      const initOptions = {
        width: this.config.width,
        height: this.config.height,
        backgroundColor: this.config.backgroundColor,
        antialias: true,
        resolution: window.devicePixelRatio || 1,
        autoDensity: true
      };

      // Try to initialize, with fallback to simpler options
      try {
        await this.app.init(initOptions);
      } catch (initError) {
        console.warn('PixiJS v8 init failed, trying v7 fallback:', initError);
        // Fallback: try to create with v7 style
        this.app = new PIXI.Application(initOptions as any);
        await new Promise((resolve, reject) => {
          this.app!.ticker.addOnce(() => resolve(void 0));
          setTimeout(() => resolve(void 0), 100); // Fallback timeout
        });
      }

      // Create layers after initialization
      this.container = new PIXI.Container();
      this.gridLayer = new PIXI.Container();
      this.noteLayer = new PIXI.Container();

      this.container.addChild(this.gridLayer);
      this.container.addChild(this.noteLayer);
      
      // Safe stage access with null checks
      if (this.app && (this.app.stage || (this.app as any).view)) {
        this.app.stage.addChild(this.container);
      }

      // Draw background grid
      this.drawGrid();
      
      this.isInitialized = true;
      console.log('MusicVisualizer initialized successfully');
    } catch (error) {
      console.error('Failed to initialize MusicVisualizer:', error);
      throw error;
    }
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
    if (!this.app) {
      throw new Error('Application not initialized. Call waitForInit() first.');
    }
    
    // Use the appropriate canvas access method
    const canvas = (this.app as any).canvas || (this.app as any).view || (this.app as any).renderer?.view;
    
    if (!canvas) {
      throw new Error('Canvas not available. Wait for initialization to complete.');
    }
    
    // Ensure canvas doesn't block other elements
    canvas.style.position = 'relative';
    canvas.style.zIndex = '1';
    canvas.style.display = 'block';
    canvas.style.width = '100%';
    canvas.style.height = '100%';
    
    return canvas as HTMLCanvasElement;
  }

  /**
   * Draw piano roll grid lines
   */
  private drawGrid(): void {
    if (!this.gridLayer) return;
    
    const graphics = new PIXI.Graphics();
    
    // Horizontal lines (one per octave)
    const midiRange = this.config.maxMidi - this.config.minMidi;
    const octaves = Math.floor(midiRange / 12);
    
    for (let i = 0; i <= octaves; i++) {
      const y = this.midiToY(this.config.minMidi + (i * 12));
      
      // Use v7 API format for lineStyle
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
    if (!this.isInitialized || !this.noteLayer) {
      console.warn('Visualizer not initialized yet, skipping note add');
      return;
    }

    try {
      const noteGraphics = new PIXI.Graphics();
      const color = this.getColorForMidi(pitch.midi);
      const y = this.midiToY(pitch.midi);
      const x = this.timeToX(pitch.time);

      // Draw note rectangle using v7 API
      noteGraphics.beginFill(color, 0.8);
      noteGraphics.drawRoundedRect(
        x,
        y - this.config.noteHeight / 2,
        5, // width (will grow as note sustains)
        this.config.noteHeight,
        2 // corner radius
      );
      noteGraphics.endFill();

      this.noteLayer.addChild(noteGraphics);

      // Store reference
      const noteId = `${pitch.time}-${pitch.midi}`;
      this.notes.set(noteId, noteGraphics);

      if (isRealtime) {
        // Fade in animation
        noteGraphics.alpha = 0;
        try {
          gsap.to(noteGraphics, {
            alpha: 1,
            duration: 0.1
          });
        } catch (error) {
          // Ignore animation errors
          noteGraphics.alpha = 1;
        }
      }
    } catch (error) {
      console.warn('Error adding note:', error);
    }
  }

  /**
   * Scroll the view (for real-time recording)
   */
  scroll(deltaTime: number): void {
    if (!this.isInitialized || !this.container) return;
    
    this.currentTime += deltaTime;
    this.container.x -= deltaTime * this.config.pixelsPerSecond;

    // Remove notes that are off-screen
    const minX = -this.container.x - 100;
    this.notes.forEach((note, id) => {
      if (note.x < minX) {
        try {
          note.destroy();
        } catch (error) {
          // Ignore destroy errors
        }
        this.notes.delete(id);
      }
    });
  }

  /**
   * Clear all notes
   */
  clear(): void {
    if (!this.isInitialized) return;
    
    this.notes.forEach(note => {
      try {
        note.destroy();
      } catch (error) {
        // Ignore destroy errors
      }
    });
    this.notes.clear();
    
    if (this.noteLayer) {
      try {
        this.noteLayer.removeChildren();
      } catch (error) {
        // Ignore removal errors
      }
    }
    
    this.currentTime = 0;
    if (this.container) {
      this.container.x = 0;
    }
  }

  /**
   * Load and display pre-recorded notes
   */
  loadNotes(pitches: PitchPoint[]): void {
    if (!this.isInitialized) return;
    
    this.clear();
    pitches.forEach(pitch => this.addNote(pitch, false));
  }

  /**
   * Resize canvas
   */
  resize(width: number, height: number): void {
    if (!this.app || !this.isInitialized) return;
    
    this.config.width = width;
    this.config.height = height;
    
    try {
      this.app.renderer.resize(width, height);
    } catch (error) {
      console.warn('Error resizing renderer:', error);
    }
    
    // Redraw grid
    if (this.gridLayer) {
      try {
        this.gridLayer.removeChildren();
      } catch (error) {
        console.warn('Error clearing grid layer:', error);
      }
    }
    this.drawGrid();
  }

  /**
   * Cleanup
   */
  destroy(): void {
    try {
      this.clear();
      
      if (this.app) {
        // Safe destruction with null checks
        try {
          this.app.destroy(true, { children: true, texture: true });
        } catch (error) {
          console.warn('Error during app destruction:', error);
        }
      }
      
      this.isInitialized = false;
      this.app = null;
      this.container = null;
      this.noteLayer = null;
      this.gridLayer = null;
    } catch (error) {
      console.warn('Error during visualizer cleanup:', error);
    }
  }
}
