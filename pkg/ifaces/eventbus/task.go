package eventbus

import "context"

// TaskEventHandler defines the handler interface for TaskCompletedEvent.
// zh: TaskEventHandler 定義處理 TaskCompletedEvent 的介面。
type TaskEventHandler interface {
	HandleTaskCompleted(ctx context.Context, event TaskCompletedEvent) error
}
