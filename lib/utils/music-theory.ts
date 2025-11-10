import { Chord, Key } from 'tonal';

export function getChordProgression(key: string): string[][] {
  const tonic = key.split(' ')[0];
  const isMajor = key.includes('Major');
  
  const keyInfo = isMajor ? Key.majorKey(tonic) : Key.minorKey(tonic);
  
  // I-V-vi-IV progression
  const progression = ['I', 'V', 'vi', 'IV'];
  
  return progression.map(roman => {
    const chord = Chord.getChord('major', tonic, roman);
    return chord.notes.map(note => note + '3'); // Octave 3
  });
}

export function transposeNote(note: string, semitones: number): string {
  const { Note } = require('tonal');
  return Note.transpose(note, `${semitones > 0 ? '+' : ''}${semitones}`);
}

