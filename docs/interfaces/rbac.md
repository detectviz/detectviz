# RBAC 模組介面定義（Role-Based Access Control Interfaces）

本文件說明 `internal/rbac/` 模組中各子元件對外暴露的介面，包含權限驗證、組織管理、角色範疇與 middleware 掛載等用途，供 service 層與 middleware 注入使用。

---

## Authorizer

RBAC 核心授權介面，提供使用者對資源執行動作的權限驗證方法。

```go
type Authorizer interface {
    Check(user *User, action string, resource string, scope Scope) (bool, error)
}
```

---

## Scope

代表權限的作用範圍，例如：
- `GlobalScope`
- `OrgScope{OrgID}`
- `ResourceScope{ResourceType, ResourceID}`

```go
type Scope interface {
    Matches(resourceScope Scope) bool
    String() string
}
```

---

## Middleware 授權中介層

可用於 HTTP Router 中，提供細粒度授權控制與注入。

```go
type Middleware interface {
    RequireRole(role string) echo.MiddlewareFunc
    RequireScope(scope Scope) echo.MiddlewareFunc
}
```

---

## OrgManager

組織/租戶管理邏輯，包括 CRUD 與查詢。

```go
type OrgManager interface {
    GetOrgByID(id int64) (*Org, error)
    ListUserOrgs(userID int64) ([]*Org, error)
    DeleteOrg(id int64) error
}
```

---

## TeamManager（可選）

支援使用者群組功能，例如小組分權管理。

```go
type TeamManager interface {
    AddUserToTeam(userID int64, teamID int64) error
    ListTeamsByUser(userID int64) ([]*Team, error)
}
```

---

## ServiceAccountProvider（plugin）

支援透過 token 驗證服務帳號身份。

```go
type ServiceAccountProvider interface {
    Authenticate(token string) (*ServiceAccount, error)
}
```

---

## TempUserProvider（plugin）

支援一次性臨時帳號或邀請流程。

```go
type TempUserProvider interface {
    CreateTemporaryUser(inviteToken string) (*User, error)
    CleanupExpiredUsers() error
}
```

---

## 測試用介面（internal/rbac/mock/）

建議每個 interface 提供對應 `mock` 與 `fake` 實作，供整合測試與單元模擬使用。

---