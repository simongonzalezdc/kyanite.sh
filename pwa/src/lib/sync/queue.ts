/**
 * Offline sync queue manager
 * Handles queued ideas and media when connection is restored
 */

import {
  getPendingIdeas,
  getPendingMedia,
  markIdeaSynced,
  markMediaSynced,
  clearSyncedItems,
} from "../db";
import { getSyncClient } from "../api/client";
import { SYNC_INTERVAL_MS } from "@/lib/constants";

export type SyncState = "idle" | "syncing" | "error";

export interface SyncQueueStatus {
  state: SyncState;
  pendingIdeas: number;
  pendingMedia: number;
  lastSync: Date | null;
  lastError: string | null;
}

type SyncCallback = (status: SyncQueueStatus) => void;

class SyncQueueManager {
  private status: SyncQueueStatus = {
    state: "idle",
    pendingIdeas: 0,
    pendingMedia: 0,
    lastSync: null,
    lastError: null,
  };
  
  private callbacks: Set<SyncCallback> = new Set();
  private syncInProgress = false;
  private syncInterval: ReturnType<typeof setInterval> | null = null;

  constructor() {
    // Listen for service worker sync messages
    if (typeof navigator !== "undefined" && "serviceWorker" in navigator) {
      navigator.serviceWorker.addEventListener("message", (event) => {
        if (event.data?.type === "SYNC_REQUESTED") {
          this.sync();
        }
      });
    }
  }

  // Subscribe to status changes
  subscribe(callback: SyncCallback): () => void {
    this.callbacks.add(callback);
    callback(this.status);
    return () => this.callbacks.delete(callback);
  }

  // Start automatic sync polling
  startAutoSync(intervalMs = SYNC_INTERVAL_MS): void {
    if (this.syncInterval) return;
    
    this.syncInterval = setInterval(() => {
      this.sync();
    }, intervalMs);
    
    // Sync immediately
    this.sync();
  }

  // Stop automatic sync polling
  stopAutoSync(): void {
    if (this.syncInterval) {
      clearInterval(this.syncInterval);
      this.syncInterval = null;
    }
  }

  // Update pending counts
  async updateCounts(): Promise<void> {
    const ideas = await getPendingIdeas();
    const media = await getPendingMedia();
    
    this.status.pendingIdeas = ideas.length;
    this.status.pendingMedia = media.length;
    
    this.notifyCallbacks();
  }

  // Trigger sync
  async sync(): Promise<boolean> {
    if (this.syncInProgress) return false;
    
    const client = getSyncClient();
    if (!client) {
      return false;
    }

    // Check if server is reachable
    const reachable = await client.ping();
    if (!reachable) {
      return false;
    }

    this.syncInProgress = true;
    this.status.state = "syncing";
    this.status.lastError = null;
    this.notifyCallbacks();

    try {
      // Get pending items
      const pendingIdeas = await getPendingIdeas();
      const pendingMedia = await getPendingMedia();

      // Sync media first (ideas may reference media paths)
      const mediaPathMap = new Map<string, string>();
      
      for (const media of pendingMedia) {
        try {
          const { path } = await client.uploadMedia(
            media.blob,
            media.filename,
            media.type
          );
          mediaPathMap.set(media.localIdeaId, path);
          await markMediaSynced(media.localIdeaId);
        } catch (err) {
          console.error("Failed to sync media:", err);
          // Continue with other items
        }
      }

      // Sync ideas
      for (const idea of pendingIdeas) {
        try {
          const mediaPath = mediaPathMap.get(idea.localId);
          
          await client.submitIdea({
            type: idea.type,
            content: idea.content,
            bpm: idea.bpm,
            ...(mediaPath && { media_path: mediaPath }),
          });
          
          await markIdeaSynced(idea.localId);
        } catch (err) {
          console.error("Failed to sync idea:", err);
          // Continue with other items
        }
      }

      // Clean up synced items
      await clearSyncedItems();

      // Update counts
      await this.updateCounts();

      this.status.state = "idle";
      this.status.lastSync = new Date();
      this.notifyCallbacks();
      
      return true;
    } catch (err) {
      console.error("Sync error:", err);
      this.status.state = "error";
      this.status.lastError = err instanceof Error ? err.message : "Sync failed";
      this.notifyCallbacks();
      return false;
    } finally {
      this.syncInProgress = false;
    }
  }

  // Register for background sync (when supported)
  async registerBackgroundSync(): Promise<boolean> {
    if (typeof navigator === "undefined" || !("serviceWorker" in navigator)) {
      return false;
    }

    try {
      const registration = await navigator.serviceWorker.ready;

      // Type-safe check for Background Sync API
      // The sync property is not in the standard ServiceWorkerRegistration type
      // but is available in browsers that support the Background Sync API
      interface BackgroundSyncManager {
        register: (tag: string) => Promise<void>;
      }
      type ServiceWorkerRegistrationWithSync = ServiceWorkerRegistration & {
        sync?: BackgroundSyncManager;
      };

      const regWithSync = registration as ServiceWorkerRegistrationWithSync;
      if (regWithSync.sync) {
        await regWithSync.sync.register("sync-ideas");
        return true;
      }
    } catch (err) {
      console.error("Failed to register background sync:", err);
    }

    return false;
  }

  private notifyCallbacks(): void {
    this.callbacks.forEach((cb) => cb({ ...this.status }));
  }
}

// Singleton instance
let instance: SyncQueueManager | null = null;

export function getSyncQueueManager(): SyncQueueManager {
  if (!instance) {
    instance = new SyncQueueManager();
  }
  return instance;
}
