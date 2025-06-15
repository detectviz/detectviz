package fake

import (
	"context"
	"time"
)

// FakeContext is a mock implementation of context.Context for testing.
// zh: FakeContext 是用於測試的 context.Context 模擬實作。
type FakeContext struct {
	parent   context.Context
	values   map[interface{}]interface{}
	deadline time.Time
	done     chan struct{}
	err      error
	canceled bool
}

// NewFakeContext creates a new fake context.
// zh: NewFakeContext 建立新的假上下文。
func NewFakeContext() *FakeContext {
	return &FakeContext{
		parent: context.Background(),
		values: make(map[interface{}]interface{}),
		done:   make(chan struct{}),
	}
}

// NewFakeContextWithParent creates a new fake context with a parent.
// zh: NewFakeContextWithParent 建立帶有父上下文的新假上下文。
func NewFakeContextWithParent(parent context.Context) *FakeContext {
	return &FakeContext{
		parent: parent,
		values: make(map[interface{}]interface{}),
		done:   make(chan struct{}),
	}
}

// NewFakeContextWithTracing creates a new fake context with tracing information.
// zh: NewFakeContextWithTracing 建立帶有追蹤資訊的新假上下文。
func NewFakeContextWithTracing(traceID, spanID string) *FakeContext {
	ctx := NewFakeContext()
	ctx.values["trace_id"] = traceID
	ctx.values["span_id"] = spanID
	return ctx
}

// NewFakeContextWithTimeout creates a new fake context with timeout.
// zh: NewFakeContextWithTimeout 建立帶有超時的新假上下文。
func NewFakeContextWithTimeout(timeout time.Duration) *FakeContext {
	ctx := NewFakeContext()
	ctx.deadline = time.Now().Add(timeout)

	// Start timeout goroutine
	go func() {
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-timer.C:
			ctx.cancel(context.DeadlineExceeded)
		case <-ctx.done:
			// Context was already canceled
		}
	}()

	return ctx
}

// Deadline returns the deadline for the context.
func (fc *FakeContext) Deadline() (deadline time.Time, ok bool) {
	if !fc.deadline.IsZero() {
		return fc.deadline, true
	}
	if fc.parent != nil {
		return fc.parent.Deadline()
	}
	return time.Time{}, false
}

// Done returns a channel that is closed when the context is canceled.
func (fc *FakeContext) Done() <-chan struct{} {
	return fc.done
}

// Err returns the error that caused the context to be canceled.
func (fc *FakeContext) Err() error {
	if fc.canceled {
		return fc.err
	}
	if fc.parent != nil {
		return fc.parent.Err()
	}
	return nil
}

// Value returns the value associated with the given key.
func (fc *FakeContext) Value(key interface{}) interface{} {
	if value, exists := fc.values[key]; exists {
		return value
	}
	if fc.parent != nil {
		return fc.parent.Value(key)
	}
	return nil
}

// SetValue sets a value in the context.
// zh: SetValue 在上下文中設置值。
func (fc *FakeContext) SetValue(key, value interface{}) {
	fc.values[key] = value
}

// SetTracing sets tracing information in the context.
// zh: SetTracing 在上下文中設置追蹤資訊。
func (fc *FakeContext) SetTracing(traceID, spanID string) {
	fc.values["trace_id"] = traceID
	fc.values["span_id"] = spanID
}

// GetTraceID returns the trace ID from the context.
// zh: GetTraceID 從上下文中返回追蹤 ID。
func (fc *FakeContext) GetTraceID() string {
	if traceID := fc.Value("trace_id"); traceID != nil {
		if traceIDStr, ok := traceID.(string); ok {
			return traceIDStr
		}
	}
	return ""
}

// GetSpanID returns the span ID from the context.
// zh: GetSpanID 從上下文中返回跨度 ID。
func (fc *FakeContext) GetSpanID() string {
	if spanID := fc.Value("span_id"); spanID != nil {
		if spanIDStr, ok := spanID.(string); ok {
			return spanIDStr
		}
	}
	return ""
}

// Cancel cancels the context with the given error.
// zh: Cancel 使用給定錯誤取消上下文。
func (fc *FakeContext) Cancel(err error) {
	fc.cancel(err)
}

// cancel is the internal cancel method.
// zh: cancel 是內部取消方法。
func (fc *FakeContext) cancel(err error) {
	if fc.canceled {
		return
	}

	fc.canceled = true
	fc.err = err
	close(fc.done)
}

// IsCanceled returns whether the context has been canceled.
// zh: IsCanceled 返回上下文是否已被取消。
func (fc *FakeContext) IsCanceled() bool {
	return fc.canceled
}

// IsExpired returns whether the context has expired (deadline exceeded).
// zh: IsExpired 返回上下文是否已過期（超過截止時間）。
func (fc *FakeContext) IsExpired() bool {
	if fc.deadline.IsZero() {
		return false
	}
	return time.Now().After(fc.deadline)
}

// GetAllValues returns all values stored in the context.
// zh: GetAllValues 返回上下文中存儲的所有值。
func (fc *FakeContext) GetAllValues() map[interface{}]interface{} {
	values := make(map[interface{}]interface{})
	for k, v := range fc.values {
		values[k] = v
	}
	return values
}

// HasValue checks if a value exists for the given key.
// zh: HasValue 檢查給定鍵是否存在值。
func (fc *FakeContext) HasValue(key interface{}) bool {
	_, exists := fc.values[key]
	if !exists && fc.parent != nil {
		return fc.parent.Value(key) != nil
	}
	return exists
}

// Clone creates a copy of the context.
// zh: Clone 建立上下文的副本。
func (fc *FakeContext) Clone() *FakeContext {
	newCtx := &FakeContext{
		parent:   fc.parent,
		values:   make(map[interface{}]interface{}),
		deadline: fc.deadline,
		done:     make(chan struct{}),
		err:      fc.err,
		canceled: false, // New context starts uncanceled
	}

	// Copy values
	for k, v := range fc.values {
		newCtx.values[k] = v
	}

	return newCtx
}

// WithValue creates a new context with the given key-value pair.
// zh: WithValue 建立帶有給定鍵值對的新上下文。
func (fc *FakeContext) WithValue(key, value interface{}) *FakeContext {
	newCtx := fc.Clone()
	newCtx.values[key] = value
	return newCtx
}

// WithTracing creates a new context with tracing information.
// zh: WithTracing 建立帶有追蹤資訊的新上下文。
func (fc *FakeContext) WithTracing(traceID, spanID string) *FakeContext {
	newCtx := fc.Clone()
	newCtx.values["trace_id"] = traceID
	newCtx.values["span_id"] = spanID
	return newCtx
}

// WithDeadline creates a new context with the given deadline.
// zh: WithDeadline 建立帶有給定截止時間的新上下文。
func (fc *FakeContext) WithDeadline(deadline time.Time) *FakeContext {
	newCtx := fc.Clone()
	newCtx.deadline = deadline

	// Start deadline goroutine
	go func() {
		duration := time.Until(deadline)
		if duration <= 0 {
			newCtx.cancel(context.DeadlineExceeded)
			return
		}

		timer := time.NewTimer(duration)
		defer timer.Stop()

		select {
		case <-timer.C:
			newCtx.cancel(context.DeadlineExceeded)
		case <-newCtx.done:
			// Context was already canceled
		}
	}()

	return newCtx
}

// WithTimeout creates a new context with the given timeout.
// zh: WithTimeout 建立帶有給定超時的新上下文。
func (fc *FakeContext) WithTimeout(timeout time.Duration) *FakeContext {
	return fc.WithDeadline(time.Now().Add(timeout))
}

// Reset resets the context to initial state.
// zh: Reset 重置上下文到初始狀態。
func (fc *FakeContext) Reset() {
	fc.values = make(map[interface{}]interface{})
	fc.deadline = time.Time{}
	fc.done = make(chan struct{})
	fc.err = nil
	fc.canceled = false
}

// FakeContextBuilder helps build fake contexts with various configurations.
// zh: FakeContextBuilder 幫助建立具有各種配置的假上下文。
type FakeContextBuilder struct {
	parent   context.Context
	values   map[interface{}]interface{}
	deadline time.Time
	traceID  string
	spanID   string
}

// NewFakeContextBuilder creates a new context builder.
// zh: NewFakeContextBuilder 建立新的上下文建構器。
func NewFakeContextBuilder() *FakeContextBuilder {
	return &FakeContextBuilder{
		values: make(map[interface{}]interface{}),
	}
}

// WithParent sets the parent context.
// zh: WithParent 設置父上下文。
func (fcb *FakeContextBuilder) WithParent(parent context.Context) *FakeContextBuilder {
	fcb.parent = parent
	return fcb
}

// WithValue adds a key-value pair.
// zh: WithValue 添加鍵值對。
func (fcb *FakeContextBuilder) WithValue(key, value interface{}) *FakeContextBuilder {
	fcb.values[key] = value
	return fcb
}

// WithTracing sets tracing information.
// zh: WithTracing 設置追蹤資訊。
func (fcb *FakeContextBuilder) WithTracing(traceID, spanID string) *FakeContextBuilder {
	fcb.traceID = traceID
	fcb.spanID = spanID
	return fcb
}

// WithDeadline sets the deadline.
// zh: WithDeadline 設置截止時間。
func (fcb *FakeContextBuilder) WithDeadline(deadline time.Time) *FakeContextBuilder {
	fcb.deadline = deadline
	return fcb
}

// WithTimeout sets a timeout.
// zh: WithTimeout 設置超時。
func (fcb *FakeContextBuilder) WithTimeout(timeout time.Duration) *FakeContextBuilder {
	fcb.deadline = time.Now().Add(timeout)
	return fcb
}

// Build creates the fake context.
// zh: Build 建立假上下文。
func (fcb *FakeContextBuilder) Build() *FakeContext {
	var ctx *FakeContext

	if fcb.parent != nil {
		ctx = NewFakeContextWithParent(fcb.parent)
	} else {
		ctx = NewFakeContext()
	}

	// Set values
	for k, v := range fcb.values {
		ctx.values[k] = v
	}

	// Set tracing
	if fcb.traceID != "" {
		ctx.values["trace_id"] = fcb.traceID
	}
	if fcb.spanID != "" {
		ctx.values["span_id"] = fcb.spanID
	}

	// Set deadline
	if !fcb.deadline.IsZero() {
		ctx.deadline = fcb.deadline

		// Start deadline goroutine
		go func() {
			duration := time.Until(fcb.deadline)
			if duration <= 0 {
				ctx.cancel(context.DeadlineExceeded)
				return
			}

			timer := time.NewTimer(duration)
			defer timer.Stop()

			select {
			case <-timer.C:
				ctx.cancel(context.DeadlineExceeded)
			case <-ctx.done:
				// Context was already canceled
			}
		}()
	}

	return ctx
}

// Helper functions for common context patterns

// ContextWithTracing creates a context with tracing information.
// zh: ContextWithTracing 建立帶有追蹤資訊的上下文。
func ContextWithTracing(traceID, spanID string) context.Context {
	return NewFakeContextWithTracing(traceID, spanID)
}

// ContextWithTimeout creates a context with timeout.
// zh: ContextWithTimeout 建立帶有超時的上下文。
func ContextWithTimeout(timeout time.Duration) context.Context {
	return NewFakeContextWithTimeout(timeout)
}

// ContextWithValues creates a context with multiple values.
// zh: ContextWithValues 建立帶有多個值的上下文。
func ContextWithValues(values map[interface{}]interface{}) context.Context {
	ctx := NewFakeContext()
	for k, v := range values {
		ctx.values[k] = v
	}
	return ctx
}
