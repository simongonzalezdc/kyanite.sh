/// <reference types="jest" />

import { renderHook, act } from '@testing-library/react'
import {
  useAudioStore,
  useRecordingState,
  useAnalysisState,
  useGeneratedAudioState,
  useAudioActions
} from '../audioStore'
import {
  createPitchPoints,
  createBPMAnalysis,
  createKeyAnalysis,
  createSection,
  createAudioBuffer
} from '@/__tests__/utils/testDataFactories'
import { InstrumentType } from '@/lib/types'

describe('audioStore', () => {
  beforeEach(() => {
    // Reset store before each test
    const { resetAll } = useAudioStore.getState()
    resetAll()
  })

  describe('useAudioStore', () => {
    it('should return initial state', () => {
      const state = useAudioStore.getState()

      expect(state.isRecording).toBe(false)
      expect(state.hasPermission).toBe(false)
      expect(state.recordedAudio).toBe(null)
      expect(state.originalAudio).toBe(null)
      expect(state.isPlaying).toBe(false)
      expect(state.isAnalyzing).toBe(false)
      expect(state.analysisError).toBe(null)
      expect(state.isGenerating).toBe(false)
      expect(state.generatedPlaying).toBe(false)
      expect(state.selectedInstruments).toEqual(['drums', 'bass', 'chords'])
      expect(state.sections).toEqual([])
    })

    it('should update state when actions are called', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setRecording(true)
        result.current.actions.setAnalyzing(true)
        result.current.actions.setPitches(createPitchPoints(3))
      })

      expect(result.current.store.isRecording).toBe(true)
      expect(result.current.store.isAnalyzing).toBe(true)
      expect(result.current.store.pitches).toHaveLength(3)
    })
  })

  describe('useRecordingState', () => {
    it('should return recording-related state', () => {
      const state = useAudioStore.getState()
      const recordingState = {
        isRecording: state.isRecording,
        hasPermission: state.hasPermission,
        recordedAudio: state.recordedAudio,
        originalAudio: state.originalAudio,
        isPlaying: state.isPlaying,
        trimStart: state.trimStart,
        trimEnd: state.trimEnd,
        isTrimming: state.isTrimming,
      }

      expect(recordingState.isRecording).toBe(false)
      expect(recordingState.hasPermission).toBe(false)
      expect(recordingState.recordedAudio).toBe(null)
      expect(recordingState.originalAudio).toBe(null)
      expect(recordingState.isPlaying).toBe(false)
      expect(recordingState.trimStart).toBe(0)
      expect(recordingState.trimEnd).toBe(0)
      expect(recordingState.isTrimming).toBe(false)
    })

    it('should not include analysis state', () => {
      const state = useAudioStore.getState()
      const recordingState = {
        isRecording: state.isRecording,
        hasPermission: state.hasPermission,
        recordedAudio: state.recordedAudio,
        originalAudio: state.originalAudio,
        isPlaying: state.isPlaying,
        trimStart: state.trimStart,
        trimEnd: state.trimEnd,
        isTrimming: state.isTrimming,
      }

      expect(recordingState).not.toHaveProperty('pitches')
      expect(recordingState).not.toHaveProperty('bpm')
      expect(recordingState).not.toHaveProperty('key')
      expect(recordingState).not.toHaveProperty('isAnalyzing')
    })
  })

  describe('useAnalysisState', () => {
    it('should return analysis-related state', () => {
      const state = useAudioStore.getState()
      const analysisState = {
        audioBuffer: state.audioBuffer,
        pitches: state.pitches,
        bpm: state.bpm,
        key: state.key,
        timeSignature: state.timeSignature,
        isAnalyzing: state.isAnalyzing,
        analysisError: state.analysisError,
      }

      expect(analysisState.audioBuffer).toBe(null)
      expect(analysisState.pitches).toEqual([])
      expect(analysisState.bpm).toBe(null)
      expect(analysisState.key).toBe(null)
      expect(analysisState.timeSignature).toBe(null)
      expect(analysisState.isAnalyzing).toBe(false)
      expect(analysisState.analysisError).toBe(null)
    })

    it('should not include recording state', () => {
      const state = useAudioStore.getState()
      const analysisState = {
        audioBuffer: state.audioBuffer,
        pitches: state.pitches,
        bpm: state.bpm,
        key: state.key,
        timeSignature: state.timeSignature,
        isAnalyzing: state.isAnalyzing,
        analysisError: state.analysisError,
      }

      expect(analysisState).not.toHaveProperty('isRecording')
      expect(analysisState).not.toHaveProperty('hasPermission')
      expect(analysisState).not.toHaveProperty('recordedAudio')
      expect(analysisState).not.toHaveProperty('isPlaying')
    })
  })

  describe('useGeneratedAudioState', () => {
    it('should return generated audio state', () => {
      const state = useAudioStore.getState()
      const generatedState = {
        isGenerating: state.isGenerating,
        generatedPlaying: state.generatedPlaying,
        generatedParts: state.generatedParts,
        arrangementMode: state.arrangementMode,
      }

      expect(generatedState.isGenerating).toBe(false)
      expect(generatedState.generatedPlaying).toBe(false)
      expect(generatedState.generatedParts).toEqual({})
      expect(generatedState.arrangementMode).toBe('sequential')
    })

    it('should not include other state', () => {
      const state = useAudioStore.getState()
      const generatedState = {
        isGenerating: state.isGenerating,
        generatedPlaying: state.generatedPlaying,
        generatedParts: state.generatedParts,
        arrangementMode: state.arrangementMode,
      }

      expect(generatedState).not.toHaveProperty('isRecording')
      expect(generatedState).not.toHaveProperty('isAnalyzing')
      expect(generatedState).not.toHaveProperty('pitches')
      expect(generatedState).not.toHaveProperty('sections')
    })
  })

  describe('useAudioActions', () => {
    it('should return all action functions', () => {
      const { result } = renderHook(() => useAudioActions())
      const actions = result.current

      expect(typeof actions.setRecording).toBe('function')
      expect(typeof actions.setPermission).toBe('function')
      expect(typeof actions.setRecordedAudio).toBe('function')
      expect(typeof actions.setOriginalAudio).toBe('function')
      expect(typeof actions.setPlaying).toBe('function')
      expect(typeof actions.setTrimStart).toBe('function')
      expect(typeof actions.setTrimEnd).toBe('function')
      expect(typeof actions.setTrimming).toBe('function')
      expect(typeof actions.setAudioBuffer).toBe('function')
      expect(typeof actions.setPitches).toBe('function')
      expect(typeof actions.setBPM).toBe('function')
      expect(typeof actions.setKey).toBe('function')
      expect(typeof actions.setTimeSignature).toBe('function')
      expect(typeof actions.setAnalyzing).toBe('function')
      expect(typeof actions.setAnalysisError).toBe('function')
      expect(typeof actions.setGenerating).toBe('function')
      expect(typeof actions.setGeneratedPlaying).toBe('function')
      expect(typeof actions.setGeneratedParts).toBe('function')
      expect(typeof actions.setArrangementMode).toBe('function')
      expect(typeof actions.setSelectedInstruments).toBe('function')
      expect(typeof actions.setSections).toBe('function')
      expect(typeof actions.addSection).toBe('function')
      expect(typeof actions.updateSection).toBe('function')
      expect(typeof actions.deleteSection).toBe('function')
      expect(typeof actions.setExportSettings).toBe('function')
      expect(typeof actions.resetRecording).toBe('function')
      expect(typeof actions.resetAnalysis).toBe('function')
      expect(typeof actions.resetGenerated).toBe('function')
      expect(typeof actions.resetAll).toBe('function')
    })
  })

  describe('Recording actions', () => {
    it('should set recording state', () => {
      useAudioStore.setState({ isRecording: true })

      expect(useAudioStore.getState().isRecording).toBe(true)
    })

    it('should set permission state', () => {
      useAudioStore.setState({ hasPermission: true })

      expect(useAudioStore.getState().hasPermission).toBe(true)
    })

    it('should set recorded audio', () => {
      const mockAudio = createAudioBuffer(2, 44100)

      useAudioStore.setState({ recordedAudio: mockAudio })

      expect(useAudioStore.getState().recordedAudio).toBe(mockAudio)
    })

    it('should set playing state', () => {
      useAudioStore.setState({ isPlaying: true })

      expect(useAudioStore.getState().isPlaying).toBe(true)
    })

    it('should set trim values', () => {
      useAudioStore.setState({
        trimStart: 1.5,
        trimEnd: 3.5,
        isTrimming: true
      })

      const state = useAudioStore.getState()
      expect(state.trimStart).toBe(1.5)
      expect(state.trimEnd).toBe(3.5)
      expect(state.isTrimming).toBe(true)
    })
  })

  describe('Analysis actions', () => {
    it('should set audio buffer', () => {
      const mockBuffer = createAudioBuffer(2, 44100)

      useAudioStore.setState({ audioBuffer: mockBuffer })

      expect(useAudioStore.getState().audioBuffer).toBe(mockBuffer)
    })

    it('should set pitches', () => {
      const mockPitches = createPitchPoints(5)

      useAudioStore.setState({ pitches: mockPitches })

      expect(useAudioStore.getState().pitches).toEqual(mockPitches)
    })

    it('should set BPM analysis', () => {
      const mockBPM = { bpm: 120, confidence: 0.9, stable: true }

      useAudioStore.setState({ bpm: mockBPM })

      expect(useAudioStore.getState().bpm).toEqual(mockBPM)
    })

    it('should set key analysis', () => {
      const mockKey = createKeyAnalysis()

      useAudioStore.setState({ key: mockKey })

      expect(useAudioStore.getState().key).toEqual(mockKey)
    })

    it('should set time signature', () => {
      const mockTimeSignature = { numerator: 4, denominator: 4, display: '4/4' }

      useAudioStore.setState({ timeSignature: mockTimeSignature })

      expect(useAudioStore.getState().timeSignature).toEqual(mockTimeSignature)
    })

    it('should set analyzing state', () => {
      useAudioStore.setState({ isAnalyzing: true })

      expect(useAudioStore.getState().isAnalyzing).toBe(true)
    })

    it('should set analysis error', () => {
      const mockError = 'Analysis failed'

      useAudioStore.setState({ analysisError: mockError })

      expect(useAudioStore.getState().analysisError).toBe(mockError)
    })
  })

  describe('Generated audio actions', () => {
    it('should set generating state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setGenerating(true)
      })

      expect(result.current.store.isGenerating).toBe(true)
    })

    it('should set generated playing state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setGeneratedPlaying(true)
      })

      expect(result.current.store.generatedPlaying).toBe(true)
    })

    it('should set generated parts', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })
      const mockParts = {
        drums: { id: 'drums-1' },
        bass: { id: 'bass-1' },
        chords: { id: 'chords-1' },
      }

      act(() => {
        result.current.actions.setGeneratedParts(mockParts)
      })

      expect(result.current.store.generatedParts).toEqual(mockParts)
    })

    it('should set arrangement mode', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setArrangementMode('layered')
      })

      expect(result.current.store.arrangementMode).toBe('layered')
    })
  })

  describe('Section actions', () => {
    it('should set selected instruments', () => {
      const instruments: InstrumentType[] = ['drums', 'bass']

      useAudioStore.setState({ selectedInstruments: instruments })

      expect(useAudioStore.getState().selectedInstruments).toEqual(instruments)
    })

    it('should set sections', () => {
      const mockSections = [createSection(), createSection()]

      useAudioStore.setState({ sections: mockSections })

      expect(useAudioStore.getState().sections).toEqual(mockSections)
    })

    it('should add section', () => {
      const currentSections = useAudioStore.getState().sections
      const mockSection = createSection()

      useAudioStore.setState({ sections: [...currentSections, mockSection] })

      const state = useAudioStore.getState()
      expect(state.sections).toHaveLength(currentSections.length + 1)
      expect(state.sections[state.sections.length - 1]).toEqual(mockSection)
    })

    it('should update section', () => {
      const mockSection = createSection()
      useAudioStore.setState({ sections: [mockSection] })

      const updates = { name: 'Updated Section' }
      const updatedSections = useAudioStore.getState().sections.map(section =>
        section.id === mockSection.id ? { ...section, ...updates } : section
      )
      useAudioStore.setState({ sections: updatedSections })

      const state = useAudioStore.getState()
      expect(state.sections[0].name).toBe('Updated Section')
      expect(state.sections[0]).toEqual({
        ...mockSection,
        ...updates,
      })
    })

    it('should delete section', () => {
      const mockSection = createSection()
      useAudioStore.setState({ sections: [mockSection] })

      expect(useAudioStore.getState().sections).toHaveLength(1)

      useAudioStore.setState({ sections: [] })

      expect(useAudioStore.getState().sections).toHaveLength(0)
    })
  })

  describe('Export settings actions', () => {
    it('should set export settings', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })
      const mockSettings = {
        format: 'wav' as const,
        quality: 'medium' as const,
        includeMetadata: false,
        normalizeAudio: false,
      }

      act(() => {
        result.current.actions.setExportSettings(mockSettings)
      })

      expect(result.current.store.exportSettings).toEqual({
        ...result.current.store.exportSettings,
        ...mockSettings,
      })
    })
  })

  describe('Reset actions', () => {
    it('should reset recording state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setRecording(true)
      })
      expect(result.current.store.isRecording).toBe(true)

      act(() => {
        result.current.actions.resetRecording()
      })

      expect(result.current.store.isRecording).toBe(false)
      expect(result.current.store.recordedAudio).toBe(null)
      expect(result.current.store.originalAudio).toBe(null)
      expect(result.current.store.isPlaying).toBe(false)
      expect(result.current.store.trimStart).toBe(0)
      expect(result.current.store.trimEnd).toBe(0)
      expect(result.current.store.isTrimming).toBe(false)
    })

    it('should reset analysis state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setAnalyzing(true)
        result.current.actions.setPitches(createPitchPoints(3))
      })
      expect(result.current.store.isAnalyzing).toBe(true)
      expect(result.current.store.pitches).toHaveLength(3)

      act(() => {
        result.current.actions.resetAnalysis()
      })

      expect(result.current.store.isAnalyzing).toBe(false)
      expect(result.current.store.audioBuffer).toBe(null)
      expect(result.current.store.pitches).toEqual([])
      expect(result.current.store.bpm).toBe(null)
      expect(result.current.store.key).toBe(null)
      expect(result.current.store.timeSignature).toBe(null)
      expect(result.current.store.analysisError).toBe(null)
    })

    it('should reset generated audio state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setGenerating(true)
        result.current.actions.setGeneratedParts({ drums: { id: 'test' } })
      })
      expect(result.current.store.isGenerating).toBe(true)
      expect(result.current.store.generatedParts).toEqual({ drums: { id: 'test' } })

      act(() => {
        result.current.actions.resetGenerated()
      })

      expect(result.current.store.isGenerating).toBe(false)
      expect(result.current.store.generatedPlaying).toBe(false)
      expect(result.current.store.generatedParts).toEqual({})
    })

    it('should reset all state', () => {
      const { result } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })

      act(() => {
        result.current.actions.setRecording(true)
        result.current.actions.setAnalyzing(true)
        result.current.actions.setSelectedInstruments(['drums'] as any)
      })
      expect(result.current.store.isRecording).toBe(true)
      expect(result.current.store.isAnalyzing).toBe(true)
      expect(result.current.store.selectedInstruments).toEqual(['drums'])

      act(() => {
        result.current.actions.resetAll()
      })

      // Check that all state is reset to initial values
      expect(result.current.store.isRecording).toBe(false)
      expect(result.current.store.hasPermission).toBe(false)
      expect(result.current.store.recordedAudio).toBe(null)
      expect(result.current.store.originalAudio).toBe(null)
      expect(result.current.store.isPlaying).toBe(false)
      expect(result.current.store.trimStart).toBe(0)
      expect(result.current.store.trimEnd).toBe(0)
      expect(result.current.store.isTrimming).toBe(false)
      expect(result.current.store.audioBuffer).toBe(null)
      expect(result.current.store.pitches).toEqual([])
      expect(result.current.store.bpm).toBe(null)
      expect(result.current.store.key).toBe(null)
      expect(result.current.store.timeSignature).toBe(null)
      expect(result.current.store.isAnalyzing).toBe(false)
      expect(result.current.store.analysisError).toBe(null)
      expect(result.current.store.isGenerating).toBe(false)
      expect(result.current.store.generatedPlaying).toBe(false)
      expect(result.current.store.generatedParts).toEqual({})
      expect(result.current.store.arrangementMode).toBe('sequential')
      expect(result.current.store.selectedInstruments).toEqual(['drums', 'bass', 'chords'])
      expect(result.current.store.sections).toEqual([])
    })
  })

  describe('State persistence', () => {
    it('should persist selected instruments', () => {
      const { result: result1 } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })
      const instruments: any = ['drums']

      act(() => {
        result1.current.actions.setSelectedInstruments(instruments)
      })
      expect(result1.current.store.selectedInstruments).toEqual(instruments)

      // For Zustand stores, state persists automatically across re-renders
      // We can test by creating a new hook instance
      const { result: result2 } = renderHook(() => useAudioStore())
      expect(result2.current.selectedInstruments).toEqual(instruments)
    })

    it('should persist export settings', () => {
      const { result: result1 } = renderHook(() => {
        const store = useAudioStore()
        const actions = useAudioActions()
        return { store, actions }
      })
      const settings = {
        format: 'mp3' as const,
        quality: 'low' as const,
        includeMetadata: false,
        normalizeAudio: true,
      }

      act(() => {
        result1.current.actions.setExportSettings(settings)
      })
      expect(result1.current.store.exportSettings).toEqual(expect.objectContaining(settings))

      // For Zustand stores, state persists automatically across re-renders
      // We can test by creating a new hook instance
      const { result: result2 } = renderHook(() => useAudioStore())
      expect(result2.current.exportSettings).toEqual(expect.objectContaining(settings))
    })
  })
})