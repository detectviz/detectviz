package alertlog

import (
	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// 確保 AlertPluginHandler 實作 AlertEventHandler
var _ event.AlertEventHandler = (*AlertPluginHandler)(nil)

// init 初始化 alertlog plugin 並註冊其事件處理器
// zh: 在載入此模組時自動註冊 AlertPluginHandler 作為告警事件處理器。
func init() {
	eventbusadapter.RegisterAlertHandler(&AlertPluginHandler{})
}
