package ai

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewQuickIdeaAgent(t *testing.T) {
	agent := NewQuickIdeaAgent()

	if agent == nil {
		t.Fatal("NewQuickIdeaAgent returned nil")
	}

	if agent.client == nil {
		t.Error("Expected client to be initialized")
	}

	if agent.model != defaultQuickIdeaModel {
		t.Errorf("Expected model %s, got %s", defaultQuickIdeaModel, agent.model)
	}

	if agent.timeout != defaultQuickIdeaTimeout {
		t.Errorf("Expected timeout %v, got %v", defaultQuickIdeaTimeout, agent.timeout)
	}
}

func TestQuickIdeaAgent_WithClient(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockQuickClient{}
	timeout := 2 * time.Second

	agentWithClient := agent.WithClient(mockClient, timeout)

	if agentWithClient == nil {
		t.Fatal("WithClient returned nil")
	}

	if agentWithClient.client != mockClient {
		t.Error("Expected client to be set to mockClient")
	}

	if agentWithClient.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, agentWithClient.timeout)
	}
}

func TestQuickIdeaAgent_WithClientNil(t *testing.T) {
	agent := NewQuickIdeaAgent()
	agentWithClient := agent.WithClient(nil, 0)

	if agentWithClient != agent {
		t.Error("Expected same agent to be returned when client is nil")
	}
}

func TestQuickIdeaAgent_GenerateUnstickMode(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockQuickClient{response: "1. First suggestion\n2. Second suggestion\n3. Third suggestion"}
	agentWithClient := agent.WithClient(mockClient, defaultQuickIdeaTimeout)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Previous line here",
		Options: map[string]string{"section": "verse"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(resp.Suggestions))
	}

	if resp.Suggestions[0] != "First suggestion" {
		t.Errorf("Expected 'First suggestion', got '%s'", resp.Suggestions[0])
	}
}

func TestQuickIdeaAgent_GenerateSparkMode(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockQuickClient{response: "1. First idea\n2. Second idea\n3. Third idea"}
	agentWithClient := agent.WithClient(mockClient, defaultQuickIdeaTimeout)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeSpark,
		Context: "love",
		Options: map[string]string{"theme": "love"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(resp.Suggestions))
	}
}

func TestQuickIdeaAgent_GenerateTweakMode(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockQuickClient{response: "1. Variation 1\n2. Variation 2\n3. Variation 3"}
	agentWithClient := agent.WithClient(mockClient, defaultQuickIdeaTimeout)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeTweak,
		Context: "Original line here",
		Options: map[string]string{"section": "verse"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 suggestions, got %d", len(resp.Suggestions))
	}
}

func TestQuickIdeaAgent_GenerateCheckMode(t *testing.T) {
	agent := NewQuickIdeaAgent()
	mockClient := &mockQuickClient{response: "STRONG\nSharpen the central metaphor"}
	agentWithClient := agent.WithClient(mockClient, defaultQuickIdeaTimeout)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeCheck,
		Context: "Line to check here",
		Options: map[string]string{"mode": "sketch"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	if resp.Rating != "STRONG" {
		t.Errorf("Expected rating 'STRONG', got '%s'", resp.Rating)
	}

	if resp.Tip != "Sharpen the central metaphor" {
		t.Errorf("Expected tip 'Sharpen the central metaphor', got '%s'", resp.Tip)
	}
}

func TestQuickIdeaAgent_GenerateWithFallback(t *testing.T) {
	agent := NewQuickIdeaAgent()
	// Use a client that returns an error
	mockClient := &mockQuickClient{err: true}
	agentWithClient := agent.WithClient(mockClient, defaultQuickIdeaTimeout)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Previous line here",
		Options: map[string]string{"section": "verse"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	// Should return fallback suggestions
	if len(resp.Suggestions) == 0 {
		t.Error("Expected fallback suggestions, got none")
	}
}

func TestQuickIdeaAgent_GenerateTimeout(t *testing.T) {
	agent := NewQuickIdeaAgent()
	// Use a very short timeout
	agentWithClient := agent.WithClient(&mockQuickClient{}, 1*time.Nanosecond)

	ctx := context.Background()
	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Previous line here",
		Options: map[string]string{"section": "verse"},
	}

	resp, err := agentWithClient.Generate(ctx, req)

	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Generate returned nil response")
	}

	// Should return fallback suggestions due to timeout
	if len(resp.Suggestions) == 0 {
		t.Error("Expected fallback suggestions, got none")
	}
}

func TestQuickIdeaAgent_GenerateInvalidMode(t *testing.T) {
	agent := NewQuickIdeaAgent()

	ctx := context.Background()
	req := QuickRequest{
		Mode:    "invalid",
		Context: "Some context",
		Options: map[string]string{},
	}

	_, err := agent.Generate(ctx, req)

	if err == nil {
		t.Error("Expected error for invalid mode")
	}
}

func TestQuickIdeaAgent_GenerateFallbackUnstick(t *testing.T) {
	agent := NewQuickIdeaAgent()

	req := QuickRequest{
		Mode:    QuickIdeaModeUnstick,
		Context: "Line one\nLine two\nLine three",
		Options: map[string]string{},
	}

	resp := agent.generateFallback(req)

	if resp == nil {
		t.Fatal("generateFallback returned nil")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 fallback suggestions, got %d", len(resp.Suggestions))
	}

	// Check that suggestions are based on the last line
	lastLine := "Line three"
	for _, suggestion := range resp.Suggestions {
		if !contains(suggestion, lastLine) {
			t.Errorf("Expected suggestion to contain '%s', got '%s'", lastLine, suggestion)
		}
	}
}

func TestQuickIdeaAgent_GenerateFallbackSpark(t *testing.T) {
	agent := NewQuickIdeaAgent()

	req := QuickRequest{
		Mode:    QuickIdeaModeSpark,
		Context: "love",
		Options: map[string]string{"theme": "love"},
	}

	resp := agent.generateFallback(req)

	if resp == nil {
		t.Fatal("generateFallback returned nil")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 fallback suggestions, got %d", len(resp.Suggestions))
	}

	// Check that suggestions contain the theme
	for _, suggestion := range resp.Suggestions {
		if !contains(suggestion, "love") {
			t.Errorf("Expected suggestion to contain 'love', got '%s'", suggestion)
		}
	}
}

func TestQuickIdeaAgent_GenerateFallbackTweak(t *testing.T) {
	agent := NewQuickIdeaAgent()

	req := QuickRequest{
		Mode:    QuickIdeaModeTweak,
		Context: "Original line here",
		Options: map[string]string{},
	}

	resp := agent.generateFallback(req)

	if resp == nil {
		t.Fatal("generateFallback returned nil")
	}

	if len(resp.Suggestions) != 3 {
		t.Errorf("Expected 3 fallback suggestions, got %d", len(resp.Suggestions))
	}

	// First suggestion should be the original line
	if resp.Suggestions[0] != "Original line here" {
		t.Errorf("Expected first suggestion to be the original line, got '%s'", resp.Suggestions[0])
	}
}

func TestQuickIdeaAgent_GenerateFallbackCheck(t *testing.T) {
	agent := NewQuickIdeaAgent()

	req := QuickRequest{
		Mode:    QuickIdeaModeCheck,
		Context: "Some line here",
		Options: map[string]string{},
	}

	resp := agent.generateFallback(req)

	if resp == nil {
		t.Fatal("generateFallback returned nil")
	}

	if resp.Rating != "OKAY" {
		t.Errorf("Expected rating 'OKAY', got '%s'", resp.Rating)
	}

	if resp.Tip != "Add vivid sensory image" {
		t.Errorf("Expected tip 'Add vivid sensory image', got '%s'", resp.Tip)
	}
}

// Mock client for testing
type mockQuickClient struct {
	response string
	err      bool
}

func (m *mockQuickClient) Generate(ctx context.Context, _, prompt string, options map[string]any) (string, error) {
	if m.err {
		return "", errors.New("mock client error")
	}

	// Simulate some delay
	time.Sleep(10 * time.Millisecond)

	return m.response, nil
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
