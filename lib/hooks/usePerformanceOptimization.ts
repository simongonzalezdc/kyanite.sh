'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import { isMobile, isTouchDevice, debounce } from '@/lib/utils';

interface PerformanceState {
  isLowEndDevice: boolean;
  reducedMotion: boolean;
  batteryLevel: number | null;
  connectionType: string;
  memoryPressure: 'low' | 'medium' | 'high';
}

interface PerformanceActions {
  shouldReduceAnimations: () => boolean;
  shouldLazyLoad: () => boolean;
  shouldReduceQuality: () => boolean;
  optimizeForDevice: <T>(highEnd: T, lowEnd: T) => T;
  measurePerformance: (name: string, fn: () => void) => number;
}

export function usePerformanceOptimization(): PerformanceState & PerformanceActions {
  const [state, setState] = useState<PerformanceState>({
    isLowEndDevice: false,
    reducedMotion: false,
    batteryLevel: null,
    connectionType: 'unknown',
    memoryPressure: 'medium'
  });

  const performanceMetrics = useRef<Map<string, number[]>>(new Map());

  // Detect device capabilities on mount
  useEffect(() => {
    const detectDeviceCapabilities = () => {
      const isLowEnd = detectLowEndDevice();
      const prefersReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
      
      setState(prev => ({
        ...prev,
        isLowEndDevice: isLowEnd,
        reducedMotion: prefersReducedMotion
      }));
      
      // Monitor battery level if available
      if ('getBattery' in navigator) {
        (navigator as any).getBattery().then((battery: any) => {
          setState(prev => ({ ...prev, batteryLevel: battery.level }));
          
          battery.addEventListener('levelchange', () => {
            setState(prev => ({ ...prev, batteryLevel: battery.level }));
          });
        });
      }
      
      // Monitor network connection
      if ('connection' in navigator) {
        const connection = (navigator as any).connection;
        setState(prev => ({ ...prev, connectionType: connection.effectiveType || 'unknown' }));
        
        connection.addEventListener('change', () => {
          setState(prev => ({ ...prev, connectionType: connection.effectiveType || 'unknown' }));
        });
      }
      
      // Monitor memory pressure if available
      if ('memory' in performance) {
        const memory = (performance as any).memory;
        const usedRatio = memory.usedJSHeapSize / memory.jsHeapSizeLimit;
        
        let pressure: 'low' | 'medium' | 'high' = 'medium';
        if (usedRatio < 0.5) pressure = 'low';
        else if (usedRatio > 0.8) pressure = 'high';
        
        setState(prev => ({ ...prev, memoryPressure: pressure }));
      }
    };

    detectDeviceCapabilities();
    
    // Listen for reduced motion preference changes
    const mediaQuery = window.matchMedia('(prefers-reduced-motion: reduce)');
    const handleChange = (e: MediaQueryListEvent) => {
      setState(prev => ({ ...prev, reducedMotion: e.matches }));
    };
    
    mediaQuery.addEventListener('change', handleChange);
    
    return () => {
      mediaQuery.removeEventListener('change', handleChange);
    };
  }, []); // This is already correct, just confirming

  // Detect low-end device based on various factors
  const detectLowEndDevice = useCallback(() => {
    // Check hardware concurrency (CPU cores)
    const cpuCores = navigator.hardwareConcurrency || 4;
    if (cpuCores <= 2) return true;

    // Check device memory if available
    if ('deviceMemory' in navigator) {
      const memory = (navigator as any).deviceMemory;
      if (memory <= 2) return true; // Less than 2GB RAM
    }

    // Check for mobile devices with limited capabilities
    if (isMobile()) {
      // Simple heuristic for older mobile devices
      const userAgent = navigator.userAgent;
      const isOldDevice = /Android [1-4]/.test(userAgent) || 
                          /iPhone OS [1-9]/.test(userAgent) ||
                          /iPad; CPU OS [1-9]/.test(userAgent);
      if (isOldDevice) return true;
    }

    return false;
  }, []);

  // Performance decision helpers
  const shouldReduceAnimations = useCallback(() => {
    return state.reducedMotion || state.isLowEndDevice || state.batteryLevel !== null && state.batteryLevel < 0.2;
  }, [state.reducedMotion, state.isLowEndDevice, state.batteryLevel]);

  const shouldLazyLoad = useCallback(() => {
    return state.isLowEndDevice || 
           state.connectionType === 'slow-2g' || 
           state.connectionType === '2g' ||
           state.memoryPressure === 'high';
  }, [state.isLowEndDevice, state.connectionType, state.memoryPressure]);

  const shouldReduceQuality = useCallback(() => {
    return state.isLowEndDevice || 
           state.connectionType === 'slow-2g' || 
           state.connectionType === '2g' ||
           state.memoryPressure === 'high' ||
           (state.batteryLevel !== null && state.batteryLevel < 0.1);
  }, [state.isLowEndDevice, state.connectionType, state.memoryPressure, state.batteryLevel]);

  // Device-optimized value selection
  const optimizeForDevice = useCallback(<T,>(highEnd: T, lowEnd: T): T => {
    return state.isLowEndDevice ? lowEnd : highEnd;
  }, [state.isLowEndDevice]);

  // Performance measurement
  const measurePerformance = useCallback((name: string, fn: () => void): number => {
    if (typeof window === 'undefined' || !('performance' in window)) {
      fn();
      return 0;
    }
    
    const start = performance.now();
    fn();
    const end = performance.now();
    const duration = end - start;
    
    // Store metrics for analysis
    if (!performanceMetrics.current.has(name)) {
      performanceMetrics.current.set(name, []);
    }
    const metrics = performanceMetrics.current.get(name)!;
    metrics.push(duration);
    
    // Keep only last 10 measurements
    if (metrics.length > 10) {
      metrics.shift();
    }
    
    // Log performance warnings
    const average = metrics.reduce((a, b) => a + b, 0) / metrics.length;
    if (average > 100) { // 100ms threshold
      console.warn(`Slow operation detected: ${name} took ${average.toFixed(2)}ms on average`);
    }
    
    return duration;
  }, []);

  return {
    ...state,
    shouldReduceAnimations,
    shouldLazyLoad,
    shouldReduceQuality,
    optimizeForDevice,
    measurePerformance
  };
}

// Hook for lazy loading components
export function useLazyLoad(threshold = 0.1) {
  const [isVisible, setIsVisible] = useState(false);
  const elementRef = useRef<HTMLElement>(null);

  useEffect(() => {
    const element = elementRef.current;
    if (!element) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.disconnect();
        }
      },
      {
        threshold,
        rootMargin: '50px' // Start loading 50px before element comes into view
      }
    );

    observer.observe(element);

    return () => observer.disconnect();
  }, [threshold]);

  return { elementRef, isVisible };
}

// Hook for optimizing animations
export function useOptimizedAnimation(
  trigger: boolean,
  options: {
    duration?: number;
    reducedMotionDuration?: number;
    disabled?: boolean;
  } = {}
) {
  const { shouldReduceAnimations } = usePerformanceOptimization();
  const [isAnimating, setIsAnimating] = useState(false);

  const {
    duration = 300,
    reducedMotionDuration = 0,
    disabled = false
  } = options;

  const effectiveDuration = shouldReduceAnimations() ? reducedMotionDuration : duration;
  const shouldAnimate = !disabled && !shouldReduceAnimations() && effectiveDuration > 0;

  useEffect(() => {
    if (trigger && shouldAnimate) {
      setIsAnimating(true);
      const timer = setTimeout(() => {
        setIsAnimating(false);
      }, effectiveDuration);

      return () => clearTimeout(timer);
    }
  }, [trigger, shouldAnimate, effectiveDuration]);

  return {
    isAnimating,
    shouldAnimate,
    duration: effectiveDuration
  };
}

// Hook for managing resource loading
export function useResourceLoader() {
  const { shouldLazyLoad, shouldReduceQuality } = usePerformanceOptimization();
  const [loadedResources, setLoadedResources] = useState<Set<string>>(new Set());

  const loadResource = useCallback(async (
    url: string, 
    type: 'image' | 'audio' | 'video' | 'data',
    priority: 'high' | 'low' = 'low'
  ): Promise<any> => {
    // Skip if already loaded
    if (loadedResources.has(url)) {
      return Promise.resolve();
    }

    // Skip low priority resources if lazy loading is enabled
    if (shouldLazyLoad() && priority === 'low') {
      return Promise.resolve();
    }

    try {
      let resource;
      
      switch (type) {
        case 'image':
          resource = new Image();
          if (shouldReduceQuality()) {
            // Add quality parameter for images
            const separator = url.includes('?') ? '&' : '?';
            resource.src = `${url}${separator}quality=low`;
          } else {
            resource.src = url;
          }
          break;
          
        case 'audio':
        case 'video':
          resource = document.createElement(type);
          resource.src = url;
          break;
          
        case 'data':
          resource = await fetch(url).then(res => res.json());
          break;
          
        default:
          throw new Error(`Unknown resource type: ${type}`);
      }

      // Wait for resource to load
      if (type !== 'data' && resource instanceof HTMLElement) {
        await new Promise((resolve, reject) => {
          resource.onload = resolve;
          resource.onerror = reject;
        });
      }

      setLoadedResources(prev => new Set([...prev, url]));
      return resource;
    } catch (error) {
      console.error(`Failed to load resource: ${url}`, error);
      throw error;
    }
  }, [loadedResources, shouldLazyLoad, shouldReduceQuality]);

  return { loadResource, loadedResources };
}

// Hook for throttling scroll events on mobile
export function useThrottledScroll(callback: () => void, delay = 16) {
  const lastScrollTime = useRef(0);
  const ticking = useRef(false);

  useEffect(() => {
    const handleScroll = () => {
      lastScrollTime.current = performance.now();
      
      if (!ticking.current) {
        requestAnimationFrame(() => {
          callback();
          ticking.current = false;
        });
        ticking.current = true;
      }
    };

    window.addEventListener('scroll', handleScroll, { passive: true });
    return () => window.removeEventListener('scroll', handleScroll);
  }, [callback]);
}