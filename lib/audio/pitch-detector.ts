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
    
    // Calculate RMS for silence detection
    const calculateRMS = (data: Float32Array): number => {
      let sum = 0;
      for (let i = 0; i < data.length; i++) {
        sum += data[i] * data[i];
      }
      return Math.sqrt(sum / data.length);
    };
    
    // Find overall RMS to set threshold
    const overallRMS = calculateRMS(channelData);
    const silenceThreshold = overallRMS * 0.1; // 10% of average RMS
    
    console.log('Starting pitch analysis...');
    console.log('Audio duration:', audioBuffer.duration, 'seconds');
    console.log('Sample rate:', this.sampleRate);
    console.log('Total samples:', channelData.length);
    console.log('RMS threshold:', silenceThreshold);

    for (let i = 0; i < channelData.length - windowSize; i += hopSize) {
      const window = channelData.slice(i, i + windowSize);
      
      // Skip silent sections
      const rms = calculateRMS(window);
      if (rms < silenceThreshold) {
        continue;
      }
      
      const frequency = this.detectPitch(window);

      // More strict frequency range for voice (80Hz - 1000Hz for typical voice)
      // Allow up to 2000Hz for higher voices
      if (frequency && frequency >= 80 && frequency <= 2000) {
        const time = i / this.sampleRate;
        const midi = this.frequencyToMidi(frequency);
        
        // Clamp MIDI to valid range
        const clampedMidi = Math.max(0, Math.min(127, midi));

        pitches.push({
          frequency,
          time,
          midi: clampedMidi,
          confidence: rms / overallRMS // Use RMS as confidence indicator
        });
      }
    }

    console.log('Pitch analysis complete. Detected', pitches.length, 'pitch points');
    
    // Filter out outliers (notes that appear very briefly)
    if (pitches.length > 0) {
      const filtered = this.filterOutliers(pitches);
      console.log('After filtering outliers:', filtered.length, 'pitch points');
      return filtered;
    }
    
    return pitches;
  }
  
  private filterOutliers(pitches: PitchPoint[]): PitchPoint[] {
    if (pitches.length < 3) return pitches;
    
    // Group consecutive similar notes
    const filtered: PitchPoint[] = [];
    const minDuration = 0.05; // Minimum note duration in seconds
    
    let currentNote = pitches[0];
    let noteStart = pitches[0].time;
    let noteCount = 1;
    
    for (let i = 1; i < pitches.length; i++) {
      const pitch = pitches[i];
      const midiDiff = Math.abs(pitch.midi - currentNote.midi);
      
      // If same note (within 1 semitone) or very close in time
      if (midiDiff < 1 && (pitch.time - currentNote.time) < 0.1) {
        noteCount++;
        currentNote = pitch;
      } else {
        // Check if previous note was long enough
        const noteDuration = currentNote.time - noteStart;
        if (noteDuration >= minDuration || noteCount >= 3) {
          filtered.push({
            ...currentNote,
            time: noteStart
          });
        }
        
        // Start new note
        currentNote = pitch;
        noteStart = pitch.time;
        noteCount = 1;
      }
    }
    
    // Add last note
    const noteDuration = currentNote.time - noteStart;
    if (noteDuration >= minDuration || noteCount >= 3) {
      filtered.push({
        ...currentNote,
        time: noteStart
      });
    }
    
    return filtered;
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

