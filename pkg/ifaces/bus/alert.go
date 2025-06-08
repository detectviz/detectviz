package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// AlertDispatcher defines a dispatcher for alert events.
// zh: AlertDispatcher 定義用於分派告警事件的介面。
type AlertDispatcher interface {
	// DispatchAlert sends an AlertTriggeredEvent to registered handlers.
	// zh: 將 AlertTriggeredEvent 傳遞給已註冊的處理器。
	DispatchAlert(ctx context.Context, event event.AlertTriggeredEvent) error
}
