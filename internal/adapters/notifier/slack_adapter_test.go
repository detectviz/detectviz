package notifieradapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

func TestSlackNotifier_Send(t *testing.T) {
	log := &fakes.FakeLogger{}
	n := NewSlackNotifier("slack", "https://example.com", log)
	err := n.Send(context.Background(), notifieriface.Message{Target: "https://hooks.slack", Title: "hi", Content: "msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) == 0 {
		t.Error("expected log entry recorded")
	}
	if n.Name() != "slack" {
		t.Errorf("unexpected Name: %s", n.Name())
	}
}

func TestSlackNotifier_Notify(t *testing.T) {
	log := &fakes.FakeLogger{}
	n := NewSlackNotifier("slack", "https://example.com", log)
	if err := n.Notify("title", "content"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) == 0 {
		t.Error("expected log entry recorded")
	}
}
