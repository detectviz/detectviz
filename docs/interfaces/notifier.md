# Notifier Interface 設計說明

> 本文件說明 `pkg/ifaces/notifier` 中的通知模組介面設計、使用情境與實作結構。Notifier 模組是 detectviz 的通用訊息推送介面，負責統一封裝告警與事件訊息，並透過 Email、Slack、Webhook 等方式傳遞至外部。

---

## 設計目標（Design Goals）

- 抽象化通知機制，支援多種傳送通道（如 email, webhook, slack）
- 支援訊息格式、通知等級、推送目標等欄位封裝
- 可由 alert 模組、scheduler 模組、eventbus 模組觸發使用
- 易於擴充與動態註冊
- 搭配 logger 記錄成功與錯誤資訊

---

## Interface 定義（Interface Definitions）

```go
type Notifier interface {
	Send(ctx context.Context, msg Message) error
}

type Message struct {
	Title   string            // 通知標題
	Content string            // 通知內容
	Level   string            // 分類：info / warning / critical
	Target  string            // 接收對象，例如 webhook URL、email
	Time    time.Time         // 發送時間
}
```

---

## 使用情境（When and where it's used）

- 告警觸發時由 AlertEvaluator 推送通知
- 任務完成後由 Scheduler 發送執行結果
- 系統異常時透過 EventBus 發送通知
- 可於 plugin 中擴充事件通知行為

---

## 實作位置與類型（Implementations）

| 名稱             | 檔案路徑                                               | 描述                         |
|------------------|--------------------------------------------------------|------------------------------|
| EmailNotifier     | `internal/adapters/notifier/email_adapter.go`          | 寄送 email 通知               |
| SlackNotifier     | `internal/adapters/notifier/slack_adapter.go`          | 發送 Slack 訊息               |
| WebhookNotifier   | `internal/adapters/notifier/webhook_adapter.go`        | 送出 HTTP POST 通知          |
| MockNotifier      | `internal/adapters/notifier/mock_adapter.go`           | 單元測試用，不實際發送        |
| MultiNotifier     | `internal/adapters/notifier/multi.go`                  | 將一則訊息傳送給多個 Notifier |
| NopNotifier       | `internal/adapters/notifier/nop.go`                    | 無動作通知器（開發或測試使用）|

---

## 設定與註冊方式（Configuration & Registration）

- 設定來源：`pkg/configtypes/notifier_config.go`
- 動態註冊位置：`internal/registry/notifier/registry.go`
- 支援多個 notifier 並以 `[]configtypes.NotifierConfig` 批次註冊
- 可注入 logger 與自定義 http client

---

## 測試結構（Testing Structure）

| 測試檔案位置                                         | 測試目標                        |
|------------------------------------------------------|---------------------------------|
| `mock_adapter_test.go`                               | 驗證介面契約與訊息結構          |
| `email_adapter_test.go`（預定）                      | 測試 email 發送與錯誤處理       |
| `slack_adapter_test.go`（預定）                      | 測試 Slack 訊息格式與傳遞邏輯   |
| `webhook_adapter_test.go`（預定）                   | 測試 HTTP POST 傳送邏輯         |
| `multi_test.go` / `nop_test.go`（補充中）            | 測試組合與靜態 fallback 行為    |

---

## 擴充方式（How to add a new Notifier）

1. 建立檔案於 `internal/adapters/notifier/{name}_adapter.go`
2. 實作 `pkg/ifaces/notifier.Notifier` 介面
3. 可搭配自定義 logger、HTTP client、retry 邏輯
4. 補上單元測試：`{name}_adapter_test.go`
5. 註冊至 `internal/registry/notifier/registry.go`
6. 擴充 config 設定於 `pkg/configtypes/notifier_config.go`

---

## 相關依賴模組（Related Modules）

- `logger`：記錄通知發送與錯誤
- `eventbus`：事件觸發來源
- `alert` / `scheduler`：主要呼叫來源
- `configtypes.NotifierConfig`：指定通道參數與開關

---