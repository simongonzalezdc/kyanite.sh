package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Dispatch to a subcommand; default to focus.
	app := "focus"
	if len(os.Args) > 1 {
		app = os.Args[1]
	}
	
	switch app {
	case "focus":
		focusMain()
	default:
		log.Fatalf("unknown app: %s", app)
	}
}

func focusMain() {
	// Run focus via go run to avoid import conflicts
	cmd := exec.Command("go", "run", "./cmd/focus")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run focus: %v", err)
	}
}

