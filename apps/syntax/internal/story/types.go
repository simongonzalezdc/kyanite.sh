package story

import (
	"time"

	"github.com/kyanite/syntax/internal/character"
	"github.com/kyanite/syntax/internal/location"
	"github.com/kyanite/syntax/internal/scene"
)

// Project represents a writing project
type Project struct {
	// Metadata
	ID              string    `yaml:"id"`
	Title           string    `yaml:"title"`
	Author          string    `yaml:"author,omitempty"`
	Genre           string    `yaml:"genre,omitempty"`
	Status          string    `yaml:"status"` // draft, revising, complete
	SchemaVersion   string    `yaml:"schema_version"`
	CreatedWith     string    `yaml:"created_with"`
	CreatedAt       time.Time `yaml:"created_at"`
	LastModified    time.Time `yaml:"last_modified"`
	TargetWordCount int       `yaml:"target_word_count,omitempty"`
	DailyWordGoal   int       `yaml:"daily_word_goal,omitempty"`

	// Cached statistics (for performance)
	TotalWords      int `yaml:"total_words"`
	TotalScenes     int `yaml:"total_scenes"`
	TotalCharacters int `yaml:"total_characters"`
	TotalLocations  int `yaml:"total_locations"`

	// Session tracking
	CurrentStreak      int `yaml:"current_streak,omitempty"`
	TotalSessions      int `yaml:"total_sessions,omitempty"`
	TotalTimeSeconds   int `yaml:"total_time_seconds,omitempty"`

	// Runtime data (not persisted)
	Directory  string                       `yaml:"-"`
	Characters map[string]*character.Character `yaml:"-"`
	Locations  map[string]*location.Location   `yaml:"-"`
	Scenes     map[string]*scene.Scene         `yaml:"-"`
}
