package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ThemeManager handles theme switching and persistence
type ThemeManager struct {
	currentTheme  *Theme
	themeIndex    int
	themeFilePath string
}

// NewThemeManager creates a new theme manager
func NewThemeManager(themeFilePath string) *ThemeManager {
	return &ThemeManager{
		currentTheme:  &MidnightJazzTheme, // Default theme
		themeIndex:    0,
		themeFilePath: themeFilePath,
	}
}

// GetCurrentTheme returns the current theme
func (tm *ThemeManager) GetCurrentTheme() *Theme {
	return tm.currentTheme
}

// SetTheme sets the current theme by name
func (tm *ThemeManager) SetTheme(themeName string) error {
	for i, theme := range AllThemes {
		if theme.Name == themeName {
			tm.currentTheme = &AllThemes[i]
			tm.themeIndex = i
			return tm.saveThemePreference()
		}
	}
	return fmt.Errorf("theme '%s' not found", themeName)
}

// SetThemeByIndex sets the current theme by index
func (tm *ThemeManager) SetThemeByIndex(index int) error {
	if index < 0 || index >= len(AllThemes) {
		return fmt.Errorf("theme index %d out of range", index)
	}
	tm.currentTheme = &AllThemes[index]
	tm.themeIndex = index
	return tm.saveThemePreference()
}

// NextTheme switches to the next theme
func (tm *ThemeManager) NextTheme() error {
	nextIndex := (tm.themeIndex + 1) % len(AllThemes)
	return tm.SetThemeByIndex(nextIndex)
}

// PreviousTheme switches to the previous theme
func (tm *ThemeManager) PreviousTheme() error {
	prevIndex := (tm.themeIndex - 1 + len(AllThemes)) % len(AllThemes)
	return tm.SetThemeByIndex(prevIndex)
}

// GetAllThemes returns all available themes
func (tm *ThemeManager) GetAllThemes() []Theme {
	return AllThemes
}

// GetThemeNames returns all theme names
func (tm *ThemeManager) GetThemeNames() []string {
	names := make([]string, len(AllThemes))
	for i, theme := range AllThemes {
		names[i] = theme.Name
	}
	return names
}

// GetThemeDescription returns the description of the current theme
func (tm *ThemeManager) GetThemeDescription() string {
	return tm.currentTheme.Description
}

// ApplyTheme applies the current theme to all global style variables
func (tm *ThemeManager) ApplyTheme() {
	colors := tm.currentTheme.Colors
	
	// Primary Colors
	Primary = colors.Primary
	Secondary = colors.Secondary
	Accent = colors.Accent
	
	// Functional Colors
	Success = colors.Success
	Warning = colors.Warning
	Error = colors.Error
	Info = colors.Info
	
	// Background & Text Colors
	Background = colors.Background
	TextPrimary = colors.TextPrimary
	TextSecondary = colors.TextSecondary
	TextMuted = colors.TextMuted
	TextAccent = colors.TextAccent
	
	// Extended Palette
	Dark1 = colors.Dark1
	Dark2 = colors.Dark2
	Dark3 = colors.Dark3
	
	// Border color
	BorderColor = colors.Dark3
	
	// Update all styles that depend on these colors
	updateStyles()
}

// saveThemePreference saves the current theme preference to file
func (tm *ThemeManager) saveThemePreference() error {
	if tm.themeFilePath == "" {
		return nil // No file path configured, skip saving
	}
	
	// Ensure directory exists
	dir := filepath.Dir(tm.themeFilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create theme directory: %w", err)
	}
	
	// Create theme preference data
	pref := ThemePreference{
		ThemeName: tm.currentTheme.Name,
		ThemeIndex: tm.themeIndex,
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(pref, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal theme preference: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(tm.themeFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write theme preference: %w", err)
	}
	
	return nil
}

// loadThemePreference loads the theme preference from file
func (tm *ThemeManager) loadThemePreference() error {
	if tm.themeFilePath == "" {
		return nil // No file path configured, skip loading
	}
	
	// Check if file exists
	if _, err := os.Stat(tm.themeFilePath); os.IsNotExist(err) {
		return nil // File doesn't exist, use default theme
	}
	
	// Read file
	data, err := os.ReadFile(tm.themeFilePath)
	if err != nil {
		return fmt.Errorf("failed to read theme preference: %w", err)
	}
	
	// Unmarshal JSON
	var pref ThemePreference
	if err := json.Unmarshal(data, &pref); err != nil {
		return fmt.Errorf("failed to unmarshal theme preference: %w", err)
	}
	
	// Apply theme if found
	if pref.ThemeName != "" {
		if err := tm.SetTheme(pref.ThemeName); err != nil {
			// If theme not found, try index
			if pref.ThemeIndex >= 0 && pref.ThemeIndex < len(AllThemes) {
				tm.SetThemeByIndex(pref.ThemeIndex)
			}
		}
	}
	
	return nil
}

// ThemePreference represents a saved theme preference
type ThemePreference struct {
	ThemeName string `json:"theme_name"`
	ThemeIndex int   `json:"theme_index"`
}

// updateStyles updates all global styles to use the current theme colors
func updateStyles() {
	// Base Styles
	Title = Title.Foreground(Primary)
	Subtitle = Subtitle.Foreground(TextSecondary)
	Text = Text.Foreground(TextPrimary)
	Muted = Muted.Foreground(TextMuted)
	Emphasis = Emphasis.Foreground(TextAccent)
	
	// Border Styles
	Border = Border.BorderForeground(BorderColor)
	BorderActive = BorderActive.BorderForeground(Primary)
	BorderThick = BorderThick.BorderForeground(Accent)
	
	// Button Styles
	ButtonPrimary = ButtonPrimary.Background(Primary).Foreground(Background)
	ButtonSecondary = ButtonSecondary.BorderForeground(Secondary).Foreground(Secondary)
	ButtonAccent = ButtonAccent.Background(Accent).Foreground(Background)
	ButtonDisabled = ButtonDisabled.Background(Dark3).Foreground(TextMuted)
	
	// Status Styles
	StatusSuccess = StatusSuccess.Foreground(Success)
	StatusWarning = StatusWarning.Foreground(Warning)
	StatusError = StatusError.Foreground(Error)
	StatusInfo = StatusInfo.Foreground(Info)
	
	// Editor Component Styles
	EditorPane = EditorPane.BorderForeground(Primary)
	PreviewPane = PreviewPane.BorderForeground(Secondary)
	StatusBar = StatusBar.Background(Dark2).Foreground(TextPrimary)
	Cursor = Cursor.Background(Primary).Foreground(Background)
	Divider = Divider.Foreground(BorderColor)
	
	// Typography Styles
	H1 = H1.Foreground(Primary)
	H2 = H2.Foreground(Secondary)
	H3 = H3.Foreground(TextAccent)
	Code = Code.Foreground(Accent).Background(Dark2)
	Quote = Quote.Foreground(TextSecondary).BorderLeft(true).BorderForeground(Primary)
	
	// List & Menu Styles
	ListItem = ListItem.Foreground(TextPrimary)
	ListItemSelected = ListItemSelected.Background(Primary).Foreground(Background)
	
	// Card & Panel Styles
	Card = Card.BorderForeground(BorderColor)
	CardHighlight = CardHighlight.BorderForeground(Accent)
	CardSuccess = CardSuccess.BorderForeground(Success)
	CardError = CardError.BorderForeground(Error)
	
	// Badge & Tag Styles
	Tag = Tag.Background(Secondary).Foreground(TextPrimary)
	Badge = Badge.Background(Accent).Foreground(Background)
}

// Init initializes the theme manager and loads saved preferences
func (tm *ThemeManager) Init() {
	// Load saved theme preference
	if err := tm.loadThemePreference(); err != nil {
		// If loading fails, use default theme
		tm.currentTheme = &MidnightJazzTheme
		tm.themeIndex = 0
	}
	
	// Apply the current theme
	tm.ApplyTheme()
}