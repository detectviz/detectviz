package notifieradapter

import (
	"context"
	"sync"

	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// MockNotifier is a test implementation of the Notifier interface.
// zh: MockNotifier 為測試用通知器，會記錄所有發送的訊息。
type MockNotifier struct {
	name     string
	messages []ifacenotifier.Message
	mu       sync.Mutex
}

// NewMockNotifier creates a new MockNotifier.
// zh: 建立新的 MockNotifier 實例。
func NewMockNotifier(name string) *MockNotifier {
	return &MockNotifier{name: name}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *MockNotifier) Name() string {
	return n.name
}

// Send appends the message to internal buffer.
// zh: 將通知訊息儲存至內部緩衝區，供測試驗證。
func (n *MockNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.messages = append(n.messages, msg)
	return nil
}

// Messages returns all messages sent.
// zh: 回傳所有已發送訊息。
func (n *MockNotifier) Messages() []ifacenotifier.Message {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]ifacenotifier.Message(nil), n.messages...)
}
