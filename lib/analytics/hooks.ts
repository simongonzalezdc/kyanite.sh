import { useCallback, useRef, useEffect, useState } from 'react';
import { useAnalytics } from '@/app/components/AnalyticsProvider';
import { AnalyticsEvent, InteractionEvent, PerformanceEvent, AudioEvent, FeatureEvent, ErrorEvent } from '@/lib/types/analytics';

/**
 * Hook for tracking custom events
 */
export function useTrackEvent() {
  const { trackEvent } = useAnalytics();

  return useCallback((event: AnalyticsEvent) => {
    trackEvent(event);
  }, [trackEvent]);
}

/**
 * Hook for tracking user interactions
 */
export function useTrackInteraction() {
  const { trackInteraction } = useAnalytics();

  return useCallback((
    action: InteractionEvent['action'],
    element: string,
    properties?: Partial<InteractionEvent>
  ) => {
    trackInteraction(action, element, properties);
  }, [trackInteraction]);
}

/**
 * Hook for tracking performance metrics
 */
export function useTrackPerformance() {
  const { trackPerformance } = useAnalytics();

  return useCallback((
    metric: PerformanceEvent['metric'],
    value: number,
    unit: PerformanceEvent['unit'],
    context?: string
  ) => {
    trackPerformance(metric, value, unit, context);
  }, [trackPerformance]);
}

/**
 * Hook for tracking audio processing events
 */
export function useTrackAudio() {
  const { trackAudio } = useAnalytics();

  return useCallback((
    action: AudioEvent['action'],
    properties?: Partial<AudioEvent>
  ) => {
    trackAudio(action, properties);
  }, [trackAudio]);
}

/**
 * Hook for tracking feature usage
 */
export function useTrackFeature() {
  const { trackFeature } = useAnalytics();

  return useCallback((
    feature: FeatureEvent['feature'],
    action: FeatureEvent['action'],
    properties?: Record<string, any>
  ) => {
    trackFeature(feature, action, properties);
  }, [trackFeature]);
}

/**
 * Hook for tracking errors
 */
export function useTrackError() {
  const { trackError } = useAnalytics();

  return useCallback((
    error: ErrorEvent['error'],
    context?: string
  ) => {
    trackError(new Error(error.message), context);
  }, [trackError]);
}

/**
 * Hook for tracking page views
 */
export function useTrackPageView() {
  const { trackPageView } = useAnalytics();

  return useCallback((
    page: string,
    referrer?: string
  ) => {
    trackPageView(page, referrer);
  }, [trackPageView]);
}

/**
 * Hook for measuring performance with automatic cleanup
 */
export function usePerformanceMeasure(name: string, enabled = true) {
  const { startPerformanceTimer, endPerformanceTimer } = useAnalytics();
  const isActive = useRef(false);

  const start = useCallback(() => {
    if (!enabled) return;
    isActive.current = true;
    startPerformanceTimer(name);
  }, [name, enabled, startPerformanceTimer]);

  const end = useCallback((context?: string) => {
    if (!enabled || !isActive.current) return 0;
    isActive.current = false;
    return endPerformanceTimer(name, context);
  }, [name, enabled, endPerformanceTimer]);

  const measure = useCallback(<T,>(fn: () => T | Promise<T>): T | Promise<T> => {
    if (!enabled) return fn();
    
    start();
    const result = fn();
    
    if (result instanceof Promise) {
      return result.finally(() => end());
    } else {
      end();
      return result;
    }
  }, [enabled, start, end]);

  return { start, end, measure };
}

/**
 * Hook for measuring audio processing performance
 */
export function useAudioPerformanceMeasure(name: string, enabled = true) {
  const { measureAudioProcessing } = useAnalytics();

  return useCallback(<T,>(fn: () => T | Promise<T>): T | Promise<T> => {
    if (!enabled) return fn();
    return measureAudioProcessing(name, fn);
  }, [name, enabled, measureAudioProcessing]);
}

/**
 * Hook for measuring component render performance
 */
export function useRenderPerformanceMeasure(name: string, enabled = true) {
  const { measureRender } = useAnalytics();

  return useCallback(<T,>(fn: () => T | Promise<T>): T | Promise<T> => {
    if (!enabled) return fn();
    return measureRender(name, fn);
  }, [name, enabled, measureRender]);
}

/**
 * Hook for tracking component lifecycle events
 */
export function useComponentTracking(componentName: string, options: {
  trackMount?: boolean;
  trackUnmount?: boolean;
  trackRender?: boolean;
  trackInteractions?: boolean;
} = {}) {
  const { trackFeature, trackInteraction } = useAnalytics();
  const { trackMount = true, trackUnmount = true, trackRender = false, trackInteractions = false } = options;
  const renderCount = useRef(0);
  const { measure } = usePerformanceMeasure(`render_${componentName}`, trackRender);

  // Track mount
  useEffect(() => {
    if (trackMount) {
      trackFeature('component', 'mount', { component: componentName });
    }

    return () => {
      if (trackUnmount) {
        trackFeature('component', 'unmount', { 
          component: componentName,
          renderCount: renderCount.current 
        });
      }
    };
  }, [componentName, trackMount, trackUnmount, trackFeature]);

  // Track renders
  useEffect(() => {
    if (trackRender) {
      renderCount.current++;
      measure(() => {
        trackFeature('component', 'render', {
          component: componentName,
          renderCount: renderCount.current
        });
      });
    }
  }); // Intentionally missing dependency array to track every render

  // Track interactions
  const trackClick = useCallback((element: string, properties?: any) => {
    if (trackInteractions) {
      trackInteraction('click', `${componentName}_${element}`, {
        component: componentName,
        ...properties
      });
    }
  }, [componentName, trackInteractions, trackInteraction]);

  const trackHover = useCallback((element: string, properties?: any) => {
    if (trackInteractions) {
      trackInteraction('hover', `${componentName}_${element}`, {
        component: componentName,
        ...properties
      });
    }
  }, [componentName, trackInteractions, trackInteraction]);

  return { trackClick, trackHover };
}

/**
 * Hook for tracking form interactions
 */
export function useFormTracking(formName: string) {
  const { trackInteraction, trackFeature } = useAnalytics();
  const fieldInteractions = useRef<Set<string>>(new Set());

  const trackFieldFocus = useCallback((fieldName: string) => {
    fieldInteractions.current.add(fieldName);
    trackInteraction('focus', `${formName}_${fieldName}`, {
      form: formName,
      field: fieldName,
      totalFieldsInteracted: fieldInteractions.current.size
    });
  }, [formName, trackInteraction]);

  const trackFieldChange = useCallback((fieldName: string, value: any) => {
    trackInteraction('change', `${formName}_${fieldName}`, {
      form: formName,
      field: fieldName,
      valueType: typeof value
    });
  }, [formName, trackInteraction]);

  const trackSubmit = useCallback((data: Record<string, any>, success: boolean) => {
    trackFeature('form', success ? 'submit_success' : 'submit_error', {
      form: formName,
      fieldCount: Object.keys(data).length,
      fieldsInteracted: fieldInteractions.current.size
    });
  }, [formName, trackFeature]);

  const trackValidation = useCallback((fieldName: string, isValid: boolean, message?: string) => {
    trackInteraction('validation', `${formName}_${fieldName}`, {
      form: formName,
      field: fieldName,
      isValid,
      hasMessage: !!message
    });
  }, [formName, trackInteraction]);

  return {
    trackFieldFocus,
    trackFieldChange,
    trackSubmit,
    trackValidation
  };
}

/**
 * Hook for tracking user journey and funnel events
 */
export function useJourneyTracking(journeyName: string) {
  const { trackFeature } = useAnalytics();
  const [currentStep, setCurrentStep] = useState<string>('');
  const stepsCompleted = useRef<Set<string>>(new Set());

  const trackStepStart = useCallback((stepName: string, properties?: any) => {
    setCurrentStep(stepName);
    trackFeature('journey', 'step_start', {
      journey: journeyName,
      step: stepName,
      ...properties
    });
  }, [journeyName, trackFeature]);

  const trackStepComplete = useCallback((stepName: string, properties?: any) => {
    stepsCompleted.current.add(stepName);
    trackFeature('journey', 'step_complete', {
      journey: journeyName,
      step: stepName,
      stepsCompleted: stepsCompleted.current.size,
      ...properties
    });
  }, [journeyName, trackFeature]);

  const trackJourneyComplete = useCallback((properties?: any) => {
    trackFeature('journey', 'complete', {
      journey: journeyName,
      totalSteps: stepsCompleted.current.size,
      ...properties
    });
  }, [journeyName, trackFeature]);

  const trackJourneyAbandon = useCallback((properties?: any) => {
    trackFeature('journey', 'abandon', {
      journey: journeyName,
      currentStep,
      stepsCompleted: stepsCompleted.current.size,
      ...properties
    });
  }, [journeyName, currentStep, trackFeature]);

  return {
    currentStep,
    stepsCompleted: Array.from(stepsCompleted.current),
    trackStepStart,
    trackStepComplete,
    trackJourneyComplete,
    trackJourneyAbandon
  };
}

/**
 * Hook for tracking A/B tests and experiments
 */
export function useExperimentTracking(experimentName: string) {
  const { trackFeature } = useAnalytics();
  const [variant, setVariant] = useState<string>('');
  const hasExposed = useRef(false);

  const expose = useCallback((variantName: string, properties?: any) => {
    if (hasExposed.current) return;
    
    setVariant(variantName);
    hasExposed.current = true;
    
    trackFeature('experiment', 'expose', {
      experiment: experimentName,
      variant: variantName,
      ...properties
    });
  }, [experimentName, trackFeature]);

  const trackConversion = useCallback((goal: string, value?: number, properties?: any) => {
    if (!hasExposed.current) return;
    
    trackFeature('experiment', 'convert', {
      experiment: experimentName,
      variant,
      goal,
      value,
      ...properties
    });
  }, [experimentName, variant, trackFeature]);

  return {
    variant,
    isExposed: hasExposed.current,
    expose,
    trackConversion
  };
}

/**
 * Hook for tracking media events (audio/video)
 */
export function useMediaTracking(mediaType: 'audio' | 'video', mediaName: string) {
  const { trackInteraction, trackFeature } = useAnalytics();
  const startTime = useRef<number>(0);
  const totalPlayTime = useRef<number>(0);
  const isPlaying = useRef(false);

  const trackPlay = useCallback(() => {
    startTime.current = Date.now();
    isPlaying.current = true;
    trackInteraction('click', `${mediaType}_play`, {
      mediaType,
      mediaName
    });
  }, [mediaType, mediaName, trackInteraction]);

  const trackPause = useCallback(() => {
    if (isPlaying.current) {
      totalPlayTime.current += Date.now() - startTime.current;
      isPlaying.current = false;
    }
    trackInteraction('click', `${mediaType}_pause`, {
      mediaType,
      mediaName,
      totalPlayTime: totalPlayTime.current
    });
  }, [mediaType, mediaName, trackInteraction]);

  const trackEnd = useCallback(() => {
    if (isPlaying.current) {
      totalPlayTime.current += Date.now() - startTime.current;
      isPlaying.current = false;
    }
    trackFeature('media', 'complete', {
      mediaType,
      mediaName,
      totalPlayTime: totalPlayTime.current
    });
  }, [mediaType, mediaName, trackFeature]);

  const trackSeek = useCallback((fromTime: number, toTime: number) => {
    trackInteraction('drag', `${mediaType}_seek`, {
      mediaType,
      mediaName,
      fromTime,
      toTime,
      seekDistance: Math.abs(toTime - fromTime)
    });
  }, [mediaType, mediaName, trackInteraction]);

  return {
    trackPlay,
    trackPause,
    trackEnd,
    trackSeek,
    getTotalPlayTime: () => totalPlayTime.current
  };
}