import { createStore } from './index';
import { PitchPoint, BPMAnalysis, KeyAnalysis, TimeSignature, InstrumentType, Section } from '@/lib/types';

// Audio recording state
export interface RecordingState {
  isRecording: boolean;
  hasPermission: boolean;
  recordedAudio: AudioBuffer | null;
  originalAudio: AudioBuffer | null;
  isPlaying: boolean;
  trimStart: number;
  trimEnd: number;
  isTrimming: boolean;
}

// Audio analysis state
export interface AnalysisState {
  audioBuffer: AudioBuffer | null;
  pitches: PitchPoint[];
  bpm: BPMAnalysis | null;
  key: KeyAnalysis | null;
  timeSignature: TimeSignature | null;
  isAnalyzing: boolean;
  analysisError: string | null;
}

// Generated audio state
export interface GeneratedAudioState {
  isGenerating: boolean;
  generatedPlaying: boolean;
  generatedParts: {
    drums?: any;
    bass?: any;
    chords?: any;
  };
  arrangementMode: 'sequential' | 'layered';
}

// Export settings
export interface ExportSettings {
  format: 'midi' | 'wav' | 'mp3' | 'stems';
  quality: 'low' | 'medium' | 'high';
  includeMetadata: boolean;
  normalizeAudio: boolean;
}

// Combined audio store state
export interface AudioStoreState extends RecordingState, AnalysisState, GeneratedAudioState {
  selectedInstruments: InstrumentType[];
  sections: Section[];
  exportSettings: ExportSettings;
  
  // Recording actions
  setRecording: (isRecording: boolean) => void;
  setPermission: (hasPermission: boolean) => void;
  setRecordedAudio: (audio: AudioBuffer | null) => void;
  setOriginalAudio: (audio: AudioBuffer | null) => void;
  setPlaying: (isPlaying: boolean) => void;
  setTrimStart: (start: number) => void;
  setTrimEnd: (end: number) => void;
  setTrimming: (isTrimming: boolean) => void;
  
  // Analysis actions
  setAudioBuffer: (buffer: AudioBuffer | null) => void;
  setPitches: (pitches: PitchPoint[]) => void;
  setBPM: (bpm: BPMAnalysis | null) => void;
  setKey: (key: KeyAnalysis | null) => void;
  setTimeSignature: (timeSignature: TimeSignature | null) => void;
  setAnalyzing: (isAnalyzing: boolean) => void;
  setAnalysisError: (error: string | null) => void;
  
  // Generated audio actions
  setGenerating: (isGenerating: boolean) => void;
  setGeneratedPlaying: (isPlaying: boolean) => void;
  setGeneratedParts: (parts: GeneratedAudioState['generatedParts']) => void;
  setArrangementMode: (mode: 'sequential' | 'layered') => void;
  
  // Instrument and section actions
  setSelectedInstruments: (instruments: InstrumentType[]) => void;
  setSections: (sections: Section[]) => void;
  addSection: (section: Section) => void;
  updateSection: (id: string, updates: Partial<Section>) => void;
  deleteSection: (id: string) => void;
  
  // Export settings actions
  setExportSettings: (settings: Partial<ExportSettings>) => void;
  
  // Reset actions
  resetRecording: () => void;
  resetAnalysis: () => void;
  resetGenerated: () => void;
  resetAll: () => void;
}

// Initial state
const initialState: Omit<AudioStoreState, 'setRecording' | 'setPermission' | 'setRecordedAudio' | 'setOriginalAudio' | 'setPlaying' | 'setTrimStart' | 'setTrimEnd' | 'setTrimming' | 'setAudioBuffer' | 'setPitches' | 'setBPM' | 'setKey' | 'setTimeSignature' | 'setAnalyzing' | 'setAnalysisError' | 'setGenerating' | 'setGeneratedPlaying' | 'setGeneratedParts' | 'setArrangementMode' | 'setSelectedInstruments' | 'setSections' | 'addSection' | 'updateSection' | 'deleteSection' | 'setExportSettings' | 'resetRecording' | 'resetAnalysis' | 'resetGenerated' | 'resetAll'> = {
  // Recording state
  isRecording: false,
  hasPermission: false,
  recordedAudio: null,
  originalAudio: null,
  isPlaying: false,
  trimStart: 0,
  trimEnd: 0,
  isTrimming: false,
  
  // Analysis state
  audioBuffer: null,
  pitches: [],
  bpm: null,
  key: null,
  timeSignature: null,
  isAnalyzing: false,
  analysisError: null,
  
  // Generated audio state
  isGenerating: false,
  generatedPlaying: false,
  generatedParts: {},
  arrangementMode: 'sequential',
  
  // Other state
  selectedInstruments: ['drums', 'bass', 'chords'],
  sections: [],
  exportSettings: {
    format: 'midi',
    quality: 'high',
    includeMetadata: true,
    normalizeAudio: true,
  },
};

// Create the audio store
export const useAudioStore = createStore<AudioStoreState>(
  (set, get) => ({
    ...initialState,
    
    // Recording actions
    setRecording: (isRecording) => set({ isRecording }),
    setPermission: (hasPermission) => set({ hasPermission }),
    setRecordedAudio: (recordedAudio) => set({ recordedAudio }),
    setOriginalAudio: (originalAudio) => set({ originalAudio }),
    setPlaying: (isPlaying) => set({ isPlaying }),
    setTrimStart: (trimStart) => set({ trimStart }),
    setTrimEnd: (trimEnd) => set({ trimEnd }),
    setTrimming: (isTrimming) => set({ isTrimming }),
    
    // Analysis actions
    setAudioBuffer: (audioBuffer) => set({ audioBuffer }),
    setPitches: (pitches) => set({ pitches }),
    setBPM: (bpm) => set({ bpm }),
    setKey: (key) => set({ key }),
    setTimeSignature: (timeSignature) => set({ timeSignature }),
    setAnalyzing: (isAnalyzing) => set({ isAnalyzing }),
    setAnalysisError: (analysisError) => set({ analysisError }),
    
    // Generated audio actions
    setGenerating: (isGenerating) => set({ isGenerating }),
    setGeneratedPlaying: (generatedPlaying) => set({ generatedPlaying }),
    setGeneratedParts: (generatedParts) => set({ generatedParts }),
    setArrangementMode: (arrangementMode) => set({ arrangementMode }),
    
    // Instrument and section actions
    setSelectedInstruments: (selectedInstruments) => set({ selectedInstruments }),
    setSections: (sections) => set({ sections }),
    addSection: (section) => set((state) => ({ 
      sections: [...state.sections, section] 
    })),
    updateSection: (id, updates) => set((state) => ({
      sections: state.sections.map(section => 
        section.id === id ? { ...section, ...updates } : section
      )
    })),
    deleteSection: (id) => set((state) => ({
      sections: state.sections.filter(section => section.id !== id)
    })),
    
    // Export settings actions
    setExportSettings: (settings) => set((state) => ({
      exportSettings: { ...state.exportSettings, ...settings }
    })),
    
    // Reset actions
    resetRecording: () => set({
      isRecording: false,
      recordedAudio: null,
      originalAudio: null,
      isPlaying: false,
      trimStart: 0,
      trimEnd: 0,
      isTrimming: false,
    }),
    
    resetAnalysis: () => set({
      audioBuffer: null,
      pitches: [],
      bpm: null,
      key: null,
      timeSignature: null,
      isAnalyzing: false,
      analysisError: null,
    }),
    
    resetGenerated: () => set({
      isGenerating: false,
      generatedPlaying: false,
      generatedParts: {},
    }),
    
    resetAll: () => set(initialState),
  }),
  {
    name: 'voxforge-audio-store',
    persist: {
      name: 'voxforge-audio-store',
      partialize: (state) => ({
        selectedInstruments: state.selectedInstruments,
        exportSettings: state.exportSettings,
        arrangementMode: state.arrangementMode,
      }),
    },
  }
);

// Selectors for optimized re-renders
export const useRecordingState = () => {
  const store = useAudioStore();
  return {
    isRecording: store.isRecording,
    hasPermission: store.hasPermission,
    recordedAudio: store.recordedAudio,
    originalAudio: store.originalAudio,
    isPlaying: store.isPlaying,
    trimStart: store.trimStart,
    trimEnd: store.trimEnd,
    isTrimming: store.isTrimming,
  };
};

export const useAnalysisState = () => {
  const store = useAudioStore();
  return {
    audioBuffer: store.audioBuffer,
    pitches: store.pitches,
    bpm: store.bpm,
    key: store.key,
    timeSignature: store.timeSignature,
    isAnalyzing: store.isAnalyzing,
    analysisError: store.analysisError,
  };
};

export const useGeneratedAudioState = () => {
  const store = useAudioStore();
  return {
    isGenerating: store.isGenerating,
    generatedPlaying: store.generatedPlaying,
    generatedParts: store.generatedParts,
    arrangementMode: store.arrangementMode,
  };
};

export const useAudioActions = () => {
  const store = useAudioStore();
  return {
    setRecording: store.setRecording,
    setPermission: store.setPermission,
    setRecordedAudio: store.setRecordedAudio,
    setOriginalAudio: store.setOriginalAudio,
    setPlaying: store.setPlaying,
    setTrimStart: store.setTrimStart,
    setTrimEnd: store.setTrimEnd,
    setTrimming: store.setTrimming,
    setAudioBuffer: store.setAudioBuffer,
    setPitches: store.setPitches,
    setBPM: store.setBPM,
    setKey: store.setKey,
    setTimeSignature: store.setTimeSignature,
    setAnalyzing: store.setAnalyzing,
    setAnalysisError: store.setAnalysisError,
    setGenerating: store.setGenerating,
    setGeneratedPlaying: store.setGeneratedPlaying,
    setGeneratedParts: store.setGeneratedParts,
    setArrangementMode: store.setArrangementMode,
    setSelectedInstruments: store.setSelectedInstruments,
    setSections: store.setSections,
    addSection: store.addSection,
    updateSection: store.updateSection,
    deleteSection: store.deleteSection,
    setExportSettings: store.setExportSettings,
    resetRecording: store.resetRecording,
    resetAnalysis: store.resetAnalysis,
    resetGenerated: store.resetGenerated,
    resetAll: store.resetAll,
  };
};