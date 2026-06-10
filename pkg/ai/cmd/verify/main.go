package main

import (
	"context"
	"fmt"
	"time"

	ai "github.com/kyanite/ai"
)

func main() {
	cfg := ai.DefaultConfig("verify")
	cfg.OllamaURL = "http://nucbox:11434"
	cfg.Model = "gemma4:12b"
	cfg.Timeout = 30 * time.Second

	brain, err := ai.New(cfg)
	if err != nil {
		fmt.Printf("FAIL: Brain init error: %v\n", err)
		return
	}
	defer brain.Close()

	ctx := context.Background()

	fmt.Printf("LLM available: %v\n", brain.IsLLMAvailable(ctx))

	resp, err := brain.Generate(ctx, "Say exactly: hello world", ai.WithMaxTokens(20))
	if err != nil {
		fmt.Printf("FAIL: LLM error: %v\n", err)
		return
	}
	fmt.Printf("LLM response length=%d content=[%s]\n", len(resp), resp)
	if resp == "" {
		fmt.Println("BUG: empty response from LLM!")
	}

	resp2, err := brain.Generate(ctx, ai.PrismPalettePrompt("sunset over mountains"), ai.WithJSONMode(), ai.WithMaxTokens(100))
	if err != nil {
		fmt.Printf("FAIL: JSON mode error: %v\n", err)
	} else {
		fmt.Printf("JSON response length=%d content=[%s]\n", len(resp2), resp2)
	}

	fmt.Printf("STT available: %v\n", brain.IsSTTAvailable())
	fmt.Printf("Memory available: %v\n", brain.IsMemoryAvailable(ctx))
	fmt.Println("=== DONE ===")
}
