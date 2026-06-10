package collaboration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/kyanite/design/icons"
)

// PresenceManager handles user presence tracking and notifications
type PresenceManager struct {
	// Core state
	presence map[string]*UserPresence
	sessions map[string]map[string]*Participant
	mu       sync.RWMutex

	// Context for cleanup
	ctx    context.Context
	cancel context.CancelFunc

	// Event channels
	presenceEvents chan PresenceEvent

	// Callbacks
	onPresenceUpdate func(PresenceUpdateEvent)
	onUserJoined     func(string, string) // sessionID, userID
	onUserLeft       func(string, string) // sessionID, userID
}

// PresenceEvent represents presence-related events
type PresenceEvent struct {
	Type      PresenceEventType      `json:"type"`
	UserID    string                 `json:"user_id"`
	SessionID string                 `json:"session_id,omitempty"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// PresenceEventType defines types of presence events
type PresenceEventType string

const (
	EventUserJoined    PresenceEventType = "user_joined"
	EventUserLeft      PresenceEventType = "user_left"
	EventStatusChanged PresenceEventType = "status_changed"
	EventCursorMoved   PresenceEventType = "cursor_moved"
	EventHeartbeat     PresenceEventType = "heartbeat"
)

// NewPresenceManager creates a new presence manager
func NewPresenceManager() *PresenceManager {
	ctx, cancel := context.WithCancel(context.Background())

	pm := &PresenceManager{
		presence:       make(map[string]*UserPresence),
		sessions:       make(map[string]map[string]*Participant),
		ctx:            ctx,
		cancel:         cancel,
		presenceEvents: make(chan PresenceEvent, 100),
	}

	// Start event processing
	go pm.processPresenceEvents()
	go pm.cleanupIdleUsers()

	return pm
}

// SetPresenceCallbacks configures callbacks for presence events
func (pm *PresenceManager) SetPresenceCallbacks(
	onPresenceUpdate func(PresenceUpdateEvent),
	onUserJoined func(string, string),
	onUserLeft func(string, string),
) {
	pm.onPresenceUpdate = onPresenceUpdate
	pm.onUserJoined = onUserJoined
	pm.onUserLeft = onUserLeft
}

// UpdateUserPresence updates a user's presence status
func (pm *PresenceManager) UpdateUserPresence(userID, username string, status PresenceStatus, sessionID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	presence, exists := pm.presence[userID]
	if !exists {
		presence = &UserPresence{
			UserID:     userID,
			Username:   username,
			Status:     status,
			LastSeen:   time.Now(),
			DeviceInfo: make(map[string]string),
		}
		pm.presence[userID] = presence
	}

	// Update presence information
	oldStatus := presence.Status
	presence.Status = status
	presence.LastSeen = time.Now()
	presence.CurrentSession = sessionID

	// Emit presence update event if status changed
	if oldStatus != status {
		pm.emitPresenceEvent(PresenceEvent{
			Type:      EventStatusChanged,
			UserID:    userID,
			SessionID: sessionID,
			Data: map[string]interface{}{
				"old_status": oldStatus,
				"new_status": status,
			},
			Timestamp: time.Now(),
		})
	}
}

// UpdateCursorPosition updates a user's cursor position
func (pm *PresenceManager) UpdateCursorPosition(userID, sessionID string, position CursorPosition) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	presence, exists := pm.presence[userID]
	if !exists {
		return
	}

	presence.LastSeen = time.Now()

	// Update cursor in session participants if session exists
	if participants, sessionExists := pm.sessions[sessionID]; sessionExists {
		if participant, participantExists := participants[userID]; participantExists {
			participant.Cursor = position
			participant.LastSeen = time.Now()
			pm.sessions[sessionID][userID] = participant
		}
	}

	// Emit cursor movement event
	pm.emitPresenceEvent(PresenceEvent{
		Type:      EventCursorMoved,
		UserID:    userID,
		SessionID: sessionID,
		Data: map[string]interface{}{
			"cursor": position,
		},
		Timestamp: time.Now(),
	})
}

// AddUserToSession adds a user to a session for presence tracking
func (pm *PresenceManager) AddUserToSession(sessionID, userID, username string, role ParticipantRole) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Initialize session participants map if needed
	if _, exists := pm.sessions[sessionID]; !exists {
		pm.sessions[sessionID] = make(map[string]*Participant)
	}

	// Create participant
	participant := &Participant{
		UserID:    userID,
		SessionID: sessionID,
		Username:  username,
		Role:      role,
		JoinedAt:  time.Now(),
		LastSeen:  time.Now(),
		IsActive:  true,
	}

	pm.sessions[sessionID][userID] = participant

	// Update user presence
	pm.UpdateUserPresence(userID, username, StatusOnline, sessionID)

	// Emit user joined event
	pm.emitPresenceEvent(PresenceEvent{
		Type:      EventUserJoined,
		UserID:    userID,
		SessionID: sessionID,
		Timestamp: time.Now(),
	})

	if pm.onUserJoined != nil {
		pm.onUserJoined(sessionID, userID)
	}
}

// RemoveUserFromSession removes a user from a session
func (pm *PresenceManager) RemoveUserFromSession(sessionID, userID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Remove from session participants
	if participants, exists := pm.sessions[sessionID]; exists {
		if participant, participantExists := participants[userID]; participantExists {
			delete(participants, userID)

			// Emit user left event
			pm.emitPresenceEvent(PresenceEvent{
				Type:      EventUserLeft,
				UserID:    userID,
				SessionID: sessionID,
				Timestamp: time.Now(),
			})

			if pm.onUserLeft != nil {
				pm.onUserLeft(sessionID, userID)
			}

			// Update user presence
			pm.UpdateUserPresence(participant.UserID, participant.Username, StatusOffline, "")
		}
	}
}

// GetUserPresence returns presence information for a user
func (pm *PresenceManager) GetUserPresence(userID string) (*UserPresence, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	presence, exists := pm.presence[userID]
	return presence, exists
}

// GetSessionParticipants returns all participants in a session
func (pm *PresenceManager) GetSessionParticipants(sessionID string) ([]*Participant, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	participants, exists := pm.sessions[sessionID]
	if !exists {
		return nil, false
	}

	result := make([]*Participant, 0, len(participants))
	for _, p := range participants {
		result = append(result, p)
	}

	return result, true
}

// GetActiveUsers returns all users currently online
func (pm *PresenceManager) GetActiveUsers() []*UserPresence {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var activeUsers []*UserPresence
	for _, presence := range pm.presence {
		if presence.Status == StatusOnline {
			activeUsers = append(activeUsers, presence)
		}
	}

	return activeUsers
}

// GetUsersByStatus returns users filtered by status
func (pm *PresenceManager) GetUsersByStatus(status PresenceStatus) []*UserPresence {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var users []*UserPresence
	for _, presence := range pm.presence {
		if presence.Status == status {
			users = append(users, presence)
		}
	}

	return users
}

// SendHeartbeat sends a heartbeat for a user to keep them online
func (pm *PresenceManager) SendHeartbeat(userID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	presence, exists := pm.presence[userID]
	if !exists {
		return
	}

	presence.LastSeen = time.Now()

	// Emit heartbeat event
	pm.emitPresenceEvent(PresenceEvent{
		Type:      EventHeartbeat,
		UserID:    userID,
		Timestamp: time.Now(),
	})
}

// SetUserDeviceInfo sets device information for a user
func (pm *PresenceManager) SetUserDeviceInfo(userID string, deviceInfo map[string]string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	presence, exists := pm.presence[userID]
	if !exists {
		return
	}

	presence.DeviceInfo = deviceInfo
	presence.LastSeen = time.Now()
}

// GetPresenceIndicator returns a presence indicator for UI display
func (pm *PresenceManager) GetPresenceIndicator(userID string) PresenceIndicator {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	presence, exists := pm.presence[userID]
	if !exists {
		return PresenceIndicator{
			Status: StatusOffline,
			Color:  "gray",
			Icon:   "[ ]",
		}
	}

	return pm.presenceStatusToIndicator(presence.Status)
}

// GetSessionPresenceIndicators returns presence indicators for all users in a session
func (pm *PresenceManager) GetSessionPresenceIndicators(sessionID string) []SessionPresenceIndicator {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	participants, exists := pm.sessions[sessionID]
	if !exists {
		return nil
	}

	var indicators []SessionPresenceIndicator
	for _, participant := range participants {
		presence, presenceExists := pm.presence[participant.UserID]
		if !presenceExists {
			continue
		}

		indicator := SessionPresenceIndicator{
			UserID:    participant.UserID,
			Username:  participant.Username,
			Indicator: pm.presenceStatusToIndicator(presence.Status),
			Cursor:    participant.Cursor,
			Role:      participant.Role,
		}

		indicators = append(indicators, indicator)
	}

	return indicators
}

// PresenceIndicator represents a visual presence indicator
type PresenceIndicator struct {
	Status  PresenceStatus `json:"status"`
	Color   string         `json:"color"`
	Icon    string         `json:"icon"`
	Tooltip string         `json:"tooltip"`
}

// SessionPresenceIndicator combines presence info with session-specific data
type SessionPresenceIndicator struct {
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	Indicator PresenceIndicator `json:"indicator"`
	Cursor    CursorPosition    `json:"cursor"`
	Role      ParticipantRole   `json:"role"`
}

// Helper methods

func (pm *PresenceManager) emitPresenceEvent(event PresenceEvent) {
	select {
	case pm.presenceEvents <- event:
	default:
		// Channel full, drop event
	}
}

func (pm *PresenceManager) processPresenceEvents() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PresenceManager] Panic in event processor: %v\n", r)
		}
	}()

	for {
		select {
		case <-pm.ctx.Done():
			return
		case event := <-pm.presenceEvents:
			pm.handlePresenceEvent(event)
		}
	}
}

func (pm *PresenceManager) handlePresenceEvent(event PresenceEvent) {
	switch event.Type {
	case EventStatusChanged:
		if pm.onPresenceUpdate != nil {
			pm.onPresenceUpdate(PresenceUpdateEvent{
				UserID:   event.UserID,
				Presence: pm.presence[event.UserID],
				Action:   "status_changed",
			})
		}
	case EventCursorMoved:
		// Cursor movements are handled internally
		// Could emit UI update events here if needed
	}
}

func (pm *PresenceManager) cleanupIdleUsers() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[PresenceManager] Panic in cleanup processor: %v\n", r)
		}
	}()

	ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-pm.ctx.Done():
			return
		case <-ticker.C:
			pm.markIdleUsersOffline()
		}
	}
}

func (pm *PresenceManager) markIdleUsersOffline() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	now := time.Now()
	idleThreshold := 5 * time.Minute // Consider users idle after 5 minutes

	var idleUsers []string
	for userID, presence := range pm.presence {
		if presence.Status == StatusOnline && now.Sub(presence.LastSeen) > idleThreshold {
			idleUsers = append(idleUsers, userID)
		}
	}

	// Mark idle users as away
	for _, userID := range idleUsers {
		presence := pm.presence[userID]
		presence.Status = StatusAway
		pm.presence[userID] = presence

		pm.emitPresenceEvent(PresenceEvent{
			Type:   EventStatusChanged,
			UserID: userID,
			Data: map[string]interface{}{
				"old_status": StatusOnline,
				"new_status": StatusAway,
				"reason":     "idle_timeout",
			},
			Timestamp: time.Now(),
		})
	}
}

func (pm *PresenceManager) presenceStatusToIndicator(status PresenceStatus) PresenceIndicator {
	switch status {
	case StatusOnline:
		return PresenceIndicator{
			Status:  StatusOnline,
			Color:   "green",
			Icon:    icons.GetIcon("online"),
			Tooltip: "Online",
		}
	case StatusAway:
		return PresenceIndicator{
			Status:  StatusAway,
			Color:   "yellow",
			Icon:    icons.GetIcon("away"),
			Tooltip: "Away",
		}
	case StatusBusy:
		return PresenceIndicator{
			Status:  StatusBusy,
			Color:   "red",
			Icon:    icons.GetIcon("busy"),
			Tooltip: "Busy",
		}
	case StatusOffline:
		return PresenceIndicator{
			Status:  StatusOffline,
			Color:   "gray",
			Icon:    icons.GetIcon("offline"),
			Tooltip: "Offline",
		}
	default:
		return PresenceIndicator{
			Status:  StatusOffline,
			Color:   "gray",
			Icon:    icons.GetIcon("offline"),
			Tooltip: "Unknown",
		}
	}
}

// Close shuts down the presence manager
func (pm *PresenceManager) Close() error {
	pm.cancel()
	return nil
}
