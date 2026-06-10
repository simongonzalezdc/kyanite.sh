package cli

import (
	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/utils"
)

// initEngine initializes and returns a task engine with the default repository.
// This helper reduces code duplication across CLI commands.
func initEngine() *engine.Engine {
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	return engine.New(repo)
}

// initEngineAndAI initializes and returns both a task engine and AI manager.
// This helper reduces code duplication across CLI commands that use AI features.
func initEngineAndAI() (*engine.Engine, *ai.Manager) {
	repo := repository.NewStoreRepository(utils.GetStoragePath())
	eng := engine.New(repo)
	aiMgr := ai.New()
	return eng, aiMgr
}
