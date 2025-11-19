package e2e

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
)

// TestTaskLifecycle tests the complete lifecycle of a task from creation to deletion
func TestTaskLifecycle(t *testing.T) {
	// Setup: Create temporary storage
	tmpDir, err := os.MkdirTemp("", "focus-e2e-lifecycle-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	// Step 1: Add a new task
	t.Log("Step 1: Creating a new task")
	newTask := models.ParsedTask{
		Description: "Write integration tests",
		Priority:    "high",
		Categories:  []string{"testing", "development"},
	}

	addedTask, err := eng.AddTask(newTask)
	if err != nil {
		t.Fatalf("Step 1 failed - could not add task: %v", err)
	}

	if addedTask.ID == "" {
		t.Fatal("Step 1 failed - task ID not generated")
	}

	if addedTask.Status != "pending" {
		t.Errorf("Step 1 failed - expected status 'pending', got '%s'", addedTask.Status)
	}

	taskID := addedTask.ID
	t.Logf("Step 1 passed - task created with ID: %s", taskID)

	// Step 2: List tasks and verify the task exists
	t.Log("Step 2: Listing all tasks")
	allTasks, err := eng.ListTasks("")
	if err != nil {
		t.Fatalf("Step 2 failed - could not list tasks: %v", err)
	}

	if len(allTasks) != 1 {
		t.Errorf("Step 2 failed - expected 1 task, got %d", len(allTasks))
	}

	if allTasks[0].ID != taskID {
		t.Error("Step 2 failed - listed task ID doesn't match created task")
	}

	t.Log("Step 2 passed - task found in list")

	// Step 3: Get specific task by ID
	t.Log("Step 3: Retrieving task by ID")
	retrievedTask, err := eng.GetTask(taskID)
	if err != nil {
		t.Fatalf("Step 3 failed - could not get task: %v", err)
	}

	if retrievedTask.Description != newTask.Description {
		t.Errorf("Step 3 failed - description mismatch: expected '%s', got '%s'",
			newTask.Description, retrievedTask.Description)
	}

	if retrievedTask.Priority != newTask.Priority {
		t.Errorf("Step 3 failed - priority mismatch: expected '%s', got '%s'",
			newTask.Priority, retrievedTask.Priority)
	}

	t.Log("Step 3 passed - task retrieved correctly")

	// Step 4: Complete the task
	t.Log("Step 4: Marking task as complete")
	err = eng.CompleteTask(taskID)
	if err != nil {
		t.Fatalf("Step 4 failed - could not complete task: %v", err)
	}

	completedTask, err := eng.GetTask(taskID)
	if err != nil {
		t.Fatalf("Step 4 failed - could not retrieve completed task: %v", err)
	}

	if completedTask.Status != "completed" {
		t.Errorf("Step 4 failed - expected status 'completed', got '%s'", completedTask.Status)
	}

	t.Log("Step 4 passed - task marked as complete")

	// Step 5: Filter completed tasks
	t.Log("Step 5: Filtering completed tasks")
	completedTasks, err := eng.ListTasks("completed")
	if err != nil {
		t.Fatalf("Step 5 failed - could not list completed tasks: %v", err)
	}

	if len(completedTasks) != 1 {
		t.Errorf("Step 5 failed - expected 1 completed task, got %d", len(completedTasks))
	}

	// Verify pending filter doesn't include completed task
	pendingTasks, err := eng.ListTasks("pending")
	if err != nil {
		t.Fatalf("Step 5 failed - could not list pending tasks: %v", err)
	}

	if len(pendingTasks) != 0 {
		t.Errorf("Step 5 failed - expected 0 pending tasks, got %d", len(pendingTasks))
	}

	t.Log("Step 5 passed - task filtering works correctly")

	// Step 6: Delete the task
	t.Log("Step 6: Deleting the task")
	err = eng.DeleteTask(taskID)
	if err != nil {
		t.Fatalf("Step 6 failed - could not delete task: %v", err)
	}

	// Verify task is deleted
	allTasks, err = eng.ListTasks("")
	if err != nil {
		t.Fatalf("Step 6 failed - could not list tasks after deletion: %v", err)
	}

	if len(allTasks) != 0 {
		t.Errorf("Step 6 failed - expected 0 tasks after deletion, got %d", len(allTasks))
	}

	// Verify GetTask returns error
	_, err = eng.GetTask(taskID)
	if err == nil {
		t.Error("Step 6 failed - GetTask should return error for deleted task")
	}

	t.Log("Step 6 passed - task deleted successfully")
	t.Log("✓ Complete task lifecycle test passed")
}

// TestMultiTaskWorkflow tests working with multiple tasks
func TestMultiTaskWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-multi-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	// Step 1: Add multiple tasks with different priorities
	t.Log("Step 1: Creating multiple tasks")
	tasks := []models.ParsedTask{
		{Description: "High priority task", Priority: "high"},
		{Description: "Medium priority task", Priority: "medium"},
		{Description: "Low priority task", Priority: "low"},
	}

	taskIDs := make([]string, 0, 3)
	for _, task := range tasks {
		added, err := eng.AddTask(task)
		if err != nil {
			t.Fatalf("Step 1 failed - could not add task: %v", err)
		}
		taskIDs = append(taskIDs, added.ID)
	}

	t.Logf("Step 1 passed - created %d tasks", len(taskIDs))

	// Step 2: List all tasks
	t.Log("Step 2: Listing all tasks")
	allTasks, err := eng.ListTasks("")
	if err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	if len(allTasks) != 3 {
		t.Errorf("Step 2 failed - expected 3 tasks, got %d", len(allTasks))
	}

	t.Log("Step 2 passed - all tasks listed")

	// Step 3: Complete some tasks
	t.Log("Step 3: Completing first two tasks")
	for i := 0; i < 2; i++ {
		if err := eng.CompleteTask(taskIDs[i]); err != nil {
			t.Fatalf("Step 3 failed - could not complete task %d: %v", i, err)
		}
	}

	// Verify counts
	completedTasks, _ := eng.ListTasks("completed")
	pendingTasks, _ := eng.ListTasks("pending")

	if len(completedTasks) != 2 {
		t.Errorf("Step 3 failed - expected 2 completed tasks, got %d", len(completedTasks))
	}

	if len(pendingTasks) != 1 {
		t.Errorf("Step 3 failed - expected 1 pending task, got %d", len(pendingTasks))
	}

	t.Log("Step 3 passed - tasks completed correctly")

	// Step 4: Delete all tasks
	t.Log("Step 4: Cleaning up all tasks")
	for _, id := range taskIDs {
		eng.DeleteTask(id)
	}

	allTasks, _ = eng.ListTasks("")
	if len(allTasks) != 0 {
		t.Errorf("Step 4 failed - expected 0 tasks after cleanup, got %d", len(allTasks))
	}

	t.Log("Step 4 passed - all tasks deleted")
	t.Log("✓ Multi-task workflow test passed")
}

// TestTaskWithDeadline tests tasks with deadlines
func TestTaskWithDeadline(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-deadline-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	// Step 1: Create task with deadline
	t.Log("Step 1: Creating task with deadline")
	tomorrow := time.Now().Add(24 * time.Hour)

	newTask := models.ParsedTask{
		Description: "Task with deadline",
		Priority:    "high",
		Deadline:    tomorrow,
	}

	addedTask, err := eng.AddTask(newTask)
	if err != nil {
		t.Fatalf("Step 1 failed: %v", err)
	}

	if addedTask.Deadline.IsZero() {
		t.Error("Step 1 failed - deadline not set")
	}

	t.Log("Step 1 passed - task created with deadline")

	// Step 2: Verify deadline persists
	t.Log("Step 2: Verifying deadline persistence")
	retrieved, err := eng.GetTask(addedTask.ID)
	if err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	if retrieved.Deadline.IsZero() {
		t.Error("Step 2 failed - deadline not persisted")
	}

	// Check deadline is roughly equal (within 1 second)
	if retrieved.Deadline.Sub(tomorrow).Abs() > time.Second {
		t.Errorf("Step 2 failed - deadline mismatch: expected %v, got %v",
			tomorrow, retrieved.Deadline)
	}

	t.Log("Step 2 passed - deadline persisted correctly")

	// Cleanup
	eng.DeleteTask(addedTask.ID)
	t.Log("✓ Task with deadline test passed")
}
