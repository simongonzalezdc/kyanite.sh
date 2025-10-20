# MCP Server Configuration for CRUSH CLI

This directory contains the MCP (Model Context Protocol) server implementation for integrating golangci-lint with AI assistants.

## Quick Start

### Start the MCP Server
```bash
# Using the built-in command
neon mcp-server

# Or run directly
go run ./cmd/mcp-server
```

### Configure MCP Client

Example configuration for MCP clients (like Claude Desktop, Cursor, etc.):

```json
{
  "mcpServers": {
    "crush-golangci-lint": {
      "command": "neon",
      "args": ["mcp-server"],
      "cwd": "C:\\Users\\Simon\\dev\\crush-cli"
    }
  }
}
```

Or using the full path:
```json
{
  "mcpServers": {
    "crush-golangci-lint": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server"],
      "cwd": "C:\\Users\\Simon\\dev\\crush-cli"
    }
  }
}
```

## Available Tools

### `golangci_lint`
Run golangci-lint on the current directory or specific files.

**Parameters:**
- `directory` (string, optional): Directory to run linting on (default: current directory)
- `args` (array, optional): Additional arguments for golangci-lint

**Example usage:**
```json
{
  "name": "golangci_lint",
  "arguments": {
    "directory": ".",
    "args": ["--timeout", "5m"]
  }
}
```

## Troubleshooting

### "Error calling 'initialize'" Issues

1. **Check golangci-lint installation:**
   ```bash
   golangci-lint version
   ```

2. **Verify working directory:**
   - Ensure the MCP server is started from the project root
   - Check that `go.mod` exists in the directory

3. **Test MCP server directly:**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize"}' | neon mcp-server
   ```

4. **Check for path issues:**
   - Use absolute paths in MCP client configuration
   - Ensure `neon` command is in your PATH

### Common MCP Client Configurations

**Cursor:**
- Edit `C:\Users\Simon\AppData\Roaming\Code\User\mcp.json`
- Add the server configuration under `servers`

**Claude Desktop:**
- Edit `~/.config/claude/claude_desktop_config.json`
- Add to `mcpServers` section

## Development

The MCP server is implemented in `cmd/mcp-server/main.go` and follows the MCP 2024-11-05 protocol specification. It exposes golangci-lint functionality through a standardized JSON-RPC interface that AI assistants can use to perform code quality checks.

Key features:
- Full MCP protocol compliance
- Error handling and graceful degradation
- Configurable working directories
- Support for additional golangci-lint arguments
- Real-time linting results streaming