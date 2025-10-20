package cli

// summaryCmd is temporarily commented out to fix unused variable warning
// Will be re-enabled when we implement proper command registration
/*
var summaryCmd = &cobra.Command{
	Use:   "report",
	Short: "📊 Get an AI-generated mission summary",
	Long:  "🌸 Receive insights about your productivity patterns from the digital realm",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("📊 Analyzing Mission Grid 📊")
		fmt.Println("~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~")

		// Initialize components
		store := store.New(utils.GetStoragePath())
		engine := engine.New(store)
		aiManager := ai.New()

		// Get all tasks
		tasks, err := engine.ListTasks("all")
		if err != nil {
			fmt.Printf("❌ Error loading missions: %v\n", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("😴 No missions found in the grid.")
			fmt.Println("📝 Upload missions with: neon add \"Your mission here\"")
			return
		}

		// Create task descriptions for AI
		taskDescriptions := make([]string, len(tasks))
		for i, task := range tasks {
			status := "pending"
			if task.Status == "completed" {
				status = "completed"
			}
			taskDescriptions[i] = fmt.Sprintf("%s (%s)", task.Description, status)
		}

		// Get AI summary
		fmt.Println("🧠 Processing grid analytics...")
		ctx := cmd.Context()
		summary, err := aiManager.SummarizeTasks(ctx, taskDescriptions)
		if err != nil {
			fmt.Printf("❌ Error generating summary: %v\n", err)
			return
		}

		fmt.Println("🌸 Mission Grid Summary:")
		fmt.Println(summary)
		fmt.Println("\n💡 Pro Tip: Regular summaries help optimize your workflow!")
	},
}
*/
