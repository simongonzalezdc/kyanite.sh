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
  private container!: PIXI.Container;
  private targetLine!: PIXI.Graphics;
  private fallingNotes: FallingNote[] = [];
  private config: RhythmGameConfig;
  private score: number = 0;
  private combo: number = 0;
  private isPlaying: boolean = false;
  private startTime: number = 0;
  private recordedNotes: Array<{ time: number; midi: number }> = [];
  private bpm: number = 120;
  private onScoreUpdate?: (score: number, combo: number) => void;
  private initPromise: Promise<void>;

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

    // Create Pixi application (v8 API - init separately)
    this.app = new PIXI.Application();
    this.initPromise = this.app.init({
      width: 800,
      height: 600,
      backgroundColor: 0x0a0a0a,
      antialias: true
    }).then(() => {
      this.container = new PIXI.Container();
      this.app.stage.addChild(this.container);

      // Draw target line
      this.drawTargetLine();

      // Setup keyboard input
      this.setupInput();

      // Start game loop - ticker should be available after init
      if (this.app.ticker) {
        this.app.ticker.add(() => this.update());
      }
    });
  }

  /**
   * Wait for initialization to complete
   */
  async waitForInit(): Promise<void> {
    await this.initPromise;
  }

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
   * Draw the target hit line
   */
  private drawTargetLine(): void {
    this.targetLine = new PIXI.Graphics();
    // Use the correct stroke API in PIXI v8
    this.targetLine.setStrokeStyle({
      width: 3,
      color: 0x3B82F6,
      alpha: 1
    });
    this.targetLine.moveTo(0, this.config.targetY);
    this.targetLine.lineTo(800, this.config.targetY);
    
    // Add glow - check if filters are available
    if (PIXI.BlurFilter) {
      this.targetLine.filters = [new PIXI.BlurFilter({ strength: 4 })];
    }
    
    // Ensure container exists before adding children
    if (this.container) {
      this.container.addChild(this.targetLine);

      // Add hit zones
      const hitZone = new PIXI.Graphics();
      hitZone.setStrokeStyle({
        width: 2,
        color: 0x3B82F6,
        alpha: 0.3
      });
      hitZone.drawRect(
        200,
        this.config.targetY - 30,
        400,
        60
      );
      this.container.addChild(hitZone);
    }
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
        if (this.container && this.container.children.includes(note.sprite)) {
          this.container.removeChild(note.sprite);
        }
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
    const feedback = new PIXI.Text({
      text,
      style: {
        fontSize: 24,
        fill: text === 'PERFECT!' ? 0xffd700 : 0x3B82F6,
        fontWeight: 'bold'
      }
    });
    
    feedback.x = x;
    feedback.y = y;
    feedback.anchor.set(0.5);
    
    // Ensure container exists before adding children
    if (this.container) {
      this.container.addChild(feedback);
      
      gsap.to(feedback, {
        y: y - 50,
        alpha: 0,
        duration: 1,
        onComplete: () => {
          if (this.container && this.container.children.includes(feedback)) {
            this.container.removeChild(feedback);
          }
          feedback.destroy();
        }
      });
    }
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
    sprite.setStrokeStyle({ width: 2, color: 0xffffff, alpha: 0.5 });
    sprite.drawCircle(0, 0, 20);
    
    // Position at top, will fall down
    const laneX = 300 + ((midi - 60) * 50); // Spread across lanes
    sprite.x = laneX;
    sprite.y = 0;
    
    // Ensure container exists before adding children
    if (this.container) {
      this.container.addChild(sprite);
    }
    
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
    const deltaTime = this.app.ticker?.deltaMS ? this.app.ticker.deltaMS / 1000 : 0.016; // Default to ~60fps

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
        
        if (this.container && this.container.children.includes(note.sprite)) {
          this.container.removeChild(note.sprite);
        }
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
   * Resize canvas
   */
  resize(width: number, height: number): void {
    // Use the correct renderer API in PIXI v8
    if (this.app.renderer) {
      this.app.renderer.resize(width, height);
    }
  }

  /**
   * Cleanup
   */
  destroy(): void {
    if (this.fallingNotes) {
      this.fallingNotes.forEach(note => {
        if (note.sprite) {
          note.sprite.destroy();
        }
      });
    }
    if (this.app) {
      // Use the correct destroy API in PIXI v8
      this.app.destroy(true);
    }
  }
}

