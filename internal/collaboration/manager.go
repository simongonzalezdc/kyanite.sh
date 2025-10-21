package collaboration

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/infra/db"
)

// CollaborationManager coordinates all collaborative features
type CollaborationManager struct {
	// Core state
	db           *db.DB
	sessions     map[string]*Session
	participants map[string]map[string]*Participant
	presence     map[string]*UserPresence

	// Synchronization
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc

	// Event channels
	events    chan CollaborationEvent
	broadcast chan CollaborationMessage

	// Callbacks for UI integration
	onSessionUpdate    func(SessionUpdateEvent)
	onPresenceUpdate   func(PresenceUpdateEvent)
	onConflictDetected func(ConflictEvent)
}

// NewCollaborationManager creates a new collaboration manager
func NewCollaborationManager(database *db.DB) *CollaborationManager {
	ctx, cancel := context.WithCancel(context.Background())

	cm := &CollaborationManager{
		db:           database,
		sessions:     make(map[string]*Session),
		participants: make(map[string]map[string]*Participant),
		presence:     make(map[string]*UserPresence),
		ctx:          ctx,
		cancel:       cancel,
		events:       make(chan CollaborationEvent, 100),
		broadcast:    make(chan CollaborationMessage, 100),
	}

	// Start event processing
	go cm.processEvents()
	go cm.processBroadcasts()

	return cm
}

// Session represents a collaborative editing session
type Session struct {
	ID         string          `json:"id"`
	DocumentID int             `json:"document_id"`
	Name       string          `json:"name"`
	CreatedBy  string          `json:"created_by"`
	CreatedAt  time.Time       `json:"created_at"`
	IsActive   bool            `json:"is_active"`
	Settings   SessionSettings `json:"settings"`
	Document   *domain.Song    `json:"document"`
	Operations []Operation     `json:"operations"`
}

// SessionSettings configures session behavior
type SessionSettings struct {
	MaxParticipants  int           `json:"max_participants"`
	AutoSaveInterval time.Duration `json:"auto_save_interval"`
	ConflictStrategy string        `json:"conflict_strategy"` // "merge", "lock", "manual"
	RequireApproval  bool          `json:"require_approval"`
}

// Participant represents a user in a collaborative session
type Participant struct {
	UserID      string                 `json:"user_id"`
	SessionID   string                 `json:"session_id"`
	Username    string                 `json:"username"`
	Role        ParticipantRole        `json:"role"`
	JoinedAt    time.Time              `json:"joined_at"`
	LastSeen    time.Time              `json:"last_seen"`
	IsActive    bool                   `json:"is_active"`
	Cursor      CursorPosition         `json:"cursor"`
	Permissions ParticipantPermissions `json:"permissions"`
}

// ParticipantRole defines user roles in a session
type ParticipantRole string

const (
	RoleOwner  ParticipantRole = "owner"
	RoleEditor ParticipantRole = "editor"
	RoleViewer ParticipantRole = "viewer"
)

// ParticipantPermissions defines what a participant can do
type ParticipantPermissions struct {
	CanEdit   bool `json:"can_edit"`
	CanInvite bool `json:"can_invite"`
	CanKick   bool `json:"can_kick"`
	CanDelete bool `json:"can_delete"`
}

// UserPresence tracks user activity across sessions
type UserPresence struct {
	UserID         string            `json:"user_id"`
	Username       string            `json:"username"`
	Status         PresenceStatus    `json:"status"`
	LastSeen       time.Time         `json:"last_seen"`
	CurrentSession string            `json:"current_session,omitempty"`
	DeviceInfo     map[string]string `json:"device_info"`
}

// PresenceStatus represents user availability
type PresenceStatus string

const (
	StatusOnline  PresenceStatus = "online"
	StatusAway    PresenceStatus = "away"
	StatusBusy    PresenceStatus = "busy"
	StatusOffline PresenceStatus = "offline"
)

// CursorPosition tracks where a user is editing
type CursorPosition struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Operation represents a document change for operational transform
type Operation struct {
	ID           string        `json:"id"`
	SessionID    string        `json:"session_id"`
	UserID       string        `json:"user_id"`
	Type         OperationType `json:"type"`
	Position     int           `json:"position"`
	Content      string        `json:"content,omitempty"`
	Length       int           `json:"length,omitempty"`
	Timestamp    time.Time     `json:"timestamp"`
	Version      int           `json:"version"`
	Dependencies []string      `json:"dependencies,omitempty"`
}

// OperationType defines types of document operations
type OperationType string

const (
	OpInsert OperationType = "insert"
	OpDelete OperationType = "delete"
	OpRetain OperationType = "retain"
)

// Collaboration Events
type CollaborationEvent struct {
	Type      EventType              `json:"type"`
	SessionID string                 `json:"session_id"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
	Timestamp time.Time              `json:"timestamp"`
}

// EventType defines different collaboration events
type EventType string

const (
	EventSessionCreated    EventType = "session_created"
	EventSessionJoined     EventType = "session_joined"
	EventSessionLeft       EventType = "session_left"
	EventParticipantJoined EventType = "participant_joined"
	EventParticipantLeft   EventType = "participant_left"
	EventOperationApplied  EventType = "operation_applied"
	EventConflictDetected  EventType = "conflict_detected"
	EventPresenceUpdated   EventType = "presence_updated"
)

// Collaboration Messages for broadcasting
type CollaborationMessage struct {
	Type      MessageType            `json:"type"`
	SessionID string                 `json:"session_id"`
	UserID    string                 `json:"user_id"`
	Data      map[string]interface{} `json:"data"`
}

// MessageType defines message types for broadcasting
type MessageType string

const (
	MsgOperation      MessageType = "operation"
	MsgCursorUpdate   MessageType = "cursor_update"
	MsgPresenceUpdate MessageType = "presence_update"
	MsgSessionUpdate  MessageType = "session_update"
)

// UI Integration Events
type SessionUpdateEvent struct {
	Session      *Session       `json:"session"`
	Action       string         `json:"action"` // "created", "joined", "left"
	Participants []*Participant `json:"participants"`
}

type PresenceUpdateEvent struct {
	UserID   string        `json:"user_id"`
	Presence *UserPresence `json:"presence"`
	Action   string        `json:"action"` // "joined", "left", "updated"
}

type ConflictEvent struct {
	SessionID   string      `json:"session_id"`
	Operations  []Operation `json:"operations"`
	Description string      `json:"description"`
}

// SetUICallbacks configures callbacks for UI integration
func (cm *CollaborationManager) SetUICallbacks(
	onSessionUpdate func(SessionUpdateEvent),
	onPresenceUpdate func(PresenceUpdateEvent),
	onConflictDetected func(ConflictEvent),
) {
	cm.onSessionUpdate = onSessionUpdate
	cm.onPresenceUpdate = onPresenceUpdate
	cm.onConflictDetected = onConflictDetected
}

// CreateSession creates a new collaborative session
func (cm *CollaborationManager) CreateSession(documentID int, name, createdBy string, settings SessionSettings) (*Session, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	sessionID := generateSessionID()
	session := &Session{
		ID:         sessionID,
		DocumentID: documentID,
		Name:       name,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		IsActive:   true,
		Settings:   settings,
		Operations: make([]Operation, 0),
	}

	// Store session
	cm.sessions[sessionID] = session
	cm.participants[sessionID] = make(map[string]*Participant)

	// Persist to database
	if err := cm.persistSession(session); err != nil {
		return nil, fmt.Errorf("failed to persist session: %w", err)
	}

	// Emit event
	cm.emitEvent(CollaborationEvent{
		Type:      EventSessionCreated,
		SessionID: sessionID,
		UserID:    createdBy,
		Timestamp: time.Now(),
	})

	return session, nil
}

// JoinSession allows a user to join an existing session
func (cm *CollaborationManager) JoinSession(sessionID, userID, username string, role ParticipantRole) (*Session, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	if !session.IsActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Check if user is already in session
	if _, participantExists := cm.participants[sessionID][userID]; participantExists {
		return nil, fmt.Errorf("user already in session: %s", userID)
	}

	// Check participant limit
	if len(cm.participants[sessionID]) >= session.Settings.MaxParticipants {
		return nil, fmt.Errorf("session is full")
	}

	// Create participant
	participant := &Participant{
		UserID:      userID,
		SessionID:   sessionID,
		Username:    username,
		Role:        role,
		JoinedAt:    time.Now(),
		LastSeen:    time.Now(),
		IsActive:    true,
		Permissions: cm.getPermissionsForRole(role),
	}

	cm.participants[sessionID][userID] = participant

	// Update user presence
	cm.updateUserPresence(userID, username, StatusOnline, sessionID)

	// Persist participant
	if err := cm.persistParticipant(participant); err != nil {
		return nil, fmt.Errorf("failed to persist participant: %w", err)
	}

	// Emit events
	cm.emitEvent(CollaborationEvent{
		Type:      EventSessionJoined,
		SessionID: sessionID,
		UserID:    userID,
		Timestamp: time.Now(),
	})

	cm.emitEvent(CollaborationEvent{
		Type:      EventParticipantJoined,
		SessionID: sessionID,
		UserID:    userID,
		Timestamp: time.Now(),
	})

	return session, nil
}

// LeaveSession removes a user from a session
func (cm *CollaborationManager) LeaveSession(sessionID, userID string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	participant, participantExists := cm.participants[sessionID][userID]
	if !participantExists {
		return fmt.Errorf("user not in session")
	}

	// Remove participant
	delete(cm.participants[sessionID], userID)

	// Update user presence
	cm.updateUserPresence(userID, participant.Username, StatusOffline, "")

	// If no participants left, end the session
	if len(cm.participants[sessionID]) == 0 {
		session.IsActive = false
		if err := cm.persistSession(session); err != nil {
			return fmt.Errorf("failed to update session: %w", err)
		}
	}

	// Remove participant from database
	if err := cm.removeParticipant(sessionID, userID); err != nil {
		return fmt.Errorf("failed to remove participant: %w", err)
	}

	// Emit events
	cm.emitEvent(CollaborationEvent{
		Type:      EventSessionLeft,
		SessionID: sessionID,
		UserID:    userID,
		Timestamp: time.Now(),
	})

	cm.emitEvent(CollaborationEvent{
		Type:      EventParticipantLeft,
		SessionID: sessionID,
		UserID:    userID,
		Timestamp: time.Now(),
	})

	return nil
}

// ApplyOperation applies a document operation to a session
func (cm *CollaborationManager) ApplyOperation(sessionID, userID string, operation Operation) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	participant, participantExists := cm.participants[sessionID][userID]
	if !participantExists {
		return fmt.Errorf("user not in session")
	}

	if !participant.Permissions.CanEdit {
		return fmt.Errorf("user does not have edit permissions")
	}

	// Set operation metadata
	operation.ID = generateOperationID()
	operation.SessionID = sessionID
	operation.UserID = userID
	operation.Timestamp = time.Now()
	operation.Version = len(session.Operations) + 1

	// Apply operational transform
	transformedOp, err := cm.transformOperation(session.Operations, operation)
	if err != nil {
		return fmt.Errorf("operation transform failed: %w", err)
	}

	// Add to session operations
	session.Operations = append(session.Operations, transformedOp)

	// Persist operation
	if err := cm.persistOperation(transformedOp); err != nil {
		return fmt.Errorf("failed to persist operation: %w", err)
	}

	// Update participant cursor/activity
	participant.LastSeen = time.Now()
	cm.participants[sessionID][userID] = participant

	// Emit events
	cm.emitEvent(CollaborationEvent{
		Type:      EventOperationApplied,
		SessionID: sessionID,
		UserID:    userID,
		Data: map[string]interface{}{
			"operation": transformedOp,
		},
		Timestamp: time.Now(),
	})

	// Broadcast operation to other participants
	cm.broadcastMessage(CollaborationMessage{
		Type:      MsgOperation,
		SessionID: sessionID,
		UserID:    userID,
		Data: map[string]interface{}{
			"operation": transformedOp,
		},
	})

	return nil
}

// UpdateCursor updates a user's cursor position
func (cm *CollaborationManager) UpdateCursor(sessionID, userID string, position CursorPosition) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Check if session exists
	_, sessionExists := cm.sessions[sessionID]
	if !sessionExists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Check if user is in session
	participant, participantExists := cm.participants[sessionID][userID]
	if !participantExists {
		return fmt.Errorf("user not in session")
	}

	participant.Cursor = position
	participant.LastSeen = time.Now()
	cm.participants[sessionID][userID] = participant

	// Broadcast cursor update
	cm.broadcastMessage(CollaborationMessage{
		Type:      MsgCursorUpdate,
		SessionID: sessionID,
		UserID:    userID,
		Data: map[string]interface{}{
			"cursor": position,
		},
	})

	return nil
}

// GetSession retrieves a session by ID
func (cm *CollaborationManager) GetSession(sessionID string) (*Session, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	session, exists := cm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session, nil
}

// GetParticipants returns all participants in a session
func (cm *CollaborationManager) GetParticipants(sessionID string) ([]*Participant, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	participants, exists := cm.participants[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	result := make([]*Participant, 0, len(participants))
	for _, p := range participants {
		result = append(result, p)
	}

	return result, nil
}

// GetUserPresence returns presence information for a user
func (cm *CollaborationManager) GetUserPresence(userID string) (*UserPresence, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	presence, exists := cm.presence[userID]
	if !exists {
		return nil, fmt.Errorf("user presence not found: %s", userID)
	}

	return presence, nil
}

// GetActiveSessions returns all active sessions
func (cm *CollaborationManager) GetActiveSessions() []*Session {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	sessions := make([]*Session, 0)
	for _, session := range cm.sessions {
		if session.IsActive {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// Close shuts down the collaboration manager
func (cm *CollaborationManager) Close() error {
	cm.cancel()

	// Mark all sessions as inactive
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, session := range cm.sessions {
		session.IsActive = false
		if err := cm.persistSession(session); err != nil {
			// Log error but continue cleanup
			continue
		}
	}

	return nil
}

// Helper methods

func (cm *CollaborationManager) emitEvent(event CollaborationEvent) {
	select {
	case cm.events <- event:
	default:
		// Channel full, drop event
	}
}

func (cm *CollaborationManager) broadcastMessage(msg CollaborationMessage) {
	select {
	case cm.broadcast <- msg:
	default:
		// Channel full, drop message
	}
}

func (cm *CollaborationManager) processEvents() {
	for {
		select {
		case <-cm.ctx.Done():
			return
		case event := <-cm.events:
			cm.handleEvent(event)
		case <-time.After(100 * time.Millisecond):
			// Timeout to prevent blocking
			continue
		}
	}
}

func (cm *CollaborationManager) processBroadcasts() {
	for {
		select {
		case <-cm.ctx.Done():
			return
		case msg := <-cm.broadcast:
			cm.handleBroadcast(msg)
		case <-time.After(100 * time.Millisecond):
			// Timeout to prevent blocking
			continue
		}
	}
}

func (cm *CollaborationManager) handleEvent(event CollaborationEvent) {
	switch event.Type {
	case EventSessionCreated:
		if cm.onSessionUpdate != nil {
			cm.onSessionUpdate(SessionUpdateEvent{
				Action: "created",
			})
		}
	case EventParticipantJoined:
		if cm.onPresenceUpdate != nil {
			cm.onPresenceUpdate(PresenceUpdateEvent{
				Action: "joined",
			})
		}
	case EventConflictDetected:
		if cm.onConflictDetected != nil {
			conflictEvent := ConflictEvent{
				Description: "Conflict detected in collaborative editing",
			}
			cm.onConflictDetected(conflictEvent)
		}
	}
}

func (cm *CollaborationManager) handleBroadcast(msg CollaborationMessage) {
	// Handle broadcasting to connected clients
	// This would integrate with WebSocket or similar real-time transport
}

func (cm *CollaborationManager) getPermissionsForRole(role ParticipantRole) ParticipantPermissions {
	switch role {
	case RoleOwner:
		return ParticipantPermissions{
			CanEdit:   true,
			CanInvite: true,
			CanKick:   true,
			CanDelete: true,
		}
	case RoleEditor:
		return ParticipantPermissions{
			CanEdit:   true,
			CanInvite: false,
			CanKick:   false,
			CanDelete: false,
		}
	case RoleViewer:
		return ParticipantPermissions{
			CanEdit:   false,
			CanInvite: false,
			CanKick:   false,
			CanDelete: false,
		}
	default:
		return ParticipantPermissions{}
	}
}

func (cm *CollaborationManager) updateUserPresence(userID, username string, status PresenceStatus, sessionID string) {
	presence := &UserPresence{
		UserID:         userID,
		Username:       username,
		Status:         status,
		LastSeen:       time.Now(),
		CurrentSession: sessionID,
		DeviceInfo:     make(map[string]string),
	}

	cm.presence[userID] = presence

	// Emit presence update event
	cm.emitEvent(CollaborationEvent{
		Type:      EventPresenceUpdated,
		UserID:    userID,
		Timestamp: time.Now(),
	})
}

// Database persistence methods
func (cm *CollaborationManager) persistSession(session *Session) error {
	// Implementation would use the database to persist session
	// For now, this is a placeholder
	return nil
}

func (cm *CollaborationManager) persistParticipant(participant *Participant) error {
	// Implementation would use the database to persist participant
	// For now, this is a placeholder
	return nil
}

func (cm *CollaborationManager) removeParticipant(sessionID, userID string) error {
	// Implementation would use the database to remove participant
	// For now, this is a placeholder
	return nil
}

func (cm *CollaborationManager) persistOperation(operation Operation) error {
	// Implementation would use the database to persist operation
	// For now, this is a placeholder
	return nil
}

// Operational Transform implementation
func (cm *CollaborationManager) transformOperation(operations []Operation, newOp Operation) (Operation, error) {
	// Simple operational transform implementation
	// In a full implementation, this would handle concurrent operations properly

	transformedOp := newOp

	// Check for conflicts with recent operations
	for _, op := range operations {
		if cm.operationsConflict(op, newOp) {
			// Handle conflict resolution
			cm.emitEvent(CollaborationEvent{
				Type:      EventConflictDetected,
				SessionID: newOp.SessionID,
				UserID:    newOp.UserID,
				Data: map[string]interface{}{
					"conflicting_operations": []Operation{op, newOp},
				},
				Timestamp: time.Now(),
			})
		}
	}

	return transformedOp, nil
}

func (cm *CollaborationManager) operationsConflict(op1, op2 Operation) bool {
	// Simple conflict detection based on position overlap
	return op1.Position <= op2.Position && op1.Position+op1.Length >= op2.Position
}

// Utility functions
func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

func generateOperationID() string {
	return fmt.Sprintf("op_%d", time.Now().UnixNano())
}
