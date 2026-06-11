package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
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

var editCmd = &cobra.Command{
	Use:   "edit [task-id]",
	Short: "Edit an existing task",
	Long: `Edit a task's description, priority, deadline, categories, or notes.
Example: focus edit abc123`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showEditError("Please provide a task ID")
			return
		}

		taskID := args[0]

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		// Get the existing task
		task, err := engine.GetTask(taskID)
		if err != nil {
			showEditError(fmt.Sprintf("Task not found: %s", taskID))
			return
		}

		// Show current task details
		showCurrentTask(task)

		// Interactive form for editing
		var (
			description string
			priority    string
			deadline    string
			categories  string
			notes       string
			status      string
			useAI       bool
		)

		// Prepare default values
		deadlineStr := ""
		if !task.Deadline.IsZero() {
			deadlineStr = task.Deadline.Format("2006-01-02")
		}

		categoriesStr := ""
		if len(task.Categories) > 0 {
			categoriesStr = ""
			for i, cat := range task.Categories {
				if i > 0 {
					categoriesStr += ", "
				}
				categoriesStr += cat
			}
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Task Description").
					Value(&description).
					Placeholder(task.Description),

				huh.NewSelect[string]().
					Title("Priority").
					Options(
						huh.NewOption("Low", "low"),
						huh.NewOption("Medium", "medium"),
						huh.NewOption("High", "high"),
					).
					Value(&priority),

				huh.NewSelect[string]().
					Title("Status").
					Options(
						huh.NewOption("Pending", "pending"),
						huh.NewOption("Completed", "completed"),
					).
					Value(&status),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Deadline (YYYY-MM-DD, leave empty to clear)").
					Value(&deadline).
					Placeholder(deadlineStr),

				huh.NewInput().
					Title("Categories (comma-separated)").
					Value(&categories).
					Placeholder(categoriesStr),

				huh.NewText().
					Title("Notes (optional)").
					Value(&notes).
					Placeholder(task.Notes).
					Lines(3),
			),
			huh.NewGroup(
				huh.NewConfirm().
					Title("Use AI to enhance description?").
					Value(&useAI),
			),
		)

		if err := form.Run(); err != nil {
			showEditError("Edit cancelled")
			return
		}

		// Use existing values if not changed
		if description == "" {
			description = task.Description
		}
		if priority == "" {
			priority = task.Priority
		}
		if status == "" {
			status = task.Status
		}
		if notes == "" {
			notes = task.Notes
		}

		// Validate new description
		if err := validation.ValidateTaskDescription(description); err != nil {
			showEditError(fmt.Sprintf("Invalid description: %s", err.Error()))
			return
		}

		// Sanitize input
		description = validation.SanitizeInput(description)

		// AI-powered enhancement if requested
		if useAI && description != task.Description {
			aiManager := di.GetContainer().GetAIManager()
			ctx, cancel := withAITimeout()
			defer cancel()
			parsedTask, err := aiManager.ParseTask(ctx, description)
			if err == nil {
				description = parsedTask.Description
				if parsedTask.Priority != "" {
					priority = parsedTask.Priority
				}
				if !parsedTask.Deadline.IsZero() {
					task.Deadline = parsedTask.Deadline
				}
				if len(parsedTask.Categories) > 0 {
					task.Categories = parsedTask.Categories
				}
			}
		}

		// Parse and update deadline
		if deadline != "" {
			deadlineTime, err := time.Parse("2006-01-02", deadline)
			if err == nil {
				task.Deadline = deadlineTime
			}
		} else if deadline == "" && deadlineStr != "" {
			// User cleared the deadline
			task.Deadline = time.Time{}
		}

		// Parse and update categories
		if categories != "" {
			task.Categories = []string{}
			for _, cat := range splitAndTrim(categories) {
				if cat != "" {
					task.Categories = append(task.Categories, cat)
				}
			}
		} else if categories == "" && categoriesStr != "" {
			// User cleared categories
			task.Categories = []string{}
		}

		// Update task fields
		task.Description = description
		task.Priority = priority
		task.Status = status
		task.Notes = notes
		task.UpdatedAt = time.Now()

		if err := engine.UpdateTask(task); err != nil {
			showEditError(fmt.Sprintf("Failed to update task: %s", err.Error()))
			return
		}

		showEditSuccess(task)
	},
}

func showCurrentTask(task models.Task) {
	title := styles.SynthwaveTitle("📝 EDITING TASK")
	fmt.Println(title)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"ID", task.ID, styles.SynthwaveCyan},
		{"DESCRIPTION", task.Description, styles.SynthwavePink},
		{"STATUS", task.Status, styles.SynthwaveGreen},
		{"PRIORITY", task.Priority, styles.SynthwaveYellow},
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
}

func showEditError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ EDIT FAILED\n\n%s", message))
	fmt.Println(errorBox)
}

func showEditSuccess(task models.Task) {
	successTitle := styles.SynthwaveTitle("✅ TASK UPDATED")
	fmt.Println(successTitle)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"ID", task.ID, styles.SynthwaveCyan},
		{"DESCRIPTION", task.Description, styles.SynthwavePink},
		{"PRIORITY", task.Priority, styles.SynthwaveYellow},
		{"UPDATED", task.UpdatedAt.Format(time.RFC822), styles.SynthwaveGreen},
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
}

func splitAndTrim(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func init() {
	rootCmd.AddCommand(editCmd)
}
