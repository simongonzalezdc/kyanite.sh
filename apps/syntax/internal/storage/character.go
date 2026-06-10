package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/kyanite/syntax/internal/character"
	"gopkg.in/yaml.v3"
)

// SaveCharacter saves a character to disk
func SaveCharacter(projectDir string, char *character.Character) error {
	if char.ID == "" {
		char.ID = GenerateCharacterID()
	}

	if char.CreatedAt.IsZero() {
		char.CreatedAt = time.Now().UTC()
	}
	char.UpdatedAt = time.Now().UTC()

	// Create frontmatter
	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Marshal YAML frontmatter
	yamlData, err := yaml.Marshal(char)
	if err != nil {
		return fmt.Errorf("failed to marshal character: %w", err)
	}
	buf.Write(yamlData)
	buf.WriteString("---\n\n")

	// Add bio content
	if char.Bio != "" {
		buf.WriteString(char.Bio)
	} else {
		buf.WriteString(fmt.Sprintf("# %s - Character Biography\n\n", char.Name))
		buf.WriteString("(Add character biography here...)\n")
	}

	// Write to file
	charPath := filepath.Join(projectDir, "characters", char.ID+".md")
	return AtomicWriteFile(charPath, buf.Bytes(), 0600)
}

// LoadCharacter loads a character from disk
func LoadCharacter(projectDir, characterID string) (*character.Character, error) {
	charPath := filepath.Join(projectDir, "characters", characterID+".md")

	data, err := os.ReadFile(charPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read character file: %w", err)
	}

	var char character.Character
	rest, err := frontmatter.Parse(bytes.NewReader(data), &char)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	char.Bio = string(rest)
	return &char, nil
}

// LoadAllCharacters loads all characters for a project
func LoadAllCharacters(projectDir string) (map[string]*character.Character, error) {
	charactersDir := filepath.Join(projectDir, "characters")
	characters := make(map[string]*character.Character)

	entries, err := os.ReadDir(charactersDir)
	if err != nil {
		if os.IsNotExist(err) {
			return characters, nil // Empty map if directory doesn't exist
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		charID := entry.Name()[:len(entry.Name())-3] // Remove .md extension
		char, err := LoadCharacter(projectDir, charID)
		if err != nil {
			continue // Skip invalid characters
		}
		characters[charID] = char
	}

	return characters, nil
}

// DeleteCharacter deletes a character
func DeleteCharacter(projectDir, characterID string) error {
	charPath := filepath.Join(projectDir, "characters", characterID+".md")
	return os.Remove(charPath)
}
