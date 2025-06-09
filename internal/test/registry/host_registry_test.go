package registry_test

import (
	"context"
	"testing"

	ifacereg "github.com/detectviz/detectviz/pkg/ifaces/registry"
	memreg "github.com/detectviz/detectviz/pkg/registry"
	"github.com/detectviz/detectviz/pkg/registry/apis/host"
)

func TestRegisterHostAndCRUD(t *testing.T) {
	ctx := context.Background()
	r := memreg.NewMemoryRegistry()

	gvk := ifacereg.GVK{
		Group:   "core",
		Version: "v1",
		Kind:    "Host",
	}

	// Register Host resource handler.
	// zh: 註冊 Host 資源處理器。
	if err := host.RegisterHost(r); err != nil {
		t.Fatalf("failed to register host: %v", err)
	}

	// Create a host resource.
	// zh: 建立一個名為 node01 的 Host 資源。
	h := &struct{ ifacereg.Resource }{&hostResource{Name: "node01"}}
	err := r.Create(ctx, gvk, h.Resource)
	if err != nil {
		t.Fatalf("failed to create host: %v", err)
	}

	// Retrieve the created host.
	// zh: 查詢剛剛建立的 Host。
	got, err := r.Get(ctx, gvk, "node01")
	if err != nil {
		t.Fatalf("failed to get host: %v", err)
	}
	if got.GetName() != "node01" {
		t.Errorf("expected name 'node01', got %q", got.GetName())
	}

	// List all hosts.
	// zh: 列出所有 Host 資源。
	list, err := r.List(ctx, gvk)
	if err != nil {
		t.Fatalf("failed to list hosts: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 host, got %d", len(list))
	}
}

// hostResource is a minimal implementation of Resource used for testing.
// zh: 測試用的簡化 Host 資源實作，僅提供名稱。
type hostResource struct {
	Name string
}

func (h *hostResource) GetName() string {
	return h.Name
}
