package collaboration

import (
	"fmt"
	"sync"
	"time"
)

// SessionManager handles collaborative session lifecycle
type SessionManager struct {
	// Core state
	sessions     map[string]*Session
	participants map[string]map[string]*Participant
	mu           sync.RWMutex

	// Callbacks
	onSessionCreated func(*Session)
	onSessionEnded   func(*Session)
	onUserJoined     func(*Session, *Participant)
	onUserLeft       func(*Session, *Participant)
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions:     make(map[string]*Session),
		participants: make(map[string]map[string]*Participant),
	}

	return sm
}

// SetSessionCallbacks configures callbacks for session events
func (sm *SessionManager) SetSessionCallbacks(
	onSessionCreated func(*Session),
	onSessionEnded func(*Session),
	onUserJoined func(*Session, *Participant),
	onUserLeft func(*Session, *Participant),
) {
	sm.onSessionCreated = onSessionCreated
	sm.onSessionEnded = onSessionEnded
	sm.onUserJoined = onUserJoined
	sm.onUserLeft = onUserLeft
}

// CreateSession creates a new collaborative session
func (sm *SessionManager) CreateSession(documentID int, name, createdBy string, settings SessionSettings) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

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
	sm.sessions[sessionID] = session
	sm.participants[sessionID] = make(map[string]*Participant)

	// Emit session created event
	if sm.onSessionCreated != nil {
		sm.onSessionCreated(session)
	}

	return session, nil
}

// JoinSession allows a user to join an existing session
func (sm *SessionManager) JoinSession(sessionID, userID, username string, role ParticipantRole) (*Session, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	if !session.IsActive {
		return nil, fmt.Errorf("session is not active")
	}

	// Check participant limit
	if len(sm.participants[sessionID]) >= session.Settings.MaxParticipants {
		return nil, fmt.Errorf("session is full (max %d participants)", session.Settings.MaxParticipants)
	}

	// Check if user is already in session
	if _, alreadyInSession := sm.participants[sessionID][userID]; alreadyInSession {
		return nil, fmt.Errorf("user already in session")
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
		Permissions: sm.getPermissionsForRole(role),
	}

	sm.participants[sessionID][userID] = participant

	// Emit user joined event
	if sm.onUserJoined != nil {
		sm.onUserJoined(session, participant)
	}

	return session, nil
}

// LeaveSession removes a user from a session
func (sm *SessionManager) LeaveSession(sessionID, userID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	participant, participantExists := sm.participants[sessionID][userID]
	if !participantExists {
		return fmt.Errorf("user not in session")
	}

	// Remove participant
	delete(sm.participants[sessionID], userID)

	// If no participants left, end the session
	if len(sm.participants[sessionID]) == 0 {
		session.IsActive = false

		// Emit session ended event
		if sm.onSessionEnded != nil {
			sm.onSessionEnded(session)
		}
	} else {
		// Emit user left event
		if sm.onUserLeft != nil {
			sm.onUserLeft(session, participant)
		}
	}

	return nil
}

// EndSession ends a collaborative session
func (sm *SessionManager) EndSession(sessionID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	session.IsActive = false

	// Emit session ended event
	if sm.onSessionEnded != nil {
		sm.onSessionEnded(session)
	}

	return nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*Session, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session, nil
}

// GetParticipants returns all participants in a session
func (sm *SessionManager) GetParticipants(sessionID string) ([]*Participant, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	participants, exists := sm.participants[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	result := make([]*Participant, 0, len(participants))
	for _, p := range participants {
		result = append(result, p)
	}

	return result, nil
}

// GetParticipant returns a specific participant
func (sm *SessionManager) GetParticipant(sessionID, userID string) (*Participant, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	participants, exists := sm.participants[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	participant, exists := participants[userID]
	if !exists {
		return nil, fmt.Errorf("participant not found: %s", userID)
	}

	return participant, nil
}

// GetActiveSessions returns all active sessions
func (sm *SessionManager) GetActiveSessions() []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0)
	for _, session := range sm.sessions {
		if session.IsActive {
			sessions = append(sessions, session)
		}
	}

	return sessions
}

// GetSessionsForUser returns all sessions a user is participating in
func (sm *SessionManager) GetSessionsForUser(userID string) []*Session {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	sessions := make([]*Session, 0)
	for _, session := range sm.sessions {
		if participants, exists := sm.participants[session.ID]; exists {
			if _, userInSession := participants[userID]; userInSession {
				sessions = append(sessions, session)
			}
		}
	}

	return sessions
}

// UpdateParticipant updates participant information
func (sm *SessionManager) UpdateParticipant(sessionID, userID string, update func(*Participant)) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	participants, exists := sm.participants[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	participant, exists := participants[userID]
	if !exists {
		return fmt.Errorf("participant not found: %s", userID)
	}

	update(participant)
	sm.participants[sessionID][userID] = participant

	return nil
}

// IsUserInSession checks if a user is in a specific session
func (sm *SessionManager) IsUserInSession(sessionID, userID string) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	participants, exists := sm.participants[sessionID]
	if !exists {
		return false
	}

	_, exists = participants[userID]
	return exists
}

// GetSessionCount returns the total number of sessions
func (sm *SessionManager) GetSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.sessions)
}

// GetActiveSessionCount returns the number of active sessions
func (sm *SessionManager) GetActiveSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	count := 0
	for _, session := range sm.sessions {
		if session.IsActive {
			count++
		}
	}

	return count
}

// CleanupInactiveSessions removes sessions that have been inactive for too long
func (sm *SessionManager) CleanupInactiveSessions(maxAge time.Duration) []*Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	now := time.Now()
	var removedSessions []*Session

	for sessionID, session := range sm.sessions {
		// Check if session is old and has no active participants
		if !session.IsActive && now.Sub(session.CreatedAt) > maxAge {
			// Remove session and its participants
			delete(sm.sessions, sessionID)
			delete(sm.participants, sessionID)
			removedSessions = append(removedSessions, session)
		}
	}

	return removedSessions
}

// Helper methods

func (sm *SessionManager) getPermissionsForRole(role ParticipantRole) ParticipantPermissions {
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

// Session invitation system

// Invitation represents an invitation to join a session
type Invitation struct {
	ID        string          `json:"id"`
	SessionID string          `json:"session_id"`
	FromUser  string          `json:"from_user"`
	ToUser    string          `json:"to_user"`
	Role      ParticipantRole `json:"role"`
	Message   string          `json:"message"`
	CreatedAt time.Time       `json:"created_at"`
	ExpiresAt time.Time       `json:"expires_at"`
	Accepted  bool            `json:"accepted"`
}

// InvitationManager handles session invitations
type InvitationManager struct {
	invitations map[string]*Invitation
	mu          sync.RWMutex
}

// NewInvitationManager creates a new invitation manager
func NewInvitationManager() *InvitationManager {
	return &InvitationManager{
		invitations: make(map[string]*Invitation),
	}
}

// CreateInvitation creates a new session invitation
func (im *InvitationManager) CreateInvitation(sessionID, fromUser, toUser string, role ParticipantRole, message string, expiry time.Duration) *Invitation {
	im.mu.Lock()
	defer im.mu.Unlock()

	invitation := &Invitation{
		ID:        generateInvitationID(),
		SessionID: sessionID,
		FromUser:  fromUser,
		ToUser:    toUser,
		Role:      role,
		Message:   message,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(expiry),
		Accepted:  false,
	}

	im.invitations[invitation.ID] = invitation
	return invitation
}

// AcceptInvitation accepts an invitation
func (im *InvitationManager) AcceptInvitation(invitationID string) (*Invitation, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	invitation, exists := im.invitations[invitationID]
	if !exists {
		return nil, fmt.Errorf("invitation not found: %s", invitationID)
	}

	if time.Now().After(invitation.ExpiresAt) {
		return nil, fmt.Errorf("invitation expired")
	}

	if invitation.Accepted {
		return nil, fmt.Errorf("invitation already accepted")
	}

	invitation.Accepted = true
	return invitation, nil
}

// DeclineInvitation declines an invitation
func (im *InvitationManager) DeclineInvitation(invitationID string) error {
	im.mu.Lock()
	defer im.mu.Unlock()

	invitation, exists := im.invitations[invitationID]
	if !exists {
		return fmt.Errorf("invitation not found: %s", invitationID)
	}

	// Mark as declined (we could add a Declined field if needed)
	_ = invitation
	return nil
}

// GetInvitation retrieves an invitation by ID
func (im *InvitationManager) GetInvitation(invitationID string) (*Invitation, error) {
	im.mu.RLock()
	defer im.mu.RUnlock()

	invitation, exists := im.invitations[invitationID]
	if !exists {
		return nil, fmt.Errorf("invitation not found: %s", invitationID)
	}

	return invitation, nil
}

// GetInvitationsForUser returns all pending invitations for a user
func (im *InvitationManager) GetInvitationsForUser(userID string) []*Invitation {
	im.mu.RLock()
	defer im.mu.RUnlock()

	var userInvitations []*Invitation
	for _, invitation := range im.invitations {
		if invitation.ToUser == userID && !invitation.Accepted && time.Now().Before(invitation.ExpiresAt) {
			userInvitations = append(userInvitations, invitation)
		}
	}

	return userInvitations
}

// GetInvitationsBySession returns all invitations for a session
func (im *InvitationManager) GetInvitationsBySession(sessionID string) []*Invitation {
	im.mu.RLock()
	defer im.mu.RUnlock()

	var sessionInvitations []*Invitation
	for _, invitation := range im.invitations {
		if invitation.SessionID == sessionID && !invitation.Accepted && time.Now().Before(invitation.ExpiresAt) {
			sessionInvitations = append(sessionInvitations, invitation)
		}
	}

	return sessionInvitations
}

// CleanupExpiredInvitations removes expired invitations
func (im *InvitationManager) CleanupExpiredInvitations() []*Invitation {
	im.mu.Lock()
	defer im.mu.Unlock()

	now := time.Now()
	var expired []*Invitation

	for id, invitation := range im.invitations {
		if now.After(invitation.ExpiresAt) {
			expired = append(expired, invitation)
			delete(im.invitations, id)
		}
	}

	return expired
}

// Utility functions

func generateInvitationID() string {
	return fmt.Sprintf("invitation_%d", time.Now().UnixNano())
}
