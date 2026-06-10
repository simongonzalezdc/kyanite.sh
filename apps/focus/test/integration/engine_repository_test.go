package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
)

// TestEngineRepositoryIntegration tests that engine and repository work together correctly
func TestEngineRepositoryIntegration(t *testing.T) {
	t.Run("complete_task_lifecycle", func(t *testing.T) {
		// Create a temporary directory for test storage
		tmpDir, err := os.MkdirTemp("", "focus-integration-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")

		// Initialize repository and engine
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		// Add a task
		task1 := models.ParsedTask{
			Description: "Integration test task",
			Priority:    "high",
			Categories:  []string{"test", "integration"},
		}

		addedTask, err := eng.AddTask(task1)
		if err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		if addedTask.ID == "" {
			t.Error("Task ID should be generated")
		}

		// List tasks to verify persistence
		tasks, err := eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 1 {
			t.Errorf("Expected 1 task, got %d", len(tasks))
		}

		// Complete the task
		err = eng.CompleteTask(addedTask.ID)
		if err != nil {
			t.Fatalf("Failed to complete task: %v", err)
		}

		// Verify task is completed
		completedTask, err := eng.GetTask(addedTask.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if completedTask.Status != "completed" {
			t.Errorf("Expected status completed, got %s", completedTask.Status)
		}

		// Delete the task
		err = eng.DeleteTask(addedTask.ID)
		if err != nil {
			t.Fatalf("Failed to delete task: %v", err)
		}

		// Verify task is deleted
		tasks, err = eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
		}
	})

	t.Run("persistence_across_instances", func(t *testing.T) {
		// Create a temporary directory for test storage
		tmpDir, err := os.MkdirTemp("", "focus-integration-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")

		// Initialize repository and engine
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)

		// Add tasks with first engine instance
		task1 := models.ParsedTask{
			Description: "Persistent task 1",
			Priority:    "medium",
		}
		task2 := models.ParsedTask{
			Description: "Persistent task 2",
			Priority:    "low",
		}

		_, err = eng.AddTask(task1)
		if err != nil {
			t.Fatalf("Failed to add task1: %v", err)
		}

		_, err = eng.AddTask(task2)
		if err != nil {
			t.Fatalf("Failed to add task2: %v", err)
		}

		// Create a new engine instance with same repository
		repo2 := repository.NewStoreRepository(storePath)
		eng2 := engine.New(repo2)

		// Verify tasks are persisted
		tasks, err := eng2.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 2 {
			t.Errorf("Expected 2 persisted tasks, got %d", len(tasks))
		}

		// Clean up
		for _, task := range tasks {
			eng2.DeleteTask(task.ID)
		}
	})

	t.Run("concurrent_operations", func(t *testing.T) {
		// Create a temporary directory for test storage
		tmpDir, err := os.MkdirTemp("", "focus-integration-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")

		// Initialize repository and engine
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)

		// Add multiple tasks concurrently
		done := make(chan bool)
		for i := 0; i < 5; i++ {
			go func(idx int) {
				task := models.ParsedTask{
					Description: "Concurrent task",
					Priority:    "medium",
				}
				_, err := eng.AddTask(task)
				if err != nil {
					t.Errorf("Failed to add concurrent task: %v", err)
				}
				done <- true
			}(i)
		}

		// Wait for all goroutines to complete
		for i := 0; i < 5; i++ {
			<-done
		}

		// Verify all tasks were added
		tasks, err := eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 5 {
			t.Errorf("Expected 5 concurrent tasks, got %d", len(tasks))
		}

		// Clean up
		for _, task := range tasks {
			eng.DeleteTask(task.ID)
		}
	})
}

// TestRecurringTaskIntegration tests recurring task generation with persistence
func TestRecurringTaskIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-recurring-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	t.Run("daily_recurrence_with_persistence", func(t *testing.T) {
		// Create a recurring task
		now := time.Now()
		endDate := now.AddDate(0, 0, 7) // 7 days from now

		recurringTask := models.ParsedTask{
			Description: "Daily standup",
			Priority:    "high",
			Categories:  []string{"work", "meetings"},
		}

		addedTask, err := eng.AddTask(recurringTask)
		if err != nil {
			t.Fatalf("Failed to add recurring task: %v", err)
		}

		// Manually set recurrence fields (in real usage, this would be done via CLI)
		addedTask.RecurrencePattern = models.RecurrenceDaily
		addedTask.RecurrenceInterval = 1
		addedTask.RecurrenceEndDate = endDate

		// Update the task through engine (so cache is updated)
		err = eng.UpdateTask(addedTask)
		if err != nil {
			t.Fatalf("Failed to update task with recurrence: %v", err)
		}

		// Verify the task has recurrence set
		retrieved, err := eng.GetTask(addedTask.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if !retrieved.IsRecurring() {
			t.Error("Task should be recurring")
		}

		if retrieved.RecurrencePattern != models.RecurrenceDaily {
			t.Errorf("Expected daily recurrence, got %s", retrieved.RecurrencePattern)
		}

		// Clean up
		eng.DeleteTask(addedTask.ID)
	})
}

// TestSubtaskHierarchyIntegration tests parent-child task relationships
func TestSubtaskHierarchyIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-subtask-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	storePath := filepath.Join(tmpDir, "tasks.json")
	repo := repository.NewStoreRepository(storePath)
	eng := engine.New(repo)

	t.Run("parent_child_relationship", func(t *testing.T) {
		// Create parent task
		parentTask := models.ParsedTask{
			Description: "Build new feature",
			Priority:    "high",
		}

		parent, err := eng.AddTask(parentTask)
		if err != nil {
			t.Fatalf("Failed to add parent task: %v", err)
		}

		// Create subtasks
		sub1ID, err := eng.AddSubtask(parent.ID, "Design API", "high", []string{}, "")
		if err != nil {
			t.Fatalf("Failed to add subtask1: %v", err)
		}

		sub2ID, err := eng.AddSubtask(parent.ID, "Implement backend", "medium", []string{}, "")
		if err != nil {
			t.Fatalf("Failed to add subtask2: %v", err)
		}

		// Verify parent has subtasks
		updatedParent, err := eng.GetTask(parent.ID)
		if err != nil {
			t.Fatalf("Failed to get parent task: %v", err)
		}

		if !updatedParent.HasSubtasks() {
			t.Error("Parent should have subtasks")
		}

		if len(updatedParent.SubtaskIDs) != 2 {
			t.Errorf("Expected 2 subtasks, got %d", len(updatedParent.SubtaskIDs))
		}

		// Verify subtasks have parent ID
		retrievedSub1, err := eng.GetTask(sub1ID)
		if err != nil {
			t.Fatalf("Failed to get subtask1: %v", err)
		}

		if !retrievedSub1.IsSubtask() {
			t.Error("Subtask1 should be marked as subtask")
		}

		if retrievedSub1.ParentID != parent.ID {
			t.Errorf("Expected parent ID %s, got %s", parent.ID, retrievedSub1.ParentID)
		}

		// Clean up
		eng.DeleteTask(sub1ID)
		eng.DeleteTask(sub2ID)
		eng.DeleteTask(parent.ID)
	})
}
