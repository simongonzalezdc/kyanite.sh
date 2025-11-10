'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import { useAccessibility } from './AccessibilityProvider';

// Keyboard shortcut interface
interface KeyboardShortcut {
  key: string;
  ctrlKey?: boolean;
  shiftKey?: boolean;
  altKey?: boolean;
  metaKey?: boolean;
  description: string;
  action: () => void;
  category: 'navigation' | 'recording' | 'playback' | 'general' | 'accessibility';
}

// Skip link interface
interface SkipLink {
  id: string;
  label: string;
  target: string;
}

// Keyboard navigation component
export function KeyboardNavigation({ children }: { children: React.ReactNode }) {
  const { announce, settings } = useAccessibility();
  const [shortcuts, setShortcuts] = useState<KeyboardShortcut[]>([]);
  const [skipLinks, setSkipLinks] = useState<SkipLink[]>([]);
  const [isHelpVisible, setIsHelpVisible] = useState(false);
  const shortcutsRef = useRef<KeyboardShortcut[]>([]);

  // Register keyboard shortcuts
  const registerShortcut = useCallback((shortcut: KeyboardShortcut) => {
    shortcutsRef.current = [...shortcutsRef.current, shortcut];
    setShortcuts([...shortcutsRef.current]);
  }, []);

  // Unregister keyboard shortcuts
  const unregisterShortcut = useCallback((key: string, modifiers: Partial<KeyboardShortcut> = {}) => {
    shortcutsRef.current = shortcutsRef.current.filter(
      s => !(s.key === key && 
            s.ctrlKey === modifiers.ctrlKey && 
            s.shiftKey === modifiers.shiftKey && 
            s.altKey === modifiers.altKey && 
            s.metaKey === modifiers.metaKey)
    );
    setShortcuts([...shortcutsRef.current]);
  }, []);

  // Register skip links
  const registerSkipLink = useCallback((skipLink: SkipLink) => {
    setSkipLinks(prev => [...prev, skipLink]);
  }, []);

  // Handle keyboard events
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      // Don't handle keyboard events when user is typing in input fields
      const target = event.target as HTMLElement;
      if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.contentEditable === 'true') {
        return;
      }

      // Check for matching shortcuts
      const matchingShortcut = shortcutsRef.current.find(shortcut => {
        return (
          shortcut.key.toLowerCase() === event.key.toLowerCase() &&
          !!shortcut.ctrlKey === event.ctrlKey &&
          !!shortcut.shiftKey === event.shiftKey &&
          !!shortcut.altKey === event.altKey &&
          !!shortcut.metaKey === event.metaKey
        );
      });

      if (matchingShortcut) {
        event.preventDefault();
        matchingShortcut.action();
        
        // Announce action to screen readers
        if (settings.screenReaderEnabled) {
          announce(`Executed: ${matchingShortcut.description}`);
        }
      }

      // Toggle help with ? key
      if (event.key === '?' && !event.ctrlKey && !event.shiftKey && !event.altKey && !event.metaKey) {
        event.preventDefault();
        setIsHelpVisible(prev => !prev);
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [announce, settings.screenReaderEnabled]);

  // Focus management for skip links
  const handleSkipLink = (targetId: string) => {
    const targetElement = document.getElementById(targetId);
    if (targetElement) {
      targetElement.focus();
      targetElement.scrollIntoView({ behavior: 'smooth' });
      
      if (settings.screenReaderEnabled) {
        announce(`Navigated to ${targetElement.getAttribute('aria-label') || targetElement.textContent}`);
      }
    }
  };

  return (
    <>
      {/* Skip Links */}
      <div className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 z-50 bg-gray-900 border border-gray-700 rounded-lg p-4">
        {skipLinks.map((link) => (
          <button
            key={link.id}
            onClick={() => handleSkipLink(link.target)}
            className="block w-full text-left px-4 py-2 text-white hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-primary-500"
          >
            {link.label}
          </button>
        ))}
      </div>

      {/* Keyboard Shortcuts Help Modal */}
      {isHelpVisible && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-75"
          role="dialog"
          aria-modal="true"
          aria-labelledby="keyboard-help-title"
        >
          <div className="bg-gray-900 border border-gray-700 rounded-lg p-6 max-w-2xl max-h-[80vh] overflow-y-auto">
            <h2 id="keyboard-help-title" className="text-xl font-bold mb-4">Keyboard Shortcuts</h2>
            
            <div className="space-y-4">
              {['navigation', 'recording', 'playback', 'general', 'accessibility'].map(category => (
                <div key={category}>
                  <h3 className="text-lg font-semibold mb-2 capitalize">{category}</h3>
                  <div className="space-y-1">
                    {shortcuts
                      .filter(shortcut => shortcut.category === category)
                      .map((shortcut, index) => (
                        <div key={index} className="flex justify-between items-center py-1">
                          <span className="text-gray-300">{shortcut.description}</span>
                          <kbd className="px-2 py-1 bg-gray-800 border border-gray-600 rounded text-sm">
                            {[
                              shortcut.ctrlKey && 'Ctrl',
                              shortcut.shiftKey && 'Shift',
                              shortcut.altKey && 'Alt',
                              shortcut.metaKey && 'Cmd',
                              shortcut.key.toUpperCase()
                            ].filter(Boolean).join(' + ')}
                          </kbd>
                        </div>
                      ))}
                  </div>
                </div>
              ))}
            </div>
            
            <div className="mt-6 flex justify-end">
              <button
                onClick={() => setIsHelpVisible(false)}
                className="px-4 py-2 bg-primary-500 hover:bg-primary-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Provide context for child components */}
      <KeyboardNavigationContext.Provider value={{
        registerShortcut,
        unregisterShortcut,
        registerSkipLink,
        shortcuts,
        skipLinks
      }}>
        {children}
      </KeyboardNavigationContext.Provider>
    </>
  );
}

// Context for keyboard navigation
const KeyboardNavigationContext = React.createContext<{
  registerShortcut: (shortcut: KeyboardShortcut) => void;
  unregisterShortcut: (key: string, modifiers?: Partial<KeyboardShortcut>) => void;
  registerSkipLink: (skipLink: SkipLink) => void;
  shortcuts: KeyboardShortcut[];
  skipLinks: SkipLink[];
}>({
  registerShortcut: () => {},
  unregisterShortcut: () => {},
  registerSkipLink: () => {},
  shortcuts: [],
  skipLinks: []
});

// Hook to use keyboard navigation
export function useKeyboardNavigation() {
  const context = React.useContext(KeyboardNavigationContext);
  if (!context) {
    throw new Error('useKeyboardNavigation must be used within a KeyboardNavigation component');
  }
  return context;
}

// Focus trap component for modals and dialogs
export function FocusTrap({ children, isActive = true }: { children: React.ReactNode; isActive?: boolean }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const previousFocusRef = useRef<HTMLElement | null>(null);

  useEffect(() => {
    if (!isActive || !containerRef.current) return;

    // Store current focus
    previousFocusRef.current = document.activeElement as HTMLElement;

    // Get focusable elements
    const focusableElements = containerRef.current.querySelectorAll(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    ) as NodeListOf<HTMLElement>;

    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    // Focus first element
    firstElement?.focus();

    // Handle tab key
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

    containerRef.current.addEventListener('keydown', handleTabKey);

    return () => {
      containerRef.current?.removeEventListener('keydown', handleTabKey);
      // Restore previous focus
      previousFocusRef.current?.focus();
    };
  }, [isActive]);

  return (
    <div ref={containerRef} tabIndex={-1}>
      {children}
    </div>
  );
}

// Focus management hook
export function useFocusManagement() {
  const { announce, settings } = useAccessibility();

  const setFocus = useCallback((element: HTMLElement | null) => {
    if (element) {
      element.focus();
      
      // Announce focus change to screen readers
      const label = element.getAttribute('aria-label') ||
                   element.getAttribute('title') ||
                   element.textContent;
      if (label && settings.screenReaderEnabled) {
        announce(`Focused on ${label}`);
      }
    }
  }, [announce, settings.screenReaderEnabled]);

  const restoreFocus = useCallback(() => {
    const previousFocus = document.querySelector('[data-previous-focus]') as HTMLElement;
    if (previousFocus) {
      previousFocus.focus();
    }
  }, []);

  const saveFocus = useCallback(() => {
    const activeElement = document.activeElement as HTMLElement;
    if (activeElement) {
      activeElement.setAttribute('data-previous-focus', 'true');
    }
  }, []);

  return { setFocus, restoreFocus, saveFocus };
}

// Visual focus indicator component
export function FocusIndicator({ children }: { children: React.ReactNode }) {
  const { settings } = useAccessibility();

  if (!settings.focusVisible) {
    return <>{children}</>;
  }

  return (
    <div className="focus-within:ring-2 focus-within:ring-primary-500 focus-within:ring-offset-2 focus-within:ring-offset-gray-900">
      {children}
    </div>
  );
}