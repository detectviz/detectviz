package versioningadapter

import (
	"github.com/detectviz/detectviz/pkg/versioning"
)

// StoreAdapter 包裝 VersionManager 作為統一注入點。
// zh: 將 VersionManager 實作轉換為可共用的 adapter 結構。
type StoreAdapter struct {
	Store versioning.VersionManager
}

// NewStoreAdapter 建立新的 StoreAdapter。
// zh: 接收任一實作 versioning.VersionManager 的實體。
func NewStoreAdapter(vm versioning.VersionManager) *StoreAdapter {
	return &StoreAdapter{Store: vm}
}

// Register 委派給內部 Store 實作。
// zh: 註冊版本物件。
func (a *StoreAdapter) Register(name string, v versioning.Versioned) error {
	return a.Store.Register(name, v)
}

// Get 委派給內部 Store 實作。
// zh: 查詢版本物件。
func (a *StoreAdapter) Get(name string) (versioning.Versioned, bool) {
	return a.Store.Get(name)
}

// Compare 委派給內部 Store 實作。
// zh: 比較兩版本字串。
func (a *StoreAdapter) Compare(aVer, bVer string) int {
	return a.Store.Compare(aVer, bVer)
}
