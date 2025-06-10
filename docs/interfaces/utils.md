

# Utils Interface

本文件說明 detectviz 專案中 `pkg/utils/` 各個工具模組的功能分類與介面定義。此處僅包含無狀態、無副作用的純函式模組，用於字串處理、編碼、時間、重試、泛型操作等輔助性邏輯。

---

## 模組分類與介面摘要

### stringutil

```go
func Slugify(s string) string
func Elide(s string, max int) string
func JoinQuoted(items []string) string
```

### timeutil

```go
func ParseRFC3339(s string) (time.Time, error)
func NowMillis() int64
func FormatDurationHuman(d time.Duration) string
```

### sliceutil

```go
func Dedup[T comparable](items []T) []T
func Chunk[T any](items []T, size int) [][]T
func Filter[T any](items []T, fn func(T) bool) []T
```

### encodeutil

```go
func JSONSafeMarshal(v any) string
func Base64UrlDecode(s string) ([]byte, error)
```

### idutil

```go
func UUIDv4() string
func NanoID() string
func ShortID(prefix string) string
```

### validateutil

```go
func IsEmail(s string) bool
func MaxLen(s string, limit int) bool
func IsASCII(s string) bool
```

### retryer

```go
type Retryer struct {
    MaxRetries int
    Backoff    func(int) time.Duration
    RetryIf    func(error) bool
}

func (r Retryer) Do(fn func() error) error
```

### pointer

```go
func String(v string) *string
func Bool(v bool) *bool
func Int(v int) *int
```

### pathutil

```go
func SanitizePath(p string) string
func IsAbsURL(s string) bool
func JoinSafe(paths ...string) string
```

---

## 命名規範與擴充方式

- 每個模組為獨立 package，不應有跨模組依賴
- 僅允許內部使用通用語言與標準庫（不可依賴 config/logger/service/plugin）
- 所有 public 函式應具備對應單元測試

---

## 未來擴充方向

- 可加入 `mathutil/`, `csvutil/`, `testkit/` 作為工具測試輔助模組
- 若模組數量超過 10，可依照職責層級劃分至 `internal/shared/` 結構中

---