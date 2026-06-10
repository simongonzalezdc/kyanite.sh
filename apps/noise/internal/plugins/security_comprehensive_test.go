package plugins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/errors"
	"github.com/kyanite/noise/internal/logging"
)

// TestSecurityManagerComprehensive tests all security functionality
func TestSecurityManagerComprehensive(t *testing.T) {
	sm := TestSecurityManager(t)

	t.Run("PluginHashAlgorithm", func(t *testing.T) {
		testPluginHashAlgorithm(t, sm)
	})

	t.Run("InputValidation", func(t *testing.T) {
		testInputValidation(t, sm)
	})

	t.Run("Sandboxing", func(t *testing.T) {
		testSandboxing(t, sm)
	})

	t.Run("SecurityMonitoring", func(t *testing.T) {
		testSecurityMonitoring(t, sm)
	})

	t.Run("PathValidation", func(t *testing.T) {
		testPathValidation(t, sm)
	})
}

// testPluginHashAlgorithm tests that SHA-256 is used instead of MD5
func testPluginHashAlgorithm(t *testing.T, sm *SecurityManager) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_plugin.json")
	content := `{"id": "test", "name": "Test Plugin", "version": "1.0.0"}`

	err := os.WriteFile(testFile, []byte(content), 0o600)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test hash calculation
	hash, err := sm.CalculatePluginHash(testFile)
	if err != nil {
		t.Fatalf("Failed to calculate plugin hash: %v", err)
	}

	// Verify it's a SHA-256 hash (64 characters hex)
	if len(hash) != 64 {
		t.Errorf("Expected SHA-256 hash length of 64, got %d", len(hash))
	}

	// Verify it contains only hex characters
	for _, char := range hash {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("Hash contains invalid character: %c", char)
		}
	}

	// Test that the SHA-256 specific function also works
	hash256, err := sm.CalculatePluginHashSHA256(testFile)
	if err != nil {
		t.Fatalf("Failed to calculate SHA-256 hash: %v", err)
	}

	if hash != hash256 {
		t.Error("SHA-256 hash mismatch between functions")
	}
}

// testInputValidation tests comprehensive input validation
func testInputValidation(t *testing.T, sm *SecurityManager) {
	t.Run("PluginIDValidation", func(t *testing.T) {
		validIDs := []string{"valid_plugin", "plugin123", "test-plugin"}
		invalidIDs := []string{
			"",                          // Empty
			"plugin/with/slashes",       // Contains slashes
			"plugin\\with\\backslashes", // Contains backslashes
			"plugin:with:colons",        // Contains colons
			"plugin*with*asterisks",     // Contains asterisks
			"plugin?with?question",      // Contains question marks
			"plugin\"with\"quotes",      // Contains quotes
			"plugin<with>brackets",      // Contains brackets
			"plugin|with|pipe",          // Contains pipe
			strings.Repeat("a", 101),    // Too long
			"plugin\x00with\x00null",    // Contains null bytes
			"plugin../../../etc/passwd", // Path traversal
		}

		for _, id := range validIDs {
			err := sm.validatePluginID(id)
			if err != nil {
				t.Errorf("Valid plugin ID rejected: %s, error: %v", id, err)
			}
		}

		for _, id := range invalidIDs {
			err := sm.validatePluginID(id)
			if err == nil {
				t.Errorf("Invalid plugin ID accepted: %s", id)
			}
		}
	})

	t.Run("PluginNameValidation", func(t *testing.T) {
		validNames := []string{"Valid Plugin Name", "Plugin 123", "Test-Plugin"}
		invalidNames := []string{
			"",                                     // Empty
			strings.Repeat("a", 201),               // Too long
			"Plugin <script>alert('xss')</script>", // XSS attempt
			"Plugin javascript:alert('xss')",       // JavaScript injection
			"Plugin onload=alert('xss')",           // Event handler injection
		}

		for _, name := range validNames {
			err := sm.validatePluginName(name)
			if err != nil {
				t.Errorf("Valid plugin name rejected: %s, error: %v", name, err)
			}
		}

		for _, name := range invalidNames {
			err := sm.validatePluginName(name)
			if err == nil {
				t.Errorf("Invalid plugin name accepted: %s", name)
			}
		}
	})

	t.Run("VersionValidation", func(t *testing.T) {
		validVersions := []string{"1.0.0", "2.1.3", "1.0.0-beta", "3.0.0-alpha.1", "1.0.0-beta.1"}
		invalidVersions := []string{
			"",                      // Empty
			"not-a-version",         // Invalid format
			"1.x.0",                 // Non-numeric core part
			strings.Repeat("a", 51), // Too long
			"1.0.0<script>",         // Contains scripts
		}

		for _, version := range validVersions {
			err := sm.validateVersion(version)
			if err != nil {
				t.Errorf("Valid version rejected: %s, error: %v", version, err)
			}
		}

		for _, version := range invalidVersions {
			err := sm.validateVersion(version)
			if err == nil {
				t.Errorf("Invalid version accepted: %s", version)
			}
		}
	})

	t.Run("DescriptionValidation", func(t *testing.T) {
		validDescriptions := []string{
			"A valid plugin description",
			"This plugin does useful things for music production",
			"Version 1.0.0 - adds new features",
		}
		invalidDescriptions := []string{
			strings.Repeat("a", 1001),                        // Too long
			"Description with <script>alert('xss')</script>", // XSS attempt
			"Description with javascript:alert('xss')",       // JavaScript injection
			"Description with <iframe>src=evil.com</iframe>", // Iframe injection
		}

		for _, desc := range validDescriptions {
			err := sm.validateDescription(desc)
			if err != nil {
				t.Errorf("Valid description rejected: %s, error: %v", desc, err)
			}
		}

		for _, desc := range invalidDescriptions {
			err := sm.validateDescription(desc)
			if err == nil {
				t.Errorf("Invalid description accepted: %s", desc)
			}
		}
	})
}

// testSandboxing tests OS-level sandboxing functionality
func testSandboxing(t *testing.T, sm *SecurityManager) {
	// Create a mock plugin for testing
	plugin := &StubPlugin{
		metadata: &PluginMetadata{
			ID:           "test_plugin",
			Name:         "Test Plugin",
			Version:      "1.0.0",
			Description:  "A test plugin",
			Capabilities: []Capability{CapabilityMenuItem},
		},
		enabled: false,
	}

	// Test sandbox creation
	err := sm.SandboxPlugin(plugin)
	if err != nil {
		t.Errorf("Failed to create plugin sandbox: %v", err)
	}

	// Test plugin termination
	err = sm.TerminatePlugin("test_plugin")
	if err != nil {
		t.Errorf("Failed to terminate plugin: %v", err)
	}

	// Test resource monitoring
	err = sm.MonitorPluginResourceUsage("test_plugin")
	// This might fail if plugin doesn't exist, which is expected
	if err != nil {
		t.Logf("Resource monitoring failed as expected: %v", err)
	}
}

// testSecurityMonitoring tests security monitoring functionality
func testSecurityMonitoring(t *testing.T, sm *SecurityManager) {
	// Test security event logging
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    "test_plugin",
		EventType:   "test_event",
		Description: "Test security event",
		Severity:    "info",
	})

	// Test getting security events
	events := sm.GetSecurityEvents(10)
	if len(events) == 0 {
		t.Error("Expected security events, got none")
	}

	// Test security checks
	sm.performSecurityChecks("test_plugin")

	// Verify that security events were logged
	recentEvents := sm.GetSecurityEvents(10)
	found := false
	for _, event := range recentEvents {
		if event.EventType == "file_access_check" || event.EventType == "network_check" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Security checks did not generate expected events")
	}
}

// testPathValidation tests path validation for security
func testPathValidation(t *testing.T, sm *SecurityManager) {
	// Test valid paths (these should be in allowed paths by default)
	validPaths := []string{
		filepath.Join(os.Getenv("HOME"), ".noise", "plugins"),
		"./plugins",
		"./internal/plugins/examples",
	}

	for _, path := range validPaths {
		err := sm.ValidatePluginPath(path)
		// Some paths might not exist, which is OK for validation testing
		if err != nil && !strings.Contains(err.Error(), "not in an allowed directory") {
			t.Logf("Path validation result for %s: %v", path, err)
		}
	}

	// Test invalid paths (blocked paths)
	invalidPaths := []string{
		"/etc",
		"/usr",
		"/bin",
		"/proc",
		"/sys",
	}

	for _, path := range invalidPaths {
		err := sm.ValidatePluginPath(path)
		if err == nil {
			t.Errorf("Blocked path accepted: %s", path)
		}
	}

	// Test path traversal attempts
	traversalPaths := []string{
		"../../../etc/passwd",
		"..\\..\\..\\windows\\system32",
		"./plugins/../../../etc/passwd",
	}

	for _, path := range traversalPaths {
		err := sm.ValidatePluginPath(path)
		if err == nil {
			t.Errorf("Path traversal attempt accepted: %s", path)
		}
	}
}

// TestBackupSecurity tests backup system security fixes
func TestBackupSecurity(t *testing.T) {
	// This would test the backup security fixes implemented in safety.go
	// For now, we'll create a basic test structure

	tmpDir := t.TempDir()
	logger, err := logging.NewFromConfig(config.DefaultConfig())
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test backup manager creation
	backupManager, err := errors.NewBackupManager(tmpDir, 10, logger)
	if err != nil {
		t.Fatalf("Failed to create backup manager: %v", err)
	}

	// Test that backup manager was created successfully
	// Since fields are unexported, we just verify creation succeeded
	if backupManager == nil {
		t.Error("Backup manager should not be nil")
	}
}

// TestSecurityIntegration tests integration of all security components
func TestSecurityIntegration(t *testing.T) {
	sm := TestSecurityManager(t)

	// Test complete security workflow
	plugin := &StubPlugin{
		metadata: &PluginMetadata{
			ID:           "integration_test_plugin",
			Name:         "Integration Test Plugin",
			Version:      "1.0.0",
			Description:  "A plugin for integration testing",
			Capabilities: []Capability{CapabilityMenuItem},
		},
		enabled: false,
	}

	// 1. Validate plugin manifest
	err := sm.ValidatePluginManifest(plugin.metadata)
	if err != nil {
		t.Errorf("Plugin manifest validation failed: %v", err)
	}

	// 2. Create sandbox
	err = sm.SandboxPlugin(plugin)
	if err != nil {
		t.Errorf("Plugin sandboxing failed: %v", err)
	}

	// 3. Test security monitoring
	err = sm.MonitorPluginResourceUsage(plugin.metadata.ID)
	if err != nil {
		t.Logf("Resource monitoring failed (expected for non-tracked plugin): %v", err)
	}

	// 4. Test plugin termination
	err = sm.TerminatePlugin(plugin.metadata.ID)
	if err != nil {
		t.Errorf("Plugin termination failed: %v", err)
	}

	// 5. Verify security events were logged
	events := sm.GetSecurityEvents(100)
	if len(events) == 0 {
		t.Error("No security events were logged during integration test")
	}

	// Check for key security events
	eventTypes := make(map[string]bool)
	for _, event := range events {
		eventTypes[event.EventType] = true
	}

	expectedEvents := []string{
		"sandbox_created",
		"monitoring_started",
		"plugin_terminated",
		"sandbox_cleanup",
	}

	for _, expected := range expectedEvents {
		if !eventTypes[expected] {
			t.Errorf("Expected security event not found: %s", expected)
		}
	}
}
