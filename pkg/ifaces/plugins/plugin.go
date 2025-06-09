package plugins

// Plugin 定義所有 plugin 應實作的基本介面。
// zh: 提供插件的初始化、終止與版本查詢能力。
type Plugin interface {
	// Name 回傳插件名稱。
	// zh: 插件的唯一識別名稱。
	Name() string

	// Version 回傳插件版本。
	// zh: 插件版本號，可與 versioning 模組比對。
	Version() string

	// Init 為插件初始化邏輯。
	// zh: 啟用插件時的初始化動作。
	Init() error

	// Close 為插件關閉邏輯。
	// zh: 終止插件時的釋放與關閉操作。
	Close() error
}
