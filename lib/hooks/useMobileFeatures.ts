'use client';

import { useState, useEffect, useCallback } from 'react';
import { 
  isMobile, 
  isTouchDevice, 
  getOrientation, 
  enableKeyboardAvoidance, 
  enablePullToRefresh,
  getViewportHeight,
  getViewportWidth
} from '@/lib/utils';

interface MobileFeaturesState {
  isMobile: boolean;
  isTouchDevice: boolean;
  orientation: 'portrait' | 'landscape';
  viewportHeight: number;
  viewportWidth: number;
  keyboardVisible: boolean;
  safeAreaInsets: {
    top: number;
    right: number;
    bottom: number;
    left: number;
  };
}

interface MobileFeaturesActions {
  triggerHaptic: (type: 'light' | 'medium' | 'heavy') => void;
  enablePullToRefresh: (callback: () => Promise<void> | void) => () => void;
  scrollToTop: () => void;
  scrollToSection: (sectionId: string) => void;
}

export function useMobileFeatures(): MobileFeaturesState & MobileFeaturesActions {
  const [state, setState] = useState<MobileFeaturesState>({
    isMobile: false,
    isTouchDevice: false,
    orientation: 'portrait',
    viewportHeight: 0,
    viewportWidth: 0,
    keyboardVisible: false,
    safeAreaInsets: { top: 0, right: 0, bottom: 0, left: 0 }
  });

  // Update mobile state on mount and resize
  useEffect(() => {
    const updateState = () => {
      setState(prev => ({
        ...prev,
        isMobile: isMobile(),
        isTouchDevice: isTouchDevice(),
        orientation: getOrientation(),
        viewportHeight: getViewportHeight(),
        viewportWidth: getViewportWidth(),
        safeAreaInsets: getSafeAreaInsets()
      }));
    };

    updateState();

    // Listen for orientation changes
    const handleOrientationChange = () => {
      setTimeout(updateState, 100); // Delay to allow for orientation change completion
    };

    window.addEventListener('resize', updateState);
    window.addEventListener('orientationchange', handleOrientationChange);

    return () => {
      window.removeEventListener('resize', updateState);
      window.removeEventListener('orientationchange', handleOrientationChange);
    };
  }, []);

  // Keyboard avoidance
  useEffect(() => {
    if (state.isMobile) {
      const cleanup = enableKeyboardAvoidance();
      return cleanup;
    }
  }, [state.isMobile]);

  // Haptic feedback
  const triggerHaptic = useCallback((type: 'light' | 'medium' | 'heavy') => {
    if (state.isTouchDevice && 'vibrate' in navigator) {
      const patterns = {
        light: [10],
        medium: [20],
        heavy: [40]
      };
      navigator.vibrate(patterns[type]);
    }
  }, [state.isTouchDevice]);

  // Pull-to-refresh
  const enablePullToRefreshFeature = useCallback((callback: () => Promise<void> | void) => {
    if (state.isTouchDevice) {
      return enablePullToRefresh(callback) || (() => {});
    }
    return () => {};
  }, [state.isTouchDevice]);

  // Scroll utilities
  const scrollToTop = useCallback(() => {
    window.scrollTo({
      top: 0,
      left: 0,
      behavior: 'smooth'
    });
  }, []);

  const scrollToSection = useCallback((sectionId: string) => {
    const element = document.getElementById(sectionId);
    if (element) {
      const offset = state.isMobile ? 80 : 100; // Account for mobile nav
      const elementPosition = element.getBoundingClientRect().top + window.scrollY;
      const offsetPosition = elementPosition - offset;

      window.scrollTo({
        top: offsetPosition,
        behavior: 'smooth'
      });
    }
  }, [state.isMobile]);

  return {
    ...state,
    triggerHaptic,
    enablePullToRefresh: enablePullToRefreshFeature,
    scrollToTop,
    scrollToSection
  };
}

// Helper function to get safe area insets
function getSafeAreaInsets() {
  if (typeof window === 'undefined') return { top: 0, right: 0, bottom: 0, left: 0 };
  
  const style = getComputedStyle(document.documentElement);
  return {
    top: parseInt(style.getPropertyValue('--safe-area-inset-top')) || 0,
    right: parseInt(style.getPropertyValue('--safe-area-inset-right')) || 0,
    bottom: parseInt(style.getPropertyValue('--safe-area-inset-bottom')) || 0,
    left: parseInt(style.getPropertyValue('--safe-area-inset-left')) || 0,
  };
}

// Hook for managing mobile navigation state
export function useMobileNavigation() {
  const [activeSection, setActiveSection] = useState('record');
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);

  const handleScroll = useCallback(() => {
    const currentScrollY = window.scrollY;
    
    // Show navigation when scrolling up or at top
    if (currentScrollY < lastScrollY || currentScrollY < 100) {
      setIsVisible(true);
    } 
    // Hide navigation when scrolling down
    else if (currentScrollY > lastScrollY && currentScrollY > 100) {
      setIsVisible(false);
    }
    
    setLastScrollY(currentScrollY);
  }, [lastScrollY]);

  useEffect(() => {
    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [handleScroll]);

  return {
    activeSection,
    setActiveSection,
    isVisible,
    setIsVisible
  };
}

// Hook for managing touch gestures
export function useTouchGestures(
  elementRef: React.RefObject<HTMLElement>,
  handlers: {
    onSwipeLeft?: () => void;
    onSwipeRight?: () => void;
    onSwipeUp?: () => void;
    onSwipeDown?: () => void;
    onTap?: () => void;
    onLongPress?: () => void;
  }
) {
  const [touchStart, setTouchStart] = useState({ x: 0, y: 0, time: 0 });
  const [longPressTimer, setLongPressTimer] = useState<NodeJS.Timeout | null>(null);

  useEffect(() => {
    const element = elementRef.current;
    if (!element || !isTouchDevice()) return;

    const threshold = 50; // Minimum distance for swipe
    const longPressDelay = 500; // Long press delay in ms

    const handleTouchStart = (e: TouchEvent) => {
      const touch = e.touches[0];
      setTouchStart({
        x: touch.clientX,
        y: touch.clientY,
        time: Date.now()
      });

      // Start long press timer
      const timer = setTimeout(() => {
        handlers.onLongPress?.();
      }, longPressDelay);
      setLongPressTimer(timer);
    };

    const handleTouchEnd = (e: TouchEvent) => {
      // Clear long press timer
      if (longPressTimer) {
        clearTimeout(longPressTimer);
        setLongPressTimer(null);
      }

      const touch = e.changedTouches[0];
      const deltaX = touch.clientX - touchStart.x;
      const deltaY = touch.clientY - touchStart.y;
      const deltaTime = Date.now() - touchStart.time;

      // Check for tap (quick touch with minimal movement)
      if (Math.abs(deltaX) < 10 && Math.abs(deltaY) < 10 && deltaTime < 200) {
        handlers.onTap?.();
        return;
      }

      // Check for swipe gestures
      if (Math.abs(deltaX) > Math.abs(deltaY)) {
        // Horizontal swipe
        if (Math.abs(deltaX) > threshold) {
          if (deltaX > 0) {
            handlers.onSwipeRight?.();
          } else {
            handlers.onSwipeLeft?.();
          }
        }
      } else {
        // Vertical swipe
        if (Math.abs(deltaY) > threshold) {
          if (deltaY > 0) {
            handlers.onSwipeDown?.();
          } else {
            handlers.onSwipeUp?.();
          }
        }
      }
    };

    element.addEventListener('touchstart', handleTouchStart, { passive: true });
    element.addEventListener('touchend', handleTouchEnd, { passive: true });

    return () => {
      element.removeEventListener('touchstart', handleTouchStart);
      element.removeEventListener('touchend', handleTouchEnd);
      if (longPressTimer) {
        clearTimeout(longPressTimer);
      }
    };
  }, [elementRef, handlers, touchStart, longPressTimer]);
}

// Hook for managing responsive breakpoints
export function useResponsiveBreakpoints() {
  const [breakpoint, setBreakpoint] = useState<'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'>('md');

  useEffect(() => {
    const updateBreakpoint = () => {
      const width = window.innerWidth;
      if (width < 480) setBreakpoint('xs');
      else if (width < 768) setBreakpoint('sm');
      else if (width < 1024) setBreakpoint('md');
      else if (width < 1280) setBreakpoint('lg');
      else if (width < 1536) setBreakpoint('xl');
      else setBreakpoint('2xl');
    };

    updateBreakpoint();
    window.addEventListener('resize', updateBreakpoint);
    return () => window.removeEventListener('resize', updateBreakpoint);
  }, []);

  return {
    breakpoint,
    isMobile: breakpoint === 'xs' || breakpoint === 'sm',
    isTablet: breakpoint === 'md',
    isDesktop: breakpoint === 'lg' || breakpoint === 'xl' || breakpoint === '2xl',
    isSmallScreen: breakpoint === 'xs' || breakpoint === 'sm' || breakpoint === 'md'
  };
}