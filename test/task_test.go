package main

import (
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/models"
	"testing"
)

func TestTaskManagement(t *testing.T) {
	// Create a temporary store for testing
	store := store.New("./test_tasks.json")
	engine := engine.New(store)

	// Test adding a task
	parsedTask := models.ParsedTask{
		Description: "Test task",
		Priority:    "medium",
		Categories:  []string{"test"},
	}

	task, err := engine.AddTask(parsedTask)
	if err != nil {
		t.Fatalf("Failed to add task: %v", err)
	}

	if task.Description != "Test task" {
		t.Errorf("Expected description 'Test task', got '%s'", task.Description)
	}

	if task.Priority != "medium" {
		t.Errorf("Expected priority 'medium', got '%s'", task.Priority)
	}

	// Test listing tasks
	tasks, err := engine.ListTasks("all")
	if err != nil {
		t.Fatalf("Failed to list tasks: %v", err)
	}

	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}

	// Test completing a task
	err = engine.CompleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to complete task: %v", err)
	}

	// Verify task is completed
	tasks, _ = engine.ListTasks("completed")
	if len(tasks) != 1 {
		t.Errorf("Expected 1 completed task, got %d", len(tasks))
	}

	// Test deleting a task
	err = engine.DeleteTask(task.ID)
	if err != nil {
		t.Fatalf("Failed to delete task: %v", err)
	}

	// Verify task is deleted
	tasks, _ = engine.ListTasks("all")
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(tasks))
	}
}
