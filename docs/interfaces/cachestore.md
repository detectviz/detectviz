# CacheStore Interface 說明文件

> 本文件說明 Detectviz 專案中的 `CacheStore` 介面設計原則、應用情境與實作結構。CacheStore 模組負責管理應用中的短期暫存資料，提升效能並減少重複操作。

---

## 設計目的（Design Purpose）

- 統一快取操作接口，避免不同模組重複實作快取邏輯
- 可支援多種後端（如 memory、Redis、Nop）
- 提供測試用無狀態快取，以簡化測試流程
- 可依據模組需求擴充 TTL、key namespace、動態註冊等能力

---

## Interface 定義（Methods）

```go
type CacheStore interface {
	Get(key string) (string, error)                  // 取得指定 key 對應值，不存在時回傳錯誤
	Set(key string, val string, ttlSeconds int) error // 寫入快取，ttl 為 0 表示永久有效
	Has(key string) bool                             // 檢查 key 是否存在
	Delete(key string) error                         // 刪除指定快取資料
	Keys(prefix string) ([]string, error)            // 回傳所有符合 prefix 的 key（便於批次操作）
}
```

---

## 使用情境（Use Cases）

- 掃描模組快取設備辨識結果，避免重複查詢
- 快取自動部署結果（如 conf 註記）
- Web UI 儲存快照狀態或快取選單資料
- 排程器記憶任務執行狀態避免重複執行

---

## 預期實作（Expected Implementations）

| 類型         | 實作檔案路徑                                             | 描述                     |
|--------------|----------------------------------------------------------|--------------------------|
| Memory       | `internal/adapters/cachestore/memory/memory.go`          | 使用 go-cache 作為記憶體快取 |
| Redis        | `internal/adapters/cachestore/redis/redis.go`            | 使用 Redis 實作跨節點 TTL 快取 |
| Noop         | `internal/adapters/cachestore/nop.go`                    | 空實作，用於禁用或測試用途 |

---

## 設定與註冊方式（Registry and Config）

- 設定由 `pkg/configtypes/cache_config.go` 提供
- 註冊邏輯於 `internal/registry/cachestore/registry.go`
- 可根據 config 中指定類型自動切換對應實作
- 支援 logger 注入以監控快取行為與錯誤

---

## 測試與 Mock（Testing & Mocking）

- 所有實作皆應具備對應 `_test.go` 單元測試
- 可使用 `MockCacheStore` 進行介面驗證
- 建議將測試實作放於與 adapter 相同目錄下或 `internal/test`

---

## 擴充建議（Extension Notes）

- 若需支援分散式環境，可實作 Redis Cluster adapter
- 若需監控命中率，可加入統計與 trace log 機制
- 可考慮支援 JSON / struct 序列化版本，以簡化應用層邏輯