package utils

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// StreamText displays text with a streaming effect and colors
func StreamText(text string, color lipgloss.Color, speed time.Duration) {
	// Create streaming style
	streamStyle := lipgloss.NewStyle().
		Foreground(color).
		Bold(true)

	// Stream character by character for words, or word by word for sentences
	words := strings.Fields(text)

	for i, word := range words {
		// Add spacing except for first word
		if i > 0 {
			fmt.Print(" ")
		}

		// Stream each character of the word
		for _, char := range word {
			fmt.Print(streamStyle.Render(string(char)))
			_ = os.Stdout.Sync() // Force immediate output
			time.Sleep(speed)
		}
	}
	fmt.Println()
}

// StreamSentence displays text with clean streaming effect
func StreamSentence(text string, baseColor lipgloss.Color, speed time.Duration) {
	words := strings.Fields(text)

	for i, word := range words {
		if i > 0 {
			fmt.Print(" ")
		}

		streamStyle := lipgloss.NewStyle().
			Foreground(baseColor).
			Bold(true)

		// Stream the word
		for _, char := range word {
			fmt.Print(streamStyle.Render(string(char)))
			os.Stdout.Sync()
			time.Sleep(speed / 2) // Faster for individual chars
		}
	}
	fmt.Println()
}

// StreamWithTypingEffect simulates clean typing
func StreamWithTypingEffect(text string, color lipgloss.Color) {
	words := strings.Fields(text)

	for i, word := range words {
		if i > 0 {
			fmt.Print(" ")
		}

		// Consistent typing speed
		for _, char := range word {
			fmt.Print(lipgloss.NewStyle().Foreground(color).Render(string(char)))
			os.Stdout.Sync()
			time.Sleep(50 * time.Millisecond) // Consistent 50ms delay
		}
	}
	fmt.Println()
}
