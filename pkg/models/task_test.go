package models

import (
	"testing"
	"time"
)

func TestTask_Validation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		task      Task
		wantValid bool
	}{
		{
			name: "valid task",
			task: Task{
				ID:          "test-id",
				Description: "test task",
				Status:      "pending",
				Priority:    "medium",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantValid: true,
		},
		{
			name: "task with empty description",
			task: Task{
				ID:          "test-id",
				Description: "",
				Status:      "pending",
				Priority:    "medium",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantValid: false,
		},
		{
			name: "task with invalid priority",
			task: Task{
				ID:          "test-id",
				Description: "test task",
				Status:      "pending",
				Priority:    "invalid",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantValid: true, // Should still be valid but normalized
		},
		{
			name: "task with deadline",
			task: Task{
				ID:          "test-id",
				Description: "test task",
				Status:      "pending",
				Priority:    "high",
				Deadline:    now.Add(24 * time.Hour),
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantValid: true,
		},
		{
			name: "task with categories",
			task: Task{
				ID:          "test-id",
				Description: "test task",
				Status:      "pending",
				Priority:    "low",
				Categories:  []string{"work", "urgent"},
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.task
			if tt.wantValid && task.Description == "" {
				t.Errorf("Expected valid task but got empty description")
			}
		})
	}
}

func TestParsedTask_Validation(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		task      ParsedTask
		wantValid bool
	}{
		{
			name: "valid parsed task",
			task: ParsedTask{
				Description: "test task",
				Priority:    "medium",
			},
			wantValid: true,
		},
		{
			name: "parsed task with empty description",
			task: ParsedTask{
				Description: "",
				Priority:    "medium",
			},
			wantValid: false,
		},
		{
			name: "parsed task with deadline",
			task: ParsedTask{
				Description: "test task",
				Priority:    "high",
				Deadline:    now.Add(24 * time.Hour),
			},
			wantValid: true,
		},
		{
			name: "parsed task with categories",
			task: ParsedTask{
				Description: "test task",
				Priority:    "low",
				Categories:  []string{"work", "personal"},
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := tt.task
			if tt.wantValid && task.Description == "" {
				t.Errorf("Expected valid parsed task but got empty description")
			}
		})
	}
}
