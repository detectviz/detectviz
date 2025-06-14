# scaffold: 初始目錄與範本建立

> 備註：原本 detectviz 舊版目錄搬移至 docs/detectviz-deprecated

- [Cursor Scaffold 指引補充](docs/README.md)



- [ ] 建立基礎目錄結構
  - `apps/`
  - `pkg/`
  - `internal/`
  - `plugins/`
  - `compositions/`
  - `scripts/`
  - `tools/`
  - `docs/`

- [ ] 建立 plugins scaffold
  - `plugins/core/auth/jwt/plugin.go`
  - `plugins/core/middleware/logging/plugin.go`
  - `plugins/community/importers/prometheus/plugin.go`
  - `plugins/community/exporters/influxdb/plugin.go`
  - `plugins/community/integrations/security/keycloak/plugin.go`
  - `plugins/tools/supportbundles/plugin.go`

- [ ] 建立 contracts interface 定義
  - `pkg/platform/contracts/importers.go`
  - `pkg/platform/contracts/exporters.go`
  - `pkg/platform/contracts/auth.go`
  - `pkg/platform/contracts/lifecycle.go`

- [ ] 建立組合與註冊框架
  - `internal/platform/registry/registry.go`
  - `internal/platform/composition/lifecycle.go`
  - `internal/platform/composition/resolver.go`
  - `pkg/config/loader/config_loader.go`

- [ ] 建立最小平台組合檔
  - `compositions/minimal-platform/composition.yaml`

- [ ] 建立 scaffold 註解 README
  - 每個 plugin 目錄包含 `README.md` 說明該 plugin 類型與功能範本說明