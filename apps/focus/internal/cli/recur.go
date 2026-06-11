package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/di"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/kyanite/focus/pkg/validation"
	"github.com/spf13/cobra"
)

var (
	recurPattern  string
	recurInterval int
	recurEndDate  string
	recurPriority string
)

var recurCmd = &cobra.Command{
	Use:   "recur [description]",
	Short: "Create a recurring task",
	Long: `Create a task that repeats on a regular schedule.
Examples:
  focus recur "Daily standup" --pattern=daily
  focus recur "Weekly review" --pattern=weekly --interval=1
  focus recur "Monthly report" --pattern=monthly --end="2024-12-31"
  focus recur "Quarterly planning" --pattern=monthly --interval=3`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showRecurError("Please provide a task description")
			return
		}

		description := strings.Join(args, " ")

		// Validate description
		if err := validation.ValidateTaskDescription(description); err != nil {
			showRecurError(fmt.Sprintf("Invalid description: %s", err.Error()))
			return
		}

		// Sanitize input
		description = validation.SanitizeInput(description)

		// Validate pattern
		var pattern models.RecurrencePattern
		switch strings.ToLower(recurPattern) {
		case "daily":
			pattern = models.RecurrenceDaily
		case "weekly":
			pattern = models.RecurrenceWeekly
		case "monthly":
			pattern = models.RecurrenceMonthly
		case "yearly":
			pattern = models.RecurrenceYearly
		case "":
			showRecurError("Please specify a recurrence pattern (daily, weekly, monthly, yearly)")
			return
		default:
			showRecurError("Invalid pattern. Use: daily, weekly, monthly, or yearly")
			return
		}

		// Validate interval
		if recurInterval <= 0 {
			recurInterval = 1
		}

		// Parse end date if provided
		var endDate time.Time
		if recurEndDate != "" {
			var err error
			endDate, err = time.Parse("2006-01-02", recurEndDate)
			if err != nil {
				showRecurError(fmt.Sprintf("Invalid end date format. Use YYYY-MM-DD: %s", err.Error()))
				return
			}
		}

		// Determine priority
		priority := recurPriority
		if priority == "" {
			priority = "medium"
		}
		if priority != "low" && priority != "medium" && priority != "high" {
			showRecurError("Priority must be: low, medium, or high")
			return
		}

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		taskEngine := engine.New(repo)

		// Try to enhance with AI if available
		aiManager := di.GetContainer().GetAIManager()
		parsedTask, err := aiManager.ParseTask(context.Background(), description)

		var categories []string
		deadline := time.Now() // Start recurring from today

		if err == nil {
			// AI parsing succeeded
			if len(parsedTask.Categories) > 0 {
				categories = parsedTask.Categories
			}
			if !parsedTask.Deadline.IsZero() {
				deadline = parsedTask.Deadline
			}
			if parsedTask.Priority != "" {
				priority = parsedTask.Priority
			}
		}

		// Create the recurring task
		task := models.ParsedTask{
			Description:        description,
			Priority:           priority,
			Deadline:           deadline,
			Categories:         categories,
			RecurrencePattern:  pattern,
			RecurrenceInterval: recurInterval,
			RecurrenceEndDate:  endDate,
		}

		createdTask, err := taskEngine.AddTask(task)
		if err != nil {
			showRecurError(fmt.Sprintf("Failed to create recurring task: %s", err.Error()))
			return
		}

		showRecurSuccess(createdTask, pattern, recurInterval, endDate)
	},
}

func showRecurSuccess(task models.Task, pattern models.RecurrencePattern, interval int, endDate time.Time) {
	title := styles.SynthwaveTitle("🔄 RECURRING TASK CREATED")
	fmt.Println(title)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"ID", task.ID[:8] + "...", styles.SynthwaveCyan},
		{"DESCRIPTION", task.Description, styles.SynthwavePink},
		{"PRIORITY", strings.ToUpper(task.Priority), getPriorityColor(task.Priority)},
		{"PATTERN", string(pattern), styles.SynthwaveYellow},
		{"INTERVAL", fmt.Sprintf("Every %d %s", interval, getIntervalUnit(pattern, interval)), styles.SynthwaveGreen},
	}

	if !endDate.IsZero() {
		details = append(details, struct {
			label string
			value string
			color lipgloss.Color
		}{"END DATE", endDate.Format("2006-01-02"), styles.SynthwaveRed})
	}

	for _, detail := range details {
		detailLine := lipgloss.NewStyle().
			Foreground(detail.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("▸ %s: %s", detail.label, detail.value))
		fmt.Println(detailLine)
	}
	fmt.Println()

	info := lipgloss.NewStyle().
		Foreground(styles.SynthwaveCyan).
		Background(styles.DarkVoid).
		Padding(1, 2).
		Render("💡 Recurring instances will be automatically generated based on the schedule")
	fmt.Println(info)
}

func getIntervalUnit(pattern models.RecurrencePattern, interval int) string {
	unit := ""
	switch pattern {
	case models.RecurrenceDaily:
		unit = "day"
	case models.RecurrenceWeekly:
		unit = "week"
	case models.RecurrenceMonthly:
		unit = "month"
	case models.RecurrenceYearly:
		unit = "year"
	}

	if interval > 1 {
		unit += "s"
	}

	return unit
}

func showRecurError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ RECURRING TASK FAILED\n\n%s", message))
	fmt.Println(errorBox)
}

func getPriorityColor(priority string) lipgloss.Color {
	switch priority {
	case "high":
		return styles.SynthwaveRed
	case "medium":
		return styles.SynthwaveYellow
	case "low":
		return styles.SynthwaveGreen
	default:
		return styles.SynthwaveCyan
	}
}

func init() {
	rootCmd.AddCommand(recurCmd)

	recurCmd.Flags().StringVar(&recurPattern, "pattern", "", "Recurrence pattern (daily, weekly, monthly, yearly)")
	recurCmd.Flags().IntVar(&recurInterval, "interval", 1, "Recurrence interval (e.g., every 2 weeks)")
	recurCmd.Flags().StringVar(&recurEndDate, "end", "", "End date for recurrence (YYYY-MM-DD)")
	recurCmd.Flags().StringVar(&recurPriority, "priority", "medium", "Task priority (low, medium, high)")
}
