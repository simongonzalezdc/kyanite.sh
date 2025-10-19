//go:build !debug

package logging

// buildDebugLogging reports whether the binary was built with the "debug" build tag.
// Release builds compile this file (with !debug constraint), so debug logging is
// disabled unless explicitly enabled through configuration or environment variables.
var buildDebugLogging = false
