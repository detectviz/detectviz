# API Architecture

本文件說明 detectviz 中 `internal/api/` 模組的設計結構、職責分工與與 handler/middleware 的整合邏輯，涵蓋 API router、response 標準格式與模組註冊方式。

---

## 模組目的

- 提供統一的 API 路由註冊與版本管理
- 整合 middleware 處理流程（auth, tracing, logger 等）
- 掛載各模組 handler，支援 `/api/v1/`, `/web/` 分路徑
- 提供 JSON 回應格式標準化、錯誤統一封裝

---

## 對應目錄結構建議

```
internal/api/
├── router.go            # 註冊 echo 路由，掛載所有 handler
├── middleware.go        # 中介層註冊鏈結（auth, logger, metrics）
├── response.go          # 統一 JSON 回應封裝
├── dtos/                # Data Transfer Object 結構定義
│   └── alert.go
├── errors/              # API 專用錯誤與錯誤碼
│   └── api_error.go
```

---

## API Routing 設計原則

- 所有 API handler 統一由 `router.go` 掛載
- 路由註冊依模組與版本組織：

```go
e := echo.New()
api := e.Group("/api/v0alpha1")
alert.RegisterRoutes(api.Group("/alert"))
rule.RegisterRoutes(api.Group("/rule"))
```

---

## Response 統一格式

所有回應應透過以下函式包裝：

```go
func JSON(c echo.Context, data any) error
func Error(c echo.Context, code int, msg string) error
```

標準回應結構：

```json
{
  "status": "success",
  "data": { ... }
}
```

錯誤格式：

```json
{
  "status": "error",
  "error": "bad_request",
  "message": "Invalid ID"
}
```

---

## 與其他模組整合關係

| 模組         | 關係說明 |
|--------------|----------|
| `handlers/*` | API handler 實作，實際處理功能邏輯 |
| `middleware` | 在 router 註冊鏈結各項中介層模組 |
| `auth`       | 由 middleware 驗證並注入 context |
| `dtos`       | 定義 API 輸入輸出格式 |
| `errors`     | 提供錯誤碼、訊息與統一轉換工具 |

---

## API 分層設計

```text
[Client Request]
     ↓
[Echo Router (router.go)]
     ↓
[Middleware 鏈 (auth, tracing, metrics)]
     ↓
[Handler 層 (handlers/*)]
     ↓
[Service, Store, Plugin 模組]
```

---

## 未來擴充建議

- [ ] 自動註冊所有 handler 至 router（plugin registry）
- [ ] JSON Schema or Swagger generator 支援
- [ ] 自定 `api-version` header 控制版本分派
