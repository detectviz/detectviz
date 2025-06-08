package event

import "context"

// MetricEventHandler defines the handler interface for MetricOverflowEvent.
// zh: MetricEventHandler 定義處理 MetricOverflowEvent 的事件處理器介面。
//
// This interface is used in the EventBus to register handlers for metric overflow conditions.
// zh: 本介面用於 EventBus 中註冊處理指標溢出事件的 handler。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type MetricEventHandler interface {
	// HandleMetricOverflow processes the MetricOverflowEvent.
	// zh: 處理監控指標溢出事件的實作函式。
	HandleMetricOverflow(ctx context.Context, event MetricOverflowEvent) error
}
