package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// TaskLoggerHandler is a sample implementation of TaskEventHandler that logs task results.
// zh: TaskLoggerHandler 是接收到任務完成事件時記錄日誌的處理器實作範例。
type TaskLoggerHandler struct{}

// HandleTaskCompleted handles task completion and logs the result.
// zh: 接收到任務完成事件後，透過 logger 模組輸出結構化日誌。
func (h *TaskLoggerHandler) HandleTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"task_id":   event.TaskID,
		"worker_id": event.WorkerID,
		"status":    event.Status,
	}).Info("[TASK] completed")
	return nil
}

// taskHandlers 是所有已註冊的 TaskEventHandler 清單
var taskHandlers []event.TaskEventHandler

// RegisterTaskHandler 用於讓 plugin 模組註冊自訂的任務事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterTaskHandler(handler event.TaskEventHandler) {
	taskHandlers = append(taskHandlers, handler)
}

// LoadPluginTaskHandlers 回傳目前已註冊的 plugin TaskEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 task handler。
func LoadPluginTaskHandlers() []event.TaskEventHandler {
	return taskHandlers
}
