'use client';

import React, { createContext, useContext, useEffect, useRef, useCallback, ReactNode } from 'react';
import { analytics } from '@/lib/analytics';
import { performanceMonitor } from '@/lib/performance/monitor';
import { AnalyticsEvent, AnalyticsSession } from '@/lib/types/analytics';

interface AnalyticsContextType {
  // Analytics methods
  trackEvent: (event: AnalyticsEvent) => void;
  trackInteraction: (action: 'click' | 'hover' | 'drag' | 'scroll' | 'keypress' | 'focus' | 'change' | 'validation', element: string, properties?: any) => void;
  trackPerformance: (metric: string, value: number, unit: string, context?: string) => void;
  trackAudio: (action: string, properties?: any) => void;
  trackFeature: (feature: string, action: string, properties?: any) => void;
  trackError: (error: Error, context?: string) => void;
  trackPageView: (page: string, referrer?: string) => void;
  
  // Performance methods
  startPerformanceTimer: (name: string, type?: string) => void;
  endPerformanceTimer: (name: string, context?: string) => number;
  measureAudioProcessing: <T>(name: string, fn: () => T | Promise<T>) => T | Promise<T>;
  measureRender: <T>(name: string, fn: () => T | Promise<T>) => T | Promise<T>;
  
  // Session and configuration
  getSession: () => AnalyticsSession | null;
  optOut: () => void;
  optIn: () => void;
  hasOptedOut: () => boolean;
  exportData: () => string;
  
  // Development only
  getPerformanceReport: () => string;
  getCurrentMetrics: () => any;
}

const AnalyticsContext = createContext<AnalyticsContextType | null>(null);

interface AnalyticsProviderProps {
  children: ReactNode;
  config?: {
    enabled?: boolean;
    debug?: boolean;
    endpoint?: string;
  };
}

export function AnalyticsProvider({ children, config }: AnalyticsProviderProps) {
  const isInitialized = useRef(false);
  const performanceTimers = useRef<Map<string, number>>(new Map());

  // Initialize analytics and performance monitoring
  useEffect(() => {
    if (isInitialized.current) return;

    // Initialize analytics with provided config
    if (config) {
      // Note: In a real implementation, you'd pass the config to analytics.initialize()
      // For now, we'll use the default configuration
    }
    
    analytics.initialize();
    performanceMonitor.initialize();
    
    // Setup global error tracking
    const handleError = (event: ErrorEvent) => {
      analytics.trackError({
        message: event.message,
        stack: event.error?.stack,
        source: 'system',
        severity: 'high',
        context: 'global_error_handler'
      });
    };
    
    // Setup unhandled promise rejection tracking
    const handleUnhandledRejection = (event: PromiseRejectionEvent) => {
      analytics.trackError({
        message: `Unhandled promise rejection: ${event.reason}`,
        stack: event.reason?.stack,
        source: 'system',
        severity: 'high',
        context: 'unhandled_promise_rejection'
      });
    };
    
    window.addEventListener('error', handleError);
    window.addEventListener('unhandledrejection', handleUnhandledRejection);
    
    // Track page visibility changes
    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        analytics.trackPageView(window.location.pathname);
      }
    };
    
    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    isInitialized.current = true;
    
    return () => {
      window.removeEventListener('error', handleError);
      window.removeEventListener('unhandledrejection', handleUnhandledRejection);
      document.removeEventListener('visibilitychange', handleVisibilityChange);
      analytics.dispose();
      performanceMonitor.dispose();
    };
  }, [config]);

  // Analytics methods
  const trackEvent = useCallback((event: AnalyticsEvent) => {
    // This would need to be implemented in the analytics service
    // For now, we'll use the specific tracking methods
    switch (event.type) {
      case 'interaction':
        analytics.trackInteraction(event.action, event.element, event);
        break;
      case 'performance':
        analytics.trackPerformance(event.metric, event.value, event.unit, event.context);
        break;
      case 'audio':
        analytics.trackAudio(event.action, event);
        break;
      case 'feature':
        analytics.trackFeature(event.feature, event.action, event.properties);
        break;
      case 'error':
        analytics.trackError(event.error, event.error.context);
        break;
      case 'pageview':
        analytics.trackPageView(event.page, event.referrer);
        break;
    }
  }, []);

  const trackInteraction = useCallback((action: 'click' | 'hover' | 'drag' | 'scroll' | 'keypress' | 'focus' | 'change' | 'validation', element: string, properties?: any) => {
    analytics.trackInteraction(action, element, properties);
  }, []);

  const trackPerformance = useCallback((metric: string, value: number, unit: string, context?: string) => {
    analytics.trackPerformance(metric as any, value, unit as any, context);
  }, []);

  const trackAudio = useCallback((action: string, properties?: any) => {
    analytics.trackAudio(action as any, properties);
  }, []);

  const trackFeature = useCallback((feature: string, action: string, properties?: any) => {
    analytics.trackFeature(feature as any, action as any, properties);
  }, []);

  const trackError = useCallback((error: Error, context?: string) => {
    analytics.trackError({
      message: error.message,
      stack: error.stack,
      source: 'user',
      severity: 'medium',
      context
    });
  }, []);

  const trackPageView = useCallback((page: string, referrer?: string) => {
    analytics.trackPageView(page, referrer);
  }, []);

  // Performance methods
  const startPerformanceTimer = useCallback((name: string, type = 'render') => {
    performanceTimers.current.set(name, performance.now());
    performanceMonitor.start(name, type as any);
  }, []);

  const endPerformanceTimer = useCallback((name: string, context?: string) => {
    const startTime = performanceTimers.current.get(name);
    if (startTime) {
      const duration = performance.now() - startTime;
      performanceTimers.current.delete(name);
      performanceMonitor.end(name, context);
      return duration;
    }
    return 0;
  }, []);

  const measureAudioProcessing = useCallback(<T,>(name: string, fn: () => T | Promise<T>) => {
    return performanceMonitor.measureAudioProcessing(name, fn);
  }, []);

  const measureRender = useCallback(<T,>(name: string, fn: () => T | Promise<T>) => {
    return performanceMonitor.measureRender(name, fn);
  }, []);

  // Session and configuration methods
  const getSession = useCallback(() => {
    return analytics.getSession();
  }, []);

  const optOut = useCallback(() => {
    analytics.optOut();
  }, []);

  const optIn = useCallback(() => {
    analytics.optIn();
  }, []);

  const hasOptedOut = useCallback(() => {
    return analytics.hasOptedOut();
  }, []);

  const exportData = useCallback(() => {
    return analytics.exportData();
  }, []);

  // Development only methods
  const getPerformanceReport = useCallback(() => {
    return performanceMonitor.getPerformanceReport();
  }, []);

  const getCurrentMetrics = useCallback(() => {
    return performanceMonitor.getCurrentMetrics();
  }, []);

  // Create the context value
  const contextValue: AnalyticsContextType = {
    trackEvent,
    trackInteraction,
    trackPerformance,
    trackAudio,
    trackFeature,
    trackError,
    trackPageView,
    startPerformanceTimer,
    endPerformanceTimer,
    measureAudioProcessing,
    measureRender,
    getSession,
    optOut,
    optIn,
    hasOptedOut,
    exportData,
    getPerformanceReport,
    getCurrentMetrics
  };

  return (
    <AnalyticsContext.Provider value={contextValue}>
      {children}
    </AnalyticsContext.Provider>
  );
}

// Hook to use analytics
export function useAnalytics() {
  const context = useContext(AnalyticsContext);
  if (!context) {
    throw new Error('useAnalytics must be used within an AnalyticsProvider');
  }
  return context;
}

// Error boundary component that automatically tracks errors
interface AnalyticsErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

interface AnalyticsErrorBoundaryProps {
  children: ReactNode;
  fallback?: ReactNode;
  onError?: (error: Error, errorInfo: React.ErrorInfo) => void;
}

class AnalyticsErrorBoundary extends React.Component<AnalyticsErrorBoundaryProps, AnalyticsErrorBoundaryState> {
  constructor(props: AnalyticsErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): AnalyticsErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // We need to get the analytics context, but we can't use hooks here
    // So we'll track the error directly
    if (typeof window !== 'undefined') {
      console.error('Analytics Error Boundary caught an error:', error, errorInfo);
    }
    
    this.props.onError?.(error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback || <div>Something went wrong.</div>;
    }

    return this.props.children;
  }
}

// Wrapper component that provides analytics to the error boundary
export function AnalyticsErrorBoundaryWithTracking({ children, fallback, onError }: AnalyticsErrorBoundaryProps) {
  const { trackError } = useAnalytics();

  const handleError = (error: Error, errorInfo: React.ErrorInfo) => {
    trackError(error, `React Error Boundary: ${errorInfo.componentStack}`);
    onError?.(error, errorInfo);
  };

  return (
    <AnalyticsErrorBoundary fallback={fallback} onError={handleError}>
      {children}
    </AnalyticsErrorBoundary>
  );
}

// Higher-order component for automatic tracking
export function withAnalytics<P extends object>(
  Component: React.ComponentType<P>,
  options: {
    trackRender?: boolean;
    trackInteractions?: boolean;
    componentName?: string;
  } = {}
) {
  const WrappedComponent = (props: P) => {
    const { trackRender, trackInteractions, componentName } = options;
    const { measureRender, trackInteraction } = useAnalytics();
    const displayName = componentName || Component.displayName || Component.name || 'Component';

    // Track render performance if enabled
    if (trackRender) {
      return measureRender(`render_${displayName}`, () => <Component {...props} />);
    }

    // Track interactions if enabled
    if (trackInteractions) {
      const handleClick = (event: React.MouseEvent) => {
        trackInteraction('click', `${displayName}_${(event.target as HTMLElement).tagName}`, {
          component: displayName
        });
      };

      return (
        <div onClick={handleClick}>
          <Component {...props} />
        </div>
      );
    }

    return <Component {...props} />;
  };

  WrappedComponent.displayName = `withAnalytics(${Component.displayName || Component.name || 'Component'})`;
  
  return WrappedComponent;
}