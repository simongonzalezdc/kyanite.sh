"use client";

import { useEffect, useState } from "react";
import { getSyncQueueManager, type SyncQueueStatus } from "@/lib/sync/queue";
import { type ConnectionState } from "@/lib/api/websocket";

interface SyncStatusIndicatorProps {
  connectionState: ConnectionState;
}

export function SyncStatusIndicator({ connectionState }: SyncStatusIndicatorProps) {
  const [syncStatus, setSyncStatus] = useState<SyncQueueStatus | null>(null);

  useEffect(() => {
    const manager = getSyncQueueManager();
    const unsubscribe = manager.subscribe(setSyncStatus);
    manager.updateCounts();
    
    // Start auto-sync when connected
    if (connectionState === "connected") {
      manager.startAutoSync(30000); // Every 30 seconds
    }

    return () => {
      unsubscribe();
      manager.stopAutoSync();
    };
  }, [connectionState]);

  if (!syncStatus) return null;

  const totalPending = syncStatus.pendingIdeas + syncStatus.pendingMedia;
  
  // Determine overall status
  let statusIcon: string;
  let statusText: string;
  let statusColor: string;

  if (connectionState === "error") {
    statusIcon = "✕";
    statusText = "Connection error";
    statusColor = "text-[var(--color-error)]";
  } else if (syncStatus.state === "syncing") {
    statusIcon = "◐";
    statusText = "Syncing...";
    statusColor = "text-[var(--color-warning)]";
  } else if (totalPending > 0 && connectionState !== "connected") {
    statusIcon = "⚠";
    statusText = `${totalPending} pending`;
    statusColor = "text-[var(--color-warning)]";
  } else if (connectionState === "connected") {
    statusIcon = "●";
    statusText = "Synced";
    statusColor = "text-[var(--color-success)]";
  } else if (connectionState === "connecting") {
    statusIcon = "◐";
    statusText = "Connecting...";
    statusColor = "text-[var(--color-text-muted)]";
  } else {
    statusIcon = "○";
    statusText = "Offline";
    statusColor = "text-[var(--color-text-muted)]";
  }

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className={`${statusColor} ${syncStatus.state === "syncing" ? "animate-spin" : ""}`}>
        {statusIcon}
      </span>
      <span className="text-[var(--color-text-muted)]">
        {statusText}
        {syncStatus.lastSync && connectionState === "connected" && totalPending === 0 && (
          <span className="ml-1 opacity-60">
            · {formatRelativeTime(syncStatus.lastSync)}
          </span>
        )}
      </span>
    </div>
  );
}

function formatRelativeTime(date: Date): string {
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  
  if (seconds < 60) return "just now";
  if (seconds < 120) return "1m ago";
  if (seconds < 3600) return `${Math.floor(seconds / 60)}m ago`;
  if (seconds < 7200) return "1h ago";
  return `${Math.floor(seconds / 3600)}h ago`;
}
