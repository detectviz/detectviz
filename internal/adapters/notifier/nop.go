package notifier

import (
	"context"

	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// NopNotifier 是一個不執行任何操作的通知器。
// zh: 常用於測試或當無通知功能需求時作為 fallback 實作。
type NopNotifier struct{}

// NewNopNotifier 建立 NopNotifier 實例。
// zh: 回傳不執行通知邏輯的 Notifier 實作。
func NewNopNotifier() notifieriface.Notifier {
	return &NopNotifier{}
}

// Name 回傳通知器名稱。
// zh: 標示此 notifier 為 nop 類型。
func (n *NopNotifier) Name() string {
	return "nop"
}

// Send 實作完整通知方法但不執行任何動作。
// zh: 忽略訊息內容並回傳 nil，適用於測試或預設無通報。
func (n *NopNotifier) Send(ctx context.Context, msg notifieriface.Message) error {
	return nil
}

// Notify 實作簡易通知方法但不執行任何動作。
// zh: 忽略標題與訊息並回傳 nil，適用於測試或不需通知的情境。
func (n *NopNotifier) Notify(title, message string) error {
	return nil
}
