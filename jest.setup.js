import '@testing-library/jest-dom'
import { TextEncoder, TextDecoder } from 'util'

// Mock Web Audio API
global.AudioContext = jest.fn().mockImplementation(() => ({
  createBuffer: jest.fn().mockReturnValue({
    numberOfChannels: 2,
    length: 44100,
    sampleRate: 44100,
    duration: 1,
    getChannelData: jest.fn().mockReturnValue(new Float32Array(44100)),
  }),
  createBufferSource: jest.fn().mockReturnValue({
    buffer: null,
    connect: jest.fn(),
    start: jest.fn(),
    stop: jest.fn(),
    onended: null,
  }),
  createAnalyser: jest.fn().mockReturnValue({
    fftSize: 2048,
    frequencyBinCount: 1024,
    getByteFrequencyData: jest.fn(),
    getByteTimeDomainData: jest.fn(),
    connect: jest.fn(),
  }),
  createGain: jest.fn().mockReturnValue({
    gain: { value: 1 },
    connect: jest.fn(),
  }),
  createScriptProcessor: jest.fn().mockReturnValue({
    bufferSize: 4096,
    connect: jest.fn(),
    onaudioprocess: null,
  }),
  decodeAudioData: jest.fn().mockResolvedValue({
    numberOfChannels: 2,
    length: 44100,
    sampleRate: 44100,
    duration: 1,
    getChannelData: jest.fn().mockReturnValue(new Float32Array(44100)),
  }),
  destination: {},
  sampleRate: 44100,
  currentTime: 0,
  state: 'running',
  close: jest.fn(),
  resume: jest.fn(),
  suspend: jest.fn(),
}))

// Mock MediaRecorder API
global.MediaRecorder = jest.fn().mockImplementation(() => ({
  start: jest.fn(),
  stop: jest.fn(),
  pause: jest.fn(),
  resume: jest.fn(),
  requestData: jest.fn(),
  ondataavailable: null,
  onstop: null,
  onerror: null,
  state: 'inactive',
  stream: {},
}))

// Mock getUserMedia
global.navigator.mediaDevices = {
  getUserMedia: jest.fn().mockResolvedValue({
    getTracks: jest.fn().mockReturnValue([{
      stop: jest.fn(),
      getSettings: jest.fn().mockReturnValue({
        deviceId: 'default',
        groupId: 'default',
        kind: 'audioinput',
        label: 'Default',
        sampleRate: 44100,
      }),
    }]),
  }),
  enumerateDevices: jest.fn().mockResolvedValue([
    {
      deviceId: 'default',
      groupId: 'default',
      kind: 'audioinput',
      label: 'Default',
    },
  ]),
}

// Mock URL.createObjectURL and revokeObjectURL
global.URL.createObjectURL = jest.fn().mockReturnValue('blob:mock-url')
global.URL.revokeObjectURL = jest.fn()

// Mock Blob
global.Blob = jest.fn().mockImplementation((content, options) => ({
  content,
  options,
  size: content ? content.length : 0,
  type: options?.type || '',
  arrayBuffer: jest.fn().mockResolvedValue(new ArrayBuffer(0)),
  text: jest.fn().mockResolvedValue(''),
  stream: jest.fn().mockReturnValue(new ReadableStream()),
  slice: jest.fn().mockReturnValue(new Blob()),
}))

// Mock File
global.File = jest.fn().mockImplementation((content, name, options) => ({
  content,
  name,
  options,
  size: content ? content.length : 0,
  type: options?.type || '',
  lastModified: Date.now(),
  arrayBuffer: jest.fn().mockResolvedValue(new ArrayBuffer(0)),
  text: jest.fn().mockResolvedValue(''),
  stream: jest.fn().mockReturnValue(new ReadableStream()),
  slice: jest.fn().mockReturnValue(new File()),
}))

// Mock fetch
global.fetch = jest.fn()

// Mock WebSocket
global.WebSocket = jest.fn().mockImplementation(() => ({
  close: jest.fn(),
  send: jest.fn(),
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  readyState: 1,
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3,
}))

// Mock ResizeObserver
global.ResizeObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}))

// Mock IntersectionObserver
global.IntersectionObserver = jest.fn().mockImplementation(() => ({
  observe: jest.fn(),
  unobserve: jest.fn(),
  disconnect: jest.fn(),
}))

// Mock requestAnimationFrame and cancelAnimationFrame
global.requestAnimationFrame = jest.fn((cb) => setTimeout(cb, 16))
global.cancelAnimationFrame = jest.fn((id) => clearTimeout(id))

// Mock performance.now
Object.defineProperty(global.performance, 'now', {
  value: jest.fn(() => Date.now()),
  writable: true,
})

// Mock TextEncoder/TextDecoder for Node.js environment
global.TextEncoder = TextEncoder
global.TextDecoder = TextDecoder

// Mock console methods in tests to reduce noise
const originalError = console.error
beforeAll(() => {
  console.error = (...args) => {
    if (
      typeof args[0] === 'string' &&
      args[0].includes('Warning: ReactDOM.render is deprecated')
    ) {
      return
    }
    originalError.call(console, ...args)
  }
})

afterAll(() => {
  console.error = originalError
})

// Add custom matchers for audio testing
expect.extend({
  toBeAudioBuffer(received) {
    const pass = received && 
      typeof received === 'object' &&
      typeof received.numberOfChannels === 'number' &&
      typeof received.length === 'number' &&
      typeof received.sampleRate === 'number' &&
      typeof received.duration === 'number' &&
      typeof received.getChannelData === 'function'

    if (pass) {
      return {
        message: () => `expected ${received} not to be a valid AudioBuffer`,
        pass: true,
      }
    } else {
      return {
        message: () => `expected ${received} to be a valid AudioBuffer`,
        pass: false,
      }
    }
  },

  toBePitchPoint(received) {
    const pass = received &&
      typeof received === 'object' &&
      typeof received.frequency === 'number' &&
      typeof received.time === 'number' &&
      typeof received.midi === 'number' &&
      received.frequency >= 20 && received.frequency <= 20000 &&
      received.time >= 0 &&
      received.midi >= 0 && received.midi <= 127

    if (pass) {
      return {
        message: () => `expected ${received} not to be a valid PitchPoint`,
        pass: true,
      }
    } else {
      return {
        message: () => `expected ${received} to be a valid PitchPoint`,
        pass: false,
      }
    }
  },

  toBeBPMAnalysis(received) {
    const pass = received &&
      typeof received === 'object' &&
      typeof received.bpm === 'number' &&
      typeof received.confidence === 'number' &&
      typeof received.stable === 'boolean' &&
      received.bpm >= 40 && received.bpm <= 250 &&
      received.confidence >= 0 && received.confidence <= 1

    if (pass) {
      return {
        message: () => `expected ${received} not to be a valid BPMAnalysis`,
        pass: true,
      }
    } else {
      return {
        message: () => `expected ${received} to be a valid BPMAnalysis`,
        pass: false,
      }
    }
  },
})