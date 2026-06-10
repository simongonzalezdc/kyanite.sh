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

