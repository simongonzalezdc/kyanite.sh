package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/config"
	"github.com/spf13/cobra"
)

// Enhanced config wizard that actually saves configuration
var enhancedConfigCmd = &cobra.Command{
	Use:   "enhanced-config",
	Short: "🔧 Enhanced configuration wizard with saving",
	Long:  "🌸 Configure focus.sh settings and save to file",
	Run:   enhancedConfigWizardHandler,
}

func enhancedConfigWizardHandler(cmd *cobra.Command, args []string) {
	fmt.Println("🔧 focus.sh Enhanced Configuration Wizard")
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println("⚙️ Let's configure your focus.sh experience:")
	fmt.Println()

	var configData struct {
		AIProvider    string
		Model         string
		DefaultTheme  string
		TimeFormat    string
		AutoSave      string
		Notifications bool
		Dashboard     string
		Editor        string
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("AI Provider").
				Description("Choose your preferred AI service").
				Options(
					huh.NewOption("🦙 Ollama (Local)", "ollama"),
					huh.NewOption("🌐 OpenRouter (Remote)", "openrouter"),
					huh.NewOption("🤖 OpenAI", "openai"),
				).
				Value(&configData.AIProvider),

			huh.NewSelect[string]().
				Title("Default Theme").
				Description("Select your visual preference").
				Options(
					huh.NewOption("🌌 Amber Night", "amber-night"),
					huh.NewOption("☀️ Light", "light"),
					huh.NewOption("⚪ Plain", "plain"),
				).
				Value(&configData.DefaultTheme),

			huh.NewSelect[string]().
				Title("Time Format").
				Description("How time should be displayed").
				Options(
					huh.NewOption("12-hour (e.g., 2:30 PM)", "12h"),
					huh.NewOption("24-hour (e.g., 14:30)", "24h"),
				).
				Value(&configData.TimeFormat),

			huh.NewInput().
				Title("AI Model").
				Placeholder("e.g., llama3, gpt-4").
				Value(&configData.Model),

			huh.NewSelect[string]().
				Title("Dashboard Mode").
				Description("Choose dashboard appearance").
				Options(
					huh.NewOption("🎨 Full (animated)", "animated"),
					huh.NewOption("📱 Compact", "compact"),
					huh.NewOption("⚡ Simple", "simple"),
				).
				Value(&configData.Dashboard),

			huh.NewInput().
				Title("Default Editor").
				Placeholder("e.g., vim, nano, code").
				Value(&configData.Editor),

			huh.NewConfirm().
				Title("Enable Notifications").
				Description("Show desktop notifications for reminders").
				Value(&configData.Notifications),
		),
	)

	err := form.Run()
	if err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0040")).
			Bold(true).
			Render(fmt.Sprintf("❌ Error: %v", err))

		fmt.Println(errorStyle)
		return
	}

	// Save configuration
	if err := saveEnhancedConfig(&configData); err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0040")).
			Bold(true).
			Render(fmt.Sprintf("❌ Failed to save configuration: %v", err))

		fmt.Println(errorStyle)
		return
	}

	fmt.Println()
	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true).
		Render("⚙️ Configuration completed and saved successfully!")

	fmt.Println(successStyle)
	fmt.Printf("📁 Configuration saved to: %s\n", config.GetConfigPath())
}

func saveEnhancedConfig(configData *struct {
	AIProvider    string
	Model         string
	DefaultTheme  string
	TimeFormat    string
	AutoSave      string
	Notifications bool
	Dashboard     string
	Editor        string
},
) error {
	// Load existing config
	cfg, err := config.LoadConfig()
	if err != nil {
		// Create new config if doesn't exist
		cfg = &config.Config{}
	}

	// Apply settings from wizard
	// Update AI settings
	if configData.AIProvider != "" {
		cfg.AI.Provider = configData.AIProvider
	}
	if configData.Model != "" {
		cfg.AI.Model = configData.Model
	}

	// Update theme
	if configData.DefaultTheme != "" {
		cfg.Theme = configData.DefaultTheme
	}

	// Update UI settings
	if configData.TimeFormat != "" {
		cfg.UI.TimeFormat = configData.TimeFormat
	}
	cfg.UI.Notifications = configData.Notifications

	// Update dashboard settings
	switch configData.Dashboard {
	case "compact":
		cfg.Dashboard.CompactMode = true
		cfg.Dashboard.ShowAnimation = false
	case "animated":
		cfg.Dashboard.ShowAnimation = true
		cfg.Dashboard.CompactMode = false
	default:
		cfg.Dashboard.ShowAnimation = false
		cfg.Dashboard.CompactMode = false
	}

	// Update notes settings
	if configData.Editor != "" {
		cfg.Notes.DefaultEditor = configData.Editor
	}

	return config.SaveConfig(cfg)
}
