# pkg/server 概覽

`pkg/server` 目錄是 Grafana 啟動與模組管理的核心邏輯之一，負責主程式啟動、模組註冊、HTTP/GRPC 伺服器初始化、與模組服務（如 distributor、storage、frontend）的協調與執行。以下是其主要功能說明：

* * *

📦 `pkg/server` 功能說明
--------------------

### 1\. **Server 啟動邏輯（`server.go`）**

*   提供 `New()` 與 `Run()` 方法建立與執行核心 `Server` 實例。
    
*   管理背景服務的啟動與關閉。
    
*   負責初始化 log、設定檔、PID 寫入與 systemd 通知等grafana-pkg-all-code。
    

### 2\. **模組伺服器（`module_server.go`）**

*   `ModuleServer` 是為了模組化 Grafana 運行的入口點設計。
    
*   支援 `init()` 啟動各個模組：如 `initDistributor()`、`initFrontend()`、`initInstrumentationServer()` 等。
    
*   用於建構 dskit-based 模組並管理模組依賴關係grafana-pkg-all-code。
    

### 3\. **模組註冊與服務管理（與 `pkg/modules` 整合）**

*   使用 `modules.Engine` 與 `services.BasicService` 管理模組生命週期。
    
*   各模組如 `Distributor`, `StorageServer` 等會用 `WithName()` 註冊名稱，方便統一管理與診斷grafana-pkg-all-code。
    

### 4\. **Instrumentation 伺服器（`instrumentation_service.go`）**

*   啟動 Prometheus metrics HTTP server，對外提供 `/metrics` endpoint。
    
*   也可設定為健康檢查與 profiling 伺服器（透過 `gorilla/mux` 實作）grafana-pkg-all-code。
    

### 5\. **模組型服務定義（如 `distributor.go`, `memberlist.go`）**

*   每個模組會實作自己的 `initXXX()` 函式，返回一個符合 `services.Service` 介面的執行單元。
    
*   e.g., `initDistributor()` 初始化一個 grpc-based 分發模組服務grafana-pkg-all-code。
    

### 6\. **Runner 與 ModuleRunner（`runner.go`, `module_runner.go`）**

*   `Runner` 是整體伺服器的執行者，注入 config、密鑰、user service 等元件。
    
*   `ModuleRunner` 是簡化版，只載入模組需要的最小依賴（例如 feature toggles）grafana-pkg-all-code。
    

* * *

📂 總結
-----

| 子檔案 | 功能概要 |
| --- | --- |
| `server.go` | 核心伺服器建立與生命週期管理 |
| `module_server.go` | 模組化伺服器，管理多個 dskit 模組 |
| `instrumentation_service.go` | 提供 Prometheus metrics HTTP 伺服器 |
| `distributor.go` | 建立 grpc distributor 服務 |
| `memberlist.go` | 建立 kv-based 成員列表協調模組 |
| `runner.go` / `module_runner.go` | 提供 Grafana 啟動執行流程 |

依照你在 Detectviz 的分層與架構規劃，pkg/server 是「平台啟動控制核心」，負責主流程組裝、模組註冊與伺服器執行邏輯，並不應視為通用可重用套件。因此，建議遷移至以下位置：

# 遷移到 Detectviz internal/bootstrap/server/ 後怎麼做


📌 遷移理由：

項目	說明
📦 高度耦合 Detectviz 啟動流程	如模組註冊、metrics exporter、服務註冊皆與內部模組密切相關
❌ 不具重用性	無 plugin 化可能、不可被外部 app 重用
✅ Clean Architecture 原則	將初始化組裝邏輯集中於 internal/bootstrap 或 internal/platform 是最佳實踐
🔄 與 internal/infra/, internal/registry/, internal/auth/ 等模組配合頻繁	放於同層級更能統一服務架構與生命週期流程


⸻

🧱 遷移後建議結構

internal/bootstrap/
├── server/
│   ├── server.go
│   ├── module_server.go
│   ├── instrumentation_service.go
│   ├── runner.go
│   ├── module_runner.go
│   ├── distributor.go
│   └── memberlist.go

若你將 modules/, services/ 一併重構，可一併放入 internal/bootstrap/modules/、internal/bootstrap/services/

⸻

✅ 延伸重構建議
	•	將 server.go 內建構的 Server 結構設計為依賴注入容器，載入：
	•	plugin registry
	•	event dispatcher
	•	API router
	•	ModuleRunner 可定義 interface 轉為 bootstrap.Runner，支援 app 啟動入口統一化

⸻

是否需要我幫你規劃 bootstrap/server 搭配 apps/main.go 的組裝流程、或重構後的模組生命週期圖？也可以產出 Mermaid 架構圖來輔助說明。