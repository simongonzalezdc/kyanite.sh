package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var subtaskCmd = &cobra.Command{
	Use:   "subtask [parent-task-id] [description]",
	Short: "Add a subtask to an existing task",
	Long:  `Create a new subtask under an existing parent task, creating a task hierarchy.`,
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		parentID := args[0]
		description := strings.Join(args[1:], " ")

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		eng := engine.New(repo)

		// Verify parent exists
		parent, err := eng.GetTask(parentID)
		if err != nil {
			return fmt.Errorf("parent task not found: %w", err)
		}

		// Get flags
		priority, _ := cmd.Flags().GetString("priority")
		categories, _ := cmd.Flags().GetStringSlice("categories")
		deadline, _ := cmd.Flags().GetString("deadline")

		// Create subtask
		subtaskID, err := eng.AddSubtask(parentID, description, priority, categories, deadline)
		if err != nil {
			return fmt.Errorf("failed to add subtask: %w", err)
		}

		subtask, _ := eng.GetTask(subtaskID)

		// Display success message
		successStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveGreen).
			Bold(true)

		taskStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan).
			Bold(true)

		parentStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwavePurple)

		fmt.Println(successStyle.Render("✓ Subtask added"))
		fmt.Printf("  %s %s\n", taskStyle.Render("ID:"), subtask.ID)
		fmt.Printf("  %s %s\n", taskStyle.Render("Description:"), subtask.Description)
		if priority != "" {
			fmt.Printf("  %s %s\n", taskStyle.Render("Priority:"), getPriorityWithColor(subtask.Priority))
		}
		fmt.Printf("  %s %s (%s)\n",
			taskStyle.Render("Parent:"),
			parent.Description,
			parentStyle.Render(parent.ID))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(subtaskCmd)
	subtaskCmd.Flags().StringP("priority", "p", "medium", "Priority level (low, medium, high)")
	subtaskCmd.Flags().StringSliceP("categories", "c", []string{}, "Categories (comma-separated)")
	subtaskCmd.Flags().StringP("deadline", "d", "", "Deadline (format: YYYY-MM-DD)")
}

func getPriorityWithColor(priority string) string {
	var color lipgloss.Color
	switch priority {
	case "high":
		color = styles.SynthwaveRed
	case "medium":
		color = styles.SynthwaveYellow
	case "low":
		color = styles.SynthwaveGreen
	default:
		color = styles.SynthwaveCyan
	}
	return lipgloss.NewStyle().Foreground(color).Render(priority)
}
