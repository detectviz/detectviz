package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// HostLoggerHandler is a sample implementation of HostEventHandler that logs discovered hosts.
// zh: HostLoggerHandler 是接收到主機註冊事件時記錄日誌的處理器實作範例。
type HostLoggerHandler struct{}

// HandleHostDiscovered handles HostDiscoveredEvent and logs host information.
// zh: 接收到主機註冊事件後，透過 logger 模組輸出結構化日誌。
func (h *HostLoggerHandler) HandleHostDiscovered(ctx context.Context, event eventbus.HostDiscoveredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"name":   event.Name,
		"source": event.Source,
		"labels": event.Labels,
	}).Info("[HOST] discovered")
	return nil
}
