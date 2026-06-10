/**
 * IndexedDB database for offline queue
 * Uses Dexie.js for type-safe IndexedDB access
 */

import Dexie, { type Table } from "dexie";

// Types matching Go sync server (internal/infra/sync/types.go)
export type IdeaType = "text" | "voice_memo" | "photo" | "tempo";

export interface PendingIdea {
  id?: number;
  localId: string;
  type: IdeaType;
  content: string;
  bpm?: number;
  createdAt: string;
  synced: boolean;
}

export interface PendingMedia {
  id?: number;
  localIdeaId: string;
  blob: Blob;
  filename: string;
  type: "voice" | "photo";
  synced: boolean;
}

export interface ServerConfig {
  key: string;
  url: string;
  deviceId: string;
  deviceName: string;
  pairedAt: string;
}

class NoiseDB extends Dexie {
  pendingIdeas!: Table<PendingIdea>;
  pendingMedia!: Table<PendingMedia>;
  serverConfig!: Table<ServerConfig>;

  constructor() {
    super("noise-pwa");

    this.version(1).stores({
      pendingIdeas: "++id, localId, type, createdAt, synced",
      pendingMedia: "++id, localIdeaId, type, synced",
      serverConfig: "key",
    });
  }
}

export const db = new NoiseDB();

// Helper functions

export async function getServerConfig(): Promise<ServerConfig | undefined> {
  return db.serverConfig.get("server");
}

export async function setServerConfig(config: Omit<ServerConfig, "key">): Promise<void> {
  await db.serverConfig.put({ ...config, key: "server" });
}

export async function clearServerConfig(): Promise<void> {
  await db.serverConfig.delete("server");
}

export async function addPendingIdea(idea: Omit<PendingIdea, "id" | "synced">): Promise<number> {
  return db.pendingIdeas.add({ ...idea, synced: false });
}

export async function getPendingIdeas(): Promise<PendingIdea[]> {
  return db.pendingIdeas.where("synced").equals(0).toArray();
}

export async function markIdeaSynced(localId: string): Promise<void> {
  await db.pendingIdeas.where("localId").equals(localId).modify({ synced: true });
}

export async function addPendingMedia(media: Omit<PendingMedia, "id" | "synced">): Promise<number> {
  return db.pendingMedia.add({ ...media, synced: false });
}

export async function getPendingMedia(): Promise<PendingMedia[]> {
  return db.pendingMedia.where("synced").equals(0).toArray();
}

export async function markMediaSynced(localIdeaId: string): Promise<void> {
  await db.pendingMedia.where("localIdeaId").equals(localIdeaId).modify({ synced: true });
}

export async function getPendingCount(): Promise<number> {
  const ideas = await db.pendingIdeas.where("synced").equals(0).count();
  const media = await db.pendingMedia.where("synced").equals(0).count();
  return ideas + media;
}

export async function clearSyncedItems(): Promise<void> {
  await db.pendingIdeas.where("synced").equals(1).delete();
  await db.pendingMedia.where("synced").equals(1).delete();
}
