package eventbusadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

func TestMetricLoggerHandler_HandleMetricOverflow(t *testing.T) {
	l := &fakes.FakeLogger{}
	ctx := ifacelogger.WithContext(context.Background(), l)
	h := &MetricLoggerHandler{}
	err := h.HandleMetricOverflow(ctx, event.MetricOverflowEvent{
		MetricName: "cpu",
		Value:      95,
		Threshold:  90,
		Instance:   "host1",
		Reason:     "high",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(l.Entries))
	}
	if l.Entries[0].Level != "WARN" {
		t.Errorf("expected WARN level, got %s", l.Entries[0].Level)
	}
}

func TestMetricHandlerRegistry(t *testing.T) {
	metricHandlers = nil
	h := &MetricLoggerHandler{}
	RegisterMetricHandler(h)
	list := LoadPluginMetricHandlers()
	if len(list) != 1 || list[0] != h {
		t.Fatalf("handler registry not working: %#v", list)
	}
}
