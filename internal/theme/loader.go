package theme

import (
	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/ui/styles"
)

// ApplyThemeToStyles maps a Theme into the ui/styles package variables so existing UI code continues to work.
func ApplyThemeToStyles(t Theme) {
	// Primary / Secondary / Accent
	styles.Primary = t.Primary
	styles.Secondary = t.Secondary
	styles.Accent = t.Accent

	// Functional
	styles.Success = t.Success

	// Background & Text
	styles.Background = t.Background
	styles.TextPrimary = t.Text
	// Keep some reasonable defaults for the rest based on provided colors
	styles.TextSecondary = t.Primary
	styles.TextMuted = t.Secondary
	styles.TextAccent = t.Accent

	// Border & dark shades derive from background (keep existing defaults if previously set)
	styles.Dark1 = t.Background
	styles.Dark2 = t.Background
	styles.Dark3 = t.Background

	styles.BorderColor = t.Primary

	// Recompute some composed styles that depend on color variables.
	// Many styles in the styles package are constructed at init-time using package vars.
	// To ensure UI picks up new colors immediately, reassign core style objects.
	styles.Title = styles.Title.Copy().Foreground(t.Primary).Background(styles.Background)
	styles.Subtitle = styles.Subtitle.Copy().Foreground(styles.TextSecondary)
	styles.Text = styles.Text.Copy().Foreground(t.Text).Background(t.Background)
	styles.Muted = styles.Muted.Copy().Foreground(styles.TextMuted)
	styles.Emphasis = styles.Emphasis.Copy().Foreground(styles.TextAccent)

	styles.Border = styles.Border.Copy().BorderForeground(styles.BorderColor)
	styles.BorderActive = styles.BorderActive.Copy().BorderForeground(styles.Primary)
	styles.BorderThick = styles.BorderThick.Copy().BorderForeground(styles.Accent)

	styles.ButtonPrimary = styles.ButtonPrimary.Copy().Background(t.Primary).Foreground(t.Background)
	styles.ButtonSecondary = styles.ButtonSecondary.Copy().BorderForeground(t.Secondary).Foreground(t.Secondary)
	styles.ButtonAccent = styles.ButtonAccent.Copy().Background(t.Accent).Foreground(t.Background)
	styles.ButtonDisabled = styles.ButtonDisabled.Copy().Background(styles.Dark3).Foreground(styles.TextMuted)

	styles.StatusSuccess = styles.StatusSuccess.Copy().Foreground(t.Success)
	styles.StatusWarning = styles.StatusWarning.Copy().Foreground(styles.Warning)
	styles.StatusError = styles.StatusError.Copy().Foreground(styles.Error)
	styles.StatusInfo = styles.StatusInfo.Copy().Foreground(styles.Info)

	styles.EditorPane = styles.EditorPane.Copy().BorderForeground(t.Primary)
	styles.PreviewPane = styles.PreviewPane.Copy().BorderForeground(t.Secondary)
	styles.StatusBar = styles.StatusBar.Copy().Background(styles.Dark2).Foreground(t.Text)
	styles.Cursor = styles.Cursor.Copy().Background(t.Primary).Foreground(t.Background)
	styles.Divider = styles.Divider.Copy().Foreground(styles.BorderColor)

	styles.H1 = styles.H1.Copy().Foreground(t.Primary)
	styles.H2 = styles.H2.Copy().Foreground(t.Secondary)
	styles.H3 = styles.H3.Copy().Foreground(styles.TextAccent)

	styles.ListItem = styles.ListItem.Copy().Foreground(t.Text)
	styles.ListItemSelected = styles.ListItemSelected.Copy().Background(t.Primary).Foreground(t.Background)

	styles.Card = styles.Card.Copy().BorderForeground(styles.BorderColor)
	styles.CardHighlight = styles.CardHighlight.Copy().BorderForeground(t.Accent)
	styles.CardSuccess = styles.CardSuccess.Copy().BorderForeground(t.Success)
	styles.CardError = styles.CardError.Copy().BorderForeground(styles.Error)

	styles.Tag = styles.Tag.Copy().Background(t.Secondary).Foreground(t.Text)
	styles.Badge = styles.Badge.Copy().Background(t.Accent).Foreground(t.Background)
}

// ApplyThemeByID sets the manager and applies the theme to styles by theme id.
func ApplyThemeByID(id string) {
	GetManager().SetTheme(id)
	ApplyThemeToStyles(GetTheme(id))
}

// InitFromConfig initializes global theme manager and applies theme from config.
func InitFromConfig(cfg *config.Config) {
	if cfg == nil {
		ApplyThemeByID(DefaultTheme)
		return
	}

	// Apply config to manager and styles. If UI.Theme is empty or invalid, fall back to default.
	GetManager().ApplyConfig(cfg)
	themeID := cfg.UI.Theme
	if themeID == "" {
		themeID = DefaultTheme
	}
	ApplyThemeByID(themeID)
}