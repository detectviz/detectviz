package event

import "context"

// HostEventHandler defines the handler interface for HostDiscoveredEvent.
// zh: HostEventHandler 定義處理 HostDiscoveredEvent 的事件處理器介面。
//
// This interface is typically registered to an EventBus dispatcher.
// zh: 本介面通常會註冊至 EventBus 分派器以接收主機註冊事件。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type HostEventHandler interface {
	// HandleHostDiscovered processes the HostDiscoveredEvent.
	// zh: 處理主機註冊事件的實作函式。
	HandleHostDiscovered(ctx context.Context, event HostDiscoveredEvent) error
}
