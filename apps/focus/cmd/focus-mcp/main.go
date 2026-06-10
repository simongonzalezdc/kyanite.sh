package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type mcpServer struct{}

type request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
}

type response struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type initializeResult struct {
	ProtocolVersion string     `json:"protocolVersion"`
	Capabilities    any        `json:"capabilities"`
	ServerInfo      serverInfo `json:"serverInfo"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type tool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"inputSchema"`
}

type toolsResult struct {
	Tools []tool `json:"tools"`
}

type callToolParams struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type callToolResult struct {
	Content []content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func newMCPServer() *mcpServer {
	return &mcpServer{}
}

func (s *mcpServer) handleRequest(req request) response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	default:
		return response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *mcpServer) handleInitialize(req request) response {
	result := initializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: map[string]any{
			"tools": map[string]any{},
		},
		ServerInfo: serverInfo{
			Name:    "crush-golangci-lint",
			Version: "1.0.0",
		},
	}

	return response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *mcpServer) handleToolsList(req request) response {
	tools := []tool{
		{
			Name:        "golangci_lint",
			Description: "Run golangci-lint on the current directory or specific files",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"directory": map[string]any{
						"type":        "string",
						"description": "Directory to run linting on (default: current directory)",
					},
					"args": map[string]any{
						"type":        "array",
						"description": "Additional arguments for golangci-lint",
						"items": map[string]any{
							"type": "string",
						},
					},
				},
			},
		},
	}

	result := toolsResult{Tools: tools}

	return response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *mcpServer) handleToolsCall(req request) response {
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		return response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &rpcError{
				Code:    -32603,
				Message: fmt.Sprintf("Internal error marshaling parameters: %v", err),
			},
		}
	}

	var params callToolParams
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: callToolResult{
				Content: []content{{Type: "text", Text: fmt.Sprintf("Error parsing tool parameters: %v", err)}},
				IsError: true,
			},
		}
	}

	switch params.Name {
	case "golangci_lint":
		return s.runGolangciLint(req, params.Arguments)
	default:
		return response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: callToolResult{
				Content: []content{{Type: "text", Text: fmt.Sprintf("Unknown tool: %s", params.Name)}},
				IsError: true,
			},
		}
	}
}

func (s *mcpServer) runGolangciLint(req request, args map[string]any) response {
	workingDir := "."
	if dir, ok := args["directory"].(string); ok && dir != "" {
		workingDir = dir
	}

	lintArgs := []string{"run"}
	if additionalArgs, ok := args["args"].([]any); ok {
		var argsList []string
		for _, arg := range additionalArgs {
			if argStr, ok := arg.(string); ok {
				argsList = append(argsList, argStr)
			}
		}
		if len(argsList) == 1 && (argsList[0] == "version" || argsList[0] == "--version" || argsList[0] == "-v") {
			lintArgs = []string{"version"}
		} else {
			lintArgs = append(lintArgs, argsList...)
		}
	}

	cmd := exec.Command("golangci-lint", lintArgs...)
	cmd.Dir = workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: callToolResult{
					Content: []content{{Type: "text", Text: fmt.Sprintf("Failed to execute golangci-lint: %v\n\nOutput:\n%s", err, string(output))}},
					IsError: true,
				},
			}
		}
	}

	outputText := string(output)
	if outputText == "" {
		outputText = "No linting issues found. ✅"
	}

	result := callToolResult{Content: []content{{Type: "text", Text: outputText}}}

	return response{JSONRPC: "2.0", ID: req.ID, Result: result}
}

func main() {
	srv := newMCPServer()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var req request
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		resp := srv.handleRequest(req)
		if err := encoder.Encode(resp); err != nil {
			break
		}
	}
}
