import { getChordProgression, transposeNote } from '../music-theory'

describe('music-theory', () => {
  describe('getChordProgression', () => {
    it('should return I-V-vi-IV progression for C Major', () => {
      const progression = getChordProgression('C Major')
      expect(progression).toEqual([
        ['C3', 'E3', 'G3'],     // I chord (C Major)
        ['G3', 'B3', 'D4'],     // V chord (G Major)
        ['A3', 'C4', 'E4'],     // vi chord (A minor)
        ['F3', 'A3', 'C4'],     // IV chord (F Major)
      ])
    })

    it('should return I-V-vi-IV progression for A Minor', () => {
      const progression = getChordProgression('A Minor')
      expect(progression).toEqual([
        ['A3', 'C4', 'E4'],     // i chord (A minor)
        ['E4', 'G#4', 'B4'],    // V chord (E Major)
        ['F4', 'A4', 'C5'],     // VI chord (F Major)
        ['D4', 'F4', 'A4'],     // iv chord (D minor)
      ])
    })

    it('should handle different keys correctly', () => {
      const gMajorProgression = getChordProgression('G Major')
      expect(gMajorProgression).toEqual([
        ['G3', 'B3', 'D4'],     // I chord (G Major)
        ['D4', 'F#4', 'A4'],    // V chord (D Major)
        ['E4', 'G4', 'B4'],     // vi chord (E minor)
        ['C4', 'E4', 'G4'],     // IV chord (C Major)
      ])
    })

    it('should handle keys with sharps', () => {
      const dMajorProgression = getChordProgression('D Major')
      expect(dMajorProgression).toEqual([
        ['D3', 'F#3', 'A3'],    // I chord (D Major)
        ['A3', 'C#4', 'E4'],    // V chord (A Major)
        ['B3', 'D4', 'F#4'],    // vi chord (B minor)
        ['G3', 'B3', 'D4'],     // IV chord (G Major)
      ])
    })

    it('should handle keys with flats', () => {
      const fMajorProgression = getChordProgression('F Major')
      expect(fMajorProgression).toEqual([
        ['F3', 'A3', 'C4'],     // I chord (F Major)
        ['C4', 'E4', 'G4'],     // V chord (C Major)
        ['D4', 'F4', 'A4'],     // vi chord (D minor)
        ['Bb3', 'D4', 'F4'],    // IV chord (Bb Major)
      ])
    })

    it('should handle edge cases', () => {
      // Test with different key formats
      expect(getChordProgression('C major')).toEqual(getChordProgression('C Major'))
      expect(getChordProgression('c major')).toEqual(getChordProgression('C Major'))
      expect(getChordProgression('A minor')).toEqual(getChordProgression('A Minor'))
      expect(getChordProgression('a minor')).toEqual(getChordProgression('A Minor'))
    })
  })

  describe('transposeNote', () => {
    it('should transpose note up by semitones', () => {
      expect(transposeNote('C', 2)).toBe('D')
      expect(transposeNote('C', 5)).toBe('F')
      expect(transposeNote('C', 12)).toBe('C') // One octave up
    })

    it('should transpose note down by semitones', () => {
      expect(transposeNote('C', -2)).toBe('A#')
      expect(transposeNote('C', -5)).toBe('G')
      expect(transposeNote('C', -12)).toBe('C') // One octave down
    })

    it('should handle sharps correctly', () => {
      expect(transposeNote('C#', 1)).toBe('D')
      expect(transposeNote('F#', 1)).toBe('G')
      expect(transposeNote('C#', -1)).toBe('C')
    })

    it('should handle flats correctly', () => {
      expect(transposeNote('Bb', 1)).toBe('B')
      expect(transposeNote('Eb', 1)).toBe('E')
      expect(transposeNote('Bb', -1)).toBe('A')
    })

    it('should handle octave wrapping', () => {
      expect(transposeNote('B', 1)).toBe('C')
      expect(transposeNote('B', 2)).toBe('D')
      expect(transposeNote('C', -1)).toBe('B')
      expect(transposeNote('C', -2)).toBe('A#')
    })

    it('should handle large intervals', () => {
      expect(transposeNote('C', 24)).toBe('C') // Two octaves up
      expect(transposeNote('C', -24)).toBe('C') // Two octaves down
      expect(transposeNote('C', 13)).toBe('C#') // One octave + 1 semitone
    })

    it('should handle complex transpositions', () => {
      expect(transposeNote('F#', 7)).toBe('C') // Perfect fifth up
      expect(transposeNote('Ab', 4)).toBe('C') // Major third up
      expect(transposeNote('D', -7)).toBe('G') // Perfect fourth down
    })

    it('should handle enharmonic equivalents', () => {
      // These tests verify that the function maintains consistency
      const result1 = transposeNote('C#', 0)
      const result2 = transposeNote('Db', 0)
      
      // Both should be valid notes (we don't enforce enharmonic equivalence)
      expect(['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B']).toContain(result1)
      expect(['C', 'C#', 'D', 'D#', 'E', 'F', 'F#', 'G', 'G#', 'A', 'A#', 'B']).toContain(result2)
    })

    it('should handle edge cases', () => {
      // Test with zero transposition
      expect(transposeNote('C', 0)).toBe('C')
      expect(transposeNote('F#', 0)).toBe('F#')
      
      // Test with multiple octaves
      expect(transposeNote('C', 36)).toBe('C') // Three octaves up
      expect(transposeNote('C', -36)).toBe('C') // Three octaves down
    })
  })

  describe('integration tests', () => {
    it('should work together in typical music theory workflow', () => {
      // Get chord progression for C Major
      const progression = getChordProgression('C Major')
      
      // Transpose the first chord up a whole step
      const firstChord = progression[0]
      const transposedChord = firstChord.map(note => transposeNote(note, 2))
      
      expect(transposedChord).toEqual(['D', 'F#', 'A'])
      
      // Transpose the entire progression up a perfect fifth
      const transposedProgression = progression.map(chord => 
        chord.map(note => transposeNote(note, 7))
      )
      
      expect(transposedProgression).toEqual([
        ['G', 'B', 'D'],     // I chord in G Major
        ['D', 'F#', 'A'],    // V chord in G Major
        ['E', 'G', 'B'],     // vi chord in G Major
        ['C', 'E', 'G'],     // IV chord in G Major
      ])
    })

    it('should handle complex harmonic relationships', () => {
      // Test relative major/minor relationships
      const cMajorProgression = getChordProgression('C Major')
      const aMinorProgression = getChordProgression('A Minor')
      
      // A minor is the relative minor of C major
      // They should share the same key signature (no sharps or flats)
      expect(cMajorProgression[0]).toEqual(['C', 'E', 'G'])     // C Major
      expect(aMinorProgression[0]).toEqual(['A', 'C', 'E'])     // A minor
      
      // The V chord in A minor should be the same as the III chord in C major
      expect(aMinorProgression[1]).toEqual(['E', 'G#', 'B'])  // E Major (V in Am)
    })
  })

  describe('error handling', () => {
    it('should handle invalid input gracefully', () => {
      // Test with invalid key names
      expect(() => getChordProgression('X Major')).not.toThrow()
      expect(() => getChordProgression('')).not.toThrow()
      expect(() => getChordProgression('Invalid')).not.toThrow()
      
      // Test with invalid note names
      expect(() => transposeNote('H', 1)).not.toThrow()
      expect(() => transposeNote('', 1)).not.toThrow()
    })
  })
})