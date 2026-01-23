"use client";

import { useEffect, useState } from "react";
import { getSyncQueueManager, type SyncQueueStatus } from "@/lib/sync/queue";
import { type ConnectionState } from "@/lib/api/websocket";
import {
  TIME_JUST_NOW_THRESHOLD,
  TIME_ONE_MINUTE_THRESHOLD,
  TIME_ONE_HOUR_THRESHOLD,
  TIME_TWO_HOURS_THRESHOLD,
  SECONDS_PER_MINUTE,
  SECONDS_PER_HOUR,
} from "@/lib/constants";

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

  if (seconds < TIME_JUST_NOW_THRESHOLD) return "just now";
  if (seconds < TIME_ONE_MINUTE_THRESHOLD) return "1m ago";
  if (seconds < TIME_ONE_HOUR_THRESHOLD) return `${Math.floor(seconds / SECONDS_PER_MINUTE)}m ago`;
  if (seconds < TIME_TWO_HOURS_THRESHOLD) return "1h ago";
  return `${Math.floor(seconds / SECONDS_PER_HOUR)}h ago`;
}
