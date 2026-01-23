"use client";

import { type ConnectionState } from "@/lib/api/websocket";

interface ConnectionStatusProps {
  state: ConnectionState;
  serverUrl?: string;
  pendingCount?: number;
}

export function ConnectionStatus({ state, serverUrl, pendingCount = 0 }: ConnectionStatusProps) {
  const statusConfig = {
    connected: {
      icon: "●",
      text: "Connected",
      color: "text-[var(--color-success)]",
    },
    connecting: {
      icon: "◐",
      text: "Connecting...",
      color: "text-[var(--color-warning)]",
    },
    disconnected: {
      icon: "○",
      text: "Offline",
      color: "text-[var(--color-text-muted)]",
    },
    error: {
      icon: "✕",
      text: "Error",
      color: "text-[var(--color-error)]",
    },
  };

  const config = statusConfig[state];

  return (
    <div className="flex items-center gap-2 text-sm">
      <span className={`${config.color} ${state === "connecting" ? "animate-pulse" : ""}`}>
        {config.icon}
      </span>
      <span className="text-[var(--color-text-muted)]">
        {config.text}
        {serverUrl && state === "connected" && (
          <span className="ml-1 opacity-60">· {new URL(serverUrl).hostname}</span>
        )}
        {pendingCount > 0 && state !== "connected" && (
          <span className="ml-1 text-[var(--color-warning)]">· {pendingCount} pending</span>
        )}
      </span>
    </div>
  );
}
