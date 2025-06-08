## Response:
`pkg/util` 是 Grafana 中負責提供 **小型、通用、無依賴工具函式與輔助模組** 的目錄，類似標準函式庫的補強工具。這些模組通常不屬於任何業務邏輯，只提供「低層抽象工具」，被廣泛重用於 `pkg/services`、`pkg/api`、`pkg/infra` 等模組中。

* * *

✅ 子模組功能與用途（依來源整理）
-----------------

| 子模組名稱 | 功能說明 |
| --- | --- |
| `cmd` | 執行 CLI command 或包裝 `exec.Command` |
| `cmputil` | 提供 systeminfo、CPU 核心、架構查詢等系統層工具 |
| `debouncer` | 實作 debounce 邏輯，避免短時間內重複觸發相同行為（如：view rebuild） |
| `errhttp` | 將錯誤轉換為 HTTP 錯誤格式（如 JSON）供 API 使用 |
| `httpclient` | 簡化 HTTP 請求操作（與 `pkg/infra/httpclient` 職責有部分重疊但更輕量） |
| `osutil` | 作業系統層工具，如檔案刪除、確認目錄是否存在 |
| `proxyutil` | 處理 HTTP Proxy header / remote addr 轉換工具 |
| `retryer` | 提供 retry 機制與 backoff 策略，常被資料寫入、API 傳送使用 |
| `ring` | 一種固定大小 queue buffer，常用於 metrics 或 stream 資料 buffer |
| `scheduler` | 排程工具，可做定期任務或延遲任務處理（類似 `cron` 簡化版） |
| `testutil` | 提供測試用工具（如產生假 user、token、context） |
| `ticker` | 定時器 helper，與 `time.Ticker` 差異是能明確控制起停 |
| `xorm` | 對 legacy ORM `xorm` 的輔助（如 SQL builder、date 過濾條件） |

* * *

🧩 調用對象範例
---------

| 使用模組 | 調用 `pkg/util` 功能 |
| --- | --- |
| `pkg/services/live` | 使用 `ring` 作為資料流記憶環 |
| `pkg/services/alerts` | 使用 `retryer` 發送告警結果到外部 API |
| `pkg/api/*` | 使用 `errhttp` 格式化錯誤輸出 |
| `pkg/login/*` | 使用 `httpclient` 呼叫外部身份認證 |
| `pkg/infra/log` | 使用 `cmputil` 取得系統資訊與版本報告 |
| `pkg/services/contexthandler` | 使用 `proxyutil` 處理 IP 與 header 資訊 |

* * *

✅ 與 `pkg/infra` 差異
------------------

| 模組 | 角色 | 內容風格 |
| --- | --- | --- |
| `pkg/util` | 工具箱 | 無依賴、stateless、小型 |
| `pkg/infra` | 平台設施 | 有狀態、有初始化（如 log, cache, tracing） |

例如：

*   `util/retryer` 是單純函式；但 `infra/db` 則會建立持久連線
    
*   `util/scheduler` 是簡化版排程；但 `infra/metrics` 會推送到 Prometheus
    

* * *

🧠 適用於你（detectviz）情境
--------------------

若你也想設計輕量、無依賴的共用工具，可以：

*   參考 `util/retryer`, `debouncer`, `ring` 的最小抽象設計
    
*   集中放在 `pkg/util/` 或 `internal/util/`，保留自訂性
    

是否需要我幫你整理 detectviz 可以沿用的 `util` 模組清單與模板？例如 `retryer`, `debouncer`, `testutil`？