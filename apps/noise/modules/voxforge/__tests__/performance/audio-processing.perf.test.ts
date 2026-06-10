import { PitchDetector } from '@/lib/audio/pitch-detector'
import { BPMDetector } from '@/lib/audio/bpm-detector'
import { createAudioBuffer, createPitchPoints } from '@/__tests__/utils/testDataFactories'
import { measurePerformance } from '@/__tests__/utils/testHelpers'

describe('Audio Processing Performance', () => {
  let pitchDetector: PitchDetector
  let bpmDetector: BPMDetector

  beforeEach(() => {
    pitchDetector = new PitchDetector(44100)
    bpmDetector = new BPMDetector()
  })

  describe('Pitch Detection Performance', () => {
    it('should process 1 second of audio within acceptable time', async () => {
      const audioBuffer = createAudioBuffer(1, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      const duration = await measurePerformance(async () => {
        return pitchDetector.analyze(audioBuffer)
      })

      // Should complete within 100ms for 1 second of audio
      expect(duration).toBeLessThan(100)
    })

    it('should process 5 seconds of audio within acceptable time', async () => {
      const audioBuffer = createAudioBuffer(5, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      const duration = await measurePerformance(async () => {
        return pitchDetector.analyze(audioBuffer)
      })

      // Should complete within 500ms for 5 seconds of audio
      expect(duration).toBeLessThan(500)
    })

    it('should handle high sample rates efficiently', async () => {
      const audioBuffer = createAudioBuffer(2, 96000) // High sample rate
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 96000)
      }

      const duration = await measurePerformance(async () => {
        return pitchDetector.analyze(audioBuffer)
      })

      // Should complete within 200ms even with high sample rate
      expect(duration).toBeLessThan(200)
    })

    it('should scale linearly with audio duration', async () => {
      const durations = [1, 2, 5, 10]
      const processingTimes = []

      for (const duration of durations) {
        const audioBuffer = createAudioBuffer(duration, 44100)
        const channelData = audioBuffer.getChannelData(0)
        for (let i = 0; i < channelData.length; i++) {
          channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
        }

        const processingTime = await measurePerformance(async () => {
          return pitchDetector.analyze(audioBuffer)
        })
        
        processingTimes.push(processingTime)
      }

      // Processing time should scale roughly linearly with duration
      // Allow some variance but should be roughly proportional
      const timePerSecond = processingTimes.map((time, index) => time / durations[index])
      const avgTimePerSecond = timePerSecond.reduce((a, b) => a + b, 0) / timePerSecond.length
      
      // Should be consistent (within 50% variance)
      const maxVariance = Math.max(...timePerSecond.map(time => Math.abs(time - avgTimePerSecond)))
      expect(maxVariance / avgTimePerSecond).toBeLessThan(0.5)
    })

    it('should handle silence efficiently', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Fill with silence
      const channelData = audioBuffer.getChannelData(0)
      channelData.fill(0)

      const duration = await measurePerformance(async () => {
        return pitchDetector.analyze(audioBuffer)
      })

      // Should process silence very quickly (early exit)
      expect(duration).toBeLessThan(50)
    })

    it('should maintain memory efficiency', async () => {
      const audioBuffer = createAudioBuffer(10, 44100) // Large buffer
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      // Measure memory before
      const memoryBefore = performance.memory?.usedJSHeapSize || 0

      await measurePerformance(async () => {
        return pitchDetector.analyze(audioBuffer)
      })

      // Measure memory after
      const memoryAfter = performance.memory?.usedJSHeapSize || 0
      const memoryIncrease = memoryAfter - memoryBefore

      // Memory increase should be reasonable (less than 50MB for 10 seconds of audio)
      expect(memoryIncrease).toBeLessThan(50 * 1024 * 1024)
    })
  })

  describe('BPM Detection Performance', () => {
    it('should process 1 second of audio within acceptable time', async () => {
      const audioBuffer = createAudioBuffer(1, 44100)
      // Create rhythmic pattern
      const channelData = audioBuffer.getChannelData(0)
      const beatInterval = 44100 / 2 // 120 BPM
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 100) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / 44100)
        } else {
          channelData[i] = 0
        }
      }

      const duration = await measurePerformance(async () => {
        return bpmDetector.analyze(audioBuffer)
      })

      // Should complete within 200ms for BPM detection
      expect(duration).toBeLessThan(200)
    })

    it('should handle complex rhythms efficiently', async () => {
      const audioBuffer = createAudioBuffer(3, 44100)
      // Create complex rhythm with varying tempos
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        const time = i / 44100
        // Varying rhythm
        if (Math.sin(time * 10) > 0.5) {
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 800 * (i % 50) / 44100)
        } else if (Math.sin(time * 5) > 0.5) {
          channelData[i] = 0.6 * Math.sin(2 * Math.PI * 600 * (i % 50) / 44100)
        } else {
          channelData[i] = 0
        }
      }

      const duration = await measurePerformance(async () => {
        return bpmDetector.analyze(audioBuffer)
      })

      // Should complete within 300ms even for complex rhythm
      expect(duration).toBeLessThan(300)
    })

    it('should handle silence efficiently', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Fill with silence
      const channelData = audioBuffer.getChannelData(0)
      channelData.fill(0)

      const duration = await measurePerformance(async () => {
        return bpmDetector.analyze(audioBuffer)
      })

      // Should process silence very quickly
      expect(duration).toBeLessThan(100)
    })

    it('should scale linearly with audio duration', async () => {
      const durations = [1, 2, 5, 10]
      const processingTimes = []

      for (const duration of durations) {
        const audioBuffer = createAudioBuffer(duration, 44100)
        // Create rhythmic pattern
        const channelData = audioBuffer.getChannelData(0)
        const beatInterval = 44100 / 2 // 120 BPM
        for (let i = 0; i < channelData.length; i++) {
          if (i % beatInterval < 100) {
            channelData[i] = 0.8 * Math.sin(2 * Math.PI * 1000 * (i % 100) / 44100)
          } else {
            channelData[i] = 0
          }
        }

        const processingTime = await measurePerformance(async () => {
          return bpmDetector.analyze(audioBuffer)
        })
        
        processingTimes.push(processingTime)
      }

      // Processing time should scale roughly linearly with duration
      const timePerSecond = processingTimes.map((time, index) => time / durations[index])
      const avgTimePerSecond = timePerSecond.reduce((a, b) => a + b, 0) / timePerSecond.length
      
      // Should be consistent (within 50% variance)
      const maxVariance = Math.max(...timePerSecond.map(time => Math.abs(time - avgTimePerSecond)))
      expect(maxVariance / avgTimePerSecond).toBeLessThan(0.5)
    })
  })

  describe('Combined Processing Performance', () => {
    it('should process pitch and BPM detection efficiently', async () => {
      const audioBuffer = createAudioBuffer(3, 44100)
      // Create test data with both pitch and rhythm
      const channelData = audioBuffer.getChannelData(0)
      const beatInterval = 44100 / 2 // 120 BPM
      for (let i = 0; i < channelData.length; i++) {
        if (i % beatInterval < 100) {
          // Combine pitch and rhythm
          channelData[i] = 0.8 * Math.sin(2 * Math.PI * 440 * i / 44100) + 
                           0.3 * Math.sin(2 * Math.PI * 1000 * (i % 100) / 44100)
        } else {
          channelData[i] = 0
        }
      }

      const duration = await measurePerformance(async () => {
        // Run both analyses
        const pitchPromise = pitchDetector.analyze(audioBuffer)
        const bpmPromise = bpmDetector.analyze(audioBuffer)
        return Promise.all([pitchPromise, bpmPromise])
      })

      // Combined processing should be efficient
      expect(duration).toBeLessThan(400)
    })

    it('should handle concurrent processing', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      const duration = await measurePerformance(async () => {
        // Run multiple analyses concurrently
        const promises = []
        for (let i = 0; i < 5; i++) {
          promises.push(pitchDetector.analyze(audioBuffer))
        }
        return Promise.all(promises)
      })

      // Concurrent processing should be efficient
      expect(duration).toBeLessThan(1000)
    })
  })

  describe('Memory and CPU Usage', () => {
    it('should maintain reasonable memory usage', async () => {
      const audioBuffer = createAudioBuffer(5, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      // Measure memory before
      const memoryBefore = performance.memory?.usedJSHeapSize || 0

      await measurePerformance(async () => {
        // Run multiple analyses
        const promises = []
        for (let i = 0; i < 10; i++) {
          promises.push(pitchDetector.analyze(audioBuffer))
        }
        return Promise.all(promises)
      })

      // Measure memory after
      const memoryAfter = performance.memory?.usedJSHeapSize || 0
      const memoryIncrease = memoryAfter - memoryBefore

      // Memory increase should be reasonable
      expect(memoryIncrease).toBeLessThan(100 * 1024 * 1024) // Less than 100MB
    })

    it('should not cause memory leaks', async () => {
      const audioBuffer = createAudioBuffer(1, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      // Run garbage collection before test
      if (global.gc) {
        global.gc()
      }

      const memoryBefore = performance.memory?.usedJSHeapSize || 0

      await measurePerformance(async () => {
        // Run many analyses to test for leaks
        for (let i = 0; i < 100; i++) {
          await pitchDetector.analyze(audioBuffer)
        }
      })

      // Force garbage collection
      if (global.gc) {
        global.gc()
      }

      const memoryAfter = performance.memory?.usedJSHeapSize || 0
      const memoryIncrease = memoryAfter - memoryBefore

      // Memory should be cleaned up properly
      expect(memoryIncrease).toBeLessThan(10 * 1024 * 1024) // Less than 10MB increase
    })

    it('should maintain reasonable CPU usage', async () => {
      const audioBuffer = createAudioBuffer(2, 44100)
      // Fill with test data
      const channelData = audioBuffer.getChannelData(0)
      for (let i = 0; i < channelData.length; i++) {
        channelData[i] = Math.sin(2 * Math.PI * 440 * i / 44100)
      }

      const startTime = performance.now()

      await measurePerformance(async () => {
        // Run CPU-intensive analysis
        for (let i = 0; i < 50; i++) {
          await pitchDetector.analyze(audioBuffer)
        }
      })

      const endTime = performance.now()
      const totalTime = endTime - startTime
      const avgTimePerAnalysis = totalTime / 50

      // Average time per analysis should be reasonable
      expect(avgTimePerAnalysis).toBeLessThan(50) // Less than 50ms per analysis
    })
  })
})