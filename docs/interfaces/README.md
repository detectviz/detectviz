# Interface 撰寫準則（Detectviz 專案）

本文件定義 Detectviz 專案中撰寫 interface 的統一風格與格式，供所有開發者遵循。

本規範旨在：
- 統一 interface 撰寫風格，提升維護與閱讀性
- 協助 AI 工具（如 Cursor / ChatGPT）產生一致、可替換的實作
- 促進模組化、可測試性與依賴反轉原則的落實

---

## 命名與結構原則

- Interface 檔案放置於 `pkg/ifaces/{module}/{name}.go`
- 命名使用 PascalCase（如 `Logger`, `CacheStore`, `EventBus`）
- 一個檔案只定義一個主要 interface
- 每個方法數量建議在 3～7 個以內，過多請拆分子模組

---

## 註解風格與格式

### 每個 interface 應具備：

- 英文主註解（簡潔描述用途）
- 對應的繁體中文補充，使用 `// zh:` 開頭

### 方法註解範例：

```go
// Info logs a message at the info level.
// zh: 記錄 info 級別的日誌訊息。
Info(msg string, fields ...any)
```

### Interface 註解範例：

```go
// Logger defines the structured logging interface for Detectviz.
// zh: Logger 定義 Detectviz 中的結構化日誌介面。
```

---

## 設計原則

- 僅定義「Detectviz 需要什麼」，不耦合第三方套件實作（如 zap, Redis）
- interface 為抽象 contract，不包含具體邏輯
- 若需支援 context、trace、TTL、分群等擴充性，應納入 method 設計
- 實作放在 `internal/adapters/{module}` 中

---

## 相關目錄規範

- 所有 interface 的中文說明與使用情境應補充於 `docs/interfaces/{name}.md`
- 若 interface 被核心流程注入，請註記於 `bootstrap.Init()` 文件或架構圖中

- 撰寫對應說明文件時，請參考 [interface-doc-template.md](./interface-doc-template.md)，以統一內容結構與敘述方式
