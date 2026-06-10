package character

import "time"

// Character represents a story character
type Character struct {
	ID            string         `yaml:"id"`
	Name          string         `yaml:"name"`
	Aliases       []string       `yaml:"aliases,omitempty"`
	Role          string         `yaml:"role,omitempty"` // protagonist, antagonist, supporting, minor
	Age           int            `yaml:"age,omitempty"`
	Occupation    string         `yaml:"occupation,omitempty"`
	Appearance    string         `yaml:"appearance,omitempty"`
	Background    string         `yaml:"background,omitempty"`
	Arc           string         `yaml:"arc,omitempty"`
	Relationships []Relationship `yaml:"relationships,omitempty"`
	CreatedAt     time.Time      `yaml:"created_at"`
	UpdatedAt     time.Time      `yaml:"updated_at"`
	Bio           string         `yaml:"-"` // Markdown content after frontmatter
}

// Relationship represents a relationship between characters
type Relationship struct {
	CharacterID string `yaml:"character_id"`
	Type        string `yaml:"type"`        // love interest, rival, mentor, etc
	Tension     string `yaml:"tension"`     // low, medium, high
	Notes       string `yaml:"notes,omitempty"`
}
