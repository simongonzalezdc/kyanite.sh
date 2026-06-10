package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
)

// TestSubtaskCompleteWorkflow tests the complete subtask workflow
func TestSubtaskCompleteWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-subtask-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	// Step 1: Create a parent task
	t.Log("Step 1: Creating parent task")
	parentTask := models.ParsedTask{
		Description: "Build new feature",
		Priority:    "high",
		Categories:  []string{"development", "feature"},
	}

	parent, err := eng.AddTask(parentTask)
	if err != nil {
		t.Fatalf("Step 1 failed - could not create parent task: %v", err)
	}

	if parent.HasSubtasks() {
		t.Error("Step 1 failed - new task should not have subtasks")
	}

	parentID := parent.ID
	t.Logf("Step 1 passed - parent task created with ID: %s", parentID)

	// Step 2: Add subtasks to parent
	t.Log("Step 2: Adding subtasks to parent")
	subtaskDefinitions := []struct {
		description string
		priority    string
	}{
		{"Design API endpoints", "high"},
		{"Implement backend logic", "high"},
		{"Write unit tests", "medium"},
		{"Update documentation", "low"},
	}

	subtaskIDs := make([]string, 0, len(subtaskDefinitions))
	for i, def := range subtaskDefinitions {
		subtaskID, err := eng.AddSubtask(parentID, def.description, def.priority, []string{}, "")
		if err != nil {
			t.Fatalf("Step 2 failed - could not add subtask %d: %v", i+1, err)
		}
		subtaskIDs = append(subtaskIDs, subtaskID)
	}

	t.Logf("Step 2 passed - added %d subtasks", len(subtaskIDs))

	// Step 3: Verify parent has subtasks
	t.Log("Step 3: Verifying parent-child relationships")
	updatedParent, err := eng.GetTask(parentID)
	if err != nil {
		t.Fatalf("Step 3 failed - could not get parent task: %v", err)
	}

	if !updatedParent.HasSubtasks() {
		t.Error("Step 3 failed - parent should have subtasks")
	}

	if len(updatedParent.SubtaskIDs) != len(subtaskDefinitions) {
		t.Errorf("Step 3 failed - expected %d subtask IDs, got %d",
			len(subtaskDefinitions), len(updatedParent.SubtaskIDs))
	}

	t.Log("Step 3 passed - parent has correct subtask references")

	// Step 4: Verify each subtask has parent reference
	t.Log("Step 4: Verifying subtask parent references")
	for i, subtaskID := range subtaskIDs {
		subtask, err := eng.GetTask(subtaskID)
		if err != nil {
			t.Errorf("Step 4 failed - could not get subtask %d: %v", i+1, err)
			continue
		}

		if !subtask.IsSubtask() {
			t.Errorf("Step 4 failed - subtask %d should be marked as subtask", i+1)
		}

		if subtask.ParentID != parentID {
			t.Errorf("Step 4 failed - subtask %d has wrong parent ID: expected %s, got %s",
				i+1, parentID, subtask.ParentID)
		}

		if subtask.Description != subtaskDefinitions[i].description {
			t.Errorf("Step 4 failed - subtask %d description mismatch", i+1)
		}
	}

	t.Log("Step 4 passed - all subtasks have correct parent references")

	// Step 5: List all tasks
	t.Log("Step 5: Listing all tasks")
	allTasks, err := eng.ListTasks("")
	if err != nil {
		t.Fatalf("Step 5 failed - could not list tasks: %v", err)
	}

	expectedTotal := 1 + len(subtaskDefinitions) // parent + subtasks
	if len(allTasks) != expectedTotal {
		t.Errorf("Step 5 failed - expected %d total tasks, got %d", expectedTotal, len(allTasks))
	}

	t.Log("Step 5 passed - all tasks listed correctly")

	// Step 6: Complete subtasks
	t.Log("Step 6: Completing subtasks")
	for i := 0; i < 2; i++ {
		if err := eng.CompleteTask(subtaskIDs[i]); err != nil {
			t.Fatalf("Step 6 failed - could not complete subtask %d: %v", i+1, err)
		}
	}

	completedSubtasks := 0
	for _, id := range subtaskIDs {
		task, _ := eng.GetTask(id)
		if task.Status == "completed" {
			completedSubtasks++
		}
	}

	if completedSubtasks != 2 {
		t.Errorf("Step 6 failed - expected 2 completed subtasks, got %d", completedSubtasks)
	}

	t.Log("Step 6 passed - subtasks completed")

	// Step 7: Complete parent task
	t.Log("Step 7: Completing parent task")
	if err := eng.CompleteTask(parentID); err != nil {
		t.Fatalf("Step 7 failed - could not complete parent: %v", err)
	}

	completedParent, _ := eng.GetTask(parentID)
	if completedParent.Status != "completed" {
		t.Error("Step 7 failed - parent should be completed")
	}

	// Verify subtasks still exist and maintain parent reference
	for _, id := range subtaskIDs {
		subtask, err := eng.GetTask(id)
		if err != nil {
			t.Errorf("Step 7 failed - subtask should still exist after parent completion")
			continue
		}

		if subtask.ParentID != parentID {
			t.Error("Step 7 failed - subtask parent reference should be preserved")
		}
	}

	t.Log("Step 7 passed - parent completed, hierarchy preserved")

	// Step 8: Delete parent and verify cascading behavior
	t.Log("Step 8: Testing deletion behavior")
	// Note: Current implementation might not cascade delete
	// This tests the actual behavior

	if err := eng.DeleteTask(parentID); err != nil {
		t.Fatalf("Step 8 failed - could not delete parent: %v", err)
	}

	// Verify parent is deleted
	_, err = eng.GetTask(parentID)
	if err == nil {
		t.Error("Step 8 failed - parent should be deleted")
	}

	// Check if subtasks still exist (documenting actual behavior)
	remainingSubtasks := 0
	for _, id := range subtaskIDs {
		_, err := eng.GetTask(id)
		if err == nil {
			remainingSubtasks++
		}
	}

	t.Logf("Step 8 - After parent deletion, %d subtasks remain", remainingSubtasks)

	// Cleanup remaining subtasks
	for _, id := range subtaskIDs {
		eng.DeleteTask(id)
	}

	t.Log("Step 8 passed - deletion behavior verified")
	t.Log("✓ Complete subtask workflow test passed")
}

// TestNestedSubtasks tests if subtasks can have their own subtasks
func TestNestedSubtasks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-nested-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	// Step 1: Create grandparent task
	t.Log("Step 1: Creating grandparent task")
	grandparent := models.ParsedTask{
		Description: "Major project milestone",
		Priority:    "high",
	}

	gp, err := eng.AddTask(grandparent)
	if err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}

	t.Log("Step 1 passed - grandparent created")

	// Step 2: Add child task
	t.Log("Step 2: Adding child task")
	childID, err := eng.AddSubtask(gp.ID, "Feature implementation", "high", []string{}, "")
	if err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	t.Log("Step 2 passed - child task added")

	// Step 3: Try to add grandchild (subtask of subtask)
	t.Log("Step 3: Attempting to add grandchild task")
	grandchildID, err := eng.AddSubtask(childID, "Detailed implementation step", "medium", []string{}, "")
	if err != nil {
		t.Logf("Step 3 - Cannot create nested subtasks: %v", err)
	} else {
		t.Log("Step 3 - Nested subtasks are supported")

		// Verify grandchild has correct parent
		grandchild, _ := eng.GetTask(grandchildID)
		if grandchild.ParentID != childID {
			t.Errorf("Grandchild has wrong parent ID: expected %s, got %s", childID, grandchild.ParentID)
		}

		// Cleanup
		eng.DeleteTask(grandchildID)
	}

	// Cleanup
	eng.DeleteTask(childID)
	eng.DeleteTask(gp.ID)

	t.Log("✓ Nested subtasks test passed")
}
