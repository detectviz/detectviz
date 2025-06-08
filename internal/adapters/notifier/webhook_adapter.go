package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// WebhookNotifier sends notifications to a specified webhook URL.
// zh: WebhookNotifier 負責將通知訊息送往指定的 webhook URL。
type WebhookNotifier struct {
	name   string
	client *http.Client
	logger ifacelogger.Logger
}

// NewWebhookNotifier creates a new WebhookNotifier.
// zh: 建立新的 WebhookNotifier 實例。
// 若未提供 client 則使用 http.DefaultClient。
func NewWebhookNotifier(name string, logger ifacelogger.Logger, client *http.Client) *WebhookNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &WebhookNotifier{
		name:   name,
		client: client,
		logger: logger,
	}
}

// Name returns the notifier name.
// zh: 回傳 notifier 名稱。
func (n *WebhookNotifier) Name() string {
	return n.name
}

// Send transmits the message as a JSON payload via HTTP POST.
// zh: 傳送訊息為 JSON 格式，透過 HTTP POST 傳送至 msg.Target。
func (n *WebhookNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to marshal webhook message", "error", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, msg.Target, bytes.NewBuffer(payload))
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to create webhook request", "error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to send webhook request", "error", err)
		return err
	}
	defer resp.Body.Close()

	n.logger.WithContext(ctx).Info("WebhookNotifier sent", "target", msg.Target, "status", resp.Status)
	return nil
}

// Notify sends a simple title-message notification.
// zh: 傳送簡易的標題與訊息格式至 webhook，將組成 JSON 輸出。
func (n *WebhookNotifier) Notify(title, message string) error {
	// TODO: 後續改為可設定或多通道支援
	msg := ifacenotifier.Message{
		Target:  "http://localhost:8080/webhook",
		Title:   title,
		Content: message,
	}
	return n.Send(context.Background(), msg)
}
