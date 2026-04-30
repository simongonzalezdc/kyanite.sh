import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Mobile-specific utilities
export function isMobile() {
  if (typeof window === 'undefined') return false;
  return window.innerWidth < 768;
}

export function isTablet() {
  if (typeof window === 'undefined') return false;
  return window.innerWidth >= 768 && window.innerWidth < 1024;
}

export function isDesktop() {
  if (typeof window === 'undefined') return false;
  return window.innerWidth >= 1024;
}

// Touch detection
export function isTouchDevice() {
  if (typeof window === 'undefined') return false;
  return 'ontouchstart' in window || navigator.maxTouchPoints > 0;
}

// Performance utilities
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}

export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}

// Safe area utilities
export function getSafeAreaInsets() {
  if (typeof window === 'undefined') return { top: 0, right: 0, bottom: 0, left: 0 };
  
  const style = getComputedStyle(document.documentElement);
  return {
    top: parseInt(style.getPropertyValue('--safe-area-inset-top')) || 0,
    right: parseInt(style.getPropertyValue('--safe-area-inset-right')) || 0,
    bottom: parseInt(style.getPropertyValue('--safe-area-inset-bottom')) || 0,
    left: parseInt(style.getPropertyValue('--safe-area-inset-left')) || 0,
  };
}

// Viewport utilities
export function getViewportHeight() {
  if (typeof window === 'undefined') return 0;
  return window.innerHeight;
}

export function getViewportWidth() {
  if (typeof window === 'undefined') return 0;
  return window.innerWidth;
}

// Orientation utilities
export function getOrientation() {
  if (typeof window === 'undefined') return 'portrait';
  return window.innerHeight > window.innerWidth ? 'portrait' : 'landscape';
}

// Haptic feedback (if supported)
export function triggerHaptic(type: 'light' | 'medium' | 'heavy' = 'light') {
  if (typeof window === 'undefined' || !('vibrate' in navigator)) return;
  
  const patterns = {
    light: [10],
    medium: [20],
    heavy: [40]
  };
  
  navigator.vibrate(patterns[type]);
}

// Pull-to-refresh utilities
export function enablePullToRefresh(callback: () => Promise<void> | void) {
  if (typeof window === 'undefined' || !isTouchDevice()) return;
  
  let startY = 0;
  let isPulling = false;
  const threshold = 80;
  
  const handleTouchStart = (e: TouchEvent) => {
    if (window.scrollY === 0) {
      startY = e.touches[0].clientY;
      isPulling = true;
    }
  };
  
  const handleTouchMove = (e: TouchEvent) => {
    if (!isPulling) return;
    
    const currentY = e.touches[0].clientY;
    const diff = currentY - startY;
    
    if (diff > threshold && window.scrollY === 0) {
      e.preventDefault();
      // Visual feedback could be added here
    }
  };
  
  const handleTouchEnd = async (e: TouchEvent) => {
    if (!isPulling) return;
    
    const currentY = e.changedTouches[0].clientY;
    const diff = currentY - startY;
    
    if (diff > threshold && window.scrollY === 0) {
      await callback();
    }
    
    isPulling = false;
  };
  
  document.addEventListener('touchstart', handleTouchStart, { passive: true });
  document.addEventListener('touchmove', handleTouchMove, { passive: false });
  document.addEventListener('touchend', handleTouchEnd, { passive: true });
  
  return () => {
    document.removeEventListener('touchstart', handleTouchStart);
    document.removeEventListener('touchmove', handleTouchMove);
    document.removeEventListener('touchend', handleTouchEnd);
  };
}

// Keyboard avoidance utilities
export function enableKeyboardAvoidance() {
  if (typeof window === 'undefined' || !isMobile()) return;
  
  let originalViewportHeight = window.innerHeight;
  
  const handleVisualViewportChange = () => {
    if ('visualViewport' in window) {
      const viewport = window.visualViewport;
      const height = viewport?.height || window.innerHeight;
      
      // Adjust body padding when keyboard is shown
      if (height < originalViewportHeight * 0.8) {
        document.body.style.paddingBottom = `${originalViewportHeight - height}px`;
      } else {
        document.body.style.paddingBottom = '';
      }
    }
  };
  
  if ('visualViewport' in window) {
    window.visualViewport?.addEventListener('resize', handleVisualViewportChange);
    
    return () => {
      window.visualViewport?.removeEventListener('resize', handleVisualViewportChange);
    };
  }
  
  // Fallback for browsers without Visual Viewport API
  const handleFocusIn = (e: FocusEvent) => {
    const target = e.target as HTMLElement;
    if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA') {
      setTimeout(() => {
        if (window.innerHeight < originalViewportHeight * 0.8) {
          document.body.style.paddingBottom = `${originalViewportHeight - window.innerHeight}px`;
        }
      }, 300);
    }
  };
  
  const handleFocusOut = () => {
    document.body.style.paddingBottom = '';
  };
  
  document.addEventListener('focusin', handleFocusIn);
  document.addEventListener('focusout', handleFocusOut);
  
  return () => {
    document.removeEventListener('focusin', handleFocusIn);
    document.removeEventListener('focusout', handleFocusOut);
  };
}

// Performance monitoring utilities
export function measurePerformance(name: string, fn: () => void) {
  if (typeof window === 'undefined' || !('performance' in window)) {
    return fn();
  }
  
  const start = performance.now();
  fn();
  const end = performance.now();
  
  console.log(`${name} took ${end - start} milliseconds`);
  return end - start;
}

// Lazy loading utilities
export function createIntersectionObserver(
  callback: (entries: IntersectionObserverEntry[]) => void,
  options?: IntersectionObserverInit
) {
  if (typeof window === 'undefined' || !('IntersectionObserver' in window)) {
    return null;
  }
  
  return new IntersectionObserver(callback, {
    rootMargin: '50px',
    threshold: 0.1,
    ...options
  });
}

// Swipe gesture utilities
export function enableSwipeGestures(
  element: HTMLElement,
  handlers: {
    onSwipeLeft?: () => void;
    onSwipeRight?: () => void;
    onSwipeUp?: () => void;
    onSwipeDown?: () => void;
  }
) {
  let startX = 0;
  let startY = 0;
  let endX = 0;
  let endY = 0;
  
  const threshold = 50; // Minimum distance for swipe
  
  const handleTouchStart = (e: TouchEvent) => {
    startX = e.touches[0].clientX;
    startY = e.touches[0].clientY;
  };
  
  const handleTouchEnd = (e: TouchEvent) => {
    endX = e.changedTouches[0].clientX;
    endY = e.changedTouches[0].clientY;
    
    const diffX = endX - startX;
    const diffY = endY - startY;
    
    // Check if swipe is horizontal or vertical
    if (Math.abs(diffX) > Math.abs(diffY)) {
      // Horizontal swipe
      if (Math.abs(diffX) > threshold) {
        if (diffX > 0) {
          handlers.onSwipeRight?.();
        } else {
          handlers.onSwipeLeft?.();
        }
      }
    } else {
      // Vertical swipe
      if (Math.abs(diffY) > threshold) {
        if (diffY > 0) {
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
  };
}