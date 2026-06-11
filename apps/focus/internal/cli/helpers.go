package cli

import (
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/di"
	"github.com/kyanite/focus/internal/engine"
)

// initEngine returns the singleton task engine from the DI container.
func initEngine() *engine.Engine {
	return di.GetContainer().GetEngine()
}

// initEngineAndAI returns the singleton engine and AI manager from the DI container.
func initEngineAndAI() (*engine.Engine, *ai.Manager) {
	return di.GetContainer().GetEngine(), di.GetContainer().GetAIManager()
}
