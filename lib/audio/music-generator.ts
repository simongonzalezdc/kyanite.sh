import * as Tone from 'tone';
import { InstrumentType } from '../types';

export class MusicGenerator {
  private kick: Tone.MembraneSynth | null = null;
  private snare: Tone.NoiseSynth | null = null;
  private hihat: Tone.MetalSynth | null = null;
  private bass: Tone.Synth | null = null;
  private chords: Tone.PolySynth | null = null;
  private initialized = false;
  private selectedInstruments: InstrumentType[] = ['drums', 'bass', 'chords'];
  private drumPattern: Tone.Pattern<string> | null = null;
  private bassPattern: Tone.Sequence<string> | null = null;
  private chordPattern: Tone.Part<[string, string[]]> | null = null;

  async initialize(): Promise<void> {
    if (this.initialized) return;
    
    await Tone.start();
    
    // Create drum synths
    this.kick = new Tone.MembraneSynth({
      pitchDecay: 0.05,
      octaves: 10,
      oscillator: { type: 'sine' },
      envelope: { attack: 0.001, decay: 0.4, sustain: 0.01, release: 1.4 }
    }).toDestination();

    this.snare = new Tone.NoiseSynth({
      noise: { type: 'white' },
      envelope: { attack: 0.005, decay: 0.1, sustain: 0 }
    }).toDestination();

    this.hihat = new Tone.MetalSynth({
      envelope: { attack: 0.001, decay: 0.1, release: 0.01 },
      harmonicity: 5.1,
      modulationIndex: 32,
      resonance: 4000,
      octaves: 1.5
    }).toDestination();

    // Create bass synth
    this.bass = new Tone.Synth({
      oscillator: { type: 'sawtooth' },
      envelope: { attack: 0.01, decay: 0.1, sustain: 0.5, release: 0.5 }
    }).toDestination();

    // Create chord synth
    this.chords = new Tone.PolySynth(Tone.Synth, {
      oscillator: { type: 'triangle' },
      envelope: { attack: 0.02, decay: 0.1, sustain: 0.8, release: 1 }
    }).toDestination();

    this.initialized = true;
  }

  setBPM(bpm: number): void {
    Tone.Transport.bpm.value = bpm;
  }

  setKey(key: string): void {
    // Key is stored for chord generation
    // Implementation in generateChords
  }

  setInstruments(instruments: InstrumentType[]): void {
    if (instruments.length === 0) {
      throw new Error('At least one instrument must be selected');
    }
    this.selectedInstruments = instruments;
  }

  generateDrums(patternType: 'simple' | 'moderate' | 'busy' = 'moderate'): void {
    if (!this.initialized || !this.kick || !this.snare || !this.hihat) return;

    // Stop existing pattern
    if (this.drumPattern) {
      this.drumPattern.stop();
      this.drumPattern.dispose();
    }

    let patternArray: string[];
    
    switch (patternType) {
      case 'simple':
        // Kick on 1 & 3, Snare on 2 & 4, Hi-hat on 8ths
        patternArray = ['kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat'];
        break;
      case 'moderate':
        // More variation
        patternArray = ['kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat'];
        break;
      case 'busy':
        // Complex pattern with fills
        patternArray = ['kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat', 'kick', 'hihat', 'snare', 'hihat'];
        break;
    }

    this.drumPattern = new Tone.Pattern((time, note) => {
      if (note === 'kick' && this.kick) {
        this.kick.triggerAttackRelease('C1', '8n', time);
      }
      if (note === 'snare' && this.snare) {
        this.snare.triggerAttack(time);
      }
      if (note === 'hihat' && this.hihat) {
        this.hihat.triggerAttack(time);
      }
    }, patternArray);

    this.drumPattern.interval = '8n';
  }

  generateBass(key: string): void {
    if (!this.initialized || !this.bass) return;

    // Stop existing pattern
    if (this.bassPattern) {
      this.bassPattern.stop();
      this.bassPattern.dispose();
    }

    // Extract tonic from key (e.g., "C Major" -> "C")
    const tonic = key.split(' ')[0];
    
    // Generate I-V-vi-IV progression in the detected key
    const { Scale } = require('tonal');
    const scale = Scale.get(`${tonic} major`).notes;
    
    // I-V-vi-IV progression
    const roots = [
      scale[0] + '2', // I
      scale[0] + '2', // I
      scale[4] + '2', // V
      scale[4] + '2', // V
      scale[5] + '2', // vi
      scale[5] + '2', // vi
      scale[3] + '2', // IV
      scale[3] + '2'  // IV
    ];

    this.bassPattern = new Tone.Sequence((time, note) => {
      if (this.bass) {
        this.bass.triggerAttackRelease(note, '4n', time);
      }
    }, roots, '4n');
  }

  generateChords(key: string): void {
    if (!this.initialized || !this.chords) return;

    // Stop existing pattern
    if (this.chordPattern) {
      this.chordPattern.stop();
      this.chordPattern.dispose();
    }

    // Extract tonic from key
    const tonic = key.split(' ')[0];
    const isMajor = key.includes('Major');
    
    const { Scale } = require('tonal');
    const scale = Scale.get(`${tonic} ${isMajor ? 'major' : 'minor'}`).notes;
    
    // I-V-vi-IV progression
    const chordNotes = [
      [scale[0] + '3', scale[2] + '3', scale[4] + '3'], // I
      [scale[4] + '3', scale[6] + '3', scale[1] + '4'], // V
      [scale[5] + '3', scale[0] + '4', scale[2] + '4'], // vi
      [scale[3] + '3', scale[5] + '3', scale[0] + '4']  // IV
    ];

    const events: [string, string[]][] = chordNotes.map((notes, i) => [
      `${i * 2}m`, // Time in measures
      notes
    ]);

    this.chordPattern = new Tone.Part((time, notes) => {
      if (this.chords) {
        this.chords.triggerAttackRelease(notes, '2n', time);
      }
    }, events);
  }

  start(): void {
    if (!this.initialized) return;

    if (this.selectedInstruments.includes('drums') && this.drumPattern) {
      this.drumPattern.start(0);
    }
    if (this.selectedInstruments.includes('bass') && this.bassPattern) {
      this.bassPattern.start(0);
    }
    if (this.selectedInstruments.includes('chords') && this.chordPattern) {
      this.chordPattern.start(0);
    }

    Tone.Transport.start();
  }

  stop(): void {
    Tone.Transport.stop();
    
    if (this.drumPattern) {
      this.drumPattern.stop();
    }
    if (this.bassPattern) {
      this.bassPattern.stop();
    }
    if (this.chordPattern) {
      this.chordPattern.stop();
    }
  }

  dispose(): void {
    this.stop();
    
    this.kick?.dispose();
    this.snare?.dispose();
    this.hihat?.dispose();
    this.bass?.dispose();
    this.chords?.dispose();
    
    this.drumPattern?.dispose();
    this.bassPattern?.dispose();
    this.chordPattern?.dispose();
    
    this.initialized = false;
  }
}

