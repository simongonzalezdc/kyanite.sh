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

// Paste pastes text from the system clipboard
func Paste() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return pasteDarwin()
	case "linux":
		return pasteLinux()
	case "windows":
		return pasteWindows()
	default:
		return "", fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
}

// copyDarwin copies text to clipboard on macOS
func copyDarwin(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// pasteDarwin pastes text from clipboard on macOS
func pasteDarwin() (string, error) {
	cmd := exec.Command("pbpaste")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
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

// pasteLinux pastes text from clipboard on Linux
func pasteLinux() (string, error) {
	// Try xclip first
	if isCommandAvailable("xclip") {
		cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
		out, err := cmd.Output()
		if err == nil {
			return string(out), nil
		}
	}

	// Try xsel
	if isCommandAvailable("xsel") {
		cmd := exec.Command("xsel", "--clipboard", "--output")
		out, err := cmd.Output()
		if err == nil {
			return string(out), nil
		}
	}

	// Try wl-clipboard (Wayland)
	if isCommandAvailable("wl-paste") {
		cmd := exec.Command("wl-paste")
		out, err := cmd.Output()
		if err == nil {
			return string(out), nil
		}
	}

	return "", fmt.Errorf("clipboard not available: no clipboard utility found\n" +
		"Please install one of the following:\n" +
		"  • X11: 'sudo apt install xclip' or 'sudo apt install xsel'\n" +
		"  • Wayland: 'sudo apt install wl-clipboard'")
}

// copyWindows copies text to clipboard on Windows
func copyWindows(text string) error {
	cmd := exec.Command("cmd", "/c", "clip")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// pasteWindows pastes text from clipboard on Windows
func pasteWindows() (string, error) {
	// Windows doesn't have a built-in paste command
	// We use PowerShell Get-Clipboard
	cmd := exec.Command("powershell", "-command", "Get-Clipboard")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// isCommandAvailable checks if a command is available in PATH
func isCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}
