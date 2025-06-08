package eventbus

import "context"

// AlertEventHandler defines the handler interface for AlertTriggeredEvent.
// zh: AlertEventHandler 定義處理 AlertTriggeredEvent 的介面。
type AlertEventHandler interface {
	HandleAlertTriggered(ctx context.Context, event AlertTriggeredEvent) error
}
