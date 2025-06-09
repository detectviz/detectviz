

# LibraryElements Interface 說明

LibraryElements 為 detectviz 提供元素式（Element-based）元件設計，可應用於平台中各種視覺化區塊、互動元件、儲存組件等場景。

## Element Interface

每個元件應實作 `Element` 介面：

```go
type Element interface {
  ID() string
  Type() string
  Data() []byte
}
```

### 方法說明：

| 方法  | 回傳型別 | 說明                     |
|-------|----------|--------------------------|
| ID    | string   | 元件唯一 ID              |
| Type  | string   | 元件分類（如 input, chart） |
| Data  | []byte   | 元件定義內容             |

---

## ElementStore Interface

元件可儲存至儲存庫，透過 `ElementStore` 管理 CRUD：

```go
type ElementStore interface {
  Save(e Element) error
  FindByID(id string) (Element, error)
  Delete(id string) error
  List() ([]Element, error)
}
```

此介面可實作為記憶體儲存、JSON 檔案、資料庫等儲存機制。

---

## ElementRenderer Interface

元件可轉譯為目標格式（HTML、Grafana JSON 等）：

```go
type ElementRenderer interface {
  Render(e Element) ([]byte, error)
}
```

Renderer 可依據元件類型與內容進行視覺化轉換，用於組裝 dashboard、模擬 preview 或產生靜態輸出。

---

## 延伸應用

- 可整合 Importer 將外部 Grafana panel 載入為 Element
- 支援 Library-based Template 重複使用與版本控管
- 結合 Layout 元件進行頁面組裝與嵌套

LibraryElements 模組設計原則為模組化、版本化、組合式，便於後續跨平台延伸。

---

## 測試支援與假件

為支援元件儲存邏輯的單元測試，detectviz 提供 `FakeElementStore` 假件（定義於 `internal/test/fakes/fake_element_service.go`），其實作 `ElementStore` 介面，並使用記憶體 map 模擬儲存行為。

### 功能特色：

- 支援 `Save`、`FindByID`、`Delete`、`List` 全部操作
- 可用於測試 ServiceAdapter、Bootstrap 元件注入等場景
- 不依賴任何外部儲存或序列化機制

```go
store := fakes.NewFakeElementStore()
_ = store.Save(BaseElement{ID_: "b1", Type_: "chart", Raw: []byte(`{"title":"B"}`)})
```

此假件為平台組件模組測試基礎，搭配 `internal/adapters/libraryelements/service_adapter.go` 可模擬實際操作流程。