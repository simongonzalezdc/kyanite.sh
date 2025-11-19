package templates

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestStore_Add(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	template := models.TaskTemplate{
		Name:        "Test Template",
		Description: "Test description",
		Priority:    "high",
		Categories:  []string{"work", "test"},
		Notes:       "Some notes",
	}

	addedTemplate, err := store.Add(template)
	if err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	if addedTemplate.ID == "" {
		t.Error("Add() should generate an ID")
	}

	if addedTemplate.Name != template.Name {
		t.Errorf("Add() name = %v, want %v", addedTemplate.Name, template.Name)
	}

	if addedTemplate.CreatedAt.IsZero() {
		t.Error("Add() should set CreatedAt")
	}

	if addedTemplate.UpdatedAt.IsZero() {
		t.Error("Add() should set UpdatedAt")
	}
}

func TestStore_Load(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	// Add some templates
	template1 := models.TaskTemplate{
		Name:        "Template 1",
		Description: "Description 1",
		Priority:    "high",
	}
	template2 := models.TaskTemplate{
		Name:        "Template 2",
		Description: "Description 2",
		Priority:    "low",
	}

	store.Add(template1)
	store.Add(template2)

	// Create new store instance and load
	newStore := NewStore(testFile)
	templates, err := newStore.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if len(templates) != 2 {
		t.Errorf("Load() returned %d templates, want 2", len(templates))
	}
}

func TestStore_Get(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	template := models.TaskTemplate{
		Name:        "Test Template",
		Description: "Test description",
	}

	added, _ := store.Add(template)

	// Get existing template
	retrieved, err := store.Get(added.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if retrieved.ID != added.ID {
		t.Errorf("Get() ID = %v, want %v", retrieved.ID, added.ID)
	}

	// Get non-existent template
	_, err = store.Get("non-existent")
	if err == nil {
		t.Error("Get() should return error for non-existent template")
	}
}

func TestStore_GetByName(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	template := models.TaskTemplate{
		Name:        "Unique Name",
		Description: "Test description",
	}

	store.Add(template)

	// Get by exact name
	retrieved, err := store.GetByName("Unique Name")
	if err != nil {
		t.Fatalf("GetByName() error = %v", err)
	}

	if retrieved.Name != template.Name {
		t.Errorf("GetByName() name = %v, want %v", retrieved.Name, template.Name)
	}

	// Get by non-existent name
	_, err = store.GetByName("Non Existent")
	if err == nil {
		t.Error("GetByName() should return error for non-existent name")
	}
}

func TestStore_Delete(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	template := models.TaskTemplate{
		Name:        "To Delete",
		Description: "Will be deleted",
	}

	added, _ := store.Add(template)

	// Delete the template
	err := store.Delete(added.ID)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Try to get deleted template
	_, err = store.Get(added.ID)
	if err == nil {
		t.Error("Get() should return error for deleted template")
	}

	// Delete non-existent template should not error (idempotent)
	err = store.Delete("non-existent")
	if err != nil {
		t.Errorf("Delete() should not error for non-existent template, got: %v", err)
	}
}

func TestStore_Save(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	templates := []models.TaskTemplate{
		{
			ID:          "1",
			Name:        "Template 1",
			Description: "Description 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "Template 2",
			Description: "Description 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	err := store.Save(templates)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Save() should create file")
	}

	// Load and verify
	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load() after Save() error = %v", err)
	}

	if len(loaded) != len(templates) {
		t.Errorf("Load() returned %d templates, want %d", len(loaded), len(templates))
	}
}

func TestStore_EmptyLoad(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "nonexistent.json")

	store := NewStore(testFile)

	templates, err := store.Load()
	if err != nil {
		t.Fatalf("Load() on non-existent file should not error, got %v", err)
	}

	if len(templates) != 0 {
		t.Errorf("Load() on non-existent file should return empty slice, got %d templates", len(templates))
	}
}

func TestTaskTemplate_ToTask(t *testing.T) {
	template := models.TaskTemplate{
		ID:                 "template-1",
		Name:               "Test Template",
		Description:        "Test task description",
		Priority:           "high",
		Categories:         []string{"work", "urgent"},
		Notes:              "Some notes",
		RecurrencePattern:  models.RecurrenceDaily,
		RecurrenceInterval: 2,
	}

	task := template.ToTask()

	if task.Description != template.Description {
		t.Errorf("ToTask() Description = %v, want %v", task.Description, template.Description)
	}

	if task.Priority != template.Priority {
		t.Errorf("ToTask() Priority = %v, want %v", task.Priority, template.Priority)
	}

	if len(task.Categories) != len(template.Categories) {
		t.Errorf("ToTask() Categories length = %v, want %v", len(task.Categories), len(template.Categories))
	}

	if task.Notes != template.Notes {
		t.Errorf("ToTask() Notes = %v, want %v", task.Notes, template.Notes)
	}

	if task.RecurrencePattern != template.RecurrencePattern {
		t.Errorf("ToTask() RecurrencePattern = %v, want %v", task.RecurrencePattern, template.RecurrencePattern)
	}

	if task.RecurrenceInterval != template.RecurrenceInterval {
		t.Errorf("ToTask() RecurrenceInterval = %v, want %v", task.RecurrenceInterval, template.RecurrenceInterval)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "templates.json")

	store := NewStore(testFile)

	// Add initial template
	template := models.TaskTemplate{
		Name:        "Concurrent Test",
		Description: "Test concurrent access",
	}
	store.Add(template)

	// Simulate concurrent reads
	done := make(chan bool)
	for i := 0; i < 5; i++ {
		go func() {
			_, err := store.Load()
			if err != nil {
				t.Errorf("Concurrent Load() error = %v", err)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}
}
