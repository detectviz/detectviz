package eventbusadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

func TestTaskLoggerHandler_HandleTaskCompleted(t *testing.T) {
	l := &fakes.FakeLogger{}
	ctx := ifacelogger.WithContext(context.Background(), l)
	h := &TaskLoggerHandler{}
	err := h.HandleTaskCompleted(ctx, event.TaskCompletedEvent{
		TaskID:   "id",
		WorkerID: "worker",
		Status:   "done",
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

func TestTaskHandlerRegistry(t *testing.T) {
	taskHandlers = nil
	h := &TaskLoggerHandler{}
	RegisterTaskHandler(h)
	list := LoadPluginTaskHandlers()
	if len(list) != 1 || list[0] != h {
		t.Fatalf("handler registry not working: %#v", list)
	}
}
