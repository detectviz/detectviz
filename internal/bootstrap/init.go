package bootstrap

import (
	"github.com/detectviz/detectviz/internal/adapters/alert/flux"
	"github.com/detectviz/detectviz/internal/adapters/alert/prom"
	"github.com/detectviz/detectviz/internal/alert"
	"github.com/detectviz/detectviz/internal/registry"

	alertregistry "github.com/detectviz/detectviz/internal/registry/alert"

	"github.com/detectviz/detectviz/pkg/config"
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
)

var (
	Config              ifconfig.ConfigProvider
	Registry            *registry.RegistryContainer
	AlertEvaluatorStore *alertregistry.AlertEvaluatorRegistry
)

// Init 系統初始化邏輯
// zh: 執行整體系統初始化，包含設定載入與核心元件註冊。
func Init() {
	// zh: 載入預設設定提供者（含 logger, notifier config 等）
	Config = config.NewDefaultProvider()

	// zh: 註冊所有平台元件（notifier、scheduler、cachestore 等）
	Registry = registry.NewRegistry(Config, "cron", Config.Logger())

	// zh: 建立 AlertEvaluator 註冊中心（含 prometheus, flux）
	AlertEvaluatorStore = alertregistry.NewDefaultAlertEvaluatorRegistry(Config.Logger())

	// zh: 註冊預設 prometheus 告警評估器至註冊中心
	AlertEvaluatorStore.Register("prometheus", prom.NewEvaluator(Config.Logger()))

	// zh: 註冊 flux 告警評估器至註冊中心
	AlertEvaluatorStore.Register("flux", flux.NewEvaluator(Config.Logger()))

	// zh: 將設定注入給 alert 模組（含 evaluator、告警處理器）
	alert.Init(Config)
}
