# 設定指南（Config Guide）

本文件說明 detectviz 架構中的設定檔設計方式、載入流程與 plugin 擴充支援，適用於 CLI、Server、Plugin 等元件。

---

## 設定檔載入順序

detectviz 採用下列順序載入設定：

1. 預設值（定義於 `pkg/config/defaults.go`）
2. `.env` 檔案（可選，位於執行目錄）
3. CLI 傳入參數（如 `--config=config.yaml`）
4. 環境變數（支援 `PLUGIN_`、`SERVER_` 等前綴）
5. 設定檔（支援 `.yaml`, `.json`, `.toml`）

---

## 設定格式與範例（YAML）

```yaml
server:
  port: 8080
  log_level: info

plugin:
  auth:
    keycloak:
      client_id: detectviz
      url: https://auth.example.com
  store:
    redis:
      addr: localhost:6379
      db: 1
```

---

## 設定命名規則建議

| 類型       | 命名範例                             |
|------------|--------------------------------------|
| server     | `server.port`, `server.log_level`    |
| plugin     | `plugin.auth.keycloak.url`           |
| feature    | `feature.rule.enabled`, `feature.web.ui` |
| cli        | `cli.output_format`, `cli.profile`   |

---

## Plugin 設定自動載入原則

- 所有 plugin 可透過 `RegisterConfigPrefix("plugin.xxx")` 註冊命名空間
- plugin 可在 `Init()` 中取得 config：
  ```go
  cfg := config.Get("plugin.auth.keycloak.url")
  ```
- plugin 應定義自己的 config struct 並自動映射：

  ```go
  type KeycloakConfig struct {
      URL string `mapstructure:"url"`
  }

  var kc KeycloakConfig
  config.UnmarshalKey("plugin.auth.keycloak", &kc)
  ```

---

## CLI 設定檔支援

- CLI 指令可讀取 config 中的共用參數：
  - e.g. `plugin.path`, `output.format`, `profile`
- CLI 設定檔載入支援 `--config`, `--env`, `--profile`

---

## 組態模組與目錄對應

| 模組位置                  | 目的                               |
|---------------------------|------------------------------------|
| `pkg/config/`             | 提供 Config 介面與統一載入工具     |
| `pkg/configtypes/`        | 各模組 Config 結構定義             |
| `config/`（可選）         | 範例設定檔或開發環境預設組態檔案   |
| `internal/bootstrap/init.go` | 在伺服器啟動階段載入設定並注入模組 |

---

## 延伸建議

- 所有 config key 應有對應預設值
- 插件 config 使用 `plugin.{type}.{id}` 命名空間
- config 可注入至 plugin/service 中，避免使用全域變數

---
