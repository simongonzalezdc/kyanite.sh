package collaboration

import (
	"testing"
	"time"
)

func TestCollaborationManager_CreateSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	tests := []struct {
		name        string
		documentID  int
		sessionName string
		createdBy   string
		settings    SessionSettings
		expectError bool
	}{
		{
			name:        "Valid session creation",
			documentID:  1,
			sessionName: "Test Session",
			createdBy:   "user1",
			settings: SessionSettings{
				MaxParticipants:  5,
				AutoSaveInterval: 30 * time.Second,
				ConflictStrategy: "merge",
				RequireApproval:  false,
			},
			expectError: false,
		},
		{
			name:        "Session with single participant limit",
			documentID:  2,
			sessionName: "Solo Session",
			createdBy:   "user2",
			settings: SessionSettings{
				MaxParticipants:  1,
				AutoSaveInterval: 60 * time.Second,
				ConflictStrategy: "lock",
				RequireApproval:  true,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := cm.CreateSession(tt.documentID, tt.sessionName, tt.createdBy, tt.settings)

			if tt.expectError {
				helper.AssertError(err)
				return
			}

			helper.AssertNoError(err)
			helper.AssertNotNil(session)
			helper.AssertEqual(tt.sessionName, session.Name)
			helper.AssertEqual(tt.createdBy, session.CreatedBy)
			helper.AssertEqual(tt.documentID, session.DocumentID)
			helper.AssertTrue(session.IsActive)

			// Verify session can be retrieved
			retrievedSession, err := cm.GetSession(session.ID)
			helper.AssertNoError(err)
			helper.AssertEqual(session.ID, retrievedSession.ID)
		})
	}
}

func TestCollaborationManager_JoinSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a test session
	session, err := cm.CreateSession(1, "Test Session", "owner", SessionSettings{
		MaxParticipants: 3,
	})
	helper.AssertNoError(err)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		username    string
		role        ParticipantRole
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid join as editor",
			sessionID:   session.ID,
			userID:      "user1",
			username:    "User One",
			role:        RoleEditor,
			expectError: false,
		},
		{
			name:        "Valid join as viewer",
			sessionID:   session.ID,
			userID:      "user2",
			username:    "User Two",
			role:        RoleViewer,
			expectError: false,
		},
		{
			name:        "Invalid session ID",
			sessionID:   "invalid-session",
			userID:      "user3",
			username:    "User Three",
			role:        RoleEditor,
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "Duplicate user join",
			sessionID:   session.ID,
			userID:      "user1",
			username:    "User One",
			role:        RoleEditor,
			expectError: true,
			errorMsg:    "user already in session",
		},
		{
			name:        "Session full",
			sessionID:   session.ID,
			userID:      "user4",
			username:    "User Four",
			role:        RoleEditor,
			expectError: true,
			errorMsg:    "session is full",
		},
	}

	// Add first two users (session max is 3, owner is already in)
	cm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	cm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinedSession, err := cm.JoinSession(tt.sessionID, tt.userID, tt.username, tt.role)

			if tt.expectError {
				helper.AssertError(err)
				if tt.errorMsg != "" {
					helper.AssertTrue(
						err.Error() == tt.errorMsg,
						"Expected error message '%s', got '%s'", tt.errorMsg, err.Error(),
					)
				}
				return
			}

			helper.AssertNoError(err)
			helper.AssertNotNil(joinedSession)
			helper.AssertEqual(session.ID, joinedSession.ID)

			// Verify participant was added
			participants, err := cm.GetParticipants(session.ID)
			helper.AssertNoError(err)

			found := false
			for _, p := range participants {
				if p.UserID == tt.userID {
					found = true
					helper.AssertEqual(tt.username, p.Username)
					helper.AssertEqual(tt.role, p.Role)
					break
				}
			}
			helper.AssertTrue(found, "Participant not found in session")
		})
	}
}

func TestCollaborationManager_LeaveSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a test session and add participants
	session, err := cm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	cm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	cm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid leave",
			sessionID:   session.ID,
			userID:      "user1",
			expectError: false,
		},
		{
			name:        "Invalid session ID",
			sessionID:   "invalid-session",
			userID:      "user2",
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "User not in session",
			sessionID:   session.ID,
			userID:      "nonexistent-user",
			expectError: true,
			errorMsg:    "user not in session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cm.LeaveSession(tt.sessionID, tt.userID)

			if tt.expectError {
				helper.AssertError(err)
				if tt.errorMsg != "" {
					helper.AssertTrue(
						err.Error() == tt.errorMsg,
						"Expected error message '%s', got '%s'", tt.errorMsg, err.Error(),
					)
				}
				return
			}

			helper.AssertNoError(err)

			// Verify participant was removed
			participants, err := cm.GetParticipants(session.ID)
			helper.AssertNoError(err)

			found := false
			for _, p := range participants {
				if p.UserID == tt.userID {
					found = true
					break
				}
			}
			helper.AssertFalse(found, "Participant should have been removed from session")
		})
	}
}

func TestCollaborationManager_ApplyOperation(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a test session and add a participant with edit permissions
	session, err := cm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	cm.JoinSession(session.ID, "editor", "Editor User", RoleEditor)
	cm.JoinSession(session.ID, "viewer", "Viewer User", RoleViewer)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		operation   Operation
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid insert operation by editor",
			sessionID:   session.ID,
			userID:      "editor",
			operation:   helper.CreateTestOperation(session.ID, "editor", OpInsert, 0, "Hello World"),
			expectError: false,
		},
		{
			name:        "Valid delete operation by editor",
			sessionID:   session.ID,
			userID:      "editor",
			operation:   helper.CreateTestOperation(session.ID, "editor", OpDelete, 5, "World"),
			expectError: false,
		},
		{
			name:        "Operation by viewer without permission",
			sessionID:   session.ID,
			userID:      "viewer",
			operation:   helper.CreateTestOperation(session.ID, "viewer", OpInsert, 0, "No Permission"),
			expectError: true,
			errorMsg:    "user does not have edit permissions",
		},
		{
			name:        "Operation in non-existent session",
			sessionID:   "invalid-session",
			userID:      "editor",
			operation:   helper.CreateTestOperation("invalid-session", "editor", OpInsert, 0, "Test"),
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "Operation by non-participant",
			sessionID:   session.ID,
			userID:      "nonparticipant",
			operation:   helper.CreateTestOperation(session.ID, "nonparticipant", OpInsert, 0, "Test"),
			expectError: true,
			errorMsg:    "user not in session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cm.ApplyOperation(tt.sessionID, tt.userID, tt.operation)

			if tt.expectError {
				helper.AssertError(err)
				if tt.errorMsg != "" {
					helper.AssertTrue(
						err.Error() == tt.errorMsg,
						"Expected error message '%s', got '%s'", tt.errorMsg, err.Error(),
					)
				}
				return
			}

			helper.AssertNoError(err)

			// Verify operation was added to session
			retrievedSession, err := cm.GetSession(session.ID)
			helper.AssertNoError(err)
			helper.AssertTrue(len(retrievedSession.Operations) > 0)

			// Check if our operation is in the list
			found := false
			for _, op := range retrievedSession.Operations {
				if op.UserID == tt.userID && op.Type == tt.operation.Type {
					found = true
					break
				}
			}
			helper.AssertTrue(found, "Operation not found in session")
		})
	}
}

func TestCollaborationManager_UpdateCursor(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a test session and add a participant
	session, err := cm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	cm.JoinSession(session.ID, "user1", "User One", RoleEditor)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		position    CursorPosition
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid cursor update",
			sessionID:   session.ID,
			userID:      "user1",
			position:    CursorPosition{Line: 5, Column: 10},
			expectError: false,
		},
		{
			name:        "Cursor update for non-existent session",
			sessionID:   "invalid-session",
			userID:      "user1",
			position:    CursorPosition{Line: 1, Column: 1},
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "Cursor update for non-participant",
			sessionID:   session.ID,
			userID:      "nonparticipant",
			position:    CursorPosition{Line: 1, Column: 1},
			expectError: true,
			errorMsg:    "user not in session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cm.UpdateCursor(tt.sessionID, tt.userID, tt.position)

			if tt.expectError {
				helper.AssertError(err)
				if tt.errorMsg != "" {
					helper.AssertTrue(
						err.Error() == tt.errorMsg,
						"Expected error message '%s', got '%s'", tt.errorMsg, err.Error(),
					)
				}
				return
			}

			helper.AssertNoError(err)

			// Verify cursor was updated
			participants, err := cm.GetParticipants(session.ID)
			helper.AssertNoError(err)

			found := false
			for _, p := range participants {
				if p.UserID == tt.userID {
					found = true
					helper.AssertEqual(tt.position.Line, p.Cursor.Line)
					helper.AssertEqual(tt.position.Column, p.Cursor.Column)
					break
				}
			}
			helper.AssertTrue(found, "Participant not found in session")
		})
	}
}

func TestCollaborationManager_ConcurrentOperations(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a test session and add multiple participants
	session, err := cm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	cm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	cm.JoinSession(session.ID, "user2", "User Two", RoleEditor)
	cm.JoinSession(session.ID, "user3", "User Three", RoleEditor)

	// Simulate concurrent operations
	done := make(chan bool, 3)

	// User 1 applies operations
	go func() {
		for i := 0; i < 5; i++ {
			op := helper.CreateTestOperation(session.ID, "user1", OpInsert, i*10, "User1Op"+string(rune('A'+i)))
			cm.ApplyOperation(session.ID, "user1", op)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// User 2 applies operations
	go func() {
		for i := 0; i < 5; i++ {
			op := helper.CreateTestOperation(session.ID, "user2", OpInsert, i*10+5, "User2Op"+string(rune('A'+i)))
			cm.ApplyOperation(session.ID, "user2", op)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// User 3 applies operations
	go func() {
		for i := 0; i < 5; i++ {
			op := helper.CreateTestOperation(session.ID, "user3", OpDelete, i*2, "del")
			cm.ApplyOperation(session.ID, "user3", op)
			time.Sleep(1 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for all operations to complete
	for i := 0; i < 3; i++ {
		<-done
	}

	// Verify all operations were applied
	retrievedSession, err := cm.GetSession(session.ID)
	helper.AssertNoError(err)
	helper.AssertTrue(len(retrievedSession.Operations) >= 10) // At least some operations should be applied

	// Verify events were generated
	events := cm.GetEvents()
	helper.AssertTrue(len(events) > 0)
}

func TestCollaborationManager_SessionLifecycle(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create a session
	session, err := cm.CreateSession(1, "Lifecycle Test", "owner", SessionSettings{})
	helper.AssertNoError(err)
	helper.AssertTrue(session.IsActive)

	// Add participants
	cm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	cm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	// Verify participants
	participants, err := cm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 2)

	// Remove one participant
	err = cm.LeaveSession(session.ID, "user1")
	helper.AssertNoError(err)

	participants, err = cm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 1)

	// Remove last participant (should end the session)
	err = cm.LeaveSession(session.ID, "user2")
	helper.AssertNoError(err)

	// Session should now be inactive
	retrievedSession, err := cm.GetSession(session.ID)
	helper.AssertNoError(err)
	helper.AssertFalse(retrievedSession.IsActive, "Session should be inactive after all participants leave")

	// Verify events were generated
	events := cm.GetEvents()
	helper.AssertTrue(len(events) > 0)
}

func TestCollaborationManager_GetActiveSessions(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Initially no active sessions
	activeSessions := cm.GetActiveSessions()
	helper.AssertLength(activeSessions, 0)

	// Create multiple sessions
	session1, _ := cm.CreateSession(1, "Session 1", "owner1", SessionSettings{})
	session2, _ := cm.CreateSession(2, "Session 2", "owner2", SessionSettings{})
	session3, _ := cm.CreateSession(3, "Session 3", "owner3", SessionSettings{})

	// Should have 3 active sessions
	activeSessions = cm.GetActiveSessions()
	helper.AssertLength(activeSessions, 3)

	// End one session by removing all participants
	cm.JoinSession(session3.ID, "user1", "User One", RoleEditor)
	cm.LeaveSession(session3.ID, "user1")

	// Should have 2 active sessions
	activeSessions = cm.GetActiveSessions()
	helper.AssertLength(activeSessions, 2)

	// Verify the right sessions are active
	sessionIDs := make(map[string]bool)
	for _, s := range activeSessions {
		sessionIDs[s.ID] = true
	}

	helper.AssertTrue(sessionIDs[session1.ID])
	helper.AssertTrue(sessionIDs[session2.ID])
	helper.AssertFalse(sessionIDs[session3.ID], "Session 3 should not be in active sessions")
}

func TestCollaborationManager_ErrorHandling(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Test operations on non-existent session
	_, err := cm.GetSession("non-existent")
	helper.AssertError(err)

	_, err = cm.GetParticipants("non-existent")
	helper.AssertError(err)

	err = cm.LeaveSession("non-existent", "user1")
	helper.AssertError(err)

	err = cm.UpdateCursor("non-existent", "user1", CursorPosition{Line: 1, Column: 1})
	helper.AssertError(err)

	// Test operations with invalid user
	session, _ := cm.CreateSession(1, "Test", "owner", SessionSettings{})

	err = cm.LeaveSession(session.ID, "non-existent-user")
	helper.AssertError(err)

	err = cm.UpdateCursor(session.ID, "non-existent-user", CursorPosition{Line: 1, Column: 1})
	helper.AssertError(err)

	op := helper.CreateTestOperation(session.ID, "non-existent-user", OpInsert, 0, "test")
	err = cm.ApplyOperation(session.ID, "non-existent-user", op)
	helper.AssertError(err)
}

func TestCollaborationManager_Close(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	cm := setup.MockCollaboration

	// Create some sessions
	session1, _ := cm.CreateSession(1, "Session 1", "owner1", SessionSettings{})
	session2, _ := cm.CreateSession(2, "Session 2", "owner2", SessionSettings{})

	// Verify sessions are active
	helper.AssertTrue(session1.IsActive)
	helper.AssertTrue(session2.IsActive)

	// Close the collaboration manager
	err := cm.Close()
	helper.AssertNoError(err)

	// Sessions should now be inactive
	retrievedSession1, _ := cm.GetSession(session1.ID)
	retrievedSession2, _ := cm.GetSession(session2.ID)

	helper.AssertFalse(retrievedSession1.IsActive, "Session 1 should be inactive after closing")
	helper.AssertFalse(retrievedSession2.IsActive, "Session 2 should be inactive after closing")
}
