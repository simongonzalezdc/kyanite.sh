package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/spf13/cobra"
)

var undoCmd = &cobra.Command{
	Use:   "undo",
	Short: "Undo the last operation",
	Long:  `Undo the last task operation (add, delete, complete, or update).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalHistory == nil {
			return fmt.Errorf("command history not initialized")
		}

		if !globalHistory.CanUndo() {
			warningStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveYellow).
				Bold(true)
			fmt.Println(warningStyle.Render("⚠ Nothing to undo"))
			return nil
		}

		lastCmd := globalHistory.GetLastCommand()

		if err := globalHistory.Undo(); err != nil {
			return fmt.Errorf("failed to undo: %w", err)
		}

		successStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveGreen).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan)

		fmt.Println(successStyle.Render("✓ Undone:"), descStyle.Render(lastCmd))

		return nil
	},
}

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Redo the last undone operation",
	Long:  `Redo an operation that was previously undone.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if globalHistory == nil {
			return fmt.Errorf("command history not initialized")
		}

		if !globalHistory.CanRedo() {
			warningStyle := lipgloss.NewStyle().
				Foreground(styles.SynthwaveYellow).
				Bold(true)
			fmt.Println(warningStyle.Render("⚠ Nothing to redo"))
			return nil
		}

		nextCmd := globalHistory.GetNextCommand()

		if err := globalHistory.Redo(); err != nil {
			return fmt.Errorf("failed to redo: %w", err)
		}

		successStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveGreen).
			Bold(true)

		descStyle := lipgloss.NewStyle().
			Foreground(styles.SynthwaveCyan)

		fmt.Println(successStyle.Render("✓ Redone:"), descStyle.Render(nextCmd))

		return nil
	},
}

func init() {
	rootCmd.AddCommand(undoCmd)
	rootCmd.AddCommand(redoCmd)
}
