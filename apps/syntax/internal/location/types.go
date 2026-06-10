package location

import "time"

// Location represents a story location
type Location struct {
	ID           string       `yaml:"id"`
	Name         string       `yaml:"name"`
	Type         string       `yaml:"type,omitempty"` // city, tavern, castle, forest, planet, etc
	Region       string       `yaml:"region,omitempty"`
	Climate      string       `yaml:"climate,omitempty"`
	Population   int          `yaml:"population,omitempty"`
	Significance string       `yaml:"significance,omitempty"`
	Connections  []Connection `yaml:"connections,omitempty"`
	CreatedAt    time.Time    `yaml:"created_at"`
	UpdatedAt    time.Time    `yaml:"updated_at"`
	Description  string       `yaml:"-"` // Markdown content after frontmatter
}

// Connection represents a connection to another location
type Connection struct {
	LocationID     string `yaml:"location_id"`
	ConnectionType string `yaml:"connection_type"` // road, river, alley, etc
	Distance       string `yaml:"distance,omitempty"`
}
