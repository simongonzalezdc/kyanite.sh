package plugins

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/puente-labs/noise/internal/ui/dimension"
	"github.com/puente-labs/noise/internal/ui/styles"
)

// PluginSettingsModel handles the plugin management interface
type PluginSettingsModel struct {
	manager PluginManager
	width   int
	height  int

	// UI state
	selectedPlugin int
	scrollOffset   int
	showDetails    bool
	detailsPlugin  string
}

// NewPluginSettingsModel creates a new plugin settings model
func NewPluginSettingsModel(manager PluginManager) *PluginSettingsModel {
	return &PluginSettingsModel{
		manager:        manager,
		selectedPlugin: 0,
		scrollOffset:   0,
		showDetails:    false,
	}
}

// Init initializes the plugin settings model
func (m *PluginSettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the plugin settings
func (m *PluginSettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selectedPlugin > 0 {
				m.selectedPlugin--
			}
		case "down", "j":
			if m.selectedPlugin < len(m.getPluginIDs())-1 {
				m.selectedPlugin++
			}
		case "enter":
			if !m.showDetails {
				pluginIDs := m.getPluginIDs()
				if m.selectedPlugin < len(pluginIDs) {
					for _, id := range pluginIDs {
						if m.selectedPlugin == 0 {
							m.detailsPlugin = id
							break
						}
						m.selectedPlugin--
					}
					m.showDetails = true
					m.selectedPlugin = 0
				}
			} else {
				// Toggle plugin enabled/disabled
				pluginIDs := m.getPluginIDs()
				if m.selectedPlugin < len(pluginIDs) {
					selectedID := pluginIDs[m.selectedPlugin]
					plugin, _ := m.manager.GetPlugin(selectedID)
					if plugin.IsEnabled() {
						if err := m.manager.DisablePlugin(selectedID); err != nil {
							// Handle error (could add error state to model)
							return m, nil
						}
					} else {
						if err := m.manager.EnablePlugin(selectedID); err != nil {
							// Handle error (could add error state to model)
							return m, nil
						}
					}
				}
			}
		case "d":
			if m.showDetails {
				// Show/hide plugin details
				m.showDetails = false
				m.detailsPlugin = ""
			}
		case "esc":
			if m.showDetails {
				m.showDetails = false
				m.detailsPlugin = ""
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the plugin settings interface
func (m *PluginSettingsModel) View() string {
	plugins := m.manager.GetPlugins()

	if len(plugins) == 0 {
		return m.renderEmptyState()
	}

	if m.showDetails && m.detailsPlugin != "" {
		return m.renderPluginDetails()
	}

	return m.renderPluginList()
}

// renderEmptyState renders the view when no plugins are loaded
func (m *PluginSettingsModel) renderEmptyState() string {
	content := strings.Builder{}

	title := styles.Title.Width(m.width).Render("Plugin Manager")
	content.WriteString(title + "\n\n")

	emptyState := styles.Muted.
		Align(lipgloss.Center).
		Width(m.width - 4).
		Padding(2).
		Render("No plugins loaded.\n\nPlugins can be placed in:\n• ~/.noise/plugins/\n• ./plugins/\n• ./internal/plugins/examples/")

	content.WriteString(emptyState)

	return content.String()
}

// renderPluginList renders the main plugin list
func (m *PluginSettingsModel) renderPluginList() string {
	content := strings.Builder{}

	title := styles.Title.Width(m.width).Render("Plugin Manager")
	content.WriteString(title + "\n\n")

	pluginIDs := m.getPluginIDs()
	start := m.scrollOffset
	end := m.scrollOffset + (m.height - 8) // Leave room for title and help

	if start > len(pluginIDs) {
		start = 0
		m.scrollOffset = 0
	}

	for i := start; i < len(pluginIDs) && i < end; i++ {
		pluginID := pluginIDs[i]
		plugin, _ := m.manager.GetPlugin(pluginID)

		var style lipgloss.Style
		if i == m.selectedPlugin {
			style = styles.ListItemSelected.Width(m.width - 4)
		} else {
			style = styles.ListItem.Width(m.width - 4)
		}

		status := "Disabled"
		statusColor := styles.Error
		if plugin.IsEnabled() {
			status = "Enabled"
			statusColor = styles.Success
		}

		statusStyled := lipgloss.NewStyle().Foreground(statusColor).Render(status)

		line := fmt.Sprintf("%s [%s] v%s",
			plugin.Metadata().Name,
			statusStyled,
			plugin.Metadata().Version)

		content.WriteString(style.Render(line) + "\n")
	}

	// Help text
	help := styles.Muted.
		Width(m.width).
		Padding(1).
		Render("↑/↓: Navigate • Enter: Toggle/Details • D: Back • Esc: Exit")

	content.WriteString("\n" + help)

	return content.String()
}

// renderPluginDetails renders detailed information about a specific plugin
func (m *PluginSettingsModel) renderPluginDetails() string {
	content := strings.Builder{}

	plugin, err := m.manager.GetPlugin(m.detailsPlugin)
	if err != nil {
		return m.renderPluginList()
	}

	metadata := plugin.Metadata()

	title := styles.Title.Width(m.width).Render(fmt.Sprintf("Plugin: %s", metadata.Name))
	content.WriteString(title + "\n\n")

	// Plugin information
	info := styles.Text.Width(m.width - 4).Padding(1).Render(fmt.Sprintf(`
ID:          %s
Version:     %s
Author:      %s
License:     %s
Description: %s
Status:      %s
Loaded:      %s
`,
		metadata.ID,
		metadata.Version,
		metadata.Author,
		metadata.License,
		metadata.Description,
		map[bool]string{true: "Enabled", false: "Disabled"}[plugin.IsEnabled()],
		metadata.LoadTime.Format("2006-01-02 15:04:05"),
	))

	content.WriteString(info + "\n")

	// Capabilities
	if len(metadata.Capabilities) > 0 {
		capabilities := strings.Join(m.formatCapabilities(metadata.Capabilities), ", ")
		caps := styles.Text.
			Width(m.width - 4).
			Padding(1).
			Render(fmt.Sprintf("Capabilities: %s", capabilities))
		content.WriteString(caps + "\n")
	}

	// Action buttons
	actions := m.renderPluginActions(plugin)
	content.WriteString(actions + "\n")

	// Help text
	help := styles.Muted.
		Width(m.width).
		Padding(1).
		Render("Enter: Toggle Enable/Disable • D: Back to List • Esc: Exit")

	content.WriteString(help)

	return content.String()
}

// renderPluginActions renders action buttons for a plugin
func (m *PluginSettingsModel) renderPluginActions(plugin Plugin) string {
	enabled := plugin.IsEnabled()

	var action string
	if enabled {
		action = "Disable Plugin"
	} else {
		action = "Enable Plugin"
	}

	return lipgloss.NewStyle().
		Width(m.width - 4).
		Padding(1).
		Render(fmt.Sprintf("Actions: [%s]", action))
}

// formatCapabilities formats capability strings for display
func (m *PluginSettingsModel) formatCapabilities(capabilities []Capability) []string {
	var formatted []string
	for _, cap := range capabilities {
		formatted = append(formatted, string(cap))
	}
	return formatted
}

// getPluginIDs returns a slice of plugin IDs in sorted order
func (m *PluginSettingsModel) getPluginIDs() []string {
	plugins := m.manager.GetPlugins()
	ids := make([]string, 0, len(plugins))

	for id := range plugins {
		ids = append(ids, id)
	}

	return ids
}

// SetDimensions sets the dimensions for the plugin settings
func (m *PluginSettingsModel) SetDimensions(width, height int) {
	dimension.Set(&m.width, &m.height, width, height)
}

func (m *PluginSettingsModel) GetDimensions() (int, int) {
	return m.width, m.height
}

// Focus sets focus on the plugin settings
func (m *PluginSettingsModel) Focus() {
	// Implementation for focus handling
}

// Blur removes focus from the plugin settings
func (m *PluginSettingsModel) Blur() {
	// Implementation for blur handling
}
