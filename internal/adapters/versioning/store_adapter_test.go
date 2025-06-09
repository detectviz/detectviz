package versioningadapter_test

import (
	"testing"

	versioningadapter "github.com/detectviz/detectviz/internal/adapters/versioning"
	"github.com/detectviz/detectviz/pkg/versioning"
)

type fakeVersioned struct {
	Version string
}

func (f fakeVersioned) GetVersion() string {
	return f.Version
}

// TestStoreAdapter_RegisterAndGet 測試版本註冊與查詢。
// zh: 驗證 StoreAdapter 是否正確包裝 Register 與 Get。
func TestStoreAdapter_RegisterAndGet(t *testing.T) {
	store := versioning.NewMemoryStore()
	adapter := versioningadapter.NewStoreAdapter(store)

	v := fakeVersioned{Version: "v2.1.0"}
	if err := adapter.Register("modA", v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := adapter.Get("modA")
	if !ok {
		t.Fatal("expected registered version")
	}
	if got.GetVersion() != "v2.1.0" {
		t.Errorf("unexpected version: %s", got.GetVersion())
	}
}

// TestStoreAdapter_Compare 測試版本比較功能。
// zh: 驗證 StoreAdapter.Compare 是否正確委派比較邏輯。
func TestStoreAdapter_Compare(t *testing.T) {
	store := versioning.NewMemoryStore()
	adapter := versioningadapter.NewStoreAdapter(store)

	if adapter.Compare("v2.0.0", "v1.9.9") != 1 {
		t.Error("expected v2.0.0 > v1.9.9")
	}
	if adapter.Compare("v1.0.0", "v1.0.0") != 0 {
		t.Error("expected v1.0.0 == v1.0.0")
	}
	if adapter.Compare("v1.0.0", "v1.2.0") != -1 {
		t.Error("expected v1.0.0 < v1.2.0")
	}
}
