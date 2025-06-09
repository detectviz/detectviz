package bootstrap

import (
	elementadapter "github.com/detectviz/detectviz/internal/adapters/libraryelements"
	"github.com/detectviz/detectviz/pkg/libraryelements"
)

// InitElementService 初始化元件服務實體。
// zh: 建立記憶體儲存版本 ElementService，並封裝為 ServiceAdapter。
func InitElementService() *elementadapter.ServiceAdapter {
	store := libraryelements.NewMemoryStore()
	return elementadapter.NewServiceAdapter(store)
}
