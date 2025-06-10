
# pkg/services/ 目錄：

----------------------- pkg/services/第一階段核心服務目錄 -----------------------

✅ 1. `pkg/services/auth/` - 認證服務
➡️ Detectviz: internal/auth/
- `pkg/services/accesscontrol/` - 權限與存取控制（RBAC, 角色、權限、授權驗證）：負責權限模型，所有資源（dashboard, datasource, alert, user, org, team）都會進行資源權限檢查，串接 RBAC。
➡️ Detectviz: internal/accesscontrol/
2. `pkg/services/authn/` - 認證服務（Authentication）
➡️ Detectviz: internal/authn/
3. `pkg/services/authz/` - 授權服務（Authorization）
➡️ Detectviz: internal/authz/
4. `pkg/services/dashboards/` - 儀表板核心服務
➡️ Detectviz: internal/dashboards/
5. `pkg/services/datasources/` - 資料來源服務
➡️ Detectviz: internal/datasources/
6. `pkg/services/sqlstore/` - 資料庫存儲服務
➡️ Detectviz: internal/sqlstore/
7. `pkg/services/user/` - 使用者管理服務
➡️ Detectviz: internal/user/
8. `pkg/services/org/` - 組織管理服務
➡️ Detectviz: internal/org/
9. `pkg/services/store/` - 存儲服務：「通用儲存層」服務，提供抽象的 Key-Value store 介面。它允許將設定、狀態、暫存資料等儲存到多種底層（如本地檔案、雲端物件儲存、DB 等）。標榜抽象與可擴展性：其他服務只需依賴 store 介面，而不用關心 backend 細節。
➡️ Detectviz: internal/store/
10. `pkg/services/provisioning/` - 配置服務：設定自動佈署/初始化（dashboard、datasource、alert 等自動建置）。dashboard、datasource、alerting 等服務都會在啟動時由 provisioning 自動初始化。
➡️ Detectviz: internal/provisioning/
11. `pkg/services/featuremgmt/` - 功能管理：動態功能開關、feature flag，讓不同租戶或不同環境可啟用不同功能，其他服務會查詢 feature flag 來決定是否啟用某功能，如 ngalert、rendering。
➡️ Detectviz: internal/featuremgmt/

----------------------- pkg/services/第二階段相關服務目錄 -----------------------

## 認證與授權相關服務 
- `pkg/services/anonymous/` - 匿名使用者
- `pkg/services/authapi/` - 認證 API 服務
- `pkg/services/extsvcauth/` - 外部服務認證
- `pkg/services/login/` - 登入服務
- `pkg/services/loginattempt/` - 登入嘗試記錄
- `pkg/services/ssosettings/` - SSO 設定服務
- `pkg/services/apikey/` - API 金鑰管理
- `pkg/services/ldap/` - LDAP 認證服務
- `pkg/services/oauthtoken/` - OAuth Token 服務
- `pkg/services/kmsproviders/` - 金鑰管理
- `pkg/services/licensing/` - 授權管理
- `pkg/services/secrets/` - 密鑰管理

## 使用者與組織服務
- `pkg/services/team/` - 團隊管理
- `pkg/services/temp_user/` - 臨時使用者
- `pkg/services/serviceaccounts/` - 服務帳號


## 系統與基礎設施服務 
- `pkg/services/frontend/` - 前端服務：載入 index.html 與靜態資源入口邏輯，支援 SPA 路由。可保留但抽象化處理，讓前端部屬方式可切換。
➡️ Detectviz: internal/frontend/
- `pkg/services/apiserver/` - API 伺服器
➡️ Detectviz: internal/apiserver/
- `pkg/services/grpcserver/` - gRPC 伺服器
➡️ Detectviz: internal/grpcserver/
- `pkg/services/contexthandler/` - 上下文處理：提供 API 請求的上下文注入，例如 user/org/session metadata 等。
➡️ Detectviz: internal/contexthandler/
- `pkg/services/datasourceproxy/` - 資料來源代理
➡️ Detectviz: internal/datasourceproxy/
- `pkg/services/pluginsintegration/` - 插件整合，處理插件錯誤追蹤
➡️ Detectviz: internal/pluginsintegration/
- `pkg/services/quota/` - 配額管理：用戶/租戶級資源配額（如 dashboard/alert/通知通道數量上限）。
➡️ Detectviz: internal/quota/


## 系統資料與事件運作服務
- `pkg/services/correlations/` - 關聯服務：解決「跨事件關聯、聚合」問題，是 alerting/monitoring/可觀測性場景下的分析核心。
➡️ Detectviz: internal/correlations/
- `pkg/services/stats/` - 監控系統運作指標（dashboard 數量、alert 數量、用戶數、API 呼叫次數等）。
➡️ Detectviz: internal/stats/
- `pkg/services/supportbundles/` - 發生問題時自動匯出診斷包，輔助除錯，依賴 `pkg/services/stats/`
➡️ Detectviz: internal/supportbundles/
- `pkg/services/caching/` - 快取服務：快取層（如 session、狀態、臨時中繼資料），提升查詢效能、減輕 DB 負載。
➡️ Detectviz: internal/caching/
- `pkg/services/cleanup/` - 清理服務
➡️ Detectviz: internal/cleanup/
- `pkg/services/hooks/` - 鉤子服務：是「事件觸發與擴展」的基礎設施，讓內建/外掛/外部系統皆能聯動。
➡️ Detectviz: internal/hooks/
- `pkg/services/live/` - 即時功能，提供 WebSocket 實現的資料串流（push 模型）資料來源如：告警通知即時流、log tail、即時測試結果等。
➡️ Detectviz: internal/live/
- `pkg/services/query/` - 查詢服務：統一查詢執行服務，包裝 panel、dashboard、alert 等查詢資料來源的流程。
➡️ Detectviz: internal/query/
- `pkg/services/queryhistory/` - 查詢歷史記錄：查詢歷史紀錄，追蹤 dashboard、panel、手動查詢的歷史與結果。
➡️ Detectviz: internal/queryhistory/
- `pkg/services/search/` - 搜尋服務：舊一代統一搜尋服務，支援 dashboard、資料夾、datasource、alert、team 等的模糊查詢。
➡️ Detectviz: internal/search/
- `pkg/services/searchV2/` - 搜尋服務 V2：新一代搜尋服務（搜尋建議/自動補全），支援更進階過濾、類型、權限、分頁與新資源型態（如 library panel）。
➡️ Detectviz: internal/searchV2/
- `pkg/services/searchusers/` - 使用者搜尋：專門針對 user、team、org 做搜尋、模糊過濾。
➡️ Detectviz: internal/searchusers/


### 共用服務 (需要抽象共用)
目前多用於 dashboards，但其設計抽象，可支援其他核心模組如 alert、notifications 等。

- `pkg/services/navtree/` - 導航樹：動態產生側邊選單項目（navigation tree）
➡️ Detectviz: internal/navtree/
- `pkg/services/preference/` - 偏好設定：儲存個人設定：主題、預設資料來源、預設資料夾等
➡️ Detectviz: internal/preference/
- `pkg/services/dashboardimport/` - 儀表板導入服務：提供 dashboard 的匯入、導入功能，支援 JSON、外部來源（Grafana.com、JSON 檔）、API 導入。
➡️ Detectviz: internal/dashboardimport/
- `pkg/services/dashboardversion/` - 儀表板版本控制：Dashboard 版本控管，追蹤每次異動、支援回溯、異動歷史瀏覽、比對。
➡️ Detectviz: internal/dashboardversion/
- `pkg/services/dashboardsnapshots/` - 儀表板快照：Dashboard 快照/分享，允許生成靜態快照，提供公開/時效性存取分享。
➡️ Detectviz: internal/dashboardsnapshots/
- `pkg/services/folder/` - 儀表板資料夾管理：Dashboard 資料夾結構，支援 dashboard 分類、資料夾層級權限、批次操作。
➡️ Detectviz: internal/folder/
- `pkg/services/libraryelements/` - 元素管理：Library Element 管理，像是可重用的查詢片段、文字模組等（新一代「元件庫」）。
➡️ Detectviz: internal/libraryelements/
- `pkg/services/librarypanels/` - 面板管理：Library Panel 管理，支援跨 dashboard 重用、統一管理圖表 Panel 設定。
➡️ Detectviz: internal/librarypanels/
- `pkg/services/annotations/` - 註解服務：允許在 dashboard 或特定圖表（panel）上添加「註解事件」。
➡️ Detectviz: internal/annotations/
- `pkg/services/tag/` - 標籤服務：為資源（如 Dashboard、Data Source、Plugin）打標籤與分類管理。
➡️ Detectviz: internal/tag/
- `pkg/services/star/` - 收藏服務：提供 dashboard 收藏（star）功能，讓使用者快速標記/收藏常用 dashboard。支援依 user 分類、查詢已收藏項目。
➡️ Detectviz: internal/star/

## 遷移至 pkg/util
- `pkg/services/validations/` - 提供統一的資料驗證工具，讓其他服務可重用。
- `pkg/services/encryption/` - 加密服務


## 遷移至 /plugins (插件化)
- `pkg/services/signingkeys/` - 管理 JWT / Token 簽名金鑰
- `pkg/services/rendering/` - 渲染服務：提供 Dashboard 轉圖片的後端渲染（使用 Chrome Headless）。
- `pkg/services/screenshot/` - 截圖服務：與 rendering/ 類似，但更偏向任務式的截圖（例如 API call 定期拍圖）
- `pkg/services/shorturls/` - 短網址：提供 Dashboard 或特定資源的短連結生成與跳轉（如 /d/abc123），在 Grafana 中屬於「非必要服務」。它並不參與任何核心監控、儀表板、查詢、告警等功能的執行流程。


## 遷移至 /internal
- `pkg/services//pkg/services/ngalert` - 新一代警報系統
➡️ Detectviz: internal/ngalert/
- `pkg/services//pkg/services/notifications/` - 通知服務
➡️ Detectviz: internal/notifications/

------------------------------ ❌ 不納入的服務 ------------------------------

### 儀表板專用服務
- `pkg/services/playlist/` - 播放列表服務：允許使用者建立「Dashboard 播放列表」，可自動輪播多個 dashboard（例如在監控大螢幕輪流顯示多個頁面）。支援設定輪播順序、每頁顯示時長等。
- `pkg/services/plugindashboards/` - 插件儀表板：Plugin 提供的範例 dashboard 管理，支援自動註冊、移除、同步等。
- `pkg/services/publicdashboards/` - 公開儀表板：公開 dashboard 分享（不需登入即可瀏覽），可控存取權限、有效期。

## 雲端與整合更新服務
- `pkg/services/cloudmigration/` - 雲端遷移
- `pkg/services/gcom/` - Grafana Cloud 整合
- `pkg/services/updatemanager/` - 更新管理：Grafana 用來檢查與通知插件或系統更新，新版 Grafana 發布檢查。


1. auth-service
- 必要核心模組 [internal/auths]
login, loginattempt, oauthtoken, ssosettings

- 可 plugin 化的模組（介面化後適合擴充、替換）[plugins/auths]
`apikey` 可以支援自訂簽名、儲存後端、權限擴充（如組織隔離）
`anonymous` 可設計為 pluggable 匿名識別 provider（如 cookie / IP 白名單）
`extsvcauth` 本質就是 plugin-like，建議模組化多種 provider
`kmsproviders` 屬 infra layer，應抽象為 KMSProvider interface
`ldap` 可設計為獨立 plugin（如 Keycloak, AzureAD 可各自 plugin 化）
`licensing`  雖非純 auth，但具延伸驗證角色，EE 下常為獨立服務或 plugin 控制功能開關
`secrets` 應 plugin 化支援 memory、DB、Vault、KMS 等後端

2. rbac-services 
- 必要核心模組 [internal/rbac]
accesscontrol, org, team, user

- 可 plugin 化的模組（介面化後適合擴充、替換）[plugins/rbac]
`serviceaccounts` 可做雲原生帳號管理
`temp_user` 可做邀請註冊、匿名用戶管理

3. system-services
- 必要核心模組 [internal/system]
apiserver, contexthandler, frontend, grpcserver, datasourceproxy, live, pluginsintegration, provisioning, query, queryhistory, search, searchusers, searchV2, cleanup, caching

- 可 plugin 化的模組（介面化後適合擴充、替換、動態支援、可開關）[plugins/]
`featuremgmt` 控管 feature flag 開關，可 plugin 化支援動態配置來源（如 config file、DB、remote）
`quota` 控制儲存用量/呼叫次數等資源限制機制，適合 plugin 化支援企業方案或多租戶 quota 策略
`supportbundles` 用於產生問題診斷壓縮包，非平台運作必要，可完全 plugin 化（甚至 EE 才需）
`correlations` 事件關聯分析視覺邏輯，可 plugin 化做為 tracing 分析延伸工具
`hooks` 模組或plugin的事件註冊點，適合設計為 middleware pattern，不建議耦合於單一模組，可 plugin 化注入 hook 處理器
`stats` 可設計為 StatCollector plugin，支援客製化統計來源（如 SQL, Prometheus）

4. dashboard-service
- 必要核心模組  [internal/dashboards]
dashboards, folder, dashboardversion, tag, librarypanels, libraryelements, preference

- 可 plugin 化的模組（介面化後適合擴充、替換、動態支援、可開關）[plugins/]
`annotations` 與資料來源、事件系統整合，可抽象為可插拔事件提供者
`dashboardsnapshots` 不影響主功能，適合做為分享 plugin，甚至拆為 SaaS 延伸服務
`dashboardimport` 可 plugin 化支援多種格式（JSON, ZIP, Git URL 等）
`star` 純個人化功能，與主流程無關，適合延後實作或 plugin 化
`navtree` 可抽象為 UI 插槽（slot）或前端 plugin 接口，非必要後端邏輯


5. plugins-service
`rendering` 是一個 I/O 邏輯清晰、易擴充、支援替代 backend 的引擎
`screenshot` 是 render-coordinator 的角色，也可抽成 plugin 模組調度渲染並擴充功能（如加浮水印）
`shorturls` 把資料儲存/產生邏輯 plugin 化（如提供短碼策略），留下核心 /goto/:uid handler 在主程式中轉發
`signingkeys` 可實作 memory/db/KMS 多種 plugin backend


## keep-services
cloudmigration, gcom, playlist, plugindashboards, publicdashboards, updatemanager