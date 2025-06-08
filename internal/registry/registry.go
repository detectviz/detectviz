package registry

import (
	"github.com/detectviz/detectviz/internal/adapters/eventbus"
	cachestoreregistry "github.com/detectviz/detectviz/internal/registry/cachestore"
	notifierregistry "github.com/detectviz/detectviz/internal/registry/notifier"
	scheduleregistry "github.com/detectviz/detectviz/internal/registry/scheduler"
	"github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/notifier"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// RegistryContainer 包含 Detectviz 所有模組註冊項目
// zh: 用於封裝所有可注入與呼叫的服務或核心元件
type RegistryContainer struct {
	EventDispatcher eventbusiface.EventDispatcher
	Scheduler       scheduler.Scheduler
	Notifier        notifier.Notifier
	CacheStore      cachestore.CacheStore
}

// NewRegistry 建立並初始化所有必要的註冊模組
// zh: 使用指定的設定與 logger 建立平台所需資源
func NewRegistry(cfg ifconfig.ConfigProvider, schedulerProvider string, log logger.Logger) *RegistryContainer {
	return &RegistryContainer{
		EventDispatcher: NewInMemoryEventDispatcher(),
		Scheduler:       scheduleregistry.ProvideScheduler(schedulerProvider, log),
		Notifier:        notifierregistry.ProvideNotifier(cfg.GetNotifierConfigs(), log),
		CacheStore:      cachestoreregistry.RegisterCacheStore(cfg.GetCacheConfig()),
	}
}

// NewInMemoryEventDispatcher 建立新的事件總線 Dispatcher（使用記憶體）
// zh: 建立並整合內建與 plugin 註冊的所有事件處理器。
func NewInMemoryEventDispatcher() eventbusiface.EventDispatcher {
	dispatcher := eventbus.NewInMemoryDispatcher()
	registerAllHandlers(dispatcher)
	return dispatcher
}

// NewKafkaEventDispatcher 建立新的事件總線 Dispatcher（使用 Kafka 作為傳輸層）
// zh: 建立 Kafka 型態的事件總線 Dispatcher，預留未來整合 Kafka handler。
func NewKafkaEventDispatcher() eventbusiface.EventDispatcher {
	// TODO: 實作 Kafka-based Dispatcher，例如 eventbus.NewKafkaDispatcher()
	// dispatcher := eventbus.NewKafkaDispatcher()
	// registerAllHandlers(dispatcher)
	panic("KafkaEventDispatcher 尚未實作")
}

// registerAllHandlers 將所有 plugin handler 註冊到指定 Dispatcher
// zh: 將 Alert、Host、Metric、Task 等處理器註冊至事件總線。
func registerAllHandlers(d eventbusiface.EventDispatcher) {
	for _, h := range eventbus.LoadPluginAlertHandlers() {
		d.RegisterAlertHandler(h)
	}
	for _, h := range eventbus.LoadPluginHostHandlers() {
		d.RegisterHostHandler(h)
	}
	for _, h := range eventbus.LoadPluginMetricHandlers() {
		d.RegisterMetricHandler(h)
	}
	for _, h := range eventbus.LoadPluginTaskHandlers() {
		d.RegisterTaskHandler(h)
	}
}
