package notifier

import (
	"context"

	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// MultiNotifier 將多個 Notifier 組合成一個 Notifier。
// zh: 每個通知請求會依序傳遞至所有註冊的 Notifier。
type MultiNotifier struct {
	notifiers []notifieriface.Notifier
}

// NewMultiNotifier 建立 MultiNotifier 實例。
// zh: 可接受多個 Notifier 作為參數，並整合為一個執行單元。
func NewMultiNotifier(list ...notifieriface.Notifier) notifieriface.Notifier {
	return &MultiNotifier{notifiers: list}
}

// Notify 傳送通知，會依序呼叫每個註冊的 Notifier。
// zh: 若某些 Notifier 回傳錯誤，最終僅回傳最後一個錯誤。
func (m *MultiNotifier) Notify(title, message string) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Notify(title, message); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Send 傳送完整通知訊息，依序傳遞至每個 Notifier。
// zh: 支援包含標籤、時間等欄位的複雜訊息傳送。
func (m *MultiNotifier) Send(ctx context.Context, msg notifieriface.Message) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(ctx, msg); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Name 回傳 Notifier 名稱。
// zh: 表示此組合型通知器的識別名稱。
func (m *MultiNotifier) Name() string {
	return "multi"
}
