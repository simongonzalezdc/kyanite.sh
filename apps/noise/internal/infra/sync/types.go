package sync

import (
	"time"
)

// IdeaType represents the type of captured idea
type IdeaType string

const (
	IdeaTypeText      IdeaType = "text"
	IdeaTypeVoiceMemo IdeaType = "voice_memo"
	IdeaTypePhoto     IdeaType = "photo"
	IdeaTypeTempo     IdeaType = "tempo"
)

// CapturedIdea represents an idea captured from the PWA
type CapturedIdea struct {
	ID        string    `json:"id"`
	Type      IdeaType  `json:"type"`
	Content   string    `json:"content"`
	MediaPath string    `json:"media_path,omitempty"`
	BPM       int       `json:"bpm,omitempty"`
	DeviceID  string    `json:"device_id"`
	SongID    string    `json:"song_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	SyncedAt  time.Time `json:"synced_at"`
}

// TempoCapture represents a tap tempo capture
type TempoCapture struct {
	BPM       int       `json:"bpm"`
	Taps      []int64   `json:"taps"` // Timestamps of taps in ms
	DeviceID  string    `json:"device_id"`
	CreatedAt time.Time `json:"created_at"`
}

// SyncEvent represents an event that occurred during sync
type SyncEvent struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
	DeviceID  string      `json:"device_id"`
}

// SyncStatus represents the current sync status
type SyncStatus struct {
	Connected        bool      `json:"connected"`
	LastSync         time.Time `json:"last_sync"`
	PendingIdeas     int       `json:"pending_ideas"`
	ConnectedDevices int       `json:"connected_devices"`
}
