## scaffold: 未完成元件實作補強

- [ ] 實作 LifecycleManager 主要操作
  - 檔案：`internal/platform/composition/lifecycle.go`
  - 區段：lines 228–247 目前為 TODO
  - 預期：補齊啟動、關閉、資源釋放邏輯

- [ ] 新增組合依賴解析器
  - 預期檔案：`internal/platform/composition/resolver.go`
  - 功能：實作 plugin depends_on 拓撲排序與組合順序解析

- [ ] 實作 YAML 設定載入器
  - 預期檔案：`pkg/config/loader/config_loader.go`
  - 功能：解析 composition.yaml 並提供 plugin 初始設定

- [ ] 建立最小 Server 入口範例
  - 預期檔案：`apps/server/main.go`
  - 功能：可啟動一個註冊組合、啟動 plugins 的 server 範例

- [ ] 修正 plugin 忽略設定參數問題
  - 檔案：`plugins/core/auth/jwt/plugin.go` → `NewJWTAuthenticator()`
  - 檔案：`plugins/community/importers/prometheus/plugin.go` → `NewPrometheusImporter()`
  - 說明：目前未解析 config；應支援對應結構與預設值 fallback