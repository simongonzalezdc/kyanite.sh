package icons

import (
	"strings"

	"github.com/rivo/uniseg"
)

// StringWidth returns the display width of a string in terminal columns.
// This properly handles Unicode characters, emoji, and combining characters.
func StringWidth(s string) int {
	return uniseg.StringWidth(s)
}

// Truncate truncates a string to fit within maxWidth columns,
// respecting grapheme cluster boundaries to avoid breaking multi-byte characters.
// If the string is truncated, suffix is appended (default "...").
func Truncate(s string, maxWidth int, suffix string) string {
	if suffix == "" {
		suffix = "..."
	}

	currentWidth := StringWidth(s)
	if currentWidth <= maxWidth {
		return s
	}

	suffixWidth := StringWidth(suffix)
	targetWidth := maxWidth - suffixWidth
	if targetWidth < 0 {
		targetWidth = 0
	}

	var result strings.Builder
	var width int
	gr := uniseg.NewGraphemes(s)

	for gr.Next() {
		cluster := gr.Str()
		clusterWidth := uniseg.StringWidth(cluster)

		if width+clusterWidth > targetWidth {
			break
		}

		width += clusterWidth
		result.WriteString(cluster)
	}

	if width < currentWidth {
		result.WriteString(suffix)
	}

	return result.String()
}

// TruncateLeft truncates a string from the left to fit within maxWidth columns.
// Useful for paths where the end is more important than the beginning.
func TruncateLeft(s string, maxWidth int, prefix string) string {
	if prefix == "" {
		prefix = "..."
	}

	currentWidth := StringWidth(s)
	if currentWidth <= maxWidth {
		return s
	}

	prefixWidth := StringWidth(prefix)
	targetWidth := maxWidth - prefixWidth
	if targetWidth < 0 {
		targetWidth = 0
	}

	// Collect all grapheme clusters
	var clusters []string
	gr := uniseg.NewGraphemes(s)
	for gr.Next() {
		clusters = append(clusters, gr.Str())
	}

	// Build from the end
	var result strings.Builder
	var width int

	for i := len(clusters) - 1; i >= 0; i-- {
		cluster := clusters[i]
		clusterWidth := uniseg.StringWidth(cluster)

		if width+clusterWidth > targetWidth {
			break
		}

		width += clusterWidth
	}

	// Calculate start index
	startIdx := len(clusters)
	width = 0
	for i := len(clusters) - 1; i >= 0; i-- {
		clusterWidth := uniseg.StringWidth(clusters[i])
		if width+clusterWidth > targetWidth {
			startIdx = i + 1
			break
		}
		width += clusterWidth
		startIdx = i
	}

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
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := strings.Repeat(" ", width-currentWidth)
	return s + padding
}

// PadLeft pads a string to the specified width with spaces on the left.
// Uses proper Unicode width calculation.
func PadLeft(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	padding := strings.Repeat(" ", width-currentWidth)
	return padding + s
}

// Center centers a string within the specified width using spaces.
// Uses proper Unicode width calculation.
func Center(s string, width int) string {
	currentWidth := StringWidth(s)
	if currentWidth >= width {
		return s
	}

	totalPadding := width - currentWidth
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	return strings.Repeat(" ", leftPadding) + s + strings.Repeat(" ", rightPadding)
}
