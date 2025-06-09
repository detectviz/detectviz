package importeradapter_test

import (
	"context"
	"testing"

	importeradapter "github.com/detectviz/detectviz/internal/adapters/importer"
	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// fakeImporter 是用來模擬 importer.Importer 的測試用假件。
// zh: 提供測試用 Importer 的假實作。
type fakeImporter struct {
	id   string
	gvk  registry.GVK
	data []registry.Resource
}

// Name 回傳 Importer 名稱。
// zh: 回傳假件的識別名稱。
func (f *fakeImporter) Name() string { return f.id }

// GVK 回傳此 Importer 處理的資源類型。
// zh: 回傳 Importer 處理的 GVK。
func (f *fakeImporter) GVK() registry.GVK { return f.gvk }

// Load 模擬資源匯入流程。
// zh: 回傳預設的測試資源清單與錯誤（固定為 nil）。
func (f *fakeImporter) Load(ctx context.Context, path string) ([]registry.Resource, error) {
	return f.data, nil
}

// TestRegistry_RegisterAndGet 測試 Importer 註冊與查詢是否正確。
// zh: 驗證 Register() 與 Get() 是否能正確儲存與查詢 Importer。
func TestRegistry_RegisterAndGet(t *testing.T) {
	r := importeradapter.NewRegistry()
	fi := &fakeImporter{id: "test", gvk: registry.GVK{Group: "core", Version: "v1", Kind: "Host"}}

	if err := r.Register(fi); err != nil {
		t.Fatalf("unexpected error on register: %v", err)
	}

	got, ok := r.Get("test")
	if !ok {
		t.Fatal("expected to find registered importer")
	}
	if got.Name() != "test" {
		t.Errorf("expected name 'test', got %s", got.Name())
	}
}

// TestRegistry_Register_Duplicate 測試重複註冊 Importer 的錯誤處理。
// zh: 驗證當 Importer 重複註冊時，是否能正確回傳錯誤。
func TestRegistry_Register_Duplicate(t *testing.T) {
	r := importeradapter.NewRegistry()
	fi := &fakeImporter{id: "dup", gvk: registry.GVK{Group: "core", Version: "v1", Kind: "Host"}}

	if err := r.Register(fi); err != nil {
		t.Fatalf("unexpected error on first register: %v", err)
	}
	if err := r.Register(fi); err == nil {
		t.Fatal("expected error on duplicate register, got nil")
	}
}

// TestRegistry_List 測試所有 Importer 是否能正確列出。
// zh: 驗證 Registry.List() 是否能回傳所有已註冊 Importer。
func TestRegistry_List(t *testing.T) {
	r := importeradapter.NewRegistry()
	r.Register(&fakeImporter{id: "a"})
	r.Register(&fakeImporter{id: "b"})

	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 importers, got %d", len(list))
	}
}
