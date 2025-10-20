package cli

import (
	"fmt"
	"github.com/kyanite/focus/internal/tui"
	"github.com/spf13/cobra"
	pterm "github.com/pterm/pterm"
	"github.com/charmbracelet/lipgloss"
)

var (
	neonBlue    = lipgloss.Color("#00FFFF")
	neonGreen   = lipgloss.Color("#00FF66")
	neonPurple  = lipgloss.Color("#B967C7")
	accentColor = lipgloss.Color("#FFC0CB")

	headerStyle = lipgloss.NewStyle().
		Foreground(neonBlue).
		Bold(true).
		AlignHorizontal(lipgloss.Center)





	boxStyle = lipgloss.NewStyle().
		Foreground(accentColor).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(neonBlue).
		Padding(1).
		Width(80)
)

var rootCmd = &cobra.Command{
	Use:   "neon",
	Short: "🌈 NEON Focus - AI-Powered Cyberpunk Task Manager",
	Long: boxStyle.Render(
		headerStyle.Render("🌌 W E L C O M E   T O   N E O N   F O C U S 🌌") + "\n\n" +
		lipgloss.NewStyle().Foreground(neonGreen).Render(
			"🚀 Launching immersive TUI dashboard experience\n" +
			"✨ Cyberpunk aesthetics meet intelligent productivity\n" +
			"🌈 Synthwave-inspired terminal magic\n\n" +
			"Controls: A(dd) | C(omplete) | F(ocus) | T(heme) | ?(help)") + "\n" +
		lipgloss.NewStyle().Foreground(neonBlue).Render(
			"────────────────────────────────────────────────────────────────────────────────") + "\n" +
		lipgloss.NewStyle().Foreground(neonPurple).Render(
			"💡 Pro Tip: Press '?' for help once inside the dashboard")),
	Run: func(cmd *cobra.Command, args []string) {
		// Default to TUI dashboard - TUI-FIRST approach
		fmt.Println("🚀 Initializing Synthwave Mission Matrix...")
		fmt.Println("⚡ Loading AI protocols and visual enhancement systems...")
		
		// Start the full-featured dashboard
		err := tui.StartMainDashboard([]tui.DashboardTask{})
		if err != nil {
			fmt.Printf("Error starting mission matrix: %v\n", err)
		}
	},
}

func Execute() error {
	// Setup initial appearance
	pterm.Println(pterm.FgGreen.Sprintf("🎮 NEON Focus initialized"))

	// Execute the root command
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(dashboardCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completeCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(chatCmd)
	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(notesCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.AddCommand(filterCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(simpleInteractiveCmd)
	rootCmd.AddCommand(testGumCmd)
	rootCmd.AddCommand(themeCmd)
	rootCmd.AddCommand(wizardCmd)
	rootCmd.AddCommand(configWizardCmd)
	rootCmd.AddCommand(editWizardCmd)
	
	// Add viper config commands
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configPathCmd)
	rootCmd.AddCommand(enhancedConfigCmd)
	rootCmd.AddCommand(unifiedDashboardCmd)
	rootCmd.AddCommand(mcpServerCmd)

	// Setup pterm styles
	pterm.Info.Prefix = pterm.Prefix{
		Text:  "💡",
		Style: pterm.Info.Prefix.Style,
	}
	pterm.Error.Prefix = pterm.Prefix{
		Text:  "❌",
		Style: pterm.Error.Prefix.Style,
	}
}
