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

    // Weight pitches by confidence and duration
    // Use only pitches with reasonable confidence
    const validPitches = pitches.filter(p => (p.confidence || 1.0) > 0.3);
    
    if (validPitches.length === 0) {
      return {
        key: 'C Major',
        tonic: 'C',
        scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
        chordProgression: ['I', 'V', 'vi', 'IV']
      };
    }

    // Convert MIDI notes to note names, weighted by confidence
    const weightedPitchClasses = new Array(12).fill(0);
    
    validPitches.forEach((pitch, i) => {
      const midi = Math.round(pitch.midi);
      if (midi < 0 || midi > 127) return;
      
      try {
        const noteName = Note.fromMidi(midi);
        if (!noteName) return;
        
        const pc = Note.chroma(noteName);
        if (pc !== null && pc >= 0 && pc < 12) {
          // Weight by confidence and duration (longer notes are more important)
          const weight = (pitch.confidence || 1.0);
          weightedPitchClasses[pc] += weight;
        }
      } catch (e) {
        // Skip invalid notes
      }
    });

    // Find most common pitch class (likely tonic)
    const maxCount = Math.max(...weightedPitchClasses);
    if (maxCount === 0) {
      return {
        key: 'C Major',
        tonic: 'C',
        scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
        chordProgression: ['I', 'V', 'vi', 'IV']
      };
    }
    
    const tonicIndex = weightedPitchClasses.indexOf(maxCount);
    const tonicNote = Note.fromMidiSharps(60 + tonicIndex);
    const tonic = Note.get(tonicNote).pc || 'C';

    // Get unique note names from pitches
    const noteNames = Array.from(new Set(
      validPitches.map(p => {
        try {
          return Note.fromMidi(Math.round(p.midi));
        } catch {
          return null;
        }
      }).filter(n => n !== null) as string[]
    ));

    if (noteNames.length === 0) {
      return {
        key: 'C Major',
        tonic: 'C',
        scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
        chordProgression: ['I', 'V', 'vi', 'IV']
      };
    }

    // Try major and minor keys
    const majorScale = Scale.get(`${tonic} major`).notes;
    const minorScale = Scale.get(`${tonic} minor`).notes;
    const majorMatches = this.countMatchingNotes(noteNames, majorScale);
    const minorMatches = this.countMatchingNotes(noteNames, minorScale);

    // Also check weighted matches
    const majorWeighted = this.countWeightedMatches(validPitches, majorScale);
    const minorWeighted = this.countWeightedMatches(validPitches, minorScale);

    // Use both simple and weighted matching
    const majorScore = majorMatches + majorWeighted;
    const minorScore = minorMatches + minorWeighted;

    const detectedKey = majorScore > minorScore 
      ? `${tonic} Major`
      : `${tonic} Minor`;

    const scale = majorScore > minorScore ? majorScale : minorScale;

    // Debug logging (only in development)
    if (process.env.NODE_ENV === 'development') {
      console.log('Key detection:', {
        tonic,
        majorScore,
        minorScore,
        detectedKey,
        noteCount: noteNames.length
      });
    }

    return {
      key: detectedKey,
      tonic: tonic,
      scale: scale,
      chordProgression: ['I', 'V', 'vi', 'IV']
    };
  }
  
  private countWeightedMatches(pitches: PitchPoint[], scale: string[]): number {
    const scalePcs = scale.map(note => Note.chroma(note));
    let weightedMatches = 0;
    
    pitches.forEach(pitch => {
      try {
        const noteName = Note.fromMidi(Math.round(pitch.midi));
        if (!noteName) return;
        
        const pc = Note.chroma(noteName);
        if (pc !== null && scalePcs.includes(pc)) {
          weightedMatches += (pitch.confidence || 1.0);
        }
      } catch {
        // Skip invalid notes
      }
    });
    
    return weightedMatches;
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

