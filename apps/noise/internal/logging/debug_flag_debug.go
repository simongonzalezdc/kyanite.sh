//go:build debug

package logging

// buildDebugLogging reports whether the binary was built with the "debug" build tag.
// Debug builds compile this file (with the debug constraint), enabling debug-level
// logging regardless of runtime configuration.
var buildDebugLogging = true
