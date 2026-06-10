package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/kyanite/syntax/internal/location"
	"gopkg.in/yaml.v3"
)

// SaveLocation saves a location to disk
func SaveLocation(projectDir string, loc *location.Location) error {
	if loc.ID == "" {
		loc.ID = GenerateLocationID()
	}

	if loc.CreatedAt.IsZero() {
		loc.CreatedAt = time.Now().UTC()
	}
	loc.UpdatedAt = time.Now().UTC()

	// Create frontmatter
	var buf bytes.Buffer
	buf.WriteString("---\n")

	// Marshal YAML frontmatter
	yamlData, err := yaml.Marshal(loc)
	if err != nil {
		return fmt.Errorf("failed to marshal location: %w", err)
	}
	buf.Write(yamlData)
	buf.WriteString("---\n\n")

	// Add description content
	if loc.Description != "" {
		buf.WriteString(loc.Description)
	} else {
		buf.WriteString(fmt.Sprintf("# %s\n\n", loc.Name))
		buf.WriteString("(Add location description here...)\n")
	}

	// Write to file
	locPath := filepath.Join(projectDir, "locations", loc.ID+".md")
	return AtomicWriteFile(locPath, buf.Bytes(), 0600)
}

// LoadLocation loads a location from disk
func LoadLocation(projectDir, locationID string) (*location.Location, error) {
	locPath := filepath.Join(projectDir, "locations", locationID+".md")

	data, err := os.ReadFile(locPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read location file: %w", err)
	}

	var loc location.Location
	rest, err := frontmatter.Parse(bytes.NewReader(data), &loc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	loc.Description = string(rest)
	return &loc, nil
}

// LoadAllLocations loads all locations for a project
func LoadAllLocations(projectDir string) (map[string]*location.Location, error) {
	locationsDir := filepath.Join(projectDir, "locations")
	locations := make(map[string]*location.Location)

	entries, err := os.ReadDir(locationsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return locations, nil
		}
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		locID := entry.Name()[:len(entry.Name())-3]
		loc, err := LoadLocation(projectDir, locID)
		if err != nil {
			continue
		}
		locations[locID] = loc
	}

	return locations, nil
}

// DeleteLocation deletes a location
func DeleteLocation(projectDir, locationID string) error {
	locPath := filepath.Join(projectDir, "locations", locationID+".md")
	return os.Remove(locPath)
}
