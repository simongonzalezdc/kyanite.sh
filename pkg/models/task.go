package models

import "time"

// Task represents a single task in the todo list
type Task struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // pending, completed
	Priority    string    `json:"priority"` // low, medium, high
	Deadline    time.Time `json:"deadline,omitempty"`
	Categories  []string  `json:"categories,omitempty"`
	Notes       string    `json:"notes,omitempty"`        // NEW: Markdown notes
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ParsedTask represents the output from AI parsing
type ParsedTask struct {
	Description string   `json:"description"`
	Deadline    time.Time `json:"deadline,omitempty"`
	Priority    string    `json:"priority"` // low, medium, high
	Categories  []string  `json:"categories,omitempty"`
	Notes       string    `json:"notes,omitempty"`
}
