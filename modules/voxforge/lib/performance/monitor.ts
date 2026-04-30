import { analytics } from '@/lib/analytics';
import { PerformanceMetrics } from '@/lib/types/analytics';

interface PerformanceEntry {
  name: string;
  startTime: number;
  duration?: number;
  type: 'audio' | 'render' | 'memory' | 'bundle';
}

class PerformanceMonitor {
  private entries: Map<string, PerformanceEntry> = new Map();
  private memoryObserver: PerformanceObserver | null = null;
  private renderObserver: PerformanceObserver | null = null;
  private isInitialized = false;
  private memoryCheckInterval: NodeJS.Timeout | null = null;

  /**
   * Initialize the performance monitor
   */
  initialize(): void {
    if (this.isInitialized) return;

    // Setup memory monitoring
    this.setupMemoryMonitoring();
    
    // Setup render performance monitoring
    this.setupRenderMonitoring();
    
    // Setup periodic memory checks
    this.setupPeriodicMemoryCheck();
    
    // Track initial bundle size
    this.trackBundleSize();
    
    this.isInitialized = true;
  }

  /**
   * Start timing a performance entry
   */
  start(name: string, type: PerformanceEntry['type'] = 'render'): void {
    this.entries.set(name, {
      name,
      startTime: performance.now(),
      type
    });
  }

  /**
   * End timing a performance entry and record it
   */
  end(name: string, context?: string): number {
    const entry = this.entries.get(name);
    if (!entry) {
      console.warn(`Performance: No start entry found for "${name}"`);
      return 0;
    }

    const duration = performance.now() - entry.startTime;
    
    // Track the performance metric
    switch (entry.type) {
      case 'audio':
        analytics.trackPerformance('audio_processing', duration, 'ms', context || name);
        break;
      case 'render':
        analytics.trackPerformance('render_time', duration, 'ms', context || name);
        break;
      default:
        analytics.trackPerformance('custom', duration, 'ms', context || name);
    }

    this.entries.delete(name);
    return duration;
  }

  /**
   * Measure audio processing performance
   */
  measureAudioProcessing<T>(
    name: string,
    fn: () => T | Promise<T>
  ): T | Promise<T> {
    this.start(name, 'audio');
    
    const result = fn();
    
    if (result instanceof Promise) {
      return result.finally(() => {
        this.end(name);
      });
    } else {
      this.end(name);
      return result;
    }
  }

  /**
   * Measure component render performance
   */
  measureRender<T>(
    name: string,
    fn: () => T | Promise<T>
  ): T | Promise<T> {
    this.start(name, 'render');
    
    const result = fn();
    
    if (result instanceof Promise) {
      return result.finally(() => {
        this.end(name);
      });
    } else {
      this.end(name);
      return result;
    }
  }

  /**
   * Setup memory monitoring using PerformanceObserver
   */
  private setupMemoryMonitoring(): void {
    if (typeof window === 'undefined' || !window.PerformanceObserver) return;

    try {
      this.memoryObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          if (entry.entryType === 'measure' && entry.name.includes('memory')) {
            const memoryEntry = entry as PerformanceMeasure;
            analytics.trackPerformance('memory_usage', memoryEntry.duration, 'bytes', memoryEntry.name);
          }
        }
      });
      
      this.memoryObserver.observe({ entryTypes: ['measure'] });
    } catch (error) {
      console.warn('Performance: Memory monitoring setup failed', error);
    }
  }

  /**
   * Setup render performance monitoring
   */
  private setupRenderMonitoring(): void {
    if (typeof window === 'undefined' || !window.PerformanceObserver) return;

    try {
      this.renderObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          if (entry.entryType === 'measure' && entry.name.includes('render')) {
            const renderEntry = entry as PerformanceMeasure;
            analytics.trackPerformance('render_time', renderEntry.duration, 'ms', renderEntry.name);
          }
        }
      });
      
      this.renderObserver.observe({ entryTypes: ['measure'] });
    } catch (error) {
      console.warn('Performance: Render monitoring setup failed', error);
    }
  }

  /**
   * Setup periodic memory usage checks
   */
  private setupPeriodicMemoryCheck(): void {
    if (typeof window === 'undefined' || !(performance as any).memory) return;

    this.memoryCheckInterval = setInterval(() => {
      const memory = (performance as any).memory;
      if (memory) {
        const usedMemory = memory.usedJSHeapSize;
        const totalMemory = memory.totalJSHeapSize;
        const memoryUsagePercentage = (usedMemory / totalMemory) * 100;
        
        analytics.trackPerformance('memory_usage', usedMemory, 'bytes', 'used_heap');
        analytics.trackPerformance('memory_usage', totalMemory, 'bytes', 'total_heap');
        analytics.trackPerformance('memory_usage', memoryUsagePercentage, 'percentage', 'memory_usage_percent');
      }
    }, 30000); // Check every 30 seconds
  }

  /**
   * Track initial bundle size and loading performance
   */
  private trackBundleSize(): void {
    if (typeof window === 'undefined') return;

    // Wait for page to fully load
    window.addEventListener('load', () => {
      setTimeout(() => {
        // Get navigation timing
        const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
        if (navigation) {
          const transferSize = navigation.transferSize || 0;
          const encodedBodySize = navigation.encodedBodySize || 0;
          const decodedBodySize = navigation.decodedBodySize || 0;
          
          analytics.trackPerformance('bundle_size', transferSize, 'bytes', 'transfer_size');
          analytics.trackPerformance('bundle_size', encodedBodySize, 'bytes', 'encoded_size');
          analytics.trackPerformance('bundle_size', decodedBodySize, 'bytes', 'decoded_size');
        }

        // Get resource timing for all loaded resources
        const resources = performance.getEntriesByType('resource');
        let totalResourceSize = 0;
        let resourceCount = 0;

        resources.forEach(resource => {
          const res = resource as PerformanceResourceTiming;
          if (res.transferSize) {
            totalResourceSize += res.transferSize;
            resourceCount++;
          }
        });

        if (resourceCount > 0) {
          analytics.trackPerformance('bundle_size', totalResourceSize, 'bytes', 'total_resources');
          analytics.trackPerformance('bundle_size', resourceCount, 'count', 'resource_count');
        }
      }, 1000);
    });
  }

  /**
   * Get current performance metrics
   */
  getCurrentMetrics(): PerformanceMetrics {
    const metrics: PerformanceMetrics = {};

    // Get Web Vitals from performance entries
    const lcpEntries = performance.getEntriesByType('largest-contentful-paint');
    if (lcpEntries.length > 0) {
      metrics.lcp = lcpEntries[lcpEntries.length - 1].startTime;
    }

    const fidEntries = performance.getEntriesByType('first-input');
    if (fidEntries.length > 0) {
      const fidEntry = fidEntries[0] as PerformanceEventTiming;
      metrics.fid = fidEntry.processingStart - fidEntry.startTime;
    }

    // Get CLS from performance entries
    let clsValue = 0;
    const clsEntries = performance.getEntriesByType('layout-shift');
    clsEntries.forEach(entry => {
      const layoutShiftEntry = entry as unknown as PerformanceEntry & { value: number; hadRecentInput: boolean };
      if (!layoutShiftEntry.hadRecentInput) {
        clsValue += layoutShiftEntry.value;
      }
    });
    metrics.cls = clsValue;

    // Get TTFB from navigation timing
    const navigation = performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming;
    if (navigation) {
      metrics.ttfb = navigation.responseStart - navigation.requestStart;
    }

    // Get memory usage if available
    if ((performance as any).memory) {
      const memory = (performance as any).memory;
      metrics.memoryUsage = memory.usedJSHeapSize;
    }

    return metrics;
  }

  /**
   * Create a performance mark for manual measurement
   */
  mark(name: string): void {
    if (typeof performance !== 'undefined' && performance.mark) {
      performance.mark(name);
    }
  }

  /**
   * Create a performance measure between two marks
   */
  measure(name: string, startMark: string, endMark?: string): number {
    if (typeof performance !== 'undefined' && performance.measure) {
      try {
        performance.measure(name, startMark, endMark);
        const measure = performance.getEntriesByName(name, 'measure')[0] as PerformanceMeasure;
        return measure.duration;
      } catch (error) {
        console.warn(`Performance: Failed to measure ${name}`, error);
        return 0;
      }
    }
    return 0;
  }

  /**
   * Get performance report for debugging
   */
  getPerformanceReport(): string {
    const metrics = this.getCurrentMetrics();
    const memory = (performance as any).memory;
    
    const report = {
      timestamp: new Date().toISOString(),
      webVitals: {
        lcp: metrics.lcp ? `${metrics.lcp.toFixed(2)}ms` : 'N/A',
        fid: metrics.fid ? `${metrics.fid.toFixed(2)}ms` : 'N/A',
        cls: metrics.cls ? metrics.cls.toFixed(4) : 'N/A',
        ttfb: metrics.ttfb ? `${metrics.ttfb.toFixed(2)}ms` : 'N/A'
      },
      memory: memory ? {
        used: `${(memory.usedJSHeapSize / 1024 / 1024).toFixed(2)}MB`,
        total: `${(memory.totalJSHeapSize / 1024 / 1024).toFixed(2)}MB`,
        limit: `${(memory.jsHeapSizeLimit / 1024 / 1024).toFixed(2)}MB`
      } : 'N/A',
      navigation: performance.getEntriesByType('navigation')[0] ? {
        domContentLoaded: `${((performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming).domContentLoadedEventEnd - (performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming).domContentLoadedEventStart).toFixed(2)}ms`,
        loadComplete: `${((performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming).loadEventEnd - (performance.getEntriesByType('navigation')[0] as PerformanceNavigationTiming).loadEventStart).toFixed(2)}ms`
      } : 'N/A'
    };
    
    return JSON.stringify(report, null, 2);
  }

  /**
   * Cleanup method
   */
  dispose(): void {
    if (this.memoryObserver) {
      this.memoryObserver.disconnect();
      this.memoryObserver = null;
    }
    
    if (this.renderObserver) {
      this.renderObserver.disconnect();
      this.renderObserver = null;
    }
    
    if (this.memoryCheckInterval) {
      clearInterval(this.memoryCheckInterval);
      this.memoryCheckInterval = null;
    }
    
    this.entries.clear();
    this.isInitialized = false;
  }
}

// Create singleton instance
export const performanceMonitor = new PerformanceMonitor();

// Export types and utilities
export type { PerformanceEntry };
export type { PerformanceMetrics };