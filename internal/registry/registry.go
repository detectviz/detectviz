package registry

import (
	"github.com/detectviz/detectviz/internal/registry/notifier"
	"github.com/detectviz/detectviz/internal/registry/scheduler"
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
	scheduleriface "github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// Registry 封裝所有平台元件註冊邏輯
// zh: 統一集中註冊 Notifier、Scheduler 等元件。
type Registry struct {
	Notifiers []notifieriface.Notifier
	Scheduler scheduleriface.Scheduler
}

// New 建立一個新的 Registry，根據設定注入所有子模組
// zh: 從設定與 logger 建立註冊中心
func New(cfgs []configtypes.NotifierConfig, schedulerName string, log logger.Logger) *Registry {
	return &Registry{
		Notifiers: notifier.NewNotifierRegistry(cfgs, log),
		Scheduler: scheduler.ProvideScheduler(schedulerName, log),
	}
}
