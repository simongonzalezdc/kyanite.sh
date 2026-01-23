/// <reference lib="webworker" />
import { defaultCache } from "@serwist/next/worker";
import { Serwist, type PrecacheEntry } from "serwist";

declare const self: ServiceWorkerGlobalScope & {
  __SW_MANIFEST: (PrecacheEntry | string)[];
};

export {};

const serwist = new Serwist({
  precacheEntries: self.__SW_MANIFEST,
  skipWaiting: true,
  clientsClaim: true,
  navigationPreload: true,
  runtimeCaching: defaultCache,
});

serwist.addEventListeners();

// Handle background sync for offline ideas
self.addEventListener("sync", (event) => {
  if (event.tag === "sync-ideas") {
    event.waitUntil(syncPendingIdeas());
  }
});

async function syncPendingIdeas(): Promise<void> {
  // This will be called when connection is restored
  // The actual sync logic is in the main app via Dexie
  console.log("[SW] Background sync triggered for ideas");
  
  // Notify all clients to sync
  const clients = await self.clients.matchAll({ type: "window" });
  for (const client of clients) {
    client.postMessage({ type: "SYNC_REQUESTED" });
  }
}
