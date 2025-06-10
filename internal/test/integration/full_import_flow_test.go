package integration_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/fakes"
	"github.com/detectviz/detectviz/pkg/libraryelements"
	"github.com/detectviz/detectviz/pkg/registry"
	"github.com/detectviz/detectviz/pkg/versioning"
)

// TestFullImportFlow 模擬 plugin 掛載、importer 註冊、匯入資源並寫入 registry 與版本控制的完整流程。
// zh: 測試從 plugin 匯入器到資源註冊與版本管理的整合流程。
func TestFullImportFlow(t *testing.T) {
	ctx := context.Background()

	// 建立假 Importer Registry 並註冊一個測試用 Importer
	importerReg := fakes.NewFakeImporterRegistry()
	importer := &fakes.FakeImporter{
		ID: "test-importer",
		GVK_: registry.GVK{
			Group:   "core",
			Version: "v1",
			Kind:    "Host",
		},
		Data: []registry.Resource{
			&libraryelements.BaseElement{
				ID_:   "elem1",
				Type_: "host",
				Raw:   []byte(`{"name":"node1"}`),
			},
		},
	}
	if err := importerReg.Register(importer); err != nil {
		t.Fatalf("failed to register importer: %v", err)
	}

	// 建立資源 Registry
	resRegistry := registry.NewRegistry()

	// 建立版本管理記憶體儲存
	versionStore := versioning.NewMemoryStore()

	// 建立匯入服務 (假設有此服務，模擬匯入邏輯)
	importerService := &fakes.FakeImporterService{
		ImporterRegistry: importerReg,
		ResourceRegistry: resRegistry,
		VersionStore:     versionStore,
	}

	// 執行匯入動作
	if err := importerService.Import(ctx, "test-importer", "fakefile.yaml"); err != nil {
		t.Fatalf("import failed: %v", err)
	}

	// 驗證 Registry 中有匯入的資源
	resources, err := resRegistry.List()
	if err != nil {
		t.Fatalf("failed to list resources: %v", err)
	}
	if len(resources) == 0 {
		t.Error("expected resources in registry after import")
	}

	// 驗證版本儲存中有紀錄版本資訊
	v, ok := versionStore.Get("test-importer")
	if !ok || v.GetVersion() == "" {
		t.Error("expected version info registered for importer")
	}
}
