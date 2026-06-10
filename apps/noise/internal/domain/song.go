package domain

import (
	"time"
)

// Song represents a complete song with metadata and content
type Song struct {
	ID         int          `json:"id" db:"id"`
	Filepath   string       `json:"filepath" db:"filepath"`
	Metadata   SongMetadata `json:"metadata" db:"metadata"`
	Sections   []Section    `json:"sections" db:"-"`
	RawContent string       `json:"-" db:"-"`
}

// SongMetadata is the YAML frontmatter
type SongMetadata struct {
	Title         string    `yaml:"title" json:"title" db:"title"`
	Artist        string    `yaml:"artist,omitempty" json:"artist,omitempty" db:"artist"`
	Key           string    `yaml:"key,omitempty" json:"key,omitempty" db:"key"`
	Tempo         int       `yaml:"tempo,omitempty" json:"tempo,omitempty" db:"tempo"`
	TimeSignature string    `yaml:"time_signature,omitempty" json:"time_signature,omitempty" db:"time_signature"`
	Structure     string    `yaml:"structure,omitempty" json:"structure,omitempty" db:"structure"`
	Tags          []string  `yaml:"tags,omitempty" json:"tags,omitempty" db:"tags"`
	CreatedAt     time.Time `yaml:"created_at" json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `yaml:"updated_at" json:"updated_at" db:"updated_at"`
}

// Section represents verse, chorus, bridge, etc.
type Section struct {
	Type   SectionType `json:"type" db:"type"`
	Number int         `json:"number" db:"number"`
	Lines  []Line      `json:"lines" db:"-"`
	Notes  string      `json:"notes,omitempty" db:"notes"`
}

// SectionType represents the type of song section
type SectionType string

const (
	// SectionVerse represents a verse section
	SectionVerse SectionType = "verse"
	// SectionChorus represents a chorus section
	SectionChorus SectionType = "chorus"
	// SectionPreChorus represents a pre-chorus section
	SectionPreChorus SectionType = "pre-chorus"
	// SectionBridge represents a bridge section
	SectionBridge SectionType = "bridge"
	// SectionOutro represents an outro section
	SectionOutro SectionType = "outro"
	// SectionIntro represents an intro section
	SectionIntro SectionType = "intro"
)

// Line is a single lyric line with analysis
type Line struct {
	Text          string   `json:"text" db:"text"`
	Syllables     int      `json:"syllables" db:"syllables"`
	RhymeScheme   string   `json:"rhyme_scheme,omitempty" db:"rhyme_scheme"`
	StressPattern string   `json:"stress_pattern,omitempty" db:"stress_pattern"`
	Flags         []string `json:"flags,omitempty" db:"flags"`
}

// QualityScore represents the 7-dimension rubric
type QualityScore struct {
	Specificity        float64           `json:"specificity" db:"specificity"`
	Originality        float64           `json:"originality" db:"originality"`
	EmotionalResonance float64           `json:"emotional_resonance" db:"emotional_resonance"`
	Prosody            float64           `json:"prosody" db:"prosody"`
	Coherence          float64           `json:"coherence" db:"coherence"`
	VoiceConsistency   float64           `json:"voice_consistency" db:"voice_consistency"`
	SurpriseFactor     float64           `json:"surprise_factor" db:"surprise_factor"`
	Total              float64           `json:"total" db:"total"`
	Feedback           map[string]string `json:"feedback" db:"feedback"`
}

// Version represents a snapshot of a song at a specific point in time
type Version struct {
	ID            int       `json:"id" db:"id"`
	SongID        int       `json:"song_id" db:"song_id"`
	Content       string    `json:"content" db:"content"`
	IsMilestone   bool      `json:"is_milestone" db:"is_milestone"`
	MilestoneName string    `json:"milestone_name,omitempty" db:"milestone_name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// WritingStats represents daily writing statistics
type WritingStats struct {
	ID               int       `json:"id" db:"id"`
	Date             time.Time `json:"date" db:"date"`
	WordsWritten     int       `json:"words_written" db:"words_written"`
	SongsCreated     int       `json:"songs_created" db:"songs_created"`
	SongsEdited      int       `json:"songs_edited" db:"songs_edited"`
	AIRequests       int       `json:"ai_requests" db:"ai_requests"`
	TimeSpentMinutes int       `json:"time_spent_minutes" db:"time_spent_minutes"`
}

// KnowledgeCard represents a knowledge base entry for RAG system
type KnowledgeCard struct {
	ID       string   `json:"id" db:"id"`
	Topic    string   `json:"topic" db:"topic"`
	Summary  string   `json:"summary" db:"summary"`
	Rules    []string `json:"rules" db:"rules"`
	Examples []string `json:"examples" db:"examples"`
	Checks   []string `json:"checks" db:"checks"`
}

// Project represents a collection of related songs
type Project struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	SongIDs     []int     `json:"song_ids" db:"song_ids"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
