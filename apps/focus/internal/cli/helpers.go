package cli

import (
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/di"
	"github.com/kyanite/focus/internal/engine"
)

// defaultContainer is the singleton DI container for CLI commands.
var defaultContainer = di.NewContainer()

// initEngine returns the singleton task engine from the DI container.
func initEngine() *engine.Engine {
	return defaultContainer.GetEngine()
}

// initEngineAndAI returns the singleton engine and AI manager from the DI container.
func initEngineAndAI() (*engine.Engine, *ai.Manager) {
	return defaultContainer.GetEngine(), defaultContainer.GetAIManager()
}
