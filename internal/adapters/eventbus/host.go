package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// HostLoggerHandler is a sample implementation of HostEventHandler that logs discovered hosts.
// zh: HostLoggerHandler 是接收到主機註冊事件時記錄日誌的處理器實作範例。
type HostLoggerHandler struct{}

// HandleHostDiscovered handles HostDiscoveredEvent and logs host information.
// zh: 接收到主機註冊事件後，透過 logger 模組輸出結構化日誌。
func (h *HostLoggerHandler) HandleHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"name":   event.Name,
		"source": event.Source,
		"labels": event.Labels,
	}).Info("[HOST] discovered")
	return nil
}

// hostHandlers 是所有已註冊的 HostEventHandler 清單
var hostHandlers []event.HostEventHandler

// RegisterHostHandler 用於讓 plugin 模組註冊自訂的主機事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterHostHandler(handler event.HostEventHandler) {
	hostHandlers = append(hostHandlers, handler)
}

// LoadPluginHostHandlers 回傳目前已註冊的 plugin HostEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 host handler。
func LoadPluginHostHandlers() []event.HostEventHandler {
	return hostHandlers
}
