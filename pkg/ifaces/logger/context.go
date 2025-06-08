package logger

import "context"

// ctxKey 是用於 context 注入與擷取 logger 的專用 key。
// zh: ctxKey 用於避免與其他 context 欄位衝突。
type ctxKey struct{}

// WithContext 將 logger 實例注入至 context 中，供下游模組擷取使用。
// zh: 建議於 middleware 或 handler 中呼叫，讓後續模組皆可透過 FromContext 擷取同一 logger。
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext 從 context 中擷取 logger 實例，若未注入則回傳 NopLogger。
// zh: 此函式可保證不會回傳 nil，避免程式 panic。
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return NopLogger{} // fallback 預設 logger，靜默處理所有輸出
}
