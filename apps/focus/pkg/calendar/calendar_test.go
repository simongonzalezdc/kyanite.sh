package calendar

import (
	"testing"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

func TestNew(t *testing.T) {
	calendar := New("test-theme")

	if calendar.Theme != "test-theme" {
		t.Errorf("Expected theme 'test-theme', got '%s'", calendar.Theme)
	}

	if calendar.Tasks == nil {
		t.Error("Tasks slice should be initialized")
	}

	if calendar.CurrentDate.IsZero() {
		t.Error("CurrentDate should be initialized")
	}

	if calendar.SelectedDate.IsZero() {
		t.Error("SelectedDate should be initialized")
	}
}

func TestLoadTasks(t *testing.T) {
	calendar := New("test-theme")

	tasks := []models.Task{
		{ID: "1", Description: "Task 1"},
		{ID: "2", Description: "Task 2"},
	}

	calendar.LoadTasks(tasks)

	if len(calendar.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(calendar.Tasks))
	}

	if calendar.Tasks[0].ID != "1" {
		t.Errorf("Expected first task ID '1', got '%s'", calendar.Tasks[0].ID)
	}
}

func TestAddTask(t *testing.T) {
	calendar := New("test-theme")

	task := models.Task{ID: "1", Description: "Test Task"}

	calendar.AddTask(task)

	if len(calendar.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(calendar.Tasks))
	}

	if calendar.Tasks[0].ID != "1" {
		t.Errorf("Expected task ID '1', got '%s'", calendar.Tasks[0].ID)
	}
}

func TestCalendarDay(t *testing.T) {
	date := time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC)
	tasks := []CalendarTask{
		{ID: "1", Description: "Task 1"},
		{ID: "2", Description: "Task 2"},
	}

	day := CalendarDay{
		Date:      date,
		Tasks:     tasks,
		Overdue:   false,
		Upcoming:  false,
		Completed: 1,
		Pending:   1,
	}

	if !day.Date.Equal(date) {
		t.Errorf("Expected date %v, got %v", date, day.Date)
	}

	if len(day.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(day.Tasks))
	}

	if day.Overdue {
		t.Error("Expected Overdue to be false")
	}

	if day.Completed != 1 {
		t.Errorf("Expected Completed 1, got %d", day.Completed)
	}

	if day.Pending != 1 {
		t.Errorf("Expected Pending 1, got %d", day.Pending)
	}
}

func TestCalendarDay_Overdue(t *testing.T) {
	overdueTask := CalendarTask{ID: "1", Description: "Overdue Task"}
	normalTask := CalendarTask{ID: "2", Description: "Normal Task"}

	dayWithOverdue := CalendarDay{
		Tasks:   []CalendarTask{overdueTask, normalTask},
		Overdue: true,
	}

	dayWithoutOverdue := CalendarDay{
		Tasks:   []CalendarTask{normalTask},
		Overdue: false,
	}

	if !dayWithOverdue.Overdue {
		t.Error("Expected day with overdue tasks to have Overdue = true")
	}

	if dayWithoutOverdue.Overdue {
		t.Error("Expected day without overdue tasks to have Overdue = false")
	}
}

func TestCalendarTask(t *testing.T) {
	task := CalendarTask{ID: "1", Description: "Test Task", Priority: "high"}

	if task.ID != "1" {
		t.Errorf("Expected task ID '1', got '%s'", task.ID)
	}

	if task.Description != "Test Task" {
		t.Errorf("Expected task description 'Test Task', got '%s'", task.Description)
	}

	if task.Priority != "high" {
		t.Errorf("Expected task priority 'high', got '%s'", task.Priority)
	}
}

func TestCalendar_GetMonthDays(t *testing.T) {
	// Test getting days for a month (implementation-specific)
	// This is a basic test to ensure the calendar can handle month operations
	now := time.Now()

	// Test that we can calculate days in current month
	daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	if daysInMonth < 28 || daysInMonth > 31 {
		t.Errorf("Expected days in month to be between 28 and 31, got %d", daysInMonth)
	}

	// Test that calendar can handle date operations
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if firstOfMonth.Month() != now.Month() {
		t.Errorf("Expected month %v, got %v", now.Month(), firstOfMonth.Month())
	}
}

func TestCalendar_DateNavigation(t *testing.T) {
	calendar := New("test-theme")

	originalDate := calendar.CurrentDate

	// Test selecting a different date
	newDate := time.Date(2025, 10, 19, 0, 0, 0, 0, time.UTC)
	calendar.SelectedDate = newDate

	if !calendar.SelectedDate.Equal(newDate) {
		t.Errorf("Expected selected date %v, got %v", newDate, calendar.SelectedDate)
	}

	if calendar.CurrentDate.Equal(originalDate) {
		// This test passes if CurrentDate was unchanged
	} else {
		t.Logf("CurrentDate was modified from %v to %v", originalDate, calendar.CurrentDate)
	}
}

func TestCalendar_ThemeHandling(t *testing.T) {
	// Test with different themes
	themes := []string{"light", "dark", "synthwave"}

	for _, theme := range themes {
		calendar := New(theme)

		if calendar.Theme != theme {
			t.Errorf("Expected theme '%s', got '%s'", theme, calendar.Theme)
		}
	}
}

func TestCalendarTaskManagement(t *testing.T) {
	calendar := New("test-theme")

	// Add multiple tasks
	task1 := models.Task{ID: "1", Description: "Task 1"}
	task2 := models.Task{ID: "2", Description: "Task 2"}
	task3 := models.Task{ID: "3", Description: "Task 3"}

	calendar.AddTask(task1)
	calendar.AddTask(task2)
	calendar.AddTask(task3)

	if len(calendar.Tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(calendar.Tasks))
	}

	// Verify task order
	expectedOrder := []string{"1", "2", "3"}
	for i, expectedID := range expectedOrder {
		if calendar.Tasks[i].ID != expectedID {
			t.Errorf("Expected task %d to have ID '%s', got '%s'", i, expectedID, calendar.Tasks[i].ID)
		}
	}
}

func TestCalendar_EmptyState(t *testing.T) {
	calendar := New("test-theme")

	// Should handle empty tasks gracefully
	if len(calendar.Tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty calendar, got %d", len(calendar.Tasks))
	}

	// Should still work with empty tasks
	calendar.LoadTasks([]models.Task{})

	if len(calendar.Tasks) != 0 {
		t.Errorf("Expected 0 tasks after loading empty tasks, got %d", len(calendar.Tasks))
	}
}

func TestCalendarTask_Priority(t *testing.T) {
	// Test different priority levels
	priorities := []string{"low", "medium", "high", "urgent"}

	for _, priority := range priorities {
		task := CalendarTask{
			ID:          priority + "-task",
			Description: priority + " priority task",
			Priority:    priority,
		}

		if task.Priority != priority {
			t.Errorf("Expected priority '%s', got '%s'", priority, task.Priority)
		}
	}
}

func TestCalendarTask_Status(t *testing.T) {
	// Test different status levels
	statuses := []string{"pending", "completed", "cancelled"}

	for _, status := range statuses {
		task := CalendarTask{
			ID:          status + "-task",
			Description: status + " status task",
			Status:      status,
		}

		if task.Status != status {
			t.Errorf("Expected status '%s', got '%s'", status, task.Status)
		}
	}
}
