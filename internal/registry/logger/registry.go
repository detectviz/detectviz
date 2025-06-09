package logger

import (
	loggeradapter "github.com/detectviz/detectviz/internal/adapters/logger"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
)

// ProvideLogger creates a logger instance using the Zap adapter implementation.
// zh: 建立 logger 實例，採用 internal/adapters/logger/zap_adapter.go 作為預設實作。
func ProvideLogger() ifacelogger.Logger {
	// Use Zap logger by default.
	// zh: 預設使用 Zap logger，實作於 internal/adapters/logger/zap_adapter.go。
	zapLogger, _ := zap.NewDevelopment()
	return loggeradapter.NewZapLogger(zapLogger.Sugar())
}
