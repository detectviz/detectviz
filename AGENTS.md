# DetectViz Agents Design

本文件說明 DetectViz 架構中具備主動行為的 Agent / Runner 設計方式與啟動流程，供 Cursor 與開發者參考。

> 本文件為 Codex / Cursor Scaffold 的參考依據之一。當你進行自動化 plugin scaffold、runner 實作或 platform lifecycle 任務時，請搭配以下文件一併參考：
>
> - `/todo.md`：目前尚未完成的 scaffold 與 runner 實作清單
> - `/README.md`：平台分層架構、plugin 分類與目錄規則
> - `/docs/`：包含 develop-guide、interfaces、foundation 等重要規範文件
>
> Cursor 將自動讀取本文件，判斷哪些模組屬於 Agent、Runner 類型，並據此補齊主動行為模組。

---

## Agent 與 Runner 定義

在 DetectViz 架構中，Agent 或 Runner 指具備主動執行能力的模組，例如定期收集資料、推送外部任務、發送通知、檢查健康狀態等。

這些模組通常實作下列 interface 中之一：

- `Importer`：資料收集來源
- `Exporter`：資料輸出目的地
- `Notifier`：事件通知觸發
- `HealthAware`：支援健康狀態回報
- `LifecycleAware`：支援啟動與關閉控制

---

## Agent 類型分類

| Agent 類型     | 說明                     | Interface 組合                  | 啟動方式            |
|----------------|--------------------------|----------------------------------|---------------------|
| ImportRunner   | 定期收集資料（如 Prometheus） | `Importer` + `LifecycleAware`     | platform 啟動時啟用 |
| ExportWorker   | 匯出內部資料至外部系統       | `Exporter` + `LifecycleAware`     | platform 啟動或排程 |
| NotifyAgent    | 接收事件並發送通知           | `Notifier`                        | event bus handler  |
| HealthReporter | Plugin 健康狀態查詢          | `HealthAware`                     | 由 lifecycle 查詢   |
| ScheduleRunner | 週期性任務控制（預留）        | `Schedulable` + `LifecycleAware`  | scheduler 控制      |

---

## 啟動與註冊流程

1. 使用者定義 `composition.yaml`，指定啟用哪些 plugin：
```yaml
community.importers.prometheus:
  enabled: true
  config:
    scrape_interval: 15s
    targets:
      - http://localhost:9090/metrics
```

2. platform 啟動後：
   - config loader 載入設定
   - registry 掃描所有 plugin 並呼叫 `Register()`
   - 組合 resolver 排序 plugins 啟動順序
   - lifecycle manager 啟動具 `LifecycleAware` 的 plugin
   - 若為 Importer，則產生 ImportRunner 開始週期性資料收集

---

## Plugin 需具備的行為

若 plugin 參與 agent 行為，則需符合下列準則：

- 實作至少一種 Agent Interface（如 Importer）
- 接收來自 config loader 的 `cfg any` 並解析
- 可被 registry 註冊（必須有 `Register()` 函式）
- 可由 lifecycle manager 控制啟動與關閉

---

## 測試與驗證建議

- 每種 Agent 類型應具備對應測試（見 `internal/test/integration/scaffold_test.go`）
- 所有 Agent 應可在 `composition.yaml` 中被切換啟用與禁用
- 健康狀態回報將整合至 `/health` 查詢 API（未來支援）

---