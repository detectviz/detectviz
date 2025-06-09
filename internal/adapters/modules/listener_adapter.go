package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
)

// ListenerAdapter wraps core.Listener to implement iface.ModuleListener.
// zh: ListenerAdapter 包裝 core.Listener，使其實作 ModuleListener 介面。
type ListenerAdapter struct {
	listener *core.Listener
}

// NewListenerAdapter constructs a new ListenerAdapter instance.
// zh: 建立新的 ListenerAdapter 實例。
func NewListenerAdapter(l *core.Listener) *ListenerAdapter {
	return &ListenerAdapter{
		listener: l,
	}
}

// Start begins the health monitoring loop.
// zh: 啟動健康狀態監控迴圈。
func (a *ListenerAdapter) Start(ctx context.Context) {
	a.listener.Start(ctx)
}

// Stop stops the health monitoring listener.
// zh: 停止健康狀態監聽器。
func (a *ListenerAdapter) Stop() {
	a.listener.Stop()
}
