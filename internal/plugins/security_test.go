package plugins

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSecurityManager_ValidatePluginPath(t *testing.T) {
	sm := TestSecurityManager(t)

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Allowed path - ./plugins",
			path:        "./plugins/test.json",
			expectError: false,
		},
		{
			name:        "Allowed path - examples",
			path:        "./internal/plugins/examples/test.json",
			expectError: false,
		},
		{
			name:        "Blocked path - /etc",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "blocked directory",
		},
		{
			name:        "Blocked path - /usr",
			path:        "/usr/bin/test.so",
			expectError: true,
			errorMsg:    "blocked directory",
		},
		{
			name:        "Blocked path - home directory",
			path:        filepath.Join(os.Getenv("HOME"), "test.json"),
			expectError: true,
			errorMsg:    "blocked directory",
		},
		{
			name:        "Non-existent path",
			path:        "/non/existent/path/test.json",
			expectError: true,
			errorMsg:    "not in an allowed directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidatePluginPath(tt.path)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecurityManager_ValidatePluginFile(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	tests := []struct {
		name        string
		fileCreator func(string) string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid JSON file",
			fileCreator: func(dir string) string {
				return CreateTestPluginManifest(t, dir, nil)
			},
			expectError: false,
		},
		{
			name: "Valid SO file",
			fileCreator: func(dir string) string {
				return CreateTestPluginFile(t, dir, "test.so", "fake plugin content")
			},
			expectError: false,
		},
		{
			name: "Valid DLL file",
			fileCreator: func(dir string) string {
				return CreateTestPluginFile(t, dir, "test.dll", "fake plugin content")
			},
			expectError: false,
		},
		{
			name: "Invalid extension - .exe",
			fileCreator: func(dir string) string {
				return CreateTestPluginFile(t, dir, "test.exe", "malicious content")
			},
			expectError: true,
			errorMsg:    "extension .exe is not allowed",
		},
		{
			name: "Invalid extension - .py",
			fileCreator: func(dir string) string {
				return CreateTestPluginFile(t, dir, "test.py", "#!/usr/bin/env python")
			},
			expectError: true,
			errorMsg:    "extension .py is not allowed",
		},
		{
			name: "File too large",
			fileCreator: func(dir string) string {
				return CreateLargePluginFile(t, dir, "large.json", 15*1024*1024) // 15MB
			},
			expectError: true,
			errorMsg:    "exceeds maximum allowed size",
		},
		{
			name: "World-writable file",
			fileCreator: func(dir string) string {
				filePath := filepath.Join(dir, "writable.json")
				if err := os.WriteFile(filePath, []byte("{}"), 0666); err != nil {
					t.Fatalf("Failed to create world-writable file: %v", err)
				}
				return filePath
			},
			expectError: true,
			errorMsg:    "unsafe permissions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := tt.fileCreator(testDir)
			err := sm.ValidatePluginFile(filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecurityManager_ValidatePluginManifest(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	tests := []struct {
		name        string
		metadata    *PluginMetadata
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid manifest",
			metadata: &PluginMetadata{
				ID:          "test_plugin",
				Name:        "Test Plugin",
				Version:     "1.0.0",
				Description: "A test plugin",
				Author:      "Test Author",
				License:     "MIT",
				Capabilities: []Capability{
					CapabilityExportFormat,
				},
			},
			expectError: false,
		},
		{
			name: "Missing ID",
			metadata: &PluginMetadata{
				Name:        "Test Plugin",
				Version:     "1.0.0",
				Description: "A test plugin",
				Author:      "Test Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "plugin ID is required",
		},
		{
			name: "Missing Name",
			metadata: &PluginMetadata{
				ID:          "test_plugin",
				Version:     "1.0.0",
				Description: "A test plugin",
				Author:      "Test Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "plugin name is required",
		},
		{
			name: "Missing Version",
			metadata: &PluginMetadata{
				ID:          "test_plugin",
				Name:        "Test Plugin",
				Description: "A test plugin",
				Author:      "Test Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "plugin version is required",
		},
		{
			name: "Invalid ID with path traversal",
			metadata: &PluginMetadata{
				ID:          "../../../etc/passwd",
				Name:        "Malicious Plugin",
				Version:     "1.0.0",
				Description: "A malicious plugin",
				Author:      "Malicious Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "contains invalid characters",
		},
		{
			name: "Suspicious description with script tag",
			metadata: &PluginMetadata{
				ID:          "suspicious_plugin",
				Name:        "Suspicious Plugin",
				Version:     "1.0.0",
				Description: "This plugin contains <script>alert('xss')</script> code",
				Author:      "Test Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "suspicious pattern",
		},
		{
			name: "Suspicious description with eval",
			metadata: &PluginMetadata{
				ID:          "eval_plugin",
				Name:        "Eval Plugin",
				Version:     "1.0.0",
				Description: "This plugin uses eval(malicious_code) to execute code",
				Author:      "Test Author",
				License:     "MIT",
			},
			expectError: true,
			errorMsg:    "suspicious pattern",
		},
		{
			name: "Invalid capability",
			metadata: &PluginMetadata{
				ID:          "invalid_cap_plugin",
				Name:        "Invalid Cap Plugin",
				Version:     "1.0.0",
				Description: "A plugin with invalid capabilities",
				Author:      "Test Author",
				License:     "MIT",
				Capabilities: []Capability{
					Capability("system_access"),
					Capability("file_delete"),
				},
			},
			expectError: true,
			errorMsg:    "invalid capability",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sm.ValidatePluginManifest(tt.metadata)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestSecurityManager_CalculatePluginHash(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Create a test file
	filePath := CreateTestPluginFile(t, testDir, "test.json", `{"test": "content"}`)

	// Calculate hash
	hash1, err := sm.CalculatePluginHash(filePath)
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}

	// Calculate hash again - should be the same
	hash2, err := sm.CalculatePluginHash(filePath)
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("Hashes should be identical: %s != %s", hash1, hash2)
	}

	// Modify the file and recalculate
	if err := os.WriteFile(filePath, []byte(`{"modified": "content"}`), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	hash3, err := sm.CalculatePluginHash(filePath)
	if err != nil {
		t.Fatalf("Failed to calculate hash: %v", err)
	}

	if hash1 == hash3 {
		t.Errorf("Hashes should be different after modification: %s == %s", hash1, hash3)
	}
}

func TestSecurityManager_VerifyPluginIntegrity(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Create a test file
	filePath := CreateTestPluginFile(t, testDir, "test.json", `{"test": "content"}`)

	// Verify integrity - should not error
	err := sm.VerifyPluginIntegrity(filePath)
	if err != nil {
		t.Errorf("Unexpected error verifying integrity: %v", err)
	}

	// Test with non-existent file
	err = sm.VerifyPluginIntegrity(filepath.Join(testDir, "nonexistent.json"))
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestSecurityManager_SandboxPlugin(t *testing.T) {
	sm := TestSecurityManager(t)

	// Test with valid plugin
	validPlugin := CreateMockPlugin("test_plugin", "Test Plugin", true)
	err := sm.SandboxPlugin(validPlugin)
	if err != nil {
		t.Errorf("Unexpected error sandboxing valid plugin: %v", err)
	}

	// Test with invalid plugin (invalid ID in metadata)
	invalidPlugin := &StubPlugin{
		metadata: &PluginMetadata{
			ID:          "../../../etc/passwd",
			Name:        "Invalid Plugin",
			Version:     "1.0.0",
			Description: "A plugin with invalid ID",
			Author:      "Test Author",
			License:     "MIT",
		},
		enabled: true,
	}

	err = sm.SandboxPlugin(invalidPlugin)
	if err == nil {
		t.Error("Expected error sandboxing invalid plugin")
	}
}

func TestSecurityManager_ScanForMaliciousPlugins(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Add test directory to allowed paths for this test
	sm.allowedPaths = append(sm.allowedPaths, testDir)

	// Create suspicious files
	CreateSuspiciousFiles(t, testDir)

	// Scan for malicious plugins
	threats := sm.ScanForMaliciousPlugins()

	if len(threats) == 0 {
		t.Error("Expected to find threats but found none")
	}

	// Check for expected threat types
	foundSuspicious := false
	foundWorldWritable := false
	foundScript := false

	for _, threat := range threats {
		if strings.Contains(threat, "Suspicious filename") {
			foundSuspicious = true
		}
		if strings.Contains(threat, "World-writable") {
			foundWorldWritable = true
		}
		if strings.Contains(threat, "script file") {
			foundScript = true
		}
	}

	if !foundSuspicious {
		t.Error("Expected to find suspicious filenames")
	}
	if !foundWorldWritable {
		t.Error("Expected to find world-writable files")
	}
	if !foundScript {
		t.Error("Expected to find script files")
	}
}

func TestSecurityManager_CleanupStalePlugins(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	// Add test directory to allowed paths for this test
	sm.allowedPaths = append(sm.allowedPaths, testDir)

	// Create a file with old timestamp
	oldFile := filepath.Join(testDir, "old_plugin.json")
	if err := os.WriteFile(oldFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create old file: %v", err)
	}

	// Set file modification time to 2 hours ago
	oldTime := time.Now().Add(-2 * time.Hour)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set file time: %v", err)
	}

	// Create a recent file
	recentFile := filepath.Join(testDir, "recent_plugin.json")
	if err := os.WriteFile(recentFile, []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create recent file: %v", err)
	}

	// Cleanup with max age of 1 hour (should remove old file but keep recent)
	err := sm.CleanupStalePlugins(1 * time.Hour)
	if err != nil {
		t.Errorf("Unexpected error during cleanup: %v", err)
	}

	// Check that old file was removed
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("Expected old file to be removed")
	}

	// Check that recent file still exists
	if _, err := os.Stat(recentFile); os.IsNotExist(err) {
		t.Error("Expected recent file to still exist")
	}
}

func TestSecurityManager_GetSecurityReport(t *testing.T) {
	sm := TestSecurityManager(t)
	manager := CreateTestPluginManager(t)

	// Add some test plugins
	plugin1 := CreateMockPlugin("secure_plugin", "Secure Plugin", true)
	plugin1.metadata.Author = "Secure Author"
	plugin2 := CreateMockPlugin("insecure_plugin", "Insecure Plugin", true)
	plugin2.metadata.Author = "" // Missing author makes it insecure

	manager.GetPlugins()["secure_plugin"] = plugin1
	manager.GetPlugins()["insecure_plugin"] = plugin2

	// Get security report
	report := sm.GetSecurityReport(manager)

	// Check report structure
	if total, ok := report["total_plugins"].(int); !ok || total != 2 {
		t.Errorf("Expected total_plugins to be 2, got %v", report["total_plugins"])
	}

	if secure, ok := report["secure_plugins"].([]string); !ok || len(secure) != 1 || secure[0] != "secure_plugin" {
		t.Errorf("Expected secure_plugins to contain only 'secure_plugin', got %v", report["secure_plugins"])
	}

	if insecure, ok := report["insecure_plugins"].([]string); !ok || len(insecure) != 1 || insecure[0] != "insecure_plugin" {
		t.Errorf("Expected insecure_plugins to contain only 'insecure_plugin', got %v", report["insecure_plugins"])
	}

	if _, ok := report["total_size_bytes"].(int64); !ok {
		t.Error("Expected total_size_bytes to be present")
	}

	if _, ok := report["scan_timestamp"].(time.Time); !ok {
		t.Error("Expected scan_timestamp to be present")
	}
}

func TestSecurityManager_isValidCapability(t *testing.T) {
	sm := TestSecurityManager(t)

	// Test valid capabilities
	validCapabilities := []Capability{
		CapabilityMenuItem,
		CapabilityScreen,
		CapabilityUIExtension,
		CapabilityEditorTool,
		CapabilityExportFormat,
		CapabilitySyntax,
		CapabilityTheoryTool,
		CapabilityChordLib,
		CapabilityAudioEffect,
		CapabilityMIDIHandler,
		CapabilityDataProvider,
		CapabilityExporter,
		CapabilityHook,
		CapabilityConfig,
	}

	for _, cap := range validCapabilities {
		if !sm.isValidCapability(cap) {
			t.Errorf("Expected capability %s to be valid", cap)
		}
	}

	// Test invalid capabilities
	invalidCapabilities := []Capability{
		Capability("system_access"),
		Capability("file_delete"),
		Capability("network_access"),
		Capability("privilege_escalation"),
	}

	for _, cap := range invalidCapabilities {
		if sm.isValidCapability(cap) {
			t.Errorf("Expected capability %s to be invalid", cap)
		}
	}
}

// Table-driven tests for various security scenarios
func TestSecurityManager_SecurityScenarios(t *testing.T) {
	sm := TestSecurityManager(t)
	testDir := CreateTestPluginDir(t)
	defer CleanupTestPluginDir(t, testDir)

	scenarios := []struct {
		name        string
		setup       func() string
		testFunc    func(sm *SecurityManager, path string) error
		expectError bool
		errorMsg    string
	}{
		{
			name: "Path traversal attempt",
			setup: func() string {
				return filepath.Join(testDir, "..", "..", "etc", "passwd")
			},
			testFunc: func(sm *SecurityManager, path string) error {
				return sm.ValidatePluginPath(path)
			},
			expectError: true,
			errorMsg:    "not in an allowed directory",
		},
		{
			name: "Executable file in plugin directory",
			setup: func() string {
				return CreateTestPluginFile(t, testDir, "malicious.exe", "fake executable")
			},
			testFunc: func(sm *SecurityManager, path string) error {
				return sm.ValidatePluginFile(path)
			},
			expectError: true,
			errorMsg:    "extension .exe is not allowed",
		},
		{
			name: "Oversized plugin file",
			setup: func() string {
				return CreateLargePluginFile(t, testDir, "oversized.json", 20*1024*1024) // 20MB
			},
			testFunc: func(sm *SecurityManager, path string) error {
				return sm.ValidatePluginFile(path)
			},
			expectError: true,
			errorMsg:    "exceeds maximum allowed size",
		},
		{
			name: "Plugin with system access capability",
			setup: func() string {
				metadata := &PluginMetadata{
					ID:          "system_plugin",
					Name:        "System Plugin",
					Version:     "1.0.0",
					Description: "A plugin that requests system access",
					Author:      "Test Author",
					License:     "MIT",
					Capabilities: []Capability{
						Capability("system_access"),
					},
				}
				return CreateTestPluginManifest(t, testDir, metadata)
			},
			testFunc: func(sm *SecurityManager, path string) error {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				var metadata PluginMetadata
				if err := json.Unmarshal(data, &metadata); err != nil {
					return err
				}
				return sm.ValidatePluginManifest(&metadata)
			},
			expectError: true,
			errorMsg:    "invalid capability",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			testPath := scenario.setup()
			err := scenario.testFunc(sm, testPath)

			if scenario.expectError {
				if err == nil {
					t.Errorf("Expected error but got none for scenario: %s", scenario.name)
					return
				}
				if scenario.errorMsg != "" && !strings.Contains(err.Error(), scenario.errorMsg) {
					t.Errorf("Expected error to contain '%s', got: %v for scenario: %s",
						scenario.errorMsg, err, scenario.name)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for scenario %s: %v", scenario.name, err)
				}
			}
		})
	}
}
