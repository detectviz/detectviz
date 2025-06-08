package notifier

import (
	"context"
	"time"
)

// Message represents a notification message to be sent.
// zh: Message 表示要傳送的通知訊息。
type Message struct {
	Title   string            // zh: 訊息標題
	Content string            // zh: 訊息內容
	Labels  map[string]string // zh: 附加標籤（例如等級、來源模組）
	Target  string            // zh: 接收對象或通道，例如 email、webhook URL
	Time    time.Time         // zh: 發送時間，預設為訊息產生時間
}

// Notifier defines an interface for sending notifications to various channels.
// zh: Notifier 定義通知傳送的介面，可擴充為多種通道實作（如 Email、Slack、Webhook）。
type Notifier interface {
	// Name returns the identifier of the notifier.
	// zh: 回傳 notifier 名稱。
	Name() string

	// Send delivers the message via this notifier channel.
	// zh: 傳送訊息至此通道。
	Send(ctx context.Context, msg Message) error
}
