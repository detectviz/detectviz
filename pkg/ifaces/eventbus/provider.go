package eventbus

// DispatcherProvider 為 eventbus 後端實作註冊器
type DispatcherProvider interface {
	// Name 傳回 provider 名稱（如 in-memory, kafka）
	Name() string

	// Build 建構對應的 EventDispatcher
	Build() EventDispatcher
}
