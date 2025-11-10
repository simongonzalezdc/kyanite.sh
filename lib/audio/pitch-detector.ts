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

