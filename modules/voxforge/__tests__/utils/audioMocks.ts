/**
 * Mock utilities for Web Audio API testing
 */

declare const jest: any

export class MockAudioBuffer {
  numberOfChannels: number
  length: number
  sampleRate: number
  duration: number
  private channelData: Float32Array[]

  constructor(numberOfChannels: number = 2, length: number = 44100, sampleRate: number = 44100) {
    this.numberOfChannels = numberOfChannels
    this.length = length
    this.sampleRate = sampleRate
    this.duration = length / sampleRate
    this.channelData = []

    for (let i = 0; i < numberOfChannels; i++) {
      this.channelData.push(new Float32Array(length))
    }
  }

  getChannelData(channel: number): Float32Array {
    return this.channelData[channel] || new Float32Array(this.length)
  }

  // Helper method to fill with test data
  fillWithSineWave(frequency: number = 440, amplitude: number = 0.5): void {
    for (let channel = 0; channel < this.numberOfChannels; channel++) {
      const data = this.channelData[channel]
      for (let i = 0; i < this.length; i++) {
        const t = i / this.sampleRate
        data[i] = amplitude * Math.sin(2 * Math.PI * frequency * t)
      }
    }
  }

  // Helper method to fill with silence
  fillWithSilence(): void {
    for (let channel = 0; channel < this.numberOfChannels; channel++) {
      const data = this.channelData[channel]
      data.fill(0)
    }
  }

  // Helper method to fill with noise
  fillWithNoise(amplitude: number = 0.1): void {
    for (let channel = 0; channel < this.numberOfChannels; channel++) {
      const data = this.channelData[channel]
      for (let i = 0; i < this.length; i++) {
        data[i] = amplitude * (Math.random() * 2 - 1)
      }
    }
  }
}

export class MockAudioContext {
  sampleRate: number = 44100
  currentTime: number = 0
  state: AudioContextState = 'running'
  destination: AudioDestinationNode = {} as AudioDestinationNode

  createBuffer(numberOfChannels: number, length: number, sampleRate: number): AudioBuffer {
    return new MockAudioBuffer(numberOfChannels, length, sampleRate) as any
  }

  createBufferSource(): AudioBufferSourceNode {
    return {
      buffer: null,
      connect: jest.fn(),
      disconnect: jest.fn(),
      start: jest.fn(),
      stop: jest.fn(),
      onended: null,
      addEventListener: jest.fn(),
      removeEventListener: jest.fn(),
    } as any
  }

  createAnalyser(): AnalyserNode {
    return {
      fftSize: 2048,
      frequencyBinCount: 1024,
      getByteFrequencyData: jest.fn(),
      getByteTimeDomainData: jest.fn(),
      getFloatFrequencyData: jest.fn(),
      getFloatTimeDomainData: jest.fn(),
      connect: jest.fn(),
      disconnect: jest.fn(),
    } as any
  }

  createGain(): GainNode {
    return {
      gain: { value: 1 },
      connect: jest.fn(),
      disconnect: jest.fn(),
    } as any
  }

  createScriptProcessor(): ScriptProcessorNode {
    return {
      bufferSize: 4096,
      connect: jest.fn(),
      disconnect: jest.fn(),
      onaudioprocess: null,
    } as any
  }

  createOscillator(): OscillatorNode {
    return {
      frequency: { value: 440 },
      type: 'sine',
      connect: jest.fn(),
      disconnect: jest.fn(),
      start: jest.fn(),
      stop: jest.fn(),
    } as any
  }

  decodeAudioData(audioData: ArrayBuffer): Promise<AudioBuffer> {
    return Promise.resolve(new MockAudioBuffer() as any)
  }

  close(): void {}
  resume(): Promise<void> {
    return Promise.resolve()
  }
  suspend(): Promise<void> {
    return Promise.resolve()
  }
}

export class MockMediaRecorder {
  state: RecordingState = 'inactive'
  stream: MediaStream
  ondataavailable: ((event: BlobEvent) => void) | null = null
  onstop: ((event: Event) => void) | null = null
  onerror: ((event: Event) => void) | null = null
  onstart: ((event: Event) => void) | null = null
  onpause: ((event: Event) => void) | null = null
  onresume: ((event: Event) => void) | null = null

  constructor(stream: MediaStream) {
    this.stream = stream
  }

  start(timeslice?: number): void {
    this.state = 'recording'
    if (this.onstart) {
      this.onstart(new Event('start'))
    }
  }

  stop(): void {
    this.state = 'inactive'
    if (this.onstop) {
      this.onstop(new Event('stop'))
    }
  }

  pause(): void {
    this.state = 'paused'
    if (this.onpause) {
      this.onpause(new Event('pause'))
    }
  }

  resume(): void {
    this.state = 'recording'
    if (this.onresume) {
      this.onresume(new Event('resume'))
    }
  }

  requestData(): void {
    if (this.ondataavailable) {
      this.ondataavailable(new BlobEvent('dataavailable', { data: new Blob() }))
    }
  }
}

// Mock getUserMedia
export const mockGetUserMedia = (granted: boolean = true) => {
  if (granted) {
    return Promise.resolve({
      getTracks: () => [{
        stop: jest.fn(),
        getSettings: () => ({
          deviceId: 'default',
          groupId: 'default',
          kind: 'audioinput',
          label: 'Default',
          sampleRate: 44100,
        }),
        getCapabilities: () => ({
          sampleRate: { min: 8000, max: 96000 },
          channelCount: { min: 1, max: 2 },
        }),
      }],
      getAudioTracks: () => [],
      getVideoTracks: () => [],
      addTrack: jest.fn(),
      removeTrack: jest.fn(),
      active: true,
      id: 'mock-stream-id',
    })
  } else {
    return Promise.reject(new Error('Permission denied'))
  }
}

// Mock MediaDevices
export const mockMediaDevices = {
  getUserMedia: jest.fn().mockImplementation(mockGetUserMedia),
  enumerateDevices: jest.fn().mockResolvedValue([
    {
      deviceId: 'default',
      groupId: 'default',
      kind: 'audioinput',
      label: 'Default',
    },
    {
      deviceId: 'secondary',
      groupId: 'secondary',
      kind: 'audioinput',
      label: 'Secondary',
    },
  ]),
  ondevicechange: null,
  getDisplayMedia: jest.fn(),
  getSupportedConstraints: jest.fn().mockReturnValue({}),
  addEventListener: jest.fn(),
  removeEventListener: jest.fn(),
  dispatchEvent: jest.fn(),
}

// Setup global mocks
export const setupAudioMocks = () => {
  global.AudioContext = MockAudioContext as any
  global.MediaRecorder = MockMediaRecorder as any
  global.navigator = {
    ...global.navigator,
    mediaDevices: mockMediaDevices,
  }
}

// Create test audio data
export const createTestAudioBuffer = (
  duration: number = 1,
  frequency: number = 440,
  sampleRate: number = 44100
): AudioBuffer => {
  const length = Math.floor(duration * sampleRate)
  const buffer = new MockAudioBuffer(2, length, sampleRate) as any
  buffer.fillWithSineWave(frequency)
  return buffer
}

// Create test pitch points
export const createTestPitchPoints = (count: number = 10): any[] => {
  const pitches = []
  for (let i = 0; i < count; i++) {
    pitches.push({
      frequency: 440 + Math.random() * 440, // Random frequency between 440-880Hz
      time: i * 0.1, // 100ms intervals
      midi: Math.floor(Math.random() * 127),
      confidence: Math.random(),
    })
  }
  return pitches
}

// Create test BPM analysis
export const createTestBPMAnalysis = (bpm: number = 120): any => {
  return {
    bpm,
    confidence: Math.random() * 0.5 + 0.5, // 0.5-1.0
    stable: Math.random() > 0.3,
  }
}

// Create test key analysis
export const createTestKeyAnalysis = (key: string = 'C Major'): any => {
  return {
    key,
    tonic: key.split(' ')[0],
    scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
    chordProgression: ['I', 'V', 'vi', 'IV'],
  }
}