import { createStore } from './index';

// Project metadata
export interface ProjectMetadata {
  id: string;
  name: string;
  description?: string;
  createdAt: Date;
  updatedAt: Date;
  lastOpenedAt: Date;
  tags: string[];
  isArchived: boolean;
  isFavorite: boolean;
  thumbnail?: string;
  duration: number; // in seconds
  tempo?: number;
  key?: string;
  timeSignature?: string;
  genre?: string;
  mood?: string;
}

// Section definition
export interface Section {
  id: string;
  name: string;
  type: 'intro' | 'verse' | 'chorus' | 'bridge' | 'outro' | 'custom';
  startTime: number; // in seconds
  endTime: number; // in seconds
  color: string;
  isMuted: boolean;
  isSolo: boolean;
  volume: number; // 0-1
  pan: number; // -1 to 1
  effects: string[];
  notes?: string;
}

// Instrument configuration
export interface Instrument {
  id: string;
  name: string;
  type: 'drums' | 'bass' | 'guitar' | 'piano' | 'synth' | 'vocal' | 'strings' | 'brass' | 'wind' | 'percussion' | 'other';
  soundFont?: string;
  preset?: string;
  volume: number; // 0-1
  pan: number; // -1 to 1
  isMuted: boolean;
  isSolo: boolean;
  effects: string[];
  midiChannel?: number;
  audioBuffer?: ArrayBuffer;
  waveform?: number[];
  sections: string[]; // Section IDs this instrument appears in
}

// Export history
export interface ExportHistory {
  id: string;
  timestamp: Date;
  format: 'midi' | 'wav' | 'mp3' | 'stems';
  settings: {
    quality?: 'low' | 'medium' | 'high';
    bitrate?: number;
    sampleRate?: number;
    includeEffects?: boolean;
    normalize?: boolean;
  };
  filePath?: string;
  fileSize?: number; // in bytes
  duration?: number; // in seconds
  status: 'pending' | 'processing' | 'completed' | 'failed';
  error?: string;
}

// Project settings
export interface ProjectSettings {
  autoSave: boolean;
  autoSaveInterval: number; // in minutes
  snapToGrid: boolean;
  gridSize: number; // in beats
  showWaveforms: boolean;
  showMidiNotes: boolean;
  showTimeline: boolean;
  showMetronome: boolean;
  metronomeVolume: number; // 0-1
  countIn: number; // in bars
  loopEnabled: boolean;
  loopStart?: number; // in seconds
  loopEnd?: number; // in seconds
  playbackSpeed: number; // 0.5-2.0
  pitchShift: number; // in semitones
}

// Combined project store state
export interface ProjectStoreState extends 
  ProjectMetadata,
  ProjectSettings {
  
  // Project data
  sections: Section[];
  instruments: Instrument[];
  exportHistory: ExportHistory[];
  
  // Current state
  currentProjectId: string | null;
  isDirty: boolean;
  isSaving: boolean;
  isLoading: boolean;
  lastSaveTime: Date | null;
  
  // Playback state
  isPlaying: boolean;
  currentTime: number; // in seconds
  selectedSectionIds: string[];
  selectedInstrumentIds: string[];
  
  // Actions
  setCurrentProjectId: (id: string | null) => void;
  setProjectMetadata: (metadata: Partial<ProjectMetadata>) => void;
  createNewProject: (name: string) => string;
  loadProject: (id: string) => Promise<boolean>;
  saveProject: () => Promise<boolean>;
  deleteProject: (id: string) => Promise<boolean>;
  duplicateProject: (id: string, name?: string) => Promise<string | null>;
  
  // Section actions
  addSection: (section: Omit<Section, 'id'>) => string;
  updateSection: (id: string, updates: Partial<Section>) => void;
  deleteSection: (id: string) => void;
  reorderSections: (fromIndex: number, toIndex: number) => void;
  selectSection: (id: string, multi?: boolean) => void;
  deselectSection: (id: string) => void;
  selectAllSections: () => void;
  deselectAllSections: () => void;
  
  // Instrument actions
  addInstrument: (instrument: Omit<Instrument, 'id'>) => string;
  updateInstrument: (id: string, updates: Partial<Instrument>) => void;
  deleteInstrument: (id: string) => void;
  assignInstrumentToSection: (instrumentId: string, sectionId: string) => void;
  removeInstrumentFromSection: (instrumentId: string, sectionId: string) => void;
  selectInstrument: (id: string, multi?: boolean) => void;
  deselectInstrument: (id: string) => void;
  selectAllInstruments: () => void;
  deselectAllInstruments: () => void;
  
  // Export actions
  addToExportHistory: (exportData: Omit<ExportHistory, 'id'>) => string;
  updateExportHistory: (id: string, updates: Partial<ExportHistory>) => void;
  deleteExportHistory: (id: string) => void;
  clearExportHistory: () => void;
  
  // Settings actions
  updateProjectSettings: (settings: Partial<ProjectSettings>) => void;
  
  // Playback actions
  setIsPlaying: (isPlaying: boolean) => void;
  setCurrentTime: (time: number) => void;
  setLoopPoints: (start?: number, end?: number) => void;
  
  // Utility actions
  markAsDirty: () => void;
  markAsClean: () => void;
  setIsSaving: (isSaving: boolean) => void;
  setIsLoading: (isLoading: boolean) => void;
  updateLastSaveTime: () => void;
  
  // Export/Import actions
  exportProject: () => string;
  importProject: (data: string) => Promise<boolean>;
  exportProjectAsMidi: () => Promise<ArrayBuffer | null>;
  exportProjectAsAudio: (format: 'wav' | 'mp3', settings?: any) => Promise<ArrayBuffer | null>;
}

// Initial state
const initialState: Omit<ProjectStoreState,
  'setCurrentProjectId' | 'setProjectMetadata' | 'createNewProject' | 'loadProject' | 'saveProject' | 'deleteProject' | 'duplicateProject' |
  'addSection' | 'updateSection' | 'deleteSection' | 'reorderSections' | 'selectSection' | 'deselectSection' | 'selectAllSections' | 'deselectAllSections' |
  'addInstrument' | 'updateInstrument' | 'deleteInstrument' | 'assignInstrumentToSection' | 'removeInstrumentFromSection' | 'selectInstrument' | 'deselectInstrument' | 'selectAllInstruments' | 'deselectAllInstruments' |
  'addToExportHistory' | 'updateExportHistory' | 'deleteExportHistory' | 'clearExportHistory' |
  'updateProjectSettings' |
  'setIsPlaying' | 'setCurrentTime' | 'setLoopPoints' |
  'markAsDirty' | 'markAsClean' | 'setIsSaving' | 'setIsLoading' | 'updateLastSaveTime' |
  'exportProject' | 'importProject' | 'exportProjectAsMidi' | 'exportProjectAsAudio'
> = {
  // Project metadata
  id: '',
  name: '',
  description: '',
  createdAt: new Date(),
  updatedAt: new Date(),
  lastOpenedAt: new Date(),
  tags: [],
  isArchived: false,
  isFavorite: false,
  thumbnail: '',
  duration: 0,
  tempo: 120,
  key: 'C',
  timeSignature: '4/4',
  genre: '',
  mood: '',
  
  // Project settings
  autoSave: true,
  autoSaveInterval: 5,
  snapToGrid: true,
  gridSize: 1,
  showWaveforms: true,
  showMidiNotes: true,
  showTimeline: true,
  showMetronome: false,
  metronomeVolume: 0.5,
  countIn: 2,
  loopEnabled: false,
  loopStart: undefined,
  loopEnd: undefined,
  playbackSpeed: 1.0,
  pitchShift: 0,
  
  // Project data
  sections: [],
  instruments: [],
  exportHistory: [],
  
  // Current state
  currentProjectId: null,
  isDirty: false,
  isSaving: false,
  isLoading: false,
  lastSaveTime: null,
  
  // Playback state
  isPlaying: false,
  currentTime: 0,
  selectedSectionIds: [],
  selectedInstrumentIds: [],
};

// Create project store
export const useProjectStore = createStore<ProjectStoreState>(
  (set, get) => ({
    ...initialState,
    
    // Project actions
    setCurrentProjectId: (currentProjectId) => set({ currentProjectId }),
    
    setProjectMetadata: (metadata) => set((state) => ({
      ...state,
      ...metadata,
      updatedAt: new Date(),
      isDirty: true,
    })),
    
    createNewProject: (name) => {
      const id = `project-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      const now = new Date();
      set({
        id,
        name,
        createdAt: now,
        updatedAt: now,
        lastOpenedAt: now,
        currentProjectId: id,
        isDirty: false,
        sections: [],
        instruments: [],
        exportHistory: [],
      });
      return id;
    },
    
    loadProject: async (id) => {
      set({ isLoading: true });
      try {
        // In a real app, this would load from a database or file
        // For now, we'll just set the current project ID
        set({
          currentProjectId: id,
          isLoading: false,
          isDirty: false,
          lastOpenedAt: new Date(),
        });
        return true;
      } catch (error) {
        console.error('Failed to load project:', error);
        set({ isLoading: false });
        return false;
      }
    },
    
    saveProject: async () => {
      const state = get();
      if (!state.currentProjectId) return false;
      
      set({ isSaving: true });
      try {
        // In a real app, this would save to a database or file
        // For now, we'll just update the save time
        set({
          isSaving: false,
          isDirty: false,
          lastSaveTime: new Date(),
          updatedAt: new Date(),
        });
        return true;
      } catch (error) {
        console.error('Failed to save project:', error);
        set({ isSaving: false });
        return false;
      }
    },
    
    deleteProject: async (id) => {
      try {
        // In a real app, this would delete from a database or file
        // For now, we'll just clear the current project if it matches
        const state = get();
        if (state.currentProjectId === id) {
          set({
            currentProjectId: null,
            isDirty: false,
            sections: [],
            instruments: [],
            exportHistory: [],
          });
        }
        return true;
      } catch (error) {
        console.error('Failed to delete project:', error);
        return false;
      }
    },
    
    duplicateProject: async (id, name) => {
      try {
        // In a real app, this would duplicate the project data
        // For now, we'll just create a new project with a suffix
        const state = get();
        const newName = name || `${state.name} (Copy)`;
        const newId = get().createNewProject(newName);
        return newId;
      } catch (error) {
        console.error('Failed to duplicate project:', error);
        return null;
      }
    },
    
    // Section actions
    addSection: (section) => {
      const id = `section-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      set((state) => ({
        sections: [...state.sections, { ...section, id }],
        isDirty: true,
        updatedAt: new Date(),
      }));
      return id;
    },
    
    updateSection: (id, updates) => set((state) => ({
      sections: state.sections.map(section =>
        section.id === id ? { ...section, ...updates } : section
      ),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    deleteSection: (id) => set((state) => ({
      sections: state.sections.filter(section => section.id !== id),
      selectedSectionIds: state.selectedSectionIds.filter(selectedId => selectedId !== id),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    reorderSections: (fromIndex, toIndex) => set((state) => {
      const newSections = [...state.sections];
      const [movedSection] = newSections.splice(fromIndex, 1);
      newSections.splice(toIndex, 0, movedSection);
      return {
        sections: newSections,
        isDirty: true,
        updatedAt: new Date(),
      };
    }),
    
    selectSection: (id, multi = false) => set((state) => ({
      selectedSectionIds: multi
        ? state.selectedSectionIds.includes(id)
          ? state.selectedSectionIds
          : [...state.selectedSectionIds, id]
        : [id],
    })),
    
    deselectSection: (id) => set((state) => ({
      selectedSectionIds: state.selectedSectionIds.filter(selectedId => selectedId !== id),
    })),
    
    selectAllSections: () => set((state) => ({
      selectedSectionIds: state.sections.map(section => section.id),
    })),
    
    deselectAllSections: () => set({ selectedSectionIds: [] }),
    
    // Instrument actions
    addInstrument: (instrument) => {
      const id = `instrument-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      set((state) => ({
        instruments: [...state.instruments, { ...instrument, id }],
        isDirty: true,
        updatedAt: new Date(),
      }));
      return id;
    },
    
    updateInstrument: (id, updates) => set((state) => ({
      instruments: state.instruments.map(instrument =>
        instrument.id === id ? { ...instrument, ...updates } : instrument
      ),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    deleteInstrument: (id) => set((state) => ({
      instruments: state.instruments.filter(instrument => instrument.id !== id),
      selectedInstrumentIds: state.selectedInstrumentIds.filter(selectedId => selectedId !== id),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    assignInstrumentToSection: (instrumentId, sectionId) => set((state) => ({
      instruments: state.instruments.map(instrument =>
        instrument.id === instrumentId
          ? { ...instrument, sections: [...new Set([...instrument.sections, sectionId])] }
          : instrument
      ),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    removeInstrumentFromSection: (instrumentId, sectionId) => set((state) => ({
      instruments: state.instruments.map(instrument =>
        instrument.id === instrumentId
          ? { ...instrument, sections: instrument.sections.filter(id => id !== sectionId) }
          : instrument
      ),
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    selectInstrument: (id, multi = false) => set((state) => ({
      selectedInstrumentIds: multi
        ? state.selectedInstrumentIds.includes(id)
          ? state.selectedInstrumentIds
          : [...state.selectedInstrumentIds, id]
        : [id],
    })),
    
    deselectInstrument: (id) => set((state) => ({
      selectedInstrumentIds: state.selectedInstrumentIds.filter(selectedId => selectedId !== id),
    })),
    
    selectAllInstruments: () => set((state) => ({
      selectedInstrumentIds: state.instruments.map(instrument => instrument.id),
    })),
    
    deselectAllInstruments: () => set({ selectedInstrumentIds: [] }),
    
    // Export actions
    addToExportHistory: (exportData) => {
      const id = `export-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      set((state) => ({
        exportHistory: [...state.exportHistory, { ...exportData, id }],
      }));
      return id;
    },
    
    updateExportHistory: (id, updates) => set((state) => ({
      exportHistory: state.exportHistory.map(item =>
        item.id === id ? { ...item, ...updates } : item
      ),
    })),
    
    deleteExportHistory: (id) => set((state) => ({
      exportHistory: state.exportHistory.filter(item => item.id !== id),
    })),
    
    clearExportHistory: () => set({ exportHistory: [] }),
    
    // Settings actions
    updateProjectSettings: (settings) => set((state) => ({
      ...state,
      ...settings,
      isDirty: true,
      updatedAt: new Date(),
    })),
    
    // Playback actions
    setIsPlaying: (isPlaying) => set({ isPlaying }),
    setCurrentTime: (currentTime) => set({ currentTime }),
    setLoopPoints: (loopStart, loopEnd) => set({ loopStart, loopEnd }),
    
    // Utility actions
    markAsDirty: () => set({ isDirty: true }),
    markAsClean: () => set({ isDirty: false }),
    setIsSaving: (isSaving) => set({ isSaving }),
    setIsLoading: (isLoading) => set({ isLoading }),
    updateLastSaveTime: () => set({ lastSaveTime: new Date() }),
    
    // Export/Import actions
    exportProject: () => {
      const state = get();
      const exportData = {
        metadata: {
          name: state.name,
          description: state.description,
          tags: state.tags,
          tempo: state.tempo,
          key: state.key,
          timeSignature: state.timeSignature,
          genre: state.genre,
          mood: state.mood,
        },
        settings: {
          autoSave: state.autoSave,
          autoSaveInterval: state.autoSaveInterval,
          snapToGrid: state.snapToGrid,
          gridSize: state.gridSize,
          showWaveforms: state.showWaveforms,
          showMidiNotes: state.showMidiNotes,
          showTimeline: state.showTimeline,
          showMetronome: state.showMetronome,
          metronomeVolume: state.metronomeVolume,
          countIn: state.countIn,
          loopEnabled: state.loopEnabled,
          playbackSpeed: state.playbackSpeed,
          pitchShift: state.pitchShift,
        },
        sections: state.sections,
        instruments: state.instruments,
      };
      return JSON.stringify(exportData, null, 2);
    },
    
    importProject: async (data) => {
      try {
        const importData = JSON.parse(data);
        
        if (importData.metadata) {
          get().setProjectMetadata(importData.metadata);
        }
        
        if (importData.settings) {
          get().updateProjectSettings(importData.settings);
        }
        
        if (importData.sections) {
          set({ sections: importData.sections });
        }
        
        if (importData.instruments) {
          set({ instruments: importData.instruments });
        }
        
        set({
          isDirty: true,
          updatedAt: new Date(),
        });
        
        return true;
      } catch (error) {
        console.error('Failed to import project:', error);
        return false;
      }
    },
    
    exportProjectAsMidi: async () => {
      try {
        // In a real app, this would convert the project to MIDI format
        // For now, we'll return null
        return null;
      } catch (error) {
        console.error('Failed to export MIDI:', error);
        return null;
      }
    },
    
    exportProjectAsAudio: async (format, settings) => {
      try {
        // In a real app, this would render the project to audio
        // For now, we'll return null
        return null;
      } catch (error) {
        console.error('Failed to export audio:', error);
        return null;
      }
    },
  }),
  {
    name: 'voxforge-project-store',
    persist: {
      name: 'voxforge-project-store',
      partialize: (state) => ({
        currentProjectId: state.currentProjectId,
        autoSave: state.autoSave,
        autoSaveInterval: state.autoSaveInterval,
        snapToGrid: state.snapToGrid,
        gridSize: state.gridSize,
        showWaveforms: state.showWaveforms,
        showMidiNotes: state.showMidiNotes,
        showTimeline: state.showTimeline,
        showMetronome: state.showMetronome,
        metronomeVolume: state.metronomeVolume,
        countIn: state.countIn,
        loopEnabled: state.loopEnabled,
        playbackSpeed: state.playbackSpeed,
        pitchShift: state.pitchShift,
      }),
    },
  }
);

// Selectors for optimized re-renders
export const useProjectMetadata = () => {
  const store = useProjectStore();
  return {
    id: store.id,
    name: store.name,
    description: store.description,
    createdAt: store.createdAt,
    updatedAt: store.updatedAt,
    lastOpenedAt: store.lastOpenedAt,
    tags: store.tags,
    isArchived: store.isArchived,
    isFavorite: store.isFavorite,
    thumbnail: store.thumbnail,
    duration: store.duration,
    tempo: store.tempo,
    key: store.key,
    timeSignature: store.timeSignature,
    genre: store.genre,
    mood: store.mood,
  };
};

export const useProjectSettings = () => {
  const store = useProjectStore();
  return {
    autoSave: store.autoSave,
    autoSaveInterval: store.autoSaveInterval,
    snapToGrid: store.snapToGrid,
    gridSize: store.gridSize,
    showWaveforms: store.showWaveforms,
    showMidiNotes: store.showMidiNotes,
    showTimeline: store.showTimeline,
    showMetronome: store.showMetronome,
    metronomeVolume: store.metronomeVolume,
    countIn: store.countIn,
    loopEnabled: store.loopEnabled,
    loopStart: store.loopStart,
    loopEnd: store.loopEnd,
    playbackSpeed: store.playbackSpeed,
    pitchShift: store.pitchShift,
  };
};

export const useProjectSections = () => {
  const store = useProjectStore();
  return {
    sections: store.sections,
    selectedSectionIds: store.selectedSectionIds,
  };
};

export const useProjectInstruments = () => {
  const store = useProjectStore();
  return {
    instruments: store.instruments,
    selectedInstrumentIds: store.selectedInstrumentIds,
  };
};

export const useProjectExportHistory = () => {
  const store = useProjectStore();
  return {
    exportHistory: store.exportHistory,
  };
};

export const useProjectState = () => {
  const store = useProjectStore();
  return {
    currentProjectId: store.currentProjectId,
    isDirty: store.isDirty,
    isSaving: store.isSaving,
    isLoading: store.isLoading,
    lastSaveTime: store.lastSaveTime,
    isPlaying: store.isPlaying,
    currentTime: store.currentTime,
  };
};

export const useProjectActions = () => {
  const store = useProjectStore();
  return {
    setCurrentProjectId: store.setCurrentProjectId,
    setProjectMetadata: store.setProjectMetadata,
    createNewProject: store.createNewProject,
    loadProject: store.loadProject,
    saveProject: store.saveProject,
    deleteProject: store.deleteProject,
    duplicateProject: store.duplicateProject,
    addSection: store.addSection,
    updateSection: store.updateSection,
    deleteSection: store.deleteSection,
    reorderSections: store.reorderSections,
    selectSection: store.selectSection,
    deselectSection: store.deselectSection,
    selectAllSections: store.selectAllSections,
    deselectAllSections: store.deselectAllSections,
    addInstrument: store.addInstrument,
    updateInstrument: store.updateInstrument,
    deleteInstrument: store.deleteInstrument,
    assignInstrumentToSection: store.assignInstrumentToSection,
    removeInstrumentFromSection: store.removeInstrumentFromSection,
    selectInstrument: store.selectInstrument,
    deselectInstrument: store.deselectInstrument,
    selectAllInstruments: store.selectAllInstruments,
    deselectAllInstruments: store.deselectAllInstruments,
    addToExportHistory: store.addToExportHistory,
    updateExportHistory: store.updateExportHistory,
    deleteExportHistory: store.deleteExportHistory,
    clearExportHistory: store.clearExportHistory,
    updateProjectSettings: store.updateProjectSettings,
    setIsPlaying: store.setIsPlaying,
    setCurrentTime: store.setCurrentTime,
    setLoopPoints: store.setLoopPoints,
    markAsDirty: store.markAsDirty,
    markAsClean: store.markAsClean,
    setIsSaving: store.setIsSaving,
    setIsLoading: store.setIsLoading,
    updateLastSaveTime: store.updateLastSaveTime,
    exportProject: store.exportProject,
    importProject: store.importProject,
    exportProjectAsMidi: store.exportProjectAsMidi,
    exportProjectAsAudio: store.exportProjectAsAudio,
  };
};