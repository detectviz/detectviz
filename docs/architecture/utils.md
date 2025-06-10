# Utils Architecture

本文件說明 `pkg/utils/` 工具模組的設計原則與分類方式，避免 utils 成為責任混雜的 catch-all 工具箱，確保每一項功能具備可測試性、模組化與維護便利性。

---

## 設計原則與邊界

- 僅放置無狀態、無副作用的純輔助函式
- 工具模組不得依賴外部服務（如 DB、Logger、Plugin）
- 僅供轉換、格式處理、泛型操作等工具性目的
- 每個模組皆應有對應測試（`_test.go`）與文件說明
- 不得存取 domain 層物件、DTO 或 config provider

---

## 建議目錄結構與模組分類

每個工具模組獨立為一個子 package，避免 `stringutil.go`, `timeutil.go` 等過度集中導致維護困難。

```
pkg/utils/
├── stringutil/       # ex: Slugify, SplitN, Elide, JoinQuoted
├── timeutil/         # ex: ParseRFC3339, NowMillis, FormatHuman
├── sliceutil/        # ex: Dedup, Chunk, Map, Filter (泛型支援)
├── encodeutil/       # ex: JSONSafeMarshal, Base64UrlDecode
├── idutil/           # ex: UUIDv4, NanoID, ShortID
├── validateutil/     # ex: IsEmail, IsASCII, MaxLen
├── retryer/          # ex: Retry(fn), WithBackoff, RetryIf(error)
├── pathutil/         # ex: SanitizePath, IsAbsURL, JoinSafe
├── pointer/          # ex: StringPtr, IntPtr, BoolPtr
```

> 每個子目錄皆為獨立 Go package，使用者以模組名稱清楚辨識功能歸屬（如 `stringutil.Slugify()`）。

---

## 泛型支援建議（Go 1.18+）

建議於 `sliceutil/` 實作通用型別操作函式：

```go
func Dedup[T comparable](items []T) []T
func Map[T any, R any](items []T, fn func(T) R) []R
func Filter[T any](items []T, fn func(T) bool) []T
```

---

## Retryer 工具模組建議

模組位置：`pkg/utils/retryer/`

核心型別：

```go
type Retryer struct {
    MaxRetries int
    Backoff    func(int) time.Duration
    RetryIf    func(error) bool
}

func (r Retryer) Do(fn func() error) error
```

適用場景：
- httpclient 重試
- 外部 plugin 探活或 webhook 通知
- rule service 較不穩定的非同步調用

---

## 測試與覆蓋建議

- 每個工具函式應設有對應 `_test.go` 檔案
- 必須涵蓋正常流程、邊界條件與錯誤情境
- 測試應避免依賴全域狀態與非純粹輸入

---

## 禁止事項

- 不可使用 `panic()`、`MustXxx()` 無容錯邏輯
- 不可在 utils 呼叫 logger、metrics、外部網路
- 不可放置配置解析器、plugin 註冊器等與應用邏輯耦合內容

---

## 延伸建議與未來方向

- 搭配 `golangci-lint` 禁止不當使用方式（如 fmt.Println）
- 可新增 `pkg/utils/testkit/` 作為測試輔助工具集
- utils 模組數量超過 10 以上可考慮拆至 `internal/shared/`

---