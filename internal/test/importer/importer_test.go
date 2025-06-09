package importer_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// TestFakeImporter_Load 測試 FakeImporter 的 Load 函式。
// zh: 驗證 FakeImporter 回傳的 GVK、名稱、與資源清單數量正確。
func TestFakeImporter_Load(t *testing.T) {
	fake := &fakes.FakeImporter{
		ID: "fake-host-importer",
		GVK_: registry.GVK{
			Group:   "core",
			Version: "v1",
			Kind:    "Host",
		},
		Data: []registry.Resource{
			&fakeResource{name: "node01"},
			&fakeResource{name: "node02"},
		},
		Err: nil,
	}

	if fake.Name() != "fake-host-importer" {
		t.Errorf("unexpected importer name: %s", fake.Name())
	}

	if fake.GVK().Kind != "Host" {
		t.Errorf("unexpected GVK kind: %s", fake.GVK().Kind)
	}

	res, err := fake.Load(context.Background(), "fake.yaml")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(res) != 2 {
		t.Errorf("expected 2 resources, got: %d", len(res))
	}
}

// fakeResource 為模擬用資源實作。
// zh: 假資源物件，用來實作 registry.Resource 介面。
type fakeResource struct {
	name string
}

// GetName 回傳 fakeResource 的名稱。
// zh: 實作 registry.Resource 的 GetName() 方法。
func (r *fakeResource) GetName() string {
	return r.name
}
