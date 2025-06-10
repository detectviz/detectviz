


# Validation Interface

本文件定義 detectviz 專案中 `pkg/validation/` 模組的使用方式與介面設計。該模組負責提供純函數形式的輸入驗證工具，避免複雜邏輯重複出現於 handler 或 service 層，並支援測試與擴充。

---

## 設計原則

- 僅提供純函數形式的靜態驗證邏輯
- 不支援 plugin 化（避免外部行為與副作用）
- 不與資料來源、密碼庫等狀態耦合
- 可搭配 schema 驗證器使用（ex: JSON Schema、DTO 驗證）

---

## 類型分類建議

```
pkg/validation/
├── basic.go         # 基本驗證（required, min, max, length）
├── regexp.go        # 格式驗證（email, IP, slug）
├── jsonschema.go    # JSON schema 驗證輔助（選用）
├── error.go         # 統一錯誤訊息結構
```

---

## 典型函式範例

```go
func IsEmail(s string) bool
func IsSlug(s string) bool
func MaxLength(s string, max int) bool
func Required(s string) bool
```

錯誤包裝範例：
```go
func ValidateEmail(field, value string) error {
    if !IsEmail(value) {
        return NewFieldError(field, "must be a valid email")
    }
    return nil
}
```

---

## 統一錯誤介面

```go
type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}
```

---

## 使用情境建議

- handler 層對輸入做欄位驗證與回傳友善錯誤
- service 層前置驗證（如 rule 名稱格式、notifier 限長）
- 避免在 handler 中手寫 regexp 與 string 判斷
- 搭配 middleware 驗證結構體（ex: binding + validation）

---

## 測試與覆蓋率

- 每個函式應搭配 `_test.go` 覆蓋正常與異常條件
- 避免動態生成錯誤訊息模板，保持訊息簡潔一致
- 支援中文/英文切換建議由外部處理（非 validation 層責任）

---

## 擴充方向

- 支援 struct tag 驗證整合（如 `validate:"required,email"`）
- 可考慮自動生成表單驗證腳本（如支援 web UI）
- 不提供 plugin 註冊機制，擴充請以新增函式為主

---