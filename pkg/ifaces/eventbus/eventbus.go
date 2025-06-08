package eventbus

import "context"

// EventDispatcher defines the interface for dispatching typed events between modules.
// zh: EventDispatcher 定義模組間明確事件分派的抽象介面。
type EventDispatcher interface {
	// -------------------------------------------------------------------------
	// AlertTriggeredEvent
	// -------------------------------------------------------------------------

	// DispatchAlertTriggered dispatches an AlertTriggeredEvent to all registered handlers.
	// zh: 將 AlertTriggeredEvent 分派給所有已註冊的處理器。
	DispatchAlertTriggered(ctx context.Context, event any) error

	// RegisterAlertHandler registers a new handler for AlertTriggeredEvent.
	// zh: 註冊一個處理 AlertTriggeredEvent 的事件處理器。
	RegisterAlertHandler(handler any)

	// -------------------------------------------------------------------------
	// TaskCompletedEvent
	// -------------------------------------------------------------------------

	// DispatchTaskCompleted dispatches a TaskCompletedEvent to all registered handlers.
	// zh: 將 TaskCompletedEvent 分派給所有已註冊的處理器。
	DispatchTaskCompleted(ctx context.Context, event any) error

	// RegisterTaskHandler registers a new handler for TaskCompletedEvent.
	// zh: 註冊一個處理 TaskCompletedEvent 的事件處理器。
	RegisterTaskHandler(handler any)
	// -------------------------------------------------------------------------
	// HostDiscoveredEvent
	// -------------------------------------------------------------------------

	// DispatchHostDiscovered dispatches a HostDiscoveredEvent to all registered handlers.
	// zh: 將 HostDiscoveredEvent 分派給所有已註冊的處理器。
	DispatchHostDiscovered(ctx context.Context, event any) error

	// RegisterHostHandler registers a new handler for HostDiscoveredEvent.
	// zh: 註冊一個處理 HostDiscoveredEvent 的事件處理器。
	RegisterHostHandler(handler any)

	// -------------------------------------------------------------------------
	// MetricOverflowEvent
	// -------------------------------------------------------------------------

	// DispatchMetricOverflow dispatches a MetricOverflowEvent to all registered handlers.
	// zh: 將 MetricOverflowEvent 分派給所有已註冊的處理器。
	DispatchMetricOverflow(ctx context.Context, event any) error

	// RegisterMetricHandler registers a new handler for MetricOverflowEvent.
	// zh: 註冊一個處理 MetricOverflowEvent 的事件處理器。
	RegisterMetricHandler(handler any)
}
