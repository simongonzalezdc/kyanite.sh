package main

import (
	"fmt"
	"os"

	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/plugins"
	"github.com/Kyanite/noise/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	// Version information (set during build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

const helpMessage = `noise.sh - AI-Powered Songwriting Terminal Interface

USAGE:
  noise [OPTIONS] [COMMAND] [ARGS]

COMMANDS:
  quick [theme]            Start rapid prototyping mode

OPTIONS:
  -v, --version    Show version information
  -h, --help       Show this help message
  --debug          Enable debug mode

ARGUMENTS:
  FILE             Song file to open (optional)
  theme            Theme for quick start (optional)

EXAMPLES:
  noise                    Start noise.sh
  noise --debug            Start with debug logging
  noise song.md            Open a specific song file
  noise quick              Start in scratch mode
  noise quick "lost love"  Start with theme brainstorm

For more information, visit: https://github.com/Kyanite/noise
`

type cliParseResult struct {
	QuickConfig *ui.QuickStartConfig
	PendingFile string
	ShowHelp    bool
	ShowVersion bool
}

// main is the entry point for the noise.sh application
func main() {
	parsedArgs, err := parseArgs(os.Args[1:])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if parsedArgs.ShowVersion {
		if _, err := fmt.Fprintf(os.Stdout, "noise.sh %s (commit: %s, built: %s)\n", version, commit, date); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error writing version information: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if parsedArgs.ShowHelp {
		if _, err := fmt.Fprint(os.Stdout, helpMessage); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error writing help information: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := run(parsedArgs); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(parsed *cliParseResult) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}

	logging.SetDefaultLogger(logger)

	logger.Infof("Starting noise.sh v%s (commit: %s, built: %s)", version, commit, date)
	logger.Debug("Configuration loaded successfully")

	var database *db.DB
	if !cfg.Dev.SkipDatabase {
		database, err = db.New(db.Config{DataDir: cfg.GetDataDir()})
		if err != nil {
			logger.Errorf("Failed to initialize database: %v", err)
			if !cfg.Dev.Debug {
				return fmt.Errorf("initialize database: %w", err)
			}
		} else {
			logger.Info("Database initialized successfully")
		}
	} else {
		logger.Debug("Skipping database initialization (dev mode)")
	}

	pluginManager := plugins.NewManager(cfg, logger)
	if err := pluginManager.LoadPlugins(); err != nil {
		logger.Warnf("Failed to load plugins: %v", err)
	} else {
		logger.Infof("Plugin system initialized with %d plugins", len(pluginManager.GetPlugins()))
	}

	rootModel := ui.NewRootModel(pluginManager)

	if parsed.QuickConfig != nil {
		rootModel.SetQuickStart(parsed.QuickConfig)
	}

	if parsed.PendingFile != "" {
		logger.Infof("Opening file: %s", parsed.PendingFile)
	}

	p := tea.NewProgram(
		rootModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	defer func() {
		if database != nil {
			if err := database.Close(); err != nil {
				logger.Errorf("Error closing database: %v", err)
			}
		}
		logger.Info("noise.sh shutdown complete")
	}()

	logger.Info("Starting TUI...")
	if _, err := p.Run(); err != nil {
		logger.Errorf("TUI error: %v", err)
		return fmt.Errorf("run TUI: %w", err)
	}

	return nil
}

func parseArgs(args []string) (*cliParseResult, error) {
	result := &cliParseResult{}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch arg {
		case "--version", "-v":
			result.ShowVersion = true
			return result, nil
		case "--help", "-h":
			result.ShowHelp = true
			return result, nil
		case "--debug":
			if err := os.Setenv("NOISE_DEV_DEBUG", "true"); err != nil {
				return nil, fmt.Errorf("set debug environment variable: %w", err)
			}
		case "quick":
			result.QuickConfig = &ui.QuickStartConfig{
				ScratchMode:    true,
				AutoBrainstorm: false,
			}

			if i+1 < len(args) && !isFlag(args[i+1]) {
				result.QuickConfig.Theme = args[i+1]
				result.QuickConfig.AutoBrainstorm = true
				i++
			}
		default:
			if !isFlag(arg) && result.PendingFile == "" {
				result.PendingFile = arg
			}
		}
	}

	return result, nil
}

// isFlag checks if a string is a command line flag
func isFlag(s string) bool {
	return len(s) > 1 && s[0] == '-'
}
