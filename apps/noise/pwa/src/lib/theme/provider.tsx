"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  useSyncExternalStore,
  ReactNode,
} from "react";
import { themes, defaultThemeId, getTheme, listThemes, type Theme } from "./tokens";

// =============================================================================
// Constants
// =============================================================================

const THEME_STORAGE_KEY = "noise-theme";

// =============================================================================
// Types
// =============================================================================

interface ThemeContextValue {
  theme: Theme;
  themeId: string;
  setTheme: (id: string) => void;
  availableThemes: Theme[];
}

// =============================================================================
// Context
// =============================================================================

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function useTheme(): ThemeContextValue {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within a ThemeProvider");
  }
  return context;
}

// =============================================================================
// CSS Variable Application
// =============================================================================

function applyThemeToDocument(theme: Theme): void {
  const root = document.documentElement;

  // Map theme properties to CSS variables
  root.style.setProperty("--color-primary", theme.primary);
  root.style.setProperty("--color-secondary", theme.secondary);
  root.style.setProperty("--color-accent", theme.accent);
  root.style.setProperty("--color-background", theme.background);
  root.style.setProperty("--color-text", theme.text);
  root.style.setProperty("--color-success", theme.success);
  root.style.setProperty("--color-warning", theme.warning);
  root.style.setProperty("--color-error", theme.error);

  // Derived colors
  root.style.setProperty("--color-text-muted", adjustOpacity(theme.text, 0.6));
  root.style.setProperty("--color-surface", adjustBrightness(theme.background, 1.5));

  // Update meta theme-color for mobile browsers
  const metaThemeColor = document.querySelector('meta[name="theme-color"]');
  if (metaThemeColor) {
    metaThemeColor.setAttribute("content", theme.primary);
  }
}

/**
 * Adjust color opacity by converting to rgba
 */
function adjustOpacity(hex: string, opacity: number): string {
  // Validate hex format
  if (!hex || !/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)) {
    return `rgba(0, 0, 0, ${opacity})`;
  }

  // Expand shorthand hex (#RGB to #RRGGBB)
  let fullHex = hex;
  if (hex.length === 4) {
    fullHex = "#" + hex[1] + hex[1] + hex[2] + hex[2] + hex[3] + hex[3];
  }

  const r = parseInt(fullHex.slice(1, 3), 16);
  const g = parseInt(fullHex.slice(3, 5), 16);
  const b = parseInt(fullHex.slice(5, 7), 16);
  return `rgba(${r}, ${g}, ${b}, ${opacity})`;
}

/**
 * Adjust color brightness
 */
function adjustBrightness(hex: string, factor: number): string {
  // Validate hex format
  if (!hex || !/^#([A-Fa-f0-9]{3}){1,2}$/.test(hex)) {
    return "#000000";
  }

  // Expand shorthand hex
  let fullHex = hex;
  if (hex.length === 4) {
    fullHex = "#" + hex[1] + hex[1] + hex[2] + hex[2] + hex[3] + hex[3];
  }

  const r = Math.min(255, Math.round(parseInt(fullHex.slice(1, 3), 16) * factor));
  const g = Math.min(255, Math.round(parseInt(fullHex.slice(3, 5), 16) * factor));
  const b = Math.min(255, Math.round(parseInt(fullHex.slice(5, 7), 16) * factor));
  return `#${r.toString(16).padStart(2, "0")}${g.toString(16).padStart(2, "0")}${b.toString(16).padStart(2, "0")}`;
}

// =============================================================================
// Provider
// =============================================================================

interface ThemeProviderProps {
  children: ReactNode;
}

// External store for mounted state to avoid setState in effects
let isMounted = false;
const mountedListeners = new Set<() => void>();

function subscribeToMounted(callback: () => void) {
  mountedListeners.add(callback);
  // Mark as mounted immediately on first subscription (client-side)
  if (!isMounted && typeof window !== "undefined") {
    isMounted = true;
    // Notify in next tick to avoid sync issues
    queueMicrotask(() => mountedListeners.forEach((cb) => cb()));
  }
  return () => mountedListeners.delete(callback);
}

function getMountedSnapshot() {
  return isMounted;
}

function getMountedServerSnapshot() {
  return false;
}

export function ThemeProvider({ children }: ThemeProviderProps) {
  // Initialize theme from localStorage (lazy initialization)
  const [themeId, setThemeId] = useState<string>(() => {
    if (typeof window === "undefined") return defaultThemeId;
    const savedTheme = localStorage.getItem(THEME_STORAGE_KEY);
    return savedTheme && themes[savedTheme] ? savedTheme : defaultThemeId;
  });

  // Track mounted state for hydration safety
  const mounted = useSyncExternalStore(
    subscribeToMounted,
    getMountedSnapshot,
    getMountedServerSnapshot
  );

  // Apply theme when it changes
  useEffect(() => {
    if (mounted) {
      const theme = getTheme(themeId);
      applyThemeToDocument(theme);
    }
  }, [themeId, mounted]);

  const setTheme = useCallback((id: string) => {
    if (themes[id]) {
      setThemeId(id);
      localStorage.setItem(THEME_STORAGE_KEY, id);
    }
  }, []);

  const theme = getTheme(themeId);
  const availableThemes = listThemes().map((id) => themes[id]);

  // Prevent flash of unstyled content
  if (!mounted) {
    return null;
  }

  return (
    <ThemeContext.Provider value={{ theme, themeId, setTheme, availableThemes }}>
      {children}
    </ThemeContext.Provider>
  );
}
