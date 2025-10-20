package ai

import (
	"context"
	"testing"
)

func TestAIManager_ParseTask(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	
	manager := New()
	
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "simple task",
			input:   "Buy milk",
			wantErr: false,
		},
		{
			name:    "task with deadline",
			input:   "Finish report by Friday",
			wantErr: false,
		},
		{
			name:    "task with priority",
			input:   "Urgent: Call dentist",
			wantErr: false,
		},
		{
			name:    "complex task",
			input:   "Schedule team meeting for next Monday to discuss Q3 goals",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.ParseTask(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil && result.Description == "" {
				t.Errorf("ParseTask() returned empty description")
			}
		})
	}
}

func TestAIManager_SuggestTasks(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	
	manager := New()
	
	tests := []struct {
		name     string
		tasks    []string
		wantErr  bool
	}{
		{
			name:    "no existing tasks",
			tasks:   []string{},
			wantErr: false,
		},
		{
			name:    "few tasks",
			tasks:   []string{"Buy milk", "Call dentist"},
			wantErr: false,
		},
		{
			name:    "many tasks",
			tasks:   []string{"Buy milk", "Call dentist", "Finish report", "Schedule meeting", "Plan vacation"},
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := manager.SuggestTasks(context.Background(), tt.tasks)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuggestTasks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			// It's okay to return empty suggestions, but if we do, it should be a slice
			if result == nil {
				t.Errorf("SuggestTasks() returned nil instead of slice")
			}
		})
	}
}

func TestAIManager_ValidateResponse(t *testing.T) {
	manager := New()
	
	tests := []struct {
		name     string
		task     *ParsedTask
		wantValid bool
	}{
		{
			name: "valid task",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
			},
			wantValid: true,
		},
		{
			name: "empty description",
			task: &ParsedTask{
				Description: "",
				Priority:    "medium",
			},
			wantValid: false,
		},
		{
			name: "too short description",
			task: &ParsedTask{
				Description: "A",
				Priority:    "medium",
			},
			wantValid: false,
		},
		{
			name: "invalid priority",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "invalid",
			},
			wantValid: true, // Should be normalized
		},
		{
			name: "valid categories",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
				Categories:  []string{"work", "personal"},
			},
			wantValid: true,
		},
		{
			name: "invalid categories",
			task: &ParsedTask{
				Description: "Test task",
				Priority:    "medium",
				Categories:  []string{"a", "verylongcategorynameexceedinglimit"},
			},
			wantValid: true, // Invalid categories should be filtered out
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, gotValid := manager.validateResponse(tt.task)
			if gotValid != tt.wantValid {
				t.Errorf("validateResponse() = %v, want %v", gotValid, tt.wantValid)
			}
			
			// Check that invalid categories are filtered out
			if result != nil && tt.name == "invalid categories" {
				if len(result.Categories) != 0 {
					t.Errorf("validateResponse() should filter out invalid categories")
				}
			}
		})
	}
}

func TestAIManager_BasicParse(t *testing.T) {
	manager := New()
	
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "short input",
			input: "Buy milk",
		},
		{
			name:  "long input",
			input: "This is a very long task description that exceeds the normal limits and should be truncated appropriately by the basic parser",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := manager.basicParse(tt.input)
			if result == nil {
				t.Errorf("basicParse() returned nil")
				return
			}
			
			if result.Description == "" {
				t.Errorf("basicParse() returned empty description")
			}
			
			if result.Priority != "medium" {
				t.Errorf("basicParse() should set default priority to medium")
			}
		})
	}
}
