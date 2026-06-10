package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/kyanite/syntax/internal/scene"
	"gopkg.in/yaml.v3"
)

// SaveScene saves a scene to disk
func SaveScene(projectDir string, sc *scene.Scene) error {
	if sc.ID == "" {
		sc.ID = GenerateSceneID(sc.Chapter, sc.SceneNumber)
	}

	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = time.Now().UTC()
	}
	sc.UpdatedAt = time.Now().UTC()

	// Calculate word count
	sc.WordCount = calculateWordCount(sc.Content)

	// Create frontmatter
	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Marshal YAML frontmatter
	yamlData, err := yaml.Marshal(sc)
	if err != nil {
		return fmt.Errorf("failed to marshal scene: %w", err)
	}
	buf.Write(yamlData)
	buf.WriteString("---\n\n")

	// Add scene content
	if sc.Content != "" {
		buf.WriteString(sc.Content)
	} else {
		buf.WriteString(fmt.Sprintf("# Chapter %d, Scene %d: %s\n\n", sc.Chapter, sc.SceneNumber, sc.Name))
		buf.WriteString("(Start writing your scene here...)\n")
	}

	// Write to file
	scenePath := filepath.Join(projectDir, "scenes", sc.ID+".md")
	return AtomicWriteFile(scenePath, buf.Bytes(), 0600)
}

// LoadScene loads a scene from disk
func LoadScene(projectDir, sceneID string) (*scene.Scene, error) {
	scenePath := filepath.Join(projectDir, "scenes", sceneID+".md")

	data, err := os.ReadFile(scenePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read scene file: %w", err)
	}

	var sc scene.Scene
	rest, err := frontmatter.Parse(bytes.NewReader(data), &sc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	sc.Content = string(rest)
	return &sc, nil
}

// LoadAllScenes loads all scenes for a project
func LoadAllScenes(projectDir string) (map[string]*scene.Scene, error) {
	scenesDir := filepath.Join(projectDir, "scenes")
	scenes := make(map[string]*scene.Scene)

	entries, err := os.ReadDir(scenesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return scenes, nil // Empty map if directory doesn't exist
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		sceneID := entry.Name()[:len(entry.Name())-3] // Remove .md extension
		sc, err := LoadScene(projectDir, sceneID)
		if err != nil {
			continue // Skip invalid scenes
		}
		scenes[sceneID] = sc
	}

	return scenes, nil
}

// DeleteScene deletes a scene
func DeleteScene(projectDir, sceneID string) error {
	scenePath := filepath.Join(projectDir, "scenes", sceneID+".md")
	return os.Remove(scenePath)
}

// calculateWordCount counts words in text (simple implementation)
func calculateWordCount(text string) int {
	// Remove markdown syntax (basic)
	cleaned := text

	// Split by whitespace
	words := strings.Fields(cleaned)

	return len(words)
}
