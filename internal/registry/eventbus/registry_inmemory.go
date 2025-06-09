package eventbus

import (
	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

type inMemoryProvider struct{}

func (p *inMemoryProvider) Name() string {
	return "in-memory"
}

func (p *inMemoryProvider) Build() eventbusiface.EventDispatcher {
	dispatcher := eventbusadapter.NewInMemoryDispatcher()

	for _, h := range eventbusadapter.LoadPluginAlertHandlers() {
		dispatcher.RegisterAlertHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginHostHandlers() {
		dispatcher.RegisterHostHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginMetricHandlers() {
		dispatcher.RegisterMetricHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginTaskHandlers() {
		dispatcher.RegisterTaskHandler(h)
	}

	return dispatcher
}

func init() {
	RegisterProvider(&inMemoryProvider{})
}
