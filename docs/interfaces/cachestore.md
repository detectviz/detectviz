# CacheStore Interface 說明文件

> 本文件為 Detectviz 專案中 `CacheStore` 介面的設計說明與使用情境整理。

## 介面用途（What it does）

CacheStore 提供統一的快取介面，用於支援多種短期暫存行為，例如狀態記憶、重複執行防護、快取計算結果等。其設計目標為：

- 支援 memory / Redis 等可替換快取後端
- 提供簡潔 API：Get、Set、Delete、Has、Keys
- 可設定 TTL 避免資料長期殘留
- 可抽換為 no-op 快取以支援測試場景

## 使用情境（When and where it's used）

- 掃描模組快取 PDU 辨識結果（避免重複掃描）
- 快取自動部署中繼資料（如已產生 conf 標記）
- Web UI 暫存狀態資訊（如選單、快照等）
- 任務排程器暫存任務執行結果防止重複執行

## 方法說明（Methods）

- `Get(key string) (string, error)`：取得指定 key 對應值，不存在時回傳錯誤
- `Set(key string, val string, ttlSeconds int)`：寫入快取，ttl 為 0 表示永久有效
- `Has(key string) bool`：檢查 key 是否存在
- `Delete(key string)`：刪除指定快取資料
- `Keys(prefix string) ([]string, error)`：回傳所有符合 prefix 的 key（便於批次操作）

## 預期實作（Expected implementations）

- `internal/adapters/cachestore/mem.go`：記憶體快取，使用 go-cache 實作
- `internal/adapters/cachestore/redis.go`：Redis 快取，支援跨節點與 TTL
- `internal/adapters/cachestore/nop.go`：空實作，用於測試、禁用快取場景

## 關聯模組與擴充性（Related & extensibility）

- 可與任務模組、掃描模組、Web UI 共享使用
- 可支援自定義 namespace 或 key prefix 規則
- 若需支援分散式環境，可實作對應 Redis Cluster adapter