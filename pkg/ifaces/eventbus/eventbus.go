package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// EventDispatcher defines the interface for dispatching typed events between modules.
// zh: EventDispatcher 定義模組間明確事件分派的抽象介面。
type EventDispatcher interface {
	// -------------------------------------------------------------------------
	// AlertTriggeredEvent
	// -------------------------------------------------------------------------

	// DispatchAlertTriggered dispatches an AlertTriggeredEvent to all registered handlers.
	// zh: 將 AlertTriggeredEvent 分派給所有已註冊的處理器。
	DispatchAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error

	// RegisterAlertHandler registers a new handler for AlertTriggeredEvent.
	// zh: 註冊一個處理 AlertTriggeredEvent 的事件處理器。
	RegisterAlertHandler(handler event.AlertEventHandler)

	// -------------------------------------------------------------------------
	// TaskCompletedEvent
	// -------------------------------------------------------------------------

	// DispatchTaskCompleted dispatches a TaskCompletedEvent to all registered handlers.
	// zh: 將 TaskCompletedEvent 分派給所有已註冊的處理器。
	DispatchTaskCompleted(ctx context.Context, e event.TaskCompletedEvent) error

	// RegisterTaskHandler registers a new handler for TaskCompletedEvent.
	// zh: 註冊一個處理 TaskCompletedEvent 的事件處理器。
	RegisterTaskHandler(handler event.TaskEventHandler)

	// -------------------------------------------------------------------------
	// HostDiscoveredEvent
	// -------------------------------------------------------------------------

	// DispatchHostDiscovered dispatches a HostDiscoveredEvent to all registered handlers.
	// zh: 將 HostDiscoveredEvent 分派給所有已註冊的處理器。
	DispatchHostDiscovered(ctx context.Context, e event.HostDiscoveredEvent) error

	// RegisterHostHandler registers a new handler for HostDiscoveredEvent.
	// zh: 註冊一個處理 HostDiscoveredEvent 的事件處理器。
	RegisterHostHandler(handler event.HostEventHandler)

	// -------------------------------------------------------------------------
	// MetricOverflowEvent
	// -------------------------------------------------------------------------

	// DispatchMetricOverflow dispatches a MetricOverflowEvent to all registered handlers.
	// zh: 將 MetricOverflowEvent 分派給所有已註冊的處理器。
	DispatchMetricOverflow(ctx context.Context, e event.MetricOverflowEvent) error

	// RegisterMetricHandler registers a new handler for MetricOverflowEvent.
	// zh: 註冊一個處理 MetricOverflowEvent 的事件處理器。
	RegisterMetricHandler(handler event.MetricEventHandler)
}

// AlertEventHandler defines a handler for AlertTriggeredEvent.
// zh: AlertEventHandler 定義處理 AlertTriggeredEvent 的事件處理器。
type AlertEventHandler interface {
	// HandleAlertTriggered processes the given alert event.
	// zh: 接收並處理 AlertTriggeredEvent。
	HandleAlertTriggered(ctx context.Context, event event.AlertTriggeredEvent) error
}

// TaskEventHandler defines a handler for TaskCompletedEvent.
// zh: TaskEventHandler 定義處理 TaskCompletedEvent 的事件處理器。
type TaskEventHandler interface {
	// HandleTaskCompleted processes the completed task event.
	// zh: 接收並處理 TaskCompletedEvent。
	HandleTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error
}

// HostEventHandler defines a handler for HostDiscoveredEvent.
// zh: HostEventHandler 定義處理 HostDiscoveredEvent 的事件處理器。
type HostEventHandler interface {
	// HandleHostDiscovered processes the discovered host event.
	// zh: 接收並處理 HostDiscoveredEvent。
	HandleHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error
}

// MetricEventHandler defines a handler for MetricOverflowEvent.
// zh: MetricEventHandler 定義處理 MetricOverflowEvent 的事件處理器。
type MetricEventHandler interface {
	// HandleMetricOverflow processes the metric overflow event.
	// zh: 接收並處理 MetricOverflowEvent。
	HandleMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error
}

// RegisterPluginTaskHandler 是 plugin 註冊 TaskEventHandler 的統一介面。
// zh: 提供 plugin 註冊 TaskCompletedEvent 處理器。
type RegisterPluginTaskHandler interface {
	RegisterTaskHandler(handler TaskEventHandler)
}

// RegisterPluginHostHandler 是 plugin 註冊 HostEventHandler 的統一介面。
// zh: 提供 plugin 註冊 HostDiscoveredEvent 處理器。
type RegisterPluginHostHandler interface {
	RegisterHostHandler(handler HostEventHandler)
}

// RegisterPluginMetricHandler 是 plugin 註冊 MetricEventHandler 的統一介面。
// zh: 提供 plugin 註冊 MetricOverflowEvent 處理器。
type RegisterPluginMetricHandler interface {
	RegisterMetricHandler(handler MetricEventHandler)
}
