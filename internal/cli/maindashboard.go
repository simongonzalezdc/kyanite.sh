package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/tui"

	"github.com/spf13/cobra"
)

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Launch the full-featured dashboard with AI and timer",
	Long:  "Access everything in one powerful interface with tasks, Pomodoro timer, AI chat, and analytics",
	Run: func(cmd *cobra.Command, args []string) {
		// Launch the full-featured Main TUI with synthwave styling
		fmt.Println("🚀 Initializing Synthwave Mission Matrix...")
		fmt.Println("⚡ Loading AI protocols and visual enhancement systems...")

		// Start the full-featured dashboard
		err := tui.StartMainDashboard([]tui.DashboardTask{})
		if err != nil {
			fmt.Printf("Error starting mission matrix: %v\n", err)
		}
	},
}

func init() {
	// dashboardCmd is now registered in root.go to avoid duplication
}
