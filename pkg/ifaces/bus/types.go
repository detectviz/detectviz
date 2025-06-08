package bus

import "time"

// Envelope defines a wrapper for dispatched events.
// zh: Envelope 用來包裝傳遞的事件與其附加資訊，供 Dispatcher 使用。
type Envelope struct {
	EventType string      // zh: 事件類型，例如 "host.discovered"
	Payload   interface{} // zh: 真實事件資料內容，可為任意型別
	Timestamp time.Time   // zh: 發送時間，用於排序或延遲處理用途
}
