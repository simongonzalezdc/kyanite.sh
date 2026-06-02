package ai

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// OllamaManager handles Ollama lifecycle management
type OllamaManager struct {
	baseURL      string
	model        string
	pullAttempts int
}

// NewOllamaManager creates a new Ollama manager
func NewOllamaManager(baseURL, model string) *OllamaManager {
	return &OllamaManager{
		baseURL:      baseURL,
		model:        model,
		pullAttempts: 0,
	}
}

// IsAvailable checks if Ollama is running and accessible
func (om *OllamaManager) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", om.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Launch starts Ollama asynchronously
func (om *OllamaManager) Launch(ctx context.Context) error {
	// Launch Ollama in background
	if err := om.launchAsync(); err != nil {
		return fmt.Errorf("failed to launch Ollama: %w", err)
	}

	// Wait for it to become available
	if err := om.waitForAvailable(ctx); err != nil {
		return fmt.Errorf("ollama did not become available: %w", err)
	}

	return nil
}

// launchAsync starts Ollama in the background
func (om *OllamaManager) launchAsync() error {
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return err
	}
	// Detach - let it run in background
	return nil
}

// waitForAvailable waits for Ollama to become available
func (om *OllamaManager) waitForAvailable(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(30 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timeout:
			return fmt.Errorf("timeout waiting for Ollama")
		case <-ticker.C:
			if om.IsAvailable() {
				return nil
			}
		}
	}
}

// IsModelAvailable checks if a specific model is available
func (om *OllamaManager) IsModelAvailable(modelName string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to query the model
	cmd := exec.CommandContext(ctx, "ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Simple check if model name appears in output
	outputStr := string(output)
	return len(outputStr) > 0 && contains(outputStr, modelName)
}

// PullModel downloads a model
func (om *OllamaManager) PullModel(modelName string) error {
	if om.pullAttempts > 0 {
		return fmt.Errorf("model pull already attempted")
	}
	om.pullAttempts++

	fmt.Printf("🤖 Pulling model %s automatically...\n", modelName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ollama", "pull", modelName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}

	return nil
}

// contains checks if a string contains a substring (helper function)
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
