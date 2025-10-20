package noise

import (
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/infra/db"
)

// TestCollaborationFramework exercises the collaboration managers through a set
// of focused sub-tests to keep complexity manageable while preserving behavior.
func TestCollaborationFramework(t *testing.T) {
	env := newCollabTestEnv(t)

	t.Run("CreateSession", env.testCreateSession)
	t.Run("JoinSession", env.testJoinSession)
	t.Run("UserPresence", env.testUserPresence)
	t.Run("DocumentOperations", env.testDocumentOperations)
	t.Run("SessionInvitations", env.testSessionInvitations)
	t.Run("SessionManagement", env.testSessionManagement)
	t.Run("Extensibility", env.testExtensibility)
	t.Run("Cleanup", env.testCleanup)
	t.Run("Summary", env.testSummary)
}

type collabTestEnv struct {
	collabManager     *collaboration.CollaborationManager
	presenceManager   *collaboration.PresenceManager
	sessionManager    *collaboration.SessionManager
	invitationManager *collaboration.InvitationManager
	conflictResolver  *collaboration.ConflictResolver
	db                *db.DB

	sessionID string

	sessionEvents  []string
	presenceEvents []string
	conflictEvents []string
}

func newCollabTestEnv(t *testing.T) *collabTestEnv {
	t.Helper()

	database, err := db.New(db.Config{DataDir: ":memory:"})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Fatalf("Failed to close test database: %v", err)
		}
	})

	env := &collabTestEnv{
		db:                database,
		collabManager:     collaboration.NewCollaborationManager(database),
		presenceManager:   collaboration.NewPresenceManager(),
		sessionManager:    collaboration.NewSessionManager(),
		invitationManager: collaboration.NewInvitationManager(),
		conflictResolver:  collaboration.NewConflictResolver(),
	}

	env.collabManager.SetUICallbacks(
		func(event collaboration.SessionUpdateEvent) {
			env.sessionEvents = append(env.sessionEvents, "session_update:"+event.Action)
		},
		func(event collaboration.PresenceUpdateEvent) {
			env.presenceEvents = append(env.presenceEvents, "presence_update:"+event.Action)
		},
		func(event collaboration.ConflictEvent) {
			env.conflictEvents = append(env.conflictEvents, "conflict:"+event.Description)
		},
	)

	env.sessionManager.SetSessionCallbacks(
		func(session *collaboration.Session) {
			env.sessionEvents = append(env.sessionEvents, "session_created")
		},
		func(session *collaboration.Session) {
			env.sessionEvents = append(env.sessionEvents, "session_ended")
		},
		func(session *collaboration.Session, participant *collaboration.Participant) {
			env.sessionEvents = append(env.sessionEvents, "user_joined:"+participant.UserID)
		},
		func(session *collaboration.Session, participant *collaboration.Participant) {
			env.sessionEvents = append(env.sessionEvents, "user_left:"+participant.UserID)
		},
	)

	return env
}

func (env *collabTestEnv) testCreateSession(t *testing.T) {
	t.Helper()

	settings := collaboration.SessionSettings{
		MaxParticipants:  5,
		AutoSaveInterval: 30 * time.Second,
		ConflictStrategy: "merge",
		RequireApproval:  false,
	}

	session, err := env.sessionManager.CreateSession(1, "Test Song Session", "user1", settings)
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	env.sessionID = session.ID
	if session.Name != "Test Song Session" {
		t.Fatalf("Expected session name 'Test Song Session', got %q", session.Name)
	}
	if !session.IsActive {
		t.Fatal("Expected session to be active")
	}
}

func (env *collabTestEnv) testJoinSession(t *testing.T) {
	t.Helper()

	session, err := env.sessionManager.JoinSession(env.sessionID, "user2", "testuser", collaboration.RoleEditor)
	if err != nil {
		t.Fatalf("Failed to join session: %v", err)
	}

	if session.ID != env.sessionID {
		t.Fatalf("Expected session ID %q, got %q", env.sessionID, session.ID)
	}

	participants, err := env.sessionManager.GetParticipants(env.sessionID)
	if err != nil {
		t.Fatalf("Failed to get participants: %v", err)
	}
	if len(participants) != 1 {
		t.Fatalf("Expected 1 participant, got %d", len(participants))
	}
	participant := participants[0]
	if participant.Username != "testuser" {
		t.Fatalf("Expected username 'testuser', got %q", participant.Username)
	}
	if participant.Role != collaboration.RoleEditor {
		t.Fatalf("Expected role 'editor', got %q", participant.Role)
	}
}

func (env *collabTestEnv) testUserPresence(t *testing.T) {
	t.Helper()

	env.presenceManager.AddUserToSession(env.sessionID, "user2", "testuser", collaboration.RoleEditor)

	presence, exists := env.presenceManager.GetUserPresence("user2")
	if !exists {
		t.Fatal("Expected user presence to exist")
	}
	if presence.Status != collaboration.StatusOnline {
		t.Fatalf("Expected status 'online', got %q", presence.Status)
	}
	if presence.CurrentSession != env.sessionID {
		t.Fatalf("Expected current session %q, got %q", env.sessionID, presence.CurrentSession)
	}

	indicator := env.presenceManager.GetPresenceIndicator("user2")
	if indicator.Status != collaboration.StatusOnline {
		t.Fatalf("Expected indicator status 'online', got %q", indicator.Status)
	}
	if indicator.Color != "green" {
		t.Fatalf("Expected indicator color 'green', got %q", indicator.Color)
	}
}

func (env *collabTestEnv) testDocumentOperations(t *testing.T) {
	t.Helper()

	op1 := collaboration.Operation{
		ID:        "op1",
		SessionID: env.sessionID,
		UserID:    "user1",
		Type:      collaboration.OpInsert,
		Position:  10,
		Content:   "Hello ",
		Timestamp: time.Now(),
		Version:   1,
	}

	op2 := collaboration.Operation{
		ID:        "op2",
		SessionID: env.sessionID,
		UserID:    "user2",
		Type:      collaboration.OpInsert,
		Position:  10,
		Content:   "Hi ",
		Timestamp: time.Now().Add(time.Millisecond),
		Version:   2,
	}

	conflicts := env.conflictResolver.DetectConflicts([]collaboration.Operation{op1, op2})
	if len(conflicts) == 0 {
		t.Fatal("Expected conflicts to be detected")
	}

	resolution, err := env.conflictResolver.ResolveConflict(conflicts[0], "merge", "user1")
	if err != nil {
		t.Fatalf("Failed to resolve conflict: %v", err)
	}
	if resolution.Strategy != "merge" {
		t.Fatalf("Expected strategy 'merge', got %q", resolution.Strategy)
	}
	if len(resolution.Changes) == 0 {
		t.Fatal("Expected resolution to contain changes")
	}
}

func (env *collabTestEnv) testSessionInvitations(t *testing.T) {
	t.Helper()

	invitation := env.invitationManager.CreateInvitation(
		env.sessionID,
		"user1",
		"user3",
		collaboration.RoleViewer,
		"Please join our collaborative session!",
		time.Hour,
	)

	if invitation.ToUser != "user3" {
		t.Fatalf("Expected invitation to user 'user3', got %q", invitation.ToUser)
	}
	if invitation.Role != collaboration.RoleViewer {
		t.Fatalf("Expected role 'viewer', got %q", invitation.Role)
	}

	accepted, err := env.invitationManager.AcceptInvitation(invitation.ID)
	if err != nil {
		t.Fatalf("Failed to accept invitation: %v", err)
	}
	if !accepted.Accepted {
		t.Fatal("Expected invitation to be marked as accepted")
	}

	invitations := env.invitationManager.GetInvitationsForUser("user3")
	if len(invitations) == 0 {
		t.Fatal("Expected to find invitations for user3")
	}
}

func (env *collabTestEnv) testSessionManagement(t *testing.T) {
	t.Helper()

	activeSessions := env.sessionManager.GetActiveSessions()
	if len(activeSessions) == 0 {
		t.Fatal("Expected active sessions")
	}

	userSessions := env.sessionManager.GetSessionsForUser("user2")
	if len(userSessions) == 0 {
		t.Fatal("Expected sessions for user2")
	}

	if err := env.sessionManager.LeaveSession(env.sessionID, "user2"); err != nil {
		t.Fatalf("Failed to leave session: %v", err)
	}

	participants, err := env.sessionManager.GetParticipants(env.sessionID)
	if err != nil {
		t.Fatalf("Failed to get participants: %v", err)
	}
	if len(participants) != 0 {
		t.Fatalf("Expected 0 participants after leaving, got %d", len(participants))
	}
}

func (env *collabTestEnv) testExtensibility(t *testing.T) {
	t.Helper()

	customStrategy := &CustomConflictStrategy{}
	env.conflictResolver.RegisterStrategy(customStrategy)

	conflict := &collaboration.Conflict{
		ID:         "test_conflict",
		SessionID:  env.sessionID,
		Operations: []collaboration.Operation{},
	}

	resolution, err := env.conflictResolver.ResolveConflict(conflict, "custom", "user1")
	if err != nil {
		t.Fatalf("Failed to use custom strategy: %v", err)
	}
	if resolution.Strategy != "custom" {
		t.Fatalf("Expected custom strategy, got %q", resolution.Strategy)
	}

	env.presenceManager.SetUserDeviceInfo("user1", map[string]string{
		"device": "laptop",
		"os":     "linux",
	})

	presence, exists := env.presenceManager.GetUserPresence("user1")
	if !exists {
		t.Fatal("Expected user presence")
	}
	if presence.DeviceInfo["device"] != "laptop" {
		t.Fatal("Expected device info to be set")
	}
}

func (env *collabTestEnv) testCleanup(t *testing.T) {
	t.Helper()

	removedSessions := env.sessionManager.CleanupInactiveSessions(time.Hour)
	if len(removedSessions) != 0 {
		t.Fatalf("Expected 0 removed sessions, got %d", len(removedSessions))
	}

	expiredInvitations := env.invitationManager.CleanupExpiredInvitations()
	if len(expiredInvitations) != 0 {
		t.Fatalf("Expected 0 expired invitations, got %d", len(expiredInvitations))
	}

	if err := env.collabManager.Close(); err != nil {
		t.Fatalf("Failed to close collaboration manager: %v", err)
	}
	if err := env.presenceManager.Close(); err != nil {
		t.Fatalf("Failed to close presence manager: %v", err)
	}
}

func (env *collabTestEnv) testSummary(t *testing.T) {
	t.Helper()

	t.Log("Collaboration Framework Test Results:")
	t.Logf("- Session Events: %v", env.sessionEvents)
	t.Logf("- Presence Events: %v", env.presenceEvents)
	t.Logf("- Conflict Events: %v", env.conflictEvents)
	t.Logf("- Total Sessions Created: %d", env.sessionManager.GetSessionCount())
	t.Logf("- Active Sessions: %d", env.sessionManager.GetActiveSessionCount())
}

// CustomConflictStrategy demonstrates framework extensibility by providing a
// no-op conflict resolution strategy used in tests.
type CustomConflictStrategy struct{}

// Name returns the identifier for the custom conflict strategy.
func (ccs *CustomConflictStrategy) Name() string {
	return "custom"
}

// Resolve applies a trivial resolution indicating the strategy was invoked.
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
