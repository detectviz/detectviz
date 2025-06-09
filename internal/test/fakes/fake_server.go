package fakes

import (
	"context"
)

// FakeServer implements the Server interface for testing purposes.
// zh: 用於測試場景的 Server interface 假實作。
type FakeServer struct {
	RunCalled      bool
	ShutdownCalled bool
	RunError       error
	ShutdownError  error
}

// Run records that it was called and returns the configured error.
// zh: 模擬 Run 行為，並標記已呼叫狀態。
func (f *FakeServer) Run(ctx context.Context) error {
	f.RunCalled = true
	return f.RunError
}

// Shutdown records that it was called and returns the configured error.
// zh: 模擬 Shutdown 行為，並標記已呼叫狀態。
func (f *FakeServer) Shutdown(ctx context.Context) error {
	f.ShutdownCalled = true
	return f.ShutdownError
}
