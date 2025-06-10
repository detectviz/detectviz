


# Services Architecture

本文件說明 detectviz 專案中 `internal/services/` 模組的設計邏輯、責任分工與擴充策略，定義 Service 作為業務邏輯處理層的角色與實作方式。

---

## 模組定位與目的

Service 層負責處理整體業務邏輯，為 API handler 與資料層（store/plugin）之間的中介層，實現以下目標：

- 將核心邏輯與輸入輸出處理解耦
- 支援注入測試用依賴物件（如 fake store）
- 管理不同儲存後端與處理流程邏輯
- 提供跨模組調度能力（ex: Rule 驗證後觸發 Notifier）

---

## 對應目錄結構建議

```
internal/services/
├── rule/
│   └── service.go
├── notifier/
│   └── service.go
├── logger/
│   └── service.go
├── metrics/
│   └── service.go
├── eventbus/
│   └── service.go
└── interfaces.go       # 定義跨模組 Service 介面
```

---

## Service 與其他層責任分工

| 層級     | 職責說明 |
|----------|----------|
| handler  | 處理輸入、驗證、回傳格式，調用 service |
| service  | 執行商業邏輯、呼叫 store、調度 plugin |
| store    | 執行資料 CRUD / 查詢等操作 |
| plugin   | 注入擴充元件，如多儲存後端、通知處理等 |

---

## Service 實作建議

每個 service 結構體建議注入下列組件：

```go
type RuleService struct {
    Store    store.RuleStore
    Notifier notifier.Service
    Logger   logger.Service
    Metrics  metrics.Service
}
```

---

## 測試建議

- 每個 service 可替換注入的 store / plugin 實作進行單元測試
- 可搭配 fake store 實作進行邏輯路徑覆蓋
- 不依賴 handler 層進行驗證測試

---

## 與其他模組整合

- handler 調用對應模組的 service 方法
- service 調用 store 並依情境調度 plugin（如通知發送）
- plugin 可由 service 注入並依條件觸發
- eventbus 可於 service 中觸發事件通知

---

## 擴充方向建議

- 支援 decorator pattern（如 AuditService, TracingService）
- 可將流程拆分為 UseCase 類別以支援更複雜邏輯
- service 與 middleware 可搭配，執行權限或租戶篩選

---