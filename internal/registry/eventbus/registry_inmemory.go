package eventbus

import (
	"github.com/detectviz/detectviz/internal/adapters/eventbus"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

type inMemoryProvider struct{}

func (p *inMemoryProvider) Name() string {
	return "in-memory"
}

func (p *inMemoryProvider) Build() eventbusiface.EventDispatcher {
	dispatcher := eventbus.NewInMemoryDispatcher()

	for _, h := range eventbus.LoadPluginAlertHandlers() {
		dispatcher.RegisterAlertHandler(h)
	}
	for _, h := range eventbus.LoadPluginHostHandlers() {
		dispatcher.RegisterHostHandler(h)
	}
	for _, h := range eventbus.LoadPluginMetricHandlers() {
		dispatcher.RegisterMetricHandler(h)
	}
	for _, h := range eventbus.LoadPluginTaskHandlers() {
		dispatcher.RegisterTaskHandler(h)
	}

	return dispatcher
}

func init() {
	RegisterProvider(&inMemoryProvider{})
}
