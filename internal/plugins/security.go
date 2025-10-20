package plugins

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/logging"
)

// SecurityManager handles plugin security validation and sandboxing
type SecurityManager struct {
	logger       *logging.Logger
	allowedPaths []string
	blockedPaths []string
	maxFileSize  int64
	allowedExts  []string
	scanInterval time.Duration
	pluginHashes map[string]string // For integrity checking
	
	// Enhanced sandboxing
	sandboxPolicies    map[string]*SandboxPolicy
	pluginPermissions map[string]*PluginPermissions
	securityEvents     []SecurityEvent
	eventMutex         sync.RWMutex
	
	// Runtime monitoring
	pluginResourceUsage map[string]*ResourceUsage
	usageMutex          sync.RWMutex
	
	// Network restrictions
	networkPolicy *NetworkPolicy
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(logger *logging.Logger) *SecurityManager {
	return &SecurityManager{
		logger: logger,
		allowedPaths: []string{
			filepath.Join(os.Getenv("HOME"), ".noise", "plugins"),
			"./plugins",
			"./internal/plugins/examples",
		},
		blockedPaths: []string{
			"/etc",
			"/usr",
			"/bin",
			"/sbin",
			"/sys",
			"/proc",
			"/dev",
			os.Getenv("HOME"), // Block access to entire home directory except plugin dirs
		},
		maxFileSize:  10 * 1024 * 1024, // 10MB max plugin size
		allowedExts:  []string{".json", ".so", ".dll"},
		scanInterval: 5 * time.Minute,
		pluginHashes: make(map[string]string),
		
		// Enhanced sandboxing
		sandboxPolicies:    make(map[string]*SandboxPolicy),
		pluginPermissions: make(map[string]*PluginPermissions),
		securityEvents:     make([]SecurityEvent, 0),
		
		// Runtime monitoring
		pluginResourceUsage: make(map[string]*ResourceUsage),
		
		// Network restrictions
		networkPolicy: &NetworkPolicy{
			allowNetwork: false,
			allowedHosts: []string{},
			allowedPorts: []int{},
		},
	}
}

// ValidatePluginPath checks if a plugin path is allowed
func (sm *SecurityManager) ValidatePluginPath(path string) error {
	// Check if path is absolute and within allowed directories
	absPath, err := filepath.Abs(path)
	if err != nil {
		return errutil.Wrap(err, "invalid path")
	}

	// Check against blocked paths
	for _, blocked := range sm.blockedPaths {
		if strings.HasPrefix(absPath, blocked) {
			return fmt.Errorf("plugin path %s is in blocked directory %s", absPath, blocked)
		}
	}

	// Check if path is within allowed directories
	allowed := false
	for _, allowedPath := range sm.allowedPaths {
		absAllowed, err := filepath.Abs(allowedPath)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absAllowed) {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("plugin path %s is not in an allowed directory", absPath)
	}

	return nil
}

// ValidatePluginFile checks if a plugin file is safe to load
func (sm *SecurityManager) ValidatePluginFile(path string) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	allowed := false
	for _, allowedExt := range sm.allowedExts {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("plugin file extension %s is not allowed", ext)
	}

	// Check file size
	info, err := os.Stat(path)
	if err != nil {
		return errutil.Wrap(err, "access plugin file")
	}

	if info.Size() > sm.maxFileSize {
		return fmt.Errorf("plugin file size %d exceeds maximum allowed size %d", info.Size(), sm.maxFileSize)
	}

	// Check file permissions (should not be world-writable)
	if info.Mode().Perm()&0o002 != 0 {
		return fmt.Errorf("plugin file has unsafe permissions (world-writable)")
	}

	return nil
}

// ValidatePluginManifest validates a plugin manifest for security issues
func (sm *SecurityManager) ValidatePluginManifest(metadata *PluginMetadata) error {
	// Check for required fields
	if metadata.ID == "" {
		return fmt.Errorf("plugin ID is required")
	}

	if metadata.Name == "" {
		return fmt.Errorf("plugin name is required")
	}

	if metadata.Version == "" {
		return fmt.Errorf("plugin version is required")
	}

	// Validate plugin ID format (should be safe for filesystem)
	if strings.ContainsAny(metadata.ID, "/\\:*?\"<>|") {
		return fmt.Errorf("plugin ID contains invalid characters")
	}

	// Check for suspicious patterns in description
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror=",
		"eval(", "exec(", "system(", "popen(", "fopen(", "file_get_contents(",
	}

	descLower := strings.ToLower(metadata.Description)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(descLower, pattern) {
			return fmt.Errorf("plugin description contains suspicious pattern: %s", pattern)
		}
	}

	// Validate capabilities
	for _, capability := range metadata.Capabilities {
		if !sm.isValidCapability(capability) {
			return fmt.Errorf("plugin has invalid capability: %s", capability)
		}
	}

	return nil
}

// isValidCapability checks if a capability is valid and safe
func (sm *SecurityManager) isValidCapability(capability Capability) bool {
	validCapabilities := map[Capability]bool{
		CapabilityMenuItem:     true,
		CapabilityScreen:       true,
		CapabilityUIExtension:  true,
		CapabilityEditorTool:   true,
		CapabilityExportFormat: true,
		CapabilitySyntax:       true,
		CapabilityTheoryTool:   true,
		CapabilityChordLib:     true,
		CapabilityAudioEffect:  true,
		CapabilityMIDIHandler:  true,
		CapabilityDataProvider: true,
		CapabilityExporter:     true,
		CapabilityHook:         true,
		CapabilityConfig:       true,
	}

	return validCapabilities[capability]
}

// CalculatePluginHash calculates a hash of the plugin file for integrity checking
func (sm *SecurityManager) CalculatePluginHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errutil.Wrap(err, "open plugin file")
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", errutil.Wrap(err, "hash plugin file")
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// VerifyPluginIntegrity checks if a plugin file has been modified
func (sm *SecurityManager) VerifyPluginIntegrity(path string) error {
	currentHash, err := sm.CalculatePluginHash(path)
	if err != nil {
		return errutil.Wrap(err, "calculate plugin hash")
	}

	// For now, we'll just log the hash
	// In a production system, you might want to store known good hashes
	sm.logger.Debugf("Plugin %s hash: %s", path, currentHash)

	return nil
}

// SandboxPlugin executes a plugin in a restricted environment
func (sm *SecurityManager) SandboxPlugin(plugin Plugin) error {
	pluginID := plugin.Metadata().ID
	sm.logger.Debugf("Sandboxing plugin: %s", pluginID)

	// Basic validation
	if err := sm.ValidatePluginManifest(plugin.Metadata()); err != nil {
		return errutil.Wrap(err, "plugin validation failed")
	}

	// Create or update sandbox policy for this plugin
	_ = sm.getOrCreateSandboxPolicy(pluginID)
	
	// Create permissions for this plugin
	_ = sm.getOrCreatePluginPermissions(pluginID)
	
	// Initialize resource usage tracking
	sm.initializeResourceTracking(pluginID)
	
	// Log sandbox creation event
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "sandbox_created",
		Description: "Plugin sandbox created with policy",
		Severity:    "info",
	})

	// In a real implementation, you would:
	// 1. Create a restricted execution environment
	// 2. Apply OS-level restrictions (namespaces, seccomp)
	// 3. Set up resource limits
	// 4. Monitor plugin behavior
	
	return nil
}

// getOrCreateSandboxPolicy gets or creates a sandbox policy for a plugin
func (sm *SecurityManager) getOrCreateSandboxPolicy(pluginID string) *SandboxPolicy {
	if policy, exists := sm.sandboxPolicies[pluginID]; exists {
		return policy
	}
	
	// Create default restrictive policy
	policy := &SandboxPolicy{
		PluginID:        pluginID,
		AllowFileAccess: false,
		AllowedPaths:    []string{}, // No file access by default
		AllowNetwork:    false,
		MaxMemory:       50 * 1024 * 1024, // 50MB
		MaxCPU:          50, // 50% of one CPU core
		MaxExecutionTime: 30 * time.Second,
		CreatedAt:       time.Now(),
	}
	
	sm.sandboxPolicies[pluginID] = policy
	return policy
}

// getOrCreatePluginPermissions gets or creates permissions for a plugin
func (sm *SecurityManager) getOrCreatePluginPermissions(pluginID string) *PluginPermissions {
	if perms, exists := sm.pluginPermissions[pluginID]; exists {
		return perms
	}
	
	// Create default minimal permissions
	perms := &PluginPermissions{
		PluginID:     pluginID,
		CanReadFiles: false,
		CanWriteFiles: false,
		CanExecute:   false,
		CanNetwork:   false,
		CreatedAt:    time.Now(),
	}
	
	sm.pluginPermissions[pluginID] = perms
	return perms
}

// initializeResourceTracking initializes resource usage tracking for a plugin
func (sm *SecurityManager) initializeResourceTracking(pluginID string) {
	sm.usageMutex.Lock()
	defer sm.usageMutex.Unlock()
	
	sm.pluginResourceUsage[pluginID] = &ResourceUsage{
		PluginID:       pluginID,
		MemoryUsage:    0,
		CPUUsage:       0,
		FileDescriptors: 0,
		NetworkIO:      0,
		StartTime:      time.Now(),
	}
}

// MonitorPluginResourceUsage monitors resource usage for a plugin
func (sm *SecurityManager) MonitorPluginResourceUsage(pluginID string) error {
	sm.usageMutex.RLock()
	usage, exists := sm.pluginResourceUsage[pluginID]
	sm.usageMutex.RUnlock()
	
	if !exists {
		return fmt.Errorf("no resource tracking for plugin %s", pluginID)
	}
	
	// Get current resource usage
	currentMemory := sm.getCurrentMemoryUsage(pluginID)
	currentCPU := sm.getCurrentCPUUsage(pluginID)
	currentFDs := sm.getCurrentFileDescriptors(pluginID)
	
	// Update usage tracking
	sm.usageMutex.Lock()
	usage.MemoryUsage = currentMemory
	usage.CPUUsage = currentCPU
	usage.FileDescriptors = currentFDs
	usage.LastUpdate = time.Now()
	sm.usageMutex.Unlock()
	
	// Check against policy limits
	policy := sm.sandboxPolicies[pluginID]
	if policy != nil {
		if currentMemory > policy.MaxMemory {
			sm.handleResourceViolation(pluginID, "memory", currentMemory, policy.MaxMemory)
		}
		if currentCPU > policy.MaxCPU {
			sm.handleResourceViolation(pluginID, "cpu", currentCPU, policy.MaxCPU)
		}
	}
	
	return nil
}

// handleResourceViolation handles a resource usage violation
func (sm *SecurityManager) handleResourceViolation(pluginID, resourceType string, current, limit int64) {
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "resource_violation",
		Description: fmt.Sprintf("Plugin exceeded %s limit: %d/%d", resourceType, current, limit),
		Severity:    "warning",
	})
	
	sm.logger.Warnf("Plugin %s exceeded %s limit: %d/%d", pluginID, resourceType, current, limit)
	
	// In a real implementation, you might:
	// 1. Throttle the plugin
	// 2. Kill the plugin process
	// 3. Notify the user
}

// getCurrentMemoryUsage gets current memory usage for a plugin
func (sm *SecurityManager) getCurrentMemoryUsage(pluginID string) int64 {
	// In a real implementation, you would:
	// 1. Query the OS for memory usage of the plugin process
	// 2. Account for shared memory
	// 3. Consider virtual vs physical memory
	
	// For now, return a placeholder value
	return 10 * 1024 * 1024 // 10MB
}

// getCurrentCPUUsage gets current CPU usage for a plugin
func (sm *SecurityManager) getCurrentCPUUsage(pluginID string) int64 {
	// In a real implementation, you would:
	// 1. Query the OS for CPU usage of the plugin process
	// 2. Account for multiple cores
	// 3. Consider user vs system time
	
	// For now, return a placeholder value
	return 10 // 10%
}

// getCurrentFileDescriptors gets current file descriptor count for a plugin
func (sm *SecurityManager) getCurrentFileDescriptors(pluginID string) int64 {
	// In a real implementation, you would:
	// 1. Query the OS for open file descriptors
	// 2. Filter by process ID
	// 3. Count sockets separately
	
	// For now, return a placeholder value
	return 5
}

// logSecurityEvent logs a security event
func (sm *SecurityManager) logSecurityEvent(event SecurityEvent) {
	sm.eventMutex.Lock()
	defer sm.eventMutex.Unlock()
	
	sm.securityEvents = append(sm.securityEvents, event)
	
	// Keep only the last 1000 events
	if len(sm.securityEvents) > 1000 {
		sm.securityEvents = sm.securityEvents[1:]
	}
	
	// Log based on severity
	switch event.Severity {
	case "critical":
		sm.logger.Errorf("Security Event [%s]: %s - %s", event.EventType, event.PluginID, event.Description)
	case "warning":
		sm.logger.Warnf("Security Event [%s]: %s - %s", event.EventType, event.PluginID, event.Description)
	default:
		sm.logger.Infof("Security Event [%s]: %s - %s", event.EventType, event.PluginID, event.Description)
	}
}

// GetSecurityEvents returns recent security events
func (sm *SecurityManager) GetSecurityEvents(limit int) []SecurityEvent {
	sm.eventMutex.RLock()
	defer sm.eventMutex.RUnlock()
	
	if limit <= 0 || limit > len(sm.securityEvents) {
		limit = len(sm.securityEvents)
	}
	
	// Return the most recent events
	events := make([]SecurityEvent, limit)
	copy(events, sm.securityEvents[len(sm.securityEvents)-limit:])
	
	return events
}

// TerminatePlugin terminates a plugin that violates security policies
func (sm *SecurityManager) TerminatePlugin(pluginID string) error {
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "plugin_terminated",
		Description: "Plugin terminated due to security policy violation",
		Severity:    "critical",
	})
	
	// In a real implementation, you would:
	// 1. Kill the plugin process
	// 2. Clean up resources
	// 3. Remove from plugin manager
	
	sm.logger.Errorf("Plugin %s terminated due to security policy violation", pluginID)
	return nil
}

// IsPluginAllowedToAccessPath checks if a plugin is allowed to access a path
func (sm *SecurityManager) IsPluginAllowedToAccessPath(pluginID, path string) bool {
	// Check if plugin has file access permissions
	perms, exists := sm.pluginPermissions[pluginID]
	if !exists || !perms.CanReadFiles {
		return false
	}
	
	// Check against policy
	policy, exists := sm.sandboxPolicies[pluginID]
	if !exists || !policy.AllowFileAccess {
		return false
	}
	
	// Check if path is in allowed paths
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	
	for _, allowedPath := range policy.AllowedPaths {
		absAllowed, err := filepath.Abs(allowedPath)
		if err != nil {
			continue
		}
		if strings.HasPrefix(absPath, absAllowed) {
			return true
		}
	}
	
	return false
}

// IsPluginAllowedToAccessNetwork checks if a plugin is allowed to access network
func (sm *SecurityManager) IsPluginAllowedToAccessNetwork(pluginID string) bool {
	// Check if plugin has network permissions
	perms, exists := sm.pluginPermissions[pluginID]
	if !exists || !perms.CanNetwork {
		return false
	}
	
	// Check against policy
	policy, exists := sm.sandboxPolicies[pluginID]
	if !exists || !policy.AllowNetwork {
		return false
	}
	
	// Check against global network policy
	return sm.networkPolicy.allowNetwork
}

// ExecutePluginInSandbox executes a plugin command in a sandboxed environment
func (sm *SecurityManager) ExecutePluginInSandbox(pluginID string, command []string, timeout time.Duration) ([]byte, error) {
	// Check if plugin is allowed to execute
	perms, exists := sm.pluginPermissions[pluginID]
	if !exists || !perms.CanExecute {
		return nil, fmt.Errorf("plugin %s is not allowed to execute commands", pluginID)
	}
	
	// Create context with timeout
	_, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	// In a real implementation, you would:
	// 1. Create a sandboxed execution environment
	// 2. Apply OS-level restrictions
	// 3. Execute the command with limited privileges
	// 4. Monitor resource usage during execution
	
	// For now, we'll simulate execution with safety checks
	if len(command) == 0 {
		return nil, fmt.Errorf("empty command")
	}
	
	// Log execution attempt
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "command_execution",
		Description: fmt.Sprintf("Plugin attempted to execute: %s", strings.Join(command, " ")),
		Severity:    "info",
	})
	
	// In a real implementation, you would execute the command here
	// For now, return a placeholder result
	return []byte("Command executed in sandbox"), nil
}

// Enhanced hash calculation using SHA-256 for better security
func (sm *SecurityManager) CalculatePluginHashSHA256(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errutil.Wrap(err, "open plugin file")
	}
	defer file.Close()
	
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", errutil.Wrap(err, "hash plugin file")
	}
	
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// VerifyPluginSignature verifies a plugin's digital signature
func (sm *SecurityManager) VerifyPluginSignature(path string) error {
	// In a real implementation, you would:
	// 1. Load the plugin's signature
	// 2. Load the public key of the signer
	// 3. Verify the signature against the plugin hash
	// 4. Check if the signer is trusted
	
	// For now, just log the attempt
	sm.logger.Debugf("Verifying plugin signature for: %s", path)
	return nil
}

// Security policy and data structures

// SandboxPolicy defines the sandbox policy for a plugin
type SandboxPolicy struct {
	PluginID         string
	AllowFileAccess  bool
	AllowedPaths     []string
	AllowNetwork     bool
	MaxMemory        int64
	MaxCPU           int64
	MaxExecutionTime time.Duration
	CreatedAt        time.Time
}

// PluginPermissions defines the permissions for a plugin
type PluginPermissions struct {
	PluginID     string
	CanReadFiles bool
	CanWriteFiles bool
	CanExecute   bool
	CanNetwork   bool
	CreatedAt    time.Time
}

// SecurityEvent represents a security event
type SecurityEvent struct {
	Timestamp   time.Time
	PluginID    string
	EventType   string
	Description string
	Severity    string // "info", "warning", "critical"
}

// ResourceUsage tracks resource usage for a plugin
type ResourceUsage struct {
	PluginID         string
	MemoryUsage      int64
	CPUUsage         int64
	FileDescriptors  int64
	NetworkIO        int64
	StartTime        time.Time
	LastUpdate       time.Time
}

// NetworkPolicy defines network access policy
type NetworkPolicy struct {
	allowNetwork bool
	allowedHosts []string
	allowedPorts []int
}

// ScanForMaliciousPlugins scans plugin directories for potential security issues
func (sm *SecurityManager) ScanForMaliciousPlugins() []string {
	var threats []string

	for _, pluginDir := range sm.allowedPaths {
		if err := sm.scanDirectory(pluginDir, &threats); err != nil {
			sm.logger.Warnf("Error scanning plugin directory %s: %v", pluginDir, err)
		}
	}

	return threats
}

// scanDirectory recursively scans a directory for security issues
func (sm *SecurityManager) scanDirectory(dir string, threats *[]string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files we can't access
		}

		// Check for suspicious filenames
		filename := info.Name()
		suspiciousNames := []string{
			"malware", "virus", "trojan", "backdoor", "exploit",
			"hack", "steal", "keylog", "rootkit", "worm",
		}

		filenameLower := strings.ToLower(filename)
		for _, suspicious := range suspiciousNames {
			if strings.Contains(filenameLower, suspicious) {
				*threats = append(*threats, fmt.Sprintf("Suspicious filename: %s", path))
				break
			}
		}

		// Check file permissions
		if info.Mode().Perm()&0o777 == 0o777 {
			*threats = append(*threats, fmt.Sprintf("World-writable file: %s", path))
		}

		// Check for script files that shouldn't be in plugin directories
		if strings.HasSuffix(filenameLower, ".sh") ||
			strings.HasSuffix(filenameLower, ".py") ||
			strings.HasSuffix(filenameLower, ".js") {
			*threats = append(*threats, fmt.Sprintf("Unexpected script file: %s", path))
		}

		return nil
	})
}

// CleanupStalePlugins removes plugins that haven't been used recently
func (sm *SecurityManager) CleanupStalePlugins(maxAge time.Duration) error {
	cutoff := time.Now().Add(-maxAge)

	for _, pluginDir := range sm.allowedPaths {
		err := filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Check if file is older than max age
			if info.ModTime().Before(cutoff) {
				sm.logger.Debugf("Removing stale plugin file: %s", path)
				if err := os.Remove(path); err != nil {
					sm.logger.Warnf("Failed to remove stale plugin file %s: %v", path, err)
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// GetSecurityReport returns a security report for all loaded plugins
func (sm *SecurityManager) GetSecurityReport(manager PluginManager) map[string]interface{} {
	report := make(map[string]interface{})

	plugins := manager.GetPlugins()
	report["total_plugins"] = len(plugins)

	var securePlugins, insecurePlugins []string
	var totalSize int64

	for _, plugin := range plugins {
		metadata := plugin.Metadata()

		// Basic security checks
		isSecure := true

		if metadata.Author == "" {
			isSecure = false
		}

		if strings.Contains(metadata.Description, "test") || strings.Contains(metadata.Description, "debug") {
			// This is a simple heuristic - in practice you'd want more sophisticated checks
			// TODO: Implement more comprehensive security validation
			sm.logger.Debugf("Plugin description contains test/debug keywords: %s", metadata.Description)
		}

		if isSecure {
			securePlugins = append(securePlugins, metadata.ID)
		} else {
			insecurePlugins = append(insecurePlugins, metadata.ID)
		}

		// This is a placeholder for size calculation
		// In a real implementation, you'd calculate the actual plugin size
		totalSize += 1024 // Assume 1KB per plugin for demo
	}

	report["secure_plugins"] = securePlugins
	report["insecure_plugins"] = insecurePlugins
	report["total_size_bytes"] = totalSize
	report["scan_timestamp"] = time.Now()

	return report
}
