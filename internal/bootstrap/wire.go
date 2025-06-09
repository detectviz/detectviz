package bootstrap

import (
	modulesadapter "github.com/detectviz/detectviz/internal/adapters/modules"
	serveradapter "github.com/detectviz/detectviz/internal/adapters/server"
	"github.com/detectviz/detectviz/internal/modules"
	"github.com/detectviz/detectviz/internal/server"
	ifaceserver "github.com/detectviz/detectviz/pkg/ifaces/server"
)

// BuildServer assembles all core components and returns a Server interface.
// zh: 組裝所有系統元件並回傳 Server interface 實例。
func BuildServer() ifaceserver.Server {
	cfg := LoadConfig()
	engine := modules.NewEngine()

	engineAdapter := modulesadapter.NewEngineAdapter(engine)
	srv := server.NewServer(cfg, cfg.Logger(), engineAdapter)

	return serveradapter.NewServerAdapter(srv)
}
