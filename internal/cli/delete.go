package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/store"
	"github.com/kyanite/focus/pkg/utils"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "remove <task-id>",
	Short: "🗑️ Erase a mission from the digital grid",
	Long:  "🌸 Permanently remove tasks from your digital realm",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌪️ Initializing Mission Erasure Protocol 🌪️")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		taskID := args[0]

		// Initialize components
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)

		// Get task first to show what we're deleting
		task, err := engine.GetTask(taskID)
		if err != nil {
			fmt.Printf("❌ Mission not found: %v\n", err)
			return
		}

		// Delete task
		err = engine.DeleteTask(taskID)
		if err != nil {
			fmt.Printf("❌ Error erasing mission: %v\n", err)
			return
		}

		fmt.Println("✅ Mission erased successfully!")
		fmt.Printf("📋 %s\n", task.Description)
		fmt.Printf("🔢 ID: %s\n", task.ID)
		fmt.Println("\n💡 Pro Tip: Deleted tasks are gone forever - confirm your choice!")
	},
}
