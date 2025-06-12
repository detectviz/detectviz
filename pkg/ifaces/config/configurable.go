package config

// Configurable 可供模組標準實作設定結構的介面。
// zh: 模組設定結構實作此介面後可統一被框架載入與驗證。
type Configurable interface {
	Validate() error
}
