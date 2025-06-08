# grafana/pkg/login/ 概覽
`pkg/login` 目錄在 Grafana 中的主要功能是：**提供登入系統的抽象與實作，特別是整合各類 OAuth 社交登入（social login）與其邏輯處理機制**。

* * *

✅ 功能說明
------

| 子模組 | 說明 |
| --- | --- |
| `social/` | 定義各種 Social 登入 provider（GitHub、GitLab、Google、Okta、AzureAD 等）的共通介面與欄位結構 |
| `social/connectors/` | 每個 provider 對應一個 `*_oauth.go` 檔，包含 OAuth2 流程實作、API 解析、權限比對、帳號同步 |
| `socialimpl/` | 將設定（grafana.ini 或 dynamic settings）轉換為 provider 實例，註冊並提供存取方法 |
| `socialtest/` | 提供 `FakeSocialService`, `MockSocialConnector` 等測試替身與驗證工具 |

* * *

🔌 定義的核心介面
----------

```go
type SocialConnector interface {
  UserInfo(ctx, client, token) (*BasicUserInfo, error)
  IsEmailAllowed(email string) bool
  IsSignupAllowed() bool
  GetOAuthInfo() *OAuthInfo
}
```

```go
type Service interface {
  GetOAuthProviders() map[string]bool
  GetConnector(string) (SocialConnector, error)
}
```

這些介面讓其他模組（如 `authn`, `middleware`, `apiserver`）能夠統一調用各種登入來源，實現身份認證與授權。

* * *

🧩 與其他模組的關聯
-----------

| 模組 | 用途 |
| --- | --- |
| `pkg/services/authn/` | 將 `login.SocialConnector` 注入為 `authn.Client`，支援 token 驗證、自動註冊 |
| `pkg/services/ssosettings` | 允許使用者透過 API 設定 OAuth provider |
| `pkg/middleware/auth.go` | 在 API 呼叫過程中處理未登入者導向社交登入入口grafana-pkg-all-code |

* * *

🧠 設計意圖與可延伸價值
-------------

Grafana 將 `login/social` 拆為獨立模組，而非混入 `auth`、`authn`，目的是：

*   模組解耦（符合 SRP）
    
*   支援多種登入策略並行
    
*   可替換或禁用單一 provider
    
*   適合 SaaS 多租戶或動態設定場景
    

* * *

🚀 對 detectviz 的啟發
------------------

若你將來想支援：

*   OAuth2 API token 交換
    
*   單一入口整合 GitHub / Google 登入
    
*   動態切換登入策略或身份驗證模組
    

你可以參考這個設計：

*   interface 抽象放 `pkg/login`
    
*   provider 建構器放 `pkg/login/connectors`
    
*   統一注入與管理邏輯放 `pkg/login/service.go`
    
*   模擬測試放 `pkg/login/test/`
    

是否需要我幫你定義 detectviz 的 login interface 與 plugin-based 登入機制草稿？

# 遷移到 Detectviz internal/auth/ 後怎麼做

✅ 是的，你可以將 pkg/login/ 遷移至 internal/auth/，而且這麼做在 Detectviz 的分層原則下是合理且建議的。以下是詳細說明：

⸻

✅ 為什麼建議遷移至 internal/auth/

🎯 1. 屬於應用邏輯（Authentication Usecase）
	•	pkg/login/ 實際負責的是登入行為與 provider 驗證邏輯（OAuth 流程、帳號綁定…）
	•	它屬於 Detectviz 內部的身份驗證流程實作，而非獨立共用模組
	•	不太可能在其他專案或模組中重用（不像 redis client、http client 等）

⸻

📦 2. interface 抽象應該再整理進 pkg/iface/auth.go
	•	若你想保留 plugin 化的可能，可將 SocialConnector interface 抽象定義移至 pkg/iface/auth.go
	•	將所有實作與 provider 流程（如 google_oauth.go）搬至 internal/auth/social/

⸻

✅ 遷移建議結構

internal/auth/
├── social/              # provider 實作（google, github, ...）
├── socialimpl/          # provider 建構與註冊
├── socialtest/          # mock 與替身
├── service.go           # 登入流程邏輯
└── handler.go           # 登入 API 或表單頁

pkg/iface/auth.go        # 定義 SocialConnector 等抽象


⸻

🧭 延伸規劃（選擇性）
	•	若將來登入行為也支援 plugin 化，可讓 provider 實作放到 plugins/login-provider/xxx
	•	但在目前階段，「Detectviz 預設支援多種 provider」這件事仍屬於 app 內部邏輯，因此維持在 internal/ 更合理

