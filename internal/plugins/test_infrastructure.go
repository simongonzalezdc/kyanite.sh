package plugins

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/logging"
)

// CreateTestPluginManager creates a plugin manager for testing
func CreateTestPluginManager(t *testing.T) *DefaultManager {
	cfg := config.DefaultConfig()
	logger, err := logging.NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewManager(cfg, logger)
}

// TestSecurityManager creates a security manager for testing
func TestSecurityManager(t *testing.T) *SecurityManager {
	logger, err := logging.NewFromConfig(config.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	return NewSecurityManager(logger)
}

// CreateTestPluginDir creates a temporary plugin directory for testing
func CreateTestPluginDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "plugin_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp plugin dir: %v", err)
	}
	return dir
}

// CleanupTestPluginDir cleans up a temporary plugin directory
func CleanupTestPluginDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup test plugin dir: %v", err)
	}
}

// CreateTestPluginManifest creates a test plugin manifest file
func CreateTestPluginManifest(t *testing.T, dir string, metadata *PluginMetadata) string {
	if metadata == nil {
		metadata = &PluginMetadata{
			ID:          "test_plugin",
			Name:        "Test Plugin",
			Version:     "1.0.0",
			Description: "A test plugin for security testing",
			Author:      "Test Author",
			License:     "MIT",
			Capabilities: []Capability{
				CapabilityExportFormat,
			},
			Enabled: true,
		}
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal plugin metadata: %v", err)
	}

	manifestPath := filepath.Join(dir, metadata.ID+".json")
	if err := os.WriteFile(manifestPath, data, 0600); err != nil {
		t.Fatalf("Failed to write plugin manifest: %v", err)
	}

	return manifestPath
}

// CreateMaliciousPluginManifest creates a malicious plugin manifest for security testing
func CreateMaliciousPluginManifest(t *testing.T, dir string, maliciousType string) string {
	var metadata *PluginMetadata

	switch maliciousType {
	case "suspicious_description":
		metadata = &PluginMetadata{
			ID:          "malicious_plugin",
			Name:        "Malicious Plugin",
			Version:     "1.0.0",
			Description: "This plugin contains <script>alert('xss')</script> and eval(malicious_code)",
			Author:      "Malicious Author",
			License:     "MIT",
			Capabilities: []Capability{
				CapabilityExportFormat,
			},
			Enabled: true,
		}
	case "invalid_id":
		metadata = &PluginMetadata{
			ID:          "../../../etc/passwd",
			Name:        "Path Traversal Plugin",
			Version:     "1.0.0",
			Description: "A plugin that tries to access system files",
			Author:      "Malicious Author",
			License:     "MIT",
			Capabilities: []Capability{
				CapabilityExportFormat,
			},
			Enabled: true,
		}
	case "missing_fields":
		metadata = &PluginMetadata{
			// Missing required ID field
			Name:        "Incomplete Plugin",
			Version:     "1.0.0",
			Description: "A plugin with missing required fields",
			Author:      "Test Author",
			License:     "MIT",
			Capabilities: []Capability{
				CapabilityExportFormat,
			},
			Enabled: true,
		}
	case "invalid_capability":
		metadata = &PluginMetadata{
			ID:          "invalid_cap_plugin",
			Name:        "Invalid Capability Plugin",
			Version:     "1.0.0",
			Description: "A plugin with invalid capabilities",
			Author:      "Test Author",
			License:     "MIT",
			Capabilities: []Capability{
				Capability("system_access"),
				Capability("file_delete"),
			},
			Enabled: true,
		}
	default:
		t.Fatalf("Unknown malicious plugin type: %s", maliciousType)
	}

	return CreateTestPluginManifest(t, dir, metadata)
}

// CreateTestPluginFile creates a test plugin file with specified content
func CreateTestPluginFile(t *testing.T, dir, filename, content string) string {
	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, []byte(content), 0600); err != nil {
		t.Fatalf("Failed to write test plugin file: %v", err)
	}
	return filePath
}

// CreateLargePluginFile creates a large plugin file for testing size limits
func CreateLargePluginFile(t *testing.T, dir, filename string, size int64) string {
	filePath := filepath.Join(dir, filename)

	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create large plugin file: %v", err)
	}
	defer file.Close()

	// Write random data to reach the desired size
	data := make([]byte, 1024)
	for i := 0; i < len(data); i++ {
		data[i] = byte(i % 256)
	}

	written := int64(0)
	for written < size {
		n, err := file.Write(data)
		if err != nil {
			t.Fatalf("Failed to write to large plugin file: %v", err)
		}
		written += int64(n)
	}

	return filePath
}

// CreateSuspiciousFiles creates suspicious files for security testing
func CreateSuspiciousFiles(t *testing.T, dir string) []string {
	var files []string

	// Create files with suspicious names
	suspiciousNames := []string{
		"malware.so",
		"virus.dll",
		"trojan.json",
		"backdoor.so",
		"exploit.sh",
		"hack.py",
		"steal.js",
		"keylog.exe",
		"rootkit.so",
		"worm.dll",
	}

	for _, name := range suspiciousNames {
		filePath := filepath.Join(dir, name)
		if err := os.WriteFile(filePath, []byte("suspicious content"), 0600); err != nil {
			t.Fatalf("Failed to create suspicious file %s: %v", name, err)
		}
		files = append(files, filePath)
	}

	// Create world-writable files
	worldWritablePath := filepath.Join(dir, "world_writable.json")
	if err := os.WriteFile(worldWritablePath, []byte("{}"), 0600); err != nil {
		t.Fatalf("Failed to create world-writable file: %v", err)
	}
	// Make it world-writable using chmod
	if err := os.Chmod(worldWritablePath, 0666); err != nil {
		t.Fatalf("Failed to make file world-writable: %v", err)
	}
	files = append(files, worldWritablePath)

	// Create script files
	scriptFiles := []string{"malicious.sh", "hack.py", "exploit.js"}
	for _, name := range scriptFiles {
		filePath := filepath.Join(dir, name)
		if err := os.WriteFile(filePath, []byte("#!/bin/bash\necho 'malicious'"), 0600); err != nil {
			t.Fatalf("Failed to create script file %s: %v", name, err)
		}
		files = append(files, filePath)
	}

	return files
}

// CalculateFileHash calculates SHA-256 hash of a file
func CalculateFileHash(t *testing.T, filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %v", filePath, err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		t.Fatalf("Failed to calculate hash for file %s: %v", filePath, err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

// CreateMockPlugin creates a mock plugin for testing
func CreateMockPlugin(id, name string, enabled bool) *StubPlugin {
	return &StubPlugin{
		metadata: &PluginMetadata{
			ID:          id,
			Name:        name,
			Version:     "1.0.0",
			Description: "A mock plugin for testing",
			Author:      "Test Author",
			License:     "MIT",
			Capabilities: []Capability{
				CapabilityExportFormat,
			},
			Enabled:  enabled,
			LoadTime: time.Now(),
		},
		enabled: enabled,
	}
}

// CreateMockMaliciousPlugin creates a mock malicious plugin for security testing
func CreateMockMaliciousPlugin(id string, maliciousType string) *MaliciousPlugin {
	return &MaliciousPlugin{
		StubPlugin: &StubPlugin{
			metadata: &PluginMetadata{
				ID:          id,
				Name:        "Malicious Plugin",
				Version:     "1.0.0",
				Description: "A malicious plugin for security testing",
				Author:      "Malicious Author",
				License:     "MIT",
				Capabilities: []Capability{
					CapabilityExportFormat,
				},
				Enabled:  true,
				LoadTime: time.Now(),
			},
			enabled: true,
		},
		maliciousType: maliciousType,
	}
}

// MaliciousPlugin is a plugin that simulates malicious behavior for testing
type MaliciousPlugin struct {
	*StubPlugin
	maliciousType string
}

// Initialize simulates malicious initialization behavior
func (p *MaliciousPlugin) Initialize(ctx *PluginContext) error {
	switch p.maliciousType {
	case "resource_exhaustion":
		// Simulate resource exhaustion
		data := make([]byte, 1024*1024) // Allocate 1MB
		for i := range data {
			data[i] = byte(i % 256)
		}
	case "file_access":
		// Simulate unauthorized file access attempt
		paths := []string{"/etc/passwd", "/etc/shadow", os.Getenv("HOME") + "/.ssh/id_rsa"}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				// File exists, in a real malicious plugin this might read the file
				break
			}
		}
	case "network_access":
		// Simulate network access attempt
		// In a real malicious plugin this might make network calls
		break
	}

	return p.StubPlugin.Initialize(ctx)
}

// Cleanup simulates malicious cleanup behavior
func (p *MaliciousPlugin) Cleanup() error {
	switch p.maliciousType {
	case "data_leak":
		// Simulate data leak during cleanup
		// In a real malicious plugin this might exfiltrate data
		break
	}

	return p.StubPlugin.Cleanup()
}

// TestPluginContext creates a plugin context for testing
func TestPluginContext(t *testing.T) *PluginContext {
	cfg := config.DefaultConfig()
	return &PluginContext{
		Context: nil, // Use nil context for testing
		Config:  cfg,
		Metadata: &PluginMetadata{
			ID:      "test_plugin",
			Name:    "Test Plugin",
			Version: "1.0.0",
		},
	}
}

// AssertPluginSecurityValidated asserts that a plugin passes security validation
func AssertPluginSecurityValidated(t *testing.T, sm *SecurityManager, plugin Plugin) {
	if err := sm.ValidatePluginManifest(plugin.Metadata()); err != nil {
		t.Errorf("Plugin security validation failed: %v", err)
	}

	if err := sm.SandboxPlugin(plugin); err != nil {
		t.Errorf("Plugin sandboxing failed: %v", err)
	}
}

// AssertPluginSecurityRejected asserts that a plugin fails security validation
func AssertPluginSecurityRejected(t *testing.T, sm *SecurityManager, plugin Plugin, expectedError string) {
	err := sm.ValidatePluginManifest(plugin.Metadata())
	if err == nil {
		t.Error("Expected plugin security validation to fail, but it passed")
		return
	}

	if expectedError != "" && !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', but got: %v", expectedError, err)
	}
}

// CreateTestPluginWithCapabilities creates a test plugin with specified capabilities
func CreateTestPluginWithCapabilities(t *testing.T, id string, capabilities []Capability) Plugin {
	return &StubPlugin{
		metadata: &PluginMetadata{
			ID:           id,
			Name:         "Test Plugin",
			Version:      "1.0.0",
			Description:  "A test plugin",
			Author:       "Test Author",
			License:      "MIT",
			Capabilities: capabilities,
			Enabled:      true,
			LoadTime:     time.Now(),
		},
		enabled: true,
	}
}

// BenchmarkPluginOperation benchmarks a plugin operation
func BenchmarkPluginOperation(b *testing.B, operation func() error) {
	var err error
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err = operation(); err != nil {
			b.Fatalf("Operation failed: %v", err)
		}
	}
}
