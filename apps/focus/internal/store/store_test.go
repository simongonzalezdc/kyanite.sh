package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestStore_New(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		filePath     string
		expectedPath string
	}{
		{
			name:         "valid file path",
			filePath:     filepath.Join(tempDir, "test.json"),
			expectedPath: filepath.Join(tempDir, "test.json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := New(tt.filePath)
			if store.filePath != tt.expectedPath {
				t.Errorf("New() = %v, want %v", store.filePath, tt.expectedPath)
			}
		})
	}
}

func TestStore_Load(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	tests := []struct {
		name      string
		setupFile func() error
		wantTasks []models.Task
		wantError bool
	}{
		{
			name: "file doesn't exist",
			setupFile: func() error {
				return nil
			},
			wantTasks: []models.Task{},
			wantError: false,
		},
		{
			name: "empty file",
			setupFile: func() error {
				return os.WriteFile(testFile, []byte{}, 0o644)
			},
			wantTasks: []models.Task{},
			wantError: false,
		},
		{
			name: "valid tasks",
			setupFile: func() error {
				tasks := []models.Task{
					{
						ID:          "1",
						Description: "test task 1",
						Status:      "pending",
						Priority:    "medium",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
					{
						ID:          "2",
						Description: "test task 2",
						Status:      "completed",
						Priority:    "high",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				data, err := json.MarshalIndent(tasks, "", "  ")
				if err != nil {
					return err
				}
				return os.WriteFile(testFile, data, 0o644)
			},
			wantTasks: []models.Task{
				{
					ID:          "1",
					Description: "test task 1",
					Status:      "pending",
					Priority:    "medium",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          "2",
					Description: "test task 2",
					Status:      "completed",
					Priority:    "high",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantError: false,
		},
		{
			name: "invalid json",
			setupFile: func() error {
				return os.WriteFile(testFile, []byte("invalid json"), 0o644)
			},
			wantTasks: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test file
			if tt.setupFile != nil {
				err := tt.setupFile()
				if err != nil {
					t.Fatalf("Failed to setup test file: %v", err)
				}
			}

			store := New(testFile)
			tasks, err := store.Load()

			if (err != nil) != tt.wantError {
				t.Errorf("Store.Load() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if len(tasks) != len(tt.wantTasks) {
					t.Errorf("Store.Load() = %v, want %v", tasks, tt.wantTasks)
					return
				}
			}
		})
	}
}

func TestStore_Save(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.json")

	store := New(testFile)

	tests := []struct {
		name      string
		tasks     []models.Task
		wantError bool
	}{
		{
			name:      "empty tasks",
			tasks:     []models.Task{},
			wantError: false,
		},
		{
			name: "single task",
			tasks: []models.Task{
				{
					ID:          "1",
					Description: "test task",
					Status:      "pending",
					Priority:    "high",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantError: false,
		},
		{
			name: "multiple tasks",
			tasks: []models.Task{
				{
					ID:          "1",
					Description: "task 1",
					Status:      "pending",
					Priority:    "low",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				{
					ID:          "2",
					Description: "task 2",
					Status:      "completed",
					Priority:    "medium",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Save(tt.tasks)

			if (err != nil) != tt.wantError {
				t.Errorf("Store.Save() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify the file was created and contains the correct data
				data, err := os.ReadFile(testFile)
				if err != nil {
					t.Errorf("Failed to read saved file: %v", err)
					return
				}

				var savedTasks []models.Task
				err = json.Unmarshal(data, &savedTasks)
				if err != nil {
					t.Errorf("Failed to unmarshal saved tasks: %v", err)
					return
				}

				if len(savedTasks) != len(tt.tasks) {
					t.Errorf("Saved %d tasks, expected %d", len(savedTasks), len(tt.tasks))
					return
				}
			}
		})
	}
}
