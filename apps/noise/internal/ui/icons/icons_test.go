package icons

import (
	"testing"
)

func TestIconSetStyles(t *testing.T) {
	// Test that all three styles are different
	if ASCII.Success == Unicode.Success {
		t.Error("ASCII and Unicode Success icons should be different")
	}
	if ASCII.Success == NerdFont.Success {
		t.Error("ASCII and NerdFont Success icons should be different")
	}
	if Unicode.Success == NerdFont.Success {
		t.Error("Unicode and NerdFont Success icons should be different")
	}
}

func TestASCIIIconsAreASCII(t *testing.T) {
	// Verify ASCII icons only contain ASCII characters
	asciiIcons := []string{
		ASCII.Success, ASCII.Error, ASCII.Warning, ASCII.Info,
		ASCII.CheckOn, ASCII.CheckOff, ASCII.RadioOn, ASCII.RadioOff,
		ASCII.ArrowLeft, ASCII.ArrowRight, ASCII.ArrowUp, ASCII.ArrowDown,
		ASCII.Bullet, ASCII.Separator, ASCII.ProgressFull, ASCII.ProgressEmpty,
		ASCII.Online, ASCII.Away, ASCII.Busy, ASCII.Offline,
		ASCII.Music, ASCII.Folder, ASCII.File, ASCII.Settings,
	}

	for _, icon := range asciiIcons {
		for _, r := range icon {
			if r > 127 {
				t.Errorf("ASCII icon contains non-ASCII character: %q (rune %d)", icon, r)
				break
			}
		}
	}
}

func TestSetStyle(t *testing.T) {
	// Save original state
	originalStyle := GetStyle()
	defer SetStyle(originalStyle)

	// Test ASCII
	SetStyle(StyleASCII)
	if GetStyle() != StyleASCII {
		t.Errorf("Expected style %s, got %s", StyleASCII, GetStyle())
	}
	if Current().Success != ASCII.Success {
		t.Errorf("Expected ASCII Success icon, got %s", Current().Success)
	}

	// Test Unicode
	SetStyle(StyleUnicode)
	if GetStyle() != StyleUnicode {
		t.Errorf("Expected style %s, got %s", StyleUnicode, GetStyle())
	}
	if Current().Success != Unicode.Success {
		t.Errorf("Expected Unicode Success icon, got %s", Current().Success)
	}

	// Test NerdFont
	SetStyle(StyleNerdFont)
	if GetStyle() != StyleNerdFont {
		t.Errorf("Expected style %s, got %s", StyleNerdFont, GetStyle())
	}
	if Current().Success != NerdFont.Success {
		t.Errorf("Expected NerdFont Success icon, got %s", Current().Success)
	}
}

func TestSetStyleInvalidDefaultsToASCII(t *testing.T) {
	// Save original state
	originalStyle := GetStyle()
	defer SetStyle(originalStyle)

	// Test invalid style defaults to ASCII
	SetStyle(IconStyle("invalid"))
	if GetStyle() != IconStyle("invalid") {
		// The style string is stored as-is
	}
	if Current().Success != ASCII.Success {
		t.Errorf("Invalid style should default to ASCII icons")
	}
}

func TestGetAndCurrent(t *testing.T) {
	// Get() and Current() should return the same thing
	get := Get()
	current := Current()

	if get != current {
		t.Error("Get() and Current() should return the same IconSet")
	}
}

func TestIconSetCompleteness(t *testing.T) {
	// Verify all icon sets have non-empty values for all fields
	sets := map[string]*IconSet{
		"ASCII":    &ASCII,
		"Unicode":  &Unicode,
		"NerdFont": &NerdFont,
	}

	for name, set := range sets {
		if set.Success == "" {
			t.Errorf("%s: Success icon is empty", name)
		}
		if set.Error == "" {
			t.Errorf("%s: Error icon is empty", name)
		}
		if set.Warning == "" {
			t.Errorf("%s: Warning icon is empty", name)
		}
		if set.Info == "" {
			t.Errorf("%s: Info icon is empty", name)
		}
		if set.CheckOn == "" {
			t.Errorf("%s: CheckOn icon is empty", name)
		}
		if set.CheckOff == "" {
			t.Errorf("%s: CheckOff icon is empty", name)
		}
		if set.ArrowLeft == "" {
			t.Errorf("%s: ArrowLeft icon is empty", name)
		}
		if set.ArrowRight == "" {
			t.Errorf("%s: ArrowRight icon is empty", name)
		}
		if set.Bullet == "" {
			t.Errorf("%s: Bullet icon is empty", name)
		}
		if set.Music == "" {
			t.Errorf("%s: Music icon is empty", name)
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	// Test that concurrent style changes don't cause races
	// This test should be run with -race flag
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				SetStyle(StyleASCII)
				_ = Current().Success
				SetStyle(StyleUnicode)
				_ = GetStyle()
				SetStyle(StyleNerdFont)
				_ = Get().Error
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestIconStyleConstants(t *testing.T) {
	// Verify style constants have expected values
	if StyleASCII != "ascii" {
		t.Errorf("StyleASCII should be 'ascii', got %s", StyleASCII)
	}
	if StyleUnicode != "unicode" {
		t.Errorf("StyleUnicode should be 'unicode', got %s", StyleUnicode)
	}
	if StyleNerdFont != "nerdfonts" {
		t.Errorf("StyleNerdFont should be 'nerdfonts', got %s", StyleNerdFont)
	}
}
