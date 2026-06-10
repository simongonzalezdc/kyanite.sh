package cli

import (
	"os"
	"os/exec"

	pterm "github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var mcpServerCmd = &cobra.Command{
	Use:   "mcp-server",
	Short: "Start the MCP server for golangci-lint integration",
	Long: `Start the Model Context Protocol (MCP) server that provides golangci-lint 
functionality to MCP clients. This allows AI assistants to run linting checks
through the standardized MCP protocol.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Execute the mcp-server binary
		serverCmd := exec.Command("go", "run", "./cmd/mcp-server")
		serverCmd.Stdin = os.Stdin
		serverCmd.Stdout = os.Stdout
		serverCmd.Stderr = os.Stderr

		if err := serverCmd.Run(); err != nil {
			pterm.Error.Printf("Failed to start MCP server: %v\n", err)
			os.Exit(1)
		}
	},
}
