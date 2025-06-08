package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

type ZapLogger struct {
	l *zap.SugaredLogger
}

// NewZapLogger 建立新的 ZapLogger，接收 zap.SugaredLogger 實例。
// zh: 建立以 zap.SugaredLogger 為基礎的日誌介面實作。
func NewZapLogger(base *zap.SugaredLogger) ifacelogger.Logger {
	return &ZapLogger{l: base}
}

// Info logs a message at InfoLevel.
// zh: 以 InfoLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Info(msg string, args ...any) {
	z.l.Infow(msg, args...)
}

// Warn logs a message at WarnLevel.
// zh: 以 WarnLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Warn(msg string, args ...any) {
	z.l.Warnw(msg, args...)
}

// Error logs a message at ErrorLevel.
// zh: 以 ErrorLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Error(msg string, args ...any) {
	z.l.Errorw(msg, args...)
}

// Debug logs a debug-level message with optional formatting.
// zh: 以 DebugLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Debug(msg string, args ...any) {
	z.l.Debugw(msg, args...)
}

func (z *ZapLogger) Sync() error {
	return z.l.Sync()
}

func (z *ZapLogger) WithFields(fields map[string]any) ifacelogger.Logger {
	return &ZapLogger{l: z.l.With(fields)}
}

func (z *ZapLogger) WithContext(ctx context.Context) ifacelogger.Logger {
	// zh: 可選擇從 context 擷取 trace_id/span_id 並加到欄位
	// 這裡暫以原 logger 回傳
	return z
}

func (z *ZapLogger) Named(name string) ifacelogger.Logger {
	return &ZapLogger{l: z.l.Named(name)}
}

// NewDefaultZap 建立預設設定的 zap.Logger 實例。
// zh: 提供簡易用於本地測試與開發的 logger 初始化方法。
func NewDefaultZap() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, _ := cfg.Build()
	return logger.Sugar()
}
