package collaboration

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/infra/db"
)

// MockDB implements a mock database for testing
type MockDB struct {
	sessions     map[string]*Session
	participants map[string]map[string]*Participant
	operations   map[string][]Operation
	mu           sync.RWMutex
	shouldFail   bool
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		sessions:     make(map[string]*Session),
		participants: make(map[string]map[string]*Participant),
		operations:   make(map[string][]Operation),
		shouldFail:   false,
	}
}

// SetFailure sets whether database operations should fail
func (m *MockDB) SetFailure(shouldFail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = shouldFail
}

// MockCollaborationManager wraps CollaborationManager for testing
type MockCollaborationManager struct {
	*CollaborationManager
	events   []CollaborationEvent
	messages []CollaborationMessage
	mu       sync.RWMutex
	testDB   *MockDB
}

// NewMockCollaborationManager creates a new mock collaboration manager
func NewMockCollaborationManager(t *testing.T) *MockCollaborationManager {
	testDB := NewMockDB()

	// Create a mock database that implements the db.DB interface
	mockDB := &db.DB{} // We'll use this as a placeholder

	cm := NewCollaborationManager(mockDB)

	mcm := &MockCollaborationManager{
		CollaborationManager: cm,
		events:               make([]CollaborationEvent, 0),
		messages:             make([]CollaborationMessage, 0),
		testDB:               testDB,
	}

	// Set up event capture
	cm.SetUICallbacks(
		func(event SessionUpdateEvent) {
			if event.Session != nil {
				mcm.captureEvent(CollaborationEvent{
					Type:      EventSessionCreated,
					SessionID: event.Session.ID,
					Timestamp: time.Now(),
				})
			}
		},
		func(event PresenceUpdateEvent) {
			if event.Presence != nil {
				mcm.captureEvent(CollaborationEvent{
					Type:      EventPresenceUpdated,
					UserID:    event.UserID,
					Timestamp: time.Now(),
				})
			}
		},
		func(event ConflictEvent) {
			mcm.captureEvent(CollaborationEvent{
				Type:      EventConflictDetected,
				SessionID: event.SessionID,
				Timestamp: time.Now(),
			})
		},
	)

	return mcm
}

// captureEvent captures events for testing
func (mcm *MockCollaborationManager) captureEvent(event CollaborationEvent) {
	mcm.mu.Lock()
	defer mcm.mu.Unlock()
	mcm.events = append(mcm.events, event)
}

// GetEvents returns all captured events
func (mcm *MockCollaborationManager) GetEvents() []CollaborationEvent {
	mcm.mu.RLock()
	defer mcm.mu.RUnlock()

	events := make([]CollaborationEvent, len(mcm.events))
	copy(events, mcm.events)
	return events
}

// ClearEvents clears all captured events
func (mcm *MockCollaborationManager) ClearEvents() {
	mcm.mu.Lock()
	defer mcm.mu.Unlock()
	mcm.events = make([]CollaborationEvent, 0)
}

// MockSessionManager wraps SessionManager for testing
type MockSessionManager struct {
	*SessionManager
	sessionEvents []SessionEvent
	mu            sync.RWMutex
}

// SessionEvent represents session events for testing
type SessionEvent struct {
	Type      string
	SessionID string
	UserID    string
	Timestamp time.Time
}

// NewMockSessionManager creates a new mock session manager
func NewMockSessionManager() *MockSessionManager {
	sm := NewSessionManager()

	msm := &MockSessionManager{
		SessionManager: sm,
		sessionEvents:  make([]SessionEvent, 0),
	}

	// Set up event capture
	sm.SetSessionCallbacks(
		func(session *Session) {
			msm.captureEvent(SessionEvent{
				Type:      "session_created",
				SessionID: session.ID,
				Timestamp: time.Now(),
			})
		},
		func(session *Session) {
			msm.captureEvent(SessionEvent{
				Type:      "session_ended",
				SessionID: session.ID,
				Timestamp: time.Now(),
			})
		},
		func(session *Session, participant *Participant) {
			msm.captureEvent(SessionEvent{
				Type:      "user_joined",
				SessionID: session.ID,
				UserID:    participant.UserID,
				Timestamp: time.Now(),
			})
		},
		func(session *Session, participant *Participant) {
			msm.captureEvent(SessionEvent{
				Type:      "user_left",
				SessionID: session.ID,
				UserID:    participant.UserID,
				Timestamp: time.Now(),
			})
		},
	)

	return msm
}

// captureEvent captures session events for testing
func (msm *MockSessionManager) captureEvent(event SessionEvent) {
	msm.mu.Lock()
	defer msm.mu.Unlock()
	msm.sessionEvents = append(msm.sessionEvents, event)
}

// GetSessionEvents returns all captured session events
func (msm *MockSessionManager) GetSessionEvents() []SessionEvent {
	msm.mu.RLock()
	defer msm.mu.RUnlock()

	events := make([]SessionEvent, len(msm.sessionEvents))
	copy(events, msm.sessionEvents)
	return events
}

// ClearSessionEvents clears all captured session events
func (msm *MockSessionManager) ClearSessionEvents() {
	msm.mu.Lock()
	defer msm.mu.Unlock()
	msm.sessionEvents = make([]SessionEvent, 0)
}

// MockPresenceManager wraps PresenceManager for testing
type MockPresenceManager struct {
	*PresenceManager
	presenceEvents []PresenceEvent
	mu             sync.RWMutex
}

// NewMockPresenceManager creates a new mock presence manager
func NewMockPresenceManager() *MockPresenceManager {
	pm := NewPresenceManager()

	mpm := &MockPresenceManager{
		PresenceManager: pm,
		presenceEvents:  make([]PresenceEvent, 0),
	}

	// Set up event capture
	pm.SetPresenceCallbacks(
		func(event PresenceUpdateEvent) {
			mpm.captureEvent(PresenceEvent{
				Type:      EventStatusChanged,
				UserID:    event.UserID,
				Timestamp: time.Now(),
			})
		},
		func(sessionID, userID string) {
			mpm.captureEvent(PresenceEvent{
				Type:      EventUserJoined,
				SessionID: sessionID,
				UserID:    userID,
				Timestamp: time.Now(),
			})
		},
		func(sessionID, userID string) {
			mpm.captureEvent(PresenceEvent{
				Type:      EventUserLeft,
				SessionID: sessionID,
				UserID:    userID,
				Timestamp: time.Now(),
			})
		},
	)

	return mpm
}

// captureEvent captures presence events for testing
func (mpm *MockPresenceManager) captureEvent(event PresenceEvent) {
	mpm.mu.Lock()
	defer mpm.mu.Unlock()
	mpm.presenceEvents = append(mpm.presenceEvents, event)
}

// GetPresenceEvents returns all captured presence events
func (mpm *MockPresenceManager) GetPresenceEvents() []PresenceEvent {
	mpm.mu.RLock()
	defer mpm.mu.RUnlock()

	events := make([]PresenceEvent, len(mpm.presenceEvents))
	copy(events, mpm.presenceEvents)
	return events
}

// ClearPresenceEvents clears all captured presence events
func (mpm *MockPresenceManager) ClearPresenceEvents() {
	mpm.mu.Lock()
	defer mpm.mu.Unlock()
	mpm.presenceEvents = make([]PresenceEvent, 0)
}

// TestHelper provides utility functions for collaboration testing
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// CreateTestSession creates a test session
func (h *TestHelper) CreateTestSession(id, name, createdBy string) *Session {
	return &Session{
		ID:         id,
		DocumentID: 1,
		Name:       name,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now(),
		IsActive:   true,
		Settings: SessionSettings{
			MaxParticipants:  5,
			AutoSaveInterval: 30 * time.Second,
			ConflictStrategy: "merge",
			RequireApproval:  false,
		},
		Operations: make([]Operation, 0),
	}
}

// CreateTestParticipant creates a test participant
func (h *TestHelper) CreateTestParticipant(userID, sessionID, username string, role ParticipantRole) *Participant {
	return &Participant{
		UserID:      userID,
		SessionID:   sessionID,
		Username:    username,
		Role:        role,
		JoinedAt:    time.Now(),
		LastSeen:    time.Now(),
		IsActive:    true,
		Cursor:      CursorPosition{Line: 0, Column: 0},
		Permissions: getPermissionsForRole(role),
	}
}

// CreateTestOperation creates a test operation
func (h *TestHelper) CreateTestOperation(sessionID, userID string, opType OperationType, position int, content string) Operation {
	return Operation{
		ID:           generateOperationID(),
		SessionID:    sessionID,
		UserID:       userID,
		Type:         opType,
		Position:     position,
		Content:      content,
		Length:       len(content),
		Timestamp:    time.Now(),
		Version:      1,
		Dependencies: make([]string, 0),
	}
}

// CreateTestSong creates a test song
func (h *TestHelper) CreateTestSong(id int, title string) *domain.Song {
	return &domain.Song{
		ID:       id,
		Filepath: "/test/song.md",
		Metadata: domain.SongMetadata{
			Title: title,
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines: []domain.Line{
					{Text: "Test line 1", Syllables: 3},
					{Text: "Test line 2", Syllables: 3},
				},
			},
		},
	}
}

// AssertNoError asserts that an error is nil
func (h *TestHelper) AssertNoError(err error) {
	if err != nil {
		h.t.Fatalf("Expected no error, got: %v", err)
	}
}

// AssertError asserts that an error is not nil
func (h *TestHelper) AssertError(err error) {
	if err == nil {
		h.t.Fatal("Expected error, got nil")
	}
}

// AssertEqual asserts that two values are equal
func (h *TestHelper) AssertEqual(expected, actual interface{}) {
	if expected != actual {
		h.t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func (h *TestHelper) AssertNotEqual(expected, actual interface{}) {
	if expected == actual {
		h.t.Fatalf("Expected %v to not equal %v", expected, actual)
	}
}

// AssertTrue asserts that a condition is true
func (h *TestHelper) AssertTrue(condition bool, msgAndArgs ...interface{}) {
	if !condition {
		if len(msgAndArgs) == 0 {
			h.t.Fatal("Expected true, got false")
		} else if len(msgAndArgs) == 1 {
			h.t.Fatalf("Expected true, got false: %s", msgAndArgs[0])
		} else {
			format := msgAndArgs[0].(string)
			args := msgAndArgs[1:]
			h.t.Fatalf("Expected true, got false: "+format, args...)
		}
	}
}

// AssertFalse asserts that a condition is false
func (h *TestHelper) AssertFalse(condition bool, msg string) {
	if condition {
		h.t.Fatalf("Expected false, got true: %s", msg)
	}
}

// AssertNotNil asserts that a value is not nil
func (h *TestHelper) AssertNotNil(value interface{}) {
	if value == nil {
		h.t.Fatal("Expected non-nil value")
	}
}

// AssertLength asserts that a slice has the expected length
func (h *TestHelper) AssertLength(slice interface{}, expected int) {
	switch v := slice.(type) {
	case []Session:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []Participant:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []Operation:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []CollaborationEvent:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []PresenceEvent:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []SessionEvent:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []Conflict:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case []DocumentChange:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	case [][]Operation:
		if len(v) != expected {
			h.t.Fatalf("Expected length %d, got %d", expected, len(v))
		}
	default:
		// Use reflection for other slice types
		sliceValue := reflect.ValueOf(slice)
		if sliceValue.Kind() == reflect.Slice {
			if sliceValue.Len() != expected {
				h.t.Fatalf("Expected length %d, got %d", expected, sliceValue.Len())
			}
		} else {
			h.t.Fatalf("Unsupported type for AssertLength: %T", slice)
		}
	}
}

// WaitForCondition waits for a condition to be true or times out
func (h *TestHelper) WaitForCondition(condition func() bool, timeout time.Duration, message string) {
	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	h.t.Fatalf("Condition not met within timeout: %s", message)
}

// TestSetup provides common setup for collaboration tests
type TestSetup struct {
	MockDB            *MockDB
	MockCollaboration *MockCollaborationManager
	MockSession       *MockSessionManager
	MockPresence      *MockPresenceManager
	TestHelper        *TestHelper
	Cleanup           func()
}

// NewTestSetup creates a new test setup
func NewTestSetup(t *testing.T) *TestSetup {
	mockDB := NewMockDB()
	mockCollaboration := NewMockCollaborationManager(t)
	mockSession := NewMockSessionManager()
	mockPresence := NewMockPresenceManager()
	testHelper := NewTestHelper(t)

	cleanup := func() {
		mockCollaboration.Close()
		mockPresence.Close()
	}

	return &TestSetup{
		MockDB:            mockDB,
		MockCollaboration: mockCollaboration,
		MockSession:       mockSession,
		MockPresence:      mockPresence,
		TestHelper:        testHelper,
		Cleanup:           cleanup,
	}
}

// Helper function to get permissions for a role (copied from manager.go)
func getPermissionsForRole(role ParticipantRole) ParticipantPermissions {
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
