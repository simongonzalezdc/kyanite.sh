package design

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestNewTokenSet(t *testing.T) {
	th := DefaultTheme()
	ts := NewTokenSet(th)

	// Verify base tokens have foreground colors set
	if ts.Title.GetForeground() == nil {
		t.Error("Title token has no foreground color")
	}
	if ts.Body.GetForeground() == nil {
		t.Error("Body token has no foreground color")
	}
	if ts.Error.GetForeground() == nil {
		t.Error("Error token has no foreground color")
	}

	// Verify Extensions map is initialized and empty
	if ts.Extensions == nil {
		t.Error("Extensions map is nil")
	}
	if len(ts.Extensions) != 0 {
		t.Errorf("Extensions map should be empty, has %d entries", len(ts.Extensions))
	}
}

func TestTokenSetExtensions(t *testing.T) {
	th := DefaultTheme()
	ts := NewTokenSet(th)

	// Test SetExt
	custom := lipgloss.NewStyle().Bold(true)
	ts = ts.SetExt("StatusBar", custom)

	if len(ts.Extensions) != 1 {
		t.Errorf("Extensions should have 1 entry, has %d", len(ts.Extensions))
	}

	// Test Ext returns a bold style
	got := ts.Ext("StatusBar")
	if !got.GetBold() {
		t.Error("Ext(StatusBar) should be bold")
	}

	// Test Ext with missing key returns a non-bold style
	zero := ts.Ext("NonExistent")
	if zero.GetBold() {
		t.Error("Ext(NonExistent) should not be bold")
	}
}

func TestTokenSetAllThemes(t *testing.T) {
	for _, name := range List() {
		t.Run(name, func(t *testing.T) {
			th := Get(name)
			ts := NewTokenSet(th)
			if ts.Extensions == nil {
				t.Error("Extensions is nil")
			}
		})
	}
}
