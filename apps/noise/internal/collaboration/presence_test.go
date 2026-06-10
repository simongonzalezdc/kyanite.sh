//go:build collaboration

// Collaboration tests are only run when the collaboration build tag is specified.
// Run with: go test -tags collaboration ./internal/collaboration/...

package collaboration

import (
	"testing"
	"time"
)

func TestPresenceManager_UpdateUserPresence(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	tests := []struct {
		name        string
		userID      string
		username    string
		status      PresenceStatus
		sessionID   string
		expectEvent bool
	}{
		{
			name:        "New user online",
			userID:      "user1",
			username:    "User One",
			status:      StatusOnline,
			sessionID:   "session1",
			expectEvent: true,
		},
		{
			name:        "User status change to away",
			userID:      "user1",
			username:    "User One",
			status:      StatusAway,
			sessionID:   "session1",
			expectEvent: true,
		},
		{
			name:        "User status change to busy",
			userID:      "user1",
			username:    "User One",
			status:      StatusBusy,
			sessionID:   "session1",
			expectEvent: true,
		},
		{
			name:        "User status change to offline",
			userID:      "user1",
			username:    "User One",
			status:      StatusOffline,
			sessionID:   "",
			expectEvent: true,
		},
		{
			name:        "Same status no event",
			userID:      "user1",
			username:    "User One",
			status:      StatusOffline,
			sessionID:   "",
			expectEvent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm.ClearPresenceEvents()

			pm.UpdateUserPresence(tt.userID, tt.username, tt.status, tt.sessionID)

			// Verify presence was updated
			presence, exists := pm.GetUserPresence(tt.userID)
			helper.AssertTrue(exists, "User presence should exist")
			helper.AssertEqual(tt.userID, presence.UserID)
			helper.AssertEqual(tt.username, presence.Username)
			helper.AssertEqual(tt.status, presence.Status)
			helper.AssertEqual(tt.sessionID, presence.CurrentSession)

			// Check if event was generated
			events := pm.GetPresenceEvents()
			if tt.expectEvent {
				helper.AssertTrue(len(events) > 0, "Expected presence event to be generated")
			} else {
				helper.AssertEqual(0, len(events))
			}
		})
	}
}

func TestPresenceManager_UpdateCursorPosition(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add user to session for cursor tracking
	pm.AddUserToSession("session1", "user1", "User One", RoleEditor)
	pm.AddUserToSession("session1", "user2", "User Two", RoleViewer)

	tests := []struct {
		name        string
		userID      string
		sessionID   string
		position    CursorPosition
		expectEvent bool
	}{
		{
			name:        "Valid cursor update",
			userID:      "user1",
			sessionID:   "session1",
			position:    CursorPosition{Line: 5, Column: 10},
			expectEvent: true,
		},
		{
			name:        "Another user cursor update",
			userID:      "user2",
			sessionID:   "session1",
			position:    CursorPosition{Line: 2, Column: 7},
			expectEvent: true,
		},
		{
			name:        "Non-existent user cursor update",
			userID:      "nonexistent",
			sessionID:   "session1",
			position:    CursorPosition{Line: 1, Column: 1},
			expectEvent: false, // Should not generate event for non-existent user
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm.ClearPresenceEvents()

			pm.UpdateCursorPosition(tt.userID, tt.sessionID, tt.position)

			// Check if event was generated
			events := pm.GetPresenceEvents()
			if tt.expectEvent {
				helper.AssertTrue(len(events) > 0, "Expected cursor event to be generated")

				// Verify event details
				found := false
				for _, event := range events {
					if event.Type == EventCursorMoved && event.UserID == tt.userID {
						found = true
						cursor, ok := event.Data["cursor"].(CursorPosition)
						helper.AssertTrue(ok, "Cursor position should be in event data")
						helper.AssertEqual(tt.position.Line, cursor.Line)
						helper.AssertEqual(tt.position.Column, cursor.Column)
						break
					}
				}
				helper.AssertTrue(found, "Expected cursor event for user")
			} else {
				// For non-existent users, no event should be generated
				helper.AssertEqual(0, len(events))
			}

			// Verify cursor was updated for existing users
			if tt.userID == "user1" || tt.userID == "user2" {
				participants, exists := pm.GetSessionParticipants("session1")
				helper.AssertTrue(exists, "Session participants should exist")

				found := false
				for _, p := range participants {
					if p.UserID == tt.userID {
						found = true
						helper.AssertEqual(tt.position.Line, p.Cursor.Line)
						helper.AssertEqual(tt.position.Column, p.Cursor.Column)
						break
					}
				}
				helper.AssertTrue(found, "Participant cursor should be updated")
			}
		})
	}
}

func TestPresenceManager_AddUserToSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		username    string
		role        ParticipantRole
		expectEvent bool
	}{
		{
			name:        "Add user as editor",
			sessionID:   "session1",
			userID:      "user1",
			username:    "User One",
			role:        RoleEditor,
			expectEvent: true,
		},
		{
			name:        "Add user as viewer",
			sessionID:   "session1",
			userID:      "user2",
			username:    "User Two",
			role:        RoleViewer,
			expectEvent: true,
		},
		{
			name:        "Add user to different session",
			sessionID:   "session2",
			userID:      "user3",
			username:    "User Three",
			role:        RoleEditor,
			expectEvent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm.ClearPresenceEvents()

			pm.AddUserToSession(tt.sessionID, tt.userID, tt.username, tt.role)

			// Verify user was added to session
			participants, exists := pm.GetSessionParticipants(tt.sessionID)
			helper.AssertTrue(exists, "Session should exist")

			found := false
			for _, p := range participants {
				if p.UserID == tt.userID {
					found = true
					helper.AssertEqual(tt.username, p.Username)
					helper.AssertEqual(tt.role, p.Role)
					helper.AssertTrue(p.IsActive)
					break
				}
			}
			helper.AssertTrue(found, "User should be in session participants")

			// Verify user presence was updated
			presence, exists := pm.GetUserPresence(tt.userID)
			helper.AssertTrue(exists, "User presence should exist")
			helper.AssertEqual(tt.userID, presence.UserID)
			helper.AssertEqual(tt.username, presence.Username)
			helper.AssertEqual(StatusOnline, presence.Status)
			helper.AssertEqual(tt.sessionID, presence.CurrentSession)

			// Check if event was generated
			events := pm.GetPresenceEvents()
			if tt.expectEvent {
				helper.AssertTrue(len(events) > 0, "Expected user joined event to be generated")

				// Verify event details
				found := false
				for _, event := range events {
					if event.Type == EventUserJoined && event.UserID == tt.userID {
						found = true
						helper.AssertEqual(tt.sessionID, event.SessionID)
						break
					}
				}
				helper.AssertTrue(found, "Expected user joined event")
			}
		})
	}
}

func TestPresenceManager_RemoveUserFromSession(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add users to sessions
	pm.AddUserToSession("session1", "user1", "User One", RoleEditor)
	pm.AddUserToSession("session1", "user2", "User Two", RoleViewer)
	pm.AddUserToSession("session2", "user3", "User Three", RoleEditor)

	tests := []struct {
		name        string
		sessionID   string
		userID      string
		expectEvent bool
	}{
		{
			name:        "Remove user from session",
			sessionID:   "session1",
			userID:      "user1",
			expectEvent: true,
		},
		{
			name:        "Remove another user from session",
			sessionID:   "session1",
			userID:      "user2",
			expectEvent: true,
		},
		{
			name:        "Remove user from different session",
			sessionID:   "session2",
			userID:      "user3",
			expectEvent: true,
		},
		{
			name:        "Remove non-existent user",
			sessionID:   "session1",
			userID:      "nonexistent",
			expectEvent: false,
		},
		{
			name:        "Remove from non-existent session",
			sessionID:   "nonexistent",
			userID:      "user1",
			expectEvent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm.ClearPresenceEvents()

			pm.RemoveUserFromSession(tt.sessionID, tt.userID)

			// Check if event was generated
			events := pm.GetPresenceEvents()
			if tt.expectEvent {
				helper.AssertTrue(len(events) > 0, "Expected user left event to be generated")

				// Verify event details
				found := false
				for _, event := range events {
					if event.Type == EventUserLeft && event.UserID == tt.userID {
						found = true
						helper.AssertEqual(tt.sessionID, event.SessionID)
						break
					}
				}
				helper.AssertTrue(found, "Expected user left event")
			} else {
				// For non-existent users/sessions, no event should be generated
				helper.AssertEqual(0, len(events))
			}

			// Verify user presence was updated to offline (for existing users)
			if tt.userID == "user1" || tt.userID == "user2" || tt.userID == "user3" {
				presence, exists := pm.GetUserPresence(tt.userID)
				helper.AssertTrue(exists, "User presence should still exist")
				helper.AssertEqual(StatusOffline, presence.Status)
				helper.AssertEqual("", presence.CurrentSession)
			}
		})
	}
}

func TestPresenceManager_GetActiveUsers(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Initially no active users
	activeUsers := pm.GetActiveUsers()
	helper.AssertLength(activeUsers, 0)

	// Add users with different statuses
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")
	pm.UpdateUserPresence("user2", "User Two", StatusAway, "session1")
	pm.UpdateUserPresence("user3", "User Three", StatusBusy, "session2")
	pm.UpdateUserPresence("user4", "User Four", StatusOffline, "")

	// Should have 1 active user (online only)
	activeUsers = pm.GetActiveUsers()
	helper.AssertLength(activeUsers, 1)
	helper.AssertEqual("user1", activeUsers[0].UserID)
	helper.AssertEqual(StatusOnline, activeUsers[0].Status)
}

func TestPresenceManager_GetUsersByStatus(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add users with different statuses
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")
	pm.UpdateUserPresence("user2", "User Two", StatusAway, "session1")
	pm.UpdateUserPresence("user3", "User Three", StatusBusy, "session2")
	pm.UpdateUserPresence("user4", "User Four", StatusOffline, "")

	tests := []struct {
		name          string
		status        PresenceStatus
		expectedCount int
		expectedUsers []string
	}{
		{
			name:          "Get online users",
			status:        StatusOnline,
			expectedCount: 1,
			expectedUsers: []string{"user1"},
		},
		{
			name:          "Get away users",
			status:        StatusAway,
			expectedCount: 1,
			expectedUsers: []string{"user2"},
		},
		{
			name:          "Get busy users",
			status:        StatusBusy,
			expectedCount: 1,
			expectedUsers: []string{"user3"},
		},
		{
			name:          "Get offline users",
			status:        StatusOffline,
			expectedCount: 1,
			expectedUsers: []string{"user4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := pm.GetUsersByStatus(tt.status)
			helper.AssertLength(users, tt.expectedCount)

			if len(tt.expectedUsers) > 0 {
				userIDs := make(map[string]bool)
				for _, user := range users {
					userIDs[user.UserID] = true
				}

				for _, expectedUserID := range tt.expectedUsers {
					helper.AssertTrue(userIDs[expectedUserID], "Expected user %s not found", expectedUserID)
				}
			}
		})
	}
}

func TestPresenceManager_SendHeartbeat(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add user
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")

	// Get initial last seen time
	presence, _ := pm.GetUserPresence("user1")
	initialLastSeen := presence.LastSeen

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Send heartbeat
	pm.ClearPresenceEvents()
	pm.SendHeartbeat("user1")

	// Verify last seen was updated
	presence, _ = pm.GetUserPresence("user1")
	helper.AssertTrue(presence.LastSeen.After(initialLastSeen), "Last seen should be updated after heartbeat")

	// Verify heartbeat event was generated
	events := pm.GetPresenceEvents()
	helper.AssertTrue(len(events) > 0, "Expected heartbeat event to be generated")

	found := false
	for _, event := range events {
		if event.Type == EventHeartbeat && event.UserID == "user1" {
			found = true
			break
		}
	}
	helper.AssertTrue(found, "Expected heartbeat event for user")

	// Test heartbeat for non-existent user
	pm.ClearPresenceEvents()
	pm.SendHeartbeat("nonexistent")

	// No event should be generated for non-existent user
	events = pm.GetPresenceEvents()
	helper.AssertEqual(0, len(events))
}

func TestPresenceManager_SetUserDeviceInfo(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add user
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")

	// Set device info
	deviceInfo := map[string]string{
		"browser":    "Chrome",
		"os":         "Windows",
		"screen":     "1920x1080",
		"connection": "wifi",
	}

	pm.SetUserDeviceInfo("user1", deviceInfo)

	// Verify device info was set
	presence, _ := pm.GetUserPresence("user1")
	helper.AssertEqual(4, len(presence.DeviceInfo))
	helper.AssertEqual("Chrome", presence.DeviceInfo["browser"])
	helper.AssertEqual("Windows", presence.DeviceInfo["os"])
	helper.AssertEqual("1920x1080", presence.DeviceInfo["screen"])
	helper.AssertEqual("wifi", presence.DeviceInfo["connection"])

	// Test device info for non-existent user (should not panic)
	pm.SetUserDeviceInfo("nonexistent", map[string]string{"test": "value"})
}

func TestPresenceManager_GetPresenceIndicator(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	tests := []struct {
		name          string
		userID        string
		status        PresenceStatus
		expectedColor string
		expectedIcon  string
	}{
		{
			name:          "Online user",
			userID:        "user1",
			status:        StatusOnline,
			expectedColor: "green",
			expectedIcon:  "●",
		},
		{
			name:          "Away user",
			userID:        "user2",
			status:        StatusAway,
			expectedColor: "yellow",
			expectedIcon:  "●",
		},
		{
			name:          "Busy user",
			userID:        "user3",
			status:        StatusBusy,
			expectedColor: "red",
			expectedIcon:  "●",
		},
		{
			name:          "Offline user",
			userID:        "user4",
			status:        StatusOffline,
			expectedColor: "gray",
			expectedIcon:  "○",
		},
		{
			name:          "Non-existent user",
			userID:        "nonexistent",
			status:        StatusOffline,
			expectedColor: "gray",
			expectedIcon:  "○",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.userID != "nonexistent" {
				pm.UpdateUserPresence(tt.userID, "Test User", tt.status, "session1")
			}

			indicator := pm.GetPresenceIndicator(tt.userID)
			helper.AssertEqual(tt.status, indicator.Status)
			helper.AssertEqual(tt.expectedColor, indicator.Color)
			helper.AssertEqual(tt.expectedIcon, indicator.Icon)
		})
	}
}

func TestPresenceManager_GetSessionPresenceIndicators(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add users to session with different statuses and roles
	pm.AddUserToSession("session1", "user1", "User One", RoleEditor)
	pm.AddUserToSession("session1", "user2", "User Two", RoleViewer)
	pm.AddUserToSession("session1", "user3", "User Three", RoleOwner)

	// Update user statuses
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")
	pm.UpdateUserPresence("user2", "User Two", StatusAway, "session1")
	pm.UpdateUserPresence("user3", "User Three", StatusBusy, "session1")

	// Update cursor positions
	pm.UpdateCursorPosition("user1", "session1", CursorPosition{Line: 5, Column: 10})
	pm.UpdateCursorPosition("user2", "session1", CursorPosition{Line: 2, Column: 7})

	// Get session presence indicators
	indicators := pm.GetSessionPresenceIndicators("session1")
	helper.AssertLength(indicators, 3)

	// Verify indicators
	userMap := make(map[string]SessionPresenceIndicator)
	for _, indicator := range indicators {
		userMap[indicator.UserID] = indicator
	}

	// Check user1
	user1Indicator := userMap["user1"]
	helper.AssertEqual("user1", user1Indicator.UserID)
	helper.AssertEqual("User One", user1Indicator.Username)
	helper.AssertEqual(StatusOnline, user1Indicator.Indicator.Status)
	helper.AssertEqual(RoleEditor, user1Indicator.Role)
	helper.AssertEqual(5, user1Indicator.Cursor.Line)
	helper.AssertEqual(10, user1Indicator.Cursor.Column)

	// Check user2
	user2Indicator := userMap["user2"]
	helper.AssertEqual("user2", user2Indicator.UserID)
	helper.AssertEqual("User Two", user2Indicator.Username)
	helper.AssertEqual(StatusAway, user2Indicator.Indicator.Status)
	helper.AssertEqual(RoleViewer, user2Indicator.Role)
	helper.AssertEqual(2, user2Indicator.Cursor.Line)
	helper.AssertEqual(7, user2Indicator.Cursor.Column)

	// Check user3
	user3Indicator := userMap["user3"]
	helper.AssertEqual("user3", user3Indicator.UserID)
	helper.AssertEqual("User Three", user3Indicator.Username)
	helper.AssertEqual(StatusBusy, user3Indicator.Indicator.Status)
	helper.AssertEqual(RoleOwner, user3Indicator.Role)

	// Test non-existent session
	indicators = pm.GetSessionPresenceIndicators("nonexistent")
	helper.AssertLength(indicators, 0)
}

func TestPresenceManager_IdleUserCleanup(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add users
	pm.UpdateUserPresence("user1", "User One", StatusOnline, "session1")
	pm.UpdateUserPresence("user2", "User Two", StatusOnline, "session1")
	pm.UpdateUserPresence("user3", "User Three", StatusOnline, "session1")

	// Manually set last seen times to simulate idle users
	presence1, _ := pm.GetUserPresence("user1")
	presence2, _ := pm.GetUserPresence("user2")
	presence3, _ := pm.GetUserPresence("user3")

	// Make user1 and user2 idle (old last seen)
	presence1.LastSeen = time.Now().Add(-10 * time.Minute)
	presence2.LastSeen = time.Now().Add(-10 * time.Minute)

	// Keep user3 active (recent last seen)
	presence3.LastSeen = time.Now().Add(-1 * time.Minute)

	// Wait for cleanup to run (cleanup runs every 30 seconds in real implementation)
	// For testing, we'll trigger the cleanup manually
	pm.markIdleUsersOffline()

	// Verify idle users were marked as away
	presence1, _ = pm.GetUserPresence("user1")
	presence2, _ = pm.GetUserPresence("user2")
	presence3, _ = pm.GetUserPresence("user3")

	helper.AssertEqual(StatusAway, presence1.Status)
	helper.AssertEqual(StatusAway, presence2.Status)
	helper.AssertEqual(StatusOnline, presence3.Status)
}

func TestPresenceManager_ConcurrentOperations(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Simulate concurrent presence updates
	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(id int) {
			userID := string(rune('a' + id))
			username := "User " + userID
			sessionID := "session1"

			pm.AddUserToSession(sessionID, userID, username, RoleEditor)

			// Simulate cursor movements
			for j := 0; j < 10; j++ {
				pm.UpdateCursorPosition(userID, sessionID, CursorPosition{Line: j, Column: j * 2})
				time.Sleep(1 * time.Millisecond)
			}

			done <- true
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Verify all users are in the session
	participants, exists := pm.GetSessionParticipants("session1")
	helper.AssertTrue(exists, "Session should exist")
	helper.AssertLength(participants, 5)

	// Verify all users have presence
	activeUsers := pm.GetActiveUsers()
	helper.AssertLength(activeUsers, 5)

	// Verify events were generated
	events := pm.GetPresenceEvents()
	helper.AssertTrue(len(events) > 0)
}

func TestPresenceManager_Close(t *testing.T) {
	setup := NewTestSetup(t)
	defer setup.Cleanup()

	helper := setup.TestHelper
	pm := setup.MockPresence

	// Add some users
	pm.AddUserToSession("session1", "user1", "User One", RoleEditor)
	pm.AddUserToSession("session1", "user2", "User Two", RoleViewer)

	// Close the presence manager
	err := pm.Close()
	helper.AssertNoError(err)
}
