# Handlers Architecture

本文件說明 detectviz 中 `internal/handlers/` 模組的架構與設計原則，並定義 handler 在整體架構中的角色、版本化策略、與其他模組的整合關係。

---

## 模組目的

- 接收並解析 HTTP/REST 請求
- 呼叫對應 service 處理業務邏輯
- 回傳統一格式的 JSON 或 HTML 回應
- 不直接操作 DB 或進行商業邏輯，僅負責調度與錯誤處理

---

## 對應目錄結構建議

```
internal/handlers/
├── alert/v0alpha1/
│   └── alert_handler.go
├── rule/v0alpha1/
│   └── rule_handler.go
├── report/v0alpha1/
│   └── report_handler.go
└── common/
    └── response.go      # JSON 包裝、統一錯誤處理
```

---

## handler 基本結構範例

```go
func GetAlertStatus(c echo.Context) error {
    user := auth.UserFromContext(c)
    if !user.HasPermission("alerts:read") {
        return response.Forbidden("not authorized")
    }

    result, err := alertService.GetStatus(c.Request().Context())
    if err != nil {
        return response.InternalError("failed to fetch")
    }

    return response.JSON(c, result)
}
```

---


## API 版本化策略

- 每個模組 handler 應置於對應的 API 版本目錄（如 `v0alpha1`, `v1beta1`）
- 路由結構對應 `/api/{version}/{module}/...`
- 版本可透過 build tag 或 registry 進行版本註冊與切換

### 版本命名建議與原則

為了對應模組開發階段與穩定性，detectviz 採用語意化版本命名如下：

| 階段       | 推薦版本命名   | 用途說明 |
|------------|----------------|----------|
| 初期實驗   | `v0alpha1`     | 僅供內部使用，API 尚不穩定 |
| 穩定測試中 | `v1beta1`      | 結構穩定但尚未正式釋出 |
| 正式公開   | `v1`, `v2`     | API 穩定、對外公開 |

#### 設計原則

- 所有 handler 應依據模組內部開發進度，自行定義版本（如：`internal/handlers/alert/v0alpha1`）
- 各模組版本可獨立控管，例如 `alert/v1`, `rule/v1beta1`
- 未來可透過 Registry 或 config 註冊對應版本與 handler

此策略有助於：
- 保持 API 漸進演進
- 支援舊版保留與新版共存
- 建立版本升級流程測試基礎

---

## 與其他模組整合關係

| 模組         | 說明 |
|--------------|------|
| `internal/services/*` | 呼叫實際業務邏輯與資源操作 |
| `internal/store/*`    | handler 不直接依賴，可透過 service 使用 |
| `internal/middleware` | context 中注入 `UserInfo`、traceID 等資訊 |
| `pkg/validation`      | 驗證輸入結構與資料完整性 |
| `pkg/utils/response`  | JSON 回應統一格式與錯誤封裝處理 |

---

## JSON 回應格式建議

```json
{
  "status": "success",
  "data": {...}
}
```

錯誤格式：
```json
{
  "status": "error",
  "error": "unauthorized",
  "message": "You are not allowed to access this resource"
}
```

---

## 單元測試建議

- 每個 handler 應可注入 mock service
- 驗證狀態碼、body 回應結構與錯誤處理分支
- 使用 `echo.New().ServeHTTP(...)` 建立虛擬請求測試

---

## 未來擴充

- [ ] 自動化 route 掛載機制
- [ ] 基於 Swagger 或 JSON Schema 的 validator 建構器
- [ ] CLI 導出 handler stub
