package sync

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/logging"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local network
	},
}

// Message represents a WebSocket message
type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Client represents a connected WebSocket client
type Client struct {
	hub      *WebSocketHub
	conn     *websocket.Conn
	send     chan Message
	deviceID string
}

// WebSocketHub manages WebSocket connections
type WebSocketHub struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	done       chan struct{}
	logger     *logging.Logger
	mu         sync.RWMutex
}

// NewWebSocketHub creates a new WebSocket hub
func NewWebSocketHub(logger *logging.Logger) *WebSocketHub {
	if logger == nil {
		logger = logging.GetDefaultLogger()
	}
	return &WebSocketHub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Message, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		done:       make(chan struct{}),
		logger:     logger,
	}
}

// Run starts the hub's main loop
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			h.logger.Infof("Client connected: %s (total: %d)", client.deviceID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			h.logger.Infof("Client disconnected: %s (total: %d)", client.deviceID, len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send buffer is full, close connection
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()

		case <-h.done:
			h.mu.Lock()
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
			h.mu.Unlock()
			return
		}
	}
}

// Stop stops the hub
func (h *WebSocketHub) Stop() {
	close(h.done)
}

// Broadcast sends a message to all connected clients
func (h *WebSocketHub) Broadcast(msg Message) {
	select {
	case h.broadcast <- msg:
	default:
		h.logger.Warn("Broadcast channel full, dropping message")
	}
}

// ClientCount returns the number of connected clients
func (h *WebSocketHub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// HandleConnection upgrades HTTP to WebSocket and handles the connection
func (h *WebSocketHub) HandleConnection(w http.ResponseWriter, r *http.Request, deviceID string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Errorf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan Message, 256),
		deviceID: deviceID,
	}

	h.register <- client

	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512 * 1024) // 512KB max message size
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.hub.logger.Warnf("WebSocket read error: %v", err)
			}
			break
		}

		// Parse message
		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.hub.logger.Warnf("Invalid message format: %v", err)
			continue
		}

		// Handle message based on type
		c.handleMessage(msg)
	}
}

// writePump writes messages to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				c.hub.logger.Errorf("Failed to marshal message: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages
func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "tempo":
		// Broadcast tempo to all clients
		c.hub.Broadcast(msg)

	case "idea":
		// Broadcast new idea to all clients
		c.hub.Broadcast(msg)

	case "ping":
		// Send pong response
		c.send <- Message{
			Type:    "pong",
			Payload: msg.Payload,
		}

	default:
		c.hub.logger.Debugf("Unknown message type: %s", msg.Type)
	}
}

// SendToDevice sends a message to a specific device
func (h *WebSocketHub) SendToDevice(deviceID string, msg Message) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.deviceID == deviceID {
			select {
			case client.send <- msg:
				return true
			default:
				return false
			}
		}
	}
	return false
}
