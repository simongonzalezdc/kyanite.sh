'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';

// Accessibility settings interface
interface AccessibilitySettings {
  // Screen reader settings
  screenReaderEnabled: boolean;
  announcements: string[];
  
  // Visual settings
  highContrastMode: boolean;
  fontSize: 'small' | 'medium' | 'large' | 'extra-large';
  reducedMotion: boolean;
  
  // Keyboard settings
  keyboardNavigation: boolean;
  focusVisible: boolean;
  
  // Audio settings
  audioDescriptions: boolean;
  visualIndicators: boolean;
  vibrationFeedback: boolean;
}

// Accessibility context interface
interface AccessibilityContextType {
  settings: AccessibilitySettings;
  updateSettings: (updates: Partial<AccessibilitySettings>) => void;
  announce: (message: string, priority?: 'polite' | 'assertive') => void;
  clearAnnouncements: () => void;
  detectScreenReader: () => boolean;
  detectReducedMotion: () => boolean;
  detectHighContrast: () => boolean;
  setFocus: (element: HTMLElement | null) => void;
  trapFocus: (container: HTMLElement | null) => () => void;
}

// Default settings
const defaultSettings: AccessibilitySettings = {
  screenReaderEnabled: false,
  announcements: [],
  highContrastMode: false,
  fontSize: 'medium',
  reducedMotion: false,
  keyboardNavigation: true,
  focusVisible: true,
  audioDescriptions: true,
  visualIndicators: true,
  vibrationFeedback: true,
};

// Create context
const AccessibilityContext = createContext<AccessibilityContextType | undefined>(undefined);

// Provider component
export function AccessibilityProvider({ children }: { children: ReactNode }) {
  const [settings, setSettings] = useState<AccessibilitySettings>(defaultSettings);
  const [previousFocus, setPreviousFocus] = useState<HTMLElement | null>(null);

  // Detect screen reader on mount
  useEffect(() => {
    const hasScreenReader = detectScreenReader();
    const prefersReducedMotion = detectReducedMotion();
    const prefersHighContrast = detectHighContrast();
    
    setSettings(prev => ({
      ...prev,
      screenReaderEnabled: hasScreenReader,
      reducedMotion: prefersReducedMotion,
      highContrastMode: prefersHighContrast,
    }));

    // Listen for system preference changes
    const mediaQueryReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    const mediaQueryHighContrast = window.matchMedia('(prefers-contrast: high)');
    
    const handleReducedMotionChange = (e: MediaQueryListEvent) => {
      setSettings(prev => ({ ...prev, reducedMotion: e.matches }));
    };
    
    const handleHighContrastChange = (e: MediaQueryListEvent) => {
      setSettings(prev => ({ ...prev, highContrastMode: e.matches }));
    };
    
    mediaQueryReducedMotion.addEventListener('change', handleReducedMotionChange);
    mediaQueryHighContrast.addEventListener('change', handleHighContrastChange);
    
    return () => {
      mediaQueryReducedMotion.removeEventListener('change', handleReducedMotionChange);
      mediaQueryHighContrast.removeEventListener('change', handleHighContrastChange);
    };
  }, []);

  // Apply settings to document
  useEffect(() => {
    const root = document.documentElement;
    
    // Apply font size
    root.setAttribute('data-font-size', settings.fontSize);
    
    // Apply high contrast
    root.setAttribute('data-high-contrast', settings.highContrastMode.toString());
    
    // Apply reduced motion
    root.setAttribute('data-reduced-motion', settings.reducedMotion.toString());
    
    // Apply focus visible
    root.setAttribute('data-focus-visible', settings.focusVisible.toString());
    
    // Apply screen reader optimizations
    root.setAttribute('data-screen-reader', settings.screenReaderEnabled.toString());
  }, [settings]);

  // Detect screen reader using multiple methods
  function detectScreenReader(): boolean {
    // Method 1: Check for screen reader specific media queries
    const hasScreenReaderMediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    
    // Method 2: Check for common screen reader browser extensions
    const hasScreenReaderExtension = 
      document.body.getAttribute('aria-live') !== null ||
      document.querySelector('[role="log"]') !== null;
    
    // Method 3: Check for voiceOver on iOS/Mac
    const hasVoiceOver = window.speechSynthesis !== undefined && 
      window.speechSynthesis.getVoices().length > 0;
    
    return hasScreenReaderMediaQuery || hasScreenReaderExtension || hasVoiceOver;
  }

  // Detect reduced motion preference
  function detectReducedMotion(): boolean {
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches;
  }

  // Detect high contrast preference
  function detectHighContrast(): boolean {
    return window.matchMedia('(prefers-contrast: high)').matches;
  }

  // Update settings
  function updateSettings(updates: Partial<AccessibilitySettings>) {
    setSettings(prev => ({ ...prev, ...updates }));
  }

  // Announce message to screen readers
  function announce(message: string, priority: 'polite' | 'assertive' = 'polite') {
    const announcementId = `announcement-${Date.now()}`;
    
    // Create live region element
    const liveRegion = document.createElement('div');
    liveRegion.id = announcementId;
    liveRegion.setAttribute('aria-live', priority);
    liveRegion.setAttribute('aria-atomic', 'true');
    liveRegion.className = 'sr-only';
    liveRegion.textContent = message;
    
    document.body.appendChild(liveRegion);
    
    // Update announcements in state
    setSettings(prev => ({
      ...prev,
      announcements: [...prev.announcements, message]
    }));
    
    // Remove after announcement is read
    setTimeout(() => {
      document.body.removeChild(liveRegion);
    }, 1000);
  }

  // Clear announcements
  function clearAnnouncements() {
    setSettings(prev => ({ ...prev, announcements: [] }));
    
    // Remove all live regions
    const liveRegions = document.querySelectorAll('[aria-live]');
    liveRegions.forEach(region => {
      if (region.id && region.id.startsWith('announcement-')) {
        region.remove();
      }
    });
  }

  // Set focus to element
  function setFocus(element: HTMLElement | null) {
    if (element) {
      // Store previous focus for restoration
      setPreviousFocus(document.activeElement as HTMLElement);
      element.focus();
    }
  }

  // Trap focus within container
  function trapFocus(container: HTMLElement | null) {
    if (!container) return () => {};
    
    const focusableElements = container.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    ) as NodeListOf<HTMLElement>;
    
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];
    
    const handleTabKey = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return;
      
      if (e.shiftKey) {
        if (document.activeElement === firstElement) {
          lastElement?.focus();
          e.preventDefault();
        }
      } else {
        if (document.activeElement === lastElement) {
          firstElement?.focus();
          e.preventDefault();
        }
      }
    };
    
    container.addEventListener('keydown', handleTabKey);
    
    // Set initial focus
    firstElement?.focus();
    
    // Return cleanup function
    return () => {
      container.removeEventListener('keydown', handleTabKey);
      previousFocus?.focus();
    };
  }

  const contextValue: AccessibilityContextType = {
    settings,
    updateSettings,
    announce,
    clearAnnouncements,
    detectScreenReader,
    detectReducedMotion,
    detectHighContrast,
    setFocus,
    trapFocus,
  };

  return (
    <AccessibilityContext.Provider value={contextValue}>
      {children}
    </AccessibilityContext.Provider>
  );
}

// Hook to use accessibility context
export function useAccessibility() {
  const context = useContext(AccessibilityContext);
  if (context === undefined) {
    throw new Error('useAccessibility must be used within an AccessibilityProvider');
  }
  return context;
}

// Export for testing
export { AccessibilityContext };