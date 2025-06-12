# Detectviz 模組概覽（Module Overview）

本文件提供 detectviz 架構中所有核心模組的概覽與功能定位，用於對齊模組職責、擴充邊界與實作範圍。

---

## 目錄結構對應核心模組

```
internal/
├── api/             # 路由注入與 API 統一回應流程
├── handlers/        # API 處理邏輯（版本化支援）
├── auth/            # 登入驗證策略模組（支援 plugin）
├── middleware/      # 中介層處理器（如 auth, cors, recovery）
├── plugins/         # 所有 plugin 註冊與 lifecycle 管理
├── services/        # 業務邏輯（rule, alert, notifier 等）
├── store/           # 資料儲存後端（redis, influx, mysql 等）
├── system/          # 系統級服務（http server, quota, grpc 等）
├── web/             # HTMX 視覺化模組
```

---

## 模組功能分類一覽表

| 模組          | 類型     | 功能說明 |
|---------------|----------|----------|
| api           | 核心框架 | 路由註冊、middleware 掛載、錯誤與回應格式處理 |
| handlers      | API 層   | 每個功能的 API handler，對應版本與 endpoint |
| middleware    | 插件模組 | 插入認證、CORS、限速、Logger 等通用中介處理 |
| auth          | 插件模組 | 驗證策略（如 local、keycloak、oauth），支援 plugin 註冊 |
| plugins       | 架構模組 | lifecycle 管理器，支援 CLI/API/Store/Auth 等 plugin 掛載 |
| services      | 邏輯層   | 真正執行邏輯，如比對規則、產生通知、整理報表 |
| store         | 後端儲存 | 實作讀寫資料來源，可 fallback（ex: logfile, influx, redis） |
| system        | 系統支援 | HTTP Server、診斷報告、GRPC、配額與 runtime 統一入口 |
| web           | 前端     | HTMX 元件與頁面組合，支援 iframe 模式嵌入與狀態維護 |

---

## Plugin 插件註冊支援

| 插件類型     | 說明 |
|--------------|------|
| AuthStrategy | 如 keycloak, saml, local password |
| Middleware   | 註冊 logger, gzip, access log |
| Store        | 註冊 redis, mysql, file, influx |
| CLI Command  | 自定指令注入至 `pkg/cmd/` |
| API Hook     | 額外 HTTP handler 註冊點 |

---

## 模組擴充建議與最佳實踐

- 每個模組皆應有對應 `docs/interfaces/xxx.md` 介面文件
- 可拆分版本化邏輯於 `v1/`, `v1beta1/` 目錄中
- 優先採用 plugin 注入實作，保持模組解耦與可測試性
- 若模組包含狀態操作，需抽象為 interface，並於 `store/` 中實作

---

## 開發對應與 scaffold 範例

- handler → 對應 `handlers/rule/v1/rule.go`
- service → 對應 `services/rule/service.go`
- store → 對應 `store/influx/rule_store.go`
- plugin → 對應 `plugins/auth/keycloak/strategy.go`

---

## 延伸文件

- [architecture-overview.md](./architecture-overview.md)
- [foundation.md](./foundation.md)
- [develop-guide.md](./develop-guide.md)
- [detectviz-apps.md](./detectviz-apps.md)
