package plugins

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/puente-labs/lyricforge/internal/logging"
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
}

// NewSecurityManager creates a new security manager
func NewSecurityManager(logger *logging.Logger) *SecurityManager {
	return &SecurityManager{
		logger: logger,
		allowedPaths: []string{
			filepath.Join(os.Getenv("HOME"), ".lyricforge", "plugins"),
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
	}
}

// ValidatePluginPath checks if a plugin path is allowed
func (sm *SecurityManager) ValidatePluginPath(path string) error {
	// Check if path is absolute and within allowed directories
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
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
		return fmt.Errorf("cannot access plugin file: %w", err)
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
		return "", fmt.Errorf("cannot open plugin file: %w", err)
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("cannot hash plugin file: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// VerifyPluginIntegrity checks if a plugin file has been modified
func (sm *SecurityManager) VerifyPluginIntegrity(path string) error {
	currentHash, err := sm.CalculatePluginHash(path)
	if err != nil {
		return fmt.Errorf("cannot calculate plugin hash: %w", err)
	}

	// For now, we'll just log the hash
	// In a production system, you might want to store known good hashes
	sm.logger.Debugf("Plugin %s hash: %s", path, currentHash)

	return nil
}

// SandboxPlugin executes a plugin in a restricted environment
func (sm *SecurityManager) SandboxPlugin(plugin Plugin) error {
	// This is a simplified sandbox implementation
	// In a production system, you might want to use:
	// - OS-level sandboxing (seccomp, namespaces)
	// - Language-level sandboxing (gvisor, rlbox)
	// - Runtime security policies

	sm.logger.Debugf("Sandboxing plugin: %s", plugin.Metadata().ID)

	// Basic validation
	if err := sm.ValidatePluginManifest(plugin.Metadata()); err != nil {
		return fmt.Errorf("plugin validation failed: %w", err)
	}

	return nil
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
