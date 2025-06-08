package event

import "time"

// AlertTriggeredEvent represents an alert event that occurred in the system.
// zh: AlertTriggeredEvent 表示系統中觸發的告警事件。
type AlertTriggeredEvent struct {
	EventID    string    // zh: 事件唯一識別碼
	Timestamp  time.Time // zh: 事件發生時間
	AlertID    string    // zh: 告警事件 ID
	RuleName   string    // zh: 告警規則名稱
	Level      string    // zh: 告警嚴重程度，例如 critical、warning
	Instance   string    // zh: 實體名稱或識別，例如設備名稱
	Metric     string    // zh: 指標名稱
	Comparison string    // zh: 比較運算符，例如 >、<
	Value      float64   // zh: 實際值
	Threshold  float64   // zh: 閾值
	Message    string    // zh: 告警訊息內容
}

// TaskCompletedEvent represents a task completion event in the system.
// zh: TaskCompletedEvent 表示系統中某個任務完成的事件。
type TaskCompletedEvent struct {
	EventID   string    // zh: 事件唯一識別碼
	Timestamp time.Time // zh: 事件發生時間
	TaskID    string    // zh: 任務識別碼
	WorkerID  string    // zh: 執行任務的工作者 ID
	Status    string    // zh: 任務完成狀態，例如 success、failed
}

// HostDiscoveredEvent represents the discovery or registration of a host in the system.
// zh: HostDiscoveredEvent 表示系統中主機被發現或註冊的事件。
type HostDiscoveredEvent struct {
	EventID   string            // zh: 事件唯一識別碼
	Timestamp time.Time         // zh: 事件發生時間
	HostID    string            // zh: 主機識別碼
	Name      string            // zh: 主機名稱
	IP        string            // zh: 主機 IP 位址
	Source    string            // zh: 來源識別，例如由哪個掃描器或子系統發現
	Labels    map[string]string // zh: 附加標籤，例如 rack、zone 等
}

// MetricOverflowEvent represents an overflow condition in a monitored metric.
// zh: MetricOverflowEvent 表示某個監控指標超出預期範圍的事件。
type MetricOverflowEvent struct {
	EventID    string    // zh: 事件唯一識別碼
	Timestamp  time.Time // zh: 事件發生時間
	MetricName string    // zh: 指標名稱
	Value      float64   // zh: 實際值
	Threshold  float64   // zh: 閾值
	Source     string    // zh: 數據來源（原欄位，保留以維持相容）
	Instance   string    // zh: 實體名稱或識別，例如設備名稱
	Reason     string    // zh: 溢出的原因說明
}
