package versioning_test

import (
	"testing"

	"github.com/detectviz/detectviz/pkg/versioning"
)

type fakeVersioned struct {
	Version string
}

func (f fakeVersioned) GetVersion() string {
	return f.Version
}

// TestMemoryStore_RegisterAndGet 測試版本註冊與查詢。
// zh: 驗證 MemoryStore 能正確註冊與取得版本物件。
func TestMemoryStore_RegisterAndGet(t *testing.T) {
	store := versioning.NewMemoryStore()
	v := fakeVersioned{Version: "v1.0.0"}

	if err := store.Register("test", v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := store.Get("test")
	if !ok {
		t.Fatal("expected to find registered version")
	}
	if got.GetVersion() != "v1.0.0" {
		t.Errorf("unexpected version: %s", got.GetVersion())
	}
}

// TestMemoryStore_Compare 測試版本比較邏輯。
// zh: 驗證版本比較是否正確排序。
func TestMemoryStore_Compare(t *testing.T) {
	store := versioning.NewMemoryStore()

	if store.Compare("v1.2.0", "v1.1.9") != 1 {
		t.Error("expected v1.2.0 > v1.1.9")
	}
	if store.Compare("v1.0.0", "v1.0.0") != 0 {
		t.Error("expected equal versions")
	}
	if store.Compare("v0.9.9", "v1.0.0") != -1 {
		t.Error("expected v0.9.9 < v1.0.0")
	}
}
