package libraryelements_test

import (
	"testing"

	"github.com/detectviz/detectviz/pkg/libraryelements"
)

// TestMemoryStore_CRUD 測試 MemoryStore 的儲存、查詢、刪除與列出功能。
// zh: 確保 MemoryStore 的基本 CRUD 行為符合預期。
func TestMemoryStore_CRUD(t *testing.T) {
	store := libraryelements.NewMemoryStore()

	elem := libraryelements.BaseElement{
		ID_:   "e1",
		Type_: "chart",
		Raw:   []byte(`{"title":"CPU Usage"}`),
	}

	// 測試儲存
	if err := store.Save(elem); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// 測試查詢
	got, err := store.FindByID("e1")
	if err != nil {
		t.Fatalf("FindByID() error: %v", err)
	}
	if got.ID() != "e1" {
		t.Errorf("expected ID 'e1', got %s", got.ID())
	}

	// 測試列出
	list, err := store.List()
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 element, got %d", len(list))
	}

	// 測試刪除
	if err := store.Delete("e1"); err != nil {
		t.Fatalf("Delete() error: %v", err)
	}
	_, err = store.FindByID("e1")
	if err == nil {
		t.Error("expected error for deleted element, got nil")
	}
}
