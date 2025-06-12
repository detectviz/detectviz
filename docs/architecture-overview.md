# Detectviz 架構總覽

本文件概述 detectviz 的核心架構設計原則、模組劃分、依賴邊界與開發策略，提供開發者與貢獻者統一理解整體運作方式與模組定位。

---

## 架構目標與設計原則

- **模組化與可插拔**：所有功能模組皆可 plugin 註冊與注入
- **清楚的責任分層**：從 handler → service → store → plugin 階層明確
- **可版本化的 API 設計**：支援 v1, v1beta1 等 handler 結構
- **統一的 CLI 與 Web 工具框架**：共用 config、auth、middleware 設計
- **Platform-first 思維**：架構即平台，讓每個專案可快速複用組裝

---

## 模組劃分與目錄結構總覽

```
detectviz/
├── apps/                  # 各應用程式主程式（如 server、cli）
│   ├── cli/
│   └── server/
├── internal/              # 封裝應用邏輯與實作模組
│   ├── api/               # HTTP 路由設定與注入點
│   ├── handlers/          # 各功能模組 API handler（版本化）
│   ├── middleware/        # 自訂中介層與 plugin middleware 註冊器
│   ├── auth/              # Auth 驗證策略模組，可 plugin 擴充
│   ├── plugins/           # 插件註冊與載入（middleware, store, auth, etc.）
│   ├── services/          # 業務邏輯核心（Rule, Notifier, Logger, etc.）
│   ├── store/             # 多資料源儲存實作（如 redis, influx, file）
│   ├── system/            # 系統整合模組（apiserver, quota, grpc, hooks）
│   └── web/               # Web 前端組件與頁面 HTMX 實作
├── pkg/                   # 共用函式與 library 模組
│   ├── config/            # 組態載入與驗證
│   ├── cmd/               # CLI 指令模組（Cobra）
│   ├── validation/        # 輸入驗證工具
│   ├── security/encryption/  # AES 加密封裝與擴充
│   ├── infra/             # infrastructure 工具層（httpclient, trace, etc.）
│   └── utils/             # 純輔助工具模組（stringutil, sliceutil, retryer 等）
└── docs/                  # 架構與介面設計文件
```

---

## 架構依賴路徑（簡化流程）

```
[Web/CLI] → [API handler] → [Service] → [Store or Plugin] → [資料源 or 外部系統]
```

- 所有資料處理流程皆可透過 service 層注入自訂 store/plugin
- Handler 層與 Service 層保持解耦，方便測試與替換
- 中介層（如 CORS、Auth）可由 plugin 擴充或切換實作

---

## Plugin 插入點總覽

| 插件類型 | 註冊點模組             | 功能說明 |
|----------|-------------------------|----------|
| Middleware Plugin | `internal/plugins/middleware` | 註冊自訂中介層 |
| Auth Strategy     | `internal/plugins/auth`      | 擴充登入驗證策略 |
| API Hooks         | `internal/plugins/apihooks`  | 擴充外部 API 介接點 |
| Store Plugin      | `internal/store/plugins`     | 提供 Redis、MySQL、Influx 等後端存取 |
| Event Plugin      | `internal/plugins/eventbus`  | 提供事件處理與日誌傳送 |

---

## 可重用 scaffold 與應用組合

detectviz 架構提供良好 scaffold，可快速組裝下列應用：

- 偵測引擎平台（規則 + 通知）
- Web 可視化與告警審查
- Plugin 擴充執行平台（for external connector）
- 簡易 Web Dashboard 工具
- 內部工具快速整合 CLI + API + 視覺介面

---

## 延伸建議

- 所有模組應撰寫對應 `docs/interfaces/xxx.md`
- 實作流程建議依 `/todo.md` 建立順序與依賴清單
- 建議每一個模組 scaffold 皆提供 `_test.go` 與單元測試規格

---