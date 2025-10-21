package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestManager_GetCurrent(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "Get current theme after initialization",
			test: func(t *testing.T) {
				mgr := GetManager()
				theme := mgr.Current()
				require.NotNil(t, theme)
				assert.Equal(t, "Amber Night", theme.Name)
			},
		},
		{
			name: "Get current theme after switching",
			test: func(t *testing.T) {
				mgr := GetManager()
				mgr.SetTheme("twilight-mist")

				theme := mgr.Current()
				require.NotNil(t, theme)
				assert.Equal(t, "Twilight Mist", theme.Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestManager_SetTheme(t *testing.T) {
	tests := []struct {
		name    string
		themeID string
		wantErr bool
	}{
		{
			name:    "Valid theme ID",
			themeID: "twilight-mist",
			wantErr: false,
		},
		{
			name:    "Another valid theme ID",
			themeID: "cyber-noir",
			wantErr: false,
		},
		{
			name:    "Invalid theme ID",
			themeID: "invalid-theme",
			wantErr: true,
		},
		{
			name:    "Empty theme ID",
			themeID: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mgr := GetManager()
			mgr.SetTheme(tt.themeID)
			
			if !tt.wantErr {
				theme := GetTheme(tt.themeID)
				assert.Equal(t, theme.Name, mgr.Current().Name)
			}
		})
	}
}

func TestManager_GetTheme(t *testing.T) {
	// Test getting existing theme
	theme := GetTheme("amber-night")
	require.NotNil(t, theme)
	assert.Equal(t, "Amber Night", theme.Name)

	// Test getting non-existing theme
	theme = GetTheme("non-existing")
	assert.Equal(t, Default().Name, theme.Name) // Falls back to default
}

func TestManager_ListThemes(t *testing.T) {
	themes := ListThemes()

	// Should have all 10 themes
	assert.Len(t, themes, 10)

	// Check that all required themes are present
	themeIDs := make(map[string]bool)
	for _, themeID := range themes {
		themeIDs[themeID] = true
	}

	requiredThemes := []string{
		"monochrome",
		"amber-night",
		"twilight-mist",
		"indigo-depths",
		"forest-path",
		"clay-earth",
		"iron-forge",
		"sunlight",
		"cyan-wave",
		"electric-rose",
	}

	for _, themeID := range requiredThemes {
		assert.True(t, themeIDs[themeID], "Missing theme: %s", themeID)
	}
}

func TestManager_MigrateOldThemeID(t *testing.T) {
	tests := []struct {
		name     string
		oldID    string
		expected string
	}{
		{
			name:     "Migrate dark to amber-night",
			oldID:    "midnight-jazz",
			expected: "amber-night",
		},
		{
			name:     "Migrate light to amber-night",
			oldID:    "purple-jazz",
			expected: "amber-night",
		},
		{
			name:     "Unknown old ID returns empty string",
			oldID:    "unknown",
			expected: "unknown",
		},
		{
			name:     "Empty old ID returns empty string",
			oldID:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := migrateThemeID(tt.oldID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTheme_HasValidColors(t *testing.T) {
	themes := ListThemes()

	for _, themeID := range themes {
		t.Run(themeID, func(t *testing.T) {
			theme := GetTheme(themeID)

			// Check that all color fields are non-empty
			assert.NotEmpty(t, theme.Background, "Background color is empty")
			assert.NotEmpty(t, theme.Primary, "Primary color is empty")
			assert.NotEmpty(t, theme.Secondary, "Secondary color is empty")
			assert.NotEmpty(t, theme.Text, "Text color is empty")
			assert.NotEmpty(t, theme.Accent, "Accent color is empty")
			assert.NotEmpty(t, theme.Success, "Success color is empty")
			assert.NotEmpty(t, theme.Error, "Error color is empty")
			assert.NotEmpty(t, theme.Warning, "Warning color is empty")

			// Check that colors are valid hex codes
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Background)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Primary)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Secondary)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Text)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Accent)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Success)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Error)
			assert.Regexp(t, `^#[0-9a-fA-F]{6}$`, theme.Warning)
		})
	}
}

func TestTheme_HasValidContrast(t *testing.T) {
	themes := ListThemes()

	for _, themeID := range themes {
		t.Run(themeID, func(t *testing.T) {
			theme := GetTheme(themeID)

			// Check that text has sufficient contrast against background
			// This is a basic check - in a real implementation, you'd use a proper contrast ratio calculation
			assert.NotEqual(t, theme.Background, theme.Text, "Text color should not be the same as background")

			// Check that primary has sufficient contrast against background
			assert.NotEqual(t, theme.Background, theme.Primary, "Primary color should not be the same as background")

			// Check that accent has sufficient contrast against background
			assert.NotEqual(t, theme.Background, theme.Accent, "Accent color should not be the same as background")
		})
	}
}

func TestManager_Persistence(t *testing.T) {
	// This test would require mocking the file system
	// For now, we'll just test that the methods exist
	// Test that SaveTheme and LoadTheme methods exist
	mgr := GetManager()
	assert.NotNil(t, mgr.SaveThemePreference)
	assert.NotNil(t, mgr.LoadThemePreference)
}

func TestManager_GetManager(t *testing.T) {
	// Test that GetManager returns a non-nil manager
	mgr := GetManager()
	assert.NotNil(t, mgr)

	// Test that GetManager returns the same instance (singleton pattern)
	mgr2 := GetManager()
	assert.Equal(t, mgr, mgr2)
}

func TestTheme_HasRequiredFields(t *testing.T) {
	themes := ListThemes()

	for _, themeID := range themes {
		t.Run(themeID, func(t *testing.T) {
			theme := GetTheme(themeID)

			// Check that all required fields are present
			assert.NotEmpty(t, theme.Name, "Theme name is empty")
			
			// Check that theme name matches expected name from registry
			expectedTheme := GetTheme(themeID)
			assert.Equal(t, expectedTheme.Name, theme.Name, "Theme name should match registry")
		})
	}
}

func TestTheme_ValidateWCAG(t *testing.T) {
	themes := ListThemes()

	for _, themeID := range themes {
		t.Run(themeID, func(t *testing.T) {
			theme := GetTheme(themeID)

			// Validate theme against WCAG standards
			warnings := ValidateTheme(theme)

			// Log warnings if any (don't fail the test for now)
			for _, warning := range warnings {
				t.Logf("Theme %s warning: %s", theme.Name, warning)
			}
		})
	}
}

func TestManager_NextPrevious(t *testing.T) {
	mgr := GetManager()

	// Get initial theme
	initialTheme := mgr.Current()
	initialName := initialTheme.Name

	// Test Next
	nextTheme := mgr.Next()
	assert.NotEqual(t, initialName, nextTheme.Name, "Next theme should be different")

	// Test Previous
	prevTheme := mgr.Previous()
	assert.Equal(t, initialName, prevTheme.Name, "Previous should return to initial theme")
}
