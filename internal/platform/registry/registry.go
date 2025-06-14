package registry

import (
	"fmt"
	"sync"

	"detectviz/pkg/platform/contracts"
)

// Manager implements the plugin registry management.
// zh: Manager 實作插件註冊管理功能。
type Manager struct {
	plugins   map[string]contracts.PluginFactory
	metadata  map[string]*contracts.PluginMetadata
	instances map[string]contracts.Plugin
	importers map[string]func(contracts.ImporterConfig) (contracts.Importer, error)
	exporters map[string]func(contracts.ExporterConfig) (contracts.Exporter, error)
	auth      map[string]func(contracts.AuthConfig) (contracts.Authenticator, error)
	mutex     sync.RWMutex
}

// NewManager creates a new registry manager.
// zh: NewManager 建立新的註冊管理器。
func NewManager() *Manager {
	return &Manager{
		plugins:   make(map[string]contracts.PluginFactory),
		metadata:  make(map[string]*contracts.PluginMetadata),
		instances: make(map[string]contracts.Plugin),
		importers: make(map[string]func(contracts.ImporterConfig) (contracts.Importer, error)),
		exporters: make(map[string]func(contracts.ExporterConfig) (contracts.Exporter, error)),
		auth:      make(map[string]func(contracts.AuthConfig) (contracts.Authenticator, error)),
	}
}

// RegisterPlugin registers a plugin factory with the registry.
// zh: RegisterPlugin 向註冊表註冊插件工廠。
func (m *Manager) RegisterPlugin(name string, factory contracts.PluginFactory) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s already registered", name)
	}

	m.plugins[name] = factory
	return nil
}

// GetPlugin retrieves a plugin instance by name.
// zh: GetPlugin 根據名稱取得插件實例。
func (m *Manager) GetPlugin(name string) (contracts.Plugin, error) {
	m.mutex.RLock()
	if instance, exists := m.instances[name]; exists {
		m.mutex.RUnlock()
		return instance, nil
	}
	m.mutex.RUnlock()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Double-check after acquiring write lock
	if instance, exists := m.instances[name]; exists {
		return instance, nil
	}

	factory, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	// Get plugin configuration from metadata
	var config any
	if metadata, exists := m.metadata[name]; exists {
		config = metadata.Config
	}

	// Create instance with actual config
	instance, err := factory(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin %s: %w", name, err)
	}

	// Initialize the plugin
	if err := instance.Init(config); err != nil {
		return nil, fmt.Errorf("failed to initialize plugin %s: %w", name, err)
	}

	m.instances[name] = instance
	return instance, nil
}

// ListPlugins returns a list of all registered plugin names.
// zh: ListPlugins 回傳所有已註冊插件名稱的清單。
func (m *Manager) ListPlugins() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

// GetPluginMetadata retrieves metadata for a specific plugin.
// zh: GetPluginMetadata 取得特定插件的元資料。
func (m *Manager) GetPluginMetadata(name string) (*contracts.PluginMetadata, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	metadata, exists := m.metadata[name]
	if !exists {
		return nil, fmt.Errorf("metadata for plugin %s not found", name)
	}

	return metadata, nil
}

// RegisterMetadata registers metadata for a plugin.
// zh: RegisterMetadata 為插件註冊元資料。
func (m *Manager) RegisterMetadata(name string, metadata *contracts.PluginMetadata) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.metadata[name] = metadata
	return nil
}

// RegisterImporter registers an importer factory.
// zh: RegisterImporter 註冊匯入器工廠。
func (m *Manager) RegisterImporter(name string, factory func(contracts.ImporterConfig) (contracts.Importer, error)) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.importers[name]; exists {
		return fmt.Errorf("importer %s already registered", name)
	}

	m.importers[name] = factory
	return nil
}

// GetImporter retrieves an importer by name.
// zh: GetImporter 根據名稱取得匯入器。
func (m *Manager) GetImporter(name string) (contracts.Importer, error) {
	m.mutex.RLock()
	factory, exists := m.importers[name]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("importer %s not found", name)
	}

	// TODO: Load actual config from configuration system
	return factory(contracts.ImporterConfig{Name: name})
}

// ListImporters returns a list of all registered importer names.
// zh: ListImporters 回傳所有已註冊匯入器名稱的清單。
func (m *Manager) ListImporters() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.importers))
	for name := range m.importers {
		names = append(names, name)
	}
	return names
}

// RegisterExporter registers an exporter factory.
// zh: RegisterExporter 註冊匯出器工廠。
func (m *Manager) RegisterExporter(name string, factory func(contracts.ExporterConfig) (contracts.Exporter, error)) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.exporters[name]; exists {
		return fmt.Errorf("exporter %s already registered", name)
	}

	m.exporters[name] = factory
	return nil
}

// GetExporter retrieves an exporter by name.
// zh: GetExporter 根據名稱取得匯出器。
func (m *Manager) GetExporter(name string) (contracts.Exporter, error) {
	m.mutex.RLock()
	factory, exists := m.exporters[name]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("exporter %s not found", name)
	}

	// TODO: Load actual config from configuration system
	return factory(contracts.ExporterConfig{Name: name})
}

// ListExporters returns a list of all registered exporter names.
// zh: ListExporters 回傳所有已註冊匯出器名稱的清單。
func (m *Manager) ListExporters() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	names := make([]string, 0, len(m.exporters))
	for name := range m.exporters {
		names = append(names, name)
	}
	return names
}

// RegisterAuthenticator registers an authenticator factory.
// zh: RegisterAuthenticator 註冊認證器工廠。
func (m *Manager) RegisterAuthenticator(name string, factory func(contracts.AuthConfig) (contracts.Authenticator, error)) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.auth[name]; exists {
		return fmt.Errorf("authenticator %s already registered", name)
	}

	m.auth[name] = factory
	return nil
}

// GetAuthenticator retrieves an authenticator by name.
// zh: GetAuthenticator 根據名稱取得認證器。
func (m *Manager) GetAuthenticator(name string) (contracts.Authenticator, error) {
	m.mutex.RLock()
	factory, exists := m.auth[name]
	m.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("authenticator %s not found", name)
	}

	// TODO: Load actual config from configuration system
	return factory(contracts.AuthConfig{Provider: name})
}
