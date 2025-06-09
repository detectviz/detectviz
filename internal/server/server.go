package server

import (
	"net/http"

	"github.com/detectviz/detectviz/pkg/ifaces/config"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// Server represents the core application server.
// zh: Server 為核心伺服器，整合設定、日誌與模組控制。
type Server struct {
	Config       config.ConfigProvider
	Logger       logger.Logger
	ModuleEngine modules.ModuleEngine
	HTTPServer   *http.Server
}

// NewServer creates a new instance of Server.
// zh: 建立新的伺服器實例，注入設定與模組控制元件。
func NewServer(
	cfg config.ConfigProvider,
	log logger.Logger,
	engine modules.ModuleEngine,
) *Server {
	return &Server{
		Config:       cfg,
		Logger:       log,
		ModuleEngine: engine,
		HTTPServer:   nil, // 後續由 runner/init 時注入
	}
}
