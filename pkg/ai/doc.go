// Package ai provides a unified inference brain for the kyanite.sh TUI suite.
//
// It consolidates LLM inference (via Ollama on the NUCBox inference server),
// speech-to-text (via local whisper.cpp), and persistent memory (via PostgreSQL)
// behind a single Brain client that all four apps (focus, noise, syntax, prism) share.
//
// Architecture:
//
//	Local machine (MacBook Air / Mac Mini):
//	  ┌─────────────────────────────┐
//	  │  kyanite TUI app             │
//	  │  STT: whisper.cpp (local)    │
//	  │  LLM: HTTP → nucbox:11434    │
//	  └──────────────┬──────────────┘
//	                 │ tailnet
//	                 ▼
//	  ┌──────────────────────────────┐
//	  │  NUCBox                      │
//	  │  Ollama (gemma4:12b on GPU)  │
//	  │  PostgreSQL (memory store)   │
//	  └──────────────────────────────┘
package ai
