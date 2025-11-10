'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { useUIStore } from '@/lib/store/uiStore';
import { useUserStore } from '@/lib/store/userStore';

// Store context for hydration state
interface StoreContextType {
  isHydrated: boolean;
  isMobileDevice: boolean;
  error: Error | null;
}

const StoreContext = createContext<StoreContextType>({
  isHydrated: false,
  isMobileDevice: false,
  error: null,
});

// Props for StoreProvider
interface StoreProviderProps {
  children: ReactNode;
  enableDevTools?: boolean;
}

// Error boundary component for store failures
class StoreErrorBoundary extends React.Component<
  { children: ReactNode; onError: (error: Error) => void },
  { hasError: boolean; error: Error | null }
> {
  constructor(props: { children: ReactNode; onError: (error: Error) => void }) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('Store Error Boundary caught an error:', error, errorInfo);
    this.props.onError(error);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
          <div className="max-w-md w-full bg-white dark:bg-gray-800 rounded-lg shadow-lg p-6">
            <div className="flex items-center mb-4">
              <div className="flex-shrink-0">
                <svg className="h-6 w-6 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
                </svg>
              </div>
              <div className="ml-3">
                <h3 className="text-sm font-medium text-gray-900 dark:text-white">
                  Application Error
                </h3>
              </div>
            </div>
            <div className="text-sm text-gray-600 dark:text-gray-300">
              <p>Something went wrong with the application state. Please refresh the page to try again.</p>
              <details className="mt-2">
                <summary className="cursor-pointer font-medium">Error details</summary>
                <pre className="mt-2 text-xs bg-gray-100 dark:bg-gray-700 p-2 rounded overflow-auto">
                  {this.state.error?.message}
                </pre>
              </details>
            </div>
            <div className="mt-4">
              <button
                onClick={() => window.location.reload()}
                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                Refresh Page
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// Store Provider component
export function StoreProvider({ children, enableDevTools = process.env.NODE_ENV === 'development' }: StoreProviderProps) {
  const [isHydrated, setIsHydrated] = useState(false);
  const [isMobileDevice, setIsMobileDevice] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  
  // Get store instances for hydration
  const uiStore = useUIStore();
  const userStore = useUserStore();

  // Detect mobile device
  useEffect(() => {
    const detectMobile = () => {
      const userAgent = typeof navigator !== 'undefined' ? navigator.userAgent : '';
      const mobile = /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(userAgent);
      setIsMobileDevice(mobile);
      
      // Update UI store with mobile detection
      uiStore.setNavigationVisible(!mobile);
    };

    detectMobile();
    
    // Listen for window resize to update mobile detection
    if (typeof window !== 'undefined') {
      window.addEventListener('resize', detectMobile);
      return () => window.removeEventListener('resize', detectMobile);
    }
  }, [uiStore]);

  // Handle store hydration
  useEffect(() => {
    const handleHydration = async () => {
      try {
        // Wait for stores to hydrate from localStorage
        await new Promise(resolve => setTimeout(resolve, 100));
        
        // Initialize session data if needed
        const { sessionId, startTime } = userStore;
        if (!sessionId) {
          userStore.setSessionId(`session-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`);
          userStore.setStartTime(new Date());
          userStore.setLastActivity(new Date());
        }
        
        // Initialize device info
        if (typeof navigator !== 'undefined') {
          userStore.setDeviceInfo({
            userAgent: navigator.userAgent,
            platform: navigator.platform,
            screenResolution: typeof window !== 'undefined' ? `${window.screen.width}x${window.screen.height}` : '',
            isMobile: isMobileDevice,
            isTouchDevice: 'ontouchstart' in window || navigator.maxTouchPoints > 0,
          });
        }
        
        setIsHydrated(true);
      } catch (err) {
        console.error('Store hydration failed:', err);
        setError(err instanceof Error ? err : new Error('Store hydration failed'));
      }
    };

    handleHydration();
  }, [userStore, isMobileDevice]);

  // Setup DevTools integration
  useEffect(() => {
    if (enableDevTools && typeof window !== 'undefined') {
      // Expose stores to window for debugging
      (window as any).__ZUSTAND_STORES__ = {
        ui: uiStore,
        user: userStore,
      };
      
      // Add custom DevTools commands
      if ((window as any).__REDUX_DEVTOOLS_EXTENSION__) {
        console.log('Zustand stores are available for debugging in window.__ZUSTAND_STORES__');
      }
    }
  }, [enableDevTools, uiStore, userStore]);

  // Handle online/offline status
  useEffect(() => {
    const handleOnline = () => {
      uiStore.setGlobalError(null);
    };

    const handleOffline = () => {
      uiStore.setGlobalError('You are currently offline. Some features may not be available.');
    };

    if (typeof window !== 'undefined') {
      window.addEventListener('online', handleOnline);
      window.addEventListener('offline', handleOffline);
      
      // Set initial online status
      if (!navigator.onLine) {
        handleOffline();
      }
      
      return () => {
        window.removeEventListener('online', handleOnline);
        window.removeEventListener('offline', handleOffline);
      };
    }
  }, [uiStore]);

  // Handle beforeunload for unsaved changes
  useEffect(() => {
    const handleBeforeUnload = (e: BeforeUnloadEvent) => {
      // Check if there are unsaved changes
      const { isDirty } = useProjectStore();
      if (isDirty) {
        const message = 'You have unsaved changes. Are you sure you want to leave?';
        e.returnValue = message;
        return message;
      }
    };

    if (typeof window !== 'undefined') {
      window.addEventListener('beforeunload', handleBeforeUnload);
      return () => window.removeEventListener('beforeunload', handleBeforeUnload);
    }
  }, []);

  const contextValue = {
    isHydrated,
    isMobileDevice,
    error,
  };

  return (
    <StoreContext.Provider value={contextValue}>
      <StoreErrorBoundary onError={setError}>
        {children}
      </StoreErrorBoundary>
    </StoreContext.Provider>
  );
}

// Hook to access store context
export function useStoreContext() {
  const context = useContext(StoreContext);
  if (!context) {
    throw new Error('useStoreContext must be used within a StoreProvider');
  }
  return context;
}

// Hook to check if stores are hydrated
export function useIsHydrated() {
  const { isHydrated } = useStoreContext();
  return isHydrated;
}

// Hook to check if on mobile device
export function useIsMobileDevice() {
  const { isMobileDevice } = useStoreContext();
  return isMobileDevice;
}

// Hook to get store error
export function useStoreError() {
  const { error } = useStoreContext();
  return error;
}

// Higher-order component for hydration-aware components
export function withHydration<P extends object>(Component: React.ComponentType<P>) {
  return function HydratedComponent(props: P) {
    const isHydrated = useIsHydrated();
    
    if (!isHydrated) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
            <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">Loading application...</p>
          </div>
        </div>
      );
    }
    
    return <Component {...props} />;
  };
}

// Component to handle loading state during hydration
export function HydrationLoader({ children }: { children: ReactNode }) {
  const isHydrated = useIsHydrated();
  
  if (!isHydrated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-sm text-gray-600 dark:text-gray-300">Loading application...</p>
        </div>
      </div>
    );
  }
  
  return <>{children}</>;
}

// Import ProjectStore for beforeunload handler
import { useProjectStore } from '@/lib/store/projectStore';