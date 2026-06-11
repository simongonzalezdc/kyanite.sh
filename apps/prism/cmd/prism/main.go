package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	ai "github.com/kyanite/ai"
	"github.com/kyanite/config"
	"github.com/kyanite/prism/internal/app"
)

var (
	Version = "1.0.0"
	Commit  = "dev"
)

func main() {
	version := flag.Bool("version", false, "Print version and exit")
	help := flag.Bool("help", false, "Print help and exit")
	flag.Parse()

	if *version {
		fmt.Printf("prism %s (%s)\n", Version, Commit)
		os.Exit(0)
	}

	if *help {
		fmt.Println("prism — color palette & WCAG contrast tool for the terminal")
		fmt.Println()
		fmt.Println("Usage: prism [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --version       Print version and exit")
		fmt.Println("  --help          Print help and exit")
		fmt.Println()
		fmt.Println("Keybindings:")
		fmt.Println("  Ctrl+A / p      Toggle AI panel")
		fmt.Println("  q               Quit")
		fmt.Println("  Arrow keys      Navigate menus")
		fmt.Println("  Enter           Select option")
		fmt.Println()
		fmt.Println("For more information, visit: https://github.com/kyanite/prism")
		os.Exit(0)
	}

	// Load config from ~/.config/kyanite/config.yaml + env vars
	root, _ := config.Load()
	cfg := ai.ConfigFromRoot(root, "prism")
	brain, _ := ai.New(cfg)

	// Create the root model
	m := app.NewModel(brain)

	// Attempt to load the most recent session (best-effort)
	m.LoadRecentSession()

	// Ensure session is saved and resources are released on exit
	defer func() {
		m.SaveCurrentSession()
		m.Close()
	}()

	// Create the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "prism: %v\n", err)
		os.Exit(1)
	}
}
