package bootstrap

import (
	versioningadapter "github.com/detectviz/detectviz/internal/adapters/versioning"
	"github.com/detectviz/detectviz/pkg/versioning"
)

// InitVersioning 初始化版本控制元件。
// zh: 建立記憶體儲存版本控制實體，並封裝為 adapter 回傳。
func InitVersioning() *versioningadapter.StoreAdapter {
	store := versioning.NewMemoryStore()
	return versioningadapter.NewStoreAdapter(store)
}
