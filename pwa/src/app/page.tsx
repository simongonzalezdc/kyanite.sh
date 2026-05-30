"use client";

import { useEffect, useState, useCallback } from "react";
import { PairingForm } from "@/components/pairing/PairingForm";
import { SyncStatusIndicator } from "@/components/sync/SyncStatus";
import { getSyncQueueManager } from "@/lib/sync/queue";
import { Button } from "@/components/ui/Button";
import { TextCapture } from "@/components/capture/TextCapture";
import { VoiceMemo } from "@/components/capture/VoiceMemo";
import { PhotoCapture } from "@/components/capture/PhotoCapture";
import { TapTempo } from "@/components/capture/TapTempo";
import { getServerConfig, clearServerConfig } from "@/lib/db";
import { initSyncClient, clearSyncClient } from "@/lib/api/client";
import { initWebSocket, clearWebSocket, type ConnectionState } from "@/lib/api/websocket";
import { SYNC_INTERVAL_MS } from "@/lib/constants";
import { ThemeSelector } from "@/components/ui/ThemeSelector";
import { useSwipe } from "@/lib/hooks/useSwipe";
import { OnboardingWizard, hasCompletedOnboarding } from "@/components/onboarding/OnboardingWizard";
import { useToast } from "@/components/ui/Toast";

type AppState = "loading" | "unpaired" | "paired";
type CaptureTab = "text" | "voice" | "photo" | "tempo";

/**
 * Safely extract hostname from URL, returning fallback on invalid URLs
 */
function getHostname(url: string | null, fallback = "server"): string {
  if (!url) return fallback;
  try {
    return new URL(url).hostname;
  } catch {
    return fallback;
  }
}

const tabs: { id: CaptureTab; icon: string; label: string }[] = [
  { id: "text", icon: "📝", label: "Text" },
  { id: "voice", icon: "🎤", label: "Voice" },
  { id: "photo", icon: "📷", label: "Photo" },
  { id: "tempo", icon: "♪", label: "Tempo" },
];

export default function Home() {
  const [appState, setAppState] = useState<AppState>("loading");
  const [connectionState, setConnectionState] = useState<ConnectionState>("disconnected");
  const [serverUrl, setServerUrl] = useState<string>("");
  const [activeTab, setActiveTab] = useState<CaptureTab>("text");
  const { addToast } = useToast();

  // Initialize onboarding state - show for first-time users once app loads
  const [showOnboarding, setShowOnboarding] = useState(() => {
    if (typeof window === "undefined") return false;
    return !hasCompletedOnboarding();
  });

  // Tab navigation helpers
  const tabIds = tabs.map((t) => t.id);
  const currentTabIndex = tabIds.indexOf(activeTab);

  const goToNextTab = useCallback(() => {
    const nextIndex = (currentTabIndex + 1) % tabIds.length;
    setActiveTab(tabIds[nextIndex]);
  }, [currentTabIndex, tabIds]);

  const goToPrevTab = useCallback(() => {
    const prevIndex = (currentTabIndex - 1 + tabIds.length) % tabIds.length;
    setActiveTab(tabIds[prevIndex]);
  }, [currentTabIndex, tabIds]);

  // Swipe gesture handlers
  const swipeHandlers = useSwipe({
    onSwipeLeft: goToNextTab,
    onSwipeRight: goToPrevTab,
    threshold: 50,
  });

  // Check for existing pairing on mount
  useEffect(() => {
    async function checkPairing() {
      try {
        const config = await getServerConfig();
        if (config) {
          setServerUrl(config.url);
          initSyncClient(config.url, config.deviceId);
          const ws = initWebSocket(config.url, config.deviceId);
          ws.onStateChange(setConnectionState);
          ws.connect();
          setAppState("paired");
        } else {
          setAppState("unpaired");
        }
      } catch (error) {
        console.error("Failed to check pairing:", error);
        addToast({
          type: "error",
          message: "Failed to check pairing status. Please try again.",
        });
        setAppState("unpaired");
      }
    }

    checkPairing();
  }, []);

  // Initialize sync queue manager when paired
  useEffect(() => {
    if (appState !== "paired") return;

    const syncManager = getSyncQueueManager();
    syncManager.updateCounts();

    // Start auto-sync if connected
    if (connectionState === "connected") {
      syncManager.startAutoSync(SYNC_INTERVAL_MS);
    }

    return () => {
      syncManager.stopAutoSync();
    };
  }, [appState, connectionState]);

  const handlePaired = useCallback((url: string, deviceId: string) => {
    setServerUrl(url);
    const ws = initWebSocket(url, deviceId);
    ws.onStateChange(setConnectionState);
    ws.connect();
    setAppState("paired");
  }, []);

  const handleUnpair = useCallback(async () => {
    clearWebSocket();
    clearSyncClient();
    await clearServerConfig();
    setServerUrl("");
    setConnectionState("disconnected");
    setAppState("unpaired");
  }, []);

  const handleCaptured = useCallback(() => {
    // Update sync queue counts
    const syncManager = getSyncQueueManager();
    syncManager.updateCounts();
    
    // Try to sync immediately if connected
    if (connectionState === "connected") {
      syncManager.sync();
    }
  }, [connectionState]);

  // Loading state
  if (appState === "loading") {
    return (
      <main className="min-h-dvh flex items-center justify-center p-4">
        <div className="text-center">
          <div className="text-5xl mb-4 animate-pulse">♪</div>
          <p className="text-[var(--color-text-muted)]">Loading...</p>
        </div>
      </main>
    );
  }

  // Unpaired - show pairing form
  if (appState === "unpaired") {
    return (
      <main className="min-h-dvh flex items-center justify-center p-4">
        {/* Onboarding wizard for first-time users */}
        {showOnboarding && (
          <OnboardingWizard onComplete={() => setShowOnboarding(false)} />
        )}
        <PairingForm onPaired={handlePaired} />
      </main>
    );
  }

  // Paired - show capture interface
  return (
    <main className="min-h-dvh flex flex-col">
      {/* Onboarding wizard for first-time users */}
      {showOnboarding && (
        <OnboardingWizard onComplete={() => setShowOnboarding(false)} />
      )}

      {/* Header */}
      <header className="flex items-center justify-between p-4 border-b border-[var(--color-surface)]">
        <div className="flex items-center gap-2">
          <span className="text-xl">♪</span>
          <span className="font-bold text-[var(--color-primary)]">noise.sh</span>
        </div>
        <div className="flex items-center gap-2">
          <SyncStatusIndicator connectionState={connectionState} />
          <ThemeSelector compact />
        </div>
      </header>

      {/* Tab bar */}
      <nav aria-label="Capture tools" className="flex border-b border-[var(--color-surface)]" role="tablist">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            role="tab"
            aria-selected={activeTab === tab.id}
            aria-controls={`panel-${tab.id}`}
            id={`tab-${tab.id}`}
            className={`flex-1 py-3 text-center transition-colors ${
              activeTab === tab.id
                ? "text-[var(--color-primary)] border-b-2 border-[var(--color-primary)]"
                : "text-[var(--color-text-muted)] hover:text-[var(--color-text)]"
            }`}
          >
            <div className="text-lg" aria-hidden="true">{tab.icon}</div>
            <div className="text-xs mt-1">{tab.label}</div>
          </button>
        ))}
      </nav>

      {/* Capture content - swipeable */}
      <div
        className="flex-1 overflow-hidden"
        role="tabpanel"
        aria-labelledby={`tab-${activeTab}`}
        id={`panel-${activeTab}`}
        {...swipeHandlers}
      >
        {activeTab === "text" && (
          <TextCapture
            onCaptured={handleCaptured}
            disabled={connectionState === "error"}
          />
        )}
        {activeTab === "voice" && (
          <VoiceMemo
            onCaptured={handleCaptured}
            disabled={connectionState === "error"}
          />
        )}
        {activeTab === "photo" && (
          <PhotoCapture
            onCaptured={handleCaptured}
            disabled={connectionState === "error"}
          />
        )}
        {activeTab === "tempo" && (
          <TapTempo
            onCaptured={handleCaptured}
            disabled={connectionState === "error"}
          />
        )}
      </div>

      {/* Footer */}
      <footer className="p-4 border-t border-[var(--color-surface)]">
        <Button
          variant="ghost"
          size="sm"
          onClick={handleUnpair}
          className="w-full text-[var(--color-text-muted)]"
        >
          Disconnect from {getHostname(serverUrl)}
        </Button>
      </footer>
    </main>
  );
}
