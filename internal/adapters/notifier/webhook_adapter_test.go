package notifieradapter

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestWebhookNotifier_Send(t *testing.T) {
	log := &fakes.FakeLogger{}
	var capturedURL string
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		capturedURL = req.URL.String()
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(nil))}, nil
	})}
	n := NewWebhookNotifier("hook", log, client)
	msg := notifieriface.Message{Target: "http://example.com", Title: "hi", Content: "msg"}
	err := n.Send(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if capturedURL != "http://example.com" {
		t.Errorf("unexpected URL: %s", capturedURL)
	}
	if len(log.Entries) == 0 {
		t.Error("expected log entry recorded")
	}
}

func TestWebhookNotifier_Notify(t *testing.T) {
	log := &fakes.FakeLogger{}
	client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(nil))}, nil
	})}
	n := NewWebhookNotifier("hook", log, client)
	if err := n.Notify("title", "content"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(log.Entries) == 0 {
		t.Error("expected log entry recorded")
	}
}
