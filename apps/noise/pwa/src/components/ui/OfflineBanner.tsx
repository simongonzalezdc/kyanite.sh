"use client";

import { useRef, useEffect, useSyncExternalStore } from "react";

/** Duration to show the "back online" banner (ms) */
const BACK_ONLINE_BANNER_DURATION_MS = 2000;

interface OfflineBannerProps {
  /** Override the online status (for testing) */
  forceOffline?: boolean;
}

// External store for online status
function subscribeToOnlineStatus(callback: () => void) {
  window.addEventListener("online", callback);
  window.addEventListener("offline", callback);
  return () => {
    window.removeEventListener("online", callback);
    window.removeEventListener("offline", callback);
  };
}

function getOnlineSnapshot() {
  return navigator.onLine;
}

function getServerSnapshot() {
  return true; // Assume online during SSR
}

// Store for showing the "back online" banner temporarily
let backOnlineBannerVisible = false;
const backOnlineBannerListeners: Set<() => void> = new Set();

function subscribeToBackOnlineBanner(callback: () => void) {
  backOnlineBannerListeners.add(callback);
  return () => backOnlineBannerListeners.delete(callback);
}

function getBackOnlineBannerSnapshot() {
  return backOnlineBannerVisible;
}

function showBackOnlineBanner() {
  backOnlineBannerVisible = true;
  backOnlineBannerListeners.forEach((cb) => cb());
  setTimeout(() => {
    backOnlineBannerVisible = false;
    backOnlineBannerListeners.forEach((cb) => cb());
  }, BACK_ONLINE_BANNER_DURATION_MS);
}

/**
 * Banner that appears when the user is offline.
 * Uses the browser's online/offline events for detection.
 */
export function OfflineBanner({ forceOffline }: OfflineBannerProps) {
  // Use useSyncExternalStore for online status (React 18+ pattern)
  const browserOnline = useSyncExternalStore(
    subscribeToOnlineStatus,
    getOnlineSnapshot,
    getServerSnapshot
  );

  // Track "back online" banner visibility
  const showBackOnline = useSyncExternalStore(
    subscribeToBackOnlineBanner,
    getBackOnlineBannerSnapshot,
    () => false
  );

  // Combine browser status with forceOffline prop
  const isOnline = forceOffline !== undefined ? !forceOffline : browserOnline;

  // Track previous online status to detect transitions
  const wasOfflineRef = useRef(false);

  // Handle transition to online - show "back online" message
  useEffect(() => {
    if (isOnline && wasOfflineRef.current) {
      showBackOnlineBanner();
    }
    wasOfflineRef.current = !isOnline;
  }, [isOnline]);

  // Show banner if offline OR if showing "back online" message
  const showBanner = !isOnline || showBackOnline;

  if (!showBanner) {
    return null;
  }

  return (
    <div
      role="status"
      aria-live="polite"
      className={`
        fixed top-0 left-0 right-0 z-40
        px-4 py-3
        flex items-center justify-center gap-2
        text-sm font-medium
        transition-all duration-300
        ${
          isOnline
            ? "bg-[var(--color-success)] text-[var(--color-background)]"
            : "bg-[var(--color-warning)] text-[var(--color-background)]"
        }
      `}
      style={{ paddingTop: "max(env(safe-area-inset-top), 12px)" }}
    >
      {isOnline ? (
        <>
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
          </svg>
          <span>Back online! Changes will sync automatically.</span>
        </>
      ) : (
        <>
          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M18.364 5.636a9 9 0 010 12.728m-3.536-3.536a4 4 0 010-5.656m-7.072 7.072a4 4 0 010-5.656m-3.536 3.536a9 9 0 010-12.728"
            />
            <line x1="4" y1="4" x2="20" y2="20" stroke="currentColor" strokeWidth={2} />
          </svg>
          <span>Working offline - changes will sync when connected</span>
        </>
      )}
    </div>
  );
}
