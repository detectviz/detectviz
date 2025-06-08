package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// MetricLoggerHandler is a sample implementation of MetricEventHandler that logs overflow events.
// zh: MetricLoggerHandler 是接收到指標溢出事件時記錄日誌的處理器實作範例。
type MetricLoggerHandler struct{}

// HandleMetricOverflow handles MetricOverflowEvent by logging the event with structured fields.
// zh: 接收到指標溢出事件後，使用結構化欄位輸出警告日誌
func (h *MetricLoggerHandler) HandleMetricOverflow(ctx context.Context, event eventbus.MetricOverflowEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"metric":    event.MetricName,
		"value":     event.Value,
		"threshold": event.Threshold,
		"instance":  event.Instance,
		"reason":    event.Reason,
	}).Warn("[METRIC] overflow detected")
	return nil
}
