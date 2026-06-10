/**
 * Test data factories for VoxForge testing
 */

import { PitchPoint, BPMAnalysis, KeyAnalysis, TimeSignature, Section, InstrumentType } from '@/lib/types'

// Factory for creating test pitch points
export const createPitchPoint = (overrides: Partial<PitchPoint> = {}): PitchPoint => ({
  frequency: 440,
  time: 0,
  midi: 69,
  confidence: 0.8,
  ...overrides,
})

// Factory for creating multiple pitch points
export const createPitchPoints = (count: number, baseFrequency: number = 440): PitchPoint[] => {
  return Array.from({ length: count }, (_, i) => ({
    frequency: baseFrequency + (Math.random() - 0.5) * 100,
    time: i * 0.1,
    midi: Math.floor(69 + (Math.random() - 0.5) * 10),
    confidence: 0.5 + Math.random() * 0.5,
  }))
}

// Factory for creating test BPM analysis
export const createBPMAnalysis = (overrides: Partial<BPMAnalysis> = {}): BPMAnalysis => ({
  bpm: 120,
  confidence: 0.8,
  stable: true,
  ...overrides,
})

// Factory for creating test key analysis
export const createKeyAnalysis = (overrides: Partial<KeyAnalysis> = {}): KeyAnalysis => ({
  key: 'C Major',
  tonic: 'C',
  scale: ['C', 'D', 'E', 'F', 'G', 'A', 'B'],
  chordProgression: ['I', 'V', 'vi', 'IV'],
  ...overrides,
})

// Factory for creating test time signature
export const createTimeSignature = (overrides: Partial<TimeSignature> = {}): TimeSignature => ({
  numerator: 4,
  denominator: 4,
  display: '4/4',
  ...overrides,
})

// Factory for creating test sections
export const createSection = (overrides: Partial<Section> = {}): Section => ({
  id: `section-${Date.now()}`,
  name: 'Test Section',
  audioBuffer: null,
  duration: 10,
  pitches: createPitchPoints(5),
  bpm: createBPMAnalysis(),
  key: createKeyAnalysis(),
  instruments: ['drums', 'bass', 'chords'],
  ...overrides,
})

// Factory for creating multiple sections
export const createSections = (count: number): Section[] => {
  return Array.from({ length: count }, (_, i) => ({
    ...createSection({
      id: `section-${i}`,
      name: `Section ${i + 1}`,
      duration: 10 + i * 5,
    }),
  }))
}

// Factory for creating test audio buffer
export const createAudioBuffer = (
  duration: number = 1,
  sampleRate: number = 44100,
  numberOfChannels: number = 2
): AudioBuffer => {
  const length = Math.floor(duration * sampleRate)
  const audioContext = new AudioContext()
  return audioContext.createBuffer(numberOfChannels, length, sampleRate)
}

// Factory for creating test instruments
export const createInstruments = (): InstrumentType[] => ['drums', 'bass', 'chords']

// Factory for creating test project data
export const createProjectData = () => ({
  name: 'Test Project',
  sections: createSections(3),
  arrangementMode: 'sequential' as const,
  masterBPM: 120,
  masterKey: 'C Major',
  createdAt: new Date(),
  duration: 30,
})

// Factory for creating test user settings
export const createUserSettings = () => ({
  audioInput: 'default',
  audioOutput: 'default',
  theme: 'dark',
  language: 'en',
  autoSave: true,
  notifications: true,
})

// Factory for creating test export settings
export const createExportSettings = () => ({
  format: 'midi' as const,
  quality: 'high' as const,
  includeMetadata: true,
  normalizeAudio: true,
})

// Factory for creating test analytics data
export const createAnalyticsData = () => ({
  sessions: [
    {
      id: 'session-1',
      startTime: new Date(),
      duration: 300,
      recordingsCount: 2,
      exportsCount: 1,
      featuresUsed: ['recording', 'analysis', 'export'],
    },
  ],
  totalRecordings: 10,
  totalExports: 5,
  averageSessionDuration: 250,
  mostUsedFeatures: ['recording', 'analysis'],
})

// Factory for creating test error data
export const createErrorData = (message: string = 'Test error') => ({
  message,
  code: 'TEST_ERROR',
  timestamp: new Date(),
  stack: 'Test stack trace',
})

// Factory for creating test performance metrics
export const createPerformanceMetrics = () => ({
  recordingLatency: 50,
  analysisTime: 2000,
  renderTime: 16,
  memoryUsage: 1024 * 1024 * 50, // 50MB
  cpuUsage: 25,
})

// Factory for creating test MIDI data
export const createMidiData = () => ({
  header: {
    formatType: 1,
    tracksCount: 4,
    ticksPerQuarter: 480,
  },
  tracks: [
    {
      name: 'Drums',
      events: [
        { type: 'noteOn', note: 36, velocity: 100, time: 0 },
        { type: 'noteOff', note: 36, velocity: 0, time: 480 },
      ],
    },
    {
      name: 'Bass',
      events: [
        { type: 'noteOn', note: 40, velocity: 100, time: 0 },
        { type: 'noteOff', note: 40, velocity: 0, time: 960 },
      ],
    },
    {
      name: 'Chords',
      events: [
        { type: 'noteOn', note: 60, velocity: 100, time: 0 },
        { type: 'noteOn', note: 64, velocity: 100, time: 0 },
        { type: 'noteOn', note: 67, velocity: 100, time: 0 },
        { type: 'noteOff', note: 60, velocity: 0, time: 1920 },
        { type: 'noteOff', note: 64, velocity: 0, time: 1920 },
        { type: 'noteOff', note: 67, velocity: 0, time: 1920 },
      ],
    },
  ],
})

// Factory for creating test waveform data
export const createWaveformData = (length: number = 1000) => {
  return Array.from({ length }, () => Math.random() * 2 - 1)
}

// Factory for creating test frequency spectrum data
export const createFrequencySpectrum = (bins: number = 1024) => {
  return Array.from({ length: bins }, (_, i) => {
    // Create a realistic frequency spectrum with some peaks
    const frequency = (i * 22050) / bins // Nyquist frequency / 2
    if (frequency > 80 && frequency < 1000) {
      return Math.random() * 255 // Higher amplitude for voice frequencies
    }
    return Math.random() * 50 // Lower amplitude for other frequencies
  })
}

// Factory for creating test audio context state
export const createAudioContextState = () => ({
  sampleRate: 44100,
  currentTime: 0,
  state: 'running' as AudioContextState,
  destination: {},
  listener: {},
})

// Factory for creating test media stream
export const createMediaStream = () => ({
  id: 'test-stream',
  active: true,
  getTracks: () => [
    {
      id: 'audio-track-1',
      kind: 'audio',
      label: 'Test Audio',
      enabled: true,
      muted: false,
      readyState: 'live',
      stop: () => {},
      getSettings: () => ({
        deviceId: 'default',
        groupId: 'default',
        kind: 'audioinput',
        label: 'Test Audio',
        sampleRate: 44100,
      }),
    },
  ],
  getAudioTracks: () => [],
  getVideoTracks: () => [],
  addTrack: () => {},
  removeTrack: () => {},
})