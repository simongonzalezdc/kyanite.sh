package cli

import (
	"github.com/kyanite/focus/pkg/styles"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/config"
	"github.com/spf13/cobra"
)

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "📖 Get configuration value",
	Long:  "🌸 Get a specific configuration value by key",
	Args:  cobra.ExactArgs(1),
	Run:   configGetHandler,
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "💾 Set configuration value",
	Long:  "🌸 Set a specific configuration value by key",
	Args:  cobra.ExactArgs(2),
	Run:   configSetHandler,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "📋 List all configuration",
	Long:  "🌸 Display all current configuration values",
	Run:   configListHandler,
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "🔄 Reset configuration to defaults",
	Long:  "🌸 Reset all configuration values to defaults",
	Run:   configResetHandler,
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "📁 Show configuration file path",
	Long:  "🌸 Display the path to the configuration file",
	Run:   configPathHandler,
}

func configGetHandler(cmd *cobra.Command, args []string) {
	key := args[0]

	root, err := config.Load()
	if err != nil {
		fmt.Println("❌ Configuration not loaded")
		return
	}
	focus := root.Focus

	switch key {
	case "ai.provider":
		fmt.Println(focus.AI.Provider)
	case "ai.model":
		fmt.Println(focus.AI.Model)
	case "ai.temperature":
		fmt.Printf("%.1f\n", focus.AI.Temperature)
	case "ai.max_tokens":
		fmt.Printf("%d\n", focus.AI.MaxTokens)
	case "ai.timeout":
		fmt.Printf("%d\n", focus.AI.Timeout)
	case "theme":
		fmt.Println(focus.Theme)
	case "dashboard.auto_refresh":
		fmt.Printf("%t\n", focus.Dashboard.AutoRefresh)
	case "dashboard.refresh_interval":
		fmt.Printf("%d\n", focus.Dashboard.RefreshInterval)
	case "dashboard.show_animation":
		fmt.Printf("%t\n", focus.Dashboard.ShowAnimation)
	case "dashboard.compact_mode":
		fmt.Printf("%t\n", focus.Dashboard.CompactMode)
	case "notes.default_editor":
		fmt.Println(focus.Notes.DefaultEditor)
	case "notes.auto_save":
		fmt.Printf("%t\n", focus.Notes.AutoSave)
	case "notes.save_interval":
		fmt.Printf("%d\n", focus.Notes.SaveInterval)
	case "ui.time_format":
		fmt.Println(focus.UI.TimeFormat)
	case "ui.date_format":
		fmt.Println(focus.UI.DateFormat)
	case "ui.show_help_tips":
		fmt.Printf("%t\n", focus.UI.ShowHelpTips)
	case "ui.notifications":
		fmt.Printf("%t\n", focus.UI.Notifications)
	case "ui.sound_effects":
		fmt.Printf("%t\n", focus.UI.SoundEffects)
	default:
		fmt.Printf("❌ Unknown configuration key: %s\n", key)
		return
	}
}

func configSetHandler(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	// Validate and convert value based on key
	var finalValue any = value

	switch key {
	case "ai.temperature":
		if temp, err := strconv.ParseFloat(value, 64); err == nil {
			finalValue = temp
		} else {
			fmt.Printf("❌ Invalid temperature value: %s\n", value)
			return
		}
	case "ai.max_tokens", "ai.timeout", "dashboard.refresh_interval", "notes.save_interval":
		if intVal, err := strconv.Atoi(value); err == nil {
			finalValue = intVal
		} else {
			fmt.Printf("❌ Invalid integer value: %s\n", value)
			return
		}
	case "theme":
		validThemes := []string{"synthwave", "light", "plain"}
		if !contains(validThemes, value) {
			fmt.Printf("❌ Invalid theme. Valid options: %s\n", strings.Join(validThemes, ", "))
			return
		}
	case "ui.time_format":
		validFormats := []string{"12h", "24h"}
		if !contains(validFormats, value) {
			fmt.Printf("❌ Invalid time format. Valid options: %s\n", strings.Join(validFormats, ", "))
			return
		}
	case "dashboard.auto_refresh", "dashboard.show_animation", "dashboard.compact_mode",
		"notes.auto_save", "ui.show_help_tips", "ui.notifications", "ui.sound_effects":
		if boolVal, err := strconv.ParseBool(value); err == nil {
			finalValue = boolVal
		} else {
			fmt.Printf("❌ Invalid boolean value: %s (use true/false)\n", value)
			return
		}
	}

	// Update configuration
	root, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load configuration: %v\n", err)
		return
	}
	switch key {
	case "ai.provider":
		root.Focus.AI.Provider = value
	case "ai.model":
		root.Focus.AI.Model = value
	case "ai.temperature":
		if v, ok := finalValue.(float64); ok {
			root.Focus.AI.Temperature = v
		}
	case "ai.max_tokens":
		if v, ok := finalValue.(int); ok {
			root.Focus.AI.MaxTokens = v
		}
	case "ai.timeout":
		if v, ok := finalValue.(int); ok {
			root.Focus.AI.Timeout = v
		}
	case "theme":
		root.Focus.Theme = value
	case "dashboard.auto_refresh":
		if v, ok := finalValue.(bool); ok {
			root.Focus.Dashboard.AutoRefresh = v
		}
	case "dashboard.refresh_interval":
		if v, ok := finalValue.(int); ok {
			root.Focus.Dashboard.RefreshInterval = v
		}
	case "dashboard.show_animation":
		if v, ok := finalValue.(bool); ok {
			root.Focus.Dashboard.ShowAnimation = v
		}
	case "dashboard.compact_mode":
		if v, ok := finalValue.(bool); ok {
			root.Focus.Dashboard.CompactMode = v
		}
	case "notes.default_editor":
		root.Focus.Notes.DefaultEditor = value
	case "notes.auto_save":
		if v, ok := finalValue.(bool); ok {
			root.Focus.Notes.AutoSave = v
		}
	case "notes.save_interval":
		if v, ok := finalValue.(int); ok {
			root.Focus.Notes.SaveInterval = v
		}
	case "ui.time_format":
		root.Focus.UI.TimeFormat = value
	case "ui.date_format":
		root.Focus.UI.DateFormat = value
	case "ui.show_help_tips":
		if v, ok := finalValue.(bool); ok {
			root.Focus.UI.ShowHelpTips = v
		}
	case "ui.notifications":
		if v, ok := finalValue.(bool); ok {
			root.Focus.UI.Notifications = v
		}
	case "ui.sound_effects":
		if v, ok := finalValue.(bool); ok {
			root.Focus.UI.SoundEffects = v
		}
	}
	if err := config.Save(root); err != nil {
		fmt.Printf("❌ Failed to save configuration: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true)

	fmt.Printf("%s Configuration updated: %s = %v\n",
		successStyle.Render("✅"), key, finalValue)
}

func configListHandler(cmd *cobra.Command, args []string) {
	root, err := config.Load()
	if err != nil {
		fmt.Println("❌ Configuration not loaded")
		return
	}
	focus := root.Focus

	headerStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true)

	sectionStyle := lipgloss.NewStyle().
		Foreground(styles.GetAccent()).
		Bold(true)

	fmt.Println(headerStyle.Render("🔧 focus.sh Configuration"))
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println(sectionStyle.Render("🤖 AI Configuration"))
	fmt.Printf("  Provider:    %s\n", focus.AI.Provider)
	fmt.Printf("  Model:       %s\n", focus.AI.Model)
	fmt.Printf("  Temperature: %.1f\n", focus.AI.Temperature)
	fmt.Printf("  Max Tokens:  %d\n", focus.AI.MaxTokens)
	fmt.Printf("  Timeout:     %d seconds\n", focus.AI.Timeout)

	fmt.Println(sectionStyle.Render("🎨 Theme Configuration"))
	fmt.Printf("  Theme:       %s\n", focus.Theme)

	fmt.Println(sectionStyle.Render("📊 Dashboard Configuration"))
	fmt.Printf("  Auto Refresh:  %t\n", focus.Dashboard.AutoRefresh)
	fmt.Printf("  Refresh Interval: %d seconds\n", focus.Dashboard.RefreshInterval)
	fmt.Printf("  Show Animation: %t\n", focus.Dashboard.ShowAnimation)
	fmt.Printf("  Compact Mode:   %t\n", focus.Dashboard.CompactMode)

	fmt.Println(sectionStyle.Render("📝 Notes Configuration"))
	fmt.Printf("  Default Editor: %s\n", focus.Notes.DefaultEditor)
	fmt.Printf("  Auto Save:      %t\n", focus.Notes.AutoSave)
	fmt.Printf("  Save Interval:  %d minutes\n", focus.Notes.SaveInterval)
	if focus.Notes.Directory != "" {
		fmt.Printf("  Directory:      %s\n", focus.Notes.Directory)
	}

	fmt.Println(sectionStyle.Render("🎮 UI Configuration"))
	fmt.Printf("  Time Format:     %s\n", focus.UI.TimeFormat)
	fmt.Printf("  Date Format:     %s\n", focus.UI.DateFormat)
	fmt.Printf("  Show Help Tips:  %t\n", focus.UI.ShowHelpTips)
	fmt.Printf("  Notifications:   %t\n", focus.UI.Notifications)
	fmt.Printf("  Sound Effects:   %t\n", focus.UI.SoundEffects)

	fmt.Println(strings.Repeat("─", 50))
}

func configResetHandler(cmd *cobra.Command, args []string) {
	confirmStyle := lipgloss.NewStyle().
		Foreground(styles.GetError()).
		Bold(true)

	fmt.Print(confirmStyle.Render("⚠️  WARNING: This will reset all configuration to defaults!"))
	fmt.Print(" Are you sure? (y/N): ")

	var response string
	_, _ = fmt.Scanln(&response) // Ignore error for user input

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("❌ Configuration reset cancelled.")
		return
	}

	// Remove config file
	configPath := config.ConfigPath()
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("❌ Failed to remove config file: %v\n", err)
		return
	}

	// Re-create a default config so the next Load() picks up defaults
	if err := config.Init(); err != nil {
		fmt.Printf("❌ Failed to recreate default configuration: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().
		Foreground(styles.GetSuccess()).
		Bold(true)

	fmt.Println(successStyle.Render("✅ Configuration reset to defaults!"))
}

func configPathHandler(cmd *cobra.Command, args []string) {
	configPath := config.ConfigPath()

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("📁 Configuration file will be created at:\n  %s\n", configPath)
	} else {
		fmt.Printf("📁 Configuration file location:\n  %s\n", configPath)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
