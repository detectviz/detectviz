package eventbus

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

// InMemoryDispatcher provides a thread-safe, in-process implementation of the EventDispatcher interface.
// zh: InMemoryDispatcher 提供執行緒安全、僅在記憶體中運作的事件分派器。
type InMemoryDispatcher struct {
	alertHandlers  []event.AlertEventHandler
	taskHandlers   []event.TaskEventHandler
	hostHandlers   []event.HostEventHandler
	metricHandlers []event.MetricEventHandler
	mu             sync.RWMutex
}

// NewInMemoryDispatcher creates a new in-memory event dispatcher instance.
// zh: 建立新的記憶體內事件分派器。
func NewInMemoryDispatcher() eventbus.EventDispatcher {
	return &InMemoryDispatcher{}
}

// DispatchAlertTriggered dispatches AlertTriggeredEvent to all registered alert handlers.
// zh: 將 AlertTriggeredEvent 傳遞給所有已註冊的告警處理器。
func (d *InMemoryDispatcher) DispatchAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.alertHandlers {
		if err := h.HandleAlertTriggered(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterAlertHandler registers an alert event handler.
// zh: 註冊 Alert 事件處理器。
func (d *InMemoryDispatcher) RegisterAlertHandler(h event.AlertEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.alertHandlers = append(d.alertHandlers, h)
}

// DispatchTaskCompleted dispatches TaskCompletedEvent to all task handlers.
// zh: 傳遞 TaskCompletedEvent 給所有註冊的任務處理器。
func (d *InMemoryDispatcher) DispatchTaskCompleted(ctx context.Context, e event.TaskCompletedEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.taskHandlers {
		if err := h.HandleTaskCompleted(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterTaskHandler registers a task event handler.
// zh: 註冊 Task 事件處理器。
func (d *InMemoryDispatcher) RegisterTaskHandler(h event.TaskEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.taskHandlers = append(d.taskHandlers, h)
}

// DispatchHostDiscovered dispatches HostDiscoveredEvent to all host handlers.
// zh: 傳遞 HostDiscoveredEvent 給所有註冊的主機處理器。
func (d *InMemoryDispatcher) DispatchHostDiscovered(ctx context.Context, e event.HostDiscoveredEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.hostHandlers {
		if err := h.HandleHostDiscovered(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterHostHandler registers a host event handler.
// zh: 註冊 Host 事件處理器。
func (d *InMemoryDispatcher) RegisterHostHandler(h event.HostEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.hostHandlers = append(d.hostHandlers, h)
}

// DispatchMetricOverflow dispatches MetricOverflowEvent to all metric handlers.
// zh: 傳遞 MetricOverflowEvent 給所有註冊的指標處理器。
func (d *InMemoryDispatcher) DispatchMetricOverflow(ctx context.Context, e event.MetricOverflowEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.metricHandlers {
		if err := h.HandleMetricOverflow(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterMetricHandler registers a metric event handler.
// zh: 註冊 Metric 事件處理器。
func (d *InMemoryDispatcher) RegisterMetricHandler(h event.MetricEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.metricHandlers = append(d.metricHandlers, h)
}
