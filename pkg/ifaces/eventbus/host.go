package eventbus

import "context"

// HostEventHandler defines the handler interface for HostDiscoveredEvent.
// zh: HostEventHandler 定義處理 HostDiscoveredEvent 的事件處理器介面。
type HostEventHandler interface {
	HandleHostDiscovered(ctx context.Context, event HostDiscoveredEvent) error
}
