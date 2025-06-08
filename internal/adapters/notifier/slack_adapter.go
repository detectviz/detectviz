package notifier

import (
	"context"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// SlackNotifier sends notifications to Slack.
// zh: SlackNotifier 負責傳送通知至 Slack。
type SlackNotifier struct {
	name       string
	webhookURL string
	logger     ifacelogger.Logger
}

// NewSlackNotifier returns a new instance of SlackNotifier.
// zh: 建立新的 SlackNotifier 實例。
func NewSlackNotifier(name string, webhookURL string, logger ifacelogger.Logger) *SlackNotifier {
	return &SlackNotifier{
		name:       name,
		webhookURL: webhookURL,
		logger:     logger,
	}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *SlackNotifier) Name() string {
	return n.name
}

// Send sends the message to Slack.
// zh: 傳送通知訊息至 Slack（尚未實作實際 API 呼叫）。
func (n *SlackNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.logger.WithContext(ctx).Info("SlackNotifier sending",
		"to", msg.Target,
		"title", msg.Title,
		"content", msg.Content,
	)
	// TODO: implement real Slack webhook logic
	return nil
}

// Notify implements simplified notification with title and message.
// zh: 使用簡易標題與訊息格式發送 Slack 通知。
func (n *SlackNotifier) Notify(title, message string) error {
	// TODO: 決定預設或從設定注入 Target
	msg := ifacenotifier.Message{
		Target:  "default-slack-channel",
		Title:   title,
		Content: message,
	}
	return n.Send(context.Background(), msg)
}
