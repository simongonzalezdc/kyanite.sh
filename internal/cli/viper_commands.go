package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kyanite/focus/pkg/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	
	cfg := config.GetConfig()
	if cfg == nil {
		fmt.Println("❌ Configuration not loaded")
		return
	}

	// Get value using viper
	viper.SetDefault(key, "")
	value := viper.GetString(key)

	if value == "" {
		// Try nested keys
		switch key {
		case "ai.provider":
			value = cfg.AI.Provider
		case "ai.model":
			value = cfg.AI.Model
		case "ai.temperature":
			value = fmt.Sprintf("%.1f", cfg.AI.Temperature)
		case "ai.max_tokens":
			value = fmt.Sprintf("%d", cfg.AI.MaxTokens)
		case "ai.timeout":
			value = fmt.Sprintf("%d", cfg.AI.Timeout)
		case "theme":
			value = cfg.Theme
		case "dashboard.auto_refresh":
			value = fmt.Sprintf("%t", cfg.Dashboard.AutoRefresh)
		case "dashboard.refresh_interval":
			value = fmt.Sprintf("%d", cfg.Dashboard.RefreshInterval)
		case "dashboard.show_animation":
			value = fmt.Sprintf("%t", cfg.Dashboard.ShowAnimation)
		case "dashboard.compact_mode":
			value = fmt.Sprintf("%t", cfg.Dashboard.CompactMode)
		case "notes.default_editor":
			value = cfg.Notes.DefaultEditor
		case "notes.auto_save":
			value = fmt.Sprintf("%t", cfg.Notes.AutoSave)
		case "notes.save_interval":
			value = fmt.Sprintf("%d", cfg.Notes.SaveInterval)
		case "ui.time_format":
			value = cfg.UI.TimeFormat
		case "ui.date_format":
			value = cfg.UI.DateFormat
		case "ui.show_help_tips":
			value = fmt.Sprintf("%t", cfg.UI.ShowHelpTips)
		case "ui.notifications":
			value = fmt.Sprintf("%t", cfg.UI.Notifications)
		case "ui.sound_effects":
			value = fmt.Sprintf("%t", cfg.UI.SoundEffects)
		default:
			fmt.Printf("❌ Unknown configuration key: %s\n", key)
			return
		}
	}

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true)

	fmt.Printf("%s %s\n", 
		successStyle.Render(fmt.Sprintf("%s:", key)), 
		value)
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
	updates := map[string]any{key: finalValue}
	if err := config.UpdateConfig(updates); err != nil {
		fmt.Printf("❌ Failed to update configuration: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true)

	fmt.Printf("%s Configuration updated: %s = %v\n", 
		successStyle.Render("✅"), key, finalValue)
}

func configListHandler(cmd *cobra.Command, args []string) {
	cfg := config.GetConfig()
	if cfg == nil {
		fmt.Println("❌ Configuration not loaded")
		return
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF71CE")).
		Bold(true)

	sectionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFF0")).
		Bold(true)

	fmt.Println(headerStyle.Render("🔧 NEON Configuration"))
	fmt.Println(strings.Repeat("─", 50))

	fmt.Println(sectionStyle.Render("🤖 AI Configuration"))
	fmt.Printf("  Provider:    %s\n", cfg.AI.Provider)
	fmt.Printf("  Model:       %s\n", cfg.AI.Model)
	fmt.Printf("  Temperature: %.1f\n", cfg.AI.Temperature)
	fmt.Printf("  Max Tokens:  %d\n", cfg.AI.MaxTokens)
	fmt.Printf("  Timeout:     %d seconds\n", cfg.AI.Timeout)

	fmt.Println(sectionStyle.Render("🎨 Theme Configuration"))
	fmt.Printf("  Theme:       %s\n", cfg.Theme)

	fmt.Println(sectionStyle.Render("📊 Dashboard Configuration"))
	fmt.Printf("  Auto Refresh:  %t\n", cfg.Dashboard.AutoRefresh)
	fmt.Printf("  Refresh Interval: %d seconds\n", cfg.Dashboard.RefreshInterval)
	fmt.Printf("  Show Animation: %t\n", cfg.Dashboard.ShowAnimation)
	fmt.Printf("  Compact Mode:   %t\n", cfg.Dashboard.CompactMode)

	fmt.Println(sectionStyle.Render("📝 Notes Configuration"))
	fmt.Printf("  Default Editor: %s\n", cfg.Notes.DefaultEditor)
	fmt.Printf("  Auto Save:      %t\n", cfg.Notes.AutoSave)
	fmt.Printf("  Save Interval:  %d minutes\n", cfg.Notes.SaveInterval)
	if cfg.Notes.Directory != "" {
		fmt.Printf("  Directory:      %s\n", cfg.Notes.Directory)
	}

	fmt.Println(sectionStyle.Render("🎮 UI Configuration"))
	fmt.Printf("  Time Format:     %s\n", cfg.UI.TimeFormat)
	fmt.Printf("  Date Format:     %s\n", cfg.UI.DateFormat)
	fmt.Printf("  Show Help Tips:  %t\n", cfg.UI.ShowHelpTips)
	fmt.Printf("  Notifications:   %t\n", cfg.UI.Notifications)
	fmt.Printf("  Sound Effects:   %t\n", cfg.UI.SoundEffects)

	fmt.Println(strings.Repeat("─", 50))
}

func configResetHandler(cmd *cobra.Command, args []string) {
	confirmStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0040")).
		Bold(true)

	fmt.Print(confirmStyle.Render("⚠️  WARNING: This will reset all configuration to defaults!"))
	fmt.Print(" Are you sure? (y/N): ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
		fmt.Println("❌ Configuration reset cancelled.")
		return
	}

	// Remove config file
	configPath := config.GetConfigPath()
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("❌ Failed to remove config file: %v\n", err)
		return
	}

	// Reload with defaults
	if err := config.ReloadConfig(); err != nil {
		fmt.Printf("❌ Failed to reload configuration: %v\n", err)
		return
	}

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF66")).
		Bold(true)

	fmt.Println(successStyle.Render("✅ Configuration reset to defaults!"))
}

func configPathHandler(cmd *cobra.Command, args []string) {
	configPath := config.GetConfigPath()
	
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
