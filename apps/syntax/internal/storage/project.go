package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
	"gopkg.in/yaml.v3"
)

const Version = "1.0.0"

// CreateProject creates a new project with the given title
func CreateProject(title, author, genre string) (*story.Project, error) {
	dataDir := GetDataDir()
	if err := EnsureDir(dataDir); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	project := &story.Project{
		ID:              GenerateProjectID(),
		Title:           title,
		Author:          author,
		Genre:           genre,
		Status:          "draft",
		SchemaVersion:   "1.0",
		CreatedWith:     fmt.Sprintf("syntax.sh v%s", Version),
		CreatedAt:       time.Now().UTC(),
		LastModified:    time.Now().UTC(),
		TargetWordCount: 80000,
		DailyWordGoal:   1000,
		Characters:      make(map[string]*character.Character),
		Locations:       make(map[string]*location.Location),
		Scenes:          make(map[string]*scene.Scene),
	}

	// Create project directory
	projectDir := filepath.Join(dataDir, project.ID)
	project.Directory = projectDir

	if err := EnsureDir(projectDir); err != nil {
		return nil, fmt.Errorf("failed to create project directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{"characters", "locations", "scenes", "outline", "stats", "exports", ".backups"}
	for _, dir := range dirs {
		if err := EnsureDir(filepath.Join(projectDir, dir)); err != nil {
			return nil, fmt.Errorf("failed to create %s directory: %w", dir, err)
		}
	}

	// Save metadata
	if err := SaveProjectMetadata(project); err != nil {
		return nil, fmt.Errorf("failed to save metadata: %w", err)
	}

	return project, nil
}

// LoadProject loads a project from its directory
func LoadProject(projectID string) (*story.Project, error) {
	dataDir := GetDataDir()
	projectDir := filepath.Join(dataDir, projectID)

	// Check if project exists
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("project not found: %s", projectID)
	}

	// Load metadata
	metadataPath := filepath.Join(projectDir, "metadata.yaml")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata: %w", err)
	}

	var project story.Project
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}

	project.Directory = projectDir
	project.Characters = make(map[string]*character.Character)
	project.Locations = make(map[string]*location.Location)
	project.Scenes = make(map[string]*scene.Scene)

	return &project, nil
}

// SaveProjectMetadata saves project metadata
func SaveProjectMetadata(project *story.Project) error {
	project.LastModified = time.Now().UTC()

	data, err := yaml.Marshal(project)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := filepath.Join(project.Directory, "metadata.yaml")
	return AtomicWriteFile(metadataPath, data, 0600)
}

// ListProjects returns all available projects
func ListProjects() ([]*story.Project, error) {
	dataDir := GetDataDir()

	// Ensure data directory exists
	if err := EnsureDir(dataDir); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dataDir)
	if err != nil {
		return nil, err
	}

	var projects []*story.Project
	for _, entry := range entries {
		if entry.IsDir() {
			project, err := LoadProject(entry.Name())
			if err != nil {
				continue // Skip invalid projects
			}
			projects = append(projects, project)
		}
	}

	return projects, nil
}
