package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	bulkAllPending bool
	bulkPriority   string
	bulkCategory   string
)

var bulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Bulk operations on multiple tasks",
	Long:  `Perform bulk operations on multiple tasks at once.`,
}

var bulkCompleteCmd = &cobra.Command{
	Use:   "complete [task-ids...]",
	Short: "Complete multiple tasks",
	Long: `Complete multiple tasks at once by ID or by filter.
Examples:
  focus bulk complete abc123 def456 ghi789
  focus bulk complete --all-pending
  focus bulk complete --priority=high
  focus bulk complete --category=work`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		var tasksToComplete []models.Task
		var err error

		// Get tasks based on input
		if bulkAllPending {
			// Complete all pending tasks
			tasksToComplete, err = engine.ListTasks("active")
			if err != nil {
				showBulkError(fmt.Sprintf("Failed to load tasks: %s", err.Error()))
				return
			}
		} else if bulkPriority != "" || bulkCategory != "" {
			// Filter tasks
			allTasks, err := engine.ListTasks("all")
			if err != nil {
				showBulkError(fmt.Sprintf("Failed to load tasks: %s", err.Error()))
				return
			}
			tasksToComplete = filterBulkTasks(allTasks)
		} else if len(args) > 0 {
			// Complete specific task IDs
			for _, id := range args {
				task, err := engine.GetTask(id)
				if err != nil {
					showBulkWarning(fmt.Sprintf("Task %s not found, skipping", id))
					continue
				}
				if task.Status != "completed" {
					tasksToComplete = append(tasksToComplete, task)
				}
			}
		} else {
			showBulkError("Please provide task IDs or use --all-pending, --priority, or --category flags")
			return
		}

		if len(tasksToComplete) == 0 {
			showBulkWarning("No tasks to complete")
			return
		}

		// Complete tasks
		completed := 0
		failed := 0
		for _, task := range tasksToComplete {
			if err := engine.CompleteTask(task.ID); err != nil {
				failed++
				showBulkWarning(fmt.Sprintf("Failed to complete task %s: %s", task.ID[:8], err.Error()))
			} else {
				completed++
			}
		}

		showBulkCompleteSuccess(completed, failed)
	},
}

var bulkDeleteCmd = &cobra.Command{
	Use:   "delete [task-ids...]",
	Short: "Delete multiple tasks",
	Long: `Delete multiple tasks at once by ID or by filter.
Examples:
  focus bulk delete abc123 def456 ghi789
  focus bulk delete --priority=low
  focus bulk delete --category=old`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		var tasksToDelete []models.Task

		// Get tasks based on input
		if bulkPriority != "" || bulkCategory != "" {
			// Filter tasks
			allTasks, err := engine.ListTasks("all")
			if err != nil {
				showBulkError(fmt.Sprintf("Failed to load tasks: %s", err.Error()))
				return
			}
			tasksToDelete = filterBulkTasks(allTasks)
		} else if len(args) > 0 {
			// Delete specific task IDs
			for _, id := range args {
				task, err := engine.GetTask(id)
				if err != nil {
					showBulkWarning(fmt.Sprintf("Task %s not found, skipping", id))
					continue
				}
				tasksToDelete = append(tasksToDelete, task)
			}
		} else {
			showBulkError("Please provide task IDs or use --priority or --category flags")
			return
		}

		if len(tasksToDelete) == 0 {
			showBulkWarning("No tasks to delete")
			return
		}

		// Show warning
		showBulkDeleteWarning(tasksToDelete)

		// Confirm deletion
		fmt.Println()
		confirmStyle := lipgloss.Style{}.
			Foreground(styles.SynthwaveYellow).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1)
		fmt.Println(confirmStyle.Render("⚠️  Type 'yes' to confirm deletion:"))

		var confirm string
		_, err := fmt.Scanln(&confirm)
		if err != nil {
			showBulkError(fmt.Sprintf("Failed to read confirmation: %s", err.Error()))
			return
		}

		if strings.ToLower(confirm) != "yes" {
			showBulkWarning("Deletion cancelled")
			return
		}

		// Delete tasks
		deleted := 0
		failed := 0
		for _, task := range tasksToDelete {
			if err := engine.DeleteTask(task.ID); err != nil {
				failed++
				showBulkWarning(fmt.Sprintf("Failed to delete task %s: %s", task.ID[:8], err.Error()))
			} else {
				deleted++
			}
		}

		showBulkDeleteSuccess(deleted, failed)
	},
}

func filterBulkTasks(tasks []models.Task) []models.Task {
	filtered := make([]models.Task, 0)

	for _, task := range tasks {
		// Skip completed tasks for bulk operations (unless explicitly targeting them)
		if task.Status == "completed" && !bulkAllPending {
			continue
		}

		// Priority filter
		if bulkPriority != "" {
			if !strings.EqualFold(task.Priority, bulkPriority) {
				continue
			}
		}

		// Category filter
		if bulkCategory != "" {
			hasCategory := false
			for _, cat := range task.Categories {
				if strings.EqualFold(cat, bulkCategory) {
					hasCategory = true
					break
				}
			}
			if !hasCategory {
				continue
			}
		}

		filtered = append(filtered, task)
	}

	return filtered
}

func showBulkCompleteSuccess(completed, failed int) {
	successTitle := styles.SynthwaveTitle("✅ BULK COMPLETE")
	fmt.Println(successTitle)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"COMPLETED", fmt.Sprintf("%d tasks", completed), styles.SynthwaveGreen},
		{"FAILED", fmt.Sprintf("%d tasks", failed), styles.SynthwaveRed},
	}

	for _, detail := range details {
		detailLine := lipgloss.Style{}.
			Foreground(detail.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("▸ %s: %s", detail.label, detail.value))
		fmt.Println(detailLine)
	}
	fmt.Println()
}

func showBulkDeleteSuccess(deleted, failed int) {
	successTitle := styles.SynthwaveTitle("🗑️  BULK DELETE")
	fmt.Println(successTitle)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"DELETED", fmt.Sprintf("%d tasks", deleted), styles.SynthwaveGreen},
		{"FAILED", fmt.Sprintf("%d tasks", failed), styles.SynthwaveRed},
	}

	for _, detail := range details {
		detailLine := lipgloss.Style{}.
			Foreground(detail.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("▸ %s: %s", detail.label, detail.value))
		fmt.Println(detailLine)
	}
	fmt.Println()
}

func showBulkDeleteWarning(tasks []models.Task) {
	warningBox := lipgloss.Style{}.
		Foreground(styles.SynthwaveYellow).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveYellow).
		Render(fmt.Sprintf("⚠️  WARNING: You are about to delete %d task(s)\n\nThis action cannot be undone!", len(tasks)))
	fmt.Println(warningBox)
	fmt.Println()

	// Show tasks to be deleted
	for i, task := range tasks {
		if i >= 5 {
			remaining := len(tasks) - 5
			moreStyle := lipgloss.Style{}.
				Foreground(styles.SynthwaveCyan).
				Background(styles.DarkVoid).
				Padding(0, 1).
				Render(fmt.Sprintf("... and %d more task(s)", remaining))
			fmt.Println(moreStyle)
			break
		}
		taskLine := lipgloss.Style{}.
			Foreground(styles.SynthwavePink).
			Background(styles.DarkVoid).
			Padding(0, 1).
			Render(fmt.Sprintf("  • [%s] %s", task.ID[:8], task.Description))
		fmt.Println(taskLine)
	}
}

func showBulkWarning(message string) {
	warningStyle := lipgloss.Style{}.
		Foreground(styles.SynthwaveYellow).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Render(fmt.Sprintf("⚠️  %s", message))
	fmt.Println(warningStyle)
}

func showBulkError(message string) {
	errorBox := lipgloss.Style{}.
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ BULK OPERATION FAILED\n\n%s", message))
	fmt.Println(errorBox)
}

func init() {
	rootCmd.AddCommand(bulkCmd)
	bulkCmd.AddCommand(bulkCompleteCmd)
	bulkCmd.AddCommand(bulkDeleteCmd)

	// Add flags for bulk complete
	bulkCompleteCmd.Flags().BoolVar(&bulkAllPending, "all-pending", false, "Complete all pending tasks")
	bulkCompleteCmd.Flags().StringVar(&bulkPriority, "priority", "", "Filter by priority (low, medium, high)")
	bulkCompleteCmd.Flags().StringVar(&bulkCategory, "category", "", "Filter by category")

	// Add flags for bulk delete
	bulkDeleteCmd.Flags().StringVar(&bulkPriority, "priority", "", "Filter by priority (low, medium, high)")
	bulkDeleteCmd.Flags().StringVar(&bulkCategory, "category", "", "Filter by category")
}
