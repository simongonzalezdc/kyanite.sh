import { 
  AnalyticsEvent, 
  AnalyticsConfig, 
  AnalyticsSession, 
  PerformanceMetrics,
  InteractionEvent,
  PerformanceEvent,
  AudioEvent,
  FeatureEvent,
  ErrorEvent,
  PageViewEvent
} from '@/lib/types/analytics';

class AnalyticsService {
  private config: AnalyticsConfig;
  private session: AnalyticsSession | null = null;
  private eventQueue: AnalyticsEvent[] = [];
  private flushTimer: NodeJS.Timeout | null = null;
  private isInitialized = false;
  private performanceObserver: PerformanceObserver | null = null;

  constructor(config?: Partial<AnalyticsConfig>) {
    // Default configuration
    this.config = {
      enabled: process.env.NODE_ENV === 'production',
      debug: process.env.NODE_ENV === 'development',
      endpoint: process.env.NEXT_PUBLIC_ANALYTICS_ENDPOINT,
      batchSize: 10,
      flushInterval: 30000, // 30 seconds
      maxRetries: 3,
      localStorageKey: 'voxforge-analytics',
      sessionTimeout: 1800000, // 30 minutes
      optOutKey: 'voxforge-analytics-optout',
      ...config
    };
  }

  /**
   * Initialize the analytics service
   */
  initialize(): void {
    if (this.isInitialized || !this.config.enabled) return;

    // Check if user has opted out
    if (this.hasOptedOut()) {
      console.info('Analytics: User has opted out of tracking');
      return;
    }

    // Initialize or restore session
    this.initializeSession();
    
    // Setup performance monitoring
    this.setupPerformanceMonitoring();
    
    // Setup periodic flush
    this.setupPeriodicFlush();
    
    // Track initial page view
    this.trackPageView(window.location.pathname);
    
    // Setup unload handler to flush events
    window.addEventListener('beforeunload', () => {
      this.flushEvents(true);
    });
    
    this.isInitialized = true;
    
    if (this.config.debug) {
      console.info('Analytics: Service initialized', { sessionId: this.session?.id });
    }
  }

  /**
   * Check if user has opted out of analytics
   */
  hasOptedOut(): boolean {
    if (typeof window === 'undefined') return false;
    return localStorage.getItem(this.config.optOutKey) === 'true';
  }

  /**
   * Allow users to opt out of analytics
   */
  optOut(): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(this.config.optOutKey, 'true');
    this.clearLocalData();
  }

  /**
   * Allow users to opt back into analytics
   */
  optIn(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(this.config.optOutKey);
    this.initialize();
  }

  /**
   * Initialize or restore analytics session
   */
  private initializeSession(): void {
    const storedSession = this.getStoredSession();
    
    if (storedSession && this.isSessionValid(storedSession)) {
      this.session = storedSession;
      this.session.lastActivity = Date.now();
    } else {
      this.session = this.createNewSession();
    }
    
    this.storeSession();
  }

  /**
   * Create a new analytics session
   */
  private createNewSession(): AnalyticsSession {
    const sessionId = this.generateAnonymousId();
    const now = Date.now();
    
    return {
      id: sessionId,
      startTime: now,
      lastActivity: now,
      pageViews: 0,
      events: 0,
      duration: 0,
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
      viewport: {
        width: typeof window !== 'undefined' ? window.innerWidth : 0,
        height: typeof window !== 'undefined' ? window.innerHeight : 0
      }
    };
  }

  /**
   * Check if a session is still valid
   */
  private isSessionValid(session: AnalyticsSession): boolean {
    return (Date.now() - session.lastActivity) < this.config.sessionTimeout;
  }

  /**
   * Generate an anonymous ID for session identification
   */
  private generateAnonymousId(): string {
    return 'anon_' + Math.random().toString(36).substr(2, 9) + '_' + Date.now().toString(36);
  }

  /**
   * Track a user interaction event
   */
  trackInteraction(action: InteractionEvent['action'], element: string, properties?: Partial<InteractionEvent>): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: InteractionEvent = {
      type: 'interaction',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      action,
      element,
      page: window.location.pathname,
      ...properties
    };
    
    this.trackEvent(event);
  }

  /**
   * Track a performance metric
   */
  trackPerformance(metric: PerformanceEvent['metric'], value: number, unit: PerformanceEvent['unit'], context?: string): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: PerformanceEvent = {
      type: 'performance',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      metric,
      value,
      unit,
      context
    };
    
    this.trackEvent(event);
  }

  /**
   * Track an audio processing event
   */
  trackAudio(action: AudioEvent['action'], properties?: Partial<AudioEvent>): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: AudioEvent = {
      type: 'audio',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      action,
      success: true,
      ...properties
    };
    
    this.trackEvent(event);
  }

  /**
   * Track feature usage
   */
  trackFeature(feature: FeatureEvent['feature'], action: FeatureEvent['action'], properties?: Record<string, any>): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: FeatureEvent = {
      type: 'feature',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      feature,
      action,
      properties
    };
    
    this.trackEvent(event);
  }

  /**
   * Track an error
   */
  trackError(error: ErrorEvent['error'], context?: string): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: ErrorEvent = {
      type: 'error',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      error: {
        ...error,
        context
      }
    };
    
    this.trackEvent(event);
  }

  /**
   * Track a page view
   */
  trackPageView(page: string, referrer?: string): void {
    if (!this.isInitialized || !this.session) return;
    
    const event: PageViewEvent = {
      type: 'pageview',
      timestamp: Date.now(),
      sessionId: this.session.id,
      version: this.getAppVersion(),
      page,
      referrer,
      userAgent: typeof navigator !== 'undefined' ? navigator.userAgent : '',
      viewport: {
        width: typeof window !== 'undefined' ? window.innerWidth : 0,
        height: typeof window !== 'undefined' ? window.innerHeight : 0
      }
    };
    
    this.trackEvent(event);
    
    // Update session page views
    this.session.pageViews++;
    this.session.lastActivity = Date.now();
    this.storeSession();
  }

  /**
   * Generic event tracking method
   */
  private trackEvent(event: AnalyticsEvent): void {
    if (!this.config.enabled) return;
    
    // Add to queue
    this.eventQueue.push(event);
    
    // Update session
    if (this.session) {
      this.session.events++;
      this.session.lastActivity = Date.now();
      this.storeSession();
    }
    
    // Flush if batch size reached
    if (this.eventQueue.length >= this.config.batchSize) {
      this.flushEvents();
    }
    
    if (this.config.debug) {
      console.debug('Analytics: Event tracked', event);
    }
  }

  /**
   * Setup performance monitoring using PerformanceObserver API
   */
  private setupPerformanceMonitoring(): void {
    if (typeof window === 'undefined' || !window.PerformanceObserver) return;
    
    try {
      // Monitor Core Web Vitals
      this.performanceObserver = new PerformanceObserver((list) => {
        for (const entry of list.getEntries()) {
          switch (entry.entryType) {
            case 'largest-contentful-paint':
              this.trackPerformance('LCP', entry.startTime, 'ms');
              break;
            case 'first-input':
              const firstInputEntry = entry as PerformanceEventTiming;
              this.trackPerformance('FID', firstInputEntry.processingStart - firstInputEntry.startTime, 'ms');
              break;
            case 'layout-shift':
              const layoutShiftEntry = entry as PerformanceEntry & { value: number; hadRecentInput: boolean };
              if (!layoutShiftEntry.hadRecentInput) {
                this.trackPerformance('CLS', layoutShiftEntry.value, 'score');
              }
              break;
            case 'navigation':
              const navEntry = entry as PerformanceNavigationTiming;
              this.trackPerformance('TTFB', navEntry.responseStart - navEntry.requestStart, 'ms');
              break;
          }
        }
      });
      
      this.performanceObserver.observe({ entryTypes: ['largest-contentful-paint', 'first-input', 'layout-shift', 'navigation'] });
    } catch (error) {
      console.warn('Analytics: Performance monitoring setup failed', error);
    }
  }

  /**
   * Setup periodic flush of events
   */
  private setupPeriodicFlush(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
    }
    
    this.flushTimer = setInterval(() => {
      this.flushEvents();
    }, this.config.flushInterval);
  }

  /**
   * Flush events to endpoint or local storage
   */
  private async flushEvents(isSync = false): Promise<void> {
    if (this.eventQueue.length === 0) return;
    
    const events = [...this.eventQueue];
    this.eventQueue = [];
    
    try {
      if (this.config.endpoint) {
        await this.sendToEndpoint(events, isSync);
      } else {
        // Fallback to local storage
        this.storeEvents(events);
      }
    } catch (error) {
      console.error('Analytics: Failed to flush events', error);
      // Re-add events to queue for retry
      this.eventQueue.unshift(...events);
    }
  }

  /**
   * Send events to analytics endpoint
   */
  private async sendToEndpoint(events: AnalyticsEvent[], isSync = false): Promise<void> {
    if (!this.config.endpoint) return;
    
    const payload = {
      events,
      session: this.session,
      timestamp: Date.now()
    };
    
    const options: RequestInit = {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
      keepalive: isSync
    };
    
    if (isSync) {
      // Use sendBeacon for synchronous requests during page unload
      if (navigator.sendBeacon) {
        navigator.sendBeacon(this.config.endpoint, JSON.stringify(payload));
        return;
      }
    }
    
    const response = await fetch(this.config.endpoint, options);
    
    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`);
    }
  }

  /**
   * Store events in local storage as fallback
   */
  private storeEvents(events: AnalyticsEvent[]): void {
    if (typeof window === 'undefined') return;
    
    try {
      const stored = localStorage.getItem(this.config.localStorageKey);
      const existingEvents = stored ? JSON.parse(stored) : [];
      const allEvents = [...existingEvents, ...events];
      
      // Keep only last 1000 events to avoid storage issues
      const trimmedEvents = allEvents.slice(-1000);
      
      localStorage.setItem(this.config.localStorageKey, JSON.stringify(trimmedEvents));
    } catch (error) {
      console.error('Analytics: Failed to store events locally', error);
    }
  }

  /**
   * Store session data
   */
  private storeSession(): void {
    if (typeof window === 'undefined' || !this.session) return;
    
    try {
      localStorage.setItem(`${this.config.localStorageKey}-session`, JSON.stringify(this.session));
    } catch (error) {
      console.error('Analytics: Failed to store session', error);
    }
  }

  /**
   * Get stored session from local storage
   */
  private getStoredSession(): AnalyticsSession | null {
    if (typeof window === 'undefined') return null;
    
    try {
      const stored = localStorage.getItem(`${this.config.localStorageKey}-session`);
      return stored ? JSON.parse(stored) : null;
    } catch (error) {
      console.error('Analytics: Failed to retrieve session', error);
      return null;
    }
  }

  /**
   * Clear all local analytics data
   */
  private clearLocalData(): void {
    if (typeof window === 'undefined') return;
    
    localStorage.removeItem(this.config.localStorageKey);
    localStorage.removeItem(`${this.config.localStorageKey}-session`);
  }

  /**
   * Get app version
   */
  private getAppVersion(): string {
    return process.env.NEXT_PUBLIC_APP_VERSION || '1.0.0';
  }

  /**
   * Get current session
   */
  getSession(): AnalyticsSession | null {
    return this.session;
  }

  /**
   * Get stored events for export
   */
  getStoredEvents(): AnalyticsEvent[] {
    if (typeof window === 'undefined') return [];
    
    try {
      const stored = localStorage.getItem(this.config.localStorageKey);
      return stored ? JSON.parse(stored) : [];
    } catch (error) {
      console.error('Analytics: Failed to retrieve stored events', error);
      return [];
    }
  }

  /**
   * Export analytics data
   */
  exportData(): string {
    const events = this.getStoredEvents();
    const session = this.session;
    
    const exportData = {
      version: this.getAppVersion(),
      exportDate: new Date().toISOString(),
      session,
      events,
      summary: {
        totalEvents: events.length,
        sessionDuration: session ? Date.now() - session.startTime : 0,
        pageViews: session?.pageViews || 0
      }
    };
    
    return JSON.stringify(exportData, null, 2);
  }

  /**
   * Cleanup method
   */
  dispose(): void {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
      this.flushTimer = null;
    }
    
    if (this.performanceObserver) {
      this.performanceObserver.disconnect();
      this.performanceObserver = null;
    }
    
    // Flush any remaining events
    this.flushEvents(true);
    
    this.isInitialized = false;
  }
}

// Create singleton instance
export const analytics = new AnalyticsService();

// Export types
export type { AnalyticsEvent, AnalyticsConfig, AnalyticsSession, PerformanceMetrics };