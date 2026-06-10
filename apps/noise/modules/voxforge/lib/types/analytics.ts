// Base analytics event type
export interface BaseAnalyticsEvent {
  type: string;
  timestamp: number;
  sessionId: string;
  userId?: string; // Anonymous identifier only
  version: string;
}

// User interaction events
export interface InteractionEvent extends BaseAnalyticsEvent {
  type: 'interaction';
  action: 'click' | 'hover' | 'drag' | 'scroll' | 'keypress' | 'focus' | 'change' | 'validation';
  element: string;
  elementId?: string;
  elementClass?: string;
  page: string;
  value?: string | number;
}

// Performance events
export interface PerformanceEvent extends BaseAnalyticsEvent {
  type: 'performance';
  metric: 'FCP' | 'LCP' | 'FID' | 'CLS' | 'TTFB' | 'audio_processing' | 'render_time' | 'memory_usage' | 'bundle_size' | 'custom';
  value: number;
  unit: 'ms' | 'score' | 'bytes' | 'percentage' | 'count';
  context?: string;
}

// Audio processing events
export interface AudioEvent extends BaseAnalyticsEvent {
  type: 'audio';
  action: 'record_start' | 'record_complete' | 'analysis_start' | 'analysis_complete' | 'generation_start' | 'generation_complete';
  duration?: number; // in seconds
  sampleRate?: number;
  bufferSize?: number;
  success: boolean;
  error?: string;
}

// Feature usage events
export interface FeatureEvent extends BaseAnalyticsEvent {
  type: 'feature';
  feature: 'recording' | 'analysis' | 'generation' | 'export' | 'visualization' | 'onboarding';
  action: 'start' | 'complete' | 'skip' | 'error';
  properties?: Record<string, any>;
}

// Error events
export interface ErrorEvent extends BaseAnalyticsEvent {
  type: 'error';
  error: {
    message: string;
    stack?: string;
    source: 'user' | 'system' | 'network';
    severity: 'low' | 'medium' | 'high' | 'critical';
    context?: string;
  };
}

// Page view events
export interface PageViewEvent extends BaseAnalyticsEvent {
  type: 'pageview';
  page: string;
  referrer?: string;
  userAgent?: string;
  viewport?: {
    width: number;
    height: number;
  };
}

// Union type for all analytics events
export type AnalyticsEvent = 
  | InteractionEvent 
  | PerformanceEvent 
  | AudioEvent 
  | FeatureEvent 
  | ErrorEvent 
  | PageViewEvent;

// Analytics configuration
export interface AnalyticsConfig {
  enabled: boolean;
  debug: boolean;
  endpoint?: string;
  batchSize: number;
  flushInterval: number; // in ms
  maxRetries: number;
  localStorageKey: string;
  sessionTimeout: number; // in ms
  optOutKey: string;
}

// Performance metrics interface
export interface PerformanceMetrics {
  fcp?: number; // First Contentful Paint
  lcp?: number; // Largest Contentful Paint
  fid?: number; // First Input Delay
  cls?: number; // Cumulative Layout Shift
  ttfb?: number; // Time to First Byte
  audioProcessing?: number;
  renderTime?: number;
  memoryUsage?: number;
}

// Analytics session info
export interface AnalyticsSession {
  id: string;
  startTime: number;
  lastActivity: number;
  pageViews: number;
  events: number;
  duration: number;
  userAgent: string;
  viewport: {
    width: number;
    height: number;
  };
}

// Analytics data export format
export interface AnalyticsExport {
  version: string;
  exportDate: string;
  sessions: AnalyticsSession[];
  events: AnalyticsEvent[];
  metrics: PerformanceMetrics[];
  summary: {
    totalSessions: number;
    totalEvents: number;
    averageSessionDuration: number;
    topFeatures: Array<{
      feature: string;
      count: number;
    }>;
    errorRate: number;
  };
}