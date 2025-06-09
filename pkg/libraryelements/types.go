package libraryelements

// BaseElement 為預設元件實作。
// zh: Element 的通用結構，可用於儲存與轉譯。
type BaseElement struct {
	ID_   string // 元件 ID
	Type_ string // 元件類型（如 chart, input）
	Raw   []byte // 元件資料內容
}

// ID 回傳元件 ID。
// zh: 實作 Element.ID()
func (e BaseElement) ID() string {
	return e.ID_
}

// Type 回傳元件類型。
// zh: 實作 Element.Type()
func (e BaseElement) Type() string {
	return e.Type_
}

// Data 回傳元件原始內容。
// zh: 實作 Element.Data()
func (e BaseElement) Data() []byte {
	return e.Raw
}
