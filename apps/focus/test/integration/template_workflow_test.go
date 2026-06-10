package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/internal/templates"
	"github.com/kyanite/focus/pkg/models"
)

// TestTemplateWorkflowIntegration tests the complete template workflow
func TestTemplateWorkflowIntegration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-template-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templatePath := filepath.Join(tmpDir, "templates.json")
	tasksPath := filepath.Join(tmpDir, "tasks.json")

	templateStore := templates.NewStore(templatePath)
	repo := repository.NewStoreRepository(tasksPath)
	eng := engine.New(repo)

	t.Run("create_and_use_template", func(t *testing.T) {
		// Create a template
		template := models.TaskTemplate{
			Name:        "Daily Standup",
			Description: "Attend daily standup meeting",
			Priority:    "high",
			Categories:  []string{"meetings", "work"},
		}

		addedTemplate, err := templateStore.Add(template)
		if err != nil {
			t.Fatalf("Failed to add template: %v", err)
		}

		if addedTemplate.ID == "" {
			t.Error("Template ID should be generated")
		}

		// Load templates to verify persistence
		loadedTemplates, err := templateStore.Load()
		if err != nil {
			t.Fatalf("Failed to load templates: %v", err)
		}

		if len(loadedTemplates) != 1 {
			t.Errorf("Expected 1 template, got %d", len(loadedTemplates))
		}

		// Get template by ID
		retrieved, err := templateStore.Get(addedTemplate.ID)
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		if retrieved.Name != "Daily Standup" {
			t.Errorf("Expected name 'Daily Standup', got %s", retrieved.Name)
		}

		// Convert template to task
		parsedTask := retrieved.ToTask()

		if parsedTask.Description != template.Description {
			t.Errorf("Expected description %s, got %s", template.Description, parsedTask.Description)
		}

		if parsedTask.Priority != template.Priority {
			t.Errorf("Expected priority %s, got %s", template.Priority, parsedTask.Priority)
		}

		// Create task from template
		createdTask, err := eng.AddTask(parsedTask)
		if err != nil {
			t.Fatalf("Failed to create task from template: %v", err)
		}

		if createdTask.ID == "" {
			t.Error("Task ID should be generated")
		}

		// Verify task was created with template properties
		if createdTask.Description != template.Description {
			t.Errorf("Task description doesn't match template")
		}

		// Clean up
		eng.DeleteTask(createdTask.ID)
		templateStore.Delete(addedTemplate.ID)
	})

	t.Run("template_reuse", func(t *testing.T) {
		// Create a template
		template := models.TaskTemplate{
			Name:        "Weekly Review",
			Description: "Review weekly progress",
			Priority:    "medium",
			Categories:  []string{"planning", "review"},
		}

		addedTemplate, err := templateStore.Add(template)
		if err != nil {
			t.Fatalf("Failed to add template: %v", err)
		}

		// Use the template multiple times
		for i := 0; i < 3; i++ {
			parsedTask := addedTemplate.ToTask()
			task, err := eng.AddTask(parsedTask)
			if err != nil {
				t.Fatalf("Failed to create task %d from template: %v", i, err)
			}

			if task.Description != template.Description {
				t.Errorf("Task %d description doesn't match template", i)
			}
		}

		// Verify 3 tasks were created
		tasks, err := eng.ListTasks("")
		if err != nil {
			t.Fatalf("Failed to list tasks: %v", err)
		}

		if len(tasks) != 3 {
			t.Errorf("Expected 3 tasks from template, got %d", len(tasks))
		}

		// Clean up
		for _, task := range tasks {
			eng.DeleteTask(task.ID)
		}
		templateStore.Delete(addedTemplate.ID)
	})

	t.Run("get_template_by_name", func(t *testing.T) {
		// Create templates
		template1 := models.TaskTemplate{
			Name:        "Code Review",
			Description: "Review pull requests",
			Priority:    "high",
		}
		template2 := models.TaskTemplate{
			Name:        "Documentation",
			Description: "Update documentation",
			Priority:    "low",
		}

		t1, err := templateStore.Add(template1)
		if err != nil {
			t.Fatalf("Failed to add template1: %v", err)
		}

		t2, err := templateStore.Add(template2)
		if err != nil {
			t.Fatalf("Failed to add template2: %v", err)
		}

		// Get template by name
		retrieved, err := templateStore.GetByName("Code Review")
		if err != nil {
			t.Fatalf("Failed to get template by name: %v", err)
		}

		if retrieved.ID != t1.ID {
			t.Errorf("Expected template ID %s, got %s", t1.ID, retrieved.ID)
		}

		// Clean up
		templateStore.Delete(t1.ID)
		templateStore.Delete(t2.ID)
	})
}
