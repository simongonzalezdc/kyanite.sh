package scene

import "time"

// Scene represents a story scene
type Scene struct {
	ID           string    `yaml:"id"`
	Chapter      int       `yaml:"chapter"`
	SceneNumber  int       `yaml:"scene_number"`
	Name         string    `yaml:"name"`
	POVCharacter string    `yaml:"pov_character,omitempty"`
	Location     string    `yaml:"location,omitempty"`
	TimeOfDay    string    `yaml:"time_of_day,omitempty"`
	Weather      string    `yaml:"weather,omitempty"`
	Characters   []string  `yaml:"characters,omitempty"`
	PlotPoints   []string  `yaml:"plot_points,omitempty"`
	Status       string    `yaml:"status"` // draft, revising, done
	WordCount    int       `yaml:"word_count"`
	Notes        string    `yaml:"notes,omitempty"`
	CreatedAt    time.Time `yaml:"created_at"`
	UpdatedAt    time.Time `yaml:"updated_at"`
	Content      string    `yaml:"-"` // Markdown content after frontmatter
}
