package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/puente-labs/noise/internal/config"
	"github.com/puente-labs/noise/internal/infra/db"
	"github.com/puente-labs/noise/internal/logging"
	"github.com/puente-labs/noise/internal/plugins"
	"github.com/puente-labs/noise/internal/ui"
)


var (
	// Version information (set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// main is the entry point for the noise.sh application
func main() {
	// Parse command line arguments
	quickConfig, err := parseArgs()
	if err != nil {
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
	logger.Infof("Starting noise.sh v%s (commit: %s, built: %s)", version, commit, date)
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

	// If quick start is requested, configure it
	if quickConfig != nil {
		rootModel.SetQuickStart(quickConfig)
	}

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
		logger.Info("noise.sh shutdown complete")
	}()

	// Run the program
	logger.Info("Starting TUI...")
	if _, err := p.Run(); err != nil {
		logger.Errorf("TUI error: %v", err)
		os.Exit(1)
	}
}

// parseArgs parses command line arguments
func parseArgs() (*ui.QuickStartConfig, error) {
	// Simple argument parsing for foundation phase
	// In a full implementation, this would use cobra CLI framework

	var quickConfig *ui.QuickStartConfig

	for i, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("noise.sh %s (commit: %s, built: %s)\n", version, commit, date)
			os.Exit(0)
		case "--help", "-h":
			printHelp()
			os.Exit(0)
		case "--debug":
			// Enable debug mode via environment variable
			os.Setenv("NOISE_DEV_DEBUG", "true")
		case "quick":
			// Quick start command - check if theme is provided
			quickConfig = &ui.QuickStartConfig{
				ScratchMode:   true,
				AutoBrainstorm: false,
			}
			
			// Check if next argument is a theme (not a flag)
			if i+1 < len(os.Args[1:]) && !isFlag(os.Args[1:][i+1]) {
				quickConfig.Theme = os.Args[1:][i+1]
				quickConfig.AutoBrainstorm = true
			}
			return quickConfig, nil
		default:
			if i == 0 && !isFlag(arg) {
				// First non-flag argument might be a file to open
				logging.Infof("Opening file: %s", arg)
			}
		}
	}

	return quickConfig, nil
}

// isFlag checks if a string is a command line flag
func isFlag(s string) bool {
	return len(s) > 1 && s[0] == '-'
}

// printHelp prints the help message
func printHelp() {
	fmt.Println("noise.sh - AI-Powered Songwriting Terminal Interface")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  noise [OPTIONS] [COMMAND] [ARGS]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  quick [theme]            Start rapid prototyping mode")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  -v, --version    Show version information")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("  --debug          Enable debug mode")
	fmt.Println()
	fmt.Println("ARGUMENTS:")
	fmt.Println("  FILE             Song file to open (optional)")
	fmt.Println("  theme            Theme for quick start (optional)")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  noise                    Start noise.sh")
	fmt.Println("  noise --debug            Start with debug logging")
	fmt.Println("  noise song.md            Open a specific song file")
	fmt.Println("  noise quick              Start in scratch mode")
	fmt.Println("  noise quick \"lost love\"  Start with theme brainstorm")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/puente-labs/noise")
}
