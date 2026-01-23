/**
 * API client for Go sync server
 * Endpoints defined in internal/infra/sync/server.go
 */

import { API_PING_TIMEOUT_MS } from "@/lib/constants";
import type { IdeaType } from "../db";

// Types matching Go sync server (internal/infra/sync/types.go)
export interface CapturedIdea {
  id: string;
  type: IdeaType;
  content: string;
  media_path?: string;
  bpm?: number;
  device_id: string;
  song_id?: string;
  created_at: string;
  synced_at: string;
}

export interface SyncStatus {
  connected: boolean;
  last_sync: string;
  pending_ideas: number;
  connected_devices: number;
}

export interface PairResponse {
  success: boolean;
  device_id?: string;
  error?: string;
}

export interface Song {
  id: string;
  title: string;
  created_at: string;
}

export class SyncClient {
  private baseUrl: string;
  private deviceId: string;

  constructor(baseUrl: string, deviceId: string = "") {
    this.baseUrl = baseUrl.replace(/\/$/, ""); // Remove trailing slash
    this.deviceId = deviceId;
  }

  setDeviceId(deviceId: string): void {
    this.deviceId = deviceId;
  }

  // GET /status - Check server status
  async getStatus(): Promise<SyncStatus> {
    const response = await fetch(`${this.baseUrl}/status`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to get status: ${response.statusText}`);
    }

    return response.json();
  }

  // POST /pair - Pair device with server
  async pair(code: string, deviceName: string): Promise<PairResponse> {
    const response = await fetch(`${this.baseUrl}/pair`, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        code,
        device_name: deviceName,
      }),
    });

    const data = await response.json();

    if (!response.ok) {
      return { success: false, error: data.error || response.statusText };
    }

    if (data.device_id) {
      this.deviceId = data.device_id;
    }

    return { success: true, device_id: data.device_id };
  }

  // GET /ideas - List captured ideas
  async getIdeas(): Promise<CapturedIdea[]> {
    const response = await fetch(`${this.baseUrl}/ideas`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to get ideas: ${response.statusText}`);
    }

    return response.json();
  }

  // POST /ideas - Submit new idea
  async submitIdea(idea: {
    type: IdeaType;
    content: string;
    bpm?: number;
    media_path?: string;
  }): Promise<CapturedIdea> {
    const response = await fetch(`${this.baseUrl}/ideas`, {
      method: "POST",
      headers: this.getHeaders(),
      body: JSON.stringify({
        ...idea,
        device_id: this.deviceId,
        created_at: new Date().toISOString(),
      }),
    });

    if (!response.ok) {
      throw new Error(`Failed to submit idea: ${response.statusText}`);
    }

    return response.json();
  }

  // POST /upload - Upload media file
  async uploadMedia(
    file: Blob,
    filename: string,
    type: "voice" | "photo"
  ): Promise<{ path: string }> {
    const formData = new FormData();
    formData.append("file", file, filename);
    formData.append("type", type);
    formData.append("device_id", this.deviceId);

    const response = await fetch(`${this.baseUrl}/upload`, {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Failed to upload media: ${response.statusText}`);
    }

    return response.json();
  }

  // GET /songs - List available songs for assignment
  async getSongs(): Promise<Song[]> {
    const response = await fetch(`${this.baseUrl}/songs`, {
      method: "GET",
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to get songs: ${response.statusText}`);
    }

    return response.json();
  }

  // Check if server is reachable
  async ping(): Promise<boolean> {
    try {
      const response = await fetch(`${this.baseUrl}/status`, {
        method: "GET",
        headers: this.getHeaders(),
        signal: AbortSignal.timeout(API_PING_TIMEOUT_MS),
      });
      return response.ok;
    } catch {
      return false;
    }
  }

  private getHeaders(): HeadersInit {
    return {
      "Content-Type": "application/json",
      ...(this.deviceId && { "X-Device-ID": this.deviceId }),
    };
  }
}

// Singleton instance - configured after pairing
let clientInstance: SyncClient | null = null;

export function getSyncClient(): SyncClient | null {
  return clientInstance;
}

export function initSyncClient(baseUrl: string, deviceId: string = ""): SyncClient {
  clientInstance = new SyncClient(baseUrl, deviceId);
  return clientInstance;
}

export function clearSyncClient(): void {
  clientInstance = null;
}
