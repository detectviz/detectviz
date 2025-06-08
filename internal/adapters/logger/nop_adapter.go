package logger

import "github.com/detectviz/detectviz/pkg/ifaces/logger"

// NewNopLogger returns a no-op logger implementation.
// zh: 回傳一個不會輸出任何內容的 logger，常用於測試或作為預設 fallback。
func NewNopLogger() logger.Logger {
	return logger.NopLogger{}
}
