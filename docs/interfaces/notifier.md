# Notifier Interface 設計說明

本文件說明 `pkg/ifaces/notifier` 中的通知模組介面設計與使用方式，適用於告警、任務完成、錯誤推送等場景。

---

## 設計目標（Design Goals）

- 抽象化通知機制，可對接多種傳送通道
- 支援通知格式、封裝欄位、等級分類
- 可於告警事件或排程任務中呼叫
- 易於擴充：如 email, webhook, Slack, Line 等
- 與 alert 與 eventbus 模組整合

---

## Interface 定義（Interface Definitions）

```go
type Notifier interface {
	Send(ctx context.Context, msg Message) error
}

type Message struct {
	Title   string
	Content string
	Level   string // info, warning, critical
	Target  string // 接收對象或管道，例如 webhook URL、email address
	Time    time.Time
}
```

---

## 使用情境（Usage）

適用於以下場景：

- 告警觸發時發送通知（搭配 alert evaluator）
- 任務完成後回報通知（如定時轉拋任務）
- 系統異常通知管理員
- 整合 eventbus 以即時推送事件

---

## 實作位置（Implementations）

| 類型                | 路徑                                              | 描述                             |
|---------------------|---------------------------------------------------|----------------------------------|
| EmailNotifier        | `internal/adapters/notifier/email_adapter.go`     | 發送通知至指定 email 地址        |
| SlackNotifier        | `internal/adapters/notifier/slack_adapter.go`     | 發送通知至 Slack 頻道            |
| WebhookNotifier      | `internal/adapters/notifier/webhook_adapter.go`   | 發送 HTTP POST 至指定網址        |
| MockNotifier         | `internal/adapters/notifier/mock_adapter.go`      | 測試用，記錄訊息不實際發送       |

---

## 測試結構（Testing Structure）

| 測試檔案位置                                            | 測試目標                          |
|---------------------------------------------------------|-----------------------------------|
| `mock_adapter_test.go`                                  | 確認訊息格式與 interface 符合     |
| `email_adapter_test.go`（未來）                         | 測試 email 發送邏輯與錯誤處理     |
| `slack_adapter_test.go`（未來）                         | 驗證 Slack 通知格式與傳送行為    |
| `webhook_adapter_test.go`（未來）                      | 測試 webhook 傳送與錯誤處理       |

---

## 擴充方式（Extension Tips）

若要新增通知通道：

1. 實作 `pkg/ifaces/notifier.Notifier` 介面
2. 建立對應 adapter，例如 `internal/adapters/notifier/telegram_adapter.go`
3. 補上錯誤處理與測試邏輯
4. 加入註冊邏輯於 `internal/bootstrap/notifier_registry.go`，支援靜態註冊或從 `[]config.NotifierConfig` 動態載入
5. 如需設定參數，擴充 `pkg/configtypes.NotifierConfig`

---

## 相關依賴模組

- logger：用於記錄通知成功與錯誤訊息
- alert / eventbus：可能的通知來源
- config：設定每個通知通道的憑證與網址

---