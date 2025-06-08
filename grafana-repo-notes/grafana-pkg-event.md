# pkg/events:
以下是 `grafana/pkg/bus` 與 `grafana/pkg/events` 兩個目錄的功能與差異說明：

* * *

✅ `pkg/bus`：**In-Process Message Bus**
--------------------------------------

### 📌 功能

*   提供 **同步發送與處理事件（Command / Message）的總線機制**，類似 CQRS 架構中的 Command Bus。
    
*   所有模組可透過 `bus.Publish(ctx, msg)` 發送訊息，由對應 handler 同步處理。
    

### 🔌 定義的介面與結構

```go
type Msg any
type HandlerFunc any

type Bus interface {
  Publish(ctx context.Context, msg Msg) error
  AddEventListener(handler HandlerFunc)
}
```

### 🧠 核心元件

*   `InProcBus`：Grafana 內部唯一實作。
    
*   `AddEventListener(fn)`：註冊處理特定訊息類型的 handler。
    
*   `Publish(ctx, msg)`：會找出該訊息類型的所有 handler 並呼叫grafana-pkg-all-code。
    

### 📦 使用場景

*   發送內部事件，例如：`FolderRenamed{FolderID, NewName}`。
    
*   減少模組之間直接依賴，達到模組解耦。
    

* * *

✅ `pkg/events`：**外部通知與系統事件格式定義**
--------------------------------

### 📌 功能

*   定義 **可傳送至外部系統（如 AMQP, Kafka, webhook）使用的事件資料結構**。
    
*   所有事件都為資料轉移物件（DTO），設計為**序列化安全且向後相容**。
    

### 🧾 事件格式範例：

```go
type OrgCreated struct {
	Timestamp time.Time
	Id        int64
	Name      string
}

type UserCreated struct {
	Timestamp time.Time
	Id        int64
	Name      string
	Login     string
	Email     string
}
```

### 📦 使用場景

*   事件可由內部觸發後寫入外部資料通道。
    
*   例如：新使用者註冊後觸發 `UserCreated` 事件，傳送到通知系統或審計系統grafana-pkg-all-code。
    

* * *

🔍 差異與關聯對照
----------

| 面向 | `pkg/bus` | `pkg/events` |
| --- | --- | --- |
| 定位 | 執行階段事件調度 | 定義靜態事件資料結構（DTO） |
| 資料流向 | 同步呼叫（in-process only） | 外部序列化使用（Kafka, webhook 等） |
| 用法 | `bus.Publish(ctx, &SomeCommand{})` | `emit(&UserCreated{...})` |
| handler 型態 | 函數註冊（透過 reflect 類型比對） | 不處理邏輯，只是 struct 定義 |
| 解耦程度 | 用於模組內部解耦 | 用於模組與外部世界解耦 |

* * *

### 🧩 組合應用建議（你可套用於 detectviz）

*   將 `bus` 作為模組間調度中心（例如 alert → reporter）
    
*   將 `events` 作為告警紀錄 / webhook / 通知格式標準層
    
*   分離處理邏輯（Handler）與資料定義（DTO），保持模組責任單一（符合 SRP）
    

是否需要我提供 bus 使用範例或在 detectviz 中建立事件發送範本？
