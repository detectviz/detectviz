# pkg/tests 與 pkg/mocks

# pkg/tests:
`pkg/tests` 是 Grafana 原始碼中的 **集中整合測試模組（integration test harness）**，目的在於提供跨模組測試與共享測試邏輯，與 `pkg/mocks` 和 `pkg/util/testutil` 搭配運作，覆蓋底層儲存行為、資源 API 行為、查詢行為等核心流程。

* * *

✅ 功能與用途
-------

| 類型 | 說明 |
| --- | --- |
| 整合測試套件 | 專門針對 `pkg/storage/unified`、`pkg/services/ngalert`、`pkg/registry` 等模組設計一系列跨模組整合測試 |
| 測試抽象封裝 | 提供類似 `RunStorageBackendTest()`、`RunSearchBackendTest()` 的標準化 test runner，傳入不同 backend 即可共用測試邏輯grafana-pkg-all-code |
| 自動命名空間隔離 | 使用 `GenerateRandomNSPrefix()` 為每次測試產生獨立 namespace，避免並行測試資料衝突 |
| 支援跳過測試項目 | `TestOptions.SkipTests["TestBlobSupport"] = true` 可自訂跳過指定測試項目grafana-pkg-all-code |

* * *

🧩 與其他測試模組的差異
-------------

| 模組 | 目的 | 範圍 | 特點 |
| --- | --- | --- | --- |
| `pkg/tests` | 整合測試 | 全系統級 | 跨模組、多元 backend 比較 |
| `pkg/mocks` | 單元測試 mock | 針對 interface | 通常搭配 `mockery` 自動產生 |
| `pkg/util/testutil` | 單元測試共用工具 | test context, 假帳號 | 輕量級、偏向 test helper |
| 各模組 `/testutil` | 模組專屬 test helper | 限模組內使用 | 如 `ngalert/testutil`, `foldertest` |

* * *

📦 使用範例：測試 Unified Storage backend
----------------------------------

```go
RunStorageBackendTest(t,
  func(ctx context.Context) resource.StorageBackend {
    return NewMyStorageBackend(...)
  },
  &TestOptions{NSPrefix: "demo-test"},
)
```

此機制可套用到：

*   自訂後端儲存實作驗證
    
*   plugin 或 fork 運作環境測試
    
*   Regression test 與可擴充性測試
    

* * *

✅ 為何放在 `pkg/tests`
------------------

*   ✅ 提供內部共用的測試邏輯（非模組特有）
    
*   ✅ 減少重複撰寫測試樣板
    
*   ✅ 方便 CI/CD 中集中呼叫與篩選
    
*   ❌ 不會被 production runtime 引用（清楚區隔）
    

* * *

若你要在 `detectviz` 中設計一套模組化儲存或 alert 測試，也可參考這種設計。需要我幫你定義 `pkg/tests/backend_test.go` 的範本嗎？適用於 plugin or unified backend 驗證。

# pkg/mocks:
`pkg/mocks` 是 Grafana 中用來放置 **mock 類型的集中目錄**，其主要目的是支援單元測試與介面替代，特別適用於 dependency injection 架構中的測試環境。

* * *

✅ 功能與用途
-------

| 類型 | 說明 |
| --- | --- |
| Interface Mock 實作 | 使用工具（如 `mockery` 或 `gomock`）自動產生介面實作檔案 |
| 測試替身 | 提供可以斷言呼叫次數、參數、行為的物件 |
| 減少循環依賴 | 測試階段引用 `mocks`，不會造成正式邏輯相依性 |

* * *

🧩 Mock 的來源與產生方式
----------------

常見產生方式：

```go
//go:generate mockery --name=SomeInterface --output=pkg/mocks --with-expecter
```

*   使用 `mockery` 工具根據介面自動產出對應 mock 檔案。
    
*   範例：`pkg/mocks/T.go` 是對 `testing.T` 風格的 interface mock 實作grafana-pkg-all-code。
    

* * *

🆚 `pkg/testutil` 與 `mocks` 的差異
-------------------------------

| 模組 | 用途 | 是否可程式產生 |
| --- | --- | --- |
| `testutil` | 提供測試工具與資料（如 test context、假帳號） | ❌ |
| `mocks` | 對 interface 的自動 mock 實作，配合 mock assert 機制 | ✅（mockery / gomock） |

* * *

🧠 為何放在 `pkg/mocks`
-------------------

1.  **集中管理 auto-generated mock**，避免與手寫測試邏輯混淆
    
2.  **在模組之上層級**：讓所有模組都能引用 mock（避免 import loop）
    
3.  **配合 make test 或 CI 流程整合**
    

* * *

📦 使用場景
-------

*   在單元測試中注入 mock：
    
    ```go
    ctrl := gomock.NewController(t)
    mockDB := mocks.NewMockDatabase(ctrl)
    service := NewMyService(mockDB)
    ```
    
*   搭配 `EXPECT()` 檢查預期呼叫：
    
    ```go
    mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(...)
    ```
    

* * *

如果你要在 `detectviz` 導入類似架構，我可以幫你規劃：

*   `pkg/ifaces/` 放 interface
    
*   `pkg/mocks/` 放 mockery 產物
    
*   `tools/mockgen.sh` 做自動更新
    

# 遷移到 Detectviz pkg/tests 與 pkg/mocks 後怎麼做

以下是 pkg/tests 與 pkg/mocks 在 Detectviz 架構下的定位與建議遷移位置：

⸻

✅ pkg/tests（整合測試邏輯）建議遷移位置：

internal/testharness/

📌 理由：

面向	說明
🔧 功能定位	提供跨模組整合測試流程（如 alert + registry、storage + api）
📦 特性	包含測試 helper、test runner、命名空間管理、資源註冊測試
🔐 不應暴露	不屬於外部 library，應內部保留以供整合測試使用
☑️ 獨立可測試	可整合至 CI，對新 backend/plugin 實作驗證有幫助

若 Detectviz 有明確測試標準化需求（如 plugin compliance 測試），也可命名為 internal/compliance/ 或 internal/itest/

⸻

✅ pkg/mocks（interface mock 產物）建議遷移位置：

internal/test/mocks/

📌 理由：

面向	說明
📦 功能定位	自動生成的 interface 替身，用於單元測試注入
🧰 適用工具	mockery, gomock, moq 等生成工具產出
⛔ 不屬於 pkg	並非共用邏輯，不應開放給 production code 引用
✅ 清楚分層	配合 pkg/iface/ interface 管理，將 mock 分層放入測試區域最清晰


⸻

🧠 補充建議：

📁 最終建議結構：

pkg/
└── iface/
    ├── alert.go
    ├── registry.go
    └── ...

internal/
├── test/
│   └── mocks/             # interface mock
├── testharness/           # 整合測試（pkg/tests）
│   └── alert_storage_test.go


⸻

✅ 是否可加上 make mock 自動產生流程？

若你打算大量使用 mock，可以新增以下：
	•	tools/mockgen.sh：產出所有 mock（如 mockery）
	•	make mock：掃描 pkg/iface/ 並生成對應 mock 到 internal/test/mocks/

⸻

是否要我幫你產出 mockgen.sh 範例腳本與目錄結構？或者直接重命名與規劃你目前的 pkg/tests 轉成 internal/testharness/？