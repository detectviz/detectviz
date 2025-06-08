package eventbus

import "context"

// MetricEventHandler defines the handler interface for MetricOverflowEvent.
// zh: MetricEventHandler 定義處理 MetricOverflowEvent 的事件處理器介面。
type MetricEventHandler interface {
	HandleMetricOverflow(ctx context.Context, event MetricOverflowEvent) error
}
