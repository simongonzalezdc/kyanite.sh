package recurrence

import (
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestEngine_GenerateInstances(t *testing.T) {
	engine := New()

	tests := []struct {
		name          string
		task          models.Task
		until         time.Time
		expectedCount int
	}{
		{
			name: "daily recurrence for 7 days",
			task: models.Task{
				ID:                 "test-1",
				Description:        "Daily task",
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
				Deadline:           time.Now(),
			},
			until:         time.Now().AddDate(0, 0, 7),
			expectedCount: 7,
		},
		{
			name: "weekly recurrence for 4 weeks",
			task: models.Task{
				ID:                 "test-2",
				Description:        "Weekly task",
				RecurrencePattern:  models.RecurrenceWeekly,
				RecurrenceInterval: 1,
				Deadline:           time.Now(),
			},
			until:         time.Now().AddDate(0, 0, 28),
			expectedCount: 4,
		},
		{
			name: "monthly recurrence for 3 months",
			task: models.Task{
				ID:                 "test-3",
				Description:        "Monthly task",
				RecurrencePattern:  models.RecurrenceMonthly,
				RecurrenceInterval: 1,
				Deadline:           time.Now(),
			},
			until:         time.Now().AddDate(0, 3, 0),
			expectedCount: 3,
		},
		{
			name: "recurrence with end date",
			task: models.Task{
				ID:                 "test-4",
				Description:        "Limited daily task",
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
				Deadline:           time.Now(),
				RecurrenceEndDate:  time.Now().AddDate(0, 0, 3),
			},
			until:         time.Now().AddDate(0, 0, 10),
			expectedCount: 3,
		},
		{
			name: "every 2 days",
			task: models.Task{
				ID:                 "test-5",
				Description:        "Every other day",
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 2,
				Deadline:           time.Now(),
			},
			until:         time.Now().AddDate(0, 0, 10),
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instances := engine.GenerateInstances(tt.task, tt.until)

			if len(instances) != tt.expectedCount {
				t.Errorf("GenerateInstances() generated %d instances, want %d", len(instances), tt.expectedCount)
			}

			// Verify all instances have parent ID set
			for i, instance := range instances {
				if instance.ParentTaskID != tt.task.ID {
					t.Errorf("Instance %d has incorrect parent ID: got %s, want %s", i, instance.ParentTaskID, tt.task.ID)
				}

				if instance.IsRecurring() {
					t.Errorf("Instance %d should not be marked as recurring", i)
				}

				if !instance.IsRecurringInstance() {
					t.Errorf("Instance %d should be marked as recurring instance", i)
				}
			}
		})
	}
}

func TestEngine_GetNextInstance(t *testing.T) {
	engine := New()

	baseTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		task          models.Task
		lastInstance  time.Time
		wantNil       bool
		checkDeadline func(time.Time) bool
	}{
		{
			name: "next daily instance",
			task: models.Task{
				ID:                 "test-1",
				Description:        "Daily task",
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
				Deadline:           baseTime,
			},
			lastInstance: baseTime,
			wantNil:      false,
			checkDeadline: func(t time.Time) bool {
				expected := baseTime.AddDate(0, 0, 1)
				return t.Year() == expected.Year() && t.Month() == expected.Month() && t.Day() == expected.Day()
			},
		},
		{
			name: "past end date",
			task: models.Task{
				ID:                 "test-2",
				Description:        "Expired task",
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
				Deadline:           baseTime,
				RecurrenceEndDate:  baseTime.AddDate(0, 0, 5),
			},
			lastInstance: baseTime.AddDate(0, 0, 6),
			wantNil:      true,
		},
		{
			name: "weekly instance",
			task: models.Task{
				ID:                 "test-3",
				Description:        "Weekly task",
				RecurrencePattern:  models.RecurrenceWeekly,
				RecurrenceInterval: 1,
				Deadline:           baseTime,
			},
			lastInstance: baseTime,
			wantNil:      false,
			checkDeadline: func(t time.Time) bool {
				expected := baseTime.AddDate(0, 0, 7)
				return t.Year() == expected.Year() && t.Month() == expected.Month() && t.Day() == expected.Day()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instance := engine.GetNextInstance(tt.task, tt.lastInstance)

			if tt.wantNil {
				if instance != nil {
					t.Error("GetNextInstance() should return nil for expired recurrence")
				}
				return
			}

			if instance == nil {
				t.Fatal("GetNextInstance() returned nil, want instance")
			}

			if instance.ParentTaskID != tt.task.ID {
				t.Errorf("Instance has incorrect parent ID: got %s, want %s", instance.ParentTaskID, tt.task.ID)
			}

			if tt.checkDeadline != nil && !tt.checkDeadline(instance.Deadline) {
				t.Errorf("Instance deadline is incorrect: got %v", instance.Deadline)
			}
		})
	}
}

func TestEngine_ShouldGenerateNext(t *testing.T) {
	engine := New()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		name             string
		task             models.Task
		lastInstanceDate time.Time
		want             bool
	}{
		{
			name: "should generate when interval passed",
			task: models.Task{
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
			},
			lastInstanceDate: today.AddDate(0, 0, -2),
			want:             true,
		},
		{
			name: "should not generate when interval not passed",
			task: models.Task{
				RecurrencePattern:  models.RecurrenceWeekly,
				RecurrenceInterval: 1,
			},
			lastInstanceDate: today.AddDate(0, 0, -3), // Next would be in 4 days (future)
			want:             false,
		},
		{
			name: "should not generate past end date",
			task: models.Task{
				RecurrencePattern:  models.RecurrenceDaily,
				RecurrenceInterval: 1,
				RecurrenceEndDate:  today.AddDate(0, 0, -5),
			},
			lastInstanceDate: today.AddDate(0, 0, -4), // Last was 4 days ago, end was 5 days ago
			want:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := engine.ShouldGenerateNext(tt.task, tt.lastInstanceDate)
			if got != tt.want {
				t.Errorf("ShouldGenerateNext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEngine_nextOccurrence(t *testing.T) {
	engine := New()

	baseTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		pattern  models.RecurrencePattern
		interval int
		expected time.Time
	}{
		{
			name:     "daily +1",
			pattern:  models.RecurrenceDaily,
			interval: 1,
			expected: baseTime.AddDate(0, 0, 1),
		},
		{
			name:     "daily +3",
			pattern:  models.RecurrenceDaily,
			interval: 3,
			expected: baseTime.AddDate(0, 0, 3),
		},
		{
			name:     "weekly +1",
			pattern:  models.RecurrenceWeekly,
			interval: 1,
			expected: baseTime.AddDate(0, 0, 7),
		},
		{
			name:     "weekly +2",
			pattern:  models.RecurrenceWeekly,
			interval: 2,
			expected: baseTime.AddDate(0, 0, 14),
		},
		{
			name:     "monthly +1",
			pattern:  models.RecurrenceMonthly,
			interval: 1,
			expected: baseTime.AddDate(0, 1, 0),
		},
		{
			name:     "monthly +3",
			pattern:  models.RecurrenceMonthly,
			interval: 3,
			expected: baseTime.AddDate(0, 3, 0),
		},
		{
			name:     "yearly +1",
			pattern:  models.RecurrenceYearly,
			interval: 1,
			expected: baseTime.AddDate(1, 0, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.nextOccurrence(baseTime, tt.pattern, tt.interval)

			if !result.Equal(tt.expected) {
				t.Errorf("nextOccurrence() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTask_IsRecurring(t *testing.T) {
	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "daily recurrence",
			task: models.Task{RecurrencePattern: models.RecurrenceDaily},
			want: true,
		},
		{
			name: "no recurrence",
			task: models.Task{RecurrencePattern: models.RecurrenceNone},
			want: false,
		},
		{
			name: "empty pattern",
			task: models.Task{RecurrencePattern: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.IsRecurring(); got != tt.want {
				t.Errorf("Task.IsRecurring() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTask_IsRecurringInstance(t *testing.T) {
	tests := []struct {
		name string
		task models.Task
		want bool
	}{
		{
			name: "has parent ID",
			task: models.Task{ParentTaskID: "parent-123"},
			want: true,
		},
		{
			name: "no parent ID",
			task: models.Task{ParentTaskID: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.task.IsRecurringInstance(); got != tt.want {
				t.Errorf("Task.IsRecurringInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}
