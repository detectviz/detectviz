## SqlStore:

> **原本 sqlstore 使用全域 xorm 變數，現已重構為具備 DI 的 SqlStore 實體，避免共享狀態、增加測試性。**

我們可以拆成兩個層次來理解：

* * *

🧠 1. 背後原本設計的問題
---------------

### ❌ 「全域 xorm 變數」的問題

假設原本是這樣寫：

```go
var xormEngine *xorm.Engine

func InitDB() {
    xormEngine = xorm.NewEngine(...)
}

func SaveUser(u *User) {
    xormEngine.Insert(u)
}
```

這樣的問題有：

*   **狀態共享**：全域變數 `xormEngine` 在多個模組共用，可能導致並發爭用、不易追蹤錯誤來源
    
*   **測試困難**：要寫單元測試時無法注入 mock，因為它是「硬編碼」的全域變數
    
*   **無法多實例運作**：如果未來你要連兩個不同資料庫，就無法同時建立兩個 `xormEngine` 例項
    

* * *

✅ 2. 重構後的解法：導入 DI + 實例物件 `SqlStore`
-----------------------------------

重構後會變成這樣的設計：

```go
type SqlStore struct {
    engine *xorm.Engine
}

func NewSqlStore(cfg DBConfig) *SqlStore {
    engine := xorm.NewEngine(...)
    return &SqlStore{engine: engine}
}

func (s *SqlStore) SaveUser(u *User) {
    s.engine.Insert(u)
}
```

並搭配依賴注入（DI）模式來呼叫：

```go
func NewUserService(store *SqlStore) *UserService {
    return &UserService{db: store}
}
```

這種設計的好處是：

*   ✅ **每個服務獨立持有 SqlStore 實體**，可依不同需求產生不同資料庫連線
    
*   ✅ **容易測試**：測試時你可以注入 mock 實作（例如 `MockSqlStore`）
    
*   ✅ **避免全域狀態污染**：所有狀態只存在於 `store` 實體中
    
*   ✅ **符合 Clean Architecture**：資料庫層可以更好地被抽象、替換
    

* * *

📌 結論
-----

這句話的意思可簡化為：

> 以前程式寫法用的是「全域資料庫物件」，重構後改為「注入的實體對象」，能避免共用副作用、支援單元測試、提升可維護性。

如你需要，我也可以畫出這種重構對應的 **Mermaid 架構差異圖** ✅

## SqlStore:
這裡所說的「**共用** 跟 **沒有共用**」的操作，指的是資料庫操作（例如 `Insert`、`Update`、`Find`）到底是透過：

*   ❌ **共用的全域變數來執行**
    
*   ✅ 還是透過 **每個服務持有自己的實例來執行**
    

* * *

❌ 共用（共享狀態）操作範例
--------------

```go
// 全域變數
var xormEngine *xorm.Engine

func SaveUser(u *User) error {
    return xormEngine.Insert(u)
}
```

*   `xormEngine` 是全域共享的
    
*   所有模組與 goroutine 都用這個變數
    
*   無法注入 mock，無法多連線配置
    

* * *

✅ 不共用（依賴注入、實例化）操作範例
-------------------

```go
type SqlStore struct {
    engine *xorm.Engine
}

func (s *SqlStore) SaveUser(u *User) error {
    return s.engine.Insert(u)
}
```

呼叫端透過依賴注入：

```go
store := NewSqlStore(cfg)
store.SaveUser(user)
```

*   每個 `SqlStore` 實體自己管理自己的 engine
    
*   若你要 mock 或用不同設定，只要 new 一個新實體
    
*   符合 SRP 與 DI 設計原則
    

* * *

🧩 差異總結
-------

| 項目 | 共用（全域） | 不共用（實體注入） |
| --- | --- | --- |
| 資源共享 | 是，全域共享一個 engine | 否，每個服務持有自己的 |
| 測試難度 | 難：無法 mock、需手動替換 | 易：可傳入 mock 實作 |
| 擴充性 | 低：只能用一組設定 | 高：可支援多資料庫、多配置 |
| 風險 | 高：全域競爭、race condition | 低：單一責任，易控管 |

* * *

若你需要，我可以幫你將這個設計差異轉成 Mermaid 架構對比圖，或提供 `SqlStore` 的建構 + 測試模式寫法 ✅

## XORM 的連線池行為:
✅ 是的，「是否共用」會直接影響 **XORM 的連線池行為**。這點很關鍵，來幫你釐清：

* * *

🔁 情境差異比較
---------

### ❌ 全域 `xormEngine`：單一連線池（共享）

```go
var engine *xorm.Engine
engine = xorm.NewEngine("mysql", dsn)
```

*   全專案所有查詢都透過同一個 `engine`
    
*   共用同一組 connection pool（例如最大連線數是 10，所有 goroutine 一起搶）
    
*   如果有 goroutine 卡住連線，**整個應用可能受影響**
    

* * *

### ✅ 每個 `SqlStore` 實例一個 `xorm.Engine`

```go
type SqlStore struct {
    engine *xorm.Engine
}

storeA := NewSqlStore(cfgA)  // 有自己的連線池
storeB := NewSqlStore(cfgB)  // 另一個不同的連線池
```

*   每個 `SqlStore` 實例對應一個獨立的連線池
    
*   適用於：
    
    *   多租戶（每個客戶一個資料庫）
        
    *   測試環境（mock engine）
        
    *   敏感資源隔離
        
*   減少因一個模組問題導致全域資源枯竭
    

* * *

🧠 重點說明
-------

| 項目 | 全域 engine | 實例 engine（DI） |
| --- | --- | --- |
| 共用連線池 | ✅ 是 | ❌ 否 |
| 風險隔離 | ❌ 無法隔離 | ✅ 可獨立運作 |
| 彈性 | ❌ 較低 | ✅ 高：可配置不同 DSN / Pool size |
| 測試可控性 | ❌ 難以 mock | ✅ 可用 stub / fake engine |

* * *

🛠 補充設定（控制連線池）
--------------

不論哪一種方式，XORM 都可以設定連線池參數：

```go
engine.SetMaxIdleConns(10)
engine.SetMaxOpenConns(50)
engine.SetConnMaxLifetime(time.Hour)
```

但只有「實例化 SqlStore」的設計能讓你根據需求去**設多組不同的連線池**。

* * *

如果你要支援：

*   多租戶
    
*   單元測試
    
*   限縮某些服務存取能力
    

✅ 改成 DI + 實例式 SqlStore 是對的方向。  
如你希望我幫你畫出 Mermaid 流程或抽象層級設計，也可以提供。

## 獨立的連線池:
不是的，**連線池（Connection Pool）不一定只有一個**，它的數量取決於你**建立了多少個 `xorm.Engine` 實例**。

* * *

✅ 關鍵觀念：每個 `xorm.Engine` 都擁有**獨立的連線池**
-------------------------------------

### ❗如果你只建立一個全域 engine：

```go
var engine *xorm.Engine
engine = xorm.NewEngine("mysql", dsn)
```

➡️ 所有資料庫操作都會共用 **同一個連線池**

* * *

### ✅ 如果你建立多個 engine（例如多個 `SqlStore` 實例）：

```go
storeA := NewSqlStore(dsnA) // 建立自己的 xorm.Engine
storeB := NewSqlStore(dsnB) // 再建一個不同的 xorm.Engine
```

➡️ 每個 `storeA.engine`、`storeB.engine` 都會有自己的**獨立連線池**  
（包含 max open connections, idle connections, etc.）

* * *

🧪 小測試驗證
--------

假設你這樣寫：

```go
a := xorm.NewEngine("mysql", "user1:pass@tcp(host1)/db1")
a.SetMaxOpenConns(5)

b := xorm.NewEngine("mysql", "user2:pass@tcp(host2)/db2")
b.SetMaxOpenConns(20)
```

你會有：

*   `a`：最多 5 條連線
    
*   `b`：最多 20 條連線  
    兩者完全獨立，互不干擾 ✅
    

* * *

📌 結論
-----

| 問題 | 回答 |
| --- | --- |
| 連線池只有一個嗎？ | ❌ 否，**每個 `xorm.Engine` 都會建立一個獨立的連線池** |
| 什麼情況下只有一個連線池？ | 當整個程式只用一個全域 `xormEngine` 時 |
| 如何讓每個模組使用不同連線池？ | 使用依賴注入，為每個模組建立自己的 `SqlStore`（內含不同 `xorm.Engine`） |

* * *

如果你要做「多租戶」、「連多個資料庫」、「測試隔離」、「分開管控連線」，那麼多個連線池的設計是必要的 ✅  
要我幫你畫出這個概念的架構圖嗎？（Mermaid）

## 連線池數量 & 資料庫連線數限制:
✅ 不會限制「幾個**連線池**」，但會限制「幾條**同時連線**」，這兩個概念一定要分清楚：

* * *

🔍 關鍵區分：**連線池數量** vs **資料庫連線數限制**
---------------------------------

| 項目 | 意思 | 誰控制 |
| --- | --- | --- |
| **連線池數量** | 你建立了幾個 `xorm.Engine` 實例，每個都有自己的 pool | 由你（應用程式）決定 |
| **最大連線數上限** | 資料庫允許的「同時連線數（connections）」 | 資料庫本身控制 |

* * *

🧠 實例說明
-------

### 你這邊：

```go
// 建立三個連線池，每個池最多開 10 條連線
storeA.engine.SetMaxOpenConns(10)
storeB.engine.SetMaxOpenConns(10)
storeC.engine.SetMaxOpenConns(10)
```

➡️ 程式理論上最多會開 **30 條連線**

### 資料庫端：

```sql
-- 以 MySQL 為例
SHOW VARIABLES LIKE 'max_connections';
```

如果顯示：

```
max_connections = 100
```

代表整個資料庫最多允許 **100 條同時連線**

* * *

❗所以如果你開太多連線池 + 每池太多連線，就會遇到：
---------------------------

*   ❌ `Too many connections`
    
*   ❌ 資料庫拒絕連線或直接關閉應用程式的連線
    

* * *

✅ 實務建議（以 MySQL 為例）
------------------

1.  **查詢資料庫最大允許連線數**
    
    ```sql
    SHOW VARIABLES LIKE 'max_connections';
    ```
    
2.  **評估應用最多會開幾個 engine**
    
    *   如果會啟用 3 個 engine
        
    *   建議 `SetMaxOpenConns(20)` × 3 → 共 60 條連線
        
3.  **合理保留 margin** 給其他應用與 DBA
    

* * *

📌 小結
-----

*   你可以建立很多連線池（很多 `xorm.Engine` 實體）
    
*   但資料庫最終會限制你「同時開啟的連線數」
    
*   **連線池不是問題，總連線數才是**
    

* * *

如果你要我幫你撰寫自動偵測機制（例如偵測 `max_connections` 後自動調整 pool size），我也可以幫你設計 ✅

## 同時連線數量:
評估「同時連線數量（max connections）」這件事其實是性能調校中非常關鍵的一環，以下是業界常見的 **評估方法與建議原則**，我會幫你拆解為 ✅實務操作 + 🧠概念理解：

* * *

✅ 評估步驟建議（應用端）
-------------

### 1\. **估算服務實例數量 × 平均每實例需要的連線數**

假設你有：

| 項目 | 數量 |
| --- | --- |
| Go 應用服務實例數 | 5 |
| 每個服務需要的資料庫連線池上限 | 20 |

➡️ 總需求 = `5 × 20 = 100` 條資料庫連線

* * *

### 2\. **依照業務流量尖峰進行壓力測試（如 locust、wrk）**

執行以下測試觀察：

*   **每秒 QPS**
    
*   **平均與最大響應時間**
    
*   **資料庫實際使用的連線數（用監控觀察）**
    

例如：

```sql
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';
```

* * *

### 3\. **設定上限 + buffer**

*   若預估尖峰需 100 條連線，建議資料庫 `max_connections` 設定為 **150–200**
    
*   預留空間給管理介面、備份、外部分析工具等非應用來源
    

* * *

🧠 背後評估邏輯：連線≠效能
---------------

過多連線反而會拖垮效能（Connection is not concurrency!）

*   資料庫不是多連線就快，它受限於：CPU、I/O、記憶體、索引設計
    
*   與其盲目拉高連線數，不如：
    
    *   做好 SQL 效能調整（慢查詢）
        
    *   降低不必要的連線時間
        
    *   控制好 pool idle timeout
        

* * *

📈 監控建議（InfluxDB + Telegraf + Grafana）
--------------------------------------

可用以下指標來決策：

| 指標名稱 | 監控意義 |
| --- | --- |
| `Threads_connected` | 現在有幾個活躍連線 |
| `Max_used_connections` | 曾經最高的連線數 |
| `Connections` | 每秒建立幾條新連線（高代表 pool 設定不當） |
| `Aborted_connects` | 連線被拒絕的次數 |

* * *

🛠 補充實用設定（以 Go `xorm` 為例）
-------------------------

```go
engine.SetMaxIdleConns(10)  // 建議設定 = MaxOpen 的一半
engine.SetMaxOpenConns(20)
engine.SetConnMaxLifetime(30 * time.Minute)
```

* * *

📌 小結建議
-------

| 類型 | 建議 |
| --- | --- |
| 小型內部系統 | 20~50 條連線即可 |
| 中型 API 服務 | 每實例 20~50 條連線 |
| 高併發系統 | 每實例壓測後調整至合理連線數，並與 DBA 協調 |
| 總連線數控管 | 留出至少 20% buffer 給 DBA、備份等外部工具 |

* * *

若你願意提供目前的服務架構（實例數、用戶數、QPS）或監控圖，我可以幫你估出合理的上限範圍 ✅
