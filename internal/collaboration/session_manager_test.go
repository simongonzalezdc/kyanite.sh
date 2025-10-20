package collaboration

import (
	"testing"
	"time"
)

func TestSessionManager_CreateSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

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
			session, err := sm.CreateSession(tt.documentID, tt.sessionName, tt.createdBy, tt.settings)

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
			retrievedSession, err := sm.GetSession(session.ID)
			helper.AssertNoError(err)
			helper.AssertEqual(session.ID, retrievedSession.ID)
		})
	}
}

func TestSessionManager_JoinSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{
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

	// Add first two users (session max is 3)
	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			joinedSession, err := sm.JoinSession(tt.sessionID, tt.userID, tt.username, tt.role)

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
			participants, err := sm.GetParticipants(session.ID)
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

func TestSessionManager_LeaveSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session and add participants
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

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
			err := sm.LeaveSession(tt.sessionID, tt.userID)

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
			participants, err := sm.GetParticipants(session.ID)
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

func TestSessionManager_EndSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session and add participants
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	// Verify session is active
	helper.AssertTrue(session.IsActive)

	// End the session
	err = sm.EndSession(session.ID)
	helper.AssertNoError(err)

	// Verify session is no longer active
	retrievedSession, err := sm.GetSession(session.ID)
	helper.AssertNoError(err)
	helper.AssertFalse(retrievedSession.IsActive, "Session should be inactive after ending")

	// Test ending non-existent session
	err = sm.EndSession("non-existent")
	helper.AssertError(err)

	// Verify session events were generated
	events := sm.GetSessionEvents()
	helper.AssertTrue(len(events) > 0)
}

func TestSessionManager_GetParticipant(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session and add participants
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid participant retrieval",
			sessionID:   session.ID,
			userID:      "user1",
			expectError: false,
		},
		{
			name:        "Invalid session ID",
			sessionID:   "invalid-session",
			userID:      "user1",
			expectError: true,
			errorMsg:    "session not found",
		},
		{
			name:        "Non-existent participant",
			sessionID:   session.ID,
			userID:      "nonexistent-user",
			expectError: true,
			errorMsg:    "participant not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			participant, err := sm.GetParticipant(tt.sessionID, tt.userID)

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
			helper.AssertNotNil(participant)
			helper.AssertEqual(tt.userID, participant.UserID)
			helper.AssertEqual(tt.sessionID, participant.SessionID)
		})
	}
}

func TestSessionManager_GetSessionsForUser(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create multiple sessions
	session1, _ := sm.CreateSession(1, "Session 1", "owner1", SessionSettings{})
	session2, _ := sm.CreateSession(2, "Session 2", "owner2", SessionSettings{})
	session3, _ := sm.CreateSession(3, "Session 3", "owner3", SessionSettings{})

	// Add user1 to sessions 1 and 2
	sm.JoinSession(session1.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session2.ID, "user1", "User One", RoleViewer)

	// Add user2 to session 3
	sm.JoinSession(session3.ID, "user2", "User Two", RoleEditor)

	tests := []struct {
		name          string
		userID        string
		expectedCount int
		expectedIDs   []string
	}{
		{
			name:          "User in multiple sessions",
			userID:        "user1",
			expectedCount: 2,
			expectedIDs:   []string{session1.ID, session2.ID},
		},
		{
			name:          "User in single session",
			userID:        "user2",
			expectedCount: 1,
			expectedIDs:   []string{session3.ID},
		},
		{
			name:          "User in no sessions",
			userID:        "nonexistent-user",
			expectedCount: 0,
			expectedIDs:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions := sm.GetSessionsForUser(tt.userID)
			helper.AssertLength(sessions, tt.expectedCount)

			if len(tt.expectedIDs) > 0 {
				sessionIDs := make(map[string]bool)
				for _, s := range sessions {
					sessionIDs[s.ID] = true
				}

				for _, expectedID := range tt.expectedIDs {
					helper.AssertTrue(sessionIDs[expectedID], "Expected session ID %s not found", expectedID)
				}
			}
		})
	}
}

func TestSessionManager_UpdateParticipant(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session and add a participant
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)

	// Update participant
	err = sm.UpdateParticipant(session.ID, "user1", func(p *Participant) {
		p.Username = "Updated Name"
		p.Cursor = CursorPosition{Line: 5, Column: 10}
	})
	helper.AssertNoError(err)

	// Verify update
	participant, err := sm.GetParticipant(session.ID, "user1")
	helper.AssertNoError(err)
	helper.AssertEqual("Updated Name", participant.Username)
	helper.AssertEqual(5, participant.Cursor.Line)
	helper.AssertEqual(10, participant.Cursor.Column)

	// Test updating non-existent participant
	err = sm.UpdateParticipant(session.ID, "nonexistent-user", func(p *Participant) {
		p.Username = "Should not update"
	})
	helper.AssertError(err)

	// Test updating in non-existent session
	err = sm.UpdateParticipant("non-existent-session", "user1", func(p *Participant) {
		p.Username = "Should not update"
	})
	helper.AssertError(err)
}

func TestSessionManager_IsUserInSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a test session and add a participant
	session, err := sm.CreateSession(1, "Test Session", "owner", SessionSettings{})
	helper.AssertNoError(err)

	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)

	tests := []struct {
		name           string
		sessionID      string
		userID         string
		expectedResult bool
	}{
		{
			name:           "User in session",
			sessionID:      session.ID,
			userID:         "user1",
			expectedResult: true,
		},
		{
			name:           "User not in session",
			sessionID:      session.ID,
			userID:         "nonexistent-user",
			expectedResult: false,
		},
		{
			name:           "Non-existent session",
			sessionID:      "non-existent-session",
			userID:         "user1",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sm.IsUserInSession(tt.sessionID, tt.userID)
			helper.AssertEqual(tt.expectedResult, result)
		})
	}
}

func TestSessionManager_SessionCounts(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Initially no sessions
	helper.AssertEqual(0, sm.GetSessionCount())
	helper.AssertEqual(0, sm.GetActiveSessionCount())

	// Create sessions
	_, _ = sm.CreateSession(1, "Session 1", "owner1", SessionSettings{})
	_, _ = sm.CreateSession(2, "Session 2", "owner2", SessionSettings{})
	session3, _ := sm.CreateSession(3, "Session 3", "owner3", SessionSettings{})

	// Should have 3 total and 3 active sessions
	helper.AssertEqual(3, sm.GetSessionCount())
	helper.AssertEqual(3, sm.GetActiveSessionCount())

	// End one session
	sm.EndSession(session3.ID)

	// Should have 3 total but only 2 active sessions
	helper.AssertEqual(3, sm.GetSessionCount())
	helper.AssertEqual(2, sm.GetActiveSessionCount())
}

func TestSessionManager_CleanupInactiveSessions(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create sessions
	_, _ = sm.CreateSession(1, "Session 1", "owner1", SessionSettings{})
	session2, _ := sm.CreateSession(2, "Session 2", "owner2", SessionSettings{})
	_, _ = sm.CreateSession(2, "Session 3", "owner3", SessionSettings{})

	// End session 2
	sm.EndSession(session2.ID)

	// Should have 3 total, 2 active sessions
	helper.AssertEqual(3, sm.GetSessionCount())
	helper.AssertEqual(2, sm.GetActiveSessionCount())

	// Cleanup with very short max age (should remove inactive session)
	removed := sm.CleanupInactiveSessions(1 * time.Millisecond)
	helper.AssertLength(removed, 1)
	helper.AssertEqual(session2.ID, removed[0].ID)

	// Should have 2 total, 2 active sessions
	helper.AssertEqual(2, sm.GetSessionCount())
	helper.AssertEqual(2, sm.GetActiveSessionCount())

	// Cleanup with longer max age (should not remove active sessions)
	removed = sm.CleanupInactiveSessions(1 * time.Hour)
	helper.AssertLength(removed, 0)

	// Should still have 2 total, 2 active sessions
	helper.AssertEqual(2, sm.GetSessionCount())
	helper.AssertEqual(2, sm.GetActiveSessionCount())
}

func TestSessionManager_SessionLifecycle(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a session
	session, err := sm.CreateSession(1, "Lifecycle Test", "owner", SessionSettings{})
	helper.AssertNoError(err)
	helper.AssertTrue(session.IsActive)

	// Add participants
	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	sm.JoinSession(session.ID, "user2", "User Two", RoleViewer)

	// Verify participants
	participants, err := sm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 2)

	// Remove one participant
	err = sm.LeaveSession(session.ID, "user1")
	helper.AssertNoError(err)

	participants, err = sm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 1)

	// Remove last participant (should not end session automatically in SessionManager)
	err = sm.LeaveSession(session.ID, "user2")
	helper.AssertNoError(err)

	// Session should still be active (unlike CollaborationManager)
	retrievedSession, err := sm.GetSession(session.ID)
	helper.AssertNoError(err)
	helper.AssertTrue(retrievedSession.IsActive)

	// Manually end the session
	err = sm.EndSession(session.ID)
	helper.AssertNoError(err)

	// Session should now be inactive
	retrievedSession, err = sm.GetSession(session.ID)
	helper.AssertNoError(err)
	helper.AssertFalse(retrievedSession.IsActive, "Session should be inactive after ending")

	// Verify events were generated
	events := sm.GetSessionEvents()
	helper.AssertTrue(len(events) > 0)
}

func TestSessionManager_ConcurrentOperations(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Create a session
	session, err := sm.CreateSession(1, "Concurrent Test", "owner", SessionSettings{})
	helper.AssertNoError(err)

	// Simulate concurrent joins
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			userID := string(rune('a' + id))
			username := "User " + userID
			sm.JoinSession(session.ID, userID, username, RoleEditor)
			done <- true
		}(i)
	}

	// Wait for all joins to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify all participants were added
	participants, err := sm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 5)

	// Simulate concurrent leaves
	for i := 0; i < 5; i++ {
		go func(id int) {
			userID := string(rune('a' + id))
			sm.LeaveSession(session.ID, userID)
			done <- true
		}(i)
	}

	// Wait for all leaves to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify all participants were removed
	participants, err = sm.GetParticipants(session.ID)
	helper.AssertNoError(err)
	helper.AssertLength(participants, 0)
}

func TestSessionManager_ErrorHandling(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	sm := setup.MockSession

	// Test operations on non-existent session
	_, err := sm.GetSession("non-existent")
	helper.AssertError(err)

	_, err = sm.GetParticipants("non-existent")
	helper.AssertError(err)

	_, err = sm.GetParticipant("non-existent", "user1")
	helper.AssertError(err)

	err = sm.LeaveSession("non-existent", "user1")
	helper.AssertError(err)

	err = sm.EndSession("non-existent")
	helper.AssertError(err)

	err = sm.UpdateParticipant("non-existent", "user1", func(p *Participant) {})
	helper.AssertError(err)

	// Test operations with invalid user
	session, _ := sm.CreateSession(1, "Test", "owner", SessionSettings{})

	err = sm.LeaveSession(session.ID, "non-existent-user")
	helper.AssertError(err)

	_, err = sm.GetParticipant(session.ID, "non-existent-user")
	helper.AssertError(err)

	err = sm.UpdateParticipant(session.ID, "non-existent-user", func(p *Participant) {})
	helper.AssertError(err)

	// Test duplicate join
	sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	_, err = sm.JoinSession(session.ID, "user1", "User One", RoleEditor)
	helper.AssertError(err)
}
