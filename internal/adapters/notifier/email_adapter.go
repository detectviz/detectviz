package notifier

import (
	"context"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// EmailNotifier sends notifications via email.
// zh: EmailNotifier 負責透過電子郵件傳送通知。
type EmailNotifier struct {
	name   string
	sender string
	logger ifacelogger.Logger
}

// NewEmailNotifier returns a new instance of EmailNotifier.
// zh: 建立新的 EmailNotifier 實例。
func NewEmailNotifier(name string, sender string, logger ifacelogger.Logger) *EmailNotifier {
	return &EmailNotifier{
		name:   name,
		sender: sender,
		logger: logger,
	}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *EmailNotifier) Name() string {
	return n.name
}

// Send sends the message as an email.
// zh: 傳送通知訊息為 email（尚未實作寄信邏輯）。
func (n *EmailNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.logger.WithContext(ctx).Info("EmailNotifier sending",
		"to", msg.Target,
		"title", msg.Title,
		"content", msg.Content,
	)
	// TODO: implement real email sending logic using SMTP or third-party API
	return nil
}
