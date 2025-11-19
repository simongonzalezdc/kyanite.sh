package e2e

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/internal/templates"
	"github.com/kyanite/focus/pkg/models"
)

// TestTemplateCompleteWorkflow tests the complete template workflow
func TestTemplateCompleteWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-template-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templatePath := filepath.Join(tmpDir, "templates.json")
	tasksPath := filepath.Join(tmpDir, "tasks.json")

	templateStore := templates.NewStore(templatePath)
	repo := repository.NewStoreRepository(tasksPath)
	eng := engine.New(repo)

	// Step 1: Create a template
	t.Log("Step 1: Creating a task template")
	template := models.TaskTemplate{
		Name:        "Weekly Review",
		Description: "Review progress and plan next week",
		Priority:    "high",
		Categories:  []string{"planning", "review", "weekly"},
	}

	createdTemplate, err := templateStore.Add(template)
	if err != nil {
		t.Fatalf("Step 1 failed - could not create template: %v", err)
	}

	if createdTemplate.ID == "" {
		t.Fatal("Step 1 failed - template ID not generated")
	}

	if createdTemplate.Name != template.Name {
		t.Errorf("Step 1 failed - name mismatch: expected '%s', got '%s'",
			template.Name, createdTemplate.Name)
	}

	templateID := createdTemplate.ID
	t.Logf("Step 1 passed - template created with ID: %s", templateID)

	// Step 2: Load templates and verify it exists
	t.Log("Step 2: Loading templates from storage")
	loadedTemplates, err := templateStore.Load()
	if err != nil {
		t.Fatalf("Step 2 failed - could not load templates: %v", err)
	}

	if len(loadedTemplates) != 1 {
		t.Errorf("Step 2 failed - expected 1 template, got %d", len(loadedTemplates))
	}

	t.Log("Step 2 passed - template found in storage")

	// Step 3: Get template by name
	t.Log("Step 3: Retrieving template by name")
	retrievedByName, err := templateStore.GetByName("Weekly Review")
	if err != nil {
		t.Fatalf("Step 3 failed - could not get template by name: %v", err)
	}

	if retrievedByName.ID != templateID {
		t.Errorf("Step 3 failed - wrong template retrieved: expected ID %s, got %s",
			templateID, retrievedByName.ID)
	}

	t.Log("Step 3 passed - template retrieved by name")

	// Step 4: Use template to create tasks
	t.Log("Step 4: Creating tasks from template")
	tasksToCreate := 3
	createdTaskIDs := make([]string, 0, tasksToCreate)

	for i := 0; i < tasksToCreate; i++ {
		parsedTask := createdTemplate.ToTask()

		task, err := eng.AddTask(parsedTask)
		if err != nil {
			t.Fatalf("Step 4 failed - could not create task %d from template: %v", i+1, err)
		}

		if task.Description != template.Description {
			t.Errorf("Step 4 failed - task %d description doesn't match template", i+1)
		}

		if task.Priority != template.Priority {
			t.Errorf("Step 4 failed - task %d priority doesn't match template", i+1)
		}

		createdTaskIDs = append(createdTaskIDs, task.ID)
	}

	t.Logf("Step 4 passed - created %d tasks from template", tasksToCreate)

	// Step 5: Verify all tasks were created
	t.Log("Step 5: Verifying all tasks exist")
	allTasks, err := eng.ListTasks("")
	if err != nil {
		t.Fatalf("Step 5 failed - could not list tasks: %v", err)
	}

	if len(allTasks) != tasksToCreate {
		t.Errorf("Step 5 failed - expected %d tasks, got %d", tasksToCreate, len(allTasks))
	}

	// Verify each task has correct attributes
	for _, task := range allTasks {
		if task.Description != template.Description {
			t.Errorf("Step 5 failed - task description doesn't match template")
		}

		if task.Priority != template.Priority {
			t.Errorf("Step 5 failed - task priority doesn't match template")
		}

		if len(task.Categories) != len(template.Categories) {
			t.Errorf("Step 5 failed - task categories count mismatch")
		}
	}

	t.Log("Step 5 passed - all tasks verified")

	// Step 6: Complete some tasks
	t.Log("Step 6: Completing tasks")
	for i := 0; i < 2; i++ {
		if err := eng.CompleteTask(createdTaskIDs[i]); err != nil {
			t.Fatalf("Step 6 failed - could not complete task %d: %v", i+1, err)
		}
	}

	completedTasks, _ := eng.ListTasks("completed")
	pendingTasks, _ := eng.ListTasks("pending")

	if len(completedTasks) != 2 {
		t.Errorf("Step 6 failed - expected 2 completed tasks, got %d", len(completedTasks))
	}

	if len(pendingTasks) != 1 {
		t.Errorf("Step 6 failed - expected 1 pending task, got %d", len(pendingTasks))
	}

	t.Log("Step 6 passed - tasks completed")

	// Step 7: Delete template
	t.Log("Step 7: Deleting template")
	err = templateStore.Delete(templateID)
	if err != nil {
		t.Fatalf("Step 7 failed - could not delete template: %v", err)
	}

	// Verify template is deleted
	loadedTemplates, err = templateStore.Load()
	if err != nil {
		t.Fatalf("Step 7 failed - could not load templates after deletion: %v", err)
	}

	if len(loadedTemplates) != 0 {
		t.Errorf("Step 7 failed - expected 0 templates after deletion, got %d", len(loadedTemplates))
	}

	// Verify tasks still exist (deleting template doesn't delete tasks)
	allTasks, _ = eng.ListTasks("")
	if len(allTasks) != tasksToCreate {
		t.Errorf("Step 7 failed - tasks should not be deleted when template is deleted")
	}

	t.Log("Step 7 passed - template deleted, tasks preserved")

	// Cleanup
	for _, id := range createdTaskIDs {
		eng.DeleteTask(id)
	}

	t.Log("✓ Complete template workflow test passed")
}

// TestMultipleTemplates tests working with multiple templates
func TestMultipleTemplates(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "focus-e2e-multi-template-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	templatePath := filepath.Join(tmpDir, "templates.json")
	templateStore := templates.NewStore(templatePath)

	// Step 1: Create multiple templates
	t.Log("Step 1: Creating multiple templates")
	templateDefinitions := []models.TaskTemplate{
		{
			Name:        "Daily Standup",
			Description: "Attend daily standup meeting",
			Priority:    "high",
			Categories:  []string{"meetings", "work"},
		},
		{
			Name:        "Code Review",
			Description: "Review team's pull requests",
			Priority:    "medium",
			Categories:  []string{"code-review", "collaboration"},
		},
		{
			Name:        "Documentation",
			Description: "Update project documentation",
			Priority:    "low",
			Categories:  []string{"docs", "maintenance"},
		},
	}

	templateIDs := make([]string, 0, 3)
	for _, tmpl := range templateDefinitions {
		created, err := templateStore.Add(tmpl)
		if err != nil {
			t.Fatalf("Step 1 failed - could not create template '%s': %v", tmpl.Name, err)
		}
		templateIDs = append(templateIDs, created.ID)
	}

	t.Logf("Step 1 passed - created %d templates", len(templateIDs))

	// Step 2: Load all templates
	t.Log("Step 2: Loading all templates")
	allTemplates, err := templateStore.Load()
	if err != nil {
		t.Fatalf("Step 2 failed: %v", err)
	}

	if len(allTemplates) != 3 {
		t.Errorf("Step 2 failed - expected 3 templates, got %d", len(allTemplates))
	}

	t.Log("Step 2 passed - all templates loaded")

	// Step 3: Get each template by name
	t.Log("Step 3: Retrieving templates by name")
	for _, tmpl := range templateDefinitions {
		retrieved, err := templateStore.GetByName(tmpl.Name)
		if err != nil {
			t.Errorf("Step 3 failed - could not get template '%s': %v", tmpl.Name, err)
			continue
		}

		if retrieved.Description != tmpl.Description {
			t.Errorf("Step 3 failed - description mismatch for '%s'", tmpl.Name)
		}
	}

	t.Log("Step 3 passed - all templates retrieved by name")

	// Step 4: Delete templates
	t.Log("Step 4: Deleting all templates")
	for _, id := range templateIDs {
		if err := templateStore.Delete(id); err != nil {
			t.Errorf("Step 4 failed - could not delete template: %v", err)
		}
	}

	allTemplates, _ = templateStore.Load()
	if len(allTemplates) != 0 {
		t.Errorf("Step 4 failed - expected 0 templates after deletion, got %d", len(allTemplates))
	}

	t.Log("Step 4 passed - all templates deleted")
	t.Log("✓ Multiple templates test passed")
}
