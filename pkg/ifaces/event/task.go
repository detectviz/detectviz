package event

import "context"

// TaskEventHandler defines the handler interface for TaskCompletedEvent.
// zh: TaskEventHandler 定義處理 TaskCompletedEvent 的介面。
//
// This interface is used to handle task execution completion events via the EventBus dispatcher.
// zh: 本介面用於透過 EventBus 分派器處理任務執行完成事件。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type TaskEventHandler interface {
	// HandleTaskCompleted processes the TaskCompletedEvent.
	// zh: 處理任務完成事件的實作函式。
	HandleTaskCompleted(ctx context.Context, event TaskCompletedEvent) error
}
