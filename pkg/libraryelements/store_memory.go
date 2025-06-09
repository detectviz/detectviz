package libraryelements

import (
	"fmt"
	"sync"
)

// MemoryStore 為記憶體版元素儲存實作。
// zh: 提供開發與測試階段使用的輕量儲存實作。
type MemoryStore struct {
	mu      sync.RWMutex
	entries map[string]Element
}

// NewMemoryStore 建立新的 MemoryStore 實體。
// zh: 初始化內部 map 結構。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		entries: make(map[string]Element),
	}
}

// Save 儲存一個元素。
// zh: 若相同 ID 已存在則覆蓋。
func (m *MemoryStore) Save(e Element) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[e.ID()] = e
	return nil
}

// FindByID 依 ID 取得元素。
// zh: 查詢對應 ID 的元素，若不存在則回傳錯誤。
func (m *MemoryStore) FindByID(id string) (Element, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[id]
	if !ok {
		return nil, fmt.Errorf("element not found: %s", id)
	}
	return e, nil
}

// Delete 移除指定 ID 元素。
// zh: 若元素不存在不報錯。
func (m *MemoryStore) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, id)
	return nil
}

// List 回傳所有儲存的元素。
// zh: 取得目前所有已註冊元素清單。
func (m *MemoryStore) List() ([]Element, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]Element, 0, len(m.entries))
	for _, e := range m.entries {
		result = append(result, e)
	}
	return result, nil
}
