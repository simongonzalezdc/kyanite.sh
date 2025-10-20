package gum

import (
	"fmt"
	"os/exec"
	"strings"
)

// Choose prompts user to select from a list of options
func Choose(options []string, prompt string) string {
	if prompt == "" {
		prompt = "Choose an option:"
	}

	args := []string{"choose", prompt}
	args = append(args, options...)

	cmd := exec.Command("gum", args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// Confirm asks user for yes/no confirmation
func Confirm(prompt string) bool {
	if prompt == "" {
		prompt = "Confirm?"
	}

	cmd := exec.Command("gum", "confirm", prompt)
	err := cmd.Run()
	return err == nil
}

// Input prompts user for text input
func Input(prompt string) string {
	if prompt == "" {
		prompt = "Enter text:"
	}

	cmd := exec.Command("gum", "input", "--placeholder", prompt)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// Filter allows fuzzy filtering of options
func Filter(options []string, prompt string) string {
	if prompt == "" {
		prompt = "Filter options:"
	}

	cmd := exec.Command("gum", "filter", "--placeholder", prompt)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return ""
	}

	go func() {
		defer func() { _ = stdin.Close() }()
		for _, option := range options {
			_, _ = stdin.Write([]byte(option + "\n"))
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(output))
}

// MultiSelect allows selecting multiple options
func MultiSelect(options []string, prompt string, limit int) []string {
	if prompt == "" {
		prompt = "Select options:"
	}

	args := []string{"filter", "--placeholder", prompt, "--no-limit"}
	if limit > 0 {
		args = []string{"filter", "--placeholder", prompt, "--limit", fmt.Sprintf("%d", limit)}
	}

	cmd := exec.Command("gum", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	go func() {
		defer func() { _ = stdin.Close() }()
		for _, option := range options {
			_, _ = stdin.Write([]byte(option + "\n"))
		}
	}()

	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	selected := strings.Split(strings.TrimSpace(string(output)), "\n")
	result := make([]string, 0, len(selected))
	for _, item := range selected {
		if strings.TrimSpace(item) != "" {
			result = append(result, strings.TrimSpace(item))
		}
	}

	return result
}

// Spin shows a spinner while running a command
func Spin(title string, command string, args ...string) error {
	args = append([]string{"spin", "--title", title, command}, args...)
	cmd := exec.Command("gum", args...)
	return cmd.Run()
}

// Style applies styling to text
func Style(text string, foreground string, background string, bold bool) string {
	args := []string{"style"}

	if foreground != "" {
		args = append(args, "--foreground", foreground)
	}
	if background != "" {
		args = append(args, "--background", background)
	}
	if bold {
		args = append(args, "--bold")
	}

	args = append(args, text)

	cmd := exec.Command("gum", args...)
	output, err := cmd.Output()
	if err != nil {
		return text
	}

	return strings.TrimSpace(string(output))
}

// Check if gum is available
func IsAvailable() bool {
	// Try to find gum in common paths
	paths := []string{
		"gum",
		// Unix/Linux paths
		"/usr/local/bin/gum",
		"/usr/bin/gum",
		"/bin/gum",
		// macOS paths
		"/opt/homebrew/bin/gum",
		// Windows paths
		"C:\\Program Files\\gum\\gum.exe",
		"C:\\Program Files (x86)\\gum\\gum.exe",
		"C:\\Users\\Simon\\AppData\\Local\\Microsoft\\WinGet\\Packages\\charmbracelet.gum\\gum.exe",
	}

	for _, path := range paths {
		if _, err := exec.LookPath(path); err == nil {
			return true
		}
	}
	return false
}
