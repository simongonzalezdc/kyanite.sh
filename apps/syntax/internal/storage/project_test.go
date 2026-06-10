package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T) (string, func()) {
	// Create a unique temporary directory for each test
	tempDir := t.TempDir()

	// Override XDG_DATA_HOME to use temp directory
	originalDataHome := os.Getenv("XDG_DATA_HOME")
	os.Setenv("XDG_DATA_HOME", tempDir)

	// Create the full project path that will be used
	projectsDir := filepath.Join(tempDir, "syntax", "projects")
	os.MkdirAll(projectsDir, 0700)

	// Return cleanup function
	return projectsDir, func() {
		if originalDataHome != "" {
			os.Setenv("XDG_DATA_HOME", originalDataHome)
		} else {
			os.Unsetenv("XDG_DATA_HOME")
		}
	}
}

func TestCreateProject(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	project, err := CreateProject("Test Project", "Test Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	if project.Title != "Test Project" {
		t.Errorf("Expected title 'Test Project', got '%s'", project.Title)
	}

	if project.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", project.Author)
	}

	if project.Genre != "Fiction" {
		t.Errorf("Expected genre 'Fiction', got '%s'", project.Genre)
	}

	if project.Status != "draft" {
		t.Errorf("Expected status 'draft', got '%s'", project.Status)
	}

	// Check that directories were created
	dirs := []string{"characters", "locations", "scenes", "outline", "stats", "exports", ".backups"}
	for _, dir := range dirs {
		dirPath := filepath.Join(project.Directory, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", dir)
		}
	}

	// Check that metadata file was created
	metadataPath := filepath.Join(project.Directory, "metadata.yaml")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Error("Expected metadata.yaml to exist")
	}
}

func TestLoadProject(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a project first
	originalProject, err := CreateProject("Test Load", "Author", "Genre")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Load the project
	loadedProject, err := LoadProject(originalProject.ID)
	if err != nil {
		t.Fatalf("Failed to load project: %v", err)
	}

	if loadedProject.Title != originalProject.Title {
		t.Errorf("Expected title '%s', got '%s'", originalProject.Title, loadedProject.Title)
	}

	if loadedProject.Author != originalProject.Author {
		t.Errorf("Expected author '%s', got '%s'", originalProject.Author, loadedProject.Author)
	}

	if loadedProject.Genre != originalProject.Genre {
		t.Errorf("Expected genre '%s', got '%s'", originalProject.Genre, loadedProject.Genre)
	}
}

func TestSaveProjectMetadata(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	project, err := CreateProject("Test Save", "Author", "Genre")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Modify project
	project.Status = "revising"
	project.TargetWordCount = 100000

	// Save metadata
	if err := SaveProjectMetadata(project); err != nil {
		t.Fatalf("Failed to save metadata: %v", err)
	}

	// Load and verify
	loaded, err := LoadProject(project.ID)
	if err != nil {
		t.Fatalf("Failed to load project: %v", err)
	}

	if loaded.Status != "revising" {
		t.Errorf("Expected status 'revising', got '%s'", loaded.Status)
	}

	if loaded.TargetWordCount != 100000 {
		t.Errorf("Expected target word count 100000, got %d", loaded.TargetWordCount)
	}
}

func TestListProjects(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple projects with unique names
	p1, err := CreateProject("TestList Project Alpha", "Author 1", "Genre 1")
	if err != nil {
		t.Fatalf("Failed to create project 1: %v", err)
	}

	p2, err := CreateProject("TestList Project Beta", "Author 2", "Genre 2")
	if err != nil {
		t.Fatalf("Failed to create project 2: %v", err)
	}

	// List projects
	projects, err := ListProjects()
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}

	if len(projects) < 2 {
		t.Errorf("Expected at least 2 projects, got %d", len(projects))
	}

	// Verify our projects are in the list
	foundP1, foundP2 := false, false
	for _, p := range projects {
		if p.ID == p1.ID {
			foundP1 = true
		}
		if p.ID == p2.ID {
			foundP2 = true
		}
	}

	if !foundP1 {
		t.Error("Project 1 not found in list")
	}
	if !foundP2 {
		t.Error("Project 2 not found in list")
	}
}
