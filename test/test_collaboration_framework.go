package noise

import (
	"fmt"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/infra/db"
)

// TestCollaborationFramework tests the complete collaboration framework
func TestCollaborationFramework(t *testing.T) {
	// Initialize in-memory database for testing
	database, err := db.New(db.Config{DataDir: ":memory:"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer database.Close()

	// Initialize collaboration components
	collaborationManager := collaboration.NewCollaborationManager(database)
	presenceManager := collaboration.NewPresenceManager()
	sessionManager := collaboration.NewSessionManager()
	invitationManager := collaboration.NewInvitationManager()
	conflictResolver := collaboration.NewConflictResolver()

	// Set up test callbacks
	var sessionEvents []string
	var presenceEvents []string
	var conflictEvents []string

	collaborationManager.SetUICallbacks(
		func(event collaboration.SessionUpdateEvent) {
			sessionEvents = append(sessionEvents, "session_update:"+event.Action)
		},
		func(event collaboration.PresenceUpdateEvent) {
			presenceEvents = append(presenceEvents, "presence_update:"+event.Action)
		},
		func(event collaboration.ConflictEvent) {
			conflictEvents = append(conflictEvents, "conflict:"+event.Description)
		},
	)

	sessionManager.SetSessionCallbacks(
		func(session *collaboration.Session) {
			sessionEvents = append(sessionEvents, "session_created")
		},
		func(session *collaboration.Session) {
			sessionEvents = append(sessionEvents, "session_ended")
		},
		func(session *collaboration.Session, participant *collaboration.Participant) {
			sessionEvents = append(sessionEvents, "user_joined:"+participant.UserID)
		},
		func(session *collaboration.Session, participant *collaboration.Participant) {
			sessionEvents = append(sessionEvents, "user_left:"+participant.UserID)
		},
	)

	// Test 1: Create a collaborative session
	t.Run("CreateSession", func(t *testing.T) {
		settings := collaboration.SessionSettings{
			MaxParticipants:  5,
			AutoSaveInterval: 30 * time.Second,
			ConflictStrategy: "merge",
			RequireApproval:  false,
		}

		session, err := sessionManager.CreateSession(1, "Test Song Session", "user1", settings)
		if err != nil {
			t.Fatalf("Failed to create session: %v", err)
		}

		if session.Name != "Test Song Session" {
			t.Errorf("Expected session name 'Test Song Session', got '%s'", session.Name)
		}

		if !session.IsActive {
			t.Error("Expected session to be active")
		}

		// Check that event was emitted
		if len(sessionEvents) == 0 || sessionEvents[len(sessionEvents)-1] != "session_created" {
			t.Error("Expected session_created event")
		}
	})

	// Test 2: User joins session
	t.Run("JoinSession", func(t *testing.T) {
		session, err := sessionManager.JoinSession("session_1", "user2", "testuser", collaboration.RoleEditor)
		if err != nil {
			t.Fatalf("Failed to join session: %v", err)
		}

		if session.ID != "session_1" {
			t.Errorf("Expected session ID 'session_1', got '%s'", session.ID)
		}

		// Check participants
		participants, err := sessionManager.GetParticipants("session_1")
		if err != nil {
			t.Fatalf("Failed to get participants: %v", err)
		}

		if len(participants) != 1 {
			t.Errorf("Expected 1 participant, got %d", len(participants))
		}

		participant := participants[0]
		if participant.Username != "testuser" {
			t.Errorf("Expected username 'testuser', got '%s'", participant.Username)
		}

		if participant.Role != collaboration.RoleEditor {
			t.Errorf("Expected role 'editor', got '%s'", participant.Role)
		}

		// Check that events were emitted
		expectedEvents := []string{"session_created", "user_joined:user2"}
		if len(sessionEvents) < len(expectedEvents) {
			t.Errorf("Expected at least %d events, got %d", len(expectedEvents), len(sessionEvents))
		}
	})

	// Test 3: User presence tracking
	t.Run("UserPresence", func(t *testing.T) {
		presenceManager.AddUserToSession("session_1", "user2", "testuser", collaboration.RoleEditor)

		presence, exists := presenceManager.GetUserPresence("user2")
		if !exists {
			t.Fatal("Expected user presence to exist")
		}

		if presence.Status != collaboration.StatusOnline {
			t.Errorf("Expected status 'online', got '%s'", presence.Status)
		}

		if presence.CurrentSession != "session_1" {
			t.Errorf("Expected current session 'session_1', got '%s'", presence.CurrentSession)
		}

		// Test presence indicators
		indicator := presenceManager.GetPresenceIndicator("user2")
		if indicator.Status != collaboration.StatusOnline {
			t.Errorf("Expected indicator status 'online', got '%s'", indicator.Status)
		}

		if indicator.Color != "green" {
			t.Errorf("Expected indicator color 'green', got '%s'", indicator.Color)
		}
	})

	// Test 4: Document operations and conflict resolution
	t.Run("DocumentOperations", func(t *testing.T) {
		// Create test operations that might conflict
		op1 := collaboration.Operation{
			ID:        "op1",
			SessionID: "session_1",
			UserID:    "user1",
			Type:      collaboration.OpInsert,
			Position:  10,
			Content:   "Hello ",
			Timestamp: time.Now(),
			Version:   1,
		}

		op2 := collaboration.Operation{
			ID:        "op2",
			SessionID: "session_1",
			UserID:    "user2",
			Type:      collaboration.OpInsert,
			Position:  10, // Same position - potential conflict
			Content:   "Hi ",
			Timestamp: time.Now().Add(time.Millisecond),
			Version:   2,
		}

		// Test conflict detection
		operations := []collaboration.Operation{op1, op2}
		conflicts := conflictResolver.DetectConflicts(operations)

		if len(conflicts) == 0 {
			t.Error("Expected conflicts to be detected")
		}

		// Test conflict resolution
		if len(conflicts) > 0 {
			conflict := conflicts[0]
			resolution, err := conflictResolver.ResolveConflict(conflict, "merge", "user1")
			if err != nil {
				t.Fatalf("Failed to resolve conflict: %v", err)
			}

			if resolution.Strategy != "merge" {
				t.Errorf("Expected strategy 'merge', got '%s'", resolution.Strategy)
			}

			if len(resolution.Changes) == 0 {
				t.Error("Expected resolution to contain changes")
			}
		}
	})

	// Test 5: Session invitations
	t.Run("SessionInvitations", func(t *testing.T) {
		invitation := invitationManager.CreateInvitation(
			"session_1",
			"user1",
			"user3",
			collaboration.RoleViewer,
			"Please join our collaborative session!",
			time.Hour,
		)

		if invitation.ToUser != "user3" {
			t.Errorf("Expected invitation to user 'user3', got '%s'", invitation.ToUser)
		}

		if invitation.Role != collaboration.RoleViewer {
			t.Errorf("Expected role 'viewer', got '%s'", invitation.Role)
		}

		// Test accepting invitation
		accepted, err := invitationManager.AcceptInvitation(invitation.ID)
		if err != nil {
			t.Fatalf("Failed to accept invitation: %v", err)
		}

		if !accepted.Accepted {
			t.Error("Expected invitation to be marked as accepted")
		}

		// Test getting invitations for user
		invitations := invitationManager.GetInvitationsForUser("user3")
		if len(invitations) == 0 {
			t.Error("Expected to find invitations for user3")
		}
	})

	// Test 6: Session management
	t.Run("SessionManagement", func(t *testing.T) {
		// Test getting active sessions
		activeSessions := sessionManager.GetActiveSessions()
		if len(activeSessions) == 0 {
			t.Error("Expected active sessions")
		}

		// Test getting sessions for user
		userSessions := sessionManager.GetSessionsForUser("user2")
		if len(userSessions) == 0 {
			t.Error("Expected sessions for user2")
		}

		// Test leaving session
		err := sessionManager.LeaveSession("session_1", "user2")
		if err != nil {
			t.Fatalf("Failed to leave session: %v", err)
		}

		// Verify user left
		participants, err := sessionManager.GetParticipants("session_1")
		if err != nil {
			t.Fatalf("Failed to get participants: %v", err)
		}

		if len(participants) != 0 {
			t.Errorf("Expected 0 participants after leaving, got %d", len(participants))
		}
	})

	// Test 7: Framework extensibility
	t.Run("Extensibility", func(t *testing.T) {
		// Test registering custom conflict resolution strategy
		customStrategy := &CustomConflictStrategy{}
		conflictResolver.RegisterStrategy(customStrategy)

		// Create a test conflict
		conflict := &collaboration.Conflict{
			ID:         "test_conflict",
			SessionID:  "session_1",
			Operations: []collaboration.Operation{},
		}

		// Test custom strategy
		resolution, err := conflictResolver.ResolveConflict(conflict, "custom", "user1")
		if err != nil {
			t.Fatalf("Failed to use custom strategy: %v", err)
		}

		if resolution.Strategy != "custom" {
			t.Errorf("Expected custom strategy, got '%s'", resolution.Strategy)
		}

		// Test presence manager extensibility
		presenceManager.SetUserDeviceInfo("user1", map[string]string{
			"device": "laptop",
			"os":     "linux",
		})

		presence, exists := presenceManager.GetUserPresence("user1")
		if !exists {
			t.Fatal("Expected user presence")
		}

		if presence.DeviceInfo["device"] != "laptop" {
			t.Error("Expected device info to be set")
		}
	})

	// Test 8: Cleanup
	t.Run("Cleanup", func(t *testing.T) {
		// Test cleanup of inactive sessions
		removedSessions := sessionManager.CleanupInactiveSessions(time.Hour)
		if len(removedSessions) != 0 {
			t.Errorf("Expected 0 removed sessions, got %d", len(removedSessions))
		}

		// Test cleanup of expired invitations
		expiredInvitations := invitationManager.CleanupExpiredInvitations()
		if len(expiredInvitations) != 0 {
			t.Errorf("Expected 0 expired invitations, got %d", len(expiredInvitations))
		}

		// Test manager cleanup
		err := collaborationManager.Close()
		if err != nil {
			t.Errorf("Failed to close collaboration manager: %v", err)
		}

		err = presenceManager.Close()
		if err != nil {
			t.Errorf("Failed to close presence manager: %v", err)
		}
	})

	// Summary
	fmt.Printf("Collaboration Framework Test Results:\n")
	fmt.Printf("- Session Events: %v\n", sessionEvents)
	fmt.Printf("- Presence Events: %v\n", presenceEvents)
	fmt.Printf("- Conflict Events: %v\n", conflictEvents)
	fmt.Printf("- Total Sessions Created: %d\n", sessionManager.GetSessionCount())
	fmt.Printf("- Active Sessions: %d\n", sessionManager.GetActiveSessionCount())

	t.Logf("Collaboration framework test completed successfully")
}

// CustomConflictStrategy demonstrates framework extensibility
type CustomConflictStrategy struct{}

func (ccs *CustomConflictStrategy) Name() string {
	return "custom"
}

func (ccs *CustomConflictStrategy) Resolve(conflict *collaboration.Conflict) (*collaboration.Resolution, error) {
	resolution := &collaboration.Resolution{
		Strategy:     "custom",
		FinalContent: "[CUSTOM_RESOLUTION]",
		Metadata: map[string]interface{}{
			"custom_resolution": true,
			"timestamp":         time.Now(),
		},
	}

	return resolution, nil
}
