package ai

// Error types for structured error handling across the inference brain.
// Apps can use errors.Is() to check for specific failure modes.

import "errors"

// Sentinel errors for the AI brain.
var (
	// ErrLLMUnavailable is returned when the Ollama server is unreachable.
	ErrLLMUnavailable = errors.New("llm server unavailable")

	// ErrLLMRequestFailed is returned when an LLM request fails after retries.
	ErrLLMRequestFailed = errors.New("llm request failed")

	// ErrSTTNotInstalled is returned when whisper.cpp binary is not found.
	ErrSTTNotInstalled = errors.New("stt: whisper.cpp not installed")

	// ErrSTTFailed is returned when whisper.cpp returns a non-zero exit code.
	ErrSTTFailed = errors.New("stt transcription failed")

	// ErrMemoryUnreachable is returned when PostgreSQL is unreachable.
	ErrMemoryUnreachable = errors.New("memory store unreachable")

	// ErrBrainNotInitialized is returned when a method is called on a nil brain.
	ErrBrainNotInitialized = errors.New("brain not initialized")

	// ErrRateLimited is returned when the server responds with 429.
	ErrRateLimited = errors.New("rate limited")

	// ErrContextCanceled is returned when the caller cancels the context.
	ErrContextCanceled = errors.New("operation canceled")
)
