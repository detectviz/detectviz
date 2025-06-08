# pkg/setting:

```bash
.
├── config_anonymous.go
├── config_auth_proxy.go
├── config_azure.go
├── config_cloud_migration.go
├── config_data_proxy.go
├── config_feature_toggles.go
├── config_featuremgmt.go
├── config_grafana_javascript_agent.go
├── config_grpc.go
├── config_jwt.go
├── config_k8s_dashboard_cleanup.go
├── config_passwordless_magic_link.go
├── config_plugins.go
├── config_quota.go
├── config_remote_cache.go
├── config_search.go
├── config_secrets_manager.go
├── config_secure_socks_proxy.go
├── config_smtp.go
├── config_storage.go
├── config_unified_alerting.go
├── config_unified_storage.go
├── config.go
├── configs_rbac.go
├── configs_zanzana.go
├── configtest
├── date_formats.go
├── expanders.go
├── provider.go
└── README.md
```


`pkg/setting` 是 Grafana 中負責 統一管理組態設定（設定檔解析、結構化欄位、環境變數覆寫、預設值與驗證） 的模組。它是整個 Grafana 啟動過程的設定核心，提供其他模組依賴的 `Cfg` 結構與各項設定細節。

* * *

✅ 功能總覽
------

| 功能 | 說明 |
| --- | --- |
| 解析 `grafana.ini` 設定檔 | 使用 `ini.v1` 解析器載入與讀取分段設定 |
| 支援環境變數覆寫 | 可透過 `GF_XXX_YYY` 覆蓋 ini 中的值 |
| 結構化設定分類 | 每個主題對應一個 `setting_xxx.go` 檔案，定義專屬 struct 與初始化方法 |
| 驗證與預設值設定 | 多數欄位會做 `MustBool`, `MustInt`, `MustDuration` 等型別轉換與下限檢查 |

* * *

🧱 設定結構 (`Cfg`)
---------------

所有設定會聚合到 `setting.Cfg` 結構中：

```go
type Cfg struct {
  Raw    *ini.File
  Logger log.Logger

  AppUrl string
  Env    string
  Quota  QuotaSettings
  Plugins PluginSettings
  ...
}
```

* * *

子模組說明（部分）
------------

| 子檔案 | 功能 |
| --- | --- |
| `setting_plugins.go` | Plugin 安裝與更新策略設定grafana-pkg-all-code |
| `setting_remote_cache.go` | Redis 等快取設定（加密、prefix）grafana-pkg-all-code |
| `setting_unified_storage.go` | Unified Storage（新儲存引擎）細節控制，如 dual writer、shardinggrafana-pkg-all-code |
| `setting_anonymous.go` | 匿名登入設定，如 org\_role、限制功能grafana-pkg-all-code |
| `setting_search.go` | Dashboard reindex 與查詢效能設定grafana-pkg-all-code |
| `setting_unified_alerting.go` | Alerting cluster 設定與狀態儲存策略grafana-pkg-all-code |
| `expanders.go` | 定義支援 `env:`、`file:` 等自訂變數展開語法grafana-pkg-all-code |

* * *

🔗 與其他模組關係
----------

| 模組 | 說明 |
| --- | --- |
| `pkg/services/...` | 多數 service 會依賴 `setting.Cfg` 中的特定欄位來決定啟動邏輯 |
| `pkg/server` | 在 `main.go` 啟動流程中會先初始化 `setting.NewCfg()` |
| `pkg/plugins` | Plugin 的來源與啟用清單來自 `setting.PluginSettings` |
| `pkg/infra/log` | Logger 初始化時會依據 `setting` 中的等級與輸出位置進行設定 |

* * *

🧠 為什麼重要？
---------

Grafana 採用 設定中心 + 結構對映 + 動態解析 的設計，使得：

- 設定變更更易控管與落地
    
- 對應模組可以直接取結構，不需要再自行解析 ini
    
- 可支援未來動態設定儲存（如 cloud 控制台、GUI 編輯）
    

* * *

如你在 `detectviz` 也需要支援 `.ini`、`.env` 或動態設定系統，這個模組是非常適合參考的。我可以幫你簡化為一套 `setting` 框架版本，是否要我整理一份範本？
