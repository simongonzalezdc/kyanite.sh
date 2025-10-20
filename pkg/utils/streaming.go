package utils

import (
	"fmt"
	"strings"
	"time"
	"os"
	
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
			os.Stdout.Sync() // Force immediate output
			time.Sleep(speed)
		}
	}
	fmt.Println()
}

// StreamSentence displays full sentences with streaming effect
func StreamSentence(text string, baseColor lipgloss.Color, speed time.Duration) {
	// Rainbow colors for streaming effect
	colors := []lipgloss.Color{
		lipgloss.Color("#FF71CE"), // Pink
		lipgloss.Color("#00FFFF"), // Cyan  
		lipgloss.Color("#00FF66"), // Green
		lipgloss.Color("#FEC107"), // Yellow
		lipgloss.Color("#B967C7"), // Purple
	}
	
	words := strings.Fields(text)
	
	for i, word := range words {
		if i > 0 {
			fmt.Print(" ")
		}
		
		// Cycle through colors
		color := colors[i%len(colors)]
		streamStyle := lipgloss.NewStyle().
			Foreground(color).
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

// StreamWithTypingEffect simulates AI typing with variable speed
func StreamWithTypingEffect(text string, color lipgloss.Color) {
	// Simulate realistic typing with variable delays
	words := strings.Fields(text)
	
	for i, word := range words {
		if i > 0 {
			fmt.Print(" ")
		}
		
		// Variable typing speed for realism
		for _, char := range word {
			// Random delay between characters (20-80ms)
			delay := time.Duration(20+int(time.Now().UnixNano()%60)) * time.Millisecond
			
			fmt.Print(lipgloss.NewStyle().
				Foreground(color).
				Bold(true).
				Render(string(char)))
			os.Stdout.Sync()
			time.Sleep(delay)
		}
	}
	fmt.Println()
}
