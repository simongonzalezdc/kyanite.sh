// Package sync provides local network synchronization between the TUI and PWA companion.
package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/logging"
)

// SyncServer provides an embedded HTTP/WebSocket server for PWA sync
type SyncServer struct {
	httpServer     *http.Server
	wsHub          *WebSocketHub
	config         *config.SyncConfig
	pairingManager *PairingManager
	mediaStore     *MediaStore
	logger         *logging.Logger

	// State
	running bool
	localIP string
	mu      sync.RWMutex

	// Callbacks
	onIdeaReceived func(idea *CapturedIdea)
	onDevicePaired func(device *PairedDevice)
}

// SyncServerConfig holds configuration for the sync server
type SyncServerConfig struct {
	Port      int
	MediaPath string
	Logger    *logging.Logger
}

// DefaultSyncServerConfig returns the default configuration
func DefaultSyncServerConfig() SyncServerConfig {
	return SyncServerConfig{
		Port:      DefaultServerPort,
		MediaPath: DefaultMediaPath,
		Logger:    logging.GetDefaultLogger(),
	}
}

// NewSyncServerWithAutoSetup creates a sync server with automatic directory setup
// This is the recommended way to create a sync server for end-user apps
func NewSyncServerWithAutoSetup(dataDir string, port int, logger *logging.Logger) (*SyncServer, error) {
	mediaPath := dataDir + "/sync/media"

	cfg := SyncServerConfig{
		Port:      port,
		MediaPath: mediaPath,
		Logger:    logger,
	}

	return NewSyncServer(cfg)
}

// NewSyncServer creates a new sync server
func NewSyncServer(cfg SyncServerConfig) (*SyncServer, error) {
	if cfg.Logger == nil {
		cfg.Logger = logging.GetDefaultLogger()
	}

	// Create components
	wsHub := NewWebSocketHub(cfg.Logger)
	pairingManager := NewPairingManager(cfg.Logger)
	mediaStore, err := NewMediaStore(cfg.MediaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create media store: %w", err)
	}

	syncConfig := &config.SyncConfig{
		Enabled:   true,
		Port:      cfg.Port,
		MediaPath: cfg.MediaPath,
	}

	server := &SyncServer{
		wsHub:          wsHub,
		config:         syncConfig,
		pairingManager: pairingManager,
		mediaStore:     mediaStore,
		logger:         cfg.Logger,
	}

	return server, nil
}

// Start starts the sync server (non-blocking)
func (s *SyncServer) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("server already running")
	}

	// Get local IP
	localIP, err := getLocalIP()
	if err != nil {
		s.logger.Warnf("Could not determine local IP: %v", err)
		localIP = "127.0.0.1"
	}
	s.localIP = localIP

	// Create router
	mux := http.NewServeMux()
	s.setupRoutes(mux)

	// Create server
	addr := fmt.Sprintf("%s:%d", localIP, s.config.Port)
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.corsMiddleware(mux),
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	}

	// Start WebSocket hub
	go s.wsHub.Run()

	// Start server
	go func() {
		s.logger.Infof("Sync server starting on http://%s", addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Errorf("Sync server error: %v", err)
		}
	}()

	s.running = true
	return nil
}

// Stop stops the sync server
func (s *SyncServer) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), ServerShutdownTimeout)
	defer cancel()

	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("server shutdown error: %w", err)
		}
	}

	s.wsHub.Stop()
	s.running = false
	s.logger.Info("Sync server stopped")
	return nil
}

// IsRunning returns whether the server is running
func (s *SyncServer) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// GetLocalURL returns the local URL for the sync server
func (s *SyncServer) GetLocalURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprintf("http://%s:%d", s.localIP, s.config.Port)
}

// GetConnectedDevices returns the number of connected devices
func (s *SyncServer) GetConnectedDevices() int {
	return s.wsHub.ClientCount()
}

// GeneratePairingCode generates a new pairing code
func (s *SyncServer) GeneratePairingCode() string {
	return s.pairingManager.GenerateCode()
}

// OnIdeaReceived sets a callback for when an idea is received
func (s *SyncServer) OnIdeaReceived(fn func(idea *CapturedIdea)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onIdeaReceived = fn
}

// OnDevicePaired sets a callback for when a device is paired
func (s *SyncServer) OnDevicePaired(fn func(device *PairedDevice)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onDevicePaired = fn
}

// BroadcastMessage sends a message to all connected clients
func (s *SyncServer) BroadcastMessage(msgType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	msg := Message{
		Type:    msgType,
		Payload: data,
	}

	s.wsHub.Broadcast(msg)
	return nil
}

// setupRoutes configures HTTP routes
func (s *SyncServer) setupRoutes(mux *http.ServeMux) {
	// API routes
	mux.HandleFunc("/api/status", s.handleStatus)
	mux.HandleFunc("/api/pair", s.handlePair)
	mux.HandleFunc("/api/ideas", s.handleIdeas)
	mux.HandleFunc("/api/ideas/upload", s.handleUpload)
	mux.HandleFunc("/api/songs", s.handleSongs)

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.handleWebSocket)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

// corsMiddleware adds CORS and security headers for PWA access
func (s *SyncServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS headers - Allow all origins for local network sync
		// This is intentional: the sync server runs on the local network
		// and needs to accept requests from PWA running on any device
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Device-ID, X-Pairing-Code")

		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// CSP for API endpoints - restrictive since this is a JSON API
		w.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// handleStatus returns server status
func (s *SyncServer) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "running",
		"version":   "1.0.0",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"connected": s.wsHub.ClientCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// handlePair handles device pairing
func (s *SyncServer) handlePair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	code := r.Header.Get("X-Pairing-Code")
	deviceName := r.Header.Get("X-Device-Name")
	if deviceName == "" {
		deviceName = "Unknown Device"
	}

	if !s.pairingManager.ValidateCode(code) {
		http.Error(w, "Invalid pairing code", http.StatusUnauthorized)
		return
	}

	// Create paired device
	device := s.pairingManager.PairDevice(deviceName)

	// Notify callback
	s.mu.RLock()
	callback := s.onDevicePaired
	s.mu.RUnlock()
	if callback != nil {
		callback(device)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"device_id": device.ID,
		"message":   "Device paired successfully",
	})
}

// handleIdeas handles idea CRUD operations
func (s *SyncServer) handleIdeas(w http.ResponseWriter, r *http.Request) {
	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		http.Error(w, "Device ID required", http.StatusUnauthorized)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// List ideas - would query database
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]CapturedIdea{})

	case http.MethodPost:
		// Create new idea
		var idea CapturedIdea
		if err := json.NewDecoder(r.Body).Decode(&idea); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		idea.DeviceID = deviceID
		idea.SyncedAt = time.Now()

		// Notify callback
		s.mu.RLock()
		callback := s.onIdeaReceived
		s.mu.RUnlock()
		if callback != nil {
			callback(&idea)
		}

		// Broadcast to other clients
		s.wsHub.Broadcast(Message{
			Type:    "idea_created",
			Payload: mustMarshal(idea),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(idea)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleUpload handles media file uploads
func (s *SyncServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	deviceID := r.Header.Get("X-Device-ID")
	if deviceID == "" {
		http.Error(w, "Device ID required", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(MaxMultipartFormSize); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file content
	data := make([]byte, header.Size)
	if _, err := file.Read(data); err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Determine file type from Content-Type header
	contentType := header.Header.Get("Content-Type")
	var mediaPath string

	// Validate file type with magic bytes
	if isAudioType(contentType) {
		if !validateAudioMagicBytes(data, contentType) {
			http.Error(w, "Invalid audio file: content does not match declared type", http.StatusBadRequest)
			return
		}
		mediaPath, err = s.mediaStore.SaveVoiceMemo(deviceID, data)
	} else if isImageType(contentType) {
		if !validateImageMagicBytes(data, contentType) {
			http.Error(w, "Invalid image file: content does not match declared type", http.StatusBadRequest)
			return
		}
		mediaPath, err = s.mediaStore.SavePhoto(deviceID, data)
	} else {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// mediaStore already returns relative paths, no need to sanitize
	mediaID := extractMediaID(mediaPath)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path":     mediaPath,
		"media_id": mediaID,
	})
}

// handleSongs returns available songs for idea assignment
func (s *SyncServer) handleSongs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This would query the database for songs
	// For now, return empty list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode([]map[string]interface{}{})
}

// handleWebSocket upgrades connection to WebSocket
func (s *SyncServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device_id")
	if deviceID == "" {
		http.Error(w, "Device ID required", http.StatusUnauthorized)
		return
	}

	s.wsHub.HandleConnection(w, r, deviceID)
}

// getLocalIP returns the local IP address for the network interface
func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable network interface found")
}

func isAudioType(contentType string) bool {
	return contentType == "audio/webm" ||
		contentType == "audio/mp3" ||
		contentType == "audio/wav" ||
		contentType == "audio/ogg" ||
		contentType == "audio/m4a"
}

func isImageType(contentType string) bool {
	return contentType == "image/jpeg" ||
		contentType == "image/png" ||
		contentType == "image/gif" ||
		contentType == "image/webp"
}

// validateImageMagicBytes checks if the file content matches the declared image type
func validateImageMagicBytes(data []byte, contentType string) bool {
	if len(data) < 4 {
		return false
	}

	switch contentType {
	case "image/jpeg":
		// JPEG magic bytes: FF D8 FF
		return data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF
	case "image/png":
		// PNG magic bytes: 89 50 4E 47 (89 PNG)
		return data[0] == 0x89 && data[1] == 0x50 && data[2] == 0x4E && data[3] == 0x47
	case "image/gif":
		// GIF magic bytes: GIF8
		return len(data) >= 4 && string(data[0:4]) == "GIF8"
	case "image/webp":
		// WebP magic bytes: RIFF....WEBP
		return len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP"
	default:
		return false
	}
}

// validateAudioMagicBytes checks if the file content matches the declared audio type
func validateAudioMagicBytes(data []byte, contentType string) bool {
	if len(data) < 4 {
		return false
	}

	switch contentType {
	case "audio/webm":
		// WebM magic bytes: 1A 45 DF A3 (EBML header)
		return data[0] == 0x1A && data[1] == 0x45 && data[2] == 0xDF && data[3] == 0xA3
	case "audio/ogg":
		// OGG magic bytes: OggS
		return string(data[0:4]) == "OggS"
	case "audio/mp3":
		// MP3 magic bytes: ID3 (for tagged MP3s) or FF FB/FF FA (for raw MP3)
		return (len(data) >= 3 && string(data[0:3]) == "ID3") ||
			(data[0] == 0xFF && (data[1] == 0xFB || data[1] == 0xFA || data[1] == 0xF3 || data[1] == 0xF2))
	case "audio/wav":
		// WAV magic bytes: RIFF....WAVE
		return len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WAVE"
	case "audio/m4a":
		// M4A magic bytes: 00 00 00 20 66 74 79 70 (ftyp M4A )
		return len(data) >= 8 && data[4] == 0x66 && data[5] == 0x74 && data[6] == 0x79 && data[7] == 0x70
	default:
		return false
	}
}

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		// Log error and return empty JSON object as fallback
		fmt.Printf("[SyncServer] Failed to marshal data: %v\n", err)
		return []byte("{}")
	}
	return data
}

// sanitizeMediaPath converts a full filesystem path to a relative path for API responses
// This prevents leaking internal directory structure (e.g., /home/user/.kyanite/... -> /media/voice/file.webm)
func sanitizeMediaPath(fullPath string) string {
	// Extract just the filename and return a relative path
	filename := fullPath
	if lastSlash := len(fullPath) - 1; lastSlash >= 0 {
		for i := lastSlash; i >= 0; i-- {
			if fullPath[i] == '/' || fullPath[i] == '\\' {
				filename = fullPath[i+1:]
				break
			}
		}
	}

	// Determine if it's a voice memo or photo based on file extension
	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		if ext == ".webm" || ext == ".mp3" || ext == ".wav" || ext == ".ogg" || ext == ".m4a" {
			return "/media/voice/" + filename
		}
	}
	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		if ext == ".jpg" || ext == ".png" || ext == ".gif" {
			return "/media/photos/" + filename
		}
	}

	// Handle .jpeg extension (5 chars)
	if len(filename) > 5 {
		ext := filename[len(filename)-5:]
		if ext == ".jpeg" {
			return "/media/photos/" + filename
		}
	}

	// Fallback: return just the filename if we can't determine type
	return "/media/" + filename
}

// extractMediaID extracts just the filename without extension for use as a media ID
func extractMediaID(fullPath string) string {
	filename := fullPath
	if lastSlash := len(fullPath) - 1; lastSlash >= 0 {
		for i := lastSlash; i >= 0; i-- {
			if fullPath[i] == '/' || fullPath[i] == '\\' {
				filename = fullPath[i+1:]
				break
			}
		}
	}

	// Remove file extension
	if lastDot := len(filename) - 1; lastDot >= 0 {
		for i := lastDot; i >= 0; i-- {
			if filename[i] == '.' {
				return filename[:i]
			}
		}
	}

	return filename
}
