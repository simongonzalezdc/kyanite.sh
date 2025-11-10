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

export interface TimeSignature {
  numerator: number;    // 4, 3, 2, 6, etc.
  denominator: number;  // 4, 8, etc.
  display: string;      // "4/4", "3/4", "2/4", "6/8", etc.
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

