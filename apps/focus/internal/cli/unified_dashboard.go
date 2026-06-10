package cli

import (
	"fmt"

	"github.com/kyanite/focus/internal/tui"
	"github.com/spf13/cobra"
)

var unifiedDashboardCmd = &cobra.Command{
	Use:   "unified",
	Short: "🌌 Launch unified dashboard with all features",
	Long:  "🌸 Comprehensive dashboard integrating all focus.sh features",
	Run: func(cmd *cobra.Command, args []string) {
		if err := tui.StartUnifiedDashboard(); err != nil {
			fmt.Printf("Error starting unified dashboard: %v\n", err)
		}
	},
}
