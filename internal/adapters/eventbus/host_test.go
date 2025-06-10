package eventbusadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

func TestHostLoggerHandler_HandleHostDiscovered(t *testing.T) {
	l := &fakes.FakeLogger{}
	ctx := ifacelogger.WithContext(context.Background(), l)
	h := &HostLoggerHandler{}
	err := h.HandleHostDiscovered(ctx, event.HostDiscoveredEvent{
		Name:   "host1",
		Source: "scanner",
		Labels: map[string]string{"env": "dev"},
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

func TestHostHandlerRegistry(t *testing.T) {
	hostHandlers = nil
	h := &HostLoggerHandler{}
	RegisterHostHandler(h)
	list := LoadPluginHostHandlers()
	if len(list) != 1 || list[0] != h {
		t.Fatalf("handler registry not working: %#v", list)
	}
}
