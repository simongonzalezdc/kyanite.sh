package icons

import (
	"strings"

	"github.com/rivo/uniseg"
)

// StringWidth returns the monospace-aware display width of a string.
// This correctly handles East Asian wide characters, combining characters,
// and emoji sequences.
func StringWidth(s string) int {
	return uniseg.StringWidth(s)
}

// Truncate truncates a string to at most maxWidth display cells,
// appending the truncation indicator if truncation occurs.
func Truncate(s string, maxWidth int, truncationIndicator string) string {
	if truncationIndicator == "" {
		truncationIndicator = "..."
	}

	sw := StringWidth(s)
	if sw <= maxWidth {
		return s
	}

	tw := StringWidth(truncationIndicator)
	targetWidth := maxWidth - tw
	if targetWidth < 0 {
		targetWidth = 0
	}

	var width int
	var idx int
	for _, r := range s {
		rw := uniseg.StringWidth(string(r))
		if width+rw > targetWidth {
			break
		}
		width += rw
		idx += len(string(r))
	}

	return s[:idx] + truncationIndicator
}

// TruncateLeft truncates a string from the left to fit within maxWidth columns.
// Useful for paths where the end is more important than the beginning.
func TruncateLeft(s string, maxWidth int, prefix string) string {
	if prefix == "" {
		prefix = "..."
	}

	sw := StringWidth(s)
	if sw <= maxWidth {
		return s
	}

	pw := StringWidth(prefix)
	targetWidth := maxWidth - pw
	if targetWidth < 0 {
		targetWidth = 0
	}

	// Collect grapheme clusters
	var clusters []string
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		clusters = append(clusters, gr.Str())
	}

	// Build from the end to find start index
	var width int
	startIdx := len(clusters)
	for i := len(clusters) - 1; i >= 0; i-- {
		cw := uniseg.StringWidth(clusters[i])
		if width+cw > targetWidth {
			startIdx = i + 1
			break
		}
		width += cw
		startIdx = i
	}

	var result strings.Builder
	if startIdx > 0 {
		result.WriteString(prefix)
	}
	for i := startIdx; i < len(clusters); i++ {
		result.WriteString(clusters[i])
	}
	return result.String()
}

// PadRight pads a string to the specified width with spaces on the right.
// Uses proper Unicode width calculation.
func PadRight(s string, width int) string {
	sw := StringWidth(s)
	if sw >= width {
		return s
	}
	return s + strings.Repeat(" ", width-sw)
}

// PadRightChar pads a string on the right to fill maxWidth display cells
// using the given pad character.
func PadRightChar(s string, maxWidth int, padChar string) string {
	sw := StringWidth(s)
	if sw >= maxWidth {
		return s
	}

	pw := uniseg.StringWidth(padChar)
	if pw == 0 {
		return s
	}

	remaining := maxWidth - sw
	result := s
	for remaining > 0 {
		if remaining >= pw {
			result += padChar
			remaining -= pw
		} else {
			break
		}
	}
	return result
}

// PadLeft pads a string to the specified width with spaces on the left.
// Uses proper Unicode width calculation.
func PadLeft(s string, width int) string {
	sw := StringWidth(s)
	if sw >= width {
		return s
	}
	return strings.Repeat(" ", width-sw) + s
}

// Center centers a string within the specified width using spaces.
// Uses proper Unicode width calculation.
func Center(s string, width int) string {
	sw := StringWidth(s)
	if sw >= width {
		return s
	}

	totalPadding := width - sw
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
}
