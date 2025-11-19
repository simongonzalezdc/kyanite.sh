package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyanite/focus/internal/command"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
)

// TestCommandHistoryIntegration tests undo/redo with actual task operations
func TestCommandHistoryIntegration(t *testing.T) {
	t.Run("add_task_undo_redo", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "focus-command-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		history := command.NewHistory()
		parsedTask := models.ParsedTask{
			Description: "Test task for undo/redo",
			Priority:    "high",
		}

		// Execute add command
		addCmd := command.NewAddTaskCommand(eng, parsedTask)
		err = history.Execute(addCmd)
		if err != nil {
			t.Fatalf("Failed to execute add command: %v", err)
		}

		// Verify task was added
		tasks, err := eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 1 {
			t.Errorf("Expected 1 task after add, got %d", len(tasks))
		}

		taskID := tasks[0].ID

		// Undo the add
		err = history.Undo()
		if err != nil {
			t.Fatalf("Failed to undo add: %v", err)
		}

		// Verify task was removed
		tasks, err = eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks after undo, got %d", len(tasks))
		}

		// Redo the add
		err = history.Redo()
		if err != nil {
			t.Fatalf("Failed to redo add: %v", err)
		}

		// Verify task was re-added
		tasks, err = eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 1 {
			t.Errorf("Expected 1 task after redo, got %d", len(tasks))
		}

		if tasks[0].ID != taskID {
			t.Errorf("Task ID changed after redo: expected %s, got %s", taskID, tasks[0].ID)
		}

		// Clean up
		eng.DeleteTask(taskID)
	})

	t.Run("complete_task_undo_redo", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "focus-command-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		history := command.NewHistory()

		// Add a task first
		parsedTask := models.ParsedTask{
			Description: "Task to complete",
			Priority:    "medium",
		}

		task, err := eng.AddTask(parsedTask)
		if err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Execute complete command
		completeCmd := command.NewCompleteTaskCommand(eng, task.ID)
		err = history.Execute(completeCmd)
		if err != nil {
			t.Fatalf("Failed to execute complete command: %v", err)
		}

		// Verify task is completed
		completedTask, err := eng.GetTask(task.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if completedTask.Status != "completed" {
			t.Errorf("Expected status completed, got %s", completedTask.Status)
		}

		// Undo the complete
		err = history.Undo()
		if err != nil {
			t.Fatalf("Failed to undo complete: %v", err)
		}

		// Verify task is pending again
		restoredTask, err := eng.GetTask(task.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if restoredTask.Status != "pending" {
			t.Errorf("Expected status pending after undo, got %s", restoredTask.Status)
		}

		// Clean up
		eng.DeleteTask(task.ID)
	})

	t.Run("delete_task_undo_redo", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "focus-command-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		history := command.NewHistory()

		// Add a task first
		parsedTask := models.ParsedTask{
			Description: "Task to delete",
			Priority:    "low",
			Categories:  []string{"test"},
		}

		task, err := eng.AddTask(parsedTask)
		if err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		taskID := task.ID

		// Execute delete command
		deleteCmd := command.NewDeleteTaskCommand(eng, taskID)
		err = history.Execute(deleteCmd)
		if err != nil {
			t.Fatalf("Failed to execute delete command: %v", err)
		}

		// Verify task is deleted
		tasks, err := eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks after delete, got %d", len(tasks))
		}

		// Undo the delete
		err = history.Undo()
		if err != nil {
			t.Fatalf("Failed to undo delete: %v", err)
		}

		// Verify task is restored
		restoredTask, err := eng.GetTask(taskID)
		if err != nil {
			t.Fatalf("Failed to get restored task: %v", err)
		}

		if restoredTask.Description != parsedTask.Description {
			t.Errorf("Restored task description doesn't match original")
		}

		// Clean up
		eng.DeleteTask(taskID)
	})

	t.Run("multiple_operations_workflow", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "focus-command-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		history := command.NewHistory()

		// Add three tasks
		task1 := models.ParsedTask{Description: "Task 1", Priority: "high"}
		task2 := models.ParsedTask{Description: "Task 2", Priority: "medium"}
		task3 := models.ParsedTask{Description: "Task 3", Priority: "low"}

		cmd1 := command.NewAddTaskCommand(eng, task1)
		cmd2 := command.NewAddTaskCommand(eng, task2)
		cmd3 := command.NewAddTaskCommand(eng, task3)

		history.Execute(cmd1)
		history.Execute(cmd2)
		history.Execute(cmd3)

		// Verify 3 tasks
		tasks, _ := eng.ListTasks("")
		if len(tasks) != 3 {
			t.Errorf("Expected 3 tasks, got %d", len(tasks))
		}

		// Undo twice
		history.Undo()
		history.Undo()

		// Verify 1 task
		tasks, _ = eng.ListTasks("")
		if len(tasks) != 1 {
			t.Errorf("Expected 1 task after 2 undos, got %d", len(tasks))
		}

		// Redo once
		history.Redo()

		// Verify 2 tasks
		tasks, _ = eng.ListTasks("")
		if len(tasks) != 2 {
			t.Errorf("Expected 2 tasks after redo, got %d", len(tasks))
		}

		// Clean up
		for _, task := range tasks {
			eng.DeleteTask(task.ID)
		}
	})

	t.Run("update_task_undo_redo", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "focus-command-test-*")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		storePath := filepath.Join(tmpDir, "tasks.json")
		repo := repository.NewStoreRepository(storePath)
		eng := engine.New(repo)
		history := command.NewHistory()

		// Add a task
		originalTask := models.ParsedTask{
			Description: "Original description",
			Priority:    "low",
		}

		task, err := eng.AddTask(originalTask)
		if err != nil {
			t.Fatalf("Failed to add task: %v", err)
		}

		// Update the task - need to create a full Task object
		updatedTask := task
		updatedTask.Description = "Updated description"
		updatedTask.Priority = "high"
		updatedTask.Categories = []string{"updated"}

		updateCmd := command.NewUpdateTaskCommand(eng, updatedTask)
		err = history.Execute(updateCmd)
		if err != nil {
			t.Fatalf("Failed to execute update command: %v", err)
		}

		// Verify task was updated
		retrieved, err := eng.GetTask(task.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if retrieved.Description != "Updated description" {
			t.Errorf("Task was not updated")
		}

		// Undo the update
		err = history.Undo()
		if err != nil {
			t.Fatalf("Failed to undo update: %v", err)
		}

		// Verify task is restored to original
		restored, err := eng.GetTask(task.ID)
		if err != nil {
			t.Fatalf("Failed to get task: %v", err)
		}

		if restored.Description != "Original description" {
			t.Errorf("Expected original description after undo, got %s", restored.Description)
		}

		// Clean up
		eng.DeleteTask(task.ID)
	})
}
