package libraryelementsadapter

import (
	"github.com/detectviz/detectviz/pkg/libraryelements"
)

// ServiceAdapter 將 ElementStore 封裝為 adapter，可用於模組注入。
// zh: 提供元件 CRUD 功能的統一注入點。
type ServiceAdapter struct {
	Store libraryelements.ElementStore
}

// NewServiceAdapter 建立新的 ServiceAdapter。
// zh: 接收任一實作 ElementStore 的實體。
func NewServiceAdapter(store libraryelements.ElementStore) *ServiceAdapter {
	return &ServiceAdapter{Store: store}
}

// Save 儲存元素。
// zh: 將元素儲存至底層儲存庫。
func (a *ServiceAdapter) Save(e libraryelements.Element) error {
	return a.Store.Save(e)
}

// FindByID 查詢元素。
// zh: 根據 ID 查詢元素。
func (a *ServiceAdapter) FindByID(id string) (libraryelements.Element, error) {
	return a.Store.FindByID(id)
}

// Delete 移除元素。
// zh: 移除指定 ID 的元素。
func (a *ServiceAdapter) Delete(id string) error {
	return a.Store.Delete(id)
}

// List 取得所有元素。
// zh: 回傳目前儲存的所有元件。
func (a *ServiceAdapter) List() ([]libraryelements.Element, error) {
	return a.Store.List()
}
