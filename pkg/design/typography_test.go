package design

import (
	"testing"
)

func TestNewTypographySet(t *testing.T) {
	th := DefaultTheme()
	ts := NewTypographySet(th)

	// Title should be bold
	if !ts.Title.GetBold() {
		t.Error("Title should be bold")
	}

	// Heading should be bold
	if !ts.Heading.GetBold() {
		t.Error("Heading should be bold")
	}

	// Body should NOT be bold
	if ts.Body.GetBold() {
		t.Error("Body should not be bold")
	}

	// Muted should be faint
	if !ts.Muted.GetFaint() {
		t.Error("Muted should be faint")
	}
}

func TestTypographySetAllThemes(t *testing.T) {
	for _, name := range List() {
		t.Run(name, func(t *testing.T) {
			th := Get(name)
			ts := NewTypographySet(th)

			if !ts.Title.GetBold() {
				t.Error("Title should always be bold")
			}
			if !ts.Heading.GetBold() {
				t.Error("Heading should always be bold")
			}
			if !ts.Muted.GetFaint() {
				t.Error("Muted should always be faint")
			}
		})
	}
}
