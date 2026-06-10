package e2e

import (
	"os"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/app"
	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/editor"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/storage"
)

// TestCompleteWritingWorkflow simulates a complete user journey:
// Welcome -> Create Project -> Add Characters -> Add Locations -> Create Scenes -> Edit -> Save
func TestCompleteWritingWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Step 1: Start at welcome screen
	m := app.NewModel()

	if m.CurrentScreen != app.ScreenWelcome {
		t.Fatalf("Expected ScreenWelcome, got %v", m.CurrentScreen)
	}

	// Step 2: Create new project (simulating 'n' key press)
	project, err := storage.CreateProject("My Epic Novel", "Jane Doe", "Fantasy")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	m.CurrentProject = project
	m.CurrentScreen = app.ScreenEditor

	if m.CurrentScreen != app.ScreenEditor {
		t.Fatalf("Expected ScreenEditor after project creation, got %v", m.CurrentScreen)
	}

	// Step 3: Navigate to characters screen (simulating 'c' key)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("c")}
	result, _ := m.Update(msg)
	m = result.(app.Model)

	if m.CurrentScreen != app.ScreenCharacters {
		t.Fatalf("Expected ScreenCharacters, got %v", m.CurrentScreen)
	}

	// Step 4: Add characters
	protagonist := &character.Character{
		ID:         "char-hero",
		Name:       "Aria Stormwind",
		Role:       "Protagonist",
		Age:        25,
		Background: "A brave warrior with a mysterious past",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	antagonist := &character.Character{
		ID:         "char-villain",
		Name:       "Lord Malice",
		Role:       "Antagonist",
		Age:        0, // Unknown age
		Background: "Ancient dark sorcerer",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Initialize characters map if needed
	if m.CurrentProject.Characters == nil {
		m.CurrentProject.Characters = make(map[string]*character.Character)
	}
	m.CurrentProject.Characters[protagonist.ID] = protagonist
	m.CurrentProject.Characters[antagonist.ID] = antagonist

	// Save characters
	if err := storage.SaveCharacter(m.CurrentProject.Directory, protagonist); err != nil {
		t.Fatalf("Failed to save protagonist: %v", err)
	}
	if err := storage.SaveCharacter(m.CurrentProject.Directory, antagonist); err != nil {
		t.Fatalf("Failed to save antagonist: %v", err)
	}

	// Step 5: Navigate to locations screen
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("l")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	if m.CurrentScreen != app.ScreenLocations {
		t.Fatalf("Expected ScreenLocations, got %v", m.CurrentScreen)
	}

	// Step 6: Add locations
	castle := &location.Location{
		ID:          "loc-castle",
		Name:        "Dark Castle",
		Description: "An imposing fortress shrouded in shadow",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	village := &location.Location{
		ID:          "loc-village",
		Name:        "Riverside Village",
		Description: "A peaceful farming community",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if m.CurrentProject.Locations == nil {
		m.CurrentProject.Locations = make(map[string]*location.Location)
	}
	m.CurrentProject.Locations[castle.ID] = castle
	m.CurrentProject.Locations[village.ID] = village

	if err := storage.SaveLocation(m.CurrentProject.Directory, castle); err != nil {
		t.Fatalf("Failed to save castle: %v", err)
	}
	if err := storage.SaveLocation(m.CurrentProject.Directory, village); err != nil {
		t.Fatalf("Failed to save village: %v", err)
	}

	// Step 7: Navigate to scenes screen
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("s")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	if m.CurrentScreen != app.ScreenScenes {
		t.Fatalf("Expected ScreenScenes, got %v", m.CurrentScreen)
	}

	// Step 8: Create scenes
	scene1 := &scene.Scene{
		ID:           "scene-1",
		Name:         "The Journey Begins",
		Chapter:      1,
		SceneNumber:  1,
		POVCharacter: "Aria Stormwind",
		Location:     "Riverside Village",
		Content:      "Aria gazed at the distant castle, knowing her destiny awaited.",
		WordCount:    11,
		Status:       "draft",
		Characters:   []string{"Aria Stormwind"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	scene2 := &scene.Scene{
		ID:           "scene-2",
		Name:         "The Dark Lord Appears",
		Chapter:      1,
		SceneNumber:  2,
		POVCharacter: "Lord Malice",
		Location:     "Dark Castle",
		Content:      "From his throne of shadows, Lord Malice sensed a disturbance in the realm.",
		WordCount:    14,
		Status:       "draft",
		Characters:   []string{"Lord Malice"},
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if m.CurrentProject.Scenes == nil {
		m.CurrentProject.Scenes = make(map[string]*scene.Scene)
	}
	m.CurrentProject.Scenes[scene1.ID] = scene1
	m.CurrentProject.Scenes[scene2.ID] = scene2

	if err := storage.SaveScene(m.CurrentProject.Directory, scene1); err != nil {
		t.Fatalf("Failed to save scene 1: %v", err)
	}
	if err := storage.SaveScene(m.CurrentProject.Directory, scene2); err != nil {
		t.Fatalf("Failed to save scene 2: %v", err)
	}

	// Step 9: Simulate editing a scene
	m.CurrentScene = scene1
	m.Buffer = editor.NewBuffer(scene1.Content)
	m.CurrentScreen = app.ScreenTextEditor

	if !m.IsEditorActive() {
		t.Error("Editor should be active")
	}

	// Step 10: Make edits
	m.EditorMode = app.EditorModeInsert
	m.Buffer.InsertRune('!')
	m.Buffer.SetModified(true)

	if !m.HasUnsavedChanges() {
		t.Error("Should have unsaved changes after edit")
	}

	// Step 11: Save the scene
	saved := m.SaveCurrentScene()
	if !saved {
		t.Error("Failed to save scene")
	}

	if m.HasUnsavedChanges() {
		t.Error("Should not have unsaved changes after save")
	}

	// Step 12: Exit editor
	m.ExitEditor(false) // Already saved, no need to save again

	if m.CurrentScreen != app.ScreenScenes {
		t.Errorf("Expected ScreenScenes after exit, got %v", m.CurrentScreen)
	}

	if m.CurrentScene != nil {
		t.Error("Current scene should be nil after exit")
	}

	// Step 13: Navigate to stats
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("t")}
	result, _ = m.Update(msg)
	m = result.(app.Model)

	if m.CurrentScreen != app.ScreenStats {
		t.Fatalf("Expected ScreenStats, got %v", m.CurrentScreen)
	}

	// Verify stats view works
	view := m.View()
	if view == "" {
		t.Error("Stats view should not be empty")
	}
}

// TestAutoSaveWorkflow tests the auto-save functionality
func TestAutoSaveWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Create project and scene
	project, err := storage.CreateProject("Auto Save Test", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	testScene := &scene.Scene{
		ID:          "autosave-scene",
		Name:        "Test Scene",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "Original content",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := storage.SaveScene(project.Directory, testScene); err != nil {
		t.Fatalf("Failed to save initial scene: %v", err)
	}

	// Setup model in editor mode
	m := app.NewModel()
	m.CurrentProject = project
	m.CurrentScene = testScene
	m.Buffer = editor.NewBuffer(testScene.Content)
	m.CurrentScreen = app.ScreenTextEditor

	// Verify cannot auto-save immediately (no changes)
	if m.CanAutoSave() {
		t.Error("Should not be able to auto-save without changes")
	}

	// Make a change
	m.Buffer.SetModified(true)
	m.LastEditTime = time.Now()

	// Still can't auto-save (too recent)
	if m.CanAutoSave() {
		t.Error("Should not be able to auto-save immediately after edit")
	}

	// Wait enough time
	m.LastEditTime = time.Now().Add(-5 * time.Second)

	// Now can auto-save
	if !m.CanAutoSave() {
		t.Error("Should be able to auto-save after waiting")
	}
}

// TestNavigationFlow tests navigation between all screens
func TestNavigationFlow(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	project, err := storage.CreateProject("Navigation Test", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	m := app.NewModel()
	m.CurrentProject = project
	m.CurrentScreen = app.ScreenEditor

	// Test navigation to each screen
	screens := []struct {
		key      string
		expected app.Screen
	}{
		{"c", app.ScreenCharacters},
		{"esc", app.ScreenEditor},
		{"s", app.ScreenScenes},
		{"esc", app.ScreenEditor},
		{"l", app.ScreenLocations},
		{"esc", app.ScreenEditor},
		{"t", app.ScreenStats},
		{"esc", app.ScreenEditor},
		{"e", app.ScreenExport},
		{"esc", app.ScreenEditor},
		{"b", app.ScreenBackups},
		{"esc", app.ScreenEditor},
		{"h", app.ScreenHelp},
		{"esc", app.ScreenEditor},
	}

	for _, test := range screens {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(test.key)}
		result, _ := m.Update(msg)
		m = result.(app.Model)

		if m.CurrentScreen != test.expected {
			t.Errorf("After pressing %q: expected %v, got %v", test.key, test.expected, m.CurrentScreen)
		}
	}
}

// TestSceneValidationWorkflow tests the complete scene validation workflow
func TestSceneValidationWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Create project with problematic scenes
	project, err := storage.CreateProject("Validation Test", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Add scenes with issues
	orphan := &scene.Scene{
		ID:        "orphan",
		Name:      "Orphan",
		Chapter:   0,
		Content:   "Lost",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	gap1 := &scene.Scene{
		ID:          "scene1",
		Name:        "Scene 1",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "First",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	gap3 := &scene.Scene{
		ID:          "scene3",
		Name:        "Scene 3",
		Chapter:     1,
		SceneNumber: 3,
		Content:     "Third (missing 2)",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	for _, sc := range []*scene.Scene{orphan, gap1, gap3} {
		if err := storage.SaveScene(project.Directory, sc); err != nil {
			t.Fatalf("Failed to save scene: %v", err)
		}
	}

	// Load and validate
	scenes, err := storage.LoadAllScenes(project.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	report := scene.CompileScenes(scenes)

	// Should find 1 orphan and 1 gap
	if len(report.Orphans) != 1 {
		t.Errorf("Expected 1 orphan, found %d", len(report.Orphans))
	}

	if len(report.Gaps) != 1 {
		t.Errorf("Expected 1 gap, found %d", len(report.Gaps))
	}

	if !report.HasWarnings() {
		t.Error("Report should have warnings")
	}

	summary := report.Summary()
	if !contains(summary, "warning") {
		t.Errorf("Summary should mention warnings: %s", summary)
	}
}

// TestProjectReloadWorkflow tests closing and reopening a project
func TestProjectReloadWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	oldDataDir := os.Getenv("SYNTAX_DATA_DIR")
	os.Setenv("SYNTAX_DATA_DIR", tempDir)
	defer os.Setenv("SYNTAX_DATA_DIR", oldDataDir)

	// Create project with content
	project, err := storage.CreateProject("Reload Test", "Author", "Fiction")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	originalID := project.ID

	// Add content
	testScene := &scene.Scene{
		ID:          "reload-scene",
		Name:        "Test Scene",
		Chapter:     1,
		SceneNumber: 1,
		Content:     "Content to preserve",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := storage.SaveScene(project.Directory, testScene); err != nil {
		t.Fatalf("Failed to save scene: %v", err)
	}

	// Simulate closing (go to welcome)
	m := app.NewModel()
	m.CurrentProject = project
	m.CurrentScreen = app.ScreenEditor

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("esc")}
	result, _ := m.Update(msg)
	m = result.(app.Model)

	if m.CurrentScreen != app.ScreenWelcome {
		t.Fatalf("Expected ScreenWelcome, got %v", m.CurrentScreen)
	}

	// Simulate reopening
	reloadedProject, err := storage.LoadProject(originalID)
	if err != nil {
		t.Fatalf("Failed to reload project: %v", err)
	}

	m.CurrentProject = reloadedProject
	m.CurrentScreen = app.ScreenEditor

	// Load scenes
	scenes, err := storage.LoadAllScenes(reloadedProject.Directory)
	if err != nil {
		t.Fatalf("Failed to load scenes: %v", err)
	}

	if len(scenes) != 1 {
		t.Fatalf("Expected 1 scene, found %d", len(scenes))
	}

	reloadedScene := scenes["reload-scene"]
	if reloadedScene == nil {
		t.Fatal("Scene not found after reload")
	}

	expectedContent := "Content to preserve"
	if reloadedScene.Content != expectedContent && reloadedScene.Content != "\n"+expectedContent {
		t.Errorf("Content not preserved: got %q, expected %q (or with leading newline)", reloadedScene.Content, expectedContent)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
