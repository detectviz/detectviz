# {Interface Name} Interface 說明文件

> 本文件為 Detectviz 專案中 `{Interface Name}` 介面的設計說明與使用情境整理。
> 
> 請將 `{Interface Name}` 替換為實際名稱，並填入每段描述內容。

---

## 介面用途（What it does）

說明此 interface 的抽象職責與在架構中的角色，例如：

- 提供設定存取抽象層
- 支援模組間的訊息廣播
- 統一快取存取與替換機制
- 與 OpenTelemetry 整合之觀測能力抽象

---

## 使用情境（When and where it's used）

列出具體模組或流程中會用到此介面的情況，例如：

- 在 `bootstrap.Init()` 中注入供 core 使用
- 被 HTTP middleware 呼叫以取得 trace context
- 在任務排程器中使用以避免重複執行

---

## 方法說明（Methods）

針對每個方法補充說明其功能與回傳值意圖，例如：

- `Get(key string) string`：回傳對應設定值，若不存在可能為空字串
- `WithContext(ctx)`：將 context 內 trace id 注入 logger 流程
- `Publish(topic string, payload any)`：廣播一筆訊息給訂閱者

---

## 預期實作（Expected implementations）

列出你預期會有哪些實作，以及建議放置目錄與用途：

- `internal/adapters/logger/zaplogger.go`
- `internal/adapters/cachestore/memcache.go`
- `internal/adapters/eventbus/inmemory.go`
- `internal/adapters/{module}/nop.go`（測試用）

---

## 關聯模組與擴充性（Related & extensibility）

如有與其他模組或 interface 的整合需求，或未來可擴充方向，請說明：

- 與 trace / metrics 的整合
- Redis/NATS 替代支援
- 可 plugin 化設計

---