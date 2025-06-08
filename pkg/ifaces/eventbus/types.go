package eventbus

// AlertTriggeredEvent represents an alert event that occurred in the system.
// zh: AlertTriggeredEvent 表示系統中觸發的告警事件。
type AlertTriggeredEvent struct {
	RuleID   string
	Severity string
	Message  string
}

// TaskCompletedEvent represents a task completion event in the system.
// zh: TaskCompletedEvent 表示系統中某個任務完成的事件。
type TaskCompletedEvent struct {
	TaskID   string
	WorkerID string
	Status   string
}

// HostDiscoveredEvent represents the discovery or registration of a host in the system.
// zh: HostDiscoveredEvent 表示系統中主機被發現或註冊的事件。
type HostDiscoveredEvent struct {
	HostID string
	Name   string
	IP     string
	Source string            // zh: 來源識別，例如由哪個掃描器或子系統發現
	Labels map[string]string // zh: 附加標籤，例如 rack、zone 等
}

// MetricOverflowEvent represents an overflow condition in a monitored metric.
// zh: MetricOverflowEvent 表示某個監控指標超出預期範圍的事件。
type MetricOverflowEvent struct {
	MetricName string
	Value      float64
	Threshold  float64
	Source     string // zh: 原欄位，保留以維持相容
	Instance   string // zh: 實體名稱或識別，例如設備名稱
	Reason     string // zh: 溢出的原因說明
}
