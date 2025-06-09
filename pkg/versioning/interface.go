package versioning

// Versioned 表示一個具版本資訊的物件。
// zh: 所有可被版本控制的資源應實作此介面。
type Versioned interface {
	// GetVersion 回傳目前版本。
	// zh: 取得當前版本字串。
	GetVersion() string
}

// VersionManager 定義版本管理行為。
// zh: 提供註冊、比較與升級版本的統一介面。
type VersionManager interface {
	// Register 註冊一個版本物件與對應標籤。
	// zh: 記錄某一版本資源的標籤與狀態。
	Register(name string, v Versioned) error

	// Get 取得指定名稱的版本資訊。
	// zh: 查詢指定名稱的資源版本。
	Get(name string) (Versioned, bool)

	// Compare 比較兩個版本字串的高低。
	// zh: 回傳 a 是否高於 b（例如 "v2.0.0" > "v1.9.0"）。
	Compare(a, b string) int
}
