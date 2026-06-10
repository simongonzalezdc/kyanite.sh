//go:build collaboration

// Collaboration tests are only run when the collaboration build tag is specified.
// Run with: go test -tags collaboration ./internal/collaboration/...

package collaboration

import (
	"testing"
	"time"
)

func TestDebugMergeStrategy(t *testing.T) {
	cr := NewConflictResolver()

	// Create a conflict with multiple operations
	now := time.Now()
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
				Timestamp: now,
			},
			{
				ID:        "op2",
				SessionID: "session1",
				UserID:    "user2",
				Type:      OpInsert,
				Position:  3,
				Content:   "Beautiful",
				Length:    9,
				Timestamp: now.Add(1 * time.Second), // Later timestamp
			},
		},
		Description: "Test conflict",
		CreatedAt:   now,
		Resolved:    false,
	}

	// Check initial state
	t.Logf("Initial state - Resolved: %v", conflict.Resolved)
	t.Logf("Initial state - Resolution: %v", conflict.Resolution)

	// Resolve with merge strategy
	resolution, err := cr.ResolveConflict(conflict, "merge", "resolver1")
	if err != nil {
		t.Fatalf("Failed to resolve conflict: %v", err)
	}

	// Check resolution
	t.Logf("Resolution: %+v", resolution)
	t.Logf("Final content: '%s'", resolution.FinalContent)

	// Check conflict state after resolution
	t.Logf("After resolution - Resolved: %v", conflict.Resolved)
	t.Logf("After resolution - Resolution: %v", conflict.Resolution)

	if !conflict.Resolved {
		t.Error("Expected conflict to be marked as resolved")
	}

	if conflict.Resolution == nil {
		t.Error("Expected conflict to have a resolution")
	}

	if resolution.ID != conflict.Resolution.ID {
		t.Errorf("Expected resolution ID to match conflict resolution ID, got %s vs %s", resolution.ID, conflict.Resolution.ID)
	}
}
