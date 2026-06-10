package cli

import (
	"github.com/kyanite/focus/pkg/styles"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/calendar"
	"github.com/kyanite/focus/pkg/config"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "📅 Calendar view for task management",
	Long:  "🌸 View and manage tasks in calendar format",
}

var calendarShowCmd = &cobra.Command{
	Use:   "show [month|week|day]",
	Short: "📅 Show calendar view",
	Long:  "🌸 Display calendar with month/week/day view",
	Args:  cobra.MaximumNArgs(1),
	Run:   calendarShowHandler,
}

var calendarTodayCmd = &cobra.Command{
	Use:   "today",
	Short: "📅 Show today's tasks",
	Long:  "🌸 Display all tasks scheduled for today",
	Run:   calendarTodayHandler,
}

var calendarAddCmd = &cobra.Command{
	Use:   "add [task] [date]",
	Short: "📅 Add task with date",
	Long:  "🌸 Add a task with a specific date to the calendar",
	Args:  cobra.MinimumNArgs(1),
	Run:   calendarAddHandler,
}

var calendarListCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "📅 List calendar tasks",
	Long:  "🌸 List tasks by date or filter criteria",
	Args:  cobra.MaximumNArgs(1),
	Run:   calendarListHandler,
}

var calendarNavCmd = &cobra.Command{
	Use:   "navigate [date]",
	Short: "📅 Navigate to specific date",
	Long:  "🌸 Jump to a specific date in the calendar",
	Args:  cobra.MaximumNArgs(1),
	Run:   calendarNavigateHandler,
}

func init() {
	// Add calendar commands to root
	rootCmd.AddCommand(calendarCmd)
	calendarCmd.AddCommand(calendarShowCmd)
	calendarCmd.AddCommand(calendarTodayCmd)
	calendarCmd.AddCommand(calendarAddCmd)
	calendarCmd.AddCommand(calendarListCmd)
	calendarCmd.AddCommand(calendarNavCmd)
}

func calendarShowHandler(cmd *cobra.Command, args []string) {
	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{Theme: "synthwave"}
	}
	theme := cfg.Theme
	if theme == "" {
		theme = "synthwave"
	}

	// Get all tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error loading tasks: %v", err))

		fmt.Println(errorStyle)
		return
	}

	// Create calendar
	cal := calendar.New(theme)
	cal.LoadTasks(tasks)

	// Get view type
	viewType := "month"
	if len(args) > 0 {
		viewType = args[0]
	}

	// Create renderer
	renderer := calendar.NewRenderer(theme, 80, 20)

	// Render appropriate view
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Render("📅 focus.sh CALENDAR")

	fmt.Println(headerStyle)
	fmt.Println(strings.Repeat("─", 40))

	switch viewType {
	case "month":
		fmt.Println(renderer.RenderMonth(cal))
	case "week":
		fmt.Println(renderer.RenderWeek(cal))
	case "day":
		fmt.Println(renderer.RenderDay(cal))
	default:
		fmt.Printf("❌ Unknown view: %s. Use month, week, or day.\n", viewType)
	}
}

func calendarTodayHandler(cmd *cobra.Command, args []string) {
	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{Theme: "synthwave"}
	}
	theme := cfg.Theme
	if theme == "" {
		theme = "synthwave"
	}

	// Get all tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error loading tasks: %v", err))

		fmt.Println(errorStyle)
		return
	}

	// Create calendar for today
	cal := calendar.New(theme)
	cal.LoadTasks(tasks)
	cal.SelectedDate = time.Now()

	// Create renderer
	renderer := calendar.NewRenderer(theme, 80, 20)

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Render("📅 TODAY'S TASKS")

	fmt.Println(headerStyle)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println(renderer.RenderDay(cal))
}

func calendarAddHandler(cmd *cobra.Command, args []string) {
	description := args[0]

	// Parse date (default to today)
	targetDate := time.Now()
	if len(args) > 1 {
		parsedDate, err := parseDate(args[1])
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Foreground(styles.GetError()).
				Bold(true).
				Render(fmt.Sprintf("❌ Invalid date format: %v", err))

			fmt.Println(errorStyle)
			fmt.Println("📅 Try formats: YYYY-MM-DD, today, tomorrow, next monday, etc.")
			return
		}
		targetDate = parsedDate
	}

	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Create task with deadline
	task := struct {
		Description string
		Priority    string
		Categories  []string
		Deadline    time.Time
	}{
		Description: description,
		Priority:    "medium",
		Categories:  []string{"calendar"},
		Deadline:    targetDate,
	}

	// Create ParsedTask compatible structure
	parsedTask := models.ParsedTask{
		Description: task.Description,
		Deadline:    task.Deadline,
		Priority:    task.Priority,
		Categories:  task.Categories,
		Notes:       "",
	}

	// Add task
	newTask, err := taskEngine.AddTask(parsedTask)
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error adding task: %v", err))

		fmt.Println(errorStyle)
		return
	}

	successStyle := lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true).
		Render("✅ Task added to calendar!")

	fmt.Println(successStyle)
	fmt.Printf("📅 Task: %s\n", newTask.Description)
	fmt.Printf("📅 Date: %s\n", targetDate.Format("2006-01-02"))
	fmt.Printf("🆔 ID: %s\n", newTask.ID)
}

func calendarListHandler(cmd *cobra.Command, args []string) {
	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Get all tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error loading tasks: %v", err))

		fmt.Println(errorStyle)
		return
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Render("📅 CALENDAR TASK LIST")

	fmt.Println(headerStyle)
	fmt.Println(strings.Repeat("─", 40))

	// Group tasks by date
	taskMap := make(map[string][]models.Task)
	for _, task := range tasks {
		if !task.Deadline.IsZero() {
			dateStr := task.Deadline.Format("2006-01-02")
			taskMap[dateStr] = append(taskMap[dateStr], task)
		}
	}

	// Sort dates
	dates := make([]string, 0, len(taskMap))
	for date := range taskMap {
		dates = append(dates, date)
	}

	// Simple bubble sort for dates
	for i := 0; i < len(dates); i++ {
		for j := i + 1; j < len(dates); j++ {
			if dates[i] > dates[j] {
				dates[i], dates[j] = dates[j], dates[i]
			}
		}
	}

	// Display tasks
	for _, date := range dates {
		fmt.Printf("\n📅 %s (%d tasks)\n", date, len(taskMap[date]))
		fmt.Println(strings.Repeat("-", 20))

		for _, task := range taskMap[date] {
			status := "⏳"
			if task.Status == "completed" {
				status = "✅"
			}

			priority := "⚪"
			switch task.Priority {
			case "high":
				priority = "🔴"
			case "medium":
				priority = "🟡"
			case "low":
				priority = "🟢"
			}

			fmt.Printf("  %s %s %s\n", status, priority, task.Description)
			if len(task.Categories) > 0 {
				fmt.Printf("     🏷️  %s\n", strings.Join(task.Categories, ", "))
			}
		}
	}
}

func calendarNavigateHandler(cmd *cobra.Command, args []string) {
	// Initialize components
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	taskEngine := engine.New(repo)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{Theme: "synthwave"}
	}
	theme := cfg.Theme
	if theme == "" {
		theme = "synthwave"
	}

	// Get all tasks
	tasks, err := taskEngine.ListTasks("all")
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error loading tasks: %v", err))

		fmt.Println(errorStyle)
		return
	}

	// Parse target date
	targetDate := time.Now()
	if len(args) > 0 {
		parsedDate, err := parseDate(args[0])
		if err != nil {
			errorStyle := lipgloss.NewStyle().
				Foreground(styles.GetError()).
				Bold(true).
				Render(fmt.Sprintf("❌ Invalid date format: %v", err))

			fmt.Println(errorStyle)
			fmt.Println("📅 Try formats: YYYY-MM-DD, today, tomorrow, next monday, etc.")
			return
		}
		targetDate = parsedDate
	}

	// Create calendar
	cal := calendar.New(theme)
	cal.LoadTasks(tasks)
	cal.SelectedDate = targetDate

	// Create renderer
	renderer := calendar.NewRenderer(theme, 80, 20)

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true).
		Render(fmt.Sprintf("📅 NAVIGATED TO %s", targetDate.Format("2006-01-02")))

	fmt.Println(headerStyle)
	fmt.Println(strings.Repeat("─", 40))
	fmt.Println(renderer.RenderDay(cal))
}

// Helper function to parse date strings
func parseDate(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(strings.ToLower(dateStr))

	now := time.Now()

	// Handle natural language dates
	switch dateStr {
	case "today":
		return now, nil
	case "tomorrow":
		return now.AddDate(0, 0, 1), nil
	case "yesterday":
		return now.AddDate(0, 0, -1), nil
	case "next week":
		return now.AddDate(0, 0, 7), nil
	case "last week":
		return now.AddDate(0, 0, -7), nil
	case "next monday", "next tuesday", "next wednesday", "next thursday", "next friday", "next saturday", "next sunday":
		targetDays := map[string]time.Weekday{
			"next monday": time.Monday, "next tuesday": time.Tuesday, "next wednesday": time.Wednesday,
			"next thursday": time.Thursday, "next friday": time.Friday, "next saturday": time.Saturday, "next sunday": time.Sunday,
		}
		targetDay := targetDays[dateStr]
		currentDay := now.Weekday()
		daysUntilTarget := (int(targetDay) - int(currentDay) + 7) % 7
		if daysUntilTarget == 0 {
			daysUntilTarget = 7 // If it's the same day, go to next week
		}
		return now.AddDate(0, 0, daysUntilTarget), nil
	}

	// Try ISO date format
	if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
		return parsed, nil
	}

	// Try MM/DD format
	if parsed, err := time.Parse("01-02-2006", dateStr); err == nil {
		return parsed, nil
	}

	// Try other formats
	formats := []string{
		"2006/01/02",
		"01/02/2006",
		"Jan 2, 2006",
		"January 2, 2006",
	}

	for _, format := range formats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
