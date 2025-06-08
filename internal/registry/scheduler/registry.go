package scheduler

import (
	adapters "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// ProvideScheduler returns a scheduler instance based on the given name.
// zh: 根據指定名稱回傳對應的排程器實作。
func ProvideScheduler(name string, log logger.Logger) scheduler.Scheduler {
	switch name {
	case "cron":
		return adapters.NewCronScheduler(log)
	case "workerpool":
		return adapters.NewWorkerPoolScheduler(4, log) // 預設 4 workers
	case "mock":
		return adapters.NewMockScheduler()
	default:
		log.Warn("unknown scheduler type, fallback to mock", "name", name)
		return adapters.NewMockScheduler()
	}
}
