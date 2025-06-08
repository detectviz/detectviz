package memory

import (
	"strings"
	"sync"
)

// MemoryCacheStore 是基於 sync.Map 的快取實作。
// zh: 適用於單元測試與本地環境的記憶體快取。
type MemoryCacheStore struct {
	store sync.Map
}

// NewMemoryCacheStore 回傳一個新的記憶體快取實例。
func NewMemoryCacheStore() *MemoryCacheStore {
	return &MemoryCacheStore{}
}

// Get 取得指定 key 的快取值。
// zh: 若 key 不存在則回傳空字串與 nil 錯誤。
func (m *MemoryCacheStore) Get(key string) (string, error) {
	val, ok := m.store.Load(key)
	if !ok {
		return "", nil
	}
	return val.(string), nil
}

// Set 寫入 key 對應的快取值，忽略 TTL。
// zh: TTL 僅為佔位參數，目前不實作過期。
func (m *MemoryCacheStore) Set(key, value string, ttl int) error {
	m.store.Store(key, value)
	return nil
}

// Has 檢查指定 key 是否存在於快取中。
// zh: 若 key 存在則回傳 true，否則 false。此實作永不錯誤，僅為符合介面定義。
func (m *MemoryCacheStore) Has(key string) (bool, error) {
	_, ok := m.store.Load(key)
	return ok, nil
}

// Delete 移除指定 key。
// zh: 刪除快取資料，成功不回傳錯誤。
func (m *MemoryCacheStore) Delete(key string) error {
	m.store.Delete(key)
	return nil
}

// Keys 回傳符合指定 prefix 的所有 key。
// zh: 支援 key 模糊查詢，用於群組操作。
func (m *MemoryCacheStore) Keys(prefix string) ([]string, error) {
	var keys []string
	m.store.Range(func(k, _ any) bool {
		s := k.(string)
		if strings.HasPrefix(s, prefix) {
			keys = append(keys, s)
		}
		return true
	})
	return keys, nil
}
