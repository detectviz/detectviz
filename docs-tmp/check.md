# `/internal/adapters` 與 `/pkg/ifaces` 對應檢查要點

- interface 目錄: `pkg/ifaces/{module}/*`
- adapter 目錄: `internal/adapters/{module}/*`

## 1. **Adapter 與 Interface 必須一一對應**

| 檢查目標 | 說明 |
| --- | --- |
| 每個 `internal/adapters/{module}/xxx_adapter.go` 是否有實作對應的 interface？ | 例如：`modules/engine_adapter.go` 應實作 `ModuleEngine` interface |
| interface 定義位置是否為 `pkg/ifaces/{module}/`？ | interface 不應寫在 internal，須統一集中在 ifaces |

* * *

## 2. **命名一致性**

| 項目 | 檢查項目 |
| --- | --- |
| 檔名一致 | adapter 檔名應為 `{模組}_adapter.go`，如 `engine_adapter.go`, `runner_adapter.go` |
| struct 命名 | Adapter struct 命名應為 `{模組}Adapter`，如 `EngineAdapter` |
| interface 命名 | interface 命名應符合 `{模組目的}`，如 `ModuleEngine`, `Server`, `ElementService` |

* * *

### 3. **語意明確與封裝完整**

| 檢查項目 | 說明 |
| --- | --- |
| 是否只封裝必要功能？ | adapter 不應包含多餘邏輯，僅做封裝與轉換 |
| 是否注入必要依賴？ | struct 內部應正確依賴被包裝實作，不應產生全域耦合 |

* * *

### 4. **介面與實作一致性**

| 檢查項目 | 說明 |
| --- | --- |
| Adapter 是否完整實作介面所有方法？ | 確保無 interface 未滿足錯誤 |
| 測試中是否使用 fake/mock 注入？ | 例如 `fake_server.go`, `fake_config.go` 需存在且可測試 interface |

* * *

### 5. **配套測試與 fake**

| 檢查項目 | 說明 |
| --- | --- |
| 是否有 fake 實作？ | 每個 interface 至少需有 `internal/test/fakes/fake_xxx.go` |
| 是否有測試樣板？ | 每個 adapter 若具邏輯轉換，應補 `adapter_test.go` 驗證注入與行為正確 |

* * *

## ✅ 推薦檢查指令 (可自動化)

```bash
# 列出所有 adapters
find internal/adapters -name '*_adapter.go'

# 比對 interface 實作關係
grep -r 'implements' internal/adapters

# 檢查是否有對應 fake
find internal/test/fakes -name 'fake_*.go'
```

* * *

