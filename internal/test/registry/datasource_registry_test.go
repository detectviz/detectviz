package registry_test

import (
	"context"
	"testing"

	ifacereg "github.com/detectviz/detectviz/pkg/ifaces/registry"
	memreg "github.com/detectviz/detectviz/pkg/registry"
	"github.com/detectviz/detectviz/pkg/registry/apis/datasource"
)

func TestRegisterDatasourceAndCRUD(t *testing.T) {
	ctx := context.Background()
	r := memreg.NewMemoryRegistry()

	gvk := ifacereg.GVK{
		Group:   "core",
		Version: "v1",
		Kind:    "Datasource",
	}

	// Register Datasource resource handler.
	// zh: 註冊 Datasource 資源處理器。
	if err := datasource.RegisterDatasource(r); err != nil {
		t.Fatalf("failed to register datasource: %v", err)
	}

	// Create a datasource resource.
	// zh: 建立一個名為 influx 的 Datasource 資源。
	d := &struct{ ifacereg.Resource }{&dsResource{Name: "influx"}}
	err := r.Create(ctx, gvk, d.Resource)
	if err != nil {
		t.Fatalf("failed to create datasource: %v", err)
	}

	// Retrieve the created datasource.
	// zh: 查詢剛剛建立的 Datasource。
	got, err := r.Get(ctx, gvk, "influx")
	if err != nil {
		t.Fatalf("failed to get datasource: %v", err)
	}
	if got.GetName() != "influx" {
		t.Errorf("expected name 'influx', got %q", got.GetName())
	}

	// List all datasources.
	// zh: 列出所有 Datasource 資源。
	list, err := r.List(ctx, gvk)
	if err != nil {
		t.Fatalf("failed to list datasources: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("expected 1 datasource, got %d", len(list))
	}
}

// dsResource is a simplified datasource used for testing.
// zh: 測試用簡化 Datasource 實作。
type dsResource struct {
	Name string
}

func (d *dsResource) GetName() string {
	return d.Name
}
