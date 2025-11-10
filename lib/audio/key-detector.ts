import { Note, Key, Scale } from 'tonal';
import { PitchPoint, KeyAnalysis } from '../types';

export class KeyDetector {
  analyze(pitches: PitchPoint[]): KeyAnalysis {
    if (pitches.length === 0) {
      return {
        key: 'C Major',
        tonic: 'C',
        scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
        chordProgression: ['I', 'V', 'vi', 'IV']
      };
    }

    // Convert MIDI notes to note names
    const noteNames = pitches.map(p => {
      const midi = Math.round(p.midi);
      return Note.fromMidi(midi);
    }).filter(note => note !== null) as string[];

    if (noteNames.length === 0) {
      return {
        key: 'C Major',
        tonic: 'C',
        scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
        chordProgression: ['I', 'V', 'vi', 'IV']
      };
    }

    // Count pitch class occurrences (C, C#, D, etc.)
    const pitchClasses = noteNames.map(note => Note.chroma(note));
    const counts = new Array(12).fill(0);
    pitchClasses.forEach(pc => {
      if (pc !== null && pc >= 0 && pc < 12) {
        counts[pc]++;
      }
    });

    // Find most common note (likely tonic)
    const maxCount = Math.max(...counts);
    const tonicIndex = counts.indexOf(maxCount);
    const tonicNote = Note.fromMidiSharps(60 + tonicIndex); // C4 + offset
    const tonic = Note.get(tonicNote).pc || 'C';

    // Try major and minor keys
    const majorKey = Key.majorKey(tonic);
    const minorKey = Key.minorKey(tonic);

    // Count matching notes
    const majorScale = Scale.get(`${tonic} major`).notes;
    const minorScale = Scale.get(`${tonic} minor`).notes;
    const majorMatches = this.countMatchingNotes(noteNames, majorScale);
    const minorMatches = this.countMatchingNotes(noteNames, minorScale);

    const detectedKey = majorMatches > minorMatches 
      ? `${tonic} Major`
      : `${tonic} Minor`;

    const scale = majorMatches > minorMatches ? majorScale : minorScale;

    return {
      key: detectedKey,
      tonic: tonic,
      scale: scale,
      chordProgression: ['I', 'V', 'vi', 'IV']
    };
  }

  private countMatchingNotes(noteNames: string[], scale: string[]): number {
    const scalePcs = scale.map(note => Note.chroma(note));
    let matches = 0;
    
    noteNames.forEach(note => {
      const pc = Note.chroma(note);
      if (pc !== null && scalePcs.includes(pc)) {
        matches++;
      }
    });
    
    return matches;
  }
}

