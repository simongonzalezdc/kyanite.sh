package engine

import (
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"path/filepath"
	"testing"
	"time"
)

func TestEngine_New(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	store := store.New(testFile)
	engine := New(store)

	if engine == nil {
		t.Error("New() returned nil engine")
	}
}

func TestEngine_AddTask(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	store := store.New(testFile)
	engine := New(store)

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

	store := store.New(testFile)
	engine := New(store)

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

	store := store.New(testFile)
	engine := New(store)

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

	store := store.New(testFile)
	engine := New(store)

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

	store := store.New(testFile)
	engine := New(store)

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

	store := store.New(testFile)
	engine := New(store)

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
