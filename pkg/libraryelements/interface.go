package libraryelements

// Element 表示可複用的組件元素。
// zh: 每個可插拔的視覺化或邏輯單元元件。
type Element interface {
	// ID 回傳唯一標識。
	// zh: 每個元件應有唯一 ID。
	ID() string

	// Type 類型分類（如 chart, input, filter）。
	// zh: 用於組裝時分類元件型態。
	Type() string

	// Data 傳回原始定義資料。
	// zh: 回傳用於儲存或轉譯的元件內容資料。
	Data() []byte
}

// ElementStore 定義元素儲存行為。
// zh: 負責管理元素 CRUD 的儲存模組介面。
type ElementStore interface {
	Save(e Element) error
	FindByID(id string) (Element, error)
	Delete(id string) error
	List() ([]Element, error)
}

// ElementRenderer 定義元素轉譯為實體輸出的行為。
// zh: 可將元件轉譯為 HTML、JSON、Grafana JSON 等形式。
type ElementRenderer interface {
	Render(e Element) ([]byte, error)
}
