package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/focus/pkg/models"
)

// Calendar represents a calendar for task management
type Calendar struct {
	CurrentDate  time.Time
	SelectedDate time.Time
	Tasks        []models.Task
	Theme        string
}

// CalendarTask represents a task in the calendar (alias for models.Task)
type CalendarTask = models.Task

// CalendarDay represents a day with tasks
type CalendarDay struct {
	Date      time.Time
	Tasks     []CalendarTask
	Overdue   bool
	Upcoming  bool
	Completed int
	Pending   int
}

// New creates a new calendar instance
func New(theme string) *Calendar {
	return &Calendar{
		CurrentDate:  time.Now(),
		SelectedDate: time.Now(),
		Tasks:        []models.Task{},
		Theme:        theme,
	}
}

// LoadTasks loads tasks into the calendar
func (c *Calendar) LoadTasks(tasks []models.Task) {
	c.Tasks = tasks
}

// AddTask adds a task to the calendar
func (c *Calendar) AddTask(task models.Task) {
	c.Tasks = append(c.Tasks, task)
}

// GetDayDetails returns details for a specific day
func (c *Calendar) GetDayDetails(date time.Time) CalendarDay {
	day := CalendarDay{
		Date:  date,
		Tasks: []CalendarTask{},
	}

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	overdue := false
	upcoming := false
	completed := 0
	pending := 0

	for _, task := range c.Tasks {
		// Check if task is on this date (inclusive of start, exclusive of end)
		if !task.Deadline.Before(startOfDay) && task.Deadline.Before(endOfDay) {
			day.Tasks = append(day.Tasks, task)

			if task.Status == "completed" {
				completed++
			} else {
				pending++
			}

			// Check if overdue
			if task.Status != "completed" && task.Deadline.Before(time.Now()) {
				overdue = true
			}
		}

		// Check if upcoming (within next 7 days but not overdue)
		if task.Status != "completed" && task.Deadline.After(time.Now()) &&
			task.Deadline.Before(time.Now().AddDate(0, 0, 7)) {
			upcoming = true
		}
	}

	day.Overdue = overdue
	day.Upcoming = upcoming
	day.Completed = completed
	day.Pending = pending

	return day
}

// GetMonthDays returns all days in the current month
func (c *Calendar) GetMonthDays() []CalendarDay {
	year, month, _ := c.SelectedDate.Date()

	// Get first day of month
	firstDay := time.Date(year, month, 1, 0, 0, 0, 0, c.SelectedDate.Location())

	// Get last day of month
	lastDay := firstDay.AddDate(0, 1, -1)

	days := []CalendarDay{}
	for d := firstDay; !d.After(lastDay); d = d.AddDate(0, 0, 1) {
		days = append(days, c.GetDayDetails(d))
	}

	return days
}

// GetWeekDays returns days in the current week
func (c *Calendar) GetWeekDays() []CalendarDay {
	// Get start of week (Sunday)
	weekday := c.SelectedDate.Weekday()
	daysSinceSunday := int(weekday)
	startOfWeek := c.SelectedDate.AddDate(0, 0, -daysSinceSunday)

	days := []CalendarDay{}
	for i := 0; i < 7; i++ {
		day := startOfWeek.AddDate(0, 0, i)
		days = append(days, c.GetDayDetails(day))
	}

	return days
}

// GetUpcomingTasks returns upcoming tasks within next 30 days
func (c *Calendar) GetUpcomingTasks() []CalendarTask {
	var upcoming []CalendarTask
	thirtyDaysFromNow := time.Now().AddDate(0, 0, 30)

	for _, task := range c.Tasks {
		if task.Status != "completed" &&
			task.Deadline.After(time.Now()) &&
			task.Deadline.Before(thirtyDaysFromNow) {
			upcoming = append(upcoming, CalendarTask(task))
		}
	}

	return upcoming
}

// GetOverdueTasks returns overdue tasks
func (c *Calendar) GetOverdueTasks() []CalendarTask {
	var overdue []CalendarTask

	for _, task := range c.Tasks {
		if task.Status != "completed" && task.Deadline.Before(time.Now()) {
			overdue = append(overdue, CalendarTask(task))
		}
	}

	return overdue
}

// GetPriorityIndicator returns priority indicator for day
func (day CalendarDay) GetPriorityIndicator() string {
	if day.Overdue {
		return "🔴"
	}
	if day.Pending > 0 {
		return "🟡"
	}
	if day.Completed > 0 {
		return "🟢"
	}
	return "⚪"
}

// Renderer handles calendar rendering
type Renderer struct {
	theme  string
	width  int
	height int
}

// NewRenderer creates a new calendar renderer
func NewRenderer(theme string, width, height int) *Renderer {
	return &Renderer{
		theme:  theme,
		width:  width,
		height: height,
	}
}

// RenderMonth renders month view
func (r *Renderer) RenderMonth(cal *Calendar) string {
	days := cal.GetMonthDays()
	year, month, _ := cal.SelectedDate.Date()
	monthName := month.String()

	var builder strings.Builder

	// Header
	header := fmt.Sprintf("📅 %s %d\n", monthName, year)
	builder.WriteString(header)
	builder.WriteString(strings.Repeat("─", len(header)) + "\n")

	// Weekday headers
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for _, day := range weekdays {
		fmt.Fprintf(&builder, " %-6s", day)
	}
	builder.WriteString("\n")

	// Calendar days - organize into proper week grid
	if len(days) > 0 {
		// Add padding for first week if month doesn't start on Sunday
		firstWeekday := days[0].Date.Weekday()
		for i := 0; i < int(firstWeekday); i++ {
			builder.WriteString("       ") // 7 spaces to match weekday header width
		}
	}

	for i, day := range days {
		dayNum := day.Date.Day()
		weekday := day.Date.Weekday()

		// Start new week if it's Sunday and not the first day
		if i > 0 && weekday == time.Sunday {
			builder.WriteString("\n")
		}

		// Determine styling based on tasks and current day
		var dayPrefix, daySuffix string
		if day.Overdue {
			dayPrefix, daySuffix = "[", "]" // Red indication
		} else if day.Pending > 0 {
			dayPrefix, daySuffix = "(", ")" // Yellow indication
		} else if day.Completed > 0 {
			dayPrefix, daySuffix = "{", "}" // Green indication
		}

		if dayNum == cal.CurrentDate.Day() {
			fmt.Fprintf(&builder, " [%2d]  ", dayNum) // 7 chars total
		} else if dayPrefix != "" {
			fmt.Fprintf(&builder, "%s%2d%s  ", dayPrefix, dayNum, daySuffix) // 7 chars total
		} else {
			fmt.Fprintf(&builder, "  %2d  ", dayNum) // 7 chars total
		}
	}

	return builder.String()
}

// RenderWeek renders week view
func (r *Renderer) RenderWeek(cal *Calendar) string {
	days := cal.GetWeekDays()

	var builder strings.Builder

	builder.WriteString("📅 Week View\n")
	builder.WriteString(strings.Repeat("─", 20) + "\n")

	for _, day := range days {
		dateStr := day.Date.Format("2006-01-02")

		fmt.Fprintf(&builder, "\n%s (%d tasks)\n", dateStr, len(day.Tasks))

		for _, task := range day.Tasks {
			priority := "⚪"
			switch task.Priority {
			case "high":
				priority = "🔴"
			case "medium":
				priority = "🟡"
			case "low":
				priority = "🟢"
			}

			status := "⏳"
			if task.Status == "completed" {
				status = "✅"
			}

			fmt.Fprintf(&builder, "  %s %s %s\n", status, priority, task.Description)
		}
	}

	return builder.String()
}

// RenderDay renders day view
func (r *Renderer) RenderDay(cal *Calendar) string {
	day := cal.GetDayDetails(cal.SelectedDate)
	dateStr := cal.SelectedDate.Format("2006-01-02")

	var builder strings.Builder

	fmt.Fprintf(&builder, "📅 %s\n", dateStr)
	builder.WriteString(strings.Repeat("─", 30) + "\n")

	fmt.Fprintf(&builder, "📊 Summary: %d completed, %d pending, %d overdue\n",
		day.Completed, day.Pending, func() int {
			if day.Overdue {
				return 1
			} else {
				return 0
			}
		}())

	if len(day.Tasks) > 0 {
		builder.WriteString("\n📋 Tasks:\n")
		for _, task := range day.Tasks {
			priority := "⚪"
			switch task.Priority {
			case "high":
				priority = "🔴"
			case "medium":
				priority = "🟡"
			case "low":
				priority = "🟢"
			}

			status := "⏳"
			if task.Status == "completed" {
				status = "✅"
			}

			fmt.Fprintf(&builder, "  %s %s %s\n", status, priority, task.Description)
			if !task.Deadline.IsZero() {
				fmt.Fprintf(&builder, "     📅 Due: %s\n", task.Deadline.Format("3:04 PM"))
			}
		}
	} else {
		builder.WriteString("\n📋 No tasks scheduled for this day\n")
	}

	return builder.String()
}

// GetPriorityColor returns color for priority
func GetPriorityColor(priority string, theme string) string {
	switch priority {
	case "high":
		return "#FF0000" // Red
	case "medium":
		return "#FFFF00" // Yellow
	case "low":
		return "#00FF00" // Green
	default:
		return "#808080" // Gray
	}
}

// GetThemeColor returns color based on theme
func GetThemeColor(element string, theme string) string {
	// Synthwave theme colors
	synthwaveColors := map[string]string{
		"background": "#0A0014",
		"primary":    "#00FFF0",
		"secondary":  "#FF71CE",
		"accent":     "#39FF14",
		"warning":    "#FF0040",
		"success":    "#00FF66",
	}

	switch theme {
	case "synthwave":
		return synthwaveColors[element]
	case "light":
		// Light theme colors
		lightColors := map[string]string{
			"background": "#FFFFFF",
			"primary":    "#000000",
			"secondary":  "#333333",
			"accent":     "#0066CC",
			"warning":    "#FF6600",
			"success":    "#009900",
		}
		return lightColors[element]
	case "plain":
		// Plain theme colors
		plainColors := map[string]string{
			"background": "#1A1A1A",
			"primary":    "#FFFFFF",
			"secondary":  "#CCCCCC",
			"accent":     "#4A90E2",
			"warning":    "#F39C12",
			"success":    "#27AE60",
		}
		return plainColors[element]
	default:
		return synthwaveColors[element]
	}
}
