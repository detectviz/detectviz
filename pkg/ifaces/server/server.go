package server

import "context"

// Server defines the lifecycle of the main application server.
// zh: 定義伺服器核心生命週期的介面。
type Server interface {
	// Run starts the server and blocks until context is cancelled or an error occurs.
	// zh: 啟動伺服器，阻塞直到收到關閉訊號或發生錯誤。
	Run(ctx context.Context) error

	// Shutdown gracefully shuts down the server and its components.
	// zh: 優雅地關閉伺服器與所有相依元件。
	Shutdown(ctx context.Context) error
}
