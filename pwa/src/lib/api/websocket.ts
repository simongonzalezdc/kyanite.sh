/**
 * WebSocket client for real-time sync with Go server
 * Protocol defined in internal/infra/sync/websocket.go
 */

export type MessageType = "idea_received" | "sync_complete" | "ping" | "pong" | "connected" | "error";

export interface WSMessage {
  type: MessageType;
  payload: unknown;
}

export type ConnectionState = "disconnected" | "connecting" | "connected" | "error";

export type MessageHandler = (message: WSMessage) => void;
export type StateChangeHandler = (state: ConnectionState) => void;

export class SyncWebSocket {
  private ws: WebSocket | null = null;
  private url: string;
  private deviceId: string;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private pingInterval: ReturnType<typeof setInterval> | null = null;
  private messageHandlers: Set<MessageHandler> = new Set();
  private stateHandlers: Set<StateChangeHandler> = new Set();
  private _state: ConnectionState = "disconnected";

  constructor(serverUrl: string, deviceId: string) {
    // Convert http:// to ws://
    this.url = serverUrl.replace(/^http/, "ws") + "/ws";
    this.deviceId = deviceId;
  }

  get state(): ConnectionState {
    return this._state;
  }

  private setState(state: ConnectionState): void {
    this._state = state;
    this.stateHandlers.forEach((handler) => handler(state));
  }

  connect(): void {
    if (this.ws?.readyState === WebSocket.OPEN) {
      return;
    }

    this.setState("connecting");

    try {
      this.ws = new WebSocket(`${this.url}?device_id=${encodeURIComponent(this.deviceId)}`);

      this.ws.onopen = () => {
        this.reconnectAttempts = 0;
        this.setState("connected");
        this.startPingInterval();
        this.notifyHandlers({ type: "connected", payload: null });
      };

      this.ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data);
          
          // Handle pong internally
          if (message.type === "pong") {
            return;
          }

          this.notifyHandlers(message);
        } catch (error) {
          console.error("Failed to parse WebSocket message:", error);
        }
      };

      this.ws.onclose = () => {
        this.stopPingInterval();
        this.setState("disconnected");
        this.attemptReconnect();
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        this.setState("error");
      };
    } catch (error) {
      console.error("Failed to create WebSocket:", error);
      this.setState("error");
      this.attemptReconnect();
    }
  }

  disconnect(): void {
    this.reconnectAttempts = this.maxReconnectAttempts; // Prevent reconnection
    this.stopPingInterval();
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.setState("disconnected");
  }

  send(message: WSMessage): boolean {
    if (this.ws?.readyState !== WebSocket.OPEN) {
      return false;
    }

    try {
      this.ws.send(JSON.stringify(message));
      return true;
    } catch (error) {
      console.error("Failed to send WebSocket message:", error);
      return false;
    }
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  onStateChange(handler: StateChangeHandler): () => void {
    this.stateHandlers.add(handler);
    // Immediately notify of current state
    handler(this._state);
    return () => this.stateHandlers.delete(handler);
  }

  private notifyHandlers(message: WSMessage): void {
    this.messageHandlers.forEach((handler) => handler(message));
  }

  private startPingInterval(): void {
    this.stopPingInterval();
    this.pingInterval = setInterval(() => {
      this.send({ type: "ping", payload: null });
    }, 30000); // Ping every 30 seconds
  }

  private stopPingInterval(): void {
    if (this.pingInterval) {
      clearInterval(this.pingInterval);
      this.pingInterval = null;
    }
  }

  private attemptReconnect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.log("Max reconnection attempts reached");
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
    
    setTimeout(() => {
      this.connect();
    }, delay);
  }
}

// Singleton instance
let wsInstance: SyncWebSocket | null = null;

export function getWebSocket(): SyncWebSocket | null {
  return wsInstance;
}

export function initWebSocket(serverUrl: string, deviceId: string): SyncWebSocket {
  if (wsInstance) {
    wsInstance.disconnect();
  }
  wsInstance = new SyncWebSocket(serverUrl, deviceId);
  return wsInstance;
}

export function clearWebSocket(): void {
  if (wsInstance) {
    wsInstance.disconnect();
    wsInstance = null;
  }
}
