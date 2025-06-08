package configtypes

// NotifierConfig defines the configuration for a notifier.
// zh: NotifierConfig 定義單一通知通道的設定結構。
type NotifierConfig struct {
	Name   string `json:"name"`   // zh: 通道名稱（email, slack, webhook 等）
	Type   string `json:"type"`   // zh: 通道類型
	Target string `json:"target"` // zh: 傳送目標（例如 email address、webhook URL）
	Enable bool   `json:"enable"` // zh: 是否啟用此通道
}
