package plugins

import "time"

// Security constants
const (
	// Sandbox resource limits
	DefaultMaxMemory        = 50 * 1024 * 1024 // 50MB
	DefaultMaxCPUPercent    = 50               // 50% of one core
	DefaultMaxExecutionTime = 30 * time.Second

	// Validation limits
	MaxPluginNameLength        = 100
	MaxPluginDescriptionLength = 1000
	MaxPluginVersionLength     = 20
	MaxPluginAuthorLength      = 100

	// Security event retention
	MaxSecurityEvents      = 1000
	SecurityEventRetention = 7 * 24 * time.Hour // 7 days
)

// Plugin loading constants
const (
	// Timeouts
	PluginLoadTimeout     = 10 * time.Second
	PluginInitTimeout     = 5 * time.Second
	PluginShutdownTimeout = 5 * time.Second

	// Retry configuration
	MaxPluginLoadRetries = 3
	LoadRetryDelay       = 1 * time.Second
)

// Sandbox constants
const (
	// Command execution
	SandboxCommandTimeout = 30 * time.Second
	MaxCommandOutputSize  = 1024 * 1024 // 1MB
)
