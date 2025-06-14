package contracts

// Plugin defines the basic interface that all plugins must implement.
// zh: Plugin 定義所有插件必須實作的基礎介面。
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Init(config any) error
	Shutdown() error
}

// LifecycleAware defines the interface for plugins that need lifecycle management.
// zh: LifecycleAware 定義需要生命週期管理的插件介面。
type LifecycleAware interface {
	OnRegister() error
	OnStart() error
	OnStop() error
	OnShutdown() error
}

// ConfigurablePlugin defines the interface for plugins that support configuration.
// zh: ConfigurablePlugin 定義支援配置的插件介面。
type ConfigurablePlugin interface {
	Plugin
	ValidateConfig(config any) error
	GetDefaultConfig() any
}

// PluginMetadata contains metadata information about a plugin.
// zh: PluginMetadata 包含插件的元資料資訊。
type PluginMetadata struct {
	Name         string         `yaml:"name" json:"name"`
	Version      string         `yaml:"version" json:"version"`
	Type         string         `yaml:"type" json:"type"`         // importer, exporter, integration, tool
	Category     string         `yaml:"category" json:"category"` // core, community, custom
	Description  string         `yaml:"description" json:"description"`
	Author       string         `yaml:"author" json:"author"`
	License      string         `yaml:"license" json:"license"`
	Dependencies []string       `yaml:"dependencies" json:"dependencies"`
	Config       map[string]any `yaml:"config" json:"config"`
	Enabled      bool           `yaml:"enabled" json:"enabled"`
}

// PluginFactory is a function type for creating plugin instances.
// zh: PluginFactory 是建立插件實例的函式類型。
type PluginFactory func(config any) (Plugin, error)

// Registry defines the interface for plugin registration and discovery.
// zh: Registry 定義插件註冊與發現的介面。
type Registry interface {
	RegisterPlugin(name string, factory PluginFactory) error
	GetPlugin(name string) (Plugin, error)
	ListPlugins() []string
	GetPluginMetadata(name string) (*PluginMetadata, error)
}
