import { 
  useAudioStore, 
  useRecordingState, 
  useAnalysisState, 
  useGeneratedAudioState, 
  useAudioActions 
} from '../audioStore'
import { renderHook, act } from '@testing-library/react'
import { 
  createPitchPoints, 
  createBPMAnalysis, 
  createKeyAnalysis, 
  createSection 
} from '@/__tests__/utils/testDataFactories'

describe('audioStore', () => {
  beforeEach(() => {
    // Reset store before each test
    const { resetAll } = useAudioStore.getState()
    resetAll()
  })

  describe('useAudioStore', () => {
    it('should return initial state', () => {
      const { result } = renderHook(() => useAudioStore())
      
      expect(result.current.isRecording).toBe(false)
      expect(result.current.hasPermission).toBe(false)
      expect(result.current.recordedAudio).toBe(null)
      expect(result.current.originalAudio).toBe(null)
      expect(result.current.isPlaying).toBe(false)
      expect(result.current.isAnalyzing).toBe(false)
      expect(result.current.analysisError).toBe(null)
      expect(result.current.isGenerating).toBe(false)
      expect(result.current.generatedPlaying).toBe(false)
      expect(result.current.selectedInstruments).toEqual(['drums', 'bass', 'chords'])
      expect(result.current.sections).toEqual([])
    })

    it('should update state when actions are called', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setRecording, setAnalyzing, setPitches } = useAudioActions()
      
      act(() => {
        setRecording(true)
        setAnalyzing(true)
        setPitches(createPitchPoints(3))
      })
      
      expect(result.current.isRecording).toBe(true)
      expect(result.current.isAnalyzing).toBe(true)
      expect(result.current.pitches).toHaveLength(3)
    })
  })

  describe('useRecordingState', () => {
    it('should return recording-related state', () => {
      const { result } = renderHook(() => useRecordingState())
      
      expect(result.current.isRecording).toBe(false)
      expect(result.current.hasPermission).toBe(false)
      expect(result.current.recordedAudio).toBe(null)
      expect(result.current.originalAudio).toBe(null)
      expect(result.current.isPlaying).toBe(false)
      expect(result.current.trimStart).toBe(0)
      expect(result.current.trimEnd).toBe(0)
      expect(result.current.isTrimming).toBe(false)
    })

    it('should not include analysis state', () => {
      const { result } = renderHook(() => useRecordingState())
      
      expect(result.current).not.toHaveProperty('pitches')
      expect(result.current).not.toHaveProperty('bpm')
      expect(result.current).not.toHaveProperty('key')
      expect(result.current).not.toHaveProperty('isAnalyzing')
    })
  })

  describe('useAnalysisState', () => {
    it('should return analysis-related state', () => {
      const { result } = renderHook(() => useAnalysisState())
      
      expect(result.current.audioBuffer).toBe(null)
      expect(result.current.pitches).toEqual([])
      expect(result.current.bpm).toBe(null)
      expect(result.current.key).toBe(null)
      expect(result.current.timeSignature).toBe(null)
      expect(result.current.isAnalyzing).toBe(false)
      expect(result.current.analysisError).toBe(null)
    })

    it('should not include recording state', () => {
      const { result } = renderHook(() => useAnalysisState())
      
      expect(result.current).not.toHaveProperty('isRecording')
      expect(result.current).not.toHaveProperty('hasPermission')
      expect(result.current).not.toHaveProperty('recordedAudio')
      expect(result.current).not.toHaveProperty('isPlaying')
    })
  })

  describe('useGeneratedAudioState', () => {
    it('should return generated audio state', () => {
      const { result } = renderHook(() => useGeneratedAudioState())
      
      expect(result.current.isGenerating).toBe(false)
      expect(result.current.generatedPlaying).toBe(false)
      expect(result.current.generatedParts).toEqual({})
      expect(result.current.arrangementMode).toBe('sequential')
    })

    it('should not include other state', () => {
      const { result } = renderHook(() => useGeneratedAudioState())
      
      expect(result.current).not.toHaveProperty('isRecording')
      expect(result.current).not.toHaveProperty('isAnalyzing')
      expect(result.current).not.toHaveProperty('pitches')
      expect(result.current).not.toHaveProperty('sections')
    })
  })

  describe('useAudioActions', () => {
    it('should return all action functions', () => {
      const { result } = renderHook(() => useAudioActions())
      
      expect(typeof result.current.setRecording).toBe('function')
      expect(typeof result.current.setPermission).toBe('function')
      expect(typeof result.current.setRecordedAudio).toBe('function')
      expect(typeof result.current.setOriginalAudio).toBe('function')
      expect(typeof result.current.setPlaying).toBe('function')
      expect(typeof result.current.setTrimStart).toBe('function')
      expect(typeof result.current.setTrimEnd).toBe('function')
      expect(typeof result.current.setTrimming).toBe('function')
      expect(typeof result.current.setAudioBuffer).toBe('function')
      expect(typeof result.current.setPitches).toBe('function')
      expect(typeof result.current.setBPM).toBe('function')
      expect(typeof result.current.setKey).toBe('function')
      expect(typeof result.current.setTimeSignature).toBe('function')
      expect(typeof result.current.setAnalyzing).toBe('function')
      expect(typeof result.current.setAnalysisError).toBe('function')
      expect(typeof result.current.setGenerating).toBe('function')
      expect(typeof result.current.setGeneratedPlaying).toBe('function')
      expect(typeof result.current.setGeneratedParts).toBe('function')
      expect(typeof result.current.setArrangementMode).toBe('function')
      expect(typeof result.current.setSelectedInstruments).toBe('function')
      expect(typeof result.current.setSections).toBe('function')
      expect(typeof result.current.addSection).toBe('function')
      expect(typeof result.current.updateSection).toBe('function')
      expect(typeof result.current.deleteSection).toBe('function')
      expect(typeof result.current.setExportSettings).toBe('function')
      expect(typeof result.current.resetRecording).toBe('function')
      expect(typeof result.current.resetAnalysis).toBe('function')
      expect(typeof result.current.resetGenerated).toBe('function')
      expect(typeof result.current.resetAll).toBe('function')
    })
  })

  describe('Recording actions', () => {
    it('should set recording state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setRecording } = useAudioActions()
      
      act(() => {
        setRecording(true)
      })
      
      expect(result.current.isRecording).toBe(true)
    })

    it('should set permission state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setPermission } = useAudioActions()
      
      act(() => {
        setPermission(true)
      })
      
      expect(result.current.hasPermission).toBe(true)
    })

    it('should set recorded audio', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setRecordedAudio } = useAudioActions()
      const mockAudio = createAudioBuffer(2, 44100)
      
      act(() => {
        setRecordedAudio(mockAudio)
      })
      
      expect(result.current.recordedAudio).toBe(mockAudio)
    })

    it('should set playing state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setPlaying } = useAudioActions()
      
      act(() => {
        setPlaying(true)
      })
      
      expect(result.current.isPlaying).toBe(true)
    })

    it('should set trim values', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setTrimStart, setTrimEnd, setTrimming } = useAudioActions()
      
      act(() => {
        setTrimStart(1.5)
        setTrimEnd(3.5)
        setTrimming(true)
      })
      
      expect(result.current.trimStart).toBe(1.5)
      expect(result.current.trimEnd).toBe(3.5)
      expect(result.current.isTrimming).toBe(true)
    })
  })

  describe('Analysis actions', () => {
    it('should set audio buffer', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setAudioBuffer } = useAudioActions()
      const mockBuffer = createAudioBuffer(2, 44100)
      
      act(() => {
        setAudioBuffer(mockBuffer)
      })
      
      expect(result.current.audioBuffer).toBe(mockBuffer)
    })

    it('should set pitches', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setPitches } = useAudioActions()
      const mockPitches = createPitchPoints(5)
      
      act(() => {
        setPitches(mockPitches)
      })
      
      expect(result.current.pitches).toEqual(mockPitches)
    })

    it('should set BPM analysis', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setBPM } = useAudioActions()
      const mockBPM = createBPMAnalysis(120)
      
      act(() => {
        setBPM(mockBPM)
      })
      
      expect(result.current.bpm).toEqual(mockBPM)
    })

    it('should set key analysis', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setKey } = useAudioActions()
      const mockKey = createKeyAnalysis('C Major')
      
      act(() => {
        setKey(mockKey)
      })
      
      expect(result.current.key).toEqual(mockKey)
    })

    it('should set time signature', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setTimeSignature } = useAudioActions()
      const mockTimeSignature = { numerator: 4, denominator: 4, display: '4/4' }
      
      act(() => {
        setTimeSignature(mockTimeSignature)
      })
      
      expect(result.current.timeSignature).toEqual(mockTimeSignature)
    })

    it('should set analyzing state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setAnalyzing } = useAudioActions()
      
      act(() => {
        setAnalyzing(true)
      })
      
      expect(result.current.isAnalyzing).toBe(true)
    })

    it('should set analysis error', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setAnalysisError } = useAudioActions()
      const mockError = 'Analysis failed'
      
      act(() => {
        setAnalysisError(mockError)
      })
      
      expect(result.current.analysisError).toBe(mockError)
    })
  })

  describe('Generated audio actions', () => {
    it('should set generating state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setGenerating } = useAudioActions()
      
      act(() => {
        setGenerating(true)
      })
      
      expect(result.current.isGenerating).toBe(true)
    })

    it('should set generated playing state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setGeneratedPlaying } = useAudioActions()
      
      act(() => {
        setGeneratedPlaying(true)
      })
      
      expect(result.current.generatedPlaying).toBe(true)
    })

    it('should set generated parts', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setGeneratedParts } = useAudioActions()
      const mockParts = {
        drums: { id: 'drums-1' },
        bass: { id: 'bass-1' },
        chords: { id: 'chords-1' },
      }
      
      act(() => {
        setGeneratedParts(mockParts)
      })
      
      expect(result.current.generatedParts).toEqual(mockParts)
    })

    it('should set arrangement mode', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setArrangementMode } = useAudioActions()
      
      act(() => {
        setArrangementMode('layered')
      })
      
      expect(result.current.arrangementMode).toBe('layered')
    })
  })

  describe('Section actions', () => {
    it('should set selected instruments', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setSelectedInstruments } = useAudioActions()
      const instruments = ['drums', 'bass']
      
      act(() => {
        setSelectedInstruments(instruments)
      })
      
      expect(result.current.selectedInstruments).toEqual(instruments)
    })

    it('should set sections', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setSections } = useAudioActions()
      const mockSections = [createSection(), createSection()]
      
      act(() => {
        setSections(mockSections)
      })
      
      expect(result.current.sections).toEqual(mockSections)
    })

    it('should add section', () => {
      const { result } = renderHook(() => useAudioStore())
      const { addSection } = useAudioActions()
      const mockSection = createSection()
      
      act(() => {
        addSection(mockSection)
      })
      
      expect(result.current.sections).toHaveLength(1)
      expect(result.current.sections[0]).toEqual(mockSection)
    })

    it('should update section', () => {
      const { result } = renderHook(() => useAudioStore())
      const { addSection, updateSection } = useAudioActions()
      const mockSection = createSection()
      
      act(() => {
        addSection(mockSection)
      })
      
      const updates = { name: 'Updated Section' }
      act(() => {
        updateSection(mockSection.id, updates)
      })
      
      expect(result.current.sections[0].name).toBe('Updated Section')
      expect(result.current.sections[0]).toEqual({
        ...mockSection,
        ...updates,
      })
    })

    it('should delete section', () => {
      const { result } = renderHook(() => useAudioStore())
      const { addSection, deleteSection } = useAudioActions()
      const mockSection = createSection()
      
      act(() => {
        addSection(mockSection)
      })
      
      expect(result.current.sections).toHaveLength(1)
      
      act(() => {
        deleteSection(mockSection.id)
      })
      
      expect(result.current.sections).toHaveLength(0)
    })
  })

  describe('Export settings actions', () => {
    it('should set export settings', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setExportSettings } = useAudioActions()
      const mockSettings = {
        format: 'wav' as const,
        quality: 'medium' as const,
        includeMetadata: false,
        normalizeAudio: false,
      }
      
      act(() => {
        setExportSettings(mockSettings)
      })
      
      expect(result.current.exportSettings).toEqual({
        ...result.current.exportSettings,
        ...mockSettings,
      })
    })
  })

  describe('Reset actions', () => {
    it('should reset recording state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setRecording, resetRecording } = useAudioActions()
      
      act(() => {
        setRecording(true)
      })
      expect(result.current.isRecording).toBe(true)
      
      act(() => {
        resetRecording()
      })
      
      expect(result.current.isRecording).toBe(false)
      expect(result.current.recordedAudio).toBe(null)
      expect(result.current.originalAudio).toBe(null)
      expect(result.current.isPlaying).toBe(false)
      expect(result.current.trimStart).toBe(0)
      expect(result.current.trimEnd).toBe(0)
      expect(result.current.isTrimming).toBe(false)
    })

    it('should reset analysis state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setAnalyzing, setPitches, resetAnalysis } = useAudioActions()
      
      act(() => {
        setAnalyzing(true)
        setPitches(createPitchPoints(3))
      })
      expect(result.current.isAnalyzing).toBe(true)
      expect(result.current.pitches).toHaveLength(3)
      
      act(() => {
        resetAnalysis()
      })
      
      expect(result.current.isAnalyzing).toBe(false)
      expect(result.current.audioBuffer).toBe(null)
      expect(result.current.pitches).toEqual([])
      expect(result.current.bpm).toBe(null)
      expect(result.current.key).toBe(null)
      expect(result.current.timeSignature).toBe(null)
      expect(result.current.analysisError).toBe(null)
    })

    it('should reset generated audio state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setGenerating, setGeneratedParts, resetGenerated } = useAudioActions()
      
      act(() => {
        setGenerating(true)
        setGeneratedParts({ drums: { id: 'test' } })
      })
      expect(result.current.isGenerating).toBe(true)
      expect(result.current.generatedParts).toEqual({ drums: { id: 'test' } })
      
      act(() => {
        resetGenerated()
      })
      
      expect(result.current.isGenerating).toBe(false)
      expect(result.current.generatedPlaying).toBe(false)
      expect(result.current.generatedParts).toEqual({})
    })

    it('should reset all state', () => {
      const { result } = renderHook(() => useAudioStore())
      const { setRecording, setAnalyzing, setSelectedInstruments, resetAll } = useAudioActions()
      
      act(() => {
        setRecording(true)
        setAnalyzing(true)
        setSelectedInstruments(['drums'])
      })
      expect(result.current.isRecording).toBe(true)
      expect(result.current.isAnalyzing).toBe(true)
      expect(result.current.selectedInstruments).toEqual(['drums'])
      
      act(() => {
        resetAll()
      })
      
      // Check that all state is reset to initial values
      expect(result.current.isRecording).toBe(false)
      expect(result.current.hasPermission).toBe(false)
      expect(result.current.recordedAudio).toBe(null)
      expect(result.current.originalAudio).toBe(null)
      expect(result.current.isPlaying).toBe(false)
      expect(result.current.trimStart).toBe(0)
      expect(result.current.trimEnd).toBe(0)
      expect(result.current.isTrimming).toBe(false)
      expect(result.current.audioBuffer).toBe(null)
      expect(result.current.pitches).toEqual([])
      expect(result.current.bpm).toBe(null)
      expect(result.current.key).toBe(null)
      expect(result.current.timeSignature).toBe(null)
      expect(result.current.isAnalyzing).toBe(false)
      expect(result.current.analysisError).toBe(null)
      expect(result.current.isGenerating).toBe(false)
      expect(result.current.generatedPlaying).toBe(false)
      expect(result.current.generatedParts).toEqual({})
      expect(result.current.arrangementMode).toBe('sequential')
      expect(result.current.selectedInstruments).toEqual(['drums', 'bass', 'chords'])
      expect(result.current.sections).toEqual([])
    })
  })

  describe('State persistence', () => {
    it('should persist selected instruments', () => {
      const { result: result1 } = renderHook(() => useAudioStore())
      const { setSelectedInstruments } = useAudioActions()
      const instruments = ['drums']
      
      act(() => {
        setSelectedInstruments(instruments)
      })
      expect(result1.current.selectedInstruments).toEqual(instruments)
      
      // Unmount and remount to test persistence
      result1.unmount()
      
      const { result: result2 } = renderHook(() => useAudioStore())
      expect(result2.current.selectedInstruments).toEqual(instruments)
    })

    it('should persist export settings', () => {
      const { result: result1 } = renderHook(() => useAudioStore())
      const { setExportSettings } = useAudioActions()
      const settings = {
        format: 'mp3' as const,
        quality: 'low' as const,
        includeMetadata: false,
        normalizeAudio: true,
      }
      
      act(() => {
        setExportSettings(settings)
      })
      expect(result1.current.exportSettings).toEqual(expect.objectContaining(settings))
      
      // Unmount and remount to test persistence
      result1.unmount()
      
      const { result: result2 } = renderHook(() => useAudioStore())
      expect(result2.current.exportSettings).toEqual(expect.objectContaining(settings))
    })
  })
})