package app

import (
	"github.com/kyanite/ai"
	kyaniteconfig "github.com/kyanite/config"
	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/logging"
)

// newBrain creates a *ai.Brain using the app's noise configuration.
// If Brain creation fails (e.g. NUCBox unreachable), nil is returned and
// logged — the app must remain usable offline.
func newBrain(cfg *config.Config) *ai.Brain {
	root, _ := kyaniteconfig.Load()
	brainCfg := ai.ConfigFromRoot(root, "noise")
	if cfg != nil && cfg.AI.BaseURL != "" {
		brainCfg.OllamaURL = cfg.AI.BaseURL
	}
	if cfg != nil && cfg.AI.Model != "" {
		brainCfg.Model = cfg.AI.Model
	}
	brain, err := ai.New(brainCfg)
	if err != nil {
		logging.GetDefaultLogger().Warnf("brain init failed (offline mode): %v", err)
	}
	return brain
}