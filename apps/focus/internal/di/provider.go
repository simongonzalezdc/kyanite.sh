package di

import (
	"sync"
)

var (
	// Global container instance (singleton)
	globalContainer *Container
	once            sync.Once
)

// GetContainer returns the global container instance
func GetContainer() *Container {
	once.Do(func() {
		globalContainer = NewContainer()
	})
	return globalContainer
}

// SetContainer allows setting a custom container (useful for testing)
func SetContainer(container *Container) {
	globalContainer = container
}

// ResetContainer resets the global container (useful for testing)
func ResetContainer() {
	if globalContainer != nil {
		_ = globalContainer.Close()
	}
	globalContainer = nil
	once = sync.Once{}
}
