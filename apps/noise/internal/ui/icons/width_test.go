package icons

import (
	"strings"
	"testing"
)

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty string", "", 0},
		{"ASCII only", "hello", 5},
		{"ASCII with spaces", "hello world", 11},
		{"CJK characters", "世界", 4},           // Each CJK char is 2 columns wide
		{"mixed ASCII and CJK", "Hello世界", 9}, // 5 + 4
		{"emoji", "👍", 2},                     // Most emoji are 2 columns
		{"combining characters", "é", 1},      // e + combining acute is 1 column
		{"tabs", "\t", 0},                     // Tab has 0 width in uniseg (terminal-specific)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringWidth(tt.input)
			if result != tt.expected {
				t.Errorf("StringWidth(%q) = %d, expected %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxWidth  int
		suffix    string
		minLength int // Minimum expected length (not exact due to multi-byte)
	}{
		{"no truncation needed", "hello", 10, "...", 5},
		{"truncate ASCII", "hello world", 8, "...", 5}, // "hello..." = 8
		{"truncate with CJK", "Hello世界", 7, "...", 4},  // Should not break CJK
		{"empty suffix", "hello world", 5, "", 5},
		{"very short max", "hello", 3, "...", 3},
		{"exact width", "hello", 5, "...", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Truncate(tt.input, tt.maxWidth, tt.suffix)
			resultWidth := StringWidth(result)

			if resultWidth > tt.maxWidth {
				t.Errorf("Truncate(%q, %d, %q) result width %d exceeds max %d",
					tt.input, tt.maxWidth, tt.suffix, resultWidth, tt.maxWidth)
			}

			if len(result) < tt.minLength {
				t.Errorf("Truncate(%q, %d, %q) result too short: %q",
					tt.input, tt.maxWidth, tt.suffix, result)
			}
		})
	}
}

func TestTruncatePreservesGraphemeClusters(t *testing.T) {
	// Test that truncation doesn't break multi-byte characters
	input := "Hello世界Test"

	for width := 1; width <= StringWidth(input); width++ {
		result := Truncate(input, width, "")
		// The result should be valid UTF-8 and not break mid-character
		for _, r := range result {
			if r == 0xFFFD { // Unicode replacement character
				t.Errorf("Truncate(%q, %d, \"\") produced invalid UTF-8", input, width)
			}
		}
	}
}

func TestTruncateLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWidth int
		prefix   string
		contains string // The result should contain this
	}{
		{"no truncation", "hello", 10, "...", "hello"},
		{"truncate path", "/usr/local/bin/app", 12, "...", "app"},
		{"with prefix", "abcdefghij", 6, "...", "hij"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateLeft(tt.input, tt.maxWidth, tt.prefix)
			resultWidth := StringWidth(result)

			if resultWidth > tt.maxWidth {
				t.Errorf("TruncateLeft(%q, %d, %q) result width %d exceeds max %d",
					tt.input, tt.maxWidth, tt.prefix, resultWidth, tt.maxWidth)
			}

			if !strings.Contains(result, tt.contains) {
				t.Errorf("TruncateLeft(%q, %d, %q) = %q, should contain %q",
					tt.input, tt.maxWidth, tt.prefix, result, tt.contains)
			}
		})
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected int // Expected display width of result
	}{
		{"pad ASCII", "hi", 10, 10},
		{"pad CJK", "世界", 10, 10},        // CJK is 4 wide, should pad to 10
		{"no pad needed", "hello", 3, 5}, // Already wider, no change
		{"exact width", "hello", 5, 5},
		{"empty string", "", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadRight(tt.input, tt.width)
			resultWidth := StringWidth(result)

			if resultWidth != tt.expected {
				t.Errorf("PadRight(%q, %d) width = %d, expected %d",
					tt.input, tt.width, resultWidth, tt.expected)
			}

			// Result should start with input
			if !strings.HasPrefix(result, tt.input) {
				t.Errorf("PadRight(%q, %d) = %q, should start with input",
					tt.input, tt.width, result)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected int
	}{
		{"pad ASCII", "hi", 10, 10},
		{"pad CJK", "世界", 10, 10},
		{"no pad needed", "hello", 3, 5},
		{"exact width", "hello", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PadLeft(tt.input, tt.width)
			resultWidth := StringWidth(result)

			if resultWidth != tt.expected {
				t.Errorf("PadLeft(%q, %d) width = %d, expected %d",
					tt.input, tt.width, resultWidth, tt.expected)
			}

			// Result should end with input
			if !strings.HasSuffix(result, tt.input) {
				t.Errorf("PadLeft(%q, %d) = %q, should end with input",
					tt.input, tt.width, result)
			}
		})
	}
}

func TestCenter(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		width    int
		expected int
	}{
		{"center ASCII", "hi", 10, 10},
		{"center CJK", "世界", 10, 10},
		{"no pad needed", "hello", 3, 5},
		{"odd padding", "hi", 9, 9}, // 3 left, 4 right or vice versa
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Center(tt.input, tt.width)
			resultWidth := StringWidth(result)

			if resultWidth != tt.expected {
				t.Errorf("Center(%q, %d) width = %d, expected %d",
					tt.input, tt.width, resultWidth, tt.expected)
			}

			// Result should contain input
			if !strings.Contains(result, tt.input) {
				t.Errorf("Center(%q, %d) = %q, should contain input",
					tt.input, tt.width, result)
			}
		})
	}
}

func TestStringWidthVsLen(t *testing.T) {
	// Demonstrate that StringWidth differs from len() for multi-byte strings
	tests := []struct {
		input       string
		lenResult   int
		widthResult int
	}{
		{"hello", 5, 5},    // Same for ASCII
		{"世界", 6, 4},       // len=6 bytes, width=4 columns
		{"Hello世界", 11, 9}, // len=11 bytes, width=9 columns
		{"é", 2, 1},        // len=2 bytes (e + combining), width=1
	}

	for _, tt := range tests {
		lenR := len(tt.input)
		widthR := StringWidth(tt.input)

		if lenR != tt.lenResult {
			t.Errorf("len(%q) = %d, expected %d", tt.input, lenR, tt.lenResult)
		}
		if widthR != tt.widthResult {
			t.Errorf("StringWidth(%q) = %d, expected %d", tt.input, widthR, tt.widthResult)
		}
	}
}

func BenchmarkStringWidth(b *testing.B) {
	inputs := []string{
		"hello world",
		"Hello世界Test",
		strings.Repeat("Hello世界", 100),
	}

	for _, input := range inputs {
		b.Run(input[:min(20, len(input))], func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				StringWidth(input)
			}
		})
	}
}

func BenchmarkTruncate(b *testing.B) {
	input := strings.Repeat("Hello世界Test", 50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Truncate(input, 100, "...")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
