package di

import (
	"sync"

	"github.com/kyanite/focus/internal/ai"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/pkg/utils"
)

// Container holds all application dependencies
type Container struct {
	// Singleton instances
	repo      repository.Repository
	engine    *engine.Engine
	aiManager *ai.Manager

	// Mutex for thread-safe singleton creation
	mu sync.Mutex

	// Configuration
	storagePath string
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{
		storagePath: utils.GetStoragePath(),
	}
}

// NewContainerWithPath creates a container with a custom storage path
func NewContainerWithPath(storagePath string) *Container {
	return &Container{
		storagePath: storagePath,
	}
}

// GetRepository returns the singleton repository instance
func (c *Container) GetRepository() repository.Repository {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.repo == nil {
		c.repo = repository.NewStoreRepository(c.storagePath)
	}
	return c.repo
}

// GetEngine returns the singleton engine instance
func (c *Container) GetEngine() *engine.Engine {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.engine == nil {
		c.engine = engine.New(c.GetRepository())
		// Wire AI manager for cache invalidation
		if c.aiManager != nil {
			c.engine.SetAIManager(c.aiManager)
		}
	}
	return c.engine
}

// GetAIManager returns the singleton AI manager instance
func (c *Container) GetAIManager() *ai.Manager {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.aiManager == nil {
		c.aiManager = ai.New()
		// If engine already exists, wire it for cache invalidation
		if c.engine != nil {
			c.engine.SetAIManager(c.aiManager)
		}
	}
	return c.aiManager
}

// Reset clears all cached instances (useful for testing)
func (c *Container) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.repo = nil
	c.engine = nil
	c.aiManager = nil
}

// Close releases all resources
func (c *Container) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.repo != nil {
		if err := c.repo.Close(); err != nil {
			return err
		}
	}
	return nil
}
