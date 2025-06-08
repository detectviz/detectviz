package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// TaskLoggerHandler is a sample implementation of TaskEventHandler that logs task results.
// zh: TaskLoggerHandler 是接收到任務完成事件時記錄日誌的處理器實作範例。
type TaskLoggerHandler struct{}

// HandleTaskCompleted handles task completion and logs the result.
// zh: 接收到任務完成事件後，透過 logger 模組輸出結構化日誌。
func (h *TaskLoggerHandler) HandleTaskCompleted(ctx context.Context, event eventbus.TaskCompletedEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"task_id":   event.TaskID,
		"worker_id": event.WorkerID,
		"status":    event.Status,
	}).Info("[TASK] completed")
	return nil
}
