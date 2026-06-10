package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/storage"
)

// TestProjectCreationAndSceneManagement tests the complete workflow of
// creating a project and managing scenes
func TestProjectCreationAndSceneManagement(t *testing.T) {
	// Setup: Create temp directory for test project
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Step 1: Create a new project
	project, err := storage.CreateProject("Test Novel", "Test Author", "Fantasy")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	if project.Title != "Test Novel" {
		t.Errorf("Project title = %q, expected %q", project.Title, "Test Novel")
	}

	// Step 2: Add scenes to the project
	scene1 := &scene.Scene{
		ID:          "scene-1",
		Name:        "Opening Scene",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "It was a dark and stormy night...",
		WordCount:   7,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	scene2 := &scene.Scene{
		ID:          "scene-2",
		Name:        "Conflict Begins",
		Chapter:     1,
		SceneNumber: 2,
		Content:     "The hero faced a difficult choice.",
		WordCount:   6,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	scene3 := &scene.Scene{
		ID:          "scene-3",
		Name:        "Resolution",
		Chapter:     2,
		SceneNumber: 1,
		Content:     "In the end, wisdom prevailed.",
		WordCount:   5,
		Status:      "done",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Step 3: Save scenes
	if err := storage.SaveScene(project.Directory, scene1); err != nil {
		t.Fatalf("Failed to save scene 1: %v", err)
	}
	if err := storage.SaveScene(project.Directory, scene2); err != nil {
		t.Fatalf("Failed to save scene 2: %v", err)
	}
	if err := storage.SaveScene(project.Directory, scene3); err != nil {
		t.Fatalf("Failed to save scene 3: %v", err)
	}

	// Step 4: Load all scenes
	scenes, err := storage.LoadAllScenes(project.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	if len(scenes) != 3 {
		t.Errorf("Loaded %d scenes, expected 3", len(scenes))
	}

	// Step 5: Compile scenes to check for issues
	report := scene.CompileScenes(scenes)

	if report.TotalScenes != 3 {
		t.Errorf("Report shows %d scenes, expected 3", report.TotalScenes)
	}

	if report.ChapterCount != 2 {
		t.Errorf("Report shows %d chapters, expected 2", report.ChapterCount)
	}

	if len(report.Issues) != 0 {
		t.Errorf("Found %d issues in valid scenes: %v", len(report.Issues), report.Issues)
	}

	// Step 6: Verify sorted order
	sortedScenes := scene.SortScenes(scenes)

	if len(sortedScenes) != 3 {
		t.Fatalf("Sorted %d scenes, expected 3", len(sortedScenes))
	}

	// Check order
	if sortedScenes[0].ID != "scene-1" {
		t.Errorf("First sorted scene = %s, expected scene-1", sortedScenes[0].ID)
	}
	if sortedScenes[1].ID != "scene-2" {
		t.Errorf("Second sorted scene = %s, expected scene-2", sortedScenes[1].ID)
	}
	if sortedScenes[2].ID != "scene-3" {
		t.Errorf("Third sorted scene = %s, expected scene-3", sortedScenes[2].ID)
	}
}

// TestSceneCompilationWithIssues tests the integration between scene
// creation and compilation with various issues
func TestSceneCompilationWithIssues(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	project, err := storage.CreateProject("Test Project", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Create scenes with various issues
	orphan := &scene.Scene{
		ID:        "orphan",
		Name:      "Orphaned Scene",
		Chapter:   0, // Missing chapter
		Content:   "Lost scene",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	scene1 := &scene.Scene{
		ID:          "scene1",
		Name:        "Scene 1",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "First scene",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	scene3 := &scene.Scene{
		ID:          "scene3",
		Name:        "Scene 3",
		Chapter:     1,
		SceneNumber: 3, // Gap at scene 2
		Content:     "Third scene",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dup1 := &scene.Scene{
		ID:          "dup1",
		Name:        "Duplicate A",
		Chapter:     2,
		SceneNumber: 1,
		Content:     "Duplicate scene A",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	dup2 := &scene.Scene{
		ID:          "dup2",
		Name:        "Duplicate B",
		Chapter:     2,
		SceneNumber: 1, // Duplicate
		Content:     "Duplicate scene B",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save all scenes
	for _, sc := range []*scene.Scene{orphan, scene1, scene3, dup1, dup2} {
		if err := storage.SaveScene(project.Directory, sc); err != nil {
			t.Fatalf("Failed to save scene %s: %v", sc.ID, err)
		}
	}

	// Load and compile
	scenes, err := storage.LoadAllScenes(project.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	report := scene.CompileScenes(scenes)

	// Verify orphan detected
	if len(report.Orphans) != 1 {
		t.Errorf("Found %d orphans, expected 1", len(report.Orphans))
	}

	// Verify gap detected
	if len(report.Gaps) != 1 {
		t.Errorf("Found %d gaps, expected 1 (scene 2)", len(report.Gaps))
	} else {
		gap := report.Gaps[0]
		if gap.Chapter != 1 || gap.SceneNumber != 2 {
			t.Errorf("Gap at chapter %d scene %d, expected chapter 1 scene 2", gap.Chapter, gap.SceneNumber)
		}
	}

	// Verify duplicate detected
	if len(report.Duplicates) != 1 {
		t.Errorf("Found %d duplicate groups, expected 1", len(report.Duplicates))
	}

	// Verify errors and warnings
	if !report.HasErrors() {
		t.Error("Report should have errors (duplicates)")
	}

	if !report.HasWarnings() {
		t.Error("Report should have warnings (orphans, gaps)")
	}

	summary := report.Summary()
	if summary == "✓ All scenes are properly organized!" {
		t.Error("Summary should indicate issues found")
	}
}

// TestProjectPersistence tests that project data persists correctly
func TestProjectPersistence(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Create and save project
	originalProject, err := storage.CreateProject("Persistent Novel", "Author Name", "SciFi")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	originalID := originalProject.ID

	// Add a scene
	testScene := &scene.Scene{
		ID:          "test-scene",
		Name:        "Test Scene",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "Test content that should persist.",
		WordCount:   5,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := storage.SaveScene(originalProject.Directory, testScene); err != nil {
		t.Fatalf("Failed to save scene: %v", err)
	}

	// Load project from disk
	loadedProject, err := storage.LoadProject(originalID)
	if err != nil {
		t.Fatalf("Failed to load project: %v", err)
	}

	// Verify project data persisted
	if loadedProject.Title != "Persistent Novel" {
		t.Errorf("Loaded title = %q, expected %q", loadedProject.Title, "Persistent Novel")
	}

	if loadedProject.Author != "Author Name" {
		t.Errorf("Loaded author = %q, expected %q", loadedProject.Author, "Author Name")
	}

	if loadedProject.Genre != "SciFi" {
		t.Errorf("Loaded genre = %q, expected %q", loadedProject.Genre, "SciFi")
	}

	// Load scenes for the project
	loadedScenes, err := storage.LoadAllScenes(loadedProject.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	if len(loadedScenes) != 1 {
		t.Fatalf("Loaded %d scenes, expected 1", len(loadedScenes))
	}

	loadedScene := loadedScenes["test-scene"]
	if loadedScene == nil {
		t.Fatal("Scene 'test-scene' not found")
	}

	// Content may have whitespace from frontmatter parsing
	expectedContent := "Test content that should persist."
	if loadedScene.Content != expectedContent && loadedScene.Content != "\n"+expectedContent {
		t.Errorf("Scene content = %q, expected %q (or with leading newline)", loadedScene.Content, expectedContent)
	}
}

// TestProjectListing tests listing all projects
func TestProjectListing(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Create multiple projects
	project1, err := storage.CreateProject("Novel One", "Author A", "Fantasy")
	if err != nil {
		t.Fatalf("Failed to create project 1: %v", err)
	}

	project2, err := storage.CreateProject("Novel Two", "Author B", "SciFi")
	if err != nil {
		t.Fatalf("Failed to create project 2: %v", err)
	}

	project3, err := storage.CreateProject("Novel Three", "Author C", "Mystery")
	if err != nil {
		t.Fatalf("Failed to create project 3: %v", err)
	}

	// List all projects
	projects, err := storage.ListProjects()
	if err != nil {
		t.Fatalf("Failed to list projects: %v", err)
	}

	// Verify all created projects are in the list
	// (there may be other projects from previous tests, so we check >= 3)
	if len(projects) < 3 {
		t.Errorf("Found %d projects, expected at least 3", len(projects))
	}

	foundProjects := make(map[string]bool)
	for _, p := range projects {
		foundProjects[p.ID] = true
	}

	if !foundProjects[project1.ID] {
		t.Error("Project 1 not found in listing")
	}
	if !foundProjects[project2.ID] {
		t.Error("Project 2 not found in listing")
	}
	if !foundProjects[project3.ID] {
		t.Error("Project 3 not found in listing")
	}
}

// TestSceneFileStructure verifies scene files are created in correct locations
func TestSceneFileStructure(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	project, err := storage.CreateProject("Structure Test", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	testScene := &scene.Scene{
		ID:          "file-test",
		Name:        "File Test Scene",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "Testing file structure",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := storage.SaveScene(project.Directory, testScene); err != nil {
		t.Fatalf("Failed to save scene: %v", err)
	}

	// Verify file exists
	scenePath := filepath.Join(project.Directory, "scenes", "file-test.md")
	if _, err := os.Stat(scenePath); os.IsNotExist(err) {
		t.Errorf("Scene file not created at %s", scenePath)
	}

	// Verify content can be read back
	loadedScenes, err := storage.LoadAllScenes(project.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	if loadedScene, ok := loadedScenes["file-test"]; !ok {
		t.Error("Scene not loaded from file")
	} else {
		expectedContent := "Testing file structure"
		if loadedScene.Content != expectedContent && loadedScene.Content != "\n"+expectedContent {
			t.Errorf("Scene content = %q, expected %q (or with leading newline)", loadedScene.Content, expectedContent)
		}
	}
}
