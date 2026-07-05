package engine

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
)

func TestEngine_New(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	if engine == nil {
		t.Error("New() returned nil engine")
	}
}

func TestEngine_AddTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	tests := []struct {
		name    string
		task    models.ParsedTask
		wantErr bool
	}{
		{
			name: "valid task",
			task: models.ParsedTask{
				Description: "Test task",
				Priority:    "medium",
			},
			wantErr: false,
		},
		{
			name: "empty description",
			task: models.ParsedTask{
				Description: "",
				Priority:    "high",
			},
			wantErr: true,
		},
		{
			name: "task with deadline",
			task: models.ParsedTask{
				Description: "Test task with deadline",
				Priority:    "high",
				Deadline:    time.Now().Add(24 * time.Hour),
			},
			wantErr: false,
		},
		{
			name: "task with categories",
			task: models.ParsedTask{
				Description: "Test task with categories",
				Priority:    "low",
				Categories:  []string{"work", "personal"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := engine.AddTask(tt.task)

			if (err != nil) != tt.wantErr {
				t.Errorf("Engine.AddTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if task.ID == "" {
					t.Error("Expected task to have an ID")
				}

				if task.Description != tt.task.Description {
					t.Errorf("Expected description %s, got %s", tt.task.Description, task.Description)
				}

				if task.Status != "pending" {
					t.Errorf("Expected status 'pending', got %s", task.Status)
				}

				if task.Priority != tt.task.Priority {
					t.Errorf("Expected priority %s, got %s", tt.task.Priority, task.Priority)
				}

				// Check for default priority
				if tt.task.Priority == "" && task.Priority != "medium" {
					t.Errorf("Expected default priority 'medium', got %s", task.Priority)
				}
			}
		})
	}
}

func TestEngine_ListTasks(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Add some test tasks
	task1 := models.ParsedTask{
		Description: "Task 1",
		Priority:    "high",
	}
	task2 := models.ParsedTask{
		Description: "Task 2",
		Priority:    "low",
	}

	engine.AddTask(task1)
	engine.AddTask(task2)

	// Complete one task for testing
	tasks, _ := engine.ListTasks("all")
	if len(tasks) > 0 {
		engine.CompleteTask(tasks[0].ID)
	}

	tests := []struct {
		name      string
		filter    string
		wantCount int
	}{
		{
			name:      "all tasks",
			filter:    "all",
			wantCount: 2,
		},
		{
			name:      "pending tasks",
			filter:    "active",
			wantCount: 1,
		},
		{
			name:      "completed tasks",
			filter:    "completed",
			wantCount: 1,
		},
		{
			name:      "no filter",
			filter:    "",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks, err := engine.ListTasks(tt.filter)
			if err != nil {
				t.Fatalf("ListTasks() error = %v", err)
			}

			if len(tasks) != tt.wantCount {
				t.Errorf("ListTasks() returned %d tasks, want %d", len(tasks), tt.wantCount)
			}
		})
	}
}

func TestEngine_CompleteTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Add a task first
	task := models.ParsedTask{
		Description: "Test task",
		Priority:    "medium",
	}
	addedTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Complete a task
	err = engine.CompleteTask(addedTask.ID)
	if err != nil {
		t.Fatalf("CompleteTask() error = %v", err)
	}

	// Check it was completed
	retrievedTask, err := engine.GetTask(addedTask.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if retrievedTask.Status != "completed" {
		t.Errorf("Expected task to be completed, got status %s", retrievedTask.Status)
	}
}

func TestEngine_DeleteTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Add a task first
	task := models.ParsedTask{
		Description: "Test task to delete",
		Priority:    "medium",
	}
	addedTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Delete the task
	err = engine.DeleteTask(addedTask.ID)
	if err != nil {
		t.Fatalf("DeleteTask() error = %v", err)
	}

	// Try to get the deleted task
	_, err = engine.GetTask(addedTask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task")
	}

	// Verify task count reduced
	tasks, err := engine.ListTasks("all")
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected no tasks after deletion, got %d", len(tasks))
	}
}

func TestEngine_GetTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Add a task first
	task := models.ParsedTask{
		Description: "Test get task",
		Priority:    "high",
	}
	addedTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Get the task
	retrievedTask, err := engine.GetTask(addedTask.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if retrievedTask.ID != addedTask.ID {
		t.Errorf("Expected task ID %s, got %s", addedTask.ID, retrievedTask.ID)
	}

	if retrievedTask.Description != addedTask.Description {
		t.Errorf("Expected task description %s, got %s", addedTask.Description, retrievedTask.Description)
	}
}

func TestEngine_updateTaskStatus(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Add a task first
	task := models.ParsedTask{
		Description: "Test status update",
		Priority:    "medium",
	}
	addedTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Update to completed
	err = engine.updateTaskStatus(addedTask.ID, "completed")
	if err != nil {
		t.Fatalf("updateTaskStatus() error = %v", err)
	}

	// Check if status was updated
	retrievedTask, err := engine.GetTask(addedTask.ID)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}

	if retrievedTask.Status != "completed" {
		t.Errorf("Expected status 'completed', got %s", retrievedTask.Status)
	}
}

func TestEngine_generateID(t *testing.T) {
	// This is a simple function, but we should at least verify it doesn't panic
	id1 := generateID()
	id2 := generateID()

	if id1 == "" {
		t.Error("generateID() returned empty string")
	}

	// Test that generated IDs are unique (or at least not identical in short time)
	if id1 == id2 {
		// This might fail in very rare cases, but unlikely
		t.Logf("Generated IDs are identical: %s", id1)
	}
}

// TestEngine_AddSubtask tests adding subtasks to parent tasks
func TestEngine_AddSubtask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Create a parent task first
	parentTask := models.ParsedTask{
		Description: "Parent task",
		Priority:    "high",
	}
	addedParent, err := engine.AddTask(parentTask)
	if err != nil {
		t.Fatalf("Failed to add parent task: %v", err)
	}

	tests := []struct {
		name          string
		parentID      string
		description   string
		priority      string
		categories    []string
		deadline      string
		wantErr       bool
		errorContains string
	}{
		{
			name:        "valid subtask",
			parentID:    addedParent.ID,
			description: "Subtask 1",
			priority:    "medium",
			categories:  []string{"work"},
			wantErr:     false,
		},
		{
			name:        "empty description",
			parentID:    addedParent.ID,
			description: "",
			priority:    "low",
			wantErr:     true,
			errorContains: "subtask description cannot be empty",
		},
		{
			name:          "non-existent parent",
			parentID:      "non-existent-id",
			description:   "Orphan subtask",
			priority:      "medium",
			wantErr:       true,
			errorContains: "parent task not found",
		},
		{
			name:        "subtask with deadline",
			parentID:    addedParent.ID,
			description: "Subtask with deadline",
			priority:    "high",
			deadline:    "2024-12-25",
			wantErr:     false,
		},
		{
			name:        "inherit parent priority",
			parentID:    addedParent.ID,
			description: "Subtask without priority",
			priority:    "",
			wantErr:     false,
		},
		{
			name:        "invalid priority normalization",
			parentID:    addedParent.ID,
			description: "Subtask with bad priority",
			priority:    "invalid",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subtaskID, err := engine.AddSubtask(tt.parentID, tt.description, tt.priority, tt.categories, tt.deadline)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddSubtask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errorContains != "" && err == nil {
					t.Errorf("Expected error containing %s, got nil", tt.errorContains)
				}
				if tt.errorContains != "" && err != nil {
					if !contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error containing %s, got %s", tt.errorContains, err.Error())
					}
				}
				return
			}

			// Verify subtask was created
			if subtaskID == "" {
				t.Error("Expected non-empty subtask ID")
			}

			// Verify subtask can be retrieved
			subtask, err := engine.GetTask(subtaskID)
			if err != nil {
				t.Errorf("Failed to get created subtask: %v", err)
			}

			if subtask.Description != tt.description {
				t.Errorf("Expected description %s, got %s", tt.description, subtask.Description)
			}

			// Check priority inheritance or normalization
			if tt.priority == "" && subtask.Priority != addedParent.Priority {
				t.Errorf("Expected inherited priority %s, got %s", addedParent.Priority, subtask.Priority)
			}
			if tt.priority == "invalid" && subtask.Priority != "medium" {
				t.Errorf("Expected normalized priority 'medium', got %s", subtask.Priority)
			}

			// Verify parent was updated
			parent, err := engine.GetTask(tt.parentID)
			if err != nil {
				t.Errorf("Failed to get parent task: %v", err)
			}

			if !parent.HasSubtasks() {
				t.Error("Expected parent to have subtasks")
			}

			// Check if this subtask ID is in parent's subtask list
			found := false
			for _, sid := range parent.SubtaskIDs {
				if sid == subtaskID {
					found = true
					break
				}
			}
			if !found {
				t.Error("Expected parent to contain this subtask ID")
			}
		})
	}
}

// TestEngine_RestoreTask tests restoring deleted tasks
func TestEngine_RestoreTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Create a task
	task := models.ParsedTask{
		Description: "Task to restore",
		Priority:    "high",
	}
	originalTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Delete the task
	err = engine.DeleteTask(originalTask.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify task is gone
	_, err = engine.GetTask(originalTask.ID)
	if err == nil {
		t.Error("Expected task to be deleted")
	}

	// Restore the task
	err = engine.RestoreTask(originalTask)
	if err != nil {
		t.Errorf("RestoreTask() error = %v", err)
	}

	// Verify task is restored
restoredTask, err := engine.GetTask(originalTask.ID)
	if err != nil {
		t.Errorf("Failed to get restored task: %v", err)
	}

	if restoredTask.ID != originalTask.ID {
		t.Errorf("Expected ID %s, got %s", originalTask.ID, restoredTask.ID)
	}

	if restoredTask.Description != originalTask.Description {
		t.Errorf("Expected description %s, got %s", originalTask.Description, restoredTask.Description)
	}

	// Test restore of already existing task
	err = engine.RestoreTask(originalTask)
	if err == nil {
		t.Error("Expected error when restoring already existing task")
	}
}

// TestEngine_UpdateTask tests updating task fields
func TestEngine_UpdateTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	repo := repository.NewStoreRepository(testFile)
	engine := New(repo)

	// Create a task
	task := models.ParsedTask{
		Description: "Original description",
		Priority:    "low",
	}
	originalTask, err := engine.AddTask(task)
	if err != nil {
		t.Fatalf("Failed to add test task: %v", err)
	}

	// Wait a moment to ensure UpdatedAt will be different
	time.Sleep(10 * time.Millisecond)

	tests := []struct {
		name          string
		taskID        string
		updateFunc    func(models.Task) models.Task
		wantErr       bool
		errorContains string
		verifyFunc    func(models.Task, models.Task) error
	}{
		{
			name: "update description",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Description = "Updated description"
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if updated.Description != "Updated description" {
					return fmt.Errorf("expected updated description, got %s", updated.Description)
				}
				return nil
			},
		},
		{
			name: "update priority",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Priority = "high"
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if updated.Priority != "high" {
					return fmt.Errorf("expected priority 'high', got %s", updated.Priority)
				}
				return nil
			},
		},
		{
			name: "update status",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Status = "completed"
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if updated.Status != "completed" {
					return fmt.Errorf("expected status 'completed', got %s", updated.Status)
				}
				return nil
			},
		},
		{
			name: "update categories",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Categories = []string{"work", "urgent"}
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if len(updated.Categories) != 2 {
					return fmt.Errorf("expected 2 categories, got %d", len(updated.Categories))
				}
				return nil
			},
		},
		{
			name: "preserve created at",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Description = "Updated with preserved created at"
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if !updated.CreatedAt.Equal(original.CreatedAt) {
					return fmt.Errorf("CreatedAt should be preserved, was %s, now %s", original.CreatedAt, updated.CreatedAt)
				}
				if !updated.UpdatedAt.After(original.UpdatedAt) {
					return fmt.Errorf("UpdatedAt should be more recent, was %s, now %s", original.UpdatedAt, updated.UpdatedAt)
				}
				return nil
			},
		},
		{
			name: "update non-existent task",
			taskID: "non-existent-id",
			updateFunc: func(t models.Task) models.Task {
				t.Description = "This should fail"
				return t
			},
			wantErr:       true,
			errorContains: "not found",
		},
		{
			name: "update with deadline",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Deadline = time.Now().Add(24 * time.Hour)
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if updated.Deadline.IsZero() {
					return fmt.Errorf("expected deadline to be set")
				}
				return nil
			},
		},
		{
			name: "update notes",
			taskID: originalTask.ID,
			updateFunc: func(t models.Task) models.Task {
				t.Notes = "Updated notes content"
				return t
			},
			wantErr: false,
			verifyFunc: func(updated models.Task, original models.Task) error {
				if updated.Notes != "Updated notes content" {
					return fmt.Errorf("expected notes to be updated")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get current task state
			currentTask, err := engine.GetTask(tt.taskID)
			if err != nil && tt.taskID != "non-existent-id" {
				t.Fatalf("Failed to get current task: %v", err)
			}

			// Apply update function
			updatedTask := tt.updateFunc(currentTask)

			// Perform update
			err = engine.UpdateTask(updatedTask)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if tt.errorContains != "" && err == nil {
					t.Errorf("Expected error containing %s, got nil", tt.errorContains)
				}
				if tt.errorContains != "" && err != nil {
					if !contains(err.Error(), tt.errorContains) {
						t.Errorf("Expected error containing %s, got %s", tt.errorContains, err.Error())
					}
				}
				return
			}

			// Verify update was successful
			retrievedTask, err := engine.GetTask(tt.taskID)
			if err != nil {
				t.Errorf("Failed to get updated task: %v", err)
			}

			// Run custom verification
			if tt.verifyFunc != nil {
				if err := tt.verifyFunc(retrievedTask, originalTask); err != nil {
					t.Error(err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
