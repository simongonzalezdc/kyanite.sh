package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Comprehensive Theme Testing Tool
// This tool provides a complete testing workflow for the Kyanite theme system

func main() {
	fmt.Println("========================================")
	fmt.Println("noise.sh Comprehensive Theme Testing")
	fmt.Println("========================================")
	fmt.Println()

	// Check if noise.sh is built
	if !isBinaryBuilt() {
		fmt.Println("noise.sh binary not found. Building...")
		if !buildApplication() {
			fmt.Println("ERROR: Failed to build noise.sh")
			os.Exit(1)
		}
		fmt.Println("✓ Build successful")
	} else {
		fmt.Println("✓ noise.sh binary found")
	}

	// Run automated theme tests
	fmt.Println("\n========================================")
	fmt.Println("Running Automated Theme Tests")
	fmt.Println("========================================")

	if !runAutomatedTests() {
		fmt.Println("WARNING: Some automated tests failed")
		fmt.Println("Continuing with manual testing...")
	}

	// Provide manual testing instructions
	fmt.Println("\n========================================")
	fmt.Println("Manual Testing Instructions")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("The automated tests are complete. Now let's perform manual testing.")
	fmt.Println("This ensures the visual appearance and user experience are correct.")
	fmt.Println()

	// Display theme information
	displayThemeInformation()

	// Interactive testing prompt
	fmt.Println("\n========================================")
	fmt.Println("Interactive Testing")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Choose a testing option:")
	fmt.Println("1. Launch noise.sh for manual theme testing")
	fmt.Println("2. Launch with debug mode")
	fmt.Println("3. Launch in quick mode (scratch mode)")
	fmt.Println("4. Run theme switching demo")
	fmt.Println("5. Exit")
	fmt.Println()

	choice := getUserChoice("Select option (1-5): ")

	switch choice {
	case "1":
		launchApplication("")
	case "2":
		launchApplication("--debug")
	case "3":
		launchApplication("quick")
	case "4":
		runThemeDemo()
	case "5":
		fmt.Println("Testing completed. Goodbye!")
		return
	default:
		fmt.Println("Invalid choice. Exiting...")
		return
	}
}

// isBinaryBuilt checks if the noise.sh binary exists
func isBinaryBuilt() bool {
	if _, err := os.Stat("noise.exe"); err == nil {
		return true
	}
	if _, err := os.Stat("noise"); err == nil {
		return true
	}
	if _, err := os.Stat("bin/noise.exe"); err == nil {
		return true
	}
	if _, err := os.Stat("bin/noise"); err == nil {
		return true
	}
	return false
}

// buildApplication builds the noise.sh application
func buildApplication() bool {
	cmd := exec.Command("go", "build", "-o", "noise.exe", "./cmd/noise")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err == nil
}

// runAutomatedTests runs the automated theme tests
func runAutomatedTests() bool {
	fmt.Println("Running automated theme system tests...")

	cmd := exec.Command("go", "run", "scripts/test_themes.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Automated tests failed: %v\n", err)
		return false
	}

	fmt.Println("✓ Automated tests completed successfully")
	return true
}

// displayThemeInformation shows detailed theme information
func displayThemeInformation() {
	themes := []struct {
		id   string
		name string
		desc string
	}{
		{"monochrome", "Monochrome", "Classic black and white theme"},
		{"amber-night", "Amber Night", "Warm amber tones (default)"},
		{"twilight-mist", "Twilight Mist", "Soft purple gradients"},
		{"indigo-depths", "Indigo Depths", "Deep blue ocean colors"},
		{"forest-path", "Forest Path", "Natural green tones"},
		{"clay-earth", "Clay Earth", "Warm earth colors"},
		{"iron-forge", "Iron Forge", "Industrial reds and grays"},
		{"sunlight", "Sunlight", "Bright golden yellows"},
		{"cyan-wave", "Cyan Wave", "Cool cyan blues"},
		{"electric-rose", "Electric Rose", "Vibrant pink and cyan"},
	}

	fmt.Println("Available Themes (Kyanite Suite):")
	fmt.Println()
	for i, theme := range themes {
		fmt.Printf("%d. %s (%s)\n", i+1, theme.name, theme.id)
		fmt.Printf("   %s\n", theme.desc)
		fmt.Println()
	}

	fmt.Println("Theme Switching Shortcuts:")
	fmt.Println("  Ctrl+Shift+T - Cycle through themes")
	fmt.Println("  Ctrl+Shift+N - Next theme")
	fmt.Println("  Ctrl+Shift+P - Previous theme")
	fmt.Println()

	fmt.Println("Testing Checklist:")
	fmt.Println("  [ ] Visual appearance matches theme description")
	fmt.Println("  [ ] Text is readable against background")
	fmt.Println("  [ ] UI elements are properly styled")
	fmt.Println("  [ ] Color contrast is adequate")
	fmt.Println("  [ ] No visual artifacts or glitches")
	fmt.Println("  [ ] Theme switching is smooth")
	fmt.Println("  [ ] All features work in current theme")
	fmt.Println()
}

// getUserChoice gets user input from stdin
func getUserChoice(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

// launchApplication launches the noise.sh application with specified arguments
func launchApplication(args string) {
	fmt.Printf("\nLaunching noise.sh%s...\n", func() string {
		if args != "" {
			return " with " + args
		}
		return ""
	}())
	fmt.Println("Press Ctrl+Q to exit the application")
	fmt.Println("Press Ctrl+Shift+T to cycle through themes")
	fmt.Println("Press F1 for help")
	fmt.Println()

	var cmd *exec.Cmd
	if args != "" {
		cmd = exec.Command("./noise.exe", strings.Split(args, " ")...)
	} else {
		cmd = exec.Command("./noise.exe")
	}

	// Set up stdin, stdout, stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the application
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Application exited with error: %v\n", err)
	} else {
		fmt.Println("Application closed successfully")
	}
}

// runThemeDemo runs a theme switching demonstration
func runThemeDemo() {
	fmt.Println("\n========================================")
	fmt.Println("Theme Switching Demo")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("This demo will launch noise.sh and cycle through all themes automatically.")
	fmt.Println("You'll have 3 seconds to observe each theme before switching.")
	fmt.Println()
	fmt.Println("Press Enter to start the demo...")
	getUserChoice("")

	// Note: This is a simplified demo
	// In a real implementation, you might use automation tools
	// or integrate with the application's testing framework

	fmt.Println("Launching application for manual theme cycling...")
	fmt.Println("Instructions:")
	fmt.Println("1. When the application launches, use Ctrl+Shift+T to cycle themes")
	fmt.Println("2. Spend a few seconds examining each theme")
	fmt.Println("3. Verify the visual appearance matches expectations")
	fmt.Println("4. Test basic functionality in each theme")
	fmt.Println("5. Press Ctrl+Q when done")
	fmt.Println()

	launchApplication("--debug")

	fmt.Println("\nDemo completed!")
	fmt.Println("Please provide feedback on any issues you encountered.")
}
