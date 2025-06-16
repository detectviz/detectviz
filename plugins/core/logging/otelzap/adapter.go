package otelzap

import (
	"context"

	"detectviz/pkg/shared/log"
)

// LoggerAdapter adapts OtelZapLogger to work with pkg/shared/log.LoggerInterface
// zh: LoggerAdapter 將 OtelZapLogger 適配到 pkg/shared/log.LoggerInterface
type LoggerAdapter struct {
	otelZapLogger *OtelZapLogger
	ctx           context.Context
}

// NewLoggerAdapter creates a new logger adapter
// zh: NewLoggerAdapter 建立新的日誌記錄器適配器
func NewLoggerAdapter(otelZapLogger *OtelZapLogger, ctx context.Context) *LoggerAdapter {
	return &LoggerAdapter{
		otelZapLogger: otelZapLogger,
		ctx:           ctx,
	}
}

// Debug implements log.LoggerInterface.Debug
// zh: Debug 實作 log.LoggerInterface.Debug
func (a *LoggerAdapter) Debug(msg string, fields ...interface{}) {
	if a.otelZapLogger != nil {
		a.otelZapLogger.Debug(msg, fields...)
	}
}

// Info implements log.LoggerInterface.Info
// zh: Info 實作 log.LoggerInterface.Info
func (a *LoggerAdapter) Info(msg string, fields ...interface{}) {
	if a.otelZapLogger != nil {
		a.otelZapLogger.Info(msg, fields...)
	}
}

// Warn implements log.LoggerInterface.Warn
// zh: Warn 實作 log.LoggerInterface.Warn
func (a *LoggerAdapter) Warn(msg string, fields ...interface{}) {
	if a.otelZapLogger != nil {
		a.otelZapLogger.Warn(msg, fields...)
	}
}

// Error implements log.LoggerInterface.Error
// zh: Error 實作 log.LoggerInterface.Error
func (a *LoggerAdapter) Error(msg string, fields ...interface{}) {
	if a.otelZapLogger != nil {
		a.otelZapLogger.Error(msg, fields...)
	}
}

// Fatal implements log.LoggerInterface.Fatal
// zh: Fatal 實作 log.LoggerInterface.Fatal
func (a *LoggerAdapter) Fatal(msg string, fields ...interface{}) {
	if a.otelZapLogger != nil {
		a.otelZapLogger.Fatal(msg, fields...)
	}
}

// Ensure LoggerAdapter implements log.LoggerInterface
var _ log.LoggerInterface = (*LoggerAdapter)(nil)

// SetAsGlobalLogger sets the OtelZap logger as the global logger
// zh: SetAsGlobalLogger 將 OtelZap 日誌記錄器設為全域日誌記錄器
func (p *OtelZapPlugin) SetAsGlobalLogger() {
	if p.logger != nil {
		adapter := NewLoggerAdapter(p.logger, context.Background())
		log.SetGlobalLogger(adapter)
	}
}
