package fakes

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// FakeImporter 為測試用匯入器，符合 importer.Importer 介面。
// zh: 提供固定 GVK 與結果的假匯入器，用於單元測試。
type FakeImporter struct {
	ID   string              // 匯入器名稱 / importer ID
	GVK_ registry.GVK        // 處理的資源類型 / handled GVK
	Data []registry.Resource // 匯入回傳的資料 / resources to return
	Err  error               // 預設錯誤回應 / default error to return
}

// Name 回傳匯入器名稱。
// zh: 實作 importer.Importer 的 Name() 方法。
func (f *FakeImporter) Name() string {
	return f.ID
}

// GVK 回傳匯入器支援的資源類型。
// zh: 實作 importer.Importer 的 GVK() 方法。
func (f *FakeImporter) GVK() registry.GVK {
	return f.GVK_
}

// Load 執行匯入邏輯，回傳測試資料與錯誤。
// zh: 實作 importer.Importer 的 Load() 方法，固定回傳預設資料與錯誤。
func (f *FakeImporter) Load(ctx context.Context, filePath string) ([]registry.Resource, error) {
	return f.Data, f.Err
}
