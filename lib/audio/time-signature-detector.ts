import { PitchPoint, TimeSignature } from '../types';

export class TimeSignatureDetector {
  analyze(pitches: PitchPoint[], bpm: number): TimeSignature {
    if (pitches.length === 0 || bpm <= 0) {
      return { numerator: 4, denominator: 4, display: '4/4' };
    }

    // Calculate beat intervals from pitch points
    // Group pitches by time to find strong beats
    const beatInterval = 60 / bpm; // seconds per beat
    
    // Analyze rhythm patterns
    // Look for repeating patterns in note onsets
    const noteOnsets: number[] = [];
    let lastMidi = -1;
    
    pitches.forEach((pitch, i) => {
      const currentMidi = Math.round(pitch.midi);
      // Detect note onset (new note different from previous)
      if (i === 0 || currentMidi !== lastMidi) {
        noteOnsets.push(pitch.time);
        lastMidi = currentMidi;
      }
    });

    if (noteOnsets.length < 4) {
      return { numerator: 4, denominator: 4, display: '4/4' };
    }

    // Calculate intervals between note onsets
    const intervals: number[] = [];
    for (let i = 1; i < noteOnsets.length; i++) {
      intervals.push(noteOnsets[i] - noteOnsets[i - 1]);
    }

    // Find most common interval (likely beat duration)
    const intervalCounts = new Map<number, number>();
    intervals.forEach(interval => {
      const rounded = Math.round(interval / 0.1) * 0.1; // Round to 0.1s
      intervalCounts.set(rounded, (intervalCounts.get(rounded) || 0) + 1);
    });

    let maxCount = 0;
    let mostCommonInterval = beatInterval;
    intervalCounts.forEach((count, interval) => {
      if (count > maxCount) {
        maxCount = count;
        mostCommonInterval = interval;
      }
    });

    // Analyze measure patterns
    // Look for strong beats every 4, 3, or 2 beats
    const measures: number[] = [];
    let currentMeasure = 0;
    let beatsInMeasure = 0;
    const expectedBeatDuration = beatInterval;

    noteOnsets.forEach((onset, i) => {
      if (i === 0) {
        currentMeasure = onset;
        beatsInMeasure = 1;
        return;
      }

      const timeSinceLastBeat = onset - noteOnsets[i - 1];
      const isNewBeat = Math.abs(timeSinceLastBeat - expectedBeatDuration) < expectedBeatDuration * 0.3;

      if (isNewBeat) {
        beatsInMeasure++;
        
        // Check if we've completed a measure (4, 3, or 2 beats)
        if (beatsInMeasure >= 4) {
          measures.push(beatsInMeasure);
          beatsInMeasure = 0;
        } else if (beatsInMeasure >= 3 && i === noteOnsets.length - 1) {
          // End of audio, might be 3/4
          measures.push(beatsInMeasure);
        }
      }
    });

    // Count measure lengths
    const measureCounts = new Map<number, number>();
    measures.forEach(beats => {
      measureCounts.set(beats, (measureCounts.get(beats) || 0) + 1);
    });

    // Determine time signature
    let numerator = 4;
    let denominator = 4;

    if (measureCounts.has(3) && measureCounts.get(3)! > measureCounts.get(4)!) {
      numerator = 3;
    } else if (measureCounts.has(2) && measureCounts.get(2)! > measureCounts.get(4)!) {
      numerator = 2;
    } else if (measureCounts.has(6)) {
      numerator = 6;
      denominator = 8;
    }

    // Default to 4/4 if unclear
    if (measureCounts.size === 0) {
      numerator = 4;
      denominator = 4;
    }

    // Debug logging (only in development)
    if (process.env.NODE_ENV === 'development') {
      console.log('Time signature detection:', {
        numerator,
        denominator,
        noteOnsets: noteOnsets.length,
        measures: measures.length,
        measureCounts: Object.fromEntries(measureCounts)
      });
    }

    return {
      numerator,
      denominator,
      display: `${numerator}/${denominator}`
    };
  }
}

