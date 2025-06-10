# Web Architecture

本文件描述 `internal/web/` 模組的整體架構、資料流程與設計原則，目標為支援基於 HTMX + Echo 的模組化 Web UI 系統，可與 detectviz 各模組無縫整合。

---

## 目的

- 提供統一模板渲染方式（可切換 layout）
- 支援 HTMX 部分更新（partial）
- 建立共用元件（如 sidebar、navbar、table）
- 搭配模組頁面，如 alert status、record dashboard 等
- 模組頁面支援 plugin 擴充（未來）

---

## 目錄結構建議

```
internal/web/
├── render/
│   ├── renderer.go         # 建立 template 引擎（echo + glob）
│   └── layout.go           # 提供多 layout 的選擇與組合邏輯
├── context.go              # 封裝 echo.Context，注入使用者等資訊
├── binding.go              # 處理 HTML 表單綁定與驗證
├── response.go             # 封裝 JSON / HTML 回應格式
├── router.go               # 註冊 Web 路由
├── navtree/                # 根據角色與模組生成左側導覽列
│   └── navtree.go
├── pages/
│   ├── alert_status.html
│   └── dashboard_config.html
├── partials/
│   ├── sidebar.html
│   ├── table.html
│   └── topbar.html
├── static/
│   ├── htmx.min.js
│   ├── tabulator.min.js
│   └── style.css
```

---

## 資料流程與組件交互

```text
client <--HTMX--> /web/alert/status (HTML Partial)
         |
         |--> router.go --> Render(ctx, "alert_status", data)
                                 |
                                 |--> layout.go + partials 組合
                                 |
                                 |--> 渲染 HTML 回傳
```

- 所有頁面可透過 `htmx-get`, `hx-target`, `hx-swap` 動態更新內容
- 權限資訊注入在 `context.go` 中，用於 navtree 權限濾除

---

## 渲染原則

- 使用 echo 的 `Render` 介面
- 支援 layout + page 拆分：layout base.html + `{{ block "body" . }}`
- 可支援 dark mode / theme layout 切換（layout 註冊點）

---

## 擴充性設計

- plugins 可註冊：
  - navtree 節點
  - 新頁面組件（例如：plugin 定義其 pages/*）
  - partial 插槽（未來支援 slot 替換）

---

## 使用技術與依賴

- [HTMX](https://htmx.org/)
- [Tabulator](https://tabulator.info/) - 表格互動元件
- [Echo](https://echo.labstack.com/) - Go HTTP Framework

---

## 未來功能規劃

- [ ] login/logout 頁面模組化
- [ ] 用戶喜好設定持久化
- [ ] form validator hook 與組件化
- [ ] web partial 渲染單元測試
- [ ] iframe 外掛支援（如 grafana embed）

---
