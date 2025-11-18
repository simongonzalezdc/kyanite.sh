package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
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
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)

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
			useAI       bool
		)

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

		// Validate new description
		if err := validation.ValidateTaskDescription(description); err != nil {
			showEditError(fmt.Sprintf("Invalid description: %s", err.Error()))
			return
		}

		// Sanitize input
		description = validation.SanitizeInput(description)

		// AI-powered enhancement if requested
		if useAI && description != task.Description {
			aiManager := ai.New()
			parsedTask, err := aiManager.ParseTask(context.Background(), description)
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

		// Update task
		task.Description = description
		task.Priority = priority

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
		Border(lipgloss.DoubleBorder()).
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

func init() {
	rootCmd.AddCommand(editCmd)
}
