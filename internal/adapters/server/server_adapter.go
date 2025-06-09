package serveradapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/server"
)

// ServerAdapter wraps core.Server to implement iface.Server.
// zh: ServerAdapter 包裝 core.Server，使其實作 Server interface。
type ServerAdapter struct {
	srv *core.Server
}

// NewServerAdapter creates a new adapter instance.
// zh: 建立新的 ServerAdapter 實例。
func NewServerAdapter(s *core.Server) *ServerAdapter {
	return &ServerAdapter{
		srv: s,
	}
}

// Run starts the server.
// zh: 啟動伺服器。
func (a *ServerAdapter) Run(ctx context.Context) error {
	return a.srv.Run(ctx)
}

// Shutdown gracefully shuts down the server.
// zh: 優雅地關閉伺服器。
func (a *ServerAdapter) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}
