package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/lyricforge/internal/config"
	"github.com/puente-labs/lyricforge/internal/infra/db"
	"github.com/puente-labs/lyricforge/internal/logging"
	"github.com/puente-labs/lyricforge/internal/plugins"
	"github.com/puente-labs/lyricforge/internal/ui"
)

var (
	// Version information (set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// main is the entry point for the LyricForge application
func main() {
	// Parse command line arguments
	if err := parseArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logging.SetDefaultLogger(logger)

	// Log startup information
	logger.Infof("Starting LyricForge v%s (commit: %s, built: %s)", version, commit, date)
	logger.Debug("Configuration loaded successfully")

	// Initialize database
	var database *db.DB
	if !cfg.Dev.SkipDatabase {
		database, err = db.New(db.Config{DataDir: cfg.GetDataDir()})
		if err != nil {
			logger.Errorf("Failed to initialize database: %v", err)
			if !cfg.Dev.Debug {
				fmt.Fprintf(os.Stderr, "Failed to initialize database. Use --debug for more details.\n")
				os.Exit(1)
			}
		} else {
			logger.Info("Database initialized successfully")
		}
	} else {
		logger.Debug("Skipping database initialization (dev mode)")
	}

	// Initialize plugin system
	pluginManager := plugins.NewManager(cfg, logger)
	if err := pluginManager.LoadPlugins(); err != nil {
		logger.Warnf("Failed to load plugins: %v", err)
	} else {
		logger.Infof("Plugin system initialized with %d plugins", len(pluginManager.GetPlugins()))
	}

	// Create root model with plugin manager
	rootModel := ui.NewRootModel(pluginManager)

	// Set up the Bubble Tea program
	p := tea.NewProgram(
		rootModel,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	// Handle shutdown gracefully
	defer func() {
		if database != nil {
			if err := database.Close(); err != nil {
				logger.Errorf("Error closing database: %v", err)
			}
		}
		logger.Info("LyricForge shutdown complete")
	}()

	// Run the program
	logger.Info("Starting TUI...")
	if _, err := p.Run(); err != nil {
		logger.Errorf("TUI error: %v", err)
		os.Exit(1)
	}
}

// parseArgs parses command line arguments
func parseArgs() error {
	// Simple argument parsing for foundation phase
	// In a full implementation, this would use cobra CLI framework

	for i, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("LyricForge %s (commit: %s, built: %s)\n", version, commit, date)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "--debug":
			// Enable debug mode via environment variable
			os.Setenv("LYRICFORGE_DEV_DEBUG", "true")
		default:
			if i == 0 && !isFlag(arg) {
				// First non-flag argument might be a file to open
				logging.Infof("Opening file: %s", arg)
			}
		}
	}

	return nil
}

// isFlag checks if a string is a command line flag
func isFlag(s string) bool {
	return len(s) > 1 && s[0] == '-'
}

// printHelp prints the help message
func printHelp() {
	fmt.Println("LyricForge - AI-Powered Songwriting Terminal Interface")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  lyricforge [OPTIONS] [FILE]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  --debug          Enable debug mode")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  FILE             Song file to open (optional)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  lyricforge                    Start LyricForge")
	fmt.Println("  lyricforge --debug            Start with debug logging")
	fmt.Println("  lyricforge song.md            Open a specific song file")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/puente-labs/lyricforge")
}
