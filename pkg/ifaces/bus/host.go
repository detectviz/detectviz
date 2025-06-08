package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// HostDispatcher defines a dispatcher for host discovery events.
// zh: HostDispatcher 定義用於分派主機發現事件的介面。
type HostDispatcher interface {
	// DispatchHostDiscovered sends a HostDiscoveredEvent to registered handlers.
	// zh: 將 HostDiscoveredEvent 傳遞給已註冊的處理器。
	DispatchHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error
}
