


# API Interface

本文件定義 detectviz 中 `internal/api/` 模組對外提供的 API 註冊、回應格式、錯誤處理與版本管理相關介面，供其他模組（如 handler, middleware, plugin）整合使用。

---

## 1. Router 註冊介面

統一註冊 Echo router：

```go
func RegisterRoutes(e *echo.Echo)
```

- 負責註冊 `/api/vX/` 路由群組與對應模組 handler。
- 每個模組 handler 應實作 `RegisterRoutes(g *echo.Group)`。

---

## 2. Middleware 鏈整合

在 `middleware.go` 中建立預設 middleware 組合：

```go
func SetupMiddlewares(e *echo.Echo, deps ...interface{}) {
    e.Use(Logger())
    e.Use(Tracing())
    e.Use(Metrics())
    e.Use(Auth(authenticator))
    e.Use(Recovery())
}
```

- 每層 middleware 為獨立模組註冊
- 支援由 registry 擴充 Plugin middleware

---

## 3. Response 包裝介面

標準 JSON 回應：

```go
func JSON(c echo.Context, data interface{}) error
func Error(c echo.Context, code int, msg string) error
```

統一格式：

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
  "error": "unauthorized",
  "message": "Access denied"
}
```

---

## 4. DTOs 定義規範

每個模組的 request/response DTO 建議放置於：

```
internal/api/dtos/{mod}.go
```

範例：

```go
type CreateAlertRequest struct {
    Name string `json:"name" validate:"required"`
    Level string `json:"level"`
}
```

搭配 validator 使用，並於 handler 中解析：

```go
var req CreateAlertRequest
if err := c.Bind(&req); err != nil {
    return api.Error(c, 400, "Invalid input")
}
```

---

## 5. 錯誤封裝規範

所有錯誤建議使用封裝類別：

```go
type APIError struct {
    Code    string
    Message string
}
```

集中定義於：

```
internal/api/errors/api_error.go
```

---

## 6. 未來擴充介面規劃

- 支援 plugin handler 註冊至 `/api/plugins/{mod}/...`
- 動態 API version routing（根據 header 決定 handler）
- 自動導出 Swagger / OpenAPI 定義