package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/config"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/logging"
)

// DefaultManager is the default implementation of PluginManager
type DefaultManager struct {
	// Configuration
	config *config.Config

	// Loaded plugins
	plugins map[string]Plugin
	mutex   sync.RWMutex

	// Plugin capabilities
	menuItems     []*MenuItem
	screens       map[string]Screen
	editorTools   []*EditorTool
	exportFormats []*ExportFormat

	// Hook system
	hooks     map[string][]HookHandler
	hookMutex sync.RWMutex

	// Plugin directories
	pluginDirs []string

	// Logger
	logger *logging.Logger

	// Security manager
	security *SecurityManager
}

// NewManager creates a new plugin manager
func NewManager(cfg *config.Config, logger *logging.Logger) *DefaultManager {
	return &DefaultManager{
		config:        cfg,
		plugins:       make(map[string]Plugin),
		menuItems:     make([]*MenuItem, 0),
		screens:       make(map[string]Screen),
		editorTools:   make([]*EditorTool, 0),
		exportFormats: make([]*ExportFormat, 0),
		hooks:         make(map[string][]HookHandler),
		pluginDirs:    getDefaultPluginDirs(cfg),
		logger:        logger,
		security:      NewSecurityManager(logger),
	}
}

// LoadPlugins discovers and loads all available plugins
func (m *DefaultManager) LoadPlugins() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, dir := range m.pluginDirs {
		if err := m.loadPluginsFromDir(dir); err != nil {
			m.logger.Warnf("Failed to load plugins from %s: %v", dir, err)
		}
	}

	// Initialize all loaded plugins
	for id, plugin := range m.plugins {
		if err := m.initializePlugin(plugin); err != nil {
			m.logger.Errorf("Failed to initialize plugin %s: %v", id, err)
			delete(m.plugins, id)
			continue
		}

		// Collect plugin capabilities
		m.collectPluginCapabilities(plugin)
	}

	m.logger.Infof("Loaded %d plugins", len(m.plugins))
	return nil
}

// loadPluginsFromDir loads plugins from a specific directory
func (m *DefaultManager) loadPluginsFromDir(dir string) error {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil // Directory doesn't exist, skip
	}

	// Validate directory path for security
	if err := m.security.ValidatePluginPath(dir); err != nil {
		m.logger.Warnf("Skipping plugin directory %s: %v", dir, err)
		return nil
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Look for .so files (compiled plugins) or .json files (plugin manifests)
		if filepath.Ext(path) == ".so" {
			return m.loadCompiledPlugin(path)
		} else if filepath.Ext(path) == ".json" {
			return m.loadManifestPlugin(path)
		}

		return nil
	})
}

// loadCompiledPlugin loads a compiled Go plugin
func (m *DefaultManager) loadCompiledPlugin(path string) error {
	// For now, we'll implement a simple file-based plugin system
	// In a production system, you might want to use Yaegi or compile plugins

	// Check if there's a corresponding .json manifest file
	manifestPath := path[:len(path)-3] + ".json"
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil // No manifest, skip
	}

	return m.loadManifestPlugin(manifestPath)
}

// loadManifestPlugin loads a plugin from its manifest file
func (m *DefaultManager) loadManifestPlugin(manifestPath string) error {
	// Validate plugin file for security
	if err := m.security.ValidatePluginFile(manifestPath); err != nil {
		return errutil.Wrap(err, "plugin security validation failed")
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return errutil.Wrap(err, "read manifest")
	}

	var metadata PluginMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return errutil.Wrap(err, "parse manifest")
	}

	// Validate plugin manifest for security
	if err := m.security.ValidatePluginManifest(&metadata); err != nil {
		return errutil.Wrap(err, "plugin manifest validation failed")
	}

	// Sandbox the plugin
	stubPlugin := &StubPlugin{
		metadata: &metadata,
		enabled:  metadata.Enabled,
	}

	if err := m.security.SandboxPlugin(stubPlugin); err != nil {
		return errutil.Wrap(err, "plugin sandboxing failed")
	}

	m.plugins[metadata.ID] = stubPlugin
	m.logger.Debugf("Loaded plugin: %s v%s", metadata.Name, metadata.Version)

	return nil
}

// initializePlugin initializes a plugin with context
func (m *DefaultManager) initializePlugin(p Plugin) error {
	ctx := &PluginContext{
		Context:  context.Background(),
		Config:   m.config,
		Metadata: p.Metadata(),
	}

	return p.Initialize(ctx)
}

// collectPluginCapabilities collects capabilities from a plugin
func (m *DefaultManager) collectPluginCapabilities(p Plugin) {
	metadata := p.Metadata()

	// This is where we would collect actual capabilities
	// For now, we'll just log what capabilities are available
	for _, capability := range metadata.Capabilities {
		m.logger.Debugf("Plugin %s provides capability: %s", metadata.ID, capability)
	}
}

// GetPlugin returns a plugin by ID
func (m *DefaultManager) GetPlugin(id string) (Plugin, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	p, exists := m.plugins[id]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", id)
	}

	return p, nil
}

// GetPlugins returns all loaded plugins
func (m *DefaultManager) GetPlugins() map[string]Plugin {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]Plugin)
	for id, p := range m.plugins {
		result[id] = p
	}

	return result
}

// GetPluginsByCapability returns plugins that have a specific capability
func (m *DefaultManager) GetPluginsByCapability(capability Capability) []Plugin {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var result []Plugin
	for _, p := range m.plugins {
		for _, cap := range p.Metadata().Capabilities {
			if cap == capability {
				result = append(result, p)
				break
			}
		}
	}

	return result
}

// EnablePlugin enables a plugin
func (m *DefaultManager) EnablePlugin(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	p, exists := m.plugins[id]
	if !exists {
		return fmt.Errorf("plugin %s not found", id)
	}

	if err := p.Enable(); err != nil {
		return errutil.Wrapf(err, "enable plugin %s", id)
	}

	m.logger.Infof("Enabled plugin: %s", id)
	return nil
}

// DisablePlugin disables a plugin
func (m *DefaultManager) DisablePlugin(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	p, exists := m.plugins[id]
	if !exists {
		return fmt.Errorf("plugin %s not found", id)
	}

	if err := p.Disable(); err != nil {
		return errutil.Wrapf(err, "disable plugin %s", id)
	}

	m.logger.Infof("Disabled plugin: %s", id)
	return nil
}

// UnloadPlugin unloads a plugin completely
func (m *DefaultManager) UnloadPlugin(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	p, exists := m.plugins[id]
	if !exists {
		return fmt.Errorf("plugin %s not found", id)
	}

	// Disable plugin first
	if p.IsEnabled() {
		if err := p.Disable(); err != nil {
			m.logger.Warnf("Failed to disable plugin %s before unload: %v", id, err)
		}
	}

	// Cleanup plugin
	if err := p.Cleanup(); err != nil {
		m.logger.Warnf("Failed to cleanup plugin %s: %v", id, err)
	}

	delete(m.plugins, id)
	m.logger.Infof("Unloaded plugin: %s", id)
	return nil
}

// RegisterHook registers a hook handler for a plugin
func (m *DefaultManager) RegisterHook(pluginID string, handler HookHandler) error {
	m.hookMutex.Lock()
	defer m.hookMutex.Unlock()

	m.hooks[pluginID] = append(m.hooks[pluginID], handler)
	return nil
}

// CallHook calls all registered hook handlers for a hook point
func (m *DefaultManager) CallHook(point HookPoint, data HookData) error {
	m.hookMutex.RLock()
	defer m.hookMutex.RUnlock()

	for pluginID, handlers := range m.hooks {
		for _, handler := range handlers {
			if err := handler(point, data); err != nil {
				m.logger.Errorf("Hook handler for plugin %s failed: %v", pluginID, err)
				// Continue with other handlers even if one fails
			}
		}
	}

	return nil
}

// GetMenuItems returns all menu items provided by plugins
func (m *DefaultManager) GetMenuItems() []*MenuItem {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]*MenuItem, len(m.menuItems))
	copy(result, m.menuItems)
	return result
}

// GetScreens returns all screens provided by plugins
func (m *DefaultManager) GetScreens() map[string]Screen {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make(map[string]Screen)
	for id, screen := range m.screens {
		result[id] = screen
	}

	return result
}

// GetEditorTools returns all editor tools provided by plugins
func (m *DefaultManager) GetEditorTools() []*EditorTool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]*EditorTool, len(m.editorTools))
	copy(result, m.editorTools)
	return result
}

// GetExportFormats returns all export formats provided by plugins
func (m *DefaultManager) GetExportFormats() []*ExportFormat {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	result := make([]*ExportFormat, len(m.exportFormats))
	copy(result, m.exportFormats)
	return result
}

// getDefaultPluginDirs returns the default plugin directories
func getDefaultPluginDirs(cfg *config.Config) []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return []string{
		filepath.Join(homeDir, ".noise", "plugins"),
		"./plugins",
		"./internal/plugins/examples",
	}
}

// StubPlugin is a placeholder plugin implementation for testing
type StubPlugin struct {
	mu       sync.Mutex // guards mutable state for concurrent Initialize/Enable/Disable
	metadata *PluginMetadata
	enabled  bool
}

func (p *StubPlugin) Metadata() *PluginMetadata {
	return p.metadata
}

func (p *StubPlugin) Initialize(ctx *PluginContext) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.metadata.LoadTime = time.Now()
	return nil
}

func (p *StubPlugin) Cleanup() error {
	return nil
}

func (p *StubPlugin) Enable() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = true
	return nil
}

func (p *StubPlugin) Disable() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.enabled = false
	return nil
}

func (p *StubPlugin) IsEnabled() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.enabled
}
