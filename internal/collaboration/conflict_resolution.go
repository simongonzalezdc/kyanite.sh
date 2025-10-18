package collaboration

import (
	"fmt"
	"strings"
	"time"
)

// ConflictResolver handles concurrent editing conflicts
type ConflictResolver struct {
	strategies map[string]ConflictStrategy
}

// ConflictStrategy defines how to resolve conflicts
type ConflictStrategy interface {
	Resolve(conflict *Conflict) (*Resolution, error)
	Name() string
}

// Conflict represents a concurrent editing conflict
type Conflict struct {
	ID          string      `json:"id"`
	SessionID   string      `json:"session_id"`
	Operations  []Operation `json:"operations"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Resolved    bool        `json:"resolved"`
	Resolution  *Resolution `json:"resolution,omitempty"`
}

// Resolution represents how a conflict was resolved
type Resolution struct {
	ID           string                 `json:"id"`
	ConflictID   string                 `json:"conflict_id"`
	Strategy     string                 `json:"strategy"`
	ResolvedBy   string                 `json:"resolved_by"`
	ResolvedAt   time.Time              `json:"resolved_at"`
	FinalContent string                 `json:"final_content"`
	Changes      []DocumentChange       `json:"changes"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// DocumentChange represents a change made during conflict resolution
type DocumentChange struct {
	Type       string    `json:"type"` // "insert", "delete", "replace"
	Position   int       `json:"position"`
	OldContent string    `json:"old_content"`
	NewContent string    `json:"new_content"`
	UserID     string    `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// MergeStrategy implements automatic conflict resolution by merging changes
type MergeStrategy struct{}

// LockStrategy implements conflict resolution by locking regions
type LockStrategy struct{}

// ManualStrategy implements conflict resolution requiring user intervention
type ManualStrategy struct{}

// NewConflictResolver creates a new conflict resolver with default strategies
func NewConflictResolver() *ConflictResolver {
	cr := &ConflictResolver{
		strategies: make(map[string]ConflictStrategy),
	}

	// Register default strategies
	cr.RegisterStrategy(&MergeStrategy{})
	cr.RegisterStrategy(&LockStrategy{})
	cr.RegisterStrategy(&ManualStrategy{})

	return cr
}

// RegisterStrategy registers a conflict resolution strategy
func (cr *ConflictResolver) RegisterStrategy(strategy ConflictStrategy) {
	cr.strategies[strategy.Name()] = strategy
}

// ResolveConflict resolves a conflict using the specified strategy
func (cr *ConflictResolver) ResolveConflict(conflict *Conflict, strategyName, userID string) (*Resolution, error) {
	strategy, exists := cr.strategies[strategyName]
	if !exists {
		return nil, fmt.Errorf("unknown conflict resolution strategy: %s", strategyName)
	}

	resolution, err := strategy.Resolve(conflict)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve conflict: %w", err)
	}

	// Update resolution metadata
	resolution.ID = generateResolutionID()
	resolution.ConflictID = conflict.ID
	resolution.ResolvedBy = userID
	resolution.ResolvedAt = time.Now()

	// Mark conflict as resolved
	conflict.Resolved = true
	conflict.Resolution = resolution

	return resolution, nil
}

// DetectConflicts analyzes operations for potential conflicts
func (cr *ConflictResolver) DetectConflicts(operations []Operation) []*Conflict {
	var conflicts []*Conflict

	// Group operations by position ranges
	operationGroups := cr.groupOperationsByPosition(operations)

	for _, group := range operationGroups {
		if len(group) > 1 {
			conflict := &Conflict{
				ID:          generateConflictID(),
				SessionID:   group[0].SessionID,
				Operations:  group,
				Description: cr.describeConflict(group),
				CreatedAt:   time.Now(),
				Resolved:    false,
			}
			conflicts = append(conflicts, conflict)
		}
	}

	return conflicts
}

// MergeStrategy implementation

func (ms *MergeStrategy) Name() string {
	return "merge"
}

func (ms *MergeStrategy) Resolve(conflict *Conflict) (*Resolution, error) {
	if len(conflict.Operations) < 2 {
		return nil, fmt.Errorf("merge strategy requires at least 2 operations")
	}

	// Simple merge strategy: prefer the most recent operation
	// In a full implementation, this would be more sophisticated
	latestOp := conflict.Operations[0]
	for _, op := range conflict.Operations[1:] {
		if op.Timestamp.After(latestOp.Timestamp) {
			latestOp = op
		}
	}

	// Create resolution based on the latest operation
	resolution := &Resolution{
		Strategy:     "merge",
		FinalContent: latestOp.Content,
		Changes: []DocumentChange{
			{
				Type:       string(latestOp.Type),
				Position:   latestOp.Position,
				NewContent: latestOp.Content,
				UserID:     latestOp.UserID,
				Timestamp:  latestOp.Timestamp,
			},
		},
		Metadata: map[string]interface{}{
			"merge_type":       "latest_wins",
			"operations_count": len(conflict.Operations),
		},
	}

	return resolution, nil
}

// LockStrategy implementation

func (ls *LockStrategy) Name() string {
	return "lock"
}

func (ls *LockStrategy) Resolve(conflict *Conflict) (*Resolution, error) {
	// Lock strategy: create locked regions that prevent further editing
	// Find the range that needs to be locked
	minPos := conflict.Operations[0].Position
	maxPos := minPos + conflict.Operations[0].Length

	for _, op := range conflict.Operations[1:] {
		if op.Position < minPos {
			minPos = op.Position
		}
		if op.Position+op.Length > maxPos {
			maxPos = op.Position + op.Length
		}
	}

	lockRange := maxPos - minPos

	resolution := &Resolution{
		Strategy:     "lock",
		FinalContent: fmt.Sprintf("[LOCKED:%d-%d]", minPos, maxPos),
		Changes: []DocumentChange{
			{
				Type:       "lock",
				Position:   minPos,
				NewContent: fmt.Sprintf("[LOCKED:%d-%d]", minPos, maxPos),
				Timestamp:  time.Now(),
			},
		},
		Metadata: map[string]interface{}{
			"lock_range_start": minPos,
			"lock_range_end":   maxPos,
			"lock_range_size":  lockRange,
		},
	}

	return resolution, nil
}

// ManualStrategy implementation

func (ms *ManualStrategy) Name() string {
	return "manual"
}

func (ms *ManualStrategy) Resolve(conflict *Conflict) (*Resolution, error) {
	// Manual strategy: present options to user for resolution
	// This would typically involve UI interaction

	resolution := &Resolution{
		Strategy:     "manual",
		FinalContent: "", // To be determined by user
		Metadata: map[string]interface{}{
			"requires_user_input": true,
			"conflict_options": []string{
				"accept_all_changes",
				"reject_all_changes",
				"merge_selectively",
			},
		},
	}

	return resolution, nil
}

// Helper methods

func (cr *ConflictResolver) groupOperationsByPosition(operations []Operation) [][]Operation {
	groups := make([][]Operation, 0)

	for _, op := range operations {
		added := false

		// Check if this operation overlaps with existing groups
		for i, group := range groups {
			if cr.operationsOverlap(group[0], op) {
				groups[i] = append(groups[i], op)
				added = true
				break
			}
		}

		// If no overlap found, create new group
		if !added {
			groups = append(groups, []Operation{op})
		}
	}

	return groups
}

func (cr *ConflictResolver) operationsOverlap(op1, op2 Operation) bool {
	// Two operations overlap if their position ranges intersect
	op1End := op1.Position + op1.Length
	op2End := op2.Position + op2.Length

	return !(op1End <= op2.Position || op2End <= op1.Position)
}

func (cr *ConflictResolver) describeConflict(operations []Operation) string {
	if len(operations) == 0 {
		return "Unknown conflict"
	}

	if len(operations) == 2 {
		return fmt.Sprintf("Conflict between %s and %s at position %d",
			operations[0].UserID, operations[1].UserID, operations[0].Position)
	}

	var users []string
	for _, op := range operations {
		users = append(users, op.UserID)
	}

	return fmt.Sprintf("Multi-user conflict involving %s at position %d",
		strings.Join(users, ", "), operations[0].Position)
}

// Advanced merge strategies

// ThreeWayMergeStrategy implements three-way merge for conflicts
type ThreeWayMergeStrategy struct{}

func (twm *ThreeWayMergeStrategy) Name() string {
	return "three_way_merge"
}

func (twm *ThreeWayMergeStrategy) Resolve(conflict *Conflict) (*Resolution, error) {
	// Three-way merge: compare operations against a common ancestor
	// For now, implement a simple version that tries to merge non-overlapping changes

	var mergedContent strings.Builder
	changes := make([]DocumentChange, 0)

	// Sort operations by position
	sortedOps := make([]Operation, len(conflict.Operations))
	copy(sortedOps, conflict.Operations)

	// Simple bubble sort by position
	for i := 0; i < len(sortedOps)-1; i++ {
		for j := 0; j < len(sortedOps)-i-1; j++ {
			if sortedOps[j].Position > sortedOps[j+1].Position {
				sortedOps[j], sortedOps[j+1] = sortedOps[j+1], sortedOps[j]
			}
		}
	}

	currentPos := 0
	for _, op := range sortedOps {
		// Skip overlapping operations (would need manual resolution)
		if op.Position < currentPos {
			continue
		}

		// Add non-overlapping operation
		if op.Type == OpInsert {
			changes = append(changes, DocumentChange{
				Type:       "insert",
				Position:   op.Position,
				NewContent: op.Content,
				UserID:     op.UserID,
				Timestamp:  op.Timestamp,
			})
			currentPos = op.Position + len(op.Content)
		}
	}

	resolution := &Resolution{
		Strategy:     "three_way_merge",
		FinalContent: mergedContent.String(),
		Changes:      changes,
		Metadata: map[string]interface{}{
			"merge_type":        "three_way",
			"operations_merged": len(changes),
		},
	}

	return resolution, nil
}

// Utility functions

func generateConflictID() string {
	return fmt.Sprintf("conflict_%d", time.Now().UnixNano())
}

func generateResolutionID() string {
	return fmt.Sprintf("resolution_%d", time.Now().UnixNano())
}
