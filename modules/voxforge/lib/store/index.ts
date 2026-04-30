import { create, StateCreator } from 'zustand';
import { devtools, persist, subscribeWithSelector } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';

// Store configuration interface
export interface StoreConfig {
  name: string;
  version?: number;
  devtools?: boolean;
  persist?: {
    name: string;
    partialize?: (state: any) => any;
    version?: number;
    onRehydrateStorage?: () => (state: any) => void;
  };
}

// Default store configuration
const defaultConfig: StoreConfig = {
  name: 'voxforge-store',
  version: 1,
  devtools: process.env.NODE_ENV === 'development',
};

// Simple store creator with optional middleware
export const createStore = <T>(
  stateCreator: StateCreator<T, [], [], T>,
  config: Partial<StoreConfig> = {}
) => {
  const finalConfig = { ...defaultConfig, ...config };
  
  // Start with the base state creator
  let creator: any = stateCreator;
  
  // Apply immer middleware for immutable updates
  creator = immer(creator);
  
  // Apply subscribeWithSelector for subscription capabilities
  creator = subscribeWithSelector(creator);
  
  // Add devtools in development
  if (finalConfig.devtools) {
    creator = devtools(creator, {
      name: finalConfig.name,
    });
  }
  
  // Add persistence if configured
  if (finalConfig.persist) {
    creator = persist(creator, {
      name: finalConfig.persist.name,
      partialize: finalConfig.persist.partialize,
      version: finalConfig.persist.version,
      onRehydrateStorage: finalConfig.persist.onRehydrateStorage,
    });
  }
  
  return create<T>(creator as StateCreator<T, [], [], T>);
};

// Store composition helper for combining multiple stores
export const combineStores = <T extends Record<string, any>>(stores: T) => {
  return create<T>((set, get) => ({
    ...Object.entries(stores).reduce((acc, [key, store]) => {
      return {
        ...acc,
        [key]: store.getState(),
      };
    }, {} as T),
  }));
};

// Type helpers for store creation
export type StoreState<T> = T extends { getState: () => infer S } ? S : never;
export type StoreActions<T> = Pick<StoreState<T>, { [K in keyof StoreState<T>]: StoreState<T>[K] extends Function ? K : never }[keyof StoreState<T>]>;

// Rehydration helper
export const createRehydrateHandler = (storeName: string) => {
  return () => (state: any) => {
    console.log(`Store ${storeName} rehydrated:`, state);
  };
};

// Error boundary for store failures
export class StoreError extends Error {
  constructor(
    message: string,
    public storeName: string,
    public originalError?: Error
  ) {
    super(message);
    this.name = 'StoreError';
  }
}

// Store health check
export const checkStoreHealth = (store: any, storeName: string) => {
  try {
    const state = store.getState();
    if (!state) {
      throw new StoreError(`Store ${storeName} has no state`, storeName);
    }
    return true;
  } catch (error) {
    console.error(`Store health check failed for ${storeName}:`, error);
    return false;
  }
};

// Export types for external use
export type { StateCreator } from 'zustand';