package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// MetricDispatcher defines a dispatcher for metric overflow events.
// zh: MetricDispatcher 定義用於分派監控指標溢出事件的介面。
type MetricDispatcher interface {
	// DispatchMetricOverflow sends a MetricOverflowEvent to registered handlers.
	// zh: 將 MetricOverflowEvent 傳遞給已註冊的處理器。
	DispatchMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error
}
