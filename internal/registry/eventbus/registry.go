package eventbus

import (
	"fmt"

	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

// DispatcherOptions allows injection of custom handlers for each event type.
// zh: DispatcherOptions 提供事件處理器的組態選項，支援外部擴充與測試注入。
type DispatcherOptions struct {
	AlertHandlers  []event.AlertEventHandler
	HostHandlers   []event.HostEventHandler
	MetricHandlers []event.MetricEventHandler
	TaskHandlers   []event.TaskEventHandler
}

// NewInMemoryEventDispatcherWithOptions constructs an in-memory dispatcher with injected handlers.
// zh: 根據傳入的 DispatcherOptions 建立記憶體型事件分派器。
func NewInMemoryEventDispatcherWithOptions(opt DispatcherOptions) eventbusiface.EventDispatcher {
	dispatcher := eventbusadapter.NewInMemoryDispatcher()

	for _, h := range opt.AlertHandlers {
		dispatcher.RegisterAlertHandler(h)
	}
	for _, h := range opt.HostHandlers {
		dispatcher.RegisterHostHandler(h)
	}
	for _, h := range opt.MetricHandlers {
		dispatcher.RegisterMetricHandler(h)
	}
	for _, h := range opt.TaskHandlers {
		dispatcher.RegisterTaskHandler(h)
	}

	return dispatcher
}

// NewEventDispatcher 建立指定 provider 的事件總線實作。
// zh: 根據指定 provider 名稱建立對應的事件處理器，可擴充支援 kafka/nats 等實作。
func NewEventDispatcher(provider string) (eventbusiface.EventDispatcher, error) {
	p, err := GetProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventbus provider: %w", err)
	}
	return p, nil
}
