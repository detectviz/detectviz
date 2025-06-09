package bootstrap

import (
	fluxadapter "github.com/detectviz/detectviz/internal/adapters/alert/flux"
	promadapter "github.com/detectviz/detectviz/internal/adapters/alert/prom"
	"github.com/detectviz/detectviz/internal/alert"
	"github.com/detectviz/detectviz/internal/registry"

	alertregistry "github.com/detectviz/detectviz/internal/registry/alert"

	"github.com/detectviz/detectviz/pkg/config"
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
)

var (
	Config ifconfig.ConfigProvider
	// Registry is the global runtime component registry.
	// zh: Registry 為平台執行期模組註冊中心。
	Registry            *registry.RegistryContainer
	AlertEvaluatorStore *alertregistry.AlertEvaluatorRegistry
)

// Init initializes the entire system.
// zh: 執行整體系統初始化，包含設定載入與核心元件註冊。
func Init() {
	initConfig()
	initRegistry()
	initAlertEvaluators()
	initAlertModule()
}

// initConfig loads the default config provider.
// zh: 載入預設設定提供者。
func initConfig() {
	Config = config.NewDefaultProvider()
}

// initRegistry sets up the runtime registry.
// zh: 註冊所有平台元件（如 notifier、scheduler 等）。
func initRegistry() {
	Registry = registry.NewRegistry(Config, "cron", Config.Logger())
}

// initAlertEvaluators registers alert evaluators (prometheus, flux).
// zh: 註冊預設告警評估器。
func initAlertEvaluators() {
	AlertEvaluatorStore = alertregistry.NewDefaultAlertEvaluatorRegistry(Config.Logger())
	AlertEvaluatorStore.Register("prometheus", promadapter.NewEvaluator(Config.Logger()))
	AlertEvaluatorStore.Register("flux", fluxadapter.NewEvaluator(Config.Logger()))
}

// initAlertModule injects config into the alert module.
// zh: 將設定注入 alert 模組。
func initAlertModule() {
	alert.Init(Config)
}
