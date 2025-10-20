package styles

import (
	"testing"
)

func TestFocusColorFunctions(t *testing.T) {
	// Test legacy color functions for backward compatibility
	// These should not panic
	_ = FocusStyle("Test")
	_ = FocusPinkColor("Test Pink")
	_ = FocusBlueColor("Test Blue")
	_ = FocusGreenColor("Test Green")
	_ = FocusPurpleColor("Test Purple")
}