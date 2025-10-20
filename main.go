package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	// Dispatch to a subcommand; default to neon.
	app := "neon"
	if len(os.Args) > 1 {
		app = os.Args[1]
	}
	
	switch app {
	case "neon":
		neonMain()
	case "todo":
		todoMain()
	default:
		log.Fatalf("unknown app: %s", app)
	}
}

func neonMain() {
	// Run neon via go run to avoid import conflicts
	cmd := exec.Command("go", "run", "./cmd/neon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run neon: %v", err)
	}
}

func todoMain() {
	// Run todo via go run to avoid import conflicts
	cmd := exec.Command("go", "run", "./cmd/todo")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run todo: %v", err)
	}
}
