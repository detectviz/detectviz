package loader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"detectviz/pkg/platform/contracts"
)

// CompositionConfig represents the structure of composition.yaml
// zh: CompositionConfig 代表 composition.yaml 的結構
type CompositionConfig struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   ConfigMetadata `yaml:"metadata"`
	Spec       ConfigSpec     `yaml:"spec"`
}

// ConfigMetadata contains metadata about the composition
// zh: ConfigMetadata 包含組合的元資料
type ConfigMetadata struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
}

// ConfigSpec contains the specification for the composition
// zh: ConfigSpec 包含組合的規格說明
type ConfigSpec struct {
	Platform         PlatformConfig       `yaml:"platform"`
	CorePlugins      []PluginConfig       `yaml:"core_plugins"`
	CommunityPlugins []PluginConfig       `yaml:"community_plugins"`
	CustomPlugins    []PluginConfig       `yaml:"custom_plugins,omitempty"`
	Applications     []ApplicationConfig  `yaml:"applications"`
	Infrastructure   InfrastructureConfig `yaml:"infrastructure"`
	Security         SecurityConfig       `yaml:"security"`
	Dependencies     DependencyConfig     `yaml:"dependencies"`
	Health           HealthConfig         `yaml:"health"`
}

// PlatformConfig contains platform-level configuration
// zh: PlatformConfig 包含平台層級的配置
type PlatformConfig struct {
	Registry    RegistryConfig  `yaml:"registry"`
	Lifecycle   LifecycleConfig `yaml:"lifecycle"`
	Composition CompositionMeta `yaml:"composition"`
}

// RegistryConfig contains registry configuration
// zh: RegistryConfig 包含註冊表配置
type RegistryConfig struct {
	Enabled bool   `yaml:"enabled"`
	Type    string `yaml:"type"`
}

// LifecycleConfig contains lifecycle management configuration
// zh: LifecycleConfig 包含生命週期管理配置
type LifecycleConfig struct {
	Enabled bool   `yaml:"enabled"`
	Timeout string `yaml:"timeout"`
}

// CompositionMeta contains composition metadata
// zh: CompositionMeta 包含組合元資料
type CompositionMeta struct {
	Enabled    bool   `yaml:"enabled"`
	Validation string `yaml:"validation"`
}

// PluginConfig represents a plugin configuration
// zh: PluginConfig 代表插件配置
type PluginConfig struct {
	Name         string         `yaml:"name"`
	Type         string         `yaml:"type"`
	Enabled      bool           `yaml:"enabled"`
	Config       map[string]any `yaml:"config"`
	Dependencies []string       `yaml:"dependencies,omitempty"`
	Priority     int            `yaml:"priority,omitempty"`
}

// ApplicationConfig represents application configuration
// zh: ApplicationConfig 代表應用程式配置
type ApplicationConfig struct {
	Name    string         `yaml:"name"`
	Enabled bool           `yaml:"enabled"`
	Config  map[string]any `yaml:"config"`
}

// InfrastructureConfig contains infrastructure settings
// zh: InfrastructureConfig 包含基礎設施設定
type InfrastructureConfig struct {
	Cache    CacheConfig    `yaml:"cache"`
	EventBus EventBusConfig `yaml:"eventbus"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// CacheConfig contains cache configuration
// zh: CacheConfig 包含快取配置
type CacheConfig struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}

// EventBusConfig contains event bus configuration
// zh: EventBusConfig 包含事件匯流排配置
type EventBusConfig struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}

// LoggingConfig contains logging configuration
// zh: LoggingConfig 包含日誌配置
type LoggingConfig struct {
	Type       string             `yaml:"type" json:"type"`               // console, file, both
	Level      string             `yaml:"level" json:"level"`             // debug, info, warn, error
	Format     string             `yaml:"format" json:"format"`           // json, text
	Output     string             `yaml:"output" json:"output"`           // stdout, stderr, file path
	FileConfig *LoggingFileConfig `yaml:"file_config" json:"file_config"` // file rotation config
	OTEL       *LoggingOTELConfig `yaml:"otel" json:"otel"`               // OpenTelemetry config
}

// LoggingFileConfig defines file rotation configuration for logging
// zh: LoggingFileConfig 定義日誌檔案輪轉配置
type LoggingFileConfig struct {
	Filename   string `yaml:"filename" json:"filename"`       // log file path
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // max size in MB
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // max backup files
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // max age in days
	Compress   bool   `yaml:"compress" json:"compress"`       // compress old files
}

// LoggingOTELConfig defines OpenTelemetry configuration for logging
// zh: LoggingOTELConfig 定義日誌的 OpenTelemetry 配置
type LoggingOTELConfig struct {
	Enabled        bool   `yaml:"enabled" json:"enabled"`
	ServiceName    string `yaml:"service_name" json:"service_name"`
	ServiceVersion string `yaml:"service_version" json:"service_version"`
	TraceIDField   string `yaml:"trace_id_field" json:"trace_id_field"`
	SpanIDField    string `yaml:"span_id_field" json:"span_id_field"`
}

// SecurityConfig contains security settings
// zh: SecurityConfig 包含安全設定
type SecurityConfig struct {
	Authentication AuthenticationConfig `yaml:"authentication"`
	Authorization  AuthorizationConfig  `yaml:"authorization"`
	Permissions    PermissionsConfig    `yaml:"permissions"`
}

// AuthenticationConfig contains authentication settings
// zh: AuthenticationConfig 包含認證設定
type AuthenticationConfig struct {
	Provider string `yaml:"provider"`
	Required bool   `yaml:"required"`
}

// AuthorizationConfig contains authorization settings
// zh: AuthorizationConfig 包含授權設定
type AuthorizationConfig struct {
	Enabled bool `yaml:"enabled"`
}

// PermissionsConfig contains permission settings
// zh: PermissionsConfig 包含權限設定
type PermissionsConfig struct {
	DefaultRole string `yaml:"default_role"`
	Roles       []Role `yaml:"roles"`
}

// Role represents a security role
// zh: Role 代表安全角色
type Role struct {
	Name        string       `yaml:"name"`
	Permissions []Permission `yaml:"permissions"`
}

// Permission represents a permission
// zh: Permission 代表權限
type Permission struct {
	Action   string `yaml:"action"`
	Resource string `yaml:"resource"`
}

// DependencyConfig contains dependency resolution settings
// zh: DependencyConfig 包含依賴解析設定
type DependencyConfig struct {
	Validation  bool `yaml:"validation"`
	AutoResolve bool `yaml:"auto_resolve"`
}

// HealthConfig contains health check settings
// zh: HealthConfig 包含健康檢查設定
type HealthConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Interval string `yaml:"interval"`
	Timeout  string `yaml:"timeout"`
}

// ConfigLoader handles loading and parsing of composition configurations.
// zh: ConfigLoader 處理組合配置的載入和解析。
type ConfigLoader struct {
	basePath string
	config   *CompositionConfig
}

// NewConfigLoader creates a new configuration loader.
// zh: NewConfigLoader 建立新的配置載入器。
func NewConfigLoader(basePath string) *ConfigLoader {
	return &ConfigLoader{
		basePath: basePath,
	}
}

// LoadComposition loads a composition configuration from file.
// zh: LoadComposition 從檔案載入組合配置。
func (cl *ConfigLoader) LoadComposition(filePath string) (*CompositionConfig, error) {
	// Handle relative paths
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(cl.basePath, filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open composition file %s: %w", filePath, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read composition file %s: %w", filePath, err)
	}

	var config CompositionConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse composition file %s: %w", filePath, err)
	}

	cl.config = &config
	return &config, nil
}

// GetPluginConfigs returns all plugin configurations with their settings.
// zh: GetPluginConfigs 回傳所有插件配置及其設定。
func (cl *ConfigLoader) GetPluginConfigs() ([]contracts.PluginMetadata, error) {
	if cl.config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}

	var plugins []contracts.PluginMetadata

	// Process core plugins
	for _, plugin := range cl.config.Spec.CorePlugins {
		if plugin.Enabled {
			metadata := contracts.PluginMetadata{
				Name:         plugin.Name,
				Type:         plugin.Type,
				Category:     "core",
				Config:       plugin.Config,
				Enabled:      plugin.Enabled,
				Dependencies: plugin.Dependencies,
			}
			plugins = append(plugins, metadata)
		}
	}

	// Process community plugins
	for _, plugin := range cl.config.Spec.CommunityPlugins {
		if plugin.Enabled {
			metadata := contracts.PluginMetadata{
				Name:         plugin.Name,
				Type:         plugin.Type,
				Category:     "community",
				Config:       plugin.Config,
				Enabled:      plugin.Enabled,
				Dependencies: plugin.Dependencies,
			}
			plugins = append(plugins, metadata)
		}
	}

	// Process custom plugins
	for _, plugin := range cl.config.Spec.CustomPlugins {
		if plugin.Enabled {
			metadata := contracts.PluginMetadata{
				Name:         plugin.Name,
				Type:         plugin.Type,
				Category:     "custom",
				Config:       plugin.Config,
				Enabled:      plugin.Enabled,
				Dependencies: plugin.Dependencies,
			}
			plugins = append(plugins, metadata)
		}
	}

	return plugins, nil
}

// GetPluginConfig returns configuration for a specific plugin.
// zh: GetPluginConfig 回傳特定插件的配置。
func (cl *ConfigLoader) GetPluginConfig(pluginName string) (map[string]any, error) {
	if cl.config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}

	// Search in core plugins
	for _, plugin := range cl.config.Spec.CorePlugins {
		if plugin.Name == pluginName {
			return plugin.Config, nil
		}
	}

	// Search in community plugins
	for _, plugin := range cl.config.Spec.CommunityPlugins {
		if plugin.Name == pluginName {
			return plugin.Config, nil
		}
	}

	// Search in custom plugins
	for _, plugin := range cl.config.Spec.CustomPlugins {
		if plugin.Name == pluginName {
			return plugin.Config, nil
		}
	}

	return nil, fmt.Errorf("plugin %s not found in configuration", pluginName)
}

// GetPlatformConfig returns platform-level configuration.
// zh: GetPlatformConfig 回傳平台層級配置。
func (cl *ConfigLoader) GetPlatformConfig() (*PlatformConfig, error) {
	if cl.config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}
	return &cl.config.Spec.Platform, nil
}

// GetInfrastructureConfig returns infrastructure configuration.
// zh: GetInfrastructureConfig 回傳基礎設施配置。
func (cl *ConfigLoader) GetInfrastructureConfig() (*InfrastructureConfig, error) {
	if cl.config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}
	return &cl.config.Spec.Infrastructure, nil
}

// GetSecurityConfig returns security configuration.
// zh: GetSecurityConfig 回傳安全配置。
func (cl *ConfigLoader) GetSecurityConfig() (*SecurityConfig, error) {
	if cl.config == nil {
		return nil, fmt.Errorf("no configuration loaded")
	}
	return &cl.config.Spec.Security, nil
}

// ValidateConfiguration validates the loaded configuration.
// zh: ValidateConfiguration 驗證已載入的配置。
func (cl *ConfigLoader) ValidateConfiguration() error {
	if cl.config == nil {
		return fmt.Errorf("no configuration loaded")
	}

	// Basic validation
	if cl.config.APIVersion == "" {
		return fmt.Errorf("apiVersion is required")
	}

	if cl.config.Kind == "" {
		return fmt.Errorf("kind is required")
	}

	if cl.config.Metadata.Name == "" {
		return fmt.Errorf("metadata.name is required")
	}

	// Validate plugin names are unique
	pluginNames := make(map[string]bool)

	for _, plugin := range cl.config.Spec.CorePlugins {
		if pluginNames[plugin.Name] {
			return fmt.Errorf("duplicate plugin name: %s", plugin.Name)
		}
		pluginNames[plugin.Name] = true
	}

	for _, plugin := range cl.config.Spec.CommunityPlugins {
		if pluginNames[plugin.Name] {
			return fmt.Errorf("duplicate plugin name: %s", plugin.Name)
		}
		pluginNames[plugin.Name] = true
	}

	for _, plugin := range cl.config.Spec.CustomPlugins {
		if pluginNames[plugin.Name] {
			return fmt.Errorf("duplicate plugin name: %s", plugin.Name)
		}
		pluginNames[plugin.Name] = true
	}

	return nil
}
