import {
  useAudioStore,
  useRecordingState,
  useAnalysisState,
  useGeneratedAudioState,
  useAudioActions
} from './audioStore';

import {
  useUIStore,
  useNavigationState,
  useModalState,
  useThemePreferences,
  useLayoutPreferences,
  useVisualizationState,
  useLoadingState,
  useUIActions
} from './uiStore';

import { 
  useUserStore,
  useOnboardingState,
  useUserPreferences,
  usePrivacySettings,
  useSessionData,
  useUserProfile,
  useUserActions
} from './userStore';

import { 
  useProjectStore,
  useProjectMetadata,
  useProjectSettings,
  useProjectSections,
  useProjectInstruments,
  useProjectExportHistory,
  useProjectState,
  useProjectActions
} from './projectStore';

// Re-export all store hooks
export {
  // Audio Store
  useAudioStore,
  useRecordingState,
  useAnalysisState,
  useGeneratedAudioState,
  useAudioActions,
  
  // UI Store
  useUIStore,
  useNavigationState,
  useModalState,
  useThemePreferences,
  useLayoutPreferences,
  useVisualizationState,
  useLoadingState,
  useUIActions,
  
  // User Store
  useUserStore,
  useOnboardingState,
  useUserPreferences,
  usePrivacySettings,
  useSessionData,
  useUserProfile,
  useUserActions,
  
  // Project Store
  useProjectStore,
  useProjectMetadata,
  useProjectSettings,
  useProjectSections,
  useProjectInstruments,
  useProjectExportHistory,
  useProjectState,
  useProjectActions,
};

// Combined hooks for common use cases

/**
 * Hook for accessing all store states
 * Useful for debugging or when you need multiple store states
 */
export const useAllStores = () => {
  const audioStore = useAudioStore();
  const uiStore = useUIStore();
  const userStore = useUserStore();
  const projectStore = useProjectStore();
  
  return {
    audio: audioStore,
    ui: uiStore,
    user: userStore,
    project: projectStore,
  };
};

/**
 * Hook for accessing all store actions
 * Useful for components that need to trigger actions across multiple stores
 */
export const useAllStoreActions = () => {
  const audioActions = useAudioActions();
  const uiActions = useUIActions();
  const userActions = useUserActions();
  const projectActions = useProjectActions();
  
  return {
    audio: audioActions,
    ui: uiActions,
    user: userActions,
    project: projectActions,
  };
};

/**
 * Hook for accessing app initialization state
 * Combines relevant states from multiple stores
 */
export const useAppInitialization = () => {
  const { hasSeenOnboarding } = useOnboardingState();
  const isHydrated = true; // Will be implemented in StoreProvider
  const { currentProjectId } = useProjectState();
  const { sessionId } = useSessionData();
  
  return {
    isInitialized: isHydrated,
    hasSeenOnboarding,
    hasActiveProject: !!currentProjectId,
    hasActiveSession: !!sessionId,
  };
};

/**
 * Hook for accessing current project context
 * Combines project metadata with relevant UI and audio states
 */
export const useCurrentProject = () => {
  const metadata = useProjectMetadata();
  const settings = useProjectSettings();
  const { sections, selectedSectionIds } = useProjectSections();
  const { instruments, selectedInstrumentIds } = useProjectInstruments();
  const { isDirty, isSaving, isLoading } = useProjectState();
  const { activeSection } = useNavigationState();
  const { isRecording } = useRecordingState();
  const { isAnalyzing } = useAnalysisState();
  
  return {
    metadata,
    settings,
    sections,
    selectedSectionIds,
    instruments,
    selectedInstrumentIds,
    isDirty,
    isSaving,
    isLoading,
    activeSection,
    isRecording,
    isAnalyzing,
  };
};

/**
 * Hook for accessing user preferences and settings
 * Combines user preferences with UI theme and layout settings
 */
export const useUserPreferencesAndSettings = () => {
  const preferences = useUserPreferences();
  const privacy = usePrivacySettings();
  const { theme } = useThemePreferences();
  const { sidebarCollapsed, showGrid } = useLayoutPreferences();
  const layout = 'horizontal'; // Default layout
  
  return {
    ...preferences,
    privacy,
    theme,
    layout,
    sidebarCollapsed,
    showGrid,
  };
};

/**
 * Hook for accessing audio and analysis state
 * Combines recording state with analysis results
 */
export const useAudioAndAnalysis = () => {
  const recording = useRecordingState();
  const analysis = useAnalysisState();
  const generated = useGeneratedAudioState();
  const { exportSettings } = useAudioStore();
  
  return {
    recording,
    analysis,
    generated,
    exportSettings,
  };
};

/**
 * Hook for accessing modal and navigation state
 * Combines modal states with navigation state
 */
export const useModalsAndNavigation = () => {
  const modals = useModalState();
  const navigation = useNavigationState();
  const { isMobileMenuOpen } = useUIStore();
  
  return {
    ...modals,
    ...navigation,
    isMobileMenuOpen,
  };
};

/**
 * Hook for accessing export functionality
 * Combines export settings with project export history
 */
export const useExportFunctionality = () => {
  const { exportSettings } = useAudioStore();
  const { exportHistory } = useProjectExportHistory();
  const isExporting = false; // Will be added to UI store if needed
  const { exportProject, exportProjectAsMidi, exportProjectAsAudio } = useProjectActions();
  
  return {
    settings: exportSettings,
    history: exportHistory,
    isExporting,
    actions: {
      exportProject,
      exportProjectAsMidi,
      exportProjectAsAudio,
    },
  };
};

/**
 * Hook for accessing onboarding and help state
 * Combines onboarding state with modal states for help components
 */
export const useOnboardingAndHelp = () => {
  const onboarding = useOnboardingState();
  const { welcomeModal, onboardingTour } = useModalState();
  const { completeOnboarding } = useUserActions();
  const { setWelcomeModal, setOnboardingTour } = useUIActions();
  
  return {
    ...onboarding,
    showWelcomeModal: welcomeModal,
    showQuickStartGuide: false, // Will be added to UI store if needed
    showOnboardingTour: onboardingTour,
    actions: {
      completeOnboarding,
      setModalState: (modal: string, isOpen: boolean) => {
        if (modal === 'welcome') setWelcomeModal(isOpen);
        if (modal === 'onboarding') setOnboardingTour(isOpen);
      },
    },
  };
};

/**
 * Hook for accessing performance and error state
 * Combines loading states with error information
 */
export const usePerformanceAndErrors = () => {
  const { isLoading, globalError: error } = useLoadingState();
  const { errorsEncountered } = useSessionData();
  const { incrementErrorsEncountered } = useUserActions();
  
  return {
    isLoading,
    error,
    errorsEncountered,
    actions: {
      incrementErrorsEncountered,
    },
  };
};

/**
 * Hook for accessing persistence state
 * Combines persistence-related states from multiple stores
 */
export const usePersistenceState = () => {
  const { autoSave, autoSaveInterval } = useProjectSettings();
  const { autoSave: userAutoSave } = useUserPreferences();
  const { isDirty, isSaving, lastSaveTime } = useProjectState();
  const isHydrated = true; // Will be implemented in StoreProvider
  
  return {
    projectAutoSave: autoSave,
    projectAutoSaveInterval: autoSaveInterval,
    userAutoSave,
    isDirty,
    isSaving,
    lastSaveTime,
    isHydrated,
  };
};

/**
 * Hook for accessing mobile-specific state
 * Combines mobile-related states from multiple stores
 */
export const useMobileState = () => {
  const { isMobileMenuOpen, activeSection } = useNavigationState();
  const isMobileDevice = false; // Will be detected in StoreProvider
  const { theme } = useThemePreferences();
  const layout = 'horizontal'; // Default layout
  const { enableHapticFeedback, enableNotifications } = useUserPreferences();
  
  return {
    isMobileMenuOpen,
    activeSection,
    isMobileDevice,
    layout,
    enableHapticFeedback,
    enableNotifications,
  };
};

/**
 * Hook for accessing analytics and telemetry state
 * Combines analytics consent with session data
 */
export const useAnalyticsState = () => {
  const { analyticsConsent, usageData, shareAnonymousData } = usePrivacySettings();
  const { sessionId, featuresUsed, recordingsCount, exportsCount } = useSessionData();
  const { addFeatureUsed, incrementRecordingsCount, incrementExportsCount } = useUserActions();
  
  return {
    consent: {
      analytics: analyticsConsent,
      usage: usageData,
      anonymous: shareAnonymousData,
    },
    session: {
      id: sessionId,
      featuresUsed,
      recordingsCount,
      exportsCount,
    },
    actions: {
      addFeatureUsed,
      incrementRecordingsCount,
      incrementExportsCount,
    },
  };
};