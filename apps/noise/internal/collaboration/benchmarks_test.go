//go:build collaboration

// Collaboration tests are only run when the collaboration build tag is specified.
// Run with: go test -tags collaboration ./internal/collaboration/...

package collaboration

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkCollaborationManager_CreateSession benchmarks session creation
func BenchmarkCollaborationManager_CreateSession(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session, err := cm.CreateSession(
			i,
			fmt.Sprintf("Session %d", i),
			fmt.Sprintf("user%d", i),
			SessionSettings{
				MaxParticipants:  5,
				AutoSaveInterval: 30 * time.Second,
				ConflictStrategy: "merge",
				RequireApproval:  false,
			},
		)
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
		_ = session
	}
}

// BenchmarkCollaborationManager_JoinSession benchmarks user joining sessions
func BenchmarkCollaborationManager_JoinSession(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session for benchmarking
	session, err := cm.CreateSession(1, "Benchmark Session", "owner", SessionSettings{
		MaxParticipants: 1000,
	})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i)
		_, err := cm.JoinSession(session.ID, userID, fmt.Sprintf("User %d", i), RoleEditor)
		if err != nil {
			b.Fatalf("Failed to join session: %v", err)
		}
	}
}

// BenchmarkCollaborationManager_ApplyOperation benchmarks operation application
func BenchmarkCollaborationManager_ApplyOperation(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session and add a user
	session, err := cm.CreateSession(1, "Benchmark Session", "owner", SessionSettings{})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	_, err = cm.JoinSession(session.ID, "user1", "User 1", RoleEditor)
	if err != nil {
		b.Fatalf("Failed to join session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		op := Operation{
			Type:      OpInsert,
			Position:  i * 10,
			Content:   fmt.Sprintf("Content %d", i),
			Length:    len(fmt.Sprintf("Content %d", i)),
			Timestamp: time.Now(),
		}
		err := cm.ApplyOperation(session.ID, "user1", op)
		if err != nil {
			b.Fatalf("Failed to apply operation: %v", err)
		}
	}
}

// BenchmarkCollaborationManager_UpdateCursor benchmarks cursor position updates
func BenchmarkCollaborationManager_UpdateCursor(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session and add a user
	session, err := cm.CreateSession(1, "Benchmark Session", "owner", SessionSettings{})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	_, err = cm.JoinSession(session.ID, "user1", "User 1", RoleEditor)
	if err != nil {
		b.Fatalf("Failed to join session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		position := CursorPosition{
			Line:   i % 100,
			Column: i % 80,
		}
		err := cm.UpdateCursor(session.ID, "user1", position)
		if err != nil {
			b.Fatalf("Failed to update cursor: %v", err)
		}
	}
}

// BenchmarkSessionManager_CreateSession benchmarks session creation in SessionManager
func BenchmarkSessionManager_CreateSession(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	sm := setup.MockSession

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		session, err := sm.CreateSession(
			i,
			fmt.Sprintf("Session %d", i),
			fmt.Sprintf("user%d", i),
			SessionSettings{
				MaxParticipants:  5,
				AutoSaveInterval: 30 * time.Second,
				ConflictStrategy: "merge",
				RequireApproval:  false,
			},
		)
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
		_ = session
	}
}

// BenchmarkSessionManager_JoinSession benchmarks user joining sessions in SessionManager
func BenchmarkSessionManager_JoinSession(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	sm := setup.MockSession

	// Create a session for benchmarking
	session, err := sm.CreateSession(1, "Benchmark Session", "owner", SessionSettings{
		MaxParticipants: 1000,
	})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i)
		_, err := sm.JoinSession(session.ID, userID, fmt.Sprintf("User %d", i), RoleEditor)
		if err != nil {
			b.Fatalf("Failed to join session: %v", err)
		}
	}
}

// BenchmarkPresenceManager_UpdateUserPresence benchmarks user presence updates
func BenchmarkPresenceManager_UpdateUserPresence(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	pm := setup.MockPresence

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i%100) // Reuse 100 users
		username := fmt.Sprintf("User %d", i%100)
		status := StatusOnline
		if i%4 == 1 {
			status = StatusAway
		} else if i%4 == 2 {
			status = StatusBusy
		} else if i%4 == 3 {
			status = StatusOffline
		}
		sessionID := fmt.Sprintf("session%d", i%10) // Reuse 10 sessions

		pm.UpdateUserPresence(userID, username, status, sessionID)
	}
}

// BenchmarkPresenceManager_UpdateCursorPosition benchmarks cursor position updates
func BenchmarkPresenceManager_UpdateCursorPosition(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	pm := setup.MockPresence

	// Add users to sessions
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user%d", i)
		sessionID := fmt.Sprintf("session%d", i%10)
		pm.AddUserToSession(sessionID, userID, fmt.Sprintf("User %d", i), RoleEditor)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i%100)
		sessionID := fmt.Sprintf("session%d", i%10)
		position := CursorPosition{
			Line:   i % 100,
			Column: i % 80,
		}
		pm.UpdateCursorPosition(userID, sessionID, position)
	}
}

// BenchmarkConflictResolver_DetectConflicts benchmarks conflict detection
func BenchmarkConflictResolver_DetectConflicts(b *testing.B) {
	cr := NewConflictResolver()

	// Create a large set of operations
	var operations []Operation
	for i := 0; i < 1000; i++ {
		operations = append(operations, Operation{
			ID:        fmt.Sprintf("op%d", i),
			SessionID: "session1",
			UserID:    fmt.Sprintf("user%d", i%10),
			Type:      OpInsert,
			Position:  i * 5,
			Content:   fmt.Sprintf("Content %d", i),
			Length:    len(fmt.Sprintf("Content %d", i)),
			Timestamp: time.Now(),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conflicts := cr.DetectConflicts(operations)
		_ = conflicts
	}
}

// BenchmarkConflictResolver_ResolveConflict benchmarks conflict resolution
func BenchmarkConflictResolver_ResolveConflict(b *testing.B) {
	cr := NewConflictResolver()

	// Create a conflict
	conflict := &Conflict{
		ID:        "conflict1",
		SessionID: "session1",
		Operations: []Operation{
			{
				ID:        "op1",
				SessionID: "session1",
				UserID:    "user1",
				Type:      OpInsert,
				Position:  0,
				Content:   "Hello",
				Length:    5,
				Timestamp: time.Now(),
			},
			{
				ID:        "op2",
				SessionID: "session1",
				UserID:    "user2",
				Type:      OpInsert,
				Position:  3,
				Content:   "World",
				Length:    5,
				Timestamp: time.Now(),
			},
		},
		Description: "Test conflict",
		CreatedAt:   time.Now(),
		Resolved:    false,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset conflict for each iteration
		conflict.Resolved = false
		conflict.Resolution = nil

		_, err := cr.ResolveConflict(conflict, "merge", "resolver1")
		if err != nil {
			b.Fatalf("Failed to resolve conflict: %v", err)
		}
	}
}

// BenchmarkCollaboration_ConcurrentOperations benchmarks concurrent collaboration operations
func BenchmarkCollaboration_ConcurrentOperations(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session
	session, err := cm.CreateSession(1, "Concurrent Benchmark", "owner", SessionSettings{
		MaxParticipants: 100,
	})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add users
	for i := 0; i < 10; i++ {
		userID := fmt.Sprintf("user%d", i)
		_, err := cm.JoinSession(session.ID, userID, fmt.Sprintf("User %d", i), RoleEditor)
		if err != nil {
			b.Fatalf("Failed to join session: %v", err)
		}
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			userID := fmt.Sprintf("user%d", i%10)
			op := Operation{
				Type:      OpInsert,
				Position:  i * 10,
				Content:   fmt.Sprintf("Content %d", i),
				Length:    len(fmt.Sprintf("Content %d", i)),
				Timestamp: time.Now(),
			}
			err := cm.ApplyOperation(session.ID, userID, op)
			if err != nil {
				b.Errorf("Failed to apply operation: %v", err)
			}
			i++
		}
	})
}

// BenchmarkCollaboration_ManyUsers benchmarks performance with many users
func BenchmarkCollaboration_ManyUsers(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create sessions
	var sessions []*Session
	for i := 0; i < 10; i++ {
		session, err := cm.CreateSession(i, fmt.Sprintf("Session %d", i), "owner", SessionSettings{
			MaxParticipants: 100,
		})
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
		sessions = append(sessions, session)
	}

	// Add many users
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user%d", i)
		sessionID := sessions[i%10].ID
		_, err := cm.JoinSession(sessionID, userID, fmt.Sprintf("User %d", i), RoleEditor)
		if err != nil {
			b.Fatalf("Failed to join session: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i%100)
		sessionID := sessions[i%10].ID

		// Get participants
		_, err := cm.GetParticipants(sessionID)
		if err != nil {
			b.Errorf("Failed to get participants: %v", err)
		}

		// Update cursor
		position := CursorPosition{Line: i % 100, Column: i % 80}
		err = cm.UpdateCursor(sessionID, userID, position)
		if err != nil {
			b.Errorf("Failed to update cursor: %v", err)
		}
	}
}

// BenchmarkPresenceManager_ManyUsers benchmarks presence tracking with many users
func BenchmarkPresenceManager_ManyUsers(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	pm := setup.MockPresence

	// Add many users to sessions
	for i := 0; i < 1000; i++ {
		userID := fmt.Sprintf("user%d", i)
		sessionID := fmt.Sprintf("session%d", i%10)
		pm.AddUserToSession(sessionID, userID, fmt.Sprintf("User %d", i), RoleEditor)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := fmt.Sprintf("user%d", i%1000)
		sessionID := fmt.Sprintf("session%d", i%10)

		// Update cursor
		position := CursorPosition{Line: i % 100, Column: i % 80}
		pm.UpdateCursorPosition(userID, sessionID, position)

		// Get presence
		_, _ = pm.GetUserPresence(userID)

		// Get session participants
		_, _ = pm.GetSessionParticipants(sessionID)
	}
}

// BenchmarkNetworkLatency benchmarks performance under simulated network latency
func BenchmarkNetworkLatency(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session
	session, err := cm.CreateSession(1, "Latency Test", "owner", SessionSettings{})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add a user
	_, err = cm.JoinSession(session.ID, "user1", "User 1", RoleEditor)
	if err != nil {
		b.Fatalf("Failed to join session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Simulate network latency
		time.Sleep(1 * time.Millisecond)

		op := Operation{
			Type:      OpInsert,
			Position:  i * 10,
			Content:   fmt.Sprintf("Content %d", i),
			Length:    len(fmt.Sprintf("Content %d", i)),
			Timestamp: time.Now(),
		}
		err := cm.ApplyOperation(session.ID, "user1", op)
		if err != nil {
			b.Errorf("Failed to apply operation: %v", err)
		}
	}
}

// BenchmarkMemoryUsage benchmarks memory usage with many operations
func BenchmarkMemoryUsage(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create a session
	session, err := cm.CreateSession(1, "Memory Test", "owner", SessionSettings{})
	if err != nil {
		b.Fatalf("Failed to create session: %v", err)
	}

	// Add a user
	_, err = cm.JoinSession(session.ID, "user1", "User 1", RoleEditor)
	if err != nil {
		b.Fatalf("Failed to join session: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Apply many operations to build up memory usage
		for j := 0; j < 100; j++ {
			op := Operation{
				Type:      OpInsert,
				Position:  (i*100 + j) * 10,
				Content:   fmt.Sprintf("Content %d-%d", i, j),
				Length:    len(fmt.Sprintf("Content %d-%d", i, j)),
				Timestamp: time.Now(),
			}
			err := cm.ApplyOperation(session.ID, "user1", op)
			if err != nil {
				b.Errorf("Failed to apply operation: %v", err)
			}
		}

		// Get session to ensure operations are loaded
		_, err := cm.GetSession(session.ID)
		if err != nil {
			b.Errorf("Failed to get session: %v", err)
		}
	}
}

// BenchmarkConcurrentSessions benchmarks performance with multiple concurrent sessions
func BenchmarkConcurrentSessions(b *testing.B) {
	setup := NewTestSetup(&testing.T{})
	defer setup.Cleanup()

	cm := setup.MockCollaboration

	// Create multiple sessions
	var sessions []*Session
	for i := 0; i < 10; i++ {
		session, err := cm.CreateSession(i, fmt.Sprintf("Session %d", i), "owner", SessionSettings{})
		if err != nil {
			b.Fatalf("Failed to create session: %v", err)
		}
		sessions = append(sessions, session)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sessionID := sessions[i%10].ID
			userID := fmt.Sprintf("user%d", i%100)

			// Join session if not already joined
			if i < 100 {
				_, err := cm.JoinSession(sessionID, userID, fmt.Sprintf("User %d", i), RoleEditor)
				if err != nil {
					// User might already be in session, that's okay
				}
			}

			// Apply operation
			op := Operation{
				Type:      OpInsert,
				Position:  i * 10,
				Content:   fmt.Sprintf("Content %d", i),
				Length:    len(fmt.Sprintf("Content %d", i)),
				Timestamp: time.Now(),
			}
			err := cm.ApplyOperation(sessionID, userID, op)
			if err != nil {
				b.Errorf("Failed to apply operation: %v", err)
			}
			i++
		}
	})
}
