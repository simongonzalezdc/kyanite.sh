import { 
  frequencyToMidi, 
  midiToFrequency, 
  downloadBlob, 
  formatDuration, 
  trimAudioBuffer 
} from '../audio-utils'
import { createAudioBuffer } from '@/__tests__/utils/testDataFactories'

// Mock URL and document for testing
const mockURL = {
  createObjectURL: jest.fn().mockReturnValue('blob:mock-url'),
  revokeObjectURL: jest.fn(),
}

const mockDocument = {
  createElement: jest.fn().mockReturnValue({
    href: '',
    download: '',
    click: jest.fn(),
  }),
  body: {
    appendChild: jest.fn(),
    removeChild: jest.fn(),
  },
}

Object.defineProperty(global, 'URL', { value: mockURL })
Object.defineProperty(global, 'document', { value: mockDocument })

describe('audio-utils', () => {
  describe('frequencyToMidi', () => {
    it('should convert A4 (440Hz) to MIDI note 69', () => {
      const midi = frequencyToMidi(440)
      expect(midi).toBe(69)
    })

    it('should convert C4 (261.63Hz) to MIDI note 60', () => {
      const midi = frequencyToMidi(261.63)
      expect(midi).toBeCloseTo(60, 0)
    })

    it('should handle frequency range correctly', () => {
      const lowFreq = 20
      const highFreq = 20000
      
      const lowMidi = frequencyToMidi(lowFreq)
      const highMidi = frequencyToMidi(highFreq)
      
      expect(lowMidi).toBeGreaterThanOrEqual(0)
      expect(highMidi).toBeGreaterThanOrEqual(0)
    })

    it('should handle edge cases', () => {
      expect(frequencyToMidi(0)).toBeLessThan(0)
      expect(frequencyToMidi(-440)).toBeNaN()
      expect(frequencyToMidi(NaN)).toBeNaN()
    })
  })

  describe('midiToFrequency', () => {
    it('should convert MIDI note 69 to A4 (440Hz)', () => {
      const frequency = midiToFrequency(69)
      expect(frequency).toBeCloseTo(440, 2)
    })

    it('should convert MIDI note 60 to C4 (261.63Hz)', () => {
      const frequency = midiToFrequency(60)
      expect(frequency).toBeCloseTo(261.63, 2)
    })

    it('should handle MIDI range correctly', () => {
      const lowMidi = 0
      const highMidi = 127
      
      const lowFreq = midiToFrequency(lowMidi)
      const highFreq = midiToFrequency(highMidi)
      
      expect(lowFreq).toBeGreaterThan(0)
      expect(highFreq).toBeGreaterThan(lowFreq)
    })

    it('should handle edge cases', () => {
      expect(midiToFrequency(-1)).toBeLessThan(8.18) // Below MIDI range
      expect(midiToFrequency(128)).toBeGreaterThan(12543) // Above MIDI range
      expect(midiToFrequency(NaN)).toBeNaN()
    })
  })

  describe('downloadBlob', () => {
    beforeEach(() => {
      jest.clearAllMocks()
    })

    it('should create download link and trigger download', () => {
      const blob = new Blob(['test content'], { type: 'text/plain' })
      const filename = 'test.txt'

      downloadBlob(blob, filename)

      expect(mockURL.createObjectURL).toHaveBeenCalledWith(blob)
      expect(mockDocument.createElement).toHaveBeenCalledWith('a')
      
      const linkElement = mockDocument.createElement.mock.results[0].value
      expect(linkElement.href).toBe('blob:mock-url')
      expect(linkElement.download).toBe(filename)
      expect(mockDocument.body.appendChild).toHaveBeenCalledWith(linkElement)
      expect(linkElement.click).toHaveBeenCalled()
      expect(mockDocument.body.removeChild).toHaveBeenCalledWith(linkElement)
      expect(mockURL.revokeObjectURL).toHaveBeenCalledWith('blob:mock-url')
    })

    it('should handle different blob types', () => {
      const audioBlob = new Blob(['audio data'], { type: 'audio/wav' })
      const filename = 'test.wav'

      downloadBlob(audioBlob, filename)

      expect(mockURL.createObjectURL).toHaveBeenCalledWith(audioBlob)
      expect(mockDocument.createElement).toHaveBeenCalledWith('a')
    })
  })

  describe('formatDuration', () => {
    it('should format seconds as MM:SS', () => {
      expect(formatDuration(0)).toBe('0:00')
      expect(formatDuration(5)).toBe('0:05')
      expect(formatDuration(65)).toBe('1:05')
      expect(formatDuration(125)).toBe('2:05')
      expect(formatDuration(3600)).toBe('60:00')
    })

    it('should handle edge cases', () => {
      expect(formatDuration(-1)).toBe('0:-01')
      expect(formatDuration(0.5)).toBe('0:00')
      expect(formatDuration(59.9)).toBe('0:59')
    })

    it('should handle large durations', () => {
      expect(formatDuration(7200)).toBe('120:00')
      expect(formatDuration(3661)).toBe('61:01')
    })
  })

  describe('trimAudioBuffer', () => {
    it('should trim audio buffer to specified range', async () => {
      const originalBuffer = createAudioBuffer(4, 44100)
      const startTime = 1
      const endTime = 3

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(2) // 3 - 1 = 2 seconds
      expect(trimmedBuffer.sampleRate).toBe(originalBuffer.sampleRate)
      expect(trimmedBuffer.numberOfChannels).toBe(originalBuffer.numberOfChannels)
      expect(trimmedBuffer.length).toBe(Math.floor(2 * trimmedBuffer.sampleRate))
    })

    it('should handle full buffer trim', async () => {
      const originalBuffer = createAudioBuffer(2, 44100)
      const startTime = 0
      const endTime = 2

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(2)
      expect(trimmedBuffer.length).toBe(originalBuffer.length)
    })

    it('should clamp values to valid range', async () => {
      const originalBuffer = createAudioBuffer(2, 44100)
      const startTime = -1 // Should be clamped to 0
      const endTime = 5 // Should be clamped to 2 (duration)

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(2)
    })

    it('should handle single channel audio', async () => {
      const originalBuffer = createAudioBuffer(2, 44100, 1)
      const startTime = 0.5
      const endTime = 1.5

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(1)
      expect(trimmedBuffer.numberOfChannels).toBe(1)
    })

    it('should handle multi-channel audio', async () => {
      const originalBuffer = createAudioBuffer(2, 44100, 2)
      const startTime = 0.5
      const endTime = 1.5

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(1)
      expect(trimmedBuffer.numberOfChannels).toBe(2)
    })

    it('should preserve audio data correctly', async () => {
      const originalBuffer = createAudioBuffer(2, 44100)
      // Fill with test data
      const originalChannel0 = originalBuffer.getChannelData(0)
      const originalChannel1 = originalBuffer.getChannelData(1)
      
      for (let i = 0; i < originalChannel0.length; i++) {
        originalChannel0[i] = Math.sin(i * 0.01)
        originalChannel1[i] = Math.cos(i * 0.01)
      }

      const startTime = 0.5
      const endTime = 1.5
      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      const trimmedChannel0 = trimmedBuffer.getChannelData(0)
      const trimmedChannel1 = trimmedBuffer.getChannelData(1)
      
      const startSample = Math.floor(startTime * originalBuffer.sampleRate)
      const endSample = Math.floor(endTime * originalBuffer.sampleRate)
      
      // Check that data is preserved correctly
      for (let i = 0; i < trimmedChannel0.length; i++) {
        expect(trimmedChannel0[i]).toBeCloseTo(originalChannel0[startSample + i], 5)
        expect(trimmedChannel1[i]).toBeCloseTo(originalChannel1[startSample + i], 5)
      }
    })

    it('should throw error for invalid range', async () => {
      const originalBuffer = createAudioBuffer(2, 44100)
      const startTime = 2
      const endTime = 1 // End before start

      await expect(trimAudioBuffer(originalBuffer, startTime, endTime))
        .rejects.toThrow('Invalid trim range: end time must be greater than start time')
    })

    it('should throw error for zero-length trim', async () => {
      const originalBuffer = createAudioBuffer(2, 44100)
      const startTime = 1
      const endTime = 1 // Same as start

      await expect(trimAudioBuffer(originalBuffer, startTime, endTime))
        .rejects.toThrow('Invalid trim range: end time must be greater than start time')
    })

    it('should handle different sample rates', async () => {
      const originalBuffer = createAudioBuffer(2, 48000)
      const startTime = 0.5
      const endTime = 1.5

      const trimmedBuffer = await trimAudioBuffer(originalBuffer, startTime, endTime)

      expect(trimmedBuffer.duration).toBe(1)
      expect(trimmedBuffer.sampleRate).toBe(48000)
      expect(trimmedBuffer.length).toBe(Math.floor(1 * 48000))
    })
  })

  describe('integration tests', () => {
    it('should work together in typical workflow', () => {
      // Test frequency to MIDI and back
      const originalFreq = 440
      const midi = frequencyToMidi(originalFreq)
      const convertedFreq = midiToFrequency(midi)
      
      expect(convertedFreq).toBeCloseTo(originalFreq, 2)

      // Test duration formatting
      const duration = 125.5
      const formatted = formatDuration(duration)
      expect(formatted).toBe('2:05')
    })
  })
})