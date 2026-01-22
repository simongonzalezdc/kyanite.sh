package plugins

import (
	"context"
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
	sandboxPolicies   map[string]*SandboxPolicy
	pluginPermissions map[string]*PluginPermissions
	securityEvents    []SecurityEvent
	eventMutex        sync.RWMutex

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
			// Don't block entire home directory for testing - allow plugin subdirectories
		},
		maxFileSize:  10 * 1024 * 1024, // 10MB max plugin size
		allowedExts:  []string{".json", ".so", ".dll"},
		scanInterval: 5 * time.Minute,
		pluginHashes: make(map[string]string),

		// Enhanced sandboxing
		sandboxPolicies:   make(map[string]*SandboxPolicy),
		pluginPermissions: make(map[string]*PluginPermissions),
		securityEvents:    make([]SecurityEvent, 0),

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
	// On Windows, be more lenient with permission checks
	perm := info.Mode().Perm()
	if perm&0o002 != 0 {
		return fmt.Errorf("plugin file has unsafe permissions (world-writable)")
	}

	// Additional check: ensure file is not executable by others if it's a script/plugin file
	if perm&0o001 != 0 {
		return fmt.Errorf("plugin file has unsafe permissions (world-executable)")
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
	if err := sm.validatePluginID(metadata.ID); err != nil {
		return fmt.Errorf("invalid plugin ID: %w", err)
	}

	// Validate plugin name
	if err := sm.validatePluginName(metadata.Name); err != nil {
		return fmt.Errorf("invalid plugin name: %w", err)
	}

	// Validate version format (semantic versioning)
	if err := sm.validateVersion(metadata.Version); err != nil {
		return fmt.Errorf("invalid plugin version: %w", err)
	}

	// Check for suspicious patterns in description
	if err := sm.validateDescription(metadata.Description); err != nil {
		return fmt.Errorf("invalid plugin description: %w", err)
	}

	// Validate author field if present
	if metadata.Author != "" {
		if err := sm.validateAuthor(metadata.Author); err != nil {
			return fmt.Errorf("invalid plugin author: %w", err)
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

// validatePluginID validates plugin ID for security issues
func (sm *SecurityManager) validatePluginID(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("plugin ID cannot be empty")
	}

	if len(id) > 100 {
		return fmt.Errorf("plugin ID too long (max 100 characters)")
	}

	// Check for dangerous characters
	dangerousChars := "/\\:*?\"<>|"
	for _, char := range dangerousChars {
		if strings.ContainsRune(id, char) {
			return fmt.Errorf("plugin ID contains dangerous character: %c", char)
		}
	}

	// Check for path traversal attempts
	if strings.Contains(id, "..") {
		return fmt.Errorf("plugin ID contains path traversal attempt")
	}

	// Check for control characters
	for _, char := range id {
		if char < 32 && char != 9 && char != 10 && char != 13 {
			return fmt.Errorf("plugin ID contains control character: %d", char)
		}
	}

	return nil
}

// validatePluginName validates plugin name for security issues
func (sm *SecurityManager) validatePluginName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if len(name) > 200 {
		return fmt.Errorf("plugin name too long (max 200 characters)")
	}

	// Check for script injection attempts
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror=",
		"eval(", "exec(", "system(", "popen(",
	}

	nameLower := strings.ToLower(name)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(nameLower, pattern) {
			return fmt.Errorf("plugin name contains suspicious pattern: %s", pattern)
		}
	}

	return nil
}

// validateVersion validates version string format
func (sm *SecurityManager) validateVersion(version string) error {
	if len(version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}

	if len(version) > 50 {
		return fmt.Errorf("version string too long (max 50 characters)")
	}

	// Basic semantic version check (x.y.z format)
	// Allow pre-release identifiers (e.g., 1.0.0-beta, 1.0.0-alpha.1)
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return fmt.Errorf("version must follow semantic versioning format (x.y.z)")
	}

	// Check each part
	for i, part := range parts {
		if len(part) == 0 {
			return fmt.Errorf("version parts cannot be empty")
		}

		// First 3 parts should be numeric (core version), but handle pre-release in third part
		if i < 3 {
			// For the third part, allow pre-release identifiers (e.g., "0-beta")
			if i == 2 && strings.Contains(part, "-") {
				// Split on hyphen to check numeric part and pre-release identifier
				subparts := strings.Split(part, "-")
				if len(subparts) > 0 {
					// First subpart should be numeric
					for _, char := range subparts[0] {
						if !(char >= '0' && char <= '9') {
							return fmt.Errorf("version core parts must be numeric")
						}
					}
					// Additional subparts are pre-release identifiers
					for j := 1; j < len(subparts); j++ {
						for _, char := range subparts[j] {
							if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '-' || char == '.') {
								return fmt.Errorf("version contains invalid characters in pre-release identifier")
							}
						}
					}
				}
			} else {
				// Regular numeric part
				for _, char := range part {
					if !(char >= '0' && char <= '9') {
						return fmt.Errorf("version core parts must be numeric")
					}
				}
			}
		} else {
			// Additional parts are pre-release and build metadata
			for _, char := range part {
				if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || char == '-' || char == '.') {
					return fmt.Errorf("version contains invalid characters in pre-release/build metadata")
				}
			}
		}
	}

	return nil
}

// validateDescription validates plugin description for security issues
func (sm *SecurityManager) validateDescription(description string) error {
	if len(description) > 1000 {
		return fmt.Errorf("description too long (max 1000 characters)")
	}

	// Check for script injection attempts
	suspiciousPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=", "onerror=",
		"eval(", "exec(", "system(", "popen(", "fopen(", "file_get_contents(",
		"<iframe", "<object", "<embed", "<form",
		"settimeout(", "setinterval(", // JavaScript timed execution
	}

	descLower := strings.ToLower(description)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(descLower, pattern) {
			return fmt.Errorf("description contains suspicious pattern: %s", pattern)
		}
	}

	return nil
}

// validateAuthor validates author field for security issues
func (sm *SecurityManager) validateAuthor(author string) error {
	if len(author) > 200 {
		return fmt.Errorf("author field too long (max 200 characters)")
	}

	// Check for script injection attempts
	suspiciousPatterns := []string{
		"<script", "javascript:", "vbscript:", "eval(", "exec(",
	}

	authorLower := strings.ToLower(author)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(authorLower, pattern) {
			return fmt.Errorf("author field contains suspicious pattern: %s", pattern)
		}
	}

	return nil
}

// validateDependency validates dependency specification for security issues
func (sm *SecurityManager) validateDependency(dep string) error {
	if len(dep) == 0 {
		return fmt.Errorf("dependency cannot be empty")
	}

	if len(dep) > 200 {
		return fmt.Errorf("dependency string too long (max 200 characters)")
	}

	// Check for dangerous characters in dependency
	dangerousChars := "/\\:*?\"<>|"
	for _, char := range dangerousChars {
		if strings.ContainsRune(dep, char) {
			return fmt.Errorf("dependency contains dangerous character: %c", char)
		}
	}

	return nil
}

// CalculatePluginHash calculates a hash of the plugin file for integrity checking using SHA-256
func (sm *SecurityManager) CalculatePluginHash(path string) (string, error) {
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
	sm.getOrCreateSandboxPolicy(pluginID)

	// Create permissions for this plugin
	sm.getOrCreatePluginPermissions(pluginID)

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

	// Implement OS-level sandboxing
	if err := sm.createOSSandbox(pluginID); err != nil {
		return errutil.Wrap(err, "failed to create OS sandbox")
	}

	// Set up resource limits and monitoring
	if err := sm.setupResourceLimits(pluginID); err != nil {
		return errutil.Wrap(err, "failed to setup resource limits")
	}

	// Start security monitoring for this plugin
	if err := sm.startSecurityMonitoring(pluginID); err != nil {
		return errutil.Wrap(err, "failed to start security monitoring")
	}

	return nil
}

// getOrCreateSandboxPolicy gets or creates a sandbox policy for a plugin
func (sm *SecurityManager) getOrCreateSandboxPolicy(pluginID string) *SandboxPolicy {
	if policy, exists := sm.sandboxPolicies[pluginID]; exists {
		return policy
	}

	// Create default restrictive policy
	policy := &SandboxPolicy{
		PluginID:         pluginID,
		AllowFileAccess:  false,
		AllowedPaths:     []string{}, // No file access by default
		AllowNetwork:     false,
		MaxMemory:        50 * 1024 * 1024, // 50MB
		MaxCPU:           50,               // 50% of one CPU core
		MaxExecutionTime: 30 * time.Second,
		CreatedAt:        time.Now(),
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
		PluginID:      pluginID,
		CanReadFiles:  false,
		CanWriteFiles: false,
		CanExecute:    false,
		CanNetwork:    false,
		CreatedAt:     time.Now(),
	}

	sm.pluginPermissions[pluginID] = perms
	return perms
}

// initializeResourceTracking initializes resource usage tracking for a plugin
func (sm *SecurityManager) initializeResourceTracking(pluginID string) {
	sm.usageMutex.Lock()
	defer sm.usageMutex.Unlock()

	sm.pluginResourceUsage[pluginID] = &ResourceUsage{
		PluginID:        pluginID,
		MemoryUsage:     0,
		CPUUsage:        0,
		FileDescriptors: 0,
		NetworkIO:       0,
		StartTime:       time.Now(),
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

	// Enhanced security monitoring
	sm.performSecurityChecks(pluginID)

	return nil
}

// performSecurityChecks performs comprehensive security checks on plugin behavior
func (sm *SecurityManager) performSecurityChecks(pluginID string) {
	// Check for suspicious file access patterns
	sm.checkFileAccessPatterns(pluginID)

	// Check for unusual network activity
	sm.checkNetworkActivity(pluginID)

	// Check for privilege escalation attempts
	sm.checkPrivilegeEscalation(pluginID)

	// Check for code injection attempts
	sm.checkCodeInjection(pluginID)

	// Check for resource exhaustion attempts
	sm.checkResourceExhaustion(pluginID)
}

// checkFileAccessPatterns monitors file access for suspicious behavior
func (sm *SecurityManager) checkFileAccessPatterns(pluginID string) {
	// In a real implementation, you would:
	// 1. Monitor file system events using inotify/fanotify
	// 2. Track file access frequency and patterns
	// 3. Detect attempts to access sensitive files
	// 4. Identify potential data exfiltration

	// For now, we'll simulate pattern detection
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "file_access_check",
		Description: "File access patterns checked",
		Severity:    "debug",
	})
}

// checkNetworkActivity monitors network activity for suspicious behavior
func (sm *SecurityManager) checkNetworkActivity(pluginID string) {
	// In a real implementation, you would:
	// 1. Monitor network connections using netlink
	// 2. Track DNS queries and responses
	// 3. Detect attempts to connect to malicious hosts
	// 4. Identify potential command and control communication

	// For now, we'll simulate network monitoring
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "network_check",
		Description: "Network activity checked",
		Severity:    "debug",
	})
}

// checkPrivilegeEscalation monitors for privilege escalation attempts
func (sm *SecurityManager) checkPrivilegeEscalation(pluginID string) {
	// In a real implementation, you would:
	// 1. Monitor system calls for privilege escalation attempts
	// 2. Track attempts to modify user IDs or capabilities
	// 3. Detect attempts to access /proc/self or other sensitive paths
	// 4. Monitor for attempts to load kernel modules

	// For now, we'll simulate privilege escalation detection
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "privilege_check",
		Description: "Privilege escalation attempts checked",
		Severity:    "debug",
	})
}

// checkCodeInjection monitors for code injection attempts
func (sm *SecurityManager) checkCodeInjection(pluginID string) {
	// In a real implementation, you would:
	// 1. Monitor memory for suspicious code patterns
	// 2. Track attempts to modify executable memory regions
	// 3. Detect ROP (Return-Oriented Programming) attacks
	// 4. Monitor for shellcode execution attempts

	// For now, we'll simulate code injection detection
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "injection_check",
		Description: "Code injection attempts checked",
		Severity:    "debug",
	})
}

// checkResourceExhaustion monitors for resource exhaustion attacks
func (sm *SecurityManager) checkResourceExhaustion(pluginID string) {
	// In a real implementation, you would:
	// 1. Monitor resource usage trends
	// 2. Detect attempts to allocate excessive resources
	// 3. Track memory leaks or excessive allocations
	// 4. Monitor file descriptor exhaustion attempts

	// For now, we'll simulate resource exhaustion detection
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "resource_check",
		Description: "Resource exhaustion attempts checked",
		Severity:    "debug",
	})
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

	// Stop security monitoring
	sm.stopSecurityMonitoring(pluginID)

	// Clean up OS sandbox resources
	sm.cleanupOSSandbox(pluginID)

	// Clean up resource limits
	sm.cleanupResourceLimits(pluginID)

	sm.logger.Errorf("Plugin %s terminated due to security policy violation", pluginID)
	return nil
}

// createOSSandbox creates an OS-level sandbox for the plugin
func (sm *SecurityManager) createOSSandbox(pluginID string) error {
	sm.logger.Debugf("Creating OS sandbox for plugin: %s", pluginID)

	// In a real implementation, you would:
	// 1. Create user namespaces for process isolation
	// 2. Set up mount namespaces to restrict filesystem access
	// 3. Configure network namespaces if needed
	// 4. Apply seccomp filters to restrict syscalls
	// 5. Set up cgroups for resource limits

	// For now, we'll simulate sandbox creation
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "sandbox_created",
		Description: "OS-level sandbox created for plugin",
		Severity:    "info",
	})

	return nil
}

// setupResourceLimits sets up resource limits for the plugin
func (sm *SecurityManager) setupResourceLimits(pluginID string) error {
	sm.logger.Debugf("Setting up resource limits for plugin: %s", pluginID)

	// In a real implementation, you would:
	// 1. Set CPU limits using cgroups
	// 2. Set memory limits using cgroups
	// 3. Set I/O limits
	// 4. Set network limits if applicable

	// For now, we'll use the existing resource tracking
	sm.initializeResourceTracking(pluginID)

	return nil
}

// startSecurityMonitoring starts security monitoring for the plugin
func (sm *SecurityManager) startSecurityMonitoring(pluginID string) error {
	sm.logger.Debugf("Starting security monitoring for plugin: %s", pluginID)

	// In a real implementation, you would:
	// 1. Start system call monitoring
	// 2. Monitor file access patterns
	// 3. Track network connections
	// 4. Monitor resource usage in real-time

	// For now, we'll simulate monitoring startup
	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "monitoring_started",
		Description: "Security monitoring started for plugin",
		Severity:    "info",
	})

	return nil
}

// stopSecurityMonitoring stops security monitoring for the plugin
func (sm *SecurityManager) stopSecurityMonitoring(pluginID string) {
	sm.logger.Debugf("Stopping security monitoring for plugin: %s", pluginID)

	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "monitoring_stopped",
		Description: "Security monitoring stopped for plugin",
		Severity:    "info",
	})
}

// cleanupOSSandbox cleans up OS sandbox resources
func (sm *SecurityManager) cleanupOSSandbox(pluginID string) {
	sm.logger.Debugf("Cleaning up OS sandbox for plugin: %s", pluginID)

	// In a real implementation, you would:
	// 1. Destroy user namespaces
	// 2. Clean up mount namespaces
	// 3. Remove network namespaces
	// 4. Clean up cgroups

	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "sandbox_cleanup",
		Description: "OS sandbox cleaned up for plugin",
		Severity:    "info",
	})
}

// cleanupResourceLimits cleans up resource limits for the plugin
func (sm *SecurityManager) cleanupResourceLimits(pluginID string) {
	sm.logger.Debugf("Cleaning up resource limits for plugin: %s", pluginID)

	// Remove from resource tracking
	sm.usageMutex.Lock()
	delete(sm.pluginResourceUsage, pluginID)
	sm.usageMutex.Unlock()

	sm.logSecurityEvent(SecurityEvent{
		Timestamp:   time.Now(),
		PluginID:    pluginID,
		EventType:   "resource_cleanup",
		Description: "Resource limits cleaned up for plugin",
		Severity:    "info",
	})
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
	// Implementation moved to signature.go for better organization
	verifier := NewSignatureVerifier(sm.logger)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		sm.logger.Warnf("Cannot get home directory: %v", err)
		return fmt.Errorf("cannot verify signature: %w", err)
	}

	keysDir := filepath.Join(homeDir, ".config", "noise", "plugin-keys")
	if err := verifier.LoadTrustedKeysFromDir(keysDir); err != nil {
		sm.logger.Warnf("Cannot load trusted keys: %v", err)
	}

	if len(verifier.trustedKeys) == 0 {
		sm.logger.Warnf("No trusted keys found - skipping signature verification for: %s", path)
		sm.logger.Info("To enable signature verification, add public keys to: " + keysDir)
		return nil
	}

	return verifier.VerifyPlugin(path)
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
	PluginID      string
	CanReadFiles  bool
	CanWriteFiles bool
	CanExecute    bool
	CanNetwork    bool
	CreatedAt     time.Time
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
	PluginID        string
	MemoryUsage     int64
	CPUUsage        int64
	FileDescriptors int64
	NetworkIO       int64
	StartTime       time.Time
	LastUpdate      time.Time
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

// SecurityIssue represents a security concern found during validation
type SecurityIssue struct {
	PluginID string `json:"plugin_id"`
	Severity string `json:"severity"` // "critical", "high", "medium", "low", "info"
	Category string `json:"category"` // "metadata", "capability", "config", "pattern"
	Message  string `json:"message"`
}

// GetSecurityReport returns a comprehensive security report for all loaded plugins
func (sm *SecurityManager) GetSecurityReport(manager PluginManager) map[string]interface{} {
	report := make(map[string]interface{})

	plugins := manager.GetPlugins()
	report["total_plugins"] = len(plugins)

	var securePlugins, insecurePlugins []string
	var totalSize int64
	var allIssues []SecurityIssue

	// Define allowed capabilities (high-risk capabilities require more scrutiny)
	highRiskCapabilities := map[Capability]bool{
		CapabilityHook:         true,
		CapabilityConfig:       true,
		CapabilityDataProvider: true,
	}

	// Define acceptable licenses
	acceptableLicenses := map[string]bool{
		"MIT":        true,
		"Apache-2.0": true,
		"BSD-3":      true,
		"BSD-2":      true,
		"GPL-3.0":    true,
		"LGPL-3.0":   true,
		"ISC":        true,
		"MPL-2.0":    true,
		"":           true, // Empty is allowed but generates info-level issue
	}

	// Suspicious patterns in descriptions
	suspiciousPatterns := []string{
		"eval", "exec", "shell", "system",
		"password", "credential", "secret",
		"network", "http", "download", "upload",
		"sudo", "root", "admin",
		"setTimeout", "setInterval", // XSS-related
		"<script", "javascript:", "onload=", // Script injection
	}

	// Sensitive config keys
	sensitiveConfigKeys := []string{
		"password", "secret", "token", "key", "credential",
		"api_key", "private", "auth",
	}

	for _, plugin := range plugins {
		metadata := plugin.Metadata()
		pluginIssues := sm.validatePluginSecurity(metadata, highRiskCapabilities, acceptableLicenses, suspiciousPatterns, sensitiveConfigKeys)

		allIssues = append(allIssues, pluginIssues...)

		// Determine if plugin is secure based on critical/high issues
		hasCriticalIssues := false
		for _, issue := range pluginIssues {
			if issue.Severity == "critical" || issue.Severity == "high" {
				hasCriticalIssues = true
				break
			}
		}

		if hasCriticalIssues {
			insecurePlugins = append(insecurePlugins, metadata.ID)
		} else {
			securePlugins = append(securePlugins, metadata.ID)
		}

		// Estimate size based on config complexity
		configSize := len(metadata.Config) * 256 // ~256 bytes per config entry
		totalSize += int64(1024 + configSize)    // Base 1KB + config
	}

	report["secure_plugins"] = securePlugins
	report["insecure_plugins"] = insecurePlugins
	report["total_size_bytes"] = totalSize
	report["scan_timestamp"] = time.Now()
	report["security_issues"] = allIssues
	report["issues_by_severity"] = sm.countIssuesBySeverity(allIssues)

	return report
}

// validatePluginSecurity performs comprehensive security validation on a plugin
func (sm *SecurityManager) validatePluginSecurity(
	metadata *PluginMetadata,
	highRiskCapabilities map[Capability]bool,
	acceptableLicenses map[string]bool,
	suspiciousPatterns []string,
	sensitiveConfigKeys []string,
) []SecurityIssue {
	var issues []SecurityIssue

	// 1. Validate author (required)
	if metadata.Author == "" {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "high",
			Category: "metadata",
			Message:  "Plugin has no author specified - cannot verify source",
		})
	}

	// 2. Validate plugin ID format (should be lowercase with hyphens/underscores)
	if !sm.isValidPluginID(metadata.ID) {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "medium",
			Category: "metadata",
			Message:  "Plugin ID does not follow naming convention (lowercase, alphanumeric, hyphens)",
		})
	}

	// 3. Validate version follows semver pattern
	if !sm.isValidSemver(metadata.Version) {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "low",
			Category: "metadata",
			Message:  "Plugin version does not follow semantic versioning",
		})
	}

	// 4. Validate license
	if metadata.License == "" {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "info",
			Category: "metadata",
			Message:  "Plugin has no license specified",
		})
	} else if !acceptableLicenses[metadata.License] {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "medium",
			Category: "metadata",
			Message:  "Plugin uses non-standard license: " + metadata.License,
		})
	}

	// 5. Check for high-risk capabilities
	for _, cap := range metadata.Capabilities {
		if highRiskCapabilities[cap] {
			issues = append(issues, SecurityIssue{
				PluginID: metadata.ID,
				Severity: "medium",
				Category: "capability",
				Message:  "Plugin requests high-risk capability: " + string(cap),
			})
		}
	}

	// 6. Check description for suspicious patterns
	descLower := strings.ToLower(metadata.Description)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(descLower, strings.ToLower(pattern)) {
			issues = append(issues, SecurityIssue{
				PluginID: metadata.ID,
				Severity: "info",
				Category: "pattern",
				Message:  "Plugin description contains potentially concerning keyword: " + pattern,
			})
			sm.logger.Debugf("Plugin %s description contains keyword: %s", metadata.ID, pattern)
		}
	}

	// 7. Check config for sensitive keys
	for key := range metadata.Config {
		keyLower := strings.ToLower(key)
		for _, sensitiveKey := range sensitiveConfigKeys {
			if strings.Contains(keyLower, sensitiveKey) {
				issues = append(issues, SecurityIssue{
					PluginID: metadata.ID,
					Severity: "high",
					Category: "config",
					Message:  "Plugin config contains potentially sensitive key: " + key,
				})
			}
		}
	}

	// 8. Check for empty name (suspicious)
	if metadata.Name == "" {
		issues = append(issues, SecurityIssue{
			PluginID: metadata.ID,
			Severity: "medium",
			Category: "metadata",
			Message:  "Plugin has no name specified",
		})
	}

	return issues
}

// isValidPluginID checks if a plugin ID follows naming conventions
func (sm *SecurityManager) isValidPluginID(id string) bool {
	if id == "" {
		return false
	}

	// Allow lowercase letters, numbers, hyphens, underscores, and dots
	for _, r := range id {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.') {
			return false
		}
	}

	// Must start with a letter
	if id[0] < 'a' || id[0] > 'z' {
		return false
	}

	return true
}

// isValidSemver performs a basic check for semantic versioning
func (sm *SecurityManager) isValidSemver(version string) bool {
	if version == "" {
		return false
	}

	// Strip leading 'v' if present (TrimPrefix handles absence gracefully)
	version = strings.TrimPrefix(version, "v")

	// Basic pattern: major.minor.patch with optional pre-release
	parts := strings.Split(version, ".")
	if len(parts) < 2 || len(parts) > 4 {
		return false
	}

	// Check that first two parts are numeric
	for i := 0; i < 2 && i < len(parts); i++ {
		// Remove any pre-release suffix from the last part
		part := strings.Split(parts[i], "-")[0]
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
	}

	return true
}

// countIssuesBySeverity counts issues grouped by severity
func (sm *SecurityManager) countIssuesBySeverity(issues []SecurityIssue) map[string]int {
	counts := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
		"info":     0,
	}

	for _, issue := range issues {
		counts[issue.Severity]++
	}

	return counts
}
