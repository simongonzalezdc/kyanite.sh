import { Chord, Interval, Key, Note } from 'tonal';

export function getChordProgression(key: string): string[][] {
  const { tonic, isMajor } = parseKeyName(key);
  if (!tonic) return [[], [], [], []];

  const tonicChroma = Note.get(tonic).chroma;
  if (tonicChroma == null) return [[], [], [], []];

  const majorKey = Key.majorKey(tonic);
  const minorKey = Key.minorKey(tonic);
  const scale = isMajor ? majorKey.scale : minorKey.natural.scale;
  const harmonicScale = isMajor ? majorKey.scale : minorKey.harmonic.scale;

  const progression = isMajor
    ? [
        { root: scale[0], quality: '' },
        { root: scale[4], quality: '' },
        { root: scale[5], quality: 'm' },
        { root: scale[3], quality: '' },
      ]
    : [
        { root: scale[0], quality: 'm' },
        { root: harmonicScale[4], quality: '' },
        { root: scale[5], quality: '' },
        { root: scale[3], quality: 'm' },
      ];

  return progression.map(({ root, quality }) => {
    const rootChroma = Note.get(root).chroma;
    if (rootChroma == null) return [];
    const rootOffset = (rootChroma - tonicChroma + 12) % 12;
    const rootOctave = 3 + Math.floor((tonicChroma + rootOffset) / 12);
    const chord = Chord.get(`${root}${quality}`).notes;
    return addChordOctaves(chord, root, rootOctave);
  });
}

export function transposeNote(note: string, semitones: number): string {
  const interval = Interval.fromSemitones(semitones);
  const transposed = Note.transpose(note, interval);
  return normalizeToSharps(transposed);
}

function parseKeyName(key: string): { tonic: string; isMajor: boolean } {
  const [rawTonic = '', rawMode = 'Major'] = key.trim().split(/\s+/);
  const tonic = normalizeTonic(rawTonic);
  const isMajor = rawMode.toLowerCase() !== 'minor';
  return { tonic, isMajor };
}

function normalizeTonic(note: string): string {
  if (!note) return '';
  return note.charAt(0).toUpperCase() + note.slice(1);
}

function addChordOctaves(notes: string[], root: string, rootOctave: number): string[] {
  const rootChroma = Note.get(root).chroma;
  if (rootChroma == null) return [];

  return notes.map(note => {
    const chroma = Note.get(note).chroma;
    if (chroma == null) return note;
    const offset = (chroma - rootChroma + 12) % 12;
    const octave = rootOctave + Math.floor((rootChroma + offset) / 12);
    return `${note}${octave}`;
  });
}

function normalizeToSharps(note: string): string {
  const parsed = Note.get(note);
  if (parsed.chroma == null) return note;
  const sharpNames = ['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B'];
  const octave = parsed.oct === undefined ? '' : parsed.oct;
  return `${sharpNames[parsed.chroma]}${octave}`;
}
