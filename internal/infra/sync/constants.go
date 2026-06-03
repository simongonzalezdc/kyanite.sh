// Package sync provides local network synchronization between the TUI and PWA companion.
package sync

import "time"

// Server configuration constants
const (
	// DefaultServerPort is the default port for the sync server
	DefaultServerPort = 8765

	// DefaultMediaPath is the default path for media storage
	DefaultMediaPath = "data/sync/media"
)

// HTTP server timeout constants
const (
	// ServerReadTimeout is the maximum duration for reading the entire request
	ServerReadTimeout = 15 * time.Second

	// ServerWriteTimeout is the maximum duration before timing out writes of the response
	ServerWriteTimeout = 15 * time.Second

	// ServerIdleTimeout is the maximum amount of time to wait for the next request
	ServerIdleTimeout = 60 * time.Second

	// ServerShutdownTimeout is the maximum duration to wait for active connections to close
	ServerShutdownTimeout = 5 * time.Second
)

// WebSocket configuration constants
const (
	// WebSocketReadBufferSize is the buffer size for WebSocket reads
	WebSocketReadBufferSize = 1024

	// WebSocketWriteBufferSize is the buffer size for WebSocket writes
	WebSocketWriteBufferSize = 1024

	// WebSocketBroadcastBuffer is the buffer size for the broadcast channel
	WebSocketBroadcastBuffer = 256

	// WebSocketMaxMessageSize is the maximum size of a WebSocket message (512KB)
	WebSocketMaxMessageSize = 512 * 1024

	// WebSocketReadDeadline is the read deadline for WebSocket connections
	WebSocketReadDeadline = 60 * time.Second

	// WebSocketWriteDeadline is the write deadline for WebSocket connections
	WebSocketWriteDeadline = 10 * time.Second

	// WebSocketPingInterval is the interval between ping messages
	WebSocketPingInterval = 30 * time.Second
)

// File upload constants
const (
	// MaxMultipartFormSize is the maximum size for multipart form data (50MB)
	MaxMultipartFormSize = 50 << 20
)

// File permission constants
const (
	// DefaultDirPermission is the default permission for directories
	DefaultDirPermission = 0o755

	// DefaultFilePermission is the default permission for files
	DefaultFilePermission = 0o644
)
