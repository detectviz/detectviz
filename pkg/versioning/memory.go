package versioning

import (
	"fmt"
	"sync"
)

// MemoryStore 提供記憶體內的版本管理。
// zh: 測試與開發用的版本儲存實作。
type MemoryStore struct {
	mu       sync.RWMutex
	registry map[string]Versioned
}

// NewMemoryStore 建立 MemoryStore 實例。
// zh: 建立新的記憶體版本管理器。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		registry: make(map[string]Versioned),
	}
}

// Register 註冊版本資源。
// zh: 記錄名稱對應的版本資訊。
func (m *MemoryStore) Register(name string, v Versioned) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.registry[name]; exists {
		return fmt.Errorf("version already registered: %s", name)
	}
	m.registry[name] = v
	return nil
}

// Get 查詢名稱對應的版本資源。
// zh: 根據名稱取得已註冊的版本物件。
func (m *MemoryStore) Get(name string) (Versioned, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.registry[name]
	return v, ok
}

// Compare 比較兩個版本字串。
// zh: 回傳 a 是否大於 b：1 表 a > b，-1 表 a < b，0 表相等。
func (m *MemoryStore) Compare(a, b string) int {
	as := normalizeVersion(a)
	bs := normalizeVersion(b)

	for i := 0; i < len(as) && i < len(bs); i++ {
		if as[i] > bs[i] {
			return 1
		} else if as[i] < bs[i] {
			return -1
		}
	}

	if len(as) > len(bs) {
		return 1
	} else if len(as) < len(bs) {
		return -1
	}
	return 0
}
