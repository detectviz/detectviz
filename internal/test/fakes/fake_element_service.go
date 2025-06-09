package fakes

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/libraryelements"
)

// FakeElementStore 為測試用 ElementStore 假件。
// zh: 提供記憶體內部儲存元件的假實作，用於元件 CRUD 測試。
type FakeElementStore struct {
	mu      sync.RWMutex
	entries map[string]libraryelements.Element
}

// NewFakeElementStore 建立 FakeElementStore。
// zh: 初始化記憶體儲存空間。
func NewFakeElementStore() *FakeElementStore {
	return &FakeElementStore{
		entries: make(map[string]libraryelements.Element),
	}
}

// Save 儲存一個元件。
func (f *FakeElementStore) Save(e libraryelements.Element) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.entries[e.ID()] = e
	return nil
}

// FindByID 查詢元件。
func (f *FakeElementStore) FindByID(id string) (libraryelements.Element, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	e, ok := f.entries[id]
	if !ok {
		return nil, fmt.Errorf("element not found: %s", id)
	}
	return e, nil
}

// Delete 移除指定元件。
func (f *FakeElementStore) Delete(id string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.entries, id)
	return nil
}

// List 回傳所有元件。
func (f *FakeElementStore) List() ([]libraryelements.Element, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	result := make([]libraryelements.Element, 0, len(f.entries))
	for _, e := range f.entries {
		result = append(result, e)
	}
	return result, nil
}
