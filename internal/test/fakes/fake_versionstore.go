package fakes

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/versioning"
)

// FakeVersionStore 為 VersionManager 的測試假件。
// zh: 用於測試流程中的記憶體版本管理器。
type FakeVersionStore struct {
	mu        sync.RWMutex
	Items     map[string]versioning.Versioned
	CompareFn func(a, b string) int // 可選：自訂比較邏輯（可為 nil）
}

// NewFakeVersionStore 建立假版本儲存實體。
// zh: 預設使用 map 儲存版本資料。
func NewFakeVersionStore() *FakeVersionStore {
	return &FakeVersionStore{
		Items: make(map[string]versioning.Versioned),
	}
}

// Register 儲存版本物件。
// zh: 模擬版本註冊行為。
func (f *FakeVersionStore) Register(name string, v versioning.Versioned) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if _, exists := f.Items[name]; exists {
		return fmt.Errorf("version already registered: %s", name)
	}
	f.Items[name] = v
	return nil
}

// Get 查詢指定名稱的版本。
// zh: 模擬查詢版本資源。
func (f *FakeVersionStore) Get(name string) (versioning.Versioned, bool) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	v, ok := f.Items[name]
	return v, ok
}

// Compare 比較兩個版本字串。
// zh: 模擬版本比較邏輯，若未提供 CompareFn 則回傳 0。
func (f *FakeVersionStore) Compare(a, b string) int {
	if f.CompareFn != nil {
		return f.CompareFn(a, b)
	}
	return 0
}
