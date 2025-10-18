package theme

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistryContains13Themes(t *testing.T) {
	ids := ListThemes()
	// Using a map to dedupe aliases
	unique := make(map[string]struct{})
	for _, id := range ids {
		unique[id] = struct{}{}
	}
	// Expect at least 13 keys including aliases; ensure primary purple-jazz exists
	_, ok := Registry["purple-jazz"]
	require.True(t, ok, "purple-jazz theme must be present")
	// Check there are at least 13 distinct entries (including aliases)
	require.GreaterOrEqual(t, len(unique), 13)
}

func TestGetThemeFallback(t *testing.T) {
	t1 := GetTheme("does-not-exist")
	require.NotNil(t, t1)
	// default should be amber-night
	require.Equal(t, Registry[DefaultTheme].Name, t1.Name)
}