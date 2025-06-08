# pkg/services/featuremgmt:
以下是根據實際代碼整理的說明：

`pkg/services/featuremgmt` 功能說明
-------------------------------

此模組負責 Feature Toggle 管理與查詢，可視為 Grafana 的「功能開關中心」，包含以下重點元件：

- `manager.go`：定義 `FeatureManager` 結構，實作功能註冊（`registerFlags`）、更新（`update`）、是否啟用等判斷邏輯。
    
- `service.go`：提供 `ProvideManagerService` 給 DI 使用，初始化 FeatureManager，並從 `custom.ini` 的 `[feature_toggles]` 讀取設定。
    
- `goff_provider.go`：提供與 [GO Feature Flag](https://github.com/thomaspoignant/go-feature-flag) 整合的 `FeatureProvider`。
    
- `openfeature.go`：支援 [OpenFeature](https://openfeature.dev/) 架構，建立 Client 與 Provider 設定。
    
- `usage_stats.go`：產生目前開啟的 Feature 對應的 Prometheus metrics。
    

整體用途為：集中管理 Grafana 所有的實驗性、預覽與企業功能是否啟用，並可搭配前後端條件呈現、版本控制、Prometheus 指標等。

