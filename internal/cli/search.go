package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	searchPriority string
	searchCategory string
	searchDeadline string
	searchStatus   string
	searchSort     string
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search and filter tasks",
	Long: `Search tasks by description and filter by priority, category, deadline, or status.
Examples:
  focus search "meeting"
  focus search --priority=high
  focus search --category=work --status=pending
  focus search --deadline=today
  focus search "report" --sort=priority`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		// Get all tasks
		tasks, err := engine.ListTasks("all")
		if err != nil {
			showSearchError(fmt.Sprintf("Failed to load tasks: %s", err.Error()))
			return
		}

		// Extract search query
		var query string
		if len(args) > 0 {
			query = strings.ToLower(strings.Join(args, " "))
		}

		// Apply filters
		filtered := filterTasks(tasks, query)

		// Sort results
		sortTasks(filtered)

		// Display results
		showSearchResults(filtered, query)
	},
}

func filterTasks(tasks []models.Task, query string) []models.Task {
	filtered := make([]models.Task, 0, len(tasks))

	for _, task := range tasks {
		// Text search filter
		if query != "" {
			descLower := strings.ToLower(task.Description)
			if !strings.Contains(descLower, query) {
				continue
			}
		}

		// Priority filter
		if searchPriority != "" {
			if !strings.EqualFold(task.Priority, searchPriority) {
				continue
			}
		}

		// Status filter
		if searchStatus != "" {
			if !strings.EqualFold(task.Status, searchStatus) {
				continue
			}
		}

		// Category filter
		if searchCategory != "" {
			hasCategory := false
			for _, cat := range task.Categories {
				if strings.EqualFold(cat, searchCategory) {
					hasCategory = true
					break
				}
			}
			if !hasCategory {
				continue
			}
		}

		// Deadline filter
		if searchDeadline != "" {
			if !matchesDeadlineFilter(task, searchDeadline) {
				continue
			}
		}

		// Task passed all filters
		filtered = append(filtered, task)
	}

	return filtered
}

func matchesDeadlineFilter(task models.Task, filter string) bool {
	if task.Deadline.IsZero() {
		return false
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	endOfWeek := today.AddDate(0, 0, 7)

	switch strings.ToLower(filter) {
	case "today":
		return !task.Deadline.Before(today) && task.Deadline.Before(tomorrow)
	case "tomorrow":
		nextDay := today.AddDate(0, 0, 1)
		dayAfter := today.AddDate(0, 0, 2)
		return !task.Deadline.Before(nextDay) && task.Deadline.Before(dayAfter)
	case "week":
		return !task.Deadline.Before(today) && task.Deadline.Before(endOfWeek)
	case "overdue":
		return task.Deadline.Before(now) && task.Status != "completed"
	case "past":
		return task.Deadline.Before(now)
	case "future":
		return task.Deadline.After(now)
	default:
		return true
	}
}

func sortTasks(tasks []models.Task) {
	if searchSort == "" {
		return
	}

	sortBy := strings.ToLower(searchSort)

	sort.Slice(tasks, func(i, j int) bool {
		switch sortBy {
		case "priority":
			return comparePriority(tasks[i].Priority, tasks[j].Priority)
		case "deadline":
			// Tasks with no deadline go to the end
			if tasks[i].Deadline.IsZero() {
				return false
			}
			if tasks[j].Deadline.IsZero() {
				return true
			}
			return tasks[i].Deadline.Before(tasks[j].Deadline)
		case "created":
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		case "updated":
			return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
		case "status":
			return tasks[i].Status < tasks[j].Status
		default:
			return tasks[i].CreatedAt.After(tasks[j].CreatedAt) // Default: newest first
		}
	})
}

func comparePriority(a, b string) bool {
	priorityValue := map[string]int{
		"high":   3,
		"medium": 2,
		"low":    1,
	}
	return priorityValue[a] > priorityValue[b]
}

func showSearchResults(tasks []models.Task, query string) {
	title := styles.SynthwaveTitle("🔍 SEARCH RESULTS")
	fmt.Println(title)
	fmt.Println()

	// Show search criteria
	var criteria strings.Builder
	if query != "" {
		criteria.WriteString(fmt.Sprintf("Query: \"%s\"", query))
	}
	if searchPriority != "" {
		if criteria.Len() > 0 {
			criteria.WriteString(" | ")
		}
		criteria.WriteString(fmt.Sprintf("Priority: %s", searchPriority))
	}
	if searchCategory != "" {
		if criteria.Len() > 0 {
			criteria.WriteString(" | ")
		}
		criteria.WriteString(fmt.Sprintf("Category: %s", searchCategory))
	}
	if searchStatus != "" {
		if criteria.Len() > 0 {
			criteria.WriteString(" | ")
		}
		criteria.WriteString(fmt.Sprintf("Status: %s", searchStatus))
	}
	if searchDeadline != "" {
		if criteria.Len() > 0 {
			criteria.WriteString(" | ")
		}
		criteria.WriteString(fmt.Sprintf("Deadline: %s", searchDeadline))
	}

	if criteria.Len() > 0 {
		criteriaStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan).
			Background(styles.DarkVoid).
			Padding(0, 1).
			Render(criteria.String())
		fmt.Println(criteriaStyle)
		fmt.Println()
	}

	// Show count
	countStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwaveYellow).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(0, 1).
		Render(fmt.Sprintf("Found %d task(s)", len(tasks)))
	fmt.Println(countStyle)
	fmt.Println()

	if len(tasks) == 0 {
		noResultsStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveRed).
			Background(styles.DarkVoid).
			Padding(1, 2).
			Render("No tasks match your search criteria")
		fmt.Println(noResultsStyle)
		return
	}

	// Display tasks
	for i, task := range tasks {
		displaySearchTask(task, i+1)
	}
}

func displaySearchTask(task models.Task, index int) {
	// Priority indicator
	var priorityIcon string
	switch task.Priority {
	case "high":
		priorityIcon = "🔴"
	case "medium":
		priorityIcon = "🟡"
	case "low":
		priorityIcon = "🟢"
	default:
		priorityIcon = "⚪"
	}

	// Status indicator
	statusIcon := "⏳"
	if task.Status == "completed" {
		statusIcon = "✅"
	}

	// Build task line
	taskLine := fmt.Sprintf("%d. %s %s [%s] %s",
		index,
		statusIcon,
		priorityIcon,
		task.ID[:8],
		task.Description,
	)

	taskStyle := lipgloss.NewStyle().
		Foreground(styles.SynthwavePink).
		Background(styles.DarkVoid).
		Padding(0, 1)

	fmt.Println(taskStyle.Render(taskLine))

	// Show metadata
	var metadata strings.Builder
	metadata.WriteString("   ")

	if !task.Deadline.IsZero() {
		deadlineStr := task.Deadline.Format("Jan 02, 2006")
		deadlineColor := styles.SynthwaveCyan
		if task.Deadline.Before(time.Now()) && task.Status != "completed" {
			deadlineColor = styles.SynthwaveRed
			deadlineStr += " (OVERDUE)"
		}
		deadlineStyle := lipgloss.NewStyle().
			Foreground(deadlineColor).
			Render("📅 " + deadlineStr)
		metadata.WriteString(deadlineStyle)
		metadata.WriteString(" ")
	}

	if len(task.Categories) > 0 {
		catStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveYellow).
			Render("🏷️  " + strings.Join(task.Categories, ", "))
		metadata.WriteString(catStyle)
	}

	if metadata.Len() > 3 {
		fmt.Println(metadata.String())
	}

	fmt.Println()
}

func showSearchError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ SEARCH FAILED\n\n%s", message))
	fmt.Println(errorBox)
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Add flags
	searchCmd.Flags().StringVar(&searchPriority, "priority", "", "Filter by priority (low, medium, high)")
	searchCmd.Flags().StringVar(&searchCategory, "category", "", "Filter by category")
	searchCmd.Flags().StringVar(&searchStatus, "status", "", "Filter by status (pending, completed)")
	searchCmd.Flags().StringVar(&searchDeadline, "deadline", "", "Filter by deadline (today, tomorrow, week, overdue, past, future)")
	searchCmd.Flags().StringVar(&searchSort, "sort", "", "Sort results (priority, deadline, created, updated, status)")
}
