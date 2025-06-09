## Registry 資源註冊系統

`Registry` 模組提供類似 Kubernetes 的資源註冊機制，透過 GVK（Group-Version-Kind）分辨不同類型的資源，並支援 CRUD 操作與 plugin 擴充，允許模組化元件以統一方式進行註冊與操作。

---

## GVK（Group-Version-Kind）識別結構

GVK 是資源的唯一識別格式，包含三個欄位：

```go
type GVK struct {
  Group   string
  Version string
  Kind    string
}
```

範例：

- Group：`core`、`datasource.grafana.com`
- Version：`v1`、`v1beta1`
- Kind：`Datasource`、`AlertRule`

---

## Interface 一覽

### Resource

所有可註冊的資源都需實作 `Resource` 介面，用以提供唯一名稱。

```go
type Resource interface {
  GetName() string
}
```

---

### ResourceHandler

針對單一資源類型的 CRUD 實作邏輯。

```go
type ResourceHandler interface {
  Get(ctx context.Context, name string) (Resource, error)
  List(ctx context.Context) ([]Resource, error)
  Create(ctx context.Context, res Resource) error
  Update(ctx context.Context, res Resource) error
  Delete(ctx context.Context, name string) error
}
```

---

### Registry

主要註冊表 interface，負責註冊 GVK 與路由 CRUD 操作至對應 handler。

```go
type Registry interface {
  Register(gvk GVK, handler ResourceHandler) error
  Get(ctx context.Context, gvk GVK, name string) (Resource, error)
  List(ctx context.Context, gvk GVK) ([]Resource, error)
  Create(ctx context.Context, gvk GVK, res Resource) error
  Update(ctx context.Context, gvk GVK, res Resource) error
  Delete(ctx context.Context, gvk GVK, name string) error
}
```

---


## 實作建議

- `pkg/registry/registry.go` 應提供記憶體版 Registry 實作，可用 `map[GVK]map[string]Resource` 儲存。
- `internal/registry/loader.go` 應負責初始化 GVK 與 schema。
- plugins 啟動時可透過 `Register(gvk, handler)` 註冊自己支持的資源類型。

---

## Schema 載入與驗證

每個資源的 schema 結構定義於 `pkg/registry/schemas/`，並由 `index.yaml` 統一維護 GVK 與對應檔案對照。系統可透過 `internal/registry/loader.go` 載入該 index 並查詢對應 schema 檔案。

schema 內容為 YAML 格式，並支援透過 CUE 語言進行驗證。驗證流程由 `pkg/registry/kinds/validator.go` 提供，支援以下功能：

- 使用 `Validate(path string, data []byte)` 對任意 YAML 資料進行結構驗證。
- 錯誤會統一包裝為 `SchemaValidationError`，方便上層處理。
- 驗證資料來源可為上傳檔案、API 請求、匯入資料等。

schema 實體與範例檔案建議如下配置：

```
pkg/registry/
├── schemas/
│   ├── host.schema.yaml
│   ├── datasource.schema.yaml
│   └── index.yaml
└── kinds/
    ├── validator.go
    └── testdata/
        ├── valid_host.yaml
        └── invalid_host.yaml
```

整體流程可結合 `Registry` 與 `Validator`，用於匯入驗證、API 資源建立前檢查等場景。

---

## 註冊範例與使用方式

系統可透過各資源模組（如 `/pkg/registry/apis/host`）提供的 `RegisterHost()` 或 `RegisterDatasource()` 方法，將對應 GVK 註冊至 Registry。可於系統初始化時批次掛載。

範例程式碼：

```go
r := registry.NewMemoryRegistry()

if err := host.RegisterHost(r); err != nil {
    log.Fatal(err)
}

if err := datasource.RegisterDatasource(r); err != nil {
    log.Fatal(err)
}
```

每個資源模組應提供：
- 對應 GVK 定義（Group-Version-Kind）
- Resource 實作
- ResourceHandler 實作（記憶體或外部來源）
- `RegisterX(r Registry)` 函式對外註冊

---

## 相關模組

| 模組名稱       | 關聯描述                       |
| -------------- | ------------------------------ |
| `versioning`   | 在每次 CRUD 時記錄版本。       |
| `importer`     | 匯入資料後透過 Registry 儲存資源。 |
| `libraryelements` | 每個 ElementKind 都可透過 Registry 管理為一種資源類型。 |

---

## Engine 與 Decoder 整合應用

當系統需要從檔案載入資源定義並進行驗證與註冊時，可使用 `internal/registry/decoder.go` 進行 YAML 解碼並取得 GVK，再透過 `internal/registry/engine.go` 完成 handler 註冊與 schema 掛載。

典型流程如下：

1. 呼叫 `DecodeAndValidate(schemaPath, filePath)` 解析 YAML，取得 GVK。
2. 透過 `Engine.RegisterHandler(gvk, handler)` 註冊 CRUD 處理器。
3. 使用 `Engine.RegisterSchema(gvk, schemaPath)` 註冊對應 schema。
4. 後續可由 `Engine.Handler()` 與 `Engine.SchemaPath()` 查詢已註冊內容。

此結構可支援：
- 外部匯入模組（如 importer）
- 元件組裝與還原（如 libraryelements）
- Plugin 動態擴充資源類型

該設計亦能保持每個資源註冊邏輯與驗證機制分離、可測試與模組化。