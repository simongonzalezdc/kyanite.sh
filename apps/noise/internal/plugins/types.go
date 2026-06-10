// Package plugins provides a plugin system for extending noise.sh functionality
package plugins

import (
	"context"
	"time"

	"github.com/Kyanite/noise/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

// PluginMetadata contains information about a plugin
type PluginMetadata struct {
	// Basic plugin information
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	License     string `json:"license"`

	// Plugin capabilities
	Capabilities []Capability `json:"capabilities"`

	// Plugin configuration
	Config map[string]interface{} `json:"config,omitempty"`

	// Plugin lifecycle
	LoadTime time.Time `json:"load_time"`
	Enabled  bool      `json:"enabled"`
}

// Capability represents what a plugin can do
type Capability string

const (
	// UI Capabilities
	CapabilityMenuItem    Capability = "menu_item"    // Can add menu items
	CapabilityScreen      Capability = "screen"       // Can provide new screens
	CapabilityUIExtension Capability = "ui_extension" // Can extend existing UI

	// Editor Capabilities
	CapabilityEditorTool   Capability = "editor_tool"   // Can add editor tools
	CapabilityExportFormat Capability = "export_format" // Can add export formats
	CapabilitySyntax       Capability = "syntax"        // Can provide syntax highlighting

	// Theory Capabilities
	CapabilityTheoryTool Capability = "theory_tool" // Can add theory analysis tools
	CapabilityChordLib   Capability = "chord_lib"   // Can provide chord libraries

	// Audio Capabilities
	CapabilityAudioEffect Capability = "audio_effect" // Can provide audio effects
	CapabilityMIDIHandler Capability = "midi_handler" // Can handle MIDI events

	// Data Capabilities
	CapabilityDataProvider Capability = "data_provider" // Can provide external data
	CapabilityExporter     Capability = "exporter"      // Can export data

	// System Capabilities
	CapabilityHook   Capability = "hook"   // Can hook into system events
	CapabilityConfig Capability = "config" // Can add configuration options
)

// PluginContext provides context and services to plugins
type PluginContext struct {
	// Application context
	Context context.Context

	// Configuration
	Config *config.Config

	// Plugin metadata
	Metadata *PluginMetadata

	// Services (will be expanded as needed)
	// Database *db.DB
	// Logger logging.Logger
}

// Plugin represents a noise.sh plugin
type Plugin interface {
	// Metadata returns information about the plugin
	Metadata() *PluginMetadata

	// Initialize sets up the plugin with the provided context
	Initialize(ctx *PluginContext) error

	// Cleanup performs cleanup when the plugin is unloaded
	Cleanup() error

	// Enable activates the plugin
	Enable() error

	// Disable deactivates the plugin
	Disable() error

	// IsEnabled returns whether the plugin is currently enabled
	IsEnabled() bool
}

// MenuItem represents a menu item that can be added by plugins
type MenuItem struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Shortcut    string         `json:"shortcut"`
	Icon        string         `json:"icon"`
	Handler     func() tea.Cmd `json:"-"`
	Children    []*MenuItem    `json:"children,omitempty"`
	Enabled     bool           `json:"enabled"`
}

// Screen represents a screen that can be provided by plugins
type Screen interface {
	// ID returns the unique identifier for this screen
	ID() string

	// Name returns the display name for this screen
	Name() string

	// Init initializes the screen
	Init() tea.Cmd

	// Update handles messages for this screen
	Update(msg tea.Msg) (tea.Model, tea.Cmd)

	// View renders the screen
	View() string

	// SetDimensions sets the screen dimensions
	SetDimensions(width, height int)

	// Focus sets focus on the screen
	Focus()

	// Blur removes focus from the screen
	Blur()
}

// EditorTool represents a tool that can be added to the editor
type EditorTool struct {
	ID          string                               `json:"id"`
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Icon        string                               `json:"icon"`
	Handler     func(content string) (string, error) `json:"-"`
	Shortcut    string                               `json:"shortcut"`
	Enabled     bool                                 `json:"enabled"`
}

// ExportFormat represents an export format that can be added by plugins
type ExportFormat struct {
	ID          string                               `json:"id"`
	Name        string                               `json:"name"`
	Description string                               `json:"description"`
	Extension   string                               `json:"extension"`
	MimeType    string                               `json:"mime_type"`
	Handler     func(content string) ([]byte, error) `json:"-"`
	Enabled     bool                                 `json:"enabled"`
}

// HookPoint represents points in the application lifecycle where plugins can hook
type HookPoint string

const (
	// Application lifecycle hooks
	HookPreInit     HookPoint = "pre_init"     // Before application initialization
	HookPostInit    HookPoint = "post_init"    // After application initialization
	HookPreShutdown HookPoint = "pre_shutdown" // Before application shutdown

	// UI lifecycle hooks
	HookScreenChange HookPoint = "screen_change" // When switching screens
	HookMenuOpen     HookPoint = "menu_open"     // When opening menu
	HookMenuClose    HookPoint = "menu_close"    // When closing menu

	// Editor lifecycle hooks
	HookContentChange HookPoint = "content_change" // When editor content changes
	HookContentSave   HookPoint = "content_save"   // When content is saved
	HookContentLoad   HookPoint = "content_load"   // When content is loaded

	// Export hooks
	HookPreExport  HookPoint = "pre_export"  // Before export
	HookPostExport HookPoint = "post_export" // After export

	// Theory hooks
	HookAnalysisComplete HookPoint = "analysis_complete" // After theory analysis

	// Audio hooks
	HookPlaybackStart HookPoint = "playback_start" // When audio playback starts
	HookPlaybackStop  HookPoint = "playback_stop"  // When audio playback stops
)

// HookData contains data passed to hook handlers
type HookData map[string]interface{}

// HookHandler is a function that handles plugin hooks
type HookHandler func(point HookPoint, data HookData) error

// PluginManager manages the plugin system
type PluginManager interface {
	// LoadPlugins discovers and loads all available plugins
	LoadPlugins() error

	// GetPlugin returns a plugin by ID
	GetPlugin(id string) (Plugin, error)

	// GetPlugins returns all loaded plugins
	GetPlugins() map[string]Plugin

	// GetPluginsByCapability returns plugins that have a specific capability
	GetPluginsByCapability(capability Capability) []Plugin

	// EnablePlugin enables a plugin
	EnablePlugin(id string) error

	// DisablePlugin disables a plugin
	DisablePlugin(id string) error

	// UnloadPlugin unloads a plugin completely
	UnloadPlugin(id string) error

	// RegisterHook registers a hook handler for a plugin
	RegisterHook(pluginID string, handler HookHandler) error

	// CallHook calls all registered hook handlers for a hook point
	CallHook(point HookPoint, data HookData) error

	// GetMenuItems returns all menu items provided by plugins
	GetMenuItems() []*MenuItem

	// GetScreens returns all screens provided by plugins
	GetScreens() map[string]Screen

	// GetEditorTools returns all editor tools provided by plugins
	GetEditorTools() []*EditorTool

	// GetExportFormats returns all export formats provided by plugins
	GetExportFormats() []*ExportFormat
}
