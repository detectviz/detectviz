# Versioning Interface 說明

`pkg/versioning/interface.go` 定義 Versioned 與 VersionManager 兩個介面，負責整合所有可版本控制的物件與版本管理策略。

## Versioned Interface

```go
type Versioned interface {
  GetVersion() string
}
```

### 說明：

| 方法        | 回傳類型 | 說明                               |
|-------------|----------|------------------------------------|
| `GetVersion()` | `string` | 回傳目前版本字串（如 "v1.2.3"）     |

### 使用場景：

- 每個匯入或註冊的資源皆可實作 Versioned，提供統一版本查詢能力。
- 可搭配匯入器、匯出器或 plugin lifecycle 判斷版本差異。

---

## VersionManager Interface

```go
type VersionManager interface {
  Register(name string, v Versioned) error
  Get(name string) (Versioned, bool)
  Compare(a, b string) int
}
```

### 方法說明：

| 方法        | 回傳類型              | 說明                                   |
|-------------|-----------------------|----------------------------------------|
| `Register()` | `error`               | 註冊一筆可版本物件與名稱              |
| `Get()`      | `(Versioned, bool)`   | 根據名稱取得先前註冊的版本物件        |
| `Compare()`  | `int`                 | 比較兩個版本字串（正值代表 a > b）     |

### 比較回傳範例：

- `Compare("v1.2.0", "v1.1.9")` → `1`
- `Compare("v1.0.0", "v1.0.0")` → `0`
- `Compare("v0.9.0", "v1.0.0")` → `-1`

---

## 延伸應用

- Plugin 掛載前檢查版本衝突
- Importer 載入資料前比對是否為較舊版本
- Dashboard 與 Library Element 支援版本還原與更新

---

## 測試支援與假件

為方便模組單元測試與注入，系統提供 `FakeVersionStore`（位於 `internal/test/fakes/fake_versionstore.go`），其實作 VersionManager 介面並提供可自訂的版本比較邏輯。

### 假件功能：

- 使用 map 儲存註冊資源
- 實作 `Register()`、`Get()`、`Compare()` 方法
- 支援注入 `CompareFn` 以模擬任意版本比較行為
- 適合測試 plugin 掛載、importer 載入前後版本差異行為

```go
store := fakes.NewFakeVersionStore()
store.CompareFn = func(a, b string) int {
  if a == "v2.0.0" && b == "v1.0.0" {
    return 1
  }
  return 0
}
```

此假件可配合 StoreAdapter、Engine 等模組測試版本邏輯，避免實際依賴記憶體儲存實作。