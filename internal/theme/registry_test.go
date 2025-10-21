package theme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryContains10Themes(t *testing.T) {
	ids := ListThemes()
	// Using a map to dedupe aliases
	unique := make(map[string]struct{})
	for _, id := range ids {
		unique[id] = struct{}{}
	}
	// Expect exactly 10 themes
	require.Equal(t, 10, len(unique))
}

func TestGetThemeFallback(t *testing.T) {
	t1 := GetTheme("does-not-exist")
	require.NotNil(t, t1)
	// default should be amber-night
	require.Equal(t, Registry[DefaultTheme].Name, t1.Name)
}
