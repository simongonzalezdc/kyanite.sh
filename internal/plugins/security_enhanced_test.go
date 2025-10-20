package plugins

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnhancedSecurityManager tests the enhanced security manager functionality
func TestEnhancedSecurityManager(t *testing.T) {
	logger, err := logging.New(logging.DefaultConfig())
	require.NoError(t, err)
	sm := NewSecurityManager(logger)
	
	t.Run("SandboxPolicyCreation", func(t *testing.T) {
		pluginID := "test_plugin"
		policy := sm.getOrCreateSandboxPolicy(pluginID)
		
		assert.Equal(t, pluginID, policy.PluginID)
		assert.False(t, policy.AllowFileAccess, "File access should be disabled by default")
		assert.False(t, policy.AllowNetwork, "Network access should be disabled by default")
		assert.Greater(t, policy.MaxMemory, int64(0), "Max memory should be set")
		assert.Greater(t, policy.MaxCPU, int64(0), "Max CPU should be set")
		assert.Greater(t, policy.MaxExecutionTime, time.Duration(0), "Max execution time should be set")
	})
	
	t.Run("PluginPermissionsCreation", func(t *testing.T) {
		pluginID := "test_plugin"
		perms := sm.getOrCreatePluginPermissions(pluginID)
		
		assert.Equal(t, pluginID, perms.PluginID)
		assert.False(t, perms.CanReadFiles, "File reading should be disabled by default")
		assert.False(t, perms.CanWriteFiles, "File writing should be disabled by default")
		assert.False(t, perms.CanExecute, "Execution should be disabled by default")
		assert.False(t, perms.CanNetwork, "Network should be disabled by default")
	})
	
	t.Run("ResourceTrackingInitialization", func(t *testing.T) {
		pluginID := "test_plugin"
		sm.initializeResourceTracking(pluginID)
		
		sm.usageMutex.RLock()
		usage, exists := sm.pluginResourceUsage[pluginID]
		sm.usageMutex.RUnlock()
		
		assert.True(t, exists, "Resource usage should be tracked")
		assert.Equal(t, pluginID, usage.PluginID)
		assert.False(t, usage.StartTime.IsZero(), "Start time should be set")
	})
	
	t.Run("SecurityEventLogging", func(t *testing.T) {
		event := SecurityEvent{
			Timestamp:   time.Now(),
			PluginID:    "test_plugin",
			EventType:   "test_event",
			Description: "Test security event",
			Severity:    "info",
		}
		
		sm.logSecurityEvent(event)
		
		events := sm.GetSecurityEvents(10)
		assert.Len(t, events, 1, "Should have one security event")
		assert.Equal(t, "test_event", events[0].EventType)
	})
	
	t.Run("ResourceMonitoring", func(t *testing.T) {
		pluginID := "test_plugin"
		sm.initializeResourceTracking(pluginID)
		
		err := sm.MonitorPluginResourceUsage(pluginID)
		assert.NoError(t, err, "Resource monitoring should not error")
		
		sm.usageMutex.RLock()
		usage := sm.pluginResourceUsage[pluginID]
		sm.usageMutex.RUnlock()
		
		assert.Greater(t, usage.MemoryUsage, int64(0), "Memory usage should be tracked")
		assert.Greater(t, usage.CPUUsage, int64(0), "CPU usage should be tracked")
	})
	
	t.Run("PathAccessControl", func(t *testing.T) {
		pluginID := "test_plugin"
		
		// Test with default permissions (should be denied)
		allowed := sm.IsPluginAllowedToAccessPath(pluginID, "/etc/passwd")
		assert.False(t, allowed, "Should not allow access to system files")
		
		// Test with allowed path
		tempDir, err := os.MkdirTemp("", "security_test")
		require.NoError(t, err)
		defer os.RemoveAll(tempDir)
		
		// Update policy to allow access to temp dir
		policy := sm.getOrCreateSandboxPolicy(pluginID)
		policy.AllowFileAccess = true
		policy.AllowedPaths = []string{tempDir}
		
		// Update permissions
		perms := sm.getOrCreatePluginPermissions(pluginID)
		perms.CanReadFiles = true
		
		testFile := filepath.Join(tempDir, "test.txt")
		allowed = sm.IsPluginAllowedToAccessPath(pluginID, testFile)
		assert.True(t, allowed, "Should allow access to files in allowed directory")
	})
	
	t.Run("NetworkAccessControl", func(t *testing.T) {
		pluginID := "test_plugin"
		
		// Test with default permissions (should be denied)
		allowed := sm.IsPluginAllowedToAccessNetwork(pluginID)
		assert.False(t, allowed, "Should not allow network access by default")
		
		// Enable network access in policy and permissions
		policy := sm.getOrCreateSandboxPolicy(pluginID)
		policy.AllowNetwork = true
		
		perms := sm.getOrCreatePluginPermissions(pluginID)
		perms.CanNetwork = true
		
		// Still should be denied due to global network policy
		allowed = sm.IsPluginAllowedToAccessNetwork(pluginID)
		assert.False(t, allowed, "Should not allow network access due to global policy")
		
		// Enable global network policy
		sm.networkPolicy.allowNetwork = true
		allowed = sm.IsPluginAllowedToAccessNetwork(pluginID)
		assert.True(t, allowed, "Should allow network access when all policies allow it")
	})
	
	t.Run("PluginTermination", func(t *testing.T) {
		pluginID := "test_plugin"
		
		err := sm.TerminatePlugin(pluginID)
		assert.NoError(t, err, "Plugin termination should not error")
		
		// Check that termination event was logged
		events := sm.GetSecurityEvents(10)
		found := false
		for _, event := range events {
			if event.PluginID == pluginID && event.EventType == "plugin_terminated" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should log plugin termination event")
	})
	
	t.Run("CommandExecutionInSandbox", func(t *testing.T) {
		pluginID := "test_plugin"
		
		// Enable execution permissions
		perms := sm.getOrCreatePluginPermissions(pluginID)
		perms.CanExecute = true
		
		// Test command execution
		output, err := sm.ExecutePluginInSandbox(pluginID, []string{"echo", "test"}, time.Second)
		assert.NoError(t, err, "Command execution should not error")
		assert.NotEmpty(t, output, "Should have output")
		
		// Check that execution event was logged
		events := sm.GetSecurityEvents(10)
		found := false
		for _, event := range events {
			if event.PluginID == pluginID && event.EventType == "command_execution" {
				found = true
				break
			}
		}
		assert.True(t, found, "Should log command execution event")
	})
	
	t.Run("SHA256HashCalculation", func(t *testing.T) {
		// Create a test file
		tempFile, err := os.CreateTemp("", "hash_test")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		testContent := "test content for hashing"
		_, err = tempFile.WriteString(testContent)
		require.NoError(t, err)
		tempFile.Close()
		
		// Calculate hash
		hash, err := sm.CalculatePluginHashSHA256(tempFile.Name())
		assert.NoError(t, err, "Hash calculation should not error")
		assert.NotEmpty(t, hash, "Hash should not be empty")
		assert.Equal(t, 64, len(hash), "SHA256 hash should be 64 characters")
	})
	
	t.Run("SignatureVerification", func(t *testing.T) {
		// Create a test file
		tempFile, err := os.CreateTemp("", "sig_test")
		require.NoError(t, err)
		defer os.Remove(tempFile.Name())
		
		// Test signature verification (will be a placeholder in current implementation)
		err = sm.VerifyPluginSignature(tempFile.Name())
		assert.NoError(t, err, "Signature verification should not error")
	})
}

// TestResourceViolationHandling tests resource violation handling
func TestResourceViolationHandling(t *testing.T) {
	logger, err := logging.New(logging.DefaultConfig())
	require.NoError(t, err)
	sm := NewSecurityManager(logger)
	
	pluginID := "test_plugin"
	sm.initializeResourceTracking(pluginID)
	
	// Set a low memory limit for testing
	policy := sm.getOrCreateSandboxPolicy(pluginID)
	policy.MaxMemory = 1024 // 1KB
	
	// Simulate a memory violation
	sm.handleResourceViolation(pluginID, "memory", 2048, 1024)
	
	// Check that violation event was logged
	events := sm.GetSecurityEvents(10)
	found := false
	for _, event := range events {
		if event.PluginID == pluginID && event.EventType == "resource_violation" {
			found = true
			assert.Contains(t, event.Description, "memory")
			break
		}
	}
	assert.True(t, found, "Should log resource violation event")
}

// TestSandboxPolicyEnforcement tests sandbox policy enforcement
func TestSandboxPolicyEnforcement(t *testing.T) {
	logger, err := logging.New(logging.DefaultConfig())
	require.NoError(t, err)
	sm := NewSecurityManager(logger)
	
	// Create a mock plugin
	plugin := &StubPlugin{
		metadata: &PluginMetadata{
			ID:          "test_plugin",
			Name:        "Test Plugin",
			Version:     "1.0.0",
			Description: "A test plugin",
			Author:      "Test Author",
		},
		enabled: true,
	}
	
	// Sandbox the plugin
	err = sm.SandboxPlugin(plugin)
	assert.NoError(t, err, "Sandboxing should not error")
	
	// Check that sandbox policy was created
	policy, exists := sm.sandboxPolicies["test_plugin"]
	assert.True(t, exists, "Sandbox policy should be created")
	assert.Equal(t, "test_plugin", policy.PluginID)
	
	// Check that permissions were created
	perms, exists := sm.pluginPermissions["test_plugin"]
	assert.True(t, exists, "Plugin permissions should be created")
	assert.Equal(t, "test_plugin", perms.PluginID)
	
	// Check that resource tracking was initialized
	sm.usageMutex.RLock()
	_, exists = sm.pluginResourceUsage["test_plugin"]
	sm.usageMutex.RUnlock()
	assert.True(t, exists, "Resource tracking should be initialized")
}

// TestSecurityEventHistory tests security event history management
func TestSecurityEventHistory(t *testing.T) {
	logger, err := logging.New(logging.DefaultConfig())
	require.NoError(t, err)
	sm := NewSecurityManager(logger)
	
	// Add multiple events
	for i := 0; i < 1050; i++ {
		event := SecurityEvent{
			Timestamp:   time.Now(),
			PluginID:    "test_plugin",
			EventType:   "test_event",
			Description: "Test security event",
			Severity:    "info",
		}
		sm.logSecurityEvent(event)
	}
	
	// Check that event history is limited
	events := sm.GetSecurityEvents(2000)
	assert.LessOrEqual(t, len(events), 1000, "Event history should be limited to 1000 events")
	
	// Check that recent events are available
	recentEvents := sm.GetSecurityEvents(10)
	assert.Len(t, recentEvents, 10, "Should get recent events")
}

// TestNetworkPolicy tests network policy configuration
func TestNetworkPolicy(t *testing.T) {
	logger, err := logging.New(logging.DefaultConfig())
	require.NoError(t, err)
	sm := NewSecurityManager(logger)
	
	// Test default network policy
	assert.False(t, sm.networkPolicy.allowNetwork, "Network should be disabled by default")
	assert.Empty(t, sm.networkPolicy.allowedHosts, "No hosts should be allowed by default")
	assert.Empty(t, sm.networkPolicy.allowedPorts, "No ports should be allowed by default")
	
	// Enable network policy
	sm.networkPolicy.allowNetwork = true
	sm.networkPolicy.allowedHosts = []string{"example.com"}
	sm.networkPolicy.allowedPorts = []int{80, 443}
	
	assert.True(t, sm.networkPolicy.allowNetwork, "Network should be enabled")
	assert.Contains(t, sm.networkPolicy.allowedHosts, "example.com")
	assert.Contains(t, sm.networkPolicy.allowedPorts, 80)
	assert.Contains(t, sm.networkPolicy.allowedPorts, 443)
}

// BenchmarkSecurityManagerOperations benchmarks security manager operations
func BenchmarkSecurityManagerOperations(b *testing.B) {
	logger, _ := logging.New(logging.DefaultConfig())
	sm := NewSecurityManager(logger)
	
	b.Run("SandboxPolicyCreation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			pluginID := fmt.Sprintf("plugin_%d", i)
			sm.getOrCreateSandboxPolicy(pluginID)
		}
	})
	
	b.Run("SecurityEventLogging", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			event := SecurityEvent{
				Timestamp:   time.Now(),
				PluginID:    "test_plugin",
				EventType:   "test_event",
				Description: "Test security event",
				Severity:    "info",
			}
			sm.logSecurityEvent(event)
		}
	})
	
	b.Run("ResourceMonitoring", func(b *testing.B) {
		pluginID := "test_plugin"
		sm.initializeResourceTracking(pluginID)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			sm.MonitorPluginResourceUsage(pluginID)
		}
	})
	
	b.Run("PathAccessControl", func(b *testing.B) {
		pluginID := "test_plugin"
		path := "/test/path"
		
		for i := 0; i < b.N; i++ {
			sm.IsPluginAllowedToAccessPath(pluginID, path)
		}
	})
}