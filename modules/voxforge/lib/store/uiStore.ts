import { createStore } from './index';

// Navigation state
export interface NavigationState {
  activeSection: string;
  isMobileMenuOpen: boolean;
  isNavigationVisible: boolean;
  breadcrumbs: Array<{ label: string; href: string }>;
}

// Modal state
export interface ModalState {
  welcomeModal: boolean;
  onboardingTour: boolean;
  settingsModal: boolean;
  exportModal: boolean;
  helpModal: boolean;
  confirmModal: {
    isOpen: boolean;
    title: string;
    message: string;
    onConfirm?: () => void;
    onCancel?: () => void;
  };
}

// Theme and layout preferences
export interface ThemePreferences {
  theme: 'light' | 'dark' | 'system';
  primaryColor: string;
  accentColor: string;
  fontSize: 'small' | 'medium' | 'large';
  reduceAnimations: boolean;
  highContrast: boolean;
}

// Layout preferences
export interface LayoutPreferences {
  sidebarCollapsed: boolean;
  sidebarWidth: number;
  showGrid: boolean;
  snapToGrid: boolean;
  zoomLevel: number;
  panelLayout: 'horizontal' | 'vertical' | 'tabs';
  showTooltips: boolean;
  showKeyboardShortcuts: boolean;
}

// Pixi.js visualization state
export interface VisualizationState {
  pixiMode: 'record' | 'edit' | 'game' | null;
  visualizerEnabled: boolean;
  pianoRollEnabled: boolean;
  rhythmGameEnabled: boolean;
  visualizerSettings: {
    particleCount: number;
    colorScheme: 'rainbow' | 'monochrome' | 'gradient';
    sensitivity: number;
  };
}

// Loading and error states
export interface LoadingState {
  isLoading: boolean;
  loadingMessage: string;
  globalError: string | null;
  warnings: Array<{ id: string; message: string; type: 'info' | 'warning' | 'error' }>;
}

// Combined UI store state
export interface UIStoreState extends 
  NavigationState, 
  ModalState, 
  ThemePreferences, 
  LayoutPreferences, 
  VisualizationState, 
  LoadingState {
  
  // Navigation actions
  setActiveSection: (section: string) => void;
  setMobileMenuOpen: (isOpen: boolean) => void;
  setNavigationVisible: (isVisible: boolean) => void;
  setBreadcrumbs: (breadcrumbs: Array<{ label: string; href: string }>) => void;
  addBreadcrumb: (breadcrumb: { label: string; href: string }) => void;
  removeBreadcrumb: (index: number) => void;
  
  // Modal actions
  setWelcomeModal: (isOpen: boolean) => void;
  setOnboardingTour: (isOpen: boolean) => void;
  setSettingsModal: (isOpen: boolean) => void;
  setExportModal: (isOpen: boolean) => void;
  setHelpModal: (isOpen: boolean) => void;
  setConfirmModal: (modal: ModalState['confirmModal']) => void;
  closeAllModals: () => void;
  
  // Theme actions
  setTheme: (theme: ThemePreferences['theme']) => void;
  setPrimaryColor: (color: string) => void;
  setAccentColor: (color: string) => void;
  setFontSize: (size: ThemePreferences['fontSize']) => void;
  setReduceAnimations: (reduce: boolean) => void;
  setHighContrast: (highContrast: boolean) => void;
  
  // Layout actions
  setSidebarCollapsed: (collapsed: boolean) => void;
  setSidebarWidth: (width: number) => void;
  setShowGrid: (show: boolean) => void;
  setSnapToGrid: (snap: boolean) => void;
  setZoomLevel: (level: number) => void;
  setPanelLayout: (layout: LayoutPreferences['panelLayout']) => void;
  setShowTooltips: (show: boolean) => void;
  setShowKeyboardShortcuts: (show: boolean) => void;
  
  // Visualization actions
  setPixiMode: (mode: VisualizationState['pixiMode']) => void;
  setVisualizerEnabled: (enabled: boolean) => void;
  setPianoRollEnabled: (enabled: boolean) => void;
  setRhythmGameEnabled: (enabled: boolean) => void;
  setVisualizerSettings: (settings: Partial<VisualizationState['visualizerSettings']>) => void;
  
  // Loading and error actions
  setLoading: (isLoading: boolean, message?: string) => void;
  setGlobalError: (error: string | null) => void;
  addWarning: (warning: { id: string; message: string; type: 'info' | 'warning' | 'error' }) => void;
  removeWarning: (id: string) => void;
  clearWarnings: () => void;
  
  // Reset actions
  resetNavigation: () => void;
  resetModals: () => void;
  resetTheme: () => void;
  resetLayout: () => void;
  resetVisualization: () => void;
  resetLoading: () => void;
  resetAll: () => void;
}

// Initial state
const initialState: Omit<UIStoreState, 
  'setActiveSection' | 'setMobileMenuOpen' | 'setNavigationVisible' | 'setBreadcrumbs' | 'addBreadcrumb' | 'removeBreadcrumb' |
  'setWelcomeModal' | 'setOnboardingTour' | 'setSettingsModal' | 'setExportModal' | 'setHelpModal' | 'setConfirmModal' | 'closeAllModals' |
  'setTheme' | 'setPrimaryColor' | 'setAccentColor' | 'setFontSize' | 'setReduceAnimations' | 'setHighContrast' |
  'setSidebarCollapsed' | 'setSidebarWidth' | 'setShowGrid' | 'setSnapToGrid' | 'setZoomLevel' | 'setPanelLayout' | 'setShowTooltips' | 'setShowKeyboardShortcuts' |
  'setPixiMode' | 'setVisualizerEnabled' | 'setPianoRollEnabled' | 'setRhythmGameEnabled' | 'setVisualizerSettings' |
  'setLoading' | 'setGlobalError' | 'addWarning' | 'removeWarning' | 'clearWarnings' |
  'resetNavigation' | 'resetModals' | 'resetTheme' | 'resetLayout' | 'resetVisualization' | 'resetLoading' | 'resetAll'
> = {
  // Navigation state
  activeSection: 'record',
  isMobileMenuOpen: false,
  isNavigationVisible: true,
  breadcrumbs: [],
  
  // Modal state
  welcomeModal: false,
  onboardingTour: false,
  settingsModal: false,
  exportModal: false,
  helpModal: false,
  confirmModal: {
    isOpen: false,
    title: '',
    message: '',
  },
  
  // Theme preferences
  theme: 'dark',
  primaryColor: '#3b82f6',
  accentColor: '#8b5cf6',
  fontSize: 'medium',
  reduceAnimations: false,
  highContrast: false,
  
  // Layout preferences
  sidebarCollapsed: false,
  sidebarWidth: 280,
  showGrid: false,
  snapToGrid: false,
  zoomLevel: 1,
  panelLayout: 'horizontal',
  showTooltips: true,
  showKeyboardShortcuts: true,
  
  // Visualization state
  pixiMode: null,
  visualizerEnabled: true,
  pianoRollEnabled: true,
  rhythmGameEnabled: true,
  visualizerSettings: {
    particleCount: 100,
    colorScheme: 'rainbow',
    sensitivity: 0.5,
  },
  
  // Loading and error states
  isLoading: false,
  loadingMessage: '',
  globalError: null,
  warnings: [],
};

// Create the UI store
export const useUIStore = createStore<UIStoreState>(
  (set, get) => ({
    ...initialState,
    
    // Navigation actions
    setActiveSection: (activeSection) => set({ activeSection }),
    setMobileMenuOpen: (isMobileMenuOpen) => set({ isMobileMenuOpen }),
    setNavigationVisible: (isNavigationVisible) => set({ isNavigationVisible }),
    setBreadcrumbs: (breadcrumbs) => set({ breadcrumbs }),
    addBreadcrumb: (breadcrumb) => set((state) => ({
      breadcrumbs: [...state.breadcrumbs, breadcrumb]
    })),
    removeBreadcrumb: (index) => set((state) => ({
      breadcrumbs: state.breadcrumbs.filter((_, i) => i !== index)
    })),
    
    // Modal actions
    setWelcomeModal: (welcomeModal) => set({ welcomeModal }),
    setOnboardingTour: (onboardingTour) => set({ onboardingTour }),
    setSettingsModal: (settingsModal) => set({ settingsModal }),
    setExportModal: (exportModal) => set({ exportModal }),
    setHelpModal: (helpModal) => set({ helpModal }),
    setConfirmModal: (confirmModal) => set({ confirmModal }),
    closeAllModals: () => set({
      welcomeModal: false,
      onboardingTour: false,
      settingsModal: false,
      exportModal: false,
      helpModal: false,
      confirmModal: { isOpen: false, title: '', message: '' },
    }),
    
    // Theme actions
    setTheme: (theme) => set({ theme }),
    setPrimaryColor: (primaryColor) => set({ primaryColor }),
    setAccentColor: (accentColor) => set({ accentColor }),
    setFontSize: (fontSize) => set({ fontSize }),
    setReduceAnimations: (reduceAnimations) => set({ reduceAnimations }),
    setHighContrast: (highContrast) => set({ highContrast }),
    
    // Layout actions
    setSidebarCollapsed: (sidebarCollapsed) => set({ sidebarCollapsed }),
    setSidebarWidth: (sidebarWidth) => set({ sidebarWidth }),
    setShowGrid: (showGrid) => set({ showGrid }),
    setSnapToGrid: (snapToGrid) => set({ snapToGrid }),
    setZoomLevel: (zoomLevel) => set({ zoomLevel }),
    setPanelLayout: (panelLayout) => set({ panelLayout }),
    setShowTooltips: (showTooltips) => set({ showTooltips }),
    setShowKeyboardShortcuts: (showKeyboardShortcuts) => set({ showKeyboardShortcuts }),
    
    // Visualization actions
    setPixiMode: (pixiMode) => set({ pixiMode }),
    setVisualizerEnabled: (visualizerEnabled) => set({ visualizerEnabled }),
    setPianoRollEnabled: (pianoRollEnabled) => set({ pianoRollEnabled }),
    setRhythmGameEnabled: (rhythmGameEnabled) => set({ rhythmGameEnabled }),
    setVisualizerSettings: (visualizerSettings) => set((state) => ({
      visualizerSettings: { ...state.visualizerSettings, ...visualizerSettings }
    })),
    
    // Loading and error actions
    setLoading: (isLoading, loadingMessage = '') => set({ isLoading, loadingMessage }),
    setGlobalError: (globalError) => set({ globalError }),
    addWarning: (warning) => set((state) => ({
      warnings: [...state.warnings, warning]
    })),
    removeWarning: (id) => set((state) => ({
      warnings: state.warnings.filter(w => w.id !== id)
    })),
    clearWarnings: () => set({ warnings: [] }),
    
    // Reset actions
    resetNavigation: () => set({
      activeSection: 'record',
      isMobileMenuOpen: false,
      isNavigationVisible: true,
      breadcrumbs: [],
    }),
    
    resetModals: () => set({
      welcomeModal: false,
      onboardingTour: false,
      settingsModal: false,
      exportModal: false,
      helpModal: false,
      confirmModal: { isOpen: false, title: '', message: '' },
    }),
    
    resetTheme: () => set({
      theme: 'dark',
      primaryColor: '#3b82f6',
      accentColor: '#8b5cf6',
      fontSize: 'medium',
      reduceAnimations: false,
      highContrast: false,
    }),
    
    resetLayout: () => set({
      sidebarCollapsed: false,
      sidebarWidth: 280,
      showGrid: false,
      snapToGrid: false,
      zoomLevel: 1,
      panelLayout: 'horizontal',
      showTooltips: true,
      showKeyboardShortcuts: true,
    }),
    
    resetVisualization: () => set({
      pixiMode: null,
      visualizerEnabled: true,
      pianoRollEnabled: true,
      rhythmGameEnabled: true,
      visualizerSettings: {
        particleCount: 100,
        colorScheme: 'rainbow',
        sensitivity: 0.5,
      },
    }),
    
    resetLoading: () => set({
      isLoading: false,
      loadingMessage: '',
      globalError: null,
      warnings: [],
    }),
    
    resetAll: () => set(initialState),
  }),
  {
    name: 'voxforge-ui-store',
    persist: {
      name: 'voxforge-ui-store',
      partialize: (state) => ({
        theme: state.theme,
        primaryColor: state.primaryColor,
        accentColor: state.accentColor,
        fontSize: state.fontSize,
        reduceAnimations: state.reduceAnimations,
        highContrast: state.highContrast,
        sidebarCollapsed: state.sidebarCollapsed,
        sidebarWidth: state.sidebarWidth,
        showGrid: state.showGrid,
        snapToGrid: state.snapToGrid,
        zoomLevel: state.zoomLevel,
        panelLayout: state.panelLayout,
        showTooltips: state.showTooltips,
        showKeyboardShortcuts: state.showKeyboardShortcuts,
        visualizerSettings: state.visualizerSettings,
      }),
    },
  }
);

// Selectors for optimized re-renders
export const useNavigationState = () => {
  const store = useUIStore();
  return {
    activeSection: store.activeSection,
    isMobileMenuOpen: store.isMobileMenuOpen,
    isNavigationVisible: store.isNavigationVisible,
    breadcrumbs: store.breadcrumbs,
  };
};

export const useModalState = () => {
  const store = useUIStore();
  return {
    welcomeModal: store.welcomeModal,
    onboardingTour: store.onboardingTour,
    settingsModal: store.settingsModal,
    exportModal: store.exportModal,
    helpModal: store.helpModal,
    confirmModal: store.confirmModal,
  };
};

export const useThemePreferences = () => {
  const store = useUIStore();
  return {
    theme: store.theme,
    primaryColor: store.primaryColor,
    accentColor: store.accentColor,
    fontSize: store.fontSize,
    reduceAnimations: store.reduceAnimations,
    highContrast: store.highContrast,
  };
};

export const useLayoutPreferences = () => {
  const store = useUIStore();
  return {
    sidebarCollapsed: store.sidebarCollapsed,
    sidebarWidth: store.sidebarWidth,
    showGrid: store.showGrid,
    snapToGrid: store.snapToGrid,
    zoomLevel: store.zoomLevel,
    panelLayout: store.panelLayout,
    showTooltips: store.showTooltips,
    showKeyboardShortcuts: store.showKeyboardShortcuts,
  };
};

export const useVisualizationState = () => {
  const store = useUIStore();
  return {
    pixiMode: store.pixiMode,
    visualizerEnabled: store.visualizerEnabled,
    pianoRollEnabled: store.pianoRollEnabled,
    rhythmGameEnabled: store.rhythmGameEnabled,
    visualizerSettings: store.visualizerSettings,
  };
};

export const useLoadingState = () => {
  const store = useUIStore();
  return {
    isLoading: store.isLoading,
    loadingMessage: store.loadingMessage,
    globalError: store.globalError,
    warnings: store.warnings,
  };
};

export const useUIActions = () => {
  const store = useUIStore();
  return {
    setActiveSection: store.setActiveSection,
    setMobileMenuOpen: store.setMobileMenuOpen,
    setNavigationVisible: store.setNavigationVisible,
    setBreadcrumbs: store.setBreadcrumbs,
    addBreadcrumb: store.addBreadcrumb,
    removeBreadcrumb: store.removeBreadcrumb,
    setWelcomeModal: store.setWelcomeModal,
    setOnboardingTour: store.setOnboardingTour,
    setSettingsModal: store.setSettingsModal,
    setExportModal: store.setExportModal,
    setHelpModal: store.setHelpModal,
    setConfirmModal: store.setConfirmModal,
    closeAllModals: store.closeAllModals,
    setTheme: store.setTheme,
    setPrimaryColor: store.setPrimaryColor,
    setAccentColor: store.setAccentColor,
    setFontSize: store.setFontSize,
    setReduceAnimations: store.setReduceAnimations,
    setHighContrast: store.setHighContrast,
    setSidebarCollapsed: store.setSidebarCollapsed,
    setSidebarWidth: store.setSidebarWidth,
    setShowGrid: store.setShowGrid,
    setSnapToGrid: store.setSnapToGrid,
    setZoomLevel: store.setZoomLevel,
    setPanelLayout: store.setPanelLayout,
    setShowTooltips: store.setShowTooltips,
    setShowKeyboardShortcuts: store.setShowKeyboardShortcuts,
    setPixiMode: store.setPixiMode,
    setVisualizerEnabled: store.setVisualizerEnabled,
    setPianoRollEnabled: store.setPianoRollEnabled,
    setRhythmGameEnabled: store.setRhythmGameEnabled,
    setVisualizerSettings: store.setVisualizerSettings,
    setLoading: store.setLoading,
    setGlobalError: store.setGlobalError,
    addWarning: store.addWarning,
    removeWarning: store.removeWarning,
    clearWarnings: store.clearWarnings,
    resetNavigation: store.resetNavigation,
    resetModals: store.resetModals,
    resetTheme: store.resetTheme,
    resetLayout: store.resetLayout,
    resetVisualization: store.resetVisualization,
    resetLoading: store.resetLoading,
    resetAll: store.resetAll,
  };
};