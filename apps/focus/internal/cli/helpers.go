package cli

import (
	"context"
	"time"

	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/di"
	"github.com/kyanite/focus/internal/engine"
)

// aiCallTimeout bounds all AI operations to prevent indefinite hangs when
// the LLM backend (Ollama) is unreachable. See docs/ARCHITECTURE-BACKLOG.md T4-02.
const aiCallTimeout = 30 * time.Second

// withAITimeout returns a context bounded by aiCallTimeout plus the cancel func.
// Always pair with `defer cancel()` to avoid leaking the timer.
func withAITimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), aiCallTimeout)
}

// initEngine returns the singleton task engine from the DI container.
func initEngine() *engine.Engine {
	return di.GetContainer().GetEngine()
}

// initEngineAndAI returns the singleton engine and AI manager from the DI container.
func initEngineAndAI() (*engine.Engine, *ai.Manager) {
	return di.GetContainer().GetEngine(), di.GetContainer().GetAIManager()
}
