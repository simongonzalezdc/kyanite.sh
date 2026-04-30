import { BPMDetector } from '../bpm-detector'
import { createAudioBuffer, createBPMAnalysis } from '@/__tests__/utils/testDataFactories'
import { setupAudioMocks } from '@/__tests__/utils/audioMocks'

// Setup mocks before tests
beforeAll(() => {
  setupAudioMocks()
})

describe('BPMDetector', () => {
  let bpmDetector: BPMDetector

  beforeEach(() => {
    bpmDetector = new BPMDetector()
  })

  describe('analyze', () => {
    it('should return BPM analysis for audio with clear rhythm', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Create a rhythmic pattern (120 BPM = 2 beats per second)
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        // Create beats every 0.5 seconds
        if (i % beatInterval < 100) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 5)
      expect(analysis.confidence).toBeGreaterThan(0)
      expect(typeof analysis.stable).toBe('boolean')
    })

    it('should handle silent audio', async () => {
      const silentBuffer = createAudioBuffer(1, 44100)
      // Fill with silence
      const channelData = silentBuffer.getChannelData(0)
      channelData.fill(0)

      const analysis = await bpmDetector.analyze(silentBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBe(120) // Default fallback
      expect(analysis.confidence).toBeLessThan(0.5)
      expect(analysis.stable).toBe(false)
    })

    it('should handle very fast tempo', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Create a fast rhythm (200 BPM = 3.33 beats per second)
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 3.33
      
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 50) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 50) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(200, 5)
      expect(analysis.confidence).toBeGreaterThan(0)
    })

    it('should handle very slow tempo', async () => {
      const audioBuffer = createAudioBuffer(4, 44100)
      // Create a slow rhythm (60 BPM = 1 beat per second)
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate // 60 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 200) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 200 * (i % 200) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(60, 5)
      expect(analysis.confidence).toBeGreaterThan(0)
    })

    it('should handle irregular rhythm', async () => {
      const audioBuffer = createAudioBuffer(3, 44100)
      // Create irregular rhythm
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      
      for (let i = 0; i < channelData.length; i++) {
        // Random beats
        if (Math.random() < 0.02) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * Math.random() / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeGreaterThanOrEqual(60)
      expect(analysis.bpm).toBeLessThanOrEqual(200)
      expect(analysis.confidence).toBeLessThan(0.7)
      expect(analysis.stable).toBe(false)
    })

    it('should handle multi-channel audio', async () => {
      const audioBuffer = createAudioBuffer(2, 44100, 2)
      // Create rhythm in both channels
      const leftChannel = audioBuffer.getChannelData(0)
      const rightChannel = audioBuffer.getChannelData(1)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < leftChannel.length; i++) {
        if (i % beatInterval < 100) {
          leftChannel[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
          rightChannel[i] = 0.6 * Math.sin(2 * Math.PI * 800 * (i % 100) / sampleRate)
        } else {
          leftChannel[i] = 0
          rightChannel[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 5)
    })

    it('should handle audio with noise', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Create rhythm with background noise
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        // Add noise
        channelData[i] = (Math.random() - 0.5) * 0.1
        
        // Add beats
        if (i % beatInterval < 100) {
          channelData[i] += 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 10) // Less precise due to noise
      expect(analysis.confidence).toBeGreaterThan(0.3)
    })

    it('should handle very short audio', async () => {
      const shortBuffer = createAudioBuffer(0.5, 44100)
      // Create a few beats
      const channelData = shortBuffer.getChannelData(0)
      const sampleRate = shortBuffer.sampleRate
      
      for (let i = 0; i < channelData.length; i++) {
        if (i < sampleRate * 0.1 || (i > sampleRate * 0.25 && i < sampleRate * 0.35)) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(shortBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeGreaterThanOrEqual(60)
      expect(analysis.bpm).toBeLessThanOrEqual(200)
      expect(analysis.confidence).toBeLessThan(0.7) // Lower confidence for short audio
    })

    it('should handle different sample rates', async () => {
      const audioBuffer = createAudioBuffer(2, 48000)
      // Create rhythm at 48000 sample rate
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 100) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 5)
    })
  })

  describe('error handling', () => {
    it('should handle AudioContext creation errors', async () => {
      // Mock AudioContext to throw an error
      const originalAudioContext = global.AudioContext
      global.AudioContext = jest.fn().mockImplementation(() => {
        throw new Error('AudioContext not supported')
      }) as any

      const audioBuffer = createAudioBuffer(1, 44100)
      const analysis = await bpmDetector.analyze(audioBuffer)

      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBe(120) // Fallback value
      expect(analysis.confidence).toBeLessThan(0.5)
      expect(analysis.stable).toBe(false)

      // Restore original AudioContext
      global.AudioContext = originalAudioContext
    })

    it('should handle processor creation errors', async () => {
      // Mock createRealTimeBpmProcessor to throw an error
      const originalCreateProcessor = require('realtime-bpm-analyzer').createRealTimeBpmProcessor
      require('realtime-bpm-analyzer').createRealTimeBpmProcessor = jest.fn().mockRejectedValue(
        new Error('Processor creation failed')
      )

      const audioBuffer = createAudioBuffer(1, 44100)
      const analysis = await bpmDetector.analyze(audioBuffer)

      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeGreaterThanOrEqual(60)
      expect(analysis.bpm).toBeLessThanOrEqual(200)
      expect(analysis.confidence).toBeLessThan(0.7)

      // Restore original function
      require('realtime-bpm-analyzer').createRealTimeBpmProcessor = originalCreateProcessor
    })
  })

  describe('edge cases', () => {
    it('should handle audio with very low amplitude', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Create very quiet rhythm
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 100) {
          channelData[i] = 0.01 * Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate)
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 10)
      expect(analysis.confidence).toBeLessThan(0.5)
    })

    it('should handle audio with clipping', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Create clipped rhythm
      const channelData = audioBuffer.getChannelData(0)
      const sampleRate = audioBuffer.sampleRate
      const beatInterval = sampleRate / 2 // 120 BPM
      
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 100) {
          // Clipped signal (values > 1.0 or < -1.0)
          channelData[i] = 2.0 * Math.sign(Math.sin(2 * Math.PI * 1000 * (i % 100) / sampleRate))
        } else {
          channelData[i] = 0
        }
      }

      const analysis = await bpmDetector.analyze(audioBuffer)
      expect(analysis).toBeBPMAnalysis()
      expect(analysis.bpm).toBeCloseTo(120, 10)
    })
  })
})