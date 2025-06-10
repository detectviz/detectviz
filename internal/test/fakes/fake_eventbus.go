package fakes

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

// FakeEventBus implements eventbus.EventDispatcher for tests.
// zh: 測試用 EventDispatcher 假實作，紀錄事件與 handler 呼叫狀態。
type FakeEventBus struct {
	AlertEvents  []event.AlertTriggeredEvent
	TaskEvents   []event.TaskCompletedEvent
	HostEvents   []event.HostDiscoveredEvent
	MetricEvents []event.MetricOverflowEvent

	AlertHandlers  []eventbus.AlertEventHandler
	TaskHandlers   []eventbus.TaskEventHandler
	HostHandlers   []eventbus.HostEventHandler
	MetricHandlers []eventbus.MetricEventHandler

	DispatchAlertErr  error
	DispatchTaskErr   error
	DispatchHostErr   error
	DispatchMetricErr error
}

func (b *FakeEventBus) DispatchAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error {
	b.AlertEvents = append(b.AlertEvents, e)
	for _, h := range b.AlertHandlers {
		if err := h.HandleAlertTriggered(ctx, e); err != nil {
			return err
		}
	}
	return b.DispatchAlertErr
}

func (b *FakeEventBus) RegisterAlertHandler(h eventbus.AlertEventHandler) {
	b.AlertHandlers = append(b.AlertHandlers, h)
}

func (b *FakeEventBus) DispatchTaskCompleted(ctx context.Context, e event.TaskCompletedEvent) error {
	b.TaskEvents = append(b.TaskEvents, e)
	for _, h := range b.TaskHandlers {
		if err := h.HandleTaskCompleted(ctx, e); err != nil {
			return err
		}
	}
	return b.DispatchTaskErr
}

func (b *FakeEventBus) RegisterTaskHandler(h eventbus.TaskEventHandler) {
	b.TaskHandlers = append(b.TaskHandlers, h)
}

func (b *FakeEventBus) DispatchHostDiscovered(ctx context.Context, e event.HostDiscoveredEvent) error {
	b.HostEvents = append(b.HostEvents, e)
	for _, h := range b.HostHandlers {
		if err := h.HandleHostDiscovered(ctx, e); err != nil {
			return err
		}
	}
	return b.DispatchHostErr
}

func (b *FakeEventBus) RegisterHostHandler(h eventbus.HostEventHandler) {
	b.HostHandlers = append(b.HostHandlers, h)
}

func (b *FakeEventBus) DispatchMetricOverflow(ctx context.Context, e event.MetricOverflowEvent) error {
	b.MetricEvents = append(b.MetricEvents, e)
	for _, h := range b.MetricHandlers {
		if err := h.HandleMetricOverflow(ctx, e); err != nil {
			return err
		}
	}
	return b.DispatchMetricErr
}

func (b *FakeEventBus) RegisterMetricHandler(h eventbus.MetricEventHandler) {
	b.MetricHandlers = append(b.MetricHandlers, h)
}
