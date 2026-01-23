package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SettingsCategory represents different categories of settings
type SettingsCategory string

const (
	CategoryAppearance SettingsCategory = "Appearance"
	CategoryEditor     SettingsCategory = "Editor"
	CategoryAutoSave   SettingsCategory = "Auto-Save"
	CategoryFileOps    SettingsCategory = "File Operations"
	CategoryAudio      SettingsCategory = "Audio"
	CategoryAI         SettingsCategory = "AI"
	CategoryAdvanced   SettingsCategory = "Advanced"
)

// SettingsItem represents a single configurable setting
type SettingsItem struct {
	ID          string
	Category    SettingsCategory
	Name        string
	Description string
	Type        SettingsType
	Value       interface{}
	Default     interface{}
	Options     []string          // For select/radio types
	Min         float64           // For number types
	Max         float64           // For number types
	Unit        string            // For number types (e.g., "seconds", "MB")
	Callback    func(interface{}) // Called when value changes
}

// SettingsType represents the type of a setting
type SettingsType int

const (
	TypeBool SettingsType = iota
	TypeInt
	TypeFloat
	TypeString
	TypeSelect
	TypeDuration
)

// SettingsModel handles the comprehensive settings screen
type SettingsModel struct {
	width        int
	height       int
	activeTab    SettingsCategory
	categoryList list.Model
	settingsList list.Model
	config       *config.Config
	settings     map[string]*SettingsItem
	focused      FocusArea
	showSaveMsg  bool
	saveMsgTimer *time.Timer

	// Per-category scroll position tracking
	categoryScrollPos map[SettingsCategory]int
}

// FocusArea represents which part of the settings UI is focused
type FocusArea int

const (
	FocusCategories FocusArea = iota
	FocusSettings
	FocusValue
)

// KeyMap defines key bindings for the settings screen
type KeyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Enter key.Binding
	Tab   key.Binding
	Back  key.Binding
	Save  key.Binding
	Reset key.Binding
	Quit  key.Binding
}

// DefaultKeyMap returns the default key mappings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:    key.NewBinding(key.WithKeys("up", "k")),
		Down:  key.NewBinding(key.WithKeys("down", "j")),
		Left:  key.NewBinding(key.WithKeys("left", "h")),
		Right: key.NewBinding(key.WithKeys("right", "l")),
		Enter: key.NewBinding(key.WithKeys("enter", " ")),
		Tab:   key.NewBinding(key.WithKeys("tab")),
		Back:  key.NewBinding(key.WithKeys("esc")),
		Save:  key.NewBinding(key.WithKeys("ctrl+s")),
		Reset: key.NewBinding(key.WithKeys("ctrl+r")),
		Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c")),
	}
}

// NewSettingsModel creates a new comprehensive settings model
func NewSettingsModel(cfg *config.Config) *SettingsModel {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	m := &SettingsModel{
		config:            cfg,
		settings:          make(map[string]*SettingsItem),
		focused:           FocusCategories,
		categoryScrollPos: make(map[SettingsCategory]int),
	}

	m.initCategories()
	m.initSettings()
	m.setupLists()

	return m
}

// initCategories initializes the settings categories
func (m *SettingsModel) initCategories() {
	categories := []SettingsCategory{
		CategoryAppearance,
		CategoryEditor,
		CategoryAutoSave,
		CategoryFileOps,
		CategoryAudio,
		CategoryAI,
		CategoryAdvanced,
	}

	items := make([]list.Item, len(categories))
	for i, cat := range categories {
		items[i] = categoryItem{category: cat}
	}

	m.categoryList = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.categoryList.Title = "Settings Categories"
	m.categoryList.SetShowStatusBar(false)
	m.categoryList.SetShowHelp(false)
	m.activeTab = CategoryAppearance
}

// initSettings initializes all settings items
func (m *SettingsModel) initSettings() {
	// Appearance settings
	m.addSetting(&SettingsItem{
		ID:          "ui.theme",
		Category:    CategoryAppearance,
		Name:        "Theme",
		Description: "Color theme for the application",
		Type:        TypeSelect,
		Value:       m.config.UI.Theme,
		Default:     "amber-night",
		Options:     []string{"monochrome", "amber-night", "twilight-mist", "indigo-depths", "forest-path", "clay-earth", "iron-forge", "sunlight", "cyan-wave", "electric-rose"},
		Callback:    m.onThemeChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "ui.font_size",
		Category:    CategoryAppearance,
		Name:        "Font Size",
		Description: "Font size for the interface",
		Type:        TypeInt,
		Value:       m.config.UI.FontSize,
		Default:     12,
		Min:         8,
		Max:         24,
		Unit:        "pt",
		Callback:    m.onFontSizeChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "ui.animations",
		Category:    CategoryAppearance,
		Name:        "Animations",
		Description: "Enable UI animations",
		Type:        TypeBool,
		Value:       m.config.UI.Animations,
		Default:     true,
		Callback:    m.onAnimationsChange,
	})

	// Editor settings
	m.addSetting(&SettingsItem{
		ID:          "ui.show_line_numbers",
		Category:    CategoryEditor,
		Name:        "Line Numbers",
		Description: "Show line numbers in editor",
		Type:        TypeBool,
		Value:       m.config.UI.ShowLineNumbers,
		Default:     true,
		Callback:    m.onLineNumbersChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "ui.word_wrap",
		Category:    CategoryEditor,
		Name:        "Word Wrap",
		Description: "Enable word wrapping in editor",
		Type:        TypeBool,
		Value:       m.config.UI.WordWrap,
		Default:     true,
		Callback:    m.onWordWrapChange,
	})

	// Additional editor settings
	m.addSetting(&SettingsItem{
		ID:          "editor.tab_size",
		Category:    CategoryEditor,
		Name:        "Tab Size",
		Description: "Number of spaces per tab",
		Type:        TypeInt,
		Value:       4,
		Default:     4,
		Min:         2,
		Max:         8,
		Callback:    m.onTabSizeChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "editor.insert_spaces",
		Category:    CategoryEditor,
		Name:        "Insert Spaces",
		Description: "Insert spaces instead of tabs",
		Type:        TypeBool,
		Value:       true,
		Default:     true,
		Callback:    m.onInsertSpacesChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "editor.auto_indent",
		Category:    CategoryEditor,
		Name:        "Auto Indent",
		Description: "Automatically indent new lines",
		Type:        TypeBool,
		Value:       true,
		Default:     true,
		Callback:    m.onAutoIndentChange,
	})

	// Auto-save settings
	m.addSetting(&SettingsItem{
		ID:          "app.auto_save",
		Category:    CategoryAutoSave,
		Name:        "Auto-Save",
		Description: "Automatically save changes",
		Type:        TypeBool,
		Value:       m.config.App.AutoSave,
		Default:     true,
		Callback:    m.onAutoSaveChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "app.auto_save_interval",
		Category:    CategoryAutoSave,
		Name:        "Auto-Save Interval",
		Description: "Time between auto-saves",
		Type:        TypeDuration,
		Value:       m.config.App.AutoSaveInterval,
		Default:     30 * time.Second,
		Min:         5,
		Max:         300,
		Unit:        "seconds",
		Callback:    m.onAutoSaveIntervalChange,
	})

	// Add more auto-save settings that integrate with AutoSaveService
	m.addSetting(&SettingsItem{
		ID:          "autosave.debounce_ms",
		Category:    CategoryAutoSave,
		Name:        "Input Debounce",
		Description: "Delay before saving after user stops typing",
		Type:        TypeInt,
		Value:       2000, // Default from AutoSaveService
		Default:     2000,
		Min:         100,
		Max:         10000,
		Unit:        "ms",
		Callback:    m.onAutoSaveDebounceChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "autosave.max_retries",
		Category:    CategoryAutoSave,
		Name:        "Max Retry Attempts",
		Description: "Maximum number of retry attempts for failed saves",
		Type:        TypeInt,
		Value:       3, // Default from AutoSaveService
		Default:     3,
		Min:         1,
		Max:         10,
		Callback:    m.onAutoSaveMaxRetriesChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "autosave.enable_versioning",
		Category:    CategoryAutoSave,
		Name:        "Version History",
		Description: "Keep version history for files",
		Type:        TypeBool,
		Value:       true, // Default from AutoSaveService
		Default:     true,
		Callback:    m.onAutoSaveVersioningChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "autosave.max_versions",
		Category:    CategoryAutoSave,
		Name:        "Max Versions",
		Description: "Maximum number of versions to keep per file",
		Type:        TypeInt,
		Value:       10, // Default from AutoSaveService
		Default:     10,
		Min:         1,
		Max:         100,
		Callback:    m.onAutoSaveMaxVersionsChange,
	})

	// File operations settings
	m.addSetting(&SettingsItem{
		ID:          "app.max_recent_files",
		Category:    CategoryFileOps,
		Name:        "Recent Files",
		Description: "Maximum number of recent files to remember",
		Type:        TypeInt,
		Value:       m.config.App.MaxRecentFiles,
		Default:     10,
		Min:         5,
		Max:         50,
		Callback:    m.onMaxRecentFilesChange,
	})

	// Additional file operation settings
	m.addSetting(&SettingsItem{
		ID:          "files.default_save_location",
		Category:    CategoryFileOps,
		Name:        "Default Save Location",
		Description: "Default directory for saving files",
		Type:        TypeString,
		Value:       m.config.GetDataDir(),
		Default:     m.config.GetDataDir(),
		Callback:    m.onDefaultSaveLocationChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "files.auto_backup",
		Category:    CategoryFileOps,
		Name:        "Auto Backup",
		Description: "Create backup files before saving",
		Type:        TypeBool,
		Value:       true,
		Default:     true,
		Callback:    m.onAutoBackupChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "files.backup_count",
		Category:    CategoryFileOps,
		Name:        "Backup Count",
		Description: "Number of backup files to keep",
		Type:        TypeInt,
		Value:       5,
		Default:     5,
		Min:         1,
		Max:         20,
		Callback:    m.onBackupCountChange,
	})

	// Audio settings
	m.addSetting(&SettingsItem{
		ID:          "audio.enabled",
		Category:    CategoryAudio,
		Name:        "Audio Support",
		Description: "Enable audio playback features",
		Type:        TypeBool,
		Value:       m.config.Audio.Enabled,
		Default:     true,
		Callback:    m.onAudioEnabledChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "audio.playback_gain",
		Category:    CategoryAudio,
		Name:        "Playback Volume",
		Description: "Audio playback volume",
		Type:        TypeFloat,
		Value:       m.config.Audio.PlaybackGain,
		Default:     0.8,
		Min:         0.1,
		Max:         2.0,
		Callback:    m.onPlaybackGainChange,
	})

	// AI settings
	m.addSetting(&SettingsItem{
		ID:          "ai.enabled",
		Category:    CategoryAI,
		Name:        "AI Assistant",
		Description: "Enable AI-powered features",
		Type:        TypeBool,
		Value:       m.config.AI.Enabled,
		Default:     true,
		Callback:    m.onAIEnabledChange,
	})

	m.addSetting(&SettingsItem{
		ID:          "ai.temperature",
		Category:    CategoryAI,
		Name:        "AI Temperature",
		Description: "Creativity/randomness of AI responses",
		Type:        TypeFloat,
		Value:       m.config.AI.Temperature,
		Default:     0.7,
		Min:         0.0,
		Max:         1.0,
		Callback:    m.onAITemperatureChange,
	})

	// Advanced settings
	m.addSetting(&SettingsItem{
		ID:          "dev.debug",
		Category:    CategoryAdvanced,
		Name:        "Debug Mode",
		Description: "Enable debug logging",
		Type:        TypeBool,
		Value:       m.config.Dev.Debug,
		Default:     false,
		Callback:    m.onDebugModeChange,
	})
}

// addSetting adds a setting to the settings map
func (m *SettingsModel) addSetting(item *SettingsItem) {
	m.settings[item.ID] = item
}

// setupLists configures the list components
func (m *SettingsModel) setupLists() {
	// Configure category list
	m.categoryList.SetSize(25, 10)

	// Configure settings list
	m.settingsList = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.settingsList.Title = "Settings"
	m.settingsList.SetShowStatusBar(false)
	m.settingsList.SetShowHelp(false)
	m.settingsList.SetSize(50, 15)

	m.updateSettingsList()
}

// updateSettingsList updates the settings list for the current category
func (m *SettingsModel) updateSettingsList() {
	var items []list.Item
	for _, setting := range m.settings {
		if setting.Category == m.activeTab {
			items = append(items, settingItem{setting: setting})
		}
	}

	m.settingsList.SetItems(items)
}

// Init initializes the settings model
func (m *SettingsModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the settings
func (m *SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.categoryList.SetSize(25, m.height-10)
		m.settingsList.SetSize(50, m.height-10)

	case tea.MouseMsg:
		// Handle mouse events for settings
		cmd := m.handleMouse(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap().Tab):
			m.cycleFocus()
			return m, nil

		case key.Matches(msg, DefaultKeyMap().Back):
			if m.focused == FocusSettings && m.settingsList.Index() >= 0 {
				// Go back to categories
				m.focused = FocusCategories
			} else {
				// Exit settings
				return m, emitSettingsChangedMsg(false)
			}

		case key.Matches(msg, DefaultKeyMap().Save):
			m.saveSettings()
			return m, nil

		case key.Matches(msg, DefaultKeyMap().Quit):
			return m, emitSettingsChangedMsg(false)
		}

		// Handle focus-specific key presses
		switch m.focused {
		case FocusCategories:
			var cmd tea.Cmd
			m.categoryList, cmd = m.categoryList.Update(msg)
			cmds = append(cmds, cmd)

			// Update active tab when selection changes
			if m.categoryList.Index() >= 0 {
				selected := m.categoryList.SelectedItem()
				if catItem, ok := selected.(categoryItem); ok {
					if catItem.category != m.activeTab {
						// Save scroll position of current category before switching
						m.categoryScrollPos[m.activeTab] = m.settingsList.Index()

						m.activeTab = catItem.category
						m.updateSettingsList()

						// Restore scroll position for new category (or start at 0)
						if savedPos, ok := m.categoryScrollPos[m.activeTab]; ok {
							m.settingsList.Select(savedPos)
						} else {
							m.settingsList.Select(0)
						}
					}
				}
			}

		case FocusSettings:
			var cmd tea.Cmd
			m.settingsList, cmd = m.settingsList.Update(msg)
			cmds = append(cmds, cmd)

			// Handle setting value changes
			if key.Matches(msg, DefaultKeyMap().Enter) && m.settingsList.Index() >= 0 {
				selected := m.settingsList.SelectedItem()
				if settingItem, ok := selected.(settingItem); ok {
					m.editSetting(settingItem.setting)
					// Update the settings list to reflect changes
					m.updateSettingsList()
				}
			} else if key.Matches(msg, DefaultKeyMap().Left) && m.settingsList.Index() >= 0 {
				// Handle left arrow for certain setting types (like decreasing values)
				selected := m.settingsList.SelectedItem()
				if settingItem, ok := selected.(settingItem); ok {
					m.decrementSetting(settingItem.setting)
					m.updateSettingsList()
				}
			} else if key.Matches(msg, DefaultKeyMap().Right) && m.settingsList.Index() >= 0 {
				// Handle right arrow for certain setting types (like increasing values)
				selected := m.settingsList.SelectedItem()
				if settingItem, ok := selected.(settingItem); ok {
					m.incrementSetting(settingItem.setting)
					m.updateSettingsList()
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

// cycleFocus cycles between focus areas
func (m *SettingsModel) cycleFocus() {
	switch m.focused {
	case FocusCategories:
		m.focused = FocusSettings
	case FocusSettings:
		m.focused = FocusCategories
	}
}

// handleMouse processes mouse events for settings
func (m *SettingsModel) handleMouse(msg tea.MouseMsg) tea.Cmd {
	// Determine which panel was clicked based on X position
	categoryPanelWidth := 27 // matches categoryStyle width

	switch msg.Button {
	case tea.MouseButtonLeft:
		if msg.Action == tea.MouseActionRelease {
			if msg.X < categoryPanelWidth {
				// Click in category panel
				m.focused = FocusCategories

				// Calculate which category was clicked
				headerOffset := 6 // header area
				itemHeight := 1   // single line per category

				clickedIdx := (msg.Y - headerOffset) / itemHeight
				if clickedIdx >= 0 && clickedIdx < len(m.categoryList.Items()) {
					// Save scroll position of current category
					m.categoryScrollPos[m.activeTab] = m.settingsList.Index()

					m.categoryList.Select(clickedIdx)

					// Update active tab
					selected := m.categoryList.SelectedItem()
					if catItem, ok := selected.(categoryItem); ok {
						m.activeTab = catItem.category
						m.updateSettingsList()

						// Restore scroll position
						if savedPos, ok := m.categoryScrollPos[m.activeTab]; ok {
							m.settingsList.Select(savedPos)
						} else {
							m.settingsList.Select(0)
						}
					}
				}
			} else {
				// Click in settings panel
				m.focused = FocusSettings

				// Calculate which setting was clicked
				headerOffset := 6
				itemHeight := 5 // settings items are larger (name + desc + value)

				clickedIdx := (msg.Y - headerOffset) / itemHeight
				if clickedIdx >= 0 && clickedIdx < len(m.settingsList.Items()) {
					m.settingsList.Select(clickedIdx)

					// Activate the setting (toggle/edit)
					selected := m.settingsList.SelectedItem()
					if settingItem, ok := selected.(settingItem); ok {
						m.editSetting(settingItem.setting)
						m.updateSettingsList()
					}
				}
			}
		}

	case tea.MouseButtonWheelUp:
		if m.focused == FocusCategories {
			if m.categoryList.Index() > 0 {
				m.categoryList.CursorUp()
			}
		} else {
			if m.settingsList.Index() > 0 {
				m.settingsList.CursorUp()
			}
		}

	case tea.MouseButtonWheelDown:
		if m.focused == FocusCategories {
			if m.categoryList.Index() < len(m.categoryList.Items())-1 {
				m.categoryList.CursorDown()
			}
		} else {
			if m.settingsList.Index() < len(m.settingsList.Items())-1 {
				m.settingsList.CursorDown()
			}
		}
	}

	return nil
}

// editSetting handles editing a specific setting
func (m *SettingsModel) editSetting(setting *SettingsItem) {
	// For now, we'll implement basic editing in a future update
	// This would open a modal or inline editor for the setting value
	// For demonstration, we'll show a simple toggle for boolean settings
	if setting.Type == TypeBool {
		if boolVal, ok := setting.Value.(bool); ok {
			setting.Value = !boolVal
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	} else if setting.Type == TypeSelect {
		// Cycle through options for select type
		currentIndex := -1
		if strVal, ok := setting.Value.(string); ok {
			for i, option := range setting.Options {
				if option == strVal {
					currentIndex = i
					break
				}
			}
		}

		if currentIndex >= 0 {
			nextIndex := (currentIndex + 1) % len(setting.Options)
			setting.Value = setting.Options[nextIndex]
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	}
}

// incrementSetting increments a numeric setting value
func (m *SettingsModel) incrementSetting(setting *SettingsItem) {
	switch setting.Type {
	case TypeInt:
		if intVal, ok := setting.Value.(int); ok && intVal < int(setting.Max) {
			setting.Value = intVal + 1
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	case TypeFloat:
		if floatVal, ok := setting.Value.(float64); ok && floatVal < setting.Max {
			setting.Value = floatVal + 0.1
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	case TypeDuration:
		if duration, ok := setting.Value.(time.Duration); ok && duration.Seconds() < setting.Max {
			setting.Value = duration + time.Second
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	}
}

// decrementSetting decrements a numeric setting value
func (m *SettingsModel) decrementSetting(setting *SettingsItem) {
	switch setting.Type {
	case TypeInt:
		if intVal, ok := setting.Value.(int); ok && intVal > int(setting.Min) {
			setting.Value = intVal - 1
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	case TypeFloat:
		if floatVal, ok := setting.Value.(float64); ok && floatVal > setting.Min {
			setting.Value = floatVal - 0.1
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	case TypeDuration:
		if duration, ok := setting.Value.(time.Duration); ok && duration.Seconds() > setting.Min {
			setting.Value = duration - time.Second
			if setting.Callback != nil {
				setting.Callback(setting.Value)
			}
		}
	}
}

// saveSettings saves the current settings
func (m *SettingsModel) saveSettings() {
	// Update config with current setting values
	for _, setting := range m.settings {
		switch setting.ID {
		case "ui.theme":
			if strVal, ok := setting.Value.(string); ok {
				m.config.UI.Theme = strVal
			}
		case "ui.font_size":
			if intVal, ok := setting.Value.(int); ok {
				m.config.UI.FontSize = intVal
			}
		case "ui.animations":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.UI.Animations = boolVal
			}
		case "ui.show_line_numbers":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.UI.ShowLineNumbers = boolVal
			}
		case "ui.word_wrap":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.UI.WordWrap = boolVal
			}
		case "app.auto_save":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.App.AutoSave = boolVal
			}
		case "app.auto_save_interval":
			if durVal, ok := setting.Value.(time.Duration); ok {
				m.config.App.AutoSaveInterval = durVal
			}
		case "app.max_recent_files":
			if intVal, ok := setting.Value.(int); ok {
				m.config.App.MaxRecentFiles = intVal
			}
		case "audio.enabled":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.Audio.Enabled = boolVal
			}
		case "audio.playback_gain":
			if floatVal, ok := setting.Value.(float64); ok {
				m.config.Audio.PlaybackGain = floatVal
			}
		case "ai.enabled":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.AI.Enabled = boolVal
			}
		case "ai.temperature":
			if floatVal, ok := setting.Value.(float64); ok {
				m.config.AI.Temperature = floatVal
			}
		case "dev.debug":
			if boolVal, ok := setting.Value.(bool); ok {
				m.config.Dev.Debug = boolVal
			}
		// Editor settings
		case "editor.tab_size":
			// Store in a way that can be accessed by editor components
		case "editor.insert_spaces":
			// Store in a way that can be accessed by editor components
		case "editor.auto_indent":
			// Store in a way that can be accessed by editor components
		// File operation settings
		case "files.default_save_location":
			// Update default save location
		case "files.auto_backup":
			// Update auto backup setting
		case "files.backup_count":
			// Update backup count
		// Auto-save settings
		case "autosave.debounce_ms":
			// Update auto-save debounce setting
		case "autosave.max_retries":
			// Update auto-save max retries
		case "autosave.enable_versioning":
			// Update auto-save versioning
		case "autosave.max_versions":
			// Update auto-save max versions
		}
	}

	// Save config to file
	if err := m.config.Save(); err != nil {
		// Handle save error - could show error message to user
		return
	}

	m.showSaveMsg = true
	if m.saveMsgTimer != nil {
		m.saveMsgTimer.Stop()
	}
	m.saveMsgTimer = time.AfterFunc(3*time.Second, func() {
		m.showSaveMsg = false
	})

	// Emit settings changed message to notify other components
	emitSettingsChangedMsg(true)
}

// Callback functions for setting changes
func (m *SettingsModel) onThemeChange(value interface{}) {
	// Theme change would require reloading the UI
}

func (m *SettingsModel) onFontSizeChange(value interface{}) {
	// Font size change would require updating the UI
}

func (m *SettingsModel) onAnimationsChange(value interface{}) {
	// Animation setting change
}

func (m *SettingsModel) onLineNumbersChange(value interface{}) {
	// Line numbers setting change
}

func (m *SettingsModel) onWordWrapChange(value interface{}) {
	// Word wrap setting change
}

func (m *SettingsModel) onTabSizeChange(value interface{}) {
	// Tab size setting change
}

func (m *SettingsModel) onInsertSpacesChange(value interface{}) {
	// Insert spaces setting change
}

func (m *SettingsModel) onAutoIndentChange(value interface{}) {
	// Auto indent setting change
}

func (m *SettingsModel) onAutoSaveChange(value interface{}) {
	// Auto-save enable/disable
}

func (m *SettingsModel) onAutoSaveIntervalChange(value interface{}) {
	// Auto-save interval change
}

func (m *SettingsModel) onAutoSaveDebounceChange(value interface{}) {
	// Auto-save debounce change
}

func (m *SettingsModel) onAutoSaveMaxRetriesChange(value interface{}) {
	// Auto-save max retries change
}

func (m *SettingsModel) onAutoSaveVersioningChange(value interface{}) {
	// Auto-save versioning change
}

func (m *SettingsModel) onAutoSaveMaxVersionsChange(value interface{}) {
	// Auto-save max versions change
}

func (m *SettingsModel) onMaxRecentFilesChange(value interface{}) {
	// Max recent files change
}

func (m *SettingsModel) onDefaultSaveLocationChange(value interface{}) {
	// Default save location change
}

func (m *SettingsModel) onAutoBackupChange(value interface{}) {
	// Auto backup change
}

func (m *SettingsModel) onBackupCountChange(value interface{}) {
	// Backup count change
}

func (m *SettingsModel) onAudioEnabledChange(value interface{}) {
	// Audio enabled change
}

func (m *SettingsModel) onPlaybackGainChange(value interface{}) {
	// Playback gain change
}

func (m *SettingsModel) onAIEnabledChange(value interface{}) {
	// AI enabled change
}

func (m *SettingsModel) onAITemperatureChange(value interface{}) {
	// AI temperature change
}

func (m *SettingsModel) onDebugModeChange(value interface{}) {
	// Debug mode change
}

// View renders the settings screen
func (m *SettingsModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Loading settings..."
	}

	// Create main layout
	mainStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Padding(1)

	// Get current theme
	t := theme.GetManager().Current()

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		MarginBottom(1)
	header := headerStyle.Render("Settings")
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Padding(0, 1)
	header = borderStyle.Render(header)
	header += "\n\n"

	// Categories section
	var categoryStyle lipgloss.Style
	if m.focused == FocusCategories {
		categoryStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Padding(1).
			MarginRight(2).
			Width(27).
			Height(m.height - 8)
	} else {
		categoryStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary).
			Padding(1).
			MarginRight(2).
			Width(27).
			Height(m.height - 8)
	}

	categoryView := categoryStyle.Render(m.categoryList.View())

	// Settings section
	var settingsStyle lipgloss.Style
	if m.focused == FocusSettings {
		settingsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Accent).
			Padding(1).
			Width(52).
			Height(m.height - 8)
	} else {
		settingsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Secondary).
			Padding(1).
			Width(52).
			Height(m.height - 8)
	}

	settingsView := settingsStyle.Render(m.settingsList.View())

	// Instructions
	instructions := m.renderInstructions()

	// Save message
	var saveMsg string
	if m.showSaveMsg {
		saveMsgStyle := lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true)
		saveMsg = saveMsgStyle.Render("✔ Settings saved successfully!")
	}

	// Combine sections
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		categoryView,
		settingsView,
	)

	content = lipgloss.JoinVertical(
		lipgloss.Top,
		header,
		content,
		instructions,
		saveMsg,
	)

	return mainStyle.Render(content)
}

// renderInstructions renders the help text
func (m *SettingsModel) renderInstructions() string {
	instructions := []string{
		"Tab: Switch focus",
		"Up/Down: Navigate",
		"Enter: Edit setting",
		"Up/Up: Adjust values",
		"Ctrl+S: Save",
		"Esc: Back/Quit",
	}

	t := theme.GetManager().Current()
	mutedStyle := lipgloss.NewStyle().
		Foreground(t.Secondary)
	return mutedStyle.Render(strings.Join(instructions, " - "))
}

// Helper types for list items
type categoryItem struct {
	category SettingsCategory
}

func (c categoryItem) FilterValue() string {
	return string(c.category)
}

// Render renders the category item for display in the list
func (c categoryItem) Render() string {
	t := theme.GetManager().Current()
	style := lipgloss.NewStyle().
		Foreground(t.Primary).
		Bold(true).
		Padding(0, 2)

	icon := ""
	return style.Render(icon + " " + string(c.category))
}

type settingItem struct {
	setting *SettingsItem
}

func (s settingItem) FilterValue() string {
	return s.setting.Name
}

// Render renders the setting item for display in the list
func (s settingItem) Render() string {
	t := theme.GetManager().Current()

	// Main setting name and description
	nameStyle := lipgloss.NewStyle().
		Foreground(t.Text).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(t.Secondary).
		Italic(true)

	valueStyle := lipgloss.NewStyle().
		Foreground(t.Accent)

	// Format the value based on type
	var valueStr string
	switch s.setting.Type {
	case TypeBool:
		if s.setting.Value.(bool) {
			valueStr = valueStyle.Render("✔ Yes")
		} else {
			valueStr = valueStyle.Render("o No")
		}
	case TypeInt:
		if s.setting.Unit != "" {
			valueStr = valueStyle.Render(fmt.Sprintf("%d %s", s.setting.Value.(int), s.setting.Unit))
		} else {
			valueStr = valueStyle.Render(fmt.Sprintf("%d", s.setting.Value.(int)))
		}
	case TypeFloat:
		if s.setting.Unit != "" {
			valueStr = valueStyle.Render(fmt.Sprintf("%.2f %s", s.setting.Value.(float64), s.setting.Unit))
		} else {
			valueStr = valueStyle.Render(fmt.Sprintf("%.2f", s.setting.Value.(float64)))
		}
	case TypeDuration:
		duration := s.setting.Value.(time.Duration)
		valueStr = valueStyle.Render(fmt.Sprintf("%.0f seconds", duration.Seconds()))
	case TypeSelect:
		valueStr = valueStyle.Render(s.setting.Value.(string))
	default:
		valueStr = valueStyle.Render(fmt.Sprintf("%v", s.setting.Value))
	}

	// Layout the setting item
	nameLine := nameStyle.Render(s.setting.Name)
	descLine := descStyle.Render(s.setting.Description)
	// Fix: valueStr is already rendered by valueStyle, don't double-render
	valueLine := valueStr

	// Combine lines
	content := nameLine + "\n" + descLine + "\n" + valueLine

	// Add padding and border
	itemStyle := lipgloss.NewStyle().
		Padding(1, 2).
		MarginBottom(1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(t.Secondary)

	return itemStyle.Render(content)
}

// Messages for settings events
type settingsChangedMsg struct {
	saved bool
}

func emitSettingsChangedMsg(saved bool) tea.Cmd {
	return func() tea.Msg {
		return settingsChangedMsg{saved: saved}
	}
}

// GetConfig returns the current configuration
func (m *SettingsModel) GetConfig() *config.Config {
	return m.config
}

// SetConfig updates the configuration and refreshes settings
func (m *SettingsModel) SetConfig(cfg *config.Config) {
	m.config = cfg
	m.refreshSettingsFromConfig()
}

// refreshSettingsFromConfig updates settings from the current config
func (m *SettingsModel) refreshSettingsFromConfig() {
	// Update setting values from config
	for _, setting := range m.settings {
		switch setting.ID {
		case "ui.theme":
			setting.Value = m.config.UI.Theme
		case "ui.font_size":
			setting.Value = m.config.UI.FontSize
		case "ui.animations":
			setting.Value = m.config.UI.Animations
		case "ui.show_line_numbers":
			setting.Value = m.config.UI.ShowLineNumbers
		case "ui.word_wrap":
			setting.Value = m.config.UI.WordWrap
		case "app.auto_save":
			setting.Value = m.config.App.AutoSave
		case "app.auto_save_interval":
			setting.Value = m.config.App.AutoSaveInterval
		case "app.max_recent_files":
			setting.Value = m.config.App.MaxRecentFiles
		case "audio.enabled":
			setting.Value = m.config.Audio.Enabled
		case "audio.playback_gain":
			setting.Value = m.config.Audio.PlaybackGain
		case "ai.enabled":
			setting.Value = m.config.AI.Enabled
		case "ai.temperature":
			setting.Value = m.config.AI.Temperature
		case "dev.debug":
			setting.Value = m.config.Dev.Debug
		}
	}
	m.updateSettingsList()
}

// GetSettingsCount returns the total number of settings
func (m *SettingsModel) GetSettingsCount() int {
	return len(m.settings)
}

// GetSettingsByCategory returns settings for a specific category
func (m *SettingsModel) GetSettingsByCategory(category SettingsCategory) []*SettingsItem {
	var items []*SettingsItem
	for _, setting := range m.settings {
		if setting.Category == category {
			items = append(items, setting)
		}
	}
	return items
}

// ResetToDefaults resets all settings to their default values
func (m *SettingsModel) ResetToDefaults() {
	for _, setting := range m.settings {
		setting.Value = setting.Default
		if setting.Callback != nil {
			setting.Callback(setting.Value)
		}
	}
	m.updateSettingsList()
}
