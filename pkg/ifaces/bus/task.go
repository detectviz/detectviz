package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// TaskDispatcher defines a dispatcher for task completed events.
// zh: TaskDispatcher 定義用於分派任務完成事件的介面。
type TaskDispatcher interface {
	// DispatchTaskCompleted sends a TaskCompletedEvent to registered handlers.
	// zh: 將 TaskCompletedEvent 傳遞給已註冊的處理器。
	DispatchTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error
}
