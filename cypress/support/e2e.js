// Import commands.js using ES2015 syntax:
import './commands'

// Add audio-specific beforeEach setup
beforeEach(() => {
  // Mock Web Audio API for Cypress
  cy.window().then((win) => {
    // Mock AudioContext
    win.AudioContext = class MockAudioContext {
      constructor() {
        this.sampleRate = 44100
        this.currentTime = 0
        this.state = 'running'
      }

      createBuffer(numberOfChannels, length, sampleRate) {
        return {
          numberOfChannels,
          length,
          sampleRate,
          duration: length / sampleRate,
          getChannelData: () => new Float32Array(length),
        }
      }

      createBufferSource() {
        return {
          buffer: null,
          connect: () => {},
          start: () => {},
          stop: () => {},
          onended: null,
        }
      }

      createAnalyser() {
        return {
          fftSize: 2048,
          frequencyBinCount: 1024,
          getByteFrequencyData: () => {},
          getByteTimeDomainData: () => {},
          connect: () => {},
        }
      }

      createGain() {
        return {
          gain: { value: 1 },
          connect: () => {},
        }
      }

      decodeAudioData() {
        return Promise.resolve({
          numberOfChannels: 2,
          length: 44100,
          sampleRate: 44100,
          duration: 1,
          getChannelData: () => new Float32Array(44100),
        })
      }

      get destination() {
        return {}
      }

      close() {}
      resume() {}
      suspend() {}
    }

    // Mock MediaRecorder
    win.MediaRecorder = class MockMediaRecorder {
      constructor() {
        this.state = 'inactive'
        this.stream = {}
        this.ondataavailable = null
        this.onstop = null
        this.onerror = null
      }

      start() {}
      stop() {}
      pause() {}
      resume() {}
      requestData() {}
    }

    // Mock getUserMedia
    win.navigator.mediaDevices.getUserMedia = () => Promise.resolve({
      getTracks: () => [{
        stop: () => {},
        getSettings: () => ({
          deviceId: 'default',
          groupId: 'default',
          kind: 'audioinput',
          label: 'Default',
          sampleRate: 44100,
        }),
      }],
    })

    // Mock URL.createObjectURL
    win.URL.createObjectURL = () => 'blob:mock-url'
    win.URL.revokeObjectURL = () => {}
  })
})