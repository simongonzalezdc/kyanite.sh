package spellcheck

import (
	"bufio"
	"os/exec"
	"strings"
)

// Checker provides spell checking functionality
type Checker struct {
	enabled  bool
	language string
}

// NewChecker creates a new spell checker
func NewChecker() *Checker {
	return &Checker{
		enabled:  false,
		language: "en_US",
	}
}

// IsEnabled returns whether spell checking is enabled
func (c *Checker) IsEnabled() bool {
	return c.enabled
}

// Toggle toggles spell checking on/off
func (c *Checker) Toggle() {
	c.enabled = !c.enabled
}

// SetEnabled sets the enabled state
func (c *Checker) SetEnabled(enabled bool) {
	c.enabled = enabled
}

// SetLanguage sets the spell check language
func (c *Checker) SetLanguage(lang string) {
	c.language = lang
}

// CheckWord checks if a single word is spelled correctly
// Returns true if word is correct, false if misspelled or if spell checker is unavailable
func (c *Checker) CheckWord(word string) bool {
	if !c.enabled || word == "" {
		return true
	}

	// Check if aspell is available
	cmd := exec.Command("aspell", "-a", "-l", c.language)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return true // Assume correct if spell checker unavailable
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return true
	}

	if err := cmd.Start(); err != nil {
		return true // Spell checker not available
	}
	defer cmd.Wait()

	// Send word to aspell
	stdin.Write([]byte(word + "\n"))
	stdin.Close()

	// Read response
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Aspell responds with '*' for correct words, '&' for misspelled
		if strings.HasPrefix(line, "*") {
			return true
		}
		if strings.HasPrefix(line, "&") || strings.HasPrefix(line, "#") {
			return false
		}
	}

	return true // Default to correct if uncertain
}

// GetSuggestions returns spelling suggestions for a misspelled word
func (c *Checker) GetSuggestions(word string) []string {
	if !c.enabled || word == "" {
		return nil
	}

	cmd := exec.Command("aspell", "-a", "-l", c.language)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	if err := cmd.Start(); err != nil {
		return nil
	}
	defer cmd.Wait()

	stdin.Write([]byte(word + "\n"))
	stdin.Close()

	var suggestions []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		// Aspell format: & word count offset: suggestion1, suggestion2, ...
		if strings.HasPrefix(line, "&") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				suggestionStr := strings.TrimSpace(parts[1])
				suggestions = strings.Split(suggestionStr, ", ")
				// Limit to first 5 suggestions
				if len(suggestions) > 5 {
					suggestions = suggestions[:5]
				}
			}
		}
	}

	return suggestions
}

// IsAvailable checks if aspell is installed and available
func IsAvailable() bool {
	cmd := exec.Command("aspell", "--version")
	err := cmd.Run()
	return err == nil
}
