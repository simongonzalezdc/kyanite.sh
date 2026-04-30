import { PitchDetector } from '../pitch-detector'
import { createAudioBuffer, createPitchPoints } from '@/__tests__/utils/testDataFactories'
import { setupAudioMocks } from '@/__tests__/utils/audioMocks'

// Setup mocks before tests
beforeAll(() => {
  setupAudioMocks()
})

describe('PitchDetector', () => {
  let pitchDetector: PitchDetector

  beforeEach(() => {
    pitchDetector = new PitchDetector(44100)
  })

  describe('constructor', () => {
    it('should create a PitchDetector with default sample rate', () => {
      const detector = new PitchDetector()
      expect(detector).toBeInstanceOf(PitchDetector)
    })

    it('should create a PitchDetector with custom sample rate', () => {
      const detector = new PitchDetector(48000)
      expect(detector).toBeInstanceOf(PitchDetector)
    })
  })

  describe('analyze', () => {
    it('should return empty array for silent audio', () => {
      const silentBuffer = createAudioBuffer(1, 44100)
      // Fill with silence
      const channelData = silentBuffer.getChannelData(0)
      channelData.fill(0)

      const pitches = pitchDetector.analyze(silentBuffer)
      expect(pitches).toEqual([])
    })

    it('should detect pitches in audio with sine wave', () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Fill with 440Hz sine wave
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      for (let i = 0; i < channelData.length; i++) {
        const t = i / sampleRate
        channelData[i] = 0.5 * Math.sin(2 * Math.PI * 440 * t)
      }

      const pitches = pitchDetector.analyze(audioBuffer)
      expect(pitches.length).toBeGreaterThan(0)
      expect(pitches[0]).toBePitchPoint()
      expect(pitches[0].frequency).toBeCloseTo(440, 0)
    })

    it('should filter out frequencies outside voice range', () => {
      const audioBuffer = createAudioBuffer(1, 44100)
      // Fill with very low frequency (50Hz)
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      for (let i = 0; i < channelData.length; i++) {
        const t = i / sampleRate
        channelData[i] = 0.5 * Math.sin(2 * Math.PI * 50 * t)
      }

      const pitches = pitchDetector.analyze(audioBuffer)
      expect(pitches).toEqual([])
    })

    it('should handle multi-channel audio', () => {
      const audioBuffer = createAudioBuffer(1, 44100, 2)
      // Fill left channel with 440Hz, right channel with 880Hz
      const leftChannel = audioBuffer.getChannelData(0)
      const rightChannel = audioBuffer.getChannelData(1)
      const sampleRate = audioBuffer.sampleRate
      
      for (let i = 0; i < leftChannel.length; i++) {
        const t = i / sampleRate
        leftChannel[i] = 0.5 * Math.sin(2 * Math.PI * 440 * t)
        rightChannel[i] = 0.5 * Math.sin(2 * Math.PI * 880 * t)
      }

      const pitches = pitchDetector.analyze(audioBuffer)
      expect(pitches.length).toBeGreaterThan(0)
      // Should detect the 440Hz from the first channel
      expect(pitches.some(p => Math.abs(p.frequency - 440) < 10)).toBe(true)
    })

    it('should filter out outlier pitches', () => {
      const audioBuffer = createAudioBuffer(3, 44100)
      // Create a stable pitch with occasional outliers
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      
      for (let i = 0; i < channelData.length; i++) {
        const t = i / sampleRate
        // Mostly 440Hz with occasional 1000Hz spikes
        if (i % 1000 < 990) {
          channelData[i] = 0.5 * Math.sin(2 * Math.PI * 440 * t)
        } else {
          channelData[i] = 0.5 * Math.sin(2 * Math.PI * 1000 * t)
        }
      }

      const pitches = pitchDetector.analyze(audioBuffer)
      expect(pitches.length).toBeGreaterThan(0)
      // Most detected pitches should be around 440Hz
      const pitchesAround440 = pitches.filter(p => Math.abs(p.frequency - 440) < 50)
      expect(pitchesAround440.length).toBeGreaterThan(pitches.length * 0.7)
    })
  })

  describe('getSimplifiedMelody', () => {
    it('should return empty array for empty input', () => {
      const melody = pitchDetector.getSimplifiedMelody([])
      expect(melody).toEqual([])
    })

    it('should remove duplicate consecutive notes', () => {
      const pitches = [
        { frequency: 440, time: 0, midi: 69, confidence: 0.8 },
        { frequency: 440, time: 0.1, midi: 69, confidence: 0.8 },
        { frequency: 523, time: 0.2, midi: 72, confidence: 0.9 },
        { frequency: 523, time: 0.3, midi: 72, confidence: 0.9 },
        { frequency: 659, time: 0.4, midi: 76, confidence: 0.7 },
      ]

      const melody = pitchDetector.getSimplifiedMelody(pitches)
      expect(melody).toEqual([69, 72, 76])
    })

    it('should round MIDI notes to nearest integer', () => {
      const pitches = [
        { frequency: 440, time: 0, midi: 69.2, confidence: 0.8 },
        { frequency: 523, time: 0.1, midi: 71.8, confidence: 0.9 },
      ]

      const melody = pitchDetector.getSimplifiedMelody(pitches)
      expect(melody).toEqual([69, 72])
    })
  })

  describe('getStats', () => {
    it('should return default stats for empty input', () => {
      const stats = pitchDetector.getStats([])
      expect(stats).toEqual({
        averageFrequency: 0,
        minFrequency: 0,
        maxFrequency: 0,
        averageMidi: 0,
        range: 'unknown',
      })
    })

    it('should calculate correct statistics for pitch data', () => {
      const pitches = createPitchPoints(5, 440)
      const stats = pitchDetector.getStats(pitches)

      expect(stats.averageFrequency).toBeGreaterThan(400)
      expect(stats.averageFrequency).toBeLessThan(500)
      expect(stats.minFrequency).toBeGreaterThan(0)
      expect(stats.maxFrequency).toBeGreaterThan(stats.minFrequency)
      expect(stats.averageMidi).toBeGreaterThan(60)
      expect(stats.range).toMatch(/^(low|mid|high)$/)
    })

    it('should classify low range correctly', () => {
      const pitches = [
        { frequency: 100, time: 0, midi: 40, confidence: 0.8 },
        { frequency: 150, time: 0.1, midi: 45, confidence: 0.9 },
      ]

      const stats = pitchDetector.getStats(pitches)
      expect(stats.range).toBe('low')
    })

    it('should classify high range correctly', () => {
      const pitches = [
        { frequency: 1000, time: 0, midi: 80, confidence: 0.8 },
        { frequency: 1200, time: 0.1, midi: 85, confidence: 0.9 },
      ]

      const stats = pitchDetector.getStats(pitches)
      expect(stats.range).toBe('high')
    })

    it('should classify mid range correctly', () => {
      const pitches = [
        { frequency: 440, time: 0, midi: 69, confidence: 0.8 },
        { frequency: 523, time: 0.1, midi: 72, confidence: 0.9 },
      ]

      const stats = pitchDetector.getStats(pitches)
      expect(stats.range).toBe('mid')
    })
  })

  describe('midiToNoteName', () => {
    it('should convert MIDI to correct note name', () => {
      expect(pitchDetector.midiToNoteName(60)).toBe('C4')
      expect(pitchDetector.midiToNoteName(69)).toBe('A4')
      expect(pitchDetector.midiToNoteName(0)).toBe('C-1')
      expect(pitchDetector.midiToNoteName(127)).toBe('G9')
    })

    it('should handle sharps correctly', () => {
      expect(pitchDetector.midiToNoteName(61)).toBe('C#4')
      expect(pitchDetector.midiToNoteName(66)).toBe('F#4')
    })

    it('should handle octave boundaries', () => {
      expect(pitchDetector.midiToNoteName(59)).toBe('B3')
      expect(pitchDetector.midiToNoteName(71)).toBe('B4')
    })
  })

  describe('frequencyToMidi', () => {
    it('should convert A4 to MIDI 69', () => {
      const midi = (pitchDetector as any).frequencyToMidi(440)
      expect(midi).toBe(69)
    })

    it('should convert C4 to MIDI 60', () => {
      const midi = (pitchDetector as any).frequencyToMidi(261.63)
      expect(midi).toBeCloseTo(60, 0)
    })

    it('should handle frequency range correctly', () => {
      const lowFreq = 20
      const highFreq = 20000
      
      const lowMidi = (pitchDetector as any).frequencyToMidi(lowFreq)
      const highMidi = (pitchDetector as any).frequencyToMidi(highFreq)
      
      expect(lowMidi).toBeGreaterThanOrEqual(0)
      expect(highMidi).toBeLessThanOrEqual(127)
    })
  })
})