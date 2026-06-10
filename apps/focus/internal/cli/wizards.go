package cli

import (
	"github.com/kyanite/focus/pkg/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/wizards"
	"github.com/spf13/cobra"
)

var wizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "🧙‍♂️ Advanced task creation wizard with Huh forms",
	Long:  "🌸 Create tasks using advanced interactive forms with validation",
	Run:   wizardHandler,
}

var editWizardCmd = &cobra.Command{
	Use:   "edit-wizard [task_id]",
	Short: "✏️ Edit task using advanced wizard",
	Long:  "🌸 Edit existing tasks with comprehensive forms",
	Args:  cobra.ExactArgs(1),
	Run:   editWizardHandler,
}

func wizardHandler(cmd *cobra.Command, args []string) {
	fmt.Println("🧙‍♂️ focus.sh Advanced Task Creation Wizard")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("📋 Please fill out the form to create your task:")
	fmt.Println()

	err := wizards.TaskCreationWizard()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render(fmt.Sprintf("❌ Error: %v", err))

		fmt.Println(errorStyle)
		return
	}

	fmt.Println()
	successStyle := lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true).
		Render("🎉 Task creation completed successfully!")

	fmt.Println(successStyle)
}

func editWizardHandler(cmd *cobra.Command, args []string) {
	taskID := args[0]

	fmt.Println("✏️ focus.sh Task Edit Wizard")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("📝 Editing task: %s\n", taskID)
	fmt.Println()

	// Here you would typically:
	// 1. Load the task from storage
	// 2. Pass it to the edit wizard
	// 3. Save the changes

	// For now, show a placeholder
	placeholderStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Italic(true).
		Render("🚧 Task loading and saving will be implemented in the next phase")

	fmt.Println(placeholderStyle)
}
