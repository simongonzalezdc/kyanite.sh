package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/utils"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "done <task-id>",
	Short: "✅ Mark a mission as completed in the digital grid",
	Long:  "🌸 Signal mission completion and update your digital karma",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("💫 Processing Mission Completion 💫")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		taskID := args[0]

		// Initialize components
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		engine := engine.New(repo)

		// Get task first to show what we're completing
		task, err := engine.GetTask(taskID)
		if err != nil {
			fmt.Printf("❌ Mission not found: %v\n", err)
			return
		}

		if task.Status == "completed" {
			fmt.Printf("✅ Mission '%s' is already completed!\n", task.Description)
			return
		}

		// Complete task
		err = engine.CompleteTask(taskID)
		if err != nil {
			fmt.Printf("❌ Error completing mission: %v\n", err)
			return
		}

		fmt.Println("🎊 Mission completed successfully!")
		fmt.Printf("📄 %s\n", task.Description)
		fmt.Printf("🔢 ID: %s\n", task.ID)
		fmt.Println("✅ Status: Completed")
		fmt.Println("\n💡 Pro Tip: Completed tasks boost your productivity!")
	},
}
