package configtypes

import "errors"

//
// CacheConfig
//

// CacheConfig defines the configuration for the cache module.
// zh: 支援記憶體與 Redis 快取，透過 backend 參數選擇快取後端。
type CacheConfig struct {
	Backend string      `json:"backend" yaml:"backend"` // zh: 快取後端類型，支援 "memory" 或 "redis"
	Redis   RedisConfig `json:"redis" yaml:"redis"`     // zh: Redis 相關設定，當 Backend 為 redis 時生效
}

// RedisConfig defines the configuration for Redis cache.
// zh: Redis 快取設定，包含連線地址、密碼與資料庫編號。
type RedisConfig struct {
	Address  string `json:"address" yaml:"address"`   // zh: Redis 伺服器連線位址 (host:port)
	Password string `json:"password" yaml:"password"` // zh: Redis 密碼，可留空
	DB       int    `json:"db" yaml:"db"`             // zh: Redis 使用的資料庫編號，預設為 0
}

// Validate implements the Configurable interface.
func (c CacheConfig) Validate() error {
	if c.Backend == "" {
		return errors.New("cache.backend is required")
	}
	if c.Backend == "redis" {
		if c.Redis.Address == "" {
			return errors.New("cache.redis.address is required when backend is redis")
		}
	}
	return nil
}

//
// NotifierConfig
//

// NotifierConfig defines the configuration for a notifier.
// zh: 定義單一通知通道的設定，包括名稱、類型、目標與啟用狀態。
type NotifierConfig struct {
	Name   string `json:"name"`   // zh: 通道名稱，如 "email", "slack", "webhook"
	Type   string `json:"type"`   // zh: 通道類型，決定通知方式
	Target string `json:"target"` // zh: 通知目標，如電子郵件地址或 webhook URL
	Enable bool   `json:"enable"` // zh: 是否啟用此通知通道
}

// Validate implements the Configurable interface.
func (c NotifierConfig) Validate() error {
	if c.Type == "" {
		return errors.New("notifier.type is required")
	}
	if c.Target == "" {
		return errors.New("notifier.target is required")
	}
	return nil
}

//
// LoggerConfig
//

// LoggerConfig defines logging configuration.
// zh: 記錄模組設定，包括等級、輸出格式與目的地。
type LoggerConfig struct {
	Level  string   `mapstructure:"level" yaml:"level"`     // zh: 記錄層級，如 "info", "debug"
	Format string   `mapstructure:"format" yaml:"format"`   // zh: 輸出格式，支援 "json", "text"
	Output []string `mapstructure:"outputs" yaml:"outputs"` // zh: 輸出管道，例如 "stdout", "file"
}

func (c LoggerConfig) Validate() error {
	if c.Level == "" {
		return errors.New("logger.level is required")
	}
	return nil
}

//
// SchedulerConfig
//

// SchedulerConfig defines scheduling related configuration.
// zh: 排程模組設定，包含時區設定。
type SchedulerConfig struct {
	Timezone string `mapstructure:"timezone" yaml:"timezone"` // zh: 排程使用的時區，例如 "Asia/Taipei"
}

func (c SchedulerConfig) Validate() error {
	if c.Timezone == "" {
		return errors.New("scheduler.timezone is required")
	}
	return nil
}

//
// AlertConfig
//

// AlertConfig defines alert levels configuration.
// zh: 警示等級設定，用於控制哪些等級的事件會觸發警報。
type AlertConfig struct {
	Levels []string `mapstructure:"levels" yaml:"levels"` // zh: 警示等級列表，如 ["critical", "warning"]
}

func (c AlertConfig) Validate() error {
	if len(c.Levels) == 0 {
		return errors.New("alert.levels must not be empty")
	}
	return nil
}

//
// BusConfig
//

// BusConfig defines event bus configuration.
// zh: 事件總線設定，包含後端類型。
type BusConfig struct {
	Backend string `mapstructure:"backend" yaml:"backend"` // zh: 事件總線後端類型，如 "kafka", "rabbitmq"
}

func (c BusConfig) Validate() error {
	if c.Backend == "" {
		return errors.New("eventbus.backend is required")
	}
	return nil
}

//
// EncryptionConfig
//

// EncryptionConfig defines encryption related configuration.
// zh: 加密模組設定，包含演算法與金鑰輪替天數。
type EncryptionConfig struct {
	Algorithm       string `mapstructure:"algorithm" yaml:"algorithm"`                 // zh: 加密演算法名稱，如 "AES", "RSA"
	KeyRotationDays int    `mapstructure:"key_rotation_days" yaml:"key_rotation_days"` // zh: 金鑰輪替週期（天數）
}

func (c EncryptionConfig) Validate() error {
	if c.Algorithm == "" {
		return errors.New("encryption.algorithm is required")
	}
	return nil
}

//
// WebConfig
//

// WebConfig defines web server related configuration.
// zh: Web 伺服器設定，包括靜態檔案路徑與 UI 啟用狀態。
type WebConfig struct {
	StaticPath string `mapstructure:"static_path" yaml:"static_path"` // zh: 靜態檔案路徑
	EnableUI   bool   `mapstructure:"enable_ui" yaml:"enable_ui"`     // zh: 是否啟用 Web UI
}

func (c WebConfig) Validate() error {
	return nil
}

//
// InfraConfig
//

// InfraConfig defines infrastructure related configuration.
// zh: 基礎設施設定，包含是否啟用重試與最大重試次數。
type InfraConfig struct {
	EnableRetry bool `mapstructure:"enable_retry" yaml:"enable_retry"` // zh: 是否啟用重試機制
	MaxRetries  int  `mapstructure:"max_retries" yaml:"max_retries"`   // zh: 最大重試次數
}

func (c InfraConfig) Validate() error {
	return nil
}

//
// ValidationConfig
//

// ValidationConfig defines configuration for validation behavior.
// zh: 驗證模組設定，包含是否啟用與嚴格模式。
type ValidationConfig struct {
	Enabled bool `mapstructure:"enabled" yaml:"enabled"` // zh: 是否啟用資料驗證
	Strict  bool `mapstructure:"strict" yaml:"strict"`   // zh: 是否啟用嚴格模式驗證
}

func (c ValidationConfig) Validate() error {
	return nil
}

//
// RBACConfig
//

// RBACConfig defines role-based access control configuration.
// zh: RBAC 設定，包含是否啟用與預設存取策略。
type RBACConfig struct {
	Enabled       bool   `mapstructure:"enabled" yaml:"enabled"`               // zh: 是否啟用 RBAC
	DefaultPolicy string `mapstructure:"default_policy" yaml:"default_policy"` // zh: 預設存取策略，如 "allow" 或 "deny"
}

func (c RBACConfig) Validate() error {
	return nil
}

//
// AuthConfig
//

// AuthConfig defines authentication configuration.
// zh: 認證模組設定，包含認證類型。
type AuthConfig struct {
	Type string `mapstructure:"type" yaml:"type"` // zh: 認證類型，如 "basic", "oauth"
}

func (c AuthConfig) Validate() error {
	if c.Type == "" {
		return errors.New("auth.type is required")
	}
	return nil
}

//
// StoreConfig
//

// StoreConfig defines storage backend configuration.
// zh: 儲存模組設定，包含驅動類型。
type StoreConfig struct {
	Driver string `mapstructure:"driver" yaml:"driver"` // zh: 儲存驅動名稱，如 "local", "s3"
}

func (c StoreConfig) Validate() error {
	if c.Driver == "" {
		return errors.New("store.driver is required")
	}
	return nil
}

//
// MiddlewareConfig
//

// MiddlewareConfig defines middleware related configuration.
// zh: 中介軟體設定，包含 CORS 與流量限制。
type MiddlewareConfig struct {
	CORS      bool `mapstructure:"cors" yaml:"cors"`             // zh: 是否啟用跨來源資源共享 (CORS)
	RateLimit int  `mapstructure:"rate_limit" yaml:"rate_limit"` // zh: 請求速率限制，單位為每秒請求數
}

func (c MiddlewareConfig) Validate() error {
	return nil
}

//
// ServicesConfig
//

// ServicesConfig defines configuration for various services.
// zh: 服務模組設定，包含配額功能是否啟用。
type ServicesConfig struct {
	QuotaEnabled bool `mapstructure:"quota_enabled" yaml:"quota_enabled"` // zh: 是否啟用配額功能
}

func (c ServicesConfig) Validate() error {
	return nil
}

//
// SystemConfig
//

// SystemConfig defines system level configuration.
// zh: 系統設定，包含診斷功能是否啟用。
type SystemConfig struct {
	DiagnosticsEnabled bool `mapstructure:"diagnostics_enabled" yaml:"diagnostics_enabled"` // zh: 是否啟用診斷功能
}

func (c SystemConfig) Validate() error {
	return nil
}

//
// ModulesConfig
//

// ModulesConfig defines configuration for enabled modules.
// zh: 模組設定，使用 map 指定各模組是否啟用。
type ModulesConfig struct {
	Enabled map[string]bool `mapstructure:"enabled" yaml:"enabled"` // zh: 模組啟用狀態對應表
}

func (c ModulesConfig) Validate() error {
	return nil
}

//
// PluginRuntimeConfig
//

// PluginRuntimeConfig defines plugin runtime configuration.
// zh: 外掛執行時設定，包含路徑與自動載入選項。
type PluginRuntimeConfig struct {
	Path string `mapstructure:"path" yaml:"path"` // zh: 外掛路徑
	Auto bool   `mapstructure:"auto" yaml:"auto"` // zh: 是否自動載入外掛
}

func (c PluginRuntimeConfig) Validate() error {
	return nil
}

//
// ServerConfig
//

// ServerConfig defines server related configuration.
// zh: 伺服器設定，包含監聽埠號。
type ServerConfig struct {
	Port int `mapstructure:"port" yaml:"port"` // zh: 伺服器監聽埠號，必須大於 0
}

func (c ServerConfig) Validate() error {
	if c.Port <= 0 {
		return errors.New("server.port must be > 0")
	}
	return nil
}

//
// MetricsConfig
//

// MetricsConfig defines configuration for metrics system.
// zh: 指標模組設定，包含啟用狀態與預設資料來源類型。
type MetricsConfig struct {
	Enabled       bool   `mapstructure:"enabled" yaml:"enabled"`               // zh: 是否啟用指標模組
	DefaultSource string `mapstructure:"default_source" yaml:"default_source"` // zh: 預設使用的資料來源，例如 "prometheus"
}

func (c MetricsConfig) Validate() error {
	if c.Enabled && c.DefaultSource == "" {
		return errors.New("metrics.default_source is required when enabled")
	}
	return nil
}

//
// RuleConfig
//

// RuleConfig defines configuration for rule engine.
// zh: 規則引擎模組設定，包含預設條件與啟用狀態。
type RuleConfig struct {
	Enabled      bool     `mapstructure:"enabled" yaml:"enabled"`             // zh: 是否啟用規則模組
	DefaultTypes []string `mapstructure:"default_types" yaml:"default_types"` // zh: 預設支援的規則類型，例如 ["threshold", "anomaly"]
}

func (c RuleConfig) Validate() error {
	return nil
}

//
// DatasourceConfig
//

// DatasourceConfig defines configuration for datasource modules.
// zh: 資料來源模組設定，包含啟用狀態與可用來源清單。
type DatasourceConfig struct {
	Enabled bool     `mapstructure:"enabled" yaml:"enabled"` // zh: 是否啟用資料來源管理模組
	Types   []string `mapstructure:"types" yaml:"types"`     // zh: 可用資料來源類型列表，例如 ["influxdb", "prometheus"]
}

func (c DatasourceConfig) Validate() error {
	return nil
}

//
// AuthProviderConfig
//

// AuthProviderConfig defines provider-specific auth configuration.
// zh: 認證供應商設定，支援如 keycloak 或其他擴充提供者。
type AuthProviderConfig struct {
	Provider string `mapstructure:"provider" yaml:"provider"` // zh: 認證供應商名稱，如 "keycloak"
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"` // zh: 認證伺服器端點
}

func (c AuthProviderConfig) Validate() error {
	if c.Provider == "" {
		return errors.New("auth_provider.provider is required")
	}
	return nil
}

//
// NodeMetadata
//

// NodeMetadata defines node-level identity and scheduling attributes.
// zh: 節點識別與調度屬性設定，用於標示節點角色、可用性與調度條件。
type NodeMetadata struct {
	Name     string            `mapstructure:"name" yaml:"name"`         // zh: 節點名稱，需唯一標示此節點
	Type     string            `mapstructure:"type" yaml:"type"`         // zh: 節點類型（角色），如 "monitor", "relay"
	Labels   map[string]string `mapstructure:"labels" yaml:"labels"`     // zh: 自訂節點標籤，用於篩選、排程與組織
	Zone     string            `mapstructure:"zone" yaml:"zone"`         // zh: 所在區域或機房，用於容錯與地理調度
	Cluster  string            `mapstructure:"cluster" yaml:"cluster"`   // zh: 所屬叢集名稱，用於節點群組識別
	Priority int               `mapstructure:"priority" yaml:"priority"` // zh: 調度優先級（數字越小越優先）
	Standby  bool              `mapstructure:"standby" yaml:"standby"`   // zh: 是否為備援節點
	Enabled  bool              `mapstructure:"enabled" yaml:"enabled"`   // zh: 是否啟用該節點
	Meta     map[string]string `mapstructure:"meta" yaml:"meta"`         // zh: 額外自訂屬性（如內部代碼、記憶體等）
}

// Validate validates the NodeMetadata struct.
func (n NodeMetadata) Validate() error {
	if n.Name == "" {
		return errors.New("node.name is required")
	}
	if n.Type == "" {
		return errors.New("node.type is required")
	}
	return nil
}

//
// PathConfig
//

// PathConfig defines system path configuration.
// zh: 系統資料與資源目錄設定
type PathConfig struct {
	Data         string `mapstructure:"data" yaml:"data"`                 // zh: 資料儲存目錄
	Logs         string `mapstructure:"logs" yaml:"logs"`                 // zh: 日誌檔案儲存路徑
	Plugins      string `mapstructure:"plugins" yaml:"plugins"`           // zh: 外掛目錄
	Provisioning string `mapstructure:"provisioning" yaml:"provisioning"` // zh: 預設設定初始化目錄
}

//
// CoreConfig
//

// CoreConfig defines core application-level configuration.
// zh: 應用層級基本設定，如模式、實例名稱與路徑。
type CoreConfig struct {
	AppMode      string       `mapstructure:"app_mode" yaml:"app_mode"`           // zh: 執行模式，如 "development", "production"
	InstanceName string       `mapstructure:"instance_name" yaml:"instance_name"` // zh: 實例名稱，預設為 ${HOSTNAME}
	HotReload    bool         `mapstructure:"hot_reload" yaml:"hot_reload"`       // zh: 是否啟用熱重載（限開發環境）
	Paths        PathConfig   `mapstructure:"paths" yaml:"paths"`                 // zh: 系統目錄路徑設定
	Node         NodeMetadata `mapstructure:"node" yaml:"node"`                   // zh: 節點描述與調度屬性
}

// Validate implements the Configurable interface.
func (c CoreConfig) Validate() error {
	if c.AppMode == "" {
		return errors.New("core.app_mode is required")
	}
	if c.InstanceName == "" {
		return errors.New("core.instance_name is required")
	}
	if err := c.Node.Validate(); err != nil {
		return err
	}
	return nil
}

//
// PluginConfig
//
// file: config_plugin.go

// PluginConfig defines external plugin configuration.
// zh: 外部外掛設定，包含是否允許未簽名外掛、掃描目錄。
type PluginConfig struct {
	AllowUnsigned bool   `mapstructure:"allow_unsigned" yaml:"allow_unsigned"` // zh: 是否允許未簽名外掛（不建議於正式環境）
	Directory     string `mapstructure:"directory" yaml:"directory"`           // zh: 外掛儲存目錄路徑
}

func (c PluginConfig) Validate() error {
	if c.Directory == "" {
		return errors.New("plugin.directory is required")
	}
	return nil
}

//
// QuotaConfig
//
// file: config_services.go

// QuotaConfig defines resource quota limitation.
// zh: 配額限制設定，用於限制使用資源上限（如每組件最多條件數等）。
type QuotaConfig struct {
	Enabled bool           `mapstructure:"enabled" yaml:"enabled"` // zh: 是否啟用配額控制
	Limits  map[string]int `mapstructure:"limits" yaml:"limits"`   // zh: 模組對應配額上限，例如 "alert.rules": 1000
}

func (c QuotaConfig) Validate() error {
	return nil
}

//
// SecretsManagerConfig
//
// file: config_encryption.go

// SecretsManagerConfig defines secure secret storage configuration.
// zh: 機敏設定儲存相關設定，用於集中加密管理或轉接 AWS KMS、Vault。
type SecretsManagerConfig struct {
	Provider string `mapstructure:"provider" yaml:"provider"` // zh: 秘密管理服務提供者，如 "vault", "aws-kms"
	Endpoint string `mapstructure:"endpoint" yaml:"endpoint"` // zh: 服務端點，如 Vault 或 KMS API URL
	Token    string `mapstructure:"token" yaml:"token"`       // zh: 用於存取秘密管理服務的 Token 或金鑰
}

func (c SecretsManagerConfig) Validate() error {
	if c.Provider != "" && c.Endpoint == "" {
		return errors.New("secrets.endpoint is required if provider is set")
	}
	return nil
}

//
// UserConfig
//
// file: config_user.go

// UserConfig defines configuration for user registration and roles.
// zh: 使用者註冊與預設角色設定。
type UserConfig struct {
	AllowSignUp   bool   `mapstructure:"allow_sign_up" yaml:"allow_sign_up"`     // zh: 是否允許開放註冊
	DefaultRole   string `mapstructure:"default_role" yaml:"default_role"`       // zh: 預設註冊使用者的角色，如 "Viewer"
	AutoAssignOrg bool   `mapstructure:"auto_assign_org" yaml:"auto_assign_org"` // zh: 是否自動指派使用者到預設組織
}

func (c UserConfig) Validate() error {
	if c.AllowSignUp && c.DefaultRole == "" {
		return errors.New("user.default_role is required when sign up is allowed")
	}
	return nil
}

//
// FeatureToggleConfig
//
// file: config_feature.go

// FeatureToggleConfig defines feature switch settings.
// zh: 功能開關設定，用於控制實驗性或進階功能是否啟用。
type FeatureToggleConfig struct {
	Flags map[string]bool `mapstructure:"flags" yaml:"flags"` // zh: 功能旗標列表，例如 "rule_editor": true
}

func (c FeatureToggleConfig) Validate() error {
	return nil
}
