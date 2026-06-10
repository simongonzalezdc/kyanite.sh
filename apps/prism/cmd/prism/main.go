package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/prism/internal/app"
)

var (
	Version = "1.0.0"
	Commit  = "dev"
)

func main() {
	// Command-line flags
	version := flag.Bool("version", false, "Display version information")
	versionShort := flag.Bool("v", false, "Display version information (shorthand)")
	help := flag.Bool("help", false, "Display help information")
	helpShort := flag.Bool("h", false, "Display help information (shorthand)")
	flag.Parse()

	// Handle version flag
	if *version || *versionShort {
		fmt.Printf("Prism.sh v%s (commit: %s)\n", Version, Commit)
		fmt.Println("A terminal-based color palette design tool")
		fmt.Println("Part of the Kyanite Suite")
		os.Exit(0)
	}

	// Handle help flag
	if *help || *helpShort {
		fmt.Println("Prism.sh - Terminal Color Palette Design Tool")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  prism           Start the interactive TUI")
		fmt.Println("  prism -version  Display version information")
		fmt.Println("  prism -help     Display this help message")
		fmt.Println()
		fmt.Println("Controls:")
		fmt.Println("  Ctrl+Q          Quit application")
		fmt.Println("  Ctrl+H          Toggle help overlay")
		fmt.Println("  Ctrl+Shift+T    Cycle through themes")
		fmt.Println("  Esc             Navigate back")
		fmt.Println("  Arrow keys      Navigate menus")
		fmt.Println("  Enter           Select option")
		fmt.Println()
		fmt.Println("For more information, visit: https://github.com/kyanite/prism")
		os.Exit(0)
	}

	// Create the root model
	m := app.NewModel()

	// Create the Bubble Tea program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running application: %v\n", err)
		os.Exit(1)
	}
}
