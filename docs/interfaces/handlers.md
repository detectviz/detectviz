# Handlers Interface

本文件定義 detectviz 中 `internal/handlers/` 模組的 handler 實作規範、版本化設計與對外介面形式，並說明與其他層（service、store、middleware）的資料交互方式。

---

## Handler 定義原則

每個 handler 應遵循以下原則：

- 僅負責解析輸入與回傳輸出，不含業務邏輯
- 僅依賴 context 與已注入資料（例如 auth.UserInfo）
- 所有副作用應交由 service 層處理
- 僅透過 response 包裝器輸出 JSON 或 HTML

---

## 介面範例

每個 handler 函式應為以下形式：

```go
func GetAlertStatus(c echo.Context) error
```

若需共用參數解析、context 注入，可定義以下 interface：

```go
type Handler interface {
    RegisterRoutes(g *echo.Group)
}
```

---

## 路由版本與模組規劃

每組 handler 應置於對應模組與版本之下：

```
internal/handlers/alert/v0alpha1/alert_handler.go
internal/handlers/rule/v1/rule_handler.go
```

並依照 `/api/v1/alert/...` 方式對應 API 路徑。

---

## 與其他模組整合

| 模組         | 說明                           |
|--------------|--------------------------------|
| `auth`       | 從 context 中讀取 UserInfo     |
| `middleware` | 提供 trace ID、user context    |
| `services/*` | 呼叫具體業務邏輯               |
| `store/*`    | 間接透過 service 操作資料儲存   |
| `response`   | 封裝統一回應（success/error）格式 |

---

## Response 建議格式

使用 `internal/handlers/common/response.go` 提供的：

```go
func JSON(c echo.Context, data interface{}) error
func Error(c echo.Context, code int, message string) error
```

標準格式如下：

```json
{
  "status": "success",
  "data": { ... }
}
```

或

```json
{
  "status": "error",
  "error": "unauthorized",
  "message": "not allowed"
}
```

---

## 未來擴充項目

- 支援自動化 DTO 驗證與 error map 回應
- 加入統一 binding validator 結構（可擴充至 Web 與 API）
- 將路由定義與版本對應拆分為註冊介面

---
