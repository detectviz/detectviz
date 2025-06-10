package eventbusadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

func TestAlertLoggerHandler_HandleAlertTriggered(t *testing.T) {
	l := &fakes.FakeLogger{}
	ctx := ifacelogger.WithContext(context.Background(), l)
	h := &AlertLoggerHandler{}
	err := h.HandleAlertTriggered(ctx, event.AlertTriggeredEvent{
		AlertID:    "a",
		RuleName:   "rule",
		Level:      "critical",
		Instance:   "inst",
		Metric:     "m",
		Comparison: "<",
		Value:      1,
		Threshold:  2,
		Message:    "msg",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(l.Entries) != 1 {
		t.Fatalf("expected 1 log entry, got %d", len(l.Entries))
	}
	if l.Entries[0].Level != "INFO" {
		t.Errorf("expected INFO level, got %s", l.Entries[0].Level)
	}
}

func TestAlertHandlerRegistry(t *testing.T) {
	alertHandlers = nil
	h := &AlertLoggerHandler{}
	RegisterAlertHandler(h)
	list := LoadPluginAlertHandlers()
	if len(list) != 1 || list[0] != h {
		t.Fatalf("handler registry not working: %#v", list)
	}
}
