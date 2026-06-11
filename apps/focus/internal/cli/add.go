package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/validation"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [task description]",
	Short: "Add a new mission to the task board",
	Long: `Create a new mission with AI-powered parsing and beautiful Kyanite theming.
Example: focus add "Complete the project by Friday"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showAddError()
			return
		}

		description := strings.Join(args, " ")

		// Validate input
		if err := validation.ValidateTaskDescription(description); err != nil {
			showValidationError(err.Error())
			return
		}

		// Sanitize input
		description = validation.SanitizeInput(description)

		// Show processing animation
		showProcessingAnimation(description)

		// Initialize components
		engine, aiManager := initEngineAndAI()

		// AI-powered task parsing
		var task models.ParsedTask
		ctx, cancel := withAITimeout()
		defer cancel()
		parsedTask, err := aiManager.ParseTask(ctx, description)
		if err != nil {
			// Fallback to basic task creation if AI fails
			fmt.Println("⚠️ AI unavailable - using heuristic priority assignment")
			task = models.ParsedTask{
				Description: description,
				Priority:    "medium",
			}
		} else {
			task = models.ParsedTask{
				Description: parsedTask.Description,
				Priority:    parsedTask.Priority,
				Deadline:    parsedTask.Deadline,
				Categories:  parsedTask.Categories,
			}
		}

		// Add the task
		addedTask, err := engine.AddTask(task)
		if err != nil {
			showAddError()
			return
		}

		// Show success with maximum impact
		showAddSuccess(addedTask)
	},
}

func showProcessingAnimation(_ string) {
	fmt.Println(styles.LoadingMessage())
	fmt.Println("Processing task...")
}

func showValidationError(message string) {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ VALIDATION ERROR\n\n%s", message))
	fmt.Println(errorBox)
}

func showAddError() {
	errorBox := lipgloss.NewStyle().
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render("❌ MISSION CREATION FAILED\n\nPlease provide a task description:\nfocus add \"Your mission here\"")
	fmt.Println(errorBox)
}

func showAddSuccess(task models.Task) {
	// Epic success message
	successTitle := styles.SynthwaveTitle("🎯 MISSION ACCEPTED")
	fmt.Println(successTitle)
	fmt.Println()

	// Task details with cyber styling
	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"MISSION ID", task.ID, styles.SynthwaveCyan},
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

	// Clean confirmation message
	confirmMsg := styles.Title("Task added successfully")
	fmt.Println(confirmMsg)
}

func init() {
	rootCmd.AddCommand(addCmd)
}
