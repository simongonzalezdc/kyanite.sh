package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// Copy copies text to the system clipboard
func Copy(text string) error {
	switch runtime.GOOS {
	case "darwin":
		return copyDarwin(text)
	case "linux":
		return copyLinux(text)
	case "windows":
		return copyWindows(text)
	default:
		return fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
}

// copyDarwin copies text to clipboard on macOS
func copyDarwin(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// copyLinux copies text to clipboard on Linux
func copyLinux(text string) error {
	// Try xclip first
	if isCommandAvailable("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try xsel
	if isCommandAvailable("xsel") {
		cmd := exec.Command("xsel", "--clipboard", "--input")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	// Try wl-clipboard (Wayland)
	if isCommandAvailable("wl-copy") {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(text)
		if err := cmd.Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("clipboard not available: no clipboard utility found\n" +
		"Please install one of the following:\n" +
		"  • X11: 'sudo apt install xclip' or 'sudo apt install xsel'\n" +
		"  • Wayland: 'sudo apt install wl-clipboard'\n" +
		"Alternatively, use the export feature to save colors to a file")
}

// copyWindows copies text to clipboard on Windows
func copyWindows(text string) error {
	cmd := exec.Command("cmd", "/c", "clip")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
