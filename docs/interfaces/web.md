# Web 模組介面設計（internal/web）

本文件說明 detectviz Web 模組的對外介面、路由註冊方式與 HTMX 事件觸發邏輯。

---

## 1. Context 擴充介面

```go
type WebContext interface {
    echo.Context
    User() *auth.UserInfo
    RenderPage(name string, data any) error
    RenderPartial(name string, data any) error
}
```

- `User()`：取得登入者資訊
- `RenderPage()`：渲染完整 layout
- `RenderPartial()`：僅回傳內容區塊供 HTMX 更新

---

## 2. 路由註冊

Web 路由集中註冊於 `internal/web/router.go`，每個模組可透過 `RegisterPageRoutes(e *echo.Echo)` 將 route 掛入主程式：

```go
func RegisterPageRoutes(e *echo.Echo) {
    e.GET("/web/alerts/status", handleAlertStatus)
    e.GET("/web/rule/config", handleRuleConfig)
}
```

---

## 3. HTMX event 對應

HTMX 觸發方式依照 HTML 標記定義：

```html
<div 
  hx-get="/web/alerts/status" 
  hx-target="#status-panel" 
  hx-swap="innerHTML">
</div>
```

支援事件種類：

| 屬性        | 說明                  |
|-------------|-----------------------|
| `hx-get`    | 發起 GET 請求         |
| `hx-post`   | 發起 POST 請求        |
| `hx-trigger`| 綁定觸發條件（如 change）|
| `hx-target` | 更新 DOM 節點 ID     |
| `hx-swap`   | 指定內容取代方式      |

---

## 4. 回傳格式建議

Web handler 應回傳：

- layout 頁面 → `c.RenderPage("alert_status", data)`
- partial 片段 → `c.RenderPartial("components/table", data)`
- 錯誤回應 → `c.JSON(400, map[string]string{"error": "無效參數"})`

---

## 5. 權限注入與導覽樹處理

- 權限判斷應在 `context.go` 中注入 `auth.UserInfo`
- `navtree.Build()` 可依角色決定 sidebar 展示項目

---

## 6. 未來規劃

- [ ] 支援表單送出驗證（hx-post + form validation）
- [ ] `hx-boost` 導入簡化頁面切換
- [ ] i18n 多語系切換機制
