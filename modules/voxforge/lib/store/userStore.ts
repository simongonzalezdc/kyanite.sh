import { createStore } from './index';

// Onboarding state
export interface OnboardingState {
  hasSeenOnboarding: boolean;
  currentStep: number;
  completedSteps: string[];
  isSkipped: boolean;
  lastCompletedStep: string | null;
  onboardingVersion: string;
}

// User preferences
export interface UserPreferences {
  language: string;
  autoSave: boolean;
  autoSaveInterval: number; // in minutes
  defaultExportFormat: 'midi' | 'wav' | 'mp3' | 'stems';
  defaultInstrumentSelection: string[];
  showKeyboardShortcuts: boolean;
  enableHapticFeedback: boolean;
  enableSoundEffects: boolean;
  enableNotifications: boolean;
  autoStartRecording: boolean;
  showAdvancedOptions: boolean;
}

// Analytics and privacy settings
export interface PrivacySettings {
  analyticsConsent: boolean;
  crashReporting: boolean;
  usageData: boolean;
  shareAnonymousData: boolean;
  dataRetentionDays: number;
  cookieConsent: boolean;
  marketingEmails: boolean;
  productUpdates: boolean;
}

// Session data
export interface SessionData {
  sessionId: string;
  startTime: Date;
  lastActivity: Date;
  sessionDuration: number; // in seconds
  recordingsCount: number;
  exportsCount: number;
  featuresUsed: string[];
  errorsEncountered: number;
  deviceInfo: {
    userAgent: string;
    platform: string;
    screenResolution: string;
    isMobile: boolean;
    isTouchDevice: boolean;
  };
}

// User profile
export interface UserProfile {
  id: string;
  name: string;
  email?: string;
  avatar?: string;
  createdAt: Date;
  lastLoginAt: Date;
  totalSessions: number;
  totalRecordingTime: number; // in seconds
  totalExports: number;
  favoriteFeatures: string[];
  skillLevel: 'beginner' | 'intermediate' | 'advanced' | 'expert';
  musicalBackground: string;
}

// Combined user store state
export interface UserStoreState extends 
  OnboardingState, 
  UserPreferences, 
  PrivacySettings, 
  SessionData,
  UserProfile {
  
  // Onboarding actions
  setHasSeenOnboarding: (hasSeen: boolean) => void;
  setCurrentStep: (step: number) => void;
  setCompletedSteps: (steps: string[]) => void;
  addCompletedStep: (step: string) => void;
  setSkipped: (skipped: boolean) => void;
  setLastCompletedStep: (step: string | null) => void;
  setOnboardingVersion: (version: string) => void;
  resetOnboarding: () => void;
  
  // Preferences actions
  setLanguage: (language: string) => void;
  setAutoSave: (autoSave: boolean) => void;
  setAutoSaveInterval: (interval: number) => void;
  setDefaultExportFormat: (format: UserPreferences['defaultExportFormat']) => void;
  setDefaultInstrumentSelection: (instruments: string[]) => void;
  setShowKeyboardShortcuts: (show: boolean) => void;
  setEnableHapticFeedback: (enable: boolean) => void;
  setEnableSoundEffects: (enable: boolean) => void;
  setEnableNotifications: (enable: boolean) => void;
  setAutoStartRecording: (autoStart: boolean) => void;
  setShowAdvancedOptions: (show: boolean) => void;
  
  // Privacy actions
  setAnalyticsConsent: (consent: boolean) => void;
  setCrashReporting: (enabled: boolean) => void;
  setUsageData: (enabled: boolean) => void;
  setShareAnonymousData: (share: boolean) => void;
  setDataRetentionDays: (days: number) => void;
  setCookieConsent: (consent: boolean) => void;
  setMarketingEmails: (enabled: boolean) => void;
  setProductUpdates: (enabled: boolean) => void;
  
  // Session actions
  setSessionId: (id: string) => void;
  setStartTime: (time: Date) => void;
  setLastActivity: (time: Date) => void;
  setSessionDuration: (duration: number) => void;
  incrementRecordingsCount: () => void;
  incrementExportsCount: () => void;
  addFeatureUsed: (feature: string) => void;
  incrementErrorsEncountered: () => void;
  setDeviceInfo: (info: Partial<SessionData['deviceInfo']>) => void;
  resetSession: () => void;
  
  // Profile actions
  setProfileId: (id: string) => void;
  setName: (name: string) => void;
  setEmail: (email: string) => void;
  setAvatar: (avatar: string) => void;
  setCreatedAt: (date: Date) => void;
  setLastLoginAt: (date: Date) => void;
  incrementTotalSessions: () => void;
  setTotalRecordingTime: (time: number) => void;
  incrementTotalExports: () => void;
  setFavoriteFeatures: (features: string[]) => void;
  addFavoriteFeature: (feature: string) => void;
  removeFavoriteFeature: (feature: string) => void;
  setSkillLevel: (level: UserProfile['skillLevel']) => void;
  setMusicalBackground: (background: string) => void;
  
  // Combined actions
  updateActivity: () => void;
  completeOnboarding: () => void;
  exportUserData: () => string;
  importUserData: (data: string) => boolean;
  clearAllData: () => void;
}

// Initial state
const initialState: Omit<UserStoreState,
  'setHasSeenOnboarding' | 'setCurrentStep' | 'setCompletedSteps' | 'addCompletedStep' | 'setSkipped' | 'setLastCompletedStep' | 'setOnboardingVersion' | 'resetOnboarding' |
  'setLanguage' | 'setAutoSave' | 'setAutoSaveInterval' | 'setDefaultExportFormat' | 'setDefaultInstrumentSelection' | 'setShowKeyboardShortcuts' | 'setEnableHapticFeedback' | 'setEnableSoundEffects' | 'setEnableNotifications' | 'setAutoStartRecording' | 'setShowAdvancedOptions' |
  'setAnalyticsConsent' | 'setCrashReporting' | 'setUsageData' | 'setShareAnonymousData' | 'setDataRetentionDays' | 'setCookieConsent' | 'setMarketingEmails' | 'setProductUpdates' |
  'setSessionId' | 'setStartTime' | 'setLastActivity' | 'setSessionDuration' | 'incrementRecordingsCount' | 'incrementExportsCount' | 'addFeatureUsed' | 'incrementErrorsEncountered' | 'setDeviceInfo' | 'resetSession' |
  'setProfileId' | 'setName' | 'setEmail' | 'setAvatar' | 'setCreatedAt' | 'setLastLoginAt' | 'incrementTotalSessions' | 'setTotalRecordingTime' | 'incrementTotalExports' | 'setFavoriteFeatures' | 'addFavoriteFeature' | 'removeFavoriteFeature' | 'setSkillLevel' | 'setMusicalBackground' |
  'updateActivity' | 'completeOnboarding' | 'exportUserData' | 'importUserData' | 'clearAllData'
> = {
  // Onboarding state
  hasSeenOnboarding: false,
  currentStep: 0,
  completedSteps: [],
  isSkipped: false,
  lastCompletedStep: null,
  onboardingVersion: '1.0.0',
  
  // User preferences
  language: 'en',
  autoSave: true,
  autoSaveInterval: 5,
  defaultExportFormat: 'midi',
  defaultInstrumentSelection: ['drums', 'bass', 'chords'],
  showKeyboardShortcuts: true,
  enableHapticFeedback: true,
  enableSoundEffects: true,
  enableNotifications: true,
  autoStartRecording: false,
  showAdvancedOptions: false,
  
  // Privacy settings
  analyticsConsent: false,
  crashReporting: true,
  usageData: false,
  shareAnonymousData: false,
  dataRetentionDays: 30,
  cookieConsent: false,
  marketingEmails: false,
  productUpdates: true,
  
  // Session data
  sessionId: '',
  startTime: new Date(),
  lastActivity: new Date(),
  sessionDuration: 0,
  recordingsCount: 0,
  exportsCount: 0,
  featuresUsed: [],
  errorsEncountered: 0,
  deviceInfo: {
    userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
    platform: typeof navigator !== 'undefined' ? navigator.platform : '',
    screenResolution: typeof window !== 'undefined' ? `${window.screen.width}x${window.screen.height}` : '',
    isMobile: false,
    isTouchDevice: false,
  },
  
  // User profile
  id: '',
  name: '',
  email: '',
  avatar: '',
  createdAt: new Date(),
  lastLoginAt: new Date(),
  totalSessions: 0,
  totalRecordingTime: 0,
  totalExports: 0,
  favoriteFeatures: [],
  skillLevel: 'beginner',
  musicalBackground: '',
};

// Create user store
export const useUserStore = createStore<UserStoreState>(
  (set, get) => ({
    ...initialState,
    
    // Onboarding actions
    setHasSeenOnboarding: (hasSeenOnboarding) => set({ hasSeenOnboarding }),
    setCurrentStep: (currentStep) => set({ currentStep }),
    setCompletedSteps: (completedSteps) => set({ completedSteps }),
    addCompletedStep: (step) => set((state) => ({
      completedSteps: [...new Set([...state.completedSteps, step])]
    })),
    setSkipped: (isSkipped) => set({ isSkipped }),
    setLastCompletedStep: (lastCompletedStep) => set({ lastCompletedStep }),
    setOnboardingVersion: (onboardingVersion) => set({ onboardingVersion }),
    resetOnboarding: () => set({
      hasSeenOnboarding: false,
      currentStep: 0,
      completedSteps: [],
      isSkipped: false,
      lastCompletedStep: null,
    }),
    
    // Preferences actions
    setLanguage: (language) => set({ language }),
    setAutoSave: (autoSave) => set({ autoSave }),
    setAutoSaveInterval: (autoSaveInterval) => set({ autoSaveInterval }),
    setDefaultExportFormat: (defaultExportFormat) => set({ defaultExportFormat }),
    setDefaultInstrumentSelection: (defaultInstrumentSelection) => set({ defaultInstrumentSelection }),
    setShowKeyboardShortcuts: (showKeyboardShortcuts) => set({ showKeyboardShortcuts }),
    setEnableHapticFeedback: (enableHapticFeedback) => set({ enableHapticFeedback }),
    setEnableSoundEffects: (enableSoundEffects) => set({ enableSoundEffects }),
    setEnableNotifications: (enableNotifications) => set({ enableNotifications }),
    setAutoStartRecording: (autoStartRecording) => set({ autoStartRecording }),
    setShowAdvancedOptions: (showAdvancedOptions) => set({ showAdvancedOptions }),
    
    // Privacy actions
    setAnalyticsConsent: (analyticsConsent) => set({ analyticsConsent }),
    setCrashReporting: (crashReporting) => set({ crashReporting }),
    setUsageData: (usageData) => set({ usageData }),
    setShareAnonymousData: (shareAnonymousData) => set({ shareAnonymousData }),
    setDataRetentionDays: (dataRetentionDays) => set({ dataRetentionDays }),
    setCookieConsent: (cookieConsent) => set({ cookieConsent }),
    setMarketingEmails: (marketingEmails) => set({ marketingEmails }),
    setProductUpdates: (productUpdates) => set({ productUpdates }),
    
    // Session actions
    setSessionId: (sessionId) => set({ sessionId }),
    setStartTime: (startTime) => set({ startTime }),
    setLastActivity: (lastActivity) => set({ lastActivity }),
    setSessionDuration: (sessionDuration) => set({ sessionDuration }),
    incrementRecordingsCount: () => set((state) => ({ 
      recordingsCount: state.recordingsCount + 1 
    })),
    incrementExportsCount: () => set((state) => ({ 
      exportsCount: state.exportsCount + 1 
    })),
    addFeatureUsed: (feature) => set((state) => ({
      featuresUsed: [...new Set([...state.featuresUsed, feature])]
    })),
    incrementErrorsEncountered: () => set((state) => ({ 
      errorsEncountered: state.errorsEncountered + 1 
    })),
    setDeviceInfo: (deviceInfo) => set((state) => ({
      deviceInfo: { ...state.deviceInfo, ...deviceInfo }
    })),
    resetSession: () => set({
      sessionId: '',
      startTime: new Date(),
      lastActivity: new Date(),
      sessionDuration: 0,
      recordingsCount: 0,
      exportsCount: 0,
      featuresUsed: [],
      errorsEncountered: 0,
    }),
    
    // Profile actions
    setProfileId: (id) => set({ id }),
    setName: (name) => set({ name }),
    setEmail: (email) => set({ email }),
    setAvatar: (avatar) => set({ avatar }),
    setCreatedAt: (createdAt) => set({ createdAt }),
    setLastLoginAt: (lastLoginAt) => set({ lastLoginAt }),
    incrementTotalSessions: () => set((state) => ({ 
      totalSessions: state.totalSessions + 1 
    })),
    setTotalRecordingTime: (totalRecordingTime) => set({ totalRecordingTime }),
    incrementTotalExports: () => set((state) => ({ 
      totalExports: state.totalExports + 1 
    })),
    setFavoriteFeatures: (favoriteFeatures) => set({ favoriteFeatures }),
    addFavoriteFeature: (feature) => set((state) => ({
      favoriteFeatures: [...new Set([...state.favoriteFeatures, feature])]
    })),
    removeFavoriteFeature: (feature) => set((state) => ({
      favoriteFeatures: state.favoriteFeatures.filter(f => f !== feature)
    })),
    setSkillLevel: (skillLevel) => set({ skillLevel }),
    setMusicalBackground: (musicalBackground) => set({ musicalBackground }),
    
    // Combined actions
    updateActivity: () => set({ lastActivity: new Date() }),
    completeOnboarding: () => set({
      hasSeenOnboarding: true,
      currentStep: 0,
      completedSteps: [],
      isSkipped: false,
      lastCompletedStep: 'completed',
    }),
    exportUserData: () => {
      const state = get();
      const exportData = {
        preferences: {
          language: state.language,
          autoSave: state.autoSave,
          autoSaveInterval: state.autoSaveInterval,
          defaultExportFormat: state.defaultExportFormat,
          defaultInstrumentSelection: state.defaultInstrumentSelection,
          showKeyboardShortcuts: state.showKeyboardShortcuts,
          enableHapticFeedback: state.enableHapticFeedback,
          enableSoundEffects: state.enableSoundEffects,
          enableNotifications: state.enableNotifications,
          autoStartRecording: state.autoStartRecording,
          showAdvancedOptions: state.showAdvancedOptions,
        },
        privacy: {
          analyticsConsent: state.analyticsConsent,
          crashReporting: state.crashReporting,
          usageData: state.usageData,
          shareAnonymousData: state.shareAnonymousData,
          dataRetentionDays: state.dataRetentionDays,
          cookieConsent: state.cookieConsent,
          marketingEmails: state.marketingEmails,
          productUpdates: state.productUpdates,
        },
        profile: {
          name: state.name,
          email: state.email,
          skillLevel: state.skillLevel,
          musicalBackground: state.musicalBackground,
          favoriteFeatures: state.favoriteFeatures,
        },
      };
      return JSON.stringify(exportData, null, 2);
    },
    importUserData: (data) => {
      try {
        const importData = JSON.parse(data);
        const state = get();
        
        if (importData.preferences) {
          // Import preferences
          Object.entries(importData.preferences).forEach(([key, value]) => {
            const setterName = `set${key.charAt(0).toUpperCase() + key.slice(1)}`;
            const setter = (state as any)[setterName];
            if (typeof setter === 'function') {
              setter(value);
            }
          });
        }
        if (importData.privacy) {
          // Import privacy settings
          Object.entries(importData.privacy).forEach(([key, value]) => {
            const setterName = `set${key.charAt(0).toUpperCase() + key.slice(1)}`;
            const setter = (state as any)[setterName];
            if (typeof setter === 'function') {
              setter(value);
            }
          });
        }
        if (importData.profile) {
          // Import profile data
          Object.entries(importData.profile).forEach(([key, value]) => {
            const setterName = `set${key.charAt(0).toUpperCase() + key.slice(1)}`;
            const setter = (state as any)[setterName];
            if (typeof setter === 'function') {
              setter(value);
            }
          });
        }
        return true;
      } catch (error) {
        console.error('Failed to import user data:', error);
        return false;
      }
    },
    clearAllData: () => set(initialState),
  }),
  {
    name: 'voxforge-user-store',
    persist: {
      name: 'voxforge-user-store',
      partialize: (state) => ({
        hasSeenOnboarding: state.hasSeenOnboarding,
        onboardingVersion: state.onboardingVersion,
        language: state.language,
        autoSave: state.autoSave,
        autoSaveInterval: state.autoSaveInterval,
        defaultExportFormat: state.defaultExportFormat,
        defaultInstrumentSelection: state.defaultInstrumentSelection,
        showKeyboardShortcuts: state.showKeyboardShortcuts,
        enableHapticFeedback: state.enableHapticFeedback,
        enableSoundEffects: state.enableSoundEffects,
        enableNotifications: state.enableNotifications,
        autoStartRecording: state.autoStartRecording,
        showAdvancedOptions: state.showAdvancedOptions,
        analyticsConsent: state.analyticsConsent,
        crashReporting: state.crashReporting,
        usageData: state.usageData,
        shareAnonymousData: state.shareAnonymousData,
        dataRetentionDays: state.dataRetentionDays,
        cookieConsent: state.cookieConsent,
        marketingEmails: state.marketingEmails,
        productUpdates: state.productUpdates,
        id: state.id,
        name: state.name,
        email: state.email,
        avatar: state.avatar,
        createdAt: state.createdAt,
        lastLoginAt: state.lastLoginAt,
        totalSessions: state.totalSessions,
        totalRecordingTime: state.totalRecordingTime,
        totalExports: state.totalExports,
        favoriteFeatures: state.favoriteFeatures,
        skillLevel: state.skillLevel,
        musicalBackground: state.musicalBackground,
      }),
    },
  }
);

// Selectors for optimized re-renders
export const useOnboardingState = () => {
  const store = useUserStore();
  return {
    hasSeenOnboarding: store.hasSeenOnboarding,
    currentStep: store.currentStep,
    completedSteps: store.completedSteps,
    isSkipped: store.isSkipped,
    lastCompletedStep: store.lastCompletedStep,
    onboardingVersion: store.onboardingVersion,
  };
};

export const useUserPreferences = () => {
  const store = useUserStore();
  return {
    language: store.language,
    autoSave: store.autoSave,
    autoSaveInterval: store.autoSaveInterval,
    defaultExportFormat: store.defaultExportFormat,
    defaultInstrumentSelection: store.defaultInstrumentSelection,
    showKeyboardShortcuts: store.showKeyboardShortcuts,
    enableHapticFeedback: store.enableHapticFeedback,
    enableSoundEffects: store.enableSoundEffects,
    enableNotifications: store.enableNotifications,
    autoStartRecording: store.autoStartRecording,
    showAdvancedOptions: store.showAdvancedOptions,
  };
};

export const usePrivacySettings = () => {
  const store = useUserStore();
  return {
    analyticsConsent: store.analyticsConsent,
    crashReporting: store.crashReporting,
    usageData: store.usageData,
    shareAnonymousData: store.shareAnonymousData,
    dataRetentionDays: store.dataRetentionDays,
    cookieConsent: store.cookieConsent,
    marketingEmails: store.marketingEmails,
    productUpdates: store.productUpdates,
  };
};

export const useSessionData = () => {
  const store = useUserStore();
  return {
    sessionId: store.sessionId,
    startTime: store.startTime,
    lastActivity: store.lastActivity,
    sessionDuration: store.sessionDuration,
    recordingsCount: store.recordingsCount,
    exportsCount: store.exportsCount,
    featuresUsed: store.featuresUsed,
    errorsEncountered: store.errorsEncountered,
    deviceInfo: store.deviceInfo,
  };
};

export const useUserProfile = () => {
  const store = useUserStore();
  return {
    id: store.id,
    name: store.name,
    email: store.email,
    avatar: store.avatar,
    createdAt: store.createdAt,
    lastLoginAt: store.lastLoginAt,
    totalSessions: store.totalSessions,
    totalRecordingTime: store.totalRecordingTime,
    totalExports: store.totalExports,
    favoriteFeatures: store.favoriteFeatures,
    skillLevel: store.skillLevel,
    musicalBackground: store.musicalBackground,
  };
};

export const useUserActions = () => {
  const store = useUserStore();
  return {
    setHasSeenOnboarding: store.setHasSeenOnboarding,
    setCurrentStep: store.setCurrentStep,
    setCompletedSteps: store.setCompletedSteps,
    addCompletedStep: store.addCompletedStep,
    setSkipped: store.setSkipped,
    setLastCompletedStep: store.setLastCompletedStep,
    setOnboardingVersion: store.setOnboardingVersion,
    resetOnboarding: store.resetOnboarding,
    setLanguage: store.setLanguage,
    setAutoSave: store.setAutoSave,
    setAutoSaveInterval: store.setAutoSaveInterval,
    setDefaultExportFormat: store.setDefaultExportFormat,
    setDefaultInstrumentSelection: store.setDefaultInstrumentSelection,
    setShowKeyboardShortcuts: store.setShowKeyboardShortcuts,
    setEnableHapticFeedback: store.setEnableHapticFeedback,
    setEnableSoundEffects: store.setEnableSoundEffects,
    setEnableNotifications: store.setEnableNotifications,
    setAutoStartRecording: store.setAutoStartRecording,
    setShowAdvancedOptions: store.setShowAdvancedOptions,
    setAnalyticsConsent: store.setAnalyticsConsent,
    setCrashReporting: store.setCrashReporting,
    setUsageData: store.setUsageData,
    setShareAnonymousData: store.setShareAnonymousData,
    setDataRetentionDays: store.setDataRetentionDays,
    setCookieConsent: store.setCookieConsent,
    setMarketingEmails: store.setMarketingEmails,
    setProductUpdates: store.setProductUpdates,
    setSessionId: store.setSessionId,
    setStartTime: store.setStartTime,
    setLastActivity: store.setLastActivity,
    setSessionDuration: store.setSessionDuration,
    incrementRecordingsCount: store.incrementRecordingsCount,
    incrementExportsCount: store.incrementExportsCount,
    addFeatureUsed: store.addFeatureUsed,
    incrementErrorsEncountered: store.incrementErrorsEncountered,
    setDeviceInfo: store.setDeviceInfo,
    resetSession: store.resetSession,
    setProfileId: store.setProfileId,
    setName: store.setName,
    setEmail: store.setEmail,
    setAvatar: store.setAvatar,
    setCreatedAt: store.setCreatedAt,
    setLastLoginAt: store.setLastLoginAt,
    incrementTotalSessions: store.incrementTotalSessions,
    setTotalRecordingTime: store.setTotalRecordingTime,
    incrementTotalExports: store.incrementTotalExports,
    setFavoriteFeatures: store.setFavoriteFeatures,
    addFavoriteFeature: store.addFavoriteFeature,
    removeFavoriteFeature: store.removeFavoriteFeature,
    setSkillLevel: store.setSkillLevel,
    setMusicalBackground: store.setMusicalBackground,
    updateActivity: store.updateActivity,
    completeOnboarding: store.completeOnboarding,
    exportUserData: store.exportUserData,
    importUserData: store.importUserData,
    clearAllData: store.clearAllData,
  };
};