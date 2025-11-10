'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import { useAccessibility } from './AccessibilityProvider';

// Live region component for dynamic content announcements
interface LiveRegionProps {
  politeness?: 'polite' | 'assertive' | 'off';
  children?: React.ReactNode;
  className?: string;
}

export function LiveRegion({ politeness = 'polite', children, className = '' }: LiveRegionProps) {
  const [content, setContent] = useState('');
  const [isAnnouncing, setIsAnnouncing] = useState(false);
  const timeoutRef = useRef<NodeJS.Timeout | null>(null);

  const announce = useCallback((message: string) => {
    setContent(message);
    setIsAnnouncing(true);
    
    // Clear previous timeout
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    
    // Reset after announcement
    timeoutRef.current = setTimeout(() => {
      setContent('');
      setIsAnnouncing(false);
    }, 1000);
  }, []);

  useEffect(() => {
    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  return (
    <>
      <div
        aria-live={politeness}
        aria-atomic="true"
        className={`sr-only ${className}`}
        role="status"
      >
        {content}
      </div>
      {children && React.isValidElement(children) &&
        React.cloneElement(children, { announce } as any)}
    </>
  );
}

// Screen reader only text component
export function ScreenReaderOnly({ children, className = '' }: { children: React.ReactNode; className?: string }) {
  return (
    <span className={`sr-only ${className}`}>
      {children}
    </span>
  );
}

// ARIA label provider component
interface AriaLabelProps {
  label: string;
  description?: string;
  children: React.ReactElement;
}

export function AriaLabel({ label, description, children }: AriaLabelProps) {
  const labelId = `label-${Math.random().toString(36).substr(2, 9)}`;
  const descriptionId = description ? `desc-${Math.random().toString(36).substr(2, 9)}` : undefined;

  if (!React.isValidElement(children)) {
    return children;
  }

  return React.cloneElement(children, {
    'aria-labelledby': labelId,
    'aria-describedby': descriptionId,
  } as any, [
    <span key="label" id={labelId} className="sr-only">
      {label}
    </span>,
    description && (
      <span key="description" id={descriptionId} className="sr-only">
        {description}
      </span>
    ),
    (children as any).props.children
  ]);
}

// Progress announcer for long-running operations
interface ProgressAnnouncerProps {
  isProcessing: boolean;
  progress?: number;
  message?: string;
  completionMessage?: string;
}

export function ProgressAnnouncer({ 
  isProcessing, 
  progress, 
  message = 'Processing', 
  completionMessage = 'Process completed' 
}: ProgressAnnouncerProps) {
  const { announce } = useAccessibility();
  const [lastProgress, setLastProgress] = useState(0);

  useEffect(() => {
    if (isProcessing && progress !== undefined) {
      // Only announce significant progress changes (every 10%)
      if (Math.floor(progress / 10) > Math.floor(lastProgress / 10)) {
        announce(`${message}: ${Math.round(progress)}% complete`);
        setLastProgress(progress);
      }
    } else if (!isProcessing && lastProgress > 0 && lastProgress < 100) {
      announce(completionMessage);
      setLastProgress(0);
    }
  }, [isProcessing, progress, lastProgress, message, completionMessage, announce]);

  return null;
}

// Screen reader navigation helper
export function ScreenReaderNavigation() {
  const { announce, settings } = useAccessibility();

  useEffect(() => {
    if (!settings.screenReaderEnabled) return;

    // Announce page load
    announce('VoxForge application loaded');

    // Announce route changes
    const handleRouteChange = () => {
      announce('Page changed');
    };

    // Listen for navigation events
    window.addEventListener('popstate', handleRouteChange);
    
    return () => {
      window.removeEventListener('popstate', handleRouteChange);
    };
  }, [announce, settings.screenReaderEnabled]);

  return null;
}

// Audio description component for visual content
interface AudioDescriptionProps {
  description: string;
  isPlaying?: boolean;
  children?: React.ReactNode;
}

export function AudioDescription({ description, isPlaying = false, children }: AudioDescriptionProps) {
  const { settings, announce } = useAccessibility();

  useEffect(() => {
    if (settings.audioDescriptions && isPlaying) {
      announce(description);
    }
  }, [isPlaying, description, settings.audioDescriptions, announce]);

  return <>{children}</>;
}

// Semantic landmark provider
export function LandmarkProvider({ children }: { children: React.ReactNode }) {
  return (
    <div role="application" aria-label="VoxForge Music Application">
      {children}
    </div>
  );
}

// Status announcer for application state changes
interface StatusAnnouncerProps {
  status: string;
  type?: 'info' | 'success' | 'warning' | 'error';
}

export function StatusAnnouncer({ status, type = 'info' }: StatusAnnouncerProps) {
  const { announce } = useAccessibility();
  const [previousStatus, setPreviousStatus] = useState('');

  useEffect(() => {
    if (status !== previousStatus) {
      const prefix = type === 'error' ? 'Error' : type === 'warning' ? 'Warning' : type === 'success' ? 'Success' : 'Information';
      announce(`${prefix}: ${status}`);
      setPreviousStatus(status);
    }
  }, [status, type, previousStatus, announce]);

  return null;
}

// Form validation announcer
interface ValidationAnnouncerProps {
  errors: string[];
  isValid: boolean;
}

export function ValidationAnnouncer({ errors, isValid }: ValidationAnnouncerProps) {
  const { announce } = useAccessibility();
  const [previousErrors, setPreviousErrors] = useState<string[]>([]);

  useEffect(() => {
    if (errors.length !== previousErrors.length || 
        !errors.every((error, index) => error === previousErrors[index])) {
      
      if (errors.length > 0) {
        announce(`Form has ${errors.length} error${errors.length > 1 ? 's' : ''}: ${errors.join(', ')}`);
      } else if (previousErrors.length > 0 && isValid) {
        announce('Form is now valid');
      }
      
      setPreviousErrors(errors);
    }
  }, [errors, isValid, previousErrors, announce]);

  return null;
}

// List announcer for dynamic lists
interface ListAnnouncerProps {
  items: any[];
  itemLabel: string;
  action?: 'added' | 'removed' | 'updated';
}

export function ListAnnouncer({ items, itemLabel, action }: ListAnnouncerProps) {
  const { announce } = useAccessibility();
  const [previousCount, setPreviousCount] = useState(0);

  useEffect(() => {
    const currentCount = items.length;
    
    if (currentCount !== previousCount) {
      if (action === 'added' && currentCount > previousCount) {
        announce(`Added ${itemLabel}. Total: ${currentCount}`);
      } else if (action === 'removed' && currentCount < previousCount) {
        announce(`Removed ${itemLabel}. Total: ${currentCount}`);
      } else {
        announce(`${itemLabel} list updated. Total: ${currentCount}`);
      }
      
      setPreviousCount(currentCount);
    }
  }, [items, itemLabel, action, previousCount, announce]);

  return null;
}

// Screen reader helper hook
export function useScreenReaderHelper() {
  const { announce, settings } = useAccessibility();

  const announceElement = useCallback((element: HTMLElement, context?: string) => {
    if (!settings.screenReaderEnabled) return;

    const label = element.getAttribute('aria-label') || 
                 element.getAttribute('title') || 
                 element.textContent || 
                 element.tagName.toLowerCase();
    
    const message = context ? `${context}: ${label}` : label;
    announce(message);
  }, [announce, settings.screenReaderEnabled]);

  const announceStateChange = useCallback((state: string, value: any) => {
    if (!settings.screenReaderEnabled) return;
    
    announce(`${state} changed to ${value}`);
  }, [announce, settings.screenReaderEnabled]);

  const announceNavigation = useCallback((destination: string) => {
    if (!settings.screenReaderEnabled) return;
    
    announce(`Navigated to ${destination}`);
  }, [announce, settings.screenReaderEnabled]);

  return {
    announceElement,
    announceStateChange,
    announceNavigation
  };
}

// Screen reader test component for development
export function ScreenReaderTest() {
  const { announce, settings } = useAccessibility();

  const testAnnouncement = (message: string) => {
    announce(`Test: ${message}`);
  };

  if (process.env.NODE_ENV !== 'development') {
    return null;
  }

  return (
    <div className="fixed bottom-4 right-4 p-4 bg-gray-800 border border-gray-700 rounded-lg z-50">
      <h3 className="text-sm font-semibold mb-2">Screen Reader Test</h3>
      <div className="space-y-2">
        <button
          onClick={() => testAnnouncement('This is a test announcement')}
          className="px-2 py-1 bg-blue-500 text-white rounded text-xs"
        >
          Test Announcement
        </button>
        <div className="text-xs text-gray-400">
          Screen Reader: {settings.screenReaderEnabled ? 'Enabled' : 'Disabled'}
        </div>
      </div>
    </div>
  );
}