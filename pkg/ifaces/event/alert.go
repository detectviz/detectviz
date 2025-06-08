package event

import "context"

// AlertEventHandler defines the handler interface for AlertTriggeredEvent.
// zh: AlertEventHandler 定義處理 AlertTriggeredEvent 的介面。
//
// Used in event dispatch systems to register alert-related handlers.
// zh: 用於事件分派系統中註冊處理告警事件的 handler。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type AlertEventHandler interface {
	// HandleAlertTriggered processes the AlertTriggeredEvent.
	// zh: 處理告警事件的實作函式。
	HandleAlertTriggered(ctx context.Context, event AlertTriggeredEvent) error
}
