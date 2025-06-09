package importer

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// Importer 定義資源匯入器介面。
// zh: 用於動態載入資源檔案並回傳符合 Registry 格式的結構。
type Importer interface {
	// Name 回傳此 Importer 的唯一名稱。
	Name() string

	// GVK 回傳此 Importer 處理的資源類型。
	GVK() registry.GVK

	// Load 執行實際匯入邏輯，回傳 Resource 清單與錯誤。
	Load(ctx context.Context, filePath string) ([]registry.Resource, error)
}
