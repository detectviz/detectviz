# Store Architecture

本文件說明 detectviz 中 `internal/store/` 模組的架構設計，包含抽象介面、plugin 擴充點、資料來源統一調用邏輯，以及與 registry, services 的整合策略。

---

- 對應介面細節請參考：docs/interfaces/store.md

## 模組目標

- 抽象資料存取邏輯（CRUD / 查詢）
- 支援多種資料來源（如 MySQL, InfluxDB, Loki, File, Redis）
- 提供統一介面供 service 層使用
- 各資料來源可由 plugin 實作並注入
- 支援 resolver 決定具體實作來源（可依條件動態切換）
- 支援多模組（如 Rule, EventBus, Notifier, Logger, Metrics）使用 store 抽象

---

## 建議目錄結構

```
internal/store/
├── rule/
│   ├── memory/rule_store.go
│   ├── cache/rule_store.go
│   ├── mysql/rule_store.go
│   ├── influxdb/rule_store.go
│   └── logfile/rule_store.go
├── eventbus/
│   └── memory/eventbus_store.go
├── notifier/
│   └── mysql/notifier_store.go
├── logger/
│   └── logfile/logger_store.go
├── metrics/
│   └── influxdb/metrics_store.go
├── interfaces.go         # 可選：統一暴露全部模組介面（或分模組定義）
├── resolver.go           # 註冊並解析各模組的具體 store 實作
├── models/               # 共用資料結構（非 DB entity）
│   └── rule.go
```

---

## Interface 設計範例

```go
type RuleStore interface {
    Create(ctx context.Context, r *model.Rule) error
    Update(ctx context.Context, r *model.Rule) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (*model.Rule, error)
    List(ctx context.Context, filter RuleFilter) ([]*model.Rule, error)
}
```

---

## 各模組介面補充說明

以下列出目前支援模組與其介面定義的設計建議，詳細 interface 定義請見 [docs/interfaces/store.md](../interfaces/store.md)。

### RuleStore

- 提供 CRUD 與啟用/停用控制
- 支援快取 read-through 策略

### NotifierStore

- 訂閱者註冊與條件查詢
- 多目標與類型查詢支援

### LoggerStore

- 僅支援 Append 與條件查詢
- 可依時間範圍清除日誌

### MetricsStore

- 寫入時序資料點（InfluxDB 為主）
- 支援時間範圍查詢與 measurement 掃描

### EventBusStore

- 提供簡化的 Publish/Subscribe 機制
- 適用於本地測試與訊息模擬

---

## Plugin 實作與註冊

plugin 實作應放於：

```
plugins/datasources/rulestore_influxdb/impl.go
plugins/datasources/notifierstore_mysql/impl.go
plugins/datasources/loggerstore_logfile/impl.go
```

註冊點放在：

```
internal/registry/store/registry.go
```

使用方式：

```go
store.RegisterRuleStore("influx", influxImpl)
```

---

## Resolver 邏輯

- 可透過 config 或使用情境選擇具體實作
- 具備 fallback 與多來源支援能力
- 呼叫端僅依賴介面，不關心實作細節

---

## CacheStore / Redis 支援

- `internal/store/cachestore` 為 Redis 專用快取封裝，可作為獨立存取元件或輔助主 store
- 通常搭配 RuleStore 實作為「read-through」快取策略
- Redis plugin 可於下列位置實作：

```
plugins/datasources/rulestore_redis/impl.go
```

- 支援行為：
  - 查詢先查快取，miss 時回源並寫入快取
  - 可針對頻繁讀取的 Rule 查詢加速

- 設計原則：
  - cache 層不應做資料權限驗證與邏輯處理
  - 可加入 TTL 控制與 eviction 訊號處理

---

## 測試支援

- 提供 `memory` fake store 實作，用於 unit test 與 dev 模式
- 可透過 resolver 注入並替換測試場景所需 backend

---

## 與其他模組整合

| 模組          | 整合說明                         |
|---------------|----------------------------------|
| `services/*`  | 呼叫 store interface 進行資料存取 |
| `handlers/*`  | 不直接依賴 store，透過 service 呼叫 |
| `registry/*`  | 負責註冊 plugin 實作             |
| `plugins/*`   | 提供具體資料來源實作             |
| `cachestore/*` | 作為快取層與 Redis 整合，用於提升讀取效率 |

---

## 擴充項目建議

- [ ] 支援 composite store：如先寫入 file，再寫入 DB
- [ ] 將查詢包裝為 query object 支援 cache 層
- [ ] 提供 streaming store interface（for Loki, Kafka）

---