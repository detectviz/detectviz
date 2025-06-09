package libraryelementsadapter_test

import (
	"testing"

	adapter "github.com/detectviz/detectviz/internal/adapters/libraryelements"
	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/libraryelements"
)

// TestServiceAdapter_CRUD tests basic CRUD behavior of ServiceAdapter using a fake store.
// zh: 使用假件測試 ServiceAdapter 的基本儲存、查詢、刪除邏輯是否正確。
func TestServiceAdapter_CRUD(t *testing.T) {
	store := fakes.NewFakeElementStore()
	svc := adapter.NewServiceAdapter(store)

	elem := libraryelements.BaseElement{
		ID_:   "btn1",
		Type_: "input",
		Raw:   []byte(`{"label":"Submit"}`),
	}

	// Save
	if err := svc.Save(elem); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	// Find
	got, err := svc.FindByID("btn1")
	if err != nil {
		t.Fatalf("FindByID error: %v", err)
	}
	if got.ID() != "btn1" {
		t.Errorf("expected ID 'btn1', got %s", got.ID())
	}

	// List
	list, err := svc.List()
	if err != nil {
		t.Fatalf("List error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 element, got %d", len(list))
	}

	// Delete
	if err := svc.Delete("btn1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}
	if _, err := svc.FindByID("btn1"); err == nil {
		t.Error("expected error for deleted element, got nil")
	}
}
